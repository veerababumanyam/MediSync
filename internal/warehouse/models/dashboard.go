// Package models provides data models for MediSync.
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserPreference represents user-specific display and formatting preferences.
type UserPreference struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"userId" db:"user_id"`
	Locale         string    `json:"locale" db:"locale"`
	NumeralSystem  string    `json:"numeralSystem" db:"numeral_system"`
	CalendarSystem string    `json:"calendarSystem" db:"calendar_system"`
	ReportLanguage string    `json:"reportLanguage" db:"report_language"`
	Timezone       string    `json:"timezone" db:"timezone"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

// ChartPosition represents the position of a pinned chart on the dashboard.
type ChartPosition struct {
	Row  int `json:"row"`
	Col  int `json:"col"`
	Size int `json:"size"`
}

// PinnedChart represents a saved visualization on a user's dashboard.
type PinnedChart struct {
	ID                   uuid.UUID      `json:"id" db:"id"`
	UserID               uuid.UUID      `json:"userId" db:"user_id"`
	Title                string         `json:"title" db:"title"`
	QueryID              *uuid.UUID     `json:"queryId" db:"query_id"`
	NaturalLanguageQuery string         `json:"naturalLanguageQuery" db:"natural_language_query"`
	SQLQuery             string         `json:"sqlQuery" db:"sql_query"`
	ChartSpec            map[string]any `json:"chartSpec" db:"chart_spec"`
	ChartType            string         `json:"chartType" db:"chart_type"`
	RefreshInterval      int            `json:"refreshInterval" db:"refresh_interval"`
	Locale               string         `json:"locale" db:"locale"`
	Position             ChartPosition  `json:"position" db:"position"`
	LastRefreshedAt      *time.Time     `json:"lastRefreshedAt" db:"last_refreshed_at"`
	IsActive             bool           `json:"isActive" db:"is_active"`
	CreatedAt            time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time      `json:"updatedAt" db:"updated_at"`
}

// AlertRule represents a user-defined condition that triggers notifications.
type AlertRule struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	UserID           uuid.UUID      `json:"userId" db:"user_id"`
	Name             string         `json:"name" db:"name"`
	Description      *string        `json:"description" db:"description"`
	MetricID         string         `json:"metricId" db:"metric_id"`
	MetricName       string         `json:"metricName" db:"metric_name"`
	Operator         string         `json:"operator" db:"operator"`
	Threshold        float64        `json:"threshold" db:"threshold"`
	CheckInterval    int            `json:"checkInterval" db:"check_interval"`
	Channels         []string       `json:"channels" db:"channels"`
	Locale           string         `json:"locale" db:"locale"`
	CooldownPeriod   int            `json:"cooldownPeriod" db:"cooldown_period"`
	LastTriggeredAt  *time.Time     `json:"lastTriggeredAt" db:"last_triggered_at"`
	LastValue        *float64       `json:"lastValue" db:"last_value"`
	IsActive         bool           `json:"isActive" db:"is_active"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`
}

// NotificationContent represents the content of a notification.
type NotificationContent struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	ActionURL string `json:"actionUrl,omitempty"`
}

// Notification represents a record of an alert delivery attempt.
type Notification struct {
	ID            uuid.UUID            `json:"id" db:"id"`
	AlertRuleID   uuid.UUID            `json:"alertRuleId" db:"alert_rule_id"`
	UserID        uuid.UUID            `json:"userId" db:"user_id"`
	Type          string               `json:"type" db:"type"`
	Status        string               `json:"status" db:"status"`
	Content       NotificationContent  `json:"content" db:"content"`
	Locale        string               `json:"locale" db:"locale"`
	MetricValue   float64              `json:"metricValue" db:"metric_value"`
	Threshold     float64              `json:"threshold" db:"threshold"`
	ErrorMessage  *string              `json:"errorMessage" db:"error_message"`
	SentAt        *time.Time           `json:"sentAt" db:"sent_at"`
	DeliveredAt   *time.Time           `json:"deliveredAt" db:"delivered_at"`
	ReadAt        *time.Time           `json:"readAt" db:"read_at"`
	CreatedAt     time.Time            `json:"createdAt" db:"created_at"`
}

// Recipient represents an email recipient for a scheduled report.
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ScheduledReport represents a recurring report configuration.
type ScheduledReport struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	UserID               uuid.UUID    `json:"userId" db:"user_id"`
	Name                 string       `json:"name" db:"name"`
	Description          *string      `json:"description" db:"description"`
	QueryID              *uuid.UUID   `json:"queryId" db:"query_id"`
	NaturalLanguageQuery string       `json:"naturalLanguageQuery" db:"natural_language_query"`
	SQLQuery             string       `json:"sqlQuery" db:"sql_query"`
	ScheduleType         string       `json:"scheduleType" db:"schedule_type"`
	ScheduleTime         string       `json:"scheduleTime" db:"schedule_time"`
	ScheduleDay          *int         `json:"scheduleDay" db:"schedule_day"`
	Recipients           []Recipient  `json:"recipients" db:"recipients"`
	Format               string       `json:"format" db:"format"`
	Locale               string       `json:"locale" db:"locale"`
	IncludeCharts        bool         `json:"includeCharts" db:"include_charts"`
	LastRunAt            *time.Time   `json:"lastRunAt" db:"last_run_at"`
	NextRunAt            *time.Time   `json:"nextRunAt" db:"next_run_at"`
	IsActive             bool         `json:"isActive" db:"is_active"`
	CreatedAt            time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time    `json:"updatedAt" db:"updated_at"`
}

// ScheduledReportRun represents an audit record of a report generation attempt.
type ScheduledReportRun struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ReportID      uuid.UUID  `json:"reportId" db:"report_id"`
	Status        string     `json:"status" db:"status"`
	FilePath      *string    `json:"filePath" db:"file_path"`
	FileSizeBytes *int64     `json:"fileSizeBytes" db:"file_size_bytes"`
	RowCount      *int       `json:"rowCount" db:"row_count"`
	ErrorMessage  *string    `json:"errorMessage" db:"error_message"`
	StartedAt     time.Time  `json:"startedAt" db:"started_at"`
	CompletedAt   *time.Time `json:"completedAt" db:"completed_at"`
}

// ChatMessage represents a message in a chat conversation.
type ChatMessage struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	SessionID       uuid.UUID      `json:"sessionId" db:"session_id"`
	UserID          uuid.UUID      `json:"userId" db:"user_id"`
	Role            string         `json:"role" db:"role"`
	Content         string         `json:"content" db:"content"`
	ChartSpec       map[string]any `json:"chartSpec" db:"chart_spec"`
	TableData       map[string]any `json:"tableData" db:"table_data"`
	DrilldownQuery  *string        `json:"drilldownQuery" db:"drilldown_query"`
	ConfidenceScore *float64       `json:"confidenceScore" db:"confidence_score"`
	Locale          string         `json:"locale" db:"locale"`
	CreatedAt       time.Time      `json:"createdAt" db:"created_at"`
}

// MarshalJSON custom marshaler for ChartPosition
func (p ChartPosition) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]int{
		"row":  p.Row,
		"col":  p.Col,
		"size": p.Size,
	})
}

// UnmarshalJSON custom unmarshaler for ChartPosition
func (p *ChartPosition) UnmarshalJSON(data []byte) error {
	var m map[string]int
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	p.Row = m["row"]
	p.Col = m["col"]
	p.Size = m["size"]
	return nil
}
