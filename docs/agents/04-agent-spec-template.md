# Agent Specification — [AGENT-ID: Agent Name]

> **How to use this template:**  
> Copy this file to `docs/agents/specs/[agent-id]-[agent-name].md`.  
> Replace all `[PLACEHOLDER]` values. Remove this usage note before committing.  
> **Convention:** Agent IDs follow the pattern `[Module Letter]-[2-digit number]` (e.g. `A-01`, `B-09`, `D-04`).

---

**Agent ID:** `[X-NN]`  
**Agent Name:** `[Human-readable name]`  
**Module:** `[A: Conversational BI | B: AI Accountant | C: Easy Reports | D: Advanced Search Analytics]`  
**Phase:** `[Phase number from PRD roadmap]`  
**Priority:** `[P0 Critical | P1 High | P2 Medium | P3 Low]`  
**HITL Required:** `[Yes – describe gate | No]`  
**Status:** `[Draft | In Development | Testing | Production]`

---

## 1. Purpose

_One to three sentences describing **what** this agent does and **why** it exists. Reference the specific PRD user story or requirement this agent satisfies._

> **Addresses:** PRD §[section] — User Story US[nn]: "[story title]"

---

## 2. Trigger

_How and when is this agent invoked?_

| Property | Value |
|----------|-------|
| **Trigger type** | `[Manual / Scheduled / Event-driven / Upstream-agent-output]` |
| **Manual trigger** | `[Button label / API endpoint / Chat command]` |
| **Scheduled trigger** | `[Cron expression, e.g. 0 2 * * *]` |
| **Event trigger** | `[Queue/topic name + event payload description]` |
| **Calling agent** | `[Agent-ID of upstream agent, or "User"]` |

---

## 3. Inputs

| Input | Type | Source | Required | Notes |
|-------|------|--------|:--------:|-------|
| `[input_name]` | `[str / int / Pydantic model / file]` | `[Warehouse / Form / API / Prior agent output]` | ✅ / ⬜ | `[Validation rules, limits]` |

### Input Pydantic Model

```python
from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime

class [AgentName]Input(BaseModel):
    """Input contract for [Agent Name]."""
    
    # [field]: [description]
    field_one: str = Field(..., description="[description]")
    field_two: Optional[int] = Field(None, description="[description]")
    user_id: str = Field(..., description="Keycloak user ID for audit logging")
    session_id: str = Field(..., description="Session ID for trace correlation")
```

---

## 4. Outputs

| Output | Type | Destination | Notes |
|--------|------|-------------|-------|
| `[output_name]` | `[str / Pydantic model / file / event]` | `[Warehouse / Queue / UI / Next agent]` | |

### Output Pydantic Model

```python
from pydantic import BaseModel, Field
from typing import Optional, List
from enum import Enum

class ConfidenceLevel(str, Enum):
    HIGH   = "high"     # >= 0.95
    MEDIUM = "medium"   # 0.70 – 0.94
    LOW    = "low"      # < 0.70

class [AgentName]Output(BaseModel):
    """Output contract for [Agent Name]."""
    
    agent_id: str = "[X-NN]"
    success: bool
    confidence: ConfidenceLevel
    confidence_score: float = Field(..., ge=0.0, le=1.0)
    result: Optional[[ResultType]] = None
    warnings: List[str] = Field(default_factory=list)
    hitl_required: bool = False
    hitl_reason: Optional[str] = None
    trace_id: str  # OpenTelemetry trace ID
    
    # [domain-specific fields]
```

---

## 5. Tool Chain

_List all OSS tools used in execution order._

| Step | Tool / Library | Version | License | Purpose |
|------|---------------|---------|---------|---------|
| 1 | `[tool name]` | `[version]` | `[license]` | `[what it does in this agent]` |

### Tool Install (add to `requirements.txt`)

```
[package-name]==[version]
```

---

## 6. Architecture Diagram

_ASCII or Mermaid diagram showing data flow through this agent._

```
Input
  ↓
[Step 1: Tool / operation]
  ↓
[Step 2: Tool / operation]
  ↓
[Decision: confidence check]
  ├─ High → Return result
  └─ Low  → HITL queue
```

---

## 7. System Prompt

_For LLM-backed agents. Skip this section for non-LLM agents (replace with "N/A — No LLM involved")._

```
You are [agent role description]. Your job is to [specific task].

RULES:
1. [Hard rule 1 — use MUST / MUST NOT language]
2. [Hard rule 2]
3. MUST NEVER [critical prohibition, e.g. perform any database writes]
4. If you are not confident, respond with: {"hitl_required": true, "reason": "[explanation]"}

CONTEXT PROVIDED:
- [What context/data is injected at runtime]
- [Schema / metric definitions / company-specific terms]

OUTPUT FORMAT:
Respond ONLY with valid JSON matching this schema:
{
  "[field1]": "[type]",
  "[field2]": "[type]",
  "confidence_score": "[float 0.0–1.0]",
  "reasoning": "[brief explanation of your answer]"
}
```

---

## 8. Guardrails

_List all safety and quality controls — in approximate check order._

| # | Guard | Type | Trigger | Action |
|---|-------|------|---------|--------|
| 1 | [Guard name] | `[Pre-execution / Post-execution / Continuous]` | [Condition] | [Response / action taken] |
| 2 | OPA Policy Check | Pre-execution | Every request | Deny and log if policy violation |
| 3 | Confidence Threshold | Post-execution | confidence_score < 0.70 | Route to HITL queue |
| 4 | Audit Log Write | Post-execution | Always | Write event to `audit_log` table |

---

## 9. HITL Gate

_Describe the human-in-the-loop checkpoint(s) for this agent._

**HITL Required:** `[Yes / No]`

| Property | Value |
|----------|-------|
| **Gate type** | `[Confidence gate / Threshold gate / Always / Never]` |
| **Trigger condition** | `[confidence < X] OR [field == Y] OR [amount > Z]` |
| **Notified role(s)** | `[accountant / finance_head / analyst]` |
| **Notification method** | `[Apprise in-app + email / In-app only]` |
| **SLA** | `[Max time before escalation, e.g. 24h]` |
| **Escalation path** | `[Who gets notified if SLA breached]` |
| **Approval actions** | `[Approve / Reject / Edit + Approve]` |
| **On approval** | `[Next agent ID triggered or result released to user]` |
| **On rejection** | `[Workflow cancelled + audit log entry]` |

---

## 10. Evaluation Criteria

_Quantitative success thresholds to be measured during testing and monitoring._

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| Accuracy / F1 | `[e.g. ≥ 95%]` | `[Test dataset + golden answers]` |
| P95 Latency | `[e.g. < 5 seconds]` | `[Langfuse trace timing]` |
| HITL Escalation Rate | `[e.g. < 15%]` | `[Langfuse + audit_log query]` |
| False Positive Rate | `[e.g. < 5%]` | `[Human review of production samples]` |
| User Correction Rate | `[e.g. < 10%]` | `[UI feedback log]` |

---

## 11. Error Handling

| Error Scenario | HTTP Code | User Message | Internal Action |
|---------------|-----------|--------------|----------------|
| OPA policy denial | 403 | "You do not have permission for this action." | Log to audit_log with denial_reason |
| LLM timeout | 504 | "Analysis is taking longer than expected. Try again." | Retry with tenacity (3×); alert on-call if all fail |
| Low confidence | 200 | "Result requires review before use." (warning banner) | Set hitl_required=true; queue for human |
| Tool/API failure | 500 | "A system error occurred. ID: [trace_id]" | Log full stack trace; Apprise alert to admin |
| Input validation error | 422 | "Invalid input: [field]: [reason]" | Return Pydantic validation error detail |

### Retry Pattern (tenacity)

```python
from tenacity import retry, stop_after_attempt, wait_exponential, retry_if_exception_type

@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=1, min=2, max=30),
    retry=retry_if_exception_type((TimeoutError, ConnectionError)),
    reraise=True
)
def [agent_function](...):
    ...
```

---

## 12. Observability

### 12.1 OpenTelemetry Tracing

```python
from opentelemetry import trace

tracer = trace.get_tracer("[agent-id]-agent")

with tracer.start_as_current_span("[agent-id]-execute") as span:
    span.set_attribute("agent.id", "[X-NN]")
    span.set_attribute("user.id", input.user_id)
    span.set_attribute("session.id", input.session_id)
    # ... agent execution ...
    span.set_attribute("agent.confidence_score", result.confidence_score)
    span.set_attribute("agent.hitl_required", result.hitl_required)
```

### 12.2 Langfuse Evaluation Logging

```python
from langfuse import Langfuse

langfuse = Langfuse()

trace = langfuse.trace(
    name="[agent-id]-[agent-name]",
    user_id=input.user_id,
    session_id=input.session_id,
    metadata={"agent_id": "[X-NN]"}
)

generation = trace.generation(
    name="[step-name]",
    model="[model-name]",
    input=prompt,
    output=llm_response,
    usage={"promptTokens": ..., "completionTokens": ...}
)
```

### 12.3 Key Metrics to Track (Prometheus/Grafana)

```
medisync_agent_requests_total{agent="[X-NN]", status="success|failure"}
medisync_agent_latency_seconds{agent="[X-NN]", quantile="0.5|0.95|0.99"}
medisync_agent_confidence_score{agent="[X-NN]"}
medisync_agent_hitl_escalations_total{agent="[X-NN]"}
```

---

## 13. Audit Log Integration

Every agent invocation **must** write to the audit_log table:

```python
from datetime import datetime, timezone
import hashlib, json

def write_audit_log(
    conn, user_id, user_role, agent_id, action,
    resource_id, data_after, ip_address, session_id, trace_id, status,
    denial_reason=None
):
    data_hash = hashlib.sha256(
        json.dumps(data_after, sort_keys=True).encode()
    ).hexdigest() if data_after else None
    
    conn.execute("""
        INSERT INTO audit_log (
            user_id, user_role, agent_id, action,
            resource_id, data_after, data_hash,
            ip_address, session_id, trace_id, status, denial_reason
        ) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
    """, (
        user_id, user_role, agent_id, action,
        resource_id, json.dumps(data_after), data_hash,
        ip_address, session_id, trace_id, status, denial_reason
    ))
```

---

## 14. Testing Checklist

### Unit Tests

- [ ] Happy path: valid input → correct output + correct Pydantic model returned
- [ ] Low confidence path: input that should trigger HITL → hitl_required=True
- [ ] OPA policy denial: unauthorized role → 403 + audit log entry
- [ ] LLM timeout: mock timeout → tenacity retries → final failure handled
- [ ] Input validation: missing required fields → 422 with field name

### Integration Tests

- [ ] End-to-end with test database snapshot
- [ ] Audit log entry verifiable after each test run
- [ ] HITL notification sent to correct role
- [ ] Downstream agent triggered correctly on success

### LLM Evaluation (LangChain + Langfuse)

- [ ] Golden dataset: minimum [N] question/answer pairs
- [ ] Accuracy ≥ [X]% on golden dataset
- [ ] No PII leakage in any LLM prompt/response (automated redaction test)
- [ ] Off-topic deflection: 100% of non-domain queries deflected

### Performance Tests

- [ ] P95 latency < [X]s under [N] concurrent requests
- [ ] No memory leak after 1,000 sequential requests

---

## 15. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | `[Python 3.11 / FastAPI service / Airflow DAG / Celery worker]` |
| **Celery Queue** | `[queue name if applicable]` |
| **Airflow DAG ID** | `[dag_id if applicable]` |
| **Docker image** | `medisync/agent-[agent-id]:latest` |
| **Env vars required** | `[VAR_NAME=description, ...]` |
| **Secrets (Vault)** | `[secret/medisync/[agent-id]/...]` |
| **DB connection** | `medisync_readonly` role (SELECT only) |
| **Depends on agents** | `[comma-separated Agent IDs, or "None"]` |
| **Consumed by agents** | `[comma-separated Agent IDs, or "User-facing"]` |

---

## 16. Changelog

| Date | Author | Change |
|------|--------|--------|
| [YYYY-MM-DD] | [Name] | Initial spec created from template |

---

*Template version: 1.0 | Source: `docs/agents/04-agent-spec-template.md`*
