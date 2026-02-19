# Agent Specification — A-02: SQL Self-Correction Agent

**Agent ID:** `A-02`  
**Agent Name:** SQL Self-Correction Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P0 Critical  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Detects SQL execution errors returned from PostgreSQL, analyses the error message and original query, rewrites the SQL, and retries — up to 3 times — before surfacing a failure to the user.

> **Addresses:** PRD §5.3 — Self-healing SQL retry loop for analyst queries.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (on DB execution error) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `original_query` | `string` | A-01 output | ✅ |
| `error_message` | `string` | Postgres error | ✅ |
| `schema_context` | `JSON` | Schema Context Cache | ✅ |
| `retry_count` | `int` | A-01 state | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `session_id` | `string` | Session store | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `corrected_sql` | `string` | Postgres executor |
| `correction_applied` | `string` | Explanation of fix |
| `success` | `bool` | A-01 |
| `trace_id` | `string` | Observability |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Genkit Flow (`sql-self-correct`) | Apache-2.0 | LLM error analysis + rewrite |
| 2 | Ollama / vLLM | MIT / Apache-2.0 | Local LLM inference |
| 3 | sqlparse AST validator | BSD | Confirm corrected SQL is still SELECT-only |
| 4 | PostgreSQL | PostgreSQL | Re-execute corrected query |

---

## 6. System Prompt

```
You are a SQL correction assistant. A query failed with the following error:

Error: {{ error_message }}

Original SQL:
{{ original_query }}

Database schema:
{{ schema_context }}

Fix the SQL so it executes correctly. Common fixes: correct table/column names,
fix JOIN conditions, resolve ambiguous column references.

RULES:
1. Output ONLY a corrected SELECT statement — no DML.
2. Do not change the business intent of the query.

OUTPUT FORMAT:
{"corrected_sql": "<fixed SELECT ...>", "correction_applied": "<brief explanation>"}
```

---

## 7. Guardrails

| # | Guard | Trigger | Action |
|---|-------|---------|--------|
| 1 | Max 3 retries | `retry_count >= 3` | Abort; return error to user |
| 2 | SQL AST read-only check | Every rewrite | Reject if DML introduced |
| 3 | Infinite loop guard | Same error repeated | Abort after identical error twice |

---

## 8. HITL Gate

Not applicable — A-02 operates transparently within A-01's execution loop.

---

## 9. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Correction success rate (error resolved in ≤3 retries) | ≥ 90% |
| False correction rate (SQL changes business intent) | < 2% |

---

## 10. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (embedded in A-01 flow) |
| **Depends on** | A-01 |
| **Consumed by** | A-01 |
