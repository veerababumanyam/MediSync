// Package a02_sql_correction provides the SQL self-correction agent.
//
// This agent analyzes SQL errors and attempts to correct queries automatically.
package a02_sql_correction

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

// AgentID is the unique identifier for this agent.
const AgentID = "a-02-sql-correction"

// Agent implements the SQL self-correction agent.
type Agent struct {
	id          string
	logger      *slog.Logger
	patterns    *ErrorPatternRegistry
	maxRetries  int
	corrections map[string]Correction
	cacheMu     sync.RWMutex
}

// AgentConfig holds configuration for the agent.
type AgentConfig struct {
	Logger     *slog.Logger
	MaxRetries int
}

// New creates a new SQL self-correction agent.
func New(cfg AgentConfig) *Agent {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}

	return &Agent{
		id:          AgentID,
		logger:      cfg.Logger.With("agent", AgentID),
		patterns:    NewErrorPatternRegistry(),
		maxRetries:  cfg.MaxRetries,
		corrections: make(map[string]Correction),
	}
}

// CorrectionRequest contains the request for SQL correction.
type CorrectionRequest struct {
	SQL         string `json:"sql"`
	Error       string `json:"error"`
	RetryCount  int    `json:"retry_count"`
	Query       string `json:"query,omitempty"`
	SchemaHints []string `json:"schema_hints,omitempty"`
}

// CorrectionResponse contains the corrected SQL.
type CorrectionResponse struct {
	OriginalSQL   string `json:"original_sql"`
	CorrectedSQL  string `json:"corrected_sql"`
	Correction    string `json:"correction"`
	RetryCount    int    `json:"retry_count"`
	ShouldRetry   bool   `json:"should_retry"`
	Confidence    float64 `json:"confidence"`
	ErrorType     string `json:"error_type"`
}

// Correction represents a single correction applied.
type Correction struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Original    string `json:"original"`
	Corrected   string `json:"corrected"`
}

// AgentCard returns the ADK agent card for discovery.
func (a *Agent) AgentCard() map[string]interface{} {
	return map[string]interface{}{
		"id":          AgentID,
		"name":        "SQL Self-Correction Agent",
		"description": "Analyzes SQL errors and attempts automatic correction",
		"capabilities": []string{
			"sql-correction",
			"error-analysis",
			"retry-logic",
		},
		"version": "1.0.0",
	}
}

// Correct attempts to correct a SQL query that produced an error.
func (a *Agent) Correct(ctx context.Context, req CorrectionRequest) (*CorrectionResponse, error) {
	a.logger.Debug("analyzing SQL error",
		"sql_length", len(req.SQL),
		"error", req.Error,
		"retry_count", req.RetryCount)

	// Check if max retries exceeded
	if req.RetryCount >= a.maxRetries {
		return &CorrectionResponse{
			OriginalSQL: req.SQL,
			ShouldRetry: false,
			RetryCount:  req.RetryCount,
			Confidence:  0,
		}, nil
	}

	// Analyze the error
	errorType := a.analyzeError(req.Error)
	a.logger.Info("error analyzed", "error_type", errorType)

	// Get correction strategy
	strategy := a.patterns.GetStrategy(errorType)
	if strategy == nil {
		return &CorrectionResponse{
			OriginalSQL: req.SQL,
			ShouldRetry: false,
			RetryCount:  req.RetryCount,
			Confidence:  0,
			ErrorType:   errorType,
		}, nil
	}

	// Apply correction
	correctedSQL, correction, confidence := strategy.Apply(req)

	response := &CorrectionResponse{
		OriginalSQL:  req.SQL,
		CorrectedSQL: correctedSQL,
		Correction:   correction,
		RetryCount:   req.RetryCount + 1,
		ShouldRetry:  correctedSQL != "" && confidence > 0.5,
		Confidence:   confidence,
		ErrorType:    errorType,
	}

	// Cache the correction for learning
	a.cacheCorrection(req.Error, Correction{
		Type:        errorType,
		Description: correction,
		Original:    req.SQL,
		Corrected:   correctedSQL,
	})

	a.logger.Info("correction applied",
		"error_type", errorType,
		"confidence", confidence,
		"should_retry", response.ShouldRetry)

	return response, nil
}

// analyzeError determines the type of SQL error.
func (a *Agent) analyzeError(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	// Check for specific error patterns
	patterns := []struct {
		pattern   string
		errorType string
	}{
		{"column", "column_not_found"},
		{"does not exist", "relation_not_found"},
		{"syntax error", "syntax_error"},
		{"invalid input syntax", "type_mismatch"},
		{"function", "function_error"},
		{"permission denied", "permission_error"},
		{"timeout", "timeout_error"},
		{"ambiguous", "ambiguous_reference"},
		{"group by", "group_by_error"},
		{"aggregate", "aggregate_error"},
		{"join", "join_error"},
		{"null value", "null_constraint"},
		{"unique violation", "unique_violation"},
		{"foreign key", "foreign_key_error"},
	}

	for _, p := range patterns {
		if strings.Contains(errLower, p.pattern) {
			return p.errorType
		}
	}

	return "unknown"
}

// cacheCorrection stores a correction for learning.
func (a *Agent) cacheCorrection(err string, correction Correction) {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	a.corrections[err] = correction
}

// GetCorrections returns all cached corrections.
func (a *Agent) GetCorrections() map[string]Correction {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	result := make(map[string]Correction)
	for k, v := range a.corrections {
		result[k] = v
	}
	return result
}

// Retry executes the correction-retry loop.
func (a *Agent) Retry(ctx context.Context, req CorrectionRequest, executor SQLExecutor) (string, error) {
	currentSQL := req.SQL
	retryCount := req.RetryCount

	for retryCount < a.maxRetries {
		result, err := a.Correct(ctx, CorrectionRequest{
			SQL:        currentSQL,
			Error:      req.Error,
			RetryCount: retryCount,
		})
		if err != nil {
			return "", err
		}

		if !result.ShouldRetry {
			break
		}

		// Try executing the corrected SQL
		execResult, err := executor.Execute(ctx, result.CorrectedSQL)
		if err == nil {
			return execResult, nil
		}

		// Update for next iteration
		currentSQL = result.CorrectedSQL
		req.Error = err.Error()
		retryCount = result.RetryCount
	}

	return "", fmt.Errorf("max retries (%d) exceeded", a.maxRetries)
}

// SQLExecutor defines the interface for SQL execution.
type SQLExecutor interface {
	Execute(ctx context.Context, sql string) (string, error)
}

// ToJSON serializes the correction response.
func (r *CorrectionResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// FromJSON deserializes a correction response.
func FromJSON(data string) (*CorrectionResponse, error) {
	var resp CorrectionResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
