-- MediSync AI Agent Core - Application Tables Migration (Rollback)
-- Version: 004
-- Description: Rollback AI Agent Core tables
-- Task: T007
--
-- This migration rolls back:
-- 1. app.query_results
-- 2. app.sql_statements
-- 3. app.queries
-- 4. app.query_sessions
-- 5. Related triggers and functions

-- ============================================================================
-- REVOKE PERMISSIONS
-- ============================================================================

REVOKE SELECT ON app.query_results FROM medisync_readonly;
REVOKE SELECT ON app.sql_statements FROM medisync_readonly;
REVOKE SELECT ON app.queries FROM medisync_readonly;
REVOKE SELECT ON app.query_sessions FROM medisync_readonly;

REVOKE SELECT, INSERT, UPDATE, DELETE ON app.query_results FROM medisync_app;
REVOKE SELECT, INSERT, UPDATE, DELETE ON app.sql_statements FROM medisync_app;
REVOKE SELECT, INSERT, UPDATE, DELETE ON app.queries FROM medisync_app;
REVOKE SELECT, INSERT, UPDATE, DELETE ON app.query_sessions FROM medisync_app;

-- ============================================================================
-- DROP TRIGGERS AND FUNCTIONS
-- ============================================================================

DROP TRIGGER IF EXISTS trg_query_sessions_updated_at ON app.query_sessions;
DROP FUNCTION IF EXISTS app.update_query_session_updated_at();

-- ============================================================================
-- DROP TABLES (in reverse dependency order)
-- ============================================================================

DROP TABLE IF EXISTS app.query_results;
DROP TABLE IF EXISTS app.sql_statements;
DROP TABLE IF EXISTS app.queries;
DROP TABLE IF EXISTS app.query_sessions;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
