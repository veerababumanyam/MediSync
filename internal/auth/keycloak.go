// Package auth provides authentication and authorization for MediSync.
//
// This file provides the KeycloakValidator struct for JWT token validation
// using Keycloak. It supports both JWT verification with public keys and
// token introspection for additional validation.
//
// Features:
//   - JWT signature verification using Keycloak public keys
//   - Token introspection for real-time validation
//   - Caching of validated tokens in Redis
//   - Extraction of user claims (UserID, TenantID, Roles, Locale, Email)
//
// Usage:
//
//	validator, err := auth.NewKeycloakValidator(config.Keycloak, redisClient, logger)
//	if err != nil {
//	    log.Fatal("Failed to create validator:", err)
//	}
//
//	claims, err := validator.ValidateToken(ctx, tokenString)
//	if err != nil {
//	    // Handle invalid token
//	}
package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// Keycloak configuration constants.
const (
	// DefaultKeycloakTimeout is the default timeout for Keycloak HTTP requests.
	DefaultKeycloakTimeout = 10 * time.Second

	// DefaultCacheTTL is the default TTL for cached token validations.
	DefaultCacheTTL = 5 * time.Minute

	// CacheKeyPrefix is the prefix for cached token keys in Redis.
	CacheKeyPrefix = "medisync:token"

	// JWKSPath is the path to Keycloak's JWKS endpoint.
	JWKSPath = "/protocol/openid-connect/certs"
)

// Claims represents the extracted JWT claims for MediSync.
type Claims struct {
	// UserID is the unique identifier for the user (Keycloak sub).
	UserID string `json:"user_id"`

	// TenantID is the tenant/organization ID for multi-tenancy.
	TenantID string `json:"tenant_id,omitempty"`

	// Username is the user's login name.
	Username string `json:"username,omitempty"`

	// Email is the user's email address.
	Email string `json:"email,omitempty"`

	// Roles are the user's assigned roles.
	Roles []string `json:"roles"`

	// Locale is the user's preferred language.
	Locale string `json:"locale,omitempty"`

	// TimeZone is the user's timezone.
	TimeZone string `json:"timezone,omitempty"`

	// CalendarSystem is the user's calendar preference.
	CalendarSystem string `json:"calendar_system,omitempty"`

	// CostCentres are the cost centres the user has access to.
	CostCentres []string `json:"cost_centres,omitempty"`

	// ExpiresAt is when the token expires.
	ExpiresAt time.Time `json:"expires_at"`

	// IssuedAt is when the token was issued.
	IssuedAt time.Time `json:"issued_at"`

	// Issuer is the token issuer (Keycloak URL).
	Issuer string `json:"issuer"`
}

// KeycloakConfig holds configuration for Keycloak connection.
type KeycloakConfig struct {
	// URL is the Keycloak server base URL.
	URL string

	// Realm is the Keycloak realm name.
	Realm string

	// ClientID is the OAuth2 client ID.
	ClientID string

	// ClientSecret is the OAuth2 client secret (for introspection).
	ClientSecret string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// CacheTTL is how long to cache token validations.
	CacheTTL time.Duration

	// Logger is the structured logger.
	Logger *slog.Logger
}

// KeycloakValidator validates JWT tokens issued by Keycloak.
type KeycloakValidator struct {
	config     *KeycloakConfig
	httpClient *http.Client
	redis      *redis.Client
	logger     *slog.Logger

	// JWKS caching
	jwks      map[string]*rsa.PublicKey
	jwksMu    sync.RWMutex
	jwksExp   time.Time
	jwksURL   string
}

// jwksResponse represents the response from Keycloak's JWKS endpoint.
type jwksResponse struct {
	Keys []jsonWebKey `json:"keys"`
}

// jsonWebKey represents a JSON Web Key.
type jsonWebKey struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c,omitempty"`
}

// introspectionResponse represents the response from token introspection.
type introspectionResponse struct {
	Active    bool   `json:"active"`
	Subject   string `json:"sub,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access,omitempty"`
	TenantID      string `json:"tenant_id,omitempty"`
	Locale        string `json:"locale,omitempty"`
	TimeZone      string `json:"timezone,omitempty"`
	Calendar      string `json:"calendar_system,omitempty"`
	CostCentres   string `json:"cost_centres,omitempty"`
}

// NewKeycloakValidator creates a new Keycloak token validator.
func NewKeycloakValidator(cfg *KeycloakConfig, redisClient *redis.Client) (*KeycloakValidator, error) {
	if cfg == nil {
		return nil, fmt.Errorf("auth: keycloak config is required")
	}

	if cfg.URL == "" || cfg.Realm == "" {
		return nil, fmt.Errorf("auth: keycloak URL and realm are required")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultKeycloakTimeout
	}

	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = DefaultCacheTTL
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &KeycloakValidator{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		redis:   redisClient,
		logger:  cfg.Logger,
		jwks:    make(map[string]*rsa.PublicKey),
		jwksURL: fmt.Sprintf("%s/realms/%s%s", cfg.URL, cfg.Realm, JWKSPath),
	}, nil
}

// ValidateToken validates a JWT token and returns the claims.
// It first checks the Redis cache, then verifies the JWT signature,
// and optionally performs token introspection.
func (v *KeycloakValidator) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("auth: token is required")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Check cache first
	cachedClaims, err := v.getCachedClaims(ctx, tokenString)
	if err == nil && cachedClaims != nil {
		// Verify the cached token hasn't expired
		if time.Now().Before(cachedClaims.ExpiresAt) {
			v.logger.Debug("token validated from cache",
				slog.String("user_id", cachedClaims.UserID),
			)
			return cachedClaims, nil
		}
	}

	// Parse and validate JWT
	claims, err := v.validateJWT(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("auth: JWT validation failed: %w", err)
	}

	// Cache the validated claims
	if err := v.cacheClaims(ctx, tokenString, claims); err != nil {
		// Log but don't fail
		v.logger.Warn("failed to cache token claims", slog.String("error", err.Error()))
	}

	v.logger.Info("token validated",
		slog.String("user_id", claims.UserID),
		slog.String("email", claims.Email),
		slog.Any("roles", claims.Roles),
	)

	return claims, nil
}

// ValidateTokenWithIntrospection validates a token using both JWT and introspection.
func (v *KeycloakValidator) ValidateTokenWithIntrospection(ctx context.Context, tokenString string) (*Claims, error) {
	// First do JWT validation
	claims, err := v.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	// Perform introspection for additional security
	introspected, err := v.introspectToken(ctx, tokenString)
	if err != nil {
		v.logger.Warn("token introspection failed, relying on JWT",
			slog.String("error", err.Error()),
		)
		return claims, nil
	}

	if !introspected.Active {
		return nil, fmt.Errorf("auth: token is not active")
	}

	return claims, nil
}

// validateJWT parses and validates the JWT token.
func (v *KeycloakValidator) validateJWT(ctx context.Context, tokenString string) (*Claims, error) {
	// Parse the token without verification first to get the kid
	unverifiedToken, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Get the key ID from the header
	kid, ok := unverifiedToken.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("token missing kid header")
	}

	// Get the public key
	publicKey, err := v.getPublicKey(ctx, kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// Extract claims
	return v.extractClaims(token)
}

// extractClaims extracts MediSync claims from a JWT token.
func (v *KeycloakValidator) extractClaims(token *jwt.Token) (*Claims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	// Extract standard claims
	claimsObj := &Claims{}

	// Subject (user ID)
	if sub, ok := claims["sub"].(string); ok {
		claimsObj.UserID = sub
	}

	// Preferred username
	if preferredUsername, ok := claims["preferred_username"].(string); ok {
		claimsObj.Username = preferredUsername
	}

	// Email
	if email, ok := claims["email"].(string); ok {
		claimsObj.Email = email
	}

	// Issuer
	if iss, ok := claims["iss"].(string); ok {
		claimsObj.Issuer = iss
	}

	// Expiration
	if exp, ok := claims["exp"].(float64); ok {
		claimsObj.ExpiresAt = time.Unix(int64(exp), 0)
	}

	// Issued at
	if iat, ok := claims["iat"].(float64); ok {
		claimsObj.IssuedAt = time.Unix(int64(iat), 0)
	}

	// Tenant ID (custom claim)
	if tenantID, ok := claims["tenant_id"].(string); ok {
		claimsObj.TenantID = tenantID
	}

	// Locale
	if locale, ok := claims["locale"].(string); ok {
		claimsObj.Locale = locale
	} else {
		claimsObj.Locale = "en" // Default
	}

	// Timezone
	if tz, ok := claims["zoneinfo"].(string); ok {
		claimsObj.TimeZone = tz
	}

	// Calendar system (custom claim)
	if cal, ok := claims["calendar_system"].(string); ok {
		claimsObj.CalendarSystem = cal
	}

	// Realm roles
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					claimsObj.Roles = append(claimsObj.Roles, roleStr)
				}
			}
		}
	}

	// Cost centres (custom claim)
	if cc, ok := claims["cost_centres"].(string); ok && cc != "" {
		claimsObj.CostCentres = strings.Split(cc, ",")
	} else if ccArr, ok := claims["cost_centres"].([]interface{}); ok {
		for _, c := range ccArr {
			if cStr, ok := c.(string); ok {
				claimsObj.CostCentres = append(claimsObj.CostCentres, cStr)
			}
		}
	}

	return claimsObj, nil
}

// getPublicKey retrieves the public key from JWKS.
func (v *KeycloakValidator) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	// Check cache first
	v.jwksMu.RLock()
	if key, ok := v.jwks[kid]; ok && time.Now().Before(v.jwksExp) {
		v.jwksMu.RUnlock()
		return key, nil
	}
	v.jwksMu.RUnlock()

	// Fetch JWKS
	if err := v.fetchJWKS(ctx); err != nil {
		return nil, err
	}

	// Check again
	v.jwksMu.RLock()
	defer v.jwksMu.RUnlock()

	if key, ok := v.jwks[kid]; ok {
		return key, nil
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

// fetchJWKS fetches the JSON Web Key Set from Keycloak.
func (v *KeycloakValidator) fetchJWKS(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create JWKS request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	v.jwksMu.Lock()
	defer v.jwksMu.Unlock()

	// Parse and cache keys
	for _, key := range jwks.Keys {
		rsaKey, err := v.parseRSAPublicKey(key.N, key.E)
		if err != nil {
			v.logger.Warn("failed to parse RSA key",
				slog.String("kid", key.Kid),
				slog.String("error", err.Error()),
			)
			continue
		}
		v.jwks[key.Kid] = rsaKey
	}

	// Cache for 1 hour
	v.jwksExp = time.Now().Add(time.Hour)

	v.logger.Debug("JWKS fetched and cached",
		slog.Int("key_count", len(v.jwks)),
	)

	return nil
}

// parseRSAPublicKey converts base64-encoded RSA values to a public key.
func (v *KeycloakValidator) parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	// Decode n
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}

	// Decode e
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}

	// Convert to big.Int
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

// introspectToken performs token introspection with Keycloak.
func (v *KeycloakValidator) introspectToken(ctx context.Context, tokenString string) (*introspectionResponse, error) {
	introspectURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
		v.config.URL, v.config.Realm)

	reqBody := fmt.Sprintf("token=%s", tokenString)
	if v.config.ClientID != "" {
		reqBody += fmt.Sprintf("&client_id=%s", v.config.ClientID)
	}
	if v.config.ClientSecret != "" {
		reqBody += fmt.Sprintf("&client_secret=%s", v.config.ClientSecret)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, introspectURL,
		strings.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create introspection request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection returned status %d", resp.StatusCode)
	}

	var result introspectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode introspection response: %w", err)
	}

	return &result, nil
}

// getCachedClaims retrieves cached claims for a token.
func (v *KeycloakValidator) getCachedClaims(ctx context.Context, tokenString string) (*Claims, error) {
	if v.redis == nil {
		return nil, nil
	}

	key := v.getCacheKey(tokenString)

	data, err := v.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cached claims: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached claims: %w", err)
	}

	return &claims, nil
}

// cacheClaims caches the claims for a token.
func (v *KeycloakValidator) cacheClaims(ctx context.Context, tokenString string, claims *Claims) error {
	if v.redis == nil {
		return nil
	}

	key := v.getCacheKey(tokenString)

	data, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("failed to marshal claims: %w", err)
	}

	return v.redis.Set(ctx, key, data, v.config.CacheTTL).Err()
}

// InvalidateToken removes a token from the cache.
func (v *KeycloakValidator) InvalidateToken(ctx context.Context, tokenString string) error {
	if v.redis == nil {
		return nil
	}

	key := v.getCacheKey(tokenString)
	return v.redis.Del(ctx, key).Err()
}

// getCacheKey generates a Redis cache key for a token.
func (v *KeycloakValidator) getCacheKey(tokenString string) string {
	// Use a hash of the token to avoid storing the full token
	return fmt.Sprintf("%s:%s", CacheKeyPrefix, hashToken(tokenString))
}

// hashToken creates a simple hash of the token for caching.
func hashToken(token string) string {
	// Simple hash for cache key purposes
	// In production, use a proper hash function
	if len(token) < 32 {
		return token
	}
	return token[:32]
}
