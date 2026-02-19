# Agent Specification — D-07: Insight Discovery & Prioritisation Agent

**Agent ID:** `D-07`  
**Agent Name:** Insight Discovery & Prioritisation Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 15  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Continuously scans the data warehouse for statistically significant insights (trends, correlations, anomalies, emerging risks) and presents the top-5 most relevant, actionable insights to each user role on the home dashboard.

> **Addresses:** PRD §6.9.4 — Surfacing prioritised insights without the user needing to ask.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled + Dashboard load |
| **Scheduled trigger** | `0 7 * * *` (7 AM daily) |
| **Event trigger** | User opens dashboard (stale check: last run > 4h) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `user_id` | `string` | JWT | ✅ |
| `role` | `string` | JWT | ✅ |
| `entity_ids` | `[]string` | OPA scope | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `insights` | `[]Insight` | Top-5 prioritised insights |
| `insight_feed_id` | `UUID` | Stored feed ID for diff tracking |

```go
type Insight struct {
    InsightID   UUID    `json:"insight_id"`
    Title       string  `json:"title"`
    Narrative   string  `json:"narrative"`
    Priority    int     `json:"priority"`  // 1-5
    Category    string  `json:"category"`  // trend|anomaly|risk|opportunity
    Confidence  float64 `json:"confidence"`
    VizConfig   VizConfig `json:"viz_config"`
}
```

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | A-13 Anomaly results | Import recent anomalies |
| 2 | A-12 Trend results | Import trend signals |
| 3 | D-07 scorer (Go) | Score by impact × confidence × relevance to role |
| 4 | Genkit Flow (`insight-narrate`) | Generate natural language narrative |
| 5 | A-05 | Validate narrative |
| 6 | A-03 | Viz config for each insight |

---

## 6. Guardrails

- Insights scoped to user's OPA-allowed entities.
- Narratives validated by A-05 before display.
- Expired insights (> 48h) auto-hidden.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Insight click-through rate | ≥ 40% |
| User dismissal rate (irrelevant) | < 20% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | A-03, A-05, A-12, A-13 |
| **Consumed by** | Daily digest, Analytics dashboard home |
