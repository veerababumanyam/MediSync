<p align="center">
  <img src="public/logo.png" alt="MediSync Logo" width="180" />
</p>

<h1 align="center">MediSync</h1>

<p align="center">
  <strong>AI-Powered Conversational BI &amp; Intelligent Accounting for Healthcare</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/status-in_development-blue" alt="Status" />
  <img src="https://img.shields.io/badge/version-1.0--alpha-orange" alt="Version" />
  <img src="https://img.shields.io/badge/backend-Go-00ADD8?logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/AI-Genkit_%2B_Agent_ADK-4285F4?logo=google" alt="Genkit" />
  <img src="https://img.shields.io/badge/frontend-React_%2B_Flutter-61DAFB?logo=react" alt="React" />
  <img src="https://img.shields.io/badge/license-OSI_Open_Source-green" alt="License" />
  <img src="https://img.shields.io/badge/i18n-EN_%7C_AR_(RTL)-blueviolet" alt="i18n" />
</p>

---

## What is MediSync?

MediSync unifies the two core data systems of a healthcare-and-pharmacy business — **HIMS** (clinic operations) and **Tally ERP** (accounting) — into a single AI-powered platform. Instead of dumping CSVs and building spreadsheets, staff simply ask a question in plain language and receive instant charts, tables, and downloadable reports.

The platform has three tightly integrated product modules, powered by a fleet of **58 AI agents**:

| Module | What it does |
|---|---|
| 🗣️ **Conversational BI Dashboard** | Chat with your data in natural language; get live charts and tables in seconds |
| 🤖 **AI Accountant** | Upload documents → OCR → AI ledger mapping → one-click sync to Tally |
| 📊 **Easy Reports** | Pre-built MIS reports, zero-code custom dashboards, automated email delivery |
| 🔍 **Advanced Search Analytics** | Autonomous AI analyst; prescriptive recommendations; deep research |

---

## Table of Contents

- [Features](#features)
- [AI Agent Ecosystem](#ai-agent-ecosystem)
- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Data Flow Diagrams](#data-flow-diagrams)
- [Security & Governance](#security--governance)
- [Internationalisation (EN / AR)](#internationalisation-en--ar)
- [Phased Roadmap](#phased-roadmap)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [KPI Targets](#kpi-targets)

---

## Features

### 🗣️ Conversational BI Dashboard

- **Natural Language Queries** — Ask *"Show me pharmacy revenue for Q1 vs Q2"* and receive a rendered bar chart inside the chat window
- **Generative UI** — Dynamic widgets (charts, tables, KPI cards) generated and streamed inside the chat via CopilotKit
- **Pre-defined Quick Prompts** — Carousel of instant-action buttons: *Today's Revenue*, *Pending Invoices*, *Low Stock Alerts*
- **Pin to Dashboard** — Pin any generated chart to a permanent, auto-refreshing main dashboard
- **Drill-Down** — Click any chart element to drill down to transaction-level detail
- **Multi-Period Comparison** — Month-over-month and year-over-year comparisons from a single natural language request
- **Export** — Download any table or chart as CSV, Excel (.xlsx), or PDF
- **Scheduled Reports** — Automated report delivery via email (daily / weekly / monthly / custom)
- **KPI Alerts** — Configurable metric thresholds; multi-channel notifications (in-app, email, SMS)

### 🤖 AI Accountant

- **Bulk Document Upload** — Drag-and-drop hundreds of PDFs, images, Excel files, or scanned handwritten invoices
- **AI OCR Extraction** — 95%+ accuracy for standard documents; field-level confidence scoring; HITL review for low-confidence extractions
- **Handwriting Recognition** — Specialized sub-agent for handwritten invoices and bills
- **Intelligent Ledger Mapping** — AI suggests the correct Tally GL ledger per transaction; learns from corrections
- **Duplicate Invoice Detection** — Flags potential duplicates before they are posted
- **Vendor Matching** — Auto-matches to Tally vendor master; creates new vendor records when needed
- **Approval Workflow** — Configurable multi-level approval chain (Accountant → Manager → Finance Head)
- **One-Click Tally Sync** — Push approved journal entries, purchase bills, and sales invoices directly into Tally
- **Real-Time Sync Dashboard** — Live connection status, sync history, manual trigger, automatic retry
- **Multi-Entity Support** — Manage and sync multiple Tally company instances
- **Bank Reconciliation** — Auto-match bank statement rows to Tally entries with confidence scoring
- **Outstanding Aging Reports** — Payables and receivables aged into 0–7, 8–30, 31–90, 90+ day buckets
- **Cash Flow Forecasting** — Project future cash position; what-if scenario modeling
- **Tax Compliance** — GST/VAT reconciliation, Input Tax Credit tracking, compliance checklist
- **Document Library** — Centralized searchable repository linked to transactions; retention policy enforcement

### 📊 Easy Reports

- **Pre-Built Report Library** — P&L, Balance Sheet, Cash Flow, Debtor Aging, Sales Analytics, Inventory, Tax, and 20+ more
- **Zero-Code Report Builder** — Drag-and-drop fields; formula builder for custom KPIs; no SQL required
- **Multi-Company Consolidation** — Merge financials across N Tally instances; eliminate inter-company transactions
- **Interactive Dashboards** — Executive, Sales, Finance, Inventory, and Operations dashboards with drill-down
- **KPI Scorecards** — Color-coded (RAG) status indicators; sparklines; MoM / YoY comparison
- **Budget vs. Actual Variance** — Variance analysis with drill-down; year-end forecasting
- **Automated Scheduling** — Email reports in PDF, Excel, HTML, or CSV to configurable distribution lists
- **Row & Column Security** — Role-based data filtering; department-level row visibility; cost/margin column masking
- **Custom Fields (UDF)** — Support for Tally User-Defined Fields in all reports
- **Mobile Reporting** — Responsive dashboards; offline pre-loaded report access

### 🔍 Advanced Search-Driven Analytics

- **Google-Like Data Search** — Auto-complete, spell-check, entity recognition, query history
- **Autonomous AI Analyst (Spotter)** — Retrieves, analyzes, compares, forecasts, and recommends without explicit instructions
- **Deep Research Agent** — Discovers hidden correlations, runs statistical analyses, produces structured research reports
- **Prescriptive Recommendations** — Specific, quantified actions with root-cause analysis and business impact estimates
- **AI Dashboard Auto-Generation** — Complete dashboards auto-generated from a search query in seconds
- **Anomaly Detection** — Continuous monitoring; highlights unusual data points with explanations
- **Semantic Layer (MetricFlow)** — Governed, version-controlled metric definitions shared across all agents and reports
- **Analyst Studio** — Python notebooks and ad-hoc SQL for data scientists building custom ML models
- **Embedded Analytics API** — REST/GraphQL APIs for embedding dashboards into clinic management software
- **Natural Language Code Generation** — *"Create a patient revenue dashboard"* → React components + SQL + styling
- **Voice Search** — Voice input on mobile; transcribe → NL search pipeline → mobile-formatted results

---

## AI Agent Ecosystem

MediSync is powered by **58 specialized AI agents** across 5 modules, coordinated by the **Google A2A Protocol**.

```
                         ┌─────────────────────────────────────┐
                         │     Agent Supervisor (A2A Protocol)  │
                         └────────────────┬────────────────────┘
              ┌──────────────┬────────────┼──────────────┬──────────────┐
              ▼              ▼            ▼              ▼              ▼
        Module A         Module B    Module C        Module D      Module E
     13 agents         16 agents    8 agents        14 agents      7 agents
   Conv. BI          AI Accountant  Easy Reports  Search Analytics  i18n
```

### Key Agents

| ID | Agent | Purpose |
|---|---|---|
| **A-01** | Text-to-SQL | Converts natural language to safe, read-only SQL; core engine of the BI dashboard |
| **A-02** | SQL Self-Correction | Detects query errors and automatically rewrites + retries (up to 3x) |
| **A-03** | Visualisation Routing | Chooses the optimal chart type (bar / line / pie / scatter / table) |
| **A-04** | Domain Terminology Normaliser | Maps healthcare and accounting synonyms to database field names |
| **A-06** | Confidence Scorer | Attaches 0–100% confidence to every AI answer; routes low-confidence to HITL |
| **A-10** | KPI Alert | Monitors metrics against thresholds; fires multi-channel notifications |
| **A-12** | Trend Forecasting | Extends historical time-series using ARIMA / Prophet / LLM-based models |
| **A-13** | Anomaly Detection | Scans all metrics on a schedule; surfaces statistical outliers |
| **B-02** | OCR Extraction | Extracts structured fields from documents with confidence scoring |
| **B-05** | Ledger Mapping | AI-suggests correct Tally GL ledger; learns from user corrections |
| **B-08** | Approval Workflow | Routes transactions through multi-level human approval chain |
| **B-09** | Tally Sync | Pushes approved data to Tally ERP via TDL XML — always human-gated |
| **B-10** | Bank Reconciliation | Matches bank statement rows to Tally entries; flags unmatched items |
| **B-15** | Cash Flow Forecasting | Projects future cash position; runs what-if scenarios |
| **C-02** | Multi-Company Consolidation | Merges financials from multiple Tally instances |
| **C-06** | Data Quality Validation | ETL-level integrity checks; alerts on anomalies |
| **D-04** | Autonomous AI Analyst | Runs full analytical workflows autonomously: retrieve → analyse → forecast → recommend |
| **D-05** | Deep Research | Discovers hidden patterns; runs regression, clustering, decomposition |
| **D-08** | Prescriptive Recommendations | Generates quantified action recommendations with impact estimates |
| **D-10** | Insight-to-Action Workflow | Converts AI insights into Tally business actions (always HITL-gated) |
| **E-01** | Language Detection & Routing | Detects query language; injects locale into all downstream AI flows |
| **E-02** | Query Translation | Translates Arabic queries to English intent before SQL generation |
| **E-03** | Localised Response Formatter | Formats numbers, dates, and currency in the user's active locale |
| **E-04** | Multilingual Report Generator | Generates RTL PDF and Excel reports in Arabic with correct fonts |

---

## Architecture Overview

```
╔═══════════════════════════════════════════════════════════════════════╗
║                      EXTERNAL DATA SOURCES                            ║
║    Tally ERP (TDL XML)              HIMS (REST API)                  ║
╚══════════════════╤════════════════════════════╤═══════════════════════╝
                   │                            │
                   ▼                            ▼
╔═══════════════════════════════════════════════════════════════════════╗
║             ETL / INGESTION LAYER  (Go + Meltano)                     ║
║   Incremental sync · Data validation (C-06) · NATS event publish     ║
╚══════════════════════════════════╤════════════════════════════════════╝
                                   │
                                   ▼
╔═══════════════════════════════════════════════════════════════════════╗
║        DATA WAREHOUSE  (PostgreSQL + pgvector + Redis)                ║
║   hims_analytics | tally_analytics | semantic | app | vectors        ║
║   medisync_readonly role — no AI agent can write to the warehouse    ║
╚══════════════════════════════════╤════════════════════════════════════╝
                                   │  READ ONLY
                                   ▼
╔═══════════════════════════════════════════════════════════════════════╗
║          AI ORCHESTRATION LAYER  (Genkit + Agent ADK + A2A)          ║
║   58 agents · OPA policy engine · Audit log · Notifications          ║
╚══════════════════════════════════╤════════════════════════════════════╝
                                   │
                                   ▼
╔═══════════════════════════════════════════════════════════════════════╗
║             API GATEWAY  (Go / go-chi)                                ║
║   Keycloak JWT · OPA authz · Rate limiting · Accept-Language · SSE   ║
╚══════════╤════════════════════════════════════════════════╤══════════╝
           │                                                │
           ▼                                                ▼
╔══════════════════════╗                     ╔══════════════════════════╗
║  React Web App        ║                     ║  Flutter Mobile App      ║
║  CopilotKit GenUI     ║                     ║  iOS + Android           ║
║  Apache ECharts       ║                     ║  PowerSync offline       ║
║  i18next (EN/AR RTL)  ║                     ║  ARB i18n (EN/AR RTL)    ║
╚══════════════════════╝                     ╚══════════════════════════╝

           ┌─────────────────────────────────────┐
           │  ACTION PLANE  (Tally write-back)    │
           │  B-08 Approval → B-09 Tally Sync    │
           │  OPA gate: finance_head only         │
           │  Human approval required — always   │
           └─────────────────────────────────────┘
```

### Architecture Principles

| Principle | What it means in practice |
|---|---|
| **Decoupled data plane** | ETL to a separate warehouse — HIMS & Tally are never hit by analytics queries |
| **Read-only intelligence plane** | All AI agents connect to Postgres with `SELECT`-only credentials |
| **HITL for all write-backs** | No AI can autonomously push data to Tally; humans must approve |
| **Policy as Code** | All authorization via OPA Rego — auditable, version-controlled, hot-reloadable |
| **Open-source only** | Every component has an OSI-approved license |
| **Go-first backend** | Performance, low memory, native concurrency for ETL and AI orchestration |
| **A2A Protocol** | Standardised agent-to-agent discovery and task delegation |
| **i18n by default** | Locale injected at API gateway and every AI flow — no retrofitting needed |

---

## Tech Stack

### Backend & Infrastructure

| Layer | Technology |
|---|---|
| Backend language | **Go 1.26** |
| HTTP router | go-chi/chi (MIT) |
| Message broker | **NATS / JetStream** (Apache-2.0) |
| Identity & Auth | **Keycloak** (Apache-2.0) — OIDC, JWT, 2FA |
| Authorization | **Open Policy Agent** (Apache-2.0) — Rego policies |
| ETL orchestration | **Meltano** (MIT) |
| Database | **PostgreSQL 18.2** with pgvector extension |
| Cache | **Redis** (BSD-3) |
| Offline sync | **PowerSync** (Apache-2.0) |
| Vector DB (alt.) | **Milvus** (Apache-2.0) |

### AI & ML

| Layer | Technology |
|---|---|
| AI flow orchestration | **Google Genkit** (Apache-2.0) |
| Multi-agent framework | **Agent ADK** (Apache-2.0) |
| Inter-agent protocol | **Google A2A Protocol** (Apache-2.0) |
| Semantic layer | **MetricFlow** (Apache-2.0) |
| Local LLM serving | **Ollama** (MIT) — Llama 4, Mistral, Gemma |
| GPU LLM serving | **vLLM** (Apache-2.0) |
| Cloud LLMs | GPT-5.2, Claude 4.6 Opus, Gemini 3 Pro (swappable via Genkit plugin) |
| OCR engine | **PaddleOCR** (Apache-2.0) |

### Frontend

| Layer | Technology |
|---|---|
| Web framework | **React 19.2.4** (MIT) |
| Generative UI | **CopilotKit** (MIT) |
| Build tool | **Vite 7.3** (MIT) |
| Charting | **Apache ECharts** (Apache-2.0) |
| Web i18n | **i18next + react-i18next** (MIT) |
| Mobile | **Flutter** (BSD-3) — iOS & Android |
| Mobile i18n | flutter_localizations + intl (ARB files) |
| Mobile offline | **PowerSync** SDK |

### Observability

| Layer | Technology |
|---|---|
| Metrics | **Prometheus** (Apache-2.0) |
| Dashboards & Alerts | **Grafana** (AGPL-3.0) |
| Log aggregation | **Loki** (AGPL-3.0) |
| AI flow tracing | Genkit built-in span tracing |

---

## Data Flow Diagrams

### Conversational BI Query Flow

```
User types query
      │
      ▼
E-01: Detect language + inject locale
      │
      ▼  [if Arabic]
E-02: Translate Arabic → English intent
      │
      ▼
A-04: Normalise domain terms
      │
      ▼
A-06: Confidence pre-check
      │                          │ confidence < 0.70
      │                          └─► HITL review queue
      ▼
A-01: Generate SQL (Genkit Flow)
  Schema context from pgvector
  Metric definitions from MetricFlow
  SQL Validator: block non-SELECT DML
  Execute via medisync_readonly role
      │                          │ DB error
      │                          └─► A-02: Self-correct (up to 3x)
      ▼
A-03: Route to chart type
      │
      ▼
E-03: Format locale-aware response
      │
      ▼
CopilotKit streaming → rendered chart widget
```

### Document-to-Tally Flow

```
Bulk document upload
      │
      ▼
B-01: Classify (Invoice / Bill / Bank Statement / ...)
      │
      ▼
B-02: OCR extraction + confidence scores
      │ [handwritten]
      └─► B-03: Handwriting recognition
      │
      ▼
B-07: Duplicate detection
      │
      ▼
B-04: Vendor matching
      │
      ▼
B-05: Ledger mapping (AI suggestion + confidence)
      │
      ▼
B-06: Sub-ledger + cost centre assignment
      │
      ▼ ◄── HITL GATE
B-08: Multi-level approval workflow
  Accountant → Manager → Finance Head
      │
      ▼ ◄── HITL GATE (explicit finance_head click)
B-09: Tally Sync via TDL XML HTTP
  OPA policy: finance_head + approved workflow + no self-approval
      │
      ▼
B-14: Immutable audit log entry
```

---

## Security & Governance

### Core Security Architecture

```
┌──────────────────────────────────────┐
│        INTELLIGENCE PLANE            │
│  AI Agents + Data Warehouse          │
│  Postgres: medisync_readonly role    │    READ ONLY
│  OPA: block all DML                  │ ──────────────► [ PostgreSQL ]
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│        ACTION PLANE                  │
│  Human-approved Tally write-backs    │    TDL XML
│  OPA: finance_head role required     │ ──────────────► [ Tally ERP ]
│  Self-approval blocked               │
└──────────────────────────────────────┘
```

### Roles & Permissions

| Role | BI Query | Upload Docs | Approve | Sync to Tally | 2FA |
|---|:---:|:---:|:---:|:---:|:---:|
| `admin` | ✅ All | ✅ | ✅ | ✅ | ✅ |
| `finance_head` | ✅ All | ✅ | ✅ | ✅ | ✅ |
| `accountant_lead` | ✅ Dept | ✅ | ✅ | ✅ | ✅ |
| `accountant` | ✅ Dept | ✅ | 1st-level | ❌ | — |
| `manager` | ✅ Dept | ❌ | ❌ | ❌ | — |
| `pharmacy_manager` | ✅ Pharm | ❌ | ❌ | ❌ | — |
| `analyst` | ✅ All | ❌ | ❌ | ❌ | — |
| `viewer` | ✅ Limited | ❌ | ❌ | ❌ | — |

### Key Security Features

- **OPA Policy as Code** — All authorization in Rego, version-controlled, hot-reloadable
- **Keycloak SSO + 2FA** — TOTP mandatory for Finance Head and Admin roles
- **Read-only DB role** — `medisync_readonly` has `GRANT SELECT` only; structural write is impossible at driver level
- **Column masking** — Patient PII and cost prices stripped from responses by role
- **Row-level security** — Users see only data for their department / cost centre
- **Immutable audit trail** — Append-only `audit_log` table; tracks who did what from which document
- **Self-approval blocked** — OPA Rego prevents the same user from submitting and approving
- **TLS 1.3** everywhere; mTLS for service-to-service calls
- **Encryption at rest** — AES-256 (LUKS) for PostgreSQL and uploaded documents

---

## Internationalisation (EN / AR)

MediSync ships with first-class **English (LTR)** and **Arabic (RTL)** from Phase 1.

```
Locale detection priority:
  1. user_preferences.locale (stored in Postgres, loaded into JWT)
  2. Accept-Language HTTP header
  3. ?lang=ar URL parameter
  4. Default: en
```

| Surface | Implementation |
|---|---|
| Web strings | i18next namespaced JSON; all `en` keys mirrored in `ar`; E-05 CI gate blocks merge on gap |
| RTL layout | Tailwind logical properties (`ms-`, `me-`, `ps-`, `pe-`); `dir="rtl"` on `<html>` |
| Mobile strings | flutter_localizations + ARB files; compile-time type-safe |
| Mobile RTL | `Directionality` widget; `EdgeInsetsDirectional` |
| AI responses | Locale injected as mandatory `ResponseLanguageInstruction` into every Genkit flow |
| PDF reports | WeasyPrint + Cairo/Noto Sans Arabic + CSS `direction: rtl` |
| Excel reports | excelize RTL sheet direction |
| Number / date formatting | `golang.org/x/text/message` and `Intl.*` APIs — no manual string concatenation |
| Hijri calendar | Phase 6 — toggle per user preference; financial periods stay Gregorian |

---

## Phased Roadmap

```
Phase  1  (Wk  1–3 ) ██░░░░░░░░  ETL infrastructure, DB schema, data validation
Phase  2  (Wk  4–7 ) ██░░░░░░░░  Core AI chat, Text-to-SQL, financial analytics
Phase  3  (Wk  8–11) ███░░░░░░░  Chat UI, dashboards, scheduled reports, KPI alerts
Phase  4  (Wk 12–15) ████░░░░░░  Document upload, OCR pipeline
Phase  5  (Wk 16–19) █████░░░░░  Ledger mapping, vendor matching, approval workflow
Phase  6  (Wk 20–22) █████░░░░░  One-click Tally sync, audit logging, multi-entity
Phase  7  (Wk 23–26) ██████░░░░  Bank reconciliation, cash flow, tax compliance
Phase  8  (Wk 27–30) ███████░░░  Pre-built reports, consolidated dashboards
Phase  9  (Wk 31–33) ███████░░░  Automated scheduling and email distribution
Phase 10  (Wk 34–36) ████████░░  Zero-code report builder, RBAC data security
Phase 11  (Wk 37–38) ████████░░  All-module integration, performance optimisation
Phase 12  (Wk 39–40) █████████░  UAT, security audit, production launch (v1.0)
Phase 13  (Wk 41–43) █████████░  Semantic layer, NL search infrastructure
Phase 14  (Wk 44–46) ██████████  Autonomous agents, Spotter, deep research
Phase 15  (Wk 47–48) ██████████  Auto-dashboarding, insight engine, prescriptive AI
Phase 16  (Wk 49–50) ██████████  Code generation, embedded analytics, Analyst Studio
Phase 17  (Wk 51–52) ██████████  Governance, HIPAA/GDPR compliance
Phase 18  (Wk 53–54) ██████████  Full integration UAT, final security audit, v2.0
```

---

## Getting Started

### Prerequisites

- Go 1.26
- Docker & Docker Compose
- Node.js 24 LTS (frontend build)
- Flutter 3.42 (mobile build)
- PostgreSQL 18.2 with `pgvector` extension
- Redis 8.6.0
- NATS Server 2.12.4

### Development Setup

```bash
# 1. Clone the repository
git clone https://github.com/your-org/medisync.git
cd medisync

# 2. Start infrastructure (Postgres, Redis, NATS, Keycloak, OPA)
docker-compose up -d

# 3. Run database migrations
go run ./cmd/migrate

# 4. Start the ETL service
go run ./cmd/etl

# 5. Start the API server
go run ./cmd/api

# 6. Start the frontend (web)
cd frontend && npm install && npm run dev

# 7. Start the mobile app
cd mobile && flutter pub get && flutter run
```

### Environment Variables

```bash
# Database
POSTGRES_DSN=postgres://medisync_app:password@localhost:5432/medisync

# AI Providers (at least one required)
GENKIT_GEMINI_API_KEY=your_gemini_key
GENKIT_OPENAI_API_KEY=your_openai_key
OLLAMA_BASE_URL=http://localhost:11434    # for local LLMs

# Identity
KEYCLOAK_BASE_URL=http://localhost:8080
KEYCLOAK_REALM=medisync
KEYCLOAK_CLIENT_ID=medisync-api

# Integrations
TALLY_HOST=http://localhost:9000          # Tally XML Gateway
HIMS_BASE_URL=https://hims.internal/api

# Cache & Messaging
REDIS_URL=redis://localhost:6379
NATS_URL=nats://localhost:4222
```

---

## Project Structure

```
medisync/
├── cmd/
│   ├── api/          ← Go API server entry point
│   ├── etl/          ← ETL service entry point
│   └── migrate/      ← Database migration runner
│
├── internal/
│   ├── agents/       ← 58 AI agent implementations (Genkit flows)
│   │   ├── module_a/ ← Conversational BI agents (A-01 – A-13)
│   │   ├── module_b/ ← AI Accountant agents (B-01 – B-16)
│   │   ├── module_c/ ← Easy Reports agents (C-01 – C-08)
│   │   ├── module_d/ ← Search Analytics agents (D-01 – D-14)
│   │   └── module_e/ ← Language & i18n agents (E-01 – E-07)
│   ├── api/          ← HTTP handlers, middleware, routing
│   ├── auth/         ← Keycloak JWT validation, OPA client
│   ├── etl/          ← Tally & HIMS connectors, transform pipeline
│   ├── warehouse/    ← PostgreSQL repository layer (sqlx)
│   ├── cache/        ← Redis client & cache strategies
│   └── notifications/ ← NATS-based notification dispatcher
│
├── policies/
│   ├── bi_readonly.rego       ← OPA: enforce SELECT-only for BI
│   ├── tally_sync.rego        ← OPA: finance_head + approval gate
│   └── row_level_security.rego ← OPA: department-scoped data access
│
├── migrations/       ← SQL schema migrations
│
├── frontend/         ← React web application
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   └── hooks/
│   └── public/
│       └── locales/
│           ├── en/   ← English translation JSONs
│           └── ar/   ← Arabic translation JSONs
│
├── mobile/           ← Flutter mobile application
│   └── lib/
│       ├── l10n/
│       │   ├── app_en.arb
│       │   └── app_ar.arb
│       └── ...
│
├── docs/
│   ├── ARCHITECTURE.md        ← Full system architecture document
│   ├── PRD.md                 ← Product Requirements Document
│   ├── DESIGN.md              ← Design system (colors, typography, components)
│   ├── i18n-architecture.md   ← i18n / localisation architecture
│   ├── OpenSourceTools.md     ← OSS toolchain reference
│   └── agents/
│       ├── BLUEPRINTS.md      ← Detailed agent blueprints
│       ├── 00-agent-backlog.md
│       ├── 01-oss-toolchain.md
│       ├── 03-governance-security.md
│       └── specs/             ← Per-agent technical specifications
│
└── public/
    └── logo.png
```

---

## Documentation

| Document | Description |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Comprehensive system architecture — layers, modules, agents, data model, security, ADRs |
| [docs/PRD.md](docs/PRD.md) | Product Requirements Document — all features, user stories, phased timeline |
| [docs/DESIGN.md](docs/DESIGN.md) | Design system — color palette, typography, glassmorphism components, Generative UI patterns |
| [docs/i18n-architecture.md](docs/i18n-architecture.md) | Full i18n architecture — locale detection, translation file structure, RTL, AI response localisation |
| [docs/agents/BLUEPRINTS.md](docs/agents/BLUEPRINTS.md) | Detailed blueprints for highest-priority agents — inputs, outputs, tool chains, guardrails, HITL gates |
| [docs/agents/00-agent-backlog.md](docs/agents/00-agent-backlog.md) | Master agent task inventory — all 58 agents with complexity, HITL, and phase assignments |
| [docs/agents/01-oss-toolchain.md](docs/agents/01-oss-toolchain.md) | OSS toolchain map — license-verified stack per layer |
| [docs/agents/03-governance-security.md](docs/agents/03-governance-security.md) | Governance & security design — OPA policies, RBAC, encryption, HITL gates |

---

## KPI Targets

| Metric | Target |
|---|---|
| Query accuracy (Text-to-SQL) | ≥ 95% business intent accuracy |
| Query latency (P95) | < 5 seconds |
| Dashboard load time | < 3 seconds |
| OCR field accuracy (standard docs) | ≥ 95% |
| OCR field accuracy (handwritten) | ≥ 90% |
| Tally sync success rate | ≥ 99.5% |
| System uptime | ≥ 99.5% |
| Agent hallucination rate | < 1% |
| Translation coverage (EN → AR) | 100% (CI-enforced) |
| Weekly active user rate (target) | 100% of management team ≥ 2×/week |
| Manual reporting time reduction | ≥ 90% |

---

## License

All components of the MediSync platform use **OSI-approved open-source licenses**. See [docs/agents/01-oss-toolchain.md](docs/agents/01-oss-toolchain.md) for the full license-verified stack.

---

<p align="center">
  <strong>MediSync</strong> — Built for healthcare and accounting teams who deserve better than spreadsheets.
  <br/>
  <em>Last Updated: February 19, 2026</em>
</p>
