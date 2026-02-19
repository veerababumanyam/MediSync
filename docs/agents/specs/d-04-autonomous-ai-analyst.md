# Agent Specification — D-04: Autonomous AI Analyst (Spotter) Agent

**Agent ID:** `D-04`  
**Agent Name:** Autonomous AI Analyst (Spotter) Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 14  
**Priority:** P1 High  
**HITL Required:** No (autonomous, alerts humans)  
**Status:** Draft

---

## 1. Purpose

Proactively and autonomously analyses the data warehouse on a scheduled basis, discovering significant trends, anomalies, and emerging risks without being asked. Generates a prioritised "Spotter Brief" report surfaced in the Analytics dashboard.

> **Addresses:** PRD §6.9.2, US28, US30 — Proactive AI-driven financial and operational intelligence.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled |
| **Scheduled trigger** | `0 6 * * *` (6 AM daily) |
| **Event trigger** | Any A-13 critical anomaly alert |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `tenant_id` | `string` | Scheduler | ✅ |
| `analysis_scope` | `[]string` | Config (default: all modules) | ✅ |
| `days_lookback` | `int` | Config (default: 30) | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `spotter_brief` | `SpotterBrief` | Prioritised summary with findings |
| `findings` | `[]Finding` | Individual insights with evidence |
| `risk_score` | `float64` | Overall period risk score (0–1) |
| `recommended_actions` | `[]Action` | D-08 recommendations linked to findings |

```go
type SpotterBrief struct {
    BriefID        UUID      `json:"brief_id"`
    TenantID       string    `json:"tenant_id"`
    GeneratedAt    time.Time `json:"generated_at"`
    RiskScore      float64   `json:"risk_score"`
    Findings       []Finding `json:"findings"`
    RecommActions  []Action  `json:"recommended_actions"`
    ModelVersion   string    `json:"model_version"`
}

type Finding struct {
    FindingID   UUID    `json:"finding_id"`
    Category    string  `json:"category"`  // anomaly|trend|risk|opportunity
    Severity    string  `json:"severity"`  // critical|high|medium|low
    Title       string  `json:"title"`
    Detail      string  `json:"detail"`
    Evidence    []Evidence `json:"evidence"` // SQL + chart config
    Confidence  float64 `json:"confidence"`
}
```

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | `robfig/cron` (Go) | Schedule trigger |
| 2 | D-13 | Run monitoring checks; import recent alerts |
| 3 | A-13 | Import anomaly signals from last 24h |
| 4 | A-12 | Trend forecasts for key KPIs |
| 5 | Genkit Flow (`spotter-analysis`) | Orchestrate multi-step reasoning |
| 6 | A-01 + A-02 | SQL evidence queries |
| 7 | A-06 | Confidence scoring per finding |
| 8 | D-08 | Prescriptive recommendations for top-3 findings |
| 9 | A-03 | Viz configs for evidence charts |
| 10 | A-05 | Hallucination check on all narrative text |
| 11 | Apprise | Push Spotter Brief notification |

### Genkit Flow: `spotter-analysis`
```
Genkit Flow:
  1. Load monitored KPI list from D-09 Semantic Layer
  2. For each KPI: run A-13 anomaly check + A-12 trend check
  3. Cluster findings (deduplicate related signals)
  4. Rank findings by (impact × confidence)
  5. Generate narrative for top-10 findings (LLM)
  6. Run A-05 on all LLM-generated text
  7. Attach SQL evidence via A-01
  8. Invoke D-08 for top-3 findings → get recommended_actions
  9. Compute composite risk_score
  10. Assemble SpotterBrief
```

---

## 6. Guardrails

- All SQL evidence queries run as `medisync_readonly`.
- No autonomous write actions — findings are advisory only.
- LLM narrative validated by A-05 before inclusion.
- Confidence < 0.7 findings labelled with explicit uncertainty badge.
- Finding evidence SQL stored immutably for audit.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Finding relevance (analyst rating) | ≥ 4/5 avg |
| Finding actionability rate | ≥ 80% |
| False alert rate (irrelevant findings) | < 10% |
| Spotter Brief generation time | < 5 min |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | A-01, A-02, A-03, A-05, A-06, A-12, A-13, D-08, D-09, D-13, Apprise |
| **Consumed by** | Analytics dashboard, Finance Head, Management |
