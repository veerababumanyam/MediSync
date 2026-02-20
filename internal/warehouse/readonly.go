// Package warehouse provides read-only database access for AI agents.
//
// This file provides the ReadOnlyClient struct which enforces that all queries
// use the medisync_readonly role and are SELECT-only. All queries are logged
// for audit purposes.
//
// Security:
//   - All queries must be SELECT statements (validated before execution)
//   - Uses medisync_readonly database role (configured via DSN)
//   - All queries are logged for audit trail
//
// Usage:
//
//	client, err := warehouse.NewReadOnlyClient(ctx, pool, logger)
//	if err != nil {
//	    log.Fatal("Failed to create readonly client:", err)
//	}
//	defer client.Close()
//
//	rows, err := client.ExecuteQuery(ctx, "SELECT * FROM patients WHERE id = $1", patientID)
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Forbidden SQL keywords that indicate write operations.
var forbiddenKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE",
	"GRANT", "REVOKE", "EXEC", "EXECUTE", "CALL", "COPY", "VACUUM",
	"REINDEX", "CLUSTER", "LOCK", "REFRESH", "REASSIGN", "ALTER",
}

// sqlValidationRegex matches SELECT statements.
var selectOnlyRegex = regexp.MustCompile(`(?i)^\s*(SELECT|WITH)\s`)

// ReadOnlyClient provides read-only database access with query validation.
type ReadOnlyClient struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// QueryAuditEntry represents an audit log entry for executed queries.
type QueryAuditEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	SQL         string    `json:"sql"`
	Params      []any     `json:"params,omitempty"`
	DurationMs  int64     `json:"duration_ms"`
	RowCount    int64     `json:"row_count"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	QuerySource string    `json:"query_source,omitempty"` // e.g., "A-01", "D-04"
}

// ReadOnlyClientConfig holds configuration for the read-only client.
type ReadOnlyClientConfig struct {
	// Pool is the PostgreSQL connection pool (must use medisync_readonly role).
	Pool *pgxpool.Pool

	// Logger is the structured logger for query auditing.
	Logger *slog.Logger
}

// NewReadOnlyClient creates a new read-only database client.
// The provided pool must be configured with the medisync_readonly role.
func NewReadOnlyClient(ctx context.Context, pool *PostgresPool, logger *slog.Logger) (*ReadOnlyClient, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	// Verify the connection is working
	if err := pool.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("warehouse: failed to verify database connection: %w", err)
	}

	logger.Info("read-only database client initialized")

	return &ReadOnlyClient{
		pool:   pool.Pool(),
		logger: logger,
	}, nil
}

// NewReadOnlyClientFromPool creates a read-only client from an existing pgxpool.Pool.
func NewReadOnlyClientFromPool(pool *pgxpool.Pool, logger *slog.Logger) (*ReadOnlyClient, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &ReadOnlyClient{
		pool:   pool,
		logger: logger,
	}, nil
}

// Close releases resources associated with the client.
// Note: This does not close the underlying pool.
func (c *ReadOnlyClient) Close() {
	c.logger.Debug("read-only client closed")
}

// ValidateSQL checks if the SQL query is SELECT-only and safe to execute.
func (c *ReadOnlyClient) ValidateSQL(sql string) error {
	// Normalize the SQL
	normalizedSQL := strings.TrimSpace(strings.ToUpper(sql))

	// Check for SELECT or WITH (CTE) statement
	if !selectOnlyRegex.MatchString(sql) {
		return fmt.Errorf("warehouse: only SELECT queries are allowed, got: %s", c.truncateSQL(sql, 50))
	}

	// Check for forbidden keywords
	for _, keyword := range forbiddenKeywords {
		// Use word boundary matching to avoid false positives
		pattern := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, keyword))
		if pattern.MatchString(normalizedSQL) {
			return fmt.Errorf("warehouse: forbidden keyword '%s' detected in query", keyword)
		}
	}

	return nil
}

// ExecuteQuery executes a SELECT query and returns the rows.
// The query is validated to ensure it is SELECT-only before execution.
// All queries are logged for audit purposes.
func (c *ReadOnlyClient) ExecuteQuery(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	startTime := time.Now()

	// Validate SQL is SELECT-only
	if err := c.ValidateSQL(sql); err != nil {
		c.logQueryAudit(ctx, QueryAuditEntry{
			Timestamp:  startTime,
			SQL:        sql,
			Params:     args,
			DurationMs: time.Since(startTime).Milliseconds(),
			Success:    false,
			Error:      err.Error(),
		})
		return nil, err
	}

	// Execute the query
	rows, err := c.pool.Query(ctx, sql, args...)

	duration := time.Since(startTime).Milliseconds()

	// Log for audit
	c.logQueryAudit(ctx, QueryAuditEntry{
		Timestamp:  startTime,
		SQL:        sql,
		Params:     args,
		DurationMs: duration,
		Success:    err == nil,
		Error:      c.errorToString(err),
	})

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to execute query: %w", err)
	}

	return rows, nil
}

// ExecuteQueryWithAudit executes a query and includes additional audit metadata.
func (c *ReadOnlyClient) ExecuteQueryWithAudit(ctx context.Context, entry QueryAuditEntry, args ...any) (pgx.Rows, error) {
	startTime := time.Now()
	entry.Timestamp = startTime
	entry.Params = args

	// Validate SQL is SELECT-only
	if err := c.ValidateSQL(entry.SQL); err != nil {
		entry.DurationMs = time.Since(startTime).Milliseconds()
		entry.Success = false
		entry.Error = err.Error()
		c.logQueryAudit(ctx, entry)
		return nil, err
	}

	// Execute the query
	rows, err := c.pool.Query(ctx, entry.SQL, args...)

	entry.DurationMs = time.Since(startTime).Milliseconds()
	entry.Success = err == nil
	entry.Error = c.errorToString(err)

	// Get row count if successful
	if err == nil {
		// Note: rows.CommandTag() might not have accurate count until fully consumed
		// We'll log 0 here and the caller can update
		entry.RowCount = 0
	}

	c.logQueryAudit(ctx, entry)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to execute query: %w", err)
	}

	return rows, nil
}

// QueryRow executes a query that returns at most one row.
func (c *ReadOnlyClient) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	startTime := time.Now()

	// Validate SQL is SELECT-only
	if err := c.ValidateSQL(sql); err != nil {
		c.logQueryAudit(ctx, QueryAuditEntry{
			Timestamp:  startTime,
			SQL:        sql,
			Params:     args,
			DurationMs: time.Since(startTime).Milliseconds(),
			Success:    false,
			Error:      err.Error(),
		})
		// Return a row that will error on Scan
		return &errorRow{err: err}
	}

	row := c.pool.QueryRow(ctx, sql, args...)

	c.logQueryAudit(ctx, QueryAuditEntry{
		Timestamp:  startTime,
		SQL:        sql,
		Params:     args,
		DurationMs: time.Since(startTime).Milliseconds(),
		Success:    true,
		RowCount:   1, // Assume 1 row for QueryRow
	})

	return row
}

// Ping checks if the database connection is healthy.
func (c *ReadOnlyClient) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// logQueryAudit logs a query execution for audit purposes.
func (c *ReadOnlyClient) logQueryAudit(ctx context.Context, entry QueryAuditEntry) {
	// Sanitize SQL for logging (truncate if too long)
	sql := c.truncateSQL(entry.SQL, 500)

	logAttrs := []slog.Attr{
		slog.String("sql", sql),
		slog.Int64("duration_ms", entry.DurationMs),
		slog.Bool("success", entry.Success),
	}

	if entry.Error != "" {
		logAttrs = append(logAttrs, slog.String("error", entry.Error))
	}

	if entry.SessionID != "" {
		logAttrs = append(logAttrs, slog.String("session_id", entry.SessionID))
	}

	if entry.UserID != "" {
		logAttrs = append(logAttrs, slog.String("user_id", entry.UserID))
	}

	if entry.QuerySource != "" {
		logAttrs = append(logAttrs, slog.String("query_source", entry.QuerySource))
	}

	if entry.RowCount > 0 {
		logAttrs = append(logAttrs, slog.Int64("row_count", entry.RowCount))
	}

	if entry.Success {
		c.logger.LogAttrs(ctx, slog.LevelInfo, "query executed", logAttrs...)
	} else {
		c.logger.LogAttrs(ctx, slog.LevelWarn, "query failed", logAttrs...)
	}
}

// truncateSQL truncates a SQL string for logging purposes.
func (c *ReadOnlyClient) truncateSQL(sql string, maxLen int) string {
	sql = strings.TrimSpace(sql)
	if len(sql) <= maxLen {
		return sql
	}
	return sql[:maxLen] + "..."
}

// errorToString converts an error to a string, returning empty string if nil.
func (c *ReadOnlyClient) errorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// errorRow is a pgx.Row implementation that always returns an error on Scan.
type errorRow struct {
	err error
}

func (r *errorRow) Scan(dest ...any) error {
	return r.err
}
