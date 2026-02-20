// Package module_a provides the A-02 SQL Self-Correction Agent.
//
// This agent detects SQL execution errors, analyzes root causes, and attempts
// automatic correction with configurable retry logic.
//
// Usage:
//
//	agent := agent_a02.New(logger)
//	correctedSQL, err := agent.Correct(ctx, sql, errMsg, attempt)
package module_a

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/medisync/medisync/internal/agents/shared"
)

const (
	// AgentID is the unique identifier for this agent.
	AgentIDA02 = "a-02-sql-correction"

	// AgentNameA02 is the human-readable name.
	AgentNameA02 = "SQL Self-Correction Agent"

	// MaxRetries is the maximum number of correction attempts.
	MaxRetries = 3
)

// CorrectionAgent provides SQL correction functionality.
type CorrectionAgent struct {
	logger *slog.Logger
}

// CorrectionConfig holds configuration for the agent.
type CorrectionConfig struct {
	Logger *slog.Logger
}

// NewCorrectionAgent creates a new A-02 SQL Self-Correction Agent.
func NewCorrectionAgent(cfg *CorrectionConfig) *CorrectionAgent {
	if cfg == nil {
		cfg = &CorrectionConfig{}
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &CorrectionAgent{
		logger: logger.With(slog.String("agent", AgentIDA02)),
	}
}

// ID returns the agent identifier.
func (a *CorrectionAgent) ID() string {
	return AgentIDA02
}

// Name returns the agent name.
func (a *CorrectionAgent) Name() string {
	return AgentNameA02
}

// Health checks if the agent is healthy.
func (a *CorrectionAgent) Health(ctx context.Context) (*shared.AgentHealth, error) {
	return &shared.AgentHealth{
		ID:        AgentIDA02,
		Name:      AgentNameA02,
		Status:    shared.AgentStatusHealthy,
		LastCheck: time.Now(),
	}, nil
}

// Correct attempts to correct a failed SQL query.
func (a *CorrectionAgent) Correct(ctx context.Context, sql string, errMsg string, attempt int) (*shared.GeneratedSQL, error) {
	a.logger.Debug("attempting SQL correction",
		slog.String("sql", sql),
		slog.String("error", errMsg),
		slog.Int("attempt", attempt),
	)

	if attempt >= MaxRetries {
		return nil, fmt.Errorf("max retries (%d) exceeded", MaxRetries)
	}

	// Analyze the error and attempt correction
	correctedSQL, confidence := a.analyzeAndCorrect(sql, errMsg)

	if correctedSQL == "" {
		return nil, fmt.Errorf("unable to correct SQL: %s", errMsg)
	}

	a.logger.Info("SQL corrected",
		slog.String("original", sql),
		slog.String("corrected", correctedSQL),
		slog.Int("attempt", attempt),
	)

	return &shared.GeneratedSQL{
		SQL:             correctedSQL,
		IsParameterized: false,
		Confidence:      confidence,
		Explanation:     fmt.Sprintf("Corrected after error: %s", errMsg),
	}, nil
}

// CanRetry determines if the error is retryable.
func (a *CorrectionAgent) CanRetry(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Non-retryable errors
	nonRetryablePatterns := []string{
		"permission denied",
		"does not exist",
		"syntax error",
		"invalid",
		"unauthorized",
	}

	for _, pattern := range nonRetryablePatterns {
		if strings.Contains(strings.ToLower(errMsg), pattern) {
			return false
		}
	}

	// Retryable errors
	retryablePatterns := []string{
		"timeout",
		"connection",
		"temporary",
		"deadlock",
		"column ambiguously defined",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errMsg), pattern) {
			return true
		}
	}

	return true // Default to retryable
}

// analyzeAndCorrect analyzes the error and attempts correction.
func (a *CorrectionAgent) analyzeAndCorrect(sql, errMsg string) (string, float64) {
	confidence := 80.0
	errLower := strings.ToLower(errMsg)

	// Common error patterns and corrections
	switch {
	// Column doesn't exist - try to find similar column
	case strings.Contains(errLower, "column") && strings.Contains(errLower, "does not exist"):
		return a.fixColumnError(sql, errMsg), 70.0

	// Table doesn't exist - try to qualify with schema
	case strings.Contains(errLower, "relation") && strings.Contains(errLower, "does not exist"):
		return a.fixTableError(sql, errMsg), 75.0

	// Ambiguous column - add table qualification
	case strings.Contains(errLower, "column reference") && strings.Contains(errLower, "ambiguous"):
		return a.fixAmbiguousColumn(sql, errMsg), 72.0

	// Syntax error near keyword
	case strings.Contains(errLower, "syntax error"):
		return a.fixSyntaxError(sql, errMsg), 65.0

	// Group by error
	case strings.Contains(errLower, "must appear in the group by clause"):
		return a.fixGroupByError(sql, errMsg), 70.0

	default:
		// Return original for unknown errors
		return "", confidence
	}
}

// fixColumnError attempts to fix column name errors.
func (a *CorrectionAgent) fixColumnError(sql, errMsg string) string {
	// Extract the column name from error message
	re := regexp.MustCompile(`column "([^"]+)" does not exist`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return sql
	}

	wrongColumn := matches[1]
	a.logger.Debug("attempting column fix",
		slog.String("wrong_column", wrongColumn),
	)

	// Common column name mappings
	columnMappings := map[string]string{
		"name":       "name_en",
		"patient":    "patient_id",
		"doctor":     "doctor_id",
		"date":       "created_at",
		"amount":     "total_amount",
		"created":    "created_at",
		"updated":    "updated_at",
	}

	if correct, ok := columnMappings[strings.ToLower(wrongColumn)]; ok {
		return strings.ReplaceAll(sql, wrongColumn, correct)
	}

	return sql
}

// fixTableError attempts to fix table name errors by adding schema qualification.
func (a *CorrectionAgent) fixTableError(sql, errMsg string) string {
	// Extract table name from error
	re := regexp.MustCompile(`relation "([^"]+)" does not exist`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return sql
	}

	table := matches[1]

	// Try to qualify with common schemas
	schemas := []string{"hims_analytics", "tally_analytics", "app", "public"}
	for _, schema := range schemas {
		qualified := fmt.Sprintf("%s.%s", schema, table)
		newSQL := strings.ReplaceAll(sql, " "+table+" ", " "+qualified+" ")
		newSQL = strings.ReplaceAll(newSQL, " "+table, " "+qualified)
		if newSQL != sql {
			return newSQL
		}
	}

	return sql
}

// fixAmbiguousColumn adds table qualification to ambiguous columns.
func (a *CorrectionAgent) fixAmbiguousColumn(sql, errMsg string) string {
	// This would need more sophisticated parsing in production
	return sql // Placeholder
}

// fixSyntaxError attempts to fix common syntax errors.
func (a *CorrectionAgent) fixSyntaxError(sql, errMsg string) string {
	// Common syntax fixes
	fixed := sql

	// Remove trailing semicolons
	fixed = strings.TrimSuffix(strings.TrimSpace(fixed), ";")

	// Fix double spaces
	re := regexp.MustCompile(`\s+`)
	fixed = re.ReplaceAllString(fixed, " ")

	return fixed
}

// fixGroupByError attempts to fix GROUP BY errors.
func (a *CorrectionAgent) fixGroupByError(sql, errMsg string) string {
	// Extract column from error message
	re := regexp.MustCompile(`column "([^"]+)" must appear`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return sql
	}

	column := matches[1]

	// Add to GROUP BY clause
	if strings.Contains(strings.ToUpper(sql), "GROUP BY") {
		// Add to existing GROUP BY
		re = regexp.MustCompile(`(?i)GROUP BY (.+?)(?:ORDER|LIMIT|$)`)
		sql = re.ReplaceAllStringFunc(sql, func(m string) string {
			return m + ", " + column
		})
	} else {
		// Add GROUP BY clause
		sql = sql + " GROUP BY " + column
	}

	return sql
}
