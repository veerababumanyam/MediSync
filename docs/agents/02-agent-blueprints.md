# MediSync — Agent Blueprints

**Version:** 1.0 | **Created:** February 19, 2026  
**Cross-ref:** [00-agent-backlog.md](./00-agent-backlog.md) | [01-oss-toolchain.md](./01-oss-toolchain.md) | [03-governance-security.md](./03-governance-security.md)

> This document provides detailed development blueprints for the highest-priority and most complex agents in the MediSync platform. Each blueprint specifies inputs, outputs, tool chain, prompt design, guardrails, HITL gates, and evaluation criteria.

---

## Blueprint Index

| # | Agent | Module | Phase | Complexity | Read/Write |
|---|-------|--------|-------|-----------|-----------|
| [A-01](#a-01-text-to-sql-agent) | Text-to-SQL Agent | Conversational BI | 2 | L2 | Read-only |
| [B-02](#b-02-ocr-extraction-agent) | OCR Extraction Agent | AI Accountant | 4 | L2 | Read-only |
| [B-05](#b-05-ledger-mapping-agent) | Ledger Mapping Agent | AI Accountant | 5 | L2 | Read-only (suggests) |
| [B-08](#b-08-approval-workflow-agent) | Approval Workflow Agent | AI Accountant | 5 | L2 | **Write (workflow state)** |
| [B-09](#b-09-tally-sync-agent) | Tally Sync Agent | AI Accountant | 6 | L2 | **Write (Tally)** |
| [B-10](#b-10-bank-reconciliation-agent) | Bank Reconciliation Agent | AI Accountant | 7 | L2 | Read-only |
| [C-06](#c-06-data-quality-validation-agent) | Data Quality Validation Agent | Easy Reports | 1 | L2 | Read-only |
| [D-04](#d-04-autonomous-ai-analyst-agent) | Autonomous AI Analyst Agent | Search Analytics | 14 | L3 | Read-only |
| [D-08](#d-08-prescriptive-recommendations-agent) | Prescriptive Recommendations Agent | Search Analytics | 15 | L3 | Read-only* |

---

## A-01 — Text-to-SQL Agent

### Purpose
Convert a natural language business question into a safe, read-only SQL query against the MediSync data warehouse, execute it, and return structured results with metadata.

### Inputs
| Field | Type | Source |
|-------|------|--------|
| `user_query` | string | Chat interface |
| `user_role` | string | JWT (Keycloak) |
| `schema_context` | JSON | Schema Context Cache (pre-loaded) |
| `semantic_context` | JSON | MetricFlow metric registry |
| `conversation_history` | list[Message] | Session memory (last N turns) |

### Outputs
| Field | Type | Description |
|-------|------|-------------|
| `sql_query` | string | Generated and validated SQL |
| `result_set` | DataFrame | Query execution result |
| `chart_type` | enum | `bar \| line \| pie \| table \| scatter` |
| `confidence_score` | float (0–1) | Model's self-reported confidence |
| `explanation` | string | Plain-English explanation of what was queried |
| `error` | string \| null | If query failed after retries |

### Tool Chain (LangChain)
```
UserQuery
  → [A-04] Domain Terminology Normaliser
  → [A-06] Confidence Pre-Checker (intent ambiguity detection)
  → LangChain SQLDatabaseChain (+ schema + semantic context)
  → SQL Validator (read-only assertion: no INSERT/UPDATE/DELETE/DROP)
  → PostgreSQL Executor
  → [A-02] Self-Correction loop (on DB error, max 3 retries)
  → [A-03] Visualization Router
  → Pydantic Response Model
  → Chat Response Formatter
```

### OSS Components
- **LangChain** `langchain-community` SQLDatabaseChain / SQL Agent
- **LlamaIndex** for schema indexing and retrieval
- **MetricFlow** for resolving metric definitions to SQL fragments
- **PostgreSQL** (read-only service account — see §03-governance)
- **Pydantic** for structured output validation

### System Prompt (Template)
```
You are an expert SQL analyst for MediSync, a healthcare and accounting platform.
You have READ-ONLY access to a PostgreSQL data warehouse containing:
  - HIMS data: patients, appointments, billing, pharmacy dispensations
  - Tally data: ledgers, vouchers, inventory, sales, receipts

RULES:
1. Generate ONLY SELECT statements. NEVER use INSERT, UPDATE, DELETE, DROP, CREATE, TRUNCATE.
2. Always apply the user's role filter: {{ user_role_filter }}
3. Use metric definitions from the Semantic Layer when available.
4. If the question is ambiguous, ask one clarifying question before generating SQL.
5. If the question is not data-related, respond: "I can only answer questions about your business data."

Available schema:
{{ schema_context }}

Available metrics:
{{ semantic_context }}
```

### Guardrails
- **SQL Injection Guard:** Parameterise all inputs; reject queries with non-SELECT DML.
- **Read-only enforcement:** OPA policy `medisync.bi.read_only` blocks any write operations at the DB connection level (separate read-only Postgres role).
- **Hallucination Guard:** If confidence < 0.70, append "Low confidence — please verify" and queue for manual review.
- **Off-topic deflection:** Classifier chain detects non-business queries and refuses gracefully.
- **Column masking:** OPA policy strips sensitive columns (patient PII, cost prices) based on `user_role`.

### HITL Gate
- **Trigger:** `confidence_score < 0.70` OR SQL touches sensitive tables (patient PII).
- **Action:** Route to manual review queue; notify user "Your query is being reviewed."

### Evaluation Criteria
| Metric | Target |
|--------|--------|
| SQL correctness (executes without error) | ≥ 98% |
| Business intent accuracy (human eval) | ≥ 95% |
| Latency (P95) | < 5 seconds |
| Hallucination rate (non-data answers returned) | < 1% |
| False positives on off-topic filter | < 2% |

---

## B-02 — OCR Extraction Agent

### Purpose
Extract structured financial fields from uploaded documents (PDFs, images, Excel, handwritten scans) with confidence scoring; flag low-confidence items for human review.

### Inputs
| Field | Type | Source |
|-------|------|--------|
| `file` | bytes | Drag-and-drop upload; bulk batch |
| `file_type` | enum | `pdf \| png \| jpg \| xlsx \| csv` |
| `extraction_schema` | JSON | Pre-configured field list |

### Outputs
| Field | Type | Description |
|-------|------|-------------|
| `document_type` | enum | `invoice \| bill \| bank_statement \| receipt \| tax_doc \| other` |
| `extracted_fields` | dict | `{vendor, amount, invoice_date, invoice_no, tax_amount, currency}` |
| `confidence_scores` | dict | Per-field confidence (0–1) |
| `overall_confidence` | float | Aggregated confidence score |
| `low_confidence_flags` | list[str] | Fields requiring manual review |
| `raw_text` | string | Full OCR text output |

### Tool Chain
```
File Upload
  → [B-01] Document Classifier (PaddleOCR layout analysis)
  → Router:
      Digital PDF → PyMuPDF text extraction
      Scanned/Image → PaddleOCR pipeline
      Handwritten → PaddleOCR (HTR model) + LLM post-processing
  → Unstructured.io segmenter (tables, headers, line items)
  → LangChain extraction chain (JSON field extraction)
  → Pydantic ExtractionResult model (field validation)
  → Confidence Scorer
  → HITL Router (if overall_confidence < 0.85)
```

### OSS Components
- **PaddleOCR** (Apache-2.0) — primary OCR + layout detection
- **Tesseract** (Apache-2.0) — fallback for simple documents
- **PyMuPDF** (AGPL-3.0) — digital PDF text layer extraction
- **Unstructured.io** (Apache-2.0) — document segmentation
- **LangChain** structured extraction chain
- **Pydantic** output validation

### Guardrails
- Maximum file size: 50 MB. Reject with user error above limit.
- Allowed MIME types validated server-side (not just extension).
- Extracted data never written to Tally without passing through B-05 → B-08 → B-09 pipeline.
- All uploaded files stored with AES-256 encryption at rest.

### HITL Gate
- **Trigger:** `overall_confidence < 0.85` OR any required field (amount, vendor, date) has `confidence < 0.70`.
- **Action:** Highlight flagged fields in the upload review UI; block sync until human confirms.

### Evaluation Criteria
| Metric | Target |
|--------|--------|
| Field extraction accuracy (printed documents) | ≥ 95% |
| Field extraction accuracy (handwritten) | ≥ 90% |
| False positive rate (wrong doc type) | < 3% |
| Processing time per page (P95) | < 10 seconds |

---

## B-05 — Ledger Mapping Agent

### Purpose
Suggest the most appropriate Tally GL ledger for each extracted transaction, with a confidence score, based on transaction context and historical mapping patterns.

### Inputs
| Field | Type | Source |
|-------|------|--------|
| `transaction` | dict | B-02 extraction output |
| `tally_chart_of_accounts` | list[Ledger] | Tally sync cache |
| `historical_mappings` | vector embeddings | Chroma (learned from past corrections) |

### Outputs
| Field | Type | Description |
|-------|------|-------------|
| `suggested_ledger` | string | Tally ledger name |
| `suggested_sub_ledger` | string \| null | Sub-ledger |
| `suggested_cost_centre` | string \| null | Cost centre |
| `confidence_score` | float (0–1) | Mapping confidence |
| `confidence_badge` | enum | `high (≥0.95) \| review (0.70–0.94) \| manual (<0.70)` |
| `alternative_mappings` | list[dict] | Top 3 alternatives with scores |
| `reasoning` | string | LLM explanation for suggested mapping |

### Tool Chain
```
Transaction (from B-02)
  → Embedding Generator (BAAI/bge-small, Apache-2.0)
  → Chroma similarity search (top-5 similar past mappings)
  → Context builder: [transaction + historical matches + chart of accounts]
  → LangChain classification chain
  → Pydantic MappingResult
  → Confidence Badge Assigner
  → Update Chroma on user approval (feedback loop)
```

### Learning Loop
When a user approves or overrides a mapping:
1. Store (transaction_embedding, confirmed_ledger) in Chroma.
2. On next run, similar transactions will be matched via vector similarity before LLM call (reduces latency + cost).

### Guardrails
- Agent only **suggests** — never writes to Tally without B-08 (approval) → B-09 (sync) chain.
- "Bulk Update Rules" option (apply correction to all similar future transactions) requires Finance Head role.

### HITL Gate
- **Always** required — no transaction is mapped without user review.
- `confidence_badge = manual` → red flag, mandatory Accountant review.
- `confidence_badge = review` → amber flag, review recommended.

---

## B-08 — Approval Workflow Agent

### Purpose
Route transactions through a configurable multi-step approval chain before any data is written to Tally. Enforce separation of duties.

### Approval Chain (default)
```
Draft → [Accountant Review] → [Finance Manager Approval] → [Finance Head Sign-off] → Ready to Sync
```

### Inputs
| Field | Type | Source |
|-------|------|--------|
| `transactions` | list[Transaction] | B-02 + B-05 pipeline output |
| `approval_policy` | JSON | OPA policy bundle |
| `requesting_user` | User | JWT |

### Outputs
| Field | Type | Description |
|-------|------|-------------|
| `workflow_id` | UUID | Unique workflow instance |
| `current_status` | enum | `draft \| pending_accountant \| pending_manager \| pending_finance \| approved \| rejected` |
| `approval_history` | list[ApprovalEvent] | Timestamped actions per approver |
| `rejection_reason` | string \| null | If rejected |
| `approved_transactions` | list[Transaction] | Passed to B-09 on approval |

### Tool Chain
```
Approved Ledger Mapping (from B-05)
  → OPA policy check: has_permission(user, 'approve_transaction')
  → Workflow State Machine (FastAPI + Postgres state table)
  → Notification Dispatcher (Apprise → Email/In-App)
  → Reminder Scheduler (Celery Beat — 24h reminder for stale approvals)
  → On final approval: emit ApprovedTransactions event → B-09
```

### Guardrails
- No user can approve their own submissions (self-approval blocked by OPA policy).
- Bulk approvals > ₹1,00,000 require Finance Head role (configurable threshold).
- All approval events written to immutable audit log (B-14).
- Rejection sends notification + reason to submitting user.

### HITL Gate — **This agent IS the HITL gate.** All decisions are human-made.

---

## B-09 — Tally Sync Agent

### Purpose
Push approved transactions from the MediSync platform into Tally ERP via TDL XML API. This is the **only** agent with write access to Tally.

### Inputs
| Field | Type | Source |
|-------|------|--------|
| `approved_transactions` | list[Transaction] | B-08 approval output |
| `tally_connection_config` | TallyConfig | Encrypted secrets store |
| `company_id` | string | Multi-entity selector |

### Outputs
| Field | Type | Description |
|-------|------|-------------|
| `sync_result` | enum | `success \| partial \| failed` |
| `tally_voucher_ids` | list[string] | Tally-assigned IDs for synced entries |
| `failed_transactions` | list[dict] | With error details |
| `sync_timestamp` | datetime | UTC |
| `audit_log_entry` | AuditEvent | Written to B-14 |

### Tool Chain
```
ApprovedTransactions (from B-08)
  → Pre-sync validation:
      - OPA policy: user must have 'sync_to_tally' permission
      - Duplicate guard: check tally_sync_log for same transaction
      - Ledger existence: verify ledger names exist in Tally chart of accounts
  → TDL XML Payload Generator (custom Python + Jinja2 template)
  → HTTP POST to Tally Gateway (localhost:9000 or remote)
  → Response Parser: extract VoucherIDs or error codes
  → Auto-retry (3x, exponential backoff) on connection timeout
  → Audit Log Writer (B-14)
  → Sync Status Dashboard Update (WebSocket push)
```

### TDL Integration Notes
- Tally uses XML-based TDL (Tally Definition Language) for API integration.
- Journal entries use `<TALLYMESSAGE>` → `<VOUCHER TYPE="Journal">` structure.
- Purchase bills use `<VOUCHER TYPE="Purchase">`.
- The agent should use a dedicated Tally Gateway service that wraps TDL calls.
- **No direct database access to Tally** — always via TDL XML.

### Guardrails
- OPA hard policy: only users with `role:finance_head` or `role:accountant_lead` can initiate sync.
- Pre-sync always checks B-08 workflow_status == 'approved' before proceeding.
- Duplicate prevention: hash(vendor + amount + date + ledger) stored in `tally_sync_log`.
- Failed syncs: transaction moved to "Sync Failed" queue with retry UI.
- All syncs written to immutable audit log regardless of success/failure.

### HITL Gate
- The "Sync Now" / "One-Click Sync" button in the UI is the **explicit human trigger**.
- No autonomous sync without explicit user action.

---

## B-10 — Bank Reconciliation Agent

### Purpose
Automatically match bank statement entries to existing Tally ledger entries; assign confidence scores; surface unmatched items for manual resolution.

### Matching Algorithm
```
For each BankRow:
  1. Exact match: amount + date (±1 day) + description similarity > 0.85  → Matched (high)
  2. Amount match: same amount, date within 7 days                         → Suggested (medium)
  3. Fuzzy match: amount ± 10, similar description                         → Possible (low)
  4. No match found                                                         → Unmatched (manual)
```

### Tool Chain
```
Bank Statement Upload (CSV/PDF)
  → B-02 (OCR extraction for scanned statements)
  → Reconciliation Matcher:
      → PostgreSQL: SELECT matching Tally entries
      → rapidfuzz (MIT): description similarity scoring
      → LangChain: handle partial payments / multi-invoice splits
  → Confidence Scorer
  → Reconciliation Dashboard Update
  → Outstanding Items Report Generator (B-11)
```

### OSS Components
- **rapidfuzz** (MIT) — fast fuzzy string matching for description comparison
- **PostgreSQL** — windowed date-range matching queries
- **LangChain** — multi-invoice split reasoning

---

## C-06 — Data Quality Validation Agent

### Purpose
Run automated data quality checks on every ETL batch before data lands in the warehouse; block bad data; alert on anomalies.

### Check Categories
| Category | Examples |
|----------|---------|
| Completeness | Required fields not null (patient_id, amount, date) |
| Uniqueness | No duplicate voucher IDs, invoice numbers |
| Referential Integrity | All ledger codes exist in chart of accounts |
| Range Validation | Invoice amounts > 0 and < configurable max |
| Temporal Consistency | Transaction dates within valid range |
| Cross-Source Reconciliation | HIMS billing total ≈ Tally receipt total for same period |

### Tool Chain
```
ETL Pipeline Output (Airflow DAG task)
  → great_expectations (Apache-2.0) expectation suite
  → Validation Runner (computes pass/fail per expectation)
  → Anomaly Detector: statistical Z-score check on daily totals
  → On failure: block data load; emit alert (Apprise)
  → On warning: load with flag; add to data quality report
  → Audit Log (all validation results appended)
```

### OSS Components
- **great_expectations** (Apache-2.0) — declarative data quality assertions
- **Apache Airflow** — orchestrates validation as a DAG task
- **Apprise** — alert dispatch on critical failures

---

## D-04 — Autonomous AI Analyst Agent

### Purpose
Act as an on-demand data analyst. Accept a high-level business question, autonomously decompose it into sub-tasks, execute multi-step analysis (retrieve → analyse → compare → forecast → recommend), and return a structured report.

### Architecture: Multi-Agent (CrewAI)
```
User Prompt
  → Orchestrator Agent (CrewAI)
      ├── Data Retrieval Agent     → [A-01] Text-to-SQL → Postgres
      ├── Statistical Analysis Agent → statsmodels + pandas
      ├── Comparison Agent         → period-over-period, entity benchmarks
      ├── Forecasting Agent        → Prophet (up to 90-day horizon)
      └── Insight Narrative Agent  → LLM synthesis + Pydantic AnalystReport
  → AnalystReport (JSON)
  → Dashboard Renderer ([D-06])
  → Langfuse trace logging
```

### Input
```json
{
  "prompt": "Analyse pharmacy margin trends for Q4, compare to Q3, identify drivers, and forecast Q1.",
  "user_role": "pharmacy_manager",
  "date_context": "current_quarter"
}
```

### Output: `AnalystReport`
```json
{
  "summary": "Pharmacy margin declined 3.2% QoQ in Q4...",
  "key_findings": ["Finding 1...", "Finding 2..."],
  "data_tables": [...],
  "charts": [{"type": "line", "title": "Margin Trend", "data": {...}}],
  "forecast": {"q1_margin_estimate": 23.4, "confidence_interval": [21.8, 25.0]},
  "recommendations": ["Renegotiate supplier X contract...", "..."],
  "confidence_score": 0.88,
  "data_sources": ["tally.ledger_entries", "hims.pharmacy_dispensations"],
  "execution_trace_id": "lf_abc123"
}
```

### Guardrails
- All sub-agents operate read-only (OPA policy).
- Role-scoped: pharmacy_manager only receives pharmacy data.
- Langfuse traces every sub-agent step for auditability and debugging.
- Confidence score aggregated from all sub-agents; if < 0.75, report flagged as "preliminary".

---

## D-08 — Prescriptive Recommendations Agent

### Purpose
Go beyond descriptive insights to deliver specific, quantified, actionable recommendations — including root-cause analysis and expected business impact.

### Architecture: ReAct + Tool-Use Pattern (LangChain + CrewAI)
```
Insight Event (from D-07 or D-04)
  → Root Cause Analyser (LangChain ReAct loop)
      Tool: query_warehouse(sql) → Postgres
      Tool: compute_statistics(data) → statsmodels
      Tool: lookup_semantic_context(metric) → MetricFlow
  → Impact Quantifier
      Tool: run_forecast(scenario) → Prophet
      Tool: compute_delta(baseline, intervention) → Python
  → Recommendation Generator (LLM synthesis)
  → Pydantic RecommendationOutput
  → OPA Policy Gate: if recommendation triggers write action → HITL required
  → Notification (Apprise) → user / finance head
```

### Output: `RecommendationOutput`
```json
{
  "insight": "Item XYZ will stock out in 5 days",
  "root_cause": "35% surge in weekly sales due to seasonal demand",
  "recommendation": "Urgent reorder of 150 units from Supplier A by Thursday",
  "expected_impact": {
    "cash_flow_improvement_inr": 250000,
    "stockout_reduction_pct": 95
  },
  "confidence_score": 0.91,
  "action_required": true,
  "requires_approval": true,
  "approver_role": "pharmacy_manager"
}
```

### Guardrails
- `action_required: true` recommendations never auto-execute; always require explicit human approval.
- Tally write actions (e.g., creating PO) route through B-08 → B-09 pipeline.
- Impact estimates include confidence intervals; agent never presents point estimates without uncertainty bounds.

---

## General Agent Development Standards

All agents must conform to these standards:

### 1. Structured Output Contract
Every agent must return a **Pydantic model** as its output. No free-form string returns.

### 2. Observability
Every agent invocation must emit:
- An **OpenTelemetry span** (trace_id, span_id, latency, status)
- A **Langfuse trace** (prompt, response, token count, confidence)

### 3. Error Handling
```python
# Standard retry pattern (tenacity MIT)
@retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=2, max=10))
async def agent_invoke(input: AgentInput) -> AgentOutput:
    ...
```

### 4. HITL Integration
Any agent that produces an irreversible action must:
1. Check OPA policy before proceeding.
2. Write to `pending_human_review` queue if approved path is unclear.
3. Never auto-proceed on timeout — default to **block and notify**.

### 5. Audit Logging
All agents that interact with financial data must call `AuditLogger.log_event(agent_id, user_id, action, data_hash)` on every invocation.

### 6. Testing Requirements
- Unit tests with synthetic data (no production PII in tests).
- SQL injection test suite for A-01, A-02.
- OCR accuracy test suite with labelled invoice samples for B-02, B-03.
- Hallucination eval set for A-01, D-04, D-08 (LangSmith / Langfuse eval datasets).
