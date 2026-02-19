# Agent Specification — A-13: Anomaly Detection Agent

**Agent ID:** `A-13`  
**Agent Name:** Anomaly Detection Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 7  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Runs a scheduled scan of all monitored metrics; surfaces statistically significant outliers using Z-score and IQR methods; generates plain-language explanations for each anomaly; notifies relevant stakeholders.

> **Addresses:** PRD §6.6, §6.9.4 — Proactive anomaly surfacing with contextual explanations.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled |
| **Scheduled trigger** | `0 * * * *` (hourly) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `monitored_metrics` | `[]MetricConfig` | Semantic layer registry | ✅ |
| `lookback_window` | `int` | Config (default: 30 days) | ✅ |
| `sensitivity` | `float64` | Config (default: 2.5 Z-score threshold) | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `anomalies` | `[]AnomalyEvent` | Metric, value, expected range, explanation |
| `anomaly_count` | `int` | Total detected |
| `notifications_sent` | `int` | Via Apprise |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Go cron scheduler | MIT | Hourly trigger |
| 2 | PostgreSQL (read-only) | PostgreSQL | Fetch metric time series |
| 3 | Z-score + IQR detector (Go) | Internal | Statistical outlier detection |
| 4 | Genkit Flow (`anomaly-explain`) | Apache-2.0 | Plain-language explanation of anomaly |
| 5 | Apprise | MIT | Multi-channel notification |

---

## 6. Guardrails

- Anomalies suppressed during scheduled maintenance windows.
- Cooldown: same anomaly not re-alerted within 2h.
- Explanations generated only for anomalies above severity threshold (Z > 3.0).

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Precision (true anomalies / all flagged) | ≥ 80% |
| Recall (true anomalies detected) | ≥ 85% |
| P95 full scan latency | < 60s |
| False positive rate | < 5% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | Semantic Layer Registry, Apprise |
| **Consumed by** | D-07 (feeds detected anomalies), Alert dashboard |
