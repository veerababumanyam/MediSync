# Agent Specification — B-11: Outstanding Items Agent

**Agent ID:** `B-11`  
**Agent Name:** Outstanding Items Agent  
**Module:** B — AI Accountant  
**Phase:** 7  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Generates outstanding-payments and outstanding-receipts reports from ledger data, classifying items into age buckets (0–7 days, 8–30 days, 30+ days) and surfacing overdue follow-up actions.

> **Addresses:** PRD §6.7.4 — Debtors and creditors ageing with automated follow-up.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output / Scheduled |
| **Calling agent** | B-10 (post-reconciliation) |
| **Scheduled trigger** | `0 7 * * *` (daily 7 AM) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `as_of_date` | `date` | System date or user input | ✅ |
| `ledger_scope` | `enum` | `payables / receivables / both` | ✅ |
| `user_role` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `outstanding_items` | `[]OutstandingItem` | Grouped by vendor/customer + age bucket |
| `total_payables` | `float64` | Sum of outstanding payables |
| `total_receivables` | `float64` | Sum of outstanding receivables |
| `overdue_items` | `[]OutstandingItem` | Items > 30 days |
| `report_file` | `*bytes` | PDF/Excel export |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL (read-only) | Fetch open ledger entries with dates |
| 2 | Go age-bucket calculator | Classify by days outstanding |
| 3 | `excelize` | Excel report generation |
| 4 | Apprise (optional) | Notify accountant of critical overdue items |

---

## 6. Guardrails

- Read-only: no writes to warehouse or Tally.
- Report data scoped to user's role-allowed ledgers (OPA).

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Report accuracy (matches manual check) | 100% |
| P95 Latency | < 10s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL read-only |
| **Consumed by** | Accountant UI, C-01 (reports) |
