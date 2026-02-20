package warehouse

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/stretchr/testify/assert"
)

// TestAlertRuleModel_Fields tests model fields (unit test)
func TestAlertRuleModel_Fields(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	userID := uuid.New()
	description := "Alert when revenue drops"
	threshold := 10000.0
	lastValue := 9500.0

	rule := models.AlertRule{
		ID:              id,
		UserID:          userID,
		Name:            "Revenue Alert",
		Description:     &description,
		MetricID:        "total_revenue",
		MetricName:      "Total Revenue",
		Operator:        "<",
		Threshold:       threshold,
		CheckInterval:   300,
		Channels:        []string{"email", "slack"},
		Locale:          "en",
		CooldownPeriod:  3600,
		LastTriggeredAt: &now,
		LastValue:       &lastValue,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	assert.Equal(t, id, rule.ID)
	assert.Equal(t, userID, rule.UserID)
	assert.Equal(t, "Revenue Alert", rule.Name)
	assert.Equal(t, "Alert when revenue drops", *rule.Description)
	assert.Equal(t, "total_revenue", rule.MetricID)
	assert.Equal(t, "<", rule.Operator)
	assert.Equal(t, 10000.0, rule.Threshold)
	assert.Equal(t, 300, rule.CheckInterval)
	assert.Equal(t, []string{"email", "slack"}, rule.Channels)
	assert.True(t, rule.IsActive)
}

// TestAlertRule_Operators tests operator values (unit test)
func TestAlertRule_Operators(t *testing.T) {
	operators := []string{">", "<", ">=", "<=", "==", "!="}

	for _, op := range operators {
		t.Run(op, func(t *testing.T) {
			rule := models.AlertRule{Operator: op}
			assert.Equal(t, op, rule.Operator)
		})
	}
}

// TestAlertRule_Channels tests notification channels (unit test)
func TestAlertRule_Channels(t *testing.T) {
	tests := []struct {
		name     string
		channels []string
	}{
		{"email only", []string{"email"}},
		{"slack only", []string{"slack"}},
		{"multiple channels", []string{"email", "slack", "sms"}},
		{"no channels", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := models.AlertRule{Channels: tt.channels}
			assert.Equal(t, tt.channels, rule.Channels)
		})
	}
}

// TestAlertRule_WithNilDescription tests nil description (unit test)
func TestAlertRule_WithNilDescription(t *testing.T) {
	rule := models.AlertRule{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "Simple Alert",
		Description: nil,
	}

	assert.Nil(t, rule.Description)
}

// TestAlertRule_WithNilLastValue tests nil last value (unit test)
func TestAlertRule_WithNilLastValue(t *testing.T) {
	rule := models.AlertRule{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		LastValue: nil,
	}

	assert.Nil(t, rule.LastValue)
}

// TestAlertRule_CheckIntervals tests common check intervals (unit test)
func TestAlertRule_CheckIntervals(t *testing.T) {
	intervals := []struct {
		name     string
		seconds  int
	}{
		{"1 minute", 60},
		{"5 minutes", 300},
		{"15 minutes", 900},
		{"30 minutes", 1800},
		{"1 hour", 3600},
		{"6 hours", 21600},
		{"1 day", 86400},
	}

	for _, tt := range intervals {
		t.Run(tt.name, func(t *testing.T) {
			rule := models.AlertRule{CheckInterval: tt.seconds}
			assert.Equal(t, tt.seconds, rule.CheckInterval)
		})
	}
}

// TestAlertRule_CooldownPeriod tests cooldown period (unit test)
func TestAlertRule_CooldownPeriod(t *testing.T) {
	rule := models.AlertRule{
		CooldownPeriod: 3600, // 1 hour
	}

	assert.Equal(t, 3600, rule.CooldownPeriod)
}

// TestAlertRule_Locale tests locale field (unit test)
func TestAlertRule_Locale(t *testing.T) {
	enRule := models.AlertRule{Locale: "en"}
	arRule := models.AlertRule{Locale: "ar"}

	assert.Equal(t, "en", enRule.Locale)
	assert.Equal(t, "ar", arRule.Locale)
}

// TestAlertRule_IsActive tests active flag (unit test)
func TestAlertRule_IsActive(t *testing.T) {
	activeRule := models.AlertRule{IsActive: true}
	inactiveRule := models.AlertRule{IsActive: false}

	assert.True(t, activeRule.IsActive)
	assert.False(t, inactiveRule.IsActive)
}

// TestAlertRule_Timestamps tests timestamp fields (unit test)
func TestAlertRule_Timestamps(t *testing.T) {
	now := time.Now()
	triggered := now.Add(-1 * time.Hour)

	rule := models.AlertRule{
		CreatedAt:       now,
		UpdatedAt:       now,
		LastTriggeredAt: &triggered,
	}

	assert.Equal(t, now, rule.CreatedAt)
	assert.Equal(t, now, rule.UpdatedAt)
	assert.Equal(t, triggered, *rule.LastTriggeredAt)
}
