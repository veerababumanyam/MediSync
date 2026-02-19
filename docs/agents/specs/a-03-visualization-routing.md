# Agent Specification — A-03: Visualization Routing Agent

**Agent ID:** `A-03`  
**Agent Name:** Visualization Routing Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Classifies the query result and intent to select the optimal chart type (bar, line, pie, table, scatter), then emits the chart-type token and layout hints to the frontend renderer.

> **Addresses:** PRD §5.2, §6.1 — Automatic chart type selection based on result shape.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (after SQL execution) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `result_set` | `[]map[string]any` | A-01 executor | ✅ |
| `user_query` | `string` | A-01 input | ✅ |
| `explanation` | `string` | A-01 LLM output | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `chart_type` | `enum` (bar/line/pie/table/scatter/kpi_card) | Frontend renderer |
| `x_axis` | `*string` | Chart config |
| `y_axis` | `[]string` | Chart config |
| `title` | `string` | Chart label |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Rule-based classifier (Go) | Internal | Fast path: result shape heuristics |
| 2 | Genkit Flow (`viz-router`) | Apache-2.0 | LLM fallback for ambiguous cases |
| 3 | Apache ECharts schema | Apache-2.0 | Validate chart config output |

### Routing Heuristics (fast path — no LLM)
- 1 dimension + 1 measure, few rows → `pie`
- 1 time dimension + 1 measure → `line`
- 1 category + 1 measure, many rows → `bar`
- Multiple measures → `bar` (grouped) or `table`
- Single numeric result → `kpi_card`
- Complex multi-column → `table`

---

## 6. Guardrails

- Fallback to `table` if classification confidence < 0.60 (always safe).
- Never emit a chart type that cannot be rendered by the frontend chart library.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Chart type accuracy (user does not switch chart) | ≥ 85% |
| Latency (rule-based path) | < 50ms |
| Latency (LLM fallback path) | < 2s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (embedded in A-01 flow) |
| **Depends on** | A-01 |
| **Consumed by** | Frontend chart renderer |
