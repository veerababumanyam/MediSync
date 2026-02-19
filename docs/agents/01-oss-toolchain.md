# MediSync — OSS Toolchain Map

**Version:** 2.1 | **Created:** February 19, 2026  
**Constraint:** Open-source only (OSI-approved licenses). Primary Backend: **Go**. AI Orchestration: **Genkit**.

---

## 1. License Compliance Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | OSI-approved open-source license — safe to use |
| ⚠️ | Source-available / non-OSI (ELv2, BSL, custom) — **excluded** |
| 🔄 | OSS replacement recommended |

---

## 2. Full OSS Stack — License-Verified

### 2.1 Backend & API Layer (Go)

| Tool | License | GitHub | Role in MediSync | Context |
|------|---------|--------|-----------------|---------|
| **Go** ✅ | BSD-3-Clause | [golang/go](https://github.com/golang/go) | Core Backend & ETL Service | Performance & Concurrency |
| **go-chi/chi** ✅ | MIT | [go-chi/chi](https://github.com/go-chi/chi) | HTTP Router | Lightweight, idiomatic |
| **sqlx** ✅ | MIT | [jmoiron/sqlx](https://github.com/jmoiron/sqlx) | Database Extensions | Enhanced `database/sql` |
| **Pydantic-Go** (or similar) | — | — | *Replaced by Go Structs + JSON tags* | Type-safe data modeling |

### 2.2 AI Orchestration & Frameworks

| Tool | License | GitHub | Role in MediSync | Note |
|------|---------|--------|-----------------|------|
| **Genkit** ✅ | Apache-2.0 | [firebase/genkit](https://github.com/firebase/genkit) | AI Orchestration (Go/TS) | Type-safe flows, observability |
| **Agent ADK** ✅ | Apache-2.0 | [google/agent-adk](https://github.com/google/agent-adk) | Multi-Agent Framework | Sophisticated agent orchestration |
| **Ollama** ✅ | MIT | [ollama/ollama](https://github.com/ollama/ollama) | Local LLM Serving | Llama 4, Mistral, Gemma |
| **vLLM** ✅ | Apache-2.0 | [vllm-project/vllm](https://github.com/vllm-project/vllm) | Production LLM Serving | High-throughput GPU |
| **Redamon** ✅ | MIT | [samugit83/redamon](https://github.com/samugit83/redamon) | Agentic Red Teaming | AI-powered offensive security |

> **Strategy:** Pivot from LangChain/Python to **Genkit** for unified AI workflows and better integration with Google Cloud / Firebase ecosystem when needed.

### 2.3 Web Dashboard Stack (React + CopilotKit)

| Tool | License | GitHub | Role in MediSync | Context |
|------|---------|--------|-----------------|---------|
| **React.js** ✅ | MIT | [facebook/react](https://github.com/facebook/react) | Host Framework | Component-based UI |
| **CopilotKit** ✅ | MIT | [CopilotKit](https://github.com/CopilotKit/CopilotKit) | GenUI Framework | AI-powered UI & Orchestration |
| **Vite** ✅ | MIT | [vitejs/vite](https://github.com/vitejs/vite) | Build Tool | Development & Build |
| **Apache ECharts** ✅ | Apache-2.0 | [apache/echarts](https://github.com/apache/echarts) | Primary Charting | Medical BI Viz |

### 2.4 Mobile App Stack (Flutter)

| Tool | License | GitHub | Role in MediSync | Note |
|------|---------|--------|-----------------|------|
| **Flutter** ✅ | BSD-3-Clause | [flutter/flutter](https://github.com/flutter/flutter) | Mobile App | iOS/Android cross-platform |

### 2.5 Enterprise Foundation (IAM & Events)

| Tool | License | GitHub | Role in MediSync | Note |
|------|---------|--------|-----------------|------|
| **Keycloak** ✅ | Apache-2.0 | [keycloak/keycloak](https://github.com/keycloak) | Identity & Auth | OIDC/SAML local IAM |
| **NATS** ✅ | Apache-2.0 | [nats-io/nats-server](https://github.com/nats-io) | Message Broker | Go-native high perf bus |
| **Meltano** ✅ | MIT | [meltano/meltano](https://github.com/meltano) | Data Ingestion | ELT Orchestration engine |

### 2.6 Observability & Logs (LGTM Stack)

| Tool | License | GitHub | Role in MediSync | Note |
|------|---------|--------|-----------------|------|
| **Grafana** ✅ | AGPL-3.0 | [grafana/grafana](https://github.com/grafana) | Metrics Dashboard | Unified viz for logs/metrics |
| **Prometheus** ✅ | Apache-2.0 | [prometheus/prometheus](https://github.com/prometheus) | Metrics Storage | Time-series DB |
| **Loki** ✅ | AGPL-3.0 | [grafana/loki](https://github.com/grafana/loki) | Log Aggregation | Horizontally scalable logs |

---

## 3. Data Infrastructure

| Tool | License | GitHub | Role in MediSync |
|------|---------|--------|-----------------|
| **PostgreSQL** ✅ | PostgreSQL | [postgres/postgres](https://github.com/postgres/postgres) | Primary Warehouse |
| **pgvector** ✅ | PostgreSQL | [pgvector/pgvector](https://github.com/pgvector/pgvector) | Vector Search |
| **Milvus** ✅ | Apache-2.0 | [milvus-io/milvus](https://github.com/milvus-io/milvus) | Dedicated Vector Storage |
| **PowerSync** ✅ | Apache-2.0 | [powersync-com](https://github.com/powersync-com) | Offline-first Sync |
| **Redis** ✅ | BSD-3-Clause | [redis/redis](https://github.com/redis/redis) | Caching & Task Broker |

---

## 4. Agent-to-Tool Mapping (Go/Genkit Context)

| Agent ID | Primary Tooling | Roles |
|----------|----------------|-------|
| **A-01 Text-to-SQL** | Go, Genkit, sqlx, Postgres | Natural language to curated SQL |
| **A-03 Visualization** | Go-echarts, Apache ECharts | Dynamic UI widget generation |
| **B-02 OCR Extraction** | PaddleOCR, Go Service | document-to-structured-data |
| **C-05 Security** | OPA (Open Policy Agent), **Redamon** | Offensive & defensive agentic security |
| **D-01 NL Search** | Genkit, pgvector | Vector embedding search |
| **E-01 Language Detection** | Genkit, `golang.org/x/text/language` | Locale detection, query lang classification |
| **E-02 Query Translation** | Genkit + GPT-5.2 / Gemini 3 Pro | Arabic→intent normalisation pre-SQL |
| **E-03 Response Formatter** | `golang.org/x/text/message`, Go `Intl` | Locale-correct number/date/currency output |
| **E-04 Report Generator (i18n)** | WeasyPrint, Cairo font, excelize, `golang.org/x/text` | RTL PDF + RTL Excel generation |
| **E-05 Translation Guard** | i18next-parser, Node.js CI script | Missing translation key enforcement |
| **Inter-Agent** | **A2A Protocol** | Standardized cross-agent communication |

---

## 5. Summary Stack View

```
┌─────────────────────────────────────────────────────────────┐
│                      MEDISYNC GO/GENKIT STACK               │
├──────────────────┬──────────────────────────────────────────┤
│ Layer            │ Primary Tools (OSS/OSI)                  │
├──────────────────┼──────────────────────────────────────────┤
│ Backend / API    │ Go (BSD-3), go-chi (MIT)                 │
│ AI Orchestration │ Genkit (Apache-2.0), Agent ADK           │
│ Agent Protocol   │ **A2A (Protocol)**                       │
│ Identity / Auth  │ Keycloak (Apache-2.0)                    │
│ Messaging        │ NATS (Apache-2.0)                        │
│ Data Ingestion   │ Meltano (MIT)                            │
│ Web Frontend     │ React (MIT) + CopilotKit (MIT)           │
│ Mobile Frontend  │ Flutter (BSD-3)                          │
│ Visualization    │ Apache ECharts (Apache-2.0)              │
│ Warehouse        │ PostgreSQL (OSI) + pgvector              │
│ Sync Engine      │ PowerSync (Apache-2.0)                   │
│ Observability    │ Grafana, Prometheus, Loki (LGTM)         │
└──────────────────┴──────────────────────────────────────────┘
```

