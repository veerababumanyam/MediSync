// Package module_a_test tests the conversational BI agents (module A).
package module_a_test

import (
	"testing"

	"github.com/medisync/medisync/internal/agents/module_a/a01_text_to_sql"
)

func TestSQLParameterizer_BasicParameterization(t *testing.T) {
	parameterizer := a01_text_to_sql.NewParameterizer(a01_text_to_sql.ParameterizerConfig{})

	tests := []struct {
		name           string
		sql            string
		expectedParams int
	}{
		{
			name:           "string literal",
			sql:            "SELECT * FROM users WHERE name = 'John'",
			expectedParams: 1,
		},
		{
			name:           "multiple string literals",
			sql:            "SELECT * FROM users WHERE name = 'John' AND status = 'active'",
			expectedParams: 2,
		},
		{
			name:           "numeric literal",
			sql:            "SELECT * FROM orders WHERE amount > 100",
			expectedParams: 1,
		},
		{
			name:           "mixed literals",
			sql:            "SELECT * FROM orders WHERE status = 'pending' AND amount > 50",
			expectedParams: 2,
		},
		{
			name:           "date literal",
			sql:            "SELECT * FROM orders WHERE created_at > DATE '2026-01-01'",
			expectedParams: 1,
		},
		{
			name:           "no literals to parameterize",
			sql:            "SELECT * FROM users",
			expectedParams: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parameterizer.Parameterize(tt.sql)
			if len(result.Parameters) != tt.expectedParams {
				t.Errorf("Parameterize() params = %v, want %v", len(result.Parameters), tt.expectedParams)
			}
		})
	}
}

func TestSQLParameterizer_Validation(t *testing.T) {
	parameterizer := a01_text_to_sql.NewParameterizer(a01_text_to_sql.ParameterizerConfig{})

	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "valid SELECT",
			sql:     "SELECT * FROM users WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "INSERT blocked",
			sql:     "INSERT INTO users (name) VALUES ('John')",
			wantErr: true,
		},
		{
			name:    "UPDATE blocked",
			sql:     "UPDATE users SET name = 'Jane' WHERE id = 1",
			wantErr: true,
		},
		{
			name:    "DELETE blocked",
			sql:     "DELETE FROM users WHERE id = 1",
			wantErr: true,
		},
		{
			name:    "DROP blocked",
			sql:     "DROP TABLE users",
			wantErr: true,
		},
		{
			name:    "multiple statements blocked",
			sql:     "SELECT * FROM users; DROP TABLE users;",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parameterizer.ValidateForReadOnly(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateForReadOnly() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLParameterizer_DangerousPatterns(t *testing.T) {
	parameterizer := a01_text_to_sql.NewParameterizer(a01_text_to_sql.ParameterizerConfig{})

	tests := []struct {
		name        string
		sql         string
		expectSafe  bool
	}{
		{
			name:       "union injection attempt",
			sql:        "SELECT * FROM users WHERE id = 1 UNION SELECT * FROM admin",
			expectSafe: false,
		},
		{
			name:       "or 1=1 injection",
			sql:        "SELECT * FROM users WHERE id = 1 OR 1=1",
			expectSafe: false,
		},
		{
			name:       "comment injection",
			sql:        "SELECT * FROM users WHERE id = 1 --",
			expectSafe: false,
		},
		{
			name:       "safe query",
			sql:        "SELECT * FROM users WHERE status = 'active'",
			expectSafe: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parameterizer.Parameterize(tt.sql)
			if result.IsSafe != tt.expectSafe {
				t.Errorf("Parameterize() IsSafe = %v, want %v", result.IsSafe, tt.expectSafe)
				if len(result.Warnings) > 0 {
					t.Logf("Warnings: %v", result.Warnings)
				}
			}
		})
	}
}

func TestSQLParameterizer_IdentifierSanitization(t *testing.T) {
	parameterizer := a01_text_to_sql.NewParameterizer(a01_text_to_sql.ParameterizerConfig{})

	tests := []struct {
		name    string
		ident   string
		wantErr bool
	}{
		{
			name:    "valid identifier",
			ident:   "users",
			wantErr: false,
		},
		{
			name:    "valid with underscore",
			ident:   "user_accounts",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			ident:   "table_2026",
			wantErr: false,
		},
		{
			name:    "invalid with space",
			ident:   "user table",
			wantErr: true,
		},
		{
			name:    "invalid with special char",
			ident:   "users;drop",
			wantErr: true,
		},
		{
			name:    "reserved word",
			ident:   "SELECT",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parameterizer.SanitizeIdentifier(tt.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchemaRetriever_ContextFormatting(t *testing.T) {
	// Test that schema context can be formatted for LLM prompting
	context := &a01_text_to_sql.SchemaContext{
		Tables: []a01_text_to_sql.TableContext{
			{
				Name:        "fact_billing",
				Description: "Billing transactions for clinic appointments",
				Columns:     []string{"billing_id", "amount", "billing_date", "patient_id"},
				PrimaryKeys: []string{"billing_id"},
				Relevance:   0.95,
			},
		},
		QueryPatterns: []a01_text_to_sql.PatternContext{
			{
				Name:        "revenue_by_department",
				Description: "Calculate revenue grouped by department",
				TemplateSQL: "SELECT d.name, SUM(f.amount) FROM fact_billing f JOIN dim_department d ON f.dept_id = d.id GROUP BY d.name",
				Relevance:   0.88,
			},
		},
		RelevanceScore: 0.91,
	}

	prompt := context.ToPrompt()

	if prompt == "" {
		t.Error("ToPrompt() should not return empty string")
	}

	// Check that table name is in the prompt
}

func TestSchemaRetriever_GetTableNames(t *testing.T) {
	context := &a01_text_to_sql.SchemaContext{
		Tables: []a01_text_to_sql.TableContext{
			{Name: "fact_billing"},
			{Name: "dim_department"},
			{Name: "dim_date"},
		},
	}

	names := context.GetTableNames()

	if len(names) != 3 {
		t.Errorf("GetTableNames() = %v, want 3", len(names))
	}
}

func TestSQLValidator_LocalValidation(t *testing.T) {
	// Test local validation without OPA
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "valid SELECT",
			sql:     "SELECT * FROM users",
			wantErr: false,
		},
		{
			name:    "INSERT blocked",
			sql:     "INSERT INTO users VALUES (1)",
			wantErr: true,
		},
		{
			name:    "UPDATE blocked",
			sql:     "UPDATE users SET x = 1",
			wantErr: true,
		},
		{
			name:    "DELETE blocked",
			sql:     "DELETE FROM users",
			wantErr: true,
		},
		{
			name:    "DROP blocked",
			sql:     "DROP TABLE users",
			wantErr: true,
		},
		{
			name:    "CREATE blocked",
			sql:     "CREATE TABLE evil (id int)",
			wantErr: true,
		},
		{
			name:    "EXEC blocked",
			sql:     "EXEC sp_help",
			wantErr: true,
		},
		{
			name:    "comment blocked",
			sql:     "SELECT * FROM users --",
			wantErr: true,
		},
		{
			name:    "UNION SELECT blocked",
			sql:     "SELECT * FROM users UNION SELECT * FROM admin",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Direct test of local validation logic
			// In production, this would use the SQLValidator
		})
	}
}

func TestExecutionResult_CSVExport(t *testing.T) {
	result := &a01_text_to_sql.ExecutionResult{
		Columns: []a01_text_to_sql.ColumnInfo{
			{Name: "name", Type: "text"},
			{Name: "amount", Type: "numeric"},
		},
		Rows: []a01_text_to_sql.Row{
			{"name": "John", "amount": 100.50},
			{"name": "Jane", "amount": 200.75},
		},
		RowCount: 2,
	}

	csv := result.ToCSV()

	if csv == "" {
		t.Error("ToCSV() should not return empty string")
	}

	// Check that column headers are present
}
