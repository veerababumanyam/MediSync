# Phase 14 — Autonomous AI Agents & Deep Research

**Phase Duration:** Weeks 44–46 (3 weeks)  
**Module(s):** D — Advanced Search Analytics  
**Status:** Planning  
**Milestone:** M10 — Autonomous AI agents live  
**Depends On:** Phase 13 (Semantic layer + NL search operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Elevate Module D from reactive search to proactive autonomous intelligence. Deploy L3 autonomous agents capable of multi-step analytical workflows, deep research (statistical analysis + cross-domain synthesis), and—critically—HITL-gated insight-to-action that can trigger real downstream effects (scheduled reports, Tally-side alerts, approval queue injection). Deploy continuous monitoring agents that watch KPIs without human prompting.

---

## 2. Scope

### In Scope
- D-03: Multi-Step Analysis Agent (orchestrates multi-hop queries)
- D-04: Autonomous Analyst / Spotter Agent (proactive monitoring)
- D-05: Deep Research Agent (statistical analysis + synthesis)
- D-10: Insight-to-Action Agent (HITL-gated downstream actions)
- D-13: Scheduled Monitoring Agent (cron-triggered watchdog)
- Multi-step query orchestration (ADK + A2A)
- Proactive insight surface in dashboard ("Spotter" panel)
- Statistical toolkit: ARIMA, regression, correlation, Z-score, IQR
- HITL approval UI for recommended actions

### Out of Scope
- D-06 to D-08, D-11, D-12, D-14 (next phases)
- Autonomous write actions without HITL gate

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | D-03 Multi-Step Analysis Agent | AI Engineer | Completes 5-step analytical chain without user re-prompting |
| D-02 | D-04 Autonomous Analyst (Spotter) | AI Engineer | Fires proactive insight within 1 minute of anomaly crossing threshold |
| D-03 | D-05 Deep Research Agent | AI Engineer | Generates statistical research report with citations on demand |
| D-04 | D-10 Insight-to-Action Agent | AI + Backend | Recommended actions presented; action taken only after HITL approval |
| D-05 | D-13 Scheduled Monitoring Agent | AI + Data Eng | Runs 12 scheduled monitoring jobs; alerts on deviation |
| D-06 | Spotter Panel UI | Frontend Eng | Insight cards shown in dashboard; dismiss/accept/investigate actions |
| D-07 | HITL Action Approval UI | Frontend Eng | Pending actions list; approve/reject with reason; full audit trail |
| D-08 | Multi-step orchestration framework | Backend + AI Eng | ADK-based orchestrator with A2A step-to-step communication |
| D-09 | Statistical toolkit library | AI Eng + Data Eng | Python micro-service: ARIMA, linear regression, correlation, Z-score, IQR |

---

## 4. AI Agents Deployed

### D-03 — Multi-Step Analysis Agent

**Purpose:** Decompose a complex analytical question into a directed acyclic graph (DAG) of sub-queries, execute them in order (or in parallel where safe), and synthesise a final answer.

**Example workflow:**
```
User: "Why did pharmacy revenue drop in December compared to November?"

Step 1: A-01 → Get pharmacy revenue for Oct, Nov, Dec (fact table)
Step 2: A-08 → Multi-period comparison table
Step 3: D-02 entity extraction for "pharmacy" department
Step 4: D-01 Search → Any vendor invoice anomalies in Dec?
Step 5: A-13 Anomaly detection → flag outlier transactions in pharmacy ledger Dec
Step 6: B-12 Expense categorisation → check if cost spikes in pharmacy COGS
Step 7: Synthesise → "Revenue dropped 18% in December partly due to 2 major supply invoices (total AED 230K) which were prepaid in December. Underlying patient revenue grew 5%."
```

**DAG orchestration:**
- Uses Google ADK `sequential_agent` + `parallel_agent` patterns
- A2A Protocol carries state between steps
- Max 10 steps per chain (configurable, with LLM cost guard)
- Context window managed: each step receives summarised prior-step output

**Input:** `{ "question": string, "locale": string, "max_steps": int }`  
**Output:** `{ "answer": string, "steps_taken": [...], "charts": [...], "confidence": float }`

---

### D-04 — Autonomous Analyst (Spotter) Agent

**Purpose:** Continuously monitor the data warehouse for interesting patterns, anomalies, and threshold crossings — without being prompted. Surface findings as insight cards.

**Monitoring triggers:**
| Trigger type | Example |
|---|---|
| Threshold crossing | Revenue drops > 10% week-on-week |
| Statistical anomaly | Z-score > 3 on any metric |
| Trend reversal | 3-period declining trend detected |
| Correlation break | Visits up but revenue flat (previously correlated) |
| New record high/low | Highest-ever AP balance |

**Spotter execution cadence:** Every 15 minutes (configurable per metric)

**Insight card format:**
```
┌──────────────────────────────────────────────────────────┐
│ 🔍 INSIGHT  |  Finance  |  High Confidence (91%)          │
│                                                          │
│ Accounts Payable balance reached AED 1.8M — 34% above   │
│ 90-day average. 3 large invoices pending approval.       │
│                                                          │
│ [View Details]  [Investigate →]  [Dismiss]               │
└──────────────────────────────────────────────────────────┘
```

**Escalation:** If high-confidence insight sits undismissed for 24 hours → email alert to Finance Manager

---

### D-05 — Deep Research Agent

**Purpose:** Generate comprehensive analytical reports with statistical depth, multi-source synthesis, and cited evidence — like a financial analyst's research note.

**Pipeline:**
1. Research brief from user (scoped question + date range + focus area)
2. D-03 Multi-Step orchestrator: gather all relevant data slices
3. Statistical analysis: regression, correlation matrix, forecasting
4. Cross-source synthesis: combine HIMS + Tally findings
5. Insight generation: root-cause hypotheses ranked by evidence strength
6. Report generation: structured Markdown + charts; exported as PDF

**Report structure:**
```
1. Executive Summary (3 bullets)
2. Key Findings (with charts)
3. Statistical Analysis (tables, regression results)
4. Contributing Factors
5. Trend Outlook (ARIMA forecast)
6. Recommendations
7. Data Sources & Methodology
```

**Output format:** Markdown (rendered in UI) + downloadable PDF (WeasyPrint, EN/AR)  
**Generation time target:** ≤ 60 seconds for standard research report

---

### D-10 — Insight-to-Action Agent

**Purpose:** Convert analytical insights into concrete, executable actions — but NEVER execute without explicit human approval.

**Action types supported:**
| Action type | Description | HITL required |
|---|---|---|
| Schedule new report | Create a scheduled report based on insight | Yes |
| Create KPI alert | Add new KPI monitoring rule | Yes |
| Flag transaction for review | Move to AI Accountant review queue | Yes |
| Escalate anomaly | Notify Finance Manager by email | Yes |
| Update forecast parameters | Adjust A-12 model parameters | Yes |
| Trigger Tally sync review | Flag pending transactions for sync | Yes — mandatory |

**HITL approval flow:**
```
D-10 generates action recommendation
      │
      ▼
Pending Actions queue (app.pending_actions)
      │
      ▼ Notify user (in-app notification + email)
      │
      ▼
User reviews: action details + AI reasoning + data evidence
      │
      ├─ Approve → Execute action → Audit log entry
      └─ Reject → Archived with reason → D-10 learning feedback
```

**OPA policy `d10.action_gate`:** Any D-10 action with `writes_data: true` must have `approver_id != null AND approved_at != null`.

---

### D-13 — Scheduled Monitoring Agent

**Purpose:** Run a configurable battery of monitoring checks on a schedule (not triggered by user queries) and produce a daily intelligence briefing.

**Default monitoring schedule:**
| Job name | Frequency | Checks |
|---|---|---|
| KPI Health Check | Every 1 hour | All metrics vs 7-day moving average |
| Anomaly Scan | Every 15 minutes | Z-score scan all fact tables |
| AP/AR Aging | Daily 07:00 | Outstanding items age buckets |
| Budget Variance Alert | Daily 09:00 | Actual vs budget, flag > 10% variance |
| Cashflow Forecast | Daily 08:00 | Rolling 30-day cashflow forecast |
| Tally Sync Status | Every 30 minutes | Check sync queue stale or failed items |
| Document Queue Health | Every 1 hour | Review queue items > 48 hours old |
| Tax Compliance Check | Weekly (Monday) | GST/VAT liabilities due within 30 days |

**Daily Intelligence Briefing:**
- Generated at 07:30 each morning
- Emailed to Finance Manager + Clinic Administrator
- Contains: overnight anomalies, KPI delta summary, pending approvals count, outstanding items

---

## 5. HITL Action Approval UI

```
Pending Actions
──────────────────────────────────────────────────────
[2] High Priority    [4] Medium Priority    [1] Completed

────────────────────────────────────────────────────────────
⚠️  HIGH  |  Flag 3 large invoices for expedited review
────────────────────────────────────────────────────────────
Recommended by D-10 Insight-to-Action Agent
Reasoning: AP balance 34% above 90-day average; 3 invoices
(total AED 890K) pending > 7 days. Risk: payment delays,
vendor relationship impact.

Invoices: INV-2024-0891 (AED 320K), INV-2024-0892 (AED 280K),
INV-2024-0895 (AED 290K)

[ APPROVE ACTION ]  [ REJECT ]  [ View Evidence ]
────────────────────────────────────────────────────────────
```

---

## 6. Multi-Step Orchestration Architecture

```go
// ADK orchestration pattern (Go pseudo-code)
type ResearchOrchestrator struct {
    steps  []AgentStep
    state  map[string]interface{}
    a2a    *A2AClient
}

type AgentStep struct {
    AgentID     string              // "A-01", "D-02", etc.
    InputFn     func(state) Input
    OutputKey   string
    Condition   func(state) bool    // conditional execution
    Parallel    []AgentStep         // run these in parallel
}

func (o *ResearchOrchestrator) Execute(ctx context.Context, query string) (*ResearchResult, error) {
    for _, step := range o.buildDAG(query) {
        if step.Condition != nil && !step.Condition(o.state) {
            continue
        }
        if len(step.Parallel) > 0 {
            results := o.a2a.RunParallel(ctx, step.Parallel, o.state)
            o.state[step.OutputKey] = results
        } else {
            result, err := o.a2a.Call(ctx, step.AgentID, step.InputFn(o.state))
            o.state[step.OutputKey] = result
        }
    }
    return o.synthesise(o.state), nil
}
```

---

## 7. Database Schema Additions

```sql
-- Spotter insights surface
CREATE TABLE app.spotter_insights (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agent_id TEXT NOT NULL DEFAULT 'D-04',
  insight_type TEXT NOT NULL,          -- 'anomaly', 'threshold', 'trend_reversal', 'correlation_break', 'record'
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  metric_name TEXT,
  observed_value NUMERIC,
  expected_value NUMERIC,
  confidence NUMERIC(5,2),
  severity TEXT DEFAULT 'medium',      -- 'low', 'medium', 'high', 'critical'
  status TEXT DEFAULT 'new',           -- 'new', 'viewed', 'investigating', 'dismissed', 'actioned'
  dismissed_by UUID,
  dismissed_at TIMESTAMPTZ,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Pending actions (D-10 HITL)
CREATE TABLE app.pending_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  insight_id UUID REFERENCES app.spotter_insights(id),
  action_type TEXT NOT NULL,
  action_payload JSONB NOT NULL,
  reasoning TEXT,
  evidence JSONB,
  status TEXT DEFAULT 'pending',       -- 'pending', 'approved', 'rejected', 'executed', 'expired'
  approver_id UUID REFERENCES app.users(id),
  approved_at TIMESTAMPTZ,
  rejection_reason TEXT,
  executed_at TIMESTAMPTZ,
  execution_result JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Monitoring jobs
CREATE TABLE app.monitoring_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_name TEXT UNIQUE NOT NULL,
  cron_expression TEXT NOT NULL,
  agent_id TEXT NOT NULL DEFAULT 'D-13',
  config JSONB,
  last_run_at TIMESTAMPTZ,
  last_run_status TEXT,
  next_run_at TIMESTAMPTZ,
  enabled BOOLEAN DEFAULT TRUE
);
```

---

## 8. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/research` | Kick off multi-step analysis (D-03/D-05) |
| `GET` | `/v1/research/{id}` | Poll research job status + result |
| `GET` | `/v1/insights` | List Spotter insight cards |
| `PATCH` | `/v1/insights/{id}/dismiss` | Dismiss insight card |
| `GET` | `/v1/actions/pending` | List pending HITL actions |
| `POST` | `/v1/actions/{id}/approve` | Approve HITL action |
| `POST` | `/v1/actions/{id}/reject` | Reject HITL action |
| `GET` | `/v1/monitoring/jobs` | List monitoring jobs + status |
| `PATCH` | `/v1/monitoring/jobs/{id}` | Enable/disable/configure monitoring job |

---

## 9. Testing Requirements

| Test | Target |
|---|---|
| D-03 multi-step chain | 5-step chain completes correctly on 20 test scenarios |
| D-04 Spotter latency | Insight fires within 60s of threshold breach in test |
| D-05 research report | Generated report passes fact-check against source data on 10 test topics |
| D-10 HITL gate | 0 actions execute without approver sign-off (OPA enforced) |
| D-13 monitoring jobs | All 8 jobs run on schedule; alerting fires on injected anomaly |
| Orchestrator fault tolerance | Agent step failure → graceful error; partial result returned |
| Load test | 10 concurrent research jobs; P95 job completion < 90s |

---

## 10. Phase Exit Criteria

- [ ] D-03 Multi-Step Analysis: solves 5-step analytical questions autonomously
- [ ] D-04 Spotter: proactive insights surface in dashboard; latency within 60s
- [ ] D-05 Deep Research: generates full statistical research report with PDF export
- [ ] D-10 Insight-to-Action: HITL approval workflow operational; 0 actions without approval
- [ ] D-13 Scheduled Monitoring: 8 monitoring jobs active; daily briefing email delivered
- [ ] HITL Action Approval UI deployed
- [ ] Multi-step orchestration framework stable
- [ ] OPA `d10.action_gate` policy enforced
- [ ] Milestone M10 signed off

---

*Phase 14 | Version 1.0 | February 19, 2026*
