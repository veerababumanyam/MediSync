# Agent Specification — D-05: Deep Research Agent

**Agent ID:** `D-05`  
**Agent Name:** Deep Research Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 14  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Executes multi-source, multi-hop research tasks on a question by combining internal data (warehouse, documents) with structured reasoning. Enables deep analytical questions like "What is the root cause of the 22% drop in pharmacy revenue this quarter?"

> **Addresses:** PRD §6.9.2 — Deep analytical research with chain-of-thought reasoning.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Deep Research" button in Analytics Chat |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `research_question` | `string` | User input | ✅ |
| `session_id` | `string` | Frontend | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `time_budget_s` | `int` | Config (default: 120s) | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `research_report` | `ResearchReport` | Structured markdown report + evidence |
| `conclusions` | `[]string` | Top-3 conclusions with confidence |
| `evidence_chain` | `[]EvidenceItem` | Per-SQL evidence supporting conclusions |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`deep-research`) | Orchestrate multi-hop reasoning |
| 2 | D-02 | Entity extraction |
| 3 | A-01 + A-02 + A-07 | Iterative SQL query generation + drill-down |
| 4 | pgvector | Semantic document search |
| 5 | A-05 | Validate conclusions |
| 6 | A-06 | Confidence scoring |

---

## 6. Guardrails

- Time-bounded: hard abort at `time_budget_s` with partial results.
- Max 20 SQL queries per research task.
- All evidence SQLs run read-only.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Research conclusion accuracy | ≥ 88% |
| P95 Completion Time | < 120s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | D-02, A-01, A-02, A-05, A-06, A-07, pgvector |
| **Consumed by** | Analytics Chat (Deep Research mode) |
