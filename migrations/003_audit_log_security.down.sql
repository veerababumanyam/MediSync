-- MediSync Audit Log Row-Level Security Migration Rollback
-- Version: 003
-- Description: Remove append-only policy from audit_log table
--
-- This rollback:
-- 1. Drops the RLS policies from audit_log table
-- 2. Removes the protective triggers
-- 3. Drops the trigger functions
-- 4. Disables Row-Level Security on the table
-- 5. Restores default privileges to roles
--
-- WARNING: After this rollback, audit_log will no longer be protected!
-- UPDATE and DELETE operations will be possible.

-- ============================================================================
-- DROP VALIDATION FUNCTION
-- ============================================================================

DROP FUNCTION IF EXISTS app.validate_audit_log_security();

-- ============================================================================
-- DROP TRIGGERS
-- ============================================================================

DROP TRIGGER IF EXISTS trg_prevent_audit_log_update ON app.audit_log;
DROP TRIGGER IF EXISTS trg_prevent_audit_log_delete ON app.audit_log;

-- ============================================================================
-- DROP TRIGGER FUNCTIONS
-- ============================================================================

DROP FUNCTION IF EXISTS app.prevent_audit_log_update();
DROP FUNCTION IF EXISTS app.prevent_audit_log_delete();

-- ============================================================================
-- DROP ROW-LEVEL SECURITY POLICIES
-- ============================================================================

DROP POLICY IF EXISTS audit_log_select_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_insert_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_update_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_delete_policy ON app.audit_log;

-- ============================================================================
-- DISABLE ROW-LEVEL SECURITY
-- ============================================================================

-- Disable forced RLS first
ALTER TABLE app.audit_log NO FORCE ROW LEVEL SECURITY;

-- Disable RLS on the table
ALTER TABLE app.audit_log DISABLE ROW LEVEL SECURITY;

-- ============================================================================
-- RESTORE PRIVILEGES
-- ============================================================================

-- Restore the default privileges from migration 002
-- medisync_readonly: SELECT only (unchanged)
-- medisync_app: Full CRUD (as per original 002 migration before RLS restrictions)
-- medisync_etl: INSERT only (as per original 002 migration)

-- Note: We restore to the state defined in 002_roles.up.sql
-- which grants full CRUD to medisync_app and INSERT to medisync_etl

-- Revoke current grants first for clean slate
REVOKE ALL ON app.audit_log FROM medisync_readonly;
REVOKE ALL ON app.audit_log FROM medisync_app;
REVOKE ALL ON app.audit_log FROM medisync_etl;

-- Re-grant as per original 002_roles.up.sql intent
-- medisync_readonly has SELECT on all tables in app schema
GRANT SELECT ON app.audit_log TO medisync_readonly;

-- medisync_app has full CRUD on app schema tables
GRANT SELECT, INSERT, UPDATE, DELETE ON app.audit_log TO medisync_app;

-- medisync_etl has INSERT only on audit_log (as specified in 002)
GRANT INSERT ON app.audit_log TO medisync_etl;

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- Verify RLS is disabled:
-- SELECT relname, relrowsecurity, relforcerowsecurity
-- FROM pg_class
-- WHERE oid = 'app.audit_log'::regclass;

-- Verify no policies exist:
-- SELECT * FROM pg_policies WHERE schemaname = 'app' AND tablename = 'audit_log';

-- Verify triggers are removed:
-- SELECT tgname FROM pg_trigger t
-- JOIN pg_class c ON t.tgrelid = c.oid
-- JOIN pg_namespace n ON c.relnamespace = n.oid
-- WHERE n.nspname = 'app' AND c.relname = 'audit_log'
-- AND NOT tgisinternal;

-- ============================================================================
-- END OF ROLLBACK MIGRATION
-- ============================================================================
