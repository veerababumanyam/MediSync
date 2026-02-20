-- MediSync Dashboard Advanced Features - Chat Messages Migration
-- Version: 014
-- Description: Create chat_messages table for conversation history
-- Task: T005
--
-- This migration establishes:
-- Messages in chat conversations, either from user or AI

-- ============================================================================
-- TABLE: app.chat_messages
-- Purpose: Stores chat messages for conversation history
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    chart_spec JSONB NULL,
    table_data JSONB NULL,
    drilldown_query TEXT NULL,
    confidence_score DECIMAL(5, 4) NULL,
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_chat_messages_role CHECK (role IN ('user', 'assistant')),
    CONSTRAINT ck_chat_messages_locale CHECK (locale IN ('en', 'ar')),
    CONSTRAINT ck_chat_messages_content_not_empty CHECK (length(trim(content)) > 0),
    CONSTRAINT ck_chat_messages_confidence CHECK (confidence_score IS NULL OR (confidence_score >= 0.0 AND confidence_score <= 1.0))
);

-- Indexes for chat_messages
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON app.chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_user_id ON app.chat_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON app.chat_messages(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_created ON app.chat_messages(session_id, created_at);

COMMENT ON TABLE app.chat_messages IS 'Chat messages for conversation history with AI assistant';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.chat_messages TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.chat_messages TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
