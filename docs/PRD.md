
# Product Requirements Document (PRD): AI-Powered Conversational BI Dashboard

## 1. Overview & Context

**Product Name:** MediSync
**logo:**public/logo.png

**Problem Statement:** The client operates a multi-faceted healthcare business (clinic, pharmacy) using HIMS for operations and Tally for accounting. Currently, data is siloed. Extracting cross-platform insights (e.g., correlating clinic footfall with realized pharmacy revenue) requires manual data dumping and spreadsheet manipulation, which is slow, error-prone, and static.

**Vision:** To build an intelligent, chat-based BI dashboard powered by an AI Agent that seamlessly merges HIMS and Tally data. The platform will act as an on-demand "AI Data Analyst," allowing users to query their business data in natural language and receive instant, accurate tables, charts, and downloadable reports.

---

## 2. Goals & Success Metrics (KPIs)

### Business Goals

* **Eliminate Manual Reporting:** Reduce the time spent generating weekly/monthly financial and operational reports by 90%.
* **Unified Visibility:** Provide a single source of truth for both operational (HIMS) and financial (Tally) health.
* **Data-Driven Decisions:** Empower non-technical staff to pull complex analytics instantly using conversational prompts.

### Success Metrics (KPIs)

* **Query Accuracy:** >95% accuracy in Text-to-SQL conversions (validating that the AI pulls the correct data).
* **System Latency:** <5 seconds to generate an answer/chart for standard queries.
* **Adoption Rate:** 100% of the target management team logging in at least twice a week.

---

## 3. Target User Personas

| Persona | Role & Needs | Example Queries |
| --- | --- | --- |
| **Clinic Manager** | Needs operational oversight, patient footfall metrics, and doctor performance. | *"Show me the patient volume for Dr. Smith this month vs. last month."* |
| **Pharmacy Lead** | Needs inventory turnover, top-selling drugs, and expiry alerts. | *"What were the top 10 selling pharmacy items last week by revenue?"* |
| **Accountant / Owner** | Needs P&L, cash flow, outstanding receivables, and combined revenue. | *"Show a bar chart comparing clinic revenue in HIMS vs. actual cash collected in Tally for Q1."* |

---

## 4. Architecture & Recommended Tech Stack

To ensure Tally and HIMS performance isn't impacted by heavy AI queries, the system will use a decoupled architecture (ETL to a Data Warehouse).

* **Data Extraction Layer (ETL):**
* **Tally:** Tally API using TDL (Tally Definition Language) to expose data via XML/JSON HTTP POST/GET requests.
* **HIMS:** Native REST APIs from the HIMS provider.
* **Orchestration:** **Go (for performance)** or Apache Airflow to sync data incrementally every 15–60 minutes.


* **Data Warehouse:** **Local PostgreSQL** (Self-hosted on-premises).
* **AI & LLM Engine:**
* **Framework:** **Genkit** (Standardizing on Go/TS for type-safe AI flows) and **Agent ADK** for multi-agent coordination.
* **Communication:** **Google A2A Protocol** for standardized inter-agent discovery and collaboration across the 14+ agent ecosystem.
* **Vector Storage:** **pgvector** (Local Postgres extension) or **Milvus** for dedicated vector indexing.
* **Cache & Message Broker:** **Redis** (Local instance).
* **UI Pattern:** **Generative UI** orchestrated via **CopilotKit** (Ref: [official CopilotKit GenUI Guide](https://github.com/CopilotKit/generative-ui/blob/main/assets/generative-ui-guide.pdf)).
* **LLM:** GPT-4o, Claude 3.5 Sonnet, or Gemini 1.5 Pro (via local gateway or strictly local models like Llama 3 on Ollama).


* **Frontend UI:** React.js for Web, **Flutter** for Mobile.
* **Visualization Library:** **Apache ECharts** or Plotly (rendered dynamically via go-echarts or React components).
* **Internationalisation (i18n):**
  * **Web:** `i18next` + `react-i18next` for namespace-based JSON translations with lazy-loading.
  * **Mobile:** Flutter `flutter_localizations` + `intl` package with ARB files for compile-time type-safe strings.
  * **RTL Layout:** Tailwind CSS logical properties (`inline-start/end`) + HTML `dir="rtl"` on root for Arabic; Flutter `Directionality` + `EdgeInsetsDirectional`.
  * **Default Locale:** `en` (English, LTR). Phase 1 also ships `ar` (Arabic, RTL).

---

## 5. Security & Identity Integration

* **Identity Management:** **Keycloak** (Self-hosted) for RBAC, OIDC, and medical-grade audit logs.
* **Agentic Security:** Automated offensive security and red-teaming via **Redamon**, paired with defensive policy enforcement using **OPA** (Open Policy Agent).

---

## 6. AI Agent Skills & Capabilities

The AI Agent is not just a chatbot; it requires specialized "Agent Skills" to function as a data analyst.

1. **Text-to-SQL Skill:** The core engine. The agent takes a natural language question, reads the Data Warehouse schema context, translates the intent into a secure SQL query via **Genkit flows**, executes it, and retrieves the raw data.
2. **Visualization Routing:** The agent determines the best way to present data using **Apache ECharts** configurations.
3. **Self-Correction (Error Handling):** If the generated SQL query fails, the agent (powered by Genkit's retry logic) reads the database error and rewrites the query automatically.
4. **Domain Expertise Skill (Healthcare & Accounting):** The system prompt will be pre-loaded with domain knowledge (e.g., knowing that "footfall" means patient visits, or "outstanding" means accounts receivable).
5. **Language Detection & Routing Skill (E-01):** All user queries pass through a pre-processing agent that detects the query language, normalises intent cross-lingually, injects the user's locale into every downstream Genkit flow, and ensures AI responses are delivered in the user's preferred language with correct locale formatting.

---

## 6. Functional Requirements

### 6.1. The Conversational Interface

* **Chat Window:** A ChatGPT-style interface where users can type freely.
* **Pre-defined Prompt Buttons:** A carousel of 5–10 quick-action buttons above the chat (e.g., *"Today's Total Revenue,"* *"Pending Tally Invoices,"* *"Low Pharmacy Stock"*). Clicking these injects the prompt instantly.
*   **Chat Window:** A ChatGPT-style interface where users can type freely.
*   **Pre-defined Prompt Buttons:** A carousel of 5–10 quick-action buttons above the chat (e.g., *"Today's Total Revenue,"* *"Pending Tally Invoices,"* *"Low Pharmacy Stock"*). Clicking these injects the prompt instantly.
*   **Dynamic Rendering:** Responses must embed rich widgets (tables, bar charts, line graphs) directly inside the chat flow, not just text.

### 6.2. Dashboard & Reporting

*   **Pin to Dashboard:** Users can click a "Pin" icon on any generated chart to save it to a permanent, auto-refreshing main dashboard (similar to PowerBI).
*   **Export & Download:** Every generated table or chart must have a "Download" button allowing export to CSV, Excel (.xlsx), or PDF.
*   **Scheduled Reports:** Automated report generation and delivery via email on configurable schedules (daily, weekly, monthly).
*   **Custom KPI Dashboard:** Allow users to create personalized dashboards with custom key performance indicators (KPIs) that auto-refresh.
*   **Multi-Period Comparison:** Enable side-by-side comparison of data across different time periods (e.g., month-to-month, year-over-year).

### 6.3. Data Integration

* **Tally Sync:** Must pull Ledgers, Vouchers, Inventory Masters, and Sales/Receipt data.
* **HIMS Sync:** Must pull Patient Demographics, Appointments, Billing, and Pharmacy Dispensation data.

### 6.4. Financial Analytics Features

* **Profitability Analysis:** 
  - Revenue breakdown by department (clinic, pharmacy)
  - Cost analysis and margin calculations
  - Product/service-wise profitability reports
  - Contribution margin analysis

* **Cash Flow Management:**
  - Daily, weekly, and monthly cash flow forecasting
  - Inflow vs. outflow analysis
  - Cash position tracking
  - Outstanding receivables aging reports
  - Payables aging reports
  - Cash cycle analysis

* **Budget vs. Actual Analysis:**
  - Compare actual spending/revenue against budgeted amounts
  - Variance analysis with drill-down capability
  - Department-wise budget tracking
  - Historical trend analysis for better forecasting

* **Tax & Compliance Analytics:**
  - GST/VAT analysis and reconciliation
  - Tax collected vs. paid tracking
  - Compliance report generation
  - Period-wise tax liability calculation

* **Receivables & Payables Management:**
  - Outstanding invoice tracking
  - Payment pending reports with due dates
  - Customer credit limits and exposure
  - Vendor payment analysis
  - Days Sales Outstanding (DSO) and Days Payable Outstanding (DPO) metrics

### 6.5. Operational Analytics Features

* **Inventory Analytics:**
  - Stock aging and obsolescence reports
  - Inventory valuation (FIFO, LIFO, Weighted Average)
  - Stock turnover ratio analysis
  - Low-stock and out-of-stock alerts
  - Inventory movement pattern analysis
  - Supplier-wise inventory analysis

* **Drill-Down Capability:**
  - Click on any metric in a chart or table to drill down to transaction-level details
  - Navigate from summary trends to individual transaction records
  - Reverse drill-down to aggregate view

* **Real-Time Alerts & Notifications:**
  - Critical business metric alerts (e.g., negative cash flow, low inventory, outstanding payments due)
  - Configurable alert thresholds
  - Multi-channel notifications (in-app, email, SMS)

* **Mobile & Responsive Access:**
  - Full mobile-optimized interface for on-the-go data access
  - Offline capability for pre-loaded dashboards
  - Touch-friendly navigation and interactions
  - Mobile-optimized report generation and export

### 6.6. Advanced Analysis Capabilities

* **Trend Analysis & Forecasting:**
  - Historical trend visualization with predictive lines
  - Seasonal pattern detection
  - Growth rate calculation
  - Forecasting for revenue and expenses using AI/ML models

* **Comparative Analysis:**
  - Compare performance across different dimensions (departments, locations, doctors, products)
  - Performance benchmarking
  - Top/bottom performers identification

* **Data Quality & Audit Trail:**
  - Data validation and reconciliation reports
  - Audit logs for all data access and modifications
  - Data lineage tracking
  - Exception reports for data anomalies

---

## 6.7. AI Accountant Module: Automated Bookkeeping & Real-Time Tally Integration

The AI Accountant module adds intelligent automation for financial data processing, document management, and transaction reconciliation. This module transforms manual accounting tasks into AI-powered workflows.

### 6.7.1 Document Upload & OCR Processing

* **Bulk Statement & Bill Upload:**
  - Drag-and-drop interface for uploading bank statements, vendor bills, and credit card statements in multiple formats (PDF, Excel, CSV, scanned images, handwritten documents).
  - Support for batch uploads (hundreds of documents at once).
  - Real-time progress tracking and status indicators.

* **AI-Powered OCR Engine:**
  - Extract transaction details from documents with **95%+ accuracy** for standard documents.
  - Recognize **handwritten invoices and bills** using advanced OCR and pattern recognition.
  - Auto-detect invoice amounts, vendor names, invoice dates, and tax amounts.
  - Identify and extract key fields from various document formats and layouts.
  - Flag low-confidence extractions for manual review.

* **Document Classification:**
  - Auto-categorize documents (invoices, bills, bank statements, tax documents, receipts).
  - Organize by vendor, date range, and document type.
  - Support for custom folder structures and tagging.

### 6.7.2 Automated Transaction Mapping & Categorization

* **Intelligent Ledger Mapping:**
  - AI automatically matches transactions to appropriate Tally ledgers based on transaction details (amount, description, vendor, date).
  - Learn from past mapping patterns to improve accuracy over time.
  - Display **confidence scores** (70%, 85%, 95%+) for each AI-suggested mapping.
  - Allow users to override and correct mappings; apply corrections to similar future transactions.

* **Sub-Ledger & Cost Center Assignment:**
  - Auto-assign sub-ledgers based on transaction context.
  - Suggest cost centers (if applicable) for departmental expense tracking.
  - Support for multi-dimensional accounting hierarchies.

* **Vendor & Bill Matching:**
  - Match uploaded bills to existing vendor records in Tally.
  - Create new vendor records if unrecognized.
  - Track vendor payment terms and discount rates.
  - Identify duplicate invoices and prevent duplicate posting.

### 6.7.3 Real-Time Tally Synchronization

* **One-Click Sync to Tally:**
  - Seamless integration with Tally via native API or TDL (Tally Definition Language).
  - Push approved transactions directly into Tally with a single click.
  - Create journal entries, purchase bills, sales invoices automatically.
  - Update inventory records based on bill content.

* **Real-Time Sync Dashboard:**
  - Live connection status indicator (connected/disconnected/connecting).
  - Sync frequency configuration (real-time, every 5/15 minutes, hourly, manual).
  - Automatic retry logic if sync fails temporarily.
  - Sync history log showing last 20+ sync events with timestamps and status.
  - Manual "Sync Now" button for immediate synchronization.

* **Multi-Entity Support:**
  - Manage multiple company instances within Tally.
  - Switch between companies quickly and sync data independently.
  - Consolidated dashboards across all entities with drill-down capability.

* **Data Integrity & Validation:**
  - Pre-sync validation checks (duplicate detection, ledger availability, amount validation).
  - Audit trail logging all synced transactions with user, timestamp, and changes.
  - Reconciliation checks to ensure data consistency between source and Tally.

### 6.7.4 Bank Reconciliation & Matching

* **Automated Bank-to-Tally Matching:**
  - Compare bank statement transactions with Tally entries.
  - Auto-match based on amount, date, and transaction description.
  - Identify outstanding payments and outstanding receipts.
  - Calculate net reconciliation difference and isolate unmatched items.

* **Smart Transaction Matching:**
  - Handle partial matches and multi-part payments.
  - Match deposits/payouts to multiple invoices (uncleared checks, electronic transfers).
  - Suggest matches with confidence scoring when automatic matching fails.
  - Manual matching interface for review and correction.

* **Reconciliation Reports:**
  - "Outstanding Payments" report showing checks/orders awaiting clearance.
  - "Outstanding Receipts" report showing deposits in transit.
  - Aged reconciliation items (0-7 days, 8-30 days, 30+ days outstanding).
  - Exception reports flagging unusual patterns (e.g., recurring unmatched transactions).

### 6.7.5 Expense & Bill Management

* **Expense Categorization:**
  - Auto-categorize expenses (office supplies, utilities, travel, meals, etc.).
  - Apply category rules to speed up manual categorization.
  - Track expense trends by category.

* **Bill Tracking & Approval Workflow:**
  - Monitor bill due dates and payment status.
  - Flag overdue bills requiring immediate attention.
  - Approval workflow (user → manager → finance → posted).
  - Payment scheduling and cash flow impact preview.

* **Duplicate Detection:**
  - Identify duplicate bills before posting to prevent overpayment.
  - Flag suspicious duplicates for review (same amount, vendor, within X days).
  - Bulk delete or merge duplicate records.

### 6.7.6 Financial Reporting & Compliance

* **Automated Report Generation:**
  - Generate P&L, Balance Sheet, Cash Flow statements on-demand.
  - Tax-ready reports (GST reconciliation, TDS summary, tax liability calculation).
  - Audit trail reports for compliance and regulatory requirements.
  - Schedule automatic report generation and email delivery.

* **Tax Compliance Tracking:**
  - GST/VAT analysis and reconciliation (Input Tax Credit tracking).
  - Tax collected vs. paid tracking.
  - Period-wise tax liability calculation.
  - Compliance checklist with filing deadlines.

* **Audit Reports & Logging:**
  - Complete audit trail of all transactions and changes.
  - User activity logging (who posted what, when).
  - Data lineage tracking (source document → posted entry).
  - Exception reports for anomalies (unusual transactions, bulk changes, deleted entries).

### 6.7.7 Document Management & Library

* **Centralized Document Repository:**
  - Store all uploaded documents linked to corresponding transactions.
  - Searchable by document type, vendor name, date range, amount.
  - Version history for edited/corrected documents.
  - Secure storage with encryption and access controls.

* **Document Linking & Cross-Reference:**
  - Link documents to specific transactions in Tally.
  - Attach supporting documents to journal entries, bills, invoices.
  - Quick view of supporting documents from transaction detail pages.
  - Document retrieval during audits and tax inquiries.

* **Compliance & Archival:**
  - Automatic archival of old documents based on retention policies.
  - Secure deletion with audit trail noting document destruction.
  - Compliance with document retention regulations (GST, tax law, corporate governance).

### 6.7.8 Analytics & Insights

* **Cash Flow Forecasting:**
  - Project future cash position based on bill payment and receipt schedules.
  - Identify potential cash shortfalls and surpluses.
  - What-if analysis for scenario planning.

* **Expense Analysis:**
  - Top expense categories and trends (YoY, month-over-month).
  - Vendor-wise spending analysis.
  - Budget vs. actual spending with variance analysis.
  - Departmental cost allocation and profitability.

* **Payables & Receivables Metrics:**
  - Days Payable Outstanding (DPO) - average payment delay.
  - Days Sales Outstanding (DSO) - average collection time.
  - Aging analysis for both payables and receivables.
  - Supplier payment patterns and discrepancy tracking.

---

## 6.8. Easy Reports Module: Automated MIS & Business Intelligence Dashboards

The Easy Reports module extends the platform with enterprise-grade reporting and business intelligence capabilities, transforming raw ERP data into actionable Management Information System (MIS) reports and interactive dashboards.

### 6.8.1 Pre-Built Report Library

* **Sales & Debtors Analytics:**
  - Sales performance by salesperson, area, channel, segment, and location.
  - Customer-wise revenue and sales trend analysis.
  - Debtor aging reports showing outstanding receivables by age bracket (0-30, 31-60, 61-90, 90+ days).
  - Top customers by revenue and sales growth analysis.
  - Territory-wise sales performance and quota tracking.

* **Profitability Reports:**
  - Profitability by customer, location, product, product group, and salesperson.
  - Margin analysis and contribution margin calculations.
  - Department-wise and cost-center profitability breakdowns.
  - Product-wise revenue vs. cost analysis.
  - Segmental profitability for multi-business ventures.

* **Targets & Budgets Analysis:**
  - Sales targets vs. actuals (by salesperson, location, product, department).
  - Expense budgets vs. actuals with variance analysis.
  - Monthly/quarterly/annual budget tracking and alerts for budget overages.
  - Achievement percentage tracking and trend forecasting.
  - What-if analysis for budget revisions.

* **Financial Statements:**
  - Trial Balance and General Ledger reports.
  - Profit & Loss Statement (P&L) with multi-period comparison.
  - Balance Sheet with cost-center wise breakdown.
  - Cash Flow Statement (operating, investing, financing activities).
  - Cost Centre Trial Balance and profitability analysis.

* **Revenue Reporting:**
  - Cost-center-wise revenue, expenses, and profitability.
  - Ledger-wise analysis with drill-down to voucher details.
  - Revenue by type (services, products, other income).
  - Recurring vs. one-time revenue tracking.
  - Revenue recognition and compliance reporting.

* **Inventory Analysis:**
  - Stock aging reports identifying slow-moving and obsolete inventory.
  - Inventory valuation (FIFO, LIFO, weighted average) with variance analysis.
  - Item-wise movements and turnover ratios.
  - Stock shortage and overstock alerts.
  - Supplier-wise inventory and lead time analysis.
  - Warehouse/location-wise inventory distribution.

### 6.8.2 Consolidated Multi-Company Reporting

* **Entity Consolidation:**
  - Combine financial data from multiple Tally company instances into unified dashboards.
  - Consolidated P&L, Balance Sheet, and Cash Flow across all entities.
  - Eliminate inter-company transactions automatically.
  - Compare performance across subsidiaries and locations.
  - Drill-down from consolidated view to individual entity details.

* **Multi-Company KPI Dashboards:**
  - Aggregate KPIs across all entities (total revenue, expenses, profitability, cash position).
  - Entity-wise performance comparison and benchmarking.
  - Group-wise reporting for holding company structures.
  - Real-time consolidation with dynamic period selection.

### 6.8.3 Zero-Code Report Builder

* **Drag-and-Drop Report Design:**
  - Business users (non-technical) can create custom reports without coding.
  - Drag fields from Tally data model to build reports visually.
  - Pre-defined templates for common report types (tables, matrices, charts).
  - Formula builder for calculated fields and custom metrics.
  - Conditional formatting and highlighting based on thresholds.

* **Report Customization:**
  - Add/remove columns, change field order, rename fields for clarity.
  - Apply filters (by date, department, salesperson, product, cost center).
  - Sort and group by multiple dimensions.
  - Add subtotals and grand totals with different aggregation methods (sum, average, count, min, max).
  - Pagination and row limiting for large datasets.

* **Ad-Hoc Reporting:**
  - Business users explore data and create on-demand reports without IT involvement.
  - Save custom reports for reuse and share with team members.
  - Version history to track report changes over time.

### 6.8.4 Interactive Dashboards & Visualization

* **Pre-Built Executive Dashboards:**
  - Sales Dashboard: top customers, revenue trend, pipeline, forecast.
  - Finance Dashboard: P&L snapshot, cash flow, budget vs. actual, expense breakdown.
  - Inventory Dashboard: stock levels, aging, movement, valuation.
  - HR Dashboard (if integrated): headcount, cost, productivity metrics.
  - Operations Dashboard: KPI summary, alerts, bottlenecks.

* **Interactive Charts & Visualizations:**
  - Bar charts, line graphs, pie charts, scatter plots for multi-dimensional data.
  - Drill-down capability: click chart element to view detailed transaction data.
  - Real-time data refresh (configurable intervals).
  - Chart type switching (bar ↔ table ↔ pie) for user preference.
  - Export charts as images or embedded in reports/emails.

* **KPI Scorecards:**
  - Single-metric or multi-metric cards showing current value, target, variance, and trend.
  - Color-coded status indicators (red/amber/green based on thresholds).
  - Historical comparison (month-over-month, year-over-year growth).
  - Sparklines showing trend direction or heatmaps for pattern detection.

### 6.8.5 Automated Report Scheduling & Distribution

* **Report Scheduling:**
  - Create schedules: daily, weekly, monthly, quarterly, yearly, or custom intervals.
  - Schedule reports to run at specific times (e.g., 6 AM Monday for weekly reports).
  - Conditional scheduling (e.g., only run if threshold violated).
  - One-time reports or recurring series.

* **Email & Distribution:**
  - Automatically email reports to distribution lists (stakeholders, managers, team leads).
  - Customize email message and subject line.
  - Attach reports in multiple formats: PDF, Excel, CSV, HTML, PowerPoint.
  - Schedule different reports for different recipients based on roles.
  - Tracking: logs showing delivery success/failure.

* **Portal & Access:**
  - Self-service report portal where users access their scheduled reports.
  - Archive of historical reports (searchable by date, report type).
  - On-demand report regeneration with updated data.
  - Mobile-friendly report viewing.

### 6.8.6 Tally Data Integration & Sync

* **Native Tally Connector:**
  - Direct integration with Tally ERP 9, TallyPrime via native API or data export.
  - Automatic data extraction (ledgers, vouchers, inventory, cost centers).
  - Real-time or scheduled data refresh (daily, hourly, custom intervals).
  - Support for Tally company databases (local or networked).
  - Handles Tally's dynamic chart of accounts and inventory structures.

* **Data Validation & Quality:**
  - Data integrity checks (duplicate ledgers, orphaned transactions).
  - Reconciliation between extracted data and Tally source.
  - Data quality reports highlighting anomalies or missing data.
  - Audit trail logging all data extraction and transformations.

### 6.8.7 Excel & External Data Integration

* **Excel Export and Refresh:**
  - Export any report/dashboard to Excel with formatting preserved.
  - Create "live" Excel reports that refresh data on-demand or schedule.
  - Excel pivot table support for further analysis by end users.
  - Two-way sync: upload corrected data back to source if permitted.

* **External Data Sources:**
  - Link Tally data with external targets (sales targets, budgets, headcount from HR system).
  - Combine Tally data with custom dimensions (e.g., sales commissions from non-Tally systems).
  - Blend data from multiple systems (Tally + HIMS) for cross-domain analytics.
  - Support for connecting to databases (MySQL, SQL Server, Oracle).

### 6.8.8 Role-Based Security & Access Control

* **User Roles & Permissions:**
  - Define custom roles (Finance Manager, Sales Lead, Accountant, Viewer, Admin).
  - Control access to specific reports and dashboards per role.
  - Restrict actions: who can view, export, modify, schedule reports.
  - Allow/deny export to Excel, PDF based on role permissions.

* **Data-Level Security:**
  - Row-level security: users see only data relevant to them (e.g., salesperson sees own sales).
  - Column-level security: hide sensitive columns (e.g., cost, margin) from viewers.
  - Ledger-level restrictions: finance team sees profit centers, operations sees cost centers.
  - Automatic filtering: apply user's territory/department filter to all reports.

* **Audit & Compliance:**
  - Log all access to reports, exports, and dashboard views.
  - Track who ran which report, when, and any data exports.
  - Compliance reports for SOX, GST, tax audits.

### 6.8.9 Custom Fields & User-Defined Dimensions

* **Custom Fields (UDFs):**
  - Support for Tally's User-Defined Fields (custom columns in ledgers, items, vouchers).
  - Include custom fields in reports and filters.
  - Create custom metrics by combining standard and custom fields.

* **Custom Dimensions:**
  - Define business dimensions not available in Tally (e.g., product category variants, region codes).
  - Link custom dimensions to Tally master records (items, ledgers).
  - Use custom dimensions for slicing and dicing data in dashboards.

### 6.8.10 Mobile & Responsive Reporting

* **Mobile-Optimized Dashboards:**
  - Responsive design adapts to mobile screens.
  - Simplified dashboard layouts for mobile (fewer charts, larger text).
  - Touch-friendly interactions (tap to drill-down, swipe to navigate).
  - Offline capability: downloaded reports accessible without internet.

* **Native Mobile Apps (Optional):**
  - iOS and Android apps for quick access to top dashboards.
  - Push notifications for critical alerts (budget exceeded, sales target missed).
  - Voice-based querying: "Show me today's revenue" (future enhancement).

### 6.8.11 Advanced Analytics & Forecasting

* **Trend Analysis:**
  - Historical trend charts (12-month, 3-year trends).
  - Monthly/quarterly/annual growth rate calculations.
  - Seasonal pattern detection and adjustment.
  - Anomaly detection (unusual spike or drop) with alerts.

* **Forecasting & Projections:**
  - Revenue and expense forecasting using historical trends.
  - Cash flow projections based on receivables/payables aging and payment patterns.
  - Budget variance forecasting (project year-end actual vs. budget).
  - "What-if" scenario modeling for decision support.

* **Performance Benchmarking:**
  - Compare KPIs against previous period (month-to-month, year-to-year).
  - Industry benchmarks (if integrated with external data).
  - Peer comparison (compare unit/location performance).
  - Target achievement tracking and variance explanation.

### 6.8.12 Cost-Center & Department Analytics

* **Cost-Center Reporting:**
  - Costs, revenues, and profitability by cost center.
  - Budget allocation and tracking per cost center.
  - Cost-center-wise headcount and productivity metrics.
  - Overhead allocation methods (fixed, variable, activity-based).

* **Department Performance:**
  - Department-wise KPIs (sales, expenses, efficiency, quality).
  - Comparative analysis across departments.
  - Department contribution to overall profitability.
  - Recharge/transfer pricing for shared services departments.

---

## 6.9. Advanced Search-Driven Analytics Module: Agentic AI Conversational Intelligence

This module brings next-generation conversational analytics capabilities, enabling users to interact with data using natural language search and autonomous AI agents that proactively analyze data and generate insights without explicit queries.

### 6.9.1 Natural Language Search-Driven Analytics

* **Intelligent Search Interface:**
  - Google-like search experience for querying data across all sources (HIMS, Tally, external databases).
  - Auto-complete suggestions based on frequently asked questions and available metrics.
  - Contextual spell-check and query correction for typos or ambiguous terms.
  - Search history and saved query management for quick re-execution.
  - Guided examples and suggested questions for new users.

* **Multi-Step Conversational Analysis:**
  - System interprets complex, multi-part questions and performs sequential analyses automatically.
  - Example: "Show me top 5 selling drugs in February, compare with January, identify price changes, and suggest reorder quantities" — answered in one query.
  - Automatic follow-up question generation for deeper exploration (drill-down suggestions).
  - Chained reasoning: combines multiple data sources and metrics to answer complex business questions.

* **Natural Language Understanding & Context:**
  - Business synonym recognition (e.g., "revenue", "sales", "turnover", "income" understood as equivalent).
  - Temporal context understanding: "last month", "this quarter", "year-to-date", "previous 90 days" automatically interpreted.
  - Entity recognition: identifies patients, doctors, products, pharmacies mentioned in queries.
  - Ambiguity resolution: offers clarifications when query intent is unclear.

### 6.9.2 Autonomous AI Agents for Self-Service Analytics

* **Agentic AI Analyst (Spotter-like Agent):**
  - Autonomous agent acting as a data analyst on your team, running complex analytical workflows without human intervention.
  - Performs multi-step analysis: data retrieval → statistical analysis → comparisons → forecasting → actionable recommendations.
  - Autonomous insight discovery: proactively suggests analyses and patterns users may need before being asked.
  - Continuous learning: improves recommendations based on user interactions and feedback.
  - Multi-source reasoning: analyzes across Tally (GL, invoices, payments), HIMS (appointments, billing, inventory), and external data simultaneously.

* **Deep Research Agent:**
  - Autonomously discovers hidden patterns, correlations, and anomalies in datasets.
  - Performs advanced statistical analyses (regression, clustering, outlier detection, seasonal decomposition).
  - Auto-generates hypotheses and validates them against data with confidence scoring.
  - Produces structured research reports with supporting evidence and visualizations.
  - Automatically flags concerning trends and unusual patterns requiring attention.

* **Autonomous Task Scheduling & Monitoring:**
  - Agents execute scheduled analyses and alert users when business-critical thresholds are violated.
  - Example: "Notify me if any clinic's revenue drops >10% month-over-month or pharmacy margin falls below 20%."
  - Business rule learning: agents autonomously generate insights aligned with user-defined business rules.
  - Delegation of recurring analysis tasks: set it and forget it, with agents running analyses autonomously.

### 6.9.3 AI-Generated Dashboards & Instant Visualization

* **Intelligent Dashboard Auto-Generation:**
  - AI automatically designs and generates complete dashboards from search queries or raw data.
  - Intelligent chart type selection: system recommends optimal visualizations (bar for rankings, line for trends, scatter for correlations, heatmap for multi-dimensional analysis).
  - Automatic layout optimization: creates professional, consistency-driven dashboards in seconds (vs. hours of manual design).
  - Data storytelling: dashboards include titles, key findings, and supporting narrative context.
  - Interactive chart editing with natural language: users can modify charts using plain English ("show top 10 only", "change to pie chart", "add forecasting trend line").

* **AI-Augmented Real-Time Dashboards:**
  - Dynamic dashboards that refresh and adapt based on latest data ingestion.
  - Contextual insights embedded in visualizations (annotations like "20% growth vs. last quarter", "highest in 3 years").
  - Automated anomaly highlighting: charts visually highlight unusual or outlier data points.
  - What-if simulation: drag sliders to see how KPIs change with different parameter values.
  - Personalization: dashboards auto-adapt to user role and frequently accessed metrics.

### 6.9.4 Autonomous Insights & Prescriptive Recommendations

* **Automated Insight Discovery:**
  - System autonomously identifies important patterns, trends, and anomalies in data.
  - Example insights: "Revenue peaked on Fridays this month", "Pharmacy margin declined 3% vs. last quarter", "Dr. Smith's patient satisfaction 12% above peer average", "Inventory for Item XYZ will stock out in 5 days".
  - Smart insight prioritization: surfaces most important/actionable insights at top of reports.
  - Confidence scoring: quantifies reliability of each insight (95% confidence vs. 78% confidence).

* **Prescriptive Recommendations Engine:**
  - Goes beyond insights to provide specific, actionable recommendations.
  - Example: "Expected inventory shortage in Item XYZ within 5 days — recommend urgent reorder of 150 units by Thursday."
  - Root cause analysis: automatically explains why a situation occurred (e.g., "shortage driven by 35% surge in weekly sales due to seasonal demand").
  - Prescriptive actions: suggests specific steps to take (reorder quantity, timing, preferred supplier, payment terms).
  - Impact estimation: quantifies expected business impact ("Implementing this recommendation will improve cash flow by ₹2.5 Lakh and reduce stockouts by 95%").

### 6.9.5 Semantic Layer & Governed Metrics Infrastructure

* **Unified Semantic Model:**
  - Single source of truth for all business definitions, metrics, and relationships across the organization.
  - Pre-defined business metrics (revenue, profit margin, inventory turnover, DSO, DPO, patient acquisition cost) pre-built and governed.
  - Metric governance: ensure consistent definitions across all dashboards, reports, and searches.
  - Version control: audit metric definition changes (who changed what, when, and why).
  - Relationship mapping: defines connections between entities (customers → orders → invoices → revenue).

* **Agent-Ready Metadata & Metric Families:**
  - Data model optimized for AI agents to understand business context and relationships.
  - Automatic relationship detection (doctors → appointments → patient outcomes, drugs → sales → inventory → reorders).
  - Metric hierarchies: revenue decomposed into Sales Revenue + Service Revenue + Other Revenue with drill-down paths.
  - Custom metric creation: formula builder for business-specific KPIs and calculated fields.
  - Dimension definitions: defines time hierarchies (year → quarter → month → week → day), geography hierarchies, product categories.

### 6.9.6 Embedded Analytics & Workflow Automation

* **Embedded Conversational Analytics:**
  - Integrate search-driven analytics directly into clinic management and pharmacy software (not a separate tool).
  - White-labeled dashboards and search experiences for specific business applications.
  - API-driven embedding: external apps call analytics engine and embed results seamlessly using REST/GraphQL APIs.
  - Mobile-responsive embedding: dashboards and search accessible on phones/tablets with same functionality as desktop.

* **Insight-Driven Workflow Automation:**
  - Convert AI-generated insights into automated business actions (e.g., "Low pharmacy stock" → auto-trigger PO creation in Tally).
  - Two-way integration: insights flow into business systems (CRM, ERP, accounting), and conversely, system changes trigger re-analysis.
  - Alert-based workflows: when specific conditions detected, automatically notify relevant teams via email/SMS/Slack.
  - Audit trails: log all actions triggered by AI recommendations for compliance and traceability.

### 6.9.7 Self-Service Analytics Workbench (Analyst Studio)

* **Analyst Studio - Python & SQL Notebooks:**
  - Python notebooks for data scientists to build custom ML models, statistical analyses, and advanced workflows.
  - Ad-hoc SQL queries for analysts to explore data and answer bespoke questions without IT bottlenecks.
  - Data mashups: combine Tally GL data with HIMS appointment data with external market data in one analysis.
  - Version control and collaboration: share notebooks with team, incorporate feedback, publish approved analyses.

* **Self-Service Analysis Publishing:**
  - Analysts save and publish analysis results as governed, discoverable assets.
  - Published analyses become queryable by business users via natural language search ("Show me the analysis about drug margins").
  - Audit trails: track who created analysis, who executed it, when, and what actions resulted.
  - Approval workflows: data governance team approves analyses before publishing to broader organization.

### 6.9.8 AI-Assisted Low-Code Development (Developer IDE Integration)

* **Natural Language to Code Generation:**
  - Developers describe desired feature in plain English: "Create a patient satisfaction scorecard showing ratings by doctor with trend lines."
  - AI generates correct embedding code, UI components, database queries, and styling automatically.
  - Multi-language support: Python, JavaScript, React, Vue, TypeScript, Java, Node.js, C#.
  - Built-in best practices: generated code follows security standards, performance optimizations, and maintainability conventions.

* **Rapid Prototyping & Development Acceleration:**
  - Speeds development from weeks to days by auto-generating boilerplate and integration code.
  - Low-code/no-code for common analytics use cases (charts, filters, aggregations, sorting).
  - Fallback to manual coding: developers can edit generated code for complex custom logic.
  - Template library: leverage pre-built components and patterns for consistency across applications.

### 6.9.9 Real-Time Data Processing & Zero-Copy In-Memory Architecture

* **High-Performance In-Memory Analytics:**
  - Zero-copy architecture: data analyzed in-place without expensive copy/move operations between systems.
  - Sub-second query response times even on large datasets (100M+ rows, terabytes of data).
  - Real-time data ingestion: updates from Tally and HIMS reflected instantly in search results and dashboards.
  - Automatic data freshness management: system maintains freshness without manual trigger/refresh operations.

* **Federated Multi-Source Querying:**
  - Query across Tally (GL, invoices, payments), HIMS (appointments, billing, inventory), and external data in a single search.
  - Automatic federated query optimization: system routes queries optimally across heterogeneous data sources.
  - Data lineage transparency: users understand which source system each metric/dimension comes from.
  - Schema autodiscovery: system auto-catalogs all available tables, columns, relationships, and data types.

### 6.9.10 Mobile Search-Driven Intelligence

* **Mobile-Optimized Search Experience:**
  - Simplified search interface optimized for small screens (phones, tablets).
  - Touch-friendly search bar with swipe navigation and gesture-based interactions.
  - Voice search capability: ask questions verbally, receive instant answers.
  - Offline caching: pre-load critical dashboards for access without internet connectivity.

* **Mobile-Specific Insights & Alerts:**
  - Push notifications to mobile when important business events occur (revenue spike, stock alert, approval needed).
  - Mobile-specific chart designs: larger text, single-column layouts, simplified interactions vs. desktop.
  - Simplified dashboard variants: fewer widgets, larger content areas optimized for mobile consumption.
  - Mini summaries and quick reports formatted for mobile screens.

### 6.9.11 Data Governance & Role-Based Search Security

* **Search-Level Access Control:**
  - Control which questions different roles can ask and data they can access.
  - Example: Pharmacy Manager searches only return pharmacy metrics/inventory data; clinic manager sees clinic-specific data.
  - Sensitive question whitelisting: restrict queries about sensitive topics (payroll, margins) to appropriate roles.
  - Query auditing: comprehensive logs of all searches, results viewed, and data downloaded by user.

* **Column-Level & Row-Level Data Masking:**
  - Column masking: hide sensitive fields (patient names, cost prices, employee salaries) from unauthorized roles.
  - Row-level filtering: users only see data relevant to their department/region/team (finance team only sees financial data).
  - Automatic filtering applied to all search results and dashboards based on user permissions.
  - Compliance enforcement: GDPR/HIPAA-compliant access controls enforced at data layer for healthcare data protection.

### 6.9.12 Advanced Analytics & Predictive Intelligence

* **Statistical Analysis & Forecasting:**
  - Trend analysis with confidence intervals and statistical significance testing.
  - Seasonal decomposition: automatically separate seasonal patterns from underlying trends.
  - Time-series forecasting: predict future revenue, expenses, and inventory using statistical/ML models.
  - Anomaly detection: automatically identify unusual patterns and assign confidence scores.

* **Comparative & Cohort Analysis:**
  - Cohort analysis: segment patients by characteristics (age, location, procedure type) and compare behaviors and outcomes.
  - Peer benchmarking: compare clinic A's metrics against similar clinics B and C for performance insights.
  - Attribution modeling: understand which factors drive revenue or patient acquisition.
  - Causal analysis: identify true cause-and-effect relationships (not just correlations) using advanced statistical methods.

---

* **Role-Based Access Control (RBAC):** 
  - Define custom roles (Admin, Manager, Accountant, Analyst, Viewer).
  - Granular permission controls at dashboard, report, and data level.
  - Department-specific data visibility (e.g., Pharmacy Manager only sees pharmacy analytics).
  - Time-based access restrictions if needed.

* **User Management:**
  - Add/edit/remove user accounts.
  - Bulk invite functionality.
  - Password and session management.
  - User activity logging and audit trails.

---

## 6.10. Multi-Language Support & Localisation (i18n / l10n)

MediSync targets a bilingual workforce where operational staff, managers, and accountants may communicate in either English or Arabic. Multi-language support is a **first-class product requirement** — not a post-launch add-on.

**Full architecture details:** See [docs/i18n-architecture.md](../i18n-architecture.md)

### 6.10.1 Supported Locales

| Code | Language | Script | Direction | Phase |
|------|----------|--------|-----------|-------|
| `en` | English | Latin | LTR | Default — Phase 1 |
| `ar` | Arabic | Arabic | **RTL** | Phase 1 |

### 6.10.2 Coverage Scope

Every user-facing surface must respect the active locale:

| Surface | i18n Requirement |
|---------|------------------|
| **All Screens & Navigation** | All labels, menus, headers, tooltips, error messages translated |
| **Chat Interface** | User can type queries in Arabic or English; AI responds in matching language |
| **AI Responses** | Explanations, chart titles, table headers, insight narratives in user's language |
| **Reports (PDF/Excel)** | Report content, column headers, and summaries in chosen locale; Arabic PDF uses RTL layout and Arabic-compatible fonts (Cairo / Noto Sans Arabic) |
| **Dashboards & Charts** | Axis labels, legends, KPI card titles, and tooltip text localised; chart layout mirrors for RTL |
| **Email & Notifications** | Scheduled report emails and KPI alert notifications delivered in user's locale |
| **Validation & Error Messages** | Form validation and system error messages in user's language |
| **OCR Extraction (AI Accountant)** | Arabic-script invoice processing; extracted field labels in UI shown in Arabic |

### 6.10.3 RTL Layout Requirements

* All layout code must use **CSS logical properties** (`margin-inline-start`, `padding-inline-end`) — no hardcoded `left`/`right` physical properties.
* HTML `dir` attribute set to `rtl` on root element when Arabic locale is active; all components mirror automatically.
* Navigation sidebar moves to right edge; breadcrumb chevrons reverse direction; progress bars fill right-to-left.
* Chat bubbles mirror: user messages appear on the left in RTL; AI messages on the right.
* Directional icons (arrows, chevrons) flip via `rtl:scale-x-[-1]` CSS utility.
* Flutter: `EdgeInsetsDirectional`, `AlignmentDirectional`, and `GlobalMaterialLocalizations` delegate used throughout.

### 6.10.4 User Language Preferences

| Setting | Options | Default |
|---------|---------|--------|
| Display Language | English, Arabic | English |
| Number Format | Western (0–9), Eastern Arabic-Indic (٠–٩) | Western |
| Calendar System | Gregorian, Hijri (display only) | Gregorian |
| Report Language | English, Arabic, Bilingual | Inherits display language |
| AI Response Language | Auto (match query language), English, Arabic | Auto |

### 6.10.5 AI & LLM Multilingual Requirements

* **System Prompt Language Injection:** Every Genkit flow that produces user-visible text receives a mandatory `LANGUAGE RULE` instruction specifying the user's locale. The LLM must respond in the specified language.
* **SQL Remains English:** Generated SQL always uses English identifiers (column/table names) regardless of query language. The E-01 Language Detection agent normalises Arabic queries to intent before SQL generation.
* **Arabic-Capable Models:** GPT-4o and Claude 3.5 Sonnet are the primary models; both have strong Arabic reasoning capability. Local (Ollama) models require minimum 70B parameter size for adequate Arabic quality.
* **Bilingual Audit Notices:** Confidence warnings and compliance audit messages remain bilingual (`en` + active locale) for regulatory traceability.

### 6.10.6 Report & Document Localisation

* **PDF:** HTML/CSS pipeline with Arabic fonts (Cairo, Noto Sans Arabic — OFL licensed). `direction: rtl` and `unicode-bidi: bidi-override` applied to Arabic reports. Number and date formatting via `Intl` API at render time.
* **Excel:** `excelize` Go library with `RightToLeft: true` sheet view for Arabic workbooks. UTF-8 encoding throughout.
* **Currency:** Locale-aware `Intl.NumberFormat` with `style: 'currency'`; currency symbol position follows locale convention (symbol on right for Arabic).
* **Dates:** `Intl.DateTimeFormat` with Gregorian calendar default; Hijri calendar opt-in for display only (financial periods always Gregorian).

---

## 8. Data Governance & Quality

* **Data Validation & Reconciliation:**
  - Automated data quality checks during ETL.
  - Missing data and outlier detection.
  - Reconciliation between source systems (HIMS vs. Tally).
  - Data anomaly alerts and exception reporting.

* **Data Security & Compliance:**
  - Role-based encryption of sensitive fields (patient info, financial data).
  - Data masking for non-authorized users.
  - HIPAA/GDPR compliance for healthcare data.
  - Audit trail for all data access and modifications.
  - Data retention policies and archival.

---

## 9. Integration & Extensibility

* **Third-Party Integrations:**
  - Email service for scheduled report delivery.
  - Google Sheets/Excel for real-time data sync.
  - SMS gateway for alert notifications.
  - Slack/Teams integration for dashboard alerts.
  - Custom webhook support for external systems.

* **API Access:**
  - RESTful API for third-party applications to query data.
  - Webhook support for real-time event notifications.
  - Data export APIs for programmatic access.

---

## 10. Non-Functional Requirements (NFRs)

* **Security:** 
  - Read-only access to the Data Warehouse (the AI must *never* have write/delete permissions).
  - Data encryption at rest and in transit (HTTPS/TLS).
  - Two-factor authentication (2FA) for user accounts.
  - API key management and rate limiting.

* **Performance:** 
  - Complex joins across HIMS and Tally data should execute in under 10 seconds.
  - Dashboard load time <3 seconds for standard views.
  - Support concurrent queries from multiple users without degradation.
  - Caching mechanisms for frequently accessed reports.

* **Scalability:**
  - Database architecture supports horizontal scaling.
  - Support for growing data volumes (years of historical data).
  - Concurrent user capacity: minimum 50 simultaneous users.

* **Availability & Reliability:**
  - 99.5% uptime SLA.
  - Automated backup and disaster recovery (RPO <1 hour, RTO <4 hours).
  - Redundancy for critical components.

* **Hallucination Guardrails:** 
  - The AI must only answer data-related questions. If asked general knowledge questions (e.g., *"What is the capital of France?"*), it should gracefully pivot back to business analytics.
  - Confidence scoring for AI-generated responses.
  - Manual review queue for low-confidence responses.

* **Internationalisation (i18n) / Localisation:**
  - Default language is **English**. Arabic (RTL) is supported from Phase 1.
  - 100% of user-facing UI strings are externalised — no hardcoded English text in source code.
  - Translation coverage check enforced in CI: all `en` keys must have corresponding `ar` keys before merging.
  - AI responses must be delivered in the user's active locale; locale is injected into every Genkit flow as a mandatory system prompt instruction.
  - RTL layout must be validated by Playwright visual regression tests before each release.
  - Arabic PDF reports must render with correct RTL fonts; verified by automated PDF snapshot tests.
  - Date, number, and currency formatting must use `Intl.*` APIs — no manual string concatenation of locale-sensitive values.

---

## 11. User Stories

* **US1:** As an Owner, I want to click a button that says "Generate Monthly P&L" so that I instantly see my profit margins without waiting for the accountant.
* **US2:** As a Pharmacy Manager, I want to ask the chat, *"Which drugs expire next month?"* and get a downloadable table so I can process returns.
* **US3:** As an Accountant, I want to ask, *"Show me HIMS billing vs. Tally receipts for yesterday"* to easily reconcile missing payments.
* **US4:** As a Clinic Manager, I want to see a dashboard showing year-over-year revenue comparison by month with drill-down capability to see patient-wise breakdown.
* **US5:** As a Finance Head, I want to set up automated monthly reports that are emailed to stakeholders so I don't need to manually compile data.
* **US6:** As a Pharmacy Lead, I want inventory aging reports and low-stock alerts so I can proactively manage stock and reduce wastage.
* **US7:** As a Manager, I want to access analytics on mobile while traveling so I can make quick business decisions on the go.
* **US8:** As an Accountant, I want to see outstanding receivables aging and days sales outstanding (DSO) metrics to improve cash collection.

### AI Accountant Module User Stories

* **US9:** As an Accountant, I want to upload 50+ bank statements and vendor bills at once and have the system automatically extract invoice details (amount, date, vendor, GL account) so I can reduce manual data entry time from hours to minutes.
* **US10:** As an Accountant, I want the system to suggest the correct Tally ledger for each transaction with a confidence score, so I can quickly review and approve mappings instead of manually selecting them.
* **US11:** As a Finance Manager, I want a one-click sync button that pushes all approved transactions directly into Tally without manual journal entry creation, so that my books stay updated in real-time.
* **US12:** As an Accountant, I want to see which invoices are outstanding, which payments have cleared, and which items don't match between my bank statement and Tally, so I can reconcile accounts faster.
* **US13:** As an Accountant, I want the system to recognize handwritten invoice amounts and vendor names from scanned PDFs with 95%+ accuracy, so I don't need to manually retype information from physical documents.
* **US14:** As a Compliance Officer, I want complete audit logs showing who posted what transaction, when, and from which source document, so I can demonstrate compliance during tax audits.
* **US15:** As a Finance Head, I want to upload a new company instance in Tally and have it automatically sync statements and bills without requiring manual setup, so multi-entity accounting is effortless.
* **US16:** As an Accountant, I want to flag duplicate invoices automatically before posting them to Tally, so I prevent overpayment and maintain accurate books.

### Easy Reports Module User Stories

* **US17:** As a Finance Manager, I want to see a pre-built P&L statement that automatically pulls the latest Tally data, so I can review financial performance without manually creating reports.
* **US18:** As a Sales Director, I want to create a custom dashboard showing sales by region, top customers, and revenue trend without coding, so my team can monitor performance in real-time.
* **US19:** As an Operations Manager, I want to schedule a monthly inventory report to be emailed automatically every 1st of the month, so stakeholders get consistent insights on stock levels.
* **US20:** As an Accountant, I want to see consolidated P&L across multiple clinic branches in a single dashboard with drill-down to individual branch details, so I can analyze group performance.
* **US21:** As a Business Owner, I want my mobile dashboard to show key KPIs (revenue, profit, cash position, inventory value) with color-coded status alerts, so I can check business health on the go.
* **US22:** As a Pharmacy Manager, I want to analyze slow-moving and obsolete inventory with aging reports, so I can plan stock clearance and reduce waste.
* **US23:** As a Finance Head, I want to see budget vs. actual spending with variance alerts and forecast year-end actuals, so I can take corrective action early.
* **US24:** As a Cost Accountant, I want department-wise profitability analysis with cost allocation and overhead distribution, so I can evaluate each department's true contribution to profit.
* **US25:** As an HR/Finance Partner, I want to link external HR data (headcount, salaries) with Tally expense data to calculate cost-per-employee metrics, so I understand labor productivity.
* **US26:** As a Compliance Officer, I want automated GST reconciliation reports and audit trail logs for all transactions, so I'm audit-ready at any time.

### Advanced Search-Driven Analytics Module User Stories

* **US27:** As a Clinic Manager, I want to search "Show me which doctors have the highest patient satisfaction and lowest complaint rate" and get instant results with statistical confidence, so I can recognize top performers.
* **US28:** As an Owner, I want an autonomous agent to analyze my data every night and alert me if any KPI drops >15%, profits decline, or cash position becomes concerning, so I'm proactively aware of issues.
* **US29:** As a Finance Manager, I want to upload raw revenue data and have the system automatically design a complete dashboard with profit trends, regional breakdown, and forecasting, so I save hours on manual dashboard creation.
* **US30:** As an Accountant, I want the agent to suggest that "Clinic B's supply costs are 12% higher than Clinic A for identical items — recommend vendor renegotiation" with expected savings quantification, so I can act on recommendations.
* **US31:** As a Pharmacy Manager, I want to ask my mobile phone "Which drugs expire next month and what's my total loss?" and get instant answers with action recommendations, so I can manage stock on the go.
* **US32:** As a Data Analyst, I want to write Python code in a notebook to build a custom ML model predicting patient no-shows, and publish results so managers can query "Which appointments have >40% no-show risk?", so analytics are self-service.
* **US33:** As a Developer, I want natural language code generation: "Create a patient revenue dashboard embedded in our clinic app," so the system generates React components, database queries, and styling automatically.
* **US34:** As a Pharmacy Lead, I want a unified search across drugs, suppliers, costs, and sales to ask "Compare margin % for all drugs by supplier and identify the most profitable suppliers," so I understand profitability drivers.

### Multi-Language User Stories

* **US35:** As an Arabic-speaking Accountant, I want to ask my questions in Arabic (*"أعطني إيرادات هذا الشهر"*) and receive charts, table headers, and explanations fully in Arabic, so I can work confidently without needing to switch to English.
* **US36:** As a Manager, I want to switch the entire application from English to Arabic in one click on the Settings screen and have all screens, menus, notifications, and reports immediately reflect the change, so my whole team can use their preferred language.
* **US37:** As a Finance Head, I want to generate a PDF report that renders correctly in Arabic with right-to-left text, mirrored table columns, and Arabic-compatible fonts, so I can share professional reports with Arabic-speaking stakeholders.
* **US38:** As an Owner, I want all KPI alert emails and scheduled report delivery emails to arrive in my preferred language (Arabic), so I can read them instantly without translation.
* **US39:** As a Clinic Manager, I want to set my calendar display to the Hijri (Islamic) calendar while keeping financial period reports in Gregorian, so date references in our internal communications feel natural without breaking accounting compatibility.

---

## 12. Phased Timeline & Milestones

* **Phase 1: ETL & Infrastructure (Weeks 1–3)**
  - Set up Cloud Database.
  - Establish API/XML connections to Tally and HIMS.
  - Write scripts to sync data reliably without locking Tally.
  - Implement data validation and quality checks.
  - Set up data warehouse with optimized schema for analytics.

* **Phase 2: AI Agent & Core Analytics (Weeks 4–7)**
  - Integrate LangChain Text-to-SQL Agent.
  - Train the agent on the specific database schema and business terminology.
  - Implement data retrieval and JSON formatting for charts.
  - Develop financial analytics modules (P&L, cash flow, profitability).
  - Implement inventory and receivables/payables analytics.
  - Build drill-down and multi-period comparison features.

* **Phase 3: Dashboard & Advanced Features (Weeks 8–11)**
  - Build the chat interface and dashboard UI.
  - Integrate charting libraries and Export-to-CSV/PDF logic.
  - Implement scheduled reports and email delivery.
  - Build custom KPI dashboard creator.
  - Add mobile-responsive design.
  - Implement role-based access control (RBAC).
  - Add real-time alerts and notifications.

* **Phase 4: AI Accountant Module - Document Processing (Weeks 12–15)**
  - Implement drag-and-drop upload interface for documents.
  - Integrate OCR engine (cloud-based: AWS Textract, Google Vision, or Azure Document Intelligence).
  - Build OCR confidence scoring and error handling.
  - Develop document classification (invoices, bills, statements, receipts).
  - Implement handwritten text recognition (95%+ accuracy target).
  - Create bulk upload batch processing pipeline.

* **Phase 5: AI Accountant Module - Transaction Intelligence (Weeks 16–19)**
  - Build automated ledger mapping engine using AI/ML.
  - Implement confidence scoring for suggested mappings.
  - Create manual override and feedback loop for learning.
  - Develop vendor matching and duplicate detection algorithms.
  - Build sub-ledger and cost center assignment logic.
  - Create transaction review and approval workflow UI.

* **Phase 6: AI Accountant Module - Real-Time Integration (Weeks 20–22)**
  - Implement Tally API/TDL integration for journal entry posting.
  - Build real-time sync status dashboard with connection indicators.
  - Develop sync frequency configuration and scheduling.
  - Implement automatic retry logic and error recovery.
  - Create sync history and audit logging.
  - Add multi-entity Tally support.

* **Phase 7: AI Accountant Module - Reconciliation & Analytics (Weeks 23–26)**
  - Build bank-to-Tally matching algorithm.
  - Create reconciliation dashboard with outstanding items view.
  - Implement outstanding payables/receivables aging reports.
  - Develop DSO, DPO, and cash flow forecasting models.
  - Create expense categorization and trending analytics.
  - Build tax compliance and audit trail reporting.

* **Phase 8: Easy Reports Module - Pre-Built Reports & Dashboards (Weeks 27–30)**
  - Build pre-built report library (Sales, Profitability, Targets & Budgets, Financial Statements, Inventory, Revenue).
  - Implement consolidated multi-company reporting framework.
  - Develop interactive dashboard templates (Executive, Sales, Finance, Inventory, Operations).
  - Create drill-down and chart interaction capabilities.
  - Implement real-time data refresh mechanisms for dashboards.

* **Phase 9: Easy Reports Module - Automation & Distribution (Weeks 31–33)**
  - Build report scheduling engine (daily, weekly, monthly, custom intervals).
  - Implement email delivery system with template customization.
  - Create report portal and history archive.
  - Develop PDF/Excel export formatting and styling.
  - Build email distribution list management and role-based delivery.

* **Phase 10: Easy Reports Module - Customization & Advanced Features (Weeks 34–36)**
  - Implement zero-code drag-and-drop report builder for business users.
  - Build formula builder for calculated fields and custom metrics.
  - Develop role-based access control and data-level security (row/column level).
  - Implement custom fields (UDFs) support from Tally.
  - Create cost-center and department analytics module.

* **Phase 11: Integration & Polish - All Modules (Weeks 37–38)**
  - Integrate all three modules (BI Dashboard + AI Accountant + Easy Reports) into unified interface.
  - Implement cross-module data linking (e.g., document → transaction → report).
  - Create comprehensive document management and library system.
  - Performance optimization and caching across all modules.
  - Webhook and API support for external integrations.
  - Build mobile-responsive design for all modules.

* **Phase 12: UAT & Launch (Weeks 39–40)**
  - User Acceptance Testing (UAT) with accountants, managers, clinic staff, and IT.
  - Refine prompts, fix edge-case queries, OCR accuracy improvements.
  - Optimize report generation performance and dashboard load times.
  - Security audit, penetration testing, and compliance validation.
  - User training and comprehensive documentation (video tutorials, FAQs, admin guides, API docs).
  - Deploy to production with monitoring, alerting, and 24/7 support setup.
  - Post-launch optimization and user feedback incorporation.

* **Phase 13: Advanced Search Analytics - Semantic Layer & Infrastructure (Weeks 41–43)**
  - Design and implement semantic data model (dimensions, measures, hierarchies).
  - Build metadata layer across Tally and HIMS (relationship definitions, metric definitions).
  - Implement search-driven query engine (natural language to SQL translation layer).
  - Build search interface frontend with auto-complete, suggestions, and query history.
  - Integrate with existing data warehouse; optimize performance for sub-second search response.
  - Create semantic model administration tools for business users to define custom metrics.

* **Phase 14: Advanced Search Analytics - AI Agents & Autonomous Analysis (Weeks 44–46)**
  - Develop autonomous agent framework (Spotter-like agent) for self-service analytics.
  - Implement multi-step analytical workflows (data retrieval → analysis → insights → recommendations).
  - Build deep research agent for pattern discovery and statistical analysis.
  - Develop autonomous scheduling and alerting system for recurring analyses.
  - Create agent-driven workflow automation (insight → action in business systems).
  - Test agents with diverse use cases and refine recommendation accuracy.

* **Phase 15: Advanced Search Analytics - Auto-Dashboarding & Insights (Weeks 47–48)**
  - Develop AI dashboard auto-generation engine (intelligent chart selection, layout optimization).
  - Build autonomous insight discovery and anomaly detection.
  - Implement prescriptive recommendation engine with impact estimation.
  - Create natural language dashboard editing (modify charts and layouts with plain English).
  - Build mobile-optimized search and dashboard components.
  - Implement real-time dashboard refresh and data freshness management.

* **Phase 16: Advanced Search Analytics - Developer Tools & Embedding (Weeks 49–50)**
  - Develop natural language to code generation (SpotterCode-like feature).
  - Create REST/GraphQL APIs for embedded analytics.
  - Build Analyst Studio (Python notebooks, ad-hoc SQL, data mashups).
  - Develop IDE integration for code generation and live preview.
  - Create zero-copy in-memory query optimization for real-time performance.
  - Build federated query engine for multi-source data querying.

* **Phase 17: Advanced Search Analytics - Governance & Security (Weeks 51–52)**
  - Implement search-level access control and role-based query restrictions.
  - Build column-level and row-level data masking for search results.
  - Develop comprehensive audit logging for all searches and data access.
  - Implement sensitive question whitelisting and approval workflows.
  - Create compliance enforcement mechanisms (HIPAA/GDPR data protection).
  - Build user permission management and role creation tools.

* **Phase 18: Advanced Search Analytics - Polish & Integration (Weeks 53–54)**
  - Integrate search analytics module with existing BI Dashboard, AI Accountant, and Easy Reports modules.
  - Cross-module data linking (document → transaction → search result → action).
  - Performance optimization across all modules (caching, indexing, query optimization).
  - UAT for search module with diverse user personas and use cases.
  - Mobile optimization and responsive design refinement.
  - Final security audit and compliance validation.
  - Full integration UAT and production deployment.

---

## 13. Implementation Considerations & Key Features Summary

### Critical Success Factors

1. **Data Integration Robustness:** Reliable ETL pipelines are essential. Failed syncs or data inconsistencies can undermine trust in the analytics platform.

2. **AI Model Training:** The Text-to-SQL agent must be fine-tuned on your specific database schema and business terminology. This requires iterative testing and refinement during Phase 2.

3. **User Adoption:** Providing quick-action buttons and pre-built use cases ensures users can start deriving value immediately without needing to learn complex query syntax.

4. **Real-Time vs. Batch:** Decide early whether analytics need real-time updates (more cost) or batch updates (every 15–60 minutes). Most SME healthcare businesses can operate with hourly refreshes.

5. **Compliance & Data Privacy:** With healthcare and financial data involved, ensure HIPAA and GDPR compliance from day one. This affects database schema, backup strategies, and user access controls.

### Key Differentiators from a Standard Business Solution

**The Platform combines three tightly-integrated modules for complete business intelligence:**

1. **Conversational BI Dashboard:**
   - No need to learn complex query builders; users can ask in natural language.
   - Cross-domain analytics seamlessly correlate healthcare operations (HIMS) with financial data (Tally).
   - Pre-built domain logic enables the system to "understand" healthcare and accounting terminology.
   - Instant visualizations and drill-down capabilities for exploration.

2. **AI Accountant Module:**
   - Automates 75%+ of manual accounting tasks (data entry, categorization, reconciliation).
   - OCR-powered document processing handles handwritten invoices with 95%+ accuracy.
   - One-click sync keeps Tally data updated in real-time without manual journal entries.
   - Eliminates duplicate invoice posting and improves payment cycles through intelligent matching.
   - Provides complete audit trails for compliance and tax readiness.

3. **Easy Reports Module:**
   - Pre-built reports eliminate the need to create reports from scratch.
   - Zero-code drag-and-drop builder empowers business users to create custom reports without IT.
   - Consolidated multi-company dashboards provide enterprise-wide visibility at a glance.
   - Automated report scheduling removes the burden of manual report generation and distribution.
   - Role-based security and data filters ensure compliance and data governance.

**Unified Experience:** All three modules share consistent design, terminology, and user interface, reducing training time and increasing adoption.

### Risk Mitigation

| Risk | Mitigation Strategy |
| --- | --- |
| **Data Quality Issues** | Implement strict validation rules during ETL; automate data reconciliation checks; intelligent duplicate detection in AI Accountant. |
| **AI Hallucination or Incorrect Queries** | Implement confidence scoring; always show the generated SQL query to users; maintain a manual review queue for low-confidence responses. |
| **OCR Accuracy** | Start with cloud-based OCR (AWS/Google/Azure) for high accuracy; implement human-in-the-loop feedback to continuously improve model. |
| **Performance Degradation at Scale** | Use database indexing strategically; implement query caching; set rate limits on API endpoints; pre-compute rolled-up reports. |
| **Report Explosion** | Set naming conventions and folder hierarchy; implement version control; archive old report versions for audit purposes. |
| **User Resistance to Adoption** | Offer comprehensive training; start with pre-built dashboards and quick-action buttons; provide quick-win use cases; gather feedback early and iterate frequently. |
| **Security Breaches** | Enforce strict RBAC; encrypt sensitive data; conduct regular security audits; maintain audit trails for all data access; implement 2FA. |
| **Governance & Compliance** | Establish data ownership model; document all custom fields and metrics; maintain change logs; conduct quarterly compliance reviews. |

---

## 14. Major Risks & Inconsistencies

1. **The "Read-Only" Contradiction**
   - **The Conflict:** Section 10 (NFRs) mandates read-only access to protect financial data integrity. However, Section 6.7 (AI Accountant) requires the system to push transactions into Tally.
   - **Impact:** This creates a significant security and governance risk. Writing to a financial ledger via AI requires a different level of auditability and risk mitigation than a BI dashboard.

2. **Ambitious OCR Accuracy**
   - **The Risk:** Section 6.7.1 targets 95%+ accuracy for handwritten invoices.
   - **Reality Check:** Handwritten accounting documents vary significantly. Relying on autonomous posting without a robust Human-in-the-Loop (HITL) verification UI could lead to widespread ledger corruption.

3. **Scope Congestion (Product-Market Fit)**
   - **The Issue:** The PRD attempts to build a Conversational BI Dashboard, an Autonomous Accountant module, and a Low-Code Report Builder simultaneously.
   - **Impact:** This dilutes engineering focus and significantly increases the complexity of the "Unified Semantic Model".

---

## 15. Next Steps

1. **Finalize Product Name & Branding:** Consider unified brand identity for all three modules.
2. **Establish Project Governance:** Define decision-makers, approval workflows, and escalation paths for feature prioritization.
3. **Plan Stakeholder Engagement:** Schedule kickoff meetings with clinic manager, pharmacy lead, accountant, and IT to validate requirements and use cases.
4. **Detail Database Schema:** Work with Tally and HIMS administrators to finalize exact fields, tables, and extraction frequencies; add locale columns to `users`, `audit_log`, and `scheduled_reports` tables.
5. **Prototype AI Prompts & Report Designs:** Draft system prompts for the LLM (including `ResponseLanguageInstruction`) and mock up key reports/dashboards with sample data in both English and Arabic.
6. **OCR & Vendor Selection:** Evaluate OCR engines (AWS Textract, Google Vision, Azure Form Recognizer) for **Arabic handwriting accuracy** in addition to standard printed documents.
7. **Set Up Development Environment:** Provision cloud infrastructure, databases, version control, and CI/CD pipelines.
8. **Hire/Allocate Resources:** Assign full-stack engineers, ML/AI engineers, QA engineers, product designer/UX specialist, and a **native Arabic-speaking QA reviewer** for translation and RTL validation.
9. **Create Detailed Technical Specifications:** Expand this PRD into API specifications, data schemas, UI wireframes, and test plans; include RTL wireframes for all key screens.
10. **Risk Assessment & Insurance:** Identify critical dependencies and mitigation strategies; assess insurance/liability needs.
11. **Seed Translation Files:** Create initial `en` and `ar` JSON/ARB translation files for all namespaces; engage professional translator for medical and accounting domain terms (see `docs/i18n-architecture.md §11`).
12. **RTL QA Baseline:** Establish Playwright RTL visual regression baselines for the 10 most-used screens before sprint 3 delivery.

---

**Document Version:** 1.2 | **Last Updated:** February 19, 2026 | **Owner:** Product Management Team

