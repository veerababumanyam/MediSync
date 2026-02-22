-- Council of AIs Consensus System: Knowledge Graph Tables
-- Migration: 001_knowledge_graph.up.sql
-- Purpose: Create tables for storing the Medical Knowledge Graph

-- Enable pgvector extension if not already enabled
CREATE EXTENSION IF NOT EXISTS vector;

-- Knowledge Graph Node Types
CREATE TYPE kg_node_type AS ENUM (
    'concept',
    'medication',
    'procedure',
    'condition',
    'organization'
);

-- Knowledge Graph Edge Types
CREATE TYPE kg_edge_type AS ENUM (
    'TREATS',
    'CAUSES',
    'CONTRAINDICATES',
    'RELATED_TO',
    'SUBSUMES',
    'PART_OF'
);

-- Knowledge Graph Nodes Table
CREATE TABLE knowledge_graph_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_type kg_node_type NOT NULL,
    concept VARCHAR(255) NOT NULL,
    definition TEXT NOT NULL,
    embedding vector(1536) NOT NULL,
    source VARCHAR(255) NOT NULL,
    source_id VARCHAR(255),
    confidence DECIMAL(5,2) NOT NULL DEFAULT 100.00,
    last_verified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edges UUID[] DEFAULT '{}',
    edge_types kg_edge_type[] DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_confidence CHECK (confidence >= 0 AND confidence <= 100)
);

-- Indexes for Knowledge Graph
CREATE INDEX idx_kg_node_type ON knowledge_graph_nodes(node_type);
CREATE INDEX idx_kg_node_concept ON knowledge_graph_nodes(concept);
CREATE INDEX idx_kg_node_source ON knowledge_graph_nodes(source);
CREATE INDEX idx_kg_node_last_verified ON knowledge_graph_nodes(last_verified);

-- IVFFlat index for vector similarity search
CREATE INDEX idx_kg_node_embedding ON knowledge_graph_nodes
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_kg_node_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_kg_node_updated_at
    BEFORE UPDATE ON knowledge_graph_nodes
    FOR EACH ROW
    EXECUTE FUNCTION update_kg_node_updated_at();

-- Comments for documentation
COMMENT ON TABLE knowledge_graph_nodes IS 'Medical Knowledge Graph nodes containing verified healthcare knowledge';
COMMENT ON COLUMN knowledge_graph_nodes.embedding IS 'Semantic embedding for similarity search (1536 dimensions)';
COMMENT ON COLUMN knowledge_graph_nodes.edges IS 'Array of connected node UUIDs';
COMMENT ON COLUMN knowledge_graph_nodes.edge_types IS 'Array of edge types parallel to edges array';
COMMENT ON COLUMN knowledge_graph_nodes.confidence IS 'Knowledge reliability score (0-100)';
