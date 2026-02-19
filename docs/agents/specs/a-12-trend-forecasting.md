# Agent Specification — A-12: Trend Forecasting Agent

**Agent ID:** `A-12`  
**Agent Name:** Trend Forecasting Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 7  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Extends historical time-series data into future forecast periods using Prophet (up to 90-day horizon), returning point estimates with confidence intervals for display in the BI chart.

> **Addresses:** PRD §6.6, §6.8.11 — AI-powered forecasting integrated into trend charts.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User requests "forecast" in query or clicks "Extend Forecast" on trend chart |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `historical_data` | `[]TimeSeriesPoint` | A-01 result set | ✅ |
| `forecast_horizon_days` | `int` | User selection (max 90) | ✅ |
| `metric_name` | `string` | Semantic layer | ✅ |
| `seasonality_hints` | `[]string` | User config (`weekly / monthly / yearly`) | ⬜ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `forecast_points` | `[]ForecastPoint` | Date + predicted value + lower/upper CI |
| `model_quality` | `ModelQuality` | MAE, MAPE, R² on holdout set |
| `confidence_score` | `float64` | Forecast reliability |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Prophet (Python sidecar) | MIT | Time-series forecasting |
| 2 | Go HTTP client | Internal | Call Python forecast sidecar |
| 3 | Confidence scorer (A-06) | Internal | Assess forecast quality |

> **Note:** Prophet is Python-based. Runs as a lightweight sidecar REST service called from the Go backend.

---

## 6. Guardrails

- Minimum 30 historical data points required (return error if fewer).
- Maximum forecast horizon: 90 days.
- Forecasts labelled "AI Estimate" in UI — never presented as guaranteed.
- Confidence intervals always displayed alongside point estimates.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| MAPE on 30-day holdout | < 15% |
| Coverage of 80% CI | ≥ 80% |
| P95 Latency | < 10s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service + Python Prophet sidecar |
| **Sidecar** | `services/forecast-sidecar/` (FastAPI, Prophet, MIT) |
| **Depends on** | A-01 |
| **Consumed by** | BI trend charts, D-04 |
