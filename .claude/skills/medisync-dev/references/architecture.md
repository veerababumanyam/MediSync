# MediSync Architecture Reference

Detailed architecture patterns for the MediSync platform.

## Layer Stack

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

## Three-Plane Architecture

### 1. Data Plane (ETL Layer)

**Purpose**: Extract, transform, and load data from external systems into the warehouse without affecting operations.

**Components**:
- **Tally Connector**: Polls Tally HTTP gateway for masters/vouchers
- **HIMS Connector**: Fetches clinical/operational data via REST API
- **Data Validator**: Schema validation, duplicate detection
- **NATS Publisher**: Emits events for real-time updates

**Key Pattern**:
```go
// ETL jobs run independently of API server
// Warehouse is the single source of truth for AI queries
type ETLJob struct {
    Source    string // "tally" or "hims"
    Frequency time.Duration
    Validator *DataValidator
}
```

### 2. Intelligence Plane (AI Layer)

**Purpose**: Read-only query and analysis using AI agents.

**Components**:
- **Genkit Flows**: Orchestrate AI agent logic
- **Agent ADK**: Multi-agent coordination framework
- **A2A Protocol**: Inter-agent communication
- **OPA Policy Engine**: Authorization checks

**Critical Security**: All agents use `medisync_readonly` database role.

### 3. Action Plane (Write-Back Layer)

**Purpose**: Human-gated writes back to external systems.

**Components**:
- **Approval Workflow**: Multi-level sign-off (Module B-08)
- **HITL Gateway**: Human-in-the-loop verification
- **Tally Sync**: Push approved transactions (Module B-09)

**Key Pattern**:
```go
// No AI can write without human approval
if !entry.IsApprovedByFinanceHead {
    return errors.New("HITL approval required")
}
```

## Module Breakdown

### Module A: Conversational BI (13 agents)

**Purpose**: Natural language to SQL, visualization, drill-down.

**Key Agents**:
- A-01: Text-to-SQL - Converts queries to safe SELECT statements
- A-02: SQL Self-Correction - Detects and fixes query errors
- A-03: Visualization Routing - Selects optimal chart type
- A-04: Domain Terminology - Maps healthcare/accounting terms
- A-06: Confidence Scorer - Scores answer reliability 0-100%

**Data Flow**:
```
User Query → E-01 (Language Detect) → A-04 (Normalize Terms)
    → A-01 (Generate SQL) → A-02 (Validate SQL)
    → Execute (readonly) → A-03 (Route to Viz)
    → E-03 (Format Response)
```

### Module B: AI Accountant (16 agents)

**Purpose**: OCR, ledger mapping, approvals, Tally sync.

**Key Agents**:
- B-01: Document Classifier - Invoice/Bill/Statement
- B-02: OCR Extraction - Field-level text extraction
- B-03: Handwriting Recognition - Script OCR
- B-05: Ledger Mapping - Suggests Tally GL accounts
- B-08: Approval Workflow - Multi-step sign-off
- B-09: Tally Sync - Pushes to Tally ERP
- B-10: Bank Reconciliation - Matches entries

**Data Flow**:
```
Document Upload → B-01 (Classify) → B-02 (OCR)
    → B-05 (Map Ledger) → B-08 (Approval Workflow)
    → B-09 (Tally Sync) [HITL GATED]
```

### Module C: Easy Reports (8 agents)

**Purpose**: Pre-built MIS reports, scheduled reports.

**Key Features**:
- Pre-built templates (P&L, Balance Sheet, Sales Analysis)
- Multi-company consolidation
- Scheduled delivery (email/app)
- Zero-code custom dashboard builder

### Module D: Search Analytics (14 agents)

**Purpose**: Autonomous AI analyst, deep research, recommendations.

**Key Agents**:
- D-04: Autonomous AI Analyst - Multi-step analysis
- D-05: Deep Research - Pattern discovery
- D-08: Prescriptive AI - Quantified recommendations

**Pattern**: Chain-of-thought reasoning with tool use.

### Module E: i18n (7 agents)

**Purpose**: Language detection, translation, localization.

**Key Agents**:
- E-01: Language Detection - en/ar classification
- E-02: Query Translation - Arabic → English intent
- E-03: Localized Formatter - Numbers, dates, currency

**Pattern**: Every AI response includes locale instructions.

## Technology Stack

### Backend

| Component | Technology | License |
|-----------|------------|---------|
| Language | Go 1.26 | BSD-3 |
| Router | go-chi/chi | MIT |
| Message Broker | NATS / JetStream | Apache-2.0 |
| Identity | Keycloak | Apache-2.0 |
| Authorization | Open Policy Agent (OPA) | Apache-2.0 |
| Database | PostgreSQL 18.2 + pgvector | PostgreSQL |
| Cache | Redis | BSD-3 |
| Offline Sync | PowerSync | Apache-2.0 |

### AI & ML

| Component | Technology | License |
|-----------|------------|---------|
| Orchestration | Google Genkit | Apache-2.0 |
| Multi-agent | Agent ADK | Apache-2.0 |
| Protocol | Google A2A | Apache-2.0 |
| Semantic Layer | MetricFlow | Apache-2.0 |
| Local LLM | Ollama (Llama 4, Mistral) | MIT |
| GPU LLM | vLLM | Apache-2.0 |
| OCR | PaddleOCR | Apache-2.0 |

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

## Data Flow Patterns

### Query Flow (Conversational BI)
```
User (EN/AR) → API Gateway → E-01 (Language Detect)
    → A-01 (Text-to-SQL) → SQL Validator
    → Warehouse (readonly) → A-03 (Visualization)
    → E-03 (Format) → Response
```

### Document Flow (AI Accountant)
```
Upload → B-01 (Classify) → B-02 (OCR)
    → B-05 (Ledger Mapping) → Draft Entry
    → B-08 (Approval Workflow) → Human Approval
    → B-09 (Tally Sync) → Confirmation
```

### Report Flow (Easy Reports)
```
Select Report → C-01 (Template Resolver)
    → Warehouse Query → Format → Export (PDF/Excel)
```

## Security Architecture

### Authentication Flow
```
Client → Keycloak → JWT → API Gateway → Extract Claims
    → OPA Policy Check → Allow/Deny
```

### Authorization Layers
1. **Keycloak**: Authentication and role assignment
2. **OPA**: Fine-grained policy enforcement
3. **Database**: Row-level security via `medisync_readonly` role
4. **HITL**: Human approval for financial writes

### Audit Trail
All financial operations logged to `audit_log` table:
- User ID
- Action type
- Timestamp
- Before/after state
- Approval chain

## Scalability Patterns

### Horizontal Scaling
- API Gateway: Stateless, can run multiple instances
- AI Agents: Genkit flows can be distributed
- ETL Jobs: Can be scheduled across workers

### Caching Strategy
- Redis for: Session data, cached query results, rate limiting
- CDN for: Static assets, i18n translation files

### Database Optimization
- Read replicas for AI queries
- Connection pooling via pgx
- Materialized views for common aggregations
