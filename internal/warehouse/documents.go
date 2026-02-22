// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// DocumentRepository handles database operations for documents.
type DocumentRepository struct {
	db *Repo
}

// NewDocumentRepository creates a new document repository.
func NewDocumentRepository(db *Repo) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// DocumentStats represents document processing statistics.
type DocumentStats struct {
	TotalDocuments    int
	PendingReview     int
	UnderReview       int
	Approved          int
	Rejected          int
	Processing        int
	HighPriorityQueue int
	AverageConfidence float64
}

// Create inserts a new document.
func (r *DocumentRepository) Create(ctx context.Context, doc *models.Document) error {
	query := `
		INSERT INTO app.documents (
			id, tenant_id, uploaded_by, status, document_type, original_filename,
			storage_path, file_size_bytes, file_format, page_count, detected_language,
			upload_id, classification_confidence, overall_confidence, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.pool.Exec(ctx, query,
		doc.ID, doc.TenantID, doc.UploadedBy, doc.Status, doc.DocumentType, doc.OriginalFilename,
		doc.StoragePath, doc.FileSizeBytes, doc.FileFormat, doc.PageCount, doc.DetectedLanguage,
		doc.UploadID, doc.ClassificationConfidence, doc.OverallConfidence, doc.CreatedAt, doc.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create document: %w", err)
	}

	return nil
}

// GetByID retrieves a document by ID.
func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	query := `
		SELECT id, tenant_id, uploaded_by, status, document_type, original_filename,
			storage_path, file_size_bytes, file_format, page_count, detected_language,
			upload_id, processing_started_at, processing_completed_at,
			classification_confidence, overall_confidence, rejection_reason,
			locked_by, locked_at, created_at, updated_at
		FROM app.documents
		WHERE id = $1
	`

	doc := &models.Document{}
	var documentType, storagePath pgtype.Text
	var processingStartedAt, processingCompletedAt, lockedAt pgtype.Timestamptz
	var lockedBy pgtype.UUID
	var uploadID pgtype.UUID

	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&doc.ID, &doc.TenantID, &doc.UploadedBy, &doc.Status, &documentType,
		&doc.OriginalFilename, &storagePath, &doc.FileSizeBytes, &doc.FileFormat,
		&doc.PageCount, &doc.DetectedLanguage, &uploadID,
		&processingStartedAt, &processingCompletedAt,
		&doc.ClassificationConfidence, &doc.OverallConfidence, &doc.RejectionReason,
		&lockedBy, &lockedAt, &doc.CreatedAt, &doc.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("warehouse: document not found")
		}
		return nil, fmt.Errorf("warehouse: failed to get document: %w", err)
	}

	if documentType.Valid {
		doc.DocumentType = models.DocumentType(documentType.String)
	}
	if storagePath.Valid {
		doc.StoragePath = storagePath.String
	}
	if uploadID.Valid {
		uid, _ := uuid.FromBytes(uploadID.Bytes[:])
		doc.UploadID = uid
	}
	if processingStartedAt.Valid {
		doc.ProcessingStartedAt = &processingStartedAt.Time
	}
	if processingCompletedAt.Valid {
		doc.ProcessingCompletedAt = &processingCompletedAt.Time
	}
	if lockedBy.Valid {
		uid, _ := uuid.FromBytes(lockedBy.Bytes[:])
		doc.LockedBy = uid
	}
	if lockedAt.Valid {
		doc.LockedAt = &lockedAt.Time
	}

	return doc, nil
}

// Update updates a document.
func (r *DocumentRepository) Update(ctx context.Context, doc *models.Document) error {
	query := `
		UPDATE app.documents SET
			status = $2, document_type = $3, processing_started_at = $4,
			processing_completed_at = $5, classification_confidence = $6,
			overall_confidence = $7, rejection_reason = $8, locked_by = $9,
			locked_at = $10, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.pool.Exec(ctx, query,
		doc.ID, doc.Status, doc.DocumentType, doc.ProcessingStartedAt,
		doc.ProcessingCompletedAt, doc.ClassificationConfidence,
		doc.OverallConfidence, doc.RejectionReason, doc.LockedBy, doc.LockedAt,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update document: %w", err)
	}

	return nil
}

// Delete removes a document.
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM app.documents WHERE id = $1`

	result, err := r.db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete document: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: document not found")
	}

	return nil
}

// List retrieves documents based on filter criteria.
func (r *DocumentRepository) List(ctx context.Context, filter *models.DocumentListFilter) (*models.DocumentListResult, error) {
	// Build base query
	baseQuery := `
		FROM app.documents
		WHERE tenant_id = $1
	`
	args := []any{filter.TenantID}
	argIdx := 2

	// Add status filter
	if len(filter.Statuses) > 0 {
		baseQuery += fmt.Sprintf(" AND status = ANY($%d)", argIdx)
		statuses := make([]string, len(filter.Statuses))
		for i, s := range filter.Statuses {
			statuses[i] = string(s)
		}
		args = append(args, statuses)
		argIdx++
	}

	// Add type filter
	if len(filter.Types) > 0 {
		baseQuery += fmt.Sprintf(" AND document_type = ANY($%d)", argIdx)
		types := make([]string, len(filter.Types))
		for i, t := range filter.Types {
			types[i] = string(t)
		}
		args = append(args, types)
		argIdx++
	}

	// Add uploaded_by filter
	if filter.UploadedBy != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND uploaded_by = $%d", argIdx)
		args = append(args, filter.UploadedBy)
		argIdx++
	}

	// Add upload_id filter
	if filter.UploadID != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND upload_id = $%d", argIdx)
		args = append(args, filter.UploadID)
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

	// Add search filter
	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND (original_filename ILIKE $%d OR document_type ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to count documents: %w", err)
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
		validSorts := map[string]bool{"created_at": true, "updated_at": true, "status": true, "original_filename": true}
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
		SELECT id, tenant_id, uploaded_by, status, document_type, original_filename,
			storage_path, file_size_bytes, file_format, page_count, detected_language,
			upload_id, processing_started_at, processing_completed_at,
			classification_confidence, overall_confidence, rejection_reason,
			locked_by, locked_at, created_at, updated_at
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, baseQuery, sortBy, sortOrder, argIdx, argIdx+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to list documents: %w", err)
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		doc := models.Document{}
		var documentType, storagePath pgtype.Text
		var processingStartedAt, processingCompletedAt, lockedAt pgtype.Timestamptz
		var lockedBy, uploadID pgtype.UUID

		err := rows.Scan(
			&doc.ID, &doc.TenantID, &doc.UploadedBy, &doc.Status, &documentType,
			&doc.OriginalFilename, &storagePath, &doc.FileSizeBytes, &doc.FileFormat,
			&doc.PageCount, &doc.DetectedLanguage, &uploadID,
			&processingStartedAt, &processingCompletedAt,
			&doc.ClassificationConfidence, &doc.OverallConfidence, &doc.RejectionReason,
			&lockedBy, &lockedAt, &doc.CreatedAt, &doc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan document: %w", err)
		}

		if documentType.Valid {
			doc.DocumentType = models.DocumentType(documentType.String)
		}
		if storagePath.Valid {
			doc.StoragePath = storagePath.String
		}
		if uploadID.Valid {
			uid, _ := uuid.FromBytes(uploadID.Bytes[:])
			doc.UploadID = uid
		}
		if processingStartedAt.Valid {
			doc.ProcessingStartedAt = &processingStartedAt.Time
		}
		if processingCompletedAt.Valid {
			doc.ProcessingCompletedAt = &processingCompletedAt.Time
		}
		if lockedBy.Valid {
			uid, _ := uuid.FromBytes(lockedBy.Bytes[:])
			doc.LockedBy = uid
		}
		if lockedAt.Valid {
			doc.LockedAt = &lockedAt.Time
		}

		documents = append(documents, doc)
	}

	return &models.DocumentListResult{
		Documents:  documents,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// AcquireLock locks a document for review.
func (r *DocumentRepository) AcquireLock(ctx context.Context, documentID, userID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE app.documents SET
			locked_by = $2, locked_at = $3, status = 'under_review', updated_at = NOW()
		WHERE id = $1 AND (locked_by IS NULL OR locked_by = $2)
	`

	result, err := r.db.pool.Exec(ctx, query, documentID, userID, now)
	if err != nil {
		return fmt.Errorf("warehouse: failed to acquire lock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: document is locked by another user")
	}

	return nil
}

// ReleaseLock releases a document lock.
func (r *DocumentRepository) ReleaseLock(ctx context.Context, documentID uuid.UUID) error {
	query := `
		UPDATE app.documents SET
			locked_by = NULL, locked_at = NULL, status = 'ready_for_review', updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.pool.Exec(ctx, query, documentID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to release lock: %w", err)
	}

	return nil
}

// GetFields retrieves all extracted fields for a document.
func (r *DocumentRepository) GetFields(ctx context.Context, documentID uuid.UUID) ([]models.ExtractedField, error) {
	query := `
		SELECT id, document_id, page_number, field_name, field_type, extracted_value,
			confidence_score, bounding_box, is_handwritten, verification_status,
			verified_by, verified_at, original_value, created_at, updated_at
		FROM app.extracted_fields
		WHERE document_id = $1
		ORDER BY page_number, field_name
	`

	rows, err := r.db.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []models.ExtractedField
	for rows.Next() {
		field := models.ExtractedField{}
		var boundingBox []byte
		var verifiedBy pgtype.UUID
		var verifiedAt pgtype.Timestamptz

		err := rows.Scan(
			&field.ID, &field.DocumentID, &field.PageNumber, &field.FieldName,
			&field.FieldType, &field.ExtractedValue, &field.ConfidenceScore,
			&boundingBox, &field.IsHandwritten, &field.VerificationStatus,
			&verifiedBy, &verifiedAt, &field.OriginalValue,
			&field.CreatedAt, &field.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan field: %w", err)
		}

		if len(boundingBox) > 0 {
			var bb models.BoundingBox
			if err := json.Unmarshal(boundingBox, &bb); err == nil {
				field.BoundingBox = &bb
			}
		}
		if verifiedBy.Valid {
			uid, _ := uuid.FromBytes(verifiedBy.Bytes[:])
			field.VerifiedBy = uid
		}
		if verifiedAt.Valid {
			field.VerifiedAt = &verifiedAt.Time
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// GetFieldByID retrieves a single field by ID.
func (r *DocumentRepository) GetFieldByID(ctx context.Context, fieldID uuid.UUID) (*models.ExtractedField, error) {
	query := `
		SELECT id, document_id, page_number, field_name, field_type, extracted_value,
			confidence_score, bounding_box, is_handwritten, verification_status,
			verified_by, verified_at, original_value, created_at, updated_at
		FROM app.extracted_fields
		WHERE id = $1
	`

	field := &models.ExtractedField{}
	var boundingBox []byte
	var verifiedBy pgtype.UUID
	var verifiedAt pgtype.Timestamptz

	err := r.db.pool.QueryRow(ctx, query, fieldID).Scan(
		&field.ID, &field.DocumentID, &field.PageNumber, &field.FieldName,
		&field.FieldType, &field.ExtractedValue, &field.ConfidenceScore,
		&boundingBox, &field.IsHandwritten, &field.VerificationStatus,
		&verifiedBy, &verifiedAt, &field.OriginalValue,
		&field.CreatedAt, &field.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("warehouse: field not found")
		}
		return nil, fmt.Errorf("warehouse: failed to get field: %w", err)
	}

	if len(boundingBox) > 0 {
		var bb models.BoundingBox
		if err := json.Unmarshal(boundingBox, &bb); err == nil {
			field.BoundingBox = &bb
		}
	}
	if verifiedBy.Valid {
		uid, _ := uuid.FromBytes(verifiedBy.Bytes[:])
		field.VerifiedBy = uid
	}
	if verifiedAt.Valid {
		field.VerifiedAt = &verifiedAt.Time
	}

	return field, nil
}

// CreateField creates a new extracted field.
func (r *DocumentRepository) CreateField(ctx context.Context, field *models.ExtractedField) error {
	var boundingBoxJSON []byte
	var err error
	if field.BoundingBox != nil {
		boundingBoxJSON, err = json.Marshal(field.BoundingBox)
		if err != nil {
			return fmt.Errorf("warehouse: failed to marshal bounding box: %w", err)
		}
	}

	query := `
		INSERT INTO app.extracted_fields (
			id, document_id, page_number, field_name, field_type, extracted_value,
			confidence_score, bounding_box, is_handwritten, verification_status,
			verified_by, verified_at, original_value, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = r.db.pool.Exec(ctx, query,
		field.ID, field.DocumentID, field.PageNumber, field.FieldName,
		field.FieldType, field.ExtractedValue, field.ConfidenceScore,
		boundingBoxJSON, field.IsHandwritten, field.VerificationStatus,
		field.VerifiedBy, field.VerifiedAt, field.OriginalValue,
		field.CreatedAt, field.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create field: %w", err)
	}

	return nil
}

// UpdateField updates an extracted field.
func (r *DocumentRepository) UpdateField(ctx context.Context, field *models.ExtractedField) error {
	query := `
		UPDATE app.extracted_fields SET
			extracted_value = $2, verification_status = $3, verified_by = $4,
			verified_at = $5, original_value = $6, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.pool.Exec(ctx, query,
		field.ID, field.ExtractedValue, field.VerificationStatus,
		field.VerifiedBy, field.VerifiedAt, field.OriginalValue,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update field: %w", err)
	}

	return nil
}

// GetLineItems retrieves line items for a document.
func (r *DocumentRepository) GetLineItems(ctx context.Context, documentID uuid.UUID) ([]models.LineItem, error) {
	query := `
		SELECT id, document_id, extracted_field_id, line_number, description,
			quantity, unit_price, amount, tax_rate, transaction_date,
			reference, debit_amount, credit_amount, balance, created_at
		FROM app.line_items
		WHERE document_id = $1
		ORDER BY line_number
	`

	rows, err := r.db.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get line items: %w", err)
	}
	defer rows.Close()

	var items []models.LineItem
	for rows.Next() {
		item := models.LineItem{}
		var transactionDate pgtype.Timestamptz
		var extractedFieldID pgtype.UUID

		err := rows.Scan(
			&item.ID, &item.DocumentID, &extractedFieldID, &item.LineNumber,
			&item.Description, &item.Quantity, &item.UnitPrice, &item.Amount,
			&item.TaxRate, &transactionDate, &item.Reference,
			&item.DebitAmount, &item.CreditAmount, &item.Balance, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan line item: %w", err)
		}

		if extractedFieldID.Valid {
			uid, _ := uuid.FromBytes(extractedFieldID.Bytes[:])
			item.ExtractedFieldID = uid
		}
		if transactionDate.Valid {
			item.TransactionDate = &transactionDate.Time
		}

		items = append(items, item)
	}

	return items, nil
}

// CreateLineItem creates a new line item.
func (r *DocumentRepository) CreateLineItem(ctx context.Context, item *models.LineItem) error {
	query := `
		INSERT INTO app.line_items (
			id, document_id, extracted_field_id, line_number, description,
			quantity, unit_price, amount, tax_rate, transaction_date,
			reference, debit_amount, credit_amount, balance, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.pool.Exec(ctx, query,
		item.ID, item.DocumentID, item.ExtractedFieldID, item.LineNumber,
		item.Description, item.Quantity, item.UnitPrice, item.Amount,
		item.TaxRate, item.TransactionDate, item.Reference,
		item.DebitAmount, item.CreditAmount, item.Balance, item.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create line item: %w", err)
	}

	return nil
}

// GetStats retrieves document processing statistics.
func (r *DocumentRepository) GetStats(ctx context.Context, tenantID uuid.UUID) (*DocumentStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'ready_for_review') as pending_review,
			COUNT(*) FILTER (WHERE status = 'under_review') as under_review,
			COUNT(*) FILTER (WHERE status = 'approved') as approved,
			COUNT(*) FILTER (WHERE status = 'rejected') as rejected,
			COUNT(*) FILTER (WHERE status IN ('uploading', 'classifying', 'extracting')) as processing,
			COALESCE(AVG(overall_confidence), 0) as avg_confidence
		FROM app.documents
		WHERE tenant_id = $1
	`

	stats := &DocumentStats{}
	err := r.db.pool.QueryRow(ctx, query, tenantID).Scan(
		&stats.TotalDocuments, &stats.PendingReview, &stats.UnderReview,
		&stats.Approved, &stats.Rejected, &stats.Processing, &stats.AverageConfidence,
	)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get stats: %w", err)
	}

	// Get high priority queue count
	highPriorityQuery := `
		SELECT COUNT(DISTINCT document_id)
		FROM app.extracted_fields
		WHERE verification_status = 'high_priority'
		AND document_id IN (SELECT id FROM app.documents WHERE tenant_id = $1 AND status IN ('ready_for_review', 'under_review'))
	`
	err = r.db.pool.QueryRow(ctx, highPriorityQuery, tenantID).Scan(&stats.HighPriorityQueue)
	if err != nil {
		stats.HighPriorityQueue = 0
	}

	return stats, nil
}

// UpdateStatus updates the document status.
func (r *DocumentRepository) UpdateStatus(ctx context.Context, documentID uuid.UUID, status models.DocumentStatus) error {
	query := `UPDATE app.documents SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.pool.Exec(ctx, query, documentID, status)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update status: %w", err)
	}
	return nil
}

// SetClassification sets the document type and classification confidence.
func (r *DocumentRepository) SetClassification(ctx context.Context, documentID uuid.UUID, docType models.DocumentType, confidence float64) error {
	query := `
		UPDATE app.documents SET
			document_type = $2, classification_confidence = $3, status = 'extracting', updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query, documentID, docType, confidence)
	if err != nil {
		return fmt.Errorf("warehouse: failed to set classification: %w", err)
	}
	return nil
}

// SetExtractionComplete marks extraction as complete.
func (r *DocumentRepository) SetExtractionComplete(ctx context.Context, documentID uuid.UUID, confidence float64) error {
	now := time.Now()
	query := `
		UPDATE app.documents SET
			status = 'ready_for_review', overall_confidence = $2,
			processing_completed_at = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.pool.Exec(ctx, query, documentID, confidence, now)
	if err != nil {
		return fmt.Errorf("warehouse: failed to mark extraction complete: %w", err)
	}
	return nil
}
