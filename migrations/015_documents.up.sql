-- MediSync Document Processing Pipeline - Documents Migration
-- Version: 015
-- Description: Create document processing tables for OCR pipeline
-- Task: T001
--
-- This migration establishes:
-- - documents: Uploaded financial documents
-- - extracted_fields: OCR-extracted field data
-- - line_items: Document line items (invoice items, transactions)
-- - document_audit_log: Audit trail for all document actions

-- ============================================================================
-- ENUM TYPES
-- ============================================================================

-- Document status enum
CREATE TYPE app.document_status AS ENUM (
    'uploading',
    'uploaded',
    'classifying',
    'extracting',
    'ready_for_review',
    'under_review',
    'reviewed',
    'approved',
    'rejected',
    'failed'
);

-- Document type enum
CREATE TYPE app.document_type AS ENUM (
    'invoice',
    'receipt',
    'bank_statement',
    'expense_report',
    'credit_note',
    'debit_note',
    'other'
);

-- File format enum
CREATE TYPE app.file_format AS ENUM (
    'pdf',
    'jpeg',
    'png',
    'tiff',
    'xlsx',
    'csv'
);

-- Field type enum
CREATE TYPE app.field_type AS ENUM (
    'string',
    'number',
    'currency',
    'date',
    'percentage',
    'identifier',
    'tax_id'
);

-- Verification status enum
CREATE TYPE app.verification_status AS ENUM (
    'pending',
    'auto_accepted',
    'needs_review',
    'high_priority',
    'manually_verified',
    'manually_corrected',
    'rejected'
);

-- Audit action enum
CREATE TYPE app.audit_action AS ENUM (
    'uploaded',
    'classified',
    'extracted',
    'review_started',
    'field_edited',
    'field_verified',
    'approved',
    'rejected',
    'reprocessed'
);

-- Actor type enum
CREATE TYPE app.actor_type AS ENUM (
    'user',
    'system'
);

-- ============================================================================
-- TABLE: app.documents
-- Purpose: Stores uploaded documents awaiting processing
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    uploaded_by UUID NOT NULL,
    status app.document_status NOT NULL DEFAULT 'uploading',
    document_type app.document_type,
    original_filename VARCHAR(255) NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    file_format app.file_format NOT NULL,
    page_count INTEGER DEFAULT 1,
    detected_language VARCHAR(2) DEFAULT 'en',
    upload_id UUID,
    processing_started_at TIMESTAMPTZ,
    processing_completed_at TIMESTAMPTZ,
    classification_confidence DECIMAL(5, 4),
    overall_confidence DECIMAL(5, 4),
    rejection_reason TEXT,
    locked_by UUID,
    locked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_documents_file_size CHECK (file_size_bytes > 0 AND file_size_bytes <= 26214400),
    CONSTRAINT ck_documents_page_count CHECK (page_count > 0 AND page_count <= 20),
    CONSTRAINT ck_documents_classification_confidence CHECK (classification_confidence IS NULL OR (classification_confidence >= 0.0 AND classification_confidence <= 1.0)),
    CONSTRAINT ck_documents_overall_confidence CHECK (overall_confidence IS NULL OR (overall_confidence >= 0.0 AND overall_confidence <= 1.0))
);

-- Indexes for documents
CREATE INDEX IF NOT EXISTS idx_documents_tenant_id ON app.documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON app.documents(status);
CREATE INDEX IF NOT EXISTS idx_documents_uploaded_by ON app.documents(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_documents_created_at ON app.documents(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_documents_tenant_status ON app.documents(tenant_id, status)
    WHERE status IN ('ready_for_review', 'under_review');
CREATE INDEX IF NOT EXISTS idx_documents_priority ON app.documents(tenant_id, overall_confidence)
    WHERE status = 'ready_for_review';
CREATE INDEX IF NOT EXISTS idx_documents_upload_id ON app.documents(upload_id)
    WHERE upload_id IS NOT NULL;

COMMENT ON TABLE app.documents IS 'Uploaded documents awaiting OCR processing and review';

-- ============================================================================
-- TABLE: app.extracted_fields
-- Purpose: Stores individual fields extracted from documents
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.extracted_fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES app.documents(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL DEFAULT 1,
    field_name VARCHAR(100) NOT NULL,
    field_type app.field_type NOT NULL DEFAULT 'string',
    extracted_value TEXT,
    confidence_score DECIMAL(5, 4),
    bounding_box JSONB,
    is_handwritten BOOLEAN DEFAULT FALSE,
    verification_status app.verification_status NOT NULL DEFAULT 'pending',
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    original_value TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_extracted_fields_page_number CHECK (page_number > 0),
    CONSTRAINT ck_extracted_fields_confidence CHECK (confidence_score IS NULL OR (confidence_score >= 0.0 AND confidence_score <= 1.0)),
    CONSTRAINT ck_extracted_fields_handwritten_confidence CHECK (is_handwritten = FALSE OR confidence_score <= 0.85)
);

-- Indexes for extracted_fields
CREATE INDEX IF NOT EXISTS idx_extracted_fields_document_id ON app.extracted_fields(document_id);
CREATE INDEX IF NOT EXISTS idx_extracted_fields_verification ON app.extracted_fields(document_id, verification_status);
CREATE INDEX IF NOT EXISTS idx_extracted_fields_low_confidence ON app.extracted_fields(document_id, confidence_score)
    WHERE confidence_score < 0.70;

COMMENT ON TABLE app.extracted_fields IS 'Individual fields extracted from documents via OCR';

-- ============================================================================
-- TABLE: app.line_items
-- Purpose: Stores line items from invoices and transactions from statements
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES app.documents(id) ON DELETE CASCADE,
    extracted_field_id UUID REFERENCES app.extracted_fields(id) ON DELETE SET NULL,
    line_number INTEGER NOT NULL,
    description TEXT,
    quantity DECIMAL(12, 4),
    unit_price DECIMAL(15, 4),
    amount DECIMAL(15, 4),
    tax_rate DECIMAL(5, 4),
    transaction_date DATE,
    reference VARCHAR(100),
    debit_amount DECIMAL(15, 4),
    credit_amount DECIMAL(15, 4),
    balance DECIMAL(15, 4),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_line_items_line_number CHECK (line_number > 0),
    CONSTRAINT ck_line_items_amounts CHECK (
        (debit_amount IS NULL OR debit_amount >= 0) AND
        (credit_amount IS NULL OR credit_amount >= 0) AND
        ((debit_amount > 0 AND credit_amount = 0) OR (credit_amount > 0 AND debit_amount = 0) OR (debit_amount IS NULL AND credit_amount IS NULL))
    )
);

-- Indexes for line_items
CREATE INDEX IF NOT EXISTS idx_line_items_document_id ON app.line_items(document_id);

COMMENT ON TABLE app.line_items IS 'Line items from invoices and transactions from bank statements';

-- ============================================================================
-- TABLE: app.document_audit_log
-- Purpose: Audit trail for all document actions
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.document_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL REFERENCES app.documents(id) ON DELETE CASCADE,
    action app.audit_action NOT NULL,
    actor_id UUID,
    actor_type app.actor_type NOT NULL DEFAULT 'user',
    field_name VARCHAR(100),
    old_value JSONB,
    new_value JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for document_audit_log
CREATE INDEX IF NOT EXISTS idx_audit_log_document_id ON app.document_audit_log(document_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_tenant_created ON app.document_audit_log(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_action ON app.document_audit_log(action);

COMMENT ON TABLE app.document_audit_log IS 'Audit trail for all document processing actions';

-- ============================================================================
-- ROW LEVEL SECURITY
-- ============================================================================

ALTER TABLE app.documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.extracted_fields ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.line_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.document_audit_log ENABLE ROW LEVEL SECURITY;

-- Users can only see documents from their tenant
CREATE POLICY documents_tenant_isolation ON app.documents
    USING (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY extracted_fields_tenant_isolation ON app.extracted_fields
    USING (document_id IN (
        SELECT id FROM app.documents WHERE tenant_id = current_setting('app.current_tenant', true)::uuid
    ));

CREATE POLICY line_items_tenant_isolation ON app.line_items
    USING (document_id IN (
        SELECT id FROM app.documents WHERE tenant_id = current_setting('app.current_tenant', true)::uuid
    ));

CREATE POLICY document_audit_log_tenant_isolation ON app.document_audit_log
    USING (tenant_id = current_setting('app.current_tenant', true)::uuid);

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.documents TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.documents TO medisync_app;

GRANT SELECT ON app.extracted_fields TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.extracted_fields TO medisync_app;

GRANT SELECT ON app.line_items TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.line_items TO medisync_app;

GRANT SELECT ON app.document_audit_log TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.document_audit_log TO medisync_app;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- TRIGGER: Update updated_at timestamp
-- ============================================================================

CREATE OR REPLACE FUNCTION app.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_documents_updated_at
    BEFORE UPDATE ON app.documents
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

CREATE TRIGGER update_extracted_fields_updated_at
    BEFORE UPDATE ON app.extracted_fields
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
