# Agent Specification — D-08: Prescriptive Recommendations Agent

**Agent ID:** `D-08`  
**Agent Name:** Prescriptive Recommendations Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 15  
**Priority:** P1 High  
**HITL Required:** No (recommends; never acts autonomously)  
**Status:** Draft

---

## 1. Purpose

Takes a detected insight or business situation and produces specific, evidence-backed recommended actions with expected outcomes. Generates structured recommendations that can be accepted and executed via D-10 (Insight-to-Action).

> **Addresses:** PRD §6.9.4, US30 — Closing the loop from insight to actionable recommendation.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Event trigger** | D-04, D-07 findings requiring recommendation |
| **Manual trigger** | User clicks "Get Recommendation" on an insight |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `finding` | `Finding` | D-04 / D-07 | ✅ |
| `context_data` | `[]ContextQuery` | SQL-fetched supporting context | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `available_actions` | `[]ActionDef` | D-10 action registry | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `recommendations` | `[]Recommendation` | Ranked recommendations with rationale |
| `expected_outcomes` | `[]Outcome` | Per-recommendation expected impact |
| `confidence` | `float64` | Recommendation confidence |

```go
type Recommendation struct {
    RecID           UUID      `json:"rec_id"`
    Title           string    `json:"title"`
    Rationale       string    `json:"rationale"`
    ActionRef       string    `json:"action_ref"`  // maps to D-10 ActionDef
    ExpectedImpact  string    `json:"expected_impact"`
    Priority        int       `json:"priority"`
    Confidence      float64   `json:"confidence"`
    EvidenceSQL     []string  `json:"evidence_sql"`
}
```

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | A-01 + A-07 | Gather supporting evidence (read-only SQL) |
| 2 | A-12 | Forecast impact of recommended action |
| 3 | Genkit Flow (`prescriptive-rec`) | Generate structured recommendations |
| 4 | A-05 Hallucination Guard | Validate recommendation narrative |
| 5 | A-06 Confidence Scoring | Score recommendations |

### Genkit Flow: `prescriptive-rec`
```
System Prompt:
  You are a hospital financial and operations advisor.
  Given a finding, available actions, and supporting data,
  produce ≤3 specific, measurable recommended actions.
  Each recommendation must include:
    - One-sentence title
    - 2-3 sentence rationale citing evidence
    - Referenced action from the action registry
    - Expected impact (quantified where possible)
  Output JSON matching the Recommendation schema. Never fabricate data.

Input:
  Finding: {finding}
  Supporting Data: {context_data}
  Available Actions: {available_actions}
```

---

## 6. Guardrails

- All evidence SQL queries run read-only.
- Recommendations reference only actions in D-10's registry.
- No autonomous execution — always surfaces for human review.
- Narratives validated by A-05 before display.
- Confidence < 0.65 recommendations flagged with explicit uncertainty.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Recommendation acceptance rate | ≥ 60% |
| Recommendation quality (analyst rating ≥ 4/5) | ≥ 75% of recommendations |
| False recommendation rate | < 10% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | A-01, A-05, A-06, A-07, A-12, D-10 |
| **Consumed by** | D-04, D-07, D-10 |
