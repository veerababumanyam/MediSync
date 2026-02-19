# Agent Specification — B-09: Tally Sync Agent

**Agent ID:** `B-09`  
**Agent Name:** Tally Sync Agent  
**Module:** B — AI Accountant  
**Phase:** 6  
**Priority:** P0 Critical  
**HITL Required:** Yes — "Sync Now" is the **explicit human trigger**. No autonomous sync.  
**Status:** Draft

---

## 1. Purpose

Pushes fully approved journal entries, purchase bills, and sales invoices from MediSync into Tally ERP via TDL XML API. This is the **only agent with write access to Tally**. Every sync requires explicit human action after full B-08 approval.

> **Addresses:** PRD §6.7.3, US11 — One-click Tally sync for approved transactions.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "Sync Now" button — explicit user action in Accountant UI |
| **Calling agent** | B-08 (provides approved transaction bundle; human then clicks Sync) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `workflow_id` | `UUID` | B-08 | ✅ |
| `approved_transactions` | `[]Transaction` | B-08 approved bundle | ✅ |
| `tally_config` | `TallyConfig` | Encrypted secrets store | ✅ |
| `company_id` | `string` | Multi-entity selector | ✅ |
| `triggered_by` | `string` | User ID (must have `sync_to_tally` perm) | ✅ |
| `session_id` | `string` | Session | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `sync_result` | `enum` | `success / partial / failed` |
| `tally_voucher_ids` | `[]string` | Tally-assigned VoucherIDs |
| `failed_transactions` | `[]SyncError` | Transactions that failed with reason |
| `sync_timestamp` | `datetime` | UTC |
| `audit_log_entry_id` | `UUID` | Reference to B-14 log entry |

---

## 5. Tool Chain

```
ApprovedTransactions (from B-08) + Explicit User Sync Action
  → OPA policy gate:
      - user must have 'sync_to_tally' permission (finance_head or accountant_lead)
      - workflow_id.status must be 'approved'
      - self_approval check (OPA)
  → Pre-sync validation:
      - Duplicate guard: hash lookup in tally_sync_log
      - Ledger existence: verify all ledger names in Tally COA
      - Amount range validation
  → TDL XML Payload Generator (Go + html/template)
  → HTTPS POST to Tally Gateway (TLS 1.3)
  → Tally response parser: extract VoucherIDs or error codes
  → Auto-retry: 3× exponential backoff (2s, 4s, 8s) on timeout
  → B-14 Audit Log Writer (always, success or failure)
  → WebSocket push: Sync Status Dashboard update
  → Apprise notification (in-app + email)
```

### TDL XML Structure (Journal Entry)
```xml
<ENVELOPE>
  <HEADER><TALLYREQUEST>Import Data</TALLYREQUEST></HEADER>
  <BODY>
    <IMPORTDATA>
      <REQUESTDESC><REPORTNAME>Vouchers</REPORTNAME></REQUESTDESC>
      <REQUESTDATA>
        <TALLYMESSAGE xmlns:UDF="TallyUDF">
          <VOUCHER VCHTYPE="Journal" ACTION="Create">
            <DATE>{{ date }}</DATE>
            <NARRATION>{{ narration }}</NARRATION>
            <ALLLEDGERENTRIES.LIST>
              <LEDGERNAME>{{ debit_ledger }}</LEDGERNAME>
              <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
              <AMOUNT>-{{ amount }}</AMOUNT>
            </ALLLEDGERENTRIES.LIST>
            <ALLLEDGERENTRIES.LIST>
              <LEDGERNAME>{{ credit_ledger }}</LEDGERNAME>
              <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
              <AMOUNT>{{ amount }}</AMOUNT>
            </ALLLEDGERENTRIES.LIST>
          </VOUCHER>
        </TALLYMESSAGE>
      </REQUESTDATA>
    </IMPORTDATA>
  </BODY>
</ENVELOPE>
```

---

## 6. Guardrails

| # | Guard | Enforcement |
|---|-------|-------------|
| 1 | OPA policy: `sync_to_tally` permission | Hard block — 403 if fails |
| 2 | Workflow fully approved | Check `workflow.status == "approved"` before proceeding |
| 3 | Duplicate prevention | SHA-256 hash in `tally_sync_log`; block if exists |
| 4 | No autonomous sync | **Only executes on explicit human "Sync Now" click** |
| 5 | Audit log always written | Written before AND after sync attempt |
| 6 | Failed sync queue | Failed transactions → "Sync Failed" UI queue with retry option |

---

## 7. HITL Gate

The "Sync Now" button in the UI is the HITL gate:

| Property | Value |
|----------|-------|
| **Type** | Always — explicit user trigger |
| **Required role** | `finance_head` or `accountant_lead` |
| **Pre-conditions** | `workflow.status == approved` AND `OPA.allow == true` |

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Sync success rate | ≥ 99% |
| Duplicate sync rate | 0% (hash guard) |
| OPA unauthorised sync rate | 0% |
| P95 Tally round-trip latency | < 10s |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Tally integration** | TDL XML over HTTP (localhost:9000 default) |
| **Secrets** | Vault: `secret/medisync/tally/connection` |
| **DB table** | `tally_sync_log` (hash deduplication) |
| **Depends on** | B-08, OPA, B-14, Apprise |
| **Consumed by** | User (Finance Head / Accountant Lead) |
| **Env vars** | `TALLY_GATEWAY_URL`, `OPA_SIDECAR_URL`, `VAULT_ADDR` |
