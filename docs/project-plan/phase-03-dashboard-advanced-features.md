# Phase 03 — Dashboard, Chat UI & Advanced Features

**Phase Duration:** Weeks 8–11 (4 weeks)  
**Module(s):** Module A (continued), Module E (E-07)  
**Status:** Planning  
**Milestone:** M3 — Chat Dashboard MVP  
**Depends On:** Phase 02 complete (AI agents operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §4.6](../ARCHITECTURE.md) | [DESIGN.md](../DESIGN.md)

---

## 1. Objectives

Build and ship the full user-facing product: the React web chat dashboard, the Flutter mobile app foundation, the pinnable dashboard system, scheduled reports, KPI alerts, and drill-down analytics. This phase delivers **Milestone M3 — Chat Dashboard MVP** — the first version real users can interact with.

---

## 2. Scope

### In Scope
- React web application (CopilotKit + Vite)
- Chat interface with streaming AI responses and inline charts
- Quick-action prompt carousel (10 pre-built prompts)
- Apache ECharts dynamic chart rendering
- Pin-to-dashboard functionality
- Auto-refreshing dashboard grid
- Export to CSV / Excel / PDF
- Multi-period comparison queries (A-08)
- Drill-down on chart elements (A-07)
- KPI alert agent (A-10) with in-app + email notifications
- Scheduled reports agent (A-09) with email delivery
- Chart-to-Dashboard Pin agent (A-11)
- RTL Arabic layout (initial implementation)
- Flutter mobile app (dashboard viewer + chat — read-only)
- E-07 Bilingual Glossary Sync agent
- Playwright RTL visual regression baseline

### Out of Scope
- Document upload / AI Accountant (Phase 4+)
- Zero-code report builder (Phase 10)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | React App (CopilotKit + Vite) | Frontend Engineer | Builds and deploys; CopilotKit provider wraps app |
| D-02 | Chat Interface Component | Frontend Engineer | Streaming SSE responses render; Arabic RTL chat bubbles correct direction |
| D-03 | Quick-Action Prompt Carousel | Frontend Engineer | 10 configurable prompt buttons inject query into chat on click |
| D-04 | DynamicChart Component | Frontend Engineer | Renders bar/line/pie/scatter/table from agent JSON payload; chart type auto-selected by A-03 |
| D-05 | KPI Card Component | Frontend Engineer | Sparkline + value + trend indicator; RTL-safe layout |
| D-06 | DrillDownTable Component | Frontend Engineer | Expandable paginated table; clicking row cell triggers A-07 drill-down query |
| D-07 | Pin-to-Dashboard (A-11) | Frontend + AI Engineer | Pin icon on chat chart → saved to `app.pinned_charts`; appears on dashboard grid |
| D-08 | Dashboard Grid (auto-refresh) | Frontend Engineer | Pinned charts render in responsive grid; refresh on configurable interval |
| D-09 | Export Button (CSV/Excel/PDF) | Backend + Frontend | Every table/chart has download btn; PDF is Arabic-RTL correct via WeasyPrint (E-04 prep) |
| D-10 | A-07 Drill-Down Context Agent | AI Engineer | Clicking chart element generates drill-down SQL and returns nested detail view |
| D-11 | A-08 Multi-Period Comparison | AI Engineer | "Compare this month vs last month" → side-by-side bar correctly across all KPIs |
| D-12 | A-09 Report Scheduling Agent | AI Engineer | Admin can create cron-based report schedule; report generated + emailed on trigger |
| D-13 | A-10 KPI Alert Agent | AI Engineer | Configurable threshold alerts for any metric; fires email + in-app notification |
| D-14 | A-11 Chart-to-Dashboard Pin Agent | AI Engineer | Pin action persisted; auto-refresh schedule stored; works across sessions |
| D-15 | E-07 Bilingual Glossary Sync | AI Engineer | Glossary terms pushed to A-04 Domain Normaliser on deploy; reviewed by medical + finance advisor |
| D-16 | Flutter Mobile App (MVP) | Frontend Engineer | Dashboard viewer, chat query input, offline pinned chart view; iOS + Android builds |
| D-17 | Playwright RTL Baselines | QA Engineer | Visual regression snapshots for all 10 primary screens in both EN and AR |

---

## 4. AI Agents Deployed

### A-07 Drill-Down Context Agent

**Trigger:** User clicks on a chart element (e.g., a bar for "March Revenue")  
**Input:** `{ metric, dimension_value, parent_query, period }`  
**Action:** Generates drill-down SQL: "Show individual transactions for clinic revenue in March 2026"  
**Output:** DrillDownTable component with transaction-level detail

**Drill-down hierarchy:**
```
Total Revenue
  └── By Department (Clinic / Pharmacy)
        └── By Doctor / Drug Category
              └── Individual Transactions
```

### A-08 Multi-Period Comparison Agent

**Trigger:** User query contains temporal comparison intent ("vs last month", "year-over-year", "compare Q1 Q2")  
**Action:** Decomposes query into 2 parallel date-range SQL queries; returns side-by-side bar chart  
**Supported comparisons:**
- Month-over-month (MoM)
- Year-over-year (YoY)
- Quarter-over-quarter (QoQ)
- Any two custom date ranges

**Auto enrichment:** Calculates % change and absolute delta; annotates chart

### A-09 Report Scheduling Agent

**Type:** Proactive (cron-triggered Genkit flow)  
**Configuration stored in:** `app.scheduled_reports`  
**Execution flow:**
```
Cron trigger (e.g., every Monday 7AM)
    │
    ▼ Load report config (type, params, locale, recipient list)
    │
    ▼ Execute A-01 with report query params
    │
    ▼ Render: PDF (WeasyPrint) | Excel (excelize) | CSV
    │
    ▼ Email via Notification Dispatcher (SMTP)
    │
    ▼ Log to audit_log; update scheduled_reports.last_run_at
```

**Supports:** Daily / Weekly / Monthly / Quarterly / Custom cron  
**Report formats:** PDF, Excel (.xlsx), CSV  
**Locale:** Each schedule stores its own locale; Arabic schedules produce AR PDFs

### A-10 KPI Alert Agent

**Type:** Proactive (runs on NATS `etl.sync.completed` event + time-based schedule)  
**Configuration:** Users define alert rules: `if metric X crosses threshold Y, notify via Z`  
**Example rules:**
- `clinic_revenue < 50000` → email + SMS to clinic_manager
- `outstanding_receivables > 200000` → in-app + email to finance_head
- `pharmacy_stock_below_reorder` → in-app to pharmacy_manager

**Notification channels:** In-app, Email (SMTP), SMS (gateway integration — stub in Phase 3), Slack (webhook — stub in Phase 3)

### A-11 Chart-to-Dashboard Pin Agent

**Trigger:** User clicks "Pin" icon on any chat-rendered chart  
**Action (HITL: No — fully automatic):**
1. Serialises ECharts config JSON + SQL query to `app.pinned_charts`
2. Associates refresh interval (default: 15 minutes)
3. Returns confirmation with link to dashboard

**Auto-refresh:** Pinned charts store the underlying SQL; on dashboard load, SQL re-executed against fresh warehouse data

### E-07 Bilingual Glossary Sync Agent

**Type:** Proactive (triggered on deploy pipeline)  
**Source:** `docs/i18n-glossary.md` — canonical bilingual domain glossary  
**Action:** On each deploy, extracts all healthcare + accounting synonym pairs and pushes to A-04's synonym map (stored in `app.domain_glossary` table)  
**HITL:** Glossary term changes require sign-off from Medical Advisor + Finance Advisor before merge

---

## 5. Frontend Architecture

### React Application Structure
```
frontend/
├── src/
│   ├── providers/
│   │   ├── CopilotKitProvider.tsx    ← Generative UI orchestration
│   │   ├── I18nProvider.tsx          ← i18next + RTL switching
│   │   └── AuthProvider.tsx          ← Keycloak JWT + role claims
│   │
│   ├── layouts/
│   │   ├── DashboardLayout.tsx       ← Pinned charts grid
│   │   ├── ChatLayout.tsx            ← Conversational interface
│   │   └── ReportsLayout.tsx         ← Scheduled reports list
│   │
│   ├── components/
│   │   ├── chat/
│   │   │   ├── ChatWindow.tsx        ← SSE streaming chat
│   │   │   ├── PromptCarousel.tsx    ← 10 quick-action prompts
│   │   │   └── MessageBubble.tsx     ← RTL-aware chat bubble
│   │   ├── charts/
│   │   │   ├── DynamicChart.tsx      ← ECharts wrapper (auto type)
│   │   │   ├── KPICard.tsx           ← Sparkline KPI card
│   │   │   └── DrillDownTable.tsx    ← Paginated drill-down table
│   │   ├── dashboard/
│   │   │   ├── DashboardGrid.tsx     ← Responsive pinned chart grid
│   │   │   └── PinButton.tsx         ← Pin icon + action
│   │   └── alerts/
│   │       └── AlertConfigPanel.tsx  ← KPI alert rule builder
│   │
│   ├── hooks/
│   │   ├── useCopilotAction.ts       ← Dynamic widget rendering
│   │   ├── useStreamingChat.ts       ← SSE response handling
│   │   └── useLocale.ts              ← Locale + RTL state
│   │
│   └── state/
│       ├── appStore.ts               ← Zustand global state
│       └── chartStore.ts             ← Pinned chart state
│
├── public/locales/
│   ├── en/ {common, dashboard, chat, reports, alerts}.json
│   └── ar/ {common, dashboard, chat, reports, alerts}.json
```

### RTL Implementation Rules
- All padding/margin use Tailwind logical properties: `ms-`, `me-`, `ps-`, `pe-`
- No hardcoded `left`/`right` in any component
- `document.documentElement.dir = locale === 'ar' ? 'rtl' : 'ltr'` on locale change
- Chat bubbles: user messages `justify-start` (RTL left = visual right), AI messages `justify-end`
- Navigation sidebar mirrors position via CSS `inset-inline-start`
- Directional icons (chevrons, arrows): `rtl:scale-x-[-1]` CSS utility

### Quick-Action Prompts (Default 10)
| # | English | Arabic |
|---|---|---|
| 1 | Today's Total Revenue | إجمالي إيرادات اليوم |
| 2 | Pending Tally Invoices | فواتير تالي المعلقة |
| 3 | Low Pharmacy Stock | مخزون الصيدلية المنخفض |
| 4 | Top 5 Selling Drugs This Week | أفضل 5 أدوية مبيعًا هذا الأسبوع |
| 5 | Outstanding Receivables | الذمم المدينة المستحقة |
| 6 | Doctor-wise Patient Visits (Month) | زيارات المرضى حسب الطبيب |
| 7 | Monthly P&L Summary | ملخص الأرباح والخسائر الشهري |
| 8 | Cash Flow This Week | التدفق النقدي هذا الأسبوع |
| 9 | Appointments for Today | مواعيد اليوم |
| 10 | Pharmacy Revenue vs Clinic Revenue | إيرادات الصيدلية مقابل العيادة |

---

## 6. Flutter Mobile App (MVP)

```
mobile/
├── lib/
│   ├── main.dart
│   ├── routing/              ← go_router routes
│   ├── features/
│   │   ├── dashboard/        ← Pinned chart viewer (offline-first)
│   │   ├── chat/             ← Chat query input + response
│   │   └── auth/             ← Keycloak login
│   ├── core/
│   │   ├── sync/             ← PowerSync offline-first
│   │   ├── i18n/             ← app_en.arb + app_ar.arb
│   │   └── theme/            ← Directionality + RTL
│   └── widgets/
│       ├── DynamicChart.dart  ← fl_chart wrapper
│       └── KPICard.dart
```

**Offline capability:** PowerSync syncs last-5 pinned dashboard configs to local SQLite; viewable without internet  
**Charts:** fl_chart for mobile-optimised visualisations  
**Auth:** Keycloak OIDC flow via `flutter_appauth`

---

## 7. Notification Infrastructure

**Email:** SMTP integration (configurable: SendGrid, AWS SES, or on-premises Postfix)  
**In-app:** WebSocket push to connected clients when `alert.kpi.threshold` NATS event fires  
**SMS:** Stub interface in Phase 3; plugged into real gateway in Phase 6 (E-06)  
**Notification table:** `app.notification_queue` — all notifications logged with delivery status

---

## 8. Export Functionality

| Format | Technology | Arabic Support |
|---|---|---|
| CSV | Go `encoding/csv` | UTF-8 encoding |
| Excel (.xlsx) | `excelize` Go library | RTL sheet direction for Arabic |
| PDF | WeasyPrint + Cairo/Noto Sans Arabic | Full RTL, Arabic-compatible fonts |

---

## 9. Testing Requirements

| Test Type | Scope | Tool |
|---|---|---|
| Component tests | All React components (unit) | Vitest + React Testing Library |
| E2E tests | Full chat query → chart render flow | Playwright |
| RTL visual regression | All 10 primary screens in Arabic | Playwright screenshots |
| A-07 drill-down accuracy | 20 drill-down click scenarios | Go agent test suite |
| A-08 comparison queries | 15 period comparison scenarios | Go agent test suite |
| A-09 scheduled reports | Email delivery of PDF/Excel for 3 report types | Integration test |
| A-10 KPI alerts | Threshold breach → notification delivery | Integration test |
| Performance | Dashboard load < 3 seconds | Lighthouse CI |
| Flutter tests | Widget tests for chat and dashboard | Flutter test framework |

---

## 10. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| CopilotKit SSE streaming latency causing UI jank | Medium | Optimistic UI updates; loading skeleton states |
| Arabic RTL layout errors in complex chart components | Medium | Playwright RTL snapshots for every chart type; dedicated Arabic QA reviewer |
| PDF Arabic font rendering (WeasyPrint) | Medium | Test with Cairo + Noto Sans Arabic fonts early; have fallback font set |
| A-09 scheduled report email delivery failures | Low | SMTP retry logic; fallback delivery log; bounce handling |

---

## 11. Phase Exit Criteria

- [ ] Chat interface renders streaming AI responses with inline charts in both EN and AR
- [ ] All 5 new Module A agents (A-07 to A-11) deployed and tested
- [ ] E-07 Bilingual Glossary Sync agent deployed and connected to A-04
- [ ] Dashboard grid with pinned charts and auto-refresh working
- [ ] Export (CSV, Excel, PDF) working — Arabic PDF renders correctly
- [ ] Playwright RTL visual regression baseline captured for all 10 screens
- [ ] Flutter mobile MVP builds for iOS and Android; offline dashboard viewing works
- [ ] KPI alert rules configurable and firing notifications correctly
- [ ] Scheduled reports generating and emailing on time
- [ ] Phase gate reviewed and signed off; stakeholder demo completed

---

*Phase 03 | Version 1.0 | February 19, 2026*
