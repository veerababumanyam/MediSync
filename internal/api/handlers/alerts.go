// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the alert rules and notification management endpoints.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// AlertsHandler handles alert rules and notification endpoints.
type AlertsHandler struct {
	logger     *slog.Logger
	alertRepo  *warehouse.AlertRuleRepository
	notifRepo  *warehouse.NotificationRepository
	userPrefRepo *warehouse.UserPreferenceRepository
}

// AlertsHandlerConfig holds configuration for the AlertsHandler.
type AlertsHandlerConfig struct {
	Logger       *slog.Logger
	AlertRepo    *warehouse.AlertRuleRepository
	NotifRepo    *warehouse.NotificationRepository
	UserPrefRepo *warehouse.UserPreferenceRepository
}

// NewAlertsHandler creates a new AlertsHandler instance.
func NewAlertsHandler(cfg AlertsHandlerConfig) *AlertsHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &AlertsHandler{
		logger:       cfg.Logger,
		alertRepo:    cfg.AlertRepo,
		notifRepo:    cfg.NotifRepo,
		userPrefRepo: cfg.UserPrefRepo,
	}
}

// RegisterRoutes registers alert routes on the given router.
func (h *AlertsHandler) RegisterRoutes(r chi.Router) {
	// Alert rules
	r.Route("/alerts/rules", func(r chi.Router) {
		r.Get("/", h.HandleListRules)
		r.Post("/", h.HandleCreateRule)
		r.Get("/{rule_id}", h.HandleGetRule)
		r.Patch("/{rule_id}", h.HandleUpdateRule)
		r.Delete("/{rule_id}", h.HandleDeleteRule)
		r.Post("/{rule_id}/toggle", h.HandleToggleRule)
		r.Post("/{rule_id}/test", h.HandleTestRule)
	})

	// Metrics for alert rules
	r.Get("/alerts/metrics", h.HandleGetMetrics)

	// Notifications
	r.Route("/notifications", func(r chi.Router) {
		r.Get("/", h.HandleListNotifications)
		r.Post("/{notification_id}/read", h.HandleMarkNotificationRead)
		r.Post("/read-all", h.HandleMarkAllNotificationsRead)
		r.Get("/unread-count", h.HandleGetUnreadCount)
	})
}

// ============================================================================
// Request/Response Types
// ============================================================================

// CreateAlertRuleRequest represents a request to create an alert rule.
type CreateAlertRuleRequest struct {
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	MetricID      string   `json:"metricId"`
	MetricName    string   `json:"metricName"`
	Operator      string   `json:"operator"` // gt, gte, lt, lte, eq
	Threshold     float64  `json:"threshold"`
	CheckInterval int      `json:"checkInterval"` // seconds
	Channels      []string `json:"channels"`      // email, in_app
	CooldownPeriod int     `json:"cooldownPeriod"` // seconds
}

// UpdateAlertRuleRequest represents a request to update an alert rule.
type UpdateAlertRuleRequest struct {
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Threshold      *float64 `json:"threshold,omitempty"`
	CheckInterval  *int     `json:"checkInterval,omitempty"`
	Channels       []string `json:"channels,omitempty"`
	CooldownPeriod *int     `json:"cooldownPeriod,omitempty"`
}

// ToggleRuleRequest represents a request to toggle an alert rule.
type ToggleRuleRequest struct {
	IsActive bool `json:"isActive"`
}

// AlertRuleResponse represents an alert rule in API responses.
type AlertRuleResponse struct {
	ID               string   `json:"id"`
	UserID           string   `json:"userId"`
	Name             string   `json:"name"`
	Description      *string  `json:"description,omitempty"`
	MetricID         string   `json:"metricId"`
	MetricName       string   `json:"metricName"`
	Operator         string   `json:"operator"`
	Threshold        float64  `json:"threshold"`
	CheckInterval    int      `json:"checkInterval"`
	Channels         []string `json:"channels"`
	Locale           string   `json:"locale"`
	CooldownPeriod   int      `json:"cooldownPeriod"`
	LastTriggeredAt  *string  `json:"lastTriggeredAt,omitempty"`
	LastValue        *float64 `json:"lastValue,omitempty"`
	IsActive         bool     `json:"isActive"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
}

// NotificationResponse represents a notification in API responses.
type NotificationResponse struct {
	ID           string                 `json:"id"`
	AlertRuleID  string                 `json:"alertRuleId"`
	UserID       string                 `json:"userId"`
	Type         string                 `json:"type"`
	Status       string                 `json:"status"`
	Content      map[string]interface{} `json:"content"`
	Locale       string                 `json:"locale"`
	MetricValue  float64                `json:"metricValue"`
	Threshold    float64                `json:"threshold"`
	ErrorMessage *string                `json:"errorMessage,omitempty"`
	SentAt       *string                `json:"sentAt,omitempty"`
	DeliveredAt  *string                `json:"deliveredAt,omitempty"`
	ReadAt       *string                `json:"readAt,omitempty"`
	CreatedAt    string                 `json:"createdAt"`
}

// MetricResponse represents an available metric for alert rules.
type MetricResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// UnreadCountResponse represents the unread notification count.
type UnreadCountResponse struct {
	Count int `json:"count"`
}

// ============================================================================
// Alert Rules Handlers
// ============================================================================

// HandleListRules handles GET /alerts/rules requests.
func (h *AlertsHandler) HandleListRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	activeOnly := r.URL.Query().Get("active_only") == "true"

	rules, err := h.alertRepo.GetByUserID(ctx, userID, activeOnly)
	if err != nil {
		h.logger.Error("failed to get alert rules",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve rules")
		return
	}

	response := make([]AlertRuleResponse, len(rules))
	for i, rule := range rules {
		response[i] = h.ruleToResponse(rule)
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleCreateRule handles POST /alerts/rules requests.
func (h *AlertsHandler) HandleCreateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateAlertRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.MetricID == "" {
		h.writeError(w, http.StatusBadRequest, "name and metricId are required")
		return
	}

	// Get user's locale
	locale := "en"
	if h.userPrefRepo != nil {
		pref, err := h.userPrefRepo.GetByUserID(ctx, userID)
		if err == nil && pref != nil {
			locale = pref.Locale
		}
	}

	rule := &models.AlertRule{
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		MetricID:         req.MetricID,
		MetricName:       req.MetricName,
		Operator:         req.Operator,
		Threshold:        req.Threshold,
		CheckInterval:    req.CheckInterval,
		Channels:         req.Channels,
		Locale:           locale,
		CooldownPeriod:   req.CooldownPeriod,
		IsActive:         true,
	}

	if err := h.alertRepo.Create(ctx, rule); err != nil {
		h.logger.Error("failed to create alert rule",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to create rule")
		return
	}

	h.writeJSON(w, http.StatusCreated, h.ruleToResponse(rule))
}

// HandleGetRule handles GET /alerts/rules/{rule_id} requests.
func (h *AlertsHandler) HandleGetRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ruleIDStr := chi.URLParam(r, "rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	rule, err := h.alertRepo.GetByID(ctx, ruleID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "rule not found")
		return
	}

	// Verify ownership
	if rule.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, h.ruleToResponse(rule))
}

// HandleUpdateRule handles PATCH /alerts/rules/{rule_id} requests.
func (h *AlertsHandler) HandleUpdateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ruleIDStr := chi.URLParam(r, "rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	rule, err := h.alertRepo.GetByID(ctx, ruleID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "rule not found")
		return
	}

	if rule.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	var req UpdateAlertRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Apply updates
	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = req.Description
	}
	if req.Threshold != nil {
		rule.Threshold = *req.Threshold
	}
	if req.CheckInterval != nil {
		rule.CheckInterval = *req.CheckInterval
	}
	if len(req.Channels) > 0 {
		rule.Channels = req.Channels
	}
	if req.CooldownPeriod != nil {
		rule.CooldownPeriod = *req.CooldownPeriod
	}

	if err := h.alertRepo.Update(ctx, rule); err != nil {
		h.logger.Error("failed to update alert rule",
			slog.Any("error", err),
			slog.String("rule_id", ruleID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to update rule")
		return
	}

	h.writeJSON(w, http.StatusOK, h.ruleToResponse(rule))
}

// HandleDeleteRule handles DELETE /alerts/rules/{rule_id} requests.
func (h *AlertsHandler) HandleDeleteRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ruleIDStr := chi.URLParam(r, "rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	if err := h.alertRepo.Delete(ctx, ruleID, userID); err != nil {
		if err.Error() == "warehouse: alert rule not found" {
			h.writeError(w, http.StatusNotFound, "rule not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to delete rule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleToggleRule handles POST /alerts/rules/{rule_id}/toggle requests.
func (h *AlertsHandler) HandleToggleRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ruleIDStr := chi.URLParam(r, "rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	var req ToggleRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.alertRepo.Toggle(ctx, ruleID, userID, req.IsActive); err != nil {
		if err.Error() == "warehouse: alert rule not found" {
			h.writeError(w, http.StatusNotFound, "rule not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to toggle rule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleTestRule handles POST /alerts/rules/{rule_id}/test requests.
func (h *AlertsHandler) HandleTestRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ruleIDStr := chi.URLParam(r, "rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	rule, err := h.alertRepo.GetByID(ctx, ruleID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "rule not found")
		return
	}

	if rule.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	// Create a test notification
	notif := &models.Notification{
		AlertRuleID: ruleID,
		UserID:      userID,
		Type:        "in_app",
		Status:      "delivered",
		Content: models.NotificationContent{
			Title:   "Test Alert: " + rule.Name,
			Message: "This is a test notification for your alert rule.",
		},
		Locale:      rule.Locale,
		MetricValue: rule.Threshold * 1.1, // Simulate threshold breach
		Threshold:   rule.Threshold,
	}

	if err := h.notifRepo.Create(ctx, notif); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to create test notification")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"notification": h.notificationToResponse(notif),
	})
}

// HandleGetMetrics handles GET /alerts/metrics requests.
func (h *AlertsHandler) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	// Return available metrics for alert rules
	metrics := []MetricResponse{
		{ID: "revenue_total", Name: "Total Revenue", Description: "Total revenue across all departments", Category: "financial"},
		{ID: "revenue_clinic", Name: "Clinic Revenue", Description: "Revenue from clinic services", Category: "financial"},
		{ID: "revenue_pharmacy", Name: "Pharmacy Revenue", Description: "Revenue from pharmacy sales", Category: "financial"},
		{ID: "patients_new", Name: "New Patients", Description: "Number of new patient registrations", Category: "operations"},
		{ID: "patients_total", Name: "Total Patients", Description: "Total number of patients", Category: "operations"},
		{ID: "appointments_today", Name: "Today's Appointments", Description: "Number of appointments today", Category: "operations"},
		{ID: "inventory_low", Name: "Low Stock Items", Description: "Number of items below reorder level", Category: "inventory"},
		{ID: "inventory_expired", Name: "Expired Items", Description: "Number of expired inventory items", Category: "inventory"},
	}

	h.writeJSON(w, http.StatusOK, metrics)
}

// ============================================================================
// Notifications Handlers
// ============================================================================

// HandleListNotifications handles GET /notifications requests.
func (h *AlertsHandler) HandleListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	unreadOnly := r.URL.Query().Get("unread_only") == "true"
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		var lim int
		if _, err := json.Number(l).Int64(); err == nil {
			if lim > 0 && lim <= 100 {
				limit = lim
			}
		}
	}

	notifications, err := h.notifRepo.GetByUserID(ctx, userID, unreadOnly, limit)
	if err != nil {
		h.logger.Error("failed to get notifications",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve notifications")
		return
	}

	response := make([]NotificationResponse, len(notifications))
	for i, notif := range notifications {
		response[i] = h.notificationToResponse(notif)
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleMarkNotificationRead handles POST /notifications/{notification_id}/read requests.
func (h *AlertsHandler) HandleMarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notifIDStr := chi.URLParam(r, "notification_id")
	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid notification ID")
		return
	}

	if err := h.notifRepo.MarkAsRead(ctx, notifID, userID); err != nil {
		if err.Error() == "warehouse: notification not found" {
			h.writeError(w, http.StatusNotFound, "notification not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleMarkAllNotificationsRead handles POST /notifications/read-all requests.
func (h *AlertsHandler) HandleMarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	count, err := h.notifRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to mark notifications as read")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"markedRead": count,
	})
}

// HandleGetUnreadCount handles GET /notifications/unread-count requests.
func (h *AlertsHandler) HandleGetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	count, err := h.notifRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to get unread count")
		return
	}

	h.writeJSON(w, http.StatusOK, UnreadCountResponse{Count: count})
}

// ============================================================================
// Helper Functions
// ============================================================================

func (h *AlertsHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (h *AlertsHandler) ruleToResponse(rule *models.AlertRule) AlertRuleResponse {
	var lastTriggeredAt *string
	if rule.LastTriggeredAt != nil {
		t := rule.LastTriggeredAt.Format("2006-01-02T15:04:05Z07:00")
		lastTriggeredAt = &t
	}

	return AlertRuleResponse{
		ID:              rule.ID.String(),
		UserID:          rule.UserID.String(),
		Name:            rule.Name,
		Description:     rule.Description,
		MetricID:        rule.MetricID,
		MetricName:      rule.MetricName,
		Operator:        rule.Operator,
		Threshold:       rule.Threshold,
		CheckInterval:   rule.CheckInterval,
		Channels:        rule.Channels,
		Locale:          rule.Locale,
		CooldownPeriod:  rule.CooldownPeriod,
		LastTriggeredAt: lastTriggeredAt,
		LastValue:       rule.LastValue,
		IsActive:        rule.IsActive,
		CreatedAt:       rule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       rule.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *AlertsHandler) notificationToResponse(notif *models.Notification) NotificationResponse {
	var sentAt, deliveredAt, readAt *string
	if notif.SentAt != nil {
		t := notif.SentAt.Format("2006-01-02T15:04:05Z07:00")
		sentAt = &t
	}
	if notif.DeliveredAt != nil {
		t := notif.DeliveredAt.Format("2006-01-02T15:04:05Z07:00")
		deliveredAt = &t
	}
	if notif.ReadAt != nil {
		t := notif.ReadAt.Format("2006-01-02T15:04:05Z07:00")
		readAt = &t
	}

	content := map[string]interface{}{
		"title":   notif.Content.Title,
		"message": notif.Content.Message,
	}
	if notif.Content.ActionURL != "" {
		content["actionUrl"] = notif.Content.ActionURL
	}

	return NotificationResponse{
		ID:           notif.ID.String(),
		AlertRuleID:  notif.AlertRuleID.String(),
		UserID:       notif.UserID.String(),
		Type:         notif.Type,
		Status:       notif.Status,
		Content:      content,
		Locale:       notif.Locale,
		MetricValue:  notif.MetricValue,
		Threshold:    notif.Threshold,
		ErrorMessage: notif.ErrorMessage,
		SentAt:       sentAt,
		DeliveredAt:  deliveredAt,
		ReadAt:       readAt,
		CreatedAt:    notif.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *AlertsHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AlertsHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}
