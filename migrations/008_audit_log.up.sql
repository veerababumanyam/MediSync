-- MediSync AI Agent Core - AI Audit Log Migration
-- Version: 008
-- Description: Create AI-specific audit log table for AI Agent Core
-- Task: T011
--
-- This migration establishes:
-- 1. app.ai_audit_log - Audit trail for AI query submission, execution, and response
--
-- Purpose:
-- - Record of all AI agent actions for compliance
-- - Track query lifecycle from submission to response
-- - Support debugging and analysis of AI behavior
--
-- Action Types:
-- - query.submit: User submitted a query
-- - sql.generate: SQL was generated
-- - sql.validate: OPA validation performed
-- - sql.execute: Query executed
-- - result.return: Result returned to user
-- - review.queue: Added to review queue
-- - confidence.score: Confidence score calculated
-- - routing.decision: Routing decision made

-- ============================================================================
-- TABLE: app.ai_audit_log
-- Purpose: Record of AI query submission, execution, and response for compliance
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.ai_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_ai_audit_log_action CHECK (action IN (
        'query.submit',
        'sql.generate',
        'sql.validate',
        'sql.execute',
        'result.return',
        'review.queue',
        'confidence.score',
        'routing.decision',
        'error.occurred',
        'session.start',
        'session.end'
    ))
);

-- Indexes for ai_audit_log
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_user_id ON app.ai_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_tenant_id ON app.ai_audit_log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_created_at ON app.ai_audit_log(created_at);
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_action ON app.ai_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_resource ON app.ai_audit_log(resource_type, resource_id);

-- Partial indexes for common queries
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_tenant_created ON app.ai_audit_log(tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_ai_audit_log_user_created ON app.ai_audit_log(user_id, created_at);

COMMENT ON TABLE app.ai_audit_log IS 'Audit trail for AI Agent Core queries - records submission, execution, and response for compliance';

COMMENT ON COLUMN app.ai_audit_log.action IS 'Action type: query.submit, sql.generate, sql.validate, sql.execute, result.return, review.queue, confidence.score, routing.decision, error.occurred, session.start, session.end';
COMMENT ON COLUMN app.ai_audit_log.resource_type IS 'Resource type: query, session, statement, result';
COMMENT ON COLUMN app.ai_audit_log.resource_id IS 'UUID reference to the resource';
COMMENT ON COLUMN app.ai_audit_log.details IS 'JSON object with action-specific details';

-- ============================================================================
-- FUNCTION: Log AI action
-- Purpose: Convenience function to log an AI agent action
-- ============================================================================

CREATE OR REPLACE FUNCTION app.log_ai_action(
    p_user_id UUID,
    p_tenant_id UUID,
    p_action VARCHAR(50),
    p_resource_type VARCHAR(50) DEFAULT NULL,
    p_resource_id UUID DEFAULT NULL,
    p_details JSONB DEFAULT '{}',
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO app.ai_audit_log (
        user_id, tenant_id, action, resource_type, resource_id,
        details, ip_address, user_agent, created_at
    ) VALUES (
        p_user_id, p_tenant_id, p_action, p_resource_type, p_resource_id,
        p_details, p_ip_address, p_user_agent, NOW()
    )
    RETURNING id INTO v_id;

    RETURN v_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.log_ai_action(UUID, UUID, VARCHAR, VARCHAR, UUID, JSONB, INET, TEXT) IS 'Log an AI agent action to the audit trail. Returns the audit log entry ID.';

-- ============================================================================
-- FUNCTION: Get audit trail for query
-- Purpose: Get the complete audit trail for a specific query
-- ============================================================================

CREATE OR REPLACE FUNCTION app.get_query_audit_trail(p_query_id UUID)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    tenant_id UUID,
    action VARCHAR(50),
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        aal.id,
        aal.user_id,
        aal.tenant_id,
        aal.action,
        aal.resource_type,
        aal.resource_id,
        aal.details,
        aal.ip_address,
        aal.user_agent,
        aal.created_at
    FROM app.ai_audit_log aal
    WHERE
        aal.resource_id = p_query_id
        OR aal.details->>'query_id' = p_query_id::TEXT
    ORDER BY aal.created_at;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION app.get_query_audit_trail(UUID) IS 'Get the complete audit trail for a specific query by ID.';

-- ============================================================================
-- FUNCTION: Get audit summary for tenant
-- Purpose: Get audit summary statistics for a tenant within a date range
-- ============================================================================

CREATE OR REPLACE FUNCTION app.get_tenant_audit_summary(
    p_tenant_id UUID,
    p_start_date TIMESTAMPTZ DEFAULT NULL,
    p_end_date TIMESTAMPTZ DEFAULT NULL
)
RETURNS TABLE (
    action VARCHAR(50),
    action_count BIGINT,
    first_occurrence TIMESTAMPTZ,
    last_occurrence TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        aal.action,
        COUNT(*) AS action_count,
        MIN(aal.created_at) AS first_occurrence,
        MAX(aal.created_at) AS last_occurrence
    FROM app.ai_audit_log aal
    WHERE
        aal.tenant_id = p_tenant_id
        AND (p_start_date IS NULL OR aal.created_at >= p_start_date)
        AND (p_end_date IS NULL OR aal.created_at <= p_end_date)
    GROUP BY aal.action
    ORDER BY action_count DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ) IS 'Get audit summary statistics for a tenant within an optional date range.';

-- ============================================================================
-- ROW-LEVEL SECURITY: Enable RLS for ai_audit_log
-- ============================================================================

-- Enable RLS
ALTER TABLE app.ai_audit_log ENABLE ROW LEVEL SECURITY;

-- Force RLS for table owner
ALTER TABLE app.ai_audit_log FORCE ROW LEVEL SECURITY;

-- Drop existing policies if any (for idempotency)
DROP POLICY IF EXISTS ai_audit_log_select_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_insert_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_update_policy ON app.ai_audit_log;
DROP POLICY IF EXISTS ai_audit_log_delete_policy ON app.ai_audit_log;

-- SELECT Policy: All roles can read audit log entries
CREATE POLICY ai_audit_log_select_policy ON app.ai_audit_log
    FOR SELECT
    TO medisync_readonly, medisync_app
    USING (true);

-- INSERT Policy: Only app role can insert new audit records
CREATE POLICY ai_audit_log_insert_policy ON app.ai_audit_log
    FOR INSERT
    TO medisync_app
    WITH CHECK (true);

-- UPDATE Policy: DENY ALL - No updates allowed
CREATE POLICY ai_audit_log_update_policy ON app.ai_audit_log
    FOR UPDATE
    TO medisync_readonly, medisync_app
    USING (false)
    WITH CHECK (false);

-- DELETE Policy: DENY ALL - No deletions allowed
CREATE POLICY ai_audit_log_delete_policy ON app.ai_audit_log
    FOR DELETE
    TO medisync_readonly, medisync_app
    USING (false);

COMMENT ON POLICY ai_audit_log_select_policy ON app.ai_audit_log IS 'Allows all application roles to read AI audit log entries';
COMMENT ON POLICY ai_audit_log_insert_policy ON app.ai_audit_log IS 'Allows app role to create new audit log entries';
COMMENT ON POLICY ai_audit_log_update_policy ON app.ai_audit_log IS 'DENIES all UPDATE operations - audit records are immutable';
COMMENT ON POLICY ai_audit_log_delete_policy ON app.ai_audit_log IS 'DENIES all DELETE operations - audit records cannot be removed';

-- ============================================================================
-- TRIGGER: Prevent UPDATE/DELETE operations (defense in depth)
-- ============================================================================

CREATE OR REPLACE FUNCTION app.prevent_ai_audit_log_modification()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        RAISE EXCEPTION 'UPDATE operations are not permitted on app.ai_audit_log. Audit records are immutable.'
            USING ERRCODE = 'prohibited_sql_statement_attempted',
                  HINT = 'AI audit log entries cannot be modified after creation for compliance reasons.';
    ELSIF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION 'DELETE operations are not permitted on app.ai_audit_log. Audit records are immutable.'
            USING ERRCODE = 'prohibited_sql_statement_attempted',
                  HINT = 'AI audit log entries cannot be deleted for compliance reasons.';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_prevent_ai_audit_log_modification ON app.ai_audit_log;
CREATE TRIGGER trg_prevent_ai_audit_log_modification
    BEFORE UPDATE OR DELETE ON app.ai_audit_log
    FOR EACH ROW
    EXECUTE FUNCTION app.prevent_ai_audit_log_modification();

COMMENT ON FUNCTION app.prevent_ai_audit_log_modification() IS 'Trigger function that prevents UPDATE and DELETE operations on ai_audit_log table';

-- ============================================================================
-- REVOKE/GRANT EXPLICIT PERMISSIONS
-- ============================================================================

-- Revoke any inherited permissions
REVOKE ALL ON app.ai_audit_log FROM PUBLIC;

-- Grant explicit permissions
GRANT SELECT ON app.ai_audit_log TO medisync_readonly;
GRANT SELECT, INSERT ON app.ai_audit_log TO medisync_app;

-- Grant execute on functions
GRANT EXECUTE ON FUNCTION app.log_ai_action(UUID, UUID, VARCHAR, VARCHAR, UUID, JSONB, INET, TEXT) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.get_query_audit_trail(UUID) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.get_query_audit_trail(UUID) TO medisync_readonly;
GRANT EXECUTE ON FUNCTION app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.get_tenant_audit_summary(UUID, TIMESTAMPTZ, TIMESTAMPTZ) TO medisync_readonly;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
