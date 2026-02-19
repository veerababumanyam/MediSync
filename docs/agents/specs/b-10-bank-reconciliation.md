# Agent Specification — B-10: Bank Reconciliation Agent

**Agent ID:** `B-10`  
**Agent Name:** Bank Reconciliation Agent  
**Module:** B — AI Accountant  
**Phase:** 7  
**Priority:** P1 High  
**HITL Required:** Yes — unmatched items always routed to human  
**Status:** Draft

---

## 1. Purpose

Automatically matches bank statement entries to existing Tally ledger entries using a 4-tier confidence matching strategy (exact → amount → fuzzy → unmatched); surfaces unmatched items for manual resolution; feeds B-11 for outstanding reports.

> **Addresses:** PRD §6.7.4, US12 — Automated bank reconciliation with exception handling.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | User uploads bank statement |
| **Scheduled trigger** | `0 6 * * *` (daily 6 AM — auto-reconcile yesterday's transactions) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `bank_statement` | `[]BankRow` | B-02 (OCR) or CSV upload | ✅ |
| `tally_ledger_entries` | `[]LedgerEntry` | Warehouse (read-only) | ✅ |
| `reconciliation_date_range` | `DateRange` | User or schedule config | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `matched` | `[]ReconciliationMatch` | Confirmed matches |
| `suggested` | `[]ReconciliationMatch` | Probable matches requiring review |
| `unmatched_bank` | `[]BankRow` | No ledger match found |
| `unmatched_ledger` | `[]LedgerEntry` | In Tally but not in bank |
| `reconciliation_score` | `float64` | % of items matched |
| `hitl_required` | `bool` | True if any unmatched |

---

## 5. Matching Algorithm

```
For each BankRow:
  Tier 1 — Exact:   amount EXACT + date ±1 day + description similarity > 0.85  → Matched (high, ≥0.95)
  Tier 2 — Amount:  same amount + date within 7 days                              → Suggested (0.80–0.94)
  Tier 3 — Fuzzy:   amount ±₹10 + similar description                             → Possible (0.60–0.79)
  Tier 4 — None:    no match found                                                 → Unmatched (manual)
```

---

## 6. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | B-02 | Internal | OCR bank statement (if scanned) |
| 2 | PostgreSQL windowed query | PostgreSQL | Find candidate ledger entries by date range + amount |
| 3 | rapidfuzz (Python sidecar) | MIT | Description similarity scoring |
| 4 | Genkit Flow (`recon-split`) | Apache-2.0 | LLM handles multi-invoice partial payment splits |
| 5 | B-11 trigger | Internal | Generate outstanding items report |

---

## 7. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | Any `tier=unmatched` OR `tier=possible` items exist |
| **Notified role** | `accountant` |
| **Approval actions** | Confirm match / Select alternative / Mark as outstanding |

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Tier 1+2 auto-match rate | ≥ 85% of transactions |
| False match rate (wrong ledger paired) | < 1% |
| P95 Latency (100-row statement) | < 30s |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | B-02, PostgreSQL, rapidfuzz sidecar |
| **Consumed by** | B-11, Reconciliation dashboard |
