// Package warehouse provides database connectivity for the MediSync data warehouse.
//
// This file provides the PostgresPool struct which wraps pgxpool.Pool with
// proper connection pool settings and health check capabilities.
//
// Usage:
//
//	pool, err := warehouse.NewPostgresPool(ctx, dsn)
//	if err != nil {
//	    log.Fatal("Failed to create connection pool:", err)
//	}
//	defer pool.Close()
//
//	// Check health
//	if err := pool.HealthCheck(ctx); err != nil {
//	    log.Warn("Database health check failed:", err)
//	}
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Default pool configuration values.
const (
	// DefaultMaxConns is the default maximum number of connections in the pool.
	DefaultMaxConns = 25

	// DefaultMinConns is the default minimum number of idle connections.
	DefaultMinConns = 5

	// DefaultMaxConnLifetime is the default maximum lifetime of a connection.
	DefaultMaxConnLifetime = 5 * time.Minute

	// DefaultMaxConnIdleTime is the default maximum idle time of a connection.
	DefaultMaxConnIdleTime = 1 * time.Minute

	// DefaultHealthCheckPeriod is the default period between health checks.
	DefaultHealthCheckPeriod = 1 * time.Minute
)

// PoolConfig holds configuration for the PostgreSQL connection pool.
type PoolConfig struct {
	// DSN is the PostgreSQL connection string.
	DSN string

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int32

	// MinConns is the minimum number of idle connections in the pool.
	MinConns int32

	// MaxConnLifetime is the maximum lifetime of a connection.
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum idle time of a connection.
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is the period between background health checks.
	HealthCheckPeriod time.Duration

	// Logger is the structured logger for connection pool events.
	Logger *slog.Logger
}

// PostgresPool wraps pgxpool.Pool with additional functionality.
type PostgresPool struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	config *PoolConfig
}

// NewPostgresPool creates a new PostgreSQL connection pool with default settings.
func NewPostgresPool(ctx context.Context, dsn string) (*PostgresPool, error) {
	return NewPostgresPoolWithConfig(ctx, &PoolConfig{
		DSN:               dsn,
		MaxConns:          DefaultMaxConns,
		MinConns:          DefaultMinConns,
		MaxConnLifetime:   DefaultMaxConnLifetime,
		MaxConnIdleTime:   DefaultMaxConnIdleTime,
		HealthCheckPeriod: DefaultHealthCheckPeriod,
		Logger:            slog.Default(),
	})
}

// NewPostgresPoolWithConfig creates a new PostgreSQL connection pool with custom settings.
func NewPostgresPoolWithConfig(ctx context.Context, cfg *PoolConfig) (*PostgresPool, error) {
	if cfg == nil || cfg.DSN == "" {
		return nil, fmt.Errorf("warehouse: DSN is required")
	}

	// Set defaults
	if cfg.MaxConns == 0 {
		cfg.MaxConns = DefaultMaxConns
	}
	if cfg.MinConns == 0 {
		cfg.MinConns = DefaultMinConns
	}
	if cfg.MaxConnLifetime == 0 {
		cfg.MaxConnLifetime = DefaultMaxConnLifetime
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = DefaultMaxConnIdleTime
	}
	if cfg.HealthCheckPeriod == 0 {
		cfg.HealthCheckPeriod = DefaultHealthCheckPeriod
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	// Parse pool config from DSN
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to parse DSN: %w", err)
	}

	// Apply custom pool settings
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Create the pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("warehouse: failed to ping database: %w", err)
	}

	cfg.Logger.Info("PostgreSQL connection pool created",
		slog.Int("max_conns", int(cfg.MaxConns)),
		slog.Int("min_conns", int(cfg.MinConns)),
		slog.Duration("max_conn_lifetime", cfg.MaxConnLifetime),
		slog.Duration("max_conn_idle_time", cfg.MaxConnIdleTime),
	)

	return &PostgresPool{
		pool:   pool,
		logger: cfg.Logger,
		config: cfg,
	}, nil
}

// Close closes the connection pool and releases all resources.
func (p *PostgresPool) Close() {
	if p.pool != nil {
		p.pool.Close()
		p.logger.Debug("PostgreSQL connection pool closed")
	}
}

// Pool returns the underlying pgxpool.Pool for direct access.
func (p *PostgresPool) Pool() *pgxpool.Pool {
	return p.pool
}

// HealthCheck verifies that the database connection is healthy.
func (p *PostgresPool) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.pool.Ping(ctx); err != nil {
		return fmt.Errorf("warehouse: health check failed: %w", err)
	}
	return nil
}

// Stats returns current connection pool statistics.
func (p *PostgresPool) Stats() *pgxpool.Stat {
	if p.pool == nil {
		return nil
	}
	return p.pool.Stat()
}

// Acquire acquires a connection from the pool.
func (p *PostgresPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to acquire connection: %w", err)
	}
	return conn, nil
}

// Exec executes a SQL statement.
func (p *PostgresPool) Exec(ctx context.Context, sql string, args ...any) (interface{}, error) {
	result, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to execute query: %w", err)
	}
	return result, nil
}

// Query executes a SQL query and returns rows.
func (p *PostgresPool) Query(ctx context.Context, sql string, args ...any) (interface{}, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to execute query: %w", err)
	}
	return rows, nil
}

// QueryRow executes a SQL query that returns at most one row.
func (p *PostgresPool) QueryRow(ctx context.Context, sql string, args ...any) interface {
	Scan(dest ...any) error
} {
	return p.pool.QueryRow(ctx, sql, args...)
}

// Begin starts a new transaction.
func (p *PostgresPool) Begin(ctx context.Context) (interface{}, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to begin transaction: %w", err)
	}
	return tx, nil
}

// SendBatch sends a batch of queries to the database.
// The batch parameter should be a *pgx.Batch.
func (p *PostgresPool) SendBatch(ctx context.Context, batch *pgx.Batch) (pgx.BatchResults, error) {
	if batch == nil {
		return nil, fmt.Errorf("warehouse: batch cannot be nil")
	}
	results := p.pool.SendBatch(ctx, batch)
	return results, nil
}

// Config returns the pool configuration.
func (p *PostgresPool) Config() *PoolConfig {
	return p.config
}
