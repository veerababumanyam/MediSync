# Agent Specification — D-01: Natural Language Search Agent

**Agent ID:** `D-01`  
**Agent Name:** Natural Language Search Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 13  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Translates free-text search queries (e.g. "show revenue from Apollo hospital last quarter") into structured lookups across the data warehouse, combining full-text search, entity recognition, and Semantic Layer metadata to return instant results.

> **Addresses:** PRD §6.9.1 — Conversational search entry point for the Advanced Analytics module.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User types in the global search bar |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `raw_query` | `string` | User input | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `session_id` | `string` | Frontend | ✅ |
| `context_entities` | `[]Entity` | D-02 output | ⬜ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `search_intent` | `SearchIntent` | Classified intent type |
| `route` | `enum` | `sql_query / entity_lookup / report / document` |
| `result` | `SearchResult` | Result set + suggested viz |
| `suggested_followups` | `[]string` | 3 suggested follow-up questions |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | D-02 | Entity extraction from query |
| 2 | D-09 Semantic Layer | Term-to-metric mapping |
| 3 | A-01 Text-to-SQL | SQL path for analytical queries |
| 4 | pgvector | Semantic document search path |
| 5 | Genkit Flow (`nl-search-router`) | Route to correct sub-agent |
| 6 | A-05 Hallucination Guard | Validate result before display |

---

## 6. Guardrails

- All results filtered through C-05 row/column security.
- No write access.
- Intent classification confidence < 0.6 → fallback to A-01.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Search result relevance (nDCG@5) | ≥ 0.85 |
| P95 Latency | < 3s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | D-02, D-09, A-01, pgvector |
| **Consumed by** | Global search bar (Web + Mobile) |
