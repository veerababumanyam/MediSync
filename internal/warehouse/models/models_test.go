package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestQuerySession_Validate(t *testing.T) {
	tests := []struct {
		name    string
		session *QuerySession
		wantErr bool
	}{
		{
			name: "valid session",
			session: &QuerySession{
				UserID:   uuid.New(),
				TenantID: uuid.New(),
				Locale:   "en",
			},
			wantErr: false,
		},
		{
			name: "valid arabic locale",
			session: &QuerySession{
				UserID:   uuid.New(),
				TenantID: uuid.New(),
				Locale:   "ar",
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			session: &QuerySession{
				UserID:   uuid.Nil,
				TenantID: uuid.New(),
				Locale:   "en",
			},
			wantErr: true,
		},
		{
			name: "missing tenant_id",
			session: &QuerySession{
				UserID:   uuid.New(),
				TenantID: uuid.Nil,
				Locale:   "en",
			},
			wantErr: true,
		},
		{
			name: "invalid locale",
			session: &QuerySession{
				UserID:   uuid.New(),
				TenantID: uuid.New(),
				Locale:   "fr",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerySession.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuerySession_NewSession(t *testing.T) {
	userID := uuid.New()
	tenantID := uuid.New()
	session := NewSession(userID, tenantID, "en")

	if session.ID == uuid.Nil {
		t.Error("NewSession() ID should not be nil")
	}
	if session.UserID != userID {
		t.Error("NewSession() UserID mismatch")
	}
	if session.TenantID != tenantID {
		t.Error("NewSession() TenantID mismatch")
	}
	if session.Locale != "en" {
		t.Error("NewSession() Locale mismatch")
	}
	if session.Metadata == nil {
		t.Error("NewSession() Metadata should be initialized")
	}
}

func TestQuerySession_JSON(t *testing.T) {
	session := &QuerySession{
		ID:       uuid.New(),
		UserID:   uuid.New(),
		TenantID: uuid.New(),
		Locale:   "en",
	}

	jsonData, err := session.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	var parsed QuerySession
	err = parsed.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	if parsed.ID != session.ID {
		t.Error("JSON roundtrip ID mismatch")
	}
	if parsed.Locale != session.Locale {
		t.Error("JSON roundtrip Locale mismatch")
	}
}

func TestNaturalLanguageQuery_Validate(t *testing.T) {
	tests := []struct {
		name   string
		query  *NaturalLanguageQuery
		wantEr bool
	}{
		{
			name: "valid query",
			query: &NaturalLanguageQuery{
				RawText:        "Show me revenue for January",
				DetectedLocale: "en",
				DetectedIntent: IntentTrend,
			},
			wantEr: false,
		},
		{
			name: "empty raw text",
			query: &NaturalLanguageQuery{
				RawText:        "",
				DetectedLocale: "en",
			},
			wantEr: true,
		},
		{
			name: "text too long",
			query: &NaturalLanguageQuery{
				RawText:        string(make([]byte, 2001)),
				DetectedLocale: "en",
			},
			wantEr: true,
		},
		{
			name: "invalid locale",
			query: &NaturalLanguageQuery{
				RawText:        "Test query",
				DetectedLocale: "de",
			},
			wantEr: true,
		},
		{
			name: "invalid intent",
			query: &NaturalLanguageQuery{
				RawText:        "Test query",
				DetectedLocale: "en",
				DetectedIntent: "invalid",
			},
			wantEr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if (err != nil) != tt.wantEr {
				t.Errorf("NaturalLanguageQuery.Validate() error = %v, wantEr %v", err, tt.wantEr)
			}
		})
	}
}

func TestSQLStatement_IsSELECT(t *testing.T) {
	tests := []struct {
		name      string
		sqlText   string
		isSELECT  bool
	}{
		{"simple select", "SELECT * FROM users", true},
		{"select with newlines", "\n  SELECT id, name\n  FROM users", true},
		{"lowercase select", "select * from users", true},
		{"mixed case", "sElEcT * FROM users", true},
		{"insert statement", "INSERT INTO users VALUES (1)", false},
		{"update statement", "UPDATE users SET name = 'test'", false},
		{"delete statement", "DELETE FROM users", false},
		{"drop statement", "DROP TABLE users", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := &SQLStatement{SQLText: tt.sqlText}
			if got := stmt.IsSELECT(); got != tt.isSELECT {
				t.Errorf("SQLStatement.IsSELECT() = %v, want %v", got, tt.isSELECT)
			}
		})
	}
}

func TestSQLStatement_IncrementRetry(t *testing.T) {
	stmt := &SQLStatement{RetryCount: 0}

	// Should succeed for first 3 increments
	for i := 0; i < MaxRetryCount; i++ {
		if err := stmt.IncrementRetry(); err != nil {
			t.Errorf("IncrementRetry() attempt %d error = %v", i+1, err)
		}
	}

	// Should fail on 4th attempt
	if err := stmt.IncrementRetry(); err == nil {
		t.Error("IncrementRetry() should fail when max retries reached")
	}

	if stmt.RetryCount != MaxRetryCount {
		t.Errorf("RetryCount = %d, want %d", stmt.RetryCount, MaxRetryCount)
	}
}

func TestSQLStatement_Parameterize(t *testing.T) {
	stmt := &SQLStatement{
		SQLText: "SELECT * FROM users WHERE tenant_id = :tenant_id",
	}

	params := map[string]any{
		"tenant_id": "123",
		"locale":    "en",
	}
	stmt.Parameterize(params)

	if !stmt.IsParameterized {
		t.Error("IsParameterized should be true after Parameterize()")
	}
	if stmt.Parameters["tenant_id"] != "123" {
		t.Error("Parameters not set correctly")
	}
}

func TestQueryResult_HasError(t *testing.T) {
	tests := []struct {
		name      string
		result    *QueryResult
		hasError  bool
	}{
		{"no error", &QueryResult{ErrorMessage: ""}, false},
		{"has error", &QueryResult{ErrorMessage: "query failed"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasError(); got != tt.hasError {
				t.Errorf("QueryResult.HasError() = %v, want %v", got, tt.hasError)
			}
		})
	}
}

func TestQueryResult_JSON(t *testing.T) {
	result := &QueryResult{
		ID:              uuid.New(),
		StatementID:     uuid.New(),
		RowCount:        2,
		Columns:         []ColumnMeta{{Name: "id", Type: "integer"}, {Name: "name", Type: "string"}},
		Data:            []map[string]any{{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}},
		ExecutionTimeMs: 150,
	}

	jsonData, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	var parsed QueryResult
	err = parsed.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	if parsed.RowCount != result.RowCount {
		t.Error("JSON roundtrip RowCount mismatch")
	}
	if len(parsed.Columns) != len(result.Columns) {
		t.Error("JSON roundtrip Columns mismatch")
	}
}

func TestDomainTerm_Matches(t *testing.T) {
	term := &DomainTerm{
		Synonym:       "revenue",
		CanonicalTerm: "total_revenue",
		Category:      CategoryAccounting,
		LocaleVariants: map[string][]string{
			"ar": {"الإيرادات", "دخل"},
		},
	}

	tests := []struct {
		text     string
		expected bool
	}{
		{"revenue", true},
		{"REVENUE", true},
		{"show me revenue", true},
		{"الإيرادات", true},
		{"دخل", true},
		{"show me الدخل for today", true},
		{"expense", false},
		{"cost", false},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			if got := term.Matches(tt.text); got != tt.expected {
				t.Errorf("DomainTerm.Matches(%q) = %v, want %v", tt.text, got, tt.expected)
			}
		})
	}
}

func TestDomainTerm_GetSynonyms(t *testing.T) {
	term := &DomainTerm{
		Synonym:       "revenue",
		CanonicalTerm: "total_revenue",
		Category:      CategoryAccounting,
		LocaleVariants: map[string][]string{
			"en": {"income", "earnings"},
			"ar": {"الإيرادات"},
		},
	}

	enSynonyms := term.GetSynonyms("en")
	if len(enSynonyms) != 2 {
		t.Errorf("GetSynonyms(en) returned %d items, want 2", len(enSynonyms))
	}

	arSynonyms := term.GetSynonyms("ar")
	if len(arSynonyms) != 1 {
		t.Errorf("GetSynonyms(ar) returned %d items, want 1", len(arSynonyms))
	}

	// Test fallback for unknown locale
	fallback := term.GetSynonyms("fr")
	if len(fallback) != 1 || fallback[0] != "revenue" {
		t.Errorf("GetSynonyms(fr) should fallback to main synonym")
	}
}

func TestDomainTerm_Validate(t *testing.T) {
	tests := []struct {
		name    string
		term    *DomainTerm
		wantErr bool
	}{
		{
			name: "valid term",
			term: &DomainTerm{
				Synonym:       "revenue",
				CanonicalTerm: "total_revenue",
				Category:      CategoryAccounting,
			},
			wantErr: false,
		},
		{
			name: "missing synonym",
			term: &DomainTerm{
				Synonym:       "",
				CanonicalTerm: "total_revenue",
				Category:      CategoryAccounting,
			},
			wantErr: true,
		},
		{
			name: "missing canonical term",
			term: &DomainTerm{
				Synonym:       "revenue",
				CanonicalTerm: "",
				Category:      CategoryAccounting,
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			term: &DomainTerm{
				Synonym:       "revenue",
				CanonicalTerm: "total_revenue",
				Category:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid locale variant",
			term: &DomainTerm{
				Synonym:        "revenue",
				CanonicalTerm:  "total_revenue",
				Category:       CategoryAccounting,
				LocaleVariants: map[string][]string{"de": {"einnahmen"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.term.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainTerm.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
