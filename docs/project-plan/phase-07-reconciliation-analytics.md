# Phase 07 — Reconciliation & Financial Analytics

**Phase Duration:** Weeks 23–26 (4 weeks)  
**Module(s):** Module B (AI Accountant), Module A (A-12, A-13)  
**Status:** Planning  
**Milestone:** M6 — Reconciliation Suite  
**Depends On:** Phase 06 complete (Tally sync operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.2](../ARCHITECTURE.md)

---

## 1. Objectives

Complete the AI Accountant module with bank reconciliation, outstanding receivables/payables analysis, cash flow forecasting, tax compliance tracking, and expense categorisation. Also deploy the advanced analytics agents: trend forecasting (A-12) and anomaly detection (A-13) for the BI module.

---

## 2. Scope

### In Scope
- B-10 Bank Reconciliation Agent
- B-11 Outstanding Items Agent
- B-12 Expense Categorisation Agent
- B-13 Tax Compliance Agent
- B-15 Cash Flow Forecasting Agent
- A-12 Trend Forecasting Agent
- A-13 Anomaly Detection Agent
- Reconciliation dashboard UI
- Outstanding payables/receivables aging UI
- Tax compliance report UI
- Cash flow forecast visualisation

### Out of Scope
- Easy Reports module (Phase 8+)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | B-10 Bank Reconciliation Agent | AI Engineer | Matches ≥ 80% of bank statement rows to Tally entries automatically; unmatched items in HITL queue |
| D-02 | B-11 Outstanding Items Agent | AI Engineer | Generates correct aging buckets (0–7, 8–30, 31–90, 90+ days) for payables + receivables |
| D-03 | B-12 Expense Categorisation Agent | AI Engineer | Auto-categorises ≥ 85% of expenses correctly; learning feedback loop |
| D-04 | B-13 Tax Compliance Agent | AI Engineer | Correct GST/VAT input credit, output tax, net liability calculations; reconciliation report matches manual |
| D-05 | B-15 Cash Flow Forecasting Agent | AI Engineer | 30-day cash flow projection with payables/receivables inputs; what-if scenarios supported |
| D-06 | A-12 Trend Forecasting Agent | AI Engineer | Time-series forecast extends historical lines; accuracy acceptable on 6-month test dataset (MAPE ≤ 15%) |
| D-07 | A-13 Anomaly Detection Agent | AI Engineer | Detects and alerts on 3 synthetic anomaly types injected into test dataset |
| D-08 | Reconciliation Dashboard | Frontend Engineer | Side-by-side bank vs Tally entries; matched/unmatched counts; difference calculation |
| D-09 | Outstanding Items UI | Frontend Engineer | Aging table (0–7, 8–30, 31–90, 90+ days) for payables and receivables; drill-down to individual invoices |
| D-10 | Cash Flow Forecast UI | Frontend Engineer | 30-day projection chart with inflows/outflows; scenario slider controls |
| D-11 | Tax Compliance UI | Frontend Engineer | GST period selector; input credit vs output tax table; net liability; export to PDF |

---

## 4. AI Agents Deployed

### B-10 Bank Reconciliation Agent

**Trigger:** User uploads bank statement (processed by B-02) + triggers reconciliation  
**Matching algorithm:**

```
For each bank statement row:
    │
    ▼ Step 1: Exact match
    │   Match criteria: amount exact + date exact + Tally narration contains bank description keyword
    │   Score: 100% → auto-match
    │
    ▼ Step 2: Near match
    │   Match criteria: amount exact + date within ±3 days
    │   Score: 85–99% → suggested match (accountant confirms)
    │
    ▼ Step 3: Fuzzy match
    │   Match criteria: amount within ±2% + date within ±7 days + description similarity > 70%
    │   Score: 60–84% → flagged match (accountant reviews)
    │
    ▼ Step 4: No match
    │   → Outstanding item; goes to B-11 outstanding items report
    │   → Accountant can: create new Tally entry, mark as excluded, or defer
```

**Output:** Reconciliation statement with:
- Matched pairs (auto + confirmed)
- Unmatched bank items (outstanding receipts)
- Unmatched Tally items (outstanding payments)
- Net reconciliation difference

**HITL:** All unmatched items and low-confidence matches require accountant review.

### B-11 Outstanding Items Agent

**Type:** Reactive (L1) — runs on demand or post-reconciliation  
**Input:** Tally ledger entries + bank reconciliation output  
**Calculates:**

| Bucket | Days Outstanding |
|---|---|
| Current | 0–7 days |
| Short-term | 8–30 days |
| Medium-term | 31–90 days |
| Long-term | 91+ days |

**Outputs:** 
- Outstanding Payments Report (checks/invoices not yet cleared)
- Outstanding Receipts Report (sales invoices unpaid)
- DSO (Days Sales Outstanding) metric
- DPO (Days Payable Outstanding) metric

### B-12 Expense Categorisation Agent

**Categories:**
| Category | Description | Keywords/Patterns |
|---|---|---|
| Office Supplies | Stationery, printing | "stationery", "supplies", "ink", "paper" |
| Utilities | Electricity, water, internet | "electricity", "water", "internet", "telecom" |
| Rent & Lease | Office/clinic rent | "rent", "lease", "property" |
| Medical Supplies | Drugs, equipment | "medical", "pharma", "drugs", "equipment" |
| Professional Fees | Consultants, lawyers | "legal", "audit", "consulting" |
| Travel | Transport, accommodation | "travel", "fuel", "hotel", "flight" |
| Payroll | Salaries, wages | "salary", "wages", "payroll" |
| Maintenance | Repairs, upkeep | "repair", "maintenance", "service" |
| Marketing | Ads, promotions | "advertising", "marketing", "promo" |
| Other | Uncategorised | default |

**Learning:** User corrections stored in `mapping_corrections` table (same as B-05); future runs for same vendor/description re-use correction.

### B-13 Tax Compliance Agent

**Scope:** GST/VAT compliance (configurable tax regime)

**Calculations:**
- **Input Tax Credit (ITC):** Sum of GST paid on purchases in period
- **Output Tax:** Sum of GST collected on sales in period
- **Net Liability:** Output Tax − ITC (if positive: tax owed; if negative: refund due)
- **Tax Reconciliation:** Compare ITC claimed vs Tally purchase entries

**Reports generated:**
- GSTR-1 summary (outward supplies)
- GSTR-3B summary (net liability)
- Input Credit register
- Tax liability register by period

**Compliance checklist:** Filing deadlines displayed; overdue periods flagged red.

### B-15 Cash Flow Forecasting Agent

**Type:** Reactive (L3) — complex financial modelling  
**Inputs:**
- Outstanding payables (from B-11) with expected payment dates
- Outstanding receivables (from B-11) with expected collection dates
- Historical payment patterns (average days to pay/collect from Tally history)
- Scheduled bank payments (if known)

**Model:**
```
For each day in forecast horizon (30 days):
    Projected Inflows  = expected receipts due + probability-weighted overdue collections
    Projected Outflows = scheduled payments due + probability-weighted upcoming bills
    Net Daily Cash Flow = Inflows - Outflows
    Running Cash Position = Previous Day Balance + Net Daily Flow
```

**What-if scenarios:** User can adjust:
- "What if all receivables over 30 days are collected this week?"
- "What if we delay payables by 2 weeks?"

**Output:** 30-day cash flow forecast chart (line chart with confidence bands) + scenario comparison

### A-12 Trend Forecasting Agent

**Type:** Reactive (L3) — time-series forecasting  
**Trigger:** User adds forecasting intent to BI query ("show revenue trend for next 6 months")  
**Models supported:**
- ARIMA / SARIMA (seasonal decomposition)
- Prophet (Facebook Prophet via Python microservice)
- LLM-based trend extrapolation (simpler, faster, lower accuracy)

**Output:** Extended chart with:
- Historical line (solid)
- Forecast line (dashed)
- Confidence interval band
- Annotation: model used, confidence, trend direction

**Accuracy target:** MAPE ≤ 15% on 6-month hold-out test set for clinic revenue metric.

### A-13 Anomaly Detection Agent

**Type:** Proactive (L3) — scheduled background scan  
**Schedule:** Runs daily at 2 AM on all monitored metrics  
**Detection methods:**
- Z-score (> 3σ from rolling mean on a metric)
- Interquartile range (IQR) outlier detection
- Seasonal-trend decomposition residuals

**Alert output:**
```
"⚠️ Anomaly detected: Pharmacy revenue dropped 42% on 15 Feb 2026, 
which is 3.7 standard deviations below the 30-day average (₹12,450 vs avg ₹21,300). 
Possible causes: ETL sync issue OR actual sales drop. Recommended: investigate."
```

**Channels:** In-app notification + email (KPI alert channel E-06)

---

## 5. Database Schema Additions

```sql
CREATE TABLE reconciliation_sessions (
    session_id      UUID PRIMARY KEY,
    entity_id       UUID REFERENCES tally_entities(entity_id),
    bank_doc_id     UUID REFERENCES documents(doc_id),
    period_start    DATE,
    period_end      DATE,
    matched_count   INTEGER,
    unmatched_bank  INTEGER,
    unmatched_tally INTEGER,
    net_difference  NUMERIC(15,2),
    status          VARCHAR(20) DEFAULT 'in_progress',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE reconciliation_matches (
    match_id        UUID PRIMARY KEY,
    session_id      UUID REFERENCES reconciliation_sessions(session_id),
    bank_row_hash   VARCHAR(64),
    tally_voucher_id VARCHAR(255),
    match_score     NUMERIC(5,2),
    match_type      VARCHAR(20), -- 'exact' | 'near' | 'fuzzy' | 'manual'
    status          VARCHAR(20), -- 'auto' | 'confirmed' | 'rejected'
    confirmed_by    UUID REFERENCES users(user_id)
);

CREATE TABLE cash_flow_forecasts (
    forecast_id     UUID PRIMARY KEY,
    entity_id       UUID REFERENCES tally_entities(entity_id),
    forecast_date   DATE,
    projected_inflow NUMERIC(15,2),
    projected_outflow NUMERIC(15,2),
    net_daily_cash  NUMERIC(15,2),
    running_balance NUMERIC(15,2),
    confidence_low  NUMERIC(15,2),
    confidence_high NUMERIC(15,2),
    model_used      VARCHAR(50),
    generated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE anomaly_alerts (
    alert_id        UUID PRIMARY KEY,
    metric_name     VARCHAR(100),
    alert_date      DATE,
    observed_value  NUMERIC(15,2),
    expected_range_low NUMERIC(15,2),
    expected_range_high NUMERIC(15,2),
    z_score         NUMERIC(5,2),
    explanation     TEXT,
    status          VARCHAR(20) DEFAULT 'open', -- 'open' | 'acknowledged' | 'resolved'
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 6. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| B-10 reconciliation | 50-transaction bank statement vs Tally ledger | ≥ 80% auto-matched; 0 false positives on exact matches |
| B-11 aging buckets | 100 transactions with varying ages | Correct bucket assignment for all 4 buckets |
| B-12 expense categorisation | 50 expenses across all categories | ≥ 85% correctly categorised |
| B-13 GST calculation | 3 test periods (known correct values from manual calculation) | 0 calculation errors; matches manual ≥ 99.9% |
| B-15 cash flow forecast | 30-day forecast vs actual (retrospective test) | MAPE ≤ 20% on retrospective test |
| A-12 trend forecast | 6-month hold-out test on clinic revenue | MAPE ≤ 15% |
| A-13 anomaly detection | 3 injected anomalies into test dataset | All 3 detected; 0 false positives on clean data |

---

## 7. Phase Exit Criteria

- [ ] Bank reconciliation agent matching ≥ 80% of transactions automatically
- [ ] Outstanding Items Agent producing correct aging reports with DSO/DPO metrics
- [ ] Expense Categorisation at ≥ 85% accuracy
- [ ] GST/VAT compliance reports generating correctly
- [ ] Cash flow forecast chart rendering with 30-day projection
- [ ] Trend Forecasting (A-12) deployed; MAPE target met
- [ ] Anomaly Detection (A-13) detecting injected anomalies
- [ ] All 7 agents passing integration tests
- [ ] Phase gate reviewed, M6 milestone declared, stakeholder demo completed

---

*Phase 07 | Version 1.0 | February 19, 2026*
