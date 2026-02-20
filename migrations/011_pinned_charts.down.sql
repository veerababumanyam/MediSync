-- MediSync Dashboard Advanced Features - Pinned Charts Migration (Rollback)
-- Version: 011
-- Description: Rollback pinned_charts table

DROP TRIGGER IF EXISTS trg_pinned_charts_updated_at ON app.pinned_charts;
DROP FUNCTION IF EXISTS app.update_pinned_charts_updated_at();
DROP TABLE IF EXISTS app.pinned_charts;
