// Package module_e_test tests the i18n agents (module E).
package module_e_test

import (
	"context"
	"testing"

	"github.com/medisync/medisync/internal/agents/module_e/e01_language"
)

func TestLanguageDetection_English(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple English query",
			query:    "Show me total revenue for January",
			expected: "en",
		},
		{
			name:     "English with numbers",
			query:    "What is the revenue for Q1 2026?",
			expected: "en",
		},
		{
			name:     "English comparison query",
			query:    "Compare revenue by department for this quarter",
			expected: "en",
		},
		{
			name:     "English trend query",
			query:    "Show me the trend of patient visits over the last 6 months",
			expected: "en",
		},
		{
			name:     "English KPI query",
			query:    "What is the total number of appointments today?",
			expected: "en",
		},
	}

	agent := e01_language.New(&e01_language.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Detect(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}
			if result.Locale != tt.expected {
				t.Errorf("Detect() locale = %v, want %v", result.Locale, tt.expected)
			}
			if result.Confidence < 0.8 {
				t.Errorf("Detect() confidence = %v, want >= 0.8", result.Confidence)
			}
		})
	}
}

func TestLanguageDetection_Arabic(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple Arabic query",
			query:    "أظهر إجمالي الإيرادات لشهر يناير",
			expected: "ar",
		},
		{
			name:     "Arabic with numbers",
			query:    "ما هو عدد زيارات المرضى في الربع الأول ٢٠٢٦",
			expected: "ar",
		},
		{
			name:     "Arabic comparison query",
			query:    "قارن الإيرادات حسب القسم لهذا الربع",
			expected: "ar",
		},
		{
			name:     "Arabic trend query",
			query:    "أظهر اتجاه زيارات المرضى خلال الأشهر الستة الماضية",
			expected: "ar",
		},
		{
			name:     "Arabic KPI query",
			query:    "كم عدد المواعيد اليوم؟",
			expected: "ar",
		},
	}

	agent := e01_language.New(&e01_language.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Detect(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}
			if result.Locale != tt.expected {
				t.Errorf("Detect() locale = %v, want %v", result.Locale, tt.expected)
			}
			if result.Confidence < 0.8 {
				t.Errorf("Detect() confidence = %v, want >= 0.8", result.Confidence)
			}
		})
	}
}

func TestLanguageDetection_MixedContent(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectedMin float64
	}{
		{
			name:        "English with Arabic numbers",
			query:       "Revenue for January ٢٠٢٦",
			expectedMin: 0.5,
		},
		{
			name:        "Arabic with English brand name",
			query:       "أظهر إيرادات MediSync",
			expectedMin: 0.7,
		},
	}

	agent := e01_language.New(&e01_language.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Detect(context.Background(), tt.query)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}
			if result.Confidence < tt.expectedMin {
				t.Errorf("Detect() confidence = %v, want >= %v", result.Confidence, tt.expectedMin)
			}
		})
	}
}

func TestLanguageDetection_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		wantErr  bool
	}{
		{
			name:    "empty query",
			query:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			query:   "   ",
			wantErr: true,
		},
		{
			name:    "numbers only",
			query:   "12345",
			wantErr: false,
		},
		{
			name:    "single word",
			query:   "revenue",
			wantErr: false,
		},
	}

	agent := e01_language.New(&e01_language.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := agent.Detect(context.Background(), tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLanguageDetection_AgentCard(t *testing.T) {
	agent := e01_language.New(&e01_language.Config{})
	// AgentCard is part of the ADK integration (T084)
	// For now, verify the agent was created successfully
	if agent == nil {
		t.Error("Failed to create agent")
	}
}

func TestLanguageDetection_ConfidenceScoring(t *testing.T) {
	agent := e01_language.New(&e01_language.Config{})

	// High confidence cases should have high scores
	highConfidenceQueries := []string{
		"Show me total revenue for January 2026",
		"أظهر إجمالي الإيرادات",
	}

	for _, query := range highConfidenceQueries {
		result, err := agent.Detect(context.Background(), query)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if result.Confidence < 0.9 {
			t.Errorf("High confidence query '%s' got confidence = %v, want >= 0.9", query, result.Confidence)
		}
	}
}
