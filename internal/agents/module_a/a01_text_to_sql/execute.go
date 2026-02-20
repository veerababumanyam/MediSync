// Package a01_text_to_sql provides the text-to-SQL agent subcomponents.
//
// This file implements SQL execution with the medisync_readonly role.
package a01_text_to_sql

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLExecutor executes validated SQL queries with the readonly role.
type SQLExecutor struct {
	readOnlyPool *pgxpool.Pool
	logger       *slog.Logger
	timeout      time.Duration
	maxRows      int
}

// SQLExecutorConfig holds configuration for the executor.
type SQLExecutorConfig struct {
	ReadOnlyPool *pgxpool.Pool
	Logger       *slog.Logger
	Timeout      time.Duration
	MaxRows      int
}

// NewSQLExecutor creates a new SQL executor.
func NewSQLExecutor(cfg SQLExecutorConfig) *SQLExecutor {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRows == 0 {
		cfg.MaxRows = 10000
	}

	return &SQLExecutor{
		readOnlyPool: cfg.ReadOnlyPool,
		logger:       cfg.Logger.With("component", "sql_executor"),
		timeout:      cfg.Timeout,
		maxRows:      cfg.MaxRows,
	}
}

// ExecutionResult contains the result of SQL execution.
type ExecutionResult struct {
	Columns        []ColumnInfo `json:"columns"`
	Rows           []Row        `json:"rows"`
	RowCount       int          `json:"row_count"`
	ExecutionTime  int64        `json:"execution_time_ms"`
	Truncated      bool         `json:"truncated"`
	ErrorMessage   string       `json:"error_message,omitempty"`
	StatementID    string       `json:"statement_id,omitempty"`
}

// ColumnInfo contains metadata about a result column.
type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Row represents a single result row.
type Row map[string]interface{}

// ExecuteRequest contains the request for SQL execution.
type ExecuteRequest struct {
	SQL        string        `json:"sql"`
	Parameters []interface{} `json:"parameters,omitempty"`
	QueryID    string        `json:"query_id,omitempty"`
}

// Execute executes a validated SQL query.
func (e *SQLExecutor) Execute(ctx context.Context, req ExecuteRequest) (*ExecutionResult, error) {
	startTime := time.Now()
	e.logger.Debug("executing SQL", "sql_length", len(req.SQL), "params", len(req.Parameters))

	// Create timeout context
	execCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	result := &ExecutionResult{
		Columns: []ColumnInfo{},
		Rows:    []Row{},
	}

	// Execute the query
	var rows pgx.Rows
	var err error

	if len(req.Parameters) > 0 {
		rows, err = e.readOnlyPool.Query(execCtx, req.SQL, req.Parameters...)
	} else {
		rows, err = e.readOnlyPool.Query(execCtx, req.SQL)
	}

	if err != nil {
		result.ErrorMessage = err.Error()
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		e.logger.Warn("SQL execution failed", "error", err, "sql", req.SQL)
		return result, fmt.Errorf("SQL execution failed: %w", err)
	}
	defer rows.Close()

	// Get column information
	fieldDescriptions := rows.FieldDescriptions()
	for _, fd := range fieldDescriptions {
		result.Columns = append(result.Columns, ColumnInfo{
			Name: string(fd.Name),
			Type: fmt.Sprintf("oid:%d", uint32(fd.DataTypeOID)),
		})
	}

	// Fetch rows
	rowCount := 0
	for rows.Next() {
		if rowCount >= e.maxRows {
			result.Truncated = true
			break
		}

		values, err := rows.Values()
		if err != nil {
			e.logger.Warn("failed to get row values", "error", err)
			continue
		}

		row := make(Row)
		for i, col := range result.Columns {
			row[col.Name] = values[i]
		}
		result.Rows = append(result.Rows, row)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		result.ErrorMessage = err.Error()
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result, fmt.Errorf("row iteration error: %w", err)
	}

	result.RowCount = len(result.Rows)
	result.ExecutionTime = time.Since(startTime).Milliseconds()

	e.logger.Info("SQL execution complete",
		"row_count", result.RowCount,
		"execution_time_ms", result.ExecutionTime,
		"truncated", result.Truncated)

	return result, nil
}

// ExecuteWithRetry executes SQL with automatic retry on transient errors.
func (e *SQLExecutor) ExecuteWithRetry(ctx context.Context, req ExecuteRequest, maxRetries int) (*ExecutionResult, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := e.Execute(ctx, req)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Check if error is retryable
		if !e.isRetryableError(err) {
			return result, err
		}

		e.logger.Warn("retrying SQL execution",
			"attempt", attempt+1,
			"max_retries", maxRetries,
			"error", err)

		// Exponential backoff
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRetryableError checks if an error is retryable.
func (e *SQLExecutor) isRetryableError(err error) bool {
	// Check for connection errors, timeouts, etc.
	errStr := err.Error()
	retryablePatterns := []string{
		"connection reset",
		"connection refused",
		"timeout",
		"deadlock",
		"try again",
	}

	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// Explain returns the execution plan for a query.
func (e *SQLExecutor) Explain(ctx context.Context, sql string) (string, error) {
	explainSQL := fmt.Sprintf("EXPLAIN (FORMAT JSON) %s", sql)

	var result string
	err := e.readOnlyPool.QueryRow(ctx, explainSQL).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("explain failed: %w", err)
	}

	return result, nil
}

// EstimateRowCount estimates the number of rows a query will return.
func (e *SQLExecutor) EstimateRowCount(ctx context.Context, sql string) (int64, error) {
	explainSQL := fmt.Sprintf("EXPLAIN (FORMAT JSON) %s", sql)

	var result string
	err := e.readOnlyPool.QueryRow(ctx, explainSQL).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("explain failed: %w", err)
	}

	// Parse the explain output to get estimated rows
	var explainOutput []struct {
		Plan struct {
			EstimatedRows float64 `json:"Plan Rows"`
		} `json:"Plan"`
	}

	if err := json.Unmarshal([]byte(result), &explainOutput); err != nil {
		return 0, fmt.Errorf("failed to parse explain output: %w", err)
	}

	if len(explainOutput) > 0 {
		return int64(explainOutput[0].Plan.EstimatedRows), nil
	}

	return 0, nil
}

// CheckConnection verifies the database connection is working.
func (e *SQLExecutor) CheckConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return e.readOnlyPool.Ping(ctx)
}

// GetPoolStats returns statistics about the connection pool.
func (e *SQLExecutor) GetPoolStats() map[string]interface{} {
	stat := e.readOnlyPool.Stat()
	return map[string]interface{}{
		"acquired_conns":       stat.AcquiredConns(),
		"idle_conns":           stat.IdleConns(),
		"total_conns":          stat.TotalConns(),
		"max_conns":            stat.MaxConns(),
		"empty_acquire_count":  stat.EmptyAcquireCount(),
		"canceled_acquire_count": stat.CanceledAcquireCount(),
	}
}

// ToJSON serializes the execution result.
func (r *ExecutionResult) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// ToCSV converts the result to CSV format.
func (r *ExecutionResult) ToCSV() string {
	if len(r.Columns) == 0 || len(r.Rows) == 0 {
		return ""
	}

	var csv string

	// Header row
	for i, col := range r.Columns {
		if i > 0 {
			csv += ","
		}
		csv += col.Name
	}
	csv += "\n"

	// Data rows
	for _, row := range r.Rows {
		for i, col := range r.Columns {
			if i > 0 {
				csv += ","
			}
			if val, ok := row[col.Name]; ok {
				csv += formatCSVValue(val)
			}
		}
		csv += "\n"
	}

	return csv
}

// formatCSVValue formats a value for CSV output.
func formatCSVValue(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		// Escape quotes and wrap in quotes if needed
		if contains(v, ",") || contains(v, "\"") || contains(v, "\n") {
			return "\"" + replaceAll(v, "\"", "\"\"") + "\""
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}

// replaceAll replaces all occurrences of old with new in s.
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := findIndex(s, old)
		if idx == -1 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}

// findIndex finds the index of substr in s.
func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
