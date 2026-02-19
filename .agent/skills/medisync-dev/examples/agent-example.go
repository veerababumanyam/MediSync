// Example: Simple MediSync Agent Implementation
// This demonstrates the basic pattern for creating an AI agent in MediSync

package module_a

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// SimpleQueryAgent demonstrates a basic Genkit flow for query processing
// This would be agent A-01 (Text-to-SQL) in the full implementation

type SimpleQueryRequest struct {
	Query     string `json:"query"`
	Locale    string `json:"locale"`
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
}

type SimpleQueryResponse struct {
	SQL        string  `json:"sql"`
	Confidence float64 `json:"confidence"`
	Explanation string `json:"explanation"`
}

type SimpleQueryAgent struct {
	model     ai.Model
	warehouse WarehouseClient
	opa       OPAClient
}

func NewSimpleQueryAgent(model ai.Model, db WarehouseClient, opa OPAClient) *SimpleQueryAgent {
	return &SimpleQueryAgent{
		model:     model,
		warehouse: db,
		opa:       opa,
	}
}

// ProcessFlow is the main Genkit flow for this agent
func (a *SimpleQueryAgent) ProcessFlow(ctx context.Context, req SimpleQueryRequest) (*SimpleQueryResponse, error) {
	// Step 1: Validate input
	if err := a.validateInput(req); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Step 2: Check authorization with OPA
	if err := a.checkAuthorization(ctx, req); err != nil {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	// Step 3: Build prompt with locale instruction
	prompt := a.buildPrompt(req)

	// Step 4: Generate SQL using LLM
	resp, err := genkit.Generate(ctx, a.model, ai.WithPrompt(prompt))
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Step 5: Extract SQL from response
	sql := a.extractSQL(resp.Text)
	if sql == "" {
		return nil, fmt.Errorf("no SQL generated")
	}

	// Step 6: Validate SQL is SELECT-only (critical security check)
	if !a.isSelectOnly(sql) {
		return nil, fmt.Errorf("generated SQL is not SELECT-only")
	}

	// Step 7: Execute query via readonly role
	result, err := a.warehouse.Query(ctx, sql, "medisync_readonly")
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	// Step 8: Calculate confidence score
	confidence := a.calculateConfidence(req.Query, sql, result)

	// Step 9: Generate explanation
	explanation := a.explainSQL(sql, result, req.Locale)

	return &SimpleQueryResponse{
		SQL:        sql,
		Confidence: confidence,
		Explanation: explanation,
	}, nil
}

func (a *SimpleQueryAgent) validateInput(req SimpleQueryRequest) error {
	if strings.TrimSpace(req.Query) == "" {
		return fmt.Errorf("query cannot be empty")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

func (a *SimpleQueryAgent) checkAuthorization(ctx context.Context, req SimpleQueryRequest) error {
	allowed, err := a.opa.Allow(ctx, "warehouse_query", map[string]interface{}{
		"user_id":    req.UserID,
		"company_id": req.CompanyID,
		"action":     "select",
	})
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("user not authorized for this query")
	}
	return nil
}

func (a *SimpleQueryAgent) buildPrompt(req SimpleQueryRequest) string {
	// Schema context would be included here
	schemaContext := `
Available tables:
- patients (id, name, date_of_birth, city, phone)
- doctors (id, name, specialization)
- appointments (id, patient_id, doctor_id, date, status)
- invoices (id, patient_id, amount, date, status)
`

	// Locale instruction for i18n
	localeInstr := "\nResponseLanguageInstruction: "
	if req.Locale == "ar" {
		localeInstr += "Respond in Arabic. Use Arabic numerals (٠-٩)."
	} else {
		localeInstr += "Respond in English."
	}

	return fmt.Sprintf(`You are a SQL generator for a healthcare database.

%s

Generate a SELECT query to answer this question:
%s

Return ONLY the SQL query, no explanation.
%s`, schemaContext, req.Query, localeInstr)
}

func (a *SimpleQueryAgent) extractSQL(response string) string {
	// Extract SQL from markdown code blocks if present
	response = strings.TrimSpace(response)

	// Try to extract from ```sql ... ``` blocks
	if strings.Contains(response, "```") {
		parts := strings.Split(response, "```")
		for i, part := range parts {
			if i > 0 && strings.HasPrefix(strings.ToLower(part), "sql") {
				return strings.TrimSpace(strings.TrimPrefix(strings.ToLower(part), "sql"))
			}
		}
	}

	return response
}

func (a *SimpleQueryAgent) isSelectOnly(sql string) bool {
	trimmed := strings.ToUpper(strings.TrimSpace(sql))
	return strings.HasPrefix(trimmed, "SELECT") &&
		!strings.Contains(trimmed, "INSERT") &&
		!strings.Contains(trimmed, "UPDATE") &&
		!strings.Contains(trimmed, "DELETE") &&
		!strings.Contains(trimmed, "DROP") &&
		!strings.Contains(trimmed, "CREATE") &&
		!strings.Contains(trimmed, "ALTER")
}

func (a *SimpleQueryAgent) calculateConfidence(query, sql string, result any) float64 {
	confidence := 1.0

	// Deduct for empty results
	if result == nil {
		confidence -= 0.3
	}

	// Deduct for SQL with uncertain patterns
	if strings.Contains(strings.ToUpper(sql), "LIKE") {
		confidence -= 0.1
	}

	return confidence
}

func (a *SimpleQueryAgent) explainSQL(sql string, result any, locale string) string {
	if locale == "ar" {
		return fmt.Sprintf("تم تنفيذ الاستعلام: %s", sql)
	}
	return fmt.Sprintf("Executed query: %s", sql)
}

// Interfaces for dependencies
type WarehouseClient interface {
	Query(ctx context.Context, sql, role string) (any, error)
}

type OPAClient interface {
	Allow(ctx context.Context, action string, input map[string]interface{}) (bool, error)
}
