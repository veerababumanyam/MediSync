-- MediSync Dashboard Advanced Features - Scheduled Reports Migration
-- Version: 013
-- Description: Create scheduled_reports and scheduled_report_runs tables
-- Task: T004
--
-- This migration establishes:
-- 1. scheduled_reports - Configuration for recurring report generation
-- 2. scheduled_report_runs - Audit records of report generation attempts

-- ============================================================================
-- TABLE: app.scheduled_reports
-- Purpose: Configuration for recurring report generation and delivery
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.scheduled_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    query_id UUID NULL,
    natural_language_query TEXT NOT NULL,
    sql_query TEXT NOT NULL,
    schedule_type VARCHAR(20) NOT NULL,
    schedule_time TIME NOT NULL,
    schedule_day INTEGER NULL,
    recipients JSONB NOT NULL,
    format VARCHAR(10) NOT NULL DEFAULT 'pdf',
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    include_charts BOOLEAN DEFAULT true,
    last_run_at TIMESTAMPTZ NULL,
    next_run_at TIMESTAMPTZ NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_scheduled_reports_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT ck_scheduled_reports_schedule_type CHECK (schedule_type IN ('daily', 'weekly', 'monthly', 'quarterly')),
    CONSTRAINT ck_scheduled_reports_format CHECK (format IN ('pdf', 'xlsx', 'csv')),
    CONSTRAINT ck_scheduled_reports_locale CHECK (locale IN ('en', 'ar')),
    CONSTRAINT ck_scheduled_reports_schedule_day CHECK (schedule_day IS NULL OR (schedule_day >= 1 AND schedule_day <= 31))
);

-- Indexes for scheduled_reports
CREATE INDEX IF NOT EXISTS idx_scheduled_reports_user_id ON app.scheduled_reports(user_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_reports_next_run ON app.scheduled_reports(is_active, next_run_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_reports_created_at ON app.scheduled_reports(created_at);

COMMENT ON TABLE app.scheduled_reports IS 'Configuration for recurring report generation and delivery';

-- ============================================================================
-- TABLE: app.scheduled_report_runs
-- Purpose: Audit records of report generation attempts
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.scheduled_report_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL,
    file_path VARCHAR(500) NULL,
    file_size_bytes BIGINT NULL,
    row_count INTEGER NULL,
    error_message TEXT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ NULL,

    -- Constraints
    CONSTRAINT ck_scheduled_report_runs_status CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    CONSTRAINT ck_scheduled_report_runs_row_count CHECK (row_count IS NULL OR row_count >= 0)
);

-- Indexes for scheduled_report_runs
CREATE INDEX IF NOT EXISTS idx_scheduled_report_runs_report_id ON app.scheduled_report_runs(report_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_report_runs_status ON app.scheduled_report_runs(status);
CREATE INDEX IF NOT EXISTS idx_scheduled_report_runs_started_at ON app.scheduled_report_runs(started_at DESC);

COMMENT ON TABLE app.scheduled_report_runs IS 'Audit records of report generation attempts';

-- ============================================================================
-- FOREIGN KEY: scheduled_report_runs -> scheduled_reports
-- ============================================================================

ALTER TABLE app.scheduled_report_runs
    ADD CONSTRAINT fk_scheduled_report_runs_report_id
    FOREIGN KEY (report_id) REFERENCES app.scheduled_reports(id) ON DELETE CASCADE;

-- ============================================================================
-- TRIGGER: Update updated_at for scheduled_reports
-- ============================================================================

CREATE OR REPLACE FUNCTION app.update_scheduled_reports_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_scheduled_reports_updated_at ON app.scheduled_reports;
CREATE TRIGGER trg_scheduled_reports_updated_at
    BEFORE UPDATE ON app.scheduled_reports
    FOR EACH ROW
    EXECUTE FUNCTION app.update_scheduled_reports_updated_at();

COMMENT ON FUNCTION app.update_scheduled_reports_updated_at() IS 'Trigger function to update updated_at timestamp on scheduled_reports';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.scheduled_reports TO medisync_readonly;
GRANT SELECT ON app.scheduled_report_runs TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.scheduled_reports TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.scheduled_report_runs TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
