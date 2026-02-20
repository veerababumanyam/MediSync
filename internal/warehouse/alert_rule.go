// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// AlertRuleRepository handles database operations for alert rules.
type AlertRuleRepository struct {
	db *Repo
}

// NewAlertRuleRepository creates a new alert rule repository.
func NewAlertRuleRepository(db *Repo) *AlertRuleRepository {
	return &AlertRuleRepository{db: db}
}

// Create inserts a new alert rule.
func (r *AlertRuleRepository) Create(ctx context.Context, rule *models.AlertRule) error {
	query := `
		INSERT INTO app.alert_rules (user_id, name, description, metric_id, metric_name, operator, threshold, check_interval, channels, locale, cooldown_period)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		rule.UserID, rule.Name, rule.Description, rule.MetricID,
		rule.MetricName, rule.Operator, rule.Threshold,
		rule.CheckInterval, rule.Channels, rule.Locale, rule.CooldownPeriod,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create alert rule: %w", err)
	}

	return nil
}

// GetByUserID retrieves alert rules for a user.
func (r *AlertRuleRepository) GetByUserID(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]*models.AlertRule, error) {
	query := `
		SELECT id, user_id, name, description, metric_id, metric_name, operator, threshold, check_interval, channels, locale, cooldown_period, last_triggered_at, last_value, is_active, created_at, updated_at
		FROM app.alert_rules
		WHERE user_id = $1
	`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*models.AlertRule
	for rows.Next() {
		rule := &models.AlertRule{}
		var description pgtype.Text
		var lastTriggeredAt pgtype.Timestamptz
		var lastValue pgtype.Float8

		err := rows.Scan(
			&rule.ID, &rule.UserID, &rule.Name, &description, &rule.MetricID,
			&rule.MetricName, &rule.Operator, &rule.Threshold,
			&rule.CheckInterval, &rule.Channels, &rule.Locale, &rule.CooldownPeriod,
			&lastTriggeredAt, &lastValue, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan alert rule: %w", err)
		}

		if description.Valid {
			rule.Description = &description.String
		}
		if lastTriggeredAt.Valid {
			rule.LastTriggeredAt = &lastTriggeredAt.Time
		}
		if lastValue.Valid {
			rule.LastValue = &lastValue.Float64
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// Update updates an alert rule.
func (r *AlertRuleRepository) Update(ctx context.Context, rule *models.AlertRule) error {
	query := `
		UPDATE app.alert_rules SET
			name = $2, description = $3, threshold = $4, check_interval = $5, channels = $6, cooldown_period = $7, updated_at = NOW()
		WHERE id = $1 AND user_id = $8
		RETURNING updated_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Threshold,
		rule.CheckInterval, rule.Channels, rule.CooldownPeriod, rule.UserID,
	).Scan(&rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to update alert rule: %w", err)
	}

	return nil
}

// Delete removes an alert rule.
func (r *AlertRuleRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM app.alert_rules WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete alert rule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: alert rule not found")
	}

	return nil
}

// Toggle enables or disables an alert rule.
func (r *AlertRuleRepository) Toggle(ctx context.Context, id, userID uuid.UUID, isActive bool) error {
	query := `UPDATE app.alert_rules SET is_active = $3, updated_at = NOW() WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID, isActive)
	if err != nil {
		return fmt.Errorf("warehouse: failed to toggle alert rule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: alert rule not found")
	}

	return nil
}

// UpdateLastTriggered updates the last_triggered_at timestamp and last_value.
func (r *AlertRuleRepository) UpdateLastTriggered(ctx context.Context, id uuid.UUID, value float64) error {
	query := `UPDATE app.alert_rules SET last_triggered_at = NOW(), last_value = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.pool.Exec(ctx, query, id, value)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update last triggered: %w", err)
	}

	return nil
}

// GetDueForCheck retrieves alert rules that are due for evaluation.
func (r *AlertRuleRepository) GetDueForCheck(ctx context.Context) ([]*models.AlertRule, error) {
	query := `
		SELECT id, user_id, name, description, metric_id, metric_name, operator, threshold, check_interval, channels, locale, cooldown_period, last_triggered_at, last_value, is_active, created_at, updated_at
		FROM app.alert_rules
		WHERE is_active = true
		  AND (last_triggered_at IS NULL
		       OR last_triggered_at + (check_interval || ' seconds')::interval <= NOW())
	`

	rows, err := r.db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get due alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*models.AlertRule
	for rows.Next() {
		rule := &models.AlertRule{}
		var description pgtype.Text
		var lastTriggeredAt pgtype.Timestamptz
		var lastValue pgtype.Float8

		err := rows.Scan(
			&rule.ID, &rule.UserID, &rule.Name, &description, &rule.MetricID,
			&rule.MetricName, &rule.Operator, &rule.Threshold,
			&rule.CheckInterval, &rule.Channels, &rule.Locale, &rule.CooldownPeriod,
			&lastTriggeredAt, &lastValue, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan alert rule: %w", err)
		}

		if description.Valid {
			rule.Description = &description.String
		}
		if lastTriggeredAt.Valid {
			rule.LastTriggeredAt = &lastTriggeredAt.Time
		}
		if lastValue.Valid {
			rule.LastValue = &lastValue.Float64
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// GetByID retrieves a single alert rule by ID.
func (r *AlertRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AlertRule, error) {
	query := `
		SELECT id, user_id, name, description, metric_id, metric_name, operator, threshold, check_interval, channels, locale, cooldown_period, last_triggered_at, last_value, is_active, created_at, updated_at
		FROM app.alert_rules
		WHERE id = $1
	`

	rule := &models.AlertRule{}
	var description pgtype.Text
	var lastTriggeredAt pgtype.Timestamptz
	var lastValue pgtype.Float8

	err := r.db.pool.QueryRow(ctx, query, id).Scan(
		&rule.ID, &rule.UserID, &rule.Name, &description, &rule.MetricID,
		&rule.MetricName, &rule.Operator, &rule.Threshold,
		&rule.CheckInterval, &rule.Channels, &rule.Locale, &rule.CooldownPeriod,
		&lastTriggeredAt, &lastValue, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get alert rule: %w", err)
	}

	if description.Valid {
		rule.Description = &description.String
	}
	if lastTriggeredAt.Valid {
		rule.LastTriggeredAt = &lastTriggeredAt.Time
	}
	if lastValue.Valid {
		rule.LastValue = &lastValue.Float64
	}

	return rule, nil
}
