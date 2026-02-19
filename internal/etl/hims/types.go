// Package hims provides types for Healthcare Information Management System data.
//
// These types represent the canonical data structures from the HIMS REST API
// that are synced to the MediSync data warehouse.
package hims

import (
	"time"
)

// Patient represents a patient from HIMS.
type Patient struct {
	PatientID          string    `json:"patient_id"`
	NameEN             string    `json:"name_en"`
	NameAR             *string   `json:"name_ar,omitempty"`
	DateOfBirth        *string   `json:"date_of_birth,omitempty"` // YYYY-MM-DD
	Gender             *string   `json:"gender,omitempty"`
	Phone              *string   `json:"phone,omitempty"`
	Email              *string   `json:"email,omitempty"`
	AddressEN          *string   `json:"address_en,omitempty"`
	AddressAR          *string   `json:"address_ar,omitempty"`
	BloodGroup         *string   `json:"blood_group,omitempty"`
	Nationality        *string   `json:"nationality,omitempty"`
	NationalID         *string   `json:"national_id,omitempty"`
	InsuranceProvider  *string   `json:"insurance_provider,omitempty"`
	InsurancePolicyNum *string   `json:"insurance_policy_number,omitempty"`
	EmergencyContact   *Contact  `json:"emergency_contact,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Contact represents contact information.
type Contact struct {
	Name  *string `json:"name,omitempty"`
	Phone  *string `json:"phone,omitempty"`
	Email  *string `json:"email,omitempty"`
}

// Doctor represents a doctor from HIMS.
type Doctor struct {
	DoctorID         string    `json:"doctor_id"`
	NameEN           string    `json:"name_en"`
	NameAR           *string   `json:"name_ar,omitempty"`
	SpecialtyEN      *string   `json:"specialty_en,omitempty"`
	SpecialtyAR      *string   `json:"specialty_ar,omitempty"`
	DepartmentEN     *string   `json:"department_en,omitempty"`
	DepartmentAR     *string   `json:"department_ar,omitempty"`
	Qualification    *string   `json:"qualification,omitempty"`
	LicenseNumber    *string   `json:"license_number,omitempty"`
	Phone            *string   `json:"phone,omitempty"`
	Email            *string   `json:"email,omitempty"`
	ConsultationFee  *float64  `json:"consultation_fee,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Drug represents a drug/medication from HIMS.
type Drug struct {
	DrugID         *string  `json:"drug_id,omitempty"` // May be nil for external items
	NameEN         string   `json:"name_en"`
	NameAR         *string  `json:"name_ar,omitempty"`
	GenericNameEN  *string  `json:"generic_name_en,omitempty"`
	GenericNameAR  *string  `json:"generic_name_ar,omitempty"`
	CategoryEN     *string  `json:"category_en,omitempty"`
	CategoryAR     *string  `json:"category_ar,omitempty"`
	DosageForm     *string  `json:"dosage_form,omitempty"`
	Strength       *string  `json:"strength,omitempty"`
	Unit           *string  `json:"unit,omitempty"`
	Manufacturer   *string  `json:"manufacturer,omitempty"`
	UnitPrice      *float64 `json:"unit_price,omitempty"`
	ReorderLevel   *int     `json:"reorder_level,omitempty"`
	IsControlled   bool     `json:"is_controlled"`
	IsActive       bool     `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Department represents a department from HIMS.
type Department struct {
	DepartmentID        string    `json:"department_id"`
	NameEN              string    `json:"name_en"`
	NameAR              *string   `json:"name_ar,omitempty"`
	Code                *string   `json:"code,omitempty"`
	ParentDepartmentID  *string   `json:"parent_department_id,omitempty"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Appointment represents an appointment from HIMS.
type Appointment struct {
	ApptID              string     `json:"appointment_id"`
	PatientID           string     `json:"patient_id"`
	DoctorID            string     `json:"doctor_id"`
	DepartmentID        *string    `json:"department_id,omitempty"`
	ApptDate            string     `json:"appointment_date"` // YYYY-MM-DD
	ApptTime            *string    `json:"appointment_time,omitempty"` // HH:MM
	ApptDatetime        *time.Time `json:"appointment_datetime,omitempty"`
	Status              string     `json:"status"` // scheduled, confirmed, checked_in, in_progress, completed, cancelled, no_show
	ApptType            *string    `json:"appointment_type,omitempty"`
	DurationMinutes     *int       `json:"duration_minutes,omitempty"`
	ChiefComplaint      *string    `json:"chief_complaint,omitempty"`
	DiagnosisCode       *string    `json:"diagnosis_code,omitempty"`
	DiagnosisDescription *string   `json:"diagnosis_description,omitempty"`
	BillingID           *string    `json:"billing_id,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	IsWalkIn            bool       `json:"is_walk_in"`
	CancellationReason  *string    `json:"cancellation_reason,omitempty"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Billing represents a billing record from HIMS.
type Billing struct {
	BillID              string     `json:"bill_id"`
	PatientID           string     `json:"patient_id"`
	ApptID              *string    `json:"appointment_id,omitempty"`
	BillDate            string     `json:"bill_date"` // YYYY-MM-DD
	BillDatetime        *time.Time `json:"bill_datetime,omitempty"`
	SubtotalAmount      float64    `json:"subtotal_amount"`
	DiscountAmount      float64    `json:"discount_amount"`
	TaxAmount           float64    `json:"tax_amount"`
	TotalAmount         float64    `json:"total_amount"`
	PaidAmount          float64    `json:"paid_amount"`
	PaymentMode         *string    `json:"payment_mode,omitempty"`
	PaymentStatus       string     `json:"payment_status"` // pending, partial, paid, cancelled, refunded
	InsuranceClaimID    *string    `json:"insurance_claim_id,omitempty"`
	InsuranceAmount     float64    `json:"insurance_amount"`
	DepartmentEN        *string    `json:"department_en,omitempty"`
	DepartmentAR        *string    `json:"department_ar,omitempty"`
	BillType            *string    `json:"bill_type,omitempty"`
	ReceiptNumber       *string    `json:"receipt_number,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// PharmacyDispensation represents a pharmacy dispensation from HIMS.
type PharmacyDispensation struct {
	DispID             string     `json:"dispensation_id"`
	DrugID             string     `json:"drug_id"`
	PatientID          string     `json:"patient_id"`
	DoctorID           *string    `json:"doctor_id,omitempty"`
	BillID             *string    `json:"bill_id,omitempty"`
	PrescriptionID     *string    `json:"prescription_id,omitempty"`
	DispDate           string     `json:"dispensation_date"` // YYYY-MM-DD
	DispDatetime       *time.Time `json:"dispensation_datetime,omitempty"`
	Quantity           int        `json:"quantity"`
	Unit               *string    `json:"unit,omitempty"`
	DosageInstructions *string    `json:"dosage_instructions,omitempty"`
	DaysSupply         *int       `json:"days_supply,omitempty"`
	UnitPrice          float64    `json:"unit_price"`
	DiscountAmount     float64    `json:"discount_amount"`
	TaxAmount          float64    `json:"tax_amount"`
	TotalAmount        float64    `json:"total_amount"`
	BatchNumber        *string    `json:"batch_number,omitempty"`
	ExpiryDate         *string    `json:"expiry_date,omitempty"` // YYYY-MM-DD
	IsSubstituted      bool       `json:"is_substituted"`
	OriginalDrugID     *string    `json:"original_drug_id,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// PagedResponse represents a paginated API response.
type PagedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	TotalCount int         `json:"total_count,omitempty"`
	TotalPages int         `json:"total_pages,omitempty"`
}

// APIError represents an error response from the HIMS API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Details != "" {
		return e.Message + ": " + e.Details
	}
	return e.Message
}
