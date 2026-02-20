// Package module_e_test tests the i18n agents (module E).
package module_e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/medisync/medisync/internal/agents/module_e/e02_translation"
)

func TestQueryTranslation_ArabicToEnglish(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		minConfidence float64
	}{
		{
			name:          "revenue query",
			query:         "أظهر إجمالي إيرادات الصيدلية هذا الشهر",
			minConfidence: 0.7,
		},
		{
			name:          "trend query",
			query:         "ما هو اتجاه زيارات المرضى خلال الأشهر الستة الماضية",
			minConfidence: 0.7,
		},
		{
			name:          "comparison query",
			query:         "قارن الإيرادات حسب القسم لهذا الربع",
			minConfidence: 0.7,
		},
		{
			name:          "breakdown query",
			query:         "أظهر توزيع الإيرادات حسب نوع القسم",
			minConfidence: 0.7,
		},
	}

	agent := e02_translation.New(&e02_translation.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Translate(context.Background(), tt.query, "ar")
			if err != nil {
				t.Fatalf("Translate() error = %v", err)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("Translate() confidence = %v, want >= %v", result.Confidence, tt.minConfidence)
			}
			if result.TranslatedText == "" {
				t.Error("Translate() translated text should not be empty")
			}
		})
	}
}

func TestQueryTranslation_DomainTerms(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		expectedInResult string
	}{
		{
			name:             "healthcare term - visits",
			query:            "أظهر عدد الزيارات",
			expectedInResult: "visits",
		},
		{
			name:             "healthcare term - patients",
			query:            "كم عدد المرضى",
			expectedInResult: "patients",
		},
		{
			name:             "accounting term - revenue",
			query:            "أظهر الإيرادات",
			expectedInResult: "revenue",
		},
		{
			name:             "healthcare term - pharmacy",
			query:            "إيرادات الصيدلية",
			expectedInResult: "pharmacy",
		},
	}

	agent := e02_translation.New(&e02_translation.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Translate(context.Background(), tt.query, "ar")
			if err != nil {
				t.Fatalf("Translate() error = %v", err)
			}
			// Check if expected term is in the translated text
			if !strings.Contains(strings.ToLower(result.TranslatedText), tt.expectedInResult) {
				t.Errorf("Translate() expected %q in translated text, got %q", tt.expectedInResult, result.TranslatedText)
			}
		})
	}
}

func TestQueryTranslation_EdgeCases(t *testing.T) {
	agent := e02_translation.New(&e02_translation.Config{})

	t.Run("empty query", func(t *testing.T) {
		result, err := agent.Translate(context.Background(), "", "ar")
		if err != nil {
			t.Errorf("Translate() should not error for empty query, got: %v", err)
		}
		// Empty query returns empty result with confidence 1.0
		if result.TranslatedText != "" {
			t.Errorf("Translate() empty query should return empty translation, got %q", result.TranslatedText)
		}
	})

	t.Run("English query normalization", func(t *testing.T) {
		result, err := agent.Translate(context.Background(), "show revenue", "en")
		if err != nil {
			t.Fatalf("Translate() error = %v", err)
		}
		// English queries should be normalized but not translated
		if result.TranslatedText == "" {
			t.Error("Translate() should return normalized text for English queries")
		}
		if result.Confidence != 1.0 {
			t.Errorf("Translate() English query confidence should be 1.0, got %v", result.Confidence)
		}
	})
}

func TestQueryTranslation_DictionaryTranslation(t *testing.T) {
	agent := e02_translation.New(&e02_translation.Config{})

	tests := []struct {
		name     string
		query    string
		contains string
	}{
		{
			name:     "clinic term",
			query:    "العيادة",
			contains: "clinic",
		},
		{
			name:     "revenue term",
			query:    "الإيرادات",
			contains: "revenue",
		},
		{
			name:     "this month term",
			query:    "هذا الشهر",
			contains: "this month",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Translate(context.Background(), tt.query, "ar")
			if err != nil {
				t.Fatalf("Translate() error = %v", err)
			}
			if !strings.Contains(strings.ToLower(result.TranslatedText), tt.contains) {
				t.Errorf("Translate() expected %q in result, got %q", tt.contains, result.TranslatedText)
			}
		})
	}
}

func TestQueryTranslation_PreservedTerms(t *testing.T) {
	agent := e02_translation.New(&e02_translation.Config{})

	result, err := agent.Translate(context.Background(), "أظهر إيرادات العيادة", "ar")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}

	// Should have preserved terms detected
	if len(result.PreservedTerms) == 0 {
		t.Error("Translate() should detect preserved domain terms")
	}
}

func TestQueryTranslation_SourceTargetLocale(t *testing.T) {
	agent := e02_translation.New(&e02_translation.Config{})

	result, err := agent.Translate(context.Background(), "أظهر الإيرادات", "ar")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}

	if result.SourceLocale != "ar" {
		t.Errorf("SourceLocale = %q, want %q", result.SourceLocale, "ar")
	}
	if result.TargetLocale != "en" {
		t.Errorf("TargetLocale = %q, want %q", result.TargetLocale, "en")
	}
}
