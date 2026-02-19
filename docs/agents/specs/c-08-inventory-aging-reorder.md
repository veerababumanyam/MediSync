# Agent Specification — C-08: Inventory Aging & Reorder Agent

**Agent ID:** `C-08`  
**Agent Name:** Inventory Aging & Reorder Agent  
**Module:** C — Easy Reports  
**Phase:** 8  
**Priority:** P2 Medium  
**HITL Required:** Yes — for reorder recommendations  
**Status:** Draft

---

## 1. Purpose

Identifies slow-moving and obsolete inventory items, computes turnover ratios, and recommends reorder quantities and timing based on consumption rates. Reorder actions always require human approval before any Tally purchase order is raised.

> **Addresses:** PRD §6.5, §6.8.1, US22 — Inventory analysis with AI reorder suggestions.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual / Scheduled |
| **Manual trigger** | "Inventory Aging" report in Easy Reports |
| **Scheduled trigger** | `0 7 * * 1` (every Monday 7 AM) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `aging_threshold_days` | `int` | Config (default: 90 days) | ✅ |
| `reorder_point_config` | `ReorderConfig` | Tenant config | ✅ |
| `company_id` | `string` | Multi-entity | ✅ |
| `user_role` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `aging_items` | `[]InventoryItem` | Items with age bucket + turnover ratio |
| `reorder_recommendations` | `[]ReorderRec` | Item, suggested quantity, timing, supplier |
| `obsolete_items` | `[]InventoryItem` | Zero movement > threshold days |
| `report_file` | `bytes` | PDF/Excel export |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | PostgreSQL (read-only) | Fetch inventory + consumption data |
| 2 | Go aging engine | Compute turnover ratios + age buckets |
| 3 | A-12 (Prophet) | Forecast future consumption rate |
| 4 | Go reorder calculator | Compute EOQ (Economic Order Quantity) |
| 5 | D-08 (optional) | Prescriptive recommendations integration |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | Any reorder recommendation generated |
| **Notified role** | `pharmacy_manager` or `manager` |
| **Approval actions** | Approve / Modify quantity / Dismiss |
| **On approval** | Reorder event routed to D-10 (Insight-to-Action) → B-08 → B-09 for PO creation |

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Stockout rate after recommendations | Decreasing trend; < 2% |
| Reorder accuracy (vs actual consumption) | ≥ 90% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | PostgreSQL, A-12, D-10 (for approved reorders) |
| **Consumed by** | Pharmacy Manager, Operations Manager |
