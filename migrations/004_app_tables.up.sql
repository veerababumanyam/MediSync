-- MediSync AI Agent Core - Application Tables Migration
-- Version: 004
-- Description: Create AI Agent Core tables for query sessions, queries, SQL statements, and results
-- Task: T007
--
-- This migration establishes:
-- 1. app.query_sessions - User chat sessions with context
-- 2. app.queries - Natural language queries from users
-- 3. app.sql_statements - Generated SQL with validation tracking
-- 4. app.query_results - Query execution results with metrics
--
-- These tables support the AI Agent Core feature (001-ai-agent-core)

-- ============================================================================
-- TABLE: app.query_sessions
-- Purpose: Represents a user's chat session containing query history and context
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.query_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Indexes for query_sessions
CREATE INDEX IF NOT EXISTS idx_query_sessions_user_id ON app.query_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_query_sessions_tenant_id ON app.query_sessions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_query_sessions_created_at ON app.query_sessions(created_at);

COMMENT ON TABLE app.query_sessions IS 'User chat sessions containing query history and context for AI Agent Core';

-- ============================================================================
-- TABLE: app.queries
-- Purpose: Natural language queries submitted by users
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES app.query_sessions(id) ON DELETE CASCADE,
    raw_text TEXT NOT NULL,
    detected_locale VARCHAR(2),
    detected_intent VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_queries_raw_text_not_empty CHECK (length(trim(raw_text)) > 0),
    CONSTRAINT ck_queries_raw_text_length CHECK (length(raw_text) <= 2000),
    CONSTRAINT ck_queries_locale CHECK (detected_locale IS NULL OR detected_locale IN ('en', 'ar')),
    CONSTRAINT ck_queries_intent CHECK (detected_intent IS NULL OR detected_intent IN ('trend', 'comparison', 'breakdown', 'kpi', 'table'))
);

-- Indexes for queries
CREATE INDEX IF NOT EXISTS idx_queries_session_id ON app.queries(session_id);
CREATE INDEX IF NOT EXISTS idx_queries_created_at ON app.queries(created_at);
CREATE INDEX IF NOT EXISTS idx_queries_detected_locale ON app.queries(detected_locale);
CREATE INDEX IF NOT EXISTS idx_queries_detected_intent ON app.queries(detected_intent);

COMMENT ON TABLE app.queries IS 'Natural language queries from users in English or Arabic';

-- ============================================================================
-- TABLE: app.sql_statements
-- Purpose: Generated SQL queries with validation and retry tracking
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.sql_statements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_id UUID NOT NULL REFERENCES app.queries(id) ON DELETE CASCADE,
    sql_text TEXT NOT NULL,
    is_parameterized BOOLEAN NOT NULL DEFAULT TRUE,
    parameters JSONB DEFAULT '{}',
    validation_status VARCHAR(20) NOT NULL,
    blocked_reason TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_sql_statements_status CHECK (validation_status IN ('valid', 'blocked')),
    CONSTRAINT ck_sql_statements_retry_count CHECK (retry_count >= 0 AND retry_count <= 3)
);

-- Indexes for sql_statements
CREATE INDEX IF NOT EXISTS idx_sql_statements_query_id ON app.sql_statements(query_id);
CREATE INDEX IF NOT EXISTS idx_sql_statements_validation_status ON app.sql_statements(validation_status);
CREATE INDEX IF NOT EXISTS idx_sql_statements_created_at ON app.sql_statements(created_at);

COMMENT ON TABLE app.sql_statements IS 'Generated read-only SQL queries with OPA validation status and retry tracking';

-- ============================================================================
-- TABLE: app.query_results
-- Purpose: Query execution results with performance metrics
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.query_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    statement_id UUID NOT NULL REFERENCES app.sql_statements(id) ON DELETE CASCADE,
    row_count INTEGER DEFAULT 0,
    columns JSONB DEFAULT '[]',
    data JSONB,
    execution_time_ms INTEGER,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_query_results_row_count CHECK (row_count >= 0),
    CONSTRAINT ck_query_results_execution_time CHECK (execution_time_ms IS NULL OR execution_time_ms >= 0)
);

-- Indexes for query_results
CREATE INDEX IF NOT EXISTS idx_query_results_statement_id ON app.query_results(statement_id);
CREATE INDEX IF NOT EXISTS idx_query_results_created_at ON app.query_results(created_at);

COMMENT ON TABLE app.query_results IS 'Query execution results with row count, column metadata, and performance metrics';

-- ============================================================================
-- TRIGGER: Update updated_at for query_sessions
-- ============================================================================

-- Create a trigger to update the updated_at column on query_sessions
CREATE OR REPLACE FUNCTION app.update_query_session_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_query_sessions_updated_at ON app.query_sessions;
CREATE TRIGGER trg_query_sessions_updated_at
    BEFORE UPDATE ON app.query_sessions
    FOR EACH ROW
    EXECUTE FUNCTION app.update_query_session_updated_at();

COMMENT ON FUNCTION app.update_query_session_updated_at() IS 'Trigger function to update updated_at timestamp on query_sessions';

-- ============================================================================
-- GRANT PERMISSIONS
-- Grant SELECT to medisync_readonly role for AI agents
-- ============================================================================

GRANT SELECT ON app.query_sessions TO medisync_readonly;
GRANT SELECT ON app.queries TO medisync_readonly;
GRANT SELECT ON app.sql_statements TO medisync_readonly;
GRANT SELECT ON app.query_results TO medisync_readonly;

-- Grant full CRUD to medisync_app role for application operations
GRANT SELECT, INSERT, UPDATE, DELETE ON app.query_sessions TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.queries TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.sql_statements TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.query_results TO medisync_app;

-- Grant sequence usage for any serial columns (if added later)
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
