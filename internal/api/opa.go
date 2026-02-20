// Package api provides OPA (Open Policy Agent) client for MediSync.
//
// This file implements the OPAClient which handles policy evaluation
// for authorization decisions throughout the application.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// OPAClient handles communication with the Open Policy Agent.
type OPAClient struct {
	url        string
	logger     *slog.Logger
	httpClient *http.Client
}

// OPAConfig holds configuration for the OPA client.
type OPAConfig struct {
	// URL is the OPA server URL.
	URL string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration
}

// OPADecision represents an OPA policy decision.
type OPADecision struct {
	// Result is the policy evaluation result.
	Result interface{} `json:"result"`

	// DecisionID is the unique identifier for this decision.
	DecisionID string `json:"decision_id"`
}

// OPAInput represents the input to an OPA policy evaluation.
type OPAInput struct {
	// User contains user information.
	User OPAUserInput `json:"user"`

	// Action is the action being performed.
	Action string `json:"action"`

	// Resource is the resource being accessed.
	Resource string `json:"resource,omitempty"`

	// Context contains additional context for the policy.
	Context map[string]interface{} `json:"context,omitempty"`
}

// OPAUserInput contains user information for policy evaluation.
type OPAUserInput struct {
	// ID is the user's unique identifier.
	ID string `json:"id"`

	// Roles are the user's assigned roles.
	Roles []string `json:"roles"`

	// TenantID is the user's tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`

	// Locale is the user's preferred locale.
	Locale string `json:"locale,omitempty"`

	// ServiceAccount indicates if this is a service account.
	ServiceAccount string `json:"service_account,omitempty"`

	// Authenticated indicates if the user is authenticated.
	Authenticated bool `json:"authenticated"`
}

// BIReadOnlyInput represents input for the bi_read_only policy.
type BIReadOnlyInput struct {
	// Query is the SQL query to validate.
	Query string `json:"query"`

	// User contains user information.
	User OPAUserInput `json:"user"`
}

// BIReadOnlyResult represents the result of bi_read_only policy evaluation.
type BIReadOnlyResult struct {
	// Allow indicates if the query is allowed.
	Allow bool `json:"allow"`

	// Reason explains why the query was allowed or denied.
	Reason string `json:"reason"`

	// Violations contains any policy violations.
	Violations []string `json:"violations,omitempty"`
}

// ColumnMaskingInput represents input for the column_masking policy.
type ColumnMaskingInput struct {
	// Columns are the columns being accessed.
	Columns []string `json:"columns"`

	// User contains user information.
	User OPAUserInput `json:"user"`

	// Schema is the database schema.
	Schema string `json:"schema"`

	// Table is the database table.
	Table string `json:"table"`
}

// ColumnMaskingResult represents the result of column_masking policy evaluation.
type ColumnMaskingResult struct {
	// Masks maps column names to their mask types.
	Masks map[string]string `json:"masks"`

	// Reason explains the masking decisions.
	Reason string `json:"reason"`
}

// Mask types
const (
	MaskNone     = "none"
	MaskPartial  = "partial"
	MaskFull     = "full"
	MaskHash     = "hash"
	MaskRedacted = "redacted"
)

// NewOPAClient creates a new OPA client.
func NewOPAClient(cfg *OPAConfig, logger *slog.Logger) (*OPAClient, error) {
	if logger == nil {
		logger = slog.Default()
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	client := &OPAClient{
		url:    cfg.URL,
		logger: logger,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}

	logger.Info("OPA client initialized",
		slog.String("url", cfg.URL),
		slog.Duration("timeout", cfg.Timeout),
	)

	return client, nil
}

// Evaluate evaluates an OPA policy and returns the decision.
func (c *OPAClient) Evaluate(ctx context.Context, path string, input interface{}) (*OPADecision, error) {
	url := fmt.Sprintf("%s/v1/data/%s", c.url, path)

	// Prepare request body
	reqBody := map[string]interface{}{
		"input": input,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OPA: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OPA request failed with status %d", resp.StatusCode)
	}

	var decision OPADecision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to decode OPA response: %w", err)
	}

	return &decision, nil
}

// Allow evaluates a policy and returns whether the action is allowed.
func (c *OPAClient) Allow(ctx context.Context, path string, input interface{}) (bool, error) {
	decision, err := c.Evaluate(ctx, path, input)
	if err != nil {
		return false, err
	}

	// Check if the result contains an "allow" field
	if result, ok := decision.Result.(map[string]interface{}); ok {
		if allow, ok := result["allow"].(bool); ok {
			return allow, nil
		}
	}

	return false, nil
}

// ValidateBIQuery validates that a SQL query is read-only.
func (c *OPAClient) ValidateBIQuery(ctx context.Context, query string, user OPAUserInput) (*BIReadOnlyResult, error) {
	input := BIReadOnlyInput{
		Query: query,
		User:  user,
	}

	decision, err := c.Evaluate(ctx, "medisync.bi_read_only", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate bi_read_only policy: %w", err)
	}

	result := &BIReadOnlyResult{
		Allow: false,
	}

	// Parse the decision result
	if res, ok := decision.Result.(map[string]interface{}); ok {
		if allow, ok := res["allow"].(bool); ok {
			result.Allow = allow
		}
		if reason, ok := res["reason"].(string); ok {
			result.Reason = reason
		}
		if violations, ok := res["violations"].([]interface{}); ok {
			for _, v := range violations {
				if s, ok := v.(string); ok {
					result.Violations = append(result.Violations, s)
				}
			}
		}
	}

	return result, nil
}

// GetColumnMasks returns the masking requirements for columns.
func (c *OPAClient) GetColumnMasks(ctx context.Context, columns []string, schema, table string, user OPAUserInput) (*ColumnMaskingResult, error) {
	input := ColumnMaskingInput{
		Columns: columns,
		User:    user,
		Schema:  schema,
		Table:   table,
	}

	decision, err := c.Evaluate(ctx, "medisync.column_masking", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate column_masking policy: %w", err)
	}

	result := &ColumnMaskingResult{
		Masks:  make(map[string]string),
		Reason: "No masking required",
	}

	// Parse the decision result
	if res, ok := decision.Result.(map[string]interface{}); ok {
		if masks, ok := res["masks"].(map[string]interface{}); ok {
			for col, mask := range masks {
				if maskStr, ok := mask.(string); ok {
					result.Masks[col] = maskStr
				}
			}
		}
		if reason, ok := res["reason"].(string); ok {
			result.Reason = reason
		}
	}

	return result, nil
}

// AllowTallySync checks if the user is allowed to sync to Tally.
func (c *OPAClient) AllowTallySync(ctx context.Context, user OPAUserInput, costCentre string) (bool, error) {
	input := OPAInput{
		User:   user,
		Action: "tally_sync",
		Context: map[string]interface{}{
			"cost_centre": costCentre,
		},
	}

	return c.Allow(ctx, "medisync.tally_sync", input)
}

// AllowReportAccess checks if the user can access a report.
func (c *OPAClient) AllowReportAccess(ctx context.Context, user OPAUserInput, reportID string) (bool, error) {
	input := OPAInput{
		User:     user,
		Action:   "report_read",
		Resource: reportID,
	}

	return c.Allow(ctx, "medisync.reports", input)
}

// Health checks if the OPA server is healthy.
func (c *OPAClient) Health(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("OPA health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OPA health check returned status %d", resp.StatusCode)
	}

	return nil
}
