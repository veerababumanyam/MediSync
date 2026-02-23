-- AnySync Database Roles Migration Rollback
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

DROP FUNCTION IF EXISTS create_AnySync_user(TEXT, TEXT, TEXT);

-- ============================================================================
-- REVOKE PRIVILEGES: AnySync_etl
-- ============================================================================

-- Revoke default privileges first
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES FROM AnySync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLES FROM AnySync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE USAGE, SELECT ON SEQUENCES FROM AnySync_etl;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE USAGE, SELECT ON SEQUENCES FROM AnySync_etl;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA hims_analytics FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA tally_analytics FROM AnySync_etl;

-- Revoke specific app schema table privileges
REVOKE ALL PRIVILEGES ON app.etl_state FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON app.etl_quarantine FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON app.etl_quality_report FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON app.audit_log FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON app.notification_queue FROM AnySync_etl;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM AnySync_etl;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM AnySync_etl;
REVOKE USAGE ON SCHEMA tally_analytics FROM AnySync_etl;
REVOKE USAGE ON SCHEMA app FROM AnySync_etl;

-- ============================================================================
-- REVOKE PRIVILEGES: AnySync_app
-- ============================================================================

-- Revoke default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON TABLES FROM AnySync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON TABLES FROM AnySync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM AnySync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM AnySync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE USAGE, SELECT ON SEQUENCES FROM AnySync_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE USAGE, SELECT ON SEQUENCES FROM AnySync_app;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM AnySync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM AnySync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA app FROM AnySync_app;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA vectors FROM AnySync_app;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM AnySync_app;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA vectors FROM AnySync_app;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM AnySync_app;
REVOKE USAGE ON SCHEMA tally_analytics FROM AnySync_app;
REVOKE USAGE ON SCHEMA app FROM AnySync_app;
REVOKE USAGE ON SCHEMA vectors FROM AnySync_app;

-- ============================================================================
-- REVOKE PRIVILEGES: AnySync_readonly
-- ============================================================================

-- Revoke default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON TABLES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON TABLES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON TABLES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON TABLES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA hims_analytics REVOKE SELECT ON SEQUENCES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA tally_analytics REVOKE SELECT ON SEQUENCES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON SEQUENCES FROM AnySync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON SEQUENCES FROM AnySync_readonly;

-- Revoke table privileges
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA hims_analytics FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA tally_analytics FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA app FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA vectors FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA hims_analytics FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA tally_analytics FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app FROM AnySync_readonly;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA vectors FROM AnySync_readonly;

-- Revoke schema usage
REVOKE USAGE ON SCHEMA hims_analytics FROM AnySync_readonly;
REVOKE USAGE ON SCHEMA tally_analytics FROM AnySync_readonly;
REVOKE USAGE ON SCHEMA app FROM AnySync_readonly;
REVOKE USAGE ON SCHEMA vectors FROM AnySync_readonly;

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
    -- Find and drop users that are members of AnySync_etl
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'AnySync_etl'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of AnySync_etl)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;

    -- Find and drop users that are members of AnySync_app
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'AnySync_app'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of AnySync_app)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;

    -- Find and drop users that are members of AnySync_readonly
    FOR r IN SELECT member::regrole::text as member_name
             FROM pg_auth_members
             WHERE roleid = 'AnySync_readonly'::regrole
    LOOP
        RAISE NOTICE 'Dropping user % (member of AnySync_readonly)', r.member_name;
        EXECUTE format('DROP ROLE IF EXISTS %I', r.member_name);
    END LOOP;
EXCEPTION
    WHEN undefined_object THEN
        -- Role doesn't exist, that's fine
        RAISE NOTICE 'Some roles do not exist, continuing...';
END
$$;

-- Now drop the main roles
DROP ROLE IF EXISTS AnySync_etl;
DROP ROLE IF EXISTS AnySync_app;
DROP ROLE IF EXISTS AnySync_readonly;

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- Uncomment to verify roles were dropped:
-- SELECT rolname FROM pg_roles WHERE rolname LIKE 'AnySync_%';

-- ============================================================================
-- END OF ROLLBACK MIGRATION
-- ============================================================================
