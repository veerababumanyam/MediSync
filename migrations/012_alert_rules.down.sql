-- MediSync Dashboard Advanced Features - Alert Rules Migration (Rollback)
-- Version: 012
-- Description: Rollback alert_rules and notifications tables

ALTER TABLE IF EXISTS app.notifications DROP CONSTRAINT IF EXISTS fk_notifications_alert_rule_id;
DROP TRIGGER IF EXISTS trg_alert_rules_updated_at ON app.alert_rules;
DROP FUNCTION IF EXISTS app.update_alert_rules_updated_at();
DROP TABLE IF EXISTS app.notifications;
DROP TABLE IF EXISTS app.alert_rules;
