---
name: medisync-dev
description: Guides developers working on the MediSync AI-powered healthcare BI platform. Provides architecture patterns, AI agent development workflows, security requirements, and i18n standards. Use when adding features, creating agents, running tests, or implementing integrations.
---

# MediSync Development Guide

MediSync is an AI-powered conversational BI platform that unifies HIMS (clinic operations) and Tally ERP (accounting) data for healthcare and pharmacy businesses.

★ Insight ─────────────────────────────────────
MediSync's three-plane architecture separates concerns:
1. **Data Plane** - ETL extracts from Tally/HIMS into a warehouse
2. **Intelligence Plane** - AI agents query with read-only access
3. **Action Plane** - Human-gated write-backs to external systems

This ensures analytics never impact operational systems.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Purpose** | Conversational BI + AI Accountant for healthcare |
| **Tech Stack** | Go 1.26, React 19, Flutter 3.42, PostgreSQL 18.2, Genkit |
| **AI Agents** | 58 specialized agents across 5 modules |
| **Security** | Keycloak auth, OPA policies, HITL gates for writes |
| **i18n** | English (LTR) and Arabic (RTL) from Phase 1 |

## Directory Structure

```
medisync/
├── cmd/                    # Entry points: api, etl, migrate
├── internal/
│   ├── agents/            # AI agent implementations (modules A-E)
│   ├── api/               # HTTP handlers, middleware
│   ├── auth/              # Keycloak JWT, OPA client
│   ├── etl/               # Tally & HIMS connectors
│   ├── warehouse/         # PostgreSQL repository
│   └── cache/             # Redis client
├── policies/              # OPA Rego policies
├── migrations/            # SQL migrations
├── frontend/              # React web app
└── mobile/                # Flutter mobile app
```

## Five Agent Modules

| Module | Focus | Location |
|--------|-------|----------|
| **A** | Conversational BI | `internal/agents/module_a/` |
| **B** | AI Accountant | `internal/agents/module_b/` |
| **C** | Easy Reports | `internal/agents/module_c/` |
| **D** | Search Analytics | `internal/agents/module_d/` |
| **E** | i18n | `internal/agents/module_e/` |

## Development Workflow

### Start Infrastructure
```bash
docker-compose up -d
```

### Run Migrations
```bash
go run ./cmd/migrate
```

### Start Services
```bash
# API server
go run ./cmd/api

# ETL service
go run ./cmd/etl

# Web frontend
cd frontend && npm run dev

# Mobile app
cd mobile && flutter run
```

### Run Tests
```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend && npm test

# Mobile tests
cd mobile && flutter test
```

## AI Agent Development

All agents follow the **Genkit Flow** pattern:

```go
func (s *AgentService) MyAgentFlow(ctx context.Context, req MyRequest) (*MyResponse, error) {
    // 1. Validate input
    // 2. Detect language (E-01)
    // 3. Process with LLM
    // 4. Score confidence (A-06)
    // 5. Format localized response (E-03)
}
```

Key principles:
- Define input/output structs explicitly
- Include confidence scoring for all AI outputs
- Route to HITL when confidence < threshold
- Use the `medisync_readonly` database role for queries

### Adding a New Agent

1. Create agent file in appropriate `internal/agents/module_X/`
2. Define Genkit flow with typed structs
3. Add confidence scoring and HITL gates
4. Register in agent supervisor
5. Add OPA policy for authorization
6. Write tests with mock data

See `references/agents.md` for detailed patterns.

## Security Requirements

### Read-Only Enforcement
```go
const dbRole = "medisync_readonly"

func isSelectOnlyQuery(sql string) bool {
    return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sql)), "SELECT")
}
```

### HITL Gates for Write-Backs
All Tally sync operations require:
1. Finance head approval
2. OPA policy verification
3. Audit log entry

### Security Checklist
- [ ] AI agents use `medisync_readonly` DB role
- [ ] SQL queries validated as SELECT-only
- [ ] Write-backs require HITL approval
- [ ] OPA policies cover new endpoints
- [ ] Audit log entries created
- [ ] PII masked based on user role

## i18n by Default

### Locale Detection Priority
1. `user_preferences.locale` (from JWT)
2. `Accept-Language` HTTP header
3. `?lang=ar` URL parameter
4. Default: `en`

### Adding Translations
1. Add key to `frontend/public/locales/en/*.json`
2. Add Arabic to `frontend/public/locales/ar/*.json`
3. Use: `const { t } = useTranslation()`
4. For mobile: add to `mobile/lib/l10n/app_*.arb`

All AI prompts include locale instructions:
```go
prompt := fmt.Sprintf(
    "ResponseLanguageInstruction: Respond in %s. Format numbers according to %s locale.",
    userLocale, userLocale,
)
```

## Common Tasks

### Database Migration
1. Create migration file in `migrations/`
2. Use `medisync_readonly` role reference for AI queries
3. Test on copy of production data
4. Run: `go run ./cmd/migrate`
5. Verify with rollback test

### Adding i18n Support
See `references/i18n.md` for comprehensive patterns.

### Writing Tests
See `references/testing.md` for unit, integration, and E2E patterns.

## Troubleshooting

| Issue | Solution |
|-------|----------|
| LLM hallucination | Check confidence scores (A-06), review domain mappings (A-04) |
| Low OCR accuracy | Check document preprocessing, route to handwriting agent (B-03) |
| Tally sync failures | Verify OPA policy, check approval workflow, review gateway logs |
| i18n issues | Verify Accept-Language/JWT locale, check translation keys, test RTL layout |

## Key Agent IDs

| ID | Name | Purpose |
|----|------|---------|
| A-01 | Text-to-SQL | NL to safe SQL |
| A-03 | Visualization Routing | Chart selection |
| A-06 | Confidence Scorer | 0-100% confidence |
| B-02 | OCR Extraction | Document field extraction |
| B-05 | Ledger Mapping | Tally GL suggestions |
| B-09 | Tally Sync | Push to Tally ERP |
| E-01 | Language Detection | en/ar classification |
| E-03 | Localized Formatter | Locale-aware formatting |

## Performance Targets

| Metric | Target |
|--------|--------|
| Query accuracy | ≥ 95% |
| Query latency (P95) | < 5 seconds |
| Dashboard load | < 3 seconds |
| OCR accuracy (standard) | ≥ 95% |
| Tally sync success | ≥ 99.5% |

## Philosophy

> "MediSync exists to liberate healthcare and accounting teams from the tyranny of spreadsheets and manual reconciliation."

When working on this codebase:
- **Security first**: Never compromise on HITL gates
- **i18n by default**: Every feature works in English and Arabic
- **Open source**: Always choose OSI-approved licenses
- **User trust**: Confidence scores, audit trails, transparency

## Detailed References

| Reference | Content |
|-----------|---------|
| `references/architecture.md` | Layer stack, module breakdown, data flows |
| `references/agents.md` | Genkit patterns, A2A protocol, 58 agent specs |
| `references/i18n.md` | Locale handling, RTL support, translation patterns |
| `references/testing.md` | Unit/integration/E2E testing strategies |
