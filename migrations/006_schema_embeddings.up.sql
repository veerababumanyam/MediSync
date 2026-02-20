-- MediSync AI Agent Core - Schema Embeddings Migration
-- Version: 006
-- Description: Create schema embeddings table for semantic search in AI Agent Core
-- Task: T009
--
-- This migration establishes:
-- 1. vectors.ai_schema_embeddings - Vector representations of schema elements for semantic search
--
-- Note: This creates a new table specifically for AI Agent Core, separate from the existing
-- vectors.schema_embeddings table which has a different structure for the general schema context.
-- The ai_schema_embeddings table follows the data-model.md specification.
--
-- HNSW index is created for fast cosine similarity search.

-- ============================================================================
-- Ensure vectors schema exists
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS vectors;

-- ============================================================================
-- Ensure pgvector extension is available
-- ============================================================================

CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================================================
-- TABLE: vectors.ai_schema_embeddings
-- Purpose: Vector representations of table/column descriptions for semantic search
-- Note: This is the AI Agent Core specific table following data-model.md spec
-- ============================================================================

CREATE TABLE IF NOT EXISTS vectors.ai_schema_embeddings (
    id SERIAL PRIMARY KEY,
    embedding_type VARCHAR(20) NOT NULL,
    entity_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    embedding vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_ai_schema_embeddings_type CHECK (embedding_type IN ('table', 'column', 'query_pattern'))
);

-- Indexes for ai_schema_embeddings
CREATE INDEX IF NOT EXISTS idx_ai_schema_embeddings_type ON vectors.ai_schema_embeddings(embedding_type);
CREATE INDEX IF NOT EXISTS idx_ai_schema_embeddings_entity_name ON vectors.ai_schema_embeddings(entity_name);

-- Create HNSW index for fast cosine similarity search
-- Note: Using vector_cosine_ops for cosine similarity
CREATE INDEX IF NOT EXISTS idx_ai_schema_embeddings_vector ON vectors.ai_schema_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);

COMMENT ON TABLE vectors.ai_schema_embeddings IS 'AI Agent Core schema embeddings for semantic search during Text-to-SQL generation';

COMMENT ON COLUMN vectors.ai_schema_embeddings.embedding_type IS 'Type of schema element: table, column, or query_pattern';
COMMENT ON COLUMN vectors.ai_schema_embeddings.entity_name IS 'Name of the table, column, or pattern identifier';
COMMENT ON COLUMN vectors.ai_schema_embeddings.description IS 'Natural language description of the schema element';
COMMENT ON COLUMN vectors.ai_schema_embeddings.metadata IS 'Additional context like data_type, sample_values, business_domain';
COMMENT ON COLUMN vectors.ai_schema_embeddings.embedding IS '1536-dimensional embedding vector (OpenAI ada-002 compatible)';

-- ============================================================================
-- FUNCTION: Update updated_at timestamp
-- Purpose: Trigger function to maintain updated_at column
-- ============================================================================

CREATE OR REPLACE FUNCTION vectors.update_ai_schema_embeddings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ai_schema_embeddings_updated_at ON vectors.ai_schema_embeddings;
CREATE TRIGGER trg_ai_schema_embeddings_updated_at
    BEFORE UPDATE ON vectors.ai_schema_embeddings
    FOR EACH ROW
    EXECUTE FUNCTION vectors.update_ai_schema_embeddings_updated_at();

COMMENT ON FUNCTION vectors.update_ai_schema_embeddings_updated_at() IS 'Trigger function to update updated_at timestamp on ai_schema_embeddings';

-- ============================================================================
-- FUNCTION: Semantic search for schema elements
-- Purpose: Find similar schema elements by embedding vector
-- ============================================================================

CREATE OR REPLACE FUNCTION vectors.search_schema_elements(
    p_query_embedding vector(1536),
    p_match_count INTEGER DEFAULT 5,
    p_embedding_type VARCHAR(20) DEFAULT NULL
)
RETURNS TABLE (
    id INTEGER,
    embedding_type VARCHAR(20),
    entity_name VARCHAR(255),
    description TEXT,
    metadata JSONB,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        ase.id,
        ase.embedding_type,
        ase.entity_name,
        ase.description,
        ase.metadata,
        1 - (ase.embedding <=> p_query_embedding) AS similarity
    FROM vectors.ai_schema_embeddings ase
    WHERE
        (p_embedding_type IS NULL OR ase.embedding_type = p_embedding_type)
        AND ase.embedding IS NOT NULL
    ORDER BY ase.embedding <=> p_query_embedding
    LIMIT p_match_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR) IS 'Search for similar schema elements using cosine similarity. Returns id, type, name, description, metadata, and similarity score.';

-- ============================================================================
-- FUNCTION: Upsert schema embedding
-- Purpose: Insert or update a schema embedding
-- ============================================================================

CREATE OR REPLACE FUNCTION vectors.upsert_schema_embedding(
    p_embedding_type VARCHAR(20),
    p_entity_name VARCHAR(255),
    p_description TEXT,
    p_metadata JSONB DEFAULT '{}',
    p_embedding vector(1536) DEFAULT NULL
)
RETURNS INTEGER AS $$
DECLARE
    v_id INTEGER;
BEGIN
    -- Check if embedding exists
    SELECT id INTO v_id
    FROM vectors.ai_schema_embeddings
    WHERE embedding_type = p_embedding_type AND entity_name = p_entity_name;

    IF v_id IS NOT NULL THEN
        -- Update existing record
        UPDATE vectors.ai_schema_embeddings
        SET
            description = p_description,
            metadata = p_metadata,
            embedding = COALESCE(p_embedding, embedding),
            updated_at = NOW()
        WHERE id = v_id;
        RETURN v_id;
    ELSE
        -- Insert new record
        INSERT INTO vectors.ai_schema_embeddings (
            embedding_type, entity_name, description, metadata, embedding, created_at, updated_at
        ) VALUES (
            p_embedding_type, p_entity_name, p_description, p_metadata, p_embedding, NOW(), NOW()
        )
        RETURNING id INTO v_id;
        RETURN v_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION vectors.upsert_schema_embedding(VARCHAR(20), VARCHAR(255), TEXT, JSONB, vector(1536)) IS 'Insert or update a schema embedding. Returns the ID of the upserted record.';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

-- Grant SELECT to medisync_readonly role for AI agents
GRANT SELECT ON vectors.ai_schema_embeddings TO medisync_readonly;

-- Grant full access to medisync_app role for embedding management
GRANT SELECT, INSERT, UPDATE, DELETE ON vectors.ai_schema_embeddings TO medisync_app;
GRANT USAGE, SELECT ON SEQUENCE vectors.ai_schema_embeddings_id_seq TO medisync_app;

-- Grant execute on functions
GRANT EXECUTE ON FUNCTION vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR) TO medisync_app;
GRANT EXECUTE ON FUNCTION vectors.search_schema_elements(vector(1536), INTEGER, VARCHAR) TO medisync_readonly;
GRANT EXECUTE ON FUNCTION vectors.upsert_schema_embedding(VARCHAR(20), VARCHAR(255), TEXT, JSONB, vector(1536)) TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
