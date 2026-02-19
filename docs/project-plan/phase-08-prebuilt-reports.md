# Phase 08 — Pre-Built Reports & Dashboards (Easy Reports — Part 1)

**Phase Duration:** Weeks 27–30 (4 weeks)  
**Module(s):** Module C (Easy Reports)  
**Status:** Planning  
**Milestone:** M7 — Reports Module (partial)  
**Depends On:** Phase 07 complete (full AI Accountant suite live)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.3](../ARCHITECTURE.md) | [PRD.md §6.8](../PRD.md)

---

## 1. Objectives

Launch the Easy Reports module with a comprehensive pre-built report library covering all major financial and operational report types, consolidated multi-company dashboards, budget vs. actual variance analysis, and inventory aging reports. Users can access professional-quality reports instantly without writing queries.

---

## 2. Scope

### In Scope
- C-01 Pre-Built Report Generator Agent (20+ report types)
- C-02 Multi-Company Consolidation Agent
- C-07 Budget vs. Actual Variance Agent
- C-08 Inventory Aging & Reorder Agent
- Pre-built report library UI
- Report collection: Sales & Debtors, Profitability, Financial Statements, Revenue, Inventory
- Interactive drill-down dashboards (Sales, Finance, Inventory, Operations)
- KPI scorecard components
- Multi-period report comparison
- Real-time data refresh for dashboards

### Out of Scope
- Scheduled email distribution (Phase 9)
- Zero-code report builder (Phase 10)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | C-01 Pre-Built Report Generator | AI Engineer | ≥ 20 report types; each renders correctly with live data; EN + AR locale |
| D-02 | C-02 Multi-Company Consolidation | AI Engineer | Consolidates 2+ Tally entities; eliminates inter-company transactions; correct group P&L |
| D-03 | C-07 Budget vs. Actual Variance | AI Engineer | Correct variance % and absolute delta; year-end forecast; department breakdown |
| D-04 | C-08 Inventory Aging & Reorder | AI Engineer | Correct slow-moving / obsolete identification; reorder suggestions with quantities |
| D-05 | Report Library UI | Frontend Engineer | Categorised list with search/filter; preview thumbnails; one-click run |
| D-06 | Report Viewer | Frontend Engineer | Full-page report with charts + tables; chart type switching; pagination |
| D-07 | Sales Dashboard | Frontend Engineer | Top customers, revenue trend, pipeline metrics; drill-down to invoice level |
| D-08 | Finance Dashboard | Frontend Engineer | P&L snapshot, cash flow, budget vs actual, expense breakdown; all real-time |
| D-09 | Inventory Dashboard | Frontend Engineer | Stock levels, aging, movement chart, valuation; low-stock alerts integrated |
| D-10 | KPI Scorecard Component | Frontend Engineer | Multi-metric cards with RAG (Red/Amber/Green) status; sparklines; YoY comparison |
| D-11 | Multi-Period Report UI | Frontend Engineer | Side-by-side period comparison; slider to select periods; delta annotation |

---

## 4. AI Agents Deployed

### C-01 Pre-Built Report Generator Agent

**Type:** Reactive (L1) — deterministic template execution  
**Mode:** Template-driven (not free-form SQL) for reliability and consistency

**Report Library (20+ reports):**

| Category | Report Name | Key Metrics |
|---|---|---|
| **Sales & Debtors** | Sales by Customer | Revenue, qty, YoY growth per customer |
| | Sales by Salesperson | Per-rep revenue, target achievement |
| | Debtor Aging | 0–30, 31–60, 61–90, 90+ days outstanding |
| | Top Customers | Revenue ranked, growth rate |
| | Territory-wise Sales | Region/branch revenue comparison |
| **Profitability** | Customer Profitability | Revenue − Cost per customer |
| | Product-wise Margin | Margin % per drug/service |
| | Department Profitability | Clinic vs Pharmacy contribution |
| | Salesperson Profitability | Revenue - COGS per rep |
| **Financial Statements** | Profit & Loss (P&L) | Revenue, COGS, gross margin, EBITDA, net profit |
| | Balance Sheet | Assets, liabilities, equity |
| | Cash Flow Statement | Operating / Investing / Financing activities |
| | Trial Balance | All ledger balances |
| | Ledger-wise Analysis | Voucher drill-down per ledger |
| **Targets & Budgets** | Sales Target vs Actual | Achievement % with variance |
| | Expense Budget vs Actual | Over/under budget by category |
| | Annual Budget Tracking | Monthly trend vs annual budget |
| **Revenue** | Cost-Centre Revenue | Revenue by cost centre |
| | Revenue by Type | Services / Products / Other Income |
| **Inventory** | Stock Aging Report | Slow-moving / obsolete stock |
| | Inventory Valuation | FIFO / Weighted Avg valuation |
| | Item Movement Report | Stock in / out by item |
| | Stock Shortage Alerts | Items below reorder level |

**Report execution flow:**
```
User selects report + date range + filters (department, entity, etc.)
    │
    ▼ C-01 loads pre-defined SQL template for report type
    │   Injects parameters (date range, filters, locale)
    │   Executes against medisync_readonly Postgres role
    │
    ▼ C-05 (Phase 10): Row/column security enforcement
    │   (In Phase 8: basic OPA role filter for department)
    │
    ▼ Render: HTML table + ECharts config
    │
    ▼ Available exports via C-03 (Phase 9): PDF | Excel | CSV
```

### C-02 Multi-Company Consolidation Agent

**Input:** List of Tally entity IDs to consolidate  
**Process:**
1. Load fact_vouchers for each entity's date range
2. Identify inter-company transactions (same ledger name exists in both companies → flag as intercompany)
3. Eliminate matched inter-company pairs
4. Aggregate remaining ledger entries into consolidated view

**Output:**
- Consolidated P&L
- Consolidated Balance Sheet
- Entity-by-entity comparison table
- Intercompany elimination summary

**Limitation Phase 8:** Automated intercompany elimination uses ledger name matching (simple). Phase 10 will refine with explicit IC ledger tags.

### C-07 Budget vs. Actual Variance Agent

**Input:** Budget data uploaded by user (Excel import) + actual data from Tally warehouse  
**Calculations:**
- Variance = Actual − Budget (absolute)
- Variance % = (Actual − Budget) / Budget × 100
- Year-end forecast = YTD Actual + (Remaining budget periods × YTD run rate)

**Alerts:** Any department with > 10% budget overrun for the month is flagged red in the dashboard.

**Schema addition:**
```sql
CREATE TABLE budget_entries (
    budget_id   UUID PRIMARY KEY,
    entity_id   UUID REFERENCES tally_entities(entity_id),
    ledger_id   UUID,
    period      DATE,        -- first day of month
    budget_amount NUMERIC(15,2),
    uploaded_by UUID REFERENCES users(user_id),
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);
```

### C-08 Inventory Aging & Reorder Agent

**Aging classification:**
| Classification | Criteria |
|---|---|
| Fast-moving | Sold ≥ 60% of stock in last 30 days |
| Slow-moving | Sold < 20% of stock in last 90 days |
| Near-obsolete | Zero movement in last 60 days |
| Dead stock | Zero movement in last 180 days + no open orders |

**Reorder suggestion logic:**
```
Reorder Qty = (Average Daily Sales × (Lead Time Days + Safety Stock Days)) − Closing Stock
Safety Stock = 1.5 × Standard Deviation of Daily Sales × √Lead Time
```

**HITL:** Reorder quantity suggestions are recommendations only. Pharmacy Manager confirms before any action.

---

## 5. Pre-Built Dashboard Specifications

### Finance Dashboard
- **P&L Snapshot KPI cards:** Total Revenue, Total Expenses, Gross Margin %, Net Profit
- **Revenue vs Expense trend:** Line chart, last 12 months
- **Budget vs Actual:** Grouped bar chart, current period
- **Cash Position:** KPI card with 30-day forecast line
- **Outstanding Receivables:** Donut chart by aging bucket

### Sales (Pharmacy) Dashboard
- **Revenue KPIs:** Today, Week, Month with MoM delta
- **Top 10 Drugs:** Horizontal bar chart by revenue
- **Sales Trend:** Line chart with comparison to last year
- **Debtor Aging:** Stacked bar by age bucket
- **Territory-wise Sales:** If location data available

### Inventory Dashboard
- **Stock Value:** KPI with trend
- **Low Stock Alert count:** Red badge
- **Inventory Aging Distribution:** Pie chart (fast/slow/near-obsolete/dead)
- **Top 10 Fast-Moving Items:** Horizontal bar
- **Items at Reorder Level:** Table with reorder suggestion column

### Operations Dashboard
- **Patient Visits (Today/Week/Month):** KPI cards
- **Doctor-wise Productivity:** Bar chart
- **Appointment Cancellation Rate:** KPI with trend
- **Avg Revenue per Patient:** KPI with MoM comparison
- **Pharmacy vs Clinic Revenue Split:** Donut chart

---

## 6. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| C-01 report accuracy | All 20+ reports run against known test dataset | Report figures match manually verified reference values |
| C-02 consolidation | 2-entity consolidation with 5 intercompany transactions | Correct elimination; consolidated P&L matches expected |
| C-07 variance | Budget vs actual for 3 departments | Variance calculations 100% accurate |
| C-08 aging | 50 inventory items of varying movement rates | Correct classification for all |
| Report EN + AR | All 20 reports in both locales | Arabic reports render RTL correctly |
| Dashboard load time | All 4 dashboards | Load in < 3 seconds with 1-year of data |
| Drill-down | Click through on Sales + Finance dashboards | Correct detail data returned |

---

## 7. Phase Exit Criteria

- [ ] All 20+ pre-built reports generating correctly with live data
- [ ] Multi-company consolidation tested with 2+ entities
- [ ] Budget vs Actual variance analysis working; budget upload mechanism functional
- [ ] Inventory aging correctly classifying stock; reorder suggestions reasonable
- [ ] All 4 pre-built dashboards (Finance, Sales, Inventory, Operations) rendering
- [ ] KPI scorecards with RAG status working
- [ ] Arabic locale reports rendering correctly (reviewer sign-off)
- [ ] Dashboard load time < 3 seconds
- [ ] Phase gate reviewed and signed off

---

*Phase 08 | Version 1.0 | February 19, 2026*
