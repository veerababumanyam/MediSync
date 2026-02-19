# Agent Specification — A-07: Drill-Down Context Agent

**Agent ID:** `A-07`  
**Agent Name:** Drill-Down Context Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 3  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Interprets a user "click" on a chart element (e.g. clicking a bar labeled "January") and generates an appropriate drill-down SQL query that breaks the selected data point into finer-grained dimensions.

> **Addresses:** PRD §6.5 — Interactive drill-down from any chart cell or bar.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | `chart.click` UI event with element metadata |
| **Calling agent** | Frontend chart component → API |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `parent_sql` | `string` | A-01 original query | ✅ |
| `clicked_dimension` | `string` | Chart element label | ✅ |
| `clicked_value` | `any` | Chart element value | ✅ |
| `drill_level` | `int` | Current depth (max 5) | ✅ |
| `schema_context` | `JSON` | Schema Context Cache | ✅ |
| `user_role` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `drill_sql` | `string` | Postgres executor |
| `drill_title` | `string` | Chart title |
| `chart_type` | `enum` | A-03 router |
| `breadcrumb` | `[]string` | UI navigation path |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`drill-down-sql`) | LLM-based sub-query generation |
| 2 | SQL AST validator | SELECT-only assertion |
| 3 | PostgreSQL executor | Run drill SQL |
| 4 | A-03 Visualization Router | Pick chart type for result |

---

## 6. Guardrails

- Maximum drill depth: 5 levels (prevent infinite recursion).
- All drill-down queries subject to same OPA read-only + column masking policies as A-01.
- User's role-scope filter inherited from parent query.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Drill intent accuracy | ≥ 90% |
| P95 Latency | < 4s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | A-01, A-03 |
| **Consumed by** | Frontend chart component |
