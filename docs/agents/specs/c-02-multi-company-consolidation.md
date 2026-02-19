# Agent Specification — C-02: Multi-Company Consolidation Agent

**Agent ID:** `C-02`  
**Agent Name:** Multi-Company Consolidation Agent  
**Module:** C — Easy Reports  
**Phase:** 8  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Merges financial statements from multiple Tally company instances into a single consolidated view, eliminating intercompany transactions and producing group-level P&L, Balance Sheet, and Cash Flow statements.

> **Addresses:** PRD §6.8.2, US20 — Multi-company consolidated financial reporting.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Consolidated View" selection in Easy Reports |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `entity_ids` | `[]string` | Tenant config (all entities to consolidate) | ✅ |
| `period` | `DateRange` | User selection | ✅ |
| `currency` | `string` | Tenant config (base currency for consolidation) | ✅ |
| `format` | `enum` | `pdf / xlsx / html` | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `consolidated_pl` | `FinancialStatement` | Group P&L |
| `consolidated_bs` | `FinancialStatement` | Group Balance Sheet |
| `intercompany_eliminations` | `[]Elimination` | Transactions eliminated with detail |
| `report_file` | `bytes` | Formatted consolidated report |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL (read-only, per entity schema) | Fetch each entity's financial data |
| 2 | Go consolidation engine | Sum, eliminate intercompany, format |
| 3 | `excelize` / PDF renderer | Export |

---

## 6. Guardrails

- Accessible only to `admin` and `finance_head` roles.
- Intercompany eliminations fully logged for audit.
- Currency conversion rates sourced from a pinned config (not external API) for audit reproducibility.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Consolidation accuracy | 100% |
| P95 Latency (5 entities, 1 year) | < 60s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | B-16, PostgreSQL per-entity schemas |
| **Consumed by** | Finance Head, Group CFO |
