# Agent Specification — A-06: Confidence Scoring Agent

**Agent ID:** `A-06`  
**Agent Name:** Confidence Scoring Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P1 High  
**HITL Required:** Yes — routes low-confidence results to human review  
**Status:** Draft

---

## 1. Purpose

Attaches a calibrated confidence score (0–1) to every AI-generated answer. Routes answers scoring below 0.70 to a manual review queue rather than returning them directly to the user.

> **Addresses:** PRD §10 (NFR) — Confidence thresholding and human-review fallback.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (post-generation) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `agent_output` | `any` | Any agent returning a result | ✅ |
| `agent_id` | `string` | Caller | ✅ |
| `raw_llm_logprobs` | `[]float64` | LLM response metadata | ⬜ |
| `execution_metadata` | `ExecMeta` | Query timing, retry count | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `confidence_score` | `float64` | 0.0–1.0 |
| `confidence_level` | `enum` | `high (≥0.95)` / `medium (0.70–0.94)` / `low (<0.70)` |
| `hitl_required` | `bool` | True if low |
| `ui_badge` | `enum` | `green` / `amber` / `red` |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Log-probability aggregator (Go) | Average token-level log-probs from LLM |
| 2 | Heuristic scorer | Penalise for retries, schema mismatches, short result sets |
| 3 | Calibration layer | Isotonic regression calibration on validation set |

---

## 6. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `confidence_score < 0.70` |
| **Notified role** | `analyst` or submitting user's manager |
| **SLA** | 4h |
| **On approval** | Result released to user |

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| ECE (Expected Calibration Error) | < 0.05 |
| HITL escalation rate in production | < 15% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (shared confidence library) |
| **Consumed by** | A-01, B-02, B-05, B-10, D-04 |
