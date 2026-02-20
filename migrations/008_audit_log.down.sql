-- MediSync AI Agent Core - AI Audit Log Migration (Rollback)
-- Version: 008
-- Description: Rollback AI-specific audit log table for AI Agent Core
-- Task: T011
--
-- This migration rolls back:
-- 1. app.ai_audit_log table
-- 2. Related triggers, policies, and functions

-- ============================================================================
-- REVOKE PERMISSIONS
-- ============================================================================

REVOKE SELECT ON app.ai_audit_log FROM medisync_readonly;
REVOKE SELECT, INSERT ON app.ai_audit_log FROM medisync_app;

REVOKE EXECUTE ON FUNCTION app.log_ai_action(UUID, UUID, VARCHAR, VARCHAR, UUID, JSONB, INET, TEXT) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.get_query_audit_trail(UUID) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.get_query_audit_trail(UUID) FROM medisync_readonly;
REVOKE EXECUTE ON FUNCTION app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ) FROM medisync_readonly;

-- ============================================================================
-- DROP POLICIES AND DISABLE RLS
-- ============================================================================

DROP POLICY IF EXISTS ai_audit_log_select_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_insert_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_update_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_delete_policy ON app.ai_audit_log;

-- ============================================================================
-- DROP TRIGGERS AND FUNCTIONS
-- ============================================================================

DROP TRIGGER IF EXISTS trg_prevent_ai_audit_log_modification ON app.ai_audit_log;
DROP FUNCTION IF EXISTS app.prevent_ai_audit_log_modification();
DROP FUNCTION IF EXISTS app.log_ai_action(UUID, UUID, VARCHAR, VARCHAR, UUID, JSONB, INET, TEXT);
DROP FUNCTION IF EXISTS app.get_query_audit_trail(UUID);
DROP FUNCTION IF EXISTS app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ);

-- ============================================================================
-- DROP TABLE
-- ============================================================================

DROP TABLE IF EXISTS app.ai_audit_log;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
