-- MediSync Initial Schema Migration
-- Version: 001
-- Description: Create all foundational schemas and tables for MediSync data warehouse
-- Schemas: hims_analytics, tally_analytics, app, vectors
--
-- This migration establishes:
-- 1. Required PostgreSQL extensions (pgvector, uuid-ossp, pg_stat_statements)
-- 2. Four main schemas with their dimension and fact tables
-- 3. Standard audit columns (_source, _source_id, _synced_at, _created_at, _updated_at)
-- 4. Appropriate indexes for query performance
-- 5. Locale support columns for i18n (EN/AR)

-- ============================================================================
-- EXTENSIONS
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgvector";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- ============================================================================
-- SCHEMA: hims_analytics
-- Healthcare Information Management System data warehouse
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS hims_analytics;

-- Dimension: Patients
CREATE TABLE hims_analytics.dim_patients (
    patient_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_patient_id VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    name_ar VARCHAR(255),
    date_of_birth DATE,
    gender VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    address_en TEXT,
    address_ar TEXT,
    blood_group VARCHAR(10),
    nationality VARCHAR(100),
    national_id VARCHAR(100),
    insurance_provider VARCHAR(255),
    insurance_policy_number VARCHAR(100),
    emergency_contact_name VARCHAR(255),
    emergency_contact_phone VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Uniqueness constraint on source system ID
    CONSTRAINT uq_dim_patients_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_patients_external_id ON hims_analytics.dim_patients(external_patient_id);
CREATE INDEX idx_dim_patients_name_en ON hims_analytics.dim_patients(name_en);
CREATE INDEX idx_dim_patients_name_ar ON hims_analytics.dim_patients(name_ar);
CREATE INDEX idx_dim_patients_phone ON hims_analytics.dim_patients(phone);
CREATE INDEX idx_dim_patients_synced_at ON hims_analytics.dim_patients(_synced_at);

-- Dimension: Doctors
CREATE TABLE hims_analytics.dim_doctors (
    doctor_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_doctor_id VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    specialty_en VARCHAR(255),
    specialty_ar VARCHAR(255),
    department_en VARCHAR(255),
    department_ar VARCHAR(255),
    qualification VARCHAR(500),
    license_number VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    consultation_fee DECIMAL(12, 2),
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_doctors_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_doctors_external_id ON hims_analytics.dim_doctors(external_doctor_id);
CREATE INDEX idx_dim_doctors_name_en ON hims_analytics.dim_doctors(name_en);
CREATE INDEX idx_dim_doctors_specialty ON hims_analytics.dim_doctors(specialty_en);
CREATE INDEX idx_dim_doctors_department ON hims_analytics.dim_doctors(department_en);
CREATE INDEX idx_dim_doctors_synced_at ON hims_analytics.dim_doctors(_synced_at);

-- Dimension: Drugs/Medications
CREATE TABLE hims_analytics.dim_drugs (
    drug_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_drug_id VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    generic_name_en VARCHAR(255),
    generic_name_ar VARCHAR(255),
    category_en VARCHAR(255),
    category_ar VARCHAR(255),
    dosage_form VARCHAR(100),
    strength VARCHAR(100),
    unit VARCHAR(50),
    manufacturer VARCHAR(255),
    unit_price DECIMAL(12, 2),
    reorder_level INTEGER,
    is_controlled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_drugs_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_drugs_external_id ON hims_analytics.dim_drugs(external_drug_id);
CREATE INDEX idx_dim_drugs_name_en ON hims_analytics.dim_drugs(name_en);
CREATE INDEX idx_dim_drugs_category ON hims_analytics.dim_drugs(category_en);
CREATE INDEX idx_dim_drugs_synced_at ON hims_analytics.dim_drugs(_synced_at);

-- Dimension: Departments
CREATE TABLE hims_analytics.dim_departments (
    department_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_department_id VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    code VARCHAR(50),
    parent_department_id UUID REFERENCES hims_analytics.dim_departments(department_id),
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_departments_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_departments_code ON hims_analytics.dim_departments(code);
CREATE INDEX idx_dim_departments_synced_at ON hims_analytics.dim_departments(_synced_at);

-- Fact: Appointments
CREATE TABLE hims_analytics.fact_appointments (
    appt_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_appt_id VARCHAR(255) NOT NULL,
    patient_id UUID NOT NULL REFERENCES hims_analytics.dim_patients(patient_id),
    doctor_id UUID NOT NULL REFERENCES hims_analytics.dim_doctors(doctor_id),
    department_id UUID REFERENCES hims_analytics.dim_departments(department_id),
    appt_date DATE NOT NULL,
    appt_time TIME,
    appt_datetime TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL, -- scheduled, confirmed, checked_in, in_progress, completed, cancelled, no_show
    appt_type VARCHAR(100), -- consultation, follow_up, procedure, emergency
    duration_minutes INTEGER,
    chief_complaint TEXT,
    diagnosis_code VARCHAR(50),
    diagnosis_description TEXT,
    billing_id UUID, -- FK to fact_billing, set after billing
    notes TEXT,
    is_walk_in BOOLEAN DEFAULT FALSE,
    cancellation_reason TEXT,
    cancelled_at TIMESTAMPTZ,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fact_appointments_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_fact_appointments_patient ON hims_analytics.fact_appointments(patient_id);
CREATE INDEX idx_fact_appointments_doctor ON hims_analytics.fact_appointments(doctor_id);
CREATE INDEX idx_fact_appointments_date ON hims_analytics.fact_appointments(appt_date);
CREATE INDEX idx_fact_appointments_status ON hims_analytics.fact_appointments(status);
CREATE INDEX idx_fact_appointments_department ON hims_analytics.fact_appointments(department_id);
CREATE INDEX idx_fact_appointments_synced_at ON hims_analytics.fact_appointments(_synced_at);
CREATE INDEX idx_fact_appointments_datetime ON hims_analytics.fact_appointments(appt_datetime);

-- Fact: Billing
CREATE TABLE hims_analytics.fact_billing (
    bill_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_bill_id VARCHAR(255) NOT NULL,
    patient_id UUID NOT NULL REFERENCES hims_analytics.dim_patients(patient_id),
    appt_id UUID REFERENCES hims_analytics.fact_appointments(appt_id),
    bill_date DATE NOT NULL,
    bill_datetime TIMESTAMPTZ,
    subtotal_amount DECIMAL(14, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(14, 2) DEFAULT 0,
    tax_amount DECIMAL(14, 2) DEFAULT 0,
    total_amount DECIMAL(14, 2) NOT NULL,
    paid_amount DECIMAL(14, 2) DEFAULT 0,
    outstanding_amount DECIMAL(14, 2) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,
    payment_mode VARCHAR(50), -- cash, card, insurance, bank_transfer, cheque
    payment_status VARCHAR(50) NOT NULL, -- pending, partial, paid, cancelled, refunded
    insurance_claim_id VARCHAR(100),
    insurance_amount DECIMAL(14, 2) DEFAULT 0,
    department_en VARCHAR(255),
    department_ar VARCHAR(255),
    bill_type VARCHAR(100), -- consultation, procedure, pharmacy, lab, imaging
    receipt_number VARCHAR(100),
    notes TEXT,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fact_billing_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_fact_billing_patient ON hims_analytics.fact_billing(patient_id);
CREATE INDEX idx_fact_billing_date ON hims_analytics.fact_billing(bill_date);
CREATE INDEX idx_fact_billing_payment_status ON hims_analytics.fact_billing(payment_status);
CREATE INDEX idx_fact_billing_payment_mode ON hims_analytics.fact_billing(payment_mode);
CREATE INDEX idx_fact_billing_department ON hims_analytics.fact_billing(department_en);
CREATE INDEX idx_fact_billing_synced_at ON hims_analytics.fact_billing(_synced_at);
CREATE INDEX idx_fact_billing_bill_type ON hims_analytics.fact_billing(bill_type);

-- Fact: Pharmacy Dispensations
CREATE TABLE hims_analytics.fact_pharmacy_dispensations (
    disp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_disp_id VARCHAR(255) NOT NULL,
    drug_id UUID NOT NULL REFERENCES hims_analytics.dim_drugs(drug_id),
    patient_id UUID NOT NULL REFERENCES hims_analytics.dim_patients(patient_id),
    doctor_id UUID REFERENCES hims_analytics.dim_doctors(doctor_id),
    bill_id UUID REFERENCES hims_analytics.fact_billing(bill_id),
    prescription_id VARCHAR(255),
    disp_date DATE NOT NULL,
    disp_datetime TIMESTAMPTZ,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50),
    dosage_instructions TEXT,
    days_supply INTEGER,
    unit_price DECIMAL(12, 2) NOT NULL,
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    tax_amount DECIMAL(12, 2) DEFAULT 0,
    total_amount DECIMAL(12, 2) NOT NULL,
    batch_number VARCHAR(100),
    expiry_date DATE,
    is_substituted BOOLEAN DEFAULT FALSE,
    original_drug_id UUID REFERENCES hims_analytics.dim_drugs(drug_id),
    notes TEXT,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'hims',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fact_pharmacy_disp_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_fact_pharmacy_disp_drug ON hims_analytics.fact_pharmacy_dispensations(drug_id);
CREATE INDEX idx_fact_pharmacy_disp_patient ON hims_analytics.fact_pharmacy_dispensations(patient_id);
CREATE INDEX idx_fact_pharmacy_disp_date ON hims_analytics.fact_pharmacy_dispensations(disp_date);
CREATE INDEX idx_fact_pharmacy_disp_prescription ON hims_analytics.fact_pharmacy_dispensations(prescription_id);
CREATE INDEX idx_fact_pharmacy_disp_synced_at ON hims_analytics.fact_pharmacy_dispensations(_synced_at);
CREATE INDEX idx_fact_pharmacy_disp_doctor ON hims_analytics.fact_pharmacy_dispensations(doctor_id);

-- ============================================================================
-- SCHEMA: tally_analytics
-- Tally ERP financial data warehouse
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS tally_analytics;

-- Dimension: Ledgers (Chart of Accounts)
CREATE TABLE tally_analytics.dim_ledgers (
    ledger_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_ledger_id VARCHAR(255) NOT NULL,
    ledger_name VARCHAR(255) NOT NULL,
    ledger_name_ar VARCHAR(255),
    ledger_group VARCHAR(255) NOT NULL,
    parent_group VARCHAR(255),
    ledger_type VARCHAR(100), -- asset, liability, income, expense, equity
    opening_balance DECIMAL(16, 2) DEFAULT 0,
    closing_balance DECIMAL(16, 2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'SAR',
    is_bank_account BOOLEAN DEFAULT FALSE,
    bank_name VARCHAR(255),
    bank_account_number VARCHAR(100),
    ifsc_code VARCHAR(50),
    gst_registration VARCHAR(100),
    pan_number VARCHAR(50),
    credit_period_days INTEGER,
    credit_limit DECIMAL(16, 2),
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'tally',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_ledgers_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_ledgers_name ON tally_analytics.dim_ledgers(ledger_name);
CREATE INDEX idx_dim_ledgers_group ON tally_analytics.dim_ledgers(ledger_group);
CREATE INDEX idx_dim_ledgers_type ON tally_analytics.dim_ledgers(ledger_type);
CREATE INDEX idx_dim_ledgers_synced_at ON tally_analytics.dim_ledgers(_synced_at);
CREATE INDEX idx_dim_ledgers_external_id ON tally_analytics.dim_ledgers(external_ledger_id);

-- Dimension: Cost Centres
CREATE TABLE tally_analytics.dim_cost_centres (
    cc_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_cc_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    code VARCHAR(50),
    parent_cc_id UUID REFERENCES tally_analytics.dim_cost_centres(cc_id),
    category VARCHAR(100),
    is_revenue_centre BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'tally',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_cost_centres_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_cost_centres_name ON tally_analytics.dim_cost_centres(name);
CREATE INDEX idx_dim_cost_centres_code ON tally_analytics.dim_cost_centres(code);
CREATE INDEX idx_dim_cost_centres_synced_at ON tally_analytics.dim_cost_centres(_synced_at);

-- Dimension: Inventory Items (Stock Items)
CREATE TABLE tally_analytics.dim_inventory_items (
    item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_item_id VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    part_number VARCHAR(100),
    category VARCHAR(255),
    sub_category VARCHAR(255),
    stock_group VARCHAR(255),
    unit VARCHAR(50),
    alternate_unit VARCHAR(50),
    conversion_factor DECIMAL(12, 4),
    gst_rate DECIMAL(5, 2),
    hsn_code VARCHAR(50),
    purchase_price DECIMAL(14, 2),
    selling_price DECIMAL(14, 2),
    mrp DECIMAL(14, 2),
    reorder_level INTEGER,
    minimum_order_qty INTEGER,
    is_batch_wise BOOLEAN DEFAULT FALSE,
    maintain_expiry BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'tally',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_dim_inventory_items_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_dim_inventory_items_name ON tally_analytics.dim_inventory_items(name_en);
CREATE INDEX idx_dim_inventory_items_category ON tally_analytics.dim_inventory_items(category);
CREATE INDEX idx_dim_inventory_items_stock_group ON tally_analytics.dim_inventory_items(stock_group);
CREATE INDEX idx_dim_inventory_items_part_number ON tally_analytics.dim_inventory_items(part_number);
CREATE INDEX idx_dim_inventory_items_synced_at ON tally_analytics.dim_inventory_items(_synced_at);

-- Fact: Vouchers (All Tally Transactions)
CREATE TABLE tally_analytics.fact_vouchers (
    voucher_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_voucher_id VARCHAR(255) NOT NULL,
    voucher_number VARCHAR(100) NOT NULL,
    voucher_type VARCHAR(100) NOT NULL, -- sales, purchase, receipt, payment, journal, contra, credit_note, debit_note
    voucher_date DATE NOT NULL,
    voucher_datetime TIMESTAMPTZ,
    ledger_id UUID NOT NULL REFERENCES tally_analytics.dim_ledgers(ledger_id),
    contra_ledger_id UUID REFERENCES tally_analytics.dim_ledgers(ledger_id),
    cost_centre_id UUID REFERENCES tally_analytics.dim_cost_centres(cc_id),
    amount DECIMAL(16, 2) NOT NULL,
    is_debit BOOLEAN NOT NULL,
    currency VARCHAR(10) DEFAULT 'SAR',
    exchange_rate DECIMAL(12, 6) DEFAULT 1,
    base_currency_amount DECIMAL(16, 2),
    narration TEXT,
    reference_number VARCHAR(255),
    reference_date DATE,
    party_name VARCHAR(255),
    bill_number VARCHAR(100),
    bill_date DATE,
    due_date DATE,
    instrument_number VARCHAR(100), -- cheque/DD number
    instrument_date DATE,
    bank_name VARCHAR(255),
    gst_registration VARCHAR(100),
    invoice_number VARCHAR(100),
    is_cancelled BOOLEAN DEFAULT FALSE,
    cancelled_date DATE,
    cancellation_reason TEXT,
    is_optional BOOLEAN DEFAULT FALSE,
    -- Inventory related (for sales/purchase vouchers)
    has_inventory BOOLEAN DEFAULT FALSE,
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'tally',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fact_vouchers_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_fact_vouchers_date ON tally_analytics.fact_vouchers(voucher_date);
CREATE INDEX idx_fact_vouchers_type ON tally_analytics.fact_vouchers(voucher_type);
CREATE INDEX idx_fact_vouchers_ledger ON tally_analytics.fact_vouchers(ledger_id);
CREATE INDEX idx_fact_vouchers_contra_ledger ON tally_analytics.fact_vouchers(contra_ledger_id);
CREATE INDEX idx_fact_vouchers_cost_centre ON tally_analytics.fact_vouchers(cost_centre_id);
CREATE INDEX idx_fact_vouchers_synced_at ON tally_analytics.fact_vouchers(_synced_at);
CREATE INDEX idx_fact_vouchers_reference ON tally_analytics.fact_vouchers(reference_number);
CREATE INDEX idx_fact_vouchers_party ON tally_analytics.fact_vouchers(party_name);
CREATE INDEX idx_fact_vouchers_voucher_number ON tally_analytics.fact_vouchers(voucher_number);

-- Fact: Stock Movements (Inventory Transactions)
CREATE TABLE tally_analytics.fact_stock_movements (
    movement_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_movement_id VARCHAR(255) NOT NULL,
    item_id UUID NOT NULL REFERENCES tally_analytics.dim_inventory_items(item_id),
    voucher_id UUID REFERENCES tally_analytics.fact_vouchers(voucher_id),
    movement_type VARCHAR(50) NOT NULL, -- purchase, sales, transfer_in, transfer_out, adjustment, production, consumption
    movement_date DATE NOT NULL,
    movement_datetime TIMESTAMPTZ,
    qty_in DECIMAL(14, 4) DEFAULT 0,
    qty_out DECIMAL(14, 4) DEFAULT 0,
    unit VARCHAR(50),
    rate DECIMAL(14, 4),
    value DECIMAL(16, 2),
    batch_number VARCHAR(100),
    expiry_date DATE,
    godown_name VARCHAR(255), -- warehouse/location
    closing_stock DECIMAL(14, 4),
    closing_value DECIMAL(16, 2),
    -- Standard audit columns
    _source VARCHAR(50) NOT NULL DEFAULT 'tally',
    _source_id VARCHAR(255) NOT NULL,
    _synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fact_stock_movements_source_id UNIQUE (_source, _source_id)
);

CREATE INDEX idx_fact_stock_movements_item ON tally_analytics.fact_stock_movements(item_id);
CREATE INDEX idx_fact_stock_movements_date ON tally_analytics.fact_stock_movements(movement_date);
CREATE INDEX idx_fact_stock_movements_type ON tally_analytics.fact_stock_movements(movement_type);
CREATE INDEX idx_fact_stock_movements_voucher ON tally_analytics.fact_stock_movements(voucher_id);
CREATE INDEX idx_fact_stock_movements_synced_at ON tally_analytics.fact_stock_movements(_synced_at);
CREATE INDEX idx_fact_stock_movements_godown ON tally_analytics.fact_stock_movements(godown_name);

-- ============================================================================
-- SCHEMA: app
-- Application-specific tables (users, preferences, workflows, audit)
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS app;

-- Users (linked to Keycloak)
CREATE TABLE app.users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    keycloak_sub VARCHAR(255) NOT NULL UNIQUE, -- Keycloak subject ID
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100),
    full_name_en VARCHAR(255),
    full_name_ar VARCHAR(255),
    role VARCHAR(100) NOT NULL, -- admin, finance_head, accountant_lead, accountant, manager, pharmacy_manager, analyst, viewer
    department VARCHAR(255),
    cost_centres TEXT[], -- Array of cost centre IDs user has access to
    locale VARCHAR(10) NOT NULL DEFAULT 'en', -- 'en' or 'ar'
    calendar_system VARCHAR(20) NOT NULL DEFAULT 'gregorian', -- 'gregorian' or 'hijri'
    timezone VARCHAR(100) DEFAULT 'Asia/Riyadh',
    phone VARCHAR(50),
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON app.users(email);
CREATE INDEX idx_users_keycloak_sub ON app.users(keycloak_sub);
CREATE INDEX idx_users_role ON app.users(role);
CREATE INDEX idx_users_department ON app.users(department);

-- User Preferences (extended settings)
CREATE TABLE app.user_preferences (
    user_id UUID PRIMARY KEY REFERENCES app.users(user_id) ON DELETE CASCADE,
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    number_format VARCHAR(50) DEFAULT 'en-US', -- locale for number formatting
    date_format VARCHAR(50) DEFAULT 'YYYY-MM-DD',
    calendar_system VARCHAR(20) DEFAULT 'gregorian',
    currency VARCHAR(10) DEFAULT 'SAR',
    report_language VARCHAR(10) DEFAULT 'en',
    ai_response_language VARCHAR(10) DEFAULT 'en',
    dashboard_layout JSONB DEFAULT '{}',
    notification_preferences JSONB DEFAULT '{"email": true, "sms": false, "push": true}',
    theme VARCHAR(20) DEFAULT 'light', -- light, dark, auto
    accessibility_options JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit Log (append-only, immutable)
CREATE TABLE app.audit_log (
    log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES app.users(user_id),
    session_id VARCHAR(255),
    action VARCHAR(100) NOT NULL, -- create, read, update, delete, export, approve, sync, login, logout
    resource VARCHAR(100) NOT NULL, -- report, dashboard, document, voucher, patient, etc.
    resource_id VARCHAR(255),
    changes_json JSONB, -- {before: {...}, after: {...}} for updates
    metadata JSONB, -- Additional context
    ip_address INET,
    user_agent TEXT,
    locale VARCHAR(10),
    request_id VARCHAR(255),
    duration_ms INTEGER,
    status VARCHAR(50) DEFAULT 'success', -- success, failure, error
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Partition audit_log by month for performance
-- Note: Partitioning will be set up separately in production
CREATE INDEX idx_audit_log_user ON app.audit_log(user_id);
CREATE INDEX idx_audit_log_action ON app.audit_log(action);
CREATE INDEX idx_audit_log_resource ON app.audit_log(resource, resource_id);
CREATE INDEX idx_audit_log_created ON app.audit_log(created_at);
CREATE INDEX idx_audit_log_session ON app.audit_log(session_id);

-- ETL Quarantine (failed/invalid records)
CREATE TABLE app.etl_quarantine (
    record_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID,
    source VARCHAR(50) NOT NULL, -- 'tally', 'hims', 'bank'
    source_table VARCHAR(255),
    source_id VARCHAR(255),
    raw_data JSONB NOT NULL, -- Original record as JSON
    raw_xml TEXT, -- For Tally XML records
    error_reason TEXT NOT NULL,
    error_code VARCHAR(50),
    error_details JSONB,
    validation_rules_failed TEXT[], -- List of failed validation rules
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    last_retry_at TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'pending', -- pending, retrying, resolved, ignored
    resolved_by UUID REFERENCES app.users(user_id),
    resolved_at TIMESTAMPTZ,
    resolution_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_etl_quarantine_source ON app.etl_quarantine(source);
CREATE INDEX idx_etl_quarantine_status ON app.etl_quarantine(status);
CREATE INDEX idx_etl_quarantine_created ON app.etl_quarantine(created_at);
CREATE INDEX idx_etl_quarantine_batch ON app.etl_quarantine(batch_id);

-- ETL Quality Reports (from C-06 agent)
CREATE TABLE app.etl_quality_report (
    report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL,
    source VARCHAR(50) NOT NULL,
    sync_started_at TIMESTAMPTZ NOT NULL,
    sync_completed_at TIMESTAMPTZ,
    total_records INTEGER NOT NULL DEFAULT 0,
    records_processed INTEGER DEFAULT 0,
    records_inserted INTEGER DEFAULT 0,
    records_updated INTEGER DEFAULT 0,
    records_quarantined INTEGER DEFAULT 0,
    -- Quality check results
    completeness_score DECIMAL(5, 2), -- 0-100%
    uniqueness_score DECIMAL(5, 2),
    referential_integrity_score DECIMAL(5, 2),
    range_validation_score DECIMAL(5, 2),
    overall_quality_score DECIMAL(5, 2),
    -- Check details
    missing_value_count INTEGER DEFAULT 0,
    duplicate_count INTEGER DEFAULT 0,
    integrity_violation_count INTEGER DEFAULT 0,
    range_violation_count INTEGER DEFAULT 0,
    anomaly_count INTEGER DEFAULT 0,
    -- Row count tracking
    previous_row_count INTEGER,
    current_row_count INTEGER,
    row_count_delta_pct DECIMAL(5, 2),
    -- Alert info
    alerts_generated INTEGER DEFAULT 0,
    alert_details JSONB,
    -- Status
    validation_passed BOOLEAN NOT NULL,
    failure_reasons TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_etl_quality_report_batch ON app.etl_quality_report(batch_id);
CREATE INDEX idx_etl_quality_report_source ON app.etl_quality_report(source);
CREATE INDEX idx_etl_quality_report_created ON app.etl_quality_report(created_at);
CREATE INDEX idx_etl_quality_report_passed ON app.etl_quality_report(validation_passed);

-- ETL State (cursor tracking for incremental sync)
CREATE TABLE app.etl_state (
    state_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(50) NOT NULL,
    entity VARCHAR(255) NOT NULL, -- 'ledgers', 'vouchers', 'patients', etc.
    last_sync_at TIMESTAMPTZ,
    last_alter_id VARCHAR(255), -- For Tally LastAlterID cursor
    last_modified_at TIMESTAMPTZ, -- For HIMS modified_since cursor
    cursor_value TEXT, -- Generic cursor storage
    cursor_type VARCHAR(50), -- 'alter_id', 'timestamp', 'offset'
    records_synced INTEGER DEFAULT 0,
    sync_status VARCHAR(50) DEFAULT 'idle', -- idle, running, completed, failed
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_etl_state_source_entity UNIQUE (source, entity)
);

CREATE INDEX idx_etl_state_source ON app.etl_state(source);
CREATE INDEX idx_etl_state_status ON app.etl_state(sync_status);

-- Notification Queue
CREATE TABLE app.notification_queue (
    notif_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES app.users(user_id),
    notification_type VARCHAR(100) NOT NULL, -- etl_alert, approval_required, report_ready, kpi_alert, system
    channel VARCHAR(50) NOT NULL, -- email, sms, push, in_app, slack
    priority VARCHAR(20) DEFAULT 'normal', -- low, normal, high, urgent
    subject_en VARCHAR(500),
    subject_ar VARCHAR(500),
    body_en TEXT,
    body_ar TEXT,
    body_html TEXT,
    metadata JSONB, -- Additional data for templates
    status VARCHAR(50) DEFAULT 'pending', -- pending, queued, sent, delivered, failed, cancelled
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_queue_user ON app.notification_queue(user_id);
CREATE INDEX idx_notification_queue_status ON app.notification_queue(status);
CREATE INDEX idx_notification_queue_scheduled ON app.notification_queue(scheduled_at);
CREATE INDEX idx_notification_queue_type ON app.notification_queue(notification_type);

-- Pinned Charts (Dashboard widgets)
CREATE TABLE app.pinned_charts (
    pin_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES app.users(user_id) ON DELETE CASCADE,
    dashboard_id UUID, -- NULL for default dashboard
    title_en VARCHAR(255) NOT NULL,
    title_ar VARCHAR(255),
    chart_type VARCHAR(50) NOT NULL, -- bar, line, pie, scatter, table, kpi_card, gauge
    chart_config JSONB NOT NULL, -- ECharts configuration
    sql_query TEXT, -- Original SQL that generated the chart
    query_params JSONB,
    data_source VARCHAR(100), -- 'hims', 'tally', 'combined'
    refresh_interval INTEGER DEFAULT 300, -- seconds, 0 = manual only
    last_refreshed_at TIMESTAMPTZ,
    position_x INTEGER DEFAULT 0,
    position_y INTEGER DEFAULT 0,
    width INTEGER DEFAULT 4, -- Grid units
    height INTEGER DEFAULT 3,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pinned_charts_user ON app.pinned_charts(user_id);
CREATE INDEX idx_pinned_charts_dashboard ON app.pinned_charts(dashboard_id);

-- Scheduled Reports
CREATE TABLE app.scheduled_reports (
    sched_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES app.users(user_id) ON DELETE CASCADE,
    report_name_en VARCHAR(255) NOT NULL,
    report_name_ar VARCHAR(255),
    report_type VARCHAR(100) NOT NULL, -- pl_statement, balance_sheet, cash_flow, aging_report, inventory, custom
    report_template_id UUID,
    params_json JSONB NOT NULL, -- Report parameters
    output_format VARCHAR(20) DEFAULT 'pdf', -- pdf, excel, html, csv
    locale VARCHAR(10) DEFAULT 'en',
    cron_expr VARCHAR(100) NOT NULL, -- Standard cron expression
    timezone VARCHAR(100) DEFAULT 'Asia/Riyadh',
    recipients JSONB, -- [{email, name, locale}]
    delivery_channel VARCHAR(50) DEFAULT 'email', -- email, slack, in_app
    is_active BOOLEAN DEFAULT TRUE,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ,
    last_run_status VARCHAR(50),
    last_error TEXT,
    run_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_scheduled_reports_user ON app.scheduled_reports(user_id);
CREATE INDEX idx_scheduled_reports_next_run ON app.scheduled_reports(next_run_at);
CREATE INDEX idx_scheduled_reports_active ON app.scheduled_reports(is_active);

-- Approval Workflows (for B-08 Approval Workflow Agent)
CREATE TABLE app.approval_workflows (
    workflow_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_type VARCHAR(100) NOT NULL, -- invoice, journal_entry, payment, bank_transfer
    document_type VARCHAR(100) NOT NULL,
    document_id VARCHAR(255) NOT NULL,
    document_data JSONB NOT NULL,
    amount DECIMAL(16, 2),
    currency VARCHAR(10) DEFAULT 'SAR',
    -- Workflow state
    current_step INTEGER DEFAULT 1,
    total_steps INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, in_progress, approved, rejected, cancelled
    -- Step details stored as array
    steps JSONB NOT NULL, -- [{step: 1, role: "accountant", status: "approved", user_id: ..., approved_at: ..., notes: ...}]
    -- Initiator
    created_by UUID NOT NULL REFERENCES app.users(user_id),
    -- Final approver
    final_approver_id UUID REFERENCES app.users(user_id),
    final_approved_at TIMESTAMPTZ,
    -- Rejection info
    rejected_by UUID REFERENCES app.users(user_id),
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    -- Tally sync status (for B-09)
    tally_sync_status VARCHAR(50), -- pending, synced, failed
    tally_sync_at TIMESTAMPTZ,
    tally_voucher_id VARCHAR(255),
    tally_error TEXT,
    -- Timestamps
    due_date DATE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approval_workflows_status ON app.approval_workflows(status);
CREATE INDEX idx_approval_workflows_type ON app.approval_workflows(workflow_type);
CREATE INDEX idx_approval_workflows_created_by ON app.approval_workflows(created_by);
CREATE INDEX idx_approval_workflows_due_date ON app.approval_workflows(due_date);
CREATE INDEX idx_approval_workflows_created ON app.approval_workflows(created_at);

-- ============================================================================
-- SCHEMA: vectors
-- Vector embeddings for semantic search (pgvector)
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS vectors;

-- Schema Embeddings (for A-01 Text-to-SQL context)
CREATE TABLE vectors.schema_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schema_name VARCHAR(100) NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    column_name VARCHAR(255),
    object_type VARCHAR(50) NOT NULL, -- 'table', 'column', 'relationship'
    description_en TEXT,
    description_ar TEXT,
    data_type VARCHAR(100),
    sample_values TEXT[], -- Example values for context
    business_context TEXT, -- Business meaning/usage
    embedding vector(1536), -- OpenAI ada-002 dimension
    model_name VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_schema_embeddings UNIQUE (schema_name, table_name, column_name)
);

-- Create HNSW index for fast similarity search
CREATE INDEX idx_schema_embeddings_vector ON vectors.schema_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);
CREATE INDEX idx_schema_embeddings_table ON vectors.schema_embeddings(schema_name, table_name);

-- Metric Embeddings (for semantic metric search)
CREATE TABLE vectors.metric_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_name VARCHAR(255) NOT NULL UNIQUE,
    metric_type VARCHAR(50) NOT NULL, -- 'simple', 'derived', 'ratio', 'cumulative'
    description_en TEXT NOT NULL,
    description_ar TEXT,
    formula TEXT, -- SQL or expression
    dimensions TEXT[], -- Associated dimensions
    filters TEXT[], -- Default filters
    business_domain VARCHAR(100), -- 'finance', 'operations', 'pharmacy', 'hr'
    tags TEXT[],
    embedding vector(1536),
    model_name VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_metric_embeddings_vector ON vectors.metric_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);
CREATE INDEX idx_metric_embeddings_domain ON vectors.metric_embeddings(business_domain);
CREATE INDEX idx_metric_embeddings_type ON vectors.metric_embeddings(metric_type);

-- Query History Embeddings (for few-shot learning)
CREATE TABLE vectors.query_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES app.users(user_id),
    natural_language_query TEXT NOT NULL,
    generated_sql TEXT NOT NULL,
    query_locale VARCHAR(10) DEFAULT 'en',
    -- Quality indicators
    execution_success BOOLEAN NOT NULL,
    result_count INTEGER,
    execution_time_ms INTEGER,
    user_feedback VARCHAR(20), -- 'positive', 'negative', 'corrected', null
    corrected_sql TEXT, -- If user corrected the SQL
    confidence_score DECIMAL(5, 2),
    -- Context
    tables_used TEXT[],
    aggregations_used TEXT[],
    filters_applied JSONB,
    -- Embedding for similar query lookup
    query_embedding vector(1536),
    model_name VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_query_history_vector ON vectors.query_history
    USING hnsw (query_embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);
CREATE INDEX idx_query_history_user ON vectors.query_history(user_id);
CREATE INDEX idx_query_history_success ON vectors.query_history(execution_success);
CREATE INDEX idx_query_history_feedback ON vectors.query_history(user_feedback);
CREATE INDEX idx_query_history_created ON vectors.query_history(created_at);

-- Document Embeddings (for B-05 Ledger Mapping context)
CREATE TABLE vectors.document_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL,
    document_type VARCHAR(100) NOT NULL, -- 'invoice', 'receipt', 'bank_statement', 'bill'
    chunk_index INTEGER DEFAULT 0,
    content TEXT NOT NULL,
    -- Extracted fields
    vendor_name VARCHAR(255),
    amount DECIMAL(16, 2),
    date DATE,
    extracted_fields JSONB,
    -- Embedding
    embedding vector(1536),
    model_name VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_document_embeddings_vector ON vectors.document_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);
CREATE INDEX idx_document_embeddings_doc ON vectors.document_embeddings(document_id);
CREATE INDEX idx_document_embeddings_type ON vectors.document_embeddings(document_type);

-- ============================================================================
-- FUNCTIONS: Updated timestamp trigger
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW._updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at triggers to all relevant tables
DO $$
DECLARE
    t record;
BEGIN
    FOR t IN
        SELECT schemaname, tablename
        FROM pg_tables
        WHERE schemaname IN ('hims_analytics', 'tally_analytics', 'app', 'vectors')
        AND tablename NOT IN ('audit_log') -- audit_log is append-only
    LOOP
        -- Check if table has _updated_at or updated_at column
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = t.schemaname
            AND table_name = t.tablename
            AND column_name IN ('_updated_at', 'updated_at')
        ) THEN
            EXECUTE format(
                'CREATE TRIGGER update_%I_%I_updated_at
                 BEFORE UPDATE ON %I.%I
                 FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()',
                t.schemaname, t.tablename, t.schemaname, t.tablename
            );
        END IF;
    END LOOP;
END;
$$;

-- ============================================================================
-- COMMENTS: Documentation for schema objects
-- ============================================================================

-- Schema comments
COMMENT ON SCHEMA hims_analytics IS 'Healthcare Information Management System analytics data warehouse';
COMMENT ON SCHEMA tally_analytics IS 'Tally ERP financial analytics data warehouse';
COMMENT ON SCHEMA app IS 'MediSync application tables (users, preferences, workflows, audit)';
COMMENT ON SCHEMA vectors IS 'Vector embeddings for semantic search and AI context (pgvector)';

-- Key table comments
COMMENT ON TABLE hims_analytics.dim_patients IS 'Patient master data from HIMS with bilingual name support';
COMMENT ON TABLE hims_analytics.fact_appointments IS 'Appointment transactions from HIMS';
COMMENT ON TABLE hims_analytics.fact_billing IS 'Billing records from HIMS with payment tracking';
COMMENT ON TABLE hims_analytics.fact_pharmacy_dispensations IS 'Pharmacy dispensation records from HIMS';

COMMENT ON TABLE tally_analytics.dim_ledgers IS 'Chart of Accounts from Tally ERP';
COMMENT ON TABLE tally_analytics.fact_vouchers IS 'All financial transactions/vouchers from Tally';
COMMENT ON TABLE tally_analytics.fact_stock_movements IS 'Inventory movements from Tally';

COMMENT ON TABLE app.audit_log IS 'Immutable, append-only audit trail for all system actions';
COMMENT ON TABLE app.etl_quarantine IS 'Quarantined records that failed ETL validation';
COMMENT ON TABLE app.etl_quality_report IS 'Data quality reports from C-06 validation agent';
COMMENT ON TABLE app.approval_workflows IS 'Multi-step approval workflows for financial transactions';

COMMENT ON TABLE vectors.schema_embeddings IS 'Embeddings of database schema for Text-to-SQL context';
COMMENT ON TABLE vectors.metric_embeddings IS 'Embeddings of business metrics for semantic search';
COMMENT ON TABLE vectors.query_history IS 'Historical NL queries with embeddings for few-shot learning';

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
