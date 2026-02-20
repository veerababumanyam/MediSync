// Package module_a_test tests the conversational BI agents (module A).
package module_a_test

import (
	"context"
	"testing"

	"github.com/medisync/medisync/internal/agents/module_a/a04_terminology"
)

// MockGlossaryRepository implements GlossaryRepository for testing.
type MockGlossaryRepository struct {
	entries []a04_terminology.GlossaryEntry
}

func (m *MockGlossaryRepository) GetAll(ctx context.Context) ([]a04_terminology.GlossaryEntry, error) {
	return m.entries, nil
}

func (m *MockGlossaryRepository) GetBySynonym(ctx context.Context, synonym string) (*a04_terminology.GlossaryEntry, error) {
	for _, e := range m.entries {
		if e.Synonym == synonym {
			return &e, nil
		}
	}
	return nil, nil
}

func (m *MockGlossaryRepository) GetByCategory(ctx context.Context, category string) ([]a04_terminology.GlossaryEntry, error) {
	var result []a04_terminology.GlossaryEntry
	for _, e := range m.entries {
		if e.Category == category {
			result = append(result, e)
		}
	}
	return result, nil
}

func newMockGlossary() *MockGlossaryRepository {
	return &MockGlossaryRepository{
		entries: []a04_terminology.GlossaryEntry{
			{
				ID:            1,
				Synonym:       "footfall",
				CanonicalTerm: "patient_visits",
				Category:      "healthcare",
				SQLFragment:   "fact_appointments.appointment_id",
				LocaleVariants: map[string][]string{
					"en": {"walk-ins", "visits", "patient traffic"},
					"ar": {"زيارات", "حضور"},
				},
			},
			{
				ID:            2,
				Synonym:       "revenue",
				CanonicalTerm: "total_billed_amount",
				Category:      "accounting",
				SQLFragment:   "fact_billing.amount",
				LocaleVariants: map[string][]string{
					"en": {"income", "earnings", "sales"},
					"ar": {"إيرادات", "دخل"},
				},
			},
			{
				ID:            3,
				Synonym:       "clinic",
				CanonicalTerm: "clinic_department",
				Category:      "healthcare",
				SQLFragment:   "dim_department.department_type = 'clinic'",
				LocaleVariants: map[string][]string{
					"en": {"outpatient", "medical center"},
					"ar": {"عيادة"},
				},
			},
		},
	}
}

func TestTerminologyNormalizer_BasicNormalization(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	tests := []struct {
		name              string
		query             string
		locale            string
		expectedMappings  int
	}{
		{
			name:             "query with footfall",
			query:            "Show me footfall for this month",
			locale:           "en",
			expectedMappings: 1,
		},
		{
			name:             "query with revenue",
			query:            "Total revenue by department",
			locale:           "en",
			expectedMappings: 1,
		},
		{
			name:             "query with multiple terms",
			query:            "Show footfall and revenue for clinic",
			locale:           "en",
			expectedMappings: 3,
		},
		{
			name:             "query with no mappings",
			query:            "Show me the data",
			locale:           "en",
			expectedMappings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Normalize(context.Background(), a04_terminology.NormalizeRequest{
				Query: tt.query,
				Locale: tt.locale,
			})
			if err != nil {
				t.Fatalf("Normalize() error = %v", err)
			}
			if len(result.AppliedMappings) != tt.expectedMappings {
				t.Errorf("Normalize() mappings = %v, want %v", len(result.AppliedMappings), tt.expectedMappings)
			}
		})
	}
}

func TestTerminologyNormalizer_ArabicTerms(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	tests := []struct {
		name             string
		query            string
		locale           string
		expectedMappings int
	}{
		{
			name:             "Arabic revenue query",
			query:            "أظهر الإيرادات",
			locale:           "ar",
			expectedMappings: 1,
		},
		{
			name:             "Arabic visits query",
			query:            "عدد الزيارات",
			locale:           "ar",
			expectedMappings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Normalize(context.Background(), a04_terminology.NormalizeRequest{
				Query: tt.query,
				Locale: tt.locale,
			})
			if err != nil {
				t.Fatalf("Normalize() error = %v", err)
			}
			if len(result.AppliedMappings) < tt.expectedMappings {
				t.Errorf("Normalize() mappings = %v, want >= %v", len(result.AppliedMappings), tt.expectedMappings)
			}
		})
	}
}

func TestTerminologyNormalizer_GetSQLHints(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	tests := []struct {
		name          string
		query         string
		expectedHints int
	}{
		{
			name:          "query with footfall",
			query:         "Show me footfall",
			expectedHints: 1,
		},
		{
			name:          "query with multiple terms",
			query:         "Revenue by clinic",
			expectedHints: 2,
		},
		{
			name:          "query with no hints",
			query:         "Show me the data",
			expectedHints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hints, err := agent.GetSQLHints(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("GetSQLHints() error = %v", err)
			}
			if len(hints) < tt.expectedHints {
				t.Errorf("GetSQLHints() hints = %v, want >= %v", len(hints), tt.expectedHints)
			}
		})
	}
}

func TestTerminologyNormalizer_ExtractDomainContext(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	tests := []struct {
		name                 string
		query                string
		expectHealthcare     bool
		expectAccounting     bool
	}{
		{
			name:             "healthcare query",
			query:            "Show me footfall for clinic",
			expectHealthcare: true,
			expectAccounting: false,
		},
		{
			name:             "accounting query",
			query:            "Total revenue for this month",
			expectHealthcare: false,
			expectAccounting: true,
		},
		{
			name:             "mixed query",
			query:            "Revenue by clinic",
			expectHealthcare: true,
			expectAccounting: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := agent.ExtractDomainContext(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("ExtractDomainContext() error = %v", err)
			}

			if tt.expectHealthcare && len(context.HealthcareTerms) == 0 {
				t.Error("Expected healthcare terms but got none")
			}
			if tt.expectAccounting && len(context.AccountingTerms) == 0 {
				t.Error("Expected accounting terms but got none")
			}
		})
	}
}

func TestTerminologyNormalizer_ConfidenceScoring(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	tests := []struct {
		name             string
		query            string
		minConfidence    float64
		maxConfidence    float64
	}{
		{
			name:          "no mappings - high confidence",
			query:         "Show me the data",
			minConfidence: 0.9,
			maxConfidence: 1.0,
		},
		{
			name:          "multiple mappings - lower confidence",
			query:         "Show footfall and revenue for clinic",
			minConfidence: 0.7,
			maxConfidence: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Normalize(context.Background(), a04_terminology.NormalizeRequest{
				Query: tt.query,
				Locale: "en",
			})
			if err != nil {
				t.Fatalf("Normalize() error = %v", err)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("Normalize() confidence = %v, want >= %v", result.Confidence, tt.minConfidence)
			}
			if result.Confidence > tt.maxConfidence {
				t.Errorf("Normalize() confidence = %v, want <= %v", result.Confidence, tt.maxConfidence)
			}
		})
	}
}

func TestTerminologyNormalizer_AgentCard(t *testing.T) {
	agent := a04_terminology.New(a04_terminology.AgentConfig{
		Glossary: newMockGlossary(),
	})

	card := agent.AgentCard()

	if card.ID != "a-04-terminology" {
		t.Errorf("AgentCard ID = %v, want a-04-terminology", card.ID)
	}
	if card.Name == "" {
		t.Error("AgentCard Name should not be empty")
	}
	if len(card.Capabilities) == 0 {
		t.Error("AgentCard should have capabilities")
	}
}
