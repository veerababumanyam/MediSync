# Agent Specification — B-05: Ledger Mapping Agent

**Agent ID:** `B-05`  
**Agent Name:** Ledger Mapping Agent  
**Module:** B — AI Accountant  
**Phase:** 5  
**Priority:** P0 Critical  
**HITL Required:** Yes — always (review before any Tally sync)  
**Status:** Draft

---

## 1. Purpose

Maps each extracted transaction to the most appropriate Tally GL ledger using vector similarity on historical mappings and LLM classification. Provides confidence-tiered suggestions and learns from user corrections.

> **Addresses:** PRD §6.7.2, US10 — AI-powered ledger suggestion with continuous improvement.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-02 / B-04 (post-extraction, post-vendor-match) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `transaction` | `ExtractionResult` | B-02/B-04 | ✅ |
| `tally_chart_of_accounts` | `[]Ledger` | Tally sync cache | ✅ |
| `historical_mappings` | `vector index` | pgvector (learned corrections) | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `session_id` | `string` | Session | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `suggested_ledger` | `string` | Primary ledger name |
| `suggested_sub_ledger` | `*string` | Sub-ledger (if applicable) |
| `suggested_cost_centre` | `*string` | Cost centre assignment |
| `confidence_score` | `float64` | 0.0–1.0 |
| `confidence_badge` | `enum` | `high (≥0.95) / review (0.70–0.94) / manual (<0.70)` |
| `alternative_mappings` | `[]Mapping` | Top 3 alternatives with scores |
| `reasoning` | `string` | LLM explanation |
| `hitl_required` | `bool` | Always true |
| `trace_id` | `string` | OTel trace |

---

## 5. Tool Chain

```
Transaction (from B-02/B-04)
  → Embedding generator (BAAI/bge-small, Apache-2.0)
  → pgvector similarity search (top-5 historical mappings)
  → Context: [transaction + historical + chart of accounts]
  → Genkit Flow (ledger-classify)
  → LLM classification + reasoning
  → Confidence badge assigner
  → MappingResult struct
  → On user approval: store embedding + confirmed ledger in pgvector (feedback loop)
```

---

## 6. System Prompt

```
You are a financial ledger mapping expert. Given this transaction and the available Tally ledgers,
suggest the most appropriate GL ledger assignment.

Transaction:
{{ transaction_json }}

Historical similar mappings (most relevant first):
{{ historical_mappings }}

Available Tally ledgers:
{{ chart_of_accounts }}

RULES:
1. Choose the ledger that best matches the transaction's nature and vendor.
2. Prefer historical patterns if confidence is high (>0.90).
3. Provide top 3 alternative mappings with individual confidence scores.
4. Explain your reasoning briefly.

OUTPUT: Valid JSON matching MappingResult schema.
```

---

## 7. Guardrails

| # | Guard | Action |
|---|-------|--------|
| 1 | Suggestions only | Agent never writes to Tally — only suggests |
| 2 | HITL always required | No transaction proceeds without human review |
| 3 | Bulk rule application | Requires `finance_head` role |
| 4 | Audit log | All suggestions + user decisions logged |

---

## 8. HITL Gate

| Property | Value |
|----------|-------|
| **Gate type** | Always |
| **Notified role** | `accountant` (or submitter) |
| **Confidence badge** | `high` → green (still requires review); `review` → amber; `manual` → red |
| **SLA** | 48h |
| **On approval** | Transaction passed to B-08 approval workflow |
| **Feedback loop** | User corrections saved to vector store for future use |

---

## 9. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Mapping accuracy (top-1 suggestion) | ≥ 90% |
| User correction rate after learning loop | Decreasing trend; < 20% after 3 months |
| P95 Latency | < 5s |
| `manual` badge rate | < 10% |

---

## 10. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Vector store** | pgvector (`ledger_mapping_embeddings` table) |
| **Embedding model** | BAAI/bge-small (Apache-2.0) via sidecar |
| **Depends on** | B-02, B-04, pgvector, Tally COA cache |
| **Consumed by** | B-08 (approval workflow) |
