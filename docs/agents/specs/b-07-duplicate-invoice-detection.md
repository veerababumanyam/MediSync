# Agent Specification — B-07: Duplicate Invoice Detection Agent

**Agent ID:** `B-07`  
**Agent Name:** Duplicate Invoice Detection Agent  
**Module:** B — AI Accountant  
**Phase:** 5  
**Priority:** P0 Critical  
**HITL Required:** Yes — flags duplicate for human confirmation before blocking  
**Status:** Draft

---

## 1. Purpose

Compares every incoming invoice against existing records using a combination of exact hash matching and fuzzy heuristics (same vendor + similar amount + similar date). Flags potential duplicates before they proceed to the approval workflow.

> **Addresses:** PRD §6.7.2, §6.7.5, US16 — Prevent duplicate invoice postings in Tally.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-02 (immediately after extraction) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `extracted_invoice` | `ExtractionResult` | B-02 | ✅ |
| `dedup_window_days` | `int` | Config (default 90) | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `is_duplicate` | `bool` | Definite duplicate detected |
| `is_suspicious` | `bool` | Fuzzy match — needs review |
| `matching_records` | `[]InvoiceRef` | Existing records that match |
| `match_type` | `enum` | `exact_hash / fuzzy / none` |
| `hitl_required` | `bool` | True if suspicious |
| `duplicate_hash` | `string` | SHA-256 of (vendor + amount + date) |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | SHA-256 hasher (Go) | Internal | Exact duplicate fingerprint |
| 2 | PostgreSQL lookup | PostgreSQL | Check `invoice_hashes` table |
| 3 | rapidfuzz (Python sidecar) | MIT | Fuzzy match for near-duplicates |

### Duplicate Hash Definition
```
SHA-256(normalize(vendor_name) + round(amount, 2) + invoice_date.format("YYYY-MM"))
```

---

## 6. Guardrails

- `is_duplicate=true` (exact hash match) → **hard block**: invoice cannot proceed without Finance Head override.
- `is_suspicious=true` (fuzzy match) → HITL required; accountant must confirm it is not a duplicate.
- All duplicate checks logged to audit_log.

---

## 7. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `is_duplicate=true` OR `is_suspicious=true` |
| **Notified role** | `accountant` + `finance_head` (for hard duplicates) |
| **Approval actions** | Confirm duplicate (discard) / Confirm not duplicate (proceed) / Finance Head override |

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Exact duplicate detection rate | 100% |
| False positive (legitimate re-invoice flagged) | < 1% |
| P95 Latency | < 500ms |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (runs inline with B-02) |
| **DB table** | `invoice_hashes` (indexed on hash) |
| **Depends on** | B-02 |
| **Consumed by** | B-05 (only if not duplicate) |
