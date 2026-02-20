// Package integration_test provides integration tests for MediSync.
package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/module_a/a01_text_to_sql"
	"github.com/medisync/medisync/internal/agents/module_a/a04_terminology"
	"github.com/medisync/medisync/internal/agents/module_a/a06_confidence"
	"github.com/medisync/medisync/internal/agents/module_e/e01_language"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// TestChatFlow_EndToEnd tests the complete query processing flow.
func TestChatFlow_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Step 1: Language Detection
	t.Run("Step1_LanguageDetection", func(t *testing.T) {
		agent := e01_language.New(&e01_language.Config{})
		result, err := agent.Detect(ctx, "Show me total revenue for January 2026")
		if err != nil {
			t.Fatalf("Language detection failed: %v", err)
		}
		if result.Locale != "en" {
			t.Errorf("Expected locale 'en', got '%s'", result.Locale)
		}
	})

	// Step 2: Domain Terminology Normalization
	t.Run("Step2_TerminologyNormalization", func(t *testing.T) {
		mockGlossary := &mockGlossaryRepo{
			entries: []a04_terminology.GlossaryEntry{
				{
					Synonym:       "revenue",
					CanonicalTerm: "total_billed_amount",
					Category:      "accounting",
					SQLFragment:   "fact_billing.amount",
				},
			},
		}
		agent := a04_terminology.New(a04_terminology.AgentConfig{
			Glossary: mockGlossary,
		})

		result, err := agent.Normalize(ctx, a04_terminology.NormalizeRequest{
			Query: "Show me revenue for January",
			Locale: "en",
		})
		if err != nil {
			t.Fatalf("Terminology normalization failed: %v", err)
		}
		if len(result.AppliedMappings) == 0 {
			t.Error("Expected terminology mappings to be applied")
		}
	})

	// Step 3: SQL Generation (mocked for integration test)
	t.Run("Step3_SQLGeneration", func(t *testing.T) {
		// In a real integration test, this would call the LLM
		// For now, we verify the parameterizer works
		parameterizer := a01_text_to_sql.NewParameterizer(a01_text_to_sql.ParameterizerConfig{})

		sql := "SELECT SUM(amount) FROM fact_billing WHERE billing_date >= '2026-01-01'"
		result := parameterizer.Parameterize(sql)

		if !result.IsSafe {
			t.Errorf("SQL should be safe, warnings: %v", result.Warnings)
		}
	})

	// Step 4: Confidence Scoring
	t.Run("Step4_ConfidenceScoring", func(t *testing.T) {
		agent := a06_confidence.New(a06_confidence.AgentConfig{})

		req := a06_confidence.ScoreRequest{
			QueryID:          "test-query-1",
			UserQuery:        "Show me total revenue for January 2026",
			GeneratedSQL:     "SELECT SUM(amount) FROM fact_billing WHERE billing_date >= '2026-01-01'",
			SchemaMatches:    []string{"fact_billing", "dim_date"},
			RetryCount:       0,
			ExecutionTime:    150 * time.Millisecond,
			RowCount:         1,
			DetectedIntent:   "kpi",
			IntentConfidence: 0.95,
			ValidationPassed: true,
		}

		result, err := agent.Score(ctx, req)
		if err != nil {
			t.Fatalf("Confidence scoring failed: %v", err)
		}

		if result.Score.Score < 70 {
			t.Errorf("Expected high confidence score, got %.1f", result.Score.Score)
		}

		if result.RoutingDecision != models.RoutingNormal {
			t.Errorf("Expected normal routing, got %s", result.RoutingDecision)
		}
	})
}

// TestChatFlow_ArabicQuery tests Arabic query processing.
func TestChatFlow_ArabicQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("Arabic language detection", func(t *testing.T) {
		agent := e01_language.New(&e01_language.Config{})
		result, err := agent.Detect(ctx, "أظهر إجمالي إيرادات الصيدلية")
		if err != nil {
			t.Fatalf("Language detection failed: %v", err)
		}
		if result.Locale != "ar" {
			t.Errorf("Expected locale 'ar', got '%s'", result.Locale)
		}
	})
}

// TestChatFlow_LowConfidenceQuery tests handling of ambiguous queries.
func TestChatFlow_LowConfidenceQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	agent := a06_confidence.New(a06_confidence.AgentConfig{})

	// Ambiguous query with retries
	req := a06_confidence.ScoreRequest{
		QueryID:          "test-query-2",
		UserQuery:        "Show me the data",
		GeneratedSQL:     "SELECT * FROM fact_billing",
		SchemaMatches:    []string{"fact_billing", "fact_appointments", "fact_payments"},
		RetryCount:       2,
		ExecutionTime:    5 * time.Second,
		RowCount:         0,
		DetectedIntent:   "table",
		IntentConfidence: 0.4,
		ValidationPassed: true,
	}

	result, err := agent.Score(ctx, req)
	if err != nil {
		t.Fatalf("Confidence scoring failed: %v", err)
	}

	if result.Score.Score >= 70 {
		t.Errorf("Expected low confidence score, got %.1f", result.Score.Score)
	}

	if !result.ShouldQueue {
		t.Error("Low confidence query should be queued for review")
	}
}

// TestChatFlow_CompletePipeline tests the full agent pipeline.
func TestChatFlow_CompletePipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would verify the complete flow:
	// 1. Request received
	// 2. Authentication validated
	// 3. Language detected
	// 4. Query translated (if needed)
	// 5. Terminology normalized
	// 6. Schema context retrieved
	// 7. SQL generated
	// 8. SQL validated
	// 9. SQL executed
	// 10. Visualization routed
	// 11. Confidence scored
	// 12. Response formatted

	t.Run("Complete pipeline timing", func(t *testing.T) {
		start := time.Now()

		// Simulate pipeline execution
		ctx := context.Background()

		// Language detection
		langAgent := e01_language.New(&e01_language.Config{})
		_, _ = langAgent.Detect(ctx, "Show me revenue")

		// Confidence scoring
		confAgent := a06_confidence.New(a06_confidence.AgentConfig{})
		_, _ = confAgent.Score(ctx, a06_confidence.ScoreRequest{
			QueryID:          "test",
			UserQuery:        "test",
			IntentConfidence: 0.9,
			ValidationPassed: true,
		})

		elapsed := time.Since(start)

		// Pipeline should complete in reasonable time
		if elapsed > 5*time.Second {
			t.Errorf("Pipeline took too long: %v", elapsed)
		}
	})
}

// Mock implementations for testing

type mockGlossaryRepo struct {
	entries []a04_terminology.GlossaryEntry
}

func (m *mockGlossaryRepo) GetAll(ctx context.Context) ([]a04_terminology.GlossaryEntry, error) {
	return m.entries, nil
}

func (m *mockGlossaryRepo) GetBySynonym(ctx context.Context, synonym string) (*a04_terminology.GlossaryEntry, error) {
	for _, e := range m.entries {
		if e.Synonym == synonym {
			return &e, nil
		}
	}
	return nil, nil
}

func (m *mockGlossaryRepo) GetByCategory(ctx context.Context, category string) ([]a04_terminology.GlossaryEntry, error) {
	var result []a04_terminology.GlossaryEntry
	for _, e := range m.entries {
		if e.Category == category {
			result = append(result, e)
		}
	}
	return result, nil
}
