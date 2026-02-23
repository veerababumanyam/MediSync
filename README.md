<p align="center">
  <img src="/Users/v13478/.gemini/antigravity/brain/d60690be-5a3c-427c-8c17-03471f519a2d/AnySync_hero_banner_1771669104407.png" alt="AnySync Hero Banner" width="100%" />
</p>

<h1 align="center">AnySync</h1>

<p align="center">
  <strong>The Interconnected Medical Brain — AI-Powered Conversational BI & Intelligent Accounting</strong>
</p>

<p align="center">
  <a href="https://github.com/veerababumanyam/AnySync/stargazers"><img src="https://img.shields.io/github/stars/veerababumanyam/AnySync?style=for-the-badge&color=0056D2&logo=github" alt="Stars" /></a>
  <a href="https://github.com/veerababumanyam/AnySync/network/members"><img src="https://img.shields.io/github/forks/veerababumanyam/AnySync?style=for-the-badge&color=00E8C6&logo=github" alt="Forks" /></a>
  <a href="https://github.com/veerababumanyam/AnySync/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-Polyform_NonCommercial-blue?style=for-the-badge" alt="License" /></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Status-In_Development-blue?style=flat-square" alt="Status" />
  <img src="https://img.shields.io/badge/Version-1.0--alpha-orange?style=flat-square" alt="Version" />
  <img src="https://img.shields.io/badge/Backend-Go_1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/AI-Google_Genkit-4285F4?style=flat-square&logo=google&logoColor=white" alt="Genkit" />
  <img src="https://img.shields.io/badge/Framework-React_19.2-61DAFB?style=flat-square&logo=react&logoColor=black" alt="React" />
  <img src="https://img.shields.io/badge/CSS-Tailwind_4.2-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white" alt="Tailwind" />
  <img src="https://img.shields.io/badge/i18n-EN_%7C_AR-blueviolet?style=flat-square" alt="i18n" />
</p>

---

## 🌟 What is AnySync?

AnySync unifies the two core data systems of a healthcare-and-pharmacy business — **HIMS** (clinic operations) and **Tally ERP** (accounting) — into a single AI-powered platform. Instead of dumping CSVs and building spreadsheets, staff simply ask a question in plain language and receive instant charts, tables, and downloadable reports.

### 🚀 Core Modules

| Module | Purpose | Impact |
|:---|:---|:---|
| 🗣️ **Conversational BI** | Natural language queries for live clinical & financial data. | Zero-SQL analytics for everyone. |
| 🤖 **AI Accountant** | Automated OCR, ledger mapping, and one-click Tally sync. | 90% reduction in manual entry. |
| 📊 **Easy Reports** | Pre-built MIS reports and zero-code custom dashboards. | Executive-ready insights in seconds. |
| 🔍 **Search Analytics** | Autonomous AI analyst for trends and recommendations. | Proactive business intelligence. |

---

## 🧠 AI Agent Ecosystem

AnySync is powered by **58 specialized AI agents** orchestrated by the **Google A2A Protocol**.

<details>
<summary><b>Click to expand Agent Details</b></summary>

| ID | Agent | Purpose |
|---|---|---|
| **A-01** | Text-to-SQL | Converts NL to safe, read-only SQL queries. |
| **B-02** | OCR Extraction | High-accuracy field extraction from documents. |
| **B-05** | Ledger Mapping | AI-suggests correct Tally GL ledger. |
| **D-04** | Autonomous Analyst | Runs end-to-end analytical workflows. |
| **E-01** | i18n Router | Detects language and injects locale context. |

*Full agent backlog available in [docs/agents/00-agent-backlog.md](docs/agents/00-agent-backlog.md).*
</details>

---

## 🏗️ Architecture & Principles

AnySync is built for **Precision, Security, and Scalability.**

### 🛠️ The Stack

- **Backend**: Go (API, ETL, Orchestration)
- **Frontend**: React (Web) + Flutter (Mobile)
- **AI**: Google Genkit + Agent ADK + WebMCP
- **Data**: PostgreSQL + pgvector + Redis
- **Auth**: Keycloak + OPA (Policy as Code)

### 💎 Core Principles

- **Decoupled Data Plane**: ETL to a separate warehouse to protect production systems.
- **Read-Only Intelligence**: AI agents have `SELECT`-only access.
- **Human-in-the-Loop (HITL)**: All Tally write-backs require manual approval.
- **i18n by Default**: First-class English (LTR) and Arabic (RTL) support.

---

## 📈 Phased Roadmap

```mermaid
gantt
    title AnySync Development Roadmap
    dateFormat  YYYY-MM-DD
    section Foundation
    ETL & DB Infrastructure    :done, 2026-02-01, 2026-02-21
    section AI Core
    Core AI Chat & Analytics   :active, 2026-02-22, 2026-03-20
    section UI/UX
    Web & Mobile Dashboards    :2026-03-21, 2026-04-15
    section AI Accountant
    OCR & Ledger Mapping       :2026-04-16, 2026-05-30
    section Finalization
    UAT & Launch               :2026-06-01, 2026-07-15
```

---

## 🚦 Getting Started

### Prerequisites
- Go 1.26+, Node 24+, Flutter 3.42+
- Docker & Docker Compose
- PostgreSQL 18.2 w/ `pgvector`

### Quick Start
```bash
# 1. Start infrastructure
docker-compose up -d

# 2. Setup database
go run ./cmd/migrate

# 3. Start services
go run ./cmd/api & go run ./cmd/etl

# 4. Launch Frontend
cd frontend && npm install && npm run dev
```

---

## 📖 Documentation

Explore our comprehensive guides:
- 🏛️ [Architecture Overview](docs/ARCHITECTURE.md)
- 📝 [Product Requirements (PRD)](docs/PRD.md)
- 🎨 [Design System (DESIGN.md)](docs/DESIGN.md)
- 🌍 [i18n Architecture](docs/i18n-architecture.md)
- 🛡️ [PII Protection & Anonymization](docs/PII-PROTECTION.md)
- 🤖 [Agent Blueprints](docs/agents/BLUEPRINTS.md)

---

<p align="center">
  <strong>AnySync</strong> — Built for healthcare and accounting teams who deserve better than spreadsheets.
  <br/>
  <em>Last Updated: February 22, 2026</em>
</p>
