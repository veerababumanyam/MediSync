// Package main provides the ETL service entry point.
//
// The ETL service orchestrates data extraction, transformation, and loading
// from source systems (Tally ERP, HIMS) into the MediSync data warehouse.
// It runs scheduled sync jobs, exposes Prometheus metrics, and provides
// health check endpoints.
//
// Usage:
//
//	go run cmd/etl/main.go
//
// Environment Variables:
//   APP_ENV - Application environment (development, staging, production)
//   DATABASE_URL - PostgreSQL connection string
//   NATS_URL - NATS server URL
//   REDIS_URL - Redis connection URL
//   TALLY_HOST - Tally ERP server hostname
//   HIMS_API_URL - HIMS API base URL
//   HIMS_API_KEY - HIMS API authentication key
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/medisync/medisync/internal/config"
	"github.com/medisync/medisync/internal/etl/orchestrator"
	"github.com/medisync/medisync/internal/etl/tally"
	"github.com/medisync/medisync/internal/etl/hims"
	"github.com/medisync/medisync/internal/events"
	"github.com/medisync/medisync/internal/warehouse"
)

const (
	// ServiceName is the name of this service.
	ServiceName = "medisync-etl"

	// ServiceVersion is the version of this service.
	ServiceVersion = "1.0.0-alpha"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	cfg.LogConfig(slog.Default())

	// Set up structured logger
	logger := setupLogger(cfg)
	logger.Info("starting ETL service",
		slog.String("service", ServiceName),
		slog.String("version", ServiceVersion),
		slog.String("environment", string(cfg.App.Environment)),
	)

	// Create dependencies
	repo, err := warehouse.NewRepo(cfg.Database, logger)
	if err != nil {
		logger.Error("failed to create warehouse repository", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()

	// Test database connection
	if err := repo.Ping(context.Background()); err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create NATS publisher
	eventsPub, err := events.NewPublisher(cfg.NATS, logger)
	if err != nil {
		logger.Warn("failed to create NATS publisher", slog.String("error", err.Error()))
		// Continue without events
		eventsPub = nil
	} else {
		defer eventsPub.Close()
	}

	// Create Tally client
	tallyClient := tally.NewClient(cfg.Tally,
		tally.WithLogger(logger),
	)

	// Create HIMS client
	himsClient := hims.NewClient(cfg.HIMS,
		hims.WithLogger(logger),
	)

	// Create orchestrator
	orch, err := orchestrator.NewService(&orchestrator.ServiceConfig{
		Config:      cfg,
		Repo:        repo,
		TallyClient: tallyClient,
		HIMSClient:  himsClient,
		Events:      eventsPub,
		Logger:      logger,
	})
	if err != nil {
		logger.Error("failed to create orchestrator", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Set up HTTP server for metrics and health checks
	server := setupHTTPServer(cfg, repo, orch, logger)

	// Start HTTP server in background
	go func() {
		addr := fmt.Sprintf(":%d", cfg.ETL.MetricsPort)
		logger.Info("HTTP server listening", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	// Start orchestrator
	ctx := context.Background()
	if err := orch.Start(ctx); err != nil {
		logger.Error("failed to start orchestrator", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Wait for shutdown signal
	waitForShutdown(logger, orch, server, repo)

	logger.Info("ETL service stopped")
}

// setupLogger configures the structured logger.
func setupLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevelFromString(cfg.App.LogLevel),
	}

	switch cfg.App.LogFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// setupHTTPServer creates the HTTP server for metrics and health checks.
func setupHTTPServer(cfg *config.Config, repo *warehouse.Repo, orch *orchestrator.Service, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check database connection
		dbStatus := "ok"
		if err := repo.Ping(r.Context()); err != nil {
			dbStatus = "error"
		}

		health := map[string]interface{}{
			"status":         "up",
			"service":        ServiceName,
			"version":        ServiceVersion,
			"environment":    string(cfg.App.Environment),
			"database":       dbStatus,
			"orchestrator":   orch.GetStats(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	})

	// Ready endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if database is ready
		if err := repo.Ping(r.Context()); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	// Metrics endpoint (Prometheus)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Collect metrics
		metrics := collectMetrics(repo, orch)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(metrics))
	})

	// Sync status endpoint
	mux.HandleFunc("/api/v1/sync/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats, err := repo.GetSyncStats(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Trigger sync endpoint (for manual sync triggering)
	mux.HandleFunc("/api/v1/sync/trigger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request
		var req struct {
			Source string `json:"source"`
			Entity string `json:"entity"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Trigger sync (implementation would be in orchestrator)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "triggered",
			"source":  req.Source,
			"entity":  req.Entity,
		})
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ETL.MetricsPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

// waitForShutdown handles graceful shutdown.
func waitForShutdown(logger *slog.Logger, orch *orchestrator.Service, server *http.Server, repo *warehouse.Repo) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-sigChan
	logger.Info("received shutdown signal", slog.String("signal", sig.String()))

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop orchestrator first
	if err := orch.Stop(); err != nil {
		logger.Error("error stopping orchestrator", slog.String("error", err.Error()))
	}

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("error shutting down HTTP server", slog.String("error", err.Error()))
	}

	// Close database connection
	repo.Close()
}

// collectMetrics collects Prometheus-style metrics.
func collectMetrics(repo *warehouse.Repo, orch *orchestrator.Service) string {
	var metrics string

	// Go metrics would use prometheus/client_golang
	// This is a simplified version

	metrics += "# HELP medisync_etl_syncs_total Total number of ETL sync operations\n"
	metrics += "# TYPE medisync_etl_syncs_total counter\n"
	metrics += fmt.Sprintf("medisync_etl_syncs_total %d\n", getOrchestratorStat(orch, "total_syncs"))

	metrics += "\n# HELP medisync_etl_syncs_failed Total number of failed ETL sync operations\n"
	metrics += "# TYPE medisync_etl_syncs_failed counter\n"
	metrics += fmt.Sprintf("medisync_etl_syncs_failed %d\n", getOrchestratorStat(orch, "failed_syncs"))

	metrics += "\n# HELP medisync_etl_up Whether the ETL service is running\n"
	metrics += "# TYPE medisync_etl_up gauge\n"
	up := 1
	if getOrchestratorStat(orch, "is_running") == 0 {
		up = 0
	}
	metrics += fmt.Sprintf("medisync_etl_up %d\n", up)

	// Get sync stats
	stats, _ := repo.GetSyncStats(context.Background())
	if stats != nil {
		metrics += "\n# HELP medisync_etl_entities_total Total number of entities being synced\n"
		metrics += "# TYPE medisync_etl_entities_total gauge\n"
		metrics += fmt.Sprintf("medisync_etl_entities_total %d\n", stats.TotalEntities)

		metrics += "\n# HELP medisync_etl_entities_running Number of entities currently syncing\n"
		metrics += "# TYPE medisync_etl_entities_running gauge\n"
		metrics += fmt.Sprintf("medisync_etl_entities_running %d\n", stats.RunningSyncs)

		metrics += "\n# HELP medisync_etl_records_synced Total records synced to warehouse\n"
		metrics += "# TYPE medisync_etl_records_synced counter\n"
		metrics += fmt.Sprintf("medisync_etl_records_synced %d\n", stats.TotalRecords)
	}

	return metrics
}

// getOrchestratorStat extracts a stat from the orchestrator.
func getOrchestratorStat(orch *orchestrator.Service, key string) int64 {
	stats := orch.GetStats()
	if val, ok := stats[key]; ok {
		if v, ok := val.(int64); ok {
			return v
		}
	}
	return 0
}

// slogLevelFromString converts a log level string to slog.Level.
func slogLevelFromString(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
