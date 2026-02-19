# Phase 01 — ETL & Infrastructure

**Phase Duration:** Weeks 1–3 (3 weeks)  
**Module(s):** Foundation (all modules depend on this)  
**Status:** Planning  
**Milestone:** M1 — Data Foundation  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md#42-etl--ingestion-layer) | [PRD.md §12](../PRD.md)

---

## 1. Objectives

Establish the complete data foundation that all subsequent phases depend on. This phase creates the on-premises data warehouse, ETL pipelines from Tally ERP and HIMS, the first AI agent (C-06 Data Quality Validation), and the base infrastructure for all services.

**No user-facing features are delivered in this phase.** The output is a reliable, validated, continuously-syncing data pipeline.

---

## 2. Scope

### In Scope
- On-premises server provisioning and OS hardening
- PostgreSQL 18.2 data warehouse with all schemas
- Tally TDL connector and HIMS REST connector
- Meltano ELT pipelines with incremental sync
- Data validation agent (C-06) with quality reporting
- Keycloak identity server (initial setup)
- NATS message broker (initial setup)
- Redis cache (initial setup)
- CI/CD pipeline scaffolding
- i18n file structure and E-05 CI translation guard
- Base Docker Compose configuration

### Out of Scope
- AI query agents (Phase 2)
- Frontend UI (Phase 3)
- Document processing (Phase 4)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | PostgreSQL 18.2 DW with all schemas | Data Engineer | All schemas (`hims_analytics`, `tally_analytics`, `semantic`, `app`, `vectors`) created with correct tables, indexes, and constraints |
| D-02 | Tally TDL XML Connector (Go) | Backend Engineer | Incremental sync of Ledgers, Vouchers, Inventory Masters, Stock Items running reliably; `_synced_at` populated |
| D-03 | HIMS REST API Connector (Go) | Backend Engineer | Incremental sync of Patients, Appointments, Billing, Pharmacy Dispensations running reliably |
| D-04 | Meltano ELT Pipelines | Data Engineer | Both connectors registered as Meltano taps; pipelines run on 15–30 min schedule without errors |
| D-05 | Go Transform Service | Backend Engineer | Schema normalisation, dedup checks, `ON CONFLICT DO UPDATE` upserts working; failed records written to `etl_quarantine` |
| D-06 | C-06 Data Quality Validation Agent | AI Engineer | Missing value checks, referential integrity, duplicate detection, anomaly alerts via NATS `etl.sync.failed` topic |
| D-07 | Keycloak Initial Setup | DevOps | OIDC realm configured; admin, analyst, viewer demo roles provisioned; JWT structure matching spec |
| D-08 | NATS / JetStream Setup | DevOps | Broker running; `etl.sync.completed` and `etl.sync.failed` topics live; JetStream persistence enabled |
| D-09 | Redis Cache Server | DevOps | Redis running, bound to localhost, AUTH configured; schema context TTL 1h policy defined |
| D-10 | CI/CD Pipeline | DevOps | GitHub Actions (or equivalent) with Go build, test, lint, E-05 translation coverage check |
| D-11 | i18n File Structure | Frontend Engineer | `frontend/public/locales/en/` and `ar/` directories with placeholder JSON for all namespaces; CI fails if `ar` key missing |
| D-12 | Docker Compose (dev) | DevOps | Single `docker-compose up` launches: Postgres + pgvector, Redis, NATS, Keycloak, ETL services |

---

## 4. AI Agents Deployed

| Agent | ID | Type | Description |
|---|---|---|---|
| Data Quality Validation | C-06 | Proactive (L2) | Runs after every ETL sync; checks missing values, referential integrity, duplicates; emits NATS alert on failure |

### C-06 Agent Detail

**Trigger:** NATS subscription to `etl.sync.completed`  
**Tools:**
- SQL query executor (read-only, `medisync_readonly` role)
- NATS publisher (emit alerts)
- Notification dispatcher (write to `app.notification_queue`)

**Quality Checks Performed:**
1. Missing value rate per column — alert if > 5% null on required fields
2. Referential integrity — `fact_appointments.patient_id` must exist in `dim_patients`
3. Duplicate detection — hash of business keys; flag exact dupes
4. Value range anomalies — negative amounts, dates in future, zero-quantity dispensations
5. Row count delta — alert if sync row count drops > 30% vs previous sync (source system issue indicator)

**Output:**
- `app.etl_quality_report` table row per sync run
- NATS event `etl.data.quality.alert` when any check fails above threshold
- Grafana `etl_quality_score` metric updated

---

## 5. Data Warehouse Schema

### 5.1 Schema: `hims_analytics`

```sql
-- Dimension Tables
dim_patients       (patient_id PK, name_en, name_ar, dob, gender, phone, ...)
dim_doctors        (doctor_id PK, name_en, name_ar, specialty, department, ...)
dim_drugs          (drug_id PK, name_en, name_ar, category, unit, ...)

-- Fact Tables
fact_appointments  (appt_id PK, patient_id FK, doctor_id FK, appt_date, status, 
                    department, billing_id, _synced_at, _source_id)
fact_billing       (bill_id PK, patient_id FK, amount, tax_amount, payment_mode,
                    bill_date, department, _synced_at, _source_id)
fact_pharmacy_disp (disp_id PK, drug_id FK, patient_id FK, quantity, sale_price,
                    disp_date, _synced_at, _source_id)
```

### 5.2 Schema: `tally_analytics`

```sql
-- Dimension Tables
dim_ledgers        (ledger_id PK, ledger_name, ledger_group, parent_group, ...)
dim_cost_centres   (cc_id PK, name, parent_cc, ...)
dim_inventory_items(item_id PK, name_en, name_ar, unit, category, ...)

-- Fact Tables
fact_vouchers      (voucher_id PK, voucher_type, ledger_id FK, amount, 
                    voucher_date, narration, cost_centre_id FK, _synced_at)
fact_stock_movements (movement_id PK, item_id FK, qty_in, qty_out, 
                      closing_stock, movement_date, _synced_at)
```

### 5.3 Schema: `app`

```sql
users              (user_id PK, email, keycloak_sub, role, department, 
                    locale VARCHAR(10) DEFAULT 'en', calendar_system DEFAULT 'gregorian')
user_preferences   (user_id PK FK, locale, number_format, calendar_system, 
                    report_language, ai_response_language)
audit_log          (log_id, user_id, action, resource, resource_id, 
                    changes_json, ip_address, locale, created_at)  -- append-only
etl_quarantine     (record_id, source, source_id, raw_json, error_reason, created_at)
notification_queue (notif_id, user_id, channel, subject, body_en, body_ar, 
                    status, scheduled_at)
pinned_charts      (pin_id, user_id, chart_config_json, refresh_interval, created_at)
scheduled_reports  (sched_id, user_id, report_type, params_json, 
                    cron_expr, locale, next_run_at)
```

### 5.4 Schema: `vectors`

```sql
schema_embeddings  (id, table_name, column_name, description, embedding vector(1536))
metric_embeddings  (id, metric_name, description, embedding vector(1536))
```

---

## 6. Infrastructure Architecture

```
On-Premises Network
├── App Server
│   ├── Go ETL Service       (port 8081)
│   ├── Go Transform Service (port 8082)  
│   ├── NATS Server          (port 4222)
│   ├── Redis                (port 6379, localhost only)
│   └── Keycloak             (port 8080)
│
├── DB Server
│   ├── PostgreSQL 18.2       (port 5432)
│   │   ├── Extensions: pgvector, pg_stat_statements, uuid-ossp
│   │   ├── Roles: medisync_readonly, medisync_app, medisync_etl
│   │   └── Schemas: hims_analytics, tally_analytics, semantic, app, vectors
│   └── (Redis also deployable here)
│
└── Source Systems
    ├── Tally ERP Server
    └── HIMS Server
```

---

## 7. ETL Pipeline Design

### Tally TDL Connector

**Protocol:** HTTP POST to Tally's built-in web server (port 9000 default)  
**Format:** XML request body with TDL report definitions  
**Extraction strategy:** Incremental using `LastAlterID` cursor stored in `app.etl_state`  
**Entities extracted per run:**
- `LedgerList` → `dim_ledgers`
- `VoucherList` (filtered by date range) → `fact_vouchers`
- `StockItemList` → `dim_inventory_items`
- `StockSummaryList` → `fact_stock_movements`

**Error handling:**
- HTTP timeout: 30 seconds per request; retry 3× with exponential backoff
- XML parse error: write raw XML to `etl_quarantine` with error reason
- Tally unavailable: emit `etl.sync.failed` → A-10 alert

### HIMS REST Connector

**Protocol:** REST/JSON (provider-specific; adapts to HIMS API contract)  
**Authentication:** API key in header (stored in HashiCorp Vault / OS keystore)  
**Sync frequencies:**
- Patient Demographics: daily (low-change data)
- Appointments: every 15 minutes
- Billing: every 15 minutes
- Pharmacy Dispensations: every 15 minutes

**Idempotency:** All inserts use `ON CONFLICT (source_id) DO UPDATE` on the `_source_id` column

---

## 8. Security Baseline

| Requirement | Implementation |
|---|---|
| DB access control | Three separate Postgres roles: `medisync_readonly` (SELECT only on analytics schemas), `medisync_app` (CRUD on app schema), `medisync_etl` (INSERT/UPDATE on analytics schemas) |
| Credential management | Tally and HIMS API credentials in environment secrets loaded from HashiCorp Vault; never in `.env` files or source code |
| Network isolation | Redis bound to `127.0.0.1` only; Postgres accepting connections from app server CIDR only |
| TLS | All HTTP services behind TLS 1.3 termination; Tally TDL over HTTPS |
| Keycloak | Realm with mandatory 2FA (TOTP) for `admin`, `finance_head`, `accountant_lead` roles |
| Audit log | `app.audit_log` table is append-only; enforced via PostgreSQL row-security policy rejecting UPDATE/DELETE |

---

## 9. Observability Setup

| Component | Configuration |
|---|---|
| Prometheus | Scrapes Go ETL service `/metrics` endpoint; `etl_sync_duration_seconds`, `etl_rows_synced_total`, `etl_quarantine_rows_total` metrics |
| Loki | Go services emit structured JSON logs; Loki ingests via Promtail agent |
| Grafana | "ETL Health" dashboard: sync lag, quarantine counts, quality scores per source |
| Alerts | PagerDuty/email alert when ETL sync fails 3× consecutive or quarantine rate > 1% |

---

## 10. Testing Requirements

| Test Type | Scope | Tool |
|---|---|---|
| Unit tests | Go connector functions, transform logic, C-06 quality checks | Go `testing` package |
| Integration tests | Full ETL run against test Tally + HIMS sandbox instances | Go integration test suite |
| Schema migration tests | All DDL applied to fresh Postgres instance without errors | `golang-migrate` |
| Data quality tests | C-06 agent assertions against injected dirty data | Go test fixtures |
| CI gate | `ar` translation keys present for all new `en` keys | E-05 script in GitHub Actions |

---

## 11. Dependencies & Pre-conditions

| Dependency | Owner | Notes |
|---|---|---|
| Tally ERP TDL access | Client IT | Must enable Tally web server (port 9000) and share TDL schema |
| HIMS API credentials | Client IT + HIMS vendor | REST API documentation and API key required |
| On-premises server provisioning | DevOps | Minimum spec: 8 vCPU, 32 GB RAM, 2 TB SSD |
| Network access between servers | Client IT | App server must reach Tally and HIMS servers over internal LAN |
| HashiCorp Vault (or OS keystore) | DevOps | For secrets management |

---

## 12. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Tally TDL API not exposing required fields | High | Early prototype with Tally admin; request custom TDL report definitions |
| HIMS API rate limiting blocking sync | Medium | Implement exponential backoff; request increased rate limits from vendor |
| On-premises server delayed provisioning | High | Begin schema design and ETL development in cloud dev environment in parallel |
| Data volume larger than estimated | Medium | pgvector and Postgres tuned from day 1; query explain plans reviewed |

---

## 13. Phase Exit Criteria

- [ ] Both ETL pipelines (Tally + HIMS) running reliably on schedule with < 1% quarantine rate
- [ ] All warehouse schemas created and populated with at least 90 days of historical data
- [ ] C-06 Data Quality Validation Agent running and sending alerts correctly
- [ ] NATS topics flowing; `etl.sync.completed` events confirmed in Grafana
- [ ] Keycloak realm live with demo users provisioned
- [ ] Docker Compose dev environment documented and tested by second engineer
- [ ] Translation file structure in place; E-05 CI check blocking PRs with missing `ar` keys
- [ ] Phase gate reviewed and signed off by Technical Lead and PM

---

*Phase 01 | Version 1.0 | February 19, 2026*
