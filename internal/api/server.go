// Package api provides the HTTP API server for MediSync.
//
// This package implements the API gateway layer using go-chi/chi router.
// It handles all HTTP routing, middleware chaining, and server lifecycle.
//
// The server implements the middleware chain as specified in the architecture:
// RequestID -> RealIP -> Logger -> Recoverer -> Timeout -> Auth -> Locale -> RateLimit
//
// Usage:
//
//	cfg := config.MustLoad()
//	server := api.NewServer(cfg, deps)
//	if err := server.Start(ctx); err != nil {
//	    log.Fatal("Server failed:", err)
//	}
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/medisync/medisync/internal/api/handlers"
	"github.com/medisync/medisync/internal/api/middleware"
	"github.com/medisync/medisync/internal/cache"
	"github.com/medisync/medisync/internal/config"
	"github.com/medisync/medisync/internal/warehouse"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP API server.
type Server struct {
	config     *config.Config
	logger     *slog.Logger
	router     *chi.Mux
	httpServer *http.Server

	// Dependencies
	db       *warehouse.Repo
	cache    *cache.Client
	keycloak *KeycloakClient
	opa      *OPAClient
	readonly *warehouse.ReadOnlyClient

	// Handlers
	chatHandler *handlers.ChatHandler
	wsHandler   *handlers.WSHandler
}

// Dependencies holds the required dependencies for the API server.
type Dependencies struct {
	DB       *warehouse.Repo
	Cache    *cache.Client
	Keycloak *KeycloakClient
	OPA      *OPAClient
}

// NewServer creates a new API server instance.
func NewServer(cfg *config.Config, deps *Dependencies) *Server {
	if deps == nil {
		deps = &Dependencies{}
	}

	logger := slog.Default()

	s := &Server{
		config:   cfg,
		logger:   logger,
		router:   chi.NewRouter(),
		db:       deps.DB,
		cache:    deps.Cache,
		keycloak: deps.Keycloak,
		opa:      deps.OPA,
	}

	// Setup middleware chain
	s.setupMiddleware()

	// Register routes
	s.registerRoutes()

	return s
}

// setupMiddleware configures the middleware chain in the correct order.
// Order: RequestID -> RealIP -> Logger -> Recoverer -> Timeout -> Auth -> Locale -> RateLimit
func (s *Server) setupMiddleware() {
	// 1. RequestID - Generate unique request ID
	s.router.Use(chimiddleware.RequestID)

	// 2. RealIP - Extract real IP from headers
	s.router.Use(chimiddleware.RealIP)

	// 3. Logger - Structured logging for requests
	s.router.Use(chimiddleware.RequestLogger(&slogLogFormatter{logger: s.logger}))

	// 4. Recoverer - Panic recovery
	s.router.Use(chimiddleware.Recoverer)

	// 5. Timeout - Request timeout (60 seconds default)
	s.router.Use(chimiddleware.Timeout(60 * time.Second))

	// 6. Auth - JWT validation via Keycloak
	if s.keycloak != nil {
		s.router.Use(middleware.AuthMiddleware(s.keycloak, s.logger))
	}

	// 7. Locale - Extract user locale from request
	s.router.Use(middleware.LocaleMiddleware(s.logger))

	// 8. RateLimit - Rate limiting per user
	if s.cache != nil {
		s.router.Use(middleware.RateLimitMiddleware(s.cache, s.logger, 60)) // 60 req/min default
	}

	// Additional standard middleware
	s.router.Use(chimiddleware.AllowContentType("application/json", "multipart/form-data"))
	s.router.Use(chimiddleware.CleanPath)
	s.router.Use(chimiddleware.StripSlashes)
}

// registerRoutes mounts all API routes.
func (s *Server) registerRoutes() {
	// Health check endpoint (no auth required)
	s.router.Get("/health", s.handleHealth)

	// Ready endpoint (no auth required)
	s.router.Get("/ready", s.handleReady)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Conversational BI endpoints (Module A)
		r.Route("/bi", func(r chi.Router) {
			r.Post("/query", s.handleBIQuery)
			r.Post("/explain", s.handleBIExplain)
		})

		// AI Accountant endpoints (Module B)
		r.Route("/accountant", func(r chi.Router) {
			r.Post("/ocr", s.handleOCRProcess)
			r.Post("/ledger-map", s.handleLedgerMap)
			r.Post("/sync", s.handleTallySync)
		})

		// Reports endpoints (Module C)
		r.Route("/reports", func(r chi.Router) {
			r.Get("/", s.handleListReports)
			r.Get("/{report_id}", s.handleGetReport)
			r.Post("/generate", s.handleGenerateReport)
		})

		// Search Analytics endpoints (Module D)
		r.Route("/analytics", func(r chi.Router) {
			r.Post("/search", s.handleAnalyticsSearch)
			r.Get("/insights", s.handleGetInsights)
		})

		// User preferences endpoints
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", s.handleGetCurrentUser)
			r.Put("/me/preferences", s.handleUpdatePreferences)
		})

		// Admin endpoints
		r.Route("/admin", func(r chi.Router) {
			r.Get("/audit-logs", s.handleGetAuditLogs)
			r.Post("/quarantine/{id}/reprocess", s.handleReprocessQuarantine)
		})
	})
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.getServerPort())

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	s.logger.Info("starting API server",
		slog.String("address", addr),
		slog.String("environment", string(s.config.App.Environment)),
	)

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("server listen error: %w", err)
		}
	}()

	// Wait for either server error or context cancellation
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.logger.Info("shutting down server due to context cancellation")
		return s.Shutdown(context.Background())
	}
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("shutting down API server")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("server shutdown error", slog.Any("error", err))
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("API server shutdown complete")
	return nil
}

// Router returns the chi router for testing purposes.
func (s *Server) Router() *chi.Mux {
	return s.router
}

// getServerPort returns the server port from config or default.
func (s *Server) getServerPort() int {
	// Check for API_PORT env var or use default
	// This could be extended to read from config
	return 8080 // Default API port
}

// ============================================================================
// HTTP Handlers
// ============================================================================

// handleHealth handles the health check endpoint.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleReady handles the readiness check endpoint.
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check database connection
	if s.db != nil {
		if err := s.db.Ping(ctx); err != nil {
			s.logger.Error("readiness check: database ping failed", slog.Any("error", err))
			s.writeError(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
	}

	// Check cache connection
	if s.cache != nil {
		if err := s.cache.Ping(ctx); err != nil {
			s.logger.Error("readiness check: cache ping failed", slog.Any("error", err))
			s.writeError(w, http.StatusServiceUnavailable, "cache unavailable")
			return
		}
	}

	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": map[string]bool{
			"database": s.db != nil,
			"cache":    s.cache != nil,
		},
	}

	s.writeJSON(w, http.StatusOK, response)
}

// Placeholder handlers - to be implemented by specific modules

func (s *Server) handleBIQuery(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "BI query endpoint not implemented")
}

func (s *Server) handleBIExplain(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "BI explain endpoint not implemented")
}

func (s *Server) handleOCRProcess(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "OCR process endpoint not implemented")
}

func (s *Server) handleLedgerMap(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Ledger map endpoint not implemented")
}

func (s *Server) handleTallySync(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Tally sync endpoint not implemented")
}

func (s *Server) handleListReports(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "List reports endpoint not implemented")
}

func (s *Server) handleGetReport(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Get report endpoint not implemented")
}

func (s *Server) handleGenerateReport(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Generate report endpoint not implemented")
}

func (s *Server) handleAnalyticsSearch(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Analytics search endpoint not implemented")
}

func (s *Server) handleGetInsights(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Get insights endpoint not implemented")
}

func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Get current user endpoint not implemented")
}

func (s *Server) handleUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Update preferences endpoint not implemented")
}

func (s *Server) handleGetAuditLogs(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Get audit logs endpoint not implemented")
}

func (s *Server) handleReprocessQuarantine(w http.ResponseWriter, r *http.Request) {
	s.writeError(w, http.StatusNotImplemented, "Reprocess quarantine endpoint not implemented")
}

// ============================================================================
// Helper Functions
// ============================================================================

// writeJSON writes a JSON response.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to write JSON response", slog.Any("error", err))
	}
}

// writeError writes an error response.
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}

// Redundant Helpers Removed (using encoding/json directly)

// ============================================================================
// Logging Formatter
// ============================================================================

// slogLogFormatter implements chi's LogFormatter interface using slog.
type slogLogFormatter struct {
	logger *slog.Logger
}

// NewLogEntry creates a new log entry for the request.
// NewLogEntry creates a new log entry for the request.
func (f *slogLogFormatter) NewLogEntry(r *http.Request) chimiddleware.LogEntry {
	return &slogLogEntry{
		logger: f.logger,
		r:      r,
	}
}

// slogLogEntry implements chi's LogEntry interface.
type slogLogEntry struct {
	logger *slog.Logger
	r      *http.Request
}

// Write logs the response status and details.
func (e *slogLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.logger.Info("request completed",
		slog.String("method", e.r.Method),
		slog.String("path", e.r.URL.Path),
		slog.Int("status", status),
		slog.Int("bytes", bytes),
		slog.Duration("elapsed", elapsed),
		slog.String("request_id", chimiddleware.GetReqID(e.r.Context())),
		slog.String("remote_addr", e.r.RemoteAddr),
	)
}

// Panic logs panic information.
func (e *slogLogEntry) Panic(v interface{}, stack []byte) {
	e.logger.Error("request panic",
		slog.Any("panic", v),
		slog.String("stack", string(stack)),
		slog.String("request_id", chimiddleware.GetReqID(e.r.Context())),
	)
}
