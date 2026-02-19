# Open Source Tools: Technical Reference

This document serves as a technical catalog of open-source tools and libraries integrated into or recommended for the MediSync ecosystem. These tools are selected based on their performance, scalability, and compatibility with the core tech stack (Go, Genkit, React, and Flutter).

## 1. Backend & AI Orchestration

| Tool | Category | Role in MediSync | Key Feature |
| :--- | :--- | :--- | :--- |
| **Go** | Language | Core Backend Service | Concurrency, High Performance |
| **Genkit** | AI Framework | AI Orchestration (Go/TS) | Primary framework for type-safe AI flows and RAG |
| **Agent ADK** | Agent Orchestration | Multi-agent Systems | Specialized toolkit for building collaborating agents |
| **A2A Protocol** | Agent Protocol | Inter-agent Communication | Standardized protocol for discovery and collaboration |
| **Redamon** | Security Framework | Agentic Red Teaming | AI-powered offensive security for agent mesh |
| **HTMX** | Library | Low-JS UI Interactivity | Server-side rendering with AJAX swaps |

### Go Ecosystem Recommendations
- **go-chi/chi**: Lightweight, idiomatic router for building HTTP services.
- **sqlx**: General-purpose extensions to `database/sql` for easier DB interactions.

---

## 2. AI-Driven Frontend & Visualization

MediSync utilizes a **GenUI-First** frontend strategy, where interfaces are dynamically constructed based on agent intelligence and medical context.

### Web Dashboard Stack (Physician Portal)
- **CopilotKit** (MIT)
  - **Role**: Core orchestration for AI-powered components and **Generative UI**.
  - **Implementation**: Runs on **React.js** and is built using **Vite** for high-performance HMR.
  - **Reference**: [CopilotKit Generative UI Guide](https://github.com/CopilotKit/generative-ui/blob/main/assets/generative-ui-guide.pdf)
- **Apache ECharts** (Apache-2.0)
  - **Role**: Primary library for heavy-duty BI dashboards and medical analytics.
  - **Integration**: [go-echarts](https://github.com/go-echarts/go-echarts) for server-side chart configuration.

### Mobile App Stack (Clinical Staff)
- **Flutter** (BSD-3)
  - **Role**: Cross-platform mobile access (iOS/Android) for real-time alerts and patient summaries.
  - **Note**: Complements the Web GenUI stack with native mobile performance.

### Specialized Visualization
- **Recharts**: Lightweight declarative JSX charts for simple responsive medical trends.
- **go-chart**: Server-side rendering (SVG/PNG) for embedded charts in automated PDF reports.
- **gonum/plot**: Statistical visualizations for deep clinical data research.

---

## 3. Enterprise Foundation & Security

MediSync requires robust, on-prem infrastructure to handle sensitive medical data and synchronize with external clinical systems.

### Identity & Access Management (IAM)
- **Keycloak** (Apache-2.0)
  - **Role**: Primary local Identity Provider for authentication (OIDC/SAML) and RBAC (Role-Based Access Control).
  - **Why**: Medical-grade security with full on-prem control over user sessions and audit logs.

### Data Ingestion & ETL Orchestration
- **Meltano** (MIT)
  - **Role**: Declarative data integration engine to sync data from **Tally**, **HIMS**, and third-party CSVs into PostgreSQL.
  - **Why**: CLI-driven, integrates seamlessly with a Go/GitOps workflow.

### Messaging & Event Streaming
- **NATS** (Apache-2.0)
  - **Role**: High-performance messaging system for real-time agent communication and WebSocket broadcasting.
  - **Why**: Extremely lightweight, written in Go, and ideal for local-first clusters.

### Observability Stack (LGTM)
- **Grafana** (AGPL-3.0) & **Prometheus** (Apache-2.0)
  - **Role**: Visualization of system metrics and long-term storage of time-series data.
- **Loki** (AGPL-3.0)
  - **Role**: Log aggregation for auditing AI accountant queries and system errors.

### Security & Secrets
- **HashiCorp Vault** (MPL-2.0 / BSL check required)
  - **Role**: Secure storage for LLM API keys and database credentials.

---

## 4. Data Persistence & Synchronization

| Tool | Category | Usage | Role | Feature |
| :--- | :--- | :--- | :--- | :--- |
| **PostgreSQL** | Database | Core Warehouse | Relational Data | ACID Compliance |
| **pgvector** | Extension | Vector Store | Semantic Search | Vector indexing for medical RAG |
| **Redis** | In-Memory | Caching | Real-time State | High-speed data retrieval |
| **Milvus** | Vector DB | Large Scale Search | Deep Vector Indexing | High-performance vector scaling |
| **PowerSync** | Sync Engine | Offline Persistence | Real-time Data Sync | Local-first data synchronization |

---

## 5. BI & Embedded Analytics

For advanced enterprise analytics without custom coding every widget:
- **Metabase** (AGPL-3.0)
  - **Role**: Primary BI tool for internal medical staff to explore raw HIMS data.
- **Apache Superset** (Apache-2.0)
  - **Role**: Enterprise-grade data exploration for high-volume pharmacist analytics.

---

## 6. Reporting & Document Automation

MediSync generates clinical-grade PDF and Excel reports for patients and regulatory bodies.
- **WeasyPrint** (BSD-3)
  - **Role**: Primary HTML-to-PDF engine with full RTL/CSS support.
- **Puppeteer** (Apache-2.0)
  - **Role**: Chrome-based PDF generation for complex JS-heavy charts (ECharts snapshotting).
- **Excelize** (BSD-3)
  - **Role**: Pure Go library for generating RTL-ready Excel spreadsheets.

---

## Technical Comparison: Charting Libraries

| Feature | Apache ECharts | Chart.js | D3.js | Plotly |
| :--- | :--- | :--- | :--- | :--- |
| **Complexity** | Medium | Low | High | Medium |
| **Flexibility** | High | Medium | Absolute | High |
| **Performance** | Excellent | Good | Variable | Good |
| **Frameworks** | All | All | All | Go|

