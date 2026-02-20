// Package module_a_test tests the conversational BI agents (module A).
package module_a_test

import (
	"context"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/module_a/a06_confidence"
	"github.com/medisync/medisync/internal/warehouse/models"
)

func TestConfidenceScoring_HighConfidence(t *testing.T) {
	agent := a06_confidence.New(a06_confidence.AgentConfig{})

	tests := []struct {
		name          string
		req           a06_confidence.ScoreRequest
		minScore      float64
		expectRouting string
	}{
		{
			name: "simple query with high intent clarity",
			req: a06_confidence.ScoreRequest{
				QueryID:          "q1",
				UserQuery:        "Show me total revenue for January 2026",
				GeneratedSQL:     "SELECT SUM(amount) FROM fact_billing WHERE billing_date >= '2026-01-01'",
				SchemaMatches:    []string{"fact_billing", "dim_date"},
				RetryCount:       0,
				ExecutionTime:    100 * time.Millisecond,
				RowCount:         1,
				DetectedIntent:   "kpi",
				IntentConfidence: 0.95,
				ValidationPassed: true,
			},
			minScore:      70,
			expectRouting: models.RoutingNormal,
		},
		{
			name: "trend query with good schema match",
			req: a06_confidence.ScoreRequest{
				QueryID:          "q2",
				UserQuery:        "Show revenue trend over last 6 months",
				GeneratedSQL:     "SELECT month, SUM(amount) FROM fact_billing GROUP BY month",
				SchemaMatches:    []string{"fact_billing", "dim_date"},
				RetryCount:       0,
				ExecutionTime:    250 * time.Millisecond,
				RowCount:         6,
				DetectedIntent:   "trend",
				IntentConfidence: 0.88,
				ValidationPassed: true,
			},
			minScore:      70,
			expectRouting: models.RoutingNormal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Score(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("Score() error = %v", err)
			}
			if result.Score.Score < tt.minScore {
				t.Errorf("Score() = %v, want >= %v", result.Score.Score, tt.minScore)
			}
			if result.RoutingDecision != tt.expectRouting {
				t.Errorf("RoutingDecision = %v, want %v", result.RoutingDecision, tt.expectRouting)
			}
		})
	}
}

func TestConfidenceScoring_LowConfidence(t *testing.T) {
	agent := a06_confidence.New(a06_confidence.AgentConfig{})

	tests := []struct {
		name            string
		req             a06_confidence.ScoreRequest
		maxScore        float64
		expectQueue     bool
		expectClarify   bool
	}{
		{
			name: "ambiguous query with retries",
			req: a06_confidence.ScoreRequest{
				QueryID:          "q1",
				UserQuery:        "Show me the data",
				GeneratedSQL:     "SELECT * FROM fact_billing",
				SchemaMatches:    []string{"fact_billing", "fact_appointments", "fact_payments"},
				RetryCount:       2,
				ExecutionTime:    5 * time.Second,
				RowCount:         0,
				DetectedIntent:   "table",
				IntentConfidence: 0.4,
				ValidationPassed: true,
			},
			maxScore:      50,
			expectQueue:   true,
			expectClarify: true,
		},
		{
			name: "failed validation",
			req: a06_confidence.ScoreRequest{
				QueryID:          "q2",
				UserQuery:        "Show revenue",
				GeneratedSQL:     "SELECT * FROM unknown_table",
				SchemaMatches:    []string{},
				RetryCount:       3,
				ExecutionTime:    0,
				RowCount:         0,
				DetectedIntent:   "kpi",
				IntentConfidence: 0.6,
				ValidationPassed: false,
			},
			maxScore:      50,
			expectQueue:   true,
			expectClarify: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Score(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("Score() error = %v", err)
			}
			if result.Score.Score > tt.maxScore {
				t.Errorf("Score() = %v, want <= %v", result.Score.Score, tt.maxScore)
			}
			if result.ShouldQueue != tt.expectQueue {
				t.Errorf("ShouldQueue = %v, want %v", result.ShouldQueue, tt.expectQueue)
			}
			if result.ShouldClarify != tt.expectClarify {
				t.Errorf("ShouldClarify = %v, want %v", result.ShouldClarify, tt.expectClarify)
			}
		})
	}
}

func TestConfidenceScoring_MediumConfidence(t *testing.T) {
	agent := a06_confidence.New(a06_confidence.AgentConfig{})

	req := a06_confidence.ScoreRequest{
		QueryID:          "q1",
		UserQuery:        "Compare revenue across all departments",
		GeneratedSQL:     "SELECT dept, SUM(amount) FROM fact_billing GROUP BY dept",
		SchemaMatches:    []string{"fact_billing", "dim_department"},
		RetryCount:       1,
		ExecutionTime:    2 * time.Second,
		RowCount:         15,
		DetectedIntent:   "comparison",
		IntentConfidence: 0.75,
		ValidationPassed: true,
	}

	result, err := agent.Score(context.Background(), req)
	if err != nil {
		t.Fatalf("Score() error = %v", err)
	}

	// Medium confidence should be in warning range
	if result.Score.Score < 50 || result.Score.Score >= 70 {
		t.Errorf("Score() = %v, want between 50 and 70", result.Score.Score)
	}
	if !result.ShouldQueue {
		t.Error("Medium confidence should queue for review")
	}
}

func TestConfidenceFactors_Calculation(t *testing.T) {
	calc := a06_confidence.NewFactorCalculator()

	tests := []struct {
		name            string
		req             a06_confidence.ScoreRequest
		checkFactor     string
		minValue        float64
		maxValue        float64
	}{
		{
			name: "intent clarity with high confidence",
			req: a06_confidence.ScoreRequest{
				UserQuery:         "Show total revenue",
				IntentConfidence:  0.95,
				DetectedIntent:    "kpi",
				ValidationPassed:  true,
			},
			checkFactor: "intent_clarity",
			minValue:    0.8,
			maxValue:    1.0,
		},
		{
			name: "retry penalty for 2 retries",
			req: a06_confidence.ScoreRequest{
				RetryCount: 2,
			},
			checkFactor: "retry_penalty",
			minValue:    0.15,
			maxValue:    0.25,
		},
		{
			name: "SQL complexity for simple query",
			req: a06_confidence.ScoreRequest{
				GeneratedSQL: "SELECT * FROM users WHERE id = 1",
			},
			checkFactor: "complexity_penalty",
			minValue:    0.0,
			maxValue:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factors := calc.Calculate(tt.req)

			var value float64
			switch tt.checkFactor {
			case "intent_clarity":
				value = factors.IntentClarity
			case "retry_penalty":
				value = factors.RetryPenalty
			case "complexity_penalty":
				value = factors.SQLComplexityPenalty
			}

			if value < tt.minValue || value > tt.maxValue {
				t.Errorf("Factor %s = %v, want between %v and %v", tt.checkFactor, value, tt.minValue, tt.maxValue)
			}
		})
	}
}

func TestRoutingLogic_Decision(t *testing.T) {
	router := a06_confidence.NewRoutingLogic()

	tests := []struct {
		score          float64
		expectedAction string
	}{
		{95, a06_confidence.ActionNormal},
		{75, a06_confidence.ActionNormal},
		{70, a06_confidence.ActionNormal},
		{65, a06_confidence.ActionWarning},
		{55, a06_confidence.ActionWarning},
		{50, a06_confidence.ActionWarning},
		{45, a06_confidence.ActionClarify},
		{30, a06_confidence.ActionClarify},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			decision := router.Decide(tt.score)
			if decision.Action != tt.expectedAction {
				t.Errorf("Decide(%v) = %v, want %v", tt.score, decision.Action, tt.expectedAction)
			}
		})
	}
}

func TestRoutingLogic_QueueDecision(t *testing.T) {
	router := a06_confidence.NewRoutingLogic()

	tests := []struct {
		score      float64
		hasRetries bool
		shouldQueue bool
	}{
		{80, false, false},
		{80, true, false},
		{60, false, false},
		{60, true, true},
		{40, false, true},
		{40, true, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			shouldQueue := router.ShouldQueueForReview(tt.score, tt.hasRetries)
			if shouldQueue != tt.shouldQueue {
				t.Errorf("ShouldQueueForReview(%v, %v) = %v, want %v", tt.score, tt.hasRetries, shouldQueue, tt.shouldQueue)
			}
		})
	}
}

func TestConfidenceScore_Model(t *testing.T) {
	factors := models.ConfidenceFactors{
		IntentClarity:        0.9,
		SchemaMatchQuality:   0.85,
		SQLComplexityPenalty: 0.05,
		RetryPenalty:         0.0,
		HallucinationRisk:    0.02,
	}

	score := models.NewConfidenceScore("q1", factors)

	if score.Score < 70 {
		t.Errorf("Score = %v, want >= 70", score.Score)
	}
	if score.RoutingDecision != models.RoutingNormal {
		t.Errorf("RoutingDecision = %v, want %v", score.RoutingDecision, models.RoutingNormal)
	}
	if !score.IsAcceptable() {
		t.Error("IsAcceptable() should be true")
	}
	if score.RequiresClarification() {
		t.Error("RequiresClarification() should be false")
	}
}

func TestConfidenceScoring_AgentCard(t *testing.T) {
	agent := a06_confidence.New(a06_confidence.AgentConfig{})
	card := agent.AgentCard()

	if card["id"] != "a-06-confidence" {
		t.Errorf("AgentCard ID = %v, want a-06-confidence", card["id"])
	}
	if card["name"] == "" {
		t.Error("AgentCard Name should not be empty")
	}
}
