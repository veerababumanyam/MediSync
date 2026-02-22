// Package handler demonstrates HTTP handler patterns in Go.
// Uses go-chi/chi for routing and standard library for HTTP.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// =============================================================================
// TYPES
// =============================================================================

// UserResponse is the JSON response for a user.
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

// CreateUserRequest is the JSON request for creating a user.
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserRequest is the JSON request for updating a user.
type UpdateUserRequest struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// UserService defines the interface for user operations.
type UserService interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// User represents a user entity (import from your domain package).
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
}

// =============================================================================
// HANDLER
// =============================================================================

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	service UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes registers user routes on the provided router.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

// Create handles POST /users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON body")
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user))
}

// Get handles GET /users/{id}.
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Missing user ID")
		return
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

// Update handles PUT /users/{id}.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Missing user ID")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON body")
		return
	}

	user, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

// Delete handles DELETE /users/{id}.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Missing user ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List handles GET /users.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters with defaults
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = toUserResponse(user)
	}

	writeJSON(w, http.StatusOK, response)
}

// =============================================================================
// HELPERS
// =============================================================================

// writeJSON writes JSON response with proper headers.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// handleServiceError maps service errors to HTTP responses.
func handleServiceError(w http.ResponseWriter, err error) {
	// Import your errors package to check error types
	// This is a simplified example

	switch {
	case isErrorType(err, "not found"):
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
	case isErrorType(err, "validation"):
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	case isErrorType(err, "already exists"):
		writeError(w, http.StatusConflict, "CONFLICT", err.Error())
	case isErrorType(err, "unauthorized"):
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	default:
		// Log internal error but don't expose details
		slog.Error("internal server error", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

// isErrorType checks if an error message contains a substring.
// In practice, use errors.Is() or errors.As() with typed errors.
func isErrorType(err error, substr string) bool {
	return strings.Contains(strings.ToLower(err.Error()), substr)
}

// toUserResponse converts a domain User to JSON response.
func toUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

// =============================================================================
// ROUTER SETUP
// =============================================================================

// NewRouter creates and configures the main router.
func NewRouter(userHandler *UserHandler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
		// Add more handlers here
	})

	return r
}

// =============================================================================
// MIDDLEWARE EXAMPLES
// =============================================================================

// AuthMiddleware validates authentication tokens.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization token")
			return
		}

		// Validate token and extract user info
		userID, err := validateToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware adds CORS headers.
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
