// Package models defines data models for the MediSync warehouse.
//
// This file contains the ConfidenceScore model.
package models

import (
	"encoding/json"
	"time"
)

// ConfidenceScore represents a numerical assessment of result accuracy.
type ConfidenceScore struct {
	// ID is the unique identifier
	ID string `json:"id"`
	// QueryID references the parent query
	QueryID string `json:"query_id"`
	// Score is the confidence percentage (0-100)
	Score float64 `json:"score"`
	// Factors contains the scoring factor breakdown
	Factors ConfidenceFactors `json:"factors"`
	// RoutingDecision is the action taken based on score
	RoutingDecision string `json:"routing_decision"`
	// CreatedAt is the scoring timestamp
	CreatedAt time.Time `json:"created_at"`
}

// ConfidenceFactors contains the individual factors that contribute to the score.
type ConfidenceFactors struct {
	// IntentClarity measures how clear the user's intent was (0-1)
	IntentClarity float64 `json:"intent_clarity"`
	// SchemaMatchQuality measures how well the schema matched the query (0-1)
	SchemaMatchQuality float64 `json:"schema_match_quality"`
	// SQLComplexityPenalty penalizes complex SQL (0-0.3)
	SQLComplexityPenalty float64 `json:"sql_complexity_penalty"`
	// RetryPenalty penalizes queries that needed correction (0-0.3)
	RetryPenalty float64 `json:"retry_penalty"`
	// HallucinationRisk measures the risk of hallucinated data (0-1)
	HallucinationRisk float64 `json:"hallucination_risk"`
}

// RoutingDecision constants
const (
	RoutingNormal   = "normal"   // Score >= 70
	RoutingWarning  = "warning"  // Score 50-69
	RoutingClarify  = "clarify"  // Score < 50
)

// NewConfidenceScore creates a new confidence score.
func NewConfidenceScore(queryID string, factors ConfidenceFactors) *ConfidenceScore {
	score := CalculateScore(factors)
	return &ConfidenceScore{
		QueryID:         queryID,
		Score:           score,
		Factors:         factors,
		RoutingDecision: DetermineRouting(score),
		CreatedAt:       time.Now(),
	}
}

// CalculateScore calculates the overall confidence score from factors.
func CalculateScore(factors ConfidenceFactors) float64 {
	// Base score from clarity and schema match
	baseScore := (factors.IntentClarity + factors.SchemaMatchQuality) / 2

	// Apply penalties
	penalty := factors.SQLComplexityPenalty + factors.RetryPenalty + factors.HallucinationRisk

	// Calculate final score (0-100 scale)
	score := (baseScore - penalty) * 100

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// DetermineRouting determines the routing decision based on score.
func DetermineRouting(score float64) string {
	if score >= 70 {
		return RoutingNormal
	}
	if score >= 50 {
		return RoutingWarning
	}
	return RoutingClarify
}

// NeedsReview returns true if the query should be added to the review queue.
func (c *ConfidenceScore) NeedsReview() bool {
	return c.Score < 70
}

// IsAcceptable returns true if the confidence is high enough for normal response.
func (c *ConfidenceScore) IsAcceptable() bool {
	return c.Score >= 70
}

// RequiresClarification returns true if the user should be asked for clarification.
func (c *ConfidenceScore) RequiresClarification() bool {
	return c.Score < 50
}

// ToJSON serializes the confidence score.
func (c *ConfidenceScore) ToJSON() string {
	data, _ := json.Marshal(c)
	return string(data)
}

// FromJSON deserializes a confidence score.
func ConfidenceFromJSON(data string) (*ConfidenceScore, error) {
	var score ConfidenceScore
	if err := json.Unmarshal([]byte(data), &score); err != nil {
		return nil, err
	}
	return &score, nil
}

// ReviewQueueEntry represents an item in the review queue.
type ReviewQueueEntry struct {
	// ID is the unique identifier
	ID string `json:"id"`
	// QueryID references the low-confidence query
	QueryID string `json:"query_id"`
	// ScoreID references the confidence score
	ScoreID string `json:"score_id"`
	// Status is the review status (pending/reviewed/resolved/dismissed)
	Status string `json:"status"`
	// ReviewedBy is the user who reviewed
	ReviewedBy string `json:"reviewed_by,omitempty"`
	// ReviewedAt is when it was reviewed
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	// Resolution contains the resolution notes
	Resolution string `json:"resolution,omitempty"`
	// CreatedAt is when the entry was created
	CreatedAt time.Time `json:"created_at"`
}

// Review status constants
const (
	ReviewStatusPending   = "pending"
	ReviewStatusReviewed  = "reviewed"
	ReviewStatusResolved  = "resolved"
	ReviewStatusDismissed = "dismissed"
)

// NewReviewQueueEntry creates a new review queue entry.
func NewReviewQueueEntry(queryID, scoreID string) *ReviewQueueEntry {
	return &ReviewQueueEntry{
		QueryID:   queryID,
		ScoreID:   scoreID,
		Status:    ReviewStatusPending,
		CreatedAt: time.Now(),
	}
}

// MarkReviewed marks the entry as reviewed.
func (e *ReviewQueueEntry) MarkReviewed(reviewedBy, resolution string) {
	now := time.Now()
	e.Status = ReviewStatusReviewed
	e.ReviewedBy = reviewedBy
	e.ReviewedAt = &now
	e.Resolution = resolution
}

// Resolve marks the entry as resolved.
func (e *ReviewQueueEntry) Resolve(reviewedBy, resolution string) {
	e.MarkReviewed(reviewedBy, resolution)
	e.Status = ReviewStatusResolved
}

// Dismiss marks the entry as dismissed.
func (e *ReviewQueueEntry) Dismiss(reviewedBy string) {
	now := time.Now()
	e.Status = ReviewStatusDismissed
	e.ReviewedBy = reviewedBy
	e.ReviewedAt = &now
}

// IsPending returns true if the entry is pending review.
func (e *ReviewQueueEntry) IsPending() bool {
	return e.Status == ReviewStatusPending
}
