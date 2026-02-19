# MediSync — Comprehensive System Architecture

**Version:** 1.0  
**Status:** Approved — Engineering Baseline  
**Last Updated:** February 19, 2026  
**Cross-ref:** [PRD.md](./PRD.md) | [DESIGN.md](./DESIGN.md) | [agents/BLUEPRINTS.md](./agents/BLUEPRINTS.md) | [i18n-architecture.md](./i18n-architecture.md) | [agents/03-governance-security.md](./agents/03-governance-security.md)

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Architecture Principles](#2-architecture-principles)
3. [High-Level System Topology](#3-high-level-system-topology)
4. [Layer-by-Layer Architecture](#4-layer-by-layer-architecture)
   - 4.1 [Data Source Layer](#41-data-source-layer)
   - 4.2 [ETL / Ingestion Layer](#42-etl--ingestion-layer)
   - 4.3 [Data Warehouse Layer](#43-data-warehouse-layer)
   - 4.4 [AI Orchestration Layer](#44-ai-orchestration-layer)
   - 4.5 [API Gateway Layer](#45-api-gateway-layer)
   - 4.6 [Frontend Layer](#46-frontend-layer)
5. [Module Architectures](#5-module-architectures)
   - 5.1 [Module A — Conversational BI Dashboard](#51-module-a--conversational-bi-dashboard)
   - 5.2 [Module B — AI Accountant](#52-module-b--ai-accountant)
   - 5.3 [Module C — Easy Reports](#53-module-c--easy-reports)
   - 5.4 [Module D — Advanced Search-Driven Analytics](#54-module-d--advanced-search-driven-analytics)
6. [AI Agent Ecosystem](#6-ai-agent-ecosystem)
   - 6.1 [Agent Inventory](#61-agent-inventory)
   - 6.2 [Inter-Agent Communication (A2A Protocol)](#62-inter-agent-communication-a2a-protocol)
   - 6.3 [Cross-Cutting Agent Services](#63-cross-cutting-agent-services)
   - 6.4 [Agent Complexity Tiers](#64-agent-complexity-tiers)
7. [Data Architecture](#7-data-architecture)
   - 7.1 [Warehouse Schema Overview](#71-warehouse-schema-overview)
   - 7.2 [Vector Storage](#72-vector-storage)
   - 7.3 [Semantic Layer](#73-semantic-layer)
   - 7.4 [Caching Strategy](#74-caching-strategy)
8. [Security & Governance Architecture](#8-security--governance-architecture)
   - 8.1 [Identity & Authentication (Keycloak)](#81-identity--authentication-keycloak)
   - 8.2 [Authorization — Policy as Code (OPA)](#82-authorization--policy-as-code-opa)
   - 8.3 [Intelligence Plane vs. Action Plane](#83-intelligence-plane-vs-action-plane)
   - 8.4 [Role Definitions](#84-role-definitions)
   - 8.5 [Data Encryption](#85-data-encryption)
   - 8.6 [HITL (Human-in-the-Loop) Gates](#86-hitl-human-in-the-loop-gates)
9. [Internationalisation Architecture (i18n)](#9-internationalisation-architecture-i18n)
10. [Messaging & Event Bus](#10-messaging--event-bus)
11. [Observability Stack](#11-observability-stack)
12. [Offline & Mobile Sync](#12-offline--mobile-sync)
13. [Technology Stack Reference](#13-technology-stack-reference)
14. [Deployment Topology](#14-deployment-topology)
15. [Phased Delivery Roadmap](#15-phased-delivery-roadmap)
16. [Architecture Decision Records (ADRs)](#16-architecture-decision-records-adrs)

---

## 1. System Overview

**MediSync** is an AI-powered, chat-based Business Intelligence platform that unifies operational data from a **Healthcare Information Management System (HIMS)** and financial data from **Tally ERP**. It surfaces cross-domain insights through a conversational interface, automates bookkeeping via an AI Accountant module, delivers enterprise-grade reporting through an Easy Reports module, and offers agentic search-driven analytics powered by 58 autonomous AI agents.

### Core Problem

| Pain Point | Current State | MediSync Solution |
|---|---|---|
| Siloed data | HIMS (operations) and Tally (finance) never speak to each other | Unified data warehouse with incremental ETL every 15–60 min |
| Manual reporting | Hours of spreadsheet manipulation per report | Natural-language query → instant chart/table in < 5 seconds |
| Manual bookkeeping | Human entry of invoices, ledger mapping, Tally posting | OCR extraction → AI ledger mapping → one-click Tally sync |
| Static dashboards | Reports are stale snapshots | Auto-refreshing, pinnable, AI-generated dashboards |
| English-only tools | Arabic-speaking staff forced to context-switch | First-class Arabic (RTL) support across all surfaces |

---

## 2. Architecture Principles

| # | Principle | Rationale |
|---|---|---|
| P-01 | **Decoupled data plane** — ETL to separate warehouse | Protects HIMS/Tally from heavy analytical queries |
| P-02 | **Read-only intelligence plane** — AI agents use a `SELECT`-only DB role | Prevents AI-caused ledger corruption |
| P-03 | **Human-in-the-loop for all write-backs** — Action plane gated by approval workflow | Regulatory compliance and auditability |
| P-04 | **Policy as Code** — All authorization via OPA Rego policies | Auditable, version-controlled, hot-reloadable rules |
| P-05 | **Open-source only** — OSI-approved licenses throughout | Vendor independence, auditability, no licence risk |
| P-06 | **Go-first backend** — Go for API, ETL, and AI orchestration glue | Performance, low memory footprint, native concurrency |
| P-07 | **Genkit for AI flows** — Type-safe, observable AI pipelines | Unified tracing, retry logic, and prompt management |
| P-08 | **A2A Protocol for inter-agent communication** — Standardized discovery and collaboration | Decoupled agent ecosystem, 58+ agents coordinating safely |
| P-09 | **i18n by default** — Locale injected at every layer | Arabic (RTL) support without architectural retrofitting |
| P-10 | **Incremental delivery** — 18 phases over 54 weeks | Risk reduction through continuous value delivery |

---

## 3. High-Level System Topology

```
╔══════════════════════════════════════════════════════════════════════════════════╗
║                           EXTERNAL DATA SOURCES                                  ║
║  ┌─────────────────────┐          ┌──────────────────────────────────────────┐  ║
║  │  Tally ERP (TDL XML)│          │  HIMS (REST API)                         │  ║
║  │  • Ledgers           │          │  • Patients / Demographics               │  ║
║  │  • Vouchers          │          │  • Appointments                          │  ║
║  │  • Inventory         │          │  • Pharmacy Dispensations                │  ║
║  │  • Sales / Receipts  │          │  • Billing                               │  ║
║  └──────────┬──────────┘          └────────────────────┬─────────────────────┘  ║
╚═════════════╪═══════════════════════════════════════════╪════════════════════════╝
              │ XML/HTTP (TDL)                            │ REST/JSON
              ▼                                           ▼
╔══════════════════════════════════════════════════════════════════════════════════╗
║                         ETL / INGESTION LAYER (Go + Meltano)                     ║
║  ┌─────────────────────────────────────────────────────────────────────────────┐ ║
║  │  Meltano ELT  →  Go ETL Service  →  Data Validation (C-06 Agent)           │ ║
║  │  Incremental sync (15–60 min)  |  Schema normalisation  |  Dedup checks    │ ║
║  └────────────────────────────────────────────┬────────────────────────────────┘ ║
╚═══════════════════════════════════════════════╪════════════════════════════════════╝
                                                │
                                                ▼
╔══════════════════════════════════════════════════════════════════════════════════╗
║                      DATA WAREHOUSE LAYER (PostgreSQL + pgvector)                ║
║   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  ║
║   │  HIMS Schema │  │ Tally Schema │  │ Semantic Layer│  │  Vector Store    │  ║
║   │  (analytics) │  │  (analytics) │  │ (MetricFlow)  │  │  (pgvector /     │  ║
║   └──────────────┘  └──────────────┘  └──────────────┘  │   Milvus)        │  ║
║                                                           └──────────────────┘  ║
║   Cache: Redis    |   Read-only role: medisync_readonly                         ║
╚══════════════════════════════════════════════════════════════════════════════════╝
                                                │ READ-ONLY
                                                ▼
╔══════════════════════════════════════════════════════════════════════════════════╗
║                    AI ORCHESTRATION LAYER (Genkit + Agent ADK)                   ║
║                                                                                  ║
║   ┌─────────────────────────────────────────────────────────────────────────┐   ║
║   │                    Agent Supervisor (A2A Protocol)                       │   ║
║   │  Module A: 13 agents  │  Module B: 16 agents  │  Module C: 8 agents     │   ║
║   │  Module D: 14 agents  │  Module E: 7 agents   │  Total: 58 agents       │   ║
║   └───────────────────────────────────────────┬─────────────────────────────┘   ║
║                                               │                                  ║
║   Cross-cutting services:                     │                                  ║
║   OPA Policy Engine | Confidence Scoring | Audit Log Writer | Notification Svc  ║
╚═══════════════════════════════════════════════╪════════════════════════════════════╝
                                                │
                                                ▼
╔══════════════════════════════════════════════════════════════════════════════════╗
║                        API GATEWAY (Go / go-chi)                                 ║
║   Keycloak JWT validation  |  Accept-Language negotiation  |  Rate limiting     ║
║   REST + GraphQL  |  WebSocket (live dashboard push)  |  OPA authz sidecar     ║
╚══════════════╤═══════════════════════════════════════════════════════╤═══════════╝
               │                                                       │
               ▼                                                       ▼
╔═════════════════════════════╗                  ╔═════════════════════════════════╗
║  WEB FRONTEND               ║                  ║  MOBILE APP (Flutter)           ║
║  React + CopilotKit         ║                  ║  iOS / Android                  ║
║  Apache ECharts             ║                  ║  PowerSync (offline-first)      ║
║  i18next (EN/AR)            ║                  ║  flutter_localizations (ARB)    ║
║  Tailwind RTL               ║                  ║  ECharts Flutter                ║
╚═════════════════════════════╝                  ╚═════════════════════════════════╝


                    ┌──────────────────────────────────┐
                    │  ACTION PLANE (Tally Write-Back)  │
                    │  B-08 Approval → B-09 Tally Sync  │
                    │  OPA gate: finance_head only       │
                    │  TDL XML HTTP POST                 │
                    └──────────────────────────────────┘
```

---

## 4. Layer-by-Layer Architecture

### 4.1 Data Source Layer

MediSync is a **read-side platform** — it never writes to source systems except through the explicitly human-approved Tally sync pipeline.

#### Tally ERP Integration
- **Protocol:** TDL (Tally Definition Language) XML over HTTP POST/GET
- **Extracted entities:** Ledgers, Vouchers (Sales, Purchase, Receipt, Payment, Journal), Inventory Masters, Sales Orders, Stock Items, Cost Centres
- **Extraction mode:** Incremental — uses `LastAlterID` or date-range filters to fetch only changed records
- **Sync interval:** Configurable; default 30 minutes; real-time mode available for Phase 6+

#### HIMS Integration
- **Protocol:** Native REST API (provider-specific; typically JSON)
- **Extracted entities:** Patient Demographics, Appointments, Doctor Schedules, Billing Records, Pharmacy Dispensations, Prescription Data
- **Sync interval:** 15 minutes for operational data; daily for demographic snapshots

### 4.2 ETL / Ingestion Layer

```
Tally TDL → [Go Tally Connector]  ─┐
                                    ├─► Meltano ELT Pipeline ─► Go Transform Service ─► PostgreSQL DW
HIMS REST  → [Go HIMS Connector]  ─┘
                                        │
                                        ├─► C-06 Data Quality Validation Agent
                                        │     • Missing value checks
                                        │     • Referential integrity
                                        │     • Duplicate detection
                                        │     • Anomaly alerts
                                        │
                                        └─► NATS (publish: etl.sync.completed)
```

**Key ETL characteristics:**
- **Idempotent inserts** using `ON CONFLICT DO UPDATE` (upsert) on business keys
- **Change data capture (CDC)** timestamps on every warehouse row (`_synced_at`, `_source`, `_source_id`)
- **Backfill mode** for historical data on first run; incremental thereafter
- **Error quarantine:** Records failing validation are written to `etl_quarantine` table with error reason; never silently dropped
- **Alerting:** NATS events consumed by A-10 KPI Alert Agent when sync fails or data quality thresholds are breached

### 4.3 Data Warehouse Layer

**Primary Store:** PostgreSQL 18.2 (self-hosted, on-premises)

```
PostgreSQL Instance
├── Schema: hims_analytics
│   ├── dim_patients
│   ├── dim_doctors
│   ├── fact_appointments
│   ├── fact_billing
│   └── fact_pharmacy_dispensations
│
├── Schema: tally_analytics
│   ├── dim_ledgers
│   ├── dim_cost_centres
│   ├── dim_inventory_items
│   ├── fact_vouchers
│   └── fact_stock_movements
│
├── Schema: semantic
│   ├── metric_definitions       ← MetricFlow governed metrics
│   ├── dimension_definitions
│   └── metric_lineage
│
├── Schema: app
│   ├── users
│   ├── user_preferences         ← locale column
│   ├── pinned_charts
│   ├── scheduled_reports
│   ├── approval_workflows
│   ├── audit_log                ← append-only, immutable
│   ├── etl_quarantine
│   └── notification_queue
│
└── Extension: pgvector
    └── Schema: vectors
        ├── schema_embeddings    ← serialised schema context for LLM
        └── metric_embeddings    ← semantic metric descriptions
```

**Access control:**
- `medisync_readonly` — `GRANT SELECT` on `hims_analytics`, `tally_analytics`, `semantic` schemas only. Used by all AI agents.
- `medisync_app` — Full CRUD on `app` schema. Used by API service accounts.
- `medisync_etl` — `INSERT/UPDATE` on `hims_analytics` and `tally_analytics`. Used by ETL service only.
- No user or agent ever connects with superuser privileges.

**Extensions:**
- `pgvector` — Vector similarity search for schema indexing and NL search
- `pg_stat_statements` — Query performance monitoring
- `uuid-ossp` — UUID primary keys

### 4.4 AI Orchestration Layer

The AI orchestration layer is built on **Google Genkit** (Apache-2.0) as the primary AI flow framework, with **Agent ADK** for multi-agent coordination.

```
User Request
     │
     ▼
┌─────────────────────────────────────────────────┐
│              Genkit Flow Router                  │
│  (Detects intent: BI query / document upload /  │
│   report request / search query)                │
└──────┬──────────────────────────────────────────┘
       │
       ├──► Module A Flow (Conversational BI)
       │       └─► A-01 Text-to-SQL Agent
       │            → A-04 Domain Normaliser
       │            → A-06 Confidence Scorer
       │            → A-03 Visualisation Router
       │
       ├──► Module B Flow (AI Accountant)
       │       └─► B-01 Document Classifier
       │            → B-02 OCR Extraction
       │            → B-05 Ledger Mapping
       │            → B-08 Approval Workflow
       │            → B-09 Tally Sync
       │
       ├──► Module C Flow (Easy Reports)
       │       └─► C-01 Report Generator
       │            → C-02 Multi-Company Consolidation
       │            → C-03 Scheduling & Distribution
       │
       └──► Module D Flow (Search Analytics)
               └─► D-01 NL Search
                    → D-04 Autonomous Analyst
                    → D-08 Prescriptive Recommendations
```

**Genkit features used:**
- **Flows** — Type-safe, traceable AI pipelines with explicit input/output schemas (Pydantic-style Go structs)
- **Plugins** — Interchangeable LLM backends (Gemini 3, GPT-5.2, Ollama/Llama 4)
- **Observability** — Built-in tracing, span logging per flow step; exported to Grafana/Loki
- **Retry logic** — Configurable retry with exponential backoff on transient LLM errors
- **Streaming** — Server-Sent Events (SSE) for streaming chart/text responses to CopilotKit

### 4.5 API Gateway Layer

Built in **Go** using `go-chi` router:

```
Client Request
     │
     ▼
┌──────────────────────────────────────────────────┐
│                 API Gateway (Go / go-chi)         │
│                                                  │
│  1. TLS termination (mTLS for service accounts)  │
│  2. JWT validation (Keycloak JWKS endpoint)      │
│  3. Accept-Language negotiation → locale ctx     │
│  4. Rate limiting (per-user token bucket)        │
│  5. OPA authz sidecar call                       │
│  6. Request routing to handler                   │
│  7. Response locale formatting (E-03 Agent)      │
└────────────────────────┬─────────────────────────┘
                         │
              ┌──────────┼──────────┐
              ▼          ▼          ▼
         REST API   GraphQL    WebSocket
         /v1/...    /graphql   /ws/chat
```

**API surface:**
| Namespace | Description |
|---|---|
| `POST /v1/chat` | Submit conversational BI query; streams response via SSE |
| `GET  /v1/dashboard/{id}` | Fetch pinned dashboard configuration |
| `POST /v1/dashboard/pin` | Pin a generated chart to dashboard |
| `POST /v1/documents/upload` | Upload documents for OCR processing |
| `GET  /v1/reports` | List available reports |
| `POST /v1/reports/generate` | Generate a report on demand |
| `POST /v1/reports/schedule` | Create a scheduled report |
| `GET  /v1/sync/status` | Check Tally sync connection status |
| `POST /v1/sync/now` | Trigger manual Tally sync (finance_head only) |
| `GET  /v1/approvals` | Fetch pending approval queue |
| `POST /v1/approvals/{id}/approve` | Approve a transaction for Tally sync |
| `WS   /ws/chat` | WebSocket for real-time streaming chat |
| `GET  /graphql` | GraphQL endpoint for dashboard and report data |

### 4.6 Frontend Layer

#### React Web Application (CopilotKit + Generative UI)

```
React App (Vite build)
├── CopilotKit Provider (Generative UI orchestration)
│   ├── CopilotChat — chat window with streaming responses
│   ├── CopilotTask — background agent task runner
│   └── useCopilotAction — dynamic widget rendering hooks
│
├── Layouts
│   ├── DashboardLayout — pinned charts grid (responsive)
│   ├── ChatLayout — conversational interface
│   ├── ReportsLayout — report builder and list
│   └── AccountantLayout — document upload, reconciliation views
│
├── Visualisation (Apache ECharts via echarts-for-react)
│   ├── DynamicChart — renders bar/line/pie/scatter/table from agent JSON
│   ├── KPICard — sparkline + value + trend indicator
│   └── DrillDownTable — expandable, paginated data table
│
├── i18n
│   ├── i18next provider (en/ar namespaces, lazy-loaded)
│   ├── RTL dir switching on <html> element
│   └── Tailwind logical properties (inline-start/end)
│
└── State Management
    ├── Zustand — global app state (user, locale, active dashboard)
    └── React Query — server state + cache (reports, chart configs)
```

#### Flutter Mobile Application

```
Flutter App
├── Routing: go_router (declarative, deep-linkable)
├── State: Riverpod (reactive, testable)
├── Offline-first sync: PowerSync (Apache-2.0)
│   └── Syncs pinned dashboard configs and last-N reports offline
├── Visualisation: SyncFusion Charts / fl_chart
├── i18n: flutter_localizations + intl ARB files
│   ├── app_en.arb
│   └── app_ar.arb
├── RTL: Directionality widget + EdgeInsetsDirectional
└── Push notifications: FCM (optional) / local_notifications
```

---

## 5. Module Architectures

### 5.1 Module A — Conversational BI Dashboard

The core module. Converts natural language into data, rendered inline as chat messages.

```
User types query in chat
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Chat Request Pipeline                      │
│                                                             │
│  E-01: Language Detection & Routing (locale, lang detect)  │
│     │                                                       │
│     ▼                                                       │
│  E-02: Query Translation (Arabic → English intent, if AR)  │
│     │                                                       │
│     ▼                                                       │
│  A-04: Domain Terminology Normalisation                     │
│        "footfall" → patient_visits                         │
│        "outstanding" → accounts_receivable                 │
│     │                                                       │
│     ▼                                                       │
│  A-06: Confidence Pre-Checker (intent ambiguity detection) │
│        If confidence < 0.70 → ask clarifying question      │
│     │                                                       │
│     ▼                                                       │
│  A-01: Text-to-SQL Agent (Genkit Flow + LangChain)         │
│        Schema context from pgvector cache                  │
│        Metric definitions from MetricFlow registry         │
│        SQL Validator: block any non-SELECT DML             │
│        PostgreSQL executor (medisync_readonly role)        │
│     │                                                       │
│     ├─── On DB error ──► A-02: Self-Correction Agent       │
│     │                         (rewrite + retry, max 3x)   │
│     ▼                                                       │
│  A-03: Visualisation Router                                 │
│        trend query     → line chart ECharts config         │
│        comparison      → bar chart                         │
│        breakdown       → pie chart                         │
│        raw data        → table                             │
│     │                                                       │
│     ▼                                                       │
│  E-03: Localised Response Formatter                         │
│        Numbers, dates, currencies in user's locale         │
│     │                                                       │
│     ▼                                                       │
│  CopilotKit streaming response → Chart/Table widget        │
└─────────────────────────────────────────────────────────────┘
```

**Key security guardrails (A-01):**
- SQL Injection Guard — parameterised inputs; all user strings treated as bind parameters
- Read-only enforcement — OPA policy `medisync.bi.read_only`; separate Postgres role
- Hallucination Guard — off-topic classifier; deflects non-business questions
- Column masking — OPA strips patient PII and cost-price columns based on `user_role`
- HITL Gate — queries with confidence < 0.70 or touching PII tables are routed to manual review queue

### 5.2 Module B — AI Accountant

Automates bookkeeping from raw document upload through to human-approved Tally posting.

```
Document Upload (PDF / Image / Excel / Handwritten Scan)
         │
         ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Document Processing Pipeline                    │
│                                                                  │
│  B-01: Document Classification Agent                            │
│        Invoice | Bank Statement | Bill | Receipt | Tax Doc      │
│     │                                                            │
│     ▼                                                            │
│  B-02: OCR Extraction Agent (PaddleOCR + Go Service)            │
│        Fields: amount, vendor, date, invoice#, tax amount       │
│        Confidence score per field; low-conf → HITL queue        │
│     │                                                            │
│     ├── [handwritten] ──► B-03: Handwriting Recognition Agent  │
│     │                                                            │
│     ▼                                                            │
│  B-04: Vendor Matching Agent                                    │
│        Match to Tally vendor master; create new if unrecognised │
│     │                                                            │
│     ▼                                                            │
│  B-07: Duplicate Invoice Detection Agent                        │
│        SHA-256 hash of (amount + vendor + date) → dupe check    │
│     │                                                            │
│     ▼                                                            │
│  B-05: Ledger Mapping Agent                                     │
│        AI suggests GL ledger with confidence score             │
│        Learns from user correction feedback                     │
│     │                                                            │
│  B-06: Sub-Ledger & Cost Centre Assignment Agent                │
│     │                                                            │
│     ▼                                                            │
│  B-08: Approval Workflow Agent                  ← HITL ALWAYS  │
│        Accountant → Manager → Finance → Posted                  │
│        Stale approval reminders via Notification Dispatcher     │
│     │                                                            │
│     ▼  (On finance_head explicit "Sync" click)                  │
│  B-09: Tally Sync Agent                         ← HITL ALWAYS  │
│        OPA gate: finance_head or accountant_lead + approved wf  │
│        Self-approval blocked                                    │
│        TDL XML → Tally HTTP POST                               │
│        Pre-sync validation: dupe detection, ledger availability │
│        Post-sync audit log entry (B-14)                         │
└──────────────────────────────────────────────────────────────────┘
```

**Bank Reconciliation sub-flow (B-10):**
```
Bank Statement Upload
     │
     ├─► B-02: OCR → extract bank transaction rows
     │
     ├─► B-10: Bank Reconciliation Agent
     │         Match bank rows ↔ Tally ledger entries
     │         Match criteria: amount + date (±3 days) + description similarity
     │         Confidence scoring per matched pair
     │         Unmatched → HITL review queue
     │
     └─► B-11: Outstanding Items Agent
               Aged payables/receivables report
               (0–7, 8–30, 31–90, 90+ days buckets)
```

### 5.3 Module C — Easy Reports

```
Report Request (user picks pre-built template OR builds custom)
         │
         ├──► Pre-built: C-01 Report Generator Agent
         │              P&L | Balance Sheet | Cash Flow | Aging | Inventory
         │
         ├──► Custom:   Zero-code drag-and-drop builder
         │              C-04 Custom Metric Formula Agent
         │              → Pydantic-validated formula evaluation
         │
         └──► Multi-company: C-02 Consolidation Agent
                              Merge N Tally instances
                              Eliminate inter-company transactions
                              Produce group P&L / Balance Sheet

Output Rendering
    │
    ├── HTML table → Apache ECharts widget (web)
    ├── PDF → WeasyPrint + Arabic fonts + RTL layout (E-04)
    ├── Excel → excelize (Go library)
    └── Email → C-03 Scheduling Agent → Notification Dispatcher

Security
    └── C-05 Row/Column Security Agent
              OPA row-level filter per user.department
              Column masking for cost/margin (manager/viewer roles)
```

### 5.4 Module D — Advanced Search-Driven Analytics

```
Natural Language Search Query
         │
         ▼
┌──────────────────────────────────────────────────────────────────┐
│                   Search Analytics Pipeline                        │
│                                                                  │
│  D-02: Entity Recognition Agent                                 │
│        Identify: doctors, drugs, suppliers, locations           │
│        Resolve to DB record IDs                                 │
│     │                                                            │
│     ▼                                                            │
│  D-01: NL Search Agent                                          │
│        Vector embeddings (pgvector) for semantic matching       │
│        Auto-complete + query history                            │
│        Ranked data results                                      │
│     │                                                            │
│     ▼                                                            │
│  D-03: Multi-Step Conversational Analysis Agent                 │
│        Decomposes complex multi-part questions                  │
│        Chains sequential SQL sub-queries                        │
│     │                                                            │
│     ▼                                                            │
│  D-04: Autonomous AI Analyst (Spotter) Agent                    │
│        Retrieves → Analyses → Compares → Forecasts              │
│        Proactively surfaces patterns without explicit query     │
│     │                                                            │
│     ▼                                                            │
│  D-08: Prescriptive Recommendations Agent                       │
│        Root-cause analysis                                      │
│        Quantified action recommendations with impact estimates  │
│     │                                                            │
│     ▼                                                            │
│  D-06: Dashboard Auto-Generation Agent                          │
│        Intelligent chart-type selection                         │
│        Layout optimisation                                      │
│        Data-storytelling annotations                            │
└──────────────────────────────────────────────────────────────────┘

Proactive Background Services:
  D-05: Deep Research Agent (scheduled, scans all data)
  D-07: Insight Discovery & Prioritisation Agent
  D-13: Scheduled Autonomous Monitoring Agent
  D-10: Insight-to-Action Workflow Agent (HITL: finance_head)
```

---

## 6. AI Agent Ecosystem

### 6.1 Agent Inventory

**Module A — Conversational BI (13 agents)**

| ID | Agent | Type | Complexity | HITL | Phase |
|---|---|---|---|---|---|
| A-01 | Text-to-SQL | Reactive | L2 | No | 2 |
| A-02 | SQL Self-Correction | Reactive | L2 | No | 2 |
| A-03 | Visualisation Routing | Reactive | L1 | No | 2 |
| A-04 | Domain Terminology Normaliser | Reactive | L1 | No | 2 |
| A-05 | Hallucination Guard | Reactive | L1 | No | 2 |
| A-06 | Confidence Scoring | Reactive | L1 | Yes (low conf.) | 2 |
| A-07 | Drill-Down Context | Reactive | L2 | No | 3 |
| A-08 | Multi-Period Comparison | Reactive | L2 | No | 3 |
| A-09 | Report Scheduling | Proactive | L2 | No | 3 |
| A-10 | KPI Alert | Proactive | L2 | No | 3 |
| A-11 | Chart-to-Dashboard Pin | Reactive | L1 | No | 3 |
| A-12 | Trend Forecasting | Reactive | L3 | No | 7 |
| A-13 | Anomaly Detection | Proactive | L3 | No | 7 |

**Module B — AI Accountant (16 agents)**

| ID | Agent | Type | Complexity | HITL | Phase |
|---|---|---|---|---|---|
| B-01 | Document Classification | Reactive | L1 | No | 4 |
| B-02 | OCR Extraction | Reactive | L2 | Yes (low conf.) | 4 |
| B-03 | Handwriting Recognition | Reactive | L2 | Yes | 4 |
| B-04 | Vendor Matching | Reactive | L2 | Yes (new vendor) | 5 |
| B-05 | Ledger Mapping | Reactive | L2 | Yes | 5 |
| B-06 | Sub-Ledger & Cost Centre Assignment | Reactive | L1 | Yes | 5 |
| B-07 | Duplicate Invoice Detection | Reactive | L2 | Yes | 5 |
| B-08 | Approval Workflow | Reactive/Proactive | L2 | Yes (always) | 5 |
| B-09 | Tally Sync | Reactive | L2 | Yes (always) | 6 |
| B-10 | Bank Reconciliation | Reactive | L2 | Yes (unmatched) | 7 |
| B-11 | Outstanding Items | Reactive | L1 | No | 7 |
| B-12 | Expense Categorisation | Reactive | L1 | Yes | 7 |
| B-13 | Tax Compliance | Reactive | L2 | No | 7 |
| B-14 | Audit Trail Logger | Reactive | L1 | No | 6 |
| B-15 | Cash Flow Forecasting | Reactive | L3 | No | 7 |
| B-16 | Multi-Entity Tally Manager | Reactive | L2 | No | 6 |

**Module C — Easy Reports (8 agents)**

| ID | Agent | Type | Complexity | HITL | Phase |
|---|---|---|---|---|---|
| C-01 | Pre-Built Report Generator | Reactive | L1 | No | 8 |
| C-02 | Multi-Company Consolidation | Reactive | L2 | No | 8 |
| C-03 | Report Scheduling & Distribution | Proactive | L2 | No | 9 |
| C-04 | Custom Metric Formula | Reactive | L2 | No | 10 |
| C-05 | Row/Column Security Enforcement | Reactive | L1 | No | 10 |
| C-06 | Data Quality Validation | Proactive | L2 | No | 1 |
| C-07 | Budget vs. Actual Variance | Reactive | L2 | No | 8 |
| C-08 | Inventory Aging & Reorder | Reactive | L2 | Yes (reorder) | 8 |

**Module D — Advanced Search Analytics (14 agents)**

| ID | Agent | Type | Complexity | HITL | Phase |
|---|---|---|---|---|---|
| D-01 | Natural Language Search | Reactive | L2 | No | 13 |
| D-02 | Entity Recognition | Reactive | L1 | No | 13 |
| D-03 | Multi-Step Conversational Analysis | Reactive | L3 | No | 14 |
| D-04 | Autonomous AI Analyst (Spotter) | Proactive | L3 | No | 14 |
| D-05 | Deep Research | Proactive | L3 | No | 14 |
| D-06 | Dashboard Auto-Generation | Reactive | L3 | No | 15 |
| D-07 | Insight Discovery & Prioritisation | Proactive | L3 | No | 15 |
| D-08 | Prescriptive Recommendations | Proactive | L3 | No | 15 |
| D-09 | Semantic Layer Management | Reactive | L2 | Yes (governance) | 13 |
| D-10 | Insight-to-Action Workflow | Proactive | L3 | Yes (always) | 14 |
| D-11 | Code Generation (SpotterCode) | Reactive | L2 | No | 16 |
| D-12 | Federated Query Optimisation | Reactive | L2 | No | 16 |
| D-13 | Scheduled Autonomous Monitoring | Proactive | L2 | No | 14 |
| D-14 | Voice/Mobile Search | Reactive | L2 | No | 15 |

**Module E — Language & Localisation (7 agents)**

| ID | Agent | Type | Complexity | HITL | Phase |
|---|---|---|---|---|---|
| E-01 | Language Detection & Routing | Reactive | L1 | No | 2 |
| E-02 | Query Translation (AR → EN intent) | Reactive | L1 | No | 2 |
| E-03 | Localised Response Formatter | Reactive | L1 | No | 2 |
| E-04 | Multilingual Report Generator | Reactive | L2 | No | 5 |
| E-05 | Translation Coverage Guard (CI) | Proactive | L1 | No | 1 |
| E-06 | Multilingual Notification | Proactive | L1 | No | 6 |

**Total: 58 agents across 5 modules**


### 6.2 Inter-Agent Communication (A2A Protocol)

All agent-to-agent communication follows the **Google A2A (Agent-to-Agent) Protocol**, providing standardised discovery, capability advertising, and task delegation.

```
Agent A needs to delegate to Agent B
         │
         ▼
┌─────────────────────────────────────────┐
│         A2A Discovery Service            │
│  GET /.well-known/agent.json            │
│  Returns: AgentCard (capabilities,      │
│           skills, endpoints, auth)      │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│         A2A Task Delegation             │
│  POST /tasks/send                       │
│  Body: { task_id, message, context }   │
│  Auth: Keycloak service account JWT     │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│         A2A Streaming Response          │
│  SSE or WebSocket streaming             │
│  Status: submitted → working → done    │
│  Artifacts: structured result payload  │
└─────────────────────────────────────────┘
```

Key agent collaboration chains:
- `A-01 → A-02 → A-03` — Core BI pipeline (SQL → self-correct → visualise)
- `E-01 → E-02 → A-01 → E-03` — Full multilingual BI query pipeline
- `B-02 → B-05 → B-08 → B-09` — Full document-to-Tally pipeline
- `D-04 → D-08 → D-10` — Autonomous insight-to-action pipeline
- `E-07 → A-04` — Glossary sync updates domain normaliser on deploy

### 6.3 Cross-Cutting Agent Services

These are shared services consumed by multiple agents, not standalone agents:

| Service | Consumers | Purpose |
|---|---|---|
| **Confidence Scoring Service** | A-06, B-02, B-05, B-10, D-04 | Standardised 0–100% score + routing logic |
| **Audit Log Writer** | B-08, B-09, B-14, D-10, C-05 | Immutable, append-only audit trail |
| **Schema Context Cache** | A-01, A-02, D-01, D-03, D-12 | Pre-loaded DB schema + semantic metadata serialised into `pgvector` embeddings |
| **OPA Policy Engine** | All agents | RBAC enforcement, read-only guardrails, Tally write-back policies |
| **Notification Dispatcher** | A-10, B-08, C-03, D-13, E-06 | Email / SMS / in-app / Slack fanout |
| **Semantic Layer Registry (MetricFlow)** | A-01, D-01, D-03, D-09 | Governed metric and dimension definitions |
| **Language & Locale Service** | All chat agents | Locale detection, cross-lingual intent normalisation, response formatting |

### 6.4 Agent Complexity Tiers

| Tier | Description | Examples |
|---|---|---|
| **L1** — Single-step | One tool call or transformation; deterministic output | A-03 (chart routing), A-05 (hallucination guard), E-01 (locale detect) |
| **L2** — Multi-step | Sequential tool chain with branching; recovers from errors | A-01 (text-to-SQL), B-02 (OCR), B-09 (Tally sync) |
| **L3** — Autonomous / long-horizon | Multi-step, self-directed, may spawn sub-agents; non-deterministic | D-04 (autonomous analyst), D-05 (deep research), A-12 (forecasting) |

---

## 7. Data Architecture

### 7.1 Warehouse Schema Overview

All warehouse tables carry standard audit columns:

```sql
_source          VARCHAR(50)    -- 'tally' | 'hims'
_source_id       VARCHAR(255)   -- natural key from source
_synced_at       TIMESTAMPTZ    -- when ETL last touched this row
_created_at      TIMESTAMPTZ
_updated_at      TIMESTAMPTZ
```

**Core fact tables:**
```
fact_vouchers          (Tally) — voucher_id, voucher_type, ledger_id, amount, date, ...
fact_appointments      (HIMS)  — appt_id, patient_id, doctor_id, status, billing_id, ...
fact_billing           (HIMS)  — bill_id, patient_id, amount, payment_mode, date, ...
fact_pharmacy_disp     (HIMS)  — disp_id, drug_id, quantity, patient_id, date, amount, ...
fact_stock_movements   (Tally) — movement_id, item_id, qty_in, qty_out, date, ...
```

**Locale support columns (app schema):**
```sql
users.locale            VARCHAR(10) DEFAULT 'en'   -- 'en' | 'ar'
users.calendar_system   VARCHAR(20) DEFAULT 'gregorian'
scheduled_reports.locale VARCHAR(10)
audit_log.locale        VARCHAR(10)
```

### 7.2 Vector Storage

**pgvector** (default) or **Milvus** (high-volume alternative):

| Collection | Content | Used By |
|---|---|---|
| `schema_embeddings` | Serialised DB schema + column descriptions | A-01, A-02, D-01 |
| `metric_embeddings` | Semantic descriptions of governed metrics | A-01, D-09 |
| `query_history` | Past user queries with result quality scores | D-01, A-01 (few-shot examples) |
| `document_embeddings` | Extracted text from uploaded documents | B-05 (ledger mapping context) |

### 7.3 Semantic Layer

**MetricFlow** governs all business metric definitions. Metrics are defined once and reused across all agents and reports:

```yaml
# Example MetricFlow metric definition
metric:
  name: clinic_revenue
  type: simple
  label: Clinic Revenue
  type_params:
    measure: billing_amount
  filter: |
    {{ Dimension('billing__department') }} = 'clinic'

metric:
  name: pharmacy_margin_pct
  type: ratio
  label: Pharmacy Margin %
  type_params:
    numerator: pharmacy_gross_profit
    denominator: pharmacy_revenue
```

Metric lineage and version changes are tracked in the `semantic.metric_definitions` table. D-09 (Semantic Layer Management Agent) enforces governance and propagates changes.

### 7.4 Caching Strategy

**Redis** is used at multiple levels:

| Cache | TTL | Content |
|---|---|---|
| Schema context cache | 1 hour | Pre-serialised Postgres schema JSON for LLM injection |
| Query result cache | 15 minutes | Identical query+params hash → cached DataFrame |
| Dashboard widget cache | 5 minutes | Rendered chart configs for auto-refreshing dashboards |
| User session cache | 8 hours | JWT decoded claims for fast authz decisions |
| Metric definition cache | 1 hour | MetricFlow metric registry JSON |

Cache invalidation is event-driven via NATS: `etl.sync.completed` events purge affected caches.

---

## 8. Security & Governance Architecture

### 8.1 Identity & Authentication (Keycloak)

Keycloak (Apache-2.0) is the sole identity provider:

```
Login Request
     │
     ▼
Keycloak (OIDC/SAML)
     │
     ├── 2FA (TOTP mandatory for finance_head, admin, accountant_lead)
     ├── JWT issued: { user_id, roles[], department, cost_centres[], locale }
     ├── Session: 8-hour access token, 30-day refresh token
     └── Force logout on role change (session invalidation via Redis)

Service Accounts per agent (minimum privilege):
     sa-bi-agent          → read:warehouse, read:semantic_layer
     sa-ocr-agent         → read:documents, write:extraction_queue
     sa-mapping-agent     → read:tally_coa, read:embeddings
     sa-approval-agent    → read:pending, write:approval_events
     sa-tally-sync        → execute:tally_sync (OPA gate required)
     sa-scheduler         → read:warehouse, write:notification_queue
```

### 8.2 Authorization — Policy as Code (OPA)

All authorization is delegated to Open Policy Agent (OPA) as a sidecar:

```
Every Agent API Request
    │
    ▼
OPA Sidecar: POST /v1/data/medisync/authz/allow
    Input: { user: JWT claims, action, resource, sql? }
    Policies:
        medisync.bi.read_only      → blocks any DML in SQL queries
        medisync.tally             → finance_head + approved workflow + no self-approval
        medisync.data.row_filter   → adds department filter for non-admin users
        medisync.columns.mask      → strips PII / cost columns by role
    Output: { allow: bool, reason: string, row_filter: {} }
```

### 8.3 Intelligence Plane vs. Action Plane

```
┌──────────────────────────────────────────┐
│          INTELLIGENCE PLANE              │
│  All AI agents + Data Warehouse          │
│  Postgres role: medisync_readonly        │
│  OPA: block all DML                      │
│  → Can never corrupt financial data      │
└─────────────────────┬────────────────────┘
                      │ READ ONLY
                      ▼
              [ Data Warehouse ]

┌──────────────────────────────────────────┐
│          ACTION PLANE                    │
│  B-08 Approval Workflow → B-09 Tally     │
│  OPA: finance_head role required         │
│  Self-approval blocked                   │
│  Immutable audit log on every action     │
└─────────────────────┬────────────────────┘
                      │ TDL XML HTTP
                      ▼
              [ Tally ERP ]
```

These planes are **architecturally separate**. The only bridge is the human-approved workflow gate — no AI agent can autonomously write to Tally.

### 8.4 Role Definitions

| Role | BI Query | Upload Docs | Approve Txns | Sync to Tally | 2FA Required |
|---|:---:|:---:|:---:|:---:|:---:|
| `admin` | ✅ All | ✅ | ✅ | ✅ | ✅ |
| `finance_head` | ✅ All | ✅ | ✅ | ✅ | ✅ |
| `accountant_lead` | ✅ Dept | ✅ | ✅ | ✅ | ✅ |
| `accountant` | ✅ Dept | ✅ | 1st-level | ❌ | — |
| `manager` | ✅ Dept | ❌ | ❌ | ❌ | — |
| `pharmacy_manager` | ✅ Pharm | ❌ | ❌ | ❌ | — |
| `analyst` | ✅ All | ❌ | ❌ | ❌ | — |
| `viewer` | ✅ Limited | ❌ | ❌ | ❌ | — |

### 8.5 Data Encryption

| Layer | Method |
|---|---|
| Data at rest | PostgreSQL TDE + disk-level AES-256 (LUKS on-premises) |
| Data in transit | TLS 1.3 on all connections; mTLS for service-to-service |
| Tally TDL payloads | TLS; credentials in environment secrets (never in source) |
| JWT tokens | RS256 signed; Keycloak JWKS public key validation |
| Document uploads | AES-256 at file level before storage; encryption key per tenant |
| Redis cache | AUTH password; bind to localhost only |
| Secrets management | HashiCorp Vault or OS-level secret store (not `.env` files in prod) |

### 8.6 HITL (Human-in-the-Loop) Gates

Agents that always require a human approval before irreversible actions:

| Agent | Trigger | Required Approver |
|---|---|---|
| B-08 (Approval Workflow) | Any transaction post into Tally queue | Accountant → Manager → Finance |
| B-09 (Tally Sync) | Tally journal/invoice/bill write-back | Finance Head (explicit "Sync Now" click) |
| B-10 (Bank Reconciliation) | Unmatched bank items over threshold | Accountant |
| D-10 (Insight-to-Action) | PO or any Tally write triggered by AI insight | Finance Head |
| E-07 (Glossary Sync) | Bilingual glossary term changes | Medical Advisor + Finance Advisor |

---

## 9. Internationalisation Architecture (i18n)

MediSync ships with first-class **English (LTR)** and **Arabic (RTL)** support from Phase 1.

### Locale Detection (Priority Order)

1. `user_preferences.locale` in Postgres (loaded into JWT `locale` claim at login)
2. `Accept-Language` HTTP header (browser/OS)
3. `?lang=ar` URL parameter (for email links and report share links)
4. Default: `en`

### Web (React + i18next)

```
frontend/public/locales/
  en/  common.json | dashboard.json | chat.json | reports.json
       accountant.json | ai-responses.json | notifications.json
  ar/  (mirrors en/ — all keys required; CI fails on gap — E-05)
```

- Namespaced lazy-loading per module
- RTL switching: `document.documentElement.dir = locale === 'ar' ? 'rtl' : 'ltr'`
- Tailwind logical properties (`ms-`, `me-`, `ps-`, `pe-`) for RTL-safe layout
- Playwright visual regression tests validate RTL layouts before each release

### Mobile (Flutter + ARB)

```
mobile/lib/l10n/
  app_en.arb    ← canonical English strings
  app_ar.arb    ← Arabic translations (compile-time type-safe)
```

- `Directionality` widget wraps app root
- `EdgeInsetsDirectional` for all padding/margin
- Hijri calendar support for users who prefer it (while keeping Gregorian for financial periods)

### AI Response Localisation

Every Genkit flow carries a mandatory `ResponseLanguageInstruction`:
```
System prompt injection (E-01/E-03):
"You must respond entirely in {{ locale }}. 
 Format all numbers as {{ number_format }}.
 Format all dates as {{ date_format }}.
 Format all currency values as {{ currency_format }}."
```

### Report Localisation (E-04)

- **PDF:** WeasyPrint + Cairo/Noto Sans Arabic fonts + `direction: rtl` CSS
- **Excel:** `excelize` Go library + RTL sheet direction setting
- Bilingual reports export both EN and AR pages in a single PDF when requested

---

## 10. Messaging & Event Bus

**NATS** (Apache-2.0) is the internal event bus for all async communication:

```
Key Topics:
  etl.sync.completed          → consumed by: cache invalidation, A-10 KPI Alert
  etl.sync.failed             → consumed by: A-10 KPI Alert, Notification Dispatcher
  document.uploaded           → consumed by: B-01 Document Classifier
  document.classified         → consumed by: B-02 OCR Extraction
  transaction.queued          → consumed by: B-08 Approval Workflow
  approval.completed          → consumed by: B-09 Tally Sync
  tally.sync.completed        → consumed by: B-14 Audit Log, Notification Dispatcher
  alert.kpi.threshold         → consumed by: Notification Dispatcher
  report.scheduled.due        → consumed by: C-03 Report Distribution
  insight.discovered          → consumed by: D-07 Insight Prioritisation
  anomaly.detected            → consumed by: Notification Dispatcher, D-10
```

NATS JetStream is used for persistent, at-least-once delivery on all write-back topics (`transaction.*`, `approval.*`, `tally.*`).

---

## 11. Observability Stack

**LGTM stack** (Loki + Grafana + Tempo + Mimir/Prometheus):

```
┌───────────────────────────────────────────────────────┐
│                   LGTM Observability                   │
│                                                       │
│  Grafana (dashboards + alerting)                      │
│      │                                                │
│      ├─── Prometheus (metrics: latency, error rates, │
│      │               query counts, agent calls)       │
│      ├─── Loki (structured JSON logs from all agents)│
│      └─── Genkit traces (per-flow span tracing)       │
└───────────────────────────────────────────────────────┘
```

**Key metrics tracked:**

| Metric | Target |
|---|---|
| Query latency P95 (A-01) | < 5 seconds |
| Dashboard load time | < 3 seconds |
| SQL accuracy (A-01) | ≥ 95% business intent accuracy |
| OCR field accuracy (B-02) | ≥ 95% for standard; ≥ 90% handwritten |
| Tally sync success rate | ≥ 99.5% |
| System uptime | ≥ 99.5% |
| Agent hallucination rate | < 1% |
| Translation coverage (E-05) | 100% (CI gate) |

---

## 12. Offline & Mobile Sync

**PowerSync** (Apache-2.0) provides offline-first capability for the Flutter app:

```
Flutter App ─── PowerSync Client SDK
                     │
                     ├── Local SQLite (on-device)
                     │   Stores: pinned dashboard configs,
                     │           last-N generated reports,
                     │           user preferences
                     │
                     └── PowerSync Service
                         Bidirectional sync with PostgreSQL
                         Reconnects and syncs delta on regain connectivity
                         Conflict resolution: server-wins for financial data
```

Offline capabilities:
- View pre-loaded dashboards and last-synced reports without internet
- Browse historical reports offline
- Online required for: new chat queries, document uploads, Tally sync

---

## 13. Technology Stack Reference

### Backend

| Component | Technology | License | Role |
|---|---|---|---|
| Core backend | Go 1.26 | BSD-3 | API server, ETL service, AI flow glue |
| HTTP router | go-chi/chi | MIT | REST routing, middleware |
| SQL driver | jmoiron/sqlx | MIT | PostgreSQL data access |
| Message broker | NATS / JetStream | Apache-2.0 | Async events, queues |
| Identity | Keycloak | Apache-2.0 | OIDC/SAML, JWT, 2FA |
| Authorization | Open Policy Agent | Apache-2.0 | Policy-as-code authz |
| ETL orchestration | Meltano | MIT | ELT pipelines |

### AI / ML

| Component | Technology | License | Role |
|---|---|---|---|
| AI flow orchestration | Google Genkit | Apache-2.0 | Type-safe AI pipelines |
| Multi-agent framework | Agent ADK | Apache-2.0 | Agent coordination |
| Inter-agent protocol | A2A Protocol | Apache-2.0 | Standardised agent communication |
| Local LLM serving | Ollama | MIT | Llama 4, Mistral, Gemma |
| Production LLM serving | vLLM | Apache-2.0 | High-throughput GPU inference |
| Semantic layer | MetricFlow | Apache-2.0 | Governed metric definitions |
| NL search | Genkit + pgvector | — | Vector semantic search |
| OCR | PaddleOCR | Apache-2.0 | Document text extraction |
| LLM providers | GPT-5.2 / Claude 4.6 / Gemini 3 Pro | — | Cloud LLM backends (swappable) |

### Data

| Component | Technology | License | Role |
|---|---|---|---|
| Data warehouse | PostgreSQL 18.2 | PostgreSQL | Primary analytics store |
| Vector search | pgvector | PostgreSQL | Embedding storage and similarity |
| Dedicated vector DB | Milvus | Apache-2.0 | High-volume vector indexing |
| Cache | Redis | BSD-3 | Query cache, session cache, task broker |
| Offline sync | PowerSync | Apache-2.0 | Mobile offline-first sync |

### Frontend

| Component | Technology | License | Role |
|---|---|---|---|
| Web framework | React 19.2.4 | MIT | Web application shell |
| Generative UI | CopilotKit | MIT | AI-powered dynamic UI orchestration |
| Build tool | Vite 7.3 | MIT | Development & production build |
| Charts | Apache ECharts | Apache-2.0 | All data visualisations |
| i18n (web) | i18next | MIT | EN/AR translations, RTL |
| Mobile | Flutter | BSD-3 | iOS/Android cross-platform app |
| i18n (mobile) | flutter_localizations + intl | BSD-3 | ARB-based type-safe strings |

### Observability

| Component | Technology | License | Role |
|---|---|---|---|
| Metrics | Prometheus | Apache-2.0 | Time-series metrics |
| Dashboards | Grafana | AGPL-3.0 | Unified observability UI |
| Logs | Loki | AGPL-3.0 | Log aggregation |
| Tracing | Genkit built-in | Apache-2.0 | Per-flow AI span tracing |

---

## 14. Deployment Topology

MediSync is deployed **on-premises** (private network, self-hosted) to satisfy healthcare data residency requirements.

```
┌────────────────────── On-Premises Network ──────────────────────────┐
│                                                                      │
│  ┌──────────────────┐   ┌──────────────────┐   ┌───────────────┐   │
│  │  App Server(s)   │   │  DB Server       │   │  GPU Server   │   │
│  │                  │   │                  │   │  (optional)   │   │
│  │  Go API          │   │  PostgreSQL 18.2 │   │               │   │
│  │  Genkit Flows    │   │  pgvector ext    │   │  vLLM         │   │
│  │  Agent Services  │   │  Redis           │   │  Ollama       │   │
│  │  Keycloak        │   │  Milvus          │   │               │   │
│  │  OPA             │   │                  │   │               │   │
│  │  NATS            │   └──────────────────┘   └───────────────┘   │
│  │  Meltano ETL     │                                              │
│  │  Grafana/Loki    │   ┌──────────────────────────────────────┐  │
│  └──────────────────┘   │  Source Systems                       │  │
│                         │  Tally ERP (existing)                 │  │
│                         │  HIMS Server (existing)               │  │
│                         └──────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
         │ HTTPS
         ▼
   Client Devices (Web browser, Flutter mobile app)
```

**Containerisation:** All services run in Docker containers, orchestrated with Docker Compose (single-node) or Kubernetes (multi-node for scale).

**Backup & Recovery:**
- PostgreSQL: pg_dump scheduled hourly; WAL archiving for point-in-time recovery
- RPO: < 1 hour | RTO: < 4 hours
- Redis: AOF persistence with hourly RDB snapshots

**High Availability (Phase 12+):**
- PostgreSQL streaming replication (primary + 1 replica)
- Go API behind load balancer (2+ instances)
- NATS cluster (3-node)
- Redis Sentinel for HA cache

---

## 15. Phased Delivery Roadmap

| Phase | Weeks | Deliverables | Key Agents |
|---|---|---|---|
| 1 | 1–3 | ETL infrastructure, DB schema, data validation | C-06 |
| 2 | 4–7 | Core AI agent, Text-to-SQL, financial analytics | A-01–A-06, E-01–E-03 |
| 3 | 8–11 | Chat UI, dashboards, charts, scheduled reports | A-07–A-11, E-07 |
| 4 | 12–15 | Document upload, OCR pipeline | B-01–B-03 |
| 5 | 16–19 | Ledger mapping, vendor matching, approval workflow | B-04–B-08, E-04 |
| 6 | 20–22 | Tally real-time sync, audit logging, multi-entity | B-09, B-14, B-16, E-06 |
| 7 | 23–26 | Bank reconciliation, cash flow forecasting, tax compliance | B-10–B-15, A-12–A-13 |
| 8 | 27–30 | Pre-built reports, consolidated dashboards | C-01, C-02, C-07, C-08 |
| 9 | 31–33 | Automated report scheduling & email delivery | C-03 |
| 10 | 34–36 | Zero-code report builder, RBAC data security | C-04, C-05 |
| 11 | 37–38 | Module integration, cross-linking, performance | All |
| 12 | 39–40 | UAT, security audit, production launch | All |
| 13 | 41–43 | Semantic layer, NL search infrastructure | D-01, D-02, D-09 |
| 14 | 44–46 | Autonomous agents, deep research, Spotter | D-03–D-05, D-10, D-13 |
| 15 | 47–48 | Auto-dashboarding, insight engine, prescriptive AI | D-06–D-08, D-14 |
| 16 | 49–50 | Code generation, embedded analytics, analyst studio | D-11, D-12 |
| 17 | 51–52 | Search governance, HIPAA/GDPR compliance | C-05 ext, D-09 |
| 18 | 53–54 | Full integration UAT, final security audit, production v2 | All |

---

## 16. Architecture Decision Records (ADRs)

### ADR-01: Go as Primary Backend Language
**Status:** Accepted  
**Decision:** Use Go for API server, ETL pipeline, and AI orchestration glue rather than Python.  
**Rationale:** Go's native concurrency (goroutines), low memory footprint, and type safety are important for ETL performance and concurrent AI agent handling. Genkit has a native Go SDK.  
**Trade-offs:** Smaller ML ecosystem than Python; mitigated by calling Python OCR/ML services as microservices.

### ADR-02: Genkit over LangChain for AI Flows
**Status:** Accepted  
**Decision:** Standardise on Google Genkit for all AI pipelines.  
**Rationale:** Genkit provides type-safe flows, built-in observability (spans per step), swappable LLM plugins, and streaming support. It integrates natively with Agent ADK for multi-agent scenarios.  
**Trade-offs:** Smaller community than LangChain; Google ecosystem alignment.

### ADR-03: Intelligence Plane / Action Plane Separation
**Status:** Accepted  
**Decision:** AI agents connect to PostgreSQL with a `SELECT`-only role. Tally write-backs happen exclusively through the human-approved B-08 → B-09 workflow over TDL XML.  
**Rationale:** Resolves the PRD §10 / §6.7.3 contradiction. Protects financial data integrity while enabling AI-assisted accounting.

### ADR-04: A2A Protocol for Agent Communication
**Status:** Accepted  
**Decision:** Use Google A2A Protocol for standardised inter-agent discovery and task delegation.  
**Rationale:** With 58 agents, a point-to-point call graph becomes unmanageable. A2A provides a standard envelope, capability advertising, and task lifecycle management.

### ADR-05: On-Premises Deployment (No Cloud-Mandatory Dependencies)
**Status:** Accepted  
**Decision:** All components can run on-premises; cloud LLM APIs are optional and swappable with local Ollama/vLLM.  
**Rationale:** Healthcare and financial data residency requirements; client preference for data sovereignty. LLM dependency externalised behind Genkit plugin interface.

### ADR-06: pgvector First, Milvus as Upgrade Path
**Status:** Accepted  
**Decision:** Start with pgvector (PostgreSQL extension) for vector storage. Migrate to Milvus if embedding volume exceeds PostgreSQL performance threshold.  
**Rationale:** Reduces operational complexity initially. pgvector is sufficient for schema embeddings and NL search in early phases.

### ADR-07: OPA as the Single Authorization System
**Status:** Accepted  
**Decision:** All permission checks, row-level filters, and column masks go through OPA Rego policies — not application code.  
**Rationale:** Centralised, auditable, version-controlled policies. Enables hot-reload of rules without service restart. Supports compliance requirements (HIPAA, GDPR).

---

*Document Version: 1.0 | Last Updated: February 19, 2026 | Owner: MediSync Engineering*
