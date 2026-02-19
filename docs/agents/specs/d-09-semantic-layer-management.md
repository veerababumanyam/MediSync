# Agent Specification — D-09: Semantic Layer Management Agent

**Agent ID:** `D-09`  
**Agent Name:** Semantic Layer Management Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 13  
**Priority:** P0 Critical  
**HITL Required:** Yes — governance approval for metric changes  
**Status:** Draft

---

## 1. Purpose

Manages the Semantic Layer — the governed dictionary of business metrics, dimensions, hierarchies, and relationships. Acts as the single source of truth that translates business terms to SQL fragments used by all BI and analytics agents.

> **Addresses:** PRD §6.9.5 — Centralised semantic layer as the BI data contract.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Event-driven |
| **Manual trigger** | Admin registers or modifies a metric definition |
| **Event trigger** | C-04 creates a validated custom metric |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `metric_definition` | `MetricDef` | C-04 or Admin UI | ✅ |
| `approver_id` | `string` | Governance workflow | ✅ (for changes) |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `metric_registry` | `[]MetricDef` | Full current metric registry snapshot |
| `metric_id` | `UUID` | Registered metric ID |
| `validation_result` | `ValidationResult` | SQL validation result |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL | Metric registry store (`semantic_layer.metrics`) |
| 2 | SQL validator (Go) | Validate metric SQL fragment |
| 3 | B-08 Approval Workflow | Route metric change for governance approval |
| 4 | B-14 Audit Log | Log all registry changes immutably |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | Any modification to an existing metric definition |
| **Notified role** | `finance_head` or `admin` |
| **Approval actions** | Approve / Reject / Modify |
| **On reject** | Previous version retained; change discarded |

---

## 7. Guardrails

- Metric SQL fragments must be SELECT-only.
- All versions of metric definitions stored (never deleted) for audit.
- Breaking changes (renamed metrics) surfaced with impact analysis.

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Registry availability | ≥ 99.9% |
| Stale metric (definition/SQL mismatch) count | 0 |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL, B-08, B-14 |
| **Consumed by** | A-01, A-04, D-01, D-06, C-04, all SQL-generating agents |
