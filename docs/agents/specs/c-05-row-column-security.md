# Agent Specification — C-05: Row/Column Security Enforcement Agent

**Agent ID:** `C-05`  
**Agent Name:** Row/Column Security Enforcement Agent  
**Module:** C — Easy Reports  
**Phase:** 10  
**Priority:** P0 Critical  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Applies RBAC-enforced row-level and column-level security filters to every query and report, ensuring users only see data within their authorised scope. This is a middleware component, not a user-facing agent.

> **Addresses:** PRD §6.8.8, §6.9.11 — Row-level security and column masking across all modules.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | Every database query or report generation request |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `user_claims` | `JWTClaims` | JWT | ✅ |
| `query` | `string` | Requesting agent | ✅ |
| `target_schema` | `string` | Query metadata | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `filtered_query` | `string` | Query with WHERE clauses injected |
| `column_blacklist` | `[]string` | Columns to strip from result |
| `allow` | `bool` | OPA decision |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | OPA sidecar | Apache-2.0 | Policy evaluation: row filter + column blacklist |
| 2 | Go query rewriter | Internal | Inject WHERE clauses into SQL |
| 3 | PostgreSQL row security policies | PostgreSQL | Database-level enforcement (defence in depth) |

### OPA Policy (row filter)
```rego
row_filter["department"] = input.user.department {
    not user_has_role(input.user, "finance_head")
    not user_has_role(input.user, "admin")
}
```

---

## 6. Guardrails

- Defence-in-depth: OPA policy **plus** Postgres row security policies (both must allow).
- Sensitive column masking: patient PII, salaries, cost prices stripped per role.
- All policy denials logged to B-14 audit trail.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Unauthorised data access rate | 0% |
| Filter injection latency overhead | < 5ms |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go middleware (applied to all DB query paths) |
| **Depends on** | OPA sidecar |
| **Consumed by** | All agents that query PostgreSQL |
