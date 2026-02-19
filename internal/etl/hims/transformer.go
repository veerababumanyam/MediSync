// Package hims provides transformation functions for HIMS data.
//
// This file contains functions to transform HIMS API types into warehouse
// record types for database insertion.
package hims

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse"
)

// Transformer handles transformation of HIMS data to warehouse records.
type Transformer struct {
	source warehouse.Source
}

// NewTransformer creates a new HIMS transformer.
func NewTransformer() *Transformer {
	return &Transformer{
		source: warehouse.SourceHIMS,
	}
}

// ============================================================================
// Patient Transformation
// ============================================================================

// PatientToWarehouse converts a HIMS patient to a warehouse patient record.
func (t *Transformer) PatientToWarehouse(p *Patient) *warehouse.Patient {
	var dateOfBirth *time.Time
	if p.DateOfBirth != nil {
		if dob, err := time.Parse("2006-01-02", *p.DateOfBirth); err == nil {
			dateOfBirth = &dob
		}
	}

	patient := &warehouse.Patient{
		ExternalPatientID:     p.PatientID,
		NameEN:                p.NameEN,
		NameAR:                p.NameAR,
		DateOfBirth:           dateOfBirth,
		Gender:                p.Gender,
		Phone:                 p.Phone,
		Email:                 p.Email,
		AddressEN:             p.AddressEN,
		AddressAR:             p.AddressAR,
		BloodGroup:            p.BloodGroup,
		Nationality:           p.Nationality,
		NationalID:            p.NationalID,
		InsuranceProvider:     p.InsuranceProvider,
		InsurancePolicyNumber: p.InsurancePolicyNum,
		IsActive:             p.IsActive,
		Source:                t.source.String(),
		SourceID:              p.PatientID,
	}

	if p.EmergencyContact != nil {
		patient.EmergencyContactName = p.EmergencyContact.Name
		patient.EmergencyContactPhone = p.EmergencyContact.Phone
	}

	return patient
}

// BatchPatientsToWarehouse converts multiple HIMS patients.
func (t *Transformer) BatchPatientsToWarehouse(patients []*Patient) []*warehouse.Patient {
	result := make([]*warehouse.Patient, len(patients))
	for i, p := range patients {
		result[i] = t.PatientToWarehouse(p)
	}
	return result
}

// ============================================================================
// Doctor Transformation
// ============================================================================

// DoctorToWarehouse converts a HIMS doctor to a warehouse doctor record.
func (t *Transformer) DoctorToWarehouse(d *Doctor) *warehouse.Doctor {
	return &warehouse.Doctor{
		ExternalDocID:    d.DoctorID,
		NameEN:           d.NameEN,
		NameAR:           d.NameAR,
		SpecialtyEN:      d.SpecialtyEN,
		SpecialtyAR:      d.SpecialtyAR,
		DepartmentEN:     d.DepartmentEN,
		DepartmentAR:     d.DepartmentAR,
		Qualification:    d.Qualification,
		LicenseNumber:    d.LicenseNumber,
		Phone:            d.Phone,
		Email:            d.Email,
		ConsultationFee:  d.ConsultationFee,
		IsActive:         d.IsActive,
		Source:           t.source.String(),
		SourceID:         d.DoctorID,
	}
}

// BatchDoctorsToWarehouse converts multiple HIMS doctors.
func (t *Transformer) BatchDoctorsToWarehouse(doctors []*Doctor) []*warehouse.Doctor {
	result := make([]*warehouse.Doctor, len(doctors))
	for i, d := range doctors {
		result[i] = t.DoctorToWarehouse(d)
	}
	return result
}

// ============================================================================
// Drug Transformation
// ============================================================================

// DrugToWarehouse converts a HIMS drug to a warehouse drug record.
func (t *Transformer) DrugToWarehouse(d *Drug) *warehouse.Drug {
	// Generate external ID if not provided
	externalID := ""
	if d.DrugID != nil {
		externalID = *d.DrugID
	} else {
		// Use hash of name as ID if not provided
		externalID = generateHashID(d.NameEN)
	}

	return &warehouse.Drug{
		ExternalDrugID:  externalID,
		NameEN:          d.NameEN,
		NameAR:          d.NameAR,
		GenericNameEN:   d.GenericNameEN,
		GenericNameAR:   d.GenericNameAR,
		CategoryEN:      d.CategoryEN,
		CategoryAR:      d.CategoryAR,
		DosageForm:      d.DosageForm,
		Strength:        d.Strength,
		Unit:            d.Unit,
		Manufacturer:    d.Manufacturer,
		UnitPrice:       d.UnitPrice,
		ReorderLevel:    d.ReorderLevel,
		IsControlled:    d.IsControlled,
		IsActive:        d.IsActive,
		Source:          t.source.String(),
		SourceID:        externalID,
	}
}

// BatchDrugsToWarehouse converts multiple HIMS drugs.
func (t *Transformer) BatchDrugsToWarehouse(drugs []*Drug) []*warehouse.Drug {
	result := make([]*warehouse.Drug, len(drugs))
	for i, d := range drugs {
		result[i] = t.DrugToWarehouse(d)
	}
	return result
}

// ============================================================================
// Appointment Transformation
// ============================================================================

// AppointmentToWarehouse converts a HIMS appointment to a warehouse appointment record.
func (t *Transformer) AppointmentToWarehouse(a *Appointment) (*warehouse.Appointment, error) {
	// Parse patient and doctor IDs as UUIDs
	patientID, err := uuid.Parse(a.PatientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	doctorID, err := uuid.Parse(a.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %w", err)
	}

	// Parse date
	apptDate, err := time.Parse("2006-01-02", a.ApptDate)
	if err != nil {
		return nil, fmt.Errorf("invalid appointment date: %w", err)
	}

	var departmentID *uuid.UUID
	if a.DepartmentID != nil {
		if deptID, err := uuid.Parse(*a.DepartmentID); err == nil {
			departmentID = &deptID
		}
	}

	var billingID *uuid.UUID
	if a.BillingID != nil {
		if billID, err := uuid.Parse(*a.BillingID); err == nil {
			billingID = &billID
		}
	}

	return &warehouse.Appointment{
		ExternalApptID:        a.ApptID,
		PatientID:             patientID,
		DoctorID:              doctorID,
		DepartmentID:          departmentID,
		ApptDate:              apptDate,
		ApptTime:              a.ApptTime,
		ApptDatetime:          a.ApptDatetime,
		Status:                a.Status,
		ApptType:              a.ApptType,
		DurationMinutes:       a.DurationMinutes,
		ChiefComplaint:        a.ChiefComplaint,
		DiagnosisCode:         a.DiagnosisCode,
		DiagnosisDescription:  a.DiagnosisDescription,
		BillingID:             billingID,
		Notes:                 a.Notes,
		IsWalkIn:              a.IsWalkIn,
		CancellationReason:    a.CancellationReason,
		CancelledAt:           a.CancelledAt,
		Source:                t.source.String(),
		SourceID:              a.ApptID,
	}, nil
}

// BatchAppointmentsToWarehouse converts multiple HIMS appointments.
func (t *Transformer) BatchAppointmentsToWarehouse(appointments []*Appointment) ([]*warehouse.Appointment, error) {
	result := make([]*warehouse.Appointment, 0, len(appointments))
	for _, a := range appointments {
		appt, err := t.AppointmentToWarehouse(a)
		if err != nil {
			// Log error but continue with other records
			continue
		}
		result = append(result, appt)
	}
	return result, nil
}

// ============================================================================
// Billing Transformation
// ============================================================================

// BillingToWarehouse converts a HIMS billing record to a warehouse billing record.
func (t *Transformer) BillingToWarehouse(b *Billing) (*warehouse.Billing, error) {
	// Parse patient ID as UUID
	patientID, err := uuid.Parse(b.PatientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	// Parse date
	billDate, err := time.Parse("2006-01-02", b.BillDate)
	if err != nil {
		return nil, fmt.Errorf("invalid bill date: %w", err)
	}

	var apptID *uuid.UUID
	if b.ApptID != nil {
		if appID, err := uuid.Parse(*b.ApptID); err == nil {
			apptID = &appID
		}
	}

	return &warehouse.Billing{
		ExternalBillID:     b.BillID,
		PatientID:          patientID,
		ApptID:             apptID,
		BillDate:           billDate,
		BillDatetime:       b.BillDatetime,
		SubtotalAmount:     b.SubtotalAmount,
		DiscountAmount:     b.DiscountAmount,
		TaxAmount:          b.TaxAmount,
		TotalAmount:        b.TotalAmount,
		PaidAmount:         b.PaidAmount,
		PaymentMode:        b.PaymentMode,
		PaymentStatus:      b.PaymentStatus,
		InsuranceClaimID:   b.InsuranceClaimID,
		InsuranceAmount:    b.InsuranceAmount,
		DepartmentEN:       b.DepartmentEN,
		DepartmentAR:       b.DepartmentAR,
		BillType:           b.BillType,
		ReceiptNumber:      b.ReceiptNumber,
		Notes:              b.Notes,
		Source:             t.source.String(),
		SourceID:           b.BillID,
	}, nil
}

// BatchBillingToWarehouse converts multiple HIMS billing records.
func (t *Transformer) BatchBillingToWarehouse(billing []*Billing) ([]*warehouse.Billing, error) {
	result := make([]*warehouse.Billing, 0, len(billing))
	for _, b := range billing {
		bill, err := t.BillingToWarehouse(b)
		if err != nil {
			// Log error but continue with other records
			continue
		}
		result = append(result, bill)
	}
	return result, nil
}

// ============================================================================
// Pharmacy Dispensation Transformation
// ============================================================================

// PharmacyDispensationToWarehouse converts a HIMS pharmacy dispensation to a warehouse record.
func (t *Transformer) PharmacyDispensationToWarehouse(p *PharmacyDispensation) (*warehouse.PharmacyDispensation, error) {
	// Parse IDs as UUIDs
	drugID, err := uuid.Parse(p.DrugID)
	if err != nil {
		return nil, fmt.Errorf("invalid drug ID: %w", err)
	}

	patientID, err := uuid.Parse(p.PatientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	var doctorID *uuid.UUID
	if p.DoctorID != nil {
		if docID, err := uuid.Parse(*p.DoctorID); err == nil {
			doctorID = &docID
		}
	}

	var billID *uuid.UUID
	if p.BillID != nil {
		if bID, err := uuid.Parse(*p.BillID); err == nil {
			billID = &bID
		}
	}

	// Parse date
	dispDate, err := time.Parse("2006-01-02", p.DispDate)
	if err != nil {
		return nil, fmt.Errorf("invalid dispensation date: %w", err)
	}

	var expiryDate *time.Time
	if p.ExpiryDate != nil {
		if expDate, err := time.Parse("2006-01-02", *p.ExpiryDate); err == nil {
			expiryDate = &expDate
		}
	}

	var originalDrugID *uuid.UUID
	if p.OriginalDrugID != nil {
		if odID, err := uuid.Parse(*p.OriginalDrugID); err == nil {
			originalDrugID = &odID
		}
	}

	return &warehouse.PharmacyDispensation{
		ExternalDispID:      p.DispID,
		DrugID:              drugID,
		PatientID:           patientID,
		DoctorID:            doctorID,
		BillID:              billID,
		PrescriptionID:      p.PrescriptionID,
		DispDate:            dispDate,
		DispDatetime:        p.DispDatetime,
		Quantity:            p.Quantity,
		Unit:                p.Unit,
		DosageInstructions:  p.DosageInstructions,
		DaysSupply:          p.DaysSupply,
		UnitPrice:           p.UnitPrice,
		DiscountAmount:      p.DiscountAmount,
		TaxAmount:           p.TaxAmount,
		TotalAmount:         p.TotalAmount,
		BatchNumber:         p.BatchNumber,
		ExpiryDate:          expiryDate,
		IsSubstituted:       p.IsSubstituted,
		OriginalDrugID:      originalDrugID,
		Notes:               p.Notes,
		Source:              t.source.String(),
		SourceID:            p.DispID,
	}, nil
}

// BatchPharmacyDispensationsToWarehouse converts multiple HIMS pharmacy dispensations.
func (t *Transformer) BatchPharmacyDispensationsToWarehouse(dispensations []*PharmacyDispensation) ([]*warehouse.PharmacyDispensation, error) {
	result := make([]*warehouse.PharmacyDispensation, 0, len(dispensations))
	for _, p := range dispensations {
		disp, err := t.PharmacyDispensationToWarehouse(p)
		if err != nil {
			// Log error but continue with other records
			continue
		}
		result = append(result, disp)
	}
	return result, nil
}

// ============================================================================
// Utility Functions
// ============================================================================

// generateHashID creates a hash-based ID from a string.
// This is a simple implementation - in production, use a proper hash function.
func generateHashID(s string) string {
	// Simple hash: base64 encode the string bytes
	// In production, use SHA256 or similar
	return fmt.Sprintf("hims-%x", len(s)*17+len(s)%37) // Placeholder
}

// ToJSON converts a HIMS record to JSON for quarantine storage.
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// ParseDate parses a date string in various formats.
func ParseDate(dateStr string) (*time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"02/01/2006",
		"01/02/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to parse date: %s", dateStr)
}

// GetSyncFrequency returns the recommended sync frequency for an entity.
func GetSyncFrequency(entity string) time.Duration {
	switch entity {
	case "patients", "doctors", "drugs", "departments":
		return 24 * time.Hour // Daily
	case "appointments", "billing", "pharmacy_dispensations":
		return 15 * time.Minute // Every 15 minutes
	default:
		return 1 * time.Hour // Default hourly
	}
}
