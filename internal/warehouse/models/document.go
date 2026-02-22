// Package models provides data structures for the MediSync warehouse.
package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DocumentStatus represents the processing status of a document.
type DocumentStatus string

const (
	DocumentStatusUploading       DocumentStatus = "uploading"
	DocumentStatusUploaded        DocumentStatus = "uploaded"
	DocumentStatusClassifying     DocumentStatus = "classifying"
	DocumentStatusExtracting      DocumentStatus = "extracting"
	DocumentStatusReadyForReview  DocumentStatus = "ready_for_review"
	DocumentStatusUnderReview     DocumentStatus = "under_review"
	DocumentStatusReviewed        DocumentStatus = "reviewed"
	DocumentStatusApproved        DocumentStatus = "approved"
	DocumentStatusRejected        DocumentStatus = "rejected"
	DocumentStatusFailed          DocumentStatus = "failed"
)

// IsValid checks if the document status is valid.
func (s DocumentStatus) IsValid() bool {
	switch s {
	case DocumentStatusUploading, DocumentStatusUploaded, DocumentStatusClassifying,
		DocumentStatusExtracting, DocumentStatusReadyForReview, DocumentStatusUnderReview,
		DocumentStatusReviewed, DocumentStatusApproved, DocumentStatusRejected, DocumentStatusFailed:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the status can transition to the target status.
func (s DocumentStatus) CanTransitionTo(target DocumentStatus) bool {
	transitions := map[DocumentStatus][]DocumentStatus{
		DocumentStatusUploading:      {DocumentStatusUploaded, DocumentStatusFailed},
		DocumentStatusUploaded:       {DocumentStatusClassifying, DocumentStatusFailed},
		DocumentStatusClassifying:    {DocumentStatusExtracting, DocumentStatusFailed},
		DocumentStatusExtracting:     {DocumentStatusReadyForReview, DocumentStatusFailed},
		DocumentStatusReadyForReview: {DocumentStatusUnderReview, DocumentStatusRejected, DocumentStatusFailed},
		DocumentStatusUnderReview:    {DocumentStatusReviewed, DocumentStatusRejected, DocumentStatusReadyForReview},
		DocumentStatusReviewed:       {DocumentStatusApproved, DocumentStatusRejected},
		DocumentStatusApproved:       {},
		DocumentStatusRejected:       {DocumentStatusUploaded}, // Can reprocess
		DocumentStatusFailed:         {DocumentStatusUploaded}, // Can retry
	}

	allowed, exists := transitions[s]
	if !exists {
		return false
	}

	for _, t := range allowed {
		if t == target {
			return true
		}
	}
	return false
}

// DocumentType represents the type of document.
type DocumentType string

const (
	DocumentTypeInvoice       DocumentType = "invoice"
	DocumentTypeReceipt       DocumentType = "receipt"
	DocumentTypeBankStatement DocumentType = "bank_statement"
	DocumentTypeExpenseReport DocumentType = "expense_report"
	DocumentTypeCreditNote    DocumentType = "credit_note"
	DocumentTypeDebitNote     DocumentType = "debit_note"
	DocumentTypeOther         DocumentType = "other"
)

// IsValid checks if the document type is valid.
func (t DocumentType) IsValid() bool {
	switch t {
	case DocumentTypeInvoice, DocumentTypeReceipt, DocumentTypeBankStatement,
		DocumentTypeExpenseReport, DocumentTypeCreditNote, DocumentTypeDebitNote,
		DocumentTypeOther:
		return true
	default:
		return false
	}
}

// FileFormat represents the file format of a document.
type FileFormat string

const (
	FileFormatPDF  FileFormat = "pdf"
	FileFormatJPEG FileFormat = "jpeg"
	FileFormatPNG  FileFormat = "png"
	FileFormatTIFF FileFormat = "tiff"
	FileFormatXLSX FileFormat = "xlsx"
	FileFormatCSV  FileFormat = "csv"
)

// IsValid checks if the file format is valid.
func (f FileFormat) IsValid() bool {
	switch f {
	case FileFormatPDF, FileFormatJPEG, FileFormatPNG, FileFormatTIFF,
		FileFormatXLSX, FileFormatCSV:
		return true
	default:
		return false
	}
}

// MIMEType returns the MIME type for the file format.
func (f FileFormat) MIMEType() string {
	switch f {
	case FileFormatPDF:
		return "application/pdf"
	case FileFormatJPEG:
		return "image/jpeg"
	case FileFormatPNG:
		return "image/png"
	case FileFormatTIFF:
		return "image/tiff"
	case FileFormatXLSX:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case FileFormatCSV:
		return "text/csv"
	default:
		return "application/octet-stream"
	}
}

// FileFormatFromMIME returns the FileFormat from a MIME type.
func FileFormatFromMIME(mime string) (FileFormat, error) {
	switch mime {
	case "application/pdf":
		return FileFormatPDF, nil
	case "image/jpeg", "image/jpg":
		return FileFormatJPEG, nil
	case "image/png":
		return FileFormatPNG, nil
	case "image/tiff", "image/tif":
		return FileFormatTIFF, nil
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return FileFormatXLSX, nil
	case "text/csv":
		return FileFormatCSV, nil
	default:
		return "", fmt.Errorf("unsupported MIME type: %s", mime)
	}
}

// Document represents an uploaded document awaiting processing.
type Document struct {
	ID                     uuid.UUID      `json:"id" db:"id"`
	TenantID               uuid.UUID      `json:"tenant_id" db:"tenant_id"`
	UploadedBy             uuid.UUID      `json:"uploaded_by" db:"uploaded_by"`
	Status                 DocumentStatus `json:"status" db:"status"`
	DocumentType           DocumentType   `json:"document_type" db:"document_type"`
	OriginalFilename       string         `json:"original_filename" db:"original_filename"`
	StoragePath            string         `json:"storage_path" db:"storage_path"`
	FileSizeBytes          int64          `json:"file_size_bytes" db:"file_size_bytes"`
	FileFormat             FileFormat     `json:"file_format" db:"file_format"`
	PageCount              int            `json:"page_count" db:"page_count"`
	DetectedLanguage       string         `json:"detected_language" db:"detected_language"`
	UploadID               uuid.UUID      `json:"upload_id" db:"upload_id"`
	ProcessingStartedAt    *time.Time     `json:"processing_started_at" db:"processing_started_at"`
	ProcessingCompletedAt  *time.Time     `json:"processing_completed_at" db:"processing_completed_at"`
	ClassificationConfidence float64       `json:"classification_confidence" db:"classification_confidence"`
	OverallConfidence      float64        `json:"overall_confidence" db:"overall_confidence"`
	RejectionReason        string         `json:"rejection_reason" db:"rejection_reason"`
	LockedBy               uuid.UUID      `json:"locked_by" db:"locked_by"`
	LockedAt               *time.Time     `json:"locked_at" db:"locked_at"`
	CreatedAt              time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" db:"updated_at"`
}

// Validate checks if the Document has valid field values.
func (d *Document) Validate() error {
	var errs []error

	if d.ID == uuid.Nil {
		errs = append(errs, errors.New("id is required"))
	}

	if d.TenantID == uuid.Nil {
		errs = append(errs, errors.New("tenant_id is required"))
	}

	if d.UploadedBy == uuid.Nil {
		errs = append(errs, errors.New("uploaded_by is required"))
	}

	if !d.Status.IsValid() {
		errs = append(errs, fmt.Errorf("invalid status: %s", d.Status))
	}

	if d.DocumentType != "" && !d.DocumentType.IsValid() {
		errs = append(errs, fmt.Errorf("invalid document_type: %s", d.DocumentType))
	}

	if d.OriginalFilename == "" {
		errs = append(errs, errors.New("original_filename is required"))
	}

	if d.StoragePath == "" {
		errs = append(errs, errors.New("storage_path is required"))
	}

	if d.FileSizeBytes <= 0 {
		errs = append(errs, errors.New("file_size_bytes must be positive"))
	}

	if d.FileSizeBytes > 26214400 { // 25MB
		errs = append(errs, errors.New("file_size_bytes exceeds maximum of 25MB"))
	}

	if !d.FileFormat.IsValid() {
		errs = append(errs, fmt.Errorf("invalid file_format: %s", d.FileFormat))
	}

	if d.PageCount < 1 || d.PageCount > 20 {
		errs = append(errs, fmt.Errorf("page_count must be between 1 and 20, got %d", d.PageCount))
	}

	if d.ClassificationConfidence < 0 || d.ClassificationConfidence > 1 {
		errs = append(errs, errors.New("classification_confidence must be between 0 and 1"))
	}

	if d.OverallConfidence < 0 || d.OverallConfidence > 1 {
		errs = append(errs, errors.New("overall_confidence must be between 0 and 1"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %w", errors.Join(errs...))
	}

	return nil
}

// IsLocked returns true if the document is currently locked for review.
func (d *Document) IsLocked() bool {
	return d.LockedBy != uuid.Nil && d.LockedAt != nil
}

// IsLockedBy returns true if the document is locked by the specified user.
func (d *Document) IsLockedBy(userID uuid.UUID) bool {
	return d.LockedBy == userID
}

// IsProcessing returns true if the document is currently being processed.
func (d *Document) IsProcessing() bool {
	switch d.Status {
	case DocumentStatusUploading, DocumentStatusClassifying, DocumentStatusExtracting:
		return true
	default:
		return false
	}
}

// IsReviewable returns true if the document can be reviewed.
func (d *Document) IsReviewable() bool {
	return d.Status == DocumentStatusReadyForReview || d.Status == DocumentStatusUnderReview
}

// NeedsReview returns true if the document needs human review.
func (d *Document) NeedsReview() bool {
	return d.Status == DocumentStatusReadyForReview
}

// NewDocument creates a new Document with the provided parameters.
func NewDocument(tenantID, uploadedBy uuid.UUID, filename, storagePath string, fileSize int64, format FileFormat) *Document {
	now := time.Now()
	return &Document{
		ID:               uuid.New(),
		TenantID:         tenantID,
		UploadedBy:       uploadedBy,
		Status:           DocumentStatusUploading,
		OriginalFilename: filename,
		StoragePath:      storagePath,
		FileSizeBytes:    fileSize,
		FileFormat:       format,
		PageCount:        1,
		DetectedLanguage: "en",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// DocumentListFilter represents filters for listing documents.
type DocumentListFilter struct {
	TenantID     uuid.UUID
	Statuses     []DocumentStatus
	Types        []DocumentType
	UploadedBy   uuid.UUID
	DateFrom     *time.Time
	DateTo       *time.Time
	UploadID     uuid.UUID
	Search       string
	Page         int
	PageSize     int
	SortBy       string
	SortOrder    string
}

// DocumentListResult represents a paginated list of documents.
type DocumentListResult struct {
	Documents  []Document `json:"documents"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}
