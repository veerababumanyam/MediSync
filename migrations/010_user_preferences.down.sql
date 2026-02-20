-- MediSync Dashboard Advanced Features - User Preferences Migration (Rollback)
-- Version: 010
-- Description: Rollback user_preferences table

DROP TRIGGER IF EXISTS trg_user_preferences_updated_at ON app.user_preferences;
DROP FUNCTION IF EXISTS app.update_user_preferences_updated_at();
DROP TABLE IF EXISTS app.user_preferences;
