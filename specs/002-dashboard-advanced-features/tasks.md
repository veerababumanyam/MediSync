# Tasks: Dashboard, Chat UI & Advanced Features with i18n

**Input**: Design documents from `/specs/002-dashboard-advanced-features/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are included per Constitution Principle V (TDD). Write tests first, ensure they fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US10)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `internal/` at repository root (Go)
- **Frontend**: `frontend/src/` (React/TypeScript)
- **Migrations**: `migrations/` at repository root
- **Tests**: `tests/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for dashboard feature

- [ ] T001 Create database migrations for dashboard tables in migrations/010_user_preferences.up.sql
- [ ] T002 [P] Create database migrations for pinned charts in migrations/011_pinned_charts.up.sql
- [ ] T003 [P] Create database migrations for alert rules in migrations/012_alert_rules.up.sql
- [ ] T004 [P] Create database migrations for scheduled reports in migrations/013_scheduled_reports.up.sql
- [ ] T005 [P] Create database migrations for chat messages in migrations/014_chat_messages.up.sql
- [ ] T006 [P] Add English translations for chat feature in frontend/src/i18n/locales/en/chat.json
- [ ] T007 [P] Add Arabic translations for chat feature in frontend/src/i18n/locales/ar/chat.json
- [ ] T008 [P] Add English translations for dashboard in frontend/src/i18n/locales/en/dashboard.json
- [ ] T009 [P] Add Arabic translations for dashboard in frontend/src/i18n/locales/ar/dashboard.json
- [ ] T010 [P] Add English translations for alerts in frontend/src/i18n/locales/en/alerts.json
- [ ] T011 [P] Add Arabic translations for alerts in frontend/src/i18n/locales/ar/alerts.json
- [ ] T012 [P] Add English translations for reports in frontend/src/i18n/locales/en/reports.json
- [ ] T013 [P] Add Arabic translations for reports in frontend/src/i18n/locales/ar/reports.json

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Database Repositories

- [ ] T014 [P] Create UserPreference repository in internal/warehouse/user_preference.go
- [ ] T015 [P] Create ChatMessage repository in internal/warehouse/chat_message.go
- [ ] T016 [P] Create PinnedChart repository in internal/warehouse/pinned_chart.go
- [ ] T017 [P] Create AlertRule repository in internal/warehouse/alert_rule.go
- [ ] T018 [P] Create Notification repository in internal/warehouse/notification.go
- [ ] T019 [P] Create ScheduledReport repository in internal/warehouse/scheduled_report.go

### API Infrastructure

- [ ] T020 [P] Create locale extraction middleware in internal/api/middleware/locale.go
- [ ] T021 [P] Create WebSocket upgrader and stream handler in internal/api/websocket/stream.go
- [ ] T022 [P] Create API client service for frontend in frontend/src/services/api.ts
- [ ] T023 [P] Create WebSocket client for streaming in frontend/src/services/websocket.ts
- [ ] T024 [P] Create locale hook for React in frontend/src/hooks/useLocale.ts

### Common UI Components

- [ ] T025 [P] Create LanguageSwitcher component in frontend/src/components/common/LanguageSwitcher.tsx
- [ ] T026 [P] Create LoadingSpinner component with RTL support in frontend/src/components/common/LoadingSpinner.tsx
- [ ] T027 [P] Create ExportButton component in frontend/src/components/common/ExportButton.tsx
- [ ] T028 Configure i18next for new namespaces in frontend/src/i18n/index.ts

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Natural Language Data Query (Priority: P1) 🎯 MVP

**Goal**: Enable users to ask questions in natural language (EN/AR) and receive visual responses with charts/tables

**Independent Test**: Ask "What is today's total revenue?" in English and "ما هي إيرادات اليوم؟" in Arabic, verify accurate chart/table responses in correct language and layout direction

### Tests for User Story 1

- [ ] T029 [P] [US1] Contract test for POST /chat/query in tests/contract/chat_api_test.go
- [ ] T030 [P] [US1] Contract test for WebSocket /chat/stream in tests/contract/chat_stream_test.go
- [ ] T031 [P] [US1] Integration test for English query flow in tests/integration/chat_flow_test.go
- [ ] T032 [P] [US1] Integration test for Arabic query flow in tests/integration/chat_arabic_test.go

### Backend Implementation for User Story 1

- [ ] T033 [US1] Create ChatHandler with SubmitQuery method in internal/api/handlers/chat.go
- [ ] T034 [US1] Create StreamHandler for WebSocket connections in internal/api/websocket/stream.go (extend T021)
- [ ] T035 [US1] Implement chat query validation and sanitization in internal/api/handlers/chat.go
- [ ] T036 [US1] Integrate with existing Module A agents (Text-to-SQL, Visualization) in internal/api/handlers/chat.go
- [ ] T037 [US1] Add locale-aware response formatting in internal/api/handlers/chat.go
- [ ] T038 [US1] Implement streaming response chunks (text, chart, table, error) in internal/api/websocket/stream.go
- [ ] T039 [US1] Register chat routes in internal/api/routes.go

### Frontend Implementation for User Story 1

- [ ] T040 [P] [US1] Create useChat hook with streaming support in frontend/src/hooks/useChat.ts
- [ ] T041 [P] [US1] Create ChatInterface component in frontend/src/components/chat/ChatInterface.tsx
- [ ] T042 [P] [US1] Create MessageList component in frontend/src/components/chat/MessageList.tsx
- [ ] T043 [P] [US1] Create QueryInput component with RTL support in frontend/src/components/chat/QueryInput.tsx
- [ ] T044 [P] [US1] Create StreamingMessage component for progressive rendering in frontend/src/components/chat/StreamingMessage.tsx
- [ ] T045 [P] [US1] Create ChartRenderer component with ECharts in frontend/src/components/charts/ChartRenderer.tsx
- [ ] T046 [P] [US1] Create TableRenderer component in frontend/src/components/charts/TableRenderer.tsx
- [ ] T047 [US1] Create ChatPage route component in frontend/src/pages/ChatPage.tsx
- [ ] T048 [US1] Wire up chat route in frontend/src/App.tsx

**Checkpoint**: User Story 1 complete - users can query data in EN/AR with streaming responses

---

## Phase 4: User Story 10 - Language and Locale Preferences (Priority: P1)

**Goal**: Enable instant language switching with all UI elements, charts, and AI responses adapting accordingly

**Independent Test**: Change language from English to Arabic and verify all UI text, charts, navigation switch instantly with RTL layout

### Tests for User Story 10

- [ ] T049 [P] [US10] Contract test for GET/PATCH /preferences in tests/contract/preferences_api_test.go
- [ ] T050 [P] [US10] Integration test for locale persistence across sessions in tests/integration/locale_test.go

### Backend Implementation for User Story 10

- [ ] T051 [US10] Create PreferencesHandler in internal/api/handlers/preferences.go
- [ ] T052 [US10] Implement GET /preferences endpoint in internal/api/handlers/preferences.go
- [ ] T053 [US10] Implement PATCH /preferences endpoint with validation in internal/api/handlers/preferences.go
- [ ] T054 [US10] Register preferences routes in internal/api/routes.go

### Frontend Implementation for User Story 10

- [ ] T055 [P] [US10] Create usePreferences hook in frontend/src/hooks/usePreferences.ts
- [ ] T056 [US10] Integrate LanguageSwitcher with preferences API in frontend/src/components/common/LanguageSwitcher.tsx
- [ ] T057 [US10] Implement instant locale switching without page reload in frontend/src/hooks/useLocale.ts
- [ ] T058 [US10] Apply RTL layout direction changes in frontend/src/App.tsx

**Checkpoint**: User Story 10 complete - instant language switching works across all components

---

## Phase 5: User Story 2 - Pin Insights to Dashboard (Priority: P1)

**Goal**: Allow users to save charts to their personal dashboard for ongoing monitoring

**Independent Test**: Click pin icon on any chart, navigate to dashboard, verify pinned chart appears with correct data and auto-refresh

### Tests for User Story 2

- [ ] T059 [P] [US2] Contract test for pinned charts CRUD in tests/contract/dashboard_api_test.go
- [ ] T060 [P] [US2] Integration test for pin-refresh flow in tests/integration/dashboard_pin_test.go

### Backend Implementation for User Story 2

- [ ] T061 [US2] Create DashboardHandler in internal/api/handlers/dashboard.go
- [ ] T062 [US2] Implement GET /dashboard/charts endpoint in internal/api/handlers/dashboard.go
- [ ] T063 [US2] Implement POST /dashboard/charts (pin chart) in internal/api/handlers/dashboard.go
- [ ] T064 [US2] Implement PATCH /dashboard/charts/{id} endpoint in internal/api/handlers/dashboard.go
- [ ] T065 [US2] Implement DELETE /dashboard/charts/{id} endpoint in internal/api/handlers/dashboard.go
- [ ] T066 [US2] Implement POST /dashboard/charts/{id}/refresh endpoint in internal/api/handlers/dashboard.go
- [ ] T067 [US2] Implement POST /dashboard/charts/reorder endpoint in internal/api/handlers/dashboard.go
- [ ] T068 [US2] Register dashboard routes in internal/api/routes.go

### Frontend Implementation for User Story 2

- [ ] T069 [P] [US2] Create useDashboard hook in frontend/src/hooks/useDashboard.ts
- [ ] T070 [P] [US2] Create DashboardGrid component with responsive layout in frontend/src/components/dashboard/DashboardGrid.tsx
- [ ] T071 [P] [US2] Create DashboardWidget component in frontend/src/components/dashboard/DashboardWidget.tsx
- [ ] T072 [P] [US2] Create WidgetSettings component in frontend/src/components/dashboard/WidgetSettings.tsx
- [ ] T073 [P] [US2] Create ChartActions component (Pin, Export) in frontend/src/components/charts/ChartActions.tsx
- [ ] T074 [US2] Create DashboardPage route component in frontend/src/pages/DashboardPage.tsx
- [ ] T075 [US2] Wire up dashboard route in frontend/src/App.tsx
- [ ] T076 [US2] Implement chart refresh scheduler in frontend/src/hooks/useDashboard.ts

**Checkpoint**: User Stories 1, 10, 2 complete - MVP ready for demo

---

## Phase 6: User Story 3 - Drill-Down Analysis (Priority: P2)

**Goal**: Enable users to click on chart data points to see underlying detailed breakdowns

**Independent Test**: Display revenue chart, click on "March Clinic Revenue" bar, verify detailed transaction table appears

### Tests for User Story 3

- [ ] T077 [P] [US3] Contract test for POST /chat/drilldown in tests/contract/drilldown_api_test.go
- [ ] T078 [P] [US3] Integration test for multi-level drill-down in tests/integration/drilldown_test.go

### Backend Implementation for User Story 3

- [ ] T079 [US3] Implement POST /chat/drilldown endpoint in internal/api/handlers/chat.go
- [ ] T080 [US3] Add drill-down query generation logic in internal/api/handlers/chat.go
- [ ] T081 [US3] Implement permission check for detail-level access in internal/api/handlers/chat.go

### Frontend Implementation for User Story 3

- [ ] T082 [P] [US3] Create DrillDownModal component in frontend/src/components/charts/DrillDownModal.tsx
- [ ] T083 [US3] Add click handlers to ChartRenderer for drill-down in frontend/src/components/charts/ChartRenderer.tsx
- [ ] T084 [US3] Implement drill-down state management in frontend/src/hooks/useChat.ts

**Checkpoint**: User Story 3 complete - multi-level data exploration works

---

## Phase 7: User Story 4 - KPI Alerts (Priority: P2)

**Goal**: Notify users when key metrics cross defined thresholds via in-app and email

**Independent Test**: Create alert rule "Notify when pharmacy stock < 50", simulate condition, verify notification received

### Tests for User Story 4

- [ ] T085 [P] [US4] Contract test for alert rules CRUD in tests/contract/alerts_api_test.go
- [ ] T086 [P] [US4] Integration test for alert evaluation in tests/integration/alert_evaluation_test.go
- [ ] T087 [P] [US4] Integration test for notification delivery in tests/integration/notification_test.go

### Backend Implementation for User Story 4

- [ ] T088 [US4] Create AlertHandler in internal/api/handlers/alerts.go
- [ ] T089 [US4] Implement GET /alerts/rules endpoint in internal/api/handlers/alerts.go
- [ ] T090 [US4] Implement POST /alerts/rules endpoint in internal/api/handlers/alerts.go
- [ ] T091 [US4] Implement PATCH /alerts/rules/{id} endpoint in internal/api/handlers/alerts.go
- [ ] T092 [US4] Implement DELETE /alerts/rules/{id} endpoint in internal/api/handlers/alerts.go
- [ ] T093 [US4] Implement POST /alerts/rules/{id}/toggle endpoint in internal/api/handlers/alerts.go
- [ ] T094 [US4] Implement POST /alerts/rules/{id}/test endpoint in internal/api/handlers/alerts.go
- [ ] T095 [US4] Implement GET /alerts/metrics endpoint in internal/api/handlers/alerts.go
- [ ] T096 [US4] Create AlertScheduler service in internal/services/alert_scheduler.go
- [ ] T097 [US4] Implement alert evaluation logic in internal/services/alert_scheduler.go
- [ ] T098 [US4] Create NotificationService in internal/services/notification.go
- [ ] T099 [US4] Implement in-app notification delivery in internal/services/notification.go
- [ ] T100 [US4] Implement email notification delivery in internal/services/notification.go
- [ ] T101 [US4] Implement GET /notifications endpoint in internal/api/handlers/alerts.go
- [ ] T102 [US4] Implement POST /notifications/{id}/read endpoint in internal/api/handlers/alerts.go
- [ ] T103 [US4] Register alert routes in internal/api/routes.go
- [ ] T104 [US4] Schedule alert evaluation via NATS JetStream in internal/services/alert_scheduler.go

### Frontend Implementation for User Story 4

- [ ] T105 [P] [US4] Create useAlerts hook in frontend/src/hooks/useAlerts.ts
- [ ] T106 [P] [US4] Create AlertList component in frontend/src/components/alerts/AlertList.tsx
- [ ] T107 [P] [US4] Create AlertForm component in frontend/src/components/alerts/AlertForm.tsx
- [ ] T108 [P] [US4] Create NotificationToast component in frontend/src/components/alerts/NotificationToast.tsx
- [ ] T109 [US4] Create AlertsPage route component in frontend/src/pages/AlertsPage.tsx
- [ ] T110 [US4] Wire up alerts route in frontend/src/App.tsx
- [ ] T111 [US4] Add notification polling/WebSocket listener in frontend/src/hooks/useAlerts.ts

**Checkpoint**: User Story 4 complete - proactive monitoring and notifications work

---

## Phase 8: User Story 5 - Scheduled Reports (Priority: P2)

**Goal**: Deliver recurring reports automatically via email in PDF, spreadsheet, or CSV format

**Independent Test**: Schedule weekly report, advance time to scheduled moment, verify report generated and emailed

### Tests for User Story 5

- [ ] T112 [P] [US5] Contract test for scheduled reports CRUD in tests/contract/reports_api_test.go
- [ ] T113 [P] [US5] Integration test for report generation in tests/integration/report_generation_test.go
- [ ] T114 [P] [US5] Integration test for PDF generation with Arabic in tests/integration/report_pdf_test.go

### Backend Implementation for User Story 5

- [ ] T115 [US5] Create ReportHandler in internal/api/handlers/reports.go
- [ ] T116 [US5] Implement GET /reports/scheduled endpoint in internal/api/handlers/reports.go
- [ ] T117 [US5] Implement POST /reports/scheduled endpoint in internal/api/handlers/reports.go
- [ ] T118 [US5] Implement PATCH /reports/scheduled/{id} endpoint in internal/api/handlers/reports.go
- [ ] T119 [US5] Implement DELETE /reports/scheduled/{id} endpoint in internal/api/handlers/reports.go
- [ ] T120 [US5] Implement POST /reports/scheduled/{id}/toggle endpoint in internal/api/handlers/reports.go
- [ ] T121 [US5] Implement POST /reports/scheduled/{id}/run endpoint in internal/api/handlers/reports.go
- [ ] T122 [US5] Implement GET /reports/scheduled/{id}/runs endpoint in internal/api/handlers/reports.go
- [ ] T123 [US5] Implement GET /reports/runs/{id}/download endpoint in internal/api/handlers/reports.go
- [ ] T124 [US5] Implement GET /reports/templates endpoint in internal/api/handlers/reports.go
- [ ] T125 [US5] Create ReportGenerator service in internal/services/report_generator.go
- [ ] T126 [US5] Implement PDF generation with Puppeteer in internal/services/pdf_generator.go
- [ ] T127 [US5] Implement spreadsheet generation with excelize in internal/services/spreadsheet_generator.go
- [ ] T128 [US5] Implement CSV generation in internal/services/csv_generator.go
- [ ] T129 [US5] Implement RTL layout for Arabic PDFs in internal/services/pdf_generator.go
- [ ] T130 [US5] Schedule report generation via NATS JetStream in internal/services/report_scheduler.go
- [ ] T131 [US5] Create EmailSender service in internal/services/email_sender.go
- [ ] T132 [US5] Register report routes in internal/api/routes.go

### Frontend Implementation for User Story 5

- [ ] T133 [P] [US5] Create useReports hook in frontend/src/hooks/useReports.ts
- [ ] T134 [P] [US5] Create ReportScheduler component in frontend/src/components/reports/ReportScheduler.tsx
- [ ] T135 [P] [US5] Create ReportHistory component in frontend/src/components/reports/ReportHistory.tsx
- [ ] T136 [US5] Create ReportsPage route component in frontend/src/pages/ReportsPage.tsx
- [ ] T137 [US5] Wire up reports route in frontend/src/App.tsx

**Checkpoint**: User Story 5 complete - automated report delivery works

---

## Phase 9: User Story 6 - Multi-Period Comparison (Priority: P3)

**Goal**: Compare metrics across different time periods side-by-side with delta annotations

**Independent Test**: Ask "Compare this month's revenue vs last month", verify side-by-side comparison chart with percentage change

### Tests for User Story 6

- [ ] T138 [P] [US6] Integration test for comparison queries in tests/integration/comparison_test.go

### Backend Implementation for User Story 6

- [ ] T139 [US6] Add comparison_period field handling in internal/api/handlers/chat.go
- [ ] T140 [US6] Implement period comparison SQL generation in internal/agents/module_a/a01_text_to_sql/
- [ ] T141 [US6] Add delta calculation logic in internal/api/handlers/chat.go

### Frontend Implementation for User Story 6

- [ ] T142 [US6] Create ComparisonChart component in frontend/src/components/charts/ComparisonChart.tsx
- [ ] T143 [US6] Add comparison period selector to QueryInput in frontend/src/components/chat/QueryInput.tsx

**Checkpoint**: User Story 6 complete - temporal comparisons work

---

## Phase 10: User Story 7 - Data Export (Priority: P3)

**Goal**: Export any chart or table to PDF, spreadsheet, or CSV format

**Independent Test**: Display any data table, click export, select spreadsheet, verify downloaded file contains correct data

### Tests for User Story 7

- [ ] T144 [P] [US7] Integration test for export operations in tests/integration/export_test.go
- [ ] T145 [P] [US7] Integration test for Arabic PDF export in tests/integration/export_arabic_test.go

### Backend Implementation for User Story 7

- [ ] T146 [US7] Implement POST /chat/export endpoint in internal/api/handlers/chat.go
- [ ] T147 [US7] Add async export for large datasets in internal/api/handlers/chat.go
- [ ] T148 [US7] Integrate PDF generator in internal/api/handlers/chat.go

### Frontend Implementation for User Story 7

- [ ] T149 [P] [US7] Create useExport hook in frontend/src/hooks/useExport.ts
- [ ] T150 [US7] Implement ExportButton with format selection in frontend/src/components/common/ExportButton.tsx
- [ ] T151 [US7] Add export progress indicator in frontend/src/components/common/ExportButton.tsx

**Checkpoint**: User Story 7 complete - data export works in all formats

---

## Phase 11: User Story 9 - Quick-Action Prompts (Priority: P3)

**Goal**: Display pre-built question suggestions for new users to get started quickly

**Independent Test**: Load chat interface, click any quick-action prompt, verify query executes and displays results

### Tests for User Story 9

- [ ] T152 [P] [US9] Contract test for quick-actions API in tests/contract/quick_actions_test.go

### Backend Implementation for User Story 9

- [ ] T153 [US9] Implement GET /dashboard/quick-actions endpoint in internal/api/handlers/dashboard.go
- [ ] T154 [US9] Implement PUT /dashboard/quick-actions endpoint (admin) in internal/api/handlers/dashboard.go

### Frontend Implementation for User Story 9

- [ ] T155 [P] [US9] Create QuickActions component in frontend/src/components/chat/QuickActions.tsx
- [ ] T156 [US9] Add quick-action carousel to ChatInterface in frontend/src/components/chat/ChatInterface.tsx

**Checkpoint**: User Story 9 complete - onboarding prompts work

---

## Phase 12: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T157 [P] Add OPA policies for new endpoints in policies/dashboard_access.rego
- [ ] T158 [P] Add audit logging for dashboard operations in internal/warehouse/audit.go
- [ ] T159 [P] Implement rate limiting for chat endpoint in internal/api/middleware/rate_limit.go
- [ ] T160 [P] Add CI check for i18n key coverage in .github/workflows/i18n-check.yml
- [ ] T161 Performance optimization for dashboard grid rendering
- [ ] T162 Add error boundary components in frontend/src/components/common/ErrorBoundary.tsx
- [ ] T163 Run quickstart.md validation and fix any issues
- [ ] T164 Update CLAUDE.md with new component patterns
- [ ] T165 Security review for WebSocket connections

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-11)**: All depend on Foundational phase completion
  - US1 (P1): Can start immediately after Phase 2
  - US10 (P1): Can run parallel with US1
  - US2 (P1): Can run parallel with US1/US10, but integrates with US1 charts
  - US3-US9 (P2-P3): Can proceed in priority order after P1 stories complete
- **Polish (Phase 12)**: Depends on all desired user stories being complete

### User Story Dependencies

| Story | Priority | Depends On | Can Run Parallel With |
|-------|----------|------------|----------------------|
| US1 - Data Query | P1 | Phase 2 only | US10, US2 |
| US10 - Locale Preferences | P1 | Phase 2 only | US1, US2 |
| US2 - Dashboard Pinning | P1 | Phase 2, US1 charts | US1, US10 |
| US3 - Drill-Down | P2 | US1 | US4, US5 |
| US4 - KPI Alerts | P2 | Phase 2 only | US3, US5 |
| US5 - Scheduled Reports | P2 | Phase 2 only | US3, US4 |
| US6 - Comparison | P3 | US1 | US7, US9 |
| US7 - Data Export | P3 | US1 | US6, US9 |
| US9 - Quick Actions | P3 | Phase 2 only | US6, US7 |

### Within Each User Story

1. Tests MUST be written and FAIL before implementation
2. Backend handlers before frontend hooks
3. Frontend hooks before components
4. Components before page routes
5. Story complete before moving to next priority

### Parallel Opportunities

- All Setup migrations (T001-T005) can run in parallel
- All i18n translation files (T006-T013) can run in parallel
- All warehouse repositories (T014-T019) can run in parallel
- All common UI components (T025-T027) can run in parallel
- Tests within a story marked [P] can run in parallel
- Different user stories (P1 group) can be worked on in parallel by different developers

---

## Parallel Example: User Story 1 (Core Chat)

```bash
# Launch all tests for User Story 1 together:
Task: "Contract test for POST /chat/query in tests/contract/chat_api_test.go"
Task: "Contract test for WebSocket /chat/stream in tests/contract/chat_stream_test.go"
Task: "Integration test for English query flow in tests/integration/chat_flow_test.go"
Task: "Integration test for Arabic query flow in tests/integration/chat_arabic_test.go"

# Launch all frontend components for User Story 1 together (after backend):
Task: "Create useChat hook with streaming support in frontend/src/hooks/useChat.ts"
Task: "Create ChatInterface component in frontend/src/components/chat/ChatInterface.tsx"
Task: "Create MessageList component in frontend/src/components/chat/MessageList.tsx"
Task: "Create QueryInput component with RTL support in frontend/src/components/chat/QueryInput.tsx"
Task: "Create StreamingMessage component for progressive rendering in frontend/src/components/chat/StreamingMessage.tsx"
Task: "Create ChartRenderer component with ECharts in frontend/src/components/charts/ChartRenderer.tsx"
Task: "Create TableRenderer component in frontend/src/components/charts/TableRenderer.tsx"
```

---

## Implementation Strategy

### MVP First (User Stories 1, 10, 2 Only)

1. Complete Phase 1: Setup (migrations + i18n)
2. Complete Phase 2: Foundational (repositories + infrastructure)
3. Complete Phase 3: User Story 1 (Natural Language Query)
4. Complete Phase 4: User Story 10 (Locale Preferences)
5. Complete Phase 5: User Story 2 (Dashboard Pinning)
6. **STOP and VALIDATE**: Test all three P1 stories independently
7. Deploy/demo MVP

**MVP delivers**: Chat queries in EN/AR + language switching + pinned dashboard

### Incremental Delivery

| Milestone | Stories | Value Delivered |
|-----------|---------|-----------------|
| MVP | US1, US10, US2 | Core conversational BI with i18n |
| v1.1 | US3, US4 | Drill-down analysis + proactive alerts |
| v1.2 | US5 | Automated report delivery |
| v1.3 | US6, US7, US9 | Advanced analytics + export + onboarding |

### Parallel Team Strategy

With 3 developers:

1. **Week 1**: All complete Setup + Foundational together
2. **Week 2-3**:
   - Developer A: User Story 1 (Chat)
   - Developer B: User Story 10 (Preferences)
   - Developer C: User Story 2 (Dashboard)
3. **Week 4**: Integration testing + MVP demo
4. **Week 5+**: P2 stories in parallel (US3, US4, US5)

---

## Summary

| Metric | Count |
|--------|-------|
| **Total Tasks** | 165 |
| **Setup Phase** | 13 |
| **Foundational Phase** | 15 |
| **User Story 1 (P1)** | 20 |
| **User Story 10 (P1)** | 10 |
| **User Story 2 (P1)** | 18 |
| **User Story 3 (P2)** | 8 |
| **User Story 4 (P2)** | 27 |
| **User Story 5 (P2)** | 26 |
| **User Story 6 (P3)** | 6 |
| **User Story 7 (P3)** | 8 |
| **User Story 9 (P3)** | 5 |
| **Polish Phase** | 9 |
| **Parallel Opportunities** | 89 tasks marked [P] |

**MVP Scope**: Tasks T001-T076 (76 tasks) = Phases 1-5 = US1 + US10 + US2

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- US8 (Mobile) is out of scope per spec.md
