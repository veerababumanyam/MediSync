package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action performed on a document.
type AuditAction string

const (
	AuditActionUploaded      AuditAction = "uploaded"
	AuditActionClassified    AuditAction = "classified"
	AuditActionExtracted     AuditAction = "extracted"
	AuditActionReviewStarted AuditAction = "review_started"
	AuditActionFieldEdited   AuditAction = "field_edited"
	AuditActionFieldVerified AuditAction = "field_verified"
	AuditActionApproved      AuditAction = "approved"
	AuditActionRejected      AuditAction = "rejected"
	AuditActionReprocessed   AuditAction = "reprocessed"
)

// IsValid checks if the audit action is valid.
func (a AuditAction) IsValid() bool {
	switch a {
	case AuditActionUploaded, AuditActionClassified, AuditActionExtracted,
		AuditActionReviewStarted, AuditActionFieldEdited, AuditActionFieldVerified,
		AuditActionApproved, AuditActionRejected, AuditActionReprocessed:
		return true
	default:
		return false
	}
}

// ActorType represents who performed the action.
type ActorType string

const (
	ActorTypeUser   ActorType = "user"
	ActorTypeSystem ActorType = "system"
)

// IsValid checks if the actor type is valid.
func (t ActorType) IsValid() bool {
	return t == ActorTypeUser || t == ActorTypeSystem
}

// DocumentAuditLog represents an audit log entry for document actions.
type DocumentAuditLog struct {
	ID         uuid.UUID   `json:"id" db:"id"`
	TenantID   uuid.UUID   `json:"tenant_id" db:"tenant_id"`
	DocumentID uuid.UUID   `json:"document_id" db:"document_id"`
	Action     AuditAction `json:"action" db:"action"`
	ActorID    uuid.UUID   `json:"actor_id" db:"actor_id"`
	ActorType  ActorType   `json:"actor_type" db:"actor_type"`
	FieldName  string      `json:"field_name" db:"field_name"`
	OldValue   any         `json:"old_value" db:"old_value"`
	NewValue   any         `json:"new_value" db:"new_value"`
	Notes      string      `json:"notes" db:"notes"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
}

// Validate checks if the DocumentAuditLog has valid field values.
func (l *DocumentAuditLog) Validate() error {
	var errs []error

	if l.ID == uuid.Nil {
		errs = append(errs, errors.New("id is required"))
	}

	if l.TenantID == uuid.Nil {
		errs = append(errs, errors.New("tenant_id is required"))
	}

	if l.DocumentID == uuid.Nil {
		errs = append(errs, errors.New("document_id is required"))
	}

	if !l.Action.IsValid() {
		errs = append(errs, fmt.Errorf("invalid action: %s", l.Action))
	}

	if !l.ActorType.IsValid() {
		errs = append(errs, fmt.Errorf("invalid actor_type: %s", l.ActorType))
	}

	// ActorID is required for user actions
	if l.ActorType == ActorTypeUser && l.ActorID == uuid.Nil {
		errs = append(errs, errors.New("actor_id is required for user actions"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %w", errors.Join(errs...))
	}

	return nil
}

// IsFieldChange returns true if this log entry represents a field value change.
func (l *DocumentAuditLog) IsFieldChange() bool {
	return l.Action == AuditActionFieldEdited || l.Action == AuditActionFieldVerified
}

// NewDocumentAuditLog creates a new audit log entry.
func NewDocumentAuditLog(tenantID, documentID uuid.UUID, action AuditAction, actorID uuid.UUID, actorType ActorType) *DocumentAuditLog {
	return &DocumentAuditLog{
		ID:         uuid.New(),
		TenantID:   tenantID,
		DocumentID: documentID,
		Action:     action,
		ActorID:    actorID,
		ActorType:  actorType,
		CreatedAt:  time.Now(),
	}
}

// NewUserAuditLog creates a new audit log entry for a user action.
func NewUserAuditLog(tenantID, documentID, userID uuid.UUID, action AuditAction) *DocumentAuditLog {
	return NewDocumentAuditLog(tenantID, documentID, action, userID, ActorTypeUser)
}

// NewSystemAuditLog creates a new audit log entry for a system action.
func NewSystemAuditLog(tenantID, documentID uuid.UUID, action AuditAction) *DocumentAuditLog {
	return NewDocumentAuditLog(tenantID, documentID, action, uuid.Nil, ActorTypeSystem)
}

// WithFieldChange adds field change details to the audit log.
func (l *DocumentAuditLog) WithFieldChange(fieldName string, oldValue, newValue any) *DocumentAuditLog {
	l.FieldName = fieldName
	l.OldValue = oldValue
	l.NewValue = newValue
	return l
}

// WithNotes adds notes to the audit log.
func (l *DocumentAuditLog) WithNotes(notes string) *DocumentAuditLog {
	l.Notes = notes
	return l
}

// AuditLogFilter represents filters for querying audit logs.
type AuditLogFilter struct {
	TenantID   uuid.UUID
	DocumentID uuid.UUID
	Actions    []AuditAction
	ActorID    uuid.UUID
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
}

// AuditLogListResult represents a paginated list of audit log entries.
type AuditLogListResult struct {
	Entries   []DocumentAuditLog `json:"entries"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
