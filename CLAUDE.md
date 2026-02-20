# MediSync — Claude AI Assistant Guide

**Project:** MediSync - AI-Powered Conversational BI & Intelligent Accounting for Healthcare
**Version:** 1.0-alpha
**Last Updated:** February 19, 2026

---

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Purpose** | AI-powered platform unifying HIMS (clinic operations) and Tally ERP (accounting) data |
| **Tech Stack** | Go 1.26 (backend), React 19 (web), Flutter 3.42 (mobile), PostgreSQL 18.2, Genkit + Agent ADK |
| **AI Agents** | 58 specialized agents across 5 modules (BI, Accounting, Reports, Analytics, i18n) |
| **Architecture** | Decoupled data plane, read-only intelligence plane, human-gated action plane |
| **i18n** | First-class English (LTR) and Arabic (RTL) support from Phase 1 |
| **Security** | Keycloak auth, OPA policy-as-code, HITL gates for all write-backs |

---

## 1. Project Overview

MediSync is an AI-powered, chat-based Business Intelligence platform for healthcare and pharmacy businesses. It solves the core problem of **siloed data** between operational systems (HIMS) and financial systems (Tally ERP) by:

1. **Conversational BI Dashboard** - Chat with your data in natural language; get live charts and tables
2. **AI Accountant** - OCR extraction → AI ledger mapping → one-click Tally sync
3. **Easy Reports** - Pre-built MIS reports and zero-code custom dashboards
4. **Advanced Search Analytics** - Autonomous AI analyst with prescriptive recommendations

### Core Philosophy

- **Decoupled data plane**: ETL to separate warehouse — HIMS & Tally are never hit by analytics queries
- **Read-only intelligence plane**: All AI agents connect with SELECT-only credentials
- **Human-in-the-loop for write-backs**: No AI can autonomously push data to Tally; humans must approve
- **Policy as code**: All authorization via OPA Rego — auditable, version-controlled
- **Open-source only**: Every component has an OSI-approved license

---

## 2. Architecture

### Layer Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                    EXTERNAL DATA SOURCES                        │
│           Tally ERP (TDL XML)     |     HIMS (REST API)         │
└────────────────────┬────────────────────────────┬────────────────┘
                     │                            │
                     ▼                            ▼
┌─────────────────────────────────────────────────────────────────┐
│              ETL / INGESTION LAYER (Go + Meltano)               │
│    Incremental sync · Data validation · NATS event publish      │
└────────────────────┬───────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│        DATA WAREHOUSE (PostgreSQL + pgvector + Redis)           │
│   medisync_readonly role — no AI agent can write                │
└────────────────────┬───────────────────────────────────────────┘
                     │  READ ONLY
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│         AI ORCHESTRATION LAYER (Genkit + Agent ADK + A2A)       │
│   58 agents · OPA policy engine · Audit log · Notifications      │
└────────────────────┬───────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────────┐
│               API GATEWAY (Go / go-chi)                         │
│   Keycloak JWT · OPA authz · Rate limiting · i18n routing       │
└────────────┬────────────────────────────────────┬───────────────┘
             │                                    │
             ▼                                    ▼
┌─────────────────────────┐      ┌───────────────────────────────┐
│   React Web App         │      │   Flutter Mobile App          │
│   CopilotKit GenUI      │      │   iOS + Android               │
│   Apache ECharts        │      │   PowerSync offline           │
│   i18next (EN/AR RTL)   │      │   ARB i18n (EN/AR RTL)        │
└─────────────────────────┘      └───────────────────────────────┘
```

### Module Breakdown

| Module | Agents | Purpose |
|--------|--------|---------|
| **A - Conversational BI** | 13 | Natural language to SQL, visualization routing, drill-down |
| **B - AI Accountant** | 16 | OCR, ledger mapping, approval workflow, Tally sync |
| **C - Easy Reports** | 8 | Pre-built reports, multi-company consolidation |
| **D - Search Analytics** | 14 | Autonomous analyst, deep research, recommendations |
| **E - i18n** | 7 | Language detection, translation, localized formatting |

---

## 3. Key Technologies

### Backend

| Component | Technology | License |
|-----------|------------|---------|
| Language | Go 1.26 | BSD-3 |
| HTTP Router | go-chi/chi | MIT |
| Message Broker | NATS / JetStream | Apache-2.0 |
| Identity | Keycloak | Apache-2.0 |
| Authorization | Open Policy Agent (OPA) | Apache-2.0 |
| ETL Orchestration | Meltano | MIT |
| Database | PostgreSQL 18.2 + pgvector | PostgreSQL |
| Cache | Redis | BSD-3 |
| Offline Sync | PowerSync | Apache-2.0 |

### AI & ML

| Component | Technology | License |
|-----------|------------|---------|
| AI Flow Orchestration | Google Genkit | Apache-2.0 |
| Multi-agent Framework | Agent ADK | Apache-2.0 |
| Inter-agent Protocol | Google A2A Protocol | Apache-2.0 |
| Semantic Layer | MetricFlow | Apache-2.0 |
| Local LLM | Ollama (Llama 4, Mistral) | MIT |
| GPU LLM | vLLM | Apache-2.0 |
| OCR Engine | PaddleOCR | Apache-2.0 |

### Frontend

| Component | Technology | License |
|-----------|------------|---------|
| Web Framework | React 19.2.4 | MIT |
| Generative UI | CopilotKit | MIT |
| Build Tool | Vite 7.3 | MIT |
| Charts | Apache ECharts | Apache-2.0 |
| Web i18n | i18next + react-i18next | MIT |
| Mobile | Flutter 3.42 | BSD-3 |
| Mobile i18n | flutter_localizations | BSD-3 |

---

## 4. Important Development Patterns

### 4.1 AI Agent Development

All agents follow the **Genkit Flow** pattern:

```go
// Example agent structure
func (s *AgentService) TextToSQLFlow(ctx context.Context, req TextToSQLRequest) (*TextToSQLResponse, error) {
    // 1. Validate input
    // 2. Detect language (E-01)
    // 3. Normalize domain terms (A-04)
    // 4. Generate SQL with confidence scoring (A-01)
    // 5. Validate SQL is SELECT-only
    // 6. Execute via medisync_readonly role
    // 7. Route to visualization (A-03)
    // 8. Format localized response (E-03)
}
```

### 4.2 Security Patterns

**Read-Only Enforcement:**
```go
// All AI agents must use the readonly role
const dbRole = "medisync_readonly"

// SQL validation happens before execution
if !isSelectOnlyQuery(sql) {
    return errors.New("only SELECT queries are allowed")
}
```

**HITL Gates for Write-Backs:**
```go
// Tally sync requires explicit human approval
func (s *SyncService) SyncToTally(ctx context.Context, entry JournalEntry) error {
    if !entry.IsApprovedByFinanceHead {
        return errors.New("finance head approval required")
    }
    // OPA policy check
    if !s.opaClient.Allow(ctx, "tally_sync", user.Roles) {
        return errors.New("unauthorized")
    }
    // Proceed with sync
}
```

### 4.3 i18n Patterns

**Locale Detection Priority:**
1. `user_preferences.locale` (from database, in JWT)
2. `Accept-Language` HTTP header
3. `?lang=ar` URL parameter
4. Default: `en`

**All AI responses include locale instruction:**
```go
prompt := fmt.Sprintf(
    "ResponseLanguageInstruction: Respond in %s. Format numbers according to %s locale.",
    userLocale,
    userLocale,
)
```

### 4.4 Error Handling

```go
// Use wrapped errors for context
if err != nil {
    return fmt.Errorf("failed to execute SQL query: %w", err)
}

// Agent-specific error responses
type AgentError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    Confidence float64 `json:"confidence"`
    Suggestion string `json:"suggestion,omitempty"`
}
```

---

## 5. Code Style Guidelines

### Go Code

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful package names: `internal/agents/module_a`, `internal/warehouse`
- Exported functions must have godoc comments
- Use structured logging (slog)
- Context must be passed to all functions

### React Code

- Functional components with hooks
- TypeScript for type safety
- Tailwind CSS with logical properties for RTL
- i18next for all user-facing strings
- CopilotKit for generative UI components

### File Organization

```
medisync/
├── cmd/                    # Entry points (api, etl, migrate)
├── internal/
│   ├── agents/            # AI agent implementations
│   │   ├── module_a/      # Conversational BI
│   │   ├── module_b/      # AI Accountant
│   │   ├── module_c/      # Easy Reports
│   │   ├── module_d/      # Search Analytics
│   │   └── module_e/      # i18n
│   ├── api/               # HTTP handlers, middleware
│   ├── auth/              # Keycloak JWT, OPA client
│   ├── etl/               # Tally & HIMS connectors
│   ├── warehouse/         # PostgreSQL repository
│   └── cache/             # Redis client
├── policies/              # OPA Rego policies
├── migrations/            # SQL migrations
├── frontend/              # React web app
├── mobile/                # Flutter mobile app
└── docs/                  # Architecture docs
```

---

## 6. Key Agents Reference

| Agent ID | Name | Purpose | Module |
|----------|------|---------|--------|
| A-01 | Text-to-SQL | Natural language to safe SQL | Conversational BI |
| A-02 | SQL Self-Correction | Detect and fix query errors | Conversational BI |
| A-03 | Visualization Routing | Choose optimal chart type | Conversational BI |
| A-04 | Domain Terminology | Map healthcare/accounting terms | Conversational BI |
| A-06 | Confidence Scorer | 0-100% confidence per answer | Conversational BI |
| B-02 | OCR Extraction | Extract fields from documents | AI Accountant |
| B-05 | Ledger Mapping | AI-suggest Tally GL ledger | AI Accountant |
| B-08 | Approval Workflow | Multi-level approval routing | AI Accountant |
| B-09 | Tally Sync | Push approved data to Tally | AI Accountant |
| D-04 | Autonomous AI Analyst | Full analytical workflows | Search Analytics |
| E-01 | Language Detection | Detect query language | i18n |
| E-02 | Query Translation | Arabic → English intent | i18n |
| E-03 | Localized Formatter | Format numbers/dates/currency | i18n |

---

## 7. Common Tasks

### Adding a New Agent

1. Create agent in appropriate `internal/agents/module_X/`
2. Define Genkit flow with input/output structs
3. Add confidence scoring and HITL gates if needed
4. Register in agent supervisor
5. Add OPA policy for authorization
6. Write tests with mock data

### Adding i18n Support

1. Add key to `frontend/public/locales/en/*.json`
2. Add Arabic translation to `frontend/public/locales/ar/*.json`
3. Use in component: `const { t } = useTranslation()`
4. For mobile: add to `mobile/lib/l10n/app_*.arb`
5. Run CI to verify no missing translations

### Database Migration

1. Create migration file in `migrations/`
2. Use `medisync_readonly` role reference for AI queries
3. Test migration on copy of production data
4. Run: `go run ./cmd/migrate`
5. Verify with rollback test

---

## 8. Testing Strategy

### Unit Tests
- Go: standard `testing` package with `testify/assert`
- React: Vitest with React Testing Library
- Mock external dependencies (Tally, HIMS, LLMs)

### Integration Tests
- Test ETL pipelines with sample Tally/HIMS data
- Test agent flows with deterministic LLM responses
- Test OPA policies with various role combinations

### End-to-End Tests
- Full query flow: natural language → SQL → chart
- Document flow: upload → OCR → approval → Tally sync
- Test both English and Arabic user flows

---

## 9. Security Checklist

Before committing code, verify:

- [ ] AI agents use `medisync_readonly` DB role
- [ ] All SQL queries are validated as SELECT-only
- [ ] Write-back operations require HITL approval
- [ ] OPA policies cover new endpoints
- [ ] PII is masked based on user role
- [ ] All external API calls use TLS 1.3+
- [ ] Sensitive config is via environment variables
- [ ] Audit log entries are created for actions

---

## 10. Troubleshooting

### Common Issues

**LLM Hallucination:**
- Check confidence scores from A-06
- Review domain terminology mappings in A-04
- Verify schema context in pgvector

**Low OCR Accuracy:**
- Check document quality preprocessing
- Review confidence scores for field-level issues
- Route to handwriting agent (B-03) if needed

**Tally Sync Failures:**
- Verify OPA policy allows user action
- Check approval workflow completion
- Review Tally XML gateway logs

**i18n Issues:**
- Verify locale is in Accept-Language or JWT
- Check translation files have matching keys
- For Arabic: verify RTL layout is applied

---

## 11. Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Query accuracy | ≥ 95% | Business intent validation |
| Query latency (P95) | < 5 seconds | End-to-end query flow |
| Dashboard load | < 3 seconds | Time to first chart |
| OCR accuracy (standard) | ≥ 95% | Field-level comparison |
| OCR accuracy (handwritten) | ≥ 90% | Field-level comparison |
| Tally sync success | ≥ 99.5% | Successful syncs / total |
| System uptime | ≥ 99.5% | Infrastructure monitoring |

---

## 12. Documentation Links

| Document | Path | Purpose |
|----------|------|---------|
| Architecture | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Full system architecture |
| PRD | [docs/PRD.md](docs/PRD.md) | Product requirements |
| Design System | [docs/DESIGN.md](docs/DESIGN.md) | Colors, typography, components |
| i18n Architecture | [docs/i18n-architecture.md](docs/i18n-architecture.md) | Internationalisation design |
| Agent Blueprints | [docs/agents/](docs/agents/) | Per-agent specifications |
| Security | [docs/agents/03-governance-security.md](docs/agents/03-governance-security.md) | OPA policies, RBAC |

---

## 13. Development Workflow

```bash
# Start infrastructure
docker-compose up -d

# Run migrations
go run ./cmd/migrate

# Start API server
go run ./cmd/api

# Start ETL service
go run ./cmd/etl

# Start web frontend
cd frontend && npm run dev

# Start mobile
cd mobile && flutter run

# Run tests
go test ./...
cd frontend && npm test
cd mobile && flutter test
```

---

## 14. Philosophy Summary

> **"MediSync exists to liberate healthcare and accounting teams from the tyranny of spreadsheets and manual reconciliation."**

When working on this codebase:
- **Security first**: Never compromise on HITL gates for write-backs
- **i18n by default**: Every feature must work in both English and Arabic
- **Open source**: Always choose OSI-approved licenses
- **Incremental delivery**: Ship value early, iterate often
- **User trust**: Confidence scores, audit trails, and transparency

---

*This guide is maintained by the MediSync engineering team. For questions or updates, please refer to the project documentation or contact the team lead.*

## Active Technologies
- Go 1.26 (backend), TypeScript/React 19 (frontend types) (001-ai-agent-core)
- PostgreSQL 18.2 + pgvector (schema embeddings), Redis (session cache) (001-ai-agent-core)
- CopilotKit 1.3.6 (streaming chat UI), Apache ECharts 5.6 (charts), i18next 24.2 (i18n) (002-dashboard-advanced-features)
- NATS JetStream (alert/report scheduling), Puppeteer (PDF generation), excelize (Excel export) (002-dashboard-advanced-features)

## Recent Changes
- 001-ai-agent-core: Added Go 1.26 (backend), TypeScript/React 19 (frontend types)
- 002-dashboard-advanced-features: Added CopilotKit, ECharts, i18next, NATS JetStream scheduling
