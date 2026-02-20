-- MediSync Dashboard Advanced Features - Pinned Charts Migration
-- Version: 011
-- Description: Create pinned_charts table for user dashboard widgets
-- Task: T002
--
-- This migration establishes:
-- Saved visualizations displayed on the user's personal dashboard

-- ============================================================================
-- TABLE: app.pinned_charts
-- Purpose: Stores pinned chart configurations for user dashboards
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.pinned_charts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    query_id UUID NULL,
    natural_language_query TEXT NOT NULL,
    sql_query TEXT NOT NULL,
    chart_spec JSONB NOT NULL,
    chart_type VARCHAR(20) NOT NULL,
    refresh_interval INTEGER DEFAULT 300,
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    position JSONB NOT NULL DEFAULT '{"row":0,"col":0,"size":1}',
    last_refreshed_at TIMESTAMPTZ NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_pinned_charts_title_not_empty CHECK (length(trim(title)) > 0),
    CONSTRAINT ck_pinned_charts_chart_type CHECK (chart_type IN ('bar', 'line', 'pie', 'table', 'kpi')),
    CONSTRAINT ck_pinned_charts_locale CHECK (locale IN ('en', 'ar')),
    CONSTRAINT ck_pinned_charts_refresh_interval CHECK (refresh_interval = 0 OR (refresh_interval >= 60 AND refresh_interval <= 3600)),
    CONSTRAINT ck_pinned_charts_position_size CHECK ((position->>'size')::int IN (1, 2, 3))
);

-- Indexes for pinned_charts
CREATE INDEX IF NOT EXISTS idx_pinned_charts_user_id ON app.pinned_charts(user_id);
CREATE INDEX IF NOT EXISTS idx_pinned_charts_user_active ON app.pinned_charts(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_pinned_charts_created_at ON app.pinned_charts(created_at);

COMMENT ON TABLE app.pinned_charts IS 'Saved visualizations displayed on user personal dashboards';

-- ============================================================================
-- TRIGGER: Update updated_at for pinned_charts
-- ============================================================================

CREATE OR REPLACE FUNCTION app.update_pinned_charts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_pinned_charts_updated_at ON app.pinned_charts;
CREATE TRIGGER trg_pinned_charts_updated_at
    BEFORE UPDATE ON app.pinned_charts
    FOR EACH ROW
    EXECUTE FUNCTION app.update_pinned_charts_updated_at();

COMMENT ON FUNCTION app.update_pinned_charts_updated_at() IS 'Trigger function to update updated_at timestamp on pinned_charts';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.pinned_charts TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.pinned_charts TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
