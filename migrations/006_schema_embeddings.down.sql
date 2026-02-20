-- MediSync AI Agent Core - Schema Embeddings Migration (Rollback)
-- Version: 006
-- Description: Rollback schema embeddings table for AI Agent Core
-- Task: T009
--
-- This migration rolls back:
-- 1. vectors.ai_schema_embeddings table
-- 2. Related triggers and functions

-- ============================================================================
-- REVOKE PERMISSIONS
-- ============================================================================

REVOKE SELECT ON vectors.ai_schema_embeddings FROM medisync_readonly;
REVOKE SELECT, INSERT, UPDATE, DELETE ON vectors.ai_schema_embeddings FROM medisync_app;
REVOKE USAGE, SELECT ON SEQUENCE vectors.ai_schema_embeddings_id_seq FROM medisync_app;

REVOKE EXECUTE ON FUNCTION vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR) FROM medisync_app;
REVOKE EXECUTE ON FUNCTION vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR) FROM medisync_readonly;
REVOKE EXECUTE ON FUNCTION vectors.upsert_schema_embedding(VARCHAR(20), VARCHAR(255), TEXT, JSONB, vector(1536)) FROM medisync_app;

-- ============================================================================
-- DROP TRIGGERS AND FUNCTIONS
-- ============================================================================

DROP TRIGGER IF EXISTS trg_ai_schema_embeddings_updated_at ON vectors.ai_schema_embeddings;
DROP FUNCTION IF EXISTS vectors.update_ai_schema_embeddings_updated_at();
DROP FUNCTION IF EXISTS vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR);
DROP FUNCTION IF EXISTS vectors.upsert_schema_embedding(VARCHAR(20), VARCHAR(255), TEXT, JSONB, vector(1536));

-- ============================================================================
-- DROP TABLE
-- ============================================================================

DROP TABLE IF EXISTS vectors.ai_schema_embeddings;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
