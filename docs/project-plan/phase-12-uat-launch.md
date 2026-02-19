# Phase 12 — UAT & Production Launch v1

**Phase Duration:** Weeks 39–40 (2 weeks)  
**Module(s):** All  
**Status:** Planning  
**Milestone:** **M8 — Production v1 Public Launch** (all 3 core modules live)  
**Depends On:** Phase 11 complete (integration polish passed)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Deliver production-ready MediSync v1 to live healthcare facilities. Conduct structured User Acceptance Testing with real end-users covering all three modules, run security penetration testing, complete production deployment on customer infrastructure, activate monitoring, deliver user training, and publish API documentation.

---

## 2. Scope

### In Scope
- Structured UAT plan with real end-users across all 3 modules
- Security penetration testing (must resolve P0 and P1 findings before launch)
- Production Docker Compose / Kubernetes deployment on customer infrastructure
- SSL/TLS certificates, domain config, reverse proxy (Nginx)
- Monitoring stack fully operational (Prometheus + Grafana + Loki + alerting)
- Backup strategy activated (PostgreSQL automated backups)
- User training sessions (operations staff, accountants, finance managers)
- API documentation published (Swagger / OpenAPI spec)
- Runbook — incident response, rollback procedure
- Stakeholder launch demo

### Out of Scope
- Module D (Advanced Search Analytics) — planned Phase 13–18
- Any new features; this phase is hardening only

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | UAT Plan | PM + QA | Document with test scenarios, testers, success criteria; sign-off from stakeholder |
| D-02 | UAT Execution | All Engineers + QA | 80% of UAT scenarios pass first attempt; all P0/P1 bugs resolved |
| D-03 | Penetration Test | Security / External Vendor | All P0 findings: 0; all P1 findings: 0; P2 findings: remediated or waived with justification |
| D-04 | Production Deployment | DevOps | Live environment deployed; all services running; health checks green |
| D-05 | SSL/TLS + Domain | DevOps | HTTPS on all endpoints; Nginx reverse proxy; TLS 1.3 minimum |
| D-06 | Monitoring Stack | DevOps | Grafana dashboards: system health, AI agent latency, Tally sync status, error rate; PagerDuty/email alerts |
| D-07 | Backup Strategy | DevOps | PostgreSQL continuous WAL archiving; daily snapshot; restore tested |
| D-08 | User Training | PM + Eng | 3 training sessions delivered: (1) BI Dashboard (2) AI Accountant (3) Reports; attendance ≥ 80% |
| D-09 | API Documentation | Backend Eng | OpenAPI 3.1 spec published; Swagger UI accessible at `/api/docs`; all endpoints documented |
| D-10 | Runbook | DevOps + Backend Eng | Incident severity matrix, on-call escalation, rollback procedure, common issue FAQ |
| D-11 | Stakeholder Launch Demo | PM | Live production demo presented to stakeholders; approval to launch |

---

## 4. UAT Test Plan

### UAT Role Matrix

| Role | Module Coverage | #Testers |
|---|---|---|
| Clinic Administrator | BI Dashboard, KPI Alerts, Scheduled Reports | 2 |
| Accountant | AI Accountant (full flow), Easy Reports | 3 |
| Finance Manager | Approval Workflow, P&L, AIs reports, Forecasts | 2 |
| IT/Operations | System settings, User management, integrations | 1 |

### UAT Scenario Groups

**Group 1 — BI Dashboard (12 scenarios)**
- S1.1: Ask natural language question in English; verify correct chart
- S1.2: Ask natural language question in Arabic; verify correct chart + RTL layout
- S1.3: Drill down from chart; verify drill-down data
- S1.4: Pin chart to dashboard; verify persists after logout
- S1.5: View multi-period comparison; verify correct date ranges
- S1.6: KPI alert fires; verify notification (email + Slack)
- S1.7: Scheduled report arrives by email; verify PDF content
- S1.8: Anomaly detected; verify alert + explanation
- S1.9: Trend forecast chart rendered; verify ARIMA confidence intervals shown
- S1.10: Filter by cost-centre; verify data scoped correctly
- S1.11: Export to Excel; verify data accuracy
- S1.12: Mobile dashboard loads; KPI cards visible on phone

**Group 2 — AI Accountant (10 scenarios)**
- S2.1: Upload invoice PDF; verify classification + extraction
- S2.2: Upload handwritten receipt; verify OCR extracts key fields
- S2.3: Duplicate invoice rejected; verify duplicate detection fires
- S2.4: HITL review queue; verify reviewer can correct fields
- S2.5: Ledger mapping suggested; verify first suggestion is correct
- S2.6: 4-level approval chain completes; verify Tally creates voucher
- S2.7: Attempted self-approval blocked; verify OPA policy fires
- S2.8: Arabic vendor invoice processed; verify Arabic text extracted correctly
- S2.9: Bank reconciliation: matched + unmatched items shown correctly
- S2.10: Audit trail shows full history of approval chain

**Group 3 — Easy Reports (8 scenarios)**
- S3.1: Generate P&L report for last quarter; verify accuracy
- S3.2: Schedule monthly P&L; verify email delivery
- S3.3: Multi-company consolidation report; verify entities merged correctly
- S3.4: Budget vs Actual variance; verify variance calculations
- S3.5: Inventory aging report; verify reorder alerts
- S3.6: Custom report built with drag-and-drop builder
- S3.7: Row-level security: user A cannot see user B's data
- S3.8: Column masking: accountant cannot see "salary" column

---

## 5. Security Penetration Test Scope

**Vectors:**
- Authentication bypass attempts (Keycloak, JWT)
- SQL injection in chat NL query
- CSRF on approval workflow
- IDOR (insecure direct object reference) on document/report access
- API key bruteforce
- OPA policy bypass attempts
- SSRF via webhook endpoint
- File upload malicious payload (AI Accountant)
- Privilege escalation (Keycloak roles)

**Severity classification:**
| Grade | Definition | Release gate |
|---|---|---|
| P0 | Critical: data breach, auth bypass, RCE | **BLOCK release** |
| P1 | High: privilege escalation, mass IDOR | **BLOCK release** |
| P2 | Medium: limited exposure, difficult to exploit | Remediate within 7 days of launch |
| P3 | Low: informational | Backlog |

---

## 6. Production Infrastructure

```
Production Stack (Docker Compose → Kubernetes migration in Phase 18)
─────────────────────────────────────────────────────────────────
nginx             (reverse proxy, TLS termination)
medisync-web      (React app, port 3000)
medisync-api      (Go backend, port 8080)
medisync-worker   (NATS consumers, background jobs)
postgres-primary  (PostgreSQL 15, port 5432)
redis             (port 6379)
nats              (port 4222)
keycloak          (port 8443)
prometheus        (port 9090)
grafana           (port 3001)
loki              (port 3100)
```

**Backup strategy:**
- PostgreSQL: continuous WAL archiving to local/S3-compatible storage
- Point-in-time recovery (PITR) target: 5-minute RPO
- Daily pg_dump snapshot retained 30 days
- Redis: RDB snapshot every 900 seconds

---

## 7. Grafana Dashboards

**Dashboards to provision:**
| Dashboard | Key Metrics |
|---|---|
| System Overview | CPU, memory, disk I/O per service |
| AI Agent Performance | P50/P95/P99 query latency per agent; error rate |
| Tally Sync Status | Sync queue depth; success/failure rate; last sync time |
| AI Accountant Throughput | Documents/hr; OCR accuracy; approval queue depth |
| Error Rate | 5xx rate by endpoint; alert on > 1% |
| User Activity | Daily active users, queries/user, peak hours |

**Alert rules (PagerDuty / email):**
- Error rate > 1% for 5 minutes → P2 alert
- Tally sync failure 3× consecutive → P1 alert
- PostgreSQL disk > 80% → P2 alert
- AI agent P95 latency > 10s for 10 minutes → P2 alert

---

## 8. User Training Plan

| Session | Audience | Duration | Content |
|---|---|---|---|
| Training 1: BI Dashboard | Clinic administrators, managers | 2 hours | Chat queries, KPI alerts, drill-down, scheduled reports, mobile |
| Training 2: AI Accountant | Accountants, accounts payable | 3 hours | Document upload, OCR review, ledger mapping, approval workflow, Tally sync |
| Training 3: Easy Reports | Finance managers, department heads | 2 hours | Pre-built reports, scheduling, custom report builder, multi-company |
| Admin training | IT/operations | 1 hour | User management, security settings, backup/restore, monitoring |

**Training materials to produce:**
- User guide PDF per module (EN + AR)
- Quick-reference card (2-page) per module (EN + AR)
- Demo recording (screen recording) per module
- Admin operations guide

---

## 9. API Documentation

**Published at:** `/api/docs` (Swagger UI)  
**Spec format:** OpenAPI 3.1  
**Endpoints documented:**

| Category | Endpoints |
|---|---|
| Authentication | `/v1/auth/token`, `/v1/auth/refresh`, `/v1/auth/logout` |
| BI Chat | `POST /v1/chat` (SSE), `GET /v1/chat/history`, `DELETE /v1/chat/session/{id}` |
| Dashboards | `GET /v1/dashboards`, `POST /v1/dashboards`, `PATCH /v1/dashboards/{id}` |
| Documents | `POST /v1/documents/upload`, `GET /v1/documents`, `GET /v1/documents/{id}` |
| Approvals | `GET /v1/approvals/queue`, `POST /v1/approvals/{id}/approve`, `POST /v1/approvals/{id}/reject` |
| Reports | `GET /v1/reports`, `POST /v1/reports/generate`, `GET /v1/reports/{id}/download` |
| Schedules | `POST /v1/schedules`, `PATCH /v1/schedules/{id}/pause` |
| Webhooks | `POST /v1/webhooks`, `GET /v1/webhooks`, `DELETE /v1/webhooks/{id}` |
| External API | `POST /v1/api/query` |

---

## 10. Rollback Plan

**Trigger criteria:** > 5% error rate on critical path for 10 minutes OR data corruption detected

**Rollback steps:**
1. Revert Docker image tags to previous stable version
2. Run `docker compose up -d` with previous tag
3. If DB migration required reverting: restore from last PITR snapshot
4. Verify health checks pass
5. Notify stakeholders
6. Post-mortem within 24 hours

**Rollback time target:** < 30 minutes

---

## 11. Phase Exit Criteria

- [ ] UAT: ≥ 80% of scenarios pass first attempt
- [ ] UAT: all P0 and P1 bugs resolved and re-verified
- [ ] Penetration test: 0 P0, 0 P1 findings remaining
- [ ] Production deployment: all services healthy, health checks green
- [ ] SSL/TLS: HTTPS on all endpoints, TLS 1.3
- [ ] Monitoring: Grafana dashboards provisioned, alerts configured
- [ ] Backup: PITR tested — restore from snapshot successful
- [ ] Training: all 4 sessions delivered
- [ ] API documentation: OpenAPI spec published at `/api/docs`
- [ ] Runbook: incident response + rollback procedure documented
- [ ] Stakeholder launch demo: presented and approved
- [ ] **Milestone M8: Production v1 LAUNCHED**

---

*Phase 12 | Version 1.0 | February 19, 2026*
