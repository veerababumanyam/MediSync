# Tally Integration Security Guide

Security policies, HITL gates, and OPA authorization for Tally synchronization.

## Security Layers

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Security Layers                             │
├─────────────────────────────────────────────────────────────────────┤
│  Layer 1: Authentication (Keycloak)                                 │
│           └─ JWT token validation                                   │
│                                                                      │
│  Layer 2: HITL Gate (Human-in-the-Loop)                            │
│           └─ Multi-level approval workflow (Agent B-08)             │
│                                                                      │
│  Layer 3: OPA Authorization (Policy Engine)                         │
│           └─ Fine-grained permission checks                         │
│                                                                      │
│  Layer 4: Tally Gateway Security                                    │
│           └─ Network isolation, TLS, optional basic auth            │
│                                                                      │
│  Layer 5: Audit Logging                                             │
│           └─ Immutable record of all sync operations                │
└─────────────────────────────────────────────────────────────────────┘
```

## HITL Gate Configuration

Human-in-the-Loop (HITL) gates are the primary security control for Tally synchronization.

### Approval Thresholds

```yaml
# config/approval_thresholds.yaml
approval_thresholds:
  auto_approve:
    max_amount: 1000  # AED/SAR/USD
    require_ocr_confidence: 0.95

  single_approval:
    max_amount: 10000
    required_roles: ["manager", "accountant"]

  dual_approval:
    max_amount: 50000
    required_roles: ["manager", "finance_manager"]

  tri_approval:
    max_amount: null  # No limit
    required_roles: ["manager", "finance_manager", "finance_head"]
```

### Approval State Machine

```
                    ┌──────────────┐
                    │    draft     │  Initial state
                    └──────┬───────┘
                           │ Submit
                           ▼
                    ┌──────────────┐
                    │   pending    │  Waiting for Level 1 approval
                    └──────┬───────┘
                           │ Level 1 approves
                           ▼
                    ┌──────────────┐
                    │  level2_ok   │  Waiting for Level 2 (if needed)
                    └──────┬───────┘
                           │ Level 2 approves
                           ▼
                    ┌──────────────┐
                    │   approved   │  Ready for Tally sync
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   syncing    │  In progress
                    └──────┬───────┘
                           │
            ┌──────────────┴──────────────┐
            ▼                             ▼
     ┌──────────────┐            ┌──────────────┐
     │   synced     │            │   failed     │
     └──────────────┘            └──────────────┘
```

### Approval Implementation

```go
type ApprovalService struct {
    repo      ApprovalRepository
    opa       OPAClient
    notifier  NotificationService
    thresholds ApprovalThresholds
}

func (s *ApprovalService) RequestApproval(ctx context.Context, entry *JournalEntry) error {
    // Determine required levels
    levels := s.getRequiredLevels(entry.TotalAmount)

    // Create approval request
    request := &ApprovalRequest{
        ID:            uuid.New().String(),
        EntryID:       entry.ID,
        CompanyID:     entry.CompanyID,
        Amount:        entry.TotalAmount,
        RequiredLevel: levels,
        CurrentLevel:  0,
        Status:        "pending",
        CreatedAt:     time.Now(),
    }

    if err := s.repo.Create(ctx, request); err != nil {
        return err
    }

    // Notify approvers
    return s.notifier.NotifyApprovers(ctx, request)
}

func (s *ApprovalService) Approve(ctx context.Context, requestID, userID, comment string) error {
    // Get request
    request, err := s.repo.Get(ctx, requestID)
    if err != nil {
        return err
    }

    // Verify user is authorized for this level
    if !s.isAuthorizedForLevel(ctx, userID, request.CurrentLevel+1) {
        return fmt.Errorf("user not authorized for approval level %d", request.CurrentLevel+1)
    }

    // Record approval
    approval := &Approval{
        RequestID: requestID,
        UserID:    userID,
        Level:     request.CurrentLevel + 1,
        Decision:  "approved",
        Comment:   comment,
        Timestamp: time.Now(),
    }

    if err := s.repo.AddApproval(ctx, approval); err != nil {
        return err
    }

    // Update request state
    request.CurrentLevel++

    if request.CurrentLevel >= request.RequiredLevel {
        request.Status = "approved"
        // Trigger Tally sync
        return s.triggerSync(ctx, request.EntryID)
    }

    return s.repo.Update(ctx, request)
}

func (s *ApprovalService) Reject(ctx context.Context, requestID, userID, reason string) error {
    request, err := s.repo.Get(ctx, requestID)
    if err != nil {
        return err
    }

    // Record rejection
    approval := &Approval{
        RequestID: requestID,
        UserID:    userID,
        Decision:  "rejected",
        Comment:   reason,
        Timestamp: time.Now(),
    }

    if err := s.repo.AddApproval(ctx, approval); err != nil {
        return err
    }

    request.Status = "rejected"
    return s.repo.Update(ctx, request)
}
```

## OPA Policy Configuration

### Policy Structure

```
policies/
├── rego/
│   ├── tally_sync.rego       # Main sync authorization
│   ├── approval.rego         # Approval workflow rules
│   └── amount_limits.rego    # Role-based amount limits
└── tests/
    └── tally_sync_test.rego  # Policy tests
```

### Main Sync Policy (tally_sync.rego)

```rego
package medisync.tally

default allow = false

# Allow sync if all conditions are met
allow {
    input.action == "sync_to_tally"
    has_required_approvals(input)
    user_has_permission(input.user_roles, input.entry)
    within_amount_limits(input.user, input.entry)
    entry_not_duplicates(input.entry)
    company_active(input.entry.company_id)
}

# Check that all required approval levels are complete
has_required_approvals(input) {
    entry := data.entries[input.entry_id]
    required := get_required_levels(entry.total_amount)
    entry.approvals_completed >= required
}

# User has tally.sync permission
user_has_permission(roles, entry) {
    "tally.sync" in roles.permissions
}

# Amount within user's limit
within_amount_limits(user, entry) {
    entry.total_amount <= user.max_sync_amount
}

# Check for duplicate vouchers (by RemoteID)
entry_not_duplicates(entry) {
    not data.sync_log[entry.remote_id]
}

# Company is active and not suspended
company_active(company_id) {
    data.companies[company_id].status == "active"
}

# Calculate required approval levels based on amount
get_required_levels(amount) := levels {
    amount < 1000
    levels := 0  # Auto-approve
}
get_required_levels(amount) := levels {
    amount >= 1000
    amount < 10000
    levels := 1
}
get_required_levels(amount) := levels {
    amount >= 10000
    amount < 50000
    levels := 2
}
get_required_levels(amount) := levels {
    amount >= 50000
    levels := 3
}
```

### Approval Policy (approval.rego)

```rego
package medisync.approval

# Can user approve at specific level?
allow_approval(user, level, entry_amount) {
    some role
    role := user.roles[_]
    can_approve_at_level(role, level)
    within_role_limit(role, entry_amount)
}

can_approve_at_level("manager", 1)
can_approve_at_level("finance_manager", 2)
can_approve_at_level("finance_head", 3)

within_role_limit("manager", amount) {
    amount < 10000
}
within_role_limit("finance_manager", amount) {
    amount < 50000
}
within_role_limit("finance_head", amount) {
    true  # No limit
}

# Prevent self-approval for creator
allow_approval(user, level, entry_amount) {
    user.id != entry.created_by
    can_approve_at_level(user.role, level)
    within_role_limit(user.role, entry_amount)
}
```

### Go OPA Client

```go
type OPAClient struct {
    client *opa.Client
}

func NewOPAClient(url string) (*OPAClient, error) {
    client, err := opa.NewClient(opa.Config{
        URL:    url,
        Poll:   false,
    })
    if err != nil {
        return nil, err
    }

    return &OPAClient{client: client}, nil
}

func (c *OPAClient) Allow(ctx context.Context, rule string, input interface{}) (bool, error) {
    resp, err := c.client.Decision(ctx, opa.DecisionOptions{
        Path:  "/medisync/tally/allow",
        Input: input,
    })
    if err != nil {
        return false, err
    }

    // Check if result is true
    if result, ok := resp.Result.(bool); ok {
        return result, nil
    }

    return false, fmt.Errorf("unexpected OPA response type")
}

func (c *OPAClient) CanApprove(ctx context.Context, userID, requestID string, level int) (bool, error) {
    input := map[string]interface{}{
        "action":     "approve",
        "user_id":    userID,
        "request_id": requestID,
        "level":      level,
    }

    return c.Allow(ctx, "/medisync/approval/allow_approval", input)
}
```

## Tally Gateway Security

### Network Isolation

```yaml
# docker-compose.yml
services:
  tally-gateway-proxy:
    image: nginx:alpine
    networks:
      - tally-network
    ports:
      - "9000:9000"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - tally

  tally:
    image: tallysolutions/tallyprime:latest
    networks:
      - tally-network
    expose:
      - "9000"
    environment:
      - TALLY_HTTP_ENABLE=true
      - TALLY_HTTP_PORT=9000

networks:
  tally-network:
    internal: true  # Isolate from external network
```

### Nginx Reverse Proxy

```nginx
# nginx.conf
events {}
http {
    upstream tally {
        server tally:9000;
    }

    server {
        listen 9000;

        # Only allow POST requests (Tally requirement)
        limit_except POST {
            deny all;
        }

        # Require Basic Auth (optional)
        auth_basic "Tally Gateway";
        auth_basic_user_file /etc/nginx/.htpasswd;

        # Rate limiting
        limit_req_zone $binary_remote_addr zone=sync:10m rate=10r/m;
        limit_req zone=sync burst=5;

        location / {
            proxy_pass http://tally;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

            # Tally requires connection close
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
    }
}
```

### IP Whitelist

```go
type IPWhitelistMiddleware struct {
    allowedIPs []string
}

func (m *IPWhitelistMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := strings.Split(r.RemoteAddr, ":")[0]

        if !m.isAllowed(ip) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (m *IPWhitelistMiddleware) isAllowed(ip string) bool {
    for _, allowed := range m.allowedIPs {
        if ip == allowed {
            return true
        }
    }
    return false
}
```

## Audit Logging

### Immutable Audit Log

```go
type AuditLog struct {
    ID           string                 `json:"id" db:"id"`
    UserID       string                 `json:"user_id" db:"user_id"`
    Action       string                 `json:"action" db:"action"`
    Resource     string                 `json:"resource" db:"resource"`
    Status       string                 `json:"status" db:"status"`
    Details      map[string]interface{} `json:"details" db:"details"`
    IP           string                 `json:"ip" db:"ip"`
    UserAgent    string                 `json:"user_agent" db:"user_agent"`
    Timestamp    time.Time              `json:"timestamp" db:"timestamp"`
}

type AuditLogger struct {
    db       *sql.DB
    table    string
    readOnly bool
}

func (l *AuditLogger) Log(ctx context.Context, entry *AuditLog) error {
    // Use INSERT with ON CONFLICT to prevent overwrites
    query := fmt.Sprintf(`
        INSERT INTO %s (id, user_id, action, resource, status, details, ip, user_agent, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO NOTHING  -- Prevent overwrites
    `, l.table)

    _, err := l.db.ExecContext(ctx, query,
        entry.ID,
        entry.UserID,
        entry.Action,
        entry.Resource,
        entry.Status,
        entry.Details,
        entry.IP,
        entry.UserAgent,
        entry.Timestamp,
    )

    return err
}

// Query only - never update/delete
func (l *AuditLogger) Query(ctx context.Context, filter AuditFilter) ([]AuditLog, error) {
    query := fmt.Sprintf(`
        SELECT id, user_id, action, resource, status, details, ip, user_agent, timestamp
        FROM %s
        WHERE timestamp >= $1 AND timestamp <= $2
        ORDER BY timestamp DESC
        LIMIT 1000
    `, l.table)

    rows, err := l.db.QueryContext(ctx, query, filter.StartTime, filter.EndTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []AuditLog
    for rows.Next() {
        var log AuditLog
        if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource,
            &log.Status, &log.Details, &log.IP, &log.UserAgent, &log.Timestamp); err != nil {
            return nil, err
        }
        logs = append(logs, log)
    }

    return logs, nil
}
```

### Audit Table Schema

```sql
CREATE TABLE audit_log (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    status TEXT NOT NULL,
    details JSONB,
    ip TEXT,
    user_agent TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Security: prevent updates/deletes
    CONSTRAINT no_update CHECK (false = TRUE)  -- Always false
);

-- Create index for querying
CREATE INDEX idx_audit_timestamp ON audit_log(timestamp DESC);

-- Add comment for documentation
COMMENT ON TABLE audit_log IS 'Immutable audit log - no UPDATE/DELETE allowed';

-- Grant read-only to applications
GRANT INSERT ON audit_log TO medisync_app;
GRANT SELECT ON audit_log TO medisync_readonly;
REVOKE UPDATE, DELETE ON audit_log FROM medisync_app;
```

## Security Checklist

Before enabling Tally sync:

- [ ] HITL approval workflow configured
- [ ] Approval thresholds defined
- [ ] OPA policies deployed and tested
- [ ] Audit logging enabled
- [ ] Network isolation in place
- [ ] Rate limiting configured
- [ ] IP whitelist (if applicable)
- [ ] TLS enabled for external access
- [ ] Monitoring/alerting for failed syncs
- [ ] Backup/recovery procedure tested
