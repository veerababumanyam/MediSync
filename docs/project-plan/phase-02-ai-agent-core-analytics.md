# Phase 02 — AI Agent & Core Analytics

**Phase Duration:** Weeks 4–7 (4 weeks)  
**Module(s):** Module A (Conversational BI), Module E (Language & Localisation)  
**Status:** Planning  
**Milestone:** M2 — First AI Query  
**Depends On:** Phase 01 complete (data warehouse seeded)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.1](../ARCHITECTURE.md) | [agents/specs/a-01-text-to-sql.md](../agents/specs/a-01-text-to-sql.md)

---

## 1. Objectives

Build and deploy the core AI intelligence layer: the Text-to-SQL agent and its supporting pipeline. By end of Phase 02, any authorised user can submit a natural-language business question and receive an accurate, secured SQL query result — even in Arabic. This is the foundational capability that all subsequent modules build upon.

---

## 2. Scope

### In Scope
- Go API gateway (core middleware stack)
- Genkit framework setup with LLM plugin configuration
- A2A inter-agent protocol scaffolding
- 6 Module A agents (A-01 through A-06)
- 3 Module E localisation agents (E-01, E-02, E-03)
- pgvector schema embeddings generation
- MetricFlow semantic layer — initial metric definitions
- OPA policy engine setup with `bi.read_only` policy
- `POST /v1/chat` API endpoint (SSE streaming)
- WebSocket `/ws/chat` endpoint
- Basic authentication (JWT validation via Keycloak)

### Out of Scope
- Frontend UI (Phase 3)
- Dashboard pinning (Phase 3)
- Scheduled reports (Phase 3)
- Document processing (Phase 4+)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | Go API Gateway core (go-chi) | Backend Engineer | JWT validation, OPA sidecar call, Accept-Language header → locale ctx, rate limiting all working |
| D-02 | Genkit Framework Setup | AI Engineer | Genkit project initialised; GPT-5.2 and local Ollama plugins configured and swappable |
| D-03 | A2A Protocol Scaffolding | AI Engineer | `/.well-known/agent.json` AgentCard endpoints for each agent; task delegation `POST /tasks/send` working A-01 → A-02 |
| D-04 | A-01 Text-to-SQL Agent | AI Engineer | Converts natural-language questions to correct SQL; `medisync_readonly` role enforced; 95%+ accuracy on test set of 50 queries |
| D-05 | A-02 SQL Self-Correction Agent | AI Engineer | Detects SQL execution error, rewrites query, retries up to 3× successfully on known error classes |
| D-06 | A-03 Visualisation Routing Agent | AI Engineer | Classifies trend→line, comparison→bar, breakdown→pie, tabular→table with >98% accuracy on test set |
| D-07 | A-04 Domain Terminology Normaliser | AI Engineer | Correctly maps ≥ 30 healthcare + accounting synonyms; "footfall"→patient_visits; tested against synonym test suite |
| D-08 | A-05 Hallucination Guard | AI Engineer | Deflects off-topic questions (≥ 20 test cases); passes valid business queries through without false positives |
| D-09 | A-06 Confidence Scoring Agent | AI Engineer | Attaches 0–100% confidence to every response; low-confidence (<70%) routes to `app.review_queue` |
| D-10 | E-01 Language Detection & Routing | AI Engineer | Detects `en`/`ar` from query text with >99% accuracy; injects locale into Genkit flow context |
| D-11 | E-02 Query Translation Agent | AI Engineer | Arabic natural-language query normalised to English intent before SQL generation; tested on 20 Arabic queries |
| D-12 | E-03 Localised Response Formatter | AI Engineer | Numbers, dates, currency in user locale; EN: 1,234.56 / AR: ١٢٣٤٫٥٦; date formats correct |
| D-13 | pgvector Schema Embeddings | Data Engineer | All warehouse table/column descriptions embedded into `vectors.schema_embeddings`; A-01 retrieves correct context |
| D-14 | MetricFlow Initial Metrics | Data Engineer | ≥ 15 core metrics defined (clinic_revenue, pharmacy_margin_pct, patient_visits, outstanding_receivables, etc.) |
| D-15 | OPA `bi.read_only` Policy | DevOps | Policy blocks any DML (INSERT/UPDATE/DELETE/DROP) in agent-generated SQL; column masking for PII works |
| D-16 | `POST /v1/chat` API (SSE) | Backend Engineer | Accepts query, streams response via SSE; latency P95 < 5 seconds for standard queries |

---

## 4. AI Agents Deployed

### A-01 Text-to-SQL Agent

**Pipeline:**
```
User Query (+ locale context from E-01)
    │
    ▼ E-02: Arabic query normalised to English intent (if AR)
    │
    ▼ A-04: Domain terminology normalisation
    │
    ▼ A-06: Confidence pre-check (ambiguity detection)
    │       If ambiguous → return clarifying question
    │
    ▼ Schema context retrieval (pgvector similarity search)
    │       Top-5 relevant tables + columns injected into prompt
    │       MetricFlow metric definitions for matching metrics
    │
    ▼ LLM: Genkit flow generates SQL
    │       System prompt: role, read-only mandate, locale rule
    │       SQL injection guard: all user strings → bind parameters
    │
    ▼ OPA check: block DML (medisync.bi.read_only policy)
    │
    ▼ PostgreSQL executor (medisync_readonly role)
    │       On error → A-02 self-correction loop (max 3 retries)
    │
    ▼ Result DataFrame → A-03 visualisation routing
    │
    ▼ E-03: Format numbers/dates/currencies in user locale
    │
    ▼ SSE stream response to client
```

**Security guardrails:**
- SQL injection: parameterised bind; user strings never interpolated into SQL
- Read-only enforcement: OPA `medisync.bi.read_only` rejects DML before execution
- Hallucination guard: A-05 runs before A-01; off-topic queries rejected
- Column masking: OPA strips `cost_price`, `patient_pii_columns` for non-authorised roles
- PII access: columns tagged as PII in schema_embeddings require `analyst` role minimum; `viewer` blocked
- HITL gate: confidence < 70% → question added to `app.review_queue` for human response

### A-02 SQL Self-Correction Agent

**Trigger:** A-01 receives a PostgreSQL error on query execution  
**Strategy:**
1. Receives error message + original query + schema context
2. LLM analyzes error (unknown column? wrong table? syntax error?)
3. Generates corrected query
4. Retries execution (up to 3 attempts)
5. If all retries fail → returns "I couldn't answer that precisely" with partial result if available

**Common correction patterns:**
- Column name mismatch: checks schema_embeddings for correct column name
- Wrong aggregation: rewrites GROUP BY / HAVING
- Date function syntax: fixes PostgreSQL-specific functions

### A-03 Visualisation Routing Agent

**Classification rules:**
| Query Pattern | Chart Type | ECharts Config Key |
|---|---|---|
| Time series, trend, over time | Line chart | `lineChart` |
| Compare N items, ranking, by X | Bar chart | `barChart` |
| Proportion, share, breakdown | Pie chart | `pieChart` |
| Geographic | Map (future) | `mapChart` |
| Single metric | KPI card | `kpiCard` |
| Multi-column data | Table | `dataTable` |
| Correlation | Scatter plot | `scatterChart` |

### A-04 Domain Terminology Normaliser

**Healthcare synonyms (examples):**
| User Term | Canonical Term | SQL Column |
|---|---|---|
| footfall / walk-ins / visits | patient_visits | `fact_appointments.count` |
| outstanding / pending payments | accounts_receivable | `fact_billing WHERE payment_status='pending'` |
| OPD / outpatient | outpatient_department | `dim_doctors.department='OPD'` |
| expiry / expires | stock_expiry_date | `dim_drugs.expiry_date` |
| dispensation / dispensed | pharmacy_dispensation | `fact_pharmacy_disp` |

**Glossary source:** `docs/i18n-glossary.md` — bilingual glossary loaded as context

### A-05 Hallucination Guard

**Off-topic detection categories:**
- General knowledge questions ("What is the capital of France?")
- Personal advice ("What should I eat today?")
- Code requests unrelated to MediSync data
- Political / controversial topics

**Response:** "I'm your MediSync data analyst. I can help answer questions about your clinic, pharmacy, and financial data. What would you like to know?"

### A-06 Confidence Scoring Agent

**Scoring factors:**
- Query intent ambiguity (multiple possible interpretations) → lower score
- Schema match quality (did pgvector find highly relevant tables?) → higher if good match
- SQL complexity (simple SELECT vs complex multi-join) → simple = higher confidence
- A-02 retries required → each retry decreases score by 10%

**Routing logic:**
- Score ≥ 70%: proceed normally, show score badge in UI
- Score 50–69%: show response with "Low confidence" warning; add to review queue
- Score < 50%: ask clarifying question instead of executing SQL

### E-01 Language Detection & Routing

**Detection method:** `lingua-go` library for language identification + Unicode script detection (Arabic Unicode range U+0600–U+06FF)  
**Locale injection:** Sets `locale`, `lang`, `rtl` fields in Genkit flow context  
**Priority order:** User `preferences.locale` → query language detection → `Accept-Language` header → default `en`

### E-02 Query Translation Agent

**Approach:** 
1. Detect Arabic query
2. Pass to LLM with system prompt: "Extract the analytic intent from this Arabic question. Return the intent in English for SQL generation. Do not translate word-for-word; extract the business meaning."
3. Original Arabic query preserved for response formatting via E-03

### E-03 Localised Response Formatter

**Formatting rules by locale:**

| Field Type | English (`en`) | Arabic (`ar`) |
|---|---|---|
| Number | 1,234,567.89 | ١٬٢٣٤٬٥٦٧٫٨٩ |
| Currency (INR) | ₹1,234.56 | ١٬٢٣٤٫٥٦ ₹ |
| Date | 19 Feb 2026 | ١٩ فبراير ٢٠٢٦ |
| Percentage | 12.5% | ١٢٫٥٪ |

**AI response language injection** (system prompt snippet):
```
LANGUAGE RULE: You MUST respond entirely in {{locale}}.
Format all numbers as {{number_format}}.
Format all dates as {{date_format}}.
Format all currency as {{currency_format}}.
```

---

## 5. API Endpoints

| Method | Path | Description | Auth |
|---|---|---|---|
| `POST` | `/v1/chat` | Submit query; stream response via SSE | JWT (any role) |
| `WS` | `/ws/chat` | WebSocket real-time chat | JWT (any role) |
| `GET` | `/v1/agents/health` | Agent ecosystem health check | JWT (admin) |

### Chat Request Schema
```json
{
  "query": "Show me total clinic revenue for January 2026",
  "locale": "en",
  "session_id": "uuid",
  "context": []
}
```

### Chat Response (SSE)
```json
// Stream event 1: thinking
{"type": "thinking", "message": "Analyzing your query..."}

// Stream event 2: SQL generated
{"type": "sql_preview", "sql": "SELECT SUM(amount) FROM fact_billing WHERE..."}

// Stream event 3: result
{"type": "result", "chart_type": "kpiCard", "data": {...}, "confidence": 92}
```

---

## 6. MetricFlow Metrics (Initial Set)

```yaml
metrics:
  - name: clinic_revenue
    label: Clinic Revenue
    description: Total billed amount from clinic appointments
    
  - name: pharmacy_revenue  
    label: Pharmacy Revenue
    description: Total revenue from pharmacy dispensations
    
  - name: total_revenue
    label: Total Revenue
    description: Clinic Revenue + Pharmacy Revenue
    
  - name: patient_visits
    label: Patient Visits
    description: Count of completed clinic appointments
    
  - name: outstanding_receivables
    label: Outstanding Receivables
    description: Sum of unpaid bills (accounts receivable)
    
  - name: pharmacy_margin_pct
    label: Pharmacy Margin %
    description: (Pharmacy Revenue - Cost) / Pharmacy Revenue × 100
    
  - name: avg_revenue_per_patient
    label: Avg Revenue per Patient
    description: Total Revenue / Unique Patient Count
    
  - name: new_patients_count
    label: New Patients
    description: Count of first-time patients in period
    
  - name: appointment_cancellation_rate
    label: Cancellation Rate %
    description: Cancelled appointments / Total scheduled × 100
```

---

## 7. OPA Policy: `medisync.bi.read_only`

```rego
package medisync.bi

import future.keywords

# Block any DML in agent-generated SQL
deny[msg] if {
    input.sql_statement
    sql_upper := upper(input.sql_statement)
    patterns := ["INSERT", "UPDATE", "DELETE", "DROP", "TRUNCATE", "ALTER", "CREATE"]
    p := patterns[_]
    contains(sql_upper, p)
    msg := sprintf("DML operation '%v' not permitted in BI queries", [p])
}

# Mask PII columns for viewer role
column_mask[col] if {
    input.user.role == "viewer"
    col := {"dim_patients.name_en", "dim_patients.name_ar", 
            "dim_patients.phone", "dim_patients.dob"}[_]
}

# Restrict cost/margin columns to manager+ roles
column_mask[col] if {
    not input.user.role in ["admin", "finance_head", "accountant_lead", "manager"]
    col := {"dim_inventory_items.cost_price", "fact_vouchers.cost_amount"}[_]
}
```

---

## 8. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| A-01 SQL accuracy | 50 representative business queries from PRD user stories | ≥ 95% generate correct SQL |
| A-02 Self-correction | 15 deliberately broken SQL statements | ≥ 80% successfully corrected |
| A-03 Chart routing | 40 query intent samples | ≥ 98% correct chart type |
| A-04 Terminology | Full domain glossary (30+ terms) | 100% mapped correctly |
| A-05 Hallucination | 20 off-topic queries + 20 valid queries | 0 false positives, 0 false negatives |
| A-06 Confidence | 30 queries of varying ambiguity | Scores correlate with actual accuracy |
| E-01 Language detection | 50 Arabic + 50 English queries | ≥ 99% accuracy |
| E-02 Arabic translation | 20 Arabic business queries | Semantic equivalence validated by bilingual reviewer |
| E-03 Number formatting | 20 numbers/dates/currencies | 100% correct locale formatting |
| Performance | P95 query latency | < 5 seconds end-to-end |

---

## 9. Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Phase 01 complete | Required | Warehouse must have data for SQL queries to return results |
| LLM API keys | Required | GPT-5.2 or Claude 4.6 API keys; or Ollama running locally |
| Arabic bilingual reviewer | Required for E-02 testing | Native Arabic speaker to validate intent translation quality |

---

## 10. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| LLM SQL accuracy below 95% on complex joins | High | Iterative prompt engineering; few-shot examples in pgvector query_history; add more MetricFlow definitions |
| Arabic intent translation losing business context | High | Test with native Arabic speaker; fallback to word-by-word + glossary if semantic translation ambiguous |
| pgvector schema embeddings slow retrieval | Medium | Tune HNSW index parameters; limit k=5 top results; benchmark before moving to Phase 3 |
| OPA policy causing false DML rejections | Low | Comprehensive SQL positive/negative test suite before go-live |

---

## 11. Phase Exit Criteria

- [ ] `POST /v1/chat` API end-to-end working with ≥ 95% SQL accuracy on test suite
- [ ] All 6 Module A agents (A-01 to A-06) deployed and tested
- [ ] All 3 Module E agents (E-01, E-02, E-03) deployed; Arabic queries returning Arabic responses correctly
- [ ] OPA `bi.read_only` policy blocking DML; PII column masking verified
- [ ] pgvector schema embeddings seeded for all warehouse tables
- [ ] MetricFlow ≥ 15 core metrics defined and tested
- [ ] Genkit flow tracing visible in observability stack
- [ ] P95 query latency < 5 seconds measured on benchmark set
- [ ] Phase gate reviewed and signed off

---

*Phase 02 | Version 1.0 | February 19, 2026*
