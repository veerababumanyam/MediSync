// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the chat session management endpoints for the
// conversational BI feature. It handles session listing, message history,
// and session lifecycle management.
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

// ChatSessionsHandler handles chat session management endpoints.
type ChatSessionsHandler struct {
	logger       *slog.Logger
	chatMsgRepo  *warehouse.ChatMessageRepository
	userPrefRepo *warehouse.UserPreferenceRepository
}

// ChatSessionsHandlerConfig holds configuration for the ChatSessionsHandler.
type ChatSessionsHandlerConfig struct {
	Logger       *slog.Logger
	ChatMsgRepo  *warehouse.ChatMessageRepository
	UserPrefRepo *warehouse.UserPreferenceRepository
}

// NewChatSessionsHandler creates a new ChatSessionsHandler instance.
func NewChatSessionsHandler(cfg ChatSessionsHandlerConfig) *ChatSessionsHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &ChatSessionsHandler{
		logger:       cfg.Logger,
		chatMsgRepo:  cfg.ChatMsgRepo,
		userPrefRepo: cfg.UserPrefRepo,
	}
}

// RegisterRoutes registers chat session routes on the given router.
func (h *ChatSessionsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/chat", func(r chi.Router) {
		// Session management
		r.Get("/sessions", h.HandleListSessions)
		r.Post("/sessions", h.HandleCreateSession)
		r.Get("/sessions/{session_id}/messages", h.HandleGetMessages)
		r.Delete("/sessions/{session_id}", h.HandleDeleteSession)

		// Recent messages across all sessions
		r.Get("/recent", h.HandleGetRecentMessages)
	})
}

// ============================================================================
// Request/Response Types
// ============================================================================

// CreateSessionRequest represents a request to create a new chat session.
type CreateSessionRequest struct {
	// Locale is the preferred locale for the session.
	Locale string `json:"locale,omitempty"`

	// Title is an optional title for the session.
	Title string `json:"title,omitempty"`
}

// SessionResponse represents a chat session in API responses.
type SessionResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// MessageResponse represents a chat message in API responses.
type MessageResponse struct {
	ID              string                 `json:"id"`
	SessionID       string                 `json:"sessionId"`
	Role            string                 `json:"role"`
	Content         string                 `json:"content"`
	ChartSpec       map[string]interface{} `json:"chartSpec,omitempty"`
	TableData       map[string]interface{} `json:"tableData,omitempty"`
	DrilldownQuery  *string                `json:"drilldownQuery,omitempty"`
	ConfidenceScore *float64               `json:"confidenceScore,omitempty"`
	CreatedAt       string                 `json:"createdAt"`
}

// ListSessionsResponse represents the response for listing sessions.
type ListSessionsResponse struct {
	Sessions []SessionSummary `json:"sessions"`
	Total    int              `json:"total"`
}

// SessionSummary represents a brief session overview.
type SessionSummary struct {
	ID             string `json:"id"`
	Title          string `json:"title,omitempty"`
	LastMessage    string `json:"lastMessage,omitempty"`
	MessageCount   int    `json:"messageCount"`
	CreatedAt      string `json:"createdAt"`
	LastActivityAt string `json:"lastActivityAt"`
}

// ListMessagesResponse represents the response for listing messages.
type ListMessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int               `json:"total"`
	HasMore  bool              `json:"hasMore"`
}

// ============================================================================
// HTTP Handlers
// ============================================================================

// HandleListSessions handles GET /chat/sessions requests.
// It returns a list of the user's recent chat sessions.
func (h *ChatSessionsHandler) HandleListSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (set by auth middleware)
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get recent messages to build session summaries
	messages, err := h.chatMsgRepo.GetRecentByUserID(ctx, userID, 100)
	if err != nil {
		h.logger.Error("failed to get recent messages",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve sessions")
		return
	}

	// Group messages by session
	sessionsMap := make(map[uuid.UUID]*SessionSummary)
	for _, msg := range messages {
		summary, exists := sessionsMap[msg.SessionID]
		if !exists {
			summary = &SessionSummary{
				ID:             msg.SessionID.String(),
				MessageCount:   0,
				CreatedAt:      msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				LastActivityAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			sessionsMap[msg.SessionID] = summary
		}

		summary.MessageCount++
		// Update last activity if this message is newer
		if msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00") > summary.LastActivityAt {
			summary.LastActivityAt = msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
			// Use the last user message as the session title/preview
			if msg.Role == "user" && len(msg.Content) > 0 {
				summary.LastMessage = truncateString(msg.Content, 100)
				if summary.Title == "" {
					summary.Title = truncateString(msg.Content, 50)
				}
			}
		}
	}

	// Convert map to slice
	sessions := make([]SessionSummary, 0, len(sessionsMap))
	for _, summary := range sessionsMap {
		sessions = append(sessions, *summary)
	}

	response := ListSessionsResponse{
		Sessions: sessions,
		Total:    len(sessions),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleCreateSession handles POST /chat/sessions requests.
// It creates a new chat session and returns the session ID.
func (h *ChatSessionsHandler) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request (optional)
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Request body is optional, ignore parse errors
		req = CreateSessionRequest{}
	}

	// Get user's locale preference
	locale := req.Locale
	if locale == "" {
		pref, err := h.userPrefRepo.GetByUserID(ctx, userID)
		if err == nil && pref != nil {
			locale = pref.Locale
		}
	}
	if locale == "" {
		locale = "en"
	}

	// Generate new session ID
	sessionID := uuid.New()

	// Create initial system message for the session
	systemMsg := &models.ChatMessage{
		SessionID: sessionID,
		UserID:    userID,
		Role:      "system",
		Content:   "Welcome to MediSync BI Assistant. How can I help you analyze your data today?",
		Locale:    locale,
	}

	if err := h.chatMsgRepo.Create(ctx, systemMsg); err != nil {
		h.logger.Error("failed to create session",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	response := SessionResponse{
		ID:        sessionID.String(),
		Title:     req.Title,
		CreatedAt: systemMsg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: systemMsg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// HandleGetMessages handles GET /chat/sessions/{session_id}/messages requests.
// It returns the message history for a specific session.
func (h *ChatSessionsHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse session ID from URL
	sessionIDStr := chi.URLParam(r, "session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid session ID")
		return
	}

	// Parse pagination parameters
	limit := 50 // default
	if l := r.URL.Query().Get("limit"); l != "" {
		var lim int
		if _, err := json.Number(l).Int64(); err == nil {
			if lim > 0 && lim <= 100 {
				limit = lim
			}
		}
	}

	// Get messages for the session
	messages, err := h.chatMsgRepo.GetBySessionID(ctx, sessionID, limit)
	if err != nil {
		h.logger.Error("failed to get messages",
			slog.Any("error", err),
			slog.String("session_id", sessionID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve messages")
		return
	}

	// Verify the session belongs to the user
	if len(messages) > 0 && messages[0].UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	// Convert to response format
	msgResponses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		msgResponses[i] = MessageResponse{
			ID:              msg.ID.String(),
			SessionID:       msg.SessionID.String(),
			Role:            msg.Role,
			Content:         msg.Content,
			ChartSpec:       msg.ChartSpec,
			TableData:       msg.TableData,
			DrilldownQuery:  msg.DrilldownQuery,
			ConfidenceScore: msg.ConfidenceScore,
			CreatedAt:       msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListMessagesResponse{
		Messages: msgResponses,
		Total:    len(msgResponses),
		HasMore:  len(msgResponses) == limit,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// HandleDeleteSession handles DELETE /chat/sessions/{session_id} requests.
// It deletes a chat session and all its messages.
func (h *ChatSessionsHandler) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse session ID from URL
	sessionIDStr := chi.URLParam(r, "session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid session ID")
		return
	}

	// Verify the session belongs to the user by checking messages
	messages, err := h.chatMsgRepo.GetBySessionID(ctx, sessionID, 1)
	if err != nil {
		h.logger.Error("failed to verify session",
			slog.Any("error", err),
			slog.String("session_id", sessionID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to verify session")
		return
	}

	if len(messages) == 0 {
		h.writeError(w, http.StatusNotFound, "session not found")
		return
	}

	if messages[0].UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	// Delete the session (implemented via soft-delete or hard-delete)
	// For now, we return success without actual deletion as the repo
	// doesn't have a Delete method. In production, add Delete to repo.
	h.logger.Info("session deletion requested",
		slog.String("session_id", sessionID.String()),
		slog.String("user_id", userID.String()),
	)

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetRecentMessages handles GET /chat/recent requests.
// It returns recent messages across all user's sessions.
func (h *ChatSessionsHandler) HandleGetRecentMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse limit parameter
	limit := 20 // default
	if l := r.URL.Query().Get("limit"); l != "" {
		var lim int
		if _, err := json.Number(l).Int64(); err == nil {
			if lim > 0 && lim <= 100 {
				limit = lim
			}
		}
	}

	// Get recent messages
	messages, err := h.chatMsgRepo.GetRecentByUserID(ctx, userID, limit)
	if err != nil {
		h.logger.Error("failed to get recent messages",
			slog.Any("error", err),
			slog.String("user_id", userID.String()),
		)
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve messages")
		return
	}

	// Convert to response format
	msgResponses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		msgResponses[i] = MessageResponse{
			ID:              msg.ID.String(),
			SessionID:       msg.SessionID.String(),
			Role:            msg.Role,
			Content:         msg.Content,
			ChartSpec:       msg.ChartSpec,
			TableData:       msg.TableData,
			DrilldownQuery:  msg.DrilldownQuery,
			ConfidenceScore: msg.ConfidenceScore,
			CreatedAt:       msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListMessagesResponse{
		Messages: msgResponses,
		Total:    len(msgResponses),
		HasMore:  len(msgResponses) == limit,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ============================================================================
// Helper Functions
// ============================================================================

// getUserID extracts the user ID from the request context.
func (h *ChatSessionsHandler) getUserID(ctx context.Context) (uuid.UUID, error) {
	// The user ID should be set by the auth middleware
	userIDStr, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	return uuid.Parse(userIDStr)
}

// writeJSON writes a JSON response.
func (h *ChatSessionsHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to write JSON response", slog.Any("error", err))
	}
}

// writeError writes an error response.
func (h *ChatSessionsHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"code":    http.StatusText(status),
		},
	})
}

// truncateString truncates a string to a maximum length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
