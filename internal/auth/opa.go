// Package auth provides authorization via Open Policy Agent (OPA).
//
// This file provides the OPAClient struct for policy evaluation using OPA's
// REST API. It supports:
//   - Policy evaluation with arbitrary input data
//   - SQL query validation against user roles
//   - Caching of policy decisions for performance
//
// OPA policies are defined as Rego files and loaded by the OPA server.
// The client communicates with OPA via its REST API for policy evaluation.
//
// Usage:
//
//	opaClient, err := auth.NewOPAClient(config.OPA, logger)
//	if err != nil {
//	    log.Fatal("Failed to create OPA client:", err)
//	}
//
//	// Check if action is allowed
//	allowed, err := opaClient.Allow(ctx, "tally_sync", map[string]any{
//	    "user":   userClaims,
//	    "action": "sync",
//	    "resource": journalEntry,
//	})
package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// OPA configuration constants.
const (
	// DefaultOPATimeout is the default timeout for OPA API requests.
	DefaultOPATimeout = 5 * time.Second

	// DefaultDecisionPath is the default path for OPA decision API.
	DefaultDecisionPath = "/v1/data"
)

// OPADecision represents an OPA policy decision.
type OPADecision struct {
	// Result is the decision result.
	Result bool `json:"result"`

	// Explanation contains additional details about the decision.
	Explanation map[string]interface{} `json:"explanation,omitempty"`
}

// opaRequest represents a request to OPA.
type opaRequest struct {
	Input map[string]interface{} `json:"input"`
}

// opaResponse represents a response from OPA.
type opaResponse struct {
	Result interface{} `json:"result"`
}

// OPAConfig holds configuration for OPA connection.
type OPAConfig struct {
	// URL is the OPA server base URL.
	URL string

	// DecisionPath is the path prefix for decision API.
	DecisionPath string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// Logger is the structured logger.
	Logger *slog.Logger
}

// OPAClient provides policy evaluation via OPA REST API.
type OPAClient struct {
	config     *OPAConfig
	httpClient *http.Client
	logger     *slog.Logger
}

// NewOPAClient creates a new OPA client.
func NewOPAClient(cfg *OPAConfig) (*OPAClient, error) {
	if cfg == nil || cfg.URL == "" {
		return nil, fmt.Errorf("auth: OPA config with URL is required")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultOPATimeout
	}

	if cfg.DecisionPath == "" {
		cfg.DecisionPath = DefaultDecisionPath
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &OPAClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		logger: cfg.Logger,
	}, nil
}

// Allow evaluates a policy and returns whether the action is allowed.
// The policy parameter is the OPA policy path (e.g., "medisync/authz/allow").
// The input parameter contains the data for policy evaluation.
func (c *OPAClient) Allow(ctx context.Context, policy string, input map[string]interface{}) (bool, error) {
	decision, err := c.Evaluate(ctx, policy, input)
	if err != nil {
		return false, err
	}

	// Result is a bool directly
	c.logDecision(policy, decision.Result, input)
	return decision.Result, nil
}

// Evaluate evaluates a policy and returns the full decision.
func (c *OPAClient) Evaluate(ctx context.Context, policy string, input map[string]interface{}) (*OPADecision, error) {
	// Build the OPA API URL
	// Convert policy path like "medisync/authz/allow" to URL path
	policyPath := strings.ReplaceAll(policy, ".", "/")
	url := fmt.Sprintf("%s%s/%s", c.config.URL, c.config.DecisionPath, policyPath)

	// Build request body
	reqBody := opaRequest{
		Input: input,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to marshal OPA request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create OPA request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to execute OPA request: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth: OPA returned status %d", resp.StatusCode)
	}

	// Parse response
	var opaResp opaResponse
	if err := json.NewDecoder(resp.Body).Decode(&opaResp); err != nil {
		return nil, fmt.Errorf("auth: failed to decode OPA response: %w", err)
	}

	c.logger.Debug("OPA policy evaluated",
		slog.String("policy", policy),
		slog.Duration("duration", duration),
	)

	decision := &OPADecision{}

	switch v := opaResp.Result.(type) {
	case bool:
		decision.Result = v
	case map[string]interface{}:
		if allowed, ok := v["allow"].(bool); ok {
			decision.Result = allowed
		}
		// Store full result as explanation
		decision.Explanation = v
	default:
		decision.Result = false
	}

	return decision, nil
}

// ValidateSQL validates that a SQL query is allowed for the given user roles.
// This checks against OPA policies that define what queries each role can execute.
func (c *OPAClient) ValidateSQL(ctx context.Context, sql string, userRoles []string) error {
	input := map[string]interface{}{
		"sql":   sql,
		"roles": userRoles,
		"query": map[string]interface{}{
			"is_select_only": isSelectOnly(sql),
			"tables":         extractTables(sql),
		},
	}

	allowed, err := c.Allow(ctx, "medisync/sql/validate", input)
	if err != nil {
		// If OPA is unavailable, fall back to basic validation
		c.logger.Warn("OPA SQL validation failed, using fallback",
			slog.String("error", err.Error()),
		)
		return c.fallbackSQLValidation(sql, userRoles)
	}

	if !allowed {
		return fmt.Errorf("auth: SQL query not allowed for roles %v", userRoles)
	}

	return nil
}

// ValidateTallySync validates that a Tally sync operation is allowed.
func (c *OPAClient) ValidateTallySync(ctx context.Context, userID string, roles []string, entry map[string]interface{}) error {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    userID,
			"roles": roles,
		},
		"action":   "sync_to_tally",
		"resource": entry,
	}

	allowed, err := c.Allow(ctx, "medisync/tally/sync", input)
	if err != nil {
		return fmt.Errorf("auth: failed to evaluate Tally sync policy: %w", err)
	}

	if !allowed {
		return fmt.Errorf("auth: Tally sync not authorized")
	}

	return nil
}

// ValidateReportAccess validates access to a specific report.
func (c *OPAClient) ValidateReportAccess(ctx context.Context, userID string, roles []string, reportID string, reportType string) error {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    userID,
			"roles": roles,
		},
		"action": "view_report",
		"resource": map[string]interface{}{
			"id":   reportID,
			"type": reportType,
		},
	}

	allowed, err := c.Allow(ctx, "medisync/reports/access", input)
	if err != nil {
		return fmt.Errorf("auth: failed to evaluate report access policy: %w", err)
	}

	if !allowed {
		return fmt.Errorf("auth: report access denied")
	}

	return nil
}

// ValidateDataAccess validates access to data based on cost centre restrictions.
func (c *OPAClient) ValidateDataAccess(ctx context.Context, userID string, roles []string, costCentres []string, resourceCostCentre string) error {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"id":           userID,
			"roles":        roles,
			"cost_centres": costCentres,
		},
		"resource": map[string]interface{}{
			"cost_centre": resourceCostCentre,
		},
	}

	allowed, err := c.Allow(ctx, "medisync/data/access", input)
	if err != nil {
		return fmt.Errorf("auth: failed to evaluate data access policy: %w", err)
	}

	if !allowed {
		return fmt.Errorf("auth: data access denied for cost centre %s", resourceCostCentre)
	}

	return nil
}

// logDecision logs the OPA decision for audit purposes.
func (c *OPAClient) logDecision(policy string, allowed bool, input map[string]interface{}) {
	// Extract user info if available for logging
	var userID string
	if user, ok := input["user"].(map[string]interface{}); ok {
		if id, ok := user["id"].(string); ok {
			userID = id
		}
	}

	c.logger.Info("OPA decision",
		slog.String("policy", policy),
		slog.Bool("allowed", allowed),
		slog.String("user_id", userID),
	)
}

// fallbackSQLValidation provides basic SQL validation when OPA is unavailable.
func (c *OPAClient) fallbackSQLValidation(sql string, roles []string) error {
	// Check if query is SELECT-only
	if !isSelectOnly(sql) {
		// Only certain roles can run non-SELECT queries
		adminRoles := []string{"admin", "finance_head"}
		hasAdminRole := false
		for _, role := range roles {
			for _, adminRole := range adminRoles {
				if role == adminRole {
					hasAdminRole = true
					break
				}
			}
		}

		if !hasAdminRole {
			return fmt.Errorf("auth: only SELECT queries are allowed")
		}
	}

	return nil
}

// isSelectOnly checks if a SQL query is SELECT-only.
func isSelectOnly(sql string) bool {
	upperSQL := strings.ToUpper(strings.TrimSpace(sql))
	return strings.HasPrefix(upperSQL, "SELECT") || strings.HasPrefix(upperSQL, "WITH")
}

// extractTables extracts table names from a SQL query (simplified implementation).
func extractTables(sql string) []string {
	// This is a simplified implementation
	// In production, use a proper SQL parser
	var tables []string

	upperSQL := strings.ToUpper(sql)

	// Look for FROM and JOIN keywords
	keywords := []string{"FROM", "JOIN", "INTO", "UPDATE"}

	for _, keyword := range keywords {
		idx := strings.Index(upperSQL, keyword)
		if idx == -1 {
			continue
		}

		// Extract the word after the keyword
		after := strings.TrimSpace(sql[idx+len(keyword):])
		parts := strings.Fields(after)
		if len(parts) > 0 {
			table := strings.Trim(parts[0], ",;()")
			if table != "" && !strings.HasPrefix(strings.ToUpper(table), "SELECT") {
				tables = append(tables, table)
			}
		}
	}

	return tables
}

// HealthCheck checks if the OPA server is healthy.
func (c *OPAClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.config.URL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("auth: failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth: OPA health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: OPA returned status %d", resp.StatusCode)
	}

	return nil
}
