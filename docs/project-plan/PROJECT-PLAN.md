# MediSync — Master Project Plan

**Version:** 1.0  
**Status:** Active  
**Created:** February 19, 2026  
**Total Duration:** 54 Weeks (≈ 13.5 months)  
**Cross-ref:** [PRD.md](../PRD.md) | [ARCHITECTURE.md](../ARCHITECTURE.md) | [agents/BLUEPRINTS.md](../agents/BLUEPRINTS.md)

---

## Executive Summary

MediSync is an AI-powered, chat-based Business Intelligence platform that unifies operational data from a **Healthcare Information Management System (HIMS)** and financial data from **Tally ERP**. The platform delivers four core modules powered by **58 autonomous AI agents**, orchestrated via **Google Genkit** and the **A2A Protocol**, with first-class English and Arabic (RTL) support.

This master project plan organises all 18 delivery phases into cohesive workstreams with clearly-defined scope, deliverables, team allocations, risks, and acceptance criteria for each phase. Every phase has a dedicated document linked below.

---

## Product Modules Summary

| Module | Description | AI Agents | Phases |
|---|---|---|---|
| **Module A** — Conversational BI Dashboard | Natural-language query → SQL → chart/table inline in chat | 13 agents (A-01 to A-13) | 2, 3, 7 |
| **Module B** — AI Accountant | Document OCR → ledger mapping → approval → Tally sync | 16 agents (B-01 to B-16) | 4, 5, 6, 7 |
| **Module C** — Easy Reports | Pre-built MIS reports, zero-code builder, scheduled delivery | 8 agents (C-01 to C-08) | 8, 9, 10 |
| **Module D** — Advanced Search Analytics | Autonomous agents, semantic layer, prescriptive AI | 14 agents (D-01 to D-14) | 13–18 |
| **Module E** — Language & Localisation | Arabic/English detection, translation, RTL formatting | 7 agents (E-01 to E-07) | 1–6 |

---

## Phase Overview

| Phase | Name | Weeks | Status | Document |
|---|---|---|---|---|
| **1** | ETL & Infrastructure | 1–3 | Planning | [phase-01-etl-infrastructure.md](./phase-01-etl-infrastructure.md) |
| **2** | AI Agent & Core Analytics | 4–7 | Planning | [phase-02-ai-agent-core-analytics.md](./phase-02-ai-agent-core-analytics.md) |
| **3** | Dashboard & Advanced Features | 8–11 | Planning | [phase-03-dashboard-advanced-features.md](./phase-03-dashboard-advanced-features.md) |
| **4** | Document Processing (AI Accountant) | 12–15 | Planning | [phase-04-document-processing.md](./phase-04-document-processing.md) |
| **5** | Transaction Intelligence | 16–19 | Planning | [phase-05-transaction-intelligence.md](./phase-05-transaction-intelligence.md) |
| **6** | Tally Real-Time Integration | 20–22 | Planning | [phase-06-tally-realtime-integration.md](./phase-06-tally-realtime-integration.md) |
| **7** | Reconciliation & Financial Analytics | 23–26 | Planning | [phase-07-reconciliation-analytics.md](./phase-07-reconciliation-analytics.md) |
| **8** | Pre-Built Reports & Dashboards | 27–30 | Planning | [phase-08-prebuilt-reports.md](./phase-08-prebuilt-reports.md) |
| **9** | Report Automation & Distribution | 31–33 | Planning | [phase-09-report-automation.md](./phase-09-report-automation.md) |
| **10** | Customization & RBAC Security | 34–36 | Planning | [phase-10-customization-rbac.md](./phase-10-customization-rbac.md) |
| **11** | Integration & Polish | 37–38 | Planning | [phase-11-integration-polish.md](./phase-11-integration-polish.md) |
| **12** | UAT & Production Launch v1 | 39–40 | Planning | [phase-12-uat-launch.md](./phase-12-uat-launch.md) |
| **13** | Semantic Layer & NL Search | 41–43 | Planning | [phase-13-semantic-nlsearch.md](./phase-13-semantic-nlsearch.md) |
| **14** | Autonomous AI Agents & Deep Research | 44–46 | Planning | [phase-14-autonomous-agents.md](./phase-14-autonomous-agents.md) |
| **15** | Auto-Dashboarding & Prescriptive AI | 47–48 | Planning | [phase-15-auto-dashboarding-insights.md](./phase-15-auto-dashboarding-insights.md) |
| **16** | Developer Tools & Embedded Analytics | 49–50 | Planning | [phase-16-developer-tools-embedding.md](./phase-16-developer-tools-embedding.md) |
| **17** | Data Governance & Compliance | 51–52 | Planning | [phase-17-governance-compliance.md](./phase-17-governance-compliance.md) |
| **18** | Final Integration UAT & Launch v2 | 53–54 | Planning | [phase-18-final-integration-launch.md](./phase-18-final-integration-launch.md) |

---

## Agent Deployment Schedule

```
Phase  1  ──  C-06 (Data Quality)
Phase  2  ──  A-01, A-02, A-03, A-04, A-05, A-06, E-01, E-02, E-03
Phase  3  ──  A-07, A-08, A-09, A-10, A-11, E-07
Phase  4  ──  B-01, B-02, B-03
Phase  5  ──  B-04, B-05, B-06, B-07, B-08, E-04
Phase  6  ──  B-09, B-14, B-16, E-06
Phase  7  ──  B-10, B-11, B-12, B-13, B-15, A-12, A-13
Phase  8  ──  C-01, C-02, C-07, C-08
Phase  9  ──  C-03
Phase 10  ──  C-04, C-05
Phase 13  ──  D-01, D-02, D-09
Phase 14  ──  D-03, D-04, D-05, D-10, D-13
Phase 15  ──  D-06, D-07, D-08, D-14
Phase 16  ──  D-11, D-12
Phase 17  ──  C-05 (extended), D-09 (governance tier)
Phase 18  ──  All 58 agents — integration hardening
```

---

## Technology Stack Summary

### Backend
| Layer | Technology |
|---|---|
| Core API | Go 1.26, go-chi |
| ETL | Meltano + Go ETL Services |
| AI Orchestration | Google Genkit (Apache-2.0) |
| Multi-Agent | Agent ADK + A2A Protocol |
| Messaging | NATS / JetStream |
| Identity | Keycloak (OIDC/SAML/2FA) |
| Authorization | Open Policy Agent (OPA) |

### Data
| Layer | Technology |
|---|---|
| Data Warehouse | PostgreSQL 18.2 (on-premises) |
| Vector Search | pgvector → Milvus (upgrade path) |
| Cache | Redis |
| Offline Sync | PowerSync |

### AI / ML
| Component | Technology |
|---|---|
| LLM (Cloud) | GPT-4o / Claude 3.5 Sonnet / Gemini Pro |
| LLM (Local) | Ollama (Llama 3), vLLM |
| OCR | PaddleOCR |
| Semantic Layer | MetricFlow |

### Frontend
| Layer | Technology |
|---|---|
| Web | React 19.2.4 + CopilotKit + Vite 7.3 |
| Charts | Apache ECharts |
| Mobile | Flutter (iOS + Android) |
| i18n Web | i18next (EN/AR) |
| i18n Mobile | flutter_localizations + ARB |

### Observability
| Component | Technology |
|---|---|
| Metrics | Prometheus |
| Dashboards | Grafana |
| Logs | Loki |
| Tracing | Genkit built-in spans |

---

## Team Structure

| Role | Responsibilities | Allocation |
|---|---|---|
| **Technical Lead / Architect** | Architecture decisions, ADRs, cross-phase consistency | Full-time |
| **Backend Engineers (Go)** | API gateway, ETL services, AI orchestration glue | 2–3 FTEs |
| **AI/ML Engineers** | Genkit flows, agent development, prompt engineering, OCR integration | 2 FTEs |
| **Frontend Engineers** | React web app, Flutter mobile, CopilotKit integration | 2 FTEs |
| **Data Engineer** | PostgreSQL schema, Meltano pipelines, pgvector tuning | 1 FTE |
| **DevOps / Platform** | Docker/K8s, CI/CD, Keycloak, OPA, Grafana/Loki | 1 FTE |
| **QA Engineer** | Playwright tests, RTL regression, API test suites | 1 FTE |
| **Arabic-Speaking QA Reviewer** | RTL layout validation, Arabic translation QA | Part-time |
| **Product Manager** | Backlog, stakeholder comms, UAT coordination | Full-time |
| **UX Designer** | Wireframes, RTL UX patterns, design system | Part-time → Full-time Phase 3+ |

---

## Key Milestones & Gate Reviews

| Milestone | Target Week | Description |
|---|---|---|
| **M1 — Data Foundation** | End Week 3 | ETL pipelines running; warehouse populated with sample Tally + HIMS data |
| **M2 — First AI Query** | End Week 7 | Core Text-to-SQL agent answering financial queries end-to-end |
| **M3 — Chat Dashboard MVP** | End Week 11 | Full chat UI with charts, pinnable dashboards, scheduled reports |
| **M4 — OCR Pipeline Live** | End Week 15 | Bulk document upload, OCR extraction, confidence scoring live |
| **M5 — AI Bookkeeping** | End Week 22 | Full document→ledger→approval→Tally sync pipeline operational |
| **M6 — Reconciliation Suite** | End Week 26 | Bank reconciliation, cash flow forecasting, tax compliance live |
| **M7 — Reports Module** | End Week 36 | All 8 Easy Reports agents live; zero-code builder deployed |
| **M8 — Production v1** | End Week 40 | Hardened, UAT-cleared, security-audited production launch |
| **M9 — Search Analytics** | End Week 43 | Semantic layer + NL search infrastructure operational |
| **M10 — Autonomous Agents** | End Week 46 | Spotter agent, deep research, autonomous monitoring live |
| **M11 — Prescriptive AI** | End Week 48 | Insight engine, auto-dashboards, prescriptive recommendations |
| **M12 — Production v2** | End Week 54 | Full 58-agent platform — final security audit and relaunch |

---

## Risk Register (Cross-Phase)

| Risk | Severity | Probability | Mitigation |
|---|---|---|---|
| ETL sync failures disrupting analytics | High | Medium | Idempotent upserts, NATS alerts, quarantine table, retry logic |
| AI hallucinations in financial context | High | Medium | A-05 Hallucination Guard, A-06 Confidence Scoring, HITL queue |
| OCR accuracy < 95% on handwritten docs | High | Medium | HITL queue for low-confidence extractions; human correction feedback loop |
| Tally write-back causing ledger corruption | Critical | Low | Intelligence/Action plane separation; OPA gate; B-08 approval always required |
| Arabic RTL layout regressions | Medium | Medium | E-05 CI translation guard; Playwright RTL visual regression suite |
| Scope creep across 4 modules | High | High | Strict phase gate reviews; dedicated PM; backlog grooming every 2 weeks |
| Performance degradation at scale | Medium | Low | Redis caching strategy; read-only Postgres role; pgvector index tuning |
| Healthcare data breach (HIPAA/GDPR) | Critical | Low | AES-256 at rest, TLS 1.3, OPA row/column masking, Keycloak 2FA |
| LLM provider API instability | Medium | Medium | Genkit plugin abstraction; fallback to local Ollama model |
| Key staff attrition during 54-week timeline | High | Medium | Documentation-first culture; knowledge sharing sessions; onboarding guides |

---

## Governance & Decision Process

- **Weekly Standup** — All engineers; blockers, progress, PR reviews
- **Bi-Weekly Phase Review** — PM + Tech Lead; backlog grooming, phase gate assessment
- **Phase Gate Sign-Off** — PM + Stakeholder approval required to advance to next phase
- **Architecture Decision Records (ADRs)** — All significant tech decisions logged in `ARCHITECTURE.md §16`
- **Security Review** — DevOps + Tech Lead; before every phase gate involving new data flows
- **RTL QA Sign-Off** — Arabic-speaking QA reviewer must sign off before any UI-heavy phase completes
- **Stakeholder Demo** — End of Phases 3, 7, 12, 18 — live demos with clinic manager, pharmacy lead, accountant, and owner

---

## Definition of Done (DoD)

A phase is considered **Done** when all of the following are true:

1. All deliverables listed in the phase document are implemented and merged to `main`
2. All AI agents for the phase have passing unit tests (≥ 80% coverage) and integration tests
3. All new API endpoints have OpenAPI spec entries
4. Translation keys for all new UI strings are present in both `en` and `ar` namespaces (E-05 CI gate passes)
5. Playwright RTL visual regression tests pass for all new screens
6. OPA policies for new data flows are written and reviewed
7. Grafana dashboards and Prometheus alerts are configured for new services
8. Phase gate review passed with PM and Technical Lead sign-off
9. Stakeholder demo completed (where applicable) with no blocking feedback

---

*Document Version: 1.0 | Last Updated: February 19, 2026 | Owner: MediSync Project Management*
