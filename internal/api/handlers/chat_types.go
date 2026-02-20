// Package handlers provides HTTP handlers for the MediSync API.
//
// This file contains the request and response types for the Chat API,
// matching the OpenAPI specification in contracts/chat-api.yaml.
package handlers

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// Request Types
// ============================================================================

// ChatRequest represents a natural language query request.
// It matches the ChatRequest schema from the OpenAPI specification.
type ChatRequest struct {
	// Query is the natural language query in English or Arabic.
	// Required field, min 1 character, max 2000 characters.
	Query string `json:"query"`

	// Locale is the preferred response locale ("en" or "ar").
	// Optional, defaults to "en".
	Locale string `json:"locale,omitempty"`

	// SessionID is the UUID for conversation continuity.
	// Optional, if not provided a new session will be created.
	SessionID string `json:"session_id,omitempty"`

	// Context provides previous query context for follow-up questions.
	// Optional, max 10 items.
	Context []QueryContext `json:"context,omitempty"`
}

// QueryContext represents context from a previous query.
type QueryContext struct {
	// Query is the previous query text.
	Query string `json:"query,omitempty"`

	// ResultSummary is a brief summary of the previous result.
	ResultSummary string `json:"result_summary,omitempty"`
}

// Validate checks if the ChatRequest has valid field values.
func (r *ChatRequest) Validate() error {
	if r.Query == "" {
		return fmt.Errorf("query is required")
	}

	if len(r.Query) > 2000 {
		return fmt.Errorf("query exceeds maximum length of 2000 characters")
	}

	if r.Locale != "" && r.Locale != "en" && r.Locale != "ar" {
		return fmt.Errorf("locale must be 'en' or 'ar', got '%s'", r.Locale)
	}

	if len(r.Context) > 10 {
		return fmt.Errorf("context exceeds maximum of 10 items")
	}

	return nil
}

// GetLocale returns the locale or default "en" if not set.
func (r *ChatRequest) GetLocale() string {
	if r.Locale == "" {
		return "en"
	}
	return r.Locale
}

// ============================================================================
// SSE Event Types
// ============================================================================

// EventType defines the types of Server-Sent Events.
type EventType string

const (
	// EventTypeThinking indicates agent processing status updates.
	EventTypeThinking EventType = "thinking"

	// EventTypeSQLPreview provides the generated SQL query for transparency.
	EventTypeSQLPreview EventType = "sql_preview"

	// EventTypeResult delivers query results with chart recommendation.
	EventTypeResult EventType = "result"

	// EventTypeError indicates an error message.
	EventTypeError EventType = "error"

	// EventTypeClarification requests user clarification for low confidence.
	EventTypeClarification EventType = "clarification"
)

// ChartType defines the supported visualization types.
type ChartType string

const (
	ChartTypeLineChart    ChartType = "lineChart"
	ChartTypeBarChart     ChartType = "barChart"
	ChartTypePieChart     ChartType = "pieChart"
	ChartTypeKPICard      ChartType = "kpiCard"
	ChartTypeDataTable    ChartType = "dataTable"
	ChartTypeScatterChart ChartType = "scatterChart"
)

// SSEEvent represents a Server-Sent Event.
// It matches the SSEEvent schema from the OpenAPI specification.
type SSEEvent struct {
	// Type is the event type (required).
	Type EventType `json:"type"`

	// Message is a human-readable message for thinking/error/clarification events.
	Message string `json:"message,omitempty"`

	// SQL is the generated SQL query for sql_preview events.
	SQL string `json:"sql,omitempty"`

	// ChartType is the recommended visualization type for result events.
	ChartType ChartType `json:"chart_type,omitempty"`

	// Data is the query result data for result events.
	Data interface{} `json:"data,omitempty"`

	// Confidence is the confidence score percentage (0-100) for result events.
	Confidence float64 `json:"confidence,omitempty"`

	// Options are clarification options for clarification events.
	Options []string `json:"options,omitempty"`
}

// ToJSON serializes the SSEEvent to JSON.
func (e *SSEEvent) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SSE event: %w", err)
	}
	return string(data), nil
}

// NewThinkingEvent creates a new thinking event.
func NewThinkingEvent(message string) *SSEEvent {
	return &SSEEvent{
		Type:    EventTypeThinking,
		Message: message,
	}
}

// NewSQLPreviewEvent creates a new SQL preview event.
func NewSQLPreviewEvent(sql string) *SSEEvent {
	return &SSEEvent{
		Type: EventTypeSQLPreview,
		SQL:  sql,
	}
}

// NewResultEvent creates a new result event.
func NewResultEvent(chartType ChartType, data interface{}, confidence float64) *SSEEvent {
	return &SSEEvent{
		Type:       EventTypeResult,
		ChartType:  chartType,
		Data:       data,
		Confidence: confidence,
	}
}

// NewErrorEvent creates a new error event.
func NewErrorEvent(message string) *SSEEvent {
	return &SSEEvent{
		Type:    EventTypeError,
		Message: message,
	}
}

// NewClarificationEvent creates a new clarification request event.
func NewClarificationEvent(message string, options []string) *SSEEvent {
	return &SSEEvent{
		Type:    EventTypeClarification,
		Message: message,
		Options: options,
	}
}

// ============================================================================
// Chart Data Types
// ============================================================================

// ChartData represents data formatted for visualization.
type ChartData struct {
	// Columns contains metadata for each column.
	Columns []ChartColumn `json:"columns,omitempty"`

	// Rows contains the data rows.
	Rows []map[string]interface{} `json:"rows,omitempty"`

	// For KPI cards, single value display.
	Value     interface{} `json:"value,omitempty"`
	Formatted string      `json:"formatted,omitempty"`

	// For time series, x-axis labels.
	Labels []string `json:"labels,omitempty"`

	// For charts with series.
	Series []ChartSeries `json:"series,omitempty"`
}

// ChartColumn represents a column in chart data.
type ChartColumn struct {
	Name string `json:"name"`
	Type string `json:"type"` // "string", "number", "date"
}

// ChartSeries represents a data series for charts.
type ChartSeries struct {
	Name   string        `json:"name"`
	Values []interface{} `json:"values"`
}

// ============================================================================
// Error Response Types
// ============================================================================

// ErrorCode defines standardized error codes matching the OpenAPI spec.
type ErrorCode string

const (
	ErrorCodeInvalidRequest  ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeRateLimited     ErrorCode = "RATE_LIMITED"
	ErrorCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrorCodeLLMUnavailable  ErrorCode = "LLM_UNAVAILABLE"
)

// ErrorResponse represents an API error response.
// It matches the ErrorResponse schema from the OpenAPI specification.
type ErrorResponse struct {
	// Code is the error code.
	Code ErrorCode `json:"code"`

	// Message is the human-readable error message.
	Message string `json:"message"`

	// RetryAfter is the seconds to wait before retry (for RATE_LIMITED).
	RetryAfter int `json:"retry_after,omitempty"`
}

// ToJSON serializes the ErrorResponse to JSON.
func (e *ErrorResponse) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("failed to marshal error response: %w", err)
	}
	return string(data), nil
}

// ============================================================================
// Agent Health Types
// ============================================================================

// AgentStatus defines the health status of an agent.
type AgentStatus string

const (
	AgentStatusHealthy   AgentStatus = "healthy"
	AgentStatusDegraded  AgentStatus = "degraded"
	AgentStatusUnhealthy AgentStatus = "unhealthy"
)

// LLMProviderStatus defines the status of the LLM provider.
type LLMProviderStatus string

const (
	LLMStatusConnected    LLMProviderStatus = "connected"
	LLMStatusDisconnected LLMProviderStatus = "disconnected"
	LLMStatusRateLimited  LLMProviderStatus = "rate_limited"
)

// AgentHealth represents the health status of a single agent.
type AgentHealth struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Status       AgentStatus  `json:"status"`
	LastCheck    string       `json:"last_check,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

// LLMProviderHealth represents the health of the LLM provider.
type LLMProviderHealth struct {
	Name   string             `json:"name,omitempty"`
	Model  string             `json:"model,omitempty"`
	Status LLMProviderStatus  `json:"status,omitempty"`
}

// AgentsHealthResponse represents the response for the agent health endpoint.
type AgentsHealthResponse struct {
	Status     AgentStatus         `json:"status"`
	Timestamp  string              `json:"timestamp"`
	Agents     []AgentHealth       `json:"agents"`
	LLMProvider *LLMProviderHealth `json:"llm_provider,omitempty"`
}

// ============================================================================
// WebSocket Types
// ============================================================================

// WSMessageType defines the types of WebSocket messages.
type WSMessageType string

const (
	// Client messages
	WSMessageTypeQuery WSMessageType = "query"
	WSMessageTypePong  WSMessageType = "pong"

	// Server messages
	WSMessageTypeThinking      WSMessageType = "thinking"
	WSMessageTypeSQLPreview    WSMessageType = "sql_preview"
	WSMessageTypeResult        WSMessageType = "result"
	WSMessageTypeError         WSMessageType = "error"
	WSMessageTypeClarification WSMessageType = "clarification"
	WSMessageTypePing          WSMessageType = "ping"
)

// WSMessage represents a WebSocket message.
type WSMessage struct {
	Type    WSMessageType     `json:"type"`
	Payload *WSMessagePayload `json:"payload,omitempty"`
}

// WSMessagePayload represents the payload of a WebSocket message.
type WSMessagePayload struct {
	Query        string                 `json:"query,omitempty"`
	Locale       string                 `json:"locale,omitempty"`
	Message      string                 `json:"message,omitempty"`
	SQL          string                 `json:"sql,omitempty"`
	ChartType    string                 `json:"chart_type,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Confidence   float64                `json:"confidence,omitempty"`
	Options      []string               `json:"options,omitempty"`
	SessionID    string                 `json:"session_id,omitempty"`
}
