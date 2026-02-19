# MediSync — Governance & Security Design

**Version:** 1.0 | **Created:** February 19, 2026  
**Cross-ref:** [PRD.md §10, §14](../PRD.md) | [00-agent-backlog.md](./00-agent-backlog.md)

---

## 1. The Core Contradiction — Resolved

The PRD contains an explicit conflict that must be architecturally resolved before any development begins:

> **PRD §10 (NFR):** "The AI must *never* have write/delete permissions" *(on the data warehouse)*
> 
> **PRD §6.7.3:** "Push approved transactions directly into Tally with a single click"

### Resolution Design

These two requirements are **compatible** when properly scoped:

| Constraint | Scope | Resolution |
|------------|-------|-----------|
| AI read-only | **Data Warehouse (PostgreSQL)** | AI agents connect with a `medisync_readonly` Postgres role that has GRANT SELECT only. No INSERT/UPDATE/DELETE is possible at the DB driver level. |
| Write to Tally | **Tally ERP via TDL XML API** | Tally sync is **not** a warehouse write — it is an integration write over HTTP/TDL. This path is controlled by a separate policy layer and **always requires explicit human approval** (B-08 → B-09 pipeline). |

**Architecture principle:** Separate the data intelligence plane (read-only warehouse) from the transactional action plane (Tally write-back). Connect them **only** through the human-approved workflow gate.

```
┌──────────────────────────────────────────┐
│          INTELLIGENCE PLANE              │
│  (AI Agents + Read-only Warehouse)       │
│  Postgres role: medisync_readonly        │
│  OPA policy: block all DML               │
└─────────────────────┬────────────────────┘
                      │ READ ONLY
                      ▼
              ┌───────────────┐
              │  Data Warehouse│
              │  (PostgreSQL)  │
              └───────────────┘
                      
┌──────────────────────────────────────────┐
│          ACTION PLANE                    │
│  (Human-approved Tally write operations) │
│  Requires: B-08 Approval + B-09 Sync     │
│  OPA policy: finance_head role only      │
└─────────────────────┬────────────────────┘
                      │ TDL XML (HTTP)
                      ▼
              ┌───────────────┐
              │   Tally ERP   │
              └───────────────┘
```

---

## 2. Authentication & Identity

### 2.1 Keycloak IAM (Apache-2.0)

All user authentication is handled by **Keycloak**:

- **SSO:** Single sign-on across all MediSync modules (BI Dashboard, AI Accountant, Easy Reports).
- **2FA:** TOTP-based 2FA mandatory for all Finance and Admin roles.
- **JWT Claims:** Every API request carries a signed JWT containing `user_id`, `roles[]`, `department`, `cost_centres[]`.
- **Session Management:** 8-hour session tokens; refresh tokens valid for 30 days; force-logout on role change.

### 2.2 Service Accounts

Each agent service has its own Keycloak service account with minimum required scope:

| Service Account | Permissions |
|----------------|------------|
| `sa-bi-agent` | `read:warehouse`, `read:semantic_layer` |
| `sa-ocr-agent` | `read:documents`, `write:extraction_queue` |
| `sa-mapping-agent` | `read:tally_coa`, `read:chroma_embeddings` |
| `sa-approval-agent` | `read:pending_transactions`, `write:approval_events` |
| `sa-tally-sync` | `execute:tally_sync` (Finance Head role + OPA gate) |
| `sa-scheduler` | `read:warehouse`, `write:notification_queue` |

---

## 3. Authorization — Policy as Code (OPA)

### 3.1 OPA Policy Architecture

All authorization decisions are delegated to **Open Policy Agent (OPA)** running as a sidecar service:

```
Agent API Request
  → API Gateway (FastAPI)
  → OPA Sidecar: POST /v1/data/medisync/authz/allow
      → Policy bundle (Rego rules)
      → User context (JWT claims)
      → Request context (resource, action)
  → Decision: allow / deny
  → If deny: return 403 + reason
```

### 3.2 Core Rego Policies

```rego
# medisync/policies/bi_readonly.rego
# Enforce read-only access to the data warehouse for all BI agents

package medisync.bi

default allow = false

# Allow SELECT queries only
allow {
    input.action == "query"
    input.query_type == "SELECT"
    not contains_dml(input.sql)
}

contains_dml(sql) {
    dml_keywords := ["INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "TRUNCATE", "ALTER"]
    keyword := dml_keywords[_]
    contains(upper(sql), keyword)
}
```

```rego
# medisync/policies/tally_sync.rego
# Only finance_head or accountant_lead may trigger Tally sync
# Self-approval is blocked

package medisync.tally

default allow = false

allow {
    input.action == "sync_to_tally"
    user_has_role(input.user, "finance_head")
    workflow_approved(input.workflow_id)
    not self_approved(input.user, input.workflow_id)
}

allow {
    input.action == "sync_to_tally"
    user_has_role(input.user, "accountant_lead")
    workflow_approved(input.workflow_id)
    not self_approved(input.user, input.workflow_id)
}

user_has_role(user, role) {
    role == user.roles[_]
}

workflow_approved(wf_id) {
    data.approval_workflows[wf_id].status == "approved"
}

self_approved(user, wf_id) {
    data.approval_workflows[wf_id].submitted_by == user.id
    data.approval_workflows[wf_id].last_approver == user.id
}
```

```rego
# medisync/policies/row_level_security.rego
# Users only see data relevant to their department / cost-centre

package medisync.data

default row_filter = {}

row_filter = filter {
    filter := {"department": input.user.department}
    not user_has_role(input.user, "finance_head")
    not user_has_role(input.user, "admin")
}

# Finance head and admin see all data
row_filter = {} {
    user_has_role(input.user, "finance_head")
}
```

### 3.3 Role Definitions

| Role | Description | Can Query BI | Can Upload Docs | Can Approve Transactions | Can Sync to Tally |
|------|-------------|:---:|:---:|:---:|:---:|
| `admin` | System administrator | ✅ | ✅ | ✅ | ✅ |
| `finance_head` | CFO / Owner | ✅ (all) | ✅ | ✅ | ✅ |
| `accountant_lead` | Senior accountant | ✅ (dept) | ✅ | ✅ | ✅ |
| `accountant` | Accountant | ✅ (dept) | ✅ | 1st-level only | ❌ |
| `manager` | Clinic / Department manager | ✅ (dept) | ❌ | ❌ | ❌ |
| `pharmacy_manager` | Pharmacy lead | ✅ (pharmacy) | ❌ | ❌ | ❌ |
| `analyst` | Data analyst | ✅ (all) | ❌ | ❌ | ❌ |
| `viewer` | Read-only stakeholder | ✅ (limited) | ❌ | ❌ | ❌ |

---

## 4. Data Security

### 4.1 Encryption

| Layer | Method | Notes |
|-------|--------|-------|
| Data at rest (Postgres) | AES-256 (pgcrypto / Postgres TDE) | Applied to patient PII columns and financial amounts |
| Data in transit | TLS 1.3 (HTTPS) | All API endpoints, DB connections, Tally TDL calls |
| Uploaded documents | AES-256 at object storage level | Documents linked via encrypted references |
| JWT secrets | RS256 (asymmetric key) | Keys rotated every 90 days via Keycloak |
| LLM context window | No PII sent to external LLMs | On-premise Ollama / vLLM for sensitive data; PII stripped before any cloud LLM call |

### 4.2 Column-Level Masking

Sensitive columns masked based on role via OPA + Postgres views:

| Column | Visible to | Masked for |
|--------|-----------|-----------|
| `patient.name`, `patient.dob`, `patient.id` | admin, finance_head | analyst, viewer, manager |
| `ledger.cost_price` | finance_head, accountant | manager, viewer, pharmacy_manager |
| `employee.salary` | admin, finance_head | all others |
| `tally.vendor.bank_account` | finance_head, accountant_lead | all others |

Implementation: Postgres **row security policies** + OPA-generated `column_blacklist` injected into SELECT queries by the query layer.

### 4.3 Patient Data (HIPAA / GDPR Compliance)

- Patient identifiers (name, DOB, NHI) are **never** passed to the LLM in plain text.
- The Text-to-SQL agent anonymises patient references in generated SQL (`WHERE patient_id = ?` not `WHERE patient_name = 'John Smith'`).
- Audit logs for any query touching patient data are retained for 7 years.
- Data retention policy: Operational data 7 years; aggregated analytics 10 years; patient PII follows applicable healthcare data law.

---

## 5. Audit Trail Design

### 5.1 Audit Log Schema (PostgreSQL — append-only)

```sql
CREATE TABLE audit_log (
    id           UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    event_time   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id      UUID NOT NULL,
    user_role    TEXT NOT NULL,
    agent_id     TEXT NOT NULL,          -- e.g. 'B-09', 'A-01'
    action       TEXT NOT NULL,          -- e.g. 'tally_sync', 'sql_query'
    resource_id  TEXT,                   -- e.g. transaction ID, document ID
    data_before  JSONB,                  -- previous state (for mutations)
    data_after   JSONB,                  -- new state (for mutations)
    data_hash    TEXT,                   -- SHA-256 of data_after
    ip_address   INET,
    session_id   UUID,
    trace_id     TEXT,                   -- OpenTelemetry trace ID
    status       TEXT NOT NULL,          -- 'success' | 'failed' | 'blocked'
    denial_reason TEXT                   -- OPA policy reason if blocked
);

-- Prevent modification and deletion
CREATE RULE no_update_audit AS ON UPDATE TO audit_log DO INSTEAD NOTHING;
CREATE RULE no_delete_audit AS ON DELETE TO audit_log DO INSTEAD NOTHING;
```

### 5.2 What is Logged

| Event Category | Logged By | Retention |
|---------------|-----------|-----------|
| All BI queries (A-01) | A-01 agent | 2 years |
| Document uploads | B-01, B-02 | 7 years |
| Ledger mapping suggestions | B-05 | 7 years |
| All approval actions | B-08 | 7 years |
| Tally sync events | B-09 | 7 years |
| All report views/exports | C-01, C-03, C-05 | 2 years |
| Policy denials (OPA) | OPA + API gateway | 7 years |
| User login/logout | Keycloak | 2 years |
| Admin role changes | Keycloak | 7 years |

---

## 6. Hallucination & Reliability Guardrails

### 6.1 AI Read-Only Boundary (BI Agents)

```
Network Policy: medisync_readonly Postgres role
  GRANT SELECT ON ALL TABLES IN SCHEMA public TO medisync_readonly;
  REVOKE INSERT, UPDATE, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA public FROM medisync_readonly;
```

The `medisync_readonly` role is enforced at the **database level** — not just in application code. Even if an agent generates an INSERT statement, the database will reject it.

### 6.2 Confidence-Based Routing

```
confidence_score >= 0.95  → Return result directly
confidence_score 0.70–0.94 → Return result with warning banner "Please verify"
confidence_score < 0.70    → Hold result; notify user "Under review"; queue for human validation
```

### 6.3 Off-Topic Deflection (A-05)

System prompt anchor: "I can only answer questions about your business data. For general knowledge questions, please use a general-purpose assistant."

Implementation: A lightweight classifier (fine-tuned DistilBERT, MIT license) runs before the main SQL agent to detect off-topic queries and short-circuit with a canned response. LLM is not invoked for off-topic queries (saves cost).

### 6.4 SQL Safety Assertions

Before any generated SQL is executed:
1. Parse SQL AST (sqlparse library, BSD license).
2. Assert statement type == SELECT.
3. Assert no subquery DML.
4. Assert no `INFORMATION_SCHEMA` access (meta-table restriction).
5. Assert query targets only whitelisted schemas.

On assertion failure: log to audit trail, return error to user, never execute.

---

## 7. Tally Write-Back — Full Security Flow

This is the highest-risk operation in MediSync. The complete security chain:

```
User clicks "Approve" (B-08)
  ↓
OPA check: user.role ∈ [finance_head, accountant_lead] AND workflow.status != 'self_approved'
  ↓ (allow)
Approval event written to audit_log
  ↓
All approvers in chain have signed off
  ↓
User clicks "Sync Now" (B-09 trigger — explicit human action required)
  ↓
OPA check: user.role == finance_head AND workflow_id.status == 'fully_approved'
  ↓ (allow)
Pre-sync validation:
  - Duplicate check (hash lookup)
  - Ledger existence check
  - Amount range validation
  ↓ (pass)
TDL XML payload generated (Jinja2 template)
  ↓
HTTPS POST to Tally Gateway (TLS 1.3)
  ↓
Parse Tally response: VoucherID or error
  ↓
Write sync result to audit_log (success or failure)
  ↓
Notify user via Apprise (in-app + email)
```

**What cannot happen:**
- ❌ AI agent does not trigger sync autonomously.
- ❌ Sync cannot proceed without fully-approved workflow.
- ❌ A user cannot approve their own submissions.
- ❌ Duplicate transactions cannot be synced twice (hash guard).

---

## 8. Data Quality as a Security Layer

Bad data is as dangerous as a security breach in a financial system. The following controls apply:

| Control | Tool | Trigger |
|---------|------|---------|
| ETL validation (C-06) | great_expectations | Every Airflow DAG run |
| Duplicate invoice detection (B-07) | Postgres hash lookup | On every document upload |
| Pre-sync validation (B-09) | Custom validator | Before every Tally sync |
| Reconciliation check | B-10 | Daily scheduled run |
| Anomaly detection | A-13, D-07 | Hourly scheduled scan |

---

## 9. Key Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| LLM generates malicious SQL via prompt injection | Medium | High | SQL AST assertion; read-only DB role; OPA policy |
| Incorrect ledger mapping → wrong Tally entry | High | High | HITL review required for all mappings; approval chain |
| OCR errors → wrong amounts in Tally | High | High | Confidence threshold HITL gate; pre-sync validation |
| Unauthorised Tally sync | Low | Critical | OPA hard policy; audit log; 2FA on finance roles |
| Patient PII exposure via LLM | Low | Critical | PII stripped from context; on-premise LLM for PHI queries |
| Duplicate invoice posting | Medium | High | SHA-256 hash deduplication; B-07 pre-check |
| Audit log tampering | Low | Critical | Postgres no-update/delete rules; append-only table |

---

## 10. Compliance Checklist

| Standard | Requirement | Implementation |
|----------|------------|----------------|
| HIPAA | PHI access control | Keycloak RBAC + OPA column masking |
| HIPAA | Audit trails | 7-year audit_log retention |
| HIPAA | Encryption at rest | pgcrypto AES-256 on PHI columns |
| GDPR | Right to access | User data export endpoint |
| GDPR | Right to erasure | Flagged for legal review (7-year financial retention may override) |
| GST/Tax | Audit-ready reports | B-13 tax compliance agent + audit_log |
| SOX-like | Separation of duties | Self-approval blocked by OPA; 3-tier approval chain |
| SOX-like | Immutable audit trail | No-update/delete Postgres rules |
