# Agent Specification — B-13: Tax Compliance Agent

**Agent ID:** `B-13`  
**Agent Name:** Tax Compliance Agent  
**Module:** B — AI Accountant  
**Phase:** 7  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Computes GST/VAT input-credit, output-tax, and net tax liability per reporting period from warehouse data. Generates compliance-ready tax reports for filing.

> **Addresses:** PRD §6.7.6, US26 — GST/tax computation and compliance reporting.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | "Generate Tax Report" button |
| **Scheduled trigger** | `0 9 1 * *` (1st of every month, 9 AM) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `period` | `DateRange` | User selection | ✅ |
| `company_id` | `string` | Multi-entity selector | ✅ |
| `tax_config` | `TaxConfig` | Tenant config (GST rates, exemptions) | ✅ |
| `user_role` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `input_credit` | `float64` | Total GST input credit |
| `output_tax` | `float64` | Total GST output liability |
| `net_liability` | `float64` | `output_tax - input_credit` |
| `line_items` | `[]TaxLineItem` | Breakdown by ledger/vendor |
| `report_file` | `bytes` | PDF/Excel compliance report |
| `filing_ready` | `bool` | All required fields populated |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL (read-only) | Aggregate tax amounts by ledger type and period |
| 2 | Go tax calculator | Apply GST rates, compute input/output/net |
| 3 | `excelize` | Excel export |
| 4 | PDF renderer | PDF report generation |

---

## 6. Guardrails

- Read-only: no Tally writes.
- Tax rates stored in versioned config — rate changes logged for historical accuracy.
- Generated reports do not auto-submit to tax authorities — user must review and file manually.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Tax calculation accuracy (vs manual spot-check) | 100% |
| Report generation P95 Latency | < 15s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL read-only, tenant tax config |
| **Consumed by** | Finance Head, Accountant |
