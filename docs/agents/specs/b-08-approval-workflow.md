# Agent Specification — B-08: Approval Workflow Agent

**Agent ID:** `B-08`  
**Agent Name:** Approval Workflow Agent  
**Module:** B — AI Accountant  
**Phase:** 5  
**Priority:** P0 Critical  
**HITL Required:** Yes — always; this agent IS the HITL gate  
**Status:** Draft

---

## 1. Purpose

Routes transactions through a configurable multi-step approval chain (Accountant → Finance Manager → Finance Head) before any data is written to Tally. Enforces separation of duties, tracks state, sends reminders, and is the sole gate that authorises B-09 to execute.

> **Addresses:** PRD §6.7.5, US14 — Approval workflow with role-based chain and audit trail.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-05 / B-06 (post mapping, after user confirms ledger suggestion) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `transactions` | `[]Transaction` | B-05 output (user-confirmed) | ✅ |
| `approval_policy` | `ApprovalPolicy` | OPA policy + tenant config | ✅ |
| `requesting_user` | `User` | JWT | ✅ |
| `company_id` | `string` | Multi-entity selector | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `workflow_id` | `UUID` | Unique workflow instance |
| `current_status` | `enum` | `draft / pending_accountant / pending_manager / pending_finance / approved / rejected` |
| `approval_history` | `[]ApprovalEvent` | Timestamped action per approver |
| `rejection_reason` | `*string` | Set on rejection |
| `approved_transactions` | `[]Transaction` | Passed to B-09 when fully approved |

```go
type ApprovalEvent struct {
    EventID    UUID      `json:"event_id"`
    WorkflowID UUID      `json:"workflow_id"`
    ActorID    string    `json:"actor_id"`
    ActorRole  string    `json:"actor_role"`
    Action     string    `json:"action"` // approve / reject / comment
    Comment    *string   `json:"comment"`
    Timestamp  time.Time `json:"timestamp"`
}
```

---

## 5. Tool Chain

```
Confirmed Mapping (from B-05)
  → OPA policy check: has_permission(user, 'submit_for_approval')
  → Workflow State Machine (Go FSM + Postgres state table)
  → Notification Dispatcher (Apprise → Email + In-App)
  → Reminder Scheduler (Redis cron — 24h for stale approvals)
  → On final approval: emit event → B-09
  → B-14 Audit Log Writer (every state transition)
```

### Approval Chain (configurable per tenant)
```
DRAFT
  → submitted by Accountant
PENDING_ACCOUNTANT
  → approved by Accountant (different from submitter)
PENDING_MANAGER
  → approved by Finance Manager
PENDING_FINANCE
  → approved by Finance Head
APPROVED
  → triggers B-09 (on explicit user sync action)
```

---

## 6. Guardrails

| # | Guard | Enforcement |
|---|-------|-------------|
| 1 | Self-approval blocked | OPA hard policy: `submitted_by != approver_id` |
| 2 | Bulk approval threshold | `sum(amounts) > ₹1,00,000` → Finance Head required |
| 3 | Immutable audit log | Every state transition written to `audit_log` (no-update/delete) |
| 4 | Stale approval reminder | Celery Beat job: 24h reminder; 48h escalation |
| 5 | Rejection notification | Submitter notified with reason on rejection |

---

## 7. HITL Gate

This agent IS the HITL gate. Every approval action is a human decision:

| Step | Required Human | Action |
|------|---------------|--------|
| 1 | `accountant` (not self) | First-level review |
| 2 | `finance_manager` | Second-level approval |
| 3 | `finance_head` | Final sign-off |

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Self-approval block rate | 100% (hard OPA policy) |
| Workflow completion within SLA (72h) | ≥ 90% |
| Stale reminder delivery | ≥ 99% |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (FastAPI-style chi router) |
| **State store** | PostgreSQL `approval_workflows` table |
| **Reminder scheduler** | Redis + `robfig/cron` |
| **Depends on** | B-05, OPA sidecar, Apprise, B-14 |
| **Consumed by** | B-09 (only after full approval) |
| **Env vars** | `OPA_SIDECAR_URL`, `APPRISE_URL`, `REDIS_DSN` |
