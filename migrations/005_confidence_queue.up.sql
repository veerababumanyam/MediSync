-- MediSync AI Agent Core - Confidence Scoring and Review Queue Migration
-- Version: 005
-- Description: Create confidence scoring and review queue tables for AI Agent Core
-- Task: T008
--
-- This migration establishes:
-- 1. app.confidence_scores - Numerical assessment of query result accuracy
-- 2. app.review_queue - Low-confidence queries requiring human review
--
-- Scoring Factors (stored in factors JSONB):
--   - intent_clarity: 0.0-1.0
--   - schema_match_quality: 0.0-1.0
--   - sql_complexity_penalty: 0.0-0.3
--   - retry_penalty: 0.0-0.3
--   - hallucination_risk: 0.0-1.0
--
-- Routing Logic:
--   - score >= 70: routing_decision = "normal"
--   - score 50-69: routing_decision = "warning"
--   - score < 50: routing_decision = "clarify" (added to review_queue)

-- ============================================================================
-- TABLE: app.confidence_scores
-- Purpose: Numerical assessment of AI-generated query result accuracy
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.confidence_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_id UUID NOT NULL REFERENCES app.queries(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL,
    factors JSONB DEFAULT '{}',
    routing_decision VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_confidence_scores_score_range CHECK (score >= 0 AND score <= 100),
    CONSTRAINT ck_confidence_scores_routing CHECK (routing_decision IN ('normal', 'warning', 'review', 'clarify'))
);

-- Indexes for confidence_scores
CREATE INDEX IF NOT EXISTS idx_confidence_scores_query_id ON app.confidence_scores(query_id);
CREATE INDEX IF NOT EXISTS idx_confidence_scores_score ON app.confidence_scores(score);
CREATE INDEX IF NOT EXISTS idx_confidence_scores_routing ON app.confidence_scores(routing_decision);
CREATE INDEX IF NOT EXISTS idx_confidence_scores_created_at ON app.confidence_scores(created_at);

COMMENT ON TABLE app.confidence_scores IS 'Confidence scores for AI-generated query results with routing decisions';

-- Add a comment explaining the factors JSONB structure
COMMENT ON COLUMN app.confidence_scores.factors IS 'JSON object with scoring factors: intent_clarity, schema_match_quality, sql_complexity_penalty, retry_penalty, hallucination_risk';

-- ============================================================================
-- TABLE: app.review_queue
-- Purpose: Low-confidence queries requiring human review
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.review_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_id UUID NOT NULL REFERENCES app.queries(id) ON DELETE CASCADE,
    score_id UUID REFERENCES app.confidence_scores(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    resolution TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_review_queue_status CHECK (status IN ('pending', 'approved', 'rejected', 'clarified', 'escalated'))
);

-- Indexes for review_queue
CREATE INDEX IF NOT EXISTS idx_review_queue_query_id ON app.review_queue(query_id);
CREATE INDEX IF NOT EXISTS idx_review_queue_score_id ON app.review_queue(score_id);
CREATE INDEX IF NOT EXISTS idx_review_queue_status ON app.review_queue(status);
CREATE INDEX IF NOT EXISTS idx_review_queue_reviewed_by ON app.review_queue(reviewed_by);
CREATE INDEX IF NOT EXISTS idx_review_queue_created_at ON app.review_queue(created_at);

COMMENT ON TABLE app.review_queue IS 'Queue of low-confidence queries requiring human review before results are returned';

-- ============================================================================
-- FUNCTION: Get pending review count
-- Purpose: Helper function to count pending reviews for a tenant
-- ============================================================================

CREATE OR REPLACE FUNCTION app.get_pending_review_count(p_tenant_id UUID)
RETURNS INTEGER AS $$
DECLARE
    v_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM app.review_queue rq
    JOIN app.queries q ON rq.query_id = q.id
    JOIN app.query_sessions qs ON q.session_id = qs.id
    WHERE qs.tenant_id = p_tenant_id
    AND rq.status = 'pending';

    RETURN v_count;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION app.get_pending_review_count(UUID) IS 'Returns the count of pending reviews for a given tenant';

-- ============================================================================
-- FUNCTION: Auto-route based on confidence score
-- Purpose: Trigger function to automatically add low-score queries to review queue
-- ============================================================================

CREATE OR REPLACE FUNCTION app.auto_route_low_confidence()
RETURNS TRIGGER AS $$
BEGIN
    -- If score < 50, automatically add to review queue
    IF NEW.score < 50 THEN
        INSERT INTO app.review_queue (query_id, score_id, status, created_at)
        VALUES (NEW.query_id, NEW.id, 'pending', NOW());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for auto-routing
DROP TRIGGER IF EXISTS trg_auto_route_low_confidence ON app.confidence_scores;
CREATE TRIGGER trg_auto_route_low_confidence
    AFTER INSERT ON app.confidence_scores
    FOR EACH ROW
    EXECUTE FUNCTION app.auto_route_low_confidence();

COMMENT ON FUNCTION app.auto_route_low_confidence() IS 'Automatically adds queries with confidence score < 50 to the review queue';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

-- Grant SELECT to medisync_readonly role for AI agents
GRANT SELECT ON app.confidence_scores TO medisync_readonly;
GRANT SELECT ON app.review_queue TO medisync_readonly;

-- Grant full CRUD to medisync_app role for application operations
GRANT SELECT, INSERT, UPDATE, DELETE ON app.confidence_scores TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.review_queue TO medisync_app;

-- Grant execute on functions
GRANT EXECUTE ON FUNCTION app.get_pending_review_count(UUID) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.get_pending_review_count(UUID) TO medisync_readonly;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
