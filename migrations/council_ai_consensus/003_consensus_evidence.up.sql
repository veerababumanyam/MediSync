-- Council of AIs Consensus System: Consensus and Evidence Tables
-- Migration: 003_consensus_evidence.up.sql
-- Purpose: Create tables for consensus records and evidence trails

-- Consensus Records Table
CREATE TABLE consensus_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deliberation_id UUID NOT NULL UNIQUE REFERENCES council_deliberations(id) ON DELETE CASCADE,
    agreement_score DECIMAL(5,2) NOT NULL,
    equivalence_groups JSONB NOT NULL DEFAULT '[]',
    threshold_met BOOLEAN NOT NULL,
    dissenting_agents UUID[] NOT NULL DEFAULT '{}',
    consensus_method VARCHAR(50) NOT NULL DEFAULT 'weighted_vote',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_agreement_score CHECK (agreement_score >= 0 AND agreement_score <= 100)
);

-- Evidence Trails Table
CREATE TABLE evidence_trails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deliberation_id UUID NOT NULL UNIQUE REFERENCES council_deliberations(id) ON DELETE CASCADE,
    node_ids UUID[] NOT NULL,
    traversal_path JSONB NOT NULL,
    relevance_scores JSONB NOT NULL DEFAULT '{}',
    hop_count INTEGER NOT NULL DEFAULT 0,
    cached_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT chk_hop_count_positive CHECK (hop_count >= 0),
    CONSTRAINT chk_expires_future CHECK (expires_at > cached_at)
);

-- Indexes for Consensus and Evidence Tables
CREATE INDEX idx_consensus_deliberation ON consensus_records(deliberation_id);
CREATE INDEX idx_consensus_threshold_met ON consensus_records(threshold_met);
CREATE INDEX idx_consensus_created_at ON consensus_records(created_at);

CREATE INDEX idx_evidence_deliberation ON evidence_trails(deliberation_id);
CREATE INDEX idx_evidence_expires_at ON evidence_trails(expires_at);
CREATE INDEX idx_evidence_cached_at ON evidence_trails(cached_at);

-- GIN index for node_ids array searches
CREATE INDEX idx_evidence_node_ids ON evidence_trails USING GIN(node_ids);

-- Comments for documentation
COMMENT ON TABLE consensus_records IS 'Captures the consensus calculation results for a deliberation';
COMMENT ON TABLE evidence_trails IS 'Records the Knowledge Graph traversal path for a deliberation';
COMMENT ON COLUMN consensus_records.equivalence_groups IS 'Groups of semantically equivalent agent responses';
COMMENT ON COLUMN evidence_trails.traversal_path IS 'Full path with edges and scores';
COMMENT ON COLUMN evidence_trails.relevance_scores IS 'Node ID → relevance score mapping';
