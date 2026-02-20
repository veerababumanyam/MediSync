// Package module_a_test tests the conversational BI agents (module A).
package module_a_test

import (
	"context"
	"testing"

	"github.com/medisync/medisync/internal/agents/module_a/a05_hallucination"
)

func TestHallucinationGuard_OnTopicQueries(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})

	tests := []struct {
		name          string
		query         string
		shouldBeOnTopic bool
		minConfidence float64
	}{
		{
			name:            "revenue query",
			query:           "Show me total revenue for January 2026",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
		{
			name:            "patient visits query",
			query:           "How many patient visits did we have last month?",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
		{
			name:            "pharmacy sales query",
			query:           "What are the pharmacy sales for this quarter?",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
		{
			name:            "doctor performance query",
			query:           "Show doctor performance metrics",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
		{
			name:            "inventory query",
			query:           "What is the current inventory level for medicines?",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
		{
			name:            "billing summary query",
			query:           "Show billing summary by department",
			shouldBeOnTopic: true,
			minConfidence:   0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
				Query: tt.query,
				Locale: "en",
			})
			if err != nil {
				t.Fatalf("Guard() error = %v", err)
			}
			if result.IsOnTopic != tt.shouldBeOnTopic {
				t.Errorf("Guard() IsOnTopic = %v, want %v", result.IsOnTopic, tt.shouldBeOnTopic)
			}
			if result.Confidence < tt.minConfidence {
				t.Errorf("Guard() Confidence = %v, want >= %v", result.Confidence, tt.minConfidence)
			}
		})
	}
}

func TestHallucinationGuard_OffTopicQueries(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})

	tests := []struct {
		name            string
		query           string
		shouldBeOnTopic bool
	}{
		{
			name:            "weather query",
			query:           "What's the weather like today?",
			shouldBeOnTopic: false,
		},
		{
			name:            "poem request",
			query:           "Write me a poem about spring",
			shouldBeOnTopic: false,
		},
		{
			name:            "programming question",
			query:           "How do I implement a binary search tree?",
			shouldBeOnTopic: false,
		},
		{
			name:            "movie recommendation",
			query:           "What's the best movie to watch tonight?",
			shouldBeOnTopic: false,
		},
		{
			name:            "travel booking",
			query:           "Book me a flight to Paris",
			shouldBeOnTopic: false,
		},
		{
			name:            "recipe request",
			query:           "Give me a recipe for chocolate cake",
			shouldBeOnTopic: false,
		},
		{
			name:            "math homework",
			query:           "Solve this equation: 2x + 5 = 15",
			shouldBeOnTopic: false,
		},
		{
			name:            "general knowledge",
			query:           "Who is the president of the United States?",
			shouldBeOnTopic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
				Query: tt.query,
				Locale: "en",
			})
			if err != nil {
				t.Fatalf("Guard() error = %v", err)
			}
			if result.IsOnTopic != tt.shouldBeOnTopic {
				t.Errorf("Guard() IsOnTopic = %v, want %v (category: %s)", result.IsOnTopic, tt.shouldBeOnTopic, result.Category)
			}
		})
	}
}

func TestHallucinationGuard_RejectionMessage(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})

	tests := []struct {
		name            string
		query           string
		locale          string
		shouldBeOnTopic bool
	}{
		{
			name:            "English rejection",
			query:           "What's the weather?",
			locale:          "en",
			shouldBeOnTopic: false,
		},
		{
			name:            "Arabic rejection",
			query:           "كيف الطقس اليوم؟",
			locale:          "ar",
			shouldBeOnTopic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
				Query: tt.query,
				Locale: tt.locale,
			})
			if err != nil {
				t.Fatalf("Guard() error = %v", err)
			}
			if result.IsOnTopic != tt.shouldBeOnTopic {
				t.Errorf("Guard() IsOnTopic = %v, want %v", result.IsOnTopic, tt.shouldBeOnTopic)
			}
			if !result.IsOnTopic && result.RejectionMessage == "" {
				t.Error("Guard() should return rejection message for off-topic queries")
			}
		})
	}
}

func TestHallucinationGuard_AmbiguousQueries(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})

	tests := []struct {
		name                   string
		query                  string
		expectClarification    bool
	}{
		{
			name:                "very short query",
			query:               "Show data",
			expectClarification: true,
		},
		{
			name:                "vague query",
			query:               "Give me the numbers",
			expectClarification: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
				Query: tt.query,
				Locale: "en",
			})
			if err != nil {
				t.Fatalf("Guard() error = %v", err)
			}
			if tt.expectClarification && !result.NeedsClarification {
				t.Error("Guard() should request clarification for ambiguous queries")
			}
		})
	}
}

func TestHallucinationGuard_AgentCard(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})
	card := agent.AgentCard()

	if card["id"] != "a-05-hallucination" {
		t.Errorf("AgentCard ID = %v, want a-05-hallucination", card["id"])
	}
	if card["name"] == "" {
		t.Error("AgentCard Name should not be empty")
	}
}

func TestHallucinationGuard_Caching(t *testing.T) {
	agent := a05_hallucination.New(a05_hallucination.AgentConfig{})

	query := "Show me revenue by department"

	// First call
	result1, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
		Query: query,
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("Guard() error = %v", err)
	}

	// Second call (should hit cache)
	result2, err := agent.Guard(context.Background(), a05_hallucination.GuardRequest{
		Query: query,
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("Guard() error = %v", err)
	}

	if result1.IsOnTopic != result2.IsOnTopic {
		t.Error("Cache should return same result")
	}
}
