# Agent Specification — D-12: Federated Query Optimisation Agent

**Agent ID:** `D-12`  
**Agent Name:** Federated Query Optimisation Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 16  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Optimises cross-source queries (Tally + HIMS + bank + inventory) by analysing query plans, recommending or applying query rewrites, managing materialised views, and routing queries to the right data source for minimum latency.

> **Addresses:** PRD §6.9.9 — Query performance and multi-source federation.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | Any multi-source SQL query; scheduled materialised view refresh |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `query` | `string` | A-01 or D-01 | ✅ |
| `query_plan` | `string` | PostgreSQL EXPLAIN output | ✅ |
| `source_metadata` | `[]SourceMeta` | Source registry | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `optimised_query` | `string` | Rewritten query |
| `source_routes` | `[]SourceRoute` | Which data fetched from which source |
| `estimated_latency_ms` | `int` | Estimated execution time |
| `mat_view_suggestions` | `[]MatViewSuggestion` | Recommended materialised views |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL EXPLAIN ANALYSE | Query plan analysis |
| 2 | Genkit Flow (`query-opt`) | LLM-assisted query rewrite (complex cases) |
| 3 | Go query planner | Rule-based routing |
| 4 | Redis | Query result caching |

---

## 6. Guardrails

- Optimised queries must produce identical results to originals (validated by test execution).
- No schema modifications without explicit admin approval.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Query latency improvement | ≥ 30% avg |
| Query correctness after rewrite | 100% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go middleware |
| **Depends on** | PostgreSQL, Redis |
| **Consumed by** | A-01, D-01, all SQL-generating agents |
