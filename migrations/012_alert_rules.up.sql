-- MediSync Dashboard Advanced Features - Alert Rules Migration
-- Version: 012
-- Description: Create alert_rules and notifications tables for KPI alerts
-- Task: T003
--
-- This migration establishes:
-- 1. alert_rules - User-defined conditions that trigger notifications
-- 2. notifications - Records of alert delivery attempts

-- ============================================================================
-- TABLE: app.alert_rules
-- Purpose: User-defined conditions that trigger notifications when metrics cross thresholds
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    metric_id VARCHAR(100) NOT NULL,
    metric_name VARCHAR(200) NOT NULL,
    operator VARCHAR(5) NOT NULL,
    threshold DECIMAL(20, 6) NOT NULL,
    check_interval INTEGER NOT NULL DEFAULT 300,
    channels JSONB NOT NULL DEFAULT '["in_app"]',
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    cooldown_period INTEGER DEFAULT 3600,
    last_triggered_at TIMESTAMPTZ NULL,
    last_value DECIMAL(20, 6) NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_alert_rules_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT ck_alert_rules_operator CHECK (operator IN ('gt', 'gte', 'lt', 'lte', 'eq')),
    CONSTRAINT ck_alert_rules_locale CHECK (locale IN ('en', 'ar')),
    CONSTRAINT ck_alert_rules_check_interval CHECK (check_interval >= 60 AND check_interval <= 86400),
    CONSTRAINT ck_alert_rules_cooldown CHECK (cooldown_period >= 0 AND cooldown_period <= 86400)
);

-- Indexes for alert_rules
CREATE INDEX IF NOT EXISTS idx_alert_rules_user_id ON app.alert_rules(user_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_user_active ON app.alert_rules(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_alert_rules_next_check ON app.alert_rules(is_active, last_triggered_at);
CREATE INDEX IF NOT EXISTS idx_alert_rules_metric_id ON app.alert_rules(metric_id);

COMMENT ON TABLE app.alert_rules IS 'User-defined conditions that trigger notifications when metrics cross thresholds';

-- ============================================================================
-- TABLE: app.notifications
-- Purpose: Records of alert delivery attempts
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_rule_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    content JSONB NOT NULL,
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    metric_value DECIMAL(20, 6) NOT NULL,
    threshold DECIMAL(20, 6) NOT NULL,
    error_message TEXT NULL,
    sent_at TIMESTAMPTZ NULL,
    delivered_at TIMESTAMPTZ NULL,
    read_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT ck_notifications_type CHECK (type IN ('in_app', 'email')),
    CONSTRAINT ck_notifications_status CHECK (status IN ('pending', 'sent', 'delivered', 'failed')),
    CONSTRAINT ck_notifications_locale CHECK (locale IN ('en', 'ar'))
);

-- Indexes for notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON app.notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_alert_rule_id ON app.notifications(alert_rule_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON app.notifications(user_id, read_at) WHERE read_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON app.notifications(created_at DESC);

COMMENT ON TABLE app.notifications IS 'Records of alert delivery attempts for user notifications';

-- ============================================================================
-- FOREIGN KEY: notifications -> alert_rules
-- ============================================================================

ALTER TABLE app.notifications
    ADD CONSTRAINT fk_notifications_alert_rule_id
    FOREIGN KEY (alert_rule_id) REFERENCES app.alert_rules(id) ON DELETE CASCADE;

-- ============================================================================
-- TRIGGER: Update updated_at for alert_rules
-- ============================================================================

CREATE OR REPLACE FUNCTION app.update_alert_rules_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_alert_rules_updated_at ON app.alert_rules;
CREATE TRIGGER trg_alert_rules_updated_at
    BEFORE UPDATE ON app.alert_rules
    FOR EACH ROW
    EXECUTE FUNCTION app.update_alert_rules_updated_at();

COMMENT ON FUNCTION app.update_alert_rules_updated_at() IS 'Trigger function to update updated_at timestamp on alert_rules';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.alert_rules TO medisync_readonly;
GRANT SELECT ON app.notifications TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.alert_rules TO medisync_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.notifications TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
