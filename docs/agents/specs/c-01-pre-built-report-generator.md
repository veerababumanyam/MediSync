# Agent Specification — C-01: Pre-Built Report Generator Agent

**Agent ID:** `C-01`  
**Agent Name:** Pre-Built Report Generator Agent  
**Module:** C — Easy Reports  
**Phase:** 8  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Generates standard financial reports (P&L, Balance Sheet, Cash Flow Statement, Debtor Aging, Creditor Aging, Stock Summary, and 20+ others) from the data warehouse on demand, in PDF/Excel/HTML format.

> **Addresses:** PRD §6.8.1, US17, US22, US23 — One-click standardised financial reports.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User selects report type + period in Easy Reports UI |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `report_type` | `enum` | User selection | ✅ |
| `period` | `DateRange` | User selection | ✅ |
| `company_id` | `string` | Multi-entity selector | ✅ |
| `format` | `enum` | `pdf / xlsx / html / csv` | ✅ |
| `user_id` | `string` | JWT | ✅ |

### Supported Report Types (subset)
`profit_loss`, `balance_sheet`, `cash_flow`, `trial_balance`, `debtor_aging`, `creditor_aging`, `stock_summary`, `bank_reconciliation_summary`, `gst_summary`, `cost_centre_summary`, `budget_vs_actual`, `ledger_vouchers`

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `report_file` | `bytes` | Formatted report document |
| `report_metadata` | `ReportMeta` | Period, generated_at, agent_id, row_count |
| `download_url` | `string` | Signed URL (15-min TTL) |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Report template registry | Internal | SQL templates per report type |
| 2 | PostgreSQL (read-only) | PostgreSQL | Execute report SQL |
| 3 | OPA sidecar | Apache-2.0 | Row/column scope filter per role |
| 4 | `excelize` | BSD | Excel generation |
| 5 | PDF renderer (`chromedp` or `wkhtmltopdf`) | MIT / LGPL | PDF rendering from HTML template |

---

## 6. Guardrails

- All queries run via `medisync_readonly` Postgres role.
- Report data scoped to user's allowed entities and cost centres (OPA).
- Sensitive columns (salaries, cost prices) masked per role.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Report accuracy (vs manual calculation) | 100% |
| P95 Generation Latency | < 20s |
| Format rendering success rate | ≥ 99.9% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL, OPA, PDF renderer |
| **Consumed by** | Easy Reports UI, A-09 (scheduling) |
