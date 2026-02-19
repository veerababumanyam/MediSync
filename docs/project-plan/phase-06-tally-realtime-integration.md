# Phase 06 — Tally Real-Time Integration

**Phase Duration:** Weeks 20–22 (3 weeks)  
**Module(s):** Module B (AI Accountant)  
**Status:** Planning  
**Milestone:** M5 — AI Bookkeeping (partial — full sync live)  
**Depends On:** Phase 05 complete (approval workflow operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §4.1](../ARCHITECTURE.md) | [ARCHITECTURE.md §8.3](../ARCHITECTURE.md)

---

## 1. Objectives

Complete the AI Accountant write-back loop: push approved transactions directly into Tally via TDL XML, build the real-time sync dashboard, add immutable audit logging, and support multiple Tally company instances. This phase delivers the **Action Plane** of MediSync — the only part of the system that writes to external financial data.

---

## 2. Scope

### In Scope
- B-09 Tally Sync Agent (TDL XML write-back)
- B-14 Audit Trail Logger Agent
- B-16 Multi-Entity Tally Manager Agent
- E-06 Multilingual Notification Agent
- Real-time Tally sync status dashboard
- OPA gate for Tally write operations
- Pre-sync validation (dupe detection, ledger availability)
- Post-sync rollback + error recovery
- Multi-entity (multi-company) Tally support
- Sync history log (last 100 sync events)
- Manual "Sync Now" button (finance_head only)

### Out of Scope
- Bank reconciliation (Phase 7)
- Easy Reports (Phase 8+)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | B-09 Tally Sync Agent | Backend + AI Engineer | Successfully posts approved journal entries/bills to Tally via TDL; pre-sync validation passes; post-sync entry verified in Tally |
| D-02 | OPA Tally Write Policy | DevOps | Only `finance_head` and `accountant_lead` can trigger sync; no self-approval; approved workflow required |
| D-03 | TDL XML Builder | Backend Engineer | Generates correct TDL XML for: journal entries, purchase bills, sales invoices, stock updates |
| D-04 | Pre-Sync Validation | Backend Engineer | Checks: ledger exists in Tally, no duplicate found, amount > 0, required fields present |
| D-05 | Post-Sync Verification | Backend Engineer | After TDL POST, queries Tally to confirm entry was created; matches posted amount |
| D-06 | B-14 Audit Trail Logger | AI Engineer | Immutable log entry for every sync: user, txn_id, timestamp, Tally voucher ID, source document ID |
| D-07 | B-16 Multi-Entity Tally Manager | AI Engineer | Switch between N Tally company instances; sync each independently; consolidated cross-entity view |
| D-08 | E-06 Multilingual Notification | AI Engineer | All approval and sync notifications delivered in user's locale (EN or AR) |
| D-09 | Sync Status Dashboard | Frontend Engineer | Live connection indicator; sync frequency config; last 100 sync events log; "Sync Now" button |
| D-10 | Sync History Log | Backend + Frontend | `app.tally_sync_history` table; UI shows: status, timestamp, vouchers created, errors |
| D-11 | Automatic Retry Logic | Backend Engineer | Failed syncs retry 3× with exponential backoff; failure after 3 retries → manual intervention alert |
| D-12 | Multi-Entity UI | Frontend Engineer | Dropdown to select/switch Tally company; per-entity sync status; consolidated dashboard |

---

## 4. AI Agents Deployed

### B-09 Tally Sync Agent

**Classification:** Reactive (L2) | HITL: Always (finance_head explicit click)

**Trigger:** `approval.completed` NATS event AND explicit user "Sync Now" / scheduled sync action.  
**Both conditions required** — approved workflow + explicit human action.

**Sync Pipeline:**
```
Finance Head clicks "Sync to Tally" (or scheduled sync trigger)
    │
    ▼ OPA gate check:
    │   - user.role must be finance_head OR accountant_lead
    │   - transaction.approval_status must be 'approved'
    │   - no self-approval in approval chain
    │   - IF any check fails → BLOCKED, user informed
    │
    ▼ Pre-sync validation (B-09 internal):
    │   1. Ledger ID exists in dim_ledgers (current Tally snapshot)
    │   2. SHA-256 duplicate check (same txn not already in Tally)
    │   3. Amount > 0 and <= configured limit
    │   4. Required fields: vendor, ledger, date, amount all present
    │   IF validation fails → BLOCKED, accountant informed with specific error
    │
    ▼ TDL XML generation (B-09):
    │   Generates Tally-format XML for entry type:
    │   - Purchase Bill: <ENVELOPE><BODY><IMPORTDATA>...</IMPORTDATA></BODY></ENVELOPE>
    │   - Sales Invoice: similarly structured
    │   - Journal Entry: <VOUCHER TYPE="Journal">
    │   - Payment Voucher: <VOUCHER TYPE="Payment">
    │
    ▼ HTTP POST to Tally web server (port 9000)
    │   Timeout: 60 seconds
    │   Content-Type: text/xml
    │
    ▼ Parse Tally XML response:
    │   Check for <LINEERROR> or <IMPORTRESULT><CREATED>
    │   Extract Tally voucher master ID
    │
    ▼ Post-sync verification:
    │   Query Tally for the created voucher; match amount + date
    │
    ▼ B-14: Write immutable audit log entry
    ▼ NATS: tally.sync.completed → Notification Dispatcher
    ▼ Update transaction_queue.approval_status = 'synced_to_tally'
```

**Error handling:**
- Tally unreachable: retry 3× with backoff (10s, 30s, 90s); after 3 failures → manual intervention alert
- TDL validation error: log error, mark transaction as `sync_failed`, notify accountant with Tally error message
- Duplicate entry in Tally: mark as `sync_failed_duplicate`, require human decision

**Rollback:** Tally does not support transactional rollback via TDL. If partial failure occurs in a batch, each entry's individual status tracked. Finance head can view failed entries and re-trigger individually.

### B-14 Audit Trail Logger Agent

**Type:** Reactive (L1)  
**Trigger:** Any of: `tally.sync.completed`, `approval.completed`, `document.reviewed`  
**Action:** Writes append-only entry to `app.audit_log`:

```json
{
  "log_id": "uuid",
  "event_type": "tally_sync",
  "user_id": "uuid",
  "user_role": "finance_head",
  "action": "synced_to_tally",
  "resource_type": "transaction",
  "resource_id": "txn_uuid",
  "tally_voucher_id": "TALLY-VCH-2026-0042",
  "source_doc_id": "doc_uuid",
  "amount": 15750.00,
  "ledger_id": "ledger_uuid",
  "ip_address": "192.168.1.100",
  "locale": "en",
  "created_at": "2026-02-19T10:30:00Z"
}
```

**Immutability enforcement:** PostgreSQL row-security policy rejects any UPDATE or DELETE on `audit_log`:
```sql
CREATE POLICY audit_log_immutable ON app.audit_log
    FOR ALL TO PUBLIC
    USING (true)
    WITH CHECK (false);
-- Only INSERT is permitted
```

### B-16 Multi-Entity Tally Manager Agent

**Purpose:** For clients running multiple Tally company instances (e.g., Clinic Branch A, Branch B, Pharmacy entity)

**Configuration:** `app.tally_entities` table stores:
- Entity name, Tally host:port, company name
- Last sync timestamp per entity
- Active/inactive status

**Capabilities:**
- Sync each entity independently on separate schedules
- Switch active entity context in UI (dropdown)
- Cross-entity consolidated read view (for BI queries)
- Entity-specific approval workflows

### E-06 Multilingual Notification Agent

**Trigger:** Any NATS event that triggers a user notification  
**Action:** Loads user's `locale` from `user_preferences`, selects matching template from `locales/en` or `locales/ar`, sends via Notification Dispatcher  
**Templates stored in:** `frontend/public/locales/{locale}/notifications.json`

**Notification types:**
- Approval request (notify next approver)
- Approval reminder (stale approval alert)
- Sync success (notify submitter)
- Sync failure (notify affected party + finance_head)
- ETL failure (notify admin)
- KPI alert (notify configured recipient)

---

## 5. TDL XML Examples

### Purchase Bill Entry
```xml
<ENVELOPE>
  <HEADER>
    <TALLYREQUEST>Import Data</TALLYREQUEST>
  </HEADER>
  <BODY>
    <IMPORTDATA>
      <REQUESTDESC>
        <REPORTNAME>Vouchers</REPORTNAME>
        <STATICVARIABLES>
          <SVCURRENTCOMPANY>{{tally_company_name}}</SVCURRENTCOMPANY>
        </STATICVARIABLES>
      </REQUESTDESC>
      <REQUESTDATA>
        <TALLYMESSAGE>
          <VOUCHER VCHTYPE="Purchase" ACTION="Create">
            <DATE>{{date_ddmmyyyy}}</DATE>
            <NARRATION>{{narration}}</NARRATION>
            <PARTYLEDGERNAME>{{vendor_ledger_name}}</PARTYLEDGERNAME>
            <ALLLEDGERENTRIES.LIST>
              <LEDGERNAME>{{gl_ledger_name}}</LEDGERNAME>
              <AMOUNT>-{{amount}}</AMOUNT>
            </ALLLEDGERENTRIES.LIST>
            <ALLLEDGERENTRIES.LIST>
              <LEDGERNAME>{{vendor_ledger_name}}</LEDGERNAME>
              <AMOUNT>{{amount}}</AMOUNT>
            </ALLLEDGERENTRIES.LIST>
          </VOUCHER>
        </TALLYMESSAGE>
      </REQUESTDATA>
    </IMPORTDATA>
  </BODY>
</ENVELOPE>
```

---

## 6. Database Schema Additions

```sql
CREATE TABLE tally_entities (
    entity_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_name     VARCHAR(255),
    tally_host      VARCHAR(255),
    tally_port      INTEGER DEFAULT 9000,
    company_name    VARCHAR(255),
    sync_schedule   VARCHAR(50) DEFAULT '*/30 * * * *',
    is_active       BOOLEAN DEFAULT TRUE,
    last_synced_at  TIMESTAMPTZ
);

CREATE TABLE tally_sync_history (
    sync_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_id       UUID REFERENCES tally_entities(entity_id),
    txn_id          UUID REFERENCES transaction_queue(txn_id),
    triggered_by    UUID REFERENCES users(user_id),
    trigger_type    VARCHAR(20),   -- 'manual' | 'scheduled'
    status          VARCHAR(20),   -- 'success' | 'failed' | 'partial'
    tally_voucher_id VARCHAR(255),
    error_message   TEXT,
    sync_duration_ms INTEGER,
    synced_at       TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 7. OPA Tally Write Policy

```rego
package medisync.tally

import future.keywords

# Allow Tally sync only for authorized roles
allow if {
    input.action == "tally_sync"
    input.user.role in ["finance_head", "accountant_lead"]
    input.transaction.approval_status == "approved"
    not self_approval(input)
}

# Block self-approval: user cannot approve their own submission
self_approval if {
    input.transaction.created_by == input.user.id
}

# Block sync if approval chain was not completed
deny[msg] if {
    input.action == "tally_sync"
    input.transaction.approval_status != "approved"
    msg := "Transaction must complete approval workflow before Tally sync"
}
```

---

## 8. Sync Status Dashboard

**Key UI components:**
- **Connection Indicator:** Green/Red/Yellow dot showing Tally connectivity (polling `/v1/sync/status` every 30 seconds)
- **Sync Frequency Setting:** Configurable dropdown (Real-time / Every 5 min / Every 15 min / Hourly / Manual only)
- **Last Sync Time:** Timestamp of last successful sync per entity
- **Sync History Feed:** Chronological list of last 100 sync events with: timestamp, user, vouchers created, status badge
- **Sync Now Button:** Manual trigger; only visible to `finance_head` and `accountant_lead`; requires confirmation dialog
- **Error Detail Drawer:** Click on failed sync event to see Tally error message + retry button

---

## 9. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| B-09 end-to-end sync | Post 10 approved transactions to Tally sandbox | 100% success rate; Tally vouchers verified |
| B-09 pre-sync validation | 5 invalid transactions (missing ledger, zero amount, duplicate) | All 5 blocked with correct error messages |
| B-09 OPA gate | Non-finance_head role attempt sync | Blocked |
| B-09 self-approval block | Finance head approves + syncs own submission | Blocked by OPA |
| B-09 retry on timeout | Simulate Tally timeout | 3 retries with backoff; alert fires after 3rd failure |
| B-14 audit log | Full approval + sync flow | Audit log contains all 5 state transitions; immutability verified |
| B-16 multi-entity | Configure 2 entities; sync independently | Sync history shows per-entity entries; no cross-contamination |
| E-06 multilingual notifications | Trigger all notification types for EN + AR users | Notifications arrive in correct language |

---

## 10. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| TDL XML structure differs between Tally ERP 9 vs TallyPrime | High | Test against both Tally versions; maintain version-specific XML templates |
| Tally network port not accessible from app server | High | Network firewall rules must be confirmed with client IT before Phase 6 begins |
| TDL batch size limits (Tally crashes on large XML payload) | Medium | Cap batch at 50 transactions per TDL POST; paginate large batches |
| Tally locks up during large sync | Medium | Sync during off-hours by default (configurable); monitor Tally response time |

---

## 11. Phase Exit Criteria

- [ ] B-09 Tally Sync Agent posting approved transactions to Tally successfully (tested against sandbox Tally instance)
- [ ] OPA gate blocking unauthorized sync attempts and self-approval
- [ ] B-14 audit log entries confirmed immutable (UPDATE/DELETE blocked)
- [ ] B-16 multi-entity support with 2+ Tally instances
- [ ] E-06 notifications delivering in correct locale for all notification types
- [ ] Sync status dashboard showing live connection status, sync history, manual trigger
- [ ] Phase gate reviewed and signed off

---

*Phase 06 | Version 1.0 | February 19, 2026*
