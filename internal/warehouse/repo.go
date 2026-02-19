// Package warehouse provides the database repository layer for the MediSync data warehouse.
//
// This package handles all database operations for the analytics schemas (hims_analytics,
// tally_analytics) and application tables (app schema). It uses pgx/v5 for PostgreSQL
// connectivity and implements idempotent upsert operations for ETL sync.
//
// All operations use the medisync_etl role for INSERT/UPDATE on analytics schemas,
// ensuring proper audit trails via _source, _source_id, and _synced_at columns.
//
// Usage:
//
//	cfg := config.MustLoad()
//	repo, err := warehouse.NewRepo(cfg.Database, logger)
//	if err != nil {
//	    log.Fatal("Failed to create warehouse:", err)
//	}
//	defer repo.Close()
//
//	ctx := context.Background()
//	err = repo.UpsertPatient(ctx, patientRecord)
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// Source represents the data source system.
type Source string

const (
	// SourceTally indicates data from Tally ERP.
	SourceTally Source = "tally"
	// SourceHIMS indicates data from Healthcare Information Management System.
	SourceHIMS Source = "hims"
	// SourceBank indicates data from bank feeds/API.
	SourceBank Source = "bank"
)

// String returns the string representation of the source.
func (s Source) String() string {
	return string(s)
}

// Repo provides database operations for the data warehouse.
type Repo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// RepoConfig holds configuration for creating a new Repo.
type RepoConfig struct {
	// DSN is the PostgreSQL connection string.
	DSN string

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int32

	// MinConns is the minimum number of idle connections in the pool.
	MinConns int32

	// MaxConnLifetime is the maximum lifetime of a connection.
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum idle time of a connection.
	MaxConnIdleTime time.Duration

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewRepo creates a new warehouse repository with a connection pool.
func NewRepo(cfg interface{}, logger *slog.Logger) (*Repo, error) {
	// Extract DSN from config
	var dsn string
	switch c := cfg.(type) {
	case string:
		dsn = c
	case map[string]interface{}:
		if url, ok := c["url"].(string); ok {
			dsn = url
		}
	}

	if dsn == "" {
		return nil, fmt.Errorf("warehouse: invalid configuration, missing DSN")
	}

	if logger == nil {
		logger = slog.Default()
	}

	// Parse connection config
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to parse DSN: %w", err)
	}

	// Set pool defaults
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.MaxConnIdleTime = 1 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to create connection pool: %w", err)
	}

	return &Repo{
		pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool.
func (r *Repo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// Ping checks if the database connection is alive.
func (r *Repo) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// Pool returns the underlying connection pool.
func (r *Repo) Pool() *pgxpool.Pool {
	return r.pool
}

// ============================================================================
// HIMS Analytics - Patients
// ============================================================================

// Patient represents a patient record from HIMS.
type Patient struct {
	PatientID              uuid.UUID `db:"patient_id"`
	ExternalPatientID      string    `db:"external_patient_id"`
	NameEN                 string    `db:"name_en"`
	NameAR                 *string   `db:"name_ar"`
	DateOfBirth            *time.Time `db:"date_of_birth"`
	Gender                 *string   `db:"gender"`
	Phone                  *string   `db:"phone"`
	Email                  *string   `db:"email"`
	AddressEN              *string   `db:"address_en"`
	AddressAR              *string   `db:"address_ar"`
	BloodGroup             *string   `db:"blood_group"`
	Nationality            *string   `db:"nationality"`
	NationalID             *string   `db:"national_id"`
	InsuranceProvider      *string   `db:"insurance_provider"`
	InsurancePolicyNumber  *string   `db:"insurance_policy_number"`
	EmergencyContactName   *string   `db:"emergency_contact_name"`
	EmergencyContactPhone  *string   `db:"emergency_contact_phone"`
	IsActive               bool      `db:"is_active"`
	Source                 string    `db:"_source"`
	SourceID               string    `db:"_source_id"`
}

// UpsertPatient inserts or updates a patient record.
func (r *Repo) UpsertPatient(ctx context.Context, p *Patient) error {
	query := `
		INSERT INTO hims_analytics.dim_patients (
			external_patient_id, name_en, name_ar, date_of_birth, gender,
			phone, email, address_en, address_ar, blood_group, nationality,
			national_id, insurance_provider, insurance_policy_number,
			emergency_contact_name, emergency_contact_phone, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			address_en = EXCLUDED.address_en,
			address_ar = EXCLUDED.address_ar,
			blood_group = EXCLUDED.blood_group,
			nationality = EXCLUDED.nationality,
			national_id = EXCLUDED.national_id,
			insurance_provider = EXCLUDED.insurance_provider,
			insurance_policy_number = EXCLUDED.insurance_policy_number,
			emergency_contact_name = EXCLUDED.emergency_contact_name,
			emergency_contact_phone = EXCLUDED.emergency_contact_phone,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING patient_id, _created_at
	`

	var patientID uuid.UUID
	var createdAt time.Time

	err := r.pool.QueryRow(ctx, query,
		p.ExternalPatientID, p.NameEN, p.NameAR, p.DateOfBirth, p.Gender,
		p.Phone, p.Email, p.AddressEN, p.AddressAR, p.BloodGroup, p.Nationality,
		p.NationalID, p.InsuranceProvider, p.InsurancePolicyNumber,
		p.EmergencyContactName, p.EmergencyContactPhone, p.IsActive,
		p.Source, p.SourceID,
	).Scan(&patientID, &createdAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert patient %s: %w", p.SourceID, err)
	}

	p.PatientID = patientID
	return nil
}

// BulkUpsertPatients inserts or updates multiple patient records in a batch.
func (r *Repo) BulkUpsertPatients(ctx context.Context, patients []*Patient) error {
	if len(patients) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO hims_analytics.dim_patients (
			external_patient_id, name_en, name_ar, date_of_birth, gender,
			phone, email, address_en, address_ar, blood_group, nationality,
			national_id, insurance_provider, insurance_policy_number,
			emergency_contact_name, emergency_contact_phone, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			address_en = EXCLUDED.address_en,
			address_ar = EXCLUDED.address_ar,
			blood_group = EXCLUDED.blood_group,
			nationality = EXCLUDED.nationality,
			national_id = EXCLUDED.national_id,
			insurance_provider = EXCLUDED.insurance_provider,
			insurance_policy_number = EXCLUDED.insurance_policy_number,
			emergency_contact_name = EXCLUDED.emergency_contact_name,
			emergency_contact_phone = EXCLUDED.emergency_contact_phone,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING patient_id
	`

	for _, p := range patients {
		batch.Queue(query,
			p.ExternalPatientID, p.NameEN, p.NameAR, p.DateOfBirth, p.Gender,
			p.Phone, p.Email, p.AddressEN, p.AddressAR, p.BloodGroup, p.Nationality,
			p.NationalID, p.InsuranceProvider, p.InsurancePolicyNumber,
			p.EmergencyContactName, p.EmergencyContactPhone, p.IsActive,
			p.Source, p.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, p := range patients {
		err := results.QueryRow().Scan(&p.PatientID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert patient %d/%d (%s): %w",
				i+1, len(patients), p.SourceID, err)
		}
	}

	r.logger.Debug("bulk upsert patients completed",
		slog.Int("count", len(patients)),
	)

	return nil
}

// ============================================================================
// HIMS Analytics - Doctors
// ============================================================================

// Doctor represents a doctor record from HIMS.
type Doctor struct {
	DoctorID        uuid.UUID  `db:"doctor_id"`
	ExternalDocID   string     `db:"external_doctor_id"`
	NameEN          string     `db:"name_en"`
	NameAR          *string    `db:"name_ar"`
	SpecialtyEN     *string    `db:"specialty_en"`
	SpecialtyAR     *string    `db:"specialty_ar"`
	DepartmentEN    *string    `db:"department_en"`
	DepartmentAR    *string    `db:"department_ar"`
	Qualification   *string    `db:"qualification"`
	LicenseNumber   *string    `db:"license_number"`
	Phone           *string    `db:"phone"`
	Email           *string    `db:"email"`
	ConsultationFee *float64   `db:"consultation_fee"`
	IsActive        bool       `db:"is_active"`
	Source          string     `db:"_source"`
	SourceID        string     `db:"_source_id"`
}

// UpsertDoctor inserts or updates a doctor record.
func (r *Repo) UpsertDoctor(ctx context.Context, d *Doctor) error {
	query := `
		INSERT INTO hims_analytics.dim_doctors (
			external_doctor_id, name_en, name_ar, specialty_en, specialty_ar,
			department_en, department_ar, qualification, license_number,
			phone, email, consultation_fee, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			specialty_en = EXCLUDED.specialty_en,
			specialty_ar = EXCLUDED.specialty_ar,
			department_en = EXCLUDED.department_en,
			department_ar = EXCLUDED.department_ar,
			qualification = EXCLUDED.qualification,
			license_number = EXCLUDED.license_number,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			consultation_fee = EXCLUDED.consultation_fee,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING doctor_id
	`

	err := r.pool.QueryRow(ctx, query,
		d.ExternalDocID, d.NameEN, d.NameAR, d.SpecialtyEN, d.SpecialtyAR,
		d.DepartmentEN, d.DepartmentAR, d.Qualification, d.LicenseNumber,
		d.Phone, d.Email, d.ConsultationFee, d.IsActive,
		d.Source, d.SourceID,
	).Scan(&d.DoctorID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert doctor %s: %w", d.SourceID, err)
	}

	return nil
}

// BulkUpsertDoctors inserts or updates multiple doctor records in a batch.
func (r *Repo) BulkUpsertDoctors(ctx context.Context, doctors []*Doctor) error {
	if len(doctors) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO hims_analytics.dim_doctors (
			external_doctor_id, name_en, name_ar, specialty_en, specialty_ar,
			department_en, department_ar, qualification, license_number,
			phone, email, consultation_fee, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			specialty_en = EXCLUDED.specialty_en,
			specialty_ar = EXCLUDED.specialty_ar,
			department_en = EXCLUDED.department_en,
			department_ar = EXCLUDED.department_ar,
			qualification = EXCLUDED.qualification,
			license_number = EXCLUDED.license_number,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			consultation_fee = EXCLUDED.consultation_fee,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING doctor_id
	`

	for _, d := range doctors {
		batch.Queue(query,
			d.ExternalDocID, d.NameEN, d.NameAR, d.SpecialtyEN, d.SpecialtyAR,
			d.DepartmentEN, d.DepartmentAR, d.Qualification, d.LicenseNumber,
			d.Phone, d.Email, d.ConsultationFee, d.IsActive,
			d.Source, d.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, d := range doctors {
		err := results.QueryRow().Scan(&d.DoctorID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert doctor %d/%d (%s): %w",
				i+1, len(doctors), d.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// HIMS Analytics - Drugs
// ============================================================================

// Drug represents a drug/medication record from HIMS.
type Drug struct {
	DrugID          uuid.UUID  `db:"drug_id"`
	ExternalDrugID  string     `db:"external_drug_id"`
	NameEN          string     `db:"name_en"`
	NameAR          *string    `db:"name_ar"`
	GenericNameEN   *string    `db:"generic_name_en"`
	GenericNameAR   *string    `db:"generic_name_ar"`
	CategoryEN      *string    `db:"category_en"`
	CategoryAR      *string    `db:"category_ar"`
	DosageForm      *string    `db:"dosage_form"`
	Strength        *string    `db:"strength"`
	Unit            *string    `db:"unit"`
	Manufacturer    *string    `db:"manufacturer"`
	UnitPrice       *float64   `db:"unit_price"`
	ReorderLevel    *int       `db:"reorder_level"`
	IsControlled    bool       `db:"is_controlled"`
	IsActive        bool       `db:"is_active"`
	Source          string     `db:"_source"`
	SourceID        string     `db:"_source_id"`
}

// UpsertDrug inserts or updates a drug record.
func (r *Repo) UpsertDrug(ctx context.Context, d *Drug) error {
	query := `
		INSERT INTO hims_analytics.dim_drugs (
			external_drug_id, name_en, name_ar, generic_name_en, generic_name_ar,
			category_en, category_ar, dosage_form, strength, unit,
			manufacturer, unit_price, reorder_level, is_controlled, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			generic_name_en = EXCLUDED.generic_name_en,
			generic_name_ar = EXCLUDED.generic_name_ar,
			category_en = EXCLUDED.category_en,
			category_ar = EXCLUDED.category_ar,
			dosage_form = EXCLUDED.dosage_form,
			strength = EXCLUDED.strength,
			unit = EXCLUDED.unit,
			manufacturer = EXCLUDED.manufacturer,
			unit_price = EXCLUDED.unit_price,
			reorder_level = EXCLUDED.reorder_level,
			is_controlled = EXCLUDED.is_controlled,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING drug_id
	`

	err := r.pool.QueryRow(ctx, query,
		d.ExternalDrugID, d.NameEN, d.NameAR, d.GenericNameEN, d.GenericNameAR,
		d.CategoryEN, d.CategoryAR, d.DosageForm, d.Strength, d.Unit,
		d.Manufacturer, d.UnitPrice, d.ReorderLevel, d.IsControlled, d.IsActive,
		d.Source, d.SourceID,
	).Scan(&d.DrugID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert drug %s: %w", d.SourceID, err)
	}

	return nil
}

// BulkUpsertDrugs inserts or updates multiple drug records in a batch.
func (r *Repo) BulkUpsertDrugs(ctx context.Context, drugs []*Drug) error {
	if len(drugs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO hims_analytics.dim_drugs (
			external_drug_id, name_en, name_ar, generic_name_en, generic_name_ar,
			category_en, category_ar, dosage_form, strength, unit,
			manufacturer, unit_price, reorder_level, is_controlled, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			generic_name_en = EXCLUDED.generic_name_en,
			generic_name_ar = EXCLUDED.generic_name_ar,
			category_en = EXCLUDED.category_en,
			category_ar = EXCLUDED.category_ar,
			dosage_form = EXCLUDED.dosage_form,
			strength = EXCLUDED.strength,
			unit = EXCLUDED.unit,
			manufacturer = EXCLUDED.manufacturer,
			unit_price = EXCLUDED.unit_price,
			reorder_level = EXCLUDED.reorder_level,
			is_controlled = EXCLUDED.is_controlled,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING drug_id
	`

	for _, d := range drugs {
		batch.Queue(query,
			d.ExternalDrugID, d.NameEN, d.NameAR, d.GenericNameEN, d.GenericNameAR,
			d.CategoryEN, d.CategoryAR, d.DosageForm, d.Strength, d.Unit,
			d.Manufacturer, d.UnitPrice, d.ReorderLevel, d.IsControlled, d.IsActive,
			d.Source, d.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, d := range drugs {
		err := results.QueryRow().Scan(&d.DrugID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert drug %d/%d (%s): %w",
				i+1, len(drugs), d.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// HIMS Analytics - Appointments
// ============================================================================

// Appointment represents an appointment record from HIMS.
type Appointment struct {
	ApptID              uuid.UUID  `db:"appt_id"`
	ExternalApptID      string     `db:"external_appt_id"`
	PatientID           uuid.UUID  `db:"patient_id"`
	DoctorID            uuid.UUID  `db:"doctor_id"`
	DepartmentID        *uuid.UUID `db:"department_id"`
	ApptDate            time.Time  `db:"appt_date"`
	ApptTime            *string    `db:"appt_time"`
	ApptDatetime        *time.Time `db:"appt_datetime"`
	Status              string     `db:"status"`
	ApptType            *string    `db:"appt_type"`
	DurationMinutes     *int       `db:"duration_minutes"`
	ChiefComplaint      *string    `db:"chief_complaint"`
	DiagnosisCode       *string    `db:"diagnosis_code"`
	DiagnosisDescription *string   `db:"diagnosis_description"`
	BillingID           *uuid.UUID `db:"billing_id"`
	Notes               *string    `db:"notes"`
	IsWalkIn            bool       `db:"is_walk_in"`
	CancellationReason  *string    `db:"cancellation_reason"`
	CancelledAt         *time.Time `db:"cancelled_at"`
	Source              string     `db:"_source"`
	SourceID            string     `db:"_source_id"`
}

// UpsertAppointment inserts or updates an appointment record.
func (r *Repo) UpsertAppointment(ctx context.Context, a *Appointment) error {
	query := `
		INSERT INTO hims_analytics.fact_appointments (
			external_appt_id, patient_id, doctor_id, department_id,
			appt_date, appt_time, appt_datetime, status, appt_type,
			duration_minutes, chief_complaint, diagnosis_code, diagnosis_description,
			billing_id, notes, is_walk_in, cancellation_reason, cancelled_at,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			patient_id = EXCLUDED.patient_id,
			doctor_id = EXCLUDED.doctor_id,
			department_id = EXCLUDED.department_id,
			appt_date = EXCLUDED.appt_date,
			appt_time = EXCLUDED.appt_time,
			appt_datetime = EXCLUDED.appt_datetime,
			status = EXCLUDED.status,
			appt_type = EXCLUDED.appt_type,
			duration_minutes = EXCLUDED.duration_minutes,
			chief_complaint = EXCLUDED.chief_complaint,
			diagnosis_code = EXCLUDED.diagnosis_code,
			diagnosis_description = EXCLUDED.diagnosis_description,
			billing_id = EXCLUDED.billing_id,
			notes = EXCLUDED.notes,
			is_walk_in = EXCLUDED.is_walk_in,
			cancellation_reason = EXCLUDED.cancellation_reason,
			cancelled_at = EXCLUDED.cancelled_at,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING appt_id
	`

	err := r.pool.QueryRow(ctx, query,
		a.ExternalApptID, a.PatientID, a.DoctorID, a.DepartmentID,
		a.ApptDate, a.ApptTime, a.ApptDatetime, a.Status, a.ApptType,
		a.DurationMinutes, a.ChiefComplaint, a.DiagnosisCode, a.DiagnosisDescription,
		a.BillingID, a.Notes, a.IsWalkIn, a.CancellationReason, a.CancelledAt,
		a.Source, a.SourceID,
	).Scan(&a.ApptID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert appointment %s: %w", a.SourceID, err)
	}

	return nil
}

// ============================================================================
// HIMS Analytics - Billing
// ============================================================================

// Billing represents a billing record from HIMS.
type Billing struct {
	BillID              uuid.UUID  `db:"bill_id"`
	ExternalBillID      string     `db:"external_bill_id"`
	PatientID           uuid.UUID  `db:"patient_id"`
	ApptID              *uuid.UUID `db:"appt_id"`
	BillDate            time.Time  `db:"bill_date"`
	BillDatetime        *time.Time `db:"bill_datetime"`
	SubtotalAmount      float64    `db:"subtotal_amount"`
	DiscountAmount      float64    `db:"discount_amount"`
	TaxAmount           float64    `db:"tax_amount"`
	TotalAmount         float64    `db:"total_amount"`
	PaidAmount          float64    `db:"paid_amount"`
	PaymentMode         *string    `db:"payment_mode"`
	PaymentStatus       string     `db:"payment_status"`
	InsuranceClaimID    *string    `db:"insurance_claim_id"`
	InsuranceAmount     float64    `db:"insurance_amount"`
	DepartmentEN        *string    `db:"department_en"`
	DepartmentAR        *string    `db:"department_ar"`
	BillType            *string    `db:"bill_type"`
	ReceiptNumber       *string    `db:"receipt_number"`
	Notes               *string    `db:"notes"`
	Source              string     `db:"_source"`
	SourceID            string     `db:"_source_id"`
}

// UpsertBilling inserts or updates a billing record.
func (r *Repo) UpsertBilling(ctx context.Context, b *Billing) error {
	query := `
		INSERT INTO hims_analytics.fact_billing (
			external_bill_id, patient_id, appt_id, bill_date, bill_datetime,
			subtotal_amount, discount_amount, tax_amount, total_amount, paid_amount,
			payment_mode, payment_status, insurance_claim_id, insurance_amount,
			department_en, department_ar, bill_type, receipt_number, notes,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			patient_id = EXCLUDED.patient_id,
			appt_id = EXCLUDED.appt_id,
			bill_date = EXCLUDED.bill_date,
			bill_datetime = EXCLUDED.bill_datetime,
			subtotal_amount = EXCLUDED.subtotal_amount,
			discount_amount = EXCLUDED.discount_amount,
			tax_amount = EXCLUDED.tax_amount,
			total_amount = EXCLUDED.total_amount,
			paid_amount = EXCLUDED.paid_amount,
			payment_mode = EXCLUDED.payment_mode,
			payment_status = EXCLUDED.payment_status,
			insurance_claim_id = EXCLUDED.insurance_claim_id,
			insurance_amount = EXCLUDED.insurance_amount,
			department_en = EXCLUDED.department_en,
			department_ar = EXCLUDED.department_ar,
			bill_type = EXCLUDED.bill_type,
			receipt_number = EXCLUDED.receipt_number,
			notes = EXCLUDED.notes,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING bill_id
	`

	err := r.pool.QueryRow(ctx, query,
		b.ExternalBillID, b.PatientID, b.ApptID, b.BillDate, b.BillDatetime,
		b.SubtotalAmount, b.DiscountAmount, b.TaxAmount, b.TotalAmount, b.PaidAmount,
		b.PaymentMode, b.PaymentStatus, b.InsuranceClaimID, b.InsuranceAmount,
		b.DepartmentEN, b.DepartmentAR, b.BillType, b.ReceiptNumber, b.Notes,
		b.Source, b.SourceID,
	).Scan(&b.BillID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert billing %s: %w", b.SourceID, err)
	}

	return nil
}

// ============================================================================
// HIMS Analytics - Pharmacy Dispensations
// ============================================================================

// PharmacyDispensation represents a pharmacy dispensation record from HIMS.
type PharmacyDispensation struct {
	DispID              uuid.UUID  `db:"disp_id"`
	ExternalDispID      string     `db:"external_disp_id"`
	DrugID              uuid.UUID  `db:"drug_id"`
	PatientID           uuid.UUID  `db:"patient_id"`
	DoctorID            *uuid.UUID `db:"doctor_id"`
	BillID              *uuid.UUID `db:"bill_id"`
	PrescriptionID      *string    `db:"prescription_id"`
	DispDate            time.Time  `db:"disp_date"`
	DispDatetime        *time.Time `db:"disp_datetime"`
	Quantity            int        `db:"quantity"`
	Unit                *string    `db:"unit"`
	DosageInstructions  *string    `db:"dosage_instructions"`
	DaysSupply          *int       `db:"days_supply"`
	UnitPrice           float64    `db:"unit_price"`
	DiscountAmount      float64    `db:"discount_amount"`
	TaxAmount           float64    `db:"tax_amount"`
	TotalAmount         float64    `db:"total_amount"`
	BatchNumber         *string    `db:"batch_number"`
	ExpiryDate          *time.Time `db:"expiry_date"`
	IsSubstituted       bool       `db:"is_substituted"`
	OriginalDrugID      *uuid.UUID `db:"original_drug_id"`
	Notes               *string    `db:"notes"`
	Source              string     `db:"_source"`
	SourceID            string     `db:"_source_id"`
}

// UpsertPharmacyDispensation inserts or updates a pharmacy dispensation record.
func (r *Repo) UpsertPharmacyDispensation(ctx context.Context, p *PharmacyDispensation) error {
	query := `
		INSERT INTO hims_analytics.fact_pharmacy_dispensations (
			external_disp_id, drug_id, patient_id, doctor_id, bill_id,
			prescription_id, disp_date, disp_datetime, quantity, unit,
			dosage_instructions, days_supply, unit_price, discount_amount, tax_amount,
			total_amount, batch_number, expiry_date, is_substituted, original_drug_id, notes,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			drug_id = EXCLUDED.drug_id,
			patient_id = EXCLUDED.patient_id,
			doctor_id = EXCLUDED.doctor_id,
			bill_id = EXCLUDED.bill_id,
			prescription_id = EXCLUDED.prescription_id,
			disp_date = EXCLUDED.disp_date,
			disp_datetime = EXCLUDED.disp_datetime,
			quantity = EXCLUDED.quantity,
			unit = EXCLUDED.unit,
			dosage_instructions = EXCLUDED.dosage_instructions,
			days_supply = EXCLUDED.days_supply,
			unit_price = EXCLUDED.unit_price,
			discount_amount = EXCLUDED.discount_amount,
			tax_amount = EXCLUDED.tax_amount,
			total_amount = EXCLUDED.total_amount,
			batch_number = EXCLUDED.batch_number,
			expiry_date = EXCLUDED.expiry_date,
			is_substituted = EXCLUDED.is_substituted,
			original_drug_id = EXCLUDED.original_drug_id,
			notes = EXCLUDED.notes,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING disp_id
	`

	err := r.pool.QueryRow(ctx, query,
		p.ExternalDispID, p.DrugID, p.PatientID, p.DoctorID, p.BillID,
		p.PrescriptionID, p.DispDate, p.DispDatetime, p.Quantity, p.Unit,
		p.DosageInstructions, p.DaysSupply, p.UnitPrice, p.DiscountAmount, p.TaxAmount,
		p.TotalAmount, p.BatchNumber, p.ExpiryDate, p.IsSubstituted, p.OriginalDrugID, p.Notes,
		p.Source, p.SourceID,
	).Scan(&p.DispID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert pharmacy dispensation %s: %w", p.SourceID, err)
	}

	return nil
}

// ============================================================================
// TALLY Analytics - Ledgers
// ============================================================================

// Ledger represents a Tally ledger (account) record.
type Ledger struct {
	LedgerID           uuid.UUID  `db:"ledger_id"`
	ExternalLedgerID   string     `db:"external_ledger_id"`
	LedgerName         string     `db:"ledger_name"`
	LedgerNameAR       *string    `db:"ledger_name_ar"`
	LedgerGroup        string     `db:"ledger_group"`
	ParentGroup        *string    `db:"parent_group"`
	LedgerType         *string    `db:"ledger_type"`
	OpeningBalance     float64    `db:"opening_balance"`
	ClosingBalance     float64    `db:"closing_balance"`
	Currency           string     `db:"currency"`
	IsBankAccount      bool       `db:"is_bank_account"`
	BankName           *string    `db:"bank_name"`
	BankAccountNumber  *string    `db:"bank_account_number"`
	IFSCCode           *string    `db:"ifsc_code"`
	GSTRegistration    *string    `db:"gst_registration"`
	PANNumber          *string    `db:"pan_number"`
	CreditPeriodDays   *int       `db:"credit_period_days"`
	CreditLimit        *float64   `db:"credit_limit"`
	IsActive           bool       `db:"is_active"`
	Source             string     `db:"_source"`
	SourceID           string     `db:"_source_id"`
}

// UpsertLedger inserts or updates a ledger record.
func (r *Repo) UpsertLedger(ctx context.Context, l *Ledger) error {
	query := `
		INSERT INTO tally_analytics.dim_ledgers (
			external_ledger_id, ledger_name, ledger_name_ar, ledger_group, parent_group,
			ledger_type, opening_balance, closing_balance, currency, is_bank_account,
			bank_name, bank_account_number, ifsc_code, gst_registration, pan_number,
			credit_period_days, credit_limit, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			ledger_name = EXCLUDED.ledger_name,
			ledger_name_ar = EXCLUDED.ledger_name_ar,
			ledger_group = EXCLUDED.ledger_group,
			parent_group = EXCLUDED.parent_group,
			ledger_type = EXCLUDED.ledger_type,
			opening_balance = EXCLUDED.opening_balance,
			closing_balance = EXCLUDED.closing_balance,
			currency = EXCLUDED.currency,
			is_bank_account = EXCLUDED.is_bank_account,
			bank_name = EXCLUDED.bank_name,
			bank_account_number = EXCLUDED.bank_account_number,
			ifsc_code = EXCLUDED.ifsc_code,
			gst_registration = EXCLUDED.gst_registration,
			pan_number = EXCLUDED.pan_number,
			credit_period_days = EXCLUDED.credit_period_days,
			credit_limit = EXCLUDED.credit_limit,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING ledger_id
	`

	err := r.pool.QueryRow(ctx, query,
		l.ExternalLedgerID, l.LedgerName, l.LedgerNameAR, l.LedgerGroup, l.ParentGroup,
		l.LedgerType, l.OpeningBalance, l.ClosingBalance, l.Currency, l.IsBankAccount,
		l.BankName, l.BankAccountNumber, l.IFSCCode, l.GSTRegistration, l.PANNumber,
		l.CreditPeriodDays, l.CreditLimit, l.IsActive,
		l.Source, l.SourceID,
	).Scan(&l.LedgerID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert ledger %s: %w", l.SourceID, err)
	}

	return nil
}

// BulkUpsertLedgers inserts or updates multiple ledger records in a batch.
func (r *Repo) BulkUpsertLedgers(ctx context.Context, ledgers []*Ledger) error {
	if len(ledgers) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO tally_analytics.dim_ledgers (
			external_ledger_id, ledger_name, ledger_name_ar, ledger_group, parent_group,
			ledger_type, opening_balance, closing_balance, currency, is_bank_account,
			bank_name, bank_account_number, ifsc_code, gst_registration, pan_number,
			credit_period_days, credit_limit, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			ledger_name = EXCLUDED.ledger_name,
			ledger_name_ar = EXCLUDED.ledger_name_ar,
			ledger_group = EXCLUDED.ledger_group,
			parent_group = EXCLUDED.parent_group,
			ledger_type = EXCLUDED.ledger_type,
			opening_balance = EXCLUDED.opening_balance,
			closing_balance = EXCLUDED.closing_balance,
			currency = EXCLUDED.currency,
			is_bank_account = EXCLUDED.is_bank_account,
			bank_name = EXCLUDED.bank_name,
			bank_account_number = EXCLUDED.bank_account_number,
			ifsc_code = EXCLUDED.ifsc_code,
			gst_registration = EXCLUDED.gst_registration,
			pan_number = EXCLUDED.pan_number,
			credit_period_days = EXCLUDED.credit_period_days,
			credit_limit = EXCLUDED.credit_limit,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING ledger_id
	`

	for _, l := range ledgers {
		batch.Queue(query,
			l.ExternalLedgerID, l.LedgerName, l.LedgerNameAR, l.LedgerGroup, l.ParentGroup,
			l.LedgerType, l.OpeningBalance, l.ClosingBalance, l.Currency, l.IsBankAccount,
			l.BankName, l.BankAccountNumber, l.IFSCCode, l.GSTRegistration, l.PANNumber,
			l.CreditPeriodDays, l.CreditLimit, l.IsActive,
			l.Source, l.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, l := range ledgers {
		err := results.QueryRow().Scan(&l.LedgerID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert ledger %d/%d (%s): %w",
				i+1, len(ledgers), l.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// TALLY Analytics - Inventory Items
// ============================================================================

// InventoryItem represents a Tally stock item record.
type InventoryItem struct {
	ItemID              uuid.UUID  `db:"item_id"`
	ExternalItemID      string     `db:"external_item_id"`
	NameEN              string     `db:"name_en"`
	NameAR              *string    `db:"name_ar"`
	PartNumber          *string    `db:"part_number"`
	Category            *string    `db:"category"`
	SubCategory         *string    `db:"sub_category"`
	StockGroup          *string    `db:"stock_group"`
	Unit                *string    `db:"unit"`
	AlternateUnit       *string    `db:"alternate_unit"`
	ConversionFactor    *float64   `db:"conversion_factor"`
	GSTRate             *float64   `db:"gst_rate"`
	HSNCode             *string    `db:"hsn_code"`
	PurchasePrice       *float64   `db:"purchase_price"`
	SellingPrice        *float64   `db:"selling_price"`
	MRP                 *float64   `db:"mrp"`
	ReorderLevel        *int       `db:"reorder_level"`
	MinimumOrderQty     *int       `db:"minimum_order_qty"`
	IsBatchWise         bool       `db:"is_batch_wise"`
	MaintainExpiry      bool       `db:"maintain_expiry"`
	IsActive            bool       `db:"is_active"`
	Source              string     `db:"_source"`
	SourceID            string     `db:"_source_id"`
}

// UpsertInventoryItem inserts or updates an inventory item record.
func (r *Repo) UpsertInventoryItem(ctx context.Context, i *InventoryItem) error {
	query := `
		INSERT INTO tally_analytics.dim_inventory_items (
			external_item_id, name_en, name_ar, part_number, category, sub_category,
			stock_group, unit, alternate_unit, conversion_factor, gst_rate, hsn_code,
			purchase_price, selling_price, mrp, reorder_level, minimum_order_qty,
			is_batch_wise, maintain_expiry, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			part_number = EXCLUDED.part_number,
			category = EXCLUDED.category,
			sub_category = EXCLUDED.sub_category,
			stock_group = EXCLUDED.stock_group,
			unit = EXCLUDED.unit,
			alternate_unit = EXCLUDED.alternate_unit,
			conversion_factor = EXCLUDED.conversion_factor,
			gst_rate = EXCLUDED.gst_rate,
			hsn_code = EXCLUDED.hsn_code,
			purchase_price = EXCLUDED.purchase_price,
			selling_price = EXCLUDED.selling_price,
			mrp = EXCLUDED.mrp,
			reorder_level = EXCLUDED.reorder_level,
			minimum_order_qty = EXCLUDED.minimum_order_qty,
			is_batch_wise = EXCLUDED.is_batch_wise,
			maintain_expiry = EXCLUDED.maintain_expiry,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING item_id
	`

	err := r.pool.QueryRow(ctx, query,
		i.ExternalItemID, i.NameEN, i.NameAR, i.PartNumber, i.Category, i.SubCategory,
		i.StockGroup, i.Unit, i.AlternateUnit, i.ConversionFactor, i.GSTRate, i.HSNCode,
		i.PurchasePrice, i.SellingPrice, i.MRP, i.ReorderLevel, i.MinimumOrderQty,
		i.IsBatchWise, i.MaintainExpiry, i.IsActive,
		i.Source, i.SourceID,
	).Scan(&i.ItemID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert inventory item %s: %w", i.SourceID, err)
	}

	return nil
}

// BulkUpsertInventoryItems inserts or updates multiple inventory item records in a batch.
func (r *Repo) BulkUpsertInventoryItems(ctx context.Context, items []*InventoryItem) error {
	if len(items) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO tally_analytics.dim_inventory_items (
			external_item_id, name_en, name_ar, part_number, category, sub_category,
			stock_group, unit, alternate_unit, conversion_factor, gst_rate, hsn_code,
			purchase_price, selling_price, mrp, reorder_level, minimum_order_qty,
			is_batch_wise, maintain_expiry, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name_en = EXCLUDED.name_en,
			name_ar = EXCLUDED.name_ar,
			part_number = EXCLUDED.part_number,
			category = EXCLUDED.category,
			sub_category = EXCLUDED.sub_category,
			stock_group = EXCLUDED.stock_group,
			unit = EXCLUDED.unit,
			alternate_unit = EXCLUDED.alternate_unit,
			conversion_factor = EXCLUDED.conversion_factor,
			gst_rate = EXCLUDED.gst_rate,
			hsn_code = EXCLUDED.hsn_code,
			purchase_price = EXCLUDED.purchase_price,
			selling_price = EXCLUDED.selling_price,
			mrp = EXCLUDED.mrp,
			reorder_level = EXCLUDED.reorder_level,
			minimum_order_qty = EXCLUDED.minimum_order_qty,
			is_batch_wise = EXCLUDED.is_batch_wise,
			maintain_expiry = EXCLUDED.maintain_expiry,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING item_id
	`

	for _, i := range items {
		batch.Queue(query,
			i.ExternalItemID, i.NameEN, i.NameAR, i.PartNumber, i.Category, i.SubCategory,
			i.StockGroup, i.Unit, i.AlternateUnit, i.ConversionFactor, i.GSTRate, i.HSNCode,
			i.PurchasePrice, i.SellingPrice, i.MRP, i.ReorderLevel, i.MinimumOrderQty,
			i.IsBatchWise, i.MaintainExpiry, i.IsActive,
			i.Source, i.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, item := range items {
		err := results.QueryRow().Scan(&item.ItemID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert inventory item %d/%d (%s): %w",
				i+1, len(items), item.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// TALLY Analytics - Vouchers
// ============================================================================

// Voucher represents a Tally voucher (transaction) record.
type Voucher struct {
	VoucherID          uuid.UUID  `db:"voucher_id"`
	ExternalVoucherID  string     `db:"external_voucher_id"`
	VoucherNumber      string     `db:"voucher_number"`
	VoucherType        string     `db:"voucher_type"`
	VoucherDate        time.Time  `db:"voucher_date"`
	VoucherDatetime    *time.Time `db:"voucher_datetime"`
	LedgerID           uuid.UUID  `db:"ledger_id"`
	ContraLedgerID     *uuid.UUID `db:"contra_ledger_id"`
	CostCentreID       *uuid.UUID `db:"cost_centre_id"`
	Amount             float64    `db:"amount"`
	IsDebit            bool       `db:"is_debit"`
	Currency           string     `db:"currency"`
	ExchangeRate       float64    `db:"exchange_rate"`
	BaseCurrencyAmount *float64   `db:"base_currency_amount"`
	Narration          *string    `db:"narration"`
	ReferenceNumber    *string    `db:"reference_number"`
	ReferenceDate      *time.Time `db:"reference_date"`
	PartyName          *string    `db:"party_name"`
	BillNumber         *string    `db:"bill_number"`
	BillDate           *time.Time `db:"bill_date"`
	DueDate            *time.Time `db:"due_date"`
	InstrumentNumber   *string    `db:"instrument_number"`
	InstrumentDate     *time.Time `db:"instrument_date"`
	BankName           *string    `db:"bank_name"`
	GSTRegistration    *string    `db:"gst_registration"`
	InvoiceNumber      *string    `db:"invoice_number"`
	IsCancelled        bool       `db:"is_cancelled"`
	CancelledDate      *time.Time `db:"cancelled_date"`
	CancellationReason *string    `db:"cancellation_reason"`
	IsOptional         bool       `db:"is_optional"`
	HasInventory       bool       `db:"has_inventory"`
	Source             string     `db:"_source"`
	SourceID           string     `db:"_source_id"`
}

// UpsertVoucher inserts or updates a voucher record.
func (r *Repo) UpsertVoucher(ctx context.Context, v *Voucher) error {
	query := `
		INSERT INTO tally_analytics.fact_vouchers (
			external_voucher_id, voucher_number, voucher_type, voucher_date, voucher_datetime,
			ledger_id, contra_ledger_id, cost_centre_id, amount, is_debit,
			currency, exchange_rate, base_currency_amount, narration, reference_number,
			reference_date, party_name, bill_number, bill_date, due_date,
			instrument_number, instrument_date, bank_name, gst_registration, invoice_number,
			is_cancelled, cancelled_date, cancellation_reason, is_optional, has_inventory,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			voucher_number = EXCLUDED.voucher_number,
			voucher_type = EXCLUDED.voucher_type,
			voucher_date = EXCLUDED.voucher_date,
			voucher_datetime = EXCLUDED.voucher_datetime,
			ledger_id = EXCLUDED.ledger_id,
			contra_ledger_id = EXCLUDED.contra_ledger_id,
			cost_centre_id = EXCLUDED.cost_centre_id,
			amount = EXCLUDED.amount,
			is_debit = EXCLUDED.is_debit,
			currency = EXCLUDED.currency,
			exchange_rate = EXCLUDED.exchange_rate,
			base_currency_amount = EXCLUDED.base_currency_amount,
			narration = EXCLUDED.narration,
			reference_number = EXCLUDED.reference_number,
			reference_date = EXCLUDED.reference_date,
			party_name = EXCLUDED.party_name,
			bill_number = EXCLUDED.bill_number,
			bill_date = EXCLUDED.bill_date,
			due_date = EXCLUDED.due_date,
			instrument_number = EXCLUDED.instrument_number,
			instrument_date = EXCLUDED.instrument_date,
			bank_name = EXCLUDED.bank_name,
			gst_registration = EXCLUDED.gst_registration,
			invoice_number = EXCLUDED.invoice_number,
			is_cancelled = EXCLUDED.is_cancelled,
			cancelled_date = EXCLUDED.cancelled_date,
			cancellation_reason = EXCLUDED.cancellation_reason,
			is_optional = EXCLUDED.is_optional,
			has_inventory = EXCLUDED.has_inventory,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING voucher_id
	`

	err := r.pool.QueryRow(ctx, query,
		v.ExternalVoucherID, v.VoucherNumber, v.VoucherType, v.VoucherDate, v.VoucherDatetime,
		v.LedgerID, v.ContraLedgerID, v.CostCentreID, v.Amount, v.IsDebit,
		v.Currency, v.ExchangeRate, v.BaseCurrencyAmount, v.Narration, v.ReferenceNumber,
		v.ReferenceDate, v.PartyName, v.BillNumber, v.BillDate, v.DueDate,
		v.InstrumentNumber, v.InstrumentDate, v.BankName, v.GSTRegistration, v.InvoiceNumber,
		v.IsCancelled, v.CancelledDate, v.CancellationReason, v.IsOptional, v.HasInventory,
		v.Source, v.SourceID,
	).Scan(&v.VoucherID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert voucher %s: %w", v.SourceID, err)
	}

	return nil
}

// BulkUpsertVouchers inserts or updates multiple voucher records in a batch.
func (r *Repo) BulkUpsertVouchers(ctx context.Context, vouchers []*Voucher) error {
	if len(vouchers) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO tally_analytics.fact_vouchers (
			external_voucher_id, voucher_number, voucher_type, voucher_date, voucher_datetime,
			ledger_id, contra_ledger_id, cost_centre_id, amount, is_debit,
			currency, exchange_rate, base_currency_amount, narration, reference_number,
			reference_date, party_name, bill_number, bill_date, due_date,
			instrument_number, instrument_date, bank_name, gst_registration, invoice_number,
			is_cancelled, cancelled_date, cancellation_reason, is_optional, has_inventory,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			voucher_number = EXCLUDED.voucher_number,
			voucher_type = EXCLUDED.voucher_type,
			voucher_date = EXCLUDED.voucher_date,
			voucher_datetime = EXCLUDED.voucher_datetime,
			ledger_id = EXCLUDED.ledger_id,
			contra_ledger_id = EXCLUDED.contra_ledger_id,
			cost_centre_id = EXCLUDED.cost_centre_id,
			amount = EXCLUDED.amount,
			is_debit = EXCLUDED.is_debit,
			currency = EXCLUDED.currency,
			exchange_rate = EXCLUDED.exchange_rate,
			base_currency_amount = EXCLUDED.base_currency_amount,
			narration = EXCLUDED.narration,
			reference_number = EXCLUDED.reference_number,
			reference_date = EXCLUDED.reference_date,
			party_name = EXCLUDED.party_name,
			bill_number = EXCLUDED.bill_number,
			bill_date = EXCLUDED.bill_date,
			due_date = EXCLUDED.due_date,
			instrument_number = EXCLUDED.instrument_number,
			instrument_date = EXCLUDED.instrument_date,
			bank_name = EXCLUDED.bank_name,
			gst_registration = EXCLUDED.gst_registration,
			invoice_number = EXCLUDED.invoice_number,
			is_cancelled = EXCLUDED.is_cancelled,
			cancelled_date = EXCLUDED.cancelled_date,
			cancellation_reason = EXCLUDED.cancellation_reason,
			is_optional = EXCLUDED.is_optional,
			has_inventory = EXCLUDED.has_inventory,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING voucher_id
	`

	for _, v := range vouchers {
		batch.Queue(query,
			v.ExternalVoucherID, v.VoucherNumber, v.VoucherType, v.VoucherDate, v.VoucherDatetime,
			v.LedgerID, v.ContraLedgerID, v.CostCentreID, v.Amount, v.IsDebit,
			v.Currency, v.ExchangeRate, v.BaseCurrencyAmount, v.Narration, v.ReferenceNumber,
			v.ReferenceDate, v.PartyName, v.BillNumber, v.BillDate, v.DueDate,
			v.InstrumentNumber, v.InstrumentDate, v.BankName, v.GSTRegistration, v.InvoiceNumber,
			v.IsCancelled, v.CancelledDate, v.CancellationReason, v.IsOptional, v.HasInventory,
			v.Source, v.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, v := range vouchers {
		err := results.QueryRow().Scan(&v.VoucherID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert voucher %d/%d (%s): %w",
				i+1, len(vouchers), v.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// TALLY Analytics - Cost Centres
// ============================================================================

// CostCentre represents a Tally cost centre record.
type CostCentre struct {
	CCID              uuid.UUID  `db:"cc_id"`
	ExternalCCID      string     `db:"external_cc_id"`
	Name              string     `db:"name"`
	NameAR            *string    `db:"name_ar"`
	Code              *string    `db:"code"`
	ParentCCID        *uuid.UUID `db:"parent_cc_id"`
	Category          *string    `db:"category"`
	IsRevenueCentre   bool       `db:"is_revenue_centre"`
	IsActive          bool       `db:"is_active"`
	Source            string     `db:"_source"`
	SourceID          string     `db:"_source_id"`
}

// UpsertCostCentre inserts or updates a cost centre record.
func (r *Repo) UpsertCostCentre(ctx context.Context, c *CostCentre) error {
	query := `
		INSERT INTO tally_analytics.dim_cost_centres (
			external_cc_id, name, name_ar, code, parent_cc_id,
			category, is_revenue_centre, is_active,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			name = EXCLUDED.name,
			name_ar = EXCLUDED.name_ar,
			code = EXCLUDED.code,
			parent_cc_id = EXCLUDED.parent_cc_id,
			category = EXCLUDED.category,
			is_revenue_centre = EXCLUDED.is_revenue_centre,
			is_active = EXCLUDED.is_active,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING cc_id
	`

	err := r.pool.QueryRow(ctx, query,
		c.ExternalCCID, c.Name, c.NameAR, c.Code, c.ParentCCID,
		c.Category, c.IsRevenueCentre, c.IsActive,
		c.Source, c.SourceID,
	).Scan(&c.CCID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert cost centre %s: %w", c.SourceID, err)
	}

	return nil
}

// ============================================================================
// TALLY Analytics - Stock Movements
// ============================================================================

// StockMovement represents a Tally stock movement record.
type StockMovement struct {
	MovementID         uuid.UUID  `db:"movement_id"`
	ExternalMovementID string     `db:"external_movement_id"`
	ItemID             uuid.UUID  `db:"item_id"`
	VoucherID          *uuid.UUID `db:"voucher_id"`
	MovementType       string     `db:"movement_type"`
	MovementDate       time.Time  `db:"movement_date"`
	MovementDatetime   *time.Time `db:"movement_datetime"`
	QtyIn              float64    `db:"qty_in"`
	QtyOut             float64    `db:"qty_out"`
	Unit               *string    `db:"unit"`
	Rate               *float64   `db:"rate"`
	Value              *float64   `db:"value"`
	BatchNumber        *string    `db:"batch_number"`
	ExpiryDate         *time.Time `db:"expiry_date"`
	GodownName         *string    `db:"godown_name"`
	ClosingStock       *float64   `db:"closing_stock"`
	ClosingValue       *float64   `db:"closing_value"`
	Source             string     `db:"_source"`
	SourceID           string     `db:"_source_id"`
}

// UpsertStockMovement inserts or updates a stock movement record.
func (r *Repo) UpsertStockMovement(ctx context.Context, s *StockMovement) error {
	query := `
		INSERT INTO tally_analytics.fact_stock_movements (
			external_movement_id, item_id, voucher_id, movement_type,
			movement_date, movement_datetime, qty_in, qty_out, unit,
			rate, value, batch_number, expiry_date, godown_name,
			closing_stock, closing_value,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			item_id = EXCLUDED.item_id,
			voucher_id = EXCLUDED.voucher_id,
			movement_type = EXCLUDED.movement_type,
			movement_date = EXCLUDED.movement_date,
			movement_datetime = EXCLUDED.movement_datetime,
			qty_in = EXCLUDED.qty_in,
			qty_out = EXCLUDED.qty_out,
			unit = EXCLUDED.unit,
			rate = EXCLUDED.rate,
			value = EXCLUDED.value,
			batch_number = EXCLUDED.batch_number,
			expiry_date = EXCLUDED.expiry_date,
			godown_name = EXCLUDED.godown_name,
			closing_stock = EXCLUDED.closing_stock,
			closing_value = EXCLUDED.closing_value,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING movement_id
	`

	err := r.pool.QueryRow(ctx, query,
		s.ExternalMovementID, s.ItemID, s.VoucherID, s.MovementType,
		s.MovementDate, s.MovementDatetime, s.QtyIn, s.QtyOut, s.Unit,
		s.Rate, s.Value, s.BatchNumber, s.ExpiryDate, s.GodownName,
		s.ClosingStock, s.ClosingValue,
		s.Source, s.SourceID,
	).Scan(&s.MovementID)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert stock movement %s: %w", s.SourceID, err)
	}

	return nil
}

// BulkUpsertStockMovements inserts or updates multiple stock movement records in a batch.
func (r *Repo) BulkUpsertStockMovements(ctx context.Context, movements []*StockMovement) error {
	if len(movements) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO tally_analytics.fact_stock_movements (
			external_movement_id, item_id, voucher_id, movement_type,
			movement_date, movement_datetime, qty_in, qty_out, unit,
			rate, value, batch_number, expiry_date, godown_name,
			closing_stock, closing_value,
			_source, _source_id, _synced_at, _created_at, _updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, NOW(), NOW(), NOW()
		)
		ON CONFLICT (_source, _source_id)
		DO UPDATE SET
			item_id = EXCLUDED.item_id,
			voucher_id = EXCLUDED.voucher_id,
			movement_type = EXCLUDED.movement_type,
			movement_date = EXCLUDED.movement_date,
			movement_datetime = EXCLUDED.movement_datetime,
			qty_in = EXCLUDED.qty_in,
			qty_out = EXCLUDED.qty_out,
			unit = EXCLUDED.unit,
			rate = EXCLUDED.rate,
			value = EXCLUDED.value,
			batch_number = EXCLUDED.batch_number,
			expiry_date = EXCLUDED.expiry_date,
			godown_name = EXCLUDED.godown_name,
			closing_stock = EXCLUDED.closing_stock,
			closing_value = EXCLUDED.closing_value,
			_synced_at = NOW(),
			_updated_at = NOW()
		RETURNING movement_id
	`

	for _, s := range movements {
		batch.Queue(query,
			s.ExternalMovementID, s.ItemID, s.VoucherID, s.MovementType,
			s.MovementDate, s.MovementDatetime, s.QtyIn, s.QtyOut, s.Unit,
			s.Rate, s.Value, s.BatchNumber, s.ExpiryDate, s.GodownName,
			s.ClosingStock, s.ClosingValue,
			s.Source, s.SourceID,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i, s := range movements {
		err := results.QueryRow().Scan(&s.MovementID)
		if err != nil {
			return fmt.Errorf("warehouse: failed to upsert stock movement %d/%d (%s): %w",
				i+1, len(movements), s.SourceID, err)
		}
	}

	return nil
}

// ============================================================================
// Generic Operations
// ============================================================================

// GetTableStats returns statistics about a table including row count.
func (r *Repo) GetTableStats(ctx context.Context, schema, table string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", schema, pgx.Identifier{table}.Sanitize())
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("warehouse: failed to get table stats for %s.%s: %w", schema, table, err)
	}
	return count, nil
}

// ExecuteInTx executes a function within a transaction.
func (r *Repo) ExecuteInTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("warehouse: failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("warehouse: failed to commit transaction: %w", err)
	}

	return nil
}

// IsRecordExists checks if a record exists by source and source_id.
func (r *Repo) IsRecordExists(ctx context.Context, schema, table string, source string, sourceID string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 FROM %s.%s
			WHERE _source = $1 AND _source_id = $2
		)
	`, pgx.Identifier{schema}.Sanitize(), pgx.Identifier{table}.Sanitize())

	err := r.pool.QueryRow(ctx, query, source, sourceID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("warehouse: failed to check record existence: %w", err)
	}

	return exists, nil
}
