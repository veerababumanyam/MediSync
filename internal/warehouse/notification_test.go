package warehouse

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/stretchr/testify/assert"
)

// TestNotificationModel_Fields tests model fields (unit test)
func TestNotificationModel_Fields(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	userID := uuid.New()
	alertID := uuid.New()
	sentAt := now.Add(-5 * time.Minute)
	deliveredAt := now.Add(-4 * time.Minute)
	readAt := now.Add(-2 * time.Minute)
	errMsg := "connection timeout"

	notification := models.Notification{
		ID:           id,
		AlertRuleID:  alertID,
		UserID:       userID,
		Type:         "email",
		Status:       "delivered",
		Content:      models.NotificationContent{Title: "Alert Triggered", Message: "Revenue dropped below threshold"},
		Locale:       "en",
		MetricValue:  9500.0,
		Threshold:    10000.0,
		ErrorMessage: &errMsg,
		SentAt:       &sentAt,
		DeliveredAt:  &deliveredAt,
		ReadAt:       &readAt,
		CreatedAt:    now,
	}

	assert.Equal(t, id, notification.ID)
	assert.Equal(t, userID, notification.UserID)
	assert.Equal(t, alertID, notification.AlertRuleID)
	assert.Equal(t, "email", notification.Type)
	assert.Equal(t, "delivered", notification.Status)
	assert.Equal(t, "Alert Triggered", notification.Content.Title)
	assert.Equal(t, 9500.0, notification.MetricValue)
	assert.Equal(t, 10000.0, notification.Threshold)
	assert.NotNil(t, notification.SentAt)
	assert.NotNil(t, notification.DeliveredAt)
	assert.NotNil(t, notification.ReadAt)
}

// TestNotificationContent_Fields tests content fields (unit test)
func TestNotificationContent_Fields(t *testing.T) {
	content := models.NotificationContent{
		Title:     "High Revenue Alert",
		Message:   "Revenue exceeded $50,000 today",
		ActionURL: "/dashboard/alerts/123",
	}

	assert.Equal(t, "High Revenue Alert", content.Title)
	assert.Equal(t, "Revenue exceeded $50,000 today", content.Message)
	assert.Equal(t, "/dashboard/alerts/123", content.ActionURL)
}

// TestNotificationContent_WithoutActionURL tests content without action URL (unit test)
func TestNotificationContent_WithoutActionURL(t *testing.T) {
	content := models.NotificationContent{
		Title:   "Simple Alert",
		Message: "A threshold was crossed",
	}

	assert.Equal(t, "Simple Alert", content.Title)
	assert.Empty(t, content.ActionURL)
}

// TestNotification_Types tests notification types (unit test)
func TestNotification_Types(t *testing.T) {
	types := []string{"email", "slack", "sms", "push", "webhook"}

	for _, nt := range types {
		t.Run(nt, func(t *testing.T) {
			notification := models.Notification{Type: nt}
			assert.Equal(t, nt, notification.Type)
		})
	}
}

// TestNotification_Statuses tests notification statuses (unit test)
func TestNotification_Statuses(t *testing.T) {
	statuses := []string{"pending", "sent", "delivered", "read", "failed"}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			notification := models.Notification{Status: status}
			assert.Equal(t, status, notification.Status)
		})
	}
}

// TestNotification_WithNilError tests nil error message (unit test)
func TestNotification_WithNilError(t *testing.T) {
	notification := models.Notification{
		ID:           uuid.New(),
		ErrorMessage: nil,
	}

	assert.Nil(t, notification.ErrorMessage)
}

// TestNotification_WithError tests with error message (unit test)
func TestNotification_WithError(t *testing.T) {
	errMsg := "SMTP connection failed"
	notification := models.Notification{
		ID:           uuid.New(),
		Status:       "failed",
		ErrorMessage: &errMsg,
	}

	assert.Equal(t, "failed", notification.Status)
	assert.Equal(t, "SMTP connection failed", *notification.ErrorMessage)
}

// TestNotification_Timestamps tests timestamp fields (unit test)
func TestNotification_Timestamps(t *testing.T) {
	now := time.Now()
	sent := now.Add(-10 * time.Minute)
	delivered := now.Add(-9 * time.Minute)
	read := now.Add(-5 * time.Minute)

	notification := models.Notification{
		CreatedAt:   now,
		SentAt:      &sent,
		DeliveredAt: &delivered,
		ReadAt:      &read,
	}

	assert.Equal(t, now, notification.CreatedAt)
	assert.Equal(t, sent, *notification.SentAt)
	assert.Equal(t, delivered, *notification.DeliveredAt)
	assert.Equal(t, read, *notification.ReadAt)
}

// TestNotification_Locale tests locale field (unit test)
func TestNotification_Locale(t *testing.T) {
	enNotification := models.Notification{Locale: "en"}
	arNotification := models.Notification{Locale: "ar"}

	assert.Equal(t, "en", enNotification.Locale)
	assert.Equal(t, "ar", arNotification.Locale)
}

// TestNotification_MetricAndThreshold tests metric value and threshold (unit test)
func TestNotification_MetricAndThreshold(t *testing.T) {
	notification := models.Notification{
		MetricValue: 75000.0,
		Threshold:   50000.0,
	}

	assert.GreaterOrEqual(t, notification.MetricValue, notification.Threshold)
}
