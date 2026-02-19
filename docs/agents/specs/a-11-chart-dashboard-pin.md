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

Accepts a "pin to dashboard" action from a user on any chart, persists the chart configuration (query, chart type, title, layout position) to the user's personal dashboard store, and configures auto-refresh.

> **Addresses:** PRD §6.2 — User-curated personal dashboards with pinned charts.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Pin to Dashboard" button click in chart UI |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `chart_config` | `ChartConfig` | A-01/A-03 output | ✅ |
| `dashboard_id` | `UUID` | User's dashboard | ✅ |
| `refresh_interval` | `*Duration` | User setting (default 15m) | ⬜ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `widget_id` | `UUID` | Dashboard store |
| `success` | `bool` | Frontend confirmation |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL | Persist widget config to `dashboard_widgets` table |
| 2 | Redis | Schedule auto-refresh job |
| 3 | WebSocket push | Notify frontend: dashboard updated |

---

## 6. Guardrails

- Maximum 20 pinned widgets per dashboard.
- Refresh interval minimum: 5 minutes (prevent DB overload).

---

## 7. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | A-01, A-03 |
| **Consumed by** | Dashboard renderer |
