-- MediSync Complete Schema Example
-- This demonstrates idiomatic PostgreSQL schema design

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "vector";

-- =============================================================================
-- UTILITY FUNCTIONS
-- =============================================================================

-- Automatic updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Soft delete function
CREATE OR REPLACE FUNCTION soft_delete()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMPANIES
-- =============================================================================

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_companies_slug ON companies(slug) WHERE deleted_at IS NULL;

-- =============================================================================
-- USERS
-- =============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(company_id, email)
);

CREATE INDEX idx_users_company ON users(company_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- =============================================================================
-- PATIENTS (HIMS Data)
-- =============================================================================

CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    patient_number VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE,
    gender VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    address JSONB,
    medical_history JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_patients_company ON patients(company_id);
CREATE INDEX idx_patients_number ON patients(company_id, patient_number);
CREATE INDEX idx_patients_name ON patients USING GIN(to_tsvector('english', name));

-- =============================================================================
-- APPOINTMENTS
-- =============================================================================

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    provider_id UUID REFERENCES users(id),
    appointment_type VARCHAR(100),
    scheduled_at TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER DEFAULT 30,
    status VARCHAR(50) DEFAULT 'scheduled',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (scheduled_at);

-- Create monthly partitions
CREATE TABLE appointments_2026_02 PARTITION OF appointments
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE appointments_2026_03 PARTITION OF appointments
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_provider ON appointments(provider_id);
CREATE INDEX idx_appointments_status ON appointments(company_id, status);

-- =============================================================================
-- TRANSACTIONS (Accounting Data)
-- =============================================================================

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    transaction_type VARCHAR(50) NOT NULL,
    reference_number VARCHAR(100),
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    transaction_date DATE NOT NULL,
    description TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_transactions_company_date ON transactions(company_id, transaction_date DESC);
CREATE INDEX idx_transactions_type ON transactions(company_id, transaction_type);
CREATE INDEX idx_transactions_amount ON transactions(company_id, amount);

-- =============================================================================
-- JOURNAL ENTRIES (Tally Sync)
-- =============================================================================

CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    entry_number VARCHAR(100),
    entry_date DATE NOT NULL,
    description TEXT,
    total_debit DECIMAL(15, 2) NOT NULL,
    total_credit DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    source VARCHAR(100), -- 'manual', 'ocr', 'import'
    source_document_id UUID,
    ocr_confidence DECIMAL(5, 4), -- 0.0000 to 1.0000
    ledger_suggestions JSONB, -- AI-suggested ledger mappings
    approval_status VARCHAR(50) DEFAULT 'pending',
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    tally_sync_status VARCHAR(50) DEFAULT 'pending',
    tally_voucher_no VARCHAR(100),
    tally_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_balanced CHECK (total_debit = total_credit)
);

CREATE INDEX idx_journal_company ON journal_entries(company_id);
CREATE INDEX idx_journal_date ON journal_entries(company_id, entry_date DESC);
CREATE INDEX idx_journal_approval ON journal_entries(company_id, approval_status);
CREATE INDEX idx_journal_tally ON journal_entries(company_id, tally_sync_status);

-- =============================================================================
-- JOURNAL ENTRY LINES
-- =============================================================================

CREATE TABLE journal_entry_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    ledger_code VARCHAR(100) NOT NULL,
    ledger_name VARCHAR(255),
    debit DECIMAL(15, 2),
    credit DECIMAL(15, 2),
    description TEXT,
    cost_center VARCHAR(100),
    metadata JSONB DEFAULT '{}',

    CONSTRAINT chk_debit_or_credit CHECK (
        (debit IS NOT NULL AND credit IS NULL) OR
        (debit IS NULL AND credit IS NOT NULL)
    )
);

CREATE INDEX idx_journal_lines_entry ON journal_entry_lines(journal_entry_id);
CREATE INDEX idx_journal_lines_ledger ON journal_entry_lines(ledger_code);

-- =============================================================================
-- DOCUMENTS (OCR Processing)
-- =============================================================================

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    document_type VARCHAR(50) NOT NULL, -- 'invoice', 'receipt', 'bank_statement'
    original_filename VARCHAR(255),
    storage_path VARCHAR(500),
    mime_type VARCHAR(100),
    file_size INTEGER,
    ocr_status VARCHAR(50) DEFAULT 'pending',
    ocr_result JSONB,
    ocr_confidence DECIMAL(5, 4),
    extracted_data JSONB,
    processing_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_company ON documents(company_id);
CREATE INDEX idx_documents_type ON documents(company_id, document_type);
CREATE INDEX idx_documents_status ON documents(ocr_status);

-- =============================================================================
-- EMBEDDINGS (Semantic Search)
-- =============================================================================

CREATE TABLE document_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(1536),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_embeddings_vector ON document_embeddings
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_embeddings_document ON document_embeddings(document_id);

-- =============================================================================
-- AUDIT LOG
-- =============================================================================

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE audit_log_2026_02 PARTITION OF audit_log
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE INDEX idx_audit_company ON audit_log(company_id);
CREATE INDEX idx_audit_user ON audit_log(user_id);
CREATE INDEX idx_audit_entity ON audit_log(entity_type, entity_id);

-- =============================================================================
-- READ-ONLY ROLE (for AI agents)
-- =============================================================================

CREATE ROLE medisync_readonly WITH LOGIN PASSWORD 'secure_password_here';

GRANT CONNECT ON DATABASE medisync TO medisync_readonly;
GRANT USAGE ON SCHEMA public TO medisync_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO medisync_readonly;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO medisync_readonly;
