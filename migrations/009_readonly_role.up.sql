-- MediSync AI Agent Core - Read-Only Role Permissions Update
-- Version: 009
-- Description: Update medisync_readonly role for AI Agent Core tables
-- Task: T012
--
-- This migration ensures:
-- 1. medisync_readonly role has SELECT on all required schemas
-- 2. medisync_readonly role has SELECT on all AI Agent Core tables
-- 3. Default privileges are set for future tables
--
-- Note: The medisync_readonly role is already created in migration 002_roles.
-- This migration adds permissions specifically for the AI Agent Core feature.

-- ============================================================================
-- ENSURE ROLE EXISTS (idempotent)
-- ============================================================================

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

-- ============================================================================
-- GRANT USAGE ON SCHEMAS
-- ============================================================================

-- Grant USAGE on public schema (if not already granted)
GRANT USAGE ON SCHEMA public TO medisync_readonly;

-- Grant USAGE on app schema (if not already granted)
GRANT USAGE ON SCHEMA app TO medisync_readonly;

-- Grant USAGE on vectors schema (if not already granted)
GRANT USAGE ON SCHEMA vectors TO medisync_readonly;

-- ============================================================================
-- GRANT SELECT ON ALL EXISTING TABLES IN PUBLIC SCHEMA
-- ============================================================================

-- Grant SELECT on all tables in public schema
DO $$
DECLARE
    tbl record;
BEGIN
    FOR tbl IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
    LOOP
        EXECUTE format('GRANT SELECT ON public.%I TO medisync_readonly', tbl.tablename);
    END LOOP;
END
$$;

-- ============================================================================
-- GRANT SELECT ON ALL AI AGENT CORE TABLES IN APP SCHEMA
-- ============================================================================

-- Grant SELECT on all tables in app schema (includes AI Agent Core tables)
DO $$
DECLARE
    tbl record;
BEGIN
    FOR tbl IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'app'
    LOOP
        EXECUTE format('GRANT SELECT ON app.%I TO medisync_readonly', tbl.tablename);
    END LOOP;
END
$$;

-- ============================================================================
-- GRANT SELECT ON ALL TABLES IN VECTORS SCHEMA
-- ============================================================================

-- Grant SELECT on all tables in vectors schema
DO $$
DECLARE
    tbl record;
BEGIN
    FOR tbl IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'vectors'
    LOOP
        EXECUTE format('GRANT SELECT ON vectors.%I TO medisync_readonly', tbl.tablename);
    END LOOP;
END
$$;

-- ============================================================================
-- GRANT SELECT ON ALL SEQUENCES (for reference lookups)
-- ============================================================================

-- Grant SELECT on all sequences in app schema
DO $$
DECLARE
    seq record;
BEGIN
    FOR seq IN
        SELECT sequencename
        FROM pg_sequences
        WHERE schemaname = 'app'
    LOOP
        EXECUTE format('GRANT SELECT ON app.%I TO medisync_readonly', seq.sequencename);
    END LOOP;
END
$$;

-- Grant SELECT on all sequences in vectors schema
DO $$
DECLARE
    seq record;
BEGIN
    FOR seq IN
        SELECT sequencename
        FROM pg_sequences
        WHERE schemaname = 'vectors'
    LOOP
        EXECUTE format('GRANT SELECT ON vectors.%I TO medisync_readonly', seq.sequencename);
    END LOOP;
END
$$;

-- ============================================================================
-- SET DEFAULT PRIVILEGES FOR FUTURE TABLES
-- ============================================================================

-- Default privileges for public schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO medisync_readonly;

-- Default privileges for app schema
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON TABLES TO medisync_readonly;

-- Default privileges for vectors schema
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT SELECT ON TABLES TO medisync_readonly;

-- Default privileges for sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON SEQUENCES TO medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors GRANT SELECT ON SEQUENCES TO medisync_readonly;

-- ============================================================================
-- VERIFICATION: List all grants to medisync_readonly
-- ============================================================================

-- Create a view for easy verification (optional, can be dropped after verification)
CREATE OR REPLACE VIEW app.readonly_role_permissions AS
SELECT
    'TABLE' AS object_type,
    n.nspname AS schema_name,
    c.relname AS object_name,
    pg_catalog.array_to_string(c.relacl, ', ') AS privileges
FROM pg_catalog.pg_class c
JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
WHERE c.relacl::text LIKE '%medisync_readonly%'
AND c.relkind = 'r'
AND n.nspname IN ('public', 'app', 'vectors')
UNION ALL
SELECT
    'SEQUENCE' AS object_type,
    n.nspname AS schema_name,
    c.relname AS object_name,
    pg_catalog.array_to_string(c.relacl, ', ') AS privileges
FROM pg_catalog.pg_class c
JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
WHERE c.relacl::text LIKE '%medisync_readonly%'
AND c.relkind = 'S'
AND n.nspname IN ('public', 'app', 'vectors')
ORDER BY schema_name, object_type, object_name;

-- Grant SELECT on the verification view to medisync_app for debugging
GRANT SELECT ON app.readonly_role_permissions TO medisync_app;

COMMENT ON VIEW app.readonly_role_permissions IS 'View showing all privileges granted to medisync_readonly role in public, app, and vectors schemas';

-- ============================================================================
-- UPDATE ROLE COMMENT
-- ============================================================================

COMMENT ON ROLE medisync_readonly IS 'Read-only access for AI agents, analysts, and reporting tools. Has SELECT on all tables in public, app, and vectors schemas. No write permissions.';

-- ============================================================================
-- VERIFICATION QUERIES (for manual verification)
-- ============================================================================

-- Uncomment to verify schema privileges:
-- SELECT nspname, has_schema_privilege('medisync_readonly', nspname, 'USAGE') as usage
-- FROM pg_namespace
-- WHERE nspname IN ('public', 'app', 'vectors');

-- Uncomment to verify table privileges:
-- SELECT * FROM app.readonly_role_permissions;

-- Uncomment to verify role attributes:
-- SELECT rolname, rolsuper, rolinherit, rolcreaterole, rolcreatedb, rolcanlogin
-- FROM pg_roles
-- WHERE rolname = 'medisync_readonly';

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
