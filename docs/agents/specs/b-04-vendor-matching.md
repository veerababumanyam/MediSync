# Agent Specification — B-04: Vendor Matching Agent

**Agent ID:** `B-04`  
**Agent Name:** Vendor Matching Agent  
**Module:** B — AI Accountant  
**Phase:** 5  
**Priority:** P1 High  
**HITL Required:** Yes — when vendor is new/unrecognised  
**Status:** Draft

---

## 1. Purpose

Matches the extracted vendor name from a document to an existing vendor in the Tally vendor master using fuzzy matching and vector similarity. If no match is found, flags for human review (with an option to create a new vendor record).

> **Addresses:** PRD §6.7.2 — Vendor master matching and new vendor creation workflow.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-02 (post-extraction) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `extracted_vendor_name` | `string` | B-02 output | ✅ |
| `tally_vendor_master` | `[]Vendor` | Tally sync cache | ✅ |
| `vendor_embeddings` | `vector index` | pgvector | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `matched_vendor_id` | `*string` | Tally vendor ID if matched |
| `match_confidence` | `float64` | 0–1 |
| `match_method` | `enum` | `exact / fuzzy / vector / none` |
| `hitl_required` | `bool` | True for new/ambiguous vendors |
| `suggested_new_vendor` | `*Vendor` | Pre-populated for human to confirm |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Exact match (Go string) | Internal | O(1) lookup |
| 2 | rapidfuzz (via Python sidecar) | MIT | Fuzzy name matching |
| 3 | pgvector similarity search | PostgreSQL | Semantic vendor name matching |
| 4 | Genkit Flow (`vendor-match`) | Apache-2.0 | LLM reconciliation for ambiguous cases |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `match_method=none` (new vendor) OR `match_confidence < 0.80` |
| **Notified role** | `accountant` |
| **Approval actions** | Confirm match / Select different / Create new vendor |

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Match accuracy (known vendors) | ≥ 98% |
| New vendor false positive rate | < 2% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | B-02, pgvector, Tally vendor sync cache |
| **Consumed by** | B-05 |
