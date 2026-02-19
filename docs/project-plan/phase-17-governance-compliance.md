# Phase 17 — Data Governance & Compliance

**Phase Duration:** Weeks 51–52 (2 weeks)  
**Module(s):** All (cross-cutting)  
**Status:** Planning  
**Depends On:** Phase 16 (All agents and APIs live; embedding surface complete)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Harden MediSync for regulated healthcare environments: enforce HIPAA and GDPR compliance controls across all 58 agents and all data pathways, extend row/column security to the search and analytics layer (Module D), implement comprehensive audit logging for every data access event system-wide, enforce sensitive-question whitelisting to prevent unauthorised PII exposure through the chat interface, and activate the full compliance reporting suite for audit evidence.

---

## 2. Scope

### In Scope
- HIPAA Business Associate Agreement (BAA) controls implemented in-product
- GDPR rights implementation: right to erasure, right of access, data portability
- Search-level access control (OPA policies applied to D-01/D-03/D-12 search results)
- Comprehensive audit logging for all data access (Module A, B, C, D, E)
- Sensitive question whitelisting (chat + search)
- C-05 Row/Column Security extended to Module D search layer
- D-09 governance tier: metric change audit trail
- Data retention policies enforced (automated PII purge)
- Compliance report suite (HIPAA evidence, GDPR DSR log, access activity report)
- Privacy notice integration (consent capture and audit)

### Out of Scope
- External regulatory submission tooling (handled by customer compliance team)
- Clinical data DLP beyond patient record masking (out of MediSync scope)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | Search-level OPA access control | AI Eng + Security | Unauthorised user receives 0 rows from D-01/D-03/D-12 beyond their permission scope |
| D-02 | Comprehensive audit log (all modules) | Backend Eng | Every data access (read + write) logged with: user, timestamp, IP, resource, result count |
| D-03 | Sensitive question whitelist + PII guard | AI Eng | PII-exposing queries blocked unless user has PII role; E-05 CI check |
| D-04 | GDPR Right to Erasure | Backend Eng | Verified deletion of all PII for a given patient across all schemas in ≤ 24h |
| D-05 | GDPR Right of Access | Backend Eng | User data export (JSON/PDF) containing all stored data for a person |
| D-06 | Data retention enforcement | Data Eng + Backend | Automated purge jobs: PII retention per policy; audit log retention 7 years |
| D-07 | HIPAA control implementation | Backend + Security | 12 HIPAA technical safeguard requirements checked off |
| D-08 | Compliance report suite | Frontend + Backend | 4 compliance reports generated on demand: access log, PII exposure, consent audit, DSR log |
| D-09 | Privacy consent capture | Frontend | Consent banner; consent stored in `app.consents`; opt-out honoured |
| D-10 | Penetration test round 2 | Security | 0 P0/P1 findings on full system (all 58 agents in scope) |

---

## 4. OPA Policy Extensions

### Search-Level Access Control

Extended from C-05 to cover Module D search. Policies applied post-query at the result-row level:

```rego
# OPA policy: search_access.rego
package medisync.search

import rego.v1

# Block search results where entity_type = 'patient' unless user has PII role
allow_result if {
    input.result.entity_type != "patient"
}

allow_result if {
    input.result.entity_type == "patient"
    "pii_viewer" in input.user.roles
}

# Block search results for cost_centre user cannot access
allow_result if {
    input.result.entity_type != "transaction"
}

allow_result if {
    input.result.entity_type == "transaction"
    input.result.cost_centre in input.user.allowed_cost_centres
}
```

**OPA integration point:** D-01 NL Search and D-12 Federated Query both evaluate `medisync.search.allow_result` on each result item before returning to client.

---

### Sensitive Question Guard

```rego
# OPA policy: chat_pii_guard.rego
package medisync.chat

import rego.v1

# Questions that trigger PII guard
pii_question_patterns := [
    ".*patient.*name.*",
    ".*patient.*address.*",
    ".*salary.*",
    ".*employee.*personal.*",
    ".*staff.*contact.*",
]

blocked if {
    some pattern in pii_question_patterns
    regex.match(pattern, lower(input.query))
    not "pii_viewer" in input.user.roles
}

# If blocked, return: 403 with explanation, not the answer
response := {
    "blocked": true,
    "reason": "This query would expose personally identifiable information. Your current role does not have PII viewer access. Contact your system administrator."
} if blocked
```

**Applies to:** A-01 (chat query) and D-01 (NL search).

---

## 5. Comprehensive Audit Log

**Target:** Every read and write data access event across all modules creates an immutable audit record.

**Audit event schema:**
```sql
-- Extended app.audit_log (existing table extended)
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  event_module TEXT;               -- 'A', 'B', 'C', 'D', 'E', 'system'
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  resource_type TEXT;              -- 'query', 'document', 'report', 'transaction', 'search'
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  resource_id TEXT;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  result_row_count INT;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  ip_address INET;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  user_agent TEXT;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  pii_accessed BOOLEAN DEFAULT FALSE;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  blocked BOOLEAN DEFAULT FALSE;
ALTER TABLE app.audit_log ADD COLUMN IF NOT EXISTS
  block_reason TEXT;
```

**Audit events captured:**
| Module | Events |
|---|---|
| A | Every chat query; every chart render; every drill-down |
| B | Document upload; OCR extraction; approval decision; Tally sync |
| C | Report generation; report download; schedule creation |
| D | Every search query; every research job; every recommendation action |
| E | Translation used; OCR language detected; notification sent |
| System | Login; logout; role change; API key creation; token revocation |

**Retention:** 7 years (HIPAA requirement: 6 years; GDPR: as long as data exists)  
**Immutability:** `audit_log` rows have no `UPDATE`/`DELETE` permission for any app role — enforced by PostgreSQL row security + OPA.

---

## 6. HIPAA Technical Safeguards Checklist

| # | Requirement | Implementation |
|---|---|---|
| 1 | Unique user identification | Keycloak UID per user; shared accounts prohibited via policy |
| 2 | Emergency access procedure | Break-glass admin role with mandatory audit log + alert |
| 3 | Automatic logoff | Session timeout 30 min (configurable); Keycloak idle timeout |
| 4 | Encryption in transit | TLS 1.3 enforced on all connections |
| 5 | Encryption at rest | AES-256 on PostgreSQL tablespace + document storage |
| 6 | Audit controls | Comprehensive audit log (see above) |
| 7 | Integrity controls | Immutable audit log; document SHA-256 hash verification |
| 8 | User authentication | Keycloak OIDC with 2FA mandatory for PHI-access roles |
| 9 | Transmission security | TLS 1.3; mutual TLS for inter-service (NATS) |
| 10 | Access control (role-based) | OPA C-05 + search-level row/column security |
| 11 | Minimum necessary | Column masking ensures users see only what their role requires |
| 12 | PHI de-identification | PII columns can be pseudonymised for analytics; raw only with pii_viewer role |

---

## 7. GDPR Implementation

### Right to Erasure (Article 17)

**Trigger:** Verified DSR (Data Subject Request) from patient or staff member

**Erasure pipeline:**
```
DSR request submitted (app.dsr_requests)
      │ Verified by DPO/admin
      ▼
Erasure job created (background worker)
      │
      ├─ hims_analytics.dim_patients → pseudonymise (hash name, NULLify contact fields)
      ├─ app.documents → delete associated patient documents
      ├─ tally_analytics → patient invoices pseudonymised (voucher preserved for finance)
      ├─ vectors.search_documents → delete patient search index entries
      ├─ app.search_history → delete rows matching patient queries
      └─ app.audit_log → retain (legal obligation) but redact PII fields
      │
      ▼
Erasure completion report generated (PDF)
Stored in app.dsr_requests.completion_report
```

**Completion target:** ≤ 24 hours from verified request.

### Right of Access (Article 15)

**Output:** Full JSON export or structured PDF of all stored data for the subject.  
**Includes:** Chat history involving subject, audit log entries, documents, transactions.

### Data Portability (Article 20)

**Output:** Machine-readable JSON export (standard format)

---

## 8. Data Retention Policies

| Data Category | Retention Period | Action |
|---|---|---|
| Patient PII (clinical records) | Per facility HIPAA policy (default 7 years) | Archive → purge |
| Financial transaction records | 7 years (UAE commercial law) | Retain in pseudonymised form |
| Chat/search query logs | 2 years | Purge (PII scrubbed, aggregates retained) |
| Audit log entries | 7 years | Retain (immutable) |
| Uploaded documents (invoices) | 7 years (UAE commercial) | Retain |
| Uploaded documents (clinical) | Per facility HIPAA policy | Retain / pseudonymise |
| Scheduled report history | 2 years | Purge |

**Automated purge jobs:** Cron-based background workers run nightly; log purge activity to audit trail.

---

## 9. Compliance Report Suite

| Report | Contents | Audience |
|---|---|---|
| Access Activity Report | All data access events by user, date range, resource type | Compliance Officer, IT Admin |
| PII Exposure Report | All queries/searches where PII was accessed; user; timestamp | DPO, Legal |
| Consent Audit | Consent capture/revocation events; current consent status | DPO |
| DSR Log | All Data Subject Requests; status; completion date; erasure confirmation | DPO, Legal |

**Generation:** On demand via `/v1/compliance/reports/{type}`; output: PDF + JSON  
**Access control:** `compliance_officer`, `dpo`, `admin` roles only

---

## 10. Database Schema Additions

```sql
-- GDPR Data Subject Requests
CREATE TABLE app.dsr_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_type TEXT NOT NULL,       -- 'erasure', 'access', 'portability', 'objection'
  subject_type TEXT NOT NULL,       -- 'patient', 'staff'
  subject_id TEXT NOT NULL,         -- patient_id or staff_id
  subject_name TEXT,
  requestor_name TEXT,
  verified_by UUID REFERENCES app.users(id),
  verified_at TIMESTAMPTZ,
  status TEXT DEFAULT 'pending',    -- 'pending', 'verified', 'processing', 'completed', 'rejected'
  completion_report BYTEA,          -- PDF confirmation
  completed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Consent records
CREATE TABLE app.consents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES app.users(id),
  consent_type TEXT NOT NULL,       -- 'analytics', 'notifications', 'data_sharing'
  granted BOOLEAN NOT NULL,
  ip_address INET,
  user_agent TEXT,
  granted_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ
);
```

---

## 11. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/v1/audit/events` | Query audit log (admin/compliance only) |
| `POST` | `/v1/dsr` | Submit Data Subject Request |
| `PATCH` | `/v1/dsr/{id}/verify` | Verify DSR (admin) |
| `GET` | `/v1/dsr/{id}/status` | Check DSR status |
| `GET` | `/v1/compliance/reports/access-activity` | Access Activity Report |
| `GET` | `/v1/compliance/reports/pii-exposure` | PII Exposure Report |
| `GET` | `/v1/compliance/reports/consent` | Consent Audit Report |
| `GET` | `/v1/compliance/reports/dsr-log` | DSR Log Report |
| `POST` | `/v1/users/{id}/consent` | Capture user consent |

---

## 12. Testing Requirements

| Test | Target |
|---|---|
| OPA search-level access control | User without PII role: 0 patient records in search results |
| OPA sensitive question guard | 20 PII-probing queries: all blocked for non-PII role |
| Audit log completeness | 10 test actions across all modules: all logged with correct metadata |
| Audit log immutability | Attempt row delete/update on audit_log: PostgreSQL error returned |
| Right to erasure | Test patient erasure: all PII fields cleared across all tables in < 24h |
| HIPAA checklist | All 12 safeguards verified by security team sign-off |
| Compliance reports | All 4 reports generate correctly for test date range |
| Data retention purge | Purge job correctly archives/deletes records beyond retention period |
| Pen test round 2 | 0 P0/P1 findings across full system |

---

## 13. Phase Exit Criteria

- [ ] Search-level OPA access control enforced (D-01, D-03, D-12)
- [ ] Sensitive question guard active (A-01 + D-01)
- [ ] Comprehensive audit log covering all 58 agents + all data access events
- [ ] Audit log immutable (no app-level DELETE/UPDATE permitted)
- [ ] GDPR Right to Erasure pipeline: completes in < 24 hours
- [ ] GDPR Right of Access: data export functional
- [ ] 12 HIPAA technical safeguards verified and documented
- [ ] Data retention policies active with automated purge jobs
- [ ] Compliance report suite: all 4 reports functional
- [ ] Pen test round 2: 0 P0/P1 findings
- [ ] Phase gate signed off by compliance officer or designated authority

---

*Phase 17 | Version 1.0 | February 19, 2026*
