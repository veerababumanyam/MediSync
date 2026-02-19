# Agent Specification — A-05: Hallucination Guard Agent

**Agent ID:** `A-05`  
**Agent Name:** Hallucination Guard Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P0 Critical  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Detects off-topic, harmful, or non-business queries before the SQL agent is invoked, and short-circuits with a canned deflection response. Prevents LLM invocation waste and protects the platform from prompt injection.

> **Addresses:** PRD §10 (NFR) — AI must only answer business data questions.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (runs before LLM call) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `user_query` | `string` | User input | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `is_on_topic` | `bool` | If false, A-01 short-circuits |
| `deflection_message` | `*string` | Shown to user when off-topic |
| `confidence` | `float64` | Classifier confidence |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Fine-tuned binary classifier (Go ONNX runtime) | Apache-2.0 | Fast on-topic / off-topic classification |
| 2 | Prompt injection pattern detector (regex) | Internal | Block jailbreak attempts |

**Model:** DistilBERT fine-tuned on MediSync domain queries (MIT base license). Served via ONNX Runtime Go binding.

### Off-topic categories deflected
- General knowledge ("What is the capital of India?")
- Medical advice ("Should I prescribe X?")
- Code generation requests (these go to D-11 if enabled)
- Prompt injection attempts ("Ignore previous instructions...")

---

## 6. Guardrails

- If classifier confidence < 0.85, treat as on-topic (pass through) to avoid over-blocking legitimate queries.
- Prompt injection patterns are regex-matched (deterministic) — not classifier-dependent.
- All deflections logged to audit_log.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Off-topic detection rate | ≥ 98% |
| False positive (on-topic blocked) | < 1% |
| Prompt injection block rate | 100% |
| Latency | < 50ms |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (inline middleware in A-01 flow) |
| **Model** | `models/hallucination-guard-distilbert.onnx` |
| **Depends on** | None |
| **Consumed by** | A-01 |
