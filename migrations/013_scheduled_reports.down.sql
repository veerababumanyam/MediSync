-- MediSync Dashboard Advanced Features - Scheduled Reports Migration (Rollback)
-- Version: 013
-- Description: Rollback scheduled_reports and scheduled_report_runs tables

ALTER TABLE IF EXISTS app.scheduled_report_runs DROP CONSTRAINT IF EXISTS fk_scheduled_report_runs_report_id;
DROP TRIGGER IF EXISTS trg_scheduled_reports_updated_at ON app.scheduled_reports;
DROP FUNCTION IF EXISTS app.update_scheduled_reports_updated_at();
DROP TABLE IF EXISTS app.scheduled_report_runs;
DROP TABLE IF EXISTS app.scheduled_reports;
