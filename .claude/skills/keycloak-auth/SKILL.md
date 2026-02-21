---
name: keycloak-auth
description: This skill should be used when the user asks to "implement Keycloak authentication", "JWT validation", "OAuth2 flows", "Keycloak integration", "role-based access control", "RBAC", "token refresh", "OpenID Connect", or mentions Keycloak-specific concepts like realms, clients, roles, or groups.
---

# Keycloak Authentication Patterns for MediSync

Keycloak provides identity management for MediSync with JWT tokens, role-based access control, and integration with OPA for fine-grained authorization.

★ Insight ─────────────────────────────────────
MediSync's auth architecture:
1. **Keycloak** - Identity provider (IdP)
2. **JWT Tokens** - Stateless authentication
3. **OPA** - Policy-as-code authorization
4. **Roles** - `admin`, `finance_head`, `analyst`, `viewer`

JWTs contain user preferences including locale for i18n.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Token Type** | JWT (RS256 signed) |
| **Token Lifetime** | 15 minutes access, 24 hours refresh |
| **Realm** | `medisync` |
| **Clients** | `medisync-web`, `medisync-mobile` |
| **Custom Claims** | `locale`, `company_id`, `permissions` |

## JWT Token Structure

### Access Token Claims

```json
{
  "exp": 1708543200,
  "iat": 1708542300,
  "iss": "https://auth.medisync.io/realms/medisync",
  "aud": "account",
  "sub": "user-uuid",
  "typ": "Bearer",
  "azp": "medisync-web",
  "realm_access": {
    "roles": ["analyst", "viewer"]
  },
  "resource_access": {
    "account": {
      "roles": ["manage-account", "view-profile"]
    }
  },
  "preferred_username": "john.doe",
  "email": "john@example.com",
  "email_verified": true,
  "name": "John Doe",
  "locale": "en",
  "company_id": "company-uuid",
  "permissions": ["read:reports", "read:dashboard"]
}
```

## Backend Integration

### JWT Validation Middleware (Go)

```go
package auth

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/coreos/go-oidc/v3/oidc"
    "github.com/go-chi/chi/v5/middleware"
)

type KeycloakConfig struct {
    IssuerURL   string
    ClientID    string
    Realm       string
}

type Authenticator struct {
    provider *oidc.Provider
    verifier *oidc.IDTokenVerifier
    config   KeycloakConfig
}

func NewAuthenticator(ctx context.Context, cfg KeycloakConfig) (*Authenticator, error) {
    provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
    if err != nil {
        return nil, fmt.Errorf("create OIDC provider: %w", err)
    }

    verifier := provider.Verifier(&oidc.Config{
        ClientID: cfg.ClientID,
    })

    return &Authenticator{
        provider: provider,
        verifier: verifier,
        config:   cfg,
    }, nil
}

// Claims represents the JWT claims
type Claims struct {
    Subject           string            `json:"sub"`
    Email             string            `json:"email"`
    PreferredUsername string            `json:"preferred_username"`
    Name              string            `json:"name"`
    Locale            string            `json:"locale"`
    CompanyID         string            `json:"company_id"`
    RealmAccess       RealmAccess       `json:"realm_access"`
    Permissions       []string          `json:"permissions"`
}

type RealmAccess struct {
    Roles []string `json:"roles"`
}

// Middleware validates JWT tokens
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeError(w, http.StatusUnauthorized, "missing authorization header")
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            writeError(w, http.StatusUnauthorized, "invalid authorization format")
            return
        }

        token, err := a.verifier.Verify(r.Context(), tokenString)
        if err != nil {
            writeError(w, http.StatusUnauthorized, "invalid token")
            return
        }

        var claims Claims
        if err := token.Claims(&claims); err != nil {
            writeError(w, http.StatusUnauthorized, "invalid claims")
            return
        }

        // Add claims to context
        ctx := context.WithValue(r.Context(), "user", &claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// GetUserFromContext extracts user claims from context
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
    user, ok := ctx.Value("user").(*Claims)
    return user, ok
}

// RequireRoles middleware checks for required roles
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := GetUserFromContext(r.Context())
            if !ok {
                writeError(w, http.StatusUnauthorized, "unauthorized")
                return
            }

            if !hasAnyRole(user.RealmAccess.Roles, roles) {
                writeError(w, http.StatusForbidden, "insufficient permissions")
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

func hasAnyRole(userRoles, requiredRoles []string) bool {
    roleSet := make(map[string]bool)
    for _, r := range userRoles {
        roleSet[r] = true
    }
    for _, req := range requiredRoles {
        if roleSet[req] {
            return true
        }
    }
    return false
}
```

### Role-Based Route Protection

```go
func (s *Server) setupRoutes() chi.Router {
    r := chi.NewRouter()

    // Public routes
    r.Group(func(r chi.Router) {
        r.Get("/health", s.HealthHandler)
        r.Post("/login", s.LoginHandler)
        r.Post("/refresh", s.RefreshHandler)
    })

    // Authenticated routes
    r.Group(func(r chi.Router) {
        r.Use(s.auth.Middleware)

        // Viewer access
        r.Get("/api/dashboard", s.DashboardHandler)
        r.Get("/api/reports", s.ReportsHandler)

        // Analyst access
        r.Group(func(r chi.Router) {
            r.Use(RequireRoles("analyst", "admin"))
            r.Post("/api/query", s.QueryHandler)
        })

        // Finance head access
        r.Group(func(r chi.Router) {
            r.Use(RequireRoles("finance_head", "admin"))
            r.Post("/api/tally/sync", s.TallySyncHandler)
            r.Post("/api/approvals", s.ApprovalHandler)
        })

        // Admin only
        r.Group(func(r chi.Router) {
            r.Use(RequireRoles("admin"))
            r.Post("/api/users", s.CreateUserHandler)
            r.Delete("/api/users/{id}", s.DeleteUserHandler)
        })
    })

    return r
}
```

## Token Management

### Token Refresh Handler

```go
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`
    TokenType    string `json:"token_type"`
}

func (a *Authenticator) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
    // Use Keycloak's token endpoint
    tokenURL := fmt.Sprintf("%s/protocol/openid-connect/token", a.config.IssuerURL)

    data := url.Values{}
    data.Set("grant_type", "refresh_token")
    data.Set("refresh_token", refreshToken)
    data.Set("client_id", a.config.ClientID)

    resp, err := http.PostForm(tokenURL, data)
    if err != nil {
        return nil, fmt.Errorf("refresh token request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("refresh token failed: %s", resp.Status)
    }

    var tokenResp TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return nil, fmt.Errorf("decode token response: %w", err)
    }

    return &tokenResp, nil
}
```

### Token Introspection (for logout validation)

```go
func (a *Authenticator) IntrospectToken(ctx context.Context, token string) (*IntrospectionResult, error) {
    introspectURL := fmt.Sprintf("%s/protocol/openid-connect/token/introspect", a.config.IssuerURL)

    data := url.Values{}
    data.Set("token", token)
    data.Set("client_id", a.config.ClientID)

    req, _ := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(data.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result IntrospectionResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

type IntrospectionResult struct {
    Active bool     `json:"active"`
    Scope  string   `json:"scope"`
    Roles  []string `json:"realm_access.roles"`
}
```

## OPA Integration

### Policy Check

```go
type OPAClient struct {
    client   *http.Client
    baseURL  string
}

type AuthzRequest struct {
    Input AuthzInput `json:"input"`
}

type AuthzInput struct {
    User        string            `json:"user"`
    Roles       []string          `json:"roles"`
    Resource    string            `json:"resource"`
    Action      string            `json:"action"`
    Context     map[string]any    `json:"context"`
}

func (o *OPAClient) Allow(ctx context.Context, input AuthzInput) (bool, error) {
    reqBody := AuthzRequest{Input: input}

    body, err := json.Marshal(reqBody)
    if err != nil {
        return false, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/v1/data/medisync/allow", bytes.NewReader(body))
    if err != nil {
        return false, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := o.client.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    var result struct {
        Decision bool `json:"result"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, err
    }

    return result.Decision, nil
}

// AuthzMiddleware checks OPA policies
func (s *Server) AuthzMiddleware(resource, action string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, _ := GetUserFromContext(r.Context())

            allowed, err := s.opa.Allow(r.Context(), AuthzInput{
                User:     user.Subject,
                Roles:    user.RealmAccess.Roles,
                Resource: resource,
                Action:   action,
                Context: map[string]any{
                    "company_id": user.CompanyID,
                },
            })

            if err != nil || !allowed {
                writeError(w, http.StatusForbidden, "access denied")
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## Frontend Integration

### React Auth Hook

```typescript
// hooks/useAuth.ts
import { useAuth as useOidcAuth } from 'oidc-react';

interface User {
  id: string;
  email: string;
  name: string;
  locale: string;
  companyId: string;
  roles: string[];
  permissions: string[];
}

export function useAuth() {
  const auth = useOidcAuth();

  const user: User | null = auth.userData ? {
    id: auth.userData.profile.sub,
    email: auth.userData.profile.email,
    name: auth.userData.profile.name,
    locale: auth.userData.profile.locale || 'en',
    companyId: auth.userData.profile.company_id,
    roles: auth.userData.profile.realm_access?.roles || [],
    permissions: auth.userData.profile.permissions || [],
  } : null;

  const hasRole = (role: string) => user?.roles.includes(role) ?? false;
  const hasPermission = (permission: string) => user?.permissions.includes(permission) ?? false;

  return {
    user,
    isAuthenticated: !!auth.userData,
    isLoading: auth.isLoading,
    login: auth.login,
    logout: auth.logout,
    hasRole,
    hasPermission,
  };
}
```

### Protected Route Component

```typescript
interface ProtectedRouteProps {
  children: React.ReactNode;
  roles?: string[];
  permissions?: string[];
}

export function ProtectedRoute({ children, roles, permissions }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, hasRole, hasPermission } = useAuth();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  if (roles && !roles.some(hasRole)) {
    return <AccessDenied />;
  }

  if (permissions && !permissions.some(hasPermission)) {
    return <AccessDenied />;
  }

  return <>{children}</>;
}
```

## Keycloak Configuration

### Realm Export (Partial)

```json
{
  "realm": "medisync",
  "enabled": true,
  "sslRequired": "external",
  "roles": {
    "realm": [
      { "name": "admin", "description": "Full administrative access" },
      { "name": "finance_head", "description": "Finance department head" },
      { "name": "analyst", "description": "Data analyst" },
      { "name": "viewer", "description": "Read-only access" }
    ]
  },
  "clients": [
    {
      "clientId": "medisync-web",
      "enabled": true,
      "protocol": "openid-connect",
      "publicClient": true,
      "redirectUris": ["https://app.medisync.io/*"],
      "webOrigins": ["https://app.medisync.io"],
      "standardFlowEnabled": true,
      "implicitFlowEnabled": false,
      "directAccessGrantsEnabled": true
    }
  ],
  "browserSecurityHeaders": {
    "contentSecurityPolicyReportOnly": "",
    "xContentTypeOptions": "nosniff",
    "xRobotsTag": "none",
    "xFrameOptions": "SAMEORIGIN",
    "contentSecurityPolicy": "frame-src 'self'; frame-ancestors 'self'; object-src 'none';",
    "xXSSProtection": "1; mode=block",
    "strictTransportSecurity": "max-age=31536000; includeSubDomains"
  }
}
```

## Additional Resources

### Reference Files
- **`references/opa-policies.md`** - OPA policy examples
- **`references/token-flows.md`** - OAuth2/OIDC flow details

### Example Files
- **`examples/middleware.go`** - Complete middleware implementation
- **`examples/opa-client.go`** - OPA client implementation
