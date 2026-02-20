// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the user preferences management endpoints for
// locale, numeral system, calendar system, report language, and timezone settings.
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

// PreferencesHandler handles user preferences management endpoints.
type PreferencesHandler struct {
	logger       *slog.Logger
	userPrefRepo *warehouse.UserPreferenceRepository
}

// PreferencesHandlerConfig holds configuration for the PreferencesHandler.
type PreferencesHandlerConfig struct {
	Logger       *slog.Logger
	UserPrefRepo *warehouse.UserPreferenceRepository
}

// NewPreferencesHandler creates a new PreferencesHandler instance.
func NewPreferencesHandler(cfg PreferencesHandlerConfig) *PreferencesHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &PreferencesHandler{
		logger:       cfg.Logger,
		userPrefRepo: cfg.UserPrefRepo,
	}
}

// RegisterRoutes registers preferences routes on the given router.
func (h *PreferencesHandler) RegisterRoutes(r chi.Router) {
	r.Route("/preferences", func(r chi.Router) {
		r.Get("/", h.HandleGetPreferences)
		r.Patch("/", h.HandleUpdatePreferences)
	})
}

// ============================================================================
// Request/Response Types
// ============================================================================

// PreferencesResponse represents user preferences in API responses.
type PreferencesResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	Locale         string `json:"locale"`
	NumeralSystem  string `json:"numeralSystem"`
	CalendarSystem string `json:"calendarSystem"`
	ReportLanguage string `json:"reportLanguage"`
	Timezone       string `json:"timezone"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// UpdatePreferencesRequest represents a request to update user preferences.
type UpdatePreferencesRequest struct {
	Locale         *string `json:"locale,omitempty"`
	NumeralSystem  *string `json:"numeralSystem,omitempty"`
	CalendarSystem *string `json:"calendarSystem,omitempty"`
	ReportLanguage *string `json:"reportLanguage,omitempty"`
	Timezone       *string `json:"timezone,omitempty"`
}

// Valid locale values
var validLocales = map[string]bool{
	"en": true,
	"ar": true,
}

// Valid numeral system values
var validNumeralSystems = map[string]bool{
	"western":    true,
	"arabic-indic": true,
}

// Valid calendar system values
var validCalendarSystems = map[string]bool{
	"gregorian": true,
	"islamic":   true,
}

// Default preferences for new users
var defaultPreferences = struct {
	Locale         string
	NumeralSystem  string
	CalendarSystem string
	ReportLanguage string
}{
	Locale:         "en",
	NumeralSystem:  "western",
	CalendarSystem: "gregorian",
	ReportLanguage: "en",
}

// ============================================================================
// HTTP Handlers
// ============================================================================

// HandleGetPreferences handles GET /preferences requests.
// Returns the current user's preferences, creating defaults if not exists.
func (h *PreferencesHandler) HandleGetPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	pref, err := h.userPrefRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Preferences don't exist yet - create defaults
		pref = &models.UserPreference{
			UserID:         userID,
			Locale:         defaultPreferences.Locale,
			NumeralSystem:  defaultPreferences.NumeralSystem,
			CalendarSystem: defaultPreferences.CalendarSystem,
			ReportLanguage: defaultPreferences.ReportLanguage,
			Timezone:       "UTC",
		}

		if err := h.userPrefRepo.Upsert(ctx, pref); err != nil {
			h.logger.Error("failed to create default preferences",
				slog.Any("error", err),
				slog.String("user_id", userID.String()),
			)
			h.writeError(w, http.StatusInternalServerError, "failed to create preferences")
			return
		}

		h.logger.Info("created default preferences for user",
			slog.String("user_id", userID.String()),
		)
	}

	h.writeJSON(w, http.StatusOK, h.preferenceToResponse(pref))
}

// HandleUpdatePreferences handles PATCH /preferences requests.
// Updates the provided fields and returns the updated preferences.
func (h *PreferencesHandler) HandleUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request fields
	if err := h.validateUpdateRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing preferences or create defaults
	pref, err := h.userPrefRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Create new preferences with defaults, then apply updates
		pref = &models.UserPreference{
			UserID:         userID,
			Locale:         defaultPreferences.Locale,
			NumeralSystem:  defaultPreferences.NumeralSystem,
			CalendarSystem: defaultPreferences.CalendarSystem,
			ReportLanguage: defaultPreferences.ReportLanguage,
			Timezone:       "UTC",
		}
	}

	// Apply updates
	if req.Locale != nil {
		pref.Locale = *req.Locale
	}
	if req.NumeralSystem != nil {
		pref.NumeralSystem = *req.NumeralSystem
	}
	if req.CalendarSystem != nil {
		pref.CalendarSystem = *req.CalendarSystem
	}
	if req.ReportLanguage != nil {
		pref.ReportLanguage = *req.ReportLanguage
	}
	if req.Timezone != nil {
		pref.Timezone = *req.Timezone
	}

	// Save to database
	if err := h.userPrefRepo.Upsert(ctx, pref); err != nil {
		h.logger.Error("failed to update preferences",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to update preferences")
		return
	}

	h.logger.Info("updated user preferences",
		slog.String("user_id", userID.String()),
		slog.String("locale", pref.Locale),
	)

	h.writeJSON(w, http.StatusOK, h.preferenceToResponse(pref))
}

// ============================================================================
// Helper Functions
// ============================================================================

func (h *PreferencesHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (h *PreferencesHandler) preferenceToResponse(pref *models.UserPreference) PreferencesResponse {
	return PreferencesResponse{
		ID:             pref.ID.String(),
		UserID:         pref.UserID.String(),
		Locale:         pref.Locale,
		NumeralSystem:  pref.NumeralSystem,
		CalendarSystem: pref.CalendarSystem,
		ReportLanguage: pref.ReportLanguage,
		Timezone:       pref.Timezone,
		CreatedAt:      pref.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      pref.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *PreferencesHandler) validateUpdateRequest(req *UpdatePreferencesRequest) error {
	if req.Locale != nil && !validLocales[*req.Locale] {
		return errors.New("invalid locale: must be 'en' or 'ar'")
	}
	if req.NumeralSystem != nil && !validNumeralSystems[*req.NumeralSystem] {
		return errors.New("invalid numeralSystem: must be 'western' or 'arabic-indic'")
	}
	if req.CalendarSystem != nil && !validCalendarSystems[*req.CalendarSystem] {
		return errors.New("invalid calendarSystem: must be 'gregorian' or 'islamic'")
	}
	if req.ReportLanguage != nil && !validLocales[*req.ReportLanguage] {
		return errors.New("invalid reportLanguage: must be 'en' or 'ar'")
	}
	return nil
}

func (h *PreferencesHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *PreferencesHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}
