# Agent Specification — B-15: Cash Flow Forecasting Agent

**Agent ID:** `B-15`  
**Agent Name:** Cash Flow Forecasting Agent  
**Module:** B — AI Accountant  
**Phase:** 7  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Projects future cash position from scheduled payables and receivables, identifies forthcoming shortfalls, and enables what-if scenario modelling (e.g. "what if large vendor payment is delayed 30 days?").

> **Addresses:** PRD §6.7.8 — Forward-looking cash flow projection and scenario modelling.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | "Cash Flow Forecast" screen in AI Accountant UI |
| **Scheduled trigger** | `0 8 * * 1` (every Monday 8 AM) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `forecast_horizon_days` | `int` | User selection (default 90) | ✅ |
| `scenario` | `*ScenarioConfig` | User what-if parameters | ⬜ |
| `company_id` | `string` | Multi-entity selector | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `daily_cash_positions` | `[]CashPosition` | Projected cash balance per day |
| `shortfall_dates` | `[]date` | Days where balance projected below min threshold |
| `shortfall_amounts` | `[]float64` | projected deficit per shortfall date |
| `scenario_comparison` | `*ScenarioResult` | Baseline vs scenario delta |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | PostgreSQL (read-only) | PostgreSQL | Fetch scheduled payables + receivables |
| 2 | Prophet (Python sidecar) | MIT | Extend recurring pattern forecasts |
| 3 | Go scenario engine | Internal | Apply what-if parameter overrides |

---

## 6. Guardrails

- Forecasts always shown with uncertainty bands (no point estimates without CI).
- What-if scenarios clearly labelled "Scenario" vs "Baseline."
- Read-only: no writes to warehouse or Tally.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| 30-day cash forecast MAPE | < 12% |
| Shortfall detection lead time | ≥ 7 days |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service + Prophet sidecar |
| **Depends on** | PostgreSQL, A-12 forecast sidecar (shared) |
| **Consumed by** | Finance Head, CFO dashboard |
