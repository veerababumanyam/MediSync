// Package warehouse provides database utilities for the MediSync data warehouse.
//
// This file provides review queue functionality for low-confidence queries.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// ReviewQueue manages low-confidence query reviews.
type ReviewQueue struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewReviewQueue creates a new review queue manager.
func NewReviewQueue(pool *pgxpool.Pool, logger *slog.Logger) *ReviewQueue {
	if logger == nil {
		logger = slog.Default()
	}
	return &ReviewQueue{
		pool:   pool,
		logger: logger.With("component", "review_queue"),
	}
}

// AddEntry adds a new entry to the review queue.
func (q *ReviewQueue) AddEntry(ctx context.Context, entry *models.ReviewQueueEntry) error {
	query := `
		INSERT INTO app.review_queue (query_id, score_id, status, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := q.pool.QueryRow(ctx, query,
		entry.QueryID,
		entry.ScoreID,
		entry.Status,
		entry.CreatedAt,
	).Scan(&entry.ID)

	if err != nil {
		return fmt.Errorf("failed to add review entry: %w", err)
	}

	q.logger.Info("review queue entry added",
		"entry_id", entry.ID,
		"query_id", entry.QueryID)

	return nil
}

// GetPending retrieves pending review entries.
func (q *ReviewQueue) GetPending(ctx context.Context, limit int) ([]models.ReviewQueueEntry, error) {
	query := `
		SELECT id, query_id, score_id, status, reviewed_by, reviewed_at, resolution, created_at
		FROM app.review_queue
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := q.pool.Query(ctx, query, models.ReviewStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	defer rows.Close()

	var entries []models.ReviewQueueEntry
	for rows.Next() {
		var entry models.ReviewQueueEntry
		err := rows.Scan(
			&entry.ID,
			&entry.QueryID,
			&entry.ScoreID,
			&entry.Status,
			&entry.ReviewedBy,
			&entry.ReviewedAt,
			&entry.Resolution,
			&entry.CreatedAt,
		)
		if err != nil {
			q.logger.Warn("failed to scan review entry", "error", err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetByQueryID retrieves a review entry by query ID.
func (q *ReviewQueue) GetByQueryID(ctx context.Context, queryID string) (*models.ReviewQueueEntry, error) {
	query := `
		SELECT id, query_id, score_id, status, reviewed_by, reviewed_at, resolution, created_at
		FROM app.review_queue
		WHERE query_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var entry models.ReviewQueueEntry
	err := q.pool.QueryRow(ctx, query, queryID).Scan(
		&entry.ID,
		&entry.QueryID,
		&entry.ScoreID,
		&entry.Status,
		&entry.ReviewedBy,
		&entry.ReviewedAt,
		&entry.Resolution,
		&entry.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("review entry not found: %w", err)
	}

	return &entry, nil
}

// UpdateStatus updates the status of a review entry.
func (q *ReviewQueue) UpdateStatus(ctx context.Context, entryID string, status string, reviewedBy string, resolution string) error {
	query := `
		UPDATE app.review_queue
		SET status = $1, reviewed_by = $2, reviewed_at = $3, resolution = $4
		WHERE id = $5
	`

	_, err := q.pool.Exec(ctx, query,
		status,
		reviewedBy,
		time.Now(),
		resolution,
		entryID,
	)

	if err != nil {
		return fmt.Errorf("failed to update review status: %w", err)
	}

	q.logger.Info("review entry updated",
		"entry_id", entryID,
		"status", status,
		"reviewed_by", reviewedBy)

	return nil
}

// Resolve marks an entry as resolved with a resolution.
func (q *ReviewQueue) Resolve(ctx context.Context, entryID string, reviewedBy string, resolution string) error {
	return q.UpdateStatus(ctx, entryID, models.ReviewStatusResolved, reviewedBy, resolution)
}

// Dismiss dismisses a review entry.
func (q *ReviewQueue) Dismiss(ctx context.Context, entryID string, reviewedBy string) error {
	return q.UpdateStatus(ctx, entryID, models.ReviewStatusDismissed, reviewedBy, "")
}

// GetStats returns statistics about the review queue.
func (q *ReviewQueue) GetStats(ctx context.Context) (*ReviewQueueStats, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'reviewed') as reviewed,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved,
			COUNT(*) FILTER (WHERE status = 'dismissed') as dismissed,
			COUNT(*) as total
		FROM app.review_queue
	`

	var stats ReviewQueueStats
	err := q.pool.QueryRow(ctx, query).Scan(
		&stats.Pending,
		&stats.Reviewed,
		&stats.Resolved,
		&stats.Dismissed,
		&stats.Total,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get review stats: %w", err)
	}

	return &stats, nil
}

// ReviewQueueStats contains statistics about the review queue.
type ReviewQueueStats struct {
	Pending   int `json:"pending"`
	Reviewed  int `json:"reviewed"`
	Resolved  int `json:"resolved"`
	Dismissed int `json:"dismissed"`
	Total     int `json:"total"`
}

// PurgeOld removes old resolved/dismissed entries.
func (q *ReviewQueue) PurgeOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	query := `
		DELETE FROM app.review_queue
		WHERE status IN ('resolved', 'dismissed')
		AND created_at < $1
	`

	result, err := q.pool.Exec(ctx, query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to purge old entries: %w", err)
	}

	deleted := result.RowsAffected()
	q.logger.Info("purged old review entries", "count", deleted)

	return deleted, nil
}

// ReviewEntryWithDetails contains a review entry with additional context.
type ReviewEntryWithDetails struct {
	Entry         models.ReviewQueueEntry `json:"entry"`
	Query         string                  `json:"query"`
	GeneratedSQL  string                  `json:"generated_sql"`
	Score         float64                 `json:"score"`
	RoutingDecision string                `json:"routing_decision"`
	Factors       models.ConfidenceFactors `json:"factors"`
}

// GetWithDetails retrieves a review entry with full context.
func (q *ReviewQueue) GetWithDetails(ctx context.Context, entryID string) (*ReviewEntryWithDetails, error) {
	query := `
		SELECT
			r.id, r.query_id, r.score_id, r.status, r.reviewed_by, r.reviewed_at, r.resolution, r.created_at,
			q.raw_text,
			s.sql_text,
			c.score, c.routing_decision, c.factors
		FROM app.review_queue r
		JOIN app.queries q ON r.query_id = q.id
		JOIN app.sql_statements s ON q.id = s.query_id
		JOIN app.confidence_scores c ON r.score_id = c.id
		WHERE r.id = $1
	`

	var details ReviewEntryWithDetails
	var factorsJSON []byte

	err := q.pool.QueryRow(ctx, query, entryID).Scan(
		&details.Entry.ID,
		&details.Entry.QueryID,
		&details.Entry.ScoreID,
		&details.Entry.Status,
		&details.Entry.ReviewedBy,
		&details.Entry.ReviewedAt,
		&details.Entry.Resolution,
		&details.Entry.CreatedAt,
		&details.Query,
		&details.GeneratedSQL,
		&details.Score,
		&details.RoutingDecision,
		&factorsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get review details: %w", err)
	}

	// Parse factors JSON
	if len(factorsJSON) > 0 {
		if err := json.Unmarshal(factorsJSON, &details.Factors); err != nil {
			q.logger.Warn("failed to parse factors JSON", "error", err)
		}
	}

	return &details, nil
}
