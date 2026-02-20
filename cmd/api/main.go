// Package main provides the entry point for the MediSync API server.
//
// The API server handles all HTTP requests including:
// - Conversational BI queries
// - AI Accountant operations
// - Reports generation
// - Analytics search
//
// Usage:
//
//	go run ./cmd/api
//
// Environment variables:
//
//	DATABASE_URL      - PostgreSQL connection string
//	REDIS_URL         - Redis connection URL
//	KEYCLOAK_URL      - Keycloak server URL
//	KEYCLOAK_REALM    - Keycloak realm name
//	KEYCLOAK_CLIENT_ID - OAuth2 client ID
//	OPA_URL           - Open Policy Agent URL
//	API_PORT          - API server port (default: 8080)
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/medisync/medisync/internal/api"
	"github.com/medisync/medisync/internal/cache"
	"github.com/medisync/medisync/internal/config"
	"github.com/medisync/medisync/internal/warehouse"
)

func main() {
	// Initialize logger
	logger := setupLogger()
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// Log configuration (with sensitive values masked)
	cfg.LogConfig(logger)

	// Create context that listens for shutdown signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Initialize dependencies
	deps, err := initializeDependencies(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize dependencies",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	// Create API server
	server := api.NewServer(cfg, deps)

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-sigCh:
		logger.Info("received shutdown signal",
			slog.String("signal", sig.String()),
		)
		cancel()

	case err := <-errCh:
		logger.Error("server error",
			slog.Any("error", err),
		)
		cancel()
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error",
			slog.Any("error", err),
		)
	}

	// Close dependencies
	closeDependencies(deps, logger)

	logger.Info("API server stopped")
}

// setupLogger creates and configures the structured logger.
func setupLogger() *slog.Logger {
	// Determine log level from environment
	level := slog.LevelInfo
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		switch envLevel {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Use JSON format in production, text in development
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// initializeDependencies creates and initializes all required dependencies.
func initializeDependencies(cfg *config.Config, logger *slog.Logger) (*api.Dependencies, error) {
	deps := &api.Dependencies{}

	// Initialize database repository
	db, err := initDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	deps.DB = db

	// Initialize Redis cache
	cacheClient, err := initCache(cfg, logger)
	if err != nil {
		logger.Warn("Redis cache not available, using defaults",
			slog.Any("error", err),
		)
		// Cache is optional, continue without it
	} else {
		deps.Cache = cacheClient
	}

	// Initialize Keycloak client
	keycloakClient, err := initKeycloak(cfg, logger)
	if err != nil {
		logger.Warn("Keycloak client not available, auth will be limited",
			slog.Any("error", err),
		)
		// Keycloak is required for production, but optional for development
		if cfg.IsProduction() {
			return nil, fmt.Errorf("keycloak: %w", err)
		}
	} else {
		deps.Keycloak = keycloakClient
	}

	// Initialize OPA client
	opaClient, err := initOPA(cfg, logger)
	if err != nil {
		logger.Warn("OPA client not available, policy decisions will be limited",
			slog.Any("error", err),
		)
		// OPA is required for production, but optional for development
		if cfg.IsProduction() {
			return nil, fmt.Errorf("opa: %w", err)
		}
	} else {
		deps.OPA = opaClient
	}

	return deps, nil
}

// initDatabase initializes the database repository.
func initDatabase(cfg *config.Config, logger *slog.Logger) (*warehouse.Repo, error) {
	dsn := cfg.DatabaseDSN()

	db, err := warehouse.NewRepo(dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database repo: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		slog.String("host", cfg.Database.Host),
		slog.String("database", cfg.Database.Name),
	)

	return db, nil
}

// initCache initializes the Redis cache client.
func initCache(cfg *config.Config, logger *slog.Logger) (*cache.Client, error) {
	cacheConfig := map[string]interface{}{
		"url":      cfg.Redis.URL,
		"host":     cfg.Redis.Host,
		"port":     cfg.Redis.Port,
		"password": cfg.Redis.Password,
		"database": cfg.Redis.Database,
	}

	cacheClient, err := cache.NewClient(cacheConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cacheClient.Ping(ctx); err != nil {
		cacheClient.Close()
		return nil, fmt.Errorf("failed to ping cache: %w", err)
	}

	logger.Info("cache connection established",
		slog.String("host", cfg.Redis.Host),
		slog.Int("database", cfg.Redis.Database),
	)

	return cacheClient, nil
}

// initKeycloak initializes the Keycloak authentication client.
func initKeycloak(cfg *config.Config, logger *slog.Logger) (*api.KeycloakClient, error) {
	keycloakClient, err := api.NewKeycloakClient(&cfg.Keycloak, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Keycloak client: %w", err)
	}

	logger.Info("Keycloak client initialized",
		slog.String("url", cfg.Keycloak.URL),
		slog.String("realm", cfg.Keycloak.Realm),
	)

	return keycloakClient, nil
}

// initOPA initializes the Open Policy Agent client.
func initOPA(cfg *config.Config, logger *slog.Logger) (*api.OPAClient, error) {
	opaURL := os.Getenv("OPA_URL")
	if opaURL == "" {
		opaURL = "http://localhost:8181"
	}

	opaClient, err := api.NewOPAClient(&api.OPAConfig{
		URL:     opaURL,
		Timeout: 5 * time.Second,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create OPA client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := opaClient.Health(ctx); err != nil {
		logger.Warn("OPA health check failed, policy decisions may be limited",
			slog.Any("error", err),
		)
	}

	logger.Info("OPA client initialized",
		slog.String("url", opaURL),
	)

	return opaClient, nil
}

// closeDependencies closes all dependencies gracefully.
func closeDependencies(deps *api.Dependencies, logger *slog.Logger) {
	if deps.Cache != nil {
		if err := deps.Cache.Close(); err != nil {
			logger.Error("failed to close cache",
				slog.Any("error", err),
			)
		}
	}

	if deps.DB != nil {
		deps.DB.Close()
	}

	logger.Debug("dependencies closed")
}
