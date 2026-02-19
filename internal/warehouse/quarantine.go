// Package warehouse provides quarantine operations for failed ETL records.
//
// Records that fail validation or cannot be processed are quarantined
// for later review and potential reprocessing.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// QuarantineStatus represents the status of a quarantined record.
type QuarantineStatus string

const (
	// QuarantineStatusPending indicates the record is pending review.
	QuarantineStatusPending QuarantineStatus = "pending"
	// QuarantineStatusRetrying indicates the record is being retried.
	QuarantineStatusRetrying QuarantineStatus = "retrying"
	// QuarantineStatusResolved indicates the record was successfully processed.
	QuarantineStatusResolved QuarantineStatus = "resolved"
	// QuarantineStatusIgnored indicates the record was marked to ignore.
	QuarantineStatusIgnored QuarantineStatus = "ignored"
)

// QuarantineRecord represents a failed record in quarantine.
type QuarantineRecord struct {
	RecordID              uuid.UUID              `db:"record_id"`
	BatchID               *uuid.UUID             `db:"batch_id"`
	Source                string                 `db:"source"`
	SourceTable           *string                `db:"source_table"`
	SourceID              *string                `db:"source_id"`
	RawData               []byte                 `db:"raw_data"`      // JSONB
	RawXML                *string                `db:"raw_xml"`
	ErrorReason           string                 `db:"error_reason"`
	ErrorCode             *string                `db:"error_code"`
	ErrorDetails          []byte                 `db:"error_details"`  // JSONB
	ValidationRulesFailed []string               `db:"validation_rules_failed"`
	RetryCount            int                    `db:"retry_count"`
	MaxRetries            int                    `db:"max_retries"`
	LastRetryAt           *time.Time             `db:"last_retry_at"`
	Status                QuarantineStatus       `db:"status"`
	ResolvedBy            *uuid.UUID             `db:"resolved_by"`
	ResolvedAt            *time.Time             `db:"resolved_at"`
	ResolutionNotes       *string                `db:"resolution_notes"`
	CreatedAt             time.Time              `db:"created_at"`
}

// QuarantineFilterOptions provides filtering options for listing quarantined records.
type QuarantineFilterOptions struct {
	Source     string
	Status     *QuarantineStatus
	BatchID    *uuid.UUID
	ErrorCode  *string
	Limit      int
	Offset     int
	OlderThan  *time.Time // For cleanup
}

// QuarantineStats provides statistics about quarantined records.
type QuarantineStats struct {
	TotalRecords      int64           `json:"total_records"`
	PendingRecords    int64           `json:"pending_records"`
	RetryingRecords   int64           `json:"retrying_records"`
	ResolvedRecords   int64           `json:"resolved_records"`
	IgnoredRecords    int64           `json:"ignored_records"`
	BySource          map[string]int64 `json:"by_source"`
	ByErrorCode       map[string]int64 `json:"by_error_code"`
	AverageRetryCount float64         `json:"average_retry_count"`
}

// Quarantine adds a failed record to quarantine.
func (r *Repo) Quarantine(ctx context.Context, record *QuarantineRecord) error {
	if record.MaxRetries == 0 {
		record.MaxRetries = 3 // Default max retries
	}

	query := `
		INSERT INTO app.etl_quarantine (
			batch_id, source, source_table, source_id, raw_data, raw_xml,
			error_reason, error_code, error_details, validation_rules_failed,
			retry_count, max_retries, last_retry_at, status,
			resolved_by, resolved_at, resolution_notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, NOW()
		)
		RETURNING record_id
	`

	err := r.pool.QueryRow(ctx, query,
		record.BatchID, record.Source, record.SourceTable, record.SourceID,
		record.RawData, record.RawXML, record.ErrorReason, record.ErrorCode,
		record.ErrorDetails, record.ValidationRulesFailed,
		record.RetryCount, record.MaxRetries, record.LastRetryAt,
		record.Status, record.ResolvedBy, record.ResolvedAt, record.ResolutionNotes,
	).Scan(&record.RecordID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to quarantine record: %w", err)
	}

	r.logger.Warn("record quarantined",
		slog.String("record_id", record.RecordID.String()),
		slog.String("source", record.Source),
		slog.String("error", record.ErrorReason),
	)

	return nil
}

// QuarantineWithError adds a record to quarantine with error details.
func (r *Repo) QuarantineWithError(ctx context.Context, source string, sourceID string,
	table string, errorReason string, errorCode string, rawData interface{}) error {

	record := &QuarantineRecord{
		Source:      source,
		SourceID:    &sourceID,
		SourceTable: &table,
		ErrorReason: errorReason,
		ErrorCode:   &errorCode,
		Status:      QuarantineStatusPending,
	}

	// Marshal raw data to JSON
	if rawData != nil {
		dataBytes, err := json.Marshal(rawData)
		if err != nil {
			r.logger.Error("failed to marshal quarantine data",
				slog.String("error", err.Error()),
			)
			// Continue without raw data
		} else {
			record.RawData = dataBytes
		}
	}

	return r.Quarantine(ctx, record)
}

// QuarantineBatch adds multiple records to quarantine in a single transaction.
func (r *Repo) QuarantineBatch(ctx context.Context, records []*QuarantineRecord) error {
	if len(records) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO app.etl_quarantine (
			batch_id, source, source_table, source_id, raw_data, raw_xml,
			error_reason, error_code, error_details, validation_rules_failed,
			retry_count, max_retries, last_retry_at, status,
			resolved_by, resolved_at, resolution_notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, NOW()
		)
		RETURNING record_id
	`

	for _, record := range records {
		if record.MaxRetries == 0 {
			record.MaxRetries = 3
		}
		batch.Queue(query,
			record.BatchID, record.Source, record.SourceTable, record.SourceID,
			record.RawData, record.RawXML, record.ErrorReason, record.ErrorCode,
			record.ErrorDetails, record.ValidationRulesFailed,
			record.RetryCount, record.MaxRetries, record.LastRetryAt,
			record.Status, record.ResolvedBy, record.ResolvedAt, record.ResolutionNotes,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := range records {
		err := results.QueryRow().Scan(&records[i].RecordID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to quarantine record %d/%d: %w",
				i+1, len(records), err)
		}
	}

	r.logger.Warn("batch records quarantined",
		slog.Int("count", len(records)),
		slog.String("source", records[0].Source),
	)

	return nil
}

// GetQuarantinedRecord retrieves a specific quarantined record by ID.
func (r *Repo) GetQuarantinedRecord(ctx context.Context, recordID uuid.UUID) (*QuarantineRecord, error) {
	query := `
		SELECT record_id, batch_id, source, source_table, source_id, raw_data, raw_xml,
		       error_reason, error_code, error_details, validation_rules_failed,
		       retry_count, max_retries, last_retry_at, status,
		       resolved_by, resolved_at, resolution_notes, created_at
		FROM app.etl_quarantine
		WHERE record_id = $1
	`

	var record QuarantineRecord
	err := r.pool.QueryRow(ctx, query, recordID).Scan(
		&record.RecordID, &record.BatchID, &record.Source, &record.SourceTable,
		&record.SourceID, &record.RawData, &record.RawXML, &record.ErrorReason,
		&record.ErrorCode, &record.ErrorDetails, &record.ValidationRulesFailed,
		&record.RetryCount, &record.MaxRetries, &record.LastRetryAt, &record.Status,
		&record.ResolvedBy, &record.ResolvedAt, &record.ResolutionNotes, &record.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("warehouse: quarantined record not found: %s", recordID)
	}

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get quarantined record: %w", err)
	}

	return &record, nil
}

// ListQuarantinedRecords retrieves quarantined records with optional filtering.
func (r *Repo) ListQuarantinedRecords(ctx context.Context, opts *QuarantineFilterOptions) ([]*QuarantineRecord, int64, error) {
	if opts == nil {
		opts = &QuarantineFilterOptions{}
	}

	// First get count
	countQuery := "SELECT COUNT(*) FROM app.etl_quarantine WHERE 1=1"
	countArgs := []interface{}{}
	argIdx := 1

	whereClause := ""
	args := []interface{}{}

	if opts.Source != "" {
		whereClause += fmt.Sprintf(" AND source = $%d", argIdx)
		args = append(args, opts.Source)
		countArgs = append(countArgs, opts.Source)
		argIdx++
	}

	if opts.Status != nil {
		whereClause += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *opts.Status)
		countArgs = append(countArgs, *opts.Status)
		argIdx++
	}

	if opts.BatchID != nil {
		whereClause += fmt.Sprintf(" AND batch_id = $%d", argIdx)
		args = append(args, *opts.BatchID)
		countArgs = append(countArgs, *opts.BatchID)
		argIdx++
	}

	if opts.ErrorCode != nil {
		whereClause += fmt.Sprintf(" AND error_code = $%d", argIdx)
		args = append(args, *opts.ErrorCode)
		countArgs = append(countArgs, *opts.ErrorCode)
		argIdx++
	}

	if opts.OlderThan != nil {
		whereClause += fmt.Sprintf(" AND created_at < $%d", argIdx)
		args = append(args, *opts.OlderThan)
		countArgs = append(countArgs, *opts.OlderThan)
		argIdx++
	}

	// Get total count
	var totalCount int64
	countQuery += whereClause
	err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("warehouse: failed to count quarantined records: %w", err)
	}

	// Get records
	query := `
		SELECT record_id, batch_id, source, source_table, source_id, raw_data, raw_xml,
		       error_reason, error_code, error_details, validation_rules_failed,
		       retry_count, max_retries, last_retry_at, status,
		       resolved_by, resolved_at, resolution_notes, created_at
		FROM app.etl_quarantine
		WHERE 1=1
	` + whereClause + `
		ORDER BY created_at DESC
	`

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, opts.Limit)
		argIdx++
	}

	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, opts.Offset)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("warehouse: failed to list quarantined records: %w", err)
	}
	defer rows.Close()

	var records []*QuarantineRecord
	for rows.Next() {
		var record QuarantineRecord
		err := rows.Scan(
			&record.RecordID, &record.BatchID, &record.Source, &record.SourceTable,
			&record.SourceID, &record.RawData, &record.RawXML, &record.ErrorReason,
			&record.ErrorCode, &record.ErrorDetails, &record.ValidationRulesFailed,
			&record.RetryCount, &record.MaxRetries, &record.LastRetryAt, &record.Status,
			&record.ResolvedBy, &record.ResolvedAt, &record.ResolutionNotes, &record.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("warehouse: failed to scan quarantined record row: %w", err)
		}
		records = append(records, &record)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("warehouse: error iterating quarantined records: %w", err)
	}

	return records, totalCount, nil
}

// UpdateQuarantineStatus updates the status of a quarantined record.
func (r *Repo) UpdateQuarantineStatus(ctx context.Context, recordID uuid.UUID,
	status QuarantineStatus, resolvedBy *uuid.UUID, resolutionNotes *string) error {

	query := `
		UPDATE app.etl_quarantine
		SET status = $1,
		    resolved_by = $2,
		    resolved_at = CASE WHEN $2 IS NOT NULL THEN NOW() ELSE resolved_at END,
		    resolution_notes = $3
		WHERE record_id = $4
	`

	result, err := r.pool.Exec(ctx, query, status, resolvedBy, resolutionNotes, recordID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update quarantine status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: quarantined record not found: %s", recordID)
	}

	return nil
}

// IncrementRetry increments the retry count for a quarantined record.
func (r *Repo) IncrementRetry(ctx context.Context, recordID uuid.UUID) error {
	query := `
		UPDATE app.etl_quarantine
		SET retry_count = retry_count + 1,
		    last_retry_at = NOW(),
		    status = CASE WHEN retry_count + 1 >= max_retries THEN 'pending' ELSE 'retrying' END
		WHERE record_id = $1
		RETURNING retry_count, status
	`

	var newRetryCount int
	var newStatus QuarantineStatus
	err := r.pool.QueryRow(ctx, query, recordID).Scan(&newRetryCount, &newStatus)
	if err != nil {
		return fmt.Errorf("warehouse: failed to increment retry: %w", err)
	}

	if newRetryCount >= 3 {
		return fmt.Errorf("warehouse: max retries exceeded for record %s", recordID)
	}

	return nil
}

// GetRetriableRecords returns records that can be retried (haven't exceeded max retries).
func (r *Repo) GetRetriableRecords(ctx context.Context, source string, limit int) ([]*QuarantineRecord, error) {
	query := `
		SELECT record_id, batch_id, source, source_table, source_id, raw_data, raw_xml,
		       error_reason, error_code, error_details, validation_rules_failed,
		       retry_count, max_retries, last_retry_at, status,
		       resolved_by, resolved_at, resolution_notes, created_at
		FROM app.etl_quarantine
		WHERE source = $1
		  AND status IN ('pending', 'retrying')
		  AND retry_count < max_retries
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, source, limit)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get retriable records: %w", err)
	}
	defer rows.Close()

	var records []*QuarantineRecord
	for rows.Next() {
		var record QuarantineRecord
		err := rows.Scan(
			&record.RecordID, &record.BatchID, &record.Source, &record.SourceTable,
			&record.SourceID, &record.RawData, &record.RawXML, &record.ErrorReason,
			&record.ErrorCode, &record.ErrorDetails, &record.ValidationRulesFailed,
			&record.RetryCount, &record.MaxRetries, &record.LastRetryAt, &record.Status,
			&record.ResolvedBy, &record.ResolvedAt, &record.ResolutionNotes, &record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan retriable record row: %w", err)
		}
		records = append(records, &record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating retriable records: %w", err)
	}

	return records, nil
}

// DeleteQuarantineRecord removes a record from quarantine (typically after successful reprocessing).
func (r *Repo) DeleteQuarantineRecord(ctx context.Context, recordID uuid.UUID) error {
	query := `DELETE FROM app.etl_quarantine WHERE record_id = $1`

	result, err := r.pool.Exec(ctx, query, recordID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete quarantined record: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: quarantined record not found: %s", recordID)
	}

	return nil
}

// CleanupOldQuarantine removes old resolved/ignored records.
func (r *Repo) CleanupOldQuarantine(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM app.etl_quarantine
		WHERE status IN ('resolved', 'ignored')
		  AND created_at < NOW() - $1::INTERVAL
	`

	interval := fmt.Sprintf("%f seconds", olderThan.Seconds())
	result, err := r.pool.Exec(ctx, query, interval)
	if err != nil {
		return 0, fmt.Errorf("warehouse: failed to cleanup old quarantine records: %w", err)
	}

	count := result.RowsAffected()
	if count > 0 {
		r.logger.Info("cleaned up old quarantine records",
			slog.Int64("count", count),
		)
	}

	return count, nil
}

// GetQuarantineStats returns statistics about quarantined records.
func (r *Repo) GetQuarantineStats(ctx context.Context, source string) (*QuarantineStats, error) {
	query := `
		SELECT
			COUNT(*) as total_records,
			COUNT(*) FILTER (WHERE status = 'pending') as pending_records,
			COUNT(*) FILTER (WHERE status = 'retrying') as retrying_records,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved_records,
			COUNT(*) FILTER (WHERE status = 'ignored') as ignored_records,
			COALESCE(AVG(retry_count), 0) as average_retry_count
		FROM app.etl_quarantine
		WHERE $1 = '' OR source = $1
	`

	var stats QuarantineStats
	stats.BySource = make(map[string]int64)
	stats.ByErrorCode = make(map[string]int64)

	err := r.pool.QueryRow(ctx, query, source).Scan(
		&stats.TotalRecords,
		&stats.PendingRecords,
		&stats.RetryingRecords,
		&stats.ResolvedRecords,
		&stats.IgnoredRecords,
		&stats.AverageRetryCount,
	)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get quarantine stats: %w", err)
	}

	// Get breakdown by source
	bySourceQuery := `
		SELECT source, COUNT(*) as count
		FROM app.etl_quarantine
		WHERE $1 = '' OR source = $1
		GROUP BY source
	`

	rows, err := r.pool.Query(ctx, bySourceQuery, source)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var src string
			var count int64
			if rows.Scan(&src, &count) == nil {
				stats.BySource[src] = count
			}
		}
	}

	// Get breakdown by error code
	byErrorCodeQuery := `
		SELECT error_code, COUNT(*) as count
		FROM app.etl_quarantine
		WHERE ($1 = '' OR source = $1) AND error_code IS NOT NULL
		GROUP BY error_code
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err = r.pool.Query(ctx, byErrorCodeQuery, source)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var code string
			var count int64
			if rows.Scan(&code, &count) == nil {
				stats.ByErrorCode[code] = count
			}
		}
	}

	return &stats, nil
}

// ResolveQuarantineBatch resolves multiple quarantine records at once.
func (r *Repo) ResolveQuarantineBatch(ctx context.Context, recordIDs []uuid.UUID,
	resolvedBy uuid.UUID, resolutionNotes string) error {

	if len(recordIDs) == 0 {
		return nil
	}

	query := `
		UPDATE app.etl_quarantine
		SET status = 'resolved',
		    resolved_by = $1,
		    resolved_at = NOW(),
		    resolution_notes = $2
		WHERE record_id = ANY($3)
	`

	result, err := r.pool.Exec(ctx, query, resolvedBy, resolutionNotes, recordIDs)
	if err != nil {
		return fmt.Errorf("warehouse: failed to resolve quarantine batch: %w", err)
	}

	r.logger.Info("resolved quarantine batch",
		slog.Int64("count", result.RowsAffected()),
		slog.String("resolved_by", resolvedBy.String()),
	)

	return nil
}
