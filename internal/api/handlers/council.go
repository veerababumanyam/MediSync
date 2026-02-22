// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the Council of AIs endpoints for multi-agent consensus.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CouncilHandler handles Council of AIs API requests.
type CouncilHandler struct {
	coordinator *council.Coordinator
	repository  council.Repository
	opaClient   OPAClient
}

// OPAClient interface for authorization.
type OPAClient interface {
	Allow(ctx context.Context, action string, user *User, resource interface{}) (bool, error)
	GetListScope(ctx context.Context, user *User) (string, error)
}

// User represents an authenticated user.
type User struct {
	ID            string
	Roles         []string
	Authenticated bool
}

// NewCouncilHandler creates a new Council handler.
func NewCouncilHandler(coordinator *council.Coordinator, repository council.Repository, opaClient OPAClient) *CouncilHandler {
	return &CouncilHandler{
		coordinator: coordinator,
		repository:  repository,
		opaClient:   opaClient,
	}
}

// CreateDeliberation handles POST /v1/council/deliberations.
// Creates a new deliberation and returns the result.
func (h *CouncilHandler) CreateDeliberation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user
	user, ok := ctx.Value("user").(*User)
	if !ok || !user.Authenticated {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request
	var req council.CreateDeliberationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Query == "" {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "invalid_request",
			Message: "Query is required",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set default threshold if not provided
	if req.ConsensusThreshold == 0 {
		req.ConsensusThreshold = council.DefaultConsensusThreshold
	}

	// Check OPA policy
	allowed, err := h.opaClient.Allow(ctx, "create", user, nil)
	if err != nil || !allowed {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "forbidden",
			Message: "Not authorized to create deliberations",
		})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Execute deliberation
	result, err := h.coordinator.Deliberate(ctx, req, user.ID)
	if err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "deliberation_failed",
			Message: err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return response
	render.JSON(w, r, DeliberationResponse{
		ID:              result.Deliberation.ID,
		Query:           result.Deliberation.QueryText,
		Status:          string(result.Deliberation.Status),
		FinalResponse:   result.Deliberation.FinalResponse,
		ConfidenceScore: result.Deliberation.ConfidenceScore,
		ConsensusRecord: result.ConsensusRecord,
		EvidenceTrail:   result.EvidenceTrail,
		AgentResponses:  result.AgentResponses,
		CreatedAt:       result.Deliberation.CreatedAt,
		CompletedAt:     result.Deliberation.CompletedAt,
	})
}

// ListDeliberations handles GET /v1/council/deliberations.
// Lists deliberations with RBAC filtering (admin sees all, user sees own).
func (h *CouncilHandler) ListDeliberations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user
	user, ok := ctx.Value("user").(*User)
	if !ok || !user.Authenticated {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin := h.isAdmin(user)

	// Parse query parameters
	opts := council.ListOptions{
		Status:   council.DeliberationStatus(r.URL.Query().Get("status")),
		Limit:    parseInt(r.URL.Query().Get("limit"), 20),
		Offset:   parseInt(r.URL.Query().Get("offset"), 0),
	}

	if from := r.URL.Query().Get("from"); from != "" {
		opts.FromDate = parseDate(from)
	}
	if to := r.URL.Query().Get("to"); to != "" {
		opts.ToDate = parseDate(to)
	}
	if flagged := r.URL.Query().Get("flagged"); flagged != "" {
		if flagged == "true" {
			t := true
			opts.Flagged = &t
		} else if flagged == "false" {
			f := false
			opts.Flagged = &f
		}
	}

	// Get deliberations
	deliberations, total, err := h.repository.ListDeliberations(ctx, user.ID, isAdmin, opts)
	if err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "query_failed",
			Message: "Failed to list deliberations",
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return response
	render.JSON(w, r, DeliberationListResponse{
		Deliberations: deliberations,
		Total:         total,
		Limit:         opts.Limit,
		Offset:        opts.Offset,
	})
}

// GetDeliberation handles GET /v1/council/deliberations/{id}.
// Returns a specific deliberation with ownership check.
func (h *CouncilHandler) GetDeliberation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user
	user, ok := ctx.Value("user").(*User)
	if !ok || !user.Authenticated {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get deliberation ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "invalid_request",
			Message: "Deliberation ID is required",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get deliberation with all related data
	result, err := h.repository.GetDeliberationWithResponses(ctx, id)
	if err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "not_found",
			Message: "Deliberation not found",
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check ownership (unless admin)
	isAdmin := h.isAdmin(user)
	if !isAdmin && result.Deliberation.UserID != user.ID {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Return response
	render.JSON(w, r, DeliberationResponse{
		ID:              result.Deliberation.ID,
		Query:           result.Deliberation.QueryText,
		Status:          string(result.Deliberation.Status),
		FinalResponse:   result.Deliberation.FinalResponse,
		ConfidenceScore: result.Deliberation.ConfidenceScore,
		ConsensusRecord: result.ConsensusRecord,
		EvidenceTrail:   result.EvidenceTrail,
		AgentResponses:  result.AgentResponses,
		CreatedAt:       result.Deliberation.CreatedAt,
		CompletedAt:     result.Deliberation.CompletedAt,
	})
}

// GetHealth handles GET /v1/council/health.
// Returns the health status of the Council system.
func (h *CouncilHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	summary := h.coordinator.GetHealthSummary()

	render.JSON(w, r, CouncilHealthResponse{
		Status:         summary.OverallStatus,
		TotalAgents:    summary.TotalAgents,
		HealthyAgents:  summary.HealthyAgents,
		DegradedAgents: summary.DegradedAgents,
		FailedAgents:   summary.FailedAgents,
		AgentStatuses:  summary.AgentStatuses,
		LastChecked:    summary.LastChecked,
	})
}

// GetEvidence handles GET /v1/council/deliberations/{id}/evidence.
// Returns the evidence trail for a deliberation.
func (h *CouncilHandler) GetEvidence(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user
	user, ok := ctx.Value("user").(*User)
	if !ok || !user.Authenticated {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	// Get deliberation to check ownership
	delim, err := h.repository.GetDeliberation(ctx, id)
	if err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "not_found",
			Message: "Deliberation not found",
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check ownership
	isAdmin := h.isAdmin(user)
	if !isAdmin && delim.UserID != user.ID {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get evidence trail
	trail, err := h.repository.GetEvidenceTrail(ctx, id)
	if err != nil {
		render.JSON(w, r, CouncilErrorResponse{
			Error:   "not_found",
			Message: "Evidence trail not found",
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, trail)
}

// RegisterRoutes registers all Council routes on the provided router.
func (h *CouncilHandler) RegisterRoutes(r chi.Router) {
	// Deliberation endpoints
	r.Post("/deliberations", h.CreateDeliberation)
	r.Get("/deliberations", h.ListDeliberations)
	r.Get("/deliberations/{id}", h.GetDeliberation)
	r.Get("/deliberations/{id}/evidence", h.GetEvidence)

	// Health endpoint
	r.Get("/health", h.GetHealth)
}

// Helper functions

func (h *CouncilHandler) isAdmin(user *User) bool {
	for _, role := range user.Roles {
		if role == "admin" || role == "superadmin" {
			return true
		}
	}
	return false
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseDate(s string) time.Time {
	// Try RFC3339 format first
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	// Try date-only format
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		return t
	}
	return time.Time{}
}

// Response types

// CouncilErrorResponse represents an error response for Council API.
type CouncilErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// DeliberationResponse represents a deliberation API response.
type DeliberationResponse struct {
	ID              string                    `json:"id"`
	Query           string                    `json:"query"`
	Status          string                    `json:"status"`
	FinalResponse   string                    `json:"final_response,omitempty"`
	ConfidenceScore float64                   `json:"confidence_score,omitempty"`
	ConsensusRecord *council.ConsensusRecord  `json:"consensus_record,omitempty"`
	EvidenceTrail   *council.EvidenceTrail    `json:"evidence_trail,omitempty"`
	AgentResponses  []*council.AgentResponse  `json:"agent_responses,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	CompletedAt     *time.Time                `json:"completed_at,omitempty"`
}

// DeliberationListResponse represents a list of deliberations.
type DeliberationListResponse struct {
	Deliberations []*council.CouncilDeliberation `json:"deliberations"`
	Total         int                            `json:"total"`
	Limit         int                            `json:"limit"`
	Offset        int                            `json:"offset"`
}

// CouncilHealthResponse represents the health status of the Council.
type CouncilHealthResponse struct {
	Status         string            `json:"status"`
	TotalAgents    int               `json:"total_agents"`
	HealthyAgents  int               `json:"healthy_agents"`
	DegradedAgents int               `json:"degraded_agents"`
	FailedAgents   int               `json:"failed_agents"`
	AgentStatuses  map[string]string `json:"agent_statuses"`
	LastChecked    time.Time         `json:"last_checked"`
}

// Context key type for type safety
type contextKey string

const userContextKey contextKey = "user"
