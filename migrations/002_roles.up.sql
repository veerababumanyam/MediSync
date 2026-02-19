-- MediSync Database Roles Migration
-- Version: 002
-- Description: Create application database roles with appropriate privileges
-- Roles: medisync_readonly, medisync_app, medisync_etl
--
-- This migration establishes:
-- 1. Three application roles with distinct permission levels
-- 2. Schema-level privileges for each role
-- 3. Table-level privileges following principle of least privilege
-- 4. Default privileges for future tables
--
-- Role Hierarchy:
-- - medisync_readonly: Read-only access to all analytics schemas (for AI agents, analysts)
-- - medisync_app: Read/write access to app schema, read access to analytics (for application)
-- - medisync_etl: Full access to analytics schemas for ETL operations (for ETL service)

-- ============================================================================
-- ROLE: medisync_readonly
-- Purpose: Read-only access for AI agents, analysts, and reporting tools
-- Access: SELECT on all schemas (hims_analytics, tally_analytics, app, vectors)
-- ============================================================================

-- Create role if not exists (idempotent)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'medisync_readonly') THEN
        CREATE ROLE medisync_readonly WITH NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
        RAISE NOTICE 'Created role: medisync_readonly';
    ELSE
        RAISE NOTICE 'Role medisync_readonly already exists, skipping creation';
    END IF;
END
$$;

-- Grant usage on all schemas
GRANT USAGE ON SCHEMA hims_analytics TO medisync_readonly;
GRANT USAGE ON SCHEMA tally_analytics TO medisync_readonly;
GRANT USAGE ON SCHEMA app TO medisync_readonly;
GRANT USAGE ON SCHEMA vectors TO medisync_readonly;

-- Grant SELECT on all existing tables in hims_analytics
GRANT SELECT ON ALL TABLES IN SCHEMA hims_analytics TO medisync_readonly;

-- Grant SELECT on all existing tables in tally_analytics
GRANT SELECT ON ALL TABLES IN SCHEMA tally_analytics TO medisync_readonly;

-- Grant SELECT on all existing tables in app
GRANT SELECT ON ALL TABLES IN SCHEMA app TO medisync_readonly;

-- Grant SELECT on all existing tables in vectors
GRANT SELECT ON ALL TABLES IN SCHEMA vectors TO medisync_readonly;

-- Set default privileges for future tables (so new tables are automatically accessible)
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics GRANT SELECT ON TABLES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics GRANT SELECT ON TABLES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON TABLES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT SELECT ON TABLES TO medisync_readonly;

-- Grant SELECT on sequences (for ID reference lookup)
GRANT SELECT ON ALL SEQUENCES IN SCHEMA hims_analytics TO medisync_readonly;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA tally_analytics TO medisync_readonly;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_readonly;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA vectors TO medisync_readonly;

ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics GRANT SELECT ON SEQUENCES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics GRANT SELECT ON SEQUENCES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON SEQUENCES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT SELECT ON SEQUENCES TO medisync_readonly;

COMMENT ON ROLE medisync_readonly IS 'Read-only access for AI agents, analysts, and reporting tools. No write permissions.';

-- ============================================================================
-- ROLE: medisync_app
-- Purpose: Application backend role for normal operations
-- Access:
--   - Full CRUD on app schema (users, preferences, workflows, etc.)
--   - Read-only on analytics schemas (hims_analytics, tally_analytics)
--   - Read/write on vectors schema (for embedding updates)
--   - INSERT-only on app.audit_log (append-only)
-- ============================================================================

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'medisync_app') THEN
        CREATE ROLE medisync_app WITH NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
        RAISE NOTICE 'Created role: medisync_app';
    ELSE
        RAISE NOTICE 'Role medisync_app already exists, skipping creation';
    END IF;
END
$$;

-- Grant usage on all schemas
GRANT USAGE ON SCHEMA hims_analytics TO medisync_app;
GRANT USAGE ON SCHEMA tally_analytics TO medisync_app;
GRANT USAGE ON SCHEMA app TO medisync_app;
GRANT USAGE ON SCHEMA vectors TO medisync_app;

-- Analytics schemas: READ-ONLY access
GRANT SELECT ON ALL TABLES IN SCHEMA hims_analytics TO medisync_app;
GRANT SELECT ON ALL TABLES IN SCHEMA tally_analytics TO medisync_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics GRANT SELECT ON TABLES TO medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics GRANT SELECT ON TABLES TO medisync_app;

-- App schema: FULL CRUD access (except audit_log which is INSERT-only)
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT USAGE, SELECT ON SEQUENCES TO medisync_app;

-- Vectors schema: READ/WRITE access (for embedding management)
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA vectors TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA vectors TO medisync_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT USAGE, SELECT ON SEQUENCES TO medisync_app;

-- Special handling for audit_log: INSERT-only (no UPDATE/DELETE)
-- This will be enforced via row-level security in 003_audit_log_security.up.sql
-- For now, the app role has full access, and RLS will restrict it

COMMENT ON ROLE medisync_app IS 'Application backend role. Full CRUD on app/vectors schemas, read-only on analytics schemas.';

-- ============================================================================
-- ROLE: medisync_etl
-- Purpose: ETL service role for data synchronization
-- Access:
--   - Full access to hims_analytics and tally_analytics (INSERT, UPDATE, DELETE, TRUNCATE)
--   - Read/write on app.etl_state, app.etl_quarantine, app.etl_quality_report
--   - INSERT on app.audit_log (for ETL audit entries)
--   - No access to user data in app schema (privacy boundary)
-- ============================================================================

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'medisync_etl') THEN
        CREATE ROLE medisync_etl WITH NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
        RAISE NOTICE 'Created role: medisync_etl';
    ELSE
        RAISE NOTICE 'Role medisync_etl already exists, skipping creation';
    END IF;
END
$$;

-- Grant usage on required schemas
GRANT USAGE ON SCHEMA hims_analytics TO medisync_etl;
GRANT USAGE ON SCHEMA tally_analytics TO medisync_etl;
GRANT USAGE ON SCHEMA app TO medisync_etl;

-- HIMS Analytics: Full DML access for ETL operations
GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA hims_analytics TO medisync_etl;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA hims_analytics TO medisync_etl;

ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO medisync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics GRANT USAGE, SELECT ON SEQUENCES TO medisync_etl;

-- Tally Analytics: Full DML access for ETL operations
GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA tally_analytics TO medisync_etl;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA tally_analytics TO medisync_etl;

ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO medisync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics GRANT USAGE, SELECT ON SEQUENCES TO medisync_etl;

-- App schema: Limited access to ETL-related tables only
-- ETL State management
GRANT SELECT, INSERT, UPDATE, DELETE ON app.etl_state TO medisync_etl;

-- ETL Quarantine management
GRANT SELECT, INSERT, UPDATE, DELETE ON app.etl_quarantine TO medisync_etl;

-- ETL Quality Reports
GRANT SELECT, INSERT ON app.etl_quality_report TO medisync_etl;

-- Audit log: INSERT only (for ETL audit trail)
GRANT INSERT ON app.audit_log TO medisync_etl;

-- Notification queue: INSERT for ETL alerts
GRANT SELECT, INSERT ON app.notification_queue TO medisync_etl;

-- Grant sequence usage for app schema tables ETL needs
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_etl;

COMMENT ON ROLE medisync_etl IS 'ETL service role. Full access to analytics schemas, limited access to ETL-related app tables.';

-- ============================================================================
-- HELPER FUNCTION: Create login user from role
-- Use this function to create actual login users that inherit from these roles
-- Example: SELECT create_medisync_user('etl_service', 'medisync_etl', 'secure_password');
-- ============================================================================

CREATE OR REPLACE FUNCTION create_medisync_user(
    p_username TEXT,
    p_role TEXT,
    p_password TEXT
) RETURNS TEXT AS $$
DECLARE
    v_allowed_roles TEXT[] := ARRAY['medisync_readonly', 'medisync_app', 'medisync_etl'];
BEGIN
    -- Validate role
    IF NOT p_role = ANY(v_allowed_roles) THEN
        RAISE EXCEPTION 'Invalid role: %. Allowed roles: %', p_role, v_allowed_roles;
    END IF;

    -- Check if user already exists
    IF EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = p_username) THEN
        RAISE EXCEPTION 'User % already exists', p_username;
    END IF;

    -- Create user with login privilege
    EXECUTE format('CREATE ROLE %I WITH LOGIN PASSWORD %L INHERIT IN ROLE %I',
                   p_username, p_password, p_role);

    RETURN format('Created user %s with role %s', p_username, p_role);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Only allow superusers to create users
REVOKE ALL ON FUNCTION create_medisync_user(TEXT, TEXT, TEXT) FROM PUBLIC;

COMMENT ON FUNCTION create_medisync_user IS 'Helper function to create login users that inherit from MediSync roles. Only callable by superusers.';

-- ============================================================================
-- VERIFICATION QUERIES (for manual verification)
-- These can be run to verify role setup is correct
-- ============================================================================

-- Uncomment to verify roles were created:
-- SELECT rolname, rolsuper, rolinherit, rolcreaterole, rolcreatedb, rolcanlogin
-- FROM pg_roles
-- WHERE rolname LIKE 'medisync_%';

-- Uncomment to verify schema privileges:
-- SELECT nspname, has_schema_privilege('medisync_readonly', nspname, 'USAGE') as readonly_usage,
--        has_schema_privilege('medisync_app', nspname, 'USAGE') as app_usage,
--        has_schema_privilege('medisync_etl', nspname, 'USAGE') as etl_usage
-- FROM pg_namespace
-- WHERE nspname IN ('hims_analytics', 'tally_analytics', 'app', 'vectors');

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
