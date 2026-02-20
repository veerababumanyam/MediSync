// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// NotificationRepository handles database operations for notifications.
type NotificationRepository struct {
	db *Repo
}

// NewNotificationRepository creates a new notification repository.
func NewNotificationRepository(db *Repo) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create inserts a new notification.
func (r *NotificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	query := `
		INSERT INTO app.notifications (alert_rule_id, user_id, type, status, content, locale, metric_value, threshold)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		notif.AlertRuleID, notif.UserID, notif.Type, notif.Status,
		notif.Content, notif.Locale, notif.MetricValue, notif.Threshold,
	).Scan(&notif.ID, &notif.CreatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create notification: %w", err)
	}

	return nil
}

// GetByUserID retrieves notifications for a user.
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit int) ([]*models.Notification, error) {
	query := `
		SELECT id, alert_rule_id, user_id, type, status, content, locale, metric_value, threshold, error_message, sent_at, delivered_at, read_at, created_at
		FROM app.notifications
		WHERE user_id = $1
	`

	if unreadOnly {
		query += " AND read_at IS NULL"
	}

	query += " ORDER BY created_at DESC LIMIT $2"

	rows, err := r.db.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notif := &models.Notification{}
		var errorMessage pgtype.Text
		var sentAt, deliveredAt, readAt pgtype.Timestamptz

		err := rows.Scan(
			&notif.ID, &notif.AlertRuleID, &notif.UserID, &notif.Type, &notif.Status,
			&notif.Content, &notif.Locale, &notif.MetricValue, &notif.Threshold,
			&errorMessage, &sentAt, &deliveredAt, &readAt, &notif.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan notification: %w", err)
		}

		if errorMessage.Valid {
			notif.ErrorMessage = &errorMessage.String
		}
		if sentAt.Valid {
			notif.SentAt = &sentAt.Time
		}
		if deliveredAt.Valid {
			notif.DeliveredAt = &deliveredAt.Time
		}
		if readAt.Valid {
			notif.ReadAt = &readAt.Time
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read.
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID uuid.UUID) error {
	query := `UPDATE app.notifications SET read_at = NOW() WHERE id = $1 AND user_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("warehouse: failed to mark notification as read: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("warehouse: notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications for a user as read.
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `UPDATE app.notifications SET read_at = NOW() WHERE user_id = $1 AND read_at IS NULL`

	result, err := r.db.pool.Exec(ctx, query, userID)
	if err != nil {
		return 0, fmt.Errorf("warehouse: failed to mark all notifications as read: %w", err)
	}

	return result.RowsAffected(), nil
}

// UpdateStatus updates the status of a notification.
func (r *NotificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMessage *string) error {
	query := `
		UPDATE app.notifications SET
			status = $2,
			error_message = $3,
			sent_at = CASE WHEN $2 = 'sent' THEN NOW() ELSE sent_at END,
			delivered_at = CASE WHEN $2 = 'delivered' THEN NOW() ELSE delivered_at END
		WHERE id = $1
	`

	_, err := r.db.pool.Exec(ctx, query, id, status, errorMessage)
	if err != nil {
		return fmt.Errorf("warehouse: failed to update notification status: %w", err)
	}

	return nil
}

// GetUnreadCount returns the count of unread notifications for a user.
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM app.notifications WHERE user_id = $1 AND read_at IS NULL`

	var count int
	err := r.db.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("warehouse: failed to get unread count: %w", err)
	}

	return count, nil
}
