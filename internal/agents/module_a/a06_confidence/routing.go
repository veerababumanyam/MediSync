// Package a06_confidence provides the confidence scoring agent.
//
// This file implements routing logic for confidence-based decisions.
package a06_confidence

import (
	"log/slog"
)

// RoutingLogic determines how to route based on confidence scores.
type RoutingLogic struct {
	normalThreshold  float64
	warningThreshold float64
	logger           *slog.Logger
}

// RoutingConfig holds configuration for routing logic.
type RoutingConfig struct {
	NormalThreshold  float64
	WarningThreshold float64
	Logger           *slog.Logger
}

// NewRoutingLogic creates a new routing logic handler.
func NewRoutingLogic() *RoutingLogic {
	return &RoutingLogic{
		normalThreshold:  70.0,
		warningThreshold: 50.0,
		logger:           slog.Default().With("component", "routing_logic"),
	}
}

// RoutingDecision represents a routing decision.
type RoutingDecision struct {
	Action       string  `json:"action"`
	Score        float64 `json:"score"`
	ShouldQueue  bool    `json:"should_queue"`
	ShouldWarn   bool    `json:"should_warn"`
	Message      string  `json:"message,omitempty"`
}

// Routing action constants
const (
	ActionNormal    = "normal"
	ActionWarning   = "warning"
	ActionClarify   = "clarify"
	ActionReview    = "review"
)

// Decide determines the routing decision based on score.
func (r *RoutingLogic) Decide(score float64) RoutingDecision {
	decision := RoutingDecision{
		Score: score,
	}

	switch {
	case score >= r.normalThreshold:
		decision.Action = ActionNormal
		decision.Message = "Result confidence is high enough for normal response"

	case score >= r.warningThreshold:
		decision.Action = ActionWarning
		decision.ShouldWarn = true
		decision.ShouldQueue = true
		decision.Message = "Result confidence is moderate, adding warning and queuing for review"

	default:
		decision.Action = ActionClarify
		decision.ShouldQueue = true
		decision.Message = "Result confidence is low, requesting user clarification"
	}

	return decision
}

// ShouldQueueForReview determines if a result should be queued for human review.
func (r *RoutingLogic) ShouldQueueForReview(score float64, hasRetries bool) bool {
	// Always queue if score is below warning threshold
	if score < r.warningThreshold {
		return true
	}

	// Queue if score is below normal and there were retries
	if score < r.normalThreshold && hasRetries {
		return true
	}

	return false
}

// ShouldRequestClarification determines if user clarification is needed.
func (r *RoutingLogic) ShouldRequestClarification(score float64, intentClarity float64) bool {
	// Need clarification if score is very low
	if score < r.warningThreshold {
		return true
	}

	// Need clarification if intent was unclear even with moderate score
	if score < r.normalThreshold && intentClarity < 0.6 {
		return true
	}

	return false
}

// GetResponseStrategy returns the strategy for responding to the user.
func (r *RoutingLogic) GetResponseStrategy(score float64) ResponseStrategy {
	switch {
	case score >= r.normalThreshold:
		return ResponseStrategy{
			ShowResult:    true,
			ShowWarning:   false,
			ShowSQL:       true,
			ShowConfidence: true,
			Priority:      "high",
		}

	case score >= r.warningThreshold:
		return ResponseStrategy{
			ShowResult:     true,
			ShowWarning:    true,
			ShowSQL:        true,
			ShowConfidence: true,
			Priority:       "medium",
		}

	default:
		return ResponseStrategy{
			ShowResult:     false,
			ShowWarning:    true,
			ShowSQL:        false,
			ShowConfidence: true,
			Priority:       "low",
		}
	}
}

// ResponseStrategy defines how to format the response.
type ResponseStrategy struct {
	ShowResult     bool   `json:"show_result"`
	ShowWarning    bool   `json:"show_warning"`
	ShowSQL        bool   `json:"show_sql"`
	ShowConfidence bool   `json:"show_confidence"`
	Priority       string `json:"priority"`
}

// BatchDecide makes routing decisions for multiple scores.
func (r *RoutingLogic) BatchDecide(scores []float64) []RoutingDecision {
	decisions := make([]RoutingDecision, len(scores))
	for i, score := range scores {
		decisions[i] = r.Decide(score)
	}
	return decisions
}

// GetQueuePriority returns the priority for queue processing.
func (r *RoutingLogic) GetQueuePriority(score float64) string {
	switch {
	case score < 30:
		return "critical"
	case score < 50:
		return "high"
	case score < 70:
		return "medium"
	default:
		return "low"
	}
}

// GetEscalationLevel returns the escalation level for human review.
func (r *RoutingLogic) GetEscalationLevel(score float64, retryCount int) string {
	// Combine score and retry count to determine escalation
	baseLevel := "none"

	switch {
	case score < 30:
		baseLevel = "senior"
	case score < 50:
		baseLevel = "standard"
	case score < 70:
		baseLevel = "junior"
	}

	// Escalate if there were multiple retries
	if retryCount >= 3 && baseLevel != "none" {
		switch baseLevel {
		case "junior":
			return "standard"
		case "standard":
			return "senior"
		}
	}

	return baseLevel
}
