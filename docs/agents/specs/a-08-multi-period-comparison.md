# Agent Specification — A-08: Multi-Period Comparison Agent

**Agent ID:** `A-08`  
**Agent Name:** Multi-Period Comparison Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 3  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Automatically constructs period-over-period comparison queries (MoM, QoQ, YoY) from a single user request, then executes and returns the combined result set with delta and percentage change calculations.

> **Addresses:** PRD §6.2, US4 — "Compare this month's revenue to last month."

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (when comparison intent detected) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `base_sql` | `string` | A-01 | ✅ |
| `comparison_type` | `enum` | Detected by A-01 LLM (`MoM / QoQ / YoY / custom`) | ✅ |
| `reference_date` | `date` | User context / system date | ✅ |
| `schema_context` | `JSON` | Schema Context Cache | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `current_period_result` | `[]map[string]any` | Current period data |
| `prior_period_result` | `[]map[string]any` | Comparison period data |
| `delta` | `[]map[string]any` | Absolute + % change per metric |
| `chart_type` | `enum` | Suggested: `grouped_bar` or `line` |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`multi-period-sql`) | Generate date-shifted comparison queries |
| 2 | PostgreSQL executor × 2 | Execute current + prior period queries |
| 3 | Go delta calculator | Compute absolute + % change |

---

## 6. Guardrails

- Both period queries subject to OPA read-only policy.
- Maximum custom date range: 366 days (prevent unbounded queries).

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Period calculation correctness | ≥ 99% |
| P95 Latency | < 8s (two queries + delta) |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | A-01 |
| **Consumed by** | User via A-01 pipeline |
