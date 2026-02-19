# Agent Specification — C-04: Custom Metric Formula Agent

**Agent ID:** `C-04`  
**Agent Name:** Custom Metric Formula Agent  
**Module:** C — Easy Reports  
**Phase:** 10  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Interprets user-defined formula expressions (e.g. "Gross Margin % = (Revenue - COGS) / Revenue * 100") and registers them as governed custom metrics in the Semantic Layer, without requiring SQL coding from the user.

> **Addresses:** PRD §6.8.3, §6.9.5 — No-code custom KPI creation.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Create Custom Metric" in Easy Reports UI |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `metric_name` | `string` | User input | ✅ |
| `formula_expression` | `string` | User input (natural language or formula syntax) | ✅ |
| `existing_metrics` | `[]Metric` | Semantic Layer Registry | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `metric_definition` | `MetricDef` | Validated metric config for Semantic Layer |
| `generated_sql` | `string` | SQL fragment for the metric |
| `validation_result` | `bool` | Test execution passed |
| `metric_id` | `UUID` | Registered metric ID |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`formula-to-metric`) | Parse expression → MetricFlow metric definition |
| 2 | SQL validator | Validate generated SQL |
| 3 | PostgreSQL (read-only) | Test execute metric on sample data |
| 4 | Semantic Layer Registry | Register validated metric |

---

## 6. Guardrails

- Generated SQL must be SELECT-only (OPA).
- Test execution runs on limited sample (max 1000 rows) before registration.
- Metric names must be unique within the tenant.
- Only `analyst`, `finance_head`, `admin` roles can create custom metrics.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Formula parse success rate | ≥ 95% |
| SQL correctness (test execution passes) | ≥ 99% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | D-09 Semantic Layer Registry, PostgreSQL |
| **Consumed by** | D-09, A-01, D-01 |
