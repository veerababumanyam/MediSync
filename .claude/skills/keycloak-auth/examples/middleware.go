// Package auth provides authentication and authorization middleware
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/oauth2"
)

// =============================================================================
// TYPES
// =============================================================================

// Config holds authentication configuration
type Config struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	Realm        string
	RedirectURL  string
}

// Claims represents the JWT token claims
type Claims struct {
	Subject           string   `json:"sub"`
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	PreferredUsername string   `json:"preferred_username"`
	Name              string   `json:"name"`
	Locale            string   `json:"locale"`
	CompanyID         string   `json:"company_id"`
	RealmAccess       RealmAccess `json:"realm_access"`
	ResourceAccess    map[string]ResourceAccess `json:"resource_access"`
	Permissions       []string `json:"permissions"`
	Exp               int64    `json:"exp"`
	Iat               int64    `json:"iat"`
}

// RealmAccess contains realm-level roles
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// ResourceAccess contains client-level roles
type ResourceAccess struct {
	Roles []string `json:"roles"`
}

// User represents an authenticated user
type User struct {
	ID           string
	Email        string
	Name         string
	Locale       string
	CompanyID    string
	Roles        []string
	Permissions  []string
}

// Authenticator handles authentication
type Authenticator struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth    *oauth2.Config
	config   Config
}

// =============================================================================
// CONSTRUCTOR
// =============================================================================

// NewAuthenticator creates a new authenticator
func NewAuthenticator(ctx context.Context, cfg Config) (*Authenticator, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	oauth := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		provider: provider,
		verifier: verifier,
		oauth:    oauth,
		config:   cfg,
	}, nil
}

// =============================================================================
// MIDDLEWARE
// =============================================================================

// contextKey is the type for context keys
type contextKey string

const userKey contextKey = "user"

// Middleware validates JWT tokens and adds user to context
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			writeJSONError(w, http.StatusUnauthorized, "INVALID_FORMAT", "Invalid authorization format")
			return
		}

		token, err := a.verifier.Verify(r.Context(), tokenString)
		if err != nil {
			slog.Warn("token verification failed", "error", err)
			writeJSONError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token verification failed")
			return
		}

		var claims Claims
		if err := token.Claims(&claims); err != nil {
			writeJSONError(w, http.StatusUnauthorized, "INVALID_CLAIMS", "Invalid token claims")
			return
		}

		// Convert claims to user
		user := &User{
			ID:          claims.Subject,
			Email:       claims.Email,
			Name:        claims.Name,
			Locale:      claims.Locale,
			CompanyID:   claims.CompanyID,
			Roles:       claims.RealmAccess.Roles,
			Permissions: claims.Permissions,
		}

		// Set default locale
		if user.Locale == "" {
			user.Locale = "en"
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userKey, user)

		// Log authenticated request
		slog.Debug("request authenticated",
			"user_id", user.ID,
			"company_id", user.CompanyID,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts user from context
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}

// MustGetUser extracts user from context or panics
func MustGetUser(ctx context.Context) *User {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		panic("user not found in context")
	}
	return user
}

// =============================================================================
// ROLE CHECKING
// =============================================================================

// RequireRoles creates middleware that checks for required roles
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			if !hasAnyRole(user.Roles, roles) {
				slog.Warn("role check failed",
					"user_id", user.ID,
					"required_roles", roles,
					"user_roles", user.Roles,
				)
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermissions creates middleware that checks for required permissions
func RequirePermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			if !hasAnyPermission(user.Permissions, permissions) {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireCompanyAccess creates middleware that checks company access
func RequireCompanyAccess() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			// Get company ID from URL or request body
			requestCompanyID := chi.URLParam(r, "company_id")
			if requestCompanyID == "" {
				requestCompanyID = r.Header.Get("X-Company-ID")
			}

			// If no company specified, use user's company
			if requestCompanyID == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user has access to this company
			if user.CompanyID != requestCompanyID && !hasRole(user.Roles, "admin") {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "Company access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func hasAnyRole(userRoles, requiredRoles []string) bool {
	for _, req := range requiredRoles {
		if hasRole(userRoles, req) {
			return true
		}
	}
	return false
}

func hasAnyPermission(userPerms, requiredPerms []string) bool {
	permSet := make(map[string]bool)
	for _, p := range userPerms {
		permSet[p] = true
	}
	for _, req := range requiredPerms {
		if permSet[req] {
			return true
		}
	}
	return false
}

// =============================================================================
// RESPONSE HELPERS
// =============================================================================

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(status),
		Code:    code,
		Message: message,
	})
}

// =============================================================================
// OAUTH FLOWS
// =============================================================================

// LoginHandler initiates the OAuth2 login flow
func (a *Authenticator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	setStateCookie(w, state)

	url := a.oauth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// CallbackHandler handles the OAuth2 callback
func (a *Authenticator) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != getStateFromCookie(r) {
		writeJSONError(w, http.StatusBadRequest, "INVALID_STATE", "Invalid state parameter")
		return
	}

	code := r.URL.Query().Get("code")
	token, err := a.oauth.Exchange(r.Context(), code)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "EXCHANGE_FAILED", "Token exchange failed")
		return
	}

	// Return tokens to client
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expires_in":    token.Expiry.Sub(time.Now()).Seconds(),
		"token_type":    token.TokenType,
	})
}

// RefreshHandler refreshes an access token
func (a *Authenticator) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	token := &oauth2.Token{
		RefreshToken: req.RefreshToken,
	}

	newToken, err := a.oauth.TokenSource(r.Context(), token).Token()
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "REFRESH_FAILED", "Token refresh failed")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  newToken.AccessToken,
		"refresh_token": newToken.RefreshToken,
		"expires_in":    newToken.Expiry.Sub(time.Now()).Seconds(),
		"token_type":    newToken.TokenType,
	})
}

// =============================================================================
// STATE MANAGEMENT
// =============================================================================

func generateState() string {
	return uuid.New().String()
}

func setStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getStateFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return ""
	}
	return cookie.Value
}
