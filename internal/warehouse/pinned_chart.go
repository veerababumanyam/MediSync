// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// PinnedChartRepository handles database operations for pinned charts.
type PinnedChartRepository struct {
	db *Repo
}

// NewPinnedChartRepository creates a new pinned chart repository.
func NewPinnedChartRepository(db *Repo) *PinnedChartRepository {
	return &PinnedChartRepository{db: db}
}

// Create inserts a new pinned chart.
func (r *PinnedChartRepository) Create(ctx context.Context, chart *models.PinnedChart) error {
	query := `
		INSERT INTO app.pinned_charts (user_id, title, query_id, natural_language_query, sql_query, chart_spec, chart_type, refresh_interval, locale, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		chart.UserID, chart.Title, chart.QueryID, chart.NaturalLanguageQuery,
		chart.SQLQuery, chart.ChartSpec, chart.ChartType,
		chart.RefreshInterval, chart.Locale, chart.Position,
	).Scan(&chart.ID, &chart.CreatedAt, &chart.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create pinned chart: %w", err)
	}

	return nil
}

// GetByUserID retrieves pinned charts for a user.
func (r *PinnedChartRepository) GetByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*models.PinnedChart, error) {
	query := `
		SELECT id, user_id, title, query_id, natural_language_query, sql_query, chart_spec, chart_type, refresh_interval, locale, position, last_refreshed_at, is_active, created_at, updated_at
		FROM app.pinned_charts
		WHERE user_id = $1
	`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY (position->>'row')::int, (position->>'col')::int"

	rows, err := r.db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get pinned charts: %w", err)
	}
	defer rows.Close()

	var charts []*models.PinnedChart
	for rows.Next() {
		chart := &models.PinnedChart{}
		var queryID pgtype.UUID
		var lastRefreshedAt pgtype.Timestamptz

		err := rows.Scan(
			&chart.ID, &chart.UserID, &chart.Title, &queryID, &chart.NaturalLanguageQuery,
			&chart.SQLQuery, &chart.ChartSpec, &chart.ChartType,
			&chart.RefreshInterval, &chart.Locale, &chart.Position,
			&lastRefreshedAt, &chart.IsActive, &chart.CreatedAt, &chart.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan pinned chart: %w", err)
		}

		if queryID.Valid {
			uid, err := uuid.FromBytes(queryID.Bytes[:])
			if err == nil {
				chart.QueryID = &uid
			}
		}
		if lastRefreshedAt.Valid {
			chart.LastRefreshedAt = &lastRefreshedAt.Time
		}

		charts = append(charts, chart)
	}

	return charts, nil
}

// Update updates a pinned chart.
func (r *PinnedChartRepository) Update(ctx context.Context, chart *models.PinnedChart) error {
	query := `
		UPDATE app.pinned_charts SET
			title = $2, refresh_interval = $3, position = $4, is_active = $5, updated_at = NOW()
		WHERE id = $1 AND user_id = $6
		RETURNING updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		chart.ID, chart.Title, chart.RefreshInterval, chart.Position, chart.IsActive, chart.UserID,
	).Scan(&chart.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update pinned chart: %w", err)
	}

	return nil
}

// Delete removes a pinned chart.
func (r *PinnedChartRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM app.pinned_charts WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete pinned chart: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: pinned chart not found")
	}

	return nil
}

// UpdatePosition updates only the position of a pinned chart.
func (r *PinnedChartRepository) UpdatePosition(ctx context.Context, id, userID uuid.UUID, position models.ChartPosition) error {
	query := `UPDATE app.pinned_charts SET position = $3, updated_at = NOW() WHERE id = $1 AND user_id = $2`

	_, err := r.db.pool.Exec(ctx, query, id, userID, position)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update chart position: %w", err)
	}

	return nil
}

// Reorder updates positions for multiple charts at once.
func (r *PinnedChartRepository) Reorder(ctx context.Context, userID uuid.UUID, positions []struct {
	ID       uuid.UUID
	Position models.ChartPosition
}) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("warehouse: failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE app.pinned_charts SET position = $3, updated_at = NOW() WHERE id = $1 AND user_id = $2`

	for _, p := range positions {
		_, err := tx.Exec(ctx, query, p.ID, userID, p.Position)
		if err != nil {
			return fmt.Errorf("warehouse: failed to reorder charts: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// UpdateLastRefreshed updates the last_refreshed_at timestamp.
func (r *PinnedChartRepository) UpdateLastRefreshed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE app.pinned_charts SET last_refreshed_at = NOW(), updated_at = NOW() WHERE id = $1`

	_, err := r.db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update last refreshed: %w", err)
	}

	return nil
}
