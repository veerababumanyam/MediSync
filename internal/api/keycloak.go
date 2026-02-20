// Package api provides Keycloak authentication client for MediSync.
//
// This file implements the KeycloakClient which handles JWT token validation
// and user claims extraction from Keycloak-issued tokens.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/medisync/medisync/internal/api/middleware"
	"github.com/medisync/medisync/internal/config"
)

// KeycloakClient handles authentication with Keycloak.
type KeycloakClient struct {
	config     *config.KeycloakConfig
	logger     *slog.Logger
	httpClient *http.Client

	// Cached JWKS for token verification
	jwks     *JWKS
	jwksMu   sync.RWMutex
	jwksTime time.Time
}

// JWKS represents a JSON Web Key Set.
type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

// JSONWebKey represents a JSON Web Key.
type JSONWebKey struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// Claims represents the JWT claims extracted from a Keycloak token.
type Claims struct {
	// Standard claims
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`

	// Keycloak-specific claims
	RealmAccess    *RealmAccess               `json:"realm_access"`
	ResourceAccess *map[string]ResourceAccess `json:"resource_access"`
	PreferredName  string                     `json:"preferred_username"`
	Email          string                     `json:"email"`
	EmailVerified  bool                       `json:"email_verified"`
	GivenName      string                     `json:"given_name"`
	FamilyName     string                     `json:"family_name"`
	Name           string                     `json:"name"`
	Locale         string                     `json:"locale"`
	TenantID       string                     `json:"tenant_id"`
	SessionID      string                     `json:"sid"`
}

// RealmAccess represents realm-level roles.
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// ResourceAccess represents client-level roles.
type ResourceAccess struct {
	Roles []string `json:"roles"`
}

// NewKeycloakClient creates a new Keycloak authentication client.
func NewKeycloakClient(cfg *config.KeycloakConfig, logger *slog.Logger) (*KeycloakClient, error) {
	if logger == nil {
		logger = slog.Default()
	}

	client := &KeycloakClient{
		config: cfg,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Fetch initial JWKS
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.fetchJWKS(ctx); err != nil {
		logger.Warn("failed to fetch initial JWKS, will retry on first request",
			slog.Any("error", err),
		)
	}

	logger.Info("Keycloak client initialized",
		slog.String("url", cfg.URL),
		slog.String("realm", cfg.Realm),
		slog.String("client_id", cfg.ClientID),
	)

	return client, nil
}

// ValidateToken validates a JWT token and returns the claims.
func (c *KeycloakClient) ValidateToken(ctx context.Context, tokenString string) (*middleware.Claims, error) {
	// Parse token without verification first to get claims
	claims, err := c.parseTokenUnverified(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate issuer
	expectedIssuer := fmt.Sprintf("%s/realms/%s", c.config.URL, c.config.Realm)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", expectedIssuer, claims.Issuer)
	}

	// Validate expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired at %d", claims.ExpiresAt)
	}

	// Validate not before
	if claims.NotBefore > 0 && time.Now().Unix() < claims.NotBefore {
		return nil, fmt.Errorf("token not valid before %d", claims.NotBefore)
	}

	// Map to middleware.Claims
	mClaims := &middleware.Claims{
		Subject:   claims.Subject,
		Email:     claims.Email,
		Name:      claims.Name,
		Locale:    claims.Locale,
		TenantID:  claims.TenantID,
		SessionID: claims.SessionID,
		AuthTime:  claims.IssuedAt,
	}

	if claims.RealmAccess != nil {
		mClaims.Roles = claims.RealmAccess.Roles
	}

	return mClaims, nil
}

// parseTokenUnverified parses a JWT token without verifying the signature.
// In production, this should be replaced with proper signature verification.
func (c *KeycloakClient) parseTokenUnverified(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format: expected 3 parts, got %d", len(parts))
	}

	// Decode the payload (second part)
	payload, err := base64Decode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &claims, nil
}

// fetchJWKS fetches the JSON Web Key Set from Keycloak.
func (c *KeycloakClient) fetchJWKS(ctx context.Context) error {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", c.config.URL, c.config.Realm)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create JWKS request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS request failed with status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	c.jwksMu.Lock()
	c.jwks = &jwks
	c.jwksTime = time.Now()
	c.jwksMu.Unlock()

	c.logger.Debug("fetched JWKS",
		slog.Int("key_count", len(jwks.Keys)),
	)

	return nil
}

// GetPublicKey returns the public key for token verification.
func (c *KeycloakClient) GetPublicKey(kid string) (*JSONWebKey, error) {
	c.jwksMu.RLock()
	defer c.jwksMu.RUnlock()

	if c.jwks == nil {
		return nil, fmt.Errorf("JWKS not available")
	}

	for _, key := range c.jwks.Keys {
		if key.Kid == kid {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

// RefreshJWKS forces a refresh of the JWKS.
func (c *KeycloakClient) RefreshJWKS(ctx context.Context) error {
	return c.fetchJWKS(ctx)
}

// GetUserInfo fetches user information from Keycloak.
func (c *KeycloakClient) GetUserInfo(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", c.config.URL, c.config.Realm)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}

	return userInfo, nil
}

// LogoutURL returns the Keycloak logout URL.
func (c *KeycloakClient) LogoutURL(redirectURI string) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout?redirect_uri=%s",
		c.config.URL, c.config.Realm, redirectURI)
}

// base64Decode decodes a base64url-encoded string.
func base64Decode(encoded string) ([]byte, error) {
	// Add padding if necessary
	switch len(encoded) % 4 {
	case 2:
		encoded += "=="
	case 3:
		encoded += "="
	}

	// Use standard base64 decoding with URL encoding
	// In production, use encoding/base64
	return []byte(encoded), nil // Placeholder - implement proper base64url decoding
}
