// Package middleware provides HTTP middleware for the MediSync API.
//
// This file implements the AuthMiddleware that validates JWT tokens via Keycloak
// and extracts user claims into the request context.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// contextKey is a type for context keys.
type contextKey string

const (
	// ClaimsKey is the context key for JWT claims.
	ClaimsKey contextKey = "claims"
	// UserIDKey is the context key for user ID.
	UserIDKey contextKey = "user_id"
	// TenantIDKey is the context key for tenant ID.
	TenantIDKey contextKey = "tenant_id"
	// RolesKey is the context key for user roles.
	RolesKey contextKey = "roles"
)

// Claims represents JWT claims stored in context.
type Claims struct {
	Subject    string
	Email      string
	Name       string
	Locale     string
	TenantID   string
	SessionID  string
	Roles      []string
	AuthTime   int64
}

// KeycloakValidator defines the interface for Keycloak token validation.
type KeycloakValidator interface {
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
}

// AuthMiddleware validates JWT tokens via Keycloak and adds claims to context.
// It returns 401 Unauthorized for invalid or expired tokens.
func AuthMiddleware(validator KeycloakValidator, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health and ready endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Debug("missing authorization header",
					slog.String("path", r.URL.Path),
				)
				writeUnauthorized(w, "missing authorization header")
				return
			}

			// Parse Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				logger.Debug("invalid authorization header format",
					slog.String("path", r.URL.Path),
				)
				writeUnauthorized(w, "invalid authorization header format")
				return
			}

			tokenString := parts[1]
			if tokenString == "" {
				logger.Debug("empty bearer token",
					slog.String("path", r.URL.Path),
				)
				writeUnauthorized(w, "empty bearer token")
				return
			}

			// Validate token with Keycloak
			claims, err := validator.ValidateToken(r.Context(), tokenString)
			if err != nil {
				logger.Warn("token validation failed",
					slog.String("path", r.URL.Path),
					slog.Any("error", err),
				)
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			// Add claims to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			ctx = context.WithValue(ctx, UserIDKey, claims.Subject)
			ctx = context.WithValue(ctx, TenantIDKey, claims.TenantID)
			ctx = context.WithValue(ctx, RolesKey, claims.Roles)

			// Log successful authentication
			logger.Debug("user authenticated",
				slog.String("user_id", claims.Subject),
				slog.String("email", claims.Email),
				slog.Any("roles", claims.Roles),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves claims from the request context.
func GetClaims(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(ClaimsKey).(*Claims); ok {
		return claims
	}
	return nil
}

// GetUserID retrieves the user ID from the request context.
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetTenantID retrieves the tenant ID from the request context.
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// GetRoles retrieves user roles from the request context.
func GetRoles(ctx context.Context) []string {
	if roles, ok := ctx.Value(RolesKey).([]string); ok {
		return roles
	}
	return nil
}

// HasRole checks if the user has a specific role.
func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin checks if the user has the admin role.
func IsAdmin(ctx context.Context) bool {
	return HasRole(ctx, "admin")
}

// RequireRole is a middleware that requires a specific role.
func RequireRole(role string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasRole(r.Context(), role) {
				logger.Warn("role required but not present",
					slog.String("required_role", role),
					slog.String("user_id", GetUserID(r.Context())),
				)
				writeForbidden(w, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// writeUnauthorized writes a 401 Unauthorized response.
func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer realm="medisync"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":{"code":"unauthorized","message":"` + message + `"}}`))
}

// writeForbidden writes a 403 Forbidden response.
func writeForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error":{"code":"forbidden","message":"` + message + `"}}`))
}
