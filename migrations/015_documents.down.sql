-- MediSync Document Processing Pipeline - Documents Migration
-- Version: 015
-- Description: Rollback document processing tables
-- Task: T001

-- ============================================================================
-- DROP TRIGGERS
-- ============================================================================

DROP TRIGGER IF EXISTS update_extracted_fields_updated_at ON app.extracted_fields;
DROP TRIGGER IF EXISTS update_documents_updated_at ON app.documents;

-- ============================================================================
-- DROP FUNCTION
-- ============================================================================

DROP FUNCTION IF EXISTS app.update_updated_at_column();

-- ============================================================================
-- DROP TABLES
-- ============================================================================

DROP TABLE IF EXISTS app.document_audit_log CASCADE;
DROP TABLE IF EXISTS app.line_items CASCADE;
DROP TABLE IF EXISTS app.extracted_fields CASCADE;
DROP TABLE IF EXISTS app.documents CASCADE;

-- ============================================================================
-- DROP ENUM TYPES
-- ============================================================================

DROP TYPE IF EXISTS app.actor_type CASCADE;
DROP TYPE IF EXISTS app.audit_action CASCADE;
DROP TYPE IF EXISTS app.verification_status CASCADE;
DROP TYPE IF EXISTS app.field_type CASCADE;
DROP TYPE IF EXISTS app.file_format CASCADE;
DROP TYPE IF EXISTS app.document_type CASCADE;
DROP TYPE IF EXISTS app.document_status CASCADE;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
