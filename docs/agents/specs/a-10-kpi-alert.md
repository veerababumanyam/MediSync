# Agent Specification — A-10: KPI Alert Agent

**Agent ID:** `A-10`  
**Agent Name:** KPI Alert Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 3  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Continuously monitors key business metrics against user-defined thresholds on a scheduled cadence. Dispatches actionable alerts via in-app notification, email, or SMS when a threshold is breached.

> **Addresses:** PRD §6.5, US28 — "Alert me when revenue drops below ₹5L per day."

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled |
| **Scheduled trigger** | `*/15 * * * *` (every 15 minutes, configurable per alert) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `alert_configs` | `[]AlertConfig` | Alert config store | ✅ |
| `current_metric_values` | `map[string]float64` | Warehouse (A-01 query) | ✅ |

```go
type AlertConfig struct {
    AlertID      string        `json:"alert_id"`
    MetricSQL    string        `json:"metric_sql"`
    Condition    string        `json:"condition"`  // e.g. "< 500000"
    Threshold    float64       `json:"threshold"`
    Channels     []string      `json:"channels"`   // email, sms, in_app
    Recipients   []string      `json:"recipients"`
    Frequency    string        `json:"frequency"`  // cron
    CreatedBy    string        `json:"created_by"`
}
```

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `triggered_alerts` | `[]AlertEvent` | Apprise dispatcher |
| `alert_log_entry` | `AuditEvent` | audit_log |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Go cron (`robfig/cron`) | MIT | Schedule metric checks |
| 2 | PostgreSQL (read-only) | PostgreSQL | Execute metric SQL |
| 3 | Threshold evaluator (Go) | Internal | Compare value to condition |
| 4 | Apprise | MIT | Multi-channel alert dispatch |

---

## 6. Guardrails

- Alert cooldown: minimum 1h between repeated alerts for same condition (configurable).
- Alerts run under the creating user's role — OPA column masking applies.
- Maximum 50 alerts per user.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Alert trigger accuracy (no false positives) | ≥ 99% |
| Alert delivery latency from breach | < 5 minutes |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | A-01 (metric queries), Apprise |
| **Consumed by** | All user roles |
