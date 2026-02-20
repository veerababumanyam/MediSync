-- MediSync AI Agent Core - Domain Terms Migration (Rollback)
-- Version: 007
-- Description: Rollback domain terms table for AI Agent Core
-- Task: T010
--
-- This migration rolls back:
-- 1. app.domain_terms table
-- 2. Related functions

-- ============================================================================
-- REVOKE PERMISSIONS
-- ============================================================================

REVOKE SELECT ON app.domain_terms FROM medisync_readonly;
REVOKE SELECT, INSERT, UPDATE, DELETE ON app.domain_terms FROM medisync_app;
REVOKE USAGE, SELECT ON SEQUENCE app.domain_terms_id_seq FROM medisync_app;

REVOKE EXECUTE ON FUNCTION app.find_canonical_term(VARCHAR) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.find_canonical_term(VARCHAR) FROM medisync_readonly;
REVOKE EXECUTE ON FUNCTION app.search_domain_terms(VARCHAR, VARCHAR, INTEGER) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.search_domain_terms(VARCHAR, VARCHAR, INTEGER) FROM medisync_readonly;
REVOKE EXECUTE ON FUNCTION app.get_domain_terms_for_locale(VARCHAR, VARCHAR) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.get_domain_terms_for_locale(VARCHAR, VARCHAR) FROM medisync_readonly;
REVOKE EXECUTE ON FUNCTION app.upsert_domain_term(VARCHAR, VARCHAR, VARCHAR, TEXT, JSONB) FROM medisync_app;

-- ============================================================================
-- DROP FUNCTIONS
-- ============================================================================

DROP FUNCTION IF EXISTS app.find_canonical_term(VARCHAR);
DROP FUNCTION IF EXISTS app.search_domain_terms(VARCHAR, VARCHAR, INTEGER);
DROP FUNCTION IF EXISTS app.get_domain_terms_for_locale(VARCHAR, VARCHAR);
DROP FUNCTION IF EXISTS app.upsert_domain_term(VARCHAR, VARCHAR, VARCHAR, TEXT, JSONB);

-- ============================================================================
-- DROP TABLE
-- ============================================================================

DROP TABLE IF EXISTS app.domain_terms;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
