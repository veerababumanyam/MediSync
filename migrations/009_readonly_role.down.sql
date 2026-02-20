-- MediSync AI Agent Core - Read-Only Role Permissions Update (Rollback)
-- Version: 009
-- Description: Revoke AI Agent Core specific permissions from medisync_readonly role
-- Task: T012
--
-- This migration rolls back:
-- 1. Default privileges for AI Agent Core tables
-- 2. Verification view
--
-- Note: This does NOT drop the medisync_readonly role itself, as it may be used
-- by other features. It only removes the default privileges added by this migration.

-- ============================================================================
-- DROP VERIFICATION VIEW
-- ============================================================================

DROP VIEW IF EXISTS app.readonly_role_permissions;

-- ============================================================================
-- REVOKE DEFAULT PRIVILEGES
-- ============================================================================

-- Revoke default privileges for public schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE SELECT ON TABLES FROM medisync_readonly;

-- Revoke default privileges for app schema
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON TABLES FROM medisync_readonly;

-- Revoke default privileges for vectors schema
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON TABLES FROM medisync_readonly;

-- Revoke default privileges for sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA app REVOKE SELECT ON SEQUENCES FROM medisync_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA vectors REVOKE SELECT ON SEQUENCES FROM medisync_readonly;

-- ============================================================================
-- Note: Individual table grants are NOT revoked here because:
-- 1. Revoking SELECT from specific tables would break other features
-- 2. The role still needs SELECT on core tables for other functionality
-- 3. If a complete rollback is needed, the DBA should handle this manually
-- ============================================================================

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
