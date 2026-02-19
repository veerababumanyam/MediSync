# Agent Specification — D-06: Dashboard Auto-Generation Agent

**Agent ID:** `D-06`  
**Agent Name:** Dashboard Auto-Generation Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 15  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Automatically generates a personalised, role-relevant analytics dashboard from the user's natural language description or from their usage history, selecting appropriate metrics, chart types, and layout.

> **Addresses:** PRD §6.9.3, US29 — One-click dashboard creation from natural language description.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Create Dashboard" in Analytics UI |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `dashboard_description` | `string` | User input | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `role` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `dashboard_config` | `DashboardConfig` | Full dashboard definition with widgets |
| `dashboard_id` | `UUID` | Persisted dashboard ID |
| `widget_configs` | `[]WidgetConfig` | Per-widget SQL + viz config |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`dashboard-gen`) | Parse description → widget requirements |
| 2 | D-09 Semantic Layer | Map requirements to available metrics |
| 3 | A-01 | Generate SQL per widget |
| 4 | A-03 Viz Routing | Select chart type per widget |
| 5 | A-11 | Persist widgets to dashboard |

---

## 6. Guardrails

- Generated widgets limited to user's OPA scope.
- Max 20 widgets per dashboard.
- Generated SQL validated before persistence.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| User acceptance rate (kept without modifications) | ≥ 75% |
| P95 Generation Latency | < 30s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | D-09, A-01, A-03, A-11 |
| **Consumed by** | Analytics UI - Dashboard Creator |
