-- MediSync Database Roles Migration Rollback
-- Version: 002
-- Description: Drop application database roles and revoke all privileges
--
-- This rollback:
-- 1. Drops the helper function for user creation
-- 2. Revokes all privileges from roles
-- 3. Revokes default privileges
-- 4. Drops the three application roles
--
-- WARNING: This will remove all role-based access control!
-- Any users inheriting from these roles will lose their privileges.

-- ============================================================================
-- DROP HELPER FUNCTION
-- ============================================================================

DROP FUNCTION IF EXISTS create_medisync_user(TEXT, TEXT, TEXT);

-- ============================================================================
-- REVOKE PRIVILEGES: medisync_etl
-- ============================================================================

-- Revoke default privileges first
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES FROM medisync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES FROM medisync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE USAGE, SELECT ON SEQUENCES FROM medisync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE USAGE, SELECT ON SEQUENCES FROM medisync_etl;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM medisync_etl;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM medisync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA hims_analytics FROM medisync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA tally_analytics FROM medisync_etl;

-- Revoke specific app schema table privileges
REVOKE ALL PRIVILEGES ON app.etl_state FROM medisync_etl;
REVOKE ALL PRIVILEGES ON app.etl_quarantine FROM medisync_etl;
REVOKE ALL PRIVILEGES ON app.etl_quality_report FROM medisync_etl;
REVOKE ALL PRIVILEGES ON app.audit_log FROM medisync_etl;
REVOKE ALL PRIVILEGES ON app.notification_queue FROM medisync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM medisync_etl;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM medisync_etl;
REVOKE USAGE ON SCHEMA tally_analytics FROM medisync_etl;
REVOKE USAGE ON SCHEMA app FROM medisync_etl;

-- ============================================================================
-- REVOKE PRIVILEGES: medisync_app
-- ============================================================================

-- Revoke default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON TABLES FROM medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON TABLES FROM medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE USAGE, SELECT ON SEQUENCES FROM medisync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE USAGE, SELECT ON SEQUENCES FROM medisync_app;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM medisync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM medisync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA app FROM medisync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA vectors FROM medisync_app;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM medisync_app;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA vectors FROM medisync_app;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM medisync_app;
REVOKE USAGE ON SCHEMA tally_analytics FROM medisync_app;
REVOKE USAGE ON SCHEMA app FROM medisync_app;
REVOKE USAGE ON SCHEMA vectors FROM medisync_app;

-- ============================================================================
-- REVOKE PRIVILEGES: medisync_readonly
-- ============================================================================

-- Revoke default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON TABLES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON TABLES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON TABLES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON TABLES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON SEQUENCES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON SEQUENCES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON SEQUENCES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON SEQUENCES FROM medisync_readonly;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA app FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA vectors FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA hims_analytics FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA tally_analytics FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM medisync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA vectors FROM medisync_readonly;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM medisync_readonly;
REVOKE USAGE ON SCHEMA tally_analytics FROM medisync_readonly;
REVOKE USAGE ON SCHEMA app FROM medisync_readonly;
REVOKE USAGE ON SCHEMA vectors FROM medisync_readonly;

-- ============================================================================
-- DROP ROLES
-- Note: Roles can only be dropped if they own no objects and have no members
-- ============================================================================

-- Drop dependent users first (users that inherit from these roles)
-- This is a safety check - in production, you may want to handle this differently
DO $$
DECLARE
    r RECORD;
BEGIN
    -- Find and drop users that are members of medisync_etl
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'medisync_etl'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of medisync_etl)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;

    -- Find and drop users that are members of medisync_app
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'medisync_app'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of medisync_app)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;

    -- Find and drop users that are members of medisync_readonly
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'medisync_readonly'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of medisync_readonly)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;
EXCEPTION
    WHEN undefined_object THEN
        -- Role doesn't exist, that's fine
        RAISE NOTICE 'Some roles do not exist, continuing...';
END
$$;

-- Now drop the main roles
DROP ROLE IF EXISTS medisync_etl;
DROP ROLE IF EXISTS medisync_app;
DROP ROLE IF EXISTS medisync_readonly;

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- Uncomment to verify roles were dropped:
-- SELECT rolname FROM pg_roles WHERE rolname LIKE 'medisync_%';

-- ============================================================================
-- END OF ROLLBACK MIGRATION
-- ============================================================================
