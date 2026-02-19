-- MediSync Audit Log Row-Level Security Migration
-- Version: 003
-- Description: Implement append-only policy for audit_log table using Row-Level Security
--
-- This migration establishes:
-- 1. Row-Level Security (RLS) on app.audit_log table
-- 2. Append-only policy: users can INSERT new records but cannot UPDATE or DELETE
-- 3. Read access for all authenticated roles
-- 4. Trigger to prevent any UPDATE operations (belt-and-suspenders approach)
--
-- Security Model:
-- - medisync_readonly: Can SELECT all audit records (for compliance/reporting)
-- - medisync_app: Can SELECT all records, INSERT new records, no UPDATE/DELETE
-- - medisync_etl: Can SELECT all records, INSERT new records, no UPDATE/DELETE
-- - superusers: Bypass RLS (for emergency maintenance only)
--
-- This ensures audit trail integrity and immutability for compliance requirements.

-- ============================================================================
-- ENABLE ROW-LEVEL SECURITY
-- ============================================================================

-- Enable RLS on the audit_log table
ALTER TABLE app.audit_log ENABLE ROW LEVEL SECURITY;

-- Force RLS for table owner as well (important for security)
-- This prevents the table owner from bypassing RLS policies
ALTER TABLE app.audit_log FORCE ROW LEVEL SECURITY;

-- ============================================================================
-- DROP EXISTING POLICIES (if any, for idempotency)
-- ============================================================================

DROP POLICY IF EXISTS audit_log_select_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_insert_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_update_policy ON app.audit_log;
DROP POLICY IF EXISTS audit_log_delete_policy ON app.audit_log;

-- ============================================================================
-- ROW-LEVEL SECURITY POLICIES
-- ============================================================================

-- SELECT Policy: All roles can read all audit log entries
-- This is essential for compliance, reporting, and debugging
CREATE POLICY audit_log_select_policy ON app.audit_log
    FOR SELECT
    TO medisync_readonly, medisync_app, medisync_etl
    USING (true);

-- INSERT Policy: app and etl roles can insert new audit records
-- All inserts are allowed - the audit log captures all actions
CREATE POLICY audit_log_insert_policy ON app.audit_log
    FOR INSERT
    TO medisync_app, medisync_etl
    WITH CHECK (true);

-- UPDATE Policy: DENY ALL - No updates allowed to audit records
-- Using a policy that always returns false
CREATE POLICY audit_log_update_policy ON app.audit_log
    FOR UPDATE
    TO medisync_readonly, medisync_app, medisync_etl
    USING (false)
    WITH CHECK (false);

-- DELETE Policy: DENY ALL - No deletions allowed to audit records
-- Using a policy that always returns false
CREATE POLICY audit_log_delete_policy ON app.audit_log
    FOR DELETE
    TO medisync_readonly, medisync_app, medisync_etl
    USING (false);

-- ============================================================================
-- TRIGGER: Additional UPDATE protection (defense in depth)
-- ============================================================================

-- Create a trigger function that prevents any UPDATE operations
-- This provides defense-in-depth in case RLS is somehow bypassed
CREATE OR REPLACE FUNCTION app.prevent_audit_log_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'UPDATE operations are not permitted on app.audit_log. Audit records are immutable.'
        USING ERRCODE = 'prohibited_sql_statement_attempted',
              HINT = 'Audit log entries cannot be modified after creation for compliance reasons.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger function that prevents any DELETE operations
CREATE OR REPLACE FUNCTION app.prevent_audit_log_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'DELETE operations are not permitted on app.audit_log. Audit records are immutable.'
        USING ERRCODE = 'prohibited_sql_statement_attempted',
              HINT = 'Audit log entries cannot be deleted for compliance reasons. Contact a DBA for archival procedures.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Drop existing triggers if they exist (for idempotency)
DROP TRIGGER IF EXISTS trg_prevent_audit_log_update ON app.audit_log;
DROP TRIGGER IF EXISTS trg_prevent_audit_log_delete ON app.audit_log;

-- Create triggers to block UPDATE and DELETE
-- These fire BEFORE the operation and prevent it from happening
CREATE TRIGGER trg_prevent_audit_log_update
    BEFORE UPDATE ON app.audit_log
    FOR EACH ROW
    EXECUTE FUNCTION app.prevent_audit_log_update();

CREATE TRIGGER trg_prevent_audit_log_delete
    BEFORE DELETE ON app.audit_log
    FOR EACH ROW
    EXECUTE FUNCTION app.prevent_audit_log_delete();

-- ============================================================================
-- REVOKE DIRECT UPDATE/DELETE PRIVILEGES
-- ============================================================================

-- Explicitly revoke UPDATE and DELETE from all application roles
-- This is another layer of defense
REVOKE UPDATE, DELETE ON app.audit_log FROM medisync_readonly;
REVOKE UPDATE, DELETE ON app.audit_log FROM medisync_app;
REVOKE UPDATE, DELETE ON app.audit_log FROM medisync_etl;

-- Ensure only INSERT and SELECT are granted
-- medisync_readonly already has only SELECT from previous migration
-- Re-grant to be explicit about the intended permissions
GRANT SELECT ON app.audit_log TO medisync_readonly;
GRANT SELECT, INSERT ON app.audit_log TO medisync_app;
GRANT SELECT, INSERT ON app.audit_log TO medisync_etl;

-- ============================================================================
-- VALIDATION FUNCTION: Helper to verify audit log integrity
-- ============================================================================

-- Create a function to validate that audit log security is properly configured
CREATE OR REPLACE FUNCTION app.validate_audit_log_security()
RETURNS TABLE (
    check_name TEXT,
    check_status TEXT,
    details TEXT
) AS $$
BEGIN
    -- Check 1: RLS is enabled
    RETURN QUERY
    SELECT
        'RLS Enabled'::TEXT,
        CASE WHEN relrowsecurity THEN 'PASS'::TEXT ELSE 'FAIL'::TEXT END,
        CASE WHEN relrowsecurity
             THEN 'Row-level security is enabled on audit_log'::TEXT
             ELSE 'WARNING: Row-level security is NOT enabled!'::TEXT
        END
    FROM pg_class
    WHERE oid = 'app.audit_log'::regclass;

    -- Check 2: RLS is forced
    RETURN QUERY
    SELECT
        'RLS Forced'::TEXT,
        CASE WHEN relforcerowsecurity THEN 'PASS'::TEXT ELSE 'FAIL'::TEXT END,
        CASE WHEN relforcerowsecurity
             THEN 'Row-level security is forced for table owner'::TEXT
             ELSE 'WARNING: RLS can be bypassed by table owner!'::TEXT
        END
    FROM pg_class
    WHERE oid = 'app.audit_log'::regclass;

    -- Check 3: Policies exist
    RETURN QUERY
    SELECT
        'Policies Exist'::TEXT,
        CASE WHEN COUNT(*) >= 4 THEN 'PASS'::TEXT ELSE 'FAIL'::TEXT END,
        format('%s policies found (expected: 4)', COUNT(*))::TEXT
    FROM pg_policies
    WHERE schemaname = 'app' AND tablename = 'audit_log';

    -- Check 4: Update trigger exists
    RETURN QUERY
    SELECT
        'Update Trigger Exists'::TEXT,
        CASE WHEN EXISTS (
            SELECT 1 FROM pg_trigger t
            JOIN pg_class c ON t.tgrelid = c.oid
            JOIN pg_namespace n ON c.relnamespace = n.oid
            WHERE n.nspname = 'app'
            AND c.relname = 'audit_log'
            AND t.tgname = 'trg_prevent_audit_log_update'
        ) THEN 'PASS'::TEXT ELSE 'FAIL'::TEXT END,
        'Trigger to prevent UPDATE operations'::TEXT;

    -- Check 5: Delete trigger exists
    RETURN QUERY
    SELECT
        'Delete Trigger Exists'::TEXT,
        CASE WHEN EXISTS (
            SELECT 1 FROM pg_trigger t
            JOIN pg_class c ON t.tgrelid = c.oid
            JOIN pg_namespace n ON c.relnamespace = n.oid
            WHERE n.nspname = 'app'
            AND c.relname = 'audit_log'
            AND t.tgname = 'trg_prevent_audit_log_delete'
        ) THEN 'PASS'::TEXT ELSE 'FAIL'::TEXT END,
        'Trigger to prevent DELETE operations'::TEXT;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.validate_audit_log_security() IS 'Validates that audit_log table has proper append-only security configured. Run SELECT * FROM app.validate_audit_log_security();';

-- ============================================================================
-- COMMENTS: Documentation for security policies
-- ============================================================================

COMMENT ON POLICY audit_log_select_policy ON app.audit_log IS
    'Allows all application roles to read audit log entries for compliance and debugging';

COMMENT ON POLICY audit_log_insert_policy ON app.audit_log IS
    'Allows app and ETL roles to create new audit log entries';

COMMENT ON POLICY audit_log_update_policy ON app.audit_log IS
    'DENIES all UPDATE operations - audit records are immutable';

COMMENT ON POLICY audit_log_delete_policy ON app.audit_log IS
    'DENIES all DELETE operations - audit records cannot be removed';

COMMENT ON FUNCTION app.prevent_audit_log_update() IS
    'Trigger function that prevents UPDATE operations on audit_log table';

COMMENT ON FUNCTION app.prevent_audit_log_delete() IS
    'Trigger function that prevents DELETE operations on audit_log table';

COMMENT ON TRIGGER trg_prevent_audit_log_update ON app.audit_log IS
    'Prevents any UPDATE operations on audit_log (defense in depth)';

COMMENT ON TRIGGER trg_prevent_audit_log_delete ON app.audit_log IS
    'Prevents any DELETE operations on audit_log (defense in depth)';

-- ============================================================================
-- VERIFICATION (run manually after migration)
-- ============================================================================

-- Run this to verify security is properly configured:
-- SELECT * FROM app.validate_audit_log_security();

-- Test INSERT (should succeed):
-- SET ROLE medisync_app;
-- INSERT INTO app.audit_log (action, resource) VALUES ('test', 'migration_test');

-- Test UPDATE (should fail):
-- SET ROLE medisync_app;
-- UPDATE app.audit_log SET action = 'modified' WHERE action = 'test';

-- Test DELETE (should fail):
-- SET ROLE medisync_app;
-- DELETE FROM app.audit_log WHERE action = 'test';

-- Reset role after testing:
-- RESET ROLE;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
