# MediSync — AI Agent Backlog

**Version:** 1.0 | **Created:** February 19, 2026 | **Source:** [PRD.md](../PRD.md)

> **Goal:** Maximize automation across all three MediSync modules — Conversational BI, AI Accountant, and Easy Reports — by identifying every task that can be delegated to an AI agent, reducing custom coding by reusing proven open-source tooling.

> **Goal:** Maximize automation across all four MediSync modules — Conversational BI, AI Accountant, Easy Reports, and Advanced Search Analytics — by identifying every task that can be delegated to an AI agent, reducing custom coding by reusing proven open-source tooling.

---

## 1. Scope & Methodology

Each task in this backlog was extracted from the PRD's Functional Requirements (§6), User Stories (§11), and Phased Timeline (§12). Tasks are classified by:

- **Module** — which product module the task belongs to
- **Agent type** — reactive (triggered by user action), proactive (self-scheduled), or hybrid
- **Complexity** — L1 (single-step), L2 (multi-step), L3 (autonomous/long-horizon)
- **Human-in-the-Loop (HITL)** — whether a human must approve before action
- **Priority** — Phase number from PRD timeline (lower = earlier)

---

## 2. Master Agent Task Inventory

### 2.1 Module A — Conversational BI Dashboard

| # | Task | PRD Ref | Agent Type | Complexity | HITL | Priority |
|---|------|---------|-----------|-----------|------|----------|
| A-01 | **Text-to-SQL Translation** — Convert natural language query to secure, read-only SQL against the data warehouse | §5.1, US1–US8 | Reactive | L2 | No | Phase 2 |
| A-02 | **SQL Self-Correction** — Detect SQL execution errors, rewrite and retry the query automatically (up to N retries) | §5.3 | Reactive | L2 | No | Phase 2 |
| A-03 | **Visualization Routing** — Classify query intent (trend → line chart, comparison → bar, breakdown → pie, tabular → table) and emit the correct chart-type token | §5.2, §6.1 | Reactive | L1 | No | Phase 2 |
| A-04 | **Domain Terminology Normalization** — Map healthcare / accounting synonyms ("footfall" → patient_visits, "outstanding" → accounts_receivable) before SQL generation | §5.4 | Reactive | L1 | No | Phase 2 |
| A-05 | **Hallucination Guard** — Detect off-topic questions and redirect users back to business analytics | §10 (NFR) | Reactive | L1 | No | Phase 2 |
| A-06 | **Confidence Scoring** — Attach a confidence score (0–100%) to every generated answer; route low-confidence answers (<70%) to manual review queue | §10 (NFR) | Reactive | L1 | Yes (low conf.) | Phase 2 |
| A-07 | **Drill-Down Context Agent** — Parse a user "click event" on a chart element and generate the appropriate drill-down SQL query | §6.5 | Reactive | L2 | No | Phase 3 |
| A-08 | **Multi-Period Comparison Agent** — Automatically construct period-over-period (MoM, YoY) comparison queries from a single user request | §6.2, US4 | Reactive | L2 | No | Phase 3 |
| A-09 | **Report Scheduling Agent** — Generate scheduled reports, format output (PDF/Excel/CSV), and trigger email delivery | §6.2, US5 | Proactive | L2 | No | Phase 3 |
| A-10 | **KPI Alert Agent** — Continuously monitor key metrics against configurable thresholds; emit notifications via email/SMS/in-app when breached | §6.5, US28 | Proactive | L2 | No | Phase 3 |
| A-11 | **Chart-to-Dashboard Pin Agent** — Accept a "pin" action, persist chart config to the user's dashboard store, and schedule auto-refresh | §6.2 | Reactive | L1 | No | Phase 3 |
| A-12 | **Trend Forecasting Agent** — Apply time-series models (ARIMA, Prophet, or LLM-based) to extend historical trend lines into future periods | §6.6, §6.8.11 | Reactive | L3 | No | Phase 7 |
| A-13 | **Anomaly Detection Agent** — Scan all monitored metrics on a configured schedule; surface statistically significant outliers with plain-language explanations | §6.6, §6.9.4 | Proactive | L3 | No | Phase 7 |

---

### 2.2 Module B — AI Accountant

| # | Task | PRD Ref | Agent Type | Complexity | HITL | Priority |
|---|------|---------|-----------|-----------|------|----------|
| B-01 | **Document Classification Agent** — Auto-categorize uploaded files into Invoice, Bank Statement, Bill, Receipt, Tax Document, or Other | §6.7.1 | Reactive | L1 | No | Phase 4 |
| B-02 | **OCR Extraction Agent** — Extract structured fields (amount, vendor, date, invoice #, tax) from PDF/image/Excel files with confidence scoring; flag low-confidence extractions | §6.7.1, US9, US13 | Reactive | L2 | Yes (low conf.) | Phase 4 |
| B-03 | **Handwriting Recognition Agent** — Specialized sub-agent for handwritten invoices; applies pattern recognition on top of standard OCR | §6.7.1, US13 | Reactive | L2 | Yes | Phase 4 |
| B-04 | **Vendor Matching Agent** — Match extracted vendor names to existing Tally vendor master records; create new vendor records if unrecognised | §6.7.2 | Reactive | L2 | Yes (new vendor) | Phase 5 |
| B-05 | **Ledger Mapping Agent** — Match each transaction to the appropriate Tally GL ledger using contextual AI; provide a confidence score per mapping; learn from user corrections | §6.7.2, US10 | Reactive | L2 | Yes (review) | Phase 5 |
| B-06 | **Sub-Ledger & Cost Centre Assignment Agent** — Suggest sub-ledger and cost centre based on transaction context and historical patterns | §6.7.2 | Reactive | L1 | Yes | Phase 5 |
| B-07 | **Duplicate Invoice Detection Agent** — Compare incoming invoices against existing records (same amount + vendor + date window); flag suspicious duplicates before posting | §6.7.2, §6.7.5, US16 | Reactive | L2 | Yes | Phase 5 |
| B-08 | **Approval Workflow Agent** — Route transactions through configurable approval chains (Accountant → Manager → Finance → Posted); track status; send reminders for stale approvals | §6.7.5, US14 | Reactive/Proactive | L2 | Yes (always) | Phase 5 |
| B-09 | **Tally Sync Agent** — Push approved journal entries, purchase bills, and sales invoices directly into Tally via TDL/API; implement pre-sync validation and auto-retry | §6.7.3, US11 | Reactive | L2 | Yes (one-click) | Phase 6 |
| B-10 | **Bank Reconciliation Agent** — Match bank statement rows to Tally ledger entries by amount + date + description; assign confidence scores; suggest matches for unresolved items | §6.7.4, US12 | Reactive | L2 | Yes (unmatched) | Phase 7 |
| B-11 | **Outstanding Items Agent** — Generate outstanding-payments and outstanding-receipts reports; classify items by age bucket (0-7, 8-30, 30+ days) | §6.7.4 | Reactive | L1 | No | Phase 7 |
| B-12 | **Expense Categorisation Agent** — Auto-assign expense categories (utilities, travel, office supplies, etc.) based on vendor, description, and learned rules | §6.7.5 | Reactive | L1 | Yes | Phase 7 |
| B-13 | **Tax Compliance Agent** — Compute GST/VAT input-credit, output-tax, and net liability per period; generate compliance-ready reports | §6.7.6, US26 | Reactive | L2 | No | Phase 7 |
| B-14 | **Audit Trail Logger Agent** — Capture every mutation (who, what, when, source document) and write immutable audit log entries | §6.7.6, US14 | Reactive | No | No | Phase 6 |
| B-15 | **Cash Flow Forecasting Agent** — Project future cash position from scheduled payables/receivables; identify shortfalls; run what-if scenarios | §6.7.8 | Reactive | L3 | No | Phase 7 |
| B-16 | **Multi-Entity Tally Manager Agent** — Switch Tally company context per entity; sync and consolidate statements independently; aggregate cross-entity dashboards | §6.7.3, US15, US20 | Reactive | L2 | No | Phase 6 |

---

### 2.3 Module C — Easy Reports

| # | Task | PRD Ref | Agent Type | Complexity | HITL | Priority |
|---|------|---------|-----------|-----------|------|----------|
| C-01 | **Pre-Built Report Generator Agent** — Produce P&L, Balance Sheet, Cash Flow, Debtor Aging, and 20+ other report types from the data warehouse on demand | §6.8.1, US17, US22, US23 | Reactive | L1 | No | Phase 8 |
| C-02 | **Multi-Company Consolidation Agent** — Merge financial statements from multiple Tally instances; eliminate inter-company transactions; produce consolidated views | §6.8.2, US20 | Reactive | L2 | No | Phase 8 |
| C-03 | **Report Scheduling & Distribution Agent** — Create, manage, and execute report schedules; format outputs (PDF, Excel, HTML, PPT); email to distribution lists | §6.8.5, US19 | Proactive | L2 | No | Phase 9 |
| C-04 | **Custom Metric Formula Agent** — Interpret user-defined formula expressions and calculate custom KPIs without SQL coding; store as governed metrics | §6.8.3, §6.9.5 | Reactive | L2 | No | Phase 10 |
| C-05 | **Row/Column Security Enforcement Agent** — Apply RBAC filters to every query; mask columns; enforce row-level data scopes based on the requesting user's role | §6.8.8, §6.9.11 | Reactive | L1 | No | Phase 10 |
| C-06 | **Data Quality Validation Agent** — Run ETL-level checks (missing values, duplicates, referential integrity); produce data quality reports; alert on anomalies | §8, §6.8.6 | Proactive | L2 | No | Phase 1 |
| C-07 | **Budget vs. Actual Variance Agent** — Compare actuals to budget targets; compute variance; forecast year-end outturn; flag overages | §6.8.1, US23 | Reactive | L2 | No | Phase 8 |
| C-08 | **Inventory Aging & Reorder Agent** — Identify slow-moving / obsolete stock; compute turnover ratios; recommend reorder quantities and timing | §6.5, §6.8.1, US22 | Reactive | L2 | Yes (reorder) | Phase 8 |

---

### 2.4 Module D — Advanced Search-Driven Analytics

| # | Task | PRD Ref | Agent Type | Complexity | HITL | Priority |
|---|------|---------|-----------|-----------|------|----------|
| D-01 | **Natural Language Search Agent** — Accept Google-style queries; return ranked data results with auto-complete, spell-check, and related suggestions | §6.9.1 | Reactive | L2 | No | Phase 13 |
| D-02 | **Entity Recognition Agent** — Identify business entities (doctors, drugs, suppliers, locations) in unstructured queries and resolve to database records | §6.9.1 | Reactive | L1 | No | Phase 13 |
| D-03 | **Multi-Step Conversational Analysis Agent** — Decompose complex multi-part questions into sequential analytical sub-tasks; chain results into a cohesive answer | §6.9.1, US27,US31 | Reactive | L3 | No | Phase 14 |
| D-04 | **Autonomous AI Analyst (Spotter) Agent** — Run multi-step analytical workflows autonomously (retrieve → analyse → compare → forecast → recommend); surface proactive insights | §6.9.2, US28, US30 | Proactive | L3 | No | Phase 14 |
| D-05 | **Deep Research Agent** — Discover hidden patterns, correlations, and anomalies; run regression/clustering/decomposition; produce structured research reports | §6.9.2 | Proactive | L3 | No | Phase 14 |
| D-06 | **Dashboard Auto-Generation Agent** — From a search query or raw dataset, automatically design and render a complete, labelled dashboard with optimal chart types | §6.9.3, US29 | Reactive | L3 | No | Phase 15 |
| D-07 | **Insight Discovery & Prioritisation Agent** — Continuously scan all data; surface the top N most actionable insights ordered by business impact; embed contextual annotations in charts | §6.9.4 | Proactive | L3 | No | Phase 15 |
| D-08 | **Prescriptive Recommendations Agent** — Go beyond insights; generate specific, quantified action recommendations with root-cause analysis and expected impact | §6.9.4, US30 | Proactive | L3 | No | Phase 15 |
| D-09 | **Semantic Layer Management Agent** — Maintain the unified metric and dimension registry; resolve metric lineage queries; enforce version-controlled metric definitions | §6.9.5 | Reactive | L2 | Yes (governance) | Phase 13 |
| D-10 | **Insight-to-Action Workflow Agent** — Detect threshold violations and automatically trigger downstream business actions (e.g. low stock → create PO in Tally) | §6.9.6, US28 | Proactive | L3 | Yes (always) | Phase 14 |
| D-11 | **Code Generation Agent (SpotterCode)** — Translate natural language feature descriptions into correct React, Python, or SQL code; respect security standards and best-practices | §6.9.8, US33 | Reactive | L2 | No | Phase 16 |
| D-12 | **Federated Query Optimisation Agent** — Route cross-source queries optimally across Tally, HIMS, and external databases; return data with lineage metadata | §6.9.9 | Reactive | L2 | No | Phase 16 |
| D-13 | **Scheduled Autonomous Monitoring Agent** — Execute user-defined recurring analysis tasks on schedule; alert when conditions are met (e.g. revenue drop >10% MoM) | §6.9.2, US28 | Proactive | L2 | No | Phase 14 |
| D-14 | **Voice/Mobile Search Agent** — Accept voice input; transcribe; pass to NL search pipeline; return mobile-formatted results | §6.9.10, US31 | Reactive | L2 | No | Phase 15 |

---

### 2.5 Module E — Language & Localisation

> **Context:** All user queries, AI responses, reports, notifications, and UI strings must be delivered in the user's preferred language. Module E provides the cross-cutting language infrastructure consumed by every other module. Default language is **English**; Phase 1 ships **Arabic (RTL)**. See [PRD §6.10](../PRD.md) and [docs/i18n-architecture.md](../i18n-architecture.md) for full specification.

| # | Task | PRD Ref | Agent Type | Complexity | HITL | Priority |
|---|------|---------|-----------|-----------|------|----------|
| E-01 | **Language Detection & Routing Agent** — Detect the language of every incoming user query; normalise intent cross-lingually; inject `locale` metadata into all downstream Genkit flows; ensure AI responses are generated in the user's active locale | §6.10.5, US35 | Reactive | L1 | No | Phase 2 |
| E-02 | **Query Translation Agent** — Translate Arabic queries to English intent representation before Text-to-SQL generation; preserve numerical entities, date ranges, and domain terms without translation | §6.10.5, US35 | Reactive | L1 | No | Phase 2 |
| E-03 | **Localised Response Formatter** — Post-process AI-generated responses to apply locale-correct number formatting, currency symbols, date strings, and directional punctuation before returning to the chat client | §6.10.4, US35, US36 | Reactive | L1 | No | Phase 2 |
| E-04 | **Multilingual Report Generator** — Generate PDF and Excel reports in the user's chosen locale (`en`, `ar`, or `both`); apply Arabic fonts (Cairo / Noto Sans Arabic), RTL layout, and locale-aware Intl formatting in PDF pipeline | §6.10.6, US37 | Reactive | L2 | No | Phase 5 |
| E-05 | **Translation Coverage Guard** — CI-level agent that compares `en` and `ar` translation file keys; fails build if any `en` keys are missing from `ar` namespace files; reports coverage percentage per namespace | §6.10, §10 NFR | Proactive | L1 | No | Phase 1 |
| E-06 | **Multilingual Notification Agent** — Deliver scheduled report emails, KPI alert messages, and approval workflow notifications in the locale stored on the user profile and `scheduled_reports` table | §6.10.2, US38 | Proactive | L1 | No | Phase 6 |
| E-07 | **Bilingual Glossary Sync Agent** — Maintain the governed bilingual medical + accounting terminology glossary (`docs/i18n-glossary.md`); inject updated glossary into A-04 Domain Terminology Normaliser and B-02 OCR post-processor on every deploy | §6.10, A-04, B-02 | Proactive | L1 | Yes (glossary review) | Phase 3 |

---

## 3. Agent Count Summary

| Module | # Agents | Max Complexity | Phases |
|--------|----------|---------------|--------|
| A — Conversational BI | 13 | L3 | 2–7 |
| B — AI Accountant | 16 | L3 | 4–7 |
| C — Easy Reports | 8 | L2 | 1, 8–10 |
| D — Advanced Search Analytics | 14 | L3 | 13–16 |
| **E — Language & Localisation** | **7** | L2 | **1–6** |
| **Total** | **58** | | |

---

## 4. HITL Summary (Human-in-the-Loop)

The following agents **always** require human approval before taking an irreversible action on Tally (write-back):

| Agent ID | Trigger | Required Approver |
|----------|---------|------------------|
| B-08 | Any transaction post | Accountant → Manager → Finance |
| B-09 | Tally sync (journal / invoice / bill creation) | Finance Head ("One-Click" explicit action) |
| B-10 | Unmatched bank items | Accountant |
| D-10 | PO or any Tally write-back triggered by insight | Finance Head |
| E-07 | Bilingual glossary updates (medical/accounting terms) | Medical Advisor + Finance Advisor |

> See [03-governance-security.md](./03-governance-security.md) for the full policy enforcement design.

---

## 5. Cross-Cutting Agent Services (Shared Infrastructure)

These are not standalone agents but shared services consumed by multiple agents:

| Service | Consumers | Purpose |
|---------|-----------|---------|
| **Confidence Scoring Service** | A-06, B-02, B-05, B-10, D-04 | Standardised 0–100% score + routing logic |
| **Audit Log Writer** | B-08, B-09, B-14, D-10, C-05 | Immutable append-only audit trail |
| **Schema Context Cache** | A-01, A-02, D-01, D-03, D-12 | Pre-loaded DB schema + semantic metadata for LLM context |
| **Policy Engine (OPA)** | All agents | Enforce RBAC, read-only BI guardrails, and write-back policies |
| **Notification Dispatcher** | A-10, B-08, C-03, D-13, **E-06** | Email/SMS/Slack/in-app fanout |
| **Semantic Layer Registry** | A-01, D-01, D-03, D-09 | Governed metric/dimension definitions |
| **Language & Locale Service (E-01/E-03)** | All chat agents (A-01–A-13), D-series, E-series | Locale detection, cross-lingual intent normalisation, response formatting |
