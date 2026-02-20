// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// UserPreferenceRepository handles database operations for user preferences.
type UserPreferenceRepository struct {
	db *Repo
}

// NewUserPreferenceRepository creates a new user preference repository.
func NewUserPreferenceRepository(db *Repo) *UserPreferenceRepository {
	return &UserPreferenceRepository{db: db}
}

// GetByUserID retrieves user preferences by user ID.
func (r *UserPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserPreference, error) {
	query := `
		SELECT id, user_id, locale, numeral_system, calendar_system, report_language, timezone, created_at, updated_at
		FROM app.user_preferences
		WHERE user_id = $1
	`

	var pref models.UserPreference
	err := r.db.pool.QueryRow(ctx, query, userID).Scan(
		&pref.ID, &pref.UserID, &pref.Locale, &pref.NumeralSystem,
		&pref.CalendarSystem, &pref.ReportLanguage, &pref.Timezone,
		&pref.CreatedAt, &pref.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get user preference: %w", err)
	}

	return &pref, nil
}

// Upsert creates or updates user preferences.
func (r *UserPreferenceRepository) Upsert(ctx context.Context, pref *models.UserPreference) error {
	query := `
		INSERT INTO app.user_preferences (user_id, locale, numeral_system, calendar_system, report_language, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			locale = EXCLUDED.locale,
			numeral_system = EXCLUDED.numeral_system,
			calendar_system = EXCLUDED.calendar_system,
			report_language = EXCLUDED.report_language,
			timezone = EXCLUDED.timezone,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		pref.UserID, pref.Locale, pref.NumeralSystem,
		pref.CalendarSystem, pref.ReportLanguage, pref.Timezone,
	).Scan(&pref.ID, &pref.CreatedAt, &pref.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert user preference: %w", err)
	}

	return nil
}

// UpdateLocale updates only the locale for a user.
func (r *UserPreferenceRepository) UpdateLocale(ctx context.Context, userID uuid.UUID, locale string) error {
	query := `
		INSERT INTO app.user_preferences (user_id, locale)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET
			locale = EXCLUDED.locale,
			updated_at = NOW()
	`

	_, err := r.db.pool.Exec(ctx, query, userID, locale)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update user locale: %w", err)
	}

	return nil
}
