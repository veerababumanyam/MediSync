// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the document processing endpoints for upload,
// classification, extraction, review, and approval workflows.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/storage"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// DocumentHandler handles document processing endpoints.
type DocumentHandler struct {
	logger        *slog.Logger
	documentRepo  *warehouse.DocumentRepository
	auditLogRepo  *warehouse.AuditLogRepository
	storageClient *storage.ObjectStorage
	natsClient    NATSPublisher
}

// NATSPublisher interface for publishing events.
type NATSPublisher interface {
	PublishDocumentUploaded(ctx context.Context, tenantID, documentID uuid.UUID) error
	PublishDocumentApproved(ctx context.Context, tenantID, documentID uuid.UUID) error
}

// DocumentHandlerConfig holds configuration for the DocumentHandler.
type DocumentHandlerConfig struct {
	Logger        *slog.Logger
	DocumentRepo  *warehouse.DocumentRepository
	AuditLogRepo  *warehouse.AuditLogRepository
	StorageClient *storage.ObjectStorage
	NATSClient    NATSPublisher
}

// NewDocumentHandler creates a new DocumentHandler instance.
func NewDocumentHandler(cfg DocumentHandlerConfig) *DocumentHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &DocumentHandler{
		logger:        cfg.Logger,
		documentRepo:  cfg.DocumentRepo,
		auditLogRepo:  cfg.AuditLogRepo,
		storageClient: cfg.StorageClient,
		natsClient:    cfg.NATSClient,
	}
}

// RegisterRoutes registers document routes on the given router.
func (h *DocumentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/documents", func(r chi.Router) {
		r.Get("/", h.HandleListDocuments)
		r.Post("/", h.HandleUploadDocument)
		r.Get("/stats", h.HandleGetStats)
		r.Get("/{document_id}", h.HandleGetDocument)
		r.Patch("/{document_id}", h.HandleUpdateDocument)
		r.Delete("/{document_id}", h.HandleDeleteDocument)
		r.Post("/{document_id}/lock", h.HandleLockDocument)
		r.Delete("/{document_id}/lock", h.HandleUnlockDocument)
		r.Get("/{document_id}/fields", h.HandleGetFields)
		r.Patch("/{document_id}/fields/{field_id}", h.HandleUpdateField)
		r.Post("/{document_id}/fields/{field_id}/verify", h.HandleVerifyField)
		r.Post("/{document_id}/approve", h.HandleApproveDocument)
		r.Post("/{document_id}/reject", h.HandleRejectDocument)
		r.Post("/{document_id}/reprocess", h.HandleReprocessDocument)
		r.Get("/{document_id}/audit-log", h.HandleGetAuditLog)
		r.Get("/{document_id}/line-items", h.HandleGetLineItems)
		r.Post("/bulk-upload", h.HandleBulkUpload)
	})
}

// ============================================================================
// Request/Response Types
// ============================================================================

// UploadDocumentRequest represents a document upload request.
type UploadDocumentRequest struct {
	UploadID string `json:"uploadId,omitempty"`
}

// UpdateDocumentRequest represents a request to update document metadata.
type UpdateDocumentRequest struct {
	DocumentType *string `json:"documentType,omitempty"`
}

// UpdateFieldRequest represents a request to update an extracted field.
type UpdateFieldRequest struct {
	Value string `json:"value"`
}

// ApproveDocumentRequest represents a request to approve a document.
type ApproveDocumentRequest struct {
	Notes string `json:"notes,omitempty"`
}

// RejectDocumentRequest represents a request to reject a document.
type RejectDocumentRequest struct {
	Reason string `json:"reason"`
}

// LockDocumentRequest represents a request to lock a document for review.
type LockDocumentRequest struct {
	// Empty body - lock is acquired by the authenticated user
}

// DocumentResponse represents a document in API responses.
type DocumentResponse struct {
	ID                     string   `json:"id"`
	TenantID               string   `json:"tenantId"`
	UploadedBy             string   `json:"uploadedBy"`
	Status                 string   `json:"status"`
	DocumentType           string   `json:"documentType,omitempty"`
	OriginalFilename       string   `json:"originalFilename"`
	FileSizeBytes          int64    `json:"fileSizeBytes"`
	FileFormat             string   `json:"fileFormat"`
	PageCount              int      `json:"pageCount"`
	DetectedLanguage       string   `json:"detectedLanguage"`
	ProcessingStartedAt    *string  `json:"processingStartedAt,omitempty"`
	ProcessingCompletedAt  *string  `json:"processingCompletedAt,omitempty"`
	ClassificationConfidence float64 `json:"classificationConfidence"`
	OverallConfidence      float64  `json:"overallConfidence"`
	RejectionReason        string   `json:"rejectionReason,omitempty"`
	IsLocked               bool     `json:"isLocked"`
	LockedBy               string   `json:"lockedBy,omitempty"`
	LockedAt               *string  `json:"lockedAt,omitempty"`
	UploadURL              string   `json:"uploadUrl,omitempty"`
	CreatedAt              string   `json:"createdAt"`
	UpdatedAt              string   `json:"updatedAt"`
}

// DocumentListResponse represents a paginated list of documents.
type DocumentListResponse struct {
	Documents  []DocumentResponse `json:"documents"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}

// FieldResponse represents an extracted field in API responses.
type FieldResponse struct {
	ID                string  `json:"id"`
	DocumentID        string  `json:"documentId"`
	PageNumber        int     `json:"pageNumber"`
	FieldName         string  `json:"fieldName"`
	FieldType         string  `json:"fieldType"`
	ExtractedValue    string  `json:"extractedValue"`
	ConfidenceScore   float64 `json:"confidenceScore"`
	BoundingBox       *BoundingBox `json:"boundingBox,omitempty"`
	IsHandwritten     bool    `json:"isHandwritten"`
	VerificationStatus string `json:"verificationStatus"`
	VerifiedBy        string  `json:"verifiedBy,omitempty"`
	VerifiedAt        *string `json:"verifiedAt,omitempty"`
	OriginalValue     string  `json:"originalValue,omitempty"`
	WasEdited         bool    `json:"wasEdited"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// BoundingBox represents the location of a field in the document.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Page   int     `json:"page,omitempty"`
}

// FieldListResponse represents a list of extracted fields.
type FieldListResponse struct {
	Fields             []FieldResponse `json:"fields"`
	TotalFields        int             `json:"totalFields"`
	FieldsNeedingReview int            `json:"fieldsNeedingReview"`
	HighPriorityCount  int             `json:"highPriorityCount"`
	AutoAcceptedCount  int             `json:"autoAcceptedCount"`
}

// DocumentStatsResponse represents document processing statistics.
type DocumentStatsResponse struct {
	TotalDocuments      int `json:"totalDocuments"`
	PendingReview       int `json:"pendingReview"`
	UnderReview         int `json:"underReview"`
	Approved            int `json:"approved"`
	Rejected            int `json:"rejected"`
	Processing          int `json:"processing"`
	HighPriorityQueue   int `json:"highPriorityQueue"`
	AverageConfidence   float64 `json:"averageConfidence"`
}

// AuditLogResponse represents an audit log entry.
type AuditLogResponse struct {
	ID         string `json:"id"`
	DocumentID string `json:"documentId"`
	Action     string `json:"action"`
	ActorID    string `json:"actorId"`
	ActorType  string `json:"actorType"`
	FieldName  string `json:"fieldName,omitempty"`
	OldValue   any    `json:"oldValue,omitempty"`
	NewValue   any    `json:"newValue,omitempty"`
	Notes      string `json:"notes,omitempty"`
	CreatedAt  string `json:"createdAt"`
}

// BulkUploadResponse represents the result of a bulk upload operation.
type BulkUploadResponse struct {
	UploadID      string   `json:"uploadId"`
	TotalFiles    int      `json:"totalFiles"`
	UploadedFiles []string `json:"uploadedFiles"`
	FailedFiles   []BulkUploadFailure `json:"failedFiles,omitempty"`
}

// BulkUploadFailure represents a failed file upload in bulk operation.
type BulkUploadFailure struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
}

// ============================================================================
// HTTP Handlers
// ============================================================================

// HandleListDocuments handles GET /documents requests.
func (h *DocumentHandler) HandleListDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filter := h.parseListFilter(r, tenantID)

	result, err := h.documentRepo.List(ctx, filter)
	if err != nil {
		h.logger.Error("failed to list documents",
			slog.Any("error", err),
			slog.String("tenant_id", tenantID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve documents")
		return
	}

	response := DocumentListResponse{
		Documents:  make([]DocumentResponse, len(result.Documents)),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	for i, doc := range result.Documents {
		response.Documents[i] = h.documentToResponse(&doc)
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleUploadDocument handles POST /documents requests.
func (h *DocumentHandler) HandleUploadDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse multipart form (max 25MB)
	err = r.ParseMultipartForm(25 << 20)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Validate file
	if err := h.validateUpload(header); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate storage path
	documentID := uuid.New()
	storagePath := h.storageClient.GetDocumentPath(tenantID, documentID, header.Filename)

	// Upload to object storage
	if err := h.storageClient.Upload(ctx, storagePath, file, header.Size); err != nil {
		h.logger.Error("failed to upload document",
			slog.Any("error", err),
			slog.String("tenant_id", tenantID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to upload document")
		return
	}

	// Detect file format
	fileFormat, err := models.FileFormatFromMIME(header.Header.Get("Content-Type"))
	if err != nil {
		fileFormat = models.FileFormatPDF // Default to PDF
	}

	// Create document record
	doc := models.NewDocument(tenantID, userID, header.Filename, storagePath, header.Size, fileFormat)

	if err := h.documentRepo.Create(ctx, doc); err != nil {
		h.logger.Error("failed to create document record",
			slog.Any("error", err),
			slog.String("tenant_id", tenantID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to create document")
		return
	}

	// Create audit log entry
	auditLog := models.NewUserAuditLog(tenantID, doc.ID, userID, models.AuditActionUploaded)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	// Publish event for async processing
	if h.natsClient != nil {
		if err := h.natsClient.PublishDocumentUploaded(ctx, tenantID, doc.ID); err != nil {
			h.logger.Warn("failed to publish document uploaded event", slog.Any("error", err))
		}
	}

	// Generate upload URL for frontend
	response := h.documentToResponse(doc)
	response.UploadURL = h.storageClient.GetPresignedURL(storagePath, 15*time.Minute)

	h.writeJSON(w, http.StatusCreated, response)
}

// HandleGetDocument handles GET /documents/{document_id} requests.
func (h *DocumentHandler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		h.logger.Error("failed to get document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	// Verify tenant access
	if doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleUpdateDocument handles PATCH /documents/{document_id} requests.
func (h *DocumentHandler) HandleUpdateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	var req UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DocumentType != nil {
		docType := models.DocumentType(*req.DocumentType)
		if !docType.IsValid() {
			h.writeError(w, http.StatusBadRequest, "invalid document type")
			return
		}
		doc.DocumentType = docType
	}

	if err := h.documentRepo.Update(ctx, doc); err != nil {
		h.logger.Error("failed to update document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to update document")
		return
	}

	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleDeleteDocument handles DELETE /documents/{document_id} requests.
func (h *DocumentHandler) HandleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	// Only allow deletion of documents in certain states
	if doc.Status == models.DocumentStatusUnderReview {
		h.writeError(w, http.StatusConflict, "cannot delete document under review")
		return
	}

	// Delete from storage
	if err := h.storageClient.Delete(ctx, doc.StoragePath); err != nil {
		h.logger.Warn("failed to delete document from storage", slog.Any("error", err))
	}

	// Delete from database
	if err := h.documentRepo.Delete(ctx, documentID); err != nil {
		h.logger.Error("failed to delete document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to delete document")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleLockDocument handles POST /documents/{document_id}/lock requests.
func (h *DocumentHandler) HandleLockDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.IsReviewable() {
		h.writeError(w, http.StatusConflict, "document cannot be locked for review")
		return
	}

	if doc.IsLocked() && !doc.IsLockedBy(userID) {
		h.writeError(w, http.StatusConflict, "document is locked by another user")
		return
	}

	// Acquire lock
	if err := h.documentRepo.AcquireLock(ctx, documentID, userID); err != nil {
		h.logger.Error("failed to lock document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to lock document")
		return
	}

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionReviewStarted)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	doc, _ = h.documentRepo.GetByID(ctx, documentID)
	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleUnlockDocument handles DELETE /documents/{document_id}/lock requests.
func (h *DocumentHandler) HandleUnlockDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.IsLocked() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if !doc.IsLockedBy(userID) {
		h.writeError(w, http.StatusForbidden, "cannot unlock document locked by another user")
		return
	}

	// Release lock
	if err := h.documentRepo.ReleaseLock(ctx, documentID); err != nil {
		h.logger.Error("failed to unlock document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to unlock document")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetFields handles GET /documents/{document_id}/fields requests.
func (h *DocumentHandler) HandleGetFields(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	fields, err := h.documentRepo.GetFields(ctx, documentID)
	if err != nil {
		h.logger.Error("failed to get fields",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve fields")
		return
	}

	response := FieldListResponse{
		Fields: make([]FieldResponse, len(fields)),
	}

	highPriority := 0
	needsReview := 0
	autoAccepted := 0

	for i, field := range fields {
		response.Fields[i] = h.fieldToResponse(&field)
		if field.IsHighPriority() {
			highPriority++
		}
		if field.NeedsReview() {
			needsReview++
		}
		if field.VerificationStatus == models.VerificationStatusAutoAccepted {
			autoAccepted++
		}
	}

	response.TotalFields = len(fields)
	response.HighPriorityCount = highPriority
	response.FieldsNeedingReview = needsReview
	response.AutoAcceptedCount = autoAccepted

	h.writeJSON(w, http.StatusOK, response)
}

// HandleUpdateField handles PATCH /documents/{document_id}/fields/{field_id} requests.
func (h *DocumentHandler) HandleUpdateField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	fieldIDStr := chi.URLParam(r, "field_id")
	fieldID, err := uuid.Parse(fieldIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid field ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.IsLockedBy(userID) {
		h.writeError(w, http.StatusForbidden, "document must be locked by you to edit fields")
		return
	}

	var req UpdateFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	field, err := h.documentRepo.GetFieldByID(ctx, fieldID)
	if err != nil || field.DocumentID != documentID {
		h.writeError(w, http.StatusNotFound, "field not found")
		return
	}

	oldValue := field.ExtractedValue
	field.SetValue(req.Value, userID)

	if err := h.documentRepo.UpdateField(ctx, field); err != nil {
		h.logger.Error("failed to update field",
			slog.Any("error", err),
			slog.String("field_id", fieldID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to update field")
		return
	}

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionFieldEdited).
		WithFieldChange(field.FieldName, oldValue, req.Value)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	h.writeJSON(w, http.StatusOK, h.fieldToResponse(field))
}

// HandleVerifyField handles POST /documents/{document_id}/fields/{field_id}/verify requests.
func (h *DocumentHandler) HandleVerifyField(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	fieldIDStr := chi.URLParam(r, "field_id")
	fieldID, err := uuid.Parse(fieldIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid field ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.IsLockedBy(userID) {
		h.writeError(w, http.StatusForbidden, "document must be locked by you to verify fields")
		return
	}

	field, err := h.documentRepo.GetFieldByID(ctx, fieldID)
	if err != nil || field.DocumentID != documentID {
		h.writeError(w, http.StatusNotFound, "field not found")
		return
	}

	field.Verify(userID)

	if err := h.documentRepo.UpdateField(ctx, field); err != nil {
		h.logger.Error("failed to verify field",
			slog.Any("error", err),
			slog.String("field_id", fieldID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to verify field")
		return
	}

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionFieldVerified).
		WithFieldChange(field.FieldName, nil, nil)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	h.writeJSON(w, http.StatusOK, h.fieldToResponse(field))
}

// HandleApproveDocument handles POST /documents/{document_id}/approve requests.
func (h *DocumentHandler) HandleApproveDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.Status.CanTransitionTo(models.DocumentStatusApproved) {
		h.writeError(w, http.StatusConflict, "document cannot be approved in current state")
		return
	}

	var req ApproveDocumentRequest
	json.NewDecoder(r.Body).Decode(&req)

	doc.Status = models.DocumentStatusApproved
	now := time.Now()
	doc.ProcessingCompletedAt = &now

	if err := h.documentRepo.Update(ctx, doc); err != nil {
		h.logger.Error("failed to approve document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to approve document")
		return
	}

	// Release any lock
	_ = h.documentRepo.ReleaseLock(ctx, documentID)

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionApproved).
		WithNotes(req.Notes)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	// Publish event
	if h.natsClient != nil {
		if err := h.natsClient.PublishDocumentApproved(ctx, tenantID, doc.ID); err != nil {
			h.logger.Warn("failed to publish document approved event", slog.Any("error", err))
		}
	}

	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleRejectDocument handles POST /documents/{document_id}/reject requests.
func (h *DocumentHandler) HandleRejectDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if !doc.Status.CanTransitionTo(models.DocumentStatusRejected) {
		h.writeError(w, http.StatusConflict, "document cannot be rejected in current state")
		return
	}

	var req RejectDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Reason == "" {
		h.writeError(w, http.StatusBadRequest, "rejection reason is required")
		return
	}

	doc.Status = models.DocumentStatusRejected
	doc.RejectionReason = req.Reason

	if err := h.documentRepo.Update(ctx, doc); err != nil {
		h.logger.Error("failed to reject document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to reject document")
		return
	}

	// Release any lock
	_ = h.documentRepo.ReleaseLock(ctx, documentID)

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionRejected).
		WithNotes(req.Reason)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleReprocessDocument handles POST /documents/{document_id}/reprocess requests.
func (h *DocumentHandler) HandleReprocessDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	if doc.Status != models.DocumentStatusRejected && doc.Status != models.DocumentStatusFailed {
		h.writeError(w, http.StatusConflict, "only rejected or failed documents can be reprocessed")
		return
	}

	// Reset document status
	doc.Status = models.DocumentStatusUploaded
	doc.RejectionReason = ""
	doc.ProcessingStartedAt = nil
	doc.ProcessingCompletedAt = nil

	if err := h.documentRepo.Update(ctx, doc); err != nil {
		h.logger.Error("failed to reprocess document",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to reprocess document")
		return
	}

	// Create audit log
	auditLog := models.NewUserAuditLog(tenantID, documentID, userID, models.AuditActionReprocessed)
	if err := h.auditLogRepo.Create(ctx, auditLog); err != nil {
		h.logger.Warn("failed to create audit log", slog.Any("error", err))
	}

	// Publish event
	if h.natsClient != nil {
		if err := h.natsClient.PublishDocumentUploaded(ctx, tenantID, doc.ID); err != nil {
			h.logger.Warn("failed to publish document reprocess event", slog.Any("error", err))
		}
	}

	h.writeJSON(w, http.StatusOK, h.documentToResponse(doc))
}

// HandleGetAuditLog handles GET /documents/{document_id}/audit-log requests.
func (h *DocumentHandler) HandleGetAuditLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	entries, err := h.auditLogRepo.GetByDocumentID(ctx, documentID)
	if err != nil {
		h.logger.Error("failed to get audit log",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve audit log")
		return
	}

	response := make([]AuditLogResponse, len(entries))
	for i, entry := range entries {
		response[i] = h.auditLogToResponse(&entry)
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleGetStats handles GET /documents/stats requests.
func (h *DocumentHandler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := h.documentRepo.GetStats(ctx, tenantID)
	if err != nil {
		h.logger.Error("failed to get document stats",
			slog.Any("error", err),
			slog.String("tenant_id", tenantID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve stats")
		return
	}

	h.writeJSON(w, http.StatusOK, DocumentStatsResponse{
		TotalDocuments:     stats.TotalDocuments,
		PendingReview:      stats.PendingReview,
		UnderReview:        stats.UnderReview,
		Approved:           stats.Approved,
		Rejected:           stats.Rejected,
		Processing:         stats.Processing,
		HighPriorityQueue:  stats.HighPriorityQueue,
		AverageConfidence:  stats.AverageConfidence,
	})
}

// HandleGetLineItems handles GET /documents/{document_id}/line-items requests.
func (h *DocumentHandler) HandleGetLineItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentID, err := h.parseDocumentID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid document ID")
		return
	}

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	doc, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc.TenantID != tenantID {
		h.writeError(w, http.StatusNotFound, "document not found")
		return
	}

	items, err := h.documentRepo.GetLineItems(ctx, documentID)
	if err != nil {
		h.logger.Error("failed to get line items",
			slog.Any("error", err),
			slog.String("document_id", documentID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve line items")
		return
	}

	h.writeJSON(w, http.StatusOK, items)
}

// HandleBulkUpload handles POST /documents/bulk-upload requests.
func (h *DocumentHandler) HandleBulkUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := h.getTenantID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse multipart form (max 100MB for bulk)
	err = r.ParseMultipartForm(100 << 20)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	uploadID := uuid.New()
	response := BulkUploadResponse{
		UploadID:      uploadID.String(),
		UploadedFiles: []string{},
		FailedFiles:   []BulkUploadFailure{},
	}

	// Process each file
	for _, files := range r.MultipartForm.File {
		for _, header := range files {
			file, err := header.Open()
			if err != nil {
				response.FailedFiles = append(response.FailedFiles, BulkUploadFailure{
					Filename: header.Filename,
					Error:    "failed to open file",
				})
				continue
			}

			// Validate
			if err := h.validateUpload(header); err != nil {
				response.FailedFiles = append(response.FailedFiles, BulkUploadFailure{
					Filename: header.Filename,
					Error:    err.Error(),
				})
				file.Close()
				continue
			}

			// Upload
			documentID := uuid.New()
			storagePath := h.storageClient.GetDocumentPath(tenantID, documentID, header.Filename)

			var fileContent []byte
			fileContent, err = io.ReadAll(file)
			file.Close()

			if err != nil {
				response.FailedFiles = append(response.FailedFiles, BulkUploadFailure{
					Filename: header.Filename,
					Error:    "failed to read file",
				})
				continue
			}

			if err := h.storageClient.UploadBytes(ctx, storagePath, fileContent); err != nil {
				response.FailedFiles = append(response.FailedFiles, BulkUploadFailure{
					Filename: header.Filename,
					Error:    "failed to upload to storage",
				})
				continue
			}

			// Create document record
			fileFormat, _ := models.FileFormatFromMIME(header.Header.Get("Content-Type"))
			doc := models.NewDocument(tenantID, userID, header.Filename, storagePath, int64(len(fileContent)), fileFormat)
			doc.UploadID = uploadID

			if err := h.documentRepo.Create(ctx, doc); err != nil {
				response.FailedFiles = append(response.FailedFiles, BulkUploadFailure{
					Filename: header.Filename,
					Error:    "failed to create document record",
				})
				continue
			}

			// Publish event
			if h.natsClient != nil {
				_ = h.natsClient.PublishDocumentUploaded(ctx, tenantID, doc.ID)
			}

			response.UploadedFiles = append(response.UploadedFiles, doc.ID.String())
		}
	}

	response.TotalFiles = len(response.UploadedFiles) + len(response.FailedFiles)
	h.writeJSON(w, http.StatusCreated, response)
}

// ============================================================================
// Helper Functions
// ============================================================================

func (h *DocumentHandler) parseDocumentID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "document_id"))
}

func (h *DocumentHandler) getTenantID(ctx context.Context) (uuid.UUID, error) {
	tenantIDStr, ok := ctx.Value("tenant_id").(string)
	if !ok {
		return uuid.Nil, errors.New("tenant ID not found in context")
	}
	return uuid.Parse(tenantIDStr)
}

func (h *DocumentHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (h *DocumentHandler) parseListFilter(r *http.Request, tenantID uuid.UUID) *models.DocumentListFilter {
	filter := &models.DocumentListFilter{
		TenantID: tenantID,
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortOrder: "desc",
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			filter.PageSize = ps
		}
	}

	if statuses := r.URL.Query()["status"]; len(statuses) > 0 {
		for _, s := range statuses {
			if status := models.DocumentStatus(s); status.IsValid() {
				filter.Statuses = append(filter.Statuses, status)
			}
		}
	}

	if types := r.URL.Query()["type"]; len(types) > 0 {
		for _, t := range types {
			if docType := models.DocumentType(t); docType.IsValid() {
				filter.Types = append(filter.Types, docType)
			}
		}
	}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = search
	}

	return filter
}

func (h *DocumentHandler) validateUpload(header *multipart.FileHeader) error {
	// Check file size (max 25MB)
	if header.Size > 25<<20 {
		return errors.New("file size exceeds maximum of 25MB")
	}

	// Check file format
	contentType := header.Header.Get("Content-Type")
	if _, err := models.FileFormatFromMIME(contentType); err != nil {
		return errors.New("unsupported file format")
	}

	return nil
}

func (h *DocumentHandler) documentToResponse(doc *models.Document) DocumentResponse {
	var processingStartedAt, processingCompletedAt, lockedAt *string
	if doc.ProcessingStartedAt != nil {
		t := doc.ProcessingStartedAt.Format(time.RFC3339)
		processingStartedAt = &t
	}
	if doc.ProcessingCompletedAt != nil {
		t := doc.ProcessingCompletedAt.Format(time.RFC3339)
		processingCompletedAt = &t
	}
	if doc.LockedAt != nil {
		t := doc.LockedAt.Format(time.RFC3339)
		lockedAt = &t
	}

	var lockedBy string
	if doc.LockedBy != uuid.Nil {
		lockedBy = doc.LockedBy.String()
	}

	return DocumentResponse{
		ID:                      doc.ID.String(),
		TenantID:                doc.TenantID.String(),
		UploadedBy:              doc.UploadedBy.String(),
		Status:                  string(doc.Status),
		DocumentType:            string(doc.DocumentType),
		OriginalFilename:        doc.OriginalFilename,
		FileSizeBytes:           doc.FileSizeBytes,
		FileFormat:              string(doc.FileFormat),
		PageCount:               doc.PageCount,
		DetectedLanguage:        doc.DetectedLanguage,
		ProcessingStartedAt:     processingStartedAt,
		ProcessingCompletedAt:   processingCompletedAt,
		ClassificationConfidence: doc.ClassificationConfidence,
		OverallConfidence:       doc.OverallConfidence,
		RejectionReason:         doc.RejectionReason,
		IsLocked:                doc.IsLocked(),
		LockedBy:                lockedBy,
		LockedAt:                lockedAt,
		CreatedAt:               doc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:               doc.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *DocumentHandler) fieldToResponse(field *models.ExtractedField) FieldResponse {
	var verifiedAt *string
	if field.VerifiedAt != nil {
		t := field.VerifiedAt.Format(time.RFC3339)
		verifiedAt = &t
	}

	var verifiedBy string
	if field.VerifiedBy != uuid.Nil {
		verifiedBy = field.VerifiedBy.String()
	}

	var boundingBox *BoundingBox
	if field.BoundingBox != nil {
		boundingBox = &BoundingBox{
			X:      field.BoundingBox.X,
			Y:      field.BoundingBox.Y,
			Width:  field.BoundingBox.Width,
			Height: field.BoundingBox.Height,
			Page:   field.BoundingBox.Page,
		}
	}

	return FieldResponse{
		ID:                 field.ID.String(),
		DocumentID:         field.DocumentID.String(),
		PageNumber:         field.PageNumber,
		FieldName:          field.FieldName,
		FieldType:          string(field.FieldType),
		ExtractedValue:     field.ExtractedValue,
		ConfidenceScore:    field.ConfidenceScore,
		BoundingBox:        boundingBox,
		IsHandwritten:      field.IsHandwritten,
		VerificationStatus: string(field.VerificationStatus),
		VerifiedBy:         verifiedBy,
		VerifiedAt:         verifiedAt,
		OriginalValue:      field.OriginalValue,
		WasEdited:          field.WasEdited(),
		CreatedAt:          field.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          field.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *DocumentHandler) auditLogToResponse(entry *models.DocumentAuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:         entry.ID.String(),
		DocumentID: entry.DocumentID.String(),
		Action:     string(entry.Action),
		ActorID:    entry.ActorID.String(),
		ActorType:  string(entry.ActorType),
		FieldName:  entry.FieldName,
		OldValue:   entry.OldValue,
		NewValue:   entry.NewValue,
		Notes:      entry.Notes,
		CreatedAt:  entry.CreatedAt.Format(time.RFC3339),
	}
}

func (h *DocumentHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *DocumentHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}
