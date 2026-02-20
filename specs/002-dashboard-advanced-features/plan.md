# Implementation Plan: Dashboard, Chat UI & Advanced Features with i18n

**Branch**: `001-dashboard-advanced-features` | **Date**: 2026-02-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-dashboard-advanced-features/spec.md`

## Summary

This feature delivers the complete user-facing product for MediSync: a conversational business intelligence dashboard enabling healthcare and pharmacy staff to query operational and financial data using natural language in both English and Arabic. Key capabilities include:

- **Natural Language Queries** with streaming AI responses and inline chart/table rendering
- **Dashboard Pinning** for persistent monitoring with auto-refresh
- **Drill-Down Analysis** for multi-level data exploration
- **KPI Alerts** via in-app and email notifications
- **Scheduled Reports** in PDF, spreadsheet, and CSV formats
- **Multi-Period Comparison** with delta annotations
- **Full i18n Support** for English (LTR) and Arabic (RTL)
- **Mobile Dashboard** with offline capability

Technical approach: React 19 frontend with CopilotKit for generative UI, Apache ECharts for visualizations, i18next for internationalization; Go 1.26 backend with existing AI agents (Module A: Text-to-SQL, Visualization Routing); PostgreSQL for persistence; Redis for session caching.

---

## Technical Context

**Language/Version**:
- Backend: Go 1.26
- Frontend: TypeScript 5.x with React 19.2.4

**Primary Dependencies**:
- **Backend**: go-chi/chi v5.2.1, pgx/v5 v5.7.4, nats.go v1.39.1, go-redis/v9 v9.4.0, pgvector-go v0.2.2
- **AI Orchestration**: Google Genkit, Agent ADK (existing Module A agents)
- **Frontend**: @copilotkit/react-core v1.3.6, @copilotkit/react-ui v1.3.6, echarts v5.6.0, echarts-for-react v3.0.2, i18next v24.2.2, react-i18next v15.4.0
- **Build**: Vite 7.3, Tailwind CSS 3.4.17

**Storage**:
- PostgreSQL 18.2 + pgvector (pinned charts, alert rules, scheduled reports, user preferences)
- Redis (session cache, real-time subscriptions)
- N/A for mobile offline: PowerSync local cache (out of scope for this phase)

**Testing**:
- Backend: Go `testing` package + `testify`
- Frontend: Vitest + React Testing Library
- E2E: Contract tests in `tests/contract/`, integration tests in `tests/integration/`

**Target Platform**:
- Backend: Linux server (Docker containers)
- Frontend: Modern browsers (Chrome 90+, Firefox 88+, Safari 14+, Edge 90+)
- Mobile: iOS 12+, Android 5+ (Flutter 3.42) - *out of scope for this phase, covered in separate mobile feature*

**Project Type**: Web application (frontend + backend) with future mobile support

**Performance Goals**:
- Query end-to-end latency (P95): < 5 seconds
- Dashboard load time: < 3 seconds
- Language switching: < 1 second (no page reload)
- Streaming response first token: < 500ms
- Export operations (10k rows): < 30 seconds

**Constraints**:
- All AI queries via `medisync_readonly` database role (SELECT-only)
- No AI can autonomously write to source systems
- PII must be masked based on user role via OPA policies
- All scheduled report times in organization's local timezone
- Arabic PDF exports require RTL-compatible fonts

**Scale/Scope**:
- Initial deployment: ~100 concurrent users per organization
- 50 screens/components across web dashboard
- Support for 10 configurable quick-action prompts
- Up to 20 pinned charts per user dashboard
- Real-time alert delivery within 2 minutes of threshold crossing

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Security First & HITL Gates ✅

| Requirement | Status | Notes |
|-------------|--------|-------|
| AI agents use `medisync_readonly` DB role | ✅ PASS | All chat queries use existing Module A agents which enforce read-only access |
| SQL queries validated as SELECT-only | ✅ PASS | Existing `a01_text_to_sql` validation layer already enforces this |
| No AI autonomous writes to Tally | ✅ PASS | This feature has no write-back operations; Tally sync is Phase 04+ |
| Column-level PII masking | ✅ PASS | OPA policies in `policies/column_masking.rego` applied via `internal/auth/opa.go` |
| TLS 1.3+ for external API calls | ✅ PASS | Standard Go HTTP client with TLS config |
| Sensitive config via env vars | ✅ PASS | Using `internal/config/config.go` pattern |

### II. Read-Only Intelligence Plane ✅

| Requirement | Status | Notes |
|-------------|--------|-------|
| AI queries via `medisync_readonly` | ✅ PASS | Existing agent infrastructure enforces this |
| OPA blocks DML at driver level | ✅ PASS | Policies in `policies/bi_read_only.rego` |
| No warehouse writes by AI | ✅ PASS | Feature only reads for display/analysis |

### III. i18n by Default ✅

| Requirement | Status | Notes |
|-------------|--------|-------|
| All UI strings via i18next | ⚠️ ATTENTION | New components must use `useTranslation()` hook - verify in implementation |
| Tailwind logical properties for RTL | ⚠️ ATTENTION | Use `ms-`, `me-`, `ps-`, `pe-` instead of `ml-`, `mr-`, `pl-`, `pr-` |
| AI responses include `ResponseLanguageInstruction` | ✅ PASS | Existing Module E agents handle this |
| CI blocks on translation gaps | ⚠️ ATTENTION | Must add i18n key validation to CI pipeline |
| Locale-aware date/number/currency formatting | ⚠️ ATTENTION | Use `Intl` APIs with user's locale |

### IV. Open Source Only ✅

| Requirement | Status | Notes |
|-------------|--------|-------|
| All dependencies OSI-approved | ✅ PASS | React (MIT), ECharts (Apache-2.0), i18next (MIT), CopilotKit (MIT), Go stdlib (BSD-3) |
| No GPL/AGPL in production | ✅ PASS | All selected licenses are permissive |

### V. Test-Driven Development ✅

| Requirement | Status | Notes |
|-------------|--------|-------|
| Unit tests for business logic | ⚠️ ATTENTION | New services require test coverage |
| Integration tests for agent flows | ⚠️ ATTENTION | Chat flow tests needed in `tests/integration/` |
| E2E tests for critical journeys | ⚠️ ATTENTION | Query → Chart → Pin flow needs E2E test |
| Mock external dependencies | ✅ PASS | Pattern exists in `tests/agents/module_a/` |

### Gate Summary

**Status**: ✅ PROCEED TO PHASE 0

All core principles can be satisfied. Items marked ⚠️ ATTENTION require verification during implementation but do not block planning. No constitution violations requiring justification.

---

## Project Structure

### Documentation (this feature)

```text
specs/002-dashboard-advanced-features/
├── plan.md              # This file
├── research.md          # Phase 0 output - technology decisions
├── data-model.md        # Phase 1 output - entity definitions
├── quickstart.md        # Phase 1 output - developer setup guide
├── contracts/           # Phase 1 output - API contracts
│   ├── chat-api.yaml
│   ├── dashboard-api.yaml
│   ├── alerts-api.yaml
│   └── reports-api.yaml
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
# Backend (Go - existing structure, new files added)
internal/
├── api/
│   ├── handlers/
│   │   ├── chat.go           # NEW: Chat query endpoints
│   │   ├── dashboard.go      # NEW: Pinned chart CRUD
│   │   ├── alerts.go         # NEW: Alert rule management
│   │   ├── reports.go        # NEW: Scheduled report config
│   │   └── export.go         # NEW: Data export endpoints
│   ├── middleware/
│   │   └── locale.go         # NEW: Locale extraction from JWT/header
│   └── websocket/
│       └── stream.go         # NEW: Streaming response handler
├── warehouse/
│   ├── pinned_chart.go       # NEW: Pinned chart repository
│   ├── alert_rule.go         # NEW: Alert rule repository
│   ├── scheduled_report.go   # NEW: Scheduled report repository
│   └── user_preference.go    # NEW: User preferences repository
├── agents/
│   └── module_a/             # EXISTING: Text-to-SQL, Visualization
└── services/
    ├── alert_scheduler.go    # NEW: Alert evaluation scheduler
    ├── report_generator.go   # NEW: Report generation service
    └── notification.go       # NEW: Email/in-app notification sender

migrations/
├── 010_pinned_charts.up.sql   # NEW
├── 011_alert_rules.up.sql     # NEW
├── 012_scheduled_reports.up.sql # NEW
└── 013_user_preferences.up.sql # NEW

# Frontend (React/TypeScript - existing structure, new files added)
frontend/src/
├── components/
│   ├── chat/
│   │   ├── ChatInterface.tsx    # NEW: Main chat container
│   │   ├── MessageList.tsx      # NEW: Message display
│   │   ├── QueryInput.tsx       # NEW: Natural language input
│   │   ├── StreamingMessage.tsx # NEW: Progressive rendering
│   │   └── QuickActions.tsx     # NEW: Quick-action prompts
│   ├── charts/
│   │   ├── ChartRenderer.tsx    # NEW: ECharts wrapper
│   │   ├── ChartActions.tsx     # NEW: Pin, Export, Drill-down
│   │   └── DrillDownModal.tsx   # NEW: Drill-down detail view
│   ├── dashboard/
│   │   ├── DashboardGrid.tsx    # NEW: Pinned charts grid
│   │   ├── DashboardWidget.tsx  # NEW: Single pinned chart
│   │   └── WidgetSettings.tsx   # NEW: Refresh/config options
│   ├── alerts/
│   │   ├── AlertList.tsx        # NEW: Alert rules management
│   │   ├── AlertForm.tsx        # NEW: Create/edit alert
│   │   └── NotificationToast.tsx # NEW: In-app notifications
│   ├── reports/
│   │   ├── ReportScheduler.tsx  # NEW: Schedule configuration
│   │   └── ReportHistory.tsx    # NEW: Past report deliveries
│   └── common/
│       ├── LanguageSwitcher.tsx # NEW: EN/AR toggle
│       ├── ExportButton.tsx     # NEW: Format selection + download
│       └── LoadingSpinner.tsx   # NEW: RTL-aware loading states
├── pages/
│   ├── ChatPage.tsx           # NEW: Chat interface route
│   ├── DashboardPage.tsx      # NEW: Dashboard view route
│   ├── AlertsPage.tsx         # NEW: Alerts management route
│   └── ReportsPage.tsx        # NEW: Reports management route
├── hooks/
│   ├── useChat.ts             # NEW: Chat state + streaming
│   ├── useDashboard.ts        # NEW: Pinned charts CRUD
│   ├── useAlerts.ts           # NEW: Alert rules management
│   ├── useLocale.ts           # NEW: Locale detection/switching
│   └── useExport.ts           # NEW: Export functionality
├── services/
│   ├── api.ts                 # NEW: API client with auth
│   ├── websocket.ts           # NEW: Streaming connection
│   └── export.ts              # NEW: File generation
├── i18n/
│   ├── index.ts               # EXISTING: i18next config
│   └── locales/
│       ├── en/
│       │   ├── chat.json      # NEW
│       │   ├── dashboard.json # NEW
│       │   ├── alerts.json    # NEW
│       │   └── reports.json   # NEW
│       └── ar/
│           ├── chat.json      # NEW
│           ├── dashboard.json # NEW
│           ├── alerts.json    # NEW
│           └── reports.json   # NEW
└── styles/
    └── globals.css            # EXISTING: Tailwind + RTL utilities

# Tests
tests/
├── contract/
│   ├── chat_api_test.go       # NEW
│   ├── dashboard_api_test.go  # NEW
│   └── alerts_api_test.go     # NEW
├── integration/
│   ├── chat_flow_test.go      # EXISTING: Extend for new features
│   └── alert_notification_test.go # NEW
└── e2e/
    └── dashboard_pin_flow.spec.ts # NEW: Playwright E2E test
```

**Structure Decision**: Using existing web application structure (Option 2 from template). Backend extends `internal/` with new handlers, services, and warehouse repositories. Frontend extends `frontend/src/` with new component directories organized by feature domain. No mobile code in this phase (Flutter app in separate feature).

---

## Complexity Tracking

> No constitution violations detected. Table empty.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | - | - |
