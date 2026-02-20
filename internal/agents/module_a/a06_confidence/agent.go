// Package a06_confidence provides the confidence scoring agent.
//
// This agent evaluates query results and assigns confidence scores
// to determine routing (normal/warning/clarify).
package a06_confidence

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/medisync/medisync/internal/warehouse/models"
)

// AgentID is the unique identifier for this agent.
const AgentID = "a-06-confidence"

// Agent implements the confidence scoring agent.
type Agent struct {
	id          string
	logger      *slog.Logger
	calculator  *FactorCalculator
	router      *RoutingLogic
	cache       map[string]*models.ConfidenceScore
	cacheMu     sync.RWMutex
}

// AgentConfig holds configuration for the agent.
type AgentConfig struct {
	Logger *slog.Logger
}

// New creates a new confidence scoring agent.
func New(cfg AgentConfig) *Agent {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &Agent{
		id:         AgentID,
		logger:     cfg.Logger.With("agent", AgentID),
		calculator: NewFactorCalculator(),
		router:     NewRoutingLogic(),
		cache:      make(map[string]*models.ConfidenceScore),
	}
}

// ScoreRequest contains the request for confidence scoring.
type ScoreRequest struct {
	QueryID           string            `json:"query_id"`
	UserQuery         string            `json:"user_query"`
	GeneratedSQL      string            `json:"generated_sql"`
	SchemaMatches     []string          `json:"schema_matches"`
	RetryCount        int               `json:"retry_count"`
	ExecutionTime     time.Duration     `json:"execution_time"`
	RowCount          int               `json:"row_count"`
	DetectedIntent    string            `json:"detected_intent"`
	IntentConfidence  float64           `json:"intent_confidence"`
	ValidationPassed  bool              `json:"validation_passed"`
	Context           map[string]string `json:"context,omitempty"`
}

// ScoreResponse contains the confidence score result.
type ScoreResponse struct {
	Score            *models.ConfidenceScore `json:"score"`
	RoutingDecision  string                  `json:"routing_decision"`
	ShouldQueue      bool                    `json:"should_queue"`
	ShouldClarify    bool                    `json:"should_clarify"`
	WarningMessage   string                  `json:"warning_message,omitempty"`
	ClarificationMsg string                  `json:"clarification_message,omitempty"`
}

// AgentCard returns the ADK agent card for discovery.
func (a *Agent) AgentCard() map[string]interface{} {
	return map[string]interface{}{
		"id":          AgentID,
		"name":        "Confidence Scoring Agent",
		"description": "Evaluates query results and assigns confidence scores",
		"capabilities": []string{
			"confidence-scoring",
			"multi-factor-evaluation",
			"routing-decision",
		},
		"version": "1.0.0",
	}
}

// Score calculates a confidence score for a query result.
func (a *Agent) Score(ctx context.Context, req ScoreRequest) (*ScoreResponse, error) {
	a.logger.Debug("scoring confidence",
		"query_id", req.QueryID,
		"retry_count", req.RetryCount,
		"intent_confidence", req.IntentConfidence)

	// Calculate individual factors
	factors := a.calculator.Calculate(req)

	// Create the confidence score
	score := models.NewConfidenceScore(req.QueryID, factors)

	// Determine routing
	response := &ScoreResponse{
		Score:           score,
		RoutingDecision: score.RoutingDecision,
		ShouldQueue:     score.NeedsReview(),
		ShouldClarify:   score.RequiresClarification(),
	}

	// Add appropriate messages
	switch score.RoutingDecision {
	case models.RoutingWarning:
		response.WarningMessage = a.getWarningMessage(score.Score)
	case models.RoutingClarify:
		response.ClarificationMsg = a.getClarificationMessage(req.UserQuery)
	}

	// Cache the score
	a.cacheScore(req.QueryID, score)

	a.logger.Info("confidence score calculated",
		"query_id", req.QueryID,
		"score", score.Score,
		"routing", score.RoutingDecision)

	return response, nil
}

// GetScore retrieves a cached confidence score.
func (a *Agent) GetScore(queryID string) (*models.ConfidenceScore, bool) {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	score, ok := a.cache[queryID]
	return score, ok
}

// cacheScore stores a confidence score in the cache.
func (a *Agent) cacheScore(queryID string, score *models.ConfidenceScore) {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()

	// Simple cache eviction
	if len(a.cache) >= 1000 {
		for k := range a.cache {
			delete(a.cache, k)
			break
		}
	}

	a.cache[queryID] = score
}

// getWarningMessage returns a warning message for low-confidence results.
func (a *Agent) getWarningMessage(score float64) string {
	return fmt.Sprintf(
		"This result has a confidence score of %.0f%%. Please verify the data before making decisions.",
		score,
	)
}

// getClarificationMessage returns a clarification request for very low confidence.
func (a *Agent) getClarificationMessage(query string) string {
	return fmt.Sprintf(
		"I'm not confident about the results for \"%s\". Could you please rephrase or provide more details?",
		truncate(query, 50),
	)
}

// BatchScore scores multiple queries in batch.
func (a *Agent) BatchScore(ctx context.Context, requests []ScoreRequest) ([]ScoreResponse, error) {
	responses := make([]ScoreResponse, len(requests))

	for i, req := range requests {
		resp, err := a.Score(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to score query %s: %w", req.QueryID, err)
		}
		responses[i] = *resp
	}

	return responses, nil
}

// GetMetrics returns confidence scoring metrics.
func (a *Agent) GetMetrics() map[string]interface{} {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	totalScores := len(a.cache)
	normalCount := 0
	warningCount := 0
	clarityCount := 0

	for _, score := range a.cache {
		switch score.RoutingDecision {
		case models.RoutingNormal:
			normalCount++
		case models.RoutingWarning:
			warningCount++
		case models.RoutingClarify:
			clarityCount++
		}
	}

	return map[string]interface{}{
		"total_scores":   totalScores,
		"normal_count":   normalCount,
		"warning_count":  warningCount,
		"clarify_count":  clarityCount,
		"normal_rate":    float64(normalCount) / float64(totalScores),
	}
}

// truncate truncates a string to maxLen characters.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ToJSON serializes the score response.
func (r *ScoreResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// FromJSON deserializes a score response.
func FromJSON(data string) (*ScoreResponse, error) {
	var resp ScoreResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}

// min returns the smaller of two floats.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
