# Agent Specification — A-11: Chart-to-Dashboard Pin Agent

**Agent ID:** `A-11`  
**Agent Name:** Chart-to-Dashboard Pin Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 3  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Accepts a "pin" action from the user on any generated chart, persists the chart configuration (query, chart type, title, filters) to the user's personal dashboard, and sets up an auto-refresh schedule.

> **Addresses:** PRD §6.2 — Pinning AI-generated charts to personalised dashboards.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User clicks "Pin to Dashboard" on a chart |
| **Calling agent** | User via Frontend API |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `chart_config` | `ChartConfig` | A-01/A-03 output | ✅ |
| `dashboard_id` | `UUID` | User dashboard store | ✅ |
| `refresh_interval` | `string` | User selection (e.g. `1h`) | ⬜ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `pinned_widget_id` | `UUID` | Dashboard renderer |
| `success` | `bool` | UI confirmation |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL INSERT | Persist chart config to `dashboard_widgets` table |
| 2 | Redis scheduler | Register auto-refresh job |
| 3 | WebSocket push | Update dashboard UI in real time |

---

## 6. Guardrails

- Max 30 pinned widgets per dashboard.
- Auto-refresh queries run under the pin owner's role (OPA masking applies).

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Pin success rate | ≥ 99.9% |
| Widget renders correctly on next load | 100% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | A-01, A-03 |
| **Consumed by** | Dashboard UI |
