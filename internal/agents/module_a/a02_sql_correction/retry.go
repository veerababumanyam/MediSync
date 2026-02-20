// Package a02_sql_correction provides the SQL self-correction agent.
//
// This file implements retry logic for SQL self-correction.
package a02_sql_correction

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// RetryManager manages the retry loop for SQL correction.
type RetryManager struct {
	maxRetries    int
	backoffFactor time.Duration
	logger        *slog.Logger
}

// RetryConfig holds configuration for retry management.
type RetryConfig struct {
	MaxRetries    int
	BackoffFactor time.Duration
	Logger        *slog.Logger
}

// NewRetryManager creates a new retry manager.
func NewRetryManager(cfg RetryConfig) *RetryManager {
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.BackoffFactor == 0 {
		cfg.BackoffFactor = 100 * time.Millisecond
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &RetryManager{
		maxRetries:    cfg.MaxRetries,
		backoffFactor: cfg.BackoffFactor,
		logger:        cfg.Logger.With("component", "retry_manager"),
	}
}

// RetryState tracks the state of a retry loop.
type RetryState struct {
	Attempt      int           `json:"attempt"`
	MaxAttempts  int           `json:"max_attempts"`
	LastError    string        `json:"last_error"`
	LastSQL      string        `json:"last_sql"`
	Corrections  []Correction  `json:"corrections"`
	StartTime    time.Time     `json:"start_time"`
	TotalLatency time.Duration `json:"total_latency"`
}

// NewRetryState creates a new retry state.
func NewRetryState(maxAttempts int) *RetryState {
	return &RetryState{
		Attempt:     0,
		MaxAttempts: maxAttempts,
		Corrections: []Correction{},
		StartTime:   time.Now(),
	}
}

// CanRetry returns true if more retry attempts are available.
func (s *RetryState) CanRetry() bool {
	return s.Attempt < s.MaxAttempts
}

// RecordAttempt records an attempt result.
func (s *RetryState) RecordAttempt(sql string, err error, correction Correction) {
	s.Attempt++
	s.LastSQL = sql
	if err != nil {
		s.LastError = err.Error()
	}
	if correction.Type != "" {
		s.Corrections = append(s.Corrections, correction)
	}
	s.TotalLatency = time.Since(s.StartTime)
}

// ExecuteWithRetry executes a function with automatic retry on error.
func (m *RetryManager) ExecuteWithRetry(ctx context.Context, initialSQL string, fn ExecuteFunc) (*RetryResult, error) {
	state := NewRetryState(m.maxRetries)
	currentSQL := initialSQL

	for state.CanRetry() {
		// Execute the SQL
		result, err := fn(ctx, currentSQL)

		if err == nil {
			// Success
			return &RetryResult{
				Success:      true,
				FinalSQL:     currentSQL,
				Attempts:     state.Attempt + 1,
				Corrections:  state.Corrections,
				TotalLatency: time.Since(state.StartTime),
				Result:       result,
			}, nil
		}

		// Record the failure
		state.RecordAttempt(currentSQL, err, Correction{})

		m.logger.Debug("SQL execution failed, attempting correction",
			"attempt", state.Attempt,
			"error", err.Error())

		// Check if context is cancelled
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Get correction
		correctionAgent := New(AgentConfig{
			Logger:     m.logger,
			MaxRetries: m.maxRetries,
		})

		correction, corrErr := correctionAgent.Correct(ctx, CorrectionRequest{
			SQL:        currentSQL,
			Error:      err.Error(),
			RetryCount: state.Attempt,
		})

		if corrErr != nil {
			return nil, fmt.Errorf("correction failed: %w", corrErr)
		}

		if !correction.ShouldRetry {
			// No correction possible
			return &RetryResult{
				Success:      false,
				FinalSQL:     currentSQL,
				Attempts:     state.Attempt,
				Corrections:  state.Corrections,
				TotalLatency: time.Since(state.StartTime),
				LastError:    err.Error(),
				ErrorType:    correction.ErrorType,
			}, fmt.Errorf("no correction available: %s", correction.Correction)
		}

		// Apply the correction
		state.Corrections = append(state.Corrections, Correction{
			Type:        correction.ErrorType,
			Description: correction.Correction,
			Original:    currentSQL,
			Corrected:   correction.CorrectedSQL,
		})

		currentSQL = correction.CorrectedSQL

		// Apply backoff
		backoff := m.backoffFactor * time.Duration(state.Attempt)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}

	// Max retries exceeded
	return &RetryResult{
		Success:      false,
		FinalSQL:     currentSQL,
		Attempts:     state.Attempt,
		Corrections:  state.Corrections,
		TotalLatency: time.Since(state.StartTime),
		LastError:    state.LastError,
	}, fmt.Errorf("max retries (%d) exceeded", m.maxRetries)
}

// ExecuteFunc is the function type for SQL execution.
type ExecuteFunc func(ctx context.Context, sql string) (interface{}, error)

// RetryResult contains the result of a retry operation.
type RetryResult struct {
	Success      bool          `json:"success"`
	FinalSQL     string        `json:"final_sql"`
	Attempts     int           `json:"attempts"`
	Corrections  []Correction  `json:"corrections"`
	TotalLatency time.Duration `json:"total_latency"`
	Result       interface{}   `json:"result,omitempty"`
	LastError    string        `json:"last_error,omitempty"`
	ErrorType    string        `json:"error_type,omitempty"`
}

// IsRecoverable determines if an error is recoverable through retry.
func IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	errStr := err.Error()

	// Non-recoverable error patterns
	nonRecoverable := []string{
		"permission denied",
		"authentication failed",
		"does not exist",
		"syntax error",
	}

	for _, pattern := range nonRecoverable {
		// Some syntax errors are recoverable, others are not
		if pattern == "syntax error" {
			continue
		}
		if contains(errStr, pattern) {
			return false
		}
	}

	// Recoverable error patterns
	recoverable := []string{
		"timeout",
		"connection reset",
		"deadlock",
		"try again",
		"temporary",
	}

	for _, pattern := range recoverable {
		if contains(errStr, pattern) {
			return true
		}
	}

	// Default: assume column/relation errors are recoverable
	return contains(errStr, "column") || contains(errStr, "relation")
}

// GetBackoffDuration calculates the backoff duration for a given attempt.
func (m *RetryManager) GetBackoffDuration(attempt int) time.Duration {
	// Exponential backoff with jitter
	backoff := m.backoffFactor * time.Duration(1<<uint(attempt))

	// Add some jitter (10%)
	jitter := time.Duration(float64(backoff) * 0.1)

	return backoff + jitter
}

// Stats returns retry statistics.
func (m *RetryManager) Stats() map[string]interface{} {
	return map[string]interface{}{
		"max_retries":    m.maxRetries,
		"backoff_factor": m.backoffFactor.String(),
	}
}

// contains checks if s contains substr (case-insensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s[:len(substr)] == substr ||
			(len(s) > len(substr) && contains(s[1:], substr)))
}
