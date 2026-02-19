# Agent Specification — B-12: Expense Categorisation Agent

**Agent ID:** `B-12`  
**Agent Name:** Expense Categorisation Agent  
**Module:** B — AI Accountant  
**Phase:** 7  
**Priority:** P2 Medium  
**HITL Required:** Yes  
**Status:** Draft

---

## 1. Purpose

Auto-assigns expense categories (utilities, travel, office supplies, medical supplies, etc.) to transactions based on vendor name, description, and learned rules from historical categorisations.

> **Addresses:** PRD §6.7.5 — Automated expense categorisation for cost analysis.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-05 (runs alongside ledger mapping) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `transaction` | `ExtractionResult` | B-02 | ✅ |
| `suggested_ledger` | `string` | B-05 | ✅ |
| `expense_categories` | `[]Category` | Config / Tally cost categories | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `expense_category` | `string` | Assigned category |
| `confidence_score` | `float64` | |
| `hitl_required` | `bool` | If confidence < 0.85 |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Rule engine (Go) | Exact vendor → category mappings |
| 2 | pgvector similarity | Historical category patterns |
| 3 | Genkit Flow (`expense-cat`) | LLM fallback for novel expenses |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `confidence < 0.85` |
| **Integrated into** | B-05 HITL review screen |

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Categorisation accuracy | ≥ 92% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (inline with B-05) |
| **Consumed by** | Cost analysis reports, B-08 |
