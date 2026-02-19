# Agent Specification — B-14: Audit Trail Logger Agent

**Agent ID:** `B-14`  
**Agent Name:** Audit Trail Logger Agent  
**Module:** B — AI Accountant  
**Phase:** 6  
**Priority:** P0 Critical  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Captures every mutation and significant event across the platform (who, what, when, source document) and writes immutable audit log entries. Provides the legal and compliance backbone for all financial operations.

> **Addresses:** PRD §6.7.6, US14 — Immutable audit trail for all AI Accountant actions.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | Called by every agent that performs financial operations |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `user_id` | `string` | Calling agent | ✅ |
| `user_role` | `string` | JWT | ✅ |
| `agent_id` | `string` | Calling agent | ✅ |
| `action` | `string` | Calling agent | ✅ |
| `resource_id` | `*string` | Transaction/document ID | ⬜ |
| `data_before` | `*JSON` | Previous state | ⬜ |
| `data_after` | `*JSON` | New state | ⬜ |
| `ip_address` | `net.IP` | API gateway | ✅ |
| `session_id` | `UUID` | Session store | ✅ |
| `trace_id` | `string` | OTel | ✅ |
| `status` | `enum` | `success / failed / blocked` | ✅ |
| `denial_reason` | `*string` | OPA | ⬜ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `audit_entry_id` | `UUID` | Stored entry ID |
| `data_hash` | `string` | SHA-256 of `data_after` |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | SHA-256 hasher (Go) | Compute tamper-evident hash |
| 2 | PostgreSQL INSERT | Append to `audit_log` (no UPDATE/DELETE rules) |

### Database Enforcement

```sql
-- Prevents any modification or deletion of audit records
CREATE RULE no_update_audit AS ON UPDATE TO audit_log DO INSTEAD NOTHING;
CREATE RULE no_delete_audit AS ON DELETE TO audit_log DO INSTEAD NOTHING;
```

---

## 6. Guardrails

- Write is synchronous — calling agent does not proceed until log is confirmed written.
- Errors in logging cause the parent operation to fail (no silent log-skip allowed).
- Log table is backed up daily to a separate encrypted store.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Log write success rate | 100% |
| Log write latency (P99) | < 50ms |
| Tamper detection (hash mismatch on read) | 100% |

---

## 8. Consumers

Called by: B-08, B-09, B-10, D-10, C-05, A-01 (PII queries), all OPA policy denials.

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go library (shared across all services — not a standalone HTTP service) |
| **DB table** | `audit_log` (append-only, no-update/delete rules) |
| **Retention** | 7 years (financial operations); 2 years (BI queries) |
