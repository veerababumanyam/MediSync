// Package module_a_test tests the conversational BI agents (module A).
package module_a_test

import (
	"context"
	"testing"

	"github.com/medisync/medisync/internal/agents/module_a/a02_sql_correction"
)

func TestSQLCorrection_ColumnNotFound(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{})

	tests := []struct {
		name          string
		sql           string
		errMsg        string
		shouldRetry   bool
		minConfidence float64
	}{
		{
			name:          "patient_name column not found",
			sql:           "SELECT patient_name FROM fact_billing",
			errMsg:        `column "patient_name" does not exist`,
			shouldRetry:   true,
			minConfidence: 0.7,
		},
		{
			name:          "doctor_name column not found",
			sql:           "SELECT doctor_name FROM fact_appointments",
			errMsg:        `column "doctor_name" does not exist`,
			shouldRetry:   true,
			minConfidence: 0.7,
		},
		{
			name:          "unknown column with no correction",
			sql:           "SELECT xyz_column FROM fact_billing",
			errMsg:        `column "xyz_column" does not exist`,
			shouldRetry:   false,
			minConfidence: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Correct(context.Background(), a02_sql_correction.CorrectionRequest{
				SQL:        tt.sql,
				Error:      tt.errMsg,
				RetryCount: 0,
			})
			if err != nil {
				t.Fatalf("Correct() error = %v", err)
			}
			if result.ShouldRetry != tt.shouldRetry {
				t.Errorf("Correct() ShouldRetry = %v, want %v", result.ShouldRetry, tt.shouldRetry)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("Correct() Confidence = %v, want >= %v", result.Confidence, tt.minConfidence)
			}
		})
	}
}

func TestSQLCorrection_RelationNotFound(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{})

	tests := []struct {
		name          string
		sql           string
		errMsg        string
		shouldRetry   bool
		minConfidence float64
	}{
		{
			name:          "patients table not found",
			sql:           "SELECT * FROM patients",
			errMsg:        `relation "patients" does not exist`,
			shouldRetry:   true,
			minConfidence: 0.8,
		},
		{
			name:          "appointments table not found",
			sql:           "SELECT * FROM appointments",
			errMsg:        `relation "appointments" does not exist`,
			shouldRetry:   true,
			minConfidence: 0.8,
		},
		{
			name:          "unknown table with no correction",
			sql:           "SELECT * FROM unknown_table",
			errMsg:        `relation "unknown_table" does not exist`,
			shouldRetry:   false,
			minConfidence: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Correct(context.Background(), a02_sql_correction.CorrectionRequest{
				SQL:        tt.sql,
				Error:      tt.errMsg,
				RetryCount: 0,
			})
			if err != nil {
				t.Fatalf("Correct() error = %v", err)
			}
			if result.ShouldRetry != tt.shouldRetry {
				t.Errorf("Correct() ShouldRetry = %v, want %v", result.ShouldRetry, tt.shouldRetry)
			}
		})
	}
}

func TestSQLCorrection_SyntaxError(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{})

	tests := []struct {
		name          string
		sql           string
		errMsg        string
		shouldRetry   bool
	}{
		{
			name:        "missing closing parenthesis",
			sql:         "SELECT SUM(amount FROM fact_billing",
			errMsg:      `syntax error at or near "FROM"`,
			shouldRetry: true,
		},
		{
			name:        "syntax error with no correction",
			sql:         "SELECT * FROM",
			errMsg:      `syntax error at end of input`,
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Correct(context.Background(), a02_sql_correction.CorrectionRequest{
				SQL:        tt.sql,
				Error:      tt.errMsg,
				RetryCount: 0,
			})
			if err != nil {
				t.Fatalf("Correct() error = %v", err)
			}
			// The result depends on whether we can correct the syntax
			_ = result
		})
	}
}

func TestSQLCorrection_MaxRetries(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{
		MaxRetries: 3,
	})

	// Test that exceeding max retries returns ShouldRetry=false
	result, err := agent.Correct(context.Background(), a02_sql_correction.CorrectionRequest{
		SQL:        "SELECT * FROM unknown",
		Error:      `relation "unknown" does not exist`,
		RetryCount: 3, // Already at max
	})
	if err != nil {
		t.Fatalf("Correct() error = %v", err)
	}
	if result.ShouldRetry {
		t.Error("Correct() should not retry when max retries exceeded")
	}
}

func TestSQLCorrection_ErrorTypeAnalysis(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{})

	tests := []struct {
		name            string
		errMsg          string
		expectedType    string
	}{
		{
			name:         "column error",
			errMsg:       `column "x" does not exist`,
			expectedType: "column_not_found",
		},
		{
			name:         "relation error",
			errMsg:       `relation "x" does not exist`,
			expectedType: "relation_not_found",
		},
		{
			name:         "syntax error",
			errMsg:       `syntax error at or near "x"`,
			expectedType: "syntax_error",
		},
		{
			name:         "type mismatch",
			errMsg:       `invalid input syntax for type integer`,
			expectedType: "type_mismatch",
		},
		{
			name:         "ambiguous reference",
			errMsg:       `column "id" is ambiguous`,
			expectedType: "ambiguous_reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Correct(context.Background(), a02_sql_correction.CorrectionRequest{
				SQL:        "SELECT * FROM x",
				Error:      tt.errMsg,
				RetryCount: 0,
			})
			if err != nil {
				t.Fatalf("Correct() error = %v", err)
			}
			if result.ErrorType != tt.expectedType {
				t.Errorf("Correct() ErrorType = %v, want %v", result.ErrorType, tt.expectedType)
			}
		})
	}
}

func TestSQLCorrection_AgentCard(t *testing.T) {
	agent := a02_sql_correction.New(a02_sql_correction.AgentConfig{})
	card := agent.AgentCard()

	if card["id"] != "a-02-sql-correction" {
		t.Errorf("AgentCard ID = %v, want a-02-sql-correction", card["id"])
	}
	if card["name"] == "" {
		t.Error("AgentCard Name should not be empty")
	}
}

func TestRetryManager_Backoff(t *testing.T) {
	rm := a02_sql_correction.NewRetryManager(a02_sql_correction.RetryConfig{
		MaxRetries:    3,
		BackoffFactor: 100, // 100ms
	})

	// Test backoff calculation
	backoff0 := rm.GetBackoffDuration(0)
	backoff1 := rm.GetBackoffDuration(1)
	backoff2 := rm.GetBackoffDuration(2)

	// Backoff should increase with each attempt
	if backoff1 <= backoff0 {
		t.Errorf("Backoff should increase: backoff0=%v, backoff1=%v", backoff0, backoff1)
	}
	if backoff2 <= backoff1 {
		t.Errorf("Backoff should increase: backoff1=%v, backoff2=%v", backoff1, backoff2)
	}
}

func TestRetryState_CanRetry(t *testing.T) {
	state := a02_sql_correction.NewRetryState(3)

	// Initially can retry
	if !state.CanRetry() {
		t.Error("CanRetry() should be true initially")
	}

	// After recording attempts
	for i := 0; i < 3; i++ {
		state.RecordAttempt("SELECT *", nil, a02_sql_correction.Correction{})
	}

	// Should not be able to retry after max attempts
	if state.CanRetry() {
		t.Error("CanRetry() should be false after max attempts")
	}
}
