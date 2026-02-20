// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the dashboard widget management endpoints for
// pinned charts, quick actions, and dashboard configuration.
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

// DashboardHandler handles dashboard management endpoints.
type DashboardHandler struct {
	logger           *slog.Logger
	pinnedChartRepo  *warehouse.PinnedChartRepository
	userPrefRepo     *warehouse.UserPreferenceRepository
}

// DashboardHandlerConfig holds configuration for the DashboardHandler.
type DashboardHandlerConfig struct {
	Logger          *slog.Logger
	PinnedChartRepo *warehouse.PinnedChartRepository
	UserPrefRepo    *warehouse.UserPreferenceRepository
}

// NewDashboardHandler creates a new DashboardHandler instance.
func NewDashboardHandler(cfg DashboardHandlerConfig) *DashboardHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &DashboardHandler{
		logger:          cfg.Logger,
		pinnedChartRepo: cfg.PinnedChartRepo,
		userPrefRepo:    cfg.UserPrefRepo,
	}
}

// RegisterRoutes registers dashboard routes on the given router.
func (h *DashboardHandler) RegisterRoutes(r chi.Router) {
	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/charts", h.HandleListCharts)
		r.Post("/charts", h.HandleCreateChart)
		r.Get("/charts/{chart_id}", h.HandleGetChart)
		r.Patch("/charts/{chart_id}", h.HandleUpdateChart)
		r.Delete("/charts/{chart_id}", h.HandleDeleteChart)
		r.Post("/charts/{chart_id}/refresh", h.HandleRefreshChart)
		r.Post("/charts/reorder", h.HandleReorderCharts)
		r.Get("/quick-actions", h.HandleGetQuickActions)
	})
}

// ============================================================================
// Request/Response Types
// ============================================================================

// CreateChartRequest represents a request to pin a new chart.
type CreateChartRequest struct {
	Title                string         `json:"title"`
	QueryID              *string        `json:"queryId,omitempty"`
	NaturalLanguageQuery string         `json:"naturalLanguageQuery"`
	SQLQuery             string         `json:"sqlQuery"`
	ChartSpec            map[string]any `json:"chartSpec"`
	ChartType            string         `json:"chartType"`
	RefreshInterval      int            `json:"refreshInterval"`
	Position             ChartPosition  `json:"position"`
}

// UpdateChartRequest represents a request to update a pinned chart.
type UpdateChartRequest struct {
	Title           *string        `json:"title,omitempty"`
	RefreshInterval *int           `json:"refreshInterval,omitempty"`
	Position        *ChartPosition `json:"position,omitempty"`
	IsActive        *bool          `json:"isActive,omitempty"`
}

// ChartPosition represents the position of a chart on the dashboard.
type ChartPosition struct {
	Row  int `json:"row"`
	Col  int `json:"col"`
	Size int `json:"size"`
}

// ReorderChartsRequest represents a request to reorder charts.
type ReorderChartsRequest struct {
	Positions []ChartPositionUpdate `json:"positions"`
}

// ChartPositionUpdate represents a position update for a single chart.
type ChartPositionUpdate struct {
	ID       string       `json:"id"`
	Position ChartPosition `json:"position"`
}

// ChartResponse represents a pinned chart in API responses.
type ChartResponse struct {
	ID                   string         `json:"id"`
	UserID               string         `json:"userId"`
	Title                string         `json:"title"`
	QueryID              *string        `json:"queryId,omitempty"`
	NaturalLanguageQuery string         `json:"naturalLanguageQuery"`
	SQLQuery             string         `json:"sqlQuery"`
	ChartSpec            map[string]any `json:"chartSpec"`
	ChartType            string         `json:"chartType"`
	RefreshInterval      int            `json:"refreshInterval"`
	Locale               string         `json:"locale"`
	Position             ChartPosition  `json:"position"`
	LastRefreshedAt      *string        `json:"lastRefreshedAt,omitempty"`
	IsActive             bool           `json:"isActive"`
	CreatedAt            string         `json:"createdAt"`
	UpdatedAt            string         `json:"updatedAt"`
}

// QuickActionResponse represents a quick action for the dashboard.
type QuickActionResponse struct {
	ID    string `json:"id"`
	Query string `json:"query"`
	Label string `json:"label"`
}

// ============================================================================
// HTTP Handlers
// ============================================================================

// HandleListCharts handles GET /dashboard/charts requests.
func (h *DashboardHandler) HandleListCharts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	activeOnly := r.URL.Query().Get("active_only") == "true"

	charts, err := h.pinnedChartRepo.GetByUserID(ctx, userID, activeOnly)
	if err != nil {
		h.logger.Error("failed to get pinned charts",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve charts")
		return
	}

	response := make([]ChartResponse, len(charts))
	for i, chart := range charts {
		response[i] = h.chartToResponse(chart)
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleCreateChart handles POST /dashboard/charts requests.
func (h *DashboardHandler) HandleCreateChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateChartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "title is required")
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

	chart := &models.PinnedChart{
		UserID:               userID,
		Title:                req.Title,
		NaturalLanguageQuery: req.NaturalLanguageQuery,
		SQLQuery:             req.SQLQuery,
		ChartSpec:            req.ChartSpec,
		ChartType:            req.ChartType,
		RefreshInterval:      req.RefreshInterval,
		Locale:               locale,
		Position: models.ChartPosition{
			Row:  req.Position.Row,
			Col:  req.Position.Col,
			Size: req.Position.Size,
		},
		IsActive: true,
	}

	if req.QueryID != nil {
		queryID, err := uuid.Parse(*req.QueryID)
		if err == nil {
			chart.QueryID = &queryID
		}
	}

	if err := h.pinnedChartRepo.Create(ctx, chart); err != nil {
		h.logger.Error("failed to create chart",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to create chart")
		return
	}

	h.writeJSON(w, http.StatusCreated, h.chartToResponse(chart))
}

// HandleGetChart handles GET /dashboard/charts/{chart_id} requests.
func (h *DashboardHandler) HandleGetChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chartIDStr := chi.URLParam(r, "chart_id")
	chartID, err := uuid.Parse(chartIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid chart ID")
		return
	}

	// Note: We need to implement GetByID in the repository
	// For now, we'll get all charts and filter
	charts, err := h.pinnedChartRepo.GetByUserID(ctx, userID, false)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve chart")
		return
	}

	for _, chart := range charts {
		if chart.ID == chartID {
			h.writeJSON(w, http.StatusOK, h.chartToResponse(chart))
			return
		}
	}

	h.writeError(w, http.StatusNotFound, "chart not found")
}

// HandleUpdateChart handles PATCH /dashboard/charts/{chart_id} requests.
func (h *DashboardHandler) HandleUpdateChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chartIDStr := chi.URLParam(r, "chart_id")
	chartID, err := uuid.Parse(chartIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid chart ID")
		return
	}

	var req UpdateChartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing chart
	charts, err := h.pinnedChartRepo.GetByUserID(ctx, userID, false)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve chart")
		return
	}

	var existingChart *models.PinnedChart
	for _, chart := range charts {
		if chart.ID == chartID {
			existingChart = chart
			break
		}
	}

	if existingChart == nil {
		h.writeError(w, http.StatusNotFound, "chart not found")
		return
	}

	// Apply updates
	if req.Title != nil {
		existingChart.Title = *req.Title
	}
	if req.RefreshInterval != nil {
		existingChart.RefreshInterval = *req.RefreshInterval
	}
	if req.Position != nil {
		existingChart.Position = models.ChartPosition{
			Row:  req.Position.Row,
			Col:  req.Position.Col,
			Size: req.Position.Size,
		}
	}
	if req.IsActive != nil {
		existingChart.IsActive = *req.IsActive
	}

	if err := h.pinnedChartRepo.Update(ctx, existingChart); err != nil {
		h.logger.Error("failed to update chart",
			slog.Any("error", err),
			slog.String("chart_id", chartID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to update chart")
		return
	}

	h.writeJSON(w, http.StatusOK, h.chartToResponse(existingChart))
}

// HandleDeleteChart handles DELETE /dashboard/charts/{chart_id} requests.
func (h *DashboardHandler) HandleDeleteChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chartIDStr := chi.URLParam(r, "chart_id")
	chartID, err := uuid.Parse(chartIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid chart ID")
		return
	}

	if err := h.pinnedChartRepo.Delete(ctx, chartID, userID); err != nil {
		if err.Error() == "warehouse: pinned chart not found" {
			h.writeError(w, http.StatusNotFound, "chart not found")
			return
		}
		h.logger.Error("failed to delete chart",
			slog.Any("error", err),
			slog.String("chart_id", chartID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to delete chart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleRefreshChart handles POST /dashboard/charts/{chart_id}/refresh requests.
func (h *DashboardHandler) HandleRefreshChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chartIDStr := chi.URLParam(r, "chart_id")
	chartID, err := uuid.Parse(chartIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid chart ID")
		return
	}

	// Verify chart exists and belongs to user
	charts, err := h.pinnedChartRepo.GetByUserID(ctx, userID, false)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve chart")
		return
	}

	found := false
	var chart *models.PinnedChart
	for _, c := range charts {
		if c.ID == chartID {
			found = true
			chart = c
			break
		}
	}

	if !found {
		h.writeError(w, http.StatusNotFound, "chart not found")
		return
	}

	// Update last refreshed timestamp
	if err := h.pinnedChartRepo.UpdateLastRefreshed(ctx, chartID); err != nil {
		h.logger.Error("failed to refresh chart",
			slog.Any("error", err),
			slog.String("chart_id", chartID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to refresh chart")
		return
	}

	// In a real implementation, you would re-execute the SQL query here
	// For now, we just update the timestamp

	h.writeJSON(w, http.StatusOK, h.chartToResponse(chart))
}

// HandleReorderCharts handles POST /dashboard/charts/reorder requests.
func (h *DashboardHandler) HandleReorderCharts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ReorderChartsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Positions) == 0 {
		h.writeError(w, http.StatusBadRequest, "no positions provided")
		return
	}

	// Convert to repository format
	positions := make([]struct {
		ID       uuid.UUID
		Position models.ChartPosition
	}, len(req.Positions))

	for i, p := range req.Positions {
		id, err := uuid.Parse(p.ID)
		if err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid chart ID in positions")
			return
		}
		positions[i] = struct {
			ID       uuid.UUID
			Position models.ChartPosition
		}{
			ID: id,
			Position: models.ChartPosition{
				Row:  p.Position.Row,
				Col:  p.Position.Col,
				Size: p.Position.Size,
			},
		}
	}

	if err := h.pinnedChartRepo.Reorder(ctx, userID, positions); err != nil {
		h.logger.Error("failed to reorder charts",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to reorder charts")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetQuickActions handles GET /dashboard/quick-actions requests.
func (h *DashboardHandler) HandleGetQuickActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locale := "en"
	if userID, err := h.getUserID(ctx); err == nil && h.userPrefRepo != nil {
		pref, err := h.userPrefRepo.GetByUserID(ctx, userID)
		if err == nil && pref != nil {
			locale = pref.Locale
		}
	}

	// Return localized quick actions
	actions := h.getQuickActions(locale)
	h.writeJSON(w, http.StatusOK, actions)
}

// ============================================================================
// Helper Functions
// ============================================================================

func (h *DashboardHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (h *DashboardHandler) chartToResponse(chart *models.PinnedChart) ChartResponse {
	var queryID *string
	if chart.QueryID != nil {
		id := chart.QueryID.String()
		queryID = &id
	}

	var lastRefreshedAt *string
	if chart.LastRefreshedAt != nil {
		t := chart.LastRefreshedAt.Format("2006-01-02T15:04:05Z07:00")
		lastRefreshedAt = &t
	}

	return ChartResponse{
		ID:                   chart.ID.String(),
		UserID:               chart.UserID.String(),
		Title:                chart.Title,
		QueryID:              queryID,
		NaturalLanguageQuery: chart.NaturalLanguageQuery,
		SQLQuery:             chart.SQLQuery,
		ChartSpec:            chart.ChartSpec,
		ChartType:            chart.ChartType,
		RefreshInterval:      chart.RefreshInterval,
		Locale:               chart.Locale,
		Position: ChartPosition{
			Row:  chart.Position.Row,
			Col:  chart.Position.Col,
			Size: chart.Position.Size,
		},
		LastRefreshedAt: lastRefreshedAt,
		IsActive:        chart.IsActive,
		CreatedAt:       chart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       chart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *DashboardHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *DashboardHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}

func (h *DashboardHandler) getQuickActions(locale string) []QuickActionResponse {
	if locale == "ar" {
		return []QuickActionResponse{
			{ID: "1", Query: "إجمالي الإيرادات هذا الشهر", Label: "إيرادات الشهر"},
			{ID: "2", Query: "أكثر 5 أقسام من حيث الإيرادات", Label: "أفضل الأقسام"},
			{ID: "3", Query: "عدد المرضى الجدد هذا الأسبوع", Label: "مرضى جدد"},
			{ID: "4", Query: "مستوى المخزون المنخفض", Label: "تنبيه المخزون"},
		}
	}

	return []QuickActionResponse{
		{ID: "1", Query: "Total revenue this month", Label: "Monthly Revenue"},
		{ID: "2", Query: "Top 5 departments by revenue", Label: "Top Departments"},
		{ID: "3", Query: "New patients this week", Label: "New Patients"},
		{ID: "4", Query: "Low inventory items", Label: "Inventory Alert"},
	}
}
