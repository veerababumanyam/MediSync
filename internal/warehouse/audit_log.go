// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// AuditLogRepository handles database operations for document audit logs.
type AuditLogRepository struct {
	db *Repo
}

// NewAuditLogRepository creates a new audit log repository.
func NewAuditLogRepository(db *Repo) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create inserts a new audit log entry.
func (r *AuditLogRepository) Create(ctx context.Context, entry *models.DocumentAuditLog) error {
	var oldValueJSON, newValueJSON []byte
	var err error

	if entry.OldValue != nil {
		oldValueJSON, err = json.Marshal(entry.OldValue)
		if err != nil {
			return fmt.Errorf("warehouse: failed to marshal old value: %w", err)
		}
	}
	if entry.NewValue != nil {
		newValueJSON, err = json.Marshal(entry.NewValue)
		if err != nil {
			return fmt.Errorf("warehouse: failed to marshal new value: %w", err)
		}
	}

	query := `
		INSERT INTO app.document_audit_log (
			id, tenant_id, document_id, action, actor_id, actor_type,
			field_name, old_value, new_value, notes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.db.pool.Exec(ctx, query,
		entry.ID, entry.TenantID, entry.DocumentID, entry.Action,
		entry.ActorID, entry.ActorType, entry.FieldName,
		oldValueJSON, newValueJSON, entry.Notes, entry.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create audit log: %w", err)
	}

	return nil
}

// GetByDocumentID retrieves all audit log entries for a document.
func (r *AuditLogRepository) GetByDocumentID(ctx context.Context, documentID uuid.UUID) ([]models.DocumentAuditLog, error) {
	query := `
		SELECT id, tenant_id, document_id, action, actor_id, actor_type,
			field_name, old_value, new_value, notes, created_at
		FROM app.document_audit_log
		WHERE document_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get audit logs: %w", err)
	}
	defer rows.Close()

	var entries []models.DocumentAuditLog
	for rows.Next() {
		entry := models.DocumentAuditLog{}
		var oldValueJSON, newValueJSON []byte
		var actorID pgtype.UUID

		err := rows.Scan(
			&entry.ID, &entry.TenantID, &entry.DocumentID, &entry.Action,
			&actorID, &entry.ActorType, &entry.FieldName,
			&oldValueJSON, &newValueJSON, &entry.Notes, &entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan audit log: %w", err)
		}

		if actorID.Valid {
			uid, _ := uuid.FromBytes(actorID.Bytes[:])
			entry.ActorID = uid
		}

		if len(oldValueJSON) > 0 {
			json.Unmarshal(oldValueJSON, &entry.OldValue)
		}
		if len(newValueJSON) > 0 {
			json.Unmarshal(newValueJSON, &entry.NewValue)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetByTenantID retrieves audit log entries for a tenant with filtering.
func (r *AuditLogRepository) GetByTenantID(ctx context.Context, filter *models.AuditLogFilter) (*models.AuditLogListResult, error) {
	// Build base query
	baseQuery := `
		FROM app.document_audit_log
		WHERE tenant_id = $1
	`
	args := []any{filter.TenantID}
	argIdx := 2

	// Add document filter
	if filter.DocumentID != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND document_id = $%d", argIdx)
		args = append(args, filter.DocumentID)
		argIdx++
	}

	// Add actions filter
	if len(filter.Actions) > 0 {
		baseQuery += fmt.Sprintf(" AND action = ANY($%d)", argIdx)
		actions := make([]string, len(filter.Actions))
		for i, a := range filter.Actions {
			actions[i] = string(a)
		}
		args = append(args, actions)
		argIdx++
	}

	// Add actor filter
	if filter.ActorID != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND actor_id = $%d", argIdx)
		args = append(args, filter.ActorID)
		argIdx++
	}

	// Add date range filter
	if filter.DateFrom != nil {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, filter.DateTo)
		argIdx++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to count audit logs: %w", err)
	}

	// Determine pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	totalPages := (total + pageSize - 1) / pageSize

	// Build sort
	sortBy := "created_at"
	if filter.SortBy != "" {
		validSorts := map[string]bool{"created_at": true, "action": true, "actor_id": true}
		if validSorts[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	dataQuery := fmt.Sprintf(`
		SELECT id, tenant_id, document_id, action, actor_id, actor_type,
			field_name, old_value, new_value, notes, created_at
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, baseQuery, sortBy, sortOrder, argIdx, argIdx+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var entries []models.DocumentAuditLog
	for rows.Next() {
		entry := models.DocumentAuditLog{}
		var oldValueJSON, newValueJSON []byte
		var actorID pgtype.UUID

		err := rows.Scan(
			&entry.ID, &entry.TenantID, &entry.DocumentID, &entry.Action,
			&actorID, &entry.ActorType, &entry.FieldName,
			&oldValueJSON, &newValueJSON, &entry.Notes, &entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan audit log: %w", err)
		}

		if actorID.Valid {
			uid, _ := uuid.FromBytes(actorID.Bytes[:])
			entry.ActorID = uid
		}

		if len(oldValueJSON) > 0 {
			json.Unmarshal(oldValueJSON, &entry.OldValue)
		}
		if len(newValueJSON) > 0 {
			json.Unmarshal(newValueJSON, &entry.NewValue)
		}

		entries = append(entries, entry)
	}

	return &models.AuditLogListResult{
		Entries:    entries,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID retrieves a single audit log entry by ID.
func (r *AuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.DocumentAuditLog, error) {
	query := `
		SELECT id, tenant_id, document_id, action, actor_id, actor_type,
			field_name, old_value, new_value, notes, created_at
		FROM app.document_audit_log
		WHERE id = $1
	`

	entry := &models.DocumentAuditLog{}
	var oldValueJSON, newValueJSON []byte
	var actorID pgtype.UUID

	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&entry.ID, &entry.TenantID, &entry.DocumentID, &entry.Action,
		&actorID, &entry.ActorType, &entry.FieldName,
		&oldValueJSON, &newValueJSON, &entry.Notes, &entry.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("warehouse: audit log not found")
		}
		return nil, fmt.Errorf("warehouse: failed to get audit log: %w", err)
	}

	if actorID.Valid {
		uid, _ := uuid.FromBytes(actorID.Bytes[:])
		entry.ActorID = uid
	}

	if len(oldValueJSON) > 0 {
		json.Unmarshal(oldValueJSON, &entry.OldValue)
	}
	if len(newValueJSON) > 0 {
		json.Unmarshal(newValueJSON, &entry.NewValue)
	}

	return entry, nil
}
