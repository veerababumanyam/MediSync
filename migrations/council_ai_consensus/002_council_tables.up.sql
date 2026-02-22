-- Council of AIs Consensus System: Core Tables
-- Migration: 002_council_tables.up.sql
-- Purpose: Create core deliberation and agent tables

-- Deliberation Status Enum
CREATE TYPE deliberation_status AS ENUM (
    'pending',
    'deliberating',
    'consensus',
    'uncertain',
    'failed'
);

-- Agent Health Status Enum
CREATE TYPE agent_health_status AS ENUM (
    'healthy',
    'degraded',
    'failed'
);

-- Council Deliberations Table
CREATE TABLE council_deliberations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_text TEXT NOT NULL,
    query_hash VARCHAR(64) NOT NULL,
    user_id UUID NOT NULL,
    status deliberation_status NOT NULL DEFAULT 'pending',
    consensus_threshold DECIMAL(3,2) NOT NULL DEFAULT 0.80,
    final_response TEXT,
    confidence_score DECIMAL(5,2),
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    CONSTRAINT chk_consensus_threshold CHECK (consensus_threshold >= 0.50 AND consensus_threshold <= 1.00),
    CONSTRAINT chk_confidence_score CHECK (confidence_score IS NULL OR (confidence_score >= 0 AND confidence_score <= 100))
);

-- Agent Instances Table
CREATE TABLE agent_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    health_status agent_health_status NOT NULL DEFAULT 'healthy',
    last_heartbeat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    config JSONB NOT NULL DEFAULT '{}',
    timeout_seconds INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_timeout_positive CHECK (timeout_seconds > 0)
);

-- Agent Responses Table
CREATE TABLE agent_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deliberation_id UUID NOT NULL REFERENCES council_deliberations(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agent_instances(id) ON DELETE RESTRICT,
    response_text TEXT NOT NULL,
    evidence_ids UUID[] NOT NULL DEFAULT '{}',
    confidence DECIMAL(5,2) NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_agent_confidence CHECK (confidence >= 0 AND confidence <= 100)
);

-- Indexes for Council Tables
CREATE INDEX idx_deliberation_user_id ON council_deliberations(user_id);
CREATE INDEX idx_deliberation_status ON council_deliberations(status);
CREATE INDEX idx_deliberation_created_at ON council_deliberations(created_at);
CREATE INDEX idx_deliberation_query_hash ON council_deliberations(query_hash);

CREATE INDEX idx_agent_health_status ON agent_instances(health_status);
CREATE INDEX idx_agent_last_heartbeat ON agent_instances(last_heartbeat);

CREATE INDEX idx_agent_response_deliberation ON agent_responses(deliberation_id);
CREATE INDEX idx_agent_response_agent ON agent_responses(agent_id);
CREATE INDEX idx_agent_response_embedding ON agent_responses
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 50);

-- Comments for documentation
COMMENT ON TABLE council_deliberations IS 'Represents a single multi-agent deliberation session';
COMMENT ON TABLE agent_instances IS 'Independent AI reasoning units participating in the Council';
COMMENT ON TABLE agent_responses IS 'Individual responses from agents in a deliberation';
COMMENT ON COLUMN council_deliberations.query_hash IS 'SHA-256 hash for query deduplication';
COMMENT ON COLUMN agent_responses.embedding IS 'Semantic embedding for equivalence detection';
