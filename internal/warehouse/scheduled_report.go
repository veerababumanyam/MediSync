// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// ScheduledReportRepository handles database operations for scheduled reports.
type ScheduledReportRepository struct {
	db *Repo
}

// NewScheduledReportRepository creates a new scheduled report repository.
func NewScheduledReportRepository(db *Repo) *ScheduledReportRepository {
	return &ScheduledReportRepository{db: db}
}

// Create inserts a new scheduled report.
func (r *ScheduledReportRepository) Create(ctx context.Context, report *models.ScheduledReport) error {
	query := `
		INSERT INTO app.scheduled_reports (user_id, name, description, query_id, natural_language_query, sql_query, schedule_type, schedule_time, schedule_day, recipients, format, locale, include_charts, next_run_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		report.UserID, report.Name, report.Description, report.QueryID,
		report.NaturalLanguageQuery, report.SQLQuery, report.ScheduleType,
		report.ScheduleTime, report.ScheduleDay, report.Recipients,
		report.Format, report.Locale, report.IncludeCharts, report.NextRunAt,
	).Scan(&report.ID, &report.CreatedAt, &report.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create scheduled report: %w", err)
	}

	return nil
}

// GetByUserID retrieves scheduled reports for a user.
func (r *ScheduledReportRepository) GetByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*models.ScheduledReport, error) {
	query := `
		SELECT id, user_id, name, description, query_id, natural_language_query, sql_query, schedule_type, schedule_time, schedule_day, recipients, format, locale, include_charts, last_run_at, next_run_at, is_active, created_at, updated_at
		FROM app.scheduled_reports
		WHERE user_id = $1
	`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY next_run_at ASC NULLS LAST, created_at DESC"

	rows, err := r.db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get scheduled reports: %w", err)
	}
	defer rows.Close()

	var reports []*models.ScheduledReport
	for rows.Next() {
		report := &models.ScheduledReport{}
		var description pgtype.Text
		var queryID pgtype.UUID
		var scheduleDay pgtype.Int4
		var lastRunAt, nextRunAt pgtype.Timestamptz

		err := rows.Scan(
			&report.ID, &report.UserID, &report.Name, &description, &queryID,
			&report.NaturalLanguageQuery, &report.SQLQuery, &report.ScheduleType,
			&report.ScheduleTime, &scheduleDay, &report.Recipients,
			&report.Format, &report.Locale, &report.IncludeCharts,
			&lastRunAt, &nextRunAt, &report.IsActive, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan scheduled report: %w", err)
		}

		if description.Valid {
			report.Description = &description.String
		}
		if queryID.Valid {
			uid, err := uuid.FromBytes(queryID.Bytes[:])
			if err == nil {
				report.QueryID = &uid
			}
		}
		if scheduleDay.Valid {
			day := int(scheduleDay.Int32)
			report.ScheduleDay = &day
		}
		if lastRunAt.Valid {
			report.LastRunAt = &lastRunAt.Time
		}
		if nextRunAt.Valid {
			report.NextRunAt = &nextRunAt.Time
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// Update updates a scheduled report.
func (r *ScheduledReportRepository) Update(ctx context.Context, report *models.ScheduledReport) error {
	query := `
		UPDATE app.scheduled_reports SET
			name = $2, description = $3, schedule_type = $4, schedule_time = $5, schedule_day = $6,
			recipients = $7, format = $8, include_charts = $9, next_run_at = $10, updated_at = NOW()
		WHERE id = $1 AND user_id = $11
		RETURNING updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		report.ID, report.Name, report.Description, report.ScheduleType,
		report.ScheduleTime, report.ScheduleDay, report.Recipients,
		report.Format, report.IncludeCharts, report.NextRunAt, report.UserID,
	).Scan(&report.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update scheduled report: %w", err)
	}

	return nil
}

// Delete removes a scheduled report.
func (r *ScheduledReportRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM app.scheduled_reports WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete scheduled report: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: scheduled report not found")
	}

	return nil
}

// Toggle enables or disables a scheduled report.
func (r *ScheduledReportRepository) Toggle(ctx context.Context, id, userID uuid.UUID, isActive bool) error {
	query := `UPDATE app.scheduled_reports SET is_active = $3, updated_at = NOW() WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID, isActive)
	if err != nil {
		return fmt.Errorf("warehouse: failed to toggle scheduled report: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: scheduled report not found")
	}

	return nil
}

// GetDue retrieves scheduled reports that are due for execution.
func (r *ScheduledReportRepository) GetDue(ctx context.Context) ([]*models.ScheduledReport, error) {
	query := `
		SELECT id, user_id, name, description, query_id, natural_language_query, sql_query, schedule_type, schedule_time, schedule_day, recipients, format, locale, include_charts, last_run_at, next_run_at, is_active, created_at, updated_at
		FROM app.scheduled_reports
		WHERE is_active = true AND next_run_at <= NOW()
	`

	rows, err := r.db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get due scheduled reports: %w", err)
	}
	defer rows.Close()

	var reports []*models.ScheduledReport
	for rows.Next() {
		report := &models.ScheduledReport{}
		var description pgtype.Text
		var queryID pgtype.UUID
		var scheduleDay pgtype.Int4
		var lastRunAt, nextRunAt pgtype.Timestamptz

		err := rows.Scan(
			&report.ID, &report.UserID, &report.Name, &description, &queryID,
			&report.NaturalLanguageQuery, &report.SQLQuery, &report.ScheduleType,
			&report.ScheduleTime, &scheduleDay, &report.Recipients,
			&report.Format, &report.Locale, &report.IncludeCharts,
			&lastRunAt, &nextRunAt, &report.IsActive, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan scheduled report: %w", err)
		}

		if description.Valid {
			report.Description = &description.String
		}
		if queryID.Valid {
			uid, err := uuid.FromBytes(queryID.Bytes[:])
			if err == nil {
				report.QueryID = &uid
			}
		}
		if scheduleDay.Valid {
			day := int(scheduleDay.Int32)
			report.ScheduleDay = &day
		}
		if lastRunAt.Valid {
			report.LastRunAt = &lastRunAt.Time
		}
		if nextRunAt.Valid {
			report.NextRunAt = &nextRunAt.Time
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// UpdateLastRun updates the last_run_at and next_run_at timestamps.
func (r *ScheduledReportRepository) UpdateLastRun(ctx context.Context, id uuid.UUID, nextRunAt *time.Time) error {
	query := `UPDATE app.scheduled_reports SET last_run_at = NOW(), next_run_at = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.pool.Exec(ctx, query, id, nextRunAt)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update last run: %w", err)
	}

	return nil
}

// CreateRun creates a new report run record.
func (r *ScheduledReportRepository) CreateRun(ctx context.Context, run *models.ScheduledReportRun) error {
	query := `
		INSERT INTO app.scheduled_report_runs (report_id, status)
		VALUES ($1, $2)
		RETURNING id, started_at
	`

	err := r.db.pool.QueryRow(ctx, query, run.ReportID, run.Status).Scan(&run.ID, &run.StartedAt)
	if err != nil {
		return fmt.Errorf("warehouse: failed to create report run: %w", err)
	}

	return nil
}

// UpdateRun updates a report run record.
func (r *ScheduledReportRepository) UpdateRun(ctx context.Context, run *models.ScheduledReportRun) error {
	query := `
		UPDATE app.scheduled_report_runs SET
			status = $2, file_path = $3, file_size_bytes = $4, row_count = $5, error_message = $6, completed_at = $7
		WHERE id = $1
	`

	_, err := r.db.pool.Exec(ctx, query,
		run.ID, run.Status, run.FilePath, run.FileSizeBytes,
		run.RowCount, run.ErrorMessage, run.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update report run: %w", err)
	}

	return nil
}

// GetRuns retrieves run history for a report.
func (r *ScheduledReportRepository) GetRuns(ctx context.Context, reportID uuid.UUID, limit int) ([]*models.ScheduledReportRun, error) {
	query := `
		SELECT id, report_id, status, file_path, file_size_bytes, row_count, error_message, started_at, completed_at
		FROM app.scheduled_report_runs
		WHERE report_id = $1
		ORDER BY started_at DESC
		LIMIT $2
	`

	rows, err := r.db.pool.Query(ctx, query, reportID, limit)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get report runs: %w", err)
	}
	defer rows.Close()

	var runs []*models.ScheduledReportRun
	for rows.Next() {
		run := &models.ScheduledReportRun{}
		var filePath pgtype.Text
		var fileSizeBytes pgtype.Int8
		var rowCount pgtype.Int4
		var errorMessage pgtype.Text
		var completedAt pgtype.Timestamptz

		err := rows.Scan(
			&run.ID, &run.ReportID, &run.Status, &filePath, &fileSizeBytes,
			&rowCount, &errorMessage, &run.StartedAt, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan report run: %w", err)
		}

		if filePath.Valid {
			run.FilePath = &filePath.String
		}
		if fileSizeBytes.Valid {
			size := fileSizeBytes.Int64
			run.FileSizeBytes = &size
		}
		if rowCount.Valid {
			count := int(rowCount.Int32)
			run.RowCount = &count
		}
		if errorMessage.Valid {
			run.ErrorMessage = &errorMessage.String
		}
		if completedAt.Valid {
			run.CompletedAt = &completedAt.Time
		}

		runs = append(runs, run)
	}

	return runs, nil
}
