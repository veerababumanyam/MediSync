package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FieldType represents the type of an extracted field.
type FieldType string

const (
	FieldTypeString      FieldType = "string"
	FieldTypeNumber      FieldType = "number"
	FieldTypeCurrency    FieldType = "currency"
	FieldTypeDate        FieldType = "date"
	FieldTypePercentage  FieldType = "percentage"
	FieldTypeIdentifier  FieldType = "identifier"
	FieldTypeTaxID       FieldType = "tax_id"
)

// IsValid checks if the field type is valid.
func (t FieldType) IsValid() bool {
	switch t {
	case FieldTypeString, FieldTypeNumber, FieldTypeCurrency,
		FieldTypeDate, FieldTypePercentage, FieldTypeIdentifier, FieldTypeTaxID:
		return true
	default:
		return false
	}
}

// VerificationStatus represents the verification status of an extracted field.
type VerificationStatus string

const (
	VerificationStatusPending           VerificationStatus = "pending"
	VerificationStatusAutoAccepted      VerificationStatus = "auto_accepted"
	VerificationStatusNeedsReview       VerificationStatus = "needs_review"
	VerificationStatusHighPriority      VerificationStatus = "high_priority"
	VerificationStatusManuallyVerified  VerificationStatus = "manually_verified"
	VerificationStatusManuallyCorrected VerificationStatus = "manually_corrected"
	VerificationStatusRejected          VerificationStatus = "rejected"
)

// IsValid checks if the verification status is valid.
func (s VerificationStatus) IsValid() bool {
	switch s {
	case VerificationStatusPending, VerificationStatusAutoAccepted,
		VerificationStatusNeedsReview, VerificationStatusHighPriority,
		VerificationStatusManuallyVerified, VerificationStatusManuallyCorrected,
		VerificationStatusRejected:
		return true
	default:
		return false
	}
}

// NeedsAttention returns true if the field requires human attention.
func (s VerificationStatus) NeedsAttention() bool {
	return s == VerificationStatusNeedsReview || s == VerificationStatusHighPriority || s == VerificationStatusPending
}

// BoundingBox represents the location of a field in the original document.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Page   int     `json:"page,omitempty"`
}

// ExtractedField represents a single field extracted from a document.
type ExtractedField struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	DocumentID        uuid.UUID         `json:"document_id" db:"document_id"`
	PageNumber        int               `json:"page_number" db:"page_number"`
	FieldName         string            `json:"field_name" db:"field_name"`
	FieldType         FieldType         `json:"field_type" db:"field_type"`
	ExtractedValue    string            `json:"extracted_value" db:"extracted_value"`
	ConfidenceScore   float64           `json:"confidence_score" db:"confidence_score"`
	BoundingBox       *BoundingBox      `json:"bounding_box" db:"bounding_box"`
	IsHandwritten     bool              `json:"is_handwritten" db:"is_handwritten"`
	VerificationStatus VerificationStatus `json:"verification_status" db:"verification_status"`
	VerifiedBy        uuid.UUID         `json:"verified_by" db:"verified_by"`
	VerifiedAt        *time.Time        `json:"verified_at" db:"verified_at"`
	OriginalValue     string            `json:"original_value" db:"original_value"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

// Validate checks if the ExtractedField has valid field values.
func (f *ExtractedField) Validate() error {
	var errs []error

	if f.ID == uuid.Nil {
		errs = append(errs, errors.New("id is required"))
	}

	if f.DocumentID == uuid.Nil {
		errs = append(errs, errors.New("document_id is required"))
	}

	if f.PageNumber < 1 {
		errs = append(errs, errors.New("page_number must be positive"))
	}

	if f.FieldName == "" {
		errs = append(errs, errors.New("field_name is required"))
	}

	if !f.FieldType.IsValid() {
		errs = append(errs, fmt.Errorf("invalid field_type: %s", f.FieldType))
	}

	if f.ConfidenceScore < 0 || f.ConfidenceScore > 1 {
		errs = append(errs, errors.New("confidence_score must be between 0 and 1"))
	}

	// Handwritten content cannot have confidence > 0.85
	if f.IsHandwritten && f.ConfidenceScore > 0.85 {
		errs = append(errs, errors.New("handwritten content cannot have confidence > 0.85"))
	}

	if !f.VerificationStatus.IsValid() {
		errs = append(errs, fmt.Errorf("invalid verification_status: %s", f.VerificationStatus))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %w", errors.Join(errs...))
	}

	return nil
}

// WasEdited returns true if the field value was manually changed.
func (f *ExtractedField) WasEdited() bool {
	return f.OriginalValue != "" && f.OriginalValue != f.ExtractedValue
}

// IsVerified returns true if the field has been verified (auto or manual).
func (f *ExtractedField) IsVerified() bool {
	switch f.VerificationStatus {
	case VerificationStatusAutoAccepted, VerificationStatusManuallyVerified, VerificationStatusManuallyCorrected:
		return true
	default:
		return false
	}
}

// NeedsReview returns true if the field needs human review.
func (f *ExtractedField) NeedsReview() bool {
	return f.VerificationStatus == VerificationStatusNeedsReview ||
		f.VerificationStatus == VerificationStatusHighPriority ||
		f.VerificationStatus == VerificationStatusPending
}

// IsHighPriority returns true if the field is high priority (low confidence).
func (f *ExtractedField) IsHighPriority() bool {
	return f.VerificationStatus == VerificationStatusHighPriority || f.ConfidenceScore < 0.70
}

// CalculateVerificationStatus determines the verification status based on confidence and handwriting.
func CalculateVerificationStatus(confidence float64, isHandwritten bool) VerificationStatus {
	// Cap confidence for handwritten content
	if isHandwritten && confidence > 0.85 {
		confidence = 0.85
	}

	if confidence >= 0.95 {
		return VerificationStatusAutoAccepted
	} else if confidence >= 0.70 {
		return VerificationStatusNeedsReview
	} else {
		return VerificationStatusHighPriority
	}
}

// NewExtractedField creates a new ExtractedField with the provided parameters.
func NewExtractedField(documentID uuid.UUID, pageNumber int, fieldName string, fieldType FieldType, value string, confidence float64, isHandwritten bool) *ExtractedField {
	now := time.Now()
	status := CalculateVerificationStatus(confidence, isHandwritten)

	// Cap confidence for handwritten content
	if isHandwritten && confidence > 0.85 {
		confidence = 0.85
	}

	return &ExtractedField{
		ID:                uuid.New(),
		DocumentID:        documentID,
		PageNumber:        pageNumber,
		FieldName:         fieldName,
		FieldType:         fieldType,
		ExtractedValue:    value,
		ConfidenceScore:   confidence,
		IsHandwritten:     isHandwritten,
		VerificationStatus: status,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// SetValue updates the field value and marks it as manually corrected.
func (f *ExtractedField) SetValue(newValue string, verifiedBy uuid.UUID) {
	if f.ExtractedValue != newValue {
		if f.OriginalValue == "" {
			f.OriginalValue = f.ExtractedValue
		}
		f.ExtractedValue = newValue
		f.VerificationStatus = VerificationStatusManuallyCorrected
	} else {
		f.VerificationStatus = VerificationStatusManuallyVerified
	}
	f.VerifiedBy = verifiedBy
	now := time.Now()
	f.VerifiedAt = &now
	f.UpdatedAt = now
}

// Verify marks the field as verified without changing the value.
func (f *ExtractedField) Verify(verifiedBy uuid.UUID) {
	f.VerificationStatus = VerificationStatusManuallyVerified
	f.VerifiedBy = verifiedBy
	now := time.Time{}
	f.VerifiedAt = &now
	f.UpdatedAt = now
}

// Common field names for different document types
const (
	// Invoice fields
	FieldNameSupplierName   = "supplier_name"
	FieldNameSupplierTaxID  = "supplier_tax_id"
	FieldNameInvoiceNumber  = "invoice_number"
	FieldNameInvoiceDate    = "invoice_date"
	FieldNameDueDate        = "due_date"
	FieldNameSubtotal       = "subtotal"
	FieldNameTaxAmount      = "tax_amount"
	FieldNameTaxRate        = "tax_rate"
	FieldNameTotal          = "total"
	FieldNameCurrency       = "currency"

	// Bank statement fields
	FieldNameBankName       = "bank_name"
	FieldNameAccountNumber  = "account_number"
	FieldNameAccountName    = "account_name"
	FieldNameStatementDate  = "statement_date"
	FieldNameOpeningBalance = "opening_balance"
	FieldNameClosingBalance = "closing_balance"

	// Receipt fields
	FieldNameReceiptNumber  = "receipt_number"
	FieldNamePaymentMethod  = "payment_method"
	FieldNameAmount         = "amount"
)
