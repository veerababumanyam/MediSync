-- MediSync Dashboard Advanced Features - User Preferences Migration
-- Version: 010
-- Description: Create user_preferences table for locale, calendar, and display settings
-- Task: T001
--
-- This migration establishes:
-- User-specific display and formatting preferences that persist across sessions and devices

-- ============================================================================
-- TABLE: app.user_preferences
-- Purpose: Stores user locale, numeral system, calendar, and timezone preferences
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    locale VARCHAR(2) NOT NULL DEFAULT 'en',
    numeral_system VARCHAR(20) NOT NULL DEFAULT 'western',
    calendar_system VARCHAR(20) NOT NULL DEFAULT 'gregorian',
    report_language VARCHAR(2) NOT NULL DEFAULT 'en',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Dubai',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT uq_user_preferences_user_id UNIQUE (user_id),
    CONSTRAINT ck_user_preferences_locale CHECK (locale IN ('en', 'ar')),
    CONSTRAINT ck_user_preferences_numeral CHECK (numeral_system IN ('western', 'eastern_arabic')),
    CONSTRAINT ck_user_preferences_calendar CHECK (calendar_system IN ('gregorian', 'hijri')),
    CONSTRAINT ck_user_preferences_report_lang CHECK (report_language IN ('en', 'ar'))
);

-- Indexes for user_preferences
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON app.user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_user_preferences_locale ON app.user_preferences(locale);

COMMENT ON TABLE app.user_preferences IS 'User-specific display and formatting preferences for dashboard and i18n';

-- ============================================================================
-- TRIGGER: Update updated_at for user_preferences
-- ============================================================================

CREATE OR REPLACE FUNCTION app.update_user_preferences_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_user_preferences_updated_at ON app.user_preferences;
CREATE TRIGGER trg_user_preferences_updated_at
    BEFORE UPDATE ON app.user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION app.update_user_preferences_updated_at();

COMMENT ON FUNCTION app.update_user_preferences_updated_at() IS 'Trigger function to update updated_at timestamp on user_preferences';

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT SELECT ON app.user_preferences TO medisync_readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON app.user_preferences TO medisync_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
