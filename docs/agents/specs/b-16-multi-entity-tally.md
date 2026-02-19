# Agent Specification — B-16: Multi-Entity Tally Manager Agent

**Agent ID:** `B-16`  
**Agent Name:** Multi-Entity Tally Manager Agent  
**Module:** B — AI Accountant  
**Phase:** 6  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Manages connections to multiple Tally company instances (for healthcare chains or multi-entity groups). Switches company context per entity, syncs and consolidates statements independently, and aggregates cross-entity dashboards.

> **Addresses:** PRD §6.7.3, US15, US20 — Multi-entity Tally management and cross-entity reporting.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | Entity switcher in AI Accountant UI |
| **Scheduled trigger** | ETL DAG per entity (Airflow/Go scheduler) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `entity_ids` | `[]string` | Tenant config | ✅ |
| `operation` | `enum` | `sync / consolidate / switch` | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `active_entity` | `string` | Currently selected entity |
| `sync_statuses` | `map[string]SyncStatus` | Per-entity sync status |
| `consolidated_financials` | `*ConsolidatedReport` | Cross-entity aggregated statements |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Tally Gateway pool | Separate TDL connections per entity |
| 2 | B-09 (per entity) | Entity-scoped Tally sync |
| 3 | Go consolidation engine | Merge + eliminate intercompany transactions |

---

## 6. Guardrails

- Each entity connection authenticated separately with its own Tally credentials.
- Cross-entity data visible only to `admin` and `finance_head` roles.
- Intercompany eliminations logged for audit.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Entity switch P95 Latency | < 2s |
| Consolidation accuracy | 100% (verified against manual consolidation) |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Secrets** | Vault: `secret/medisync/tally/{entity_id}/connection` |
| **Depends on** | B-09, Tally Gateway per-entity |
| **Consumed by** | Finance Head, Admin |
