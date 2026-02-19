# Agent Specification — A-01: Text-to-SQL Agent

**Agent ID:** `A-01`  
**Agent Name:** Text-to-SQL Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P0 Critical  
**HITL Required:** Yes — confidence < 0.70 or PII table access  
**Status:** Draft

---

## 1. Purpose

Converts a natural language business question into a safe, read-only SQL query against the MediSync data warehouse, executes it, and returns structured results with a chart-type suggestion and explanation.

> **Addresses:** PRD §5.1 — User Stories US1–US8: "As a user I want to ask questions in plain English and get instant data answers."

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User submits a query in the BI chat interface |
| **Calling agent** | User |

---

## 3. Inputs

| Input | Type | Source | Required | Notes |
|-------|------|--------|:--------:|-------|
| `user_query` | `string` | Chat interface | ✅ | Max 500 chars |
| `user_role` | `string` | JWT (Keycloak) | ✅ | Used for column masking |
| `schema_context` | `JSON` | Schema Context Cache | ✅ | Pre-loaded DB schema |
| `semantic_context` | `JSON` | MetricFlow / Semantic Layer Registry | ✅ | Metric definitions |
| `conversation_history` | `list[Message]` | Session memory | ⬜ | Last N turns for follow-up queries |
| `user_id` | `string` | JWT | ✅ | Audit logging |
| `session_id` | `string` | Session store | ✅ | Trace correlation |

```go
type TextToSQLInput struct {
    UserQuery           string        `json:"user_query" validate:"required,max=500"`
    UserRole            string        `json:"user_role" validate:"required"`
    SchemaContext       SchemaContext `json:"schema_context" validate:"required"`
    SemanticContext     []Metric      `json:"semantic_context"`
    ConversationHistory []Message     `json:"conversation_history"`
    UserID              string        `json:"user_id" validate:"required"`
    SessionID           string        `json:"session_id" validate:"required"`
}
```

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `sql_query` | `string` | Logged; executed against Postgres |
| `result_set` | `[]map[string]any` | BI chart renderer |
| `chart_type` | `enum` | Frontend viz router |
| `confidence_score` | `float64` | Confidence router |
| `explanation` | `string` | Chat response UI |
| `hitl_required` | `bool` | HITL queue |
| `error` | `*string` | Error handler |

```go
type TextToSQLOutput struct {
    AgentID         string         `json:"agent_id"` // "A-01"
    Success         bool           `json:"success"`
    SQLQuery        string         `json:"sql_query"`
    ResultSet       []map[string]any `json:"result_set"`
    ChartType       ChartType      `json:"chart_type"` // bar|line|pie|table|scatter
    ConfidenceScore float64        `json:"confidence_score"`
    Explanation     string         `json:"explanation"`
    HITLRequired    bool           `json:"hitl_required"`
    HITLReason      *string        `json:"hitl_reason,omitempty"`
    TraceID         string         `json:"trace_id"`
    Error           *string        `json:"error,omitempty"`
}
```

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | A-04 (Domain Normaliser) | Internal | Normalise healthcare/accounting synonyms |
| 2 | Genkit Flow (`text-to-sql`) | Apache-2.0 | LLM-based SQL generation |
| 3 | Ollama / vLLM | MIT / Apache-2.0 | Local LLM inference |
| 4 | sqlparse (via CGO or subprocess) | BSD | SQL AST assertion — SELECT only |
| 5 | PostgreSQL (`medisync_readonly` role) | PostgreSQL | Execute validated query |
| 6 | A-02 (Self-Correction) | Internal | Retry on DB error (max 3) |
| 7 | A-03 (Visualization Router) | Internal | Emit chart-type token |
| 8 | OPA sidecar | Apache-2.0 | Column masking + read-only policy |

---

## 6. Architecture Diagram

```
UserQuery
  → [A-04] Domain Terminology Normaliser
  → OPA Policy Check (read-only, role filter)
  → Genkit Flow: text-to-sql
      LLM (Ollama/vLLM) + schema_context + semantic_context
  → SQL AST Validator (SELECT-only assertion)
  → PostgreSQL Executor (medisync_readonly)
  → [A-02] Self-Correction loop (on error, max 3 retries)
  → [A-03] Visualization Router
  → TextToSQLOutput struct
```

---

## 7. System Prompt

```
You are an expert SQL analyst for MediSync, a healthcare and accounting platform.
You have READ-ONLY access to a PostgreSQL data warehouse containing:
  - HIMS data: patients, appointments, billing, pharmacy dispensations
  - Tally data: ledgers, vouchers, inventory, sales, receipts

RULES:
1. Generate ONLY SELECT statements. NEVER use INSERT, UPDATE, DELETE, DROP, CREATE, TRUNCATE, ALTER.
2. Always apply the user's role filter: {{ user_role_filter }}
3. Use metric definitions from the Semantic Layer when available.
4. If the question is ambiguous, ask one clarifying question before generating SQL.
5. If the question is not data-related, respond with: {"off_topic": true}

Available schema:
{{ schema_context }}

Available metrics:
{{ semantic_context }}

OUTPUT FORMAT:
Respond ONLY with valid JSON:
{
  "sql": "<SELECT ...>",
  "chart_type": "bar|line|pie|table|scatter",
  "explanation": "<plain English explanation>",
  "confidence_score": <float 0.0–1.0>,
  "reasoning": "<why this SQL answers the question>"
}
```

---

## 8. Guardrails

| # | Guard | Type | Trigger | Action |
|---|-------|------|---------|--------|
| 1 | SQL AST read-only assertion | Pre-execution | Every query | Reject if non-SELECT DML detected |
| 2 | OPA policy `medisync.bi.read_only` | Pre-execution | Every request | Hard-block at DB connection level |
| 3 | Off-topic classifier | Pre-execution | LLM response `off_topic=true` | Return canned deflection; no SQL executed |
| 4 | Column masking | Pre-execution | Always | OPA strips PII/cost-price columns per role |
| 5 | Confidence gate | Post-execution | `confidence_score < 0.70` | Set `hitl_required=true`; add warning banner |
| 6 | Audit log write | Post-execution | Always | Append to `audit_log` table |

---

## 9. HITL Gate

| Property | Value |
|----------|-------|
| **Gate type** | Confidence gate |
| **Trigger condition** | `confidence_score < 0.70` OR query touches patient PII tables |
| **Notified role(s)** | `analyst`, `finance_head` |
| **Notification method** | In-app banner: "Your query is being reviewed" |
| **SLA** | 4h |
| **Escalation path** | `admin` notified if unreviewed after 4h |
| **Approval actions** | Approve / Edit + Approve / Reject |
| **On approval** | Result released to user |
| **On rejection** | User notified with reason |

---

## 10. Evaluation Criteria

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| SQL correctness (executes without error) | ≥ 98% | Golden dataset CI run |
| Business intent accuracy (human eval) | ≥ 95% | Langfuse eval dataset (min 200 Q&A pairs) |
| P95 Latency | < 5 seconds | Genkit / OpenTelemetry trace |
| Hallucination rate | < 1% | Automated off-topic + wrong-table detection |
| Off-topic false positive rate | < 2% | Human-labelled off-topic test set |

---

## 11. Error Handling

| Error Scenario | Status | User Message | Internal Action |
|---------------|--------|--------------|----------------|
| OPA policy denial | 403 | "You do not have permission to access this data." | Log to audit_log |
| LLM timeout | 504 | "Analysis is taking longer than expected. Please try again." | Retry 3× exponential backoff |
| SQL execution error | 200 | "Could not process your query. Refining..." | Trigger A-02 self-correction |
| All retries failed | 500 | "Query failed. Reference: [trace_id]" | Alert on-call via Apprise |
| Off-topic query | 200 | "I can only answer questions about your business data." | No LLM invoked; classifier short-circuit |

---

## 12. Observability

- **OpenTelemetry:** Span `a-01-text-to-sql` with attributes: `user.role`, `confidence_score`, `chart_type`, `sql_length`, `rows_returned`
- **Genkit:** Built-in flow tracing captures prompt, response, token usage
- **Metrics:** `medisync_agent_requests_total{agent="A-01"}`, `medisync_agent_latency_seconds{agent="A-01"}`

---

## 13. Audit Log

Logged fields: `agent_id=A-01`, `action=sql_query`, `resource_id=session_id`, `data_after={sql, chart_type, confidence_score}`, `status=success|failed|blocked`

---

## 14. Testing Checklist

- [ ] Happy path: "What was last month's revenue?" → valid SELECT + correct result
- [ ] Off-topic: "What is the capital of France?" → deflected, no SQL
- [ ] SQL injection attempt in user_query → rejected by AST guard
- [ ] PII table query by `viewer` role → OPA blocks, 403
- [ ] Low confidence query → `hitl_required=true` in response
- [ ] LLM timeout → retries 3× → graceful failure
- [ ] Follow-up query using conversation_history → correct contextual SQL

---

## 15. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (chi router) |
| **Genkit Flow ID** | `text-to-sql` |
| **DB connection** | `medisync_readonly` role |
| **Depends on agents** | A-02, A-03, A-04 |
| **Consumed by** | User (chat UI), D-04 |
| **Env vars** | `GENKIT_OLLAMA_URL`, `POSTGRES_READONLY_DSN`, `OPA_SIDECAR_URL` |
