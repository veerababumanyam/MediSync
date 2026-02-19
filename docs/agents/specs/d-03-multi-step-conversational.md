# Agent Specification — D-03: Multi-Step Conversational Analysis Agent

**Agent ID:** `D-03`  
**Agent Name:** Multi-Step Conversational Analysis Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 14  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Enables an iterative dialogue for analytical exploration — maintaining conversation context across turns so that follow-up queries ("drill into Q2", "break it down by region", "why did it drop?") build on prior responses without user repeating context.

> **Addresses:** PRD §6.9.1, US27, US31 — Stateful conversational analytics with multi-turn reasoning.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | User message in the Analytics Chat interface |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `message` | `string` | User input | ✅ |
| `session_id` | `string` | Frontend | ✅ |
| `conversation_history` | `[]Turn` | Redis session store | ✅ |
| `user_id` | `string` | JWT | ✅ |

```go
type Turn struct {
    Role    string `json:"role"`    // user | assistant
    Content string `json:"content"`
    SQL     string `json:"sql,omitempty"`
    ResultSummary string `json:"result_summary,omitempty"`
}
```

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `response_text` | `string` | Natural language answer |
| `viz_config` | `VizConfig` | Chart configuration (if applicable) |
| `sql_used` | `string` | SQL executed (shown in debug panel) |
| `updated_history` | `[]Turn` | Context to store in Redis |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Redis | Load/save conversation history (TTL 2h) |
| 2 | D-02 | Entity extraction from new turn |
| 3 | Genkit Flow (`conv-analysis`) | Context-aware query understanding + response |
| 4 | A-01 + A-02 | SQL generation + self-correction |
| 5 | A-07 Drill-Down | Decompose drill follow-ups |
| 6 | A-03 Viz Routing | Select appropriate chart |
| 7 | A-05 Hallucination Guard | Validate response |

---

## 6. Guardrails

- Max 20 turns per session before context reset prompt.
- Context window managed by Genkit conversation history trimming (last 10 turns).
- SQL limited to user's OPA scope at every turn.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Context retention (follow-up accuracy) | ≥ 92% |
| P95 Response Latency | < 8s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | Redis, D-02, A-01, A-02, A-03, A-05, A-07 |
| **Consumed by** | Analytics Chat UI |
