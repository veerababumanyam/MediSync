// Package a01_text_to_sql provides the text-to-SQL agent subcomponents.
//
// This file implements OPA SQL validation for security enforcement.
package a01_text_to_sql

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

// SQLValidator validates SQL queries against security policies.
type SQLValidator struct {
	opaClient   OPAClient
	logger      *slog.Logger
	cache       map[string]ValidationResult
	cacheMu     sync.RWMutex
	strictMode  bool
}

// OPAClient defines the interface for OPA policy evaluation.
type OPAClient interface {
	// Evaluate evaluates a policy with the given input.
	Evaluate(ctx context.Context, policy string, input map[string]interface{}) (map[string]interface{}, error)
}

// ValidationResult contains the result of SQL validation.
type ValidationResult struct {
	IsValid       bool     `json:"is_valid"`
	BlockedReason string   `json:"blocked_reason,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
	AllowedOps    []string `json:"allowed_ops,omitempty"`
	PolicyName    string   `json:"policy_name"`
	Score         float64  `json:"score"`
}

// SQLValidationInput is the input for OPA policy evaluation.
type SQLValidationInput struct {
	SQL         string            `json:"sql"`
	UserRoles   []string          `json:"user_roles"`
	TenantID    string            `json:"tenant_id"`
	QueryType   string            `json:"query_type"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SQLValidationConfig holds configuration for the validator.
type SQLValidationConfig struct {
	OPAClient  OPAClient
	Logger     *slog.Logger
	StrictMode bool
	CacheSize  int
}

// NewSQLValidator creates a new SQL validator.
func NewSQLValidator(cfg SQLValidationConfig) *SQLValidator {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.CacheSize == 0 {
		cfg.CacheSize = 1000
	}

	return &SQLValidator{
		opaClient:  cfg.OPAClient,
		logger:     cfg.Logger.With("component", "sql_validator"),
		cache:      make(map[string]ValidationResult, cfg.CacheSize),
		strictMode: cfg.StrictMode,
	}
}

// Validate validates a SQL query against OPA policies.
func (v *SQLValidator) Validate(ctx context.Context, sql string, input SQLValidationInput) (*ValidationResult, error) {
	v.logger.Debug("validating SQL", "sql_length", len(sql))

	// Check cache first
	cacheKey := v.getCacheKey(sql, input)
	if result, ok := v.getFromCache(cacheKey); ok {
		v.logger.Debug("validation cache hit")
		return &result, nil
	}

	// Perform local validation first (fast path)
	if err := v.localValidation(sql); err != nil {
		result := &ValidationResult{
			IsValid:       false,
			BlockedReason: err.Error(),
			PolicyName:    "local_validation",
			Score:         0,
		}
		v.addToCache(cacheKey, *result)
		return result, nil
	}

	// Prepare input for OPA
	opaInput := map[string]interface{}{
		"sql":        sql,
		"user_roles": input.UserRoles,
		"tenant_id":  input.TenantID,
		"query_type": input.QueryType,
		"metadata":   input.Metadata,
	}

	// Evaluate OPA policy
	opaResult, err := v.opaClient.Evaluate(ctx, "bi.read_only", opaInput)
	if err != nil {
		v.logger.Warn("OPA evaluation failed", "error", err)
		// In strict mode, fail closed
		if v.strictMode {
			return nil, fmt.Errorf("OPA evaluation failed: %w", err)
		}
		// In non-strict mode, fall back to local validation
		return v.fallbackValidation(sql)
	}

	// Parse OPA result
	result := v.parseOPAResult(opaResult)

	// Add to cache
	v.addToCache(cacheKey, *result)

	v.logger.Info("SQL validation complete",
		"is_valid", result.IsValid,
		"score", result.Score,
		"policy", result.PolicyName)

	return result, nil
}

// localValidation performs fast local validation checks.
func (v *SQLValidator) localValidation(sql string) error {
	sqlUpper := strings.ToUpper(sql)

	// Check for forbidden operations
	forbiddenOps := []struct {
		pattern string
		message string
	}{
		{`\bINSERT\b`, "INSERT operations are not allowed"},
		{`\bUPDATE\b`, "UPDATE operations are not allowed"},
		{`\bDELETE\b`, "DELETE operations are not allowed"},
		{`\bDROP\b`, "DROP operations are not allowed"},
		{`\bTRUNCATE\b`, "TRUNCATE operations are not allowed"},
		{`\bALTER\b`, "ALTER operations are not allowed"},
		{`\bCREATE\b`, "CREATE operations are not allowed"},
		{`\bGRANT\b`, "GRANT operations are not allowed"},
		{`\bREVOKE\b`, "REVOKE operations are not allowed"},
		{`\bEXEC\b`, "EXEC operations are not allowed"},
		{`\bEXECUTE\b`, "EXECUTE operations are not allowed"},
		{`\bCALL\b`, "CALL operations are not allowed"},
	}

	for _, op := range forbiddenOps {
		matched, _ := regexp.MatchString(op.pattern, sqlUpper)
		if matched {
			return fmt.Errorf("%s", op.message)
		}
	}

	// Check for multiple statements
	if strings.Count(sql, ";") > 1 {
		return fmt.Errorf("multiple SQL statements are not allowed")
	}

	// Check for suspicious patterns
	suspiciousPatterns := []struct {
		pattern string
		message string
	}{
		{`--`, "SQL comments are not allowed"},
		{`/\*`, "SQL comment blocks are not allowed"},
		{`\bUNION\b.*\bSELECT\b`, "UNION SELECT is not allowed"},
		{`;.*\bDROP\b`, "Statement chaining detected"},
	}

	for _, sp := range suspiciousPatterns {
		matched, _ := regexp.MatchString(sp.pattern, sqlUpper)
		if matched {
			return fmt.Errorf("%s", sp.message)
		}
	}

	// Verify SELECT is the first keyword (after trimming)
	trimmed := strings.TrimSpace(sqlUpper)
	if !strings.HasPrefix(trimmed, "SELECT") {
		return fmt.Errorf("only SELECT queries are allowed")
	}

	return nil
}

// fallbackValidation performs validation when OPA is unavailable.
func (v *SQLValidator) fallbackValidation(sql string) (*ValidationResult, error) {
	v.logger.Warn("using fallback validation (OPA unavailable)")

	result := &ValidationResult{
		IsValid:    true,
		PolicyName: "fallback",
		Score:      0.7, // Lower confidence for fallback
		Warnings:   []string{"OPA policy evaluation unavailable, using fallback validation"},
	}

	if err := v.localValidation(sql); err != nil {
		result.IsValid = false
		result.BlockedReason = err.Error()
		result.Score = 0
	}

	return result, nil
}

// parseOPAResult extracts validation result from OPA response.
func (v *SQLValidator) parseOPAResult(opaResult map[string]interface{}) *ValidationResult {
	result := &ValidationResult{
		IsValid:    false,
		PolicyName: "bi.read_only",
		Warnings:   []string{},
	}

	// Extract allow decision
	if allow, ok := opaResult["allow"].(bool); ok {
		result.IsValid = allow
	}

	// Extract score
	if score, ok := opaResult["score"].(float64); ok {
		result.Score = score
	}

	// Extract blocked reason
	if reason, ok := opaResult["blocked_reason"].(string); ok {
		result.BlockedReason = reason
	}

	// Extract warnings
	if warnings, ok := opaResult["warnings"].([]interface{}); ok {
		for _, w := range warnings {
			if ws, ok := w.(string); ok {
				result.Warnings = append(result.Warnings, ws)
			}
		}
	}

	// Extract allowed operations
	if ops, ok := opaResult["allowed_ops"].([]interface{}); ok {
		for _, op := range ops {
			if ops, ok := op.(string); ok {
				result.AllowedOps = append(result.AllowedOps, ops)
			}
		}
	}

	return result
}

// getCacheKey generates a cache key for a validation request.
func (v *SQLValidator) getCacheKey(sql string, input SQLValidationInput) string {
	// Simple hash based on SQL and roles
	roles := strings.Join(input.UserRoles, ",")
	return fmt.Sprintf("%s:%s:%s", input.TenantID, roles, sql)
}

// getFromCache retrieves a cached validation result.
func (v *SQLValidator) getFromCache(key string) (ValidationResult, bool) {
	v.cacheMu.RLock()
	defer v.cacheMu.RUnlock()
	result, ok := v.cache[key]
	return result, ok
}

// addToCache adds a validation result to the cache.
func (v *SQLValidator) addToCache(key string, result ValidationResult) {
	v.cacheMu.Lock()
	defer v.cacheMu.Unlock()

	// Simple cache eviction when full
	if len(v.cache) >= 1000 {
		// Remove a random entry (in production, use LRU)
		for k := range v.cache {
			delete(v.cache, k)
			break
		}
	}

	v.cache[key] = result
}

// ClearCache clears the validation cache.
func (v *SQLValidator) ClearCache() {
	v.cacheMu.Lock()
	defer v.cacheMu.Unlock()
	v.cache = make(map[string]ValidationResult)
}

// ValidateBatch validates multiple SQL statements.
func (v *SQLValidator) ValidateBatch(ctx context.Context, statements []string, input SQLValidationInput) ([]ValidationResult, error) {
	results := make([]ValidationResult, len(statements))

	for i, sql := range statements {
		result, err := v.Validate(ctx, sql, input)
		if err != nil {
			return nil, fmt.Errorf("validation failed for statement %d: %w", i, err)
		}
		results[i] = *result
	}

	return results, nil
}

// GetValidationMetrics returns metrics about validation performance.
func (v *SQLValidator) GetValidationMetrics() map[string]interface{} {
	v.cacheMu.RLock()
	defer v.cacheMu.RUnlock()

	return map[string]interface{}{
		"cache_size":    len(v.cache),
		"strict_mode":   v.strictMode,
	}
}

// ToJSON serializes the validation result.
func (r *ValidationResult) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// FromJSON deserializes a validation result.
func (r *ValidationResult) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), r)
}
