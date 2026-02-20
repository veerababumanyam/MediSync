# Quickstart: Dashboard, Chat UI & Advanced Features

**Feature**: 002-dashboard-advanced-features
**Last Updated**: 2026-02-20

This guide helps developers quickly set up and start working on the Dashboard, Chat UI & Advanced Features implementation.

---

## Prerequisites

- **Go 1.26+** installed
- **Node.js 20+** and **pnpm** (or npm) installed
- **Docker** and **Docker Compose** for local infrastructure
- **PostgreSQL 18.2** client tools
- Access to the MediSync repository

---

## Quick Setup

### 1. Start Infrastructure

```bash
# From repository root
docker-compose up -d postgres redis nats

# Verify services are running
docker-compose ps
```

### 2. Run Database Migrations

```bash
# Run all migrations including new dashboard tables
go run ./cmd/migrate

# Verify new tables exist
psql -d medisync -c "\dt user_preferences pinned_charts alert_rules notifications scheduled_reports scheduled_report_runs chat_messages"
```

### 3. Start Backend API

```bash
# Set environment variables (or use .env file)
export DATABASE_URL="postgres://medisync:password@localhost:5432/medisync?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export NATS_URL="nats://localhost:4222"
export KEYCLOAK_URL="http://localhost:8080"

# Start API server
go run ./cmd/api

# API will be available at http://localhost:8080/api/v1
```

### 4. Start Frontend Development Server

```bash
cd frontend

# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Frontend will be available at http://localhost:5173
```

---

## Project Structure Overview

```
medisync/
├── cmd/
│   └── api/main.go              # API entry point
├── internal/
│   ├── api/handlers/            # HTTP handlers (ADD: chat, dashboard, alerts, reports)
│   ├── warehouse/               # Database repositories (ADD: pinned_chart, alert_rule)
│   └── agents/module_a/         # Existing AI agents (Text-to-SQL, Visualization)
├── frontend/
│   └── src/
│       ├── components/          # React components (ADD: chat/, dashboard/, alerts/)
│       ├── pages/               # Page components (ADD: ChatPage, DashboardPage)
│       ├── hooks/               # Custom hooks (ADD: useChat, useDashboard)
│       └── i18n/locales/        # Translation files (ADD: chat.json, dashboard.json)
├── migrations/
│   ├── 010_user_preferences*.sql   # User preferences table
│   ├── 011_pinned_charts*.sql      # Pinned charts table
│   └── ...                         # Additional migrations
└── specs/002-dashboard-advanced-features/
    ├── spec.md                  # Feature specification
    ├── plan.md                  # This implementation plan
    ├── research.md              # Technology decisions
    ├── data-model.md            # Entity definitions
    ├── quickstart.md            # This file
    └── contracts/               # API contracts
```

---

## Key Files to Create

### Backend (Go)

| File | Purpose |
|------|---------|
| `internal/api/handlers/chat.go` | Chat query endpoints |
| `internal/api/handlers/dashboard.go` | Pinned chart CRUD |
| `internal/api/handlers/alerts.go` | Alert rule management |
| `internal/api/handlers/reports.go` | Scheduled reports |
| `internal/api/websocket/stream.go` | WebSocket streaming |
| `internal/warehouse/pinned_chart.go` | Pinned chart repository |
| `internal/warehouse/alert_rule.go` | Alert rule repository |
| `internal/services/alert_scheduler.go` | Alert evaluation |
| `internal/services/report_generator.go` | Report generation |

### Frontend (React/TypeScript)

| File | Purpose |
|------|---------|
| `frontend/src/components/chat/ChatInterface.tsx` | Main chat container |
| `frontend/src/components/chat/QueryInput.tsx` | Natural language input |
| `frontend/src/components/charts/ChartRenderer.tsx` | ECharts wrapper |
| `frontend/src/components/dashboard/DashboardGrid.tsx` | Pinned charts grid |
| `frontend/src/hooks/useChat.ts` | Chat state + streaming |
| `frontend/src/hooks/useLocale.ts` | Locale detection/switching |
| `frontend/src/i18n/locales/en/chat.json` | English translations |
| `frontend/src/i18n/locales/ar/chat.json` | Arabic translations |

---

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/dashboard-chat 001-dashboard-advanced-features
```

### 2. Run Tests

```bash
# Backend tests
go test ./internal/... -v

# Frontend tests
cd frontend && pnpm test

# E2E tests (requires running services)
pnpm test:e2e
```

### 3. Check Code Quality

```bash
# Go linting
golangci-lint run

# Frontend linting
cd frontend && pnpm lint

# Type checking
pnpm type-check
```

### 4. Verify i18n Coverage

```bash
# Check for missing translations
cd frontend
node scripts/check-i18n.js

# Should output: ✅ All keys present in en and ar
```

---

## API Quick Reference

### Chat Endpoints

```bash
# Submit a query
curl -X POST http://localhost:8080/api/v1/chat/query \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "What is today'\''s revenue?", "session_id": "uuid"}'

# WebSocket streaming
wscat -c "ws://localhost:8080/api/v1/chat/stream?session_id=uuid" \
  -H "Authorization: Bearer $TOKEN"
```

### Dashboard Endpoints

```bash
# List pinned charts
curl http://localhost:8080/api/v1/dashboard/charts \
  -H "Authorization: Bearer $TOKEN"

# Pin a chart
curl -X POST http://localhost:8080/api/v1/dashboard/charts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"message_id": "uuid", "title": "Revenue Chart"}'
```

### Alert Endpoints

```bash
# Create alert rule
curl -X POST http://localhost:8080/api/v1/alerts/rules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Low Stock Alert",
    "metric_id": "inventory.stock.level",
    "operator": "lt",
    "threshold": 50,
    "check_interval": 300,
    "channels": ["in_app", "email"]
  }'
```

---

## Common Development Tasks

### Adding a New Component

1. Create component file in appropriate directory:
   ```bash
   touch frontend/src/components/chat/NewComponent.tsx
   ```

2. Use i18n for all user-facing text:
   ```tsx
   import { useTranslation } from 'react-i18next';

   export function NewComponent() {
     const { t } = useTranslation('chat');
     return <div>{t('newComponent.title')}</div>;
   }
   ```

3. Add translations to both locales:
   ```bash
   # frontend/src/i18n/locales/en/chat.json
   echo '{"newComponent": {"title": "New Component"}}' >> ...

   # frontend/src/i18n/locales/ar/chat.json
   echo '{"newComponent": {"title": "مكون جديد"}}' >> ...
   ```

### Adding a New API Endpoint

1. Create handler in `internal/api/handlers/`:
   ```go
   func (h *Handler) GetNewEndpoint(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```

2. Register route in `internal/api/routes.go`:
   ```go
   r.Get("/new-endpoint", h.GetNewEndpoint)
   ```

3. Add OPA policy if needed in `policies/`:
   ```rego
   allow {
       input.path == ["api", "v1", "new-endpoint"]
       input.method == "GET"
   }
   ```

### Adding RTL Support

1. Use Tailwind logical properties:
   ```tsx
   // ❌ Bad - breaks RTL
   <div className="ml-4 pl-2">

   // ✅ Good - RTL-aware
   <div className="ms-4 ps-2">
   ```

2. Use `dir` attribute for chart containers:
   ```tsx
   <div dir={locale === 'ar' ? 'rtl' : 'ltr'}>
     <ChartRenderer spec={spec} />
   </div>
   ```

---

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check connection
psql -h localhost -U medisync -d medisync -c "SELECT 1"
```

### Redis Connection Issues

```bash
# Check Redis is running
docker-compose ps redis

# Test connection
redis-cli ping
# Should return: PONG
```

### i18n Keys Not Loading

```bash
# Verify locale files are valid JSON
cat frontend/src/i18n/locales/en/chat.json | jq .

# Check i18next configuration
cat frontend/src/i18n/index.ts
```

### WebSocket Connection Fails

```bash
# Check WebSocket upgrade headers
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: test" \
  -H "Sec-WebSocket-Version: 13" \
  http://localhost:8080/api/v1/chat/stream
```

---

## Related Documentation

- [Feature Specification](./spec.md)
- [Implementation Plan](./plan.md)
- [Data Model](./data-model.md)
- [API Contracts](./contracts/)
- [Research Notes](./research.md)

---

## Getting Help

1. Check existing documentation in `docs/` directory
2. Review `CLAUDE.md` for project patterns
3. Check constitution in `.specify/memory/constitution.md`
4. Ask in team chat or create an issue

---

*Happy coding! 🚀*
