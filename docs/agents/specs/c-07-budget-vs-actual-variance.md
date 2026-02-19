# Agent Specification — C-07: Budget vs. Actual Variance Agent

**Agent ID:** `C-07`  
**Agent Name:** Budget vs. Actual Variance Agent  
**Module:** C — Easy Reports  
**Phase:** 8  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Compares actual financial performance against budget targets, computes variance (absolute + %), forecasts year-end outturn, and flags departments or line items with material overages.

> **Addresses:** PRD §6.8.1, US23 — Budget vs actual comparison with variance analysis.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | "Budget vs Actual" report in Easy Reports |
| **Scheduled trigger** | `0 9 * * 1` (every Monday) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `period` | `DateRange` | User selection | ✅ |
| `budget_version` | `string` | Budget store (default: approved budget) | ✅ |
| `company_id` | `string` | Multi-entity selector | ✅ |
| `variance_threshold_pct` | `float64` | Config (default 10%) | ⬜ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `variance_lines` | `[]VarianceLine` | Per-line item: budget, actual, variance %, flag |
| `flagged_items` | `[]VarianceLine` | Items exceeding threshold |
| `yend_forecast` | `float64` | Projected year-end position |
| `report_file` | `bytes` | PDF/Excel export |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL (read-only) | Fetch actuals from warehouse |
| 2 | Budget store query | Fetch approved budget figures |
| 3 | Go variance engine | Compute variance + YE projection |
| 4 | A-12 forecast (optional) | Extend current trend for YE estimate |
| 5 | `excelize` / PDF renderer | Export |

---

## 6. Guardrails

- Read-only; no writes.
- Budget data access scoped to user's cost-centre (OPA).
- Flagged items display context (3-month trend) alongside variance.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Calculation accuracy | 100% |
| P95 Latency | < 15s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL, Budget store, A-12 (optional) |
| **Consumed by** | Finance Head, Dept Managers |
