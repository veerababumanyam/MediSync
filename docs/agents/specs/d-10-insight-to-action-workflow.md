# Agent Specification — D-10: Insight-to-Action Workflow Agent

**Agent ID:** `D-10`  
**Agent Name:** Insight-to-Action Workflow Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 14  
**Priority:** P1 High  
**HITL Required:** Yes — always  
**Status:** Draft

---

## 1. Purpose

Bridges the gap between an AI recommendation and an actual business action. Receives accepted recommendations and routes them to the appropriate downstream agent or human workflow for execution, with mandatory human approval at every step.

> **Addresses:** PRD §6.9.6, US28 — Closing the insight-to-action loop.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Event trigger** | User accepts a D-08 recommendation |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `recommendation` | `Recommendation` | D-08 | ✅ |
| `accepted_by` | `string` | User (JWT) | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `workflow_id` | `UUID` | Initiated workflow ID |
| `routed_to` | `string` | Target agent or team |
| `status` | `enum` | `pending_approval / approved / rejected / completed` |

### Action Registry (subset)

| Action ID | Description | Routes to |
|-----------|-------------|-----------|
| `create_purchase_order` | Raise PO in Tally | B-08 → B-09 |
| `raise_expense_approval` | Trigger expense approval chain | B-08 |
| `generate_report` | Create and distribute a report | C-01 → C-03 |
| `create_budget_revision` | Propose budget amendment | Finance Head |
| `escalate_to_manager` | Route insight to manager with context | Apprise notification |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Action registry lookup | Resolve `action_ref` → route |
| 2 | B-08 Approval Workflow (if writing) | Route through approval chain |
| 3 | B-14 Audit Log | Log insight → action linkage |
| 4 | Apprise | Notify relevant stakeholder |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | Every action initiation |
| **Notified role** | Role matching the action type |
| **Approval actions** | Approve / Modify / Reject |
| **On approve** | Action routed to downstream agent |
| **On reject** | Workflow closed; logged with reason |

> **Design intent:** D-10 never executes actions directly. It is a routing and oversight layer.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Action completion rate (approved → completed) | ≥ 95% |
| HITL approval turnaround | < 2 business hours avg |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | B-08, B-09, B-14, C-01, Apprise |
| **Consumed by** | D-08, C-08 |
