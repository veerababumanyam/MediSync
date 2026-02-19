# Agent Specification — D-13: Scheduled Autonomous Monitoring Agent

**Agent ID:** `D-13`  
**Agent Name:** Scheduled Autonomous Monitoring Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 14  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Runs a library of user-defined or system-defined monitoring rules on a schedule, evaluating KPIs against thresholds, detecting deviations, and triggering alerts or D-04 Spotter analysis when conditions are met.

> **Addresses:** PRD §6.9.2, US28 — Continuous proactive monitoring without manual queries.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled + Event-driven |
| **Scheduled trigger** | Configurable per monitor rule (default: hourly) |
| **Event trigger** | New ETL batch completes |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `monitor_rules` | `[]MonitorRule` | Monitor config store | ✅ |
| `tenant_id` | `string` | Scheduler | ✅ |

```go
type MonitorRule struct {
    RuleID       UUID    `json:"rule_id"`
    Name         string  `json:"name"`
    MetricID     UUID    `json:"metric_id"`    // from D-09
    Comparator   string  `json:"comparator"`   // gt|lt|eq|pct_change
    Threshold    float64 `json:"threshold"`
    Cron         string  `json:"cron"`
    AlertSeverity string `json:"alert_severity"` // critical|high|medium
    AlertChannels []string `json:"alert_channels"`
}
```

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `triggered_alerts` | `[]Alert` | Alerts fired in this run |
| `evaluation_results` | `[]EvalResult` | Per-rule result (pass/fail/score) |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | `robfig/cron` (Go) | Schedule triggers |
| 2 | A-01 | Execute metric SQL for each rule |
| 3 | Go evaluator | Apply threshold logic |
| 4 | A-10 KPI Alert | Route critical alerts |
| 5 | D-04 trigger | Initiate Spotter Brief on critical cluster |
| 6 | B-14 | Log all evaluation results |

---

## 6. Guardrails

- All metric queries run read-only.
- Duplicate alert suppression: same rule cannot fire again within its cooldown period.
- Max 500 rules per tenant.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Alert false positive rate | < 5% |
| Rule evaluation on-time rate | ≥ 99% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | A-01, A-10, B-14, D-04 |
| **Consumed by** | D-04, Monitoring dashboard |
