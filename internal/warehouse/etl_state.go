// Package warehouse provides ETL state management for incremental sync.
//
// This file handles cursor tracking for incremental data synchronization
// from source systems (Tally, HIMS) to the data warehouse.
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SyncStatus represents the current status of a sync job.
type SyncStatus string

const (
	// SyncStatusIdle indicates the sync is not currently running.
	SyncStatusIdle SyncStatus = "idle"
	// SyncStatusRunning indicates the sync is currently in progress.
	SyncStatusRunning SyncStatus = "running"
	// SyncStatusCompleted indicates the sync completed successfully.
	SyncStatusCompleted SyncStatus = "completed"
	// SyncStatusFailed indicates the sync failed with an error.
	SyncStatusFailed SyncStatus = "failed"
)

// CursorType represents the type of cursor being used for incremental sync.
type CursorType string

const (
	// CursorTypeAlterID is used for Tally LastAlterID-based incremental sync.
	CursorTypeAlterID CursorType = "alter_id"
	// CursorTypeTimestamp is used for HIMS modified_since-based incremental sync.
	CursorTypeTimestamp CursorType = "timestamp"
	// CursorTypeOffset is used for offset-based pagination.
	CursorTypeOffset CursorType = "offset"
)

// ETLState represents the current state of an ETL sync job for an entity.
type ETLState struct {
	StateID         uuid.UUID  `db:"state_id"`
	Source          string     `db:"source"`
	Entity          string     `db:"entity"`
	LastSyncAt      *time.Time `db:"last_sync_at"`
	LastAlterID     *string    `db:"last_alter_id"`      // For Tally
	LastModifiedAt  *time.Time `db:"last_modified_at"`   // For HIMS
	CursorValue     *string    `db:"cursor_value"`       // Generic cursor
	CursorType      *string    `db:"cursor_type"`
	RecordsSynced   int        `db:"records_synced"`
	SyncStatus      SyncStatus `db:"sync_status"`
	ErrorMessage    *string    `db:"error_message"`
	Metadata        []byte     `db:"metadata"` // JSONB stored as bytes
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

// GetETLState retrieves the current ETL state for a source and entity.
func (r *Repo) GetETLState(ctx context.Context, source string, entity string) (*ETLState, error) {
	query := `
		SELECT state_id, source, entity, last_sync_at, last_alter_id,
		       last_modified_at, cursor_value, cursor_type, records_synced,
		       sync_status, error_message, metadata, created_at, updated_at
		FROM app.etl_state
		WHERE source = $1 AND entity = $2
	`

	var state ETLState
	err := r.pool.QueryRow(ctx, query, source, entity).Scan(
		&state.StateID, &state.Source, &state.Entity, &state.LastSyncAt,
		&state.LastAlterID, &state.LastModifiedAt, &state.CursorValue,
		&state.CursorType, &state.RecordsSynced, &state.SyncStatus,
		&state.ErrorMessage, &state.Metadata, &state.CreatedAt, &state.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		// No state exists yet, return a new empty state
		return &ETLState{
			StateID:    uuid.New(),
			Source:     source,
			Entity:     entity,
			SyncStatus: SyncStatusIdle,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get ETL state for %s/%s: %w", source, entity, err)
	}

	return &state, nil
}

// ListETLStates retrieves all ETL states, optionally filtered by source or status.
func (r *Repo) ListETLStates(ctx context.Context, source string, status *SyncStatus) ([]*ETLState, error) {
	query := `
		SELECT state_id, source, entity, last_sync_at, last_alter_id,
		       last_modified_at, cursor_value, cursor_type, records_synced,
		       sync_status, error_message, metadata, created_at, updated_at
		FROM app.etl_state
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if source != "" {
		query += fmt.Sprintf(" AND source = $%d", argIdx)
		args = append(args, source)
		argIdx++
	}

	if status != nil {
		query += fmt.Sprintf(" AND sync_status = $%d", argIdx)
		args = append(args, *status)
		argIdx++
	}

	query += " ORDER BY source, entity"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to list ETL states: %w", err)
	}
	defer rows.Close()

	var states []*ETLState
	for rows.Next() {
		var state ETLState
		err := rows.Scan(
			&state.StateID, &state.Source, &state.Entity, &state.LastSyncAt,
			&state.LastAlterID, &state.LastModifiedAt, &state.CursorValue,
			&state.CursorType, &state.RecordsSynced, &state.SyncStatus,
			&state.ErrorMessage, &state.Metadata, &state.CreatedAt, &state.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan ETL state row: %w", err)
		}
		states = append(states, &state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating ETL state rows: %w", err)
	}

	return states, nil
}

// UpdateETLState updates the ETL state for a sync job.
func (r *Repo) UpdateETLState(ctx context.Context, state *ETLState) error {
	query := `
		INSERT INTO app.etl_state (
			source, entity, last_sync_at, last_alter_id, last_modified_at,
			cursor_value, cursor_type, records_synced, sync_status,
			error_message, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)
		ON CONFLICT (source, entity)
		DO UPDATE SET
			last_sync_at = EXCLUDED.last_sync_at,
			last_alter_id = EXCLUDED.last_alter_id,
			last_modified_at = EXCLUDED.last_modified_at,
			cursor_value = EXCLUDED.cursor_value,
			cursor_type = EXCLUDED.cursor_type,
			records_synced = EXCLUDED.records_synced,
			sync_status = EXCLUDED.sync_status,
			error_message = EXCLUDED.error_message,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
		RETURNING state_id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		state.Source, state.Entity, state.LastSyncAt, state.LastAlterID,
		state.LastModifiedAt, state.CursorValue, state.CursorType,
		state.RecordsSynced, state.SyncStatus, state.ErrorMessage, state.Metadata,
	).Scan(&state.StateID, &state.CreatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update ETL state for %s/%s: %w", state.Source, state.Entity, err)
	}

	return nil
}

// SetSyncRunning sets the sync status to running for a source/entity.
func (r *Repo) SetSyncRunning(ctx context.Context, source string, entity string) error {
	query := `
		INSERT INTO app.etl_state (
			source, entity, sync_status, records_synced, created_at, updated_at
		) VALUES ($1, $2, $3, 0, NOW(), NOW())
		ON CONFLICT (source, entity)
		DO UPDATE SET
			sync_status = $3,
			error_message = NULL,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query, source, entity, SyncStatusRunning)
	if err != nil {
		return fmt.Errorf("warehouse: failed to set sync running for %s/%s: %w", source, entity, err)
	}

	return nil
}

// SetSyncCompleted sets the sync status to completed and updates the cursor.
func (r *Repo) SetSyncCompleted(ctx context.Context, source string, entity string,
	recordsSynced int, lastSyncAt time.Time, cursorValue *string, cursorType *string) error {

	query := `
		UPDATE app.etl_state
		SET sync_status = $1,
		    records_synced = records_synced + $2,
		    last_sync_at = $3,
		    cursor_value = $4,
		    cursor_type = $5,
		    error_message = NULL,
		    updated_at = NOW()
		WHERE source = $6 AND entity = $7
	`

	result, err := r.pool.Exec(ctx, query, SyncStatusCompleted, recordsSynced, lastSyncAt,
		cursorValue, cursorType, source, entity)
	if err != nil {
		return fmt.Errorf("warehouse: failed to set sync completed for %s/%s: %w", source, entity, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: no ETL state found for %s/%s", source, entity)
	}

	return nil
}

// SetSyncFailed sets the sync status to failed with an error message.
func (r *Repo) SetSyncFailed(ctx context.Context, source string, entity string, errMsg string) error {
	query := `
		UPDATE app.etl_state
		SET sync_status = $1,
		    error_message = $2,
		    updated_at = NOW()
		WHERE source = $3 AND entity = $4
	`

	result, err := r.pool.Exec(ctx, query, SyncStatusFailed, errMsg, source, entity)
	if err != nil {
		return fmt.Errorf("warehouse: failed to set sync failed for %s/%s: %w", source, entity, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: no ETL state found for %s/%s", source, entity)
	}

	return nil
}

// UpdateTallyAlterID updates the LastAlterID cursor for Tally sync.
func (r *Repo) UpdateTallyAlterID(ctx context.Context, entity string, alterID string, recordsSynced int) error {
	cursorType := string(CursorTypeAlterID)
	return r.SetSyncCompleted(ctx, SourceTally.String(), entity, recordsSynced, time.Now(), &alterID, &cursorType)
}

// UpdateHIMSCursor updates the modified_since cursor for HIMS sync.
func (r *Repo) UpdateHIMSCursor(ctx context.Context, entity string, lastModified time.Time, recordsSynced int) error {
	lastModStr := lastModified.Format(time.RFC3339Nano)
	cursorType := string(CursorTypeTimestamp)
	return r.SetSyncCompleted(ctx, SourceHIMS.String(), entity, recordsSynced, time.Now(), &lastModStr, &cursorType)
}

// GetTallyAlterID retrieves the current LastAlterID for a Tally entity.
func (r *Repo) GetTallyAlterID(ctx context.Context, entity string) (string, error) {
	state, err := r.GetETLState(ctx, SourceTally.String(), entity)
	if err != nil {
		return "", err
	}

	if state.LastAlterID != nil {
		return *state.LastAlterID, nil
	}

	return "", nil
}

// GetHIMSCursor retrieves the current modified_since cursor for a HIMS entity.
func (r *Repo) GetHIMSCursor(ctx context.Context, entity string) (*time.Time, error) {
	state, err := r.GetETLState(ctx, SourceHIMS.String(), entity)
	if err != nil {
		return nil, err
	}

	return state.LastModifiedAt, nil
}

// ResetETLState resets the sync state for a source/entity (for full re-sync).
func (r *Repo) ResetETLState(ctx context.Context, source string, entity string) error {
	query := `
		UPDATE app.etl_state
		SET last_alter_id = NULL,
		    last_modified_at = NULL,
		    cursor_value = NULL,
		    records_synced = 0,
		    sync_status = $1,
		    error_message = NULL,
		    updated_at = NOW()
		WHERE source = $2 AND entity = $3
	`

	result, err := r.pool.Exec(ctx, query, SyncStatusIdle, source, entity)
	if err != nil {
		return fmt.Errorf("warehouse: failed to reset ETL state for %s/%s: %w", source, entity, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: no ETL state found for %s/%s", source, entity)
	}

	r.logger.Info("reset ETL state",
		slog.String("source", source),
		slog.String("entity", entity),
	)

	return nil
}

// IsSyncRunning checks if a sync is currently running for a source/entity.
func (r *Repo) IsSyncRunning(ctx context.Context, source string, entity string) (bool, error) {
	state, err := r.GetETLState(ctx, source, entity)
	if err != nil {
		return false, err
	}

	// If state was never created (StateID is nil UUID), it's not running
	if state.StateID == uuid.Nil {
		return false, nil
	}

	return state.SyncStatus == SyncStatusRunning, nil
}

// GetStaleSyncs returns syncs that have been in "running" state for too long.
// A sync is considered stale if it's been running for more than the specified duration.
func (r *Repo) GetStaleSyncs(ctx context.Context, staleDuration time.Duration) ([]*ETLState, error) {
	query := `
		SELECT state_id, source, entity, last_sync_at, last_alter_id,
		       last_modified_at, cursor_value, cursor_type, records_synced,
		       sync_status, error_message, metadata, created_at, updated_at
		FROM app.etl_state
		WHERE sync_status = $1
		  AND updated_at < NOW() - $2::INTERVAL
		ORDER BY updated_at ASC
	`

	staleInterval := fmt.Sprintf("%f seconds", staleDuration.Seconds())
	rows, err := r.pool.Query(ctx, query, SyncStatusRunning, staleInterval)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get stale syncs: %w", err)
	}
	defer rows.Close()

	var states []*ETLState
	for rows.Next() {
		var state ETLState
		err := rows.Scan(
			&state.StateID, &state.Source, &state.Entity, &state.LastSyncAt,
			&state.LastAlterID, &state.LastModifiedAt, &state.CursorValue,
			&state.CursorType, &state.RecordsSynced, &state.SyncStatus,
			&state.ErrorMessage, &state.Metadata, &state.CreatedAt, &state.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan stale sync row: %w", err)
		}
		states = append(states, &state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating stale sync rows: %w", err)
	}

	return states, nil
}

// CleanupOldStates removes ETL state entries for entities that haven't been synced
// in a long time. Returns the number of entries removed.
func (r *Repo) CleanupOldStates(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM app.etl_state
		WHERE (source, entity) NOT IN (
			SELECT source, entity
			FROM app.etl_state
			WHERE updated_at > NOW() - $1::INTERVAL
		)
	`

	interval := fmt.Sprintf("%f seconds", olderThan.Seconds())
	result, err := r.pool.Exec(ctx, query, interval)
	if err != nil {
		return 0, fmt.Errorf("warehouse: failed to cleanup old ETL states: %w", err)
	}

	count := result.RowsAffected()
	if count > 0 {
		r.logger.Info("cleaned up old ETL states",
			slog.Int64("count", count),
		)
	}

	return count, nil
}

// WithETLLock executes a function while holding an ETL sync lock for a source/entity.
// If a sync is already running, it returns ErrSyncInProgress.
func (r *Repo) WithETLLock(ctx context.Context, source string, entity string, fn func() error) error {
	running, err := r.IsSyncRunning(ctx, source, entity)
	if err != nil {
		return err
	}

	if running {
		return fmt.Errorf("sync already in progress for %s/%s", source, entity)
	}

	// Set status to running
	if err := r.SetSyncRunning(ctx, source, entity); err != nil {
		return err
	}

	// Execute the function
	execErr := fn()

	// Update final status
	if execErr != nil {
		_ = r.SetSyncFailed(ctx, source, entity, execErr.Error())
		return execErr
	}

	// Success - let caller set the completed state with cursor
	return nil
}

// GetSyncStats returns summary statistics for all sync jobs.
type SyncStats struct {
	TotalEntities    int            `json:"total_entities"`
	RunningSyncs     int            `json:"running_syncs"`
	CompletedSyncs   int            `json:"completed_syncs"`
	FailedSyncs      int            `json:"failed_syncs"`
	IdleSyncs        int            `json:"idle_syncs"`
	TotalRecords     int64          `json:"total_records"`
	LastSyncAt       *time.Time     `json:"last_sync_at"`
	PendingTime      *time.Duration `json:"pending_time"` // Time since last completed sync
}

// GetSyncStats retrieves overall sync statistics.
func (r *Repo) GetSyncStats(ctx context.Context) (*SyncStats, error) {
	query := `
		SELECT
			COUNT(*) as total_entities,
			COUNT(*) FILTER (WHERE sync_status = 'running') as running_syncs,
			COUNT(*) FILTER (WHERE sync_status = 'completed') as completed_syncs,
			COUNT(*) FILTER (WHERE sync_status = 'failed') as failed_syncs,
			COUNT(*) FILTER (WHERE sync_status = 'idle') as idle_syncs,
			COALESCE(SUM(records_synced), 0) as total_records,
			MAX(last_sync_at) as last_sync_at
		FROM app.etl_state
	`

	var stats SyncStats
	err := r.pool.QueryRow(ctx, query).Scan(
		&stats.TotalEntities,
		&stats.RunningSyncs,
		&stats.CompletedSyncs,
		&stats.FailedSyncs,
		&stats.IdleSyncs,
		&stats.TotalRecords,
		&stats.LastSyncAt,
	)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get sync stats: %w", err)
	}

	return &stats, nil
}
