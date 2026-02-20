-- MediSync AI Agent Core - Confidence Scoring and Review Queue Migration (Rollback)
-- Version: 005
-- Description: Rollback confidence scoring and review queue tables
-- Task: T008
--
-- This migration rolls back:
-- 1. app.review_queue
-- 2. app.confidence_scores
-- 3. Related triggers and functions

-- ============================================================================
-- REVOKE PERMISSIONS
-- ============================================================================

REVOKE SELECT ON app.review_queue FROM medisync_readonly;
REVOKE SELECT ON app.confidence_scores FROM medisync_readonly;

REVOKE SELECT, INSERT, UPDATE, DELETE ON app.review_queue FROM medisync_app;
REVOKE SELECT, INSERT, UPDATE, DELETE ON app.confidence_scores FROM medisync_app;

REVOKE EXECUTE ON FUNCTION app.get_pending_review_count(UUID) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION app.get_pending_review_count(UUID) FROM medisync_readonly;

-- ============================================================================
-- DROP TRIGGERS AND FUNCTIONS
-- ============================================================================

DROP TRIGGER IF EXISTS trg_auto_route_low_confidence ON app.confidence_scores;
DROP FUNCTION IF EXISTS app.auto_route_low_confidence();
DROP FUNCTION IF EXISTS app.get_pending_review_count(UUID);

-- ============================================================================
-- DROP TABLES (in reverse dependency order)
-- ============================================================================

DROP TABLE IF EXISTS app.review_queue;
DROP TABLE IF EXISTS app.confidence_scores;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
