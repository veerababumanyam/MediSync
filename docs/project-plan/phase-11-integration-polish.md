# Phase 11 — Integration & Polish (All Modules)

**Phase Duration:** Weeks 37–38 (2 weeks)  
**Module(s):** All (A, B, C, E)  
**Status:** Planning  
**Depends On:** Phase 10 complete (all three core modules live)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Unify Modules A, B, and C into a single cohesive product experience. This phase is about integration quality, not new features: cross-module data linking, performance hardening, mobile completeness, webhook/API surface, and fixing integration gaps discovered during multi-module testing sprints.

---

## 2. Scope

### In Scope
- Unified navigation shell (all 3 modules in one interface)
- Cross-module data linking (document → transaction → report → chat)
- Document management library (full implementation)
- Performance optimisation across all modules
- API surface for external integrations (webhooks, REST)
- Mobile app completeness (AI Accountant + Reports on Flutter)
- Slack/Teams webhook stubs → real integrations
- SMS gateway integration (real)
- End-to-end integration test suite
- Accessibility (WCAG 2.1 AA) review
- RTL final polish pass (all 3 modules)
- Cross-module OPA security review

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | Unified Navigation Shell | Frontend Engineer | Single app with sidebar nav covering BI Dashboard, AI Accountant, Reports; seamless routing |
| D-02 | Cross-Module Document Linking | Backend + Frontend | Document in accountant → linked to transaction → linked to Tally voucher → queryable in chat; deep link works |
| D-03 | Document Management Library | Frontend + Backend | Full CRUD document library; search by type/vendor/date/amount; linked to all associated transactions |
| D-04 | Performance Audit & Fix | All Engineers | P95 query latency < 5s; dashboard load < 3s; OOB (no regression on benchmarks from Phase 2/3) |
| D-05 | Redis Cache Optimisation | Backend + Data | Schema context cache tuned; query result cache hit rate ≥ 60% on repeated queries |
| D-06 | Webhook API | Backend Engineer | `POST /v1/webhooks` — configurable webhooks for: new report ready, sync completed, KPI alert, approval request |
| D-07 | External REST API | Backend Engineer | `/v1/api/query` — authenticated external clients can submit BI queries programmatically |
| D-08 | Mobile App — AI Accountant features | Frontend Engineer | Document upload (camera capture), review queue, transaction status on Flutter |
| D-09 | Mobile App — Reports features | Frontend Engineer | Report portal, scheduled report list, KPI dashboard on Flutter (read-only) |
| D-10 | Slack Integration | Backend Engineer | Webhook to post KPI alerts + scheduled report notifications to Slack channel |
| D-11 | SMS Gateway (real) | Backend Engineer | SMS notifications for critical KPI alerts (plugged into real SMS provider) |
| D-12 | RTL Final Polish | Frontend + QA | All 3 modules pass Playwright RTL regression suite; Arabic QA reviewer sign-off |
| D-13 | WCAG 2.1 AA Audit | QA | Keyboard navigation, screen reader, colour contrast passing for primary screens |
| D-14 | Integration Test Suite | QA | End-to-end test: document upload → OCR → ledger map → approve → Tally sync → report includes transaction |

---

## 4. Cross-Module Data Linking

**Deep link chain:**

```
Document (doc_id)
    │ B-02 OCR
    ▼
Extracted Transaction (extraction_id)
    │ B-05 Ledger Mapping
    ▼
Transaction Queue Entry (txn_id)
    │ B-08 Approval + B-09 Tally Sync
    ▼
Tally Voucher (tally_voucher_id)
    │ Indexed in tally_analytics.fact_vouchers
    ▼
BI Chat Query (A-01 can reference voucher by date/amount)
    │
    ▼
Report (C-01 P&L includes this voucher's ledger amount)
```

**UI manifestations:**
- In Document Library: "View Transaction" button → opens transaction detail
- In Transaction Detail: "Source Document" link → opens document preview
- In Tally Sync History: "View Document" and "View in Report" links
- In BI Chat: clicking a transaction in drill-down table → "View Source Document" tooltip
- In Report: footnote shows data quality score and last sync time

---

## 5. Performance Targets (Gate)

| Metric | Target | Measurement Method |
|---|---|---|
| Chat query P95 latency | < 5 seconds | k6 load test (50 concurrent users) |
| Dashboard load (4 dashboards) | < 3 seconds | Lighthouse CI |
| Report generation (P&L 12 months) | < 10 seconds | Integration test timer |
| Document upload (10 files) | Progress indicator; complete in reasonable time | Manual test |
| Redis cache hit rate | ≥ 60% | Prometheus `redis_cache_hits_total / total` |
| Postgres query P99 (A-01) | < 3 seconds | `pg_stat_statements` |

---

## 6. Webhook API

**Events supported:**
| Event | Payload |
|---|---|
| `report.ready` | `{ schedule_id, report_type, period, download_url }` |
| `sync.completed` | `{ entity_id, vouchers_created, status }` |
| `kpi.alert` | `{ metric_name, value, threshold, direction }` |
| `approval.requested` | `{ txn_id, approver_role, amount, vendor }` |
| `anomaly.detected` | `{ metric_name, z_score, observed, expected }` |

**Configuration:** Admin sets webhook URL + secret (HMAC-SHA256 signature on payload); retry 3× on failure.

---

## 7. External REST API

`GET/POST /v1/api/query` — Programmatic BI query access for embedded analytics or external integrations.

```json
// POST /v1/api/query
{
  "query": "Total pharmacy revenue last month",
  "locale": "en",
  "format": "json"
}

// Response
{
  "result": {...data...},
  "sql_used": "SELECT SUM(amount) FROM...",
  "chart_type": "kpiCard",
  "confidence": 92,
  "generated_at": "2026-02-19T10:00:00Z"
}
```

**Rate limiting:** 100 requests/minute per API key  
**Auth:** API key in header (`X-MediSync-API-Key`); scoped to requesting user's permissions

---

## 8. End-to-End Integration Test

**Full workflow test (automated):**
1. Upload test invoice PDF (known values)
2. Verify B-01 classifies as "invoice"
3. Verify B-02 extracts correct amount/vendor/date (within ±1%)
4. Verify B-04 matches to known vendor in ledger
5. Verify B-05 suggests correct GL ledger (first suggestion)
6. Accountant approves (simulated API call)
7. Manager approves (simulated)
8. Finance Head triggers Tally sync (simulated — Tally sandbox)
9. Verify Tally voucher created with correct values
10. Verify B-14 audit log entry created
11. Run P&L report for the period; verify transaction appears in correct ledger row
12. Chat query: "What did we pay [vendor] this month?"; verify correct answer referencing the transaction

---

## 9. Testing Requirements

| Test | Scope | Target |
|---|---|---|
| End-to-end integration | Full workflow test above | All 12 steps pass |
| RTL regression | All modules, all RTL screens | 0 new RTL regressions vs Phase 10 baseline |
| Performance | k6 load test (50 concurrent) | All SLOs met |
| Webhook delivery | All 5 event types | Delivered to mock endpoint with correct payload + HMAC |
| External API | 10 query types via REST | Results match chat interface |
| Mobile completeness | AI Accountant + Reports on Flutter | Core flows usable on iOS + Android |
| WCAG | Primary screens (10 screens) | WCAG 2.1 AA: 0 critical, < 3 minor |

---

## 10. Phase Exit Criteria

- [ ] Unified navigation shell with all 3 modules accessible
- [ ] Cross-module deep linking working: document → transaction → voucher → report
- [ ] Performance SLOs met (P95 chat < 5s, dashboard < 3s)
- [ ] Webhook API and External REST API live with documentation
- [ ] Mobile app covers AI Accountant + Reports features
- [ ] Slack and SMS notifications working with real integrations
- [ ] RTL final polish pass signed off by Arabic QA reviewer
- [ ] Full end-to-end integration test suite passing
- [ ] Phase gate reviewed and signed off

---

*Phase 11 | Version 1.0 | February 19, 2026*
