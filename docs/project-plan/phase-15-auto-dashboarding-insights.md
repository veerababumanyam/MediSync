# Phase 15 — Auto-Dashboarding & Prescriptive AI

**Phase Duration:** Weeks 47–48 (2 weeks)  
**Module(s):** D — Advanced Search Analytics  
**Status:** Planning  
**Milestone:** M11 — Full Module D analytics intelligence live  
**Depends On:** Phase 14 (Autonomous agents, Spotter, HITL action framework operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Complete Module D's intelligence layer with four final capabilities: AI-generated dashboards on demand, a proactive insight-discovery engine that identifies opportunities before users ask, an AI prescriptive engine that recommends concrete actions with quantified outcomes, and voice/mobile-first search for clinical staff on the move.

---

## 2. Scope

### In Scope
- D-06: Dashboard Auto-Generation Agent
- D-07: Insight Discovery Agent
- D-08: Prescriptive Recommendations Agent
- D-14: Voice & Mobile Search Agent
- AI-generated dashboards from natural language (via chat or Spotter)
- Prescriptive action cards with ROI estimates
- Voice search (Flutter mobile)
- Proactive opportunity surfacing (not just anomaly detection)
- Integration with D-04 Spotter (discovery → prescription pipeline)

### Out of Scope
- D-11, D-12 (developer tools — Phase 16)
- Module D governance layer (Phase 17)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | D-06 Dashboard Auto-Generation Agent | AI + Frontend Eng | NL request → saved dashboard in ≤ 15s; user can edit afterwards |
| D-02 | D-07 Insight Discovery Agent | AI Eng + Data Eng | Surfaces 3+ non-obvious opportunities/risks per business day |
| D-03 | D-08 Prescriptive Recommendations Agent | AI Eng | Recommendations include: action, expected outcome, evidence, confidence |
| D-04 | D-14 Voice & Mobile Search | AI Eng + Frontend Eng | Voice query on Flutter → search result in < 5s; EN + AR |
| D-05 | Prescriptive panel in dashboard | Frontend Eng | "AI Recommendations" panel on main dashboard; cards with approve/dismiss |
| D-06 | Dashboard library | Frontend Eng | All AI-generated dashboards saved to user library; shareable |
| D-07 | Opportunity finder UI | Frontend Eng | "Opportunities & Risks" section in Spotter panel |

---

## 4. AI Agents Deployed

### D-06 — Dashboard Auto-Generation Agent

**Purpose:** Generate a fully configured, visually coherent dashboard from a natural language description — no drag-and-drop required.

**Trigger modes:**
1. **Chat:** "Create a dashboard showing pharmacy performance for Q1"
2. **Spotter:** "Generate a dashboard tracking this anomaly context"
3. **Scheduled:** Automatically generate a recommended dashboard for a new user based on their role

**Generation pipeline:**
```
NL description
      │ D-09 semantic model lookup
      ▼
Metric selection (top N metrics relevant to request)
      │ A-03 Visualisation Routing
      ▼
Chart type assignment per metric
      │ A-11 Pin-to-Dashboard API
      ▼
Layout engine (auto-arrange: KPI row + 2-col charts + table)
      │
      ▼
Dashboard saved to user library
      │
      ▼
User can refine via natural language: "Move revenue to top left"
```

**Layout algorithm:**
- Row 1: 4 KPI cards (key metrics)
- Row 2–3: 2-column chart grid (bar/line/area)
- Row 4: data table (drill-down summary)

**Output:** Saved dashboard with `dashboard_id`; opens in dashboard view; all panels editable

---

### D-07 — Insight Discovery Agent

**Purpose:** Continuously scan for non-obvious opportunities and risk patterns — beyond threshold alerts. Think of this as the proactive business intelligence layer.

**Discovery categories:**
| Category | Example discovery |
|---|---|
| Cost optimisation | "Lab reagent costs up 22% while test volumes flat — possible waste or pricing issue" |
| Revenue opportunity | "OPD appointments up 30% on Thu-Fri — capacity may be underutilised Mon-Wed" |
| Collection risk | "30% of pharmacy revenue from 3 accounts with payment history > 45 days" |
| Efficiency gain | "Average document approval time 4.2 days — 8 invoices awaiting approval > 7 days" |
| Supplier risk | "Top vendor (40% of AP) has had 3 price increases in 6 months" |
| Clinical-financial correlation | "High LOS patients in ICU averaging 2.4× revenue contribution — consider dedicated billing track" |

**Discovery cadence:** Hourly scan; new discoveries deduplicated (≤ 5 net-new cards per day to avoid insight fatigue)  
**Confidence scoring:** Each discovery includes evidence chain + confidence %  
**Lifecycle:** new → viewed → investigating → actioned/dismissed

---

### D-08 — Prescriptive Recommendations Agent

**Purpose:** Take insights from D-04 (anomalies) and D-07 (opportunities) and generate concrete, quantified recommendations with expected outcomes.

**Recommendation structure:**
```
┌─────────────────────────────────────────────────────────────┐
│  💡 RECOMMENDATION  |  Finance  |  High Impact             │
│                                                             │
│  Issue: AP balance AED 1.8M — 34% above 90-day avg.        │
│                                                             │
│  Recommended Action:                                        │
│  Expedite approval for INV-0891, 0892, 0895               │
│  (total AED 890K — all vendor Gulf Medical)                 │
│                                                             │
│  Expected Outcome:                                          │
│  • Reduce AP balance by AED 890K (49%) within 3 days       │
│  • Restore vendor relationship — overdue 7 days            │
│  • Avoid 1.5% late payment penalty = AED 13,350 saving     │
│                                                             │
│  Confidence: 87% ██████████░  Based on 18 similar past     │
│  scenarios                                                  │
│                                                             │
│  [ TAKE ACTION ]  [ DISMISS ]  [ Learn More ]              │
└─────────────────────────────────────────────────────────────┘
```

**Quantification engine:** Uses historical data to estimate outcome ranges; expresses as "Expected: AED X saving / risk reduction"  
**HITL gate:** "TAKE ACTION" routes to D-10 Insight-to-Action pending actions queue; never executes directly  
**Learning:** Tracks whether actioned recommendations produced the expected outcome; feeds accuracy improvement

---

### D-14 — Voice & Mobile Search Agent

**Purpose:** Enable clinical staff and managers on mobile devices to search and query using voice input.

**Platform:** Flutter mobile (iOS + Android)  
**Languages:** English + Arabic voice recognition

**Voice search pipeline:**
```
Voice input (Flutter mic)
      │ Device STT (WhisperKit on iOS / on-device model on Android)
      ▼
Transcription text
      │ E-01 Language Detection
      ▼
D-01 NL Search (or A-01 BI chat if query type is analytic)
      │
      ▼
Result rendered as:
  • Voice summary (TTS) — 1–2 sentence answer
  • Visual card (search result or KPI card)
      │
      ▼
User can tap to expand to full result panel
```

**Voice query types:**
| Query type | Example | Route |
|---|---|---|
| Search | "Find invoices from Al Noor last month" | D-01 |
| Metric lookup | "What is pharmacy revenue today?" | D-09 + A-01 |
| Trend question | "Is our revenue growing?" | A-12 + D-08 |
| Action | "Show me pending approvals" | Module B approval queue |

**Arabic voice notes:**
- STT: On-device model with Arabic medical/financial vocabulary fine-tune
- TTS: Arabic text-to-speech with correct diacritics
- Numbers: Arabic-Indic numerals rendered in TTS

**Performance target:** Voice to result (including TTS) < 5 seconds on mobile (Wi-Fi)

---

## 5. AI-Generated Dashboard Examples

**"Pharmacy Performance Q1":**
```
Row 1: [Pharmacy Revenue AED 2.4M ↑12%] [Gross Margin 34% ↓2pp] [Rx Count 8,400 ↑5%] [Cost per Rx AED 286 ↑1%]
Row 2: [Monthly Pharmacy Revenue — Bar]        [Margin Trend — Line]
Row 3: [Top 10 Drugs by Revenue — Horizontal bar]  [Supplier Cost Breakdown — Pie]
Row 4: [Raw transaction table — sortable]
```

**Dashboard auto-title generation:** LLM generates title + subtitle from metric set; e.g., "Pharmacy Performance — Q1 2026: Revenue growing but margin pressure."

---

## 6. Prescriptive Learning Feedback Loop

```
Recommendation given → User approves action
                             │
                    Action executed
                             │
                    7-day outcome measured (actual vs predicted)
                             │
                    Stored in app.recommendation_outcomes
                             │
                    D-08 fine-tune: if recommendation was off by > 20%
                    → downweight similar future recommendations
```

---

## 7. Database Schema Additions

```sql
-- Insight Discovery catalogue
CREATE TABLE app.discovery_insights (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agent_id TEXT DEFAULT 'D-07',
  category TEXT NOT NULL,       -- 'cost_optimisation', 'revenue_opportunity', 'collection_risk', etc.
  title TEXT NOT NULL,
  summary TEXT NOT NULL,
  evidence JSONB NOT NULL,
  impact_estimate JSONB,        -- { "type": "saving", "min": 5000, "max": 20000, "currency": "AED" }
  confidence NUMERIC(5,2),
  status TEXT DEFAULT 'new',
  dedup_hash TEXT UNIQUE,       -- prevent duplicate discoveries
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Prescriptive recommendations
CREATE TABLE app.recommendations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  insight_id UUID REFERENCES app.discovery_insights(id),
  agent_id TEXT DEFAULT 'D-08',
  recommended_action TEXT NOT NULL,
  expected_outcome JSONB,
  confidence NUMERIC(5,2),
  status TEXT DEFAULT 'pending',    -- 'pending', 'actioned', 'dismissed', 'expired'
  outcome_actual JSONB,
  outcome_measured_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Generated dashboards
CREATE TABLE app.generated_dashboards (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES app.users(id),
  agent_id TEXT DEFAULT 'D-06',
  title TEXT NOT NULL,
  subtitle TEXT,
  prompt_used TEXT,
  layout_config JSONB NOT NULL,
  status TEXT DEFAULT 'active',
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 8. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/dashboards/generate` | Generate dashboard from NL description (D-06) |
| `GET` | `/v1/insights/discoveries` | List D-07 opportunity/risk discoveries |
| `GET` | `/v1/recommendations` | List D-08 prescriptive recommendations |
| `POST` | `/v1/recommendations/{id}/action` | Route recommendation to D-10 pending actions |
| `POST` | `/v1/search/voice` | Voice search (binary audio input → NL search) |
| `GET` | `/v1/dashboards/generated` | List user's AI-generated dashboards |

---

## 9. Testing Requirements

| Test | Target |
|---|---|
| D-06 dashboard generation | 10 different NL prompts → valid dashboard in ≤ 15s each |
| D-06 layout correctness | KPI row + chart grid + table present in all generated dashboards |
| D-07 discovery non-duplication | dedup hash prevents same discovery appearing twice within 24h |
| D-08 prescriptive accuracy | Expected outcome within 30% of actual on 10 historical test cases |
| D-14 voice EN | 20 voice queries; STT accuracy ≥ 95%; result correct ≥ 85% |
| D-14 voice AR | 10 Arabic voice queries; STT accuracy ≥ 90%; result correct ≥ 80% |
| Voice latency | Voice to rendered result < 5s (mobile Wi-Fi simulated) |

---

## 10. Phase Exit Criteria

- [ ] D-06: AI-generated dashboards from NL in ≤ 15s; dashboards saved to library
- [ ] D-07: Insight Discovery running hourly; ≥ 3 non-obvious discoveries per business day
- [ ] D-08: Prescriptive recommendations with quantified outcomes; HITL gate to D-10
- [ ] D-14: Voice search on Flutter (EN + AR); voice-to-result < 5s
- [ ] Prescriptive panel live on main dashboard
- [ ] AI-generated dashboard library accessible to all users
- [ ] D-04 → D-07 → D-08 pipeline connected end-to-end
- [ ] Milestone M11 — Full Module D analytics intelligence signed off

---

*Phase 15 | Version 1.0 | February 19, 2026*
