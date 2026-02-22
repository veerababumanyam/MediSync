---
name: opa-policy
description: This skill should be used when the user asks to "create OPA policies", "write Rego rules", "implement authorization", "add policy-as-code", "role-based access control", "RBAC policies", "OPA integration", or mentions Open Policy Agent or Rego policy language for MediSync.
---

# OPA Policy-as-Code for MediSync

Open Policy Agent (OPA) with Rego provides fine-grained, auditable authorization for MediSync. All access control decisions flow through OPA policies, ensuring consistent security across the platform.

★ Insight ─────────────────────────────────────
MediSync authorization architecture:
1. **Keycloak** - Identity and role assignment
2. **OPA** - Policy evaluation engine
3. **Rego** - Declarative policy language
4. **Audit trail** - All decisions logged
5. **HITL gates** - Human approval for writes
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Engine** | Open Policy Agent (latest) |
| **Language** | Rego |
| **Decision API** | POST /v1/data/medisync/allow |
| **Input Format** | JSON with user, resource, action, context |
| **Bundle** | Policies loaded from /policies directory |

## Project Structure

```
medisync/
├── policies/
│   ├── medisync/
│   │   ├── main.rego           # Entry point
│   │   ├── roles.rego          # Role definitions
│   │   ├── resources.rego      # Resource patterns
│   │   ├── dashboard.rego      # Dashboard policies
│   │   ├── reports.rego        # Reports policies
│   │   ├── tally.rego          # Tally sync policies
│   │   ├── documents.rego      # Document policies
│   │   └── util.rego           # Helper functions
│   └── data.json               # Static policy data
├── internal/
│   └── auth/
│       └── opa.go              # OPA client
└── tests/
    └── policies/
        ├── dashboard_test.rego
        └── tally_test.rego
```

## Core Policy Structure

### Main Entry Point

```rego
# policies/medisync/main.rego
package medisync

import future.keywords.if
import future.keywords.in

# Main allow decision
default allow := false

allow if {
    # User must be authenticated
    input.user != ""

    # Check role-based access
    rbac.allow

    # Check resource-specific rules
    resource_allowed
}

# Role-Based Access Control
rbac := {
    "allow": true,
} if {
    some role in input.roles
    role_grants_permission(role, input.resource, input.action)
}
```

### Role Definitions

```rego
# policies/medisync/roles.rego
package medisync.roles

import future.keywords.if
import future.keywords.in

# Role hierarchy
role_hierarchy := {
    "admin": ["finance_head", "analyst", "viewer"],
    "finance_head": ["analyst", "viewer"],
    "analyst": ["viewer"],
    "viewer": [],
}

# Role permissions map
role_permissions := {
    "admin": {
        "dashboard": ["read", "write", "delete"],
        "reports": ["read", "write", "delete", "export"],
        "tally": ["read", "write", "sync", "approve"],
        "documents": ["read", "write", "delete", "approve"],
        "users": ["read", "write", "delete"],
    },
    "finance_head": {
        "dashboard": ["read"],
        "reports": ["read", "write", "export"],
        "tally": ["read", "write", "sync", "approve"],
        "documents": ["read", "write", "approve"],
        "users": ["read"],
    },
    "analyst": {
        "dashboard": ["read"],
        "reports": ["read", "write", "export"],
        "tally": ["read"],
        "documents": ["read", "write"],
    },
    "viewer": {
        "dashboard": ["read"],
        "reports": ["read"],
        "tally": ["read"],
        "documents": ["read"],
    },
}

# Get effective permissions (including inherited roles)
get_effective_permissions(role) := permissions if {
    some permissions in [role_permissions[role]]
    some inherited_role in role_hierarchy[role]
    inherited_perms := get_effective_permissions(inherited_role)
    permissions := permissions | inherited_perms
}

# Check if role grants permission
role_grants_permission(role, resource, action) if {
    perms := role_permissions[role]
    resource_perms := perms[resource]
    action in resource_perms
}
```

### Resource Policies

```rego
# policies/medisync/dashboard.rego
package medisync.dashboard

import future.keywords.if

# Dashboard read access
allow_read if {
    input.action == "read"
    has_any_role(["viewer", "analyst", "finance_head", "admin"])
}

# Dashboard write (create/edit widgets)
allow_write if {
    input.action == "write"
    has_any_role(["analyst", "finance_head", "admin"])
}

# Dashboard delete
allow_delete if {
    input.action == "delete"
    has_any_role(["admin"])
}

# Company isolation
company_isolated if {
    # User can only see their company's data
    input.context.company_id == input.user_company_id
}

has_any_role(required_roles) if {
    some role in required_roles
    role in input.roles
}
```

### Tally Sync Policies

```rego
# policies/medisync/tally.rego
package medisync.tally

import future.keywords.if
import future.keywords.in

# Tally sync requires approval
allow_sync if {
    input.action == "sync"

    # Must be finance_head or admin
    has_any_role(["finance_head", "admin"])

    # Must have human approval
    input.context.approved == true

    # Approval must be recent (within 24 hours)
    approval_is_recent
}

# Tally read access
allow_read if {
    input.action == "read"
    has_any_role(["viewer", "analyst", "finance_head", "admin"])
}

# Approval workflow
can_approve if {
    input.action == "approve"
    has_any_role(["finance_head", "admin"])

    # Cannot approve own entries
    input.context.entry_creator != input.user
}

approval_is_recent if {
    now := time.now_ns()
    approval_time := input.context.approval_timestamp
    (now - approval_time) / 3600000000000 < 24  # hours
}

has_any_role(required_roles) if {
    some role in required_roles
    role in input.roles
}
```

### Document Policies

```rego
# policies/medisync/documents.rego
package medisync.documents

import future.keywords.if
import future.keywords.in

# Document upload
allow_upload if {
    input.action == "upload"
    has_any_role(["analyst", "finance_head", "admin"])
}

# Document read
allow_read if {
    input.action == "read"
    has_any_role(["viewer", "analyst", "finance_head", "admin"])

    # Company isolation
    input.context.company_id == input.user_company_id
}

# Document delete
allow_delete if {
    input.action == "delete"
    has_any_role(["admin"])

    # Cannot delete if document has been synced to Tally
    not input.context.synced_to_tally
}

# Document approval (OCR confidence-based)
can_approve if {
    input.action == "approve"
    has_any_role(["finance_head", "admin"])

    # High confidence documents can be auto-approved
    # Low confidence requires explicit approval
    confidence_ok
}

confidence_ok if {
    # High confidence (> 95%) - any finance role can approve
    input.context.ocr_confidence > 0.95
    has_any_role(["analyst", "finance_head", "admin"])
}

confidence_ok if {
    # Medium confidence (> 80%) - finance_head or admin
    input.context.ocr_confidence > 0.80
    has_any_role(["finance_head", "admin"])
}

confidence_ok if {
    # Low confidence (<= 80%) - admin only
    input.context.ocr_confidence <= 0.80
    has_any_role(["admin"])
}

has_any_role(required_roles) if {
    some role in required_roles
    role in input.roles
}
```

### Utility Functions

```rego
# policies/medisync/util.rego
package medisync.util

import future.keywords.if
import future.keywords.in

# Check if user has specific role
has_role(role) if {
    role in input.roles
}

# Check if user has any of the specified roles
has_any_role(roles) if {
    some role in roles
    has_role(role)
}

# Check if user has all specified roles
has_all_roles(roles) if {
    all([has_role(role) | some role in roles])
}

# Company isolation check
same_company if {
    input.user_company_id == input.context.company_id
}

# Time-based checks
is_business_hours if {
    now := time.now_ns()
    hour := time.clock(now)[0]
    hour >= 9
    hour < 18
}

# Working day check
is_working_day if {
    now := time.now_ns()
    weekday := time.weekday(now)
    weekday != "Saturday"
    weekday != "Sunday"
}
```

## Go Integration

### OPA Client

```go
// internal/auth/opa.go
package auth

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type OPAClient struct {
    client  *http.Client
    baseURL string
}

func NewOPAClient(baseURL string) *OPAClient {
    return &OPAClient{
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
        baseURL: baseURL,
    }
}

type AuthzInput struct {
    User        string         `json:"user"`
    Roles       []string       `json:"roles"`
    Resource    string         `json:"resource"`
    Action      string         `json:"action"`
    Context     map[string]any `json:"context"`
    UserCompanyID string       `json:"user_company_id"`
}

type AuthzRequest struct {
    Input AuthzInput `json:"input"`
}

type AuthzResponse struct {
    Result bool `json:"result"`
}

func (o *OPAClient) Allow(ctx context.Context, input AuthzInput) (bool, error) {
    reqBody := AuthzRequest{Input: input}

    body, err := json.Marshal(reqBody)
    if err != nil {
        return false, fmt.Errorf("marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST",
        o.baseURL+"/v1/data/medisync/allow",
        bytes.NewReader(body))
    if err != nil {
        return false, fmt.Errorf("create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := o.client.Do(req)
    if err != nil {
        return false, fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return false, fmt.Errorf("OPA error: %s", string(body))
    }

    var result AuthzResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, fmt.Errorf("decode response: %w", err)
    }

    return result.Result, nil
}
```

### Authorization Middleware

```go
// internal/api/middleware/authz.go
package middleware

import (
    "net/http"

    "medisync/internal/auth"
)

func AuthzMiddleware(opa *auth.OPAClient, resource, action string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := auth.GetUserFromContext(r.Context())
            if !ok {
                writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
                return
            }

            allowed, err := opa.Allow(r.Context(), auth.AuthzInput{
                User:          user.ID,
                Roles:         user.Roles,
                Resource:      resource,
                Action:        action,
                UserCompanyID: user.CompanyID,
                Context: map[string]any{
                    "company_id": user.CompanyID,
                },
            })

            if err != nil || !allowed {
                writeError(w, http.StatusForbidden, "FORBIDDEN", "Access denied")
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## Testing Policies

### Unit Tests

```rego
# tests/policies/tally_test.rego
package medisync.tally_test

import future.keywords.if
import data.medisync.tally

# Test case: finance_head can sync with approval
test_sync_with_approval if {
    allow_sync with input as {
        "user": "user1",
        "roles": ["finance_head"],
        "action": "sync",
        "context": {
            "approved": true,
            "approval_timestamp": time.now_ns(),
        },
    }
}

# Test case: analyst cannot sync
test_sync_denied_for_analyst if {
    not allow_sync with input as {
        "user": "user2",
        "roles": ["analyst"],
        "action": "sync",
        "context": {
            "approved": true,
        },
    }
}

# Test case: cannot approve own entry
test_cannot_approve_own_entry if {
    not can_approve with input as {
        "user": "user1",
        "roles": ["finance_head"],
        "action": "approve",
        "context": {
            "entry_creator": "user1",
        },
    }
}
```

### Running Tests

```bash
# Run OPA tests
opa test policies/ tests/

# Run with coverage
opa test --coverage policies/ tests/

# Run specific test
opa test policies/ tests/ -r test_sync_with_approval
```

## Additional Resources

### Reference Files
- **`references/policy-patterns.md`** - Advanced Rego patterns
- **`references/audit-logging.md`** - Decision audit trail setup

### Example Files
- **`examples/audit-policy.rego`** - Comprehensive audit logging policy
- **`examples/data-filtering.rego`** - Row-level security patterns

### Scripts
- **`scripts/test-policies.sh`** - Run policy test suite
- **`scripts/bundle-policies.sh`** - Create OPA bundle for deployment
