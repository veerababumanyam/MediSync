# Agent Specification — D-02: Entity Recognition Agent

**Agent ID:** `D-02`  
**Agent Name:** Entity Recognition Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 13  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Extracts named entities from user queries and documents — patients, vendors, ledger heads, cost centres, date ranges, amounts — and resolves them to canonical IDs in the master data store.

> **Addresses:** PRD §6.9.1 — Entity-aware search and analytical query processing.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Event trigger** | Any query or document from D-01, D-03, A-01 |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `text` | `string` | Query / document excerpt | ✅ |
| `context` | `string` | Domain context hint | ⬜ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `entities` | `[]Entity` | Extracted entities with type, value, canonical_id, confidence |
| `resolved_count` | `int` | Successfully resolved count |
| `unresolved` | `[]string` | Entities that couldn't be resolved |

### Entity Types
`PATIENT`, `VENDOR`, `LEDGER_HEAD`, `COST_CENTRE`, `DATE_RANGE`, `AMOUNT`, `LOCATION`, `DOCTOR`, `DEPARTMENT`

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`ner-extract`) | LLM-based NER with domain prompt |
| 2 | Synonym registry (`config/synonym_registry.yaml`) | Alias resolution |
| 3 | pgvector similarity | Fuzzy canonical ID resolution |
| 4 | Master data lookup (PostgreSQL) | Exact ID resolution |

---

## 6. Guardrails

- Entity lookup scoped to user's accessible entities (OPA).
- Confidence < 0.7 → entity marked `unresolved`; query continues without it.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Entity extraction F1 score | ≥ 0.92 |
| Resolution accuracy | ≥ 95% |
| P95 Latency | < 500ms |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | Genkit, pgvector, synonym registry |
| **Consumed by** | D-01, D-03, A-01 |
