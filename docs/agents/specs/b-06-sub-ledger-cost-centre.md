# Agent Specification — B-06: Sub-Ledger & Cost Centre Assignment Agent

**Agent ID:** `B-06`  
**Agent Name:** Sub-Ledger & Cost Centre Assignment Agent  
**Module:** B — AI Accountant  
**Phase:** 5  
**Priority:** P2 Medium  
**HITL Required:** Yes  
**Status:** Draft

---

## 1. Purpose

Suggests sub-ledger and cost centre assignments for transactions based on context (department, vendor category, project codes) and historical assignment patterns.

> **Addresses:** PRD §6.7.2 — Sub-ledger and cost centre allocation.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-05 (after ledger suggestion) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `transaction` | `ExtractionResult` | B-02 | ✅ |
| `suggested_ledger` | `string` | B-05 output | ✅ |
| `tally_cost_centres` | `[]CostCentre` | Tally sync cache | ✅ |
| `user_department` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `suggested_sub_ledger` | `*string` | Sub-ledger name |
| `suggested_cost_centre` | `*string` | Cost centre |
| `confidence_score` | `float64` | |
| `hitl_required` | `bool` | True when multiple cost centres plausible |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | pgvector similarity | Match historical sub-ledger + cost centre patterns |
| 2 | Genkit Flow (`cost-centre-assign`) | LLM assignment |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `confidence < 0.85` OR multiple cost centres plausible |
| **Notified role** | `accountant` |
| **Integrated into** | B-05 HITL review screen (shown alongside ledger suggestion) |

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Cost centre accuracy | ≥ 90% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (runs inline with B-05) |
| **Depends on** | B-05, Tally COA cache |
| **Consumed by** | B-08 |
