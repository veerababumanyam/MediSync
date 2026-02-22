# Implementation Plan: Council of AIs Consensus System

**Branch**: `001-council-ai-consensus` | **Date**: 2026-02-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-council-ai-consensus/spec.md`

## Summary

Implement a multi-agent consensus system that eradicates AI hallucinations by requiring agreement among independent agent instances, each grounded by Graph-of-Thoughts retrieval from a Medical Knowledge Graph. The system uses weighted voting with semantic similarity clustering to calculate consensus, includes full audit trails for HIPAA compliance, and provides transparent evidence exploration for user trust.

## Technical Context

**Language/Version**: Go 1.26 (backend), TypeScript 5.9 (frontend)
**Primary Dependencies**: go-chi/chi, Google Genkit, Agent ADK, pgvector, Redis, NATS JetStream
**Storage**: PostgreSQL 18.2 + pgvector (Knowledge Graph, audit trails), Redis (evidence cache)
**Testing**: Go testing + testify (unit), Vitest (frontend), Playwright (e2e)
**Target Platform**: Linux server (containerized), Web browser
**Project Type**: Web application (backend + frontend)
**Performance Goals**: 95% of queries <10s, 99.5% availability, 1000 req/min throughput
**Constraints**: <3s agent timeout, 5-min evidence cache, 7-year audit retention
**Scale/Scope**: 1000 queries/day, 3 agent instances, 10K+ Knowledge Graph nodes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Compliance |
|-----------|--------|------------|
| I. Security First & HITL Gates | ✅ PASS | Read-only Knowledge Graph access via medisync_readonly role; no autonomous writes |
| II. Read-Only Intelligence Plane | ✅ PASS | All AI queries use SELECT-only against data warehouse; consensus is computation-only |
| III. i18n by Default | ✅ PASS | Response text uses i18next keys; confidence indicators localized |
| IV. Open Source Only | ✅ PASS | All dependencies (Genkit, pgvector, NATS) are Apache-2.0/MIT/BSD |
| V. Test-Driven Development | ✅ PASS | Unit tests required for consensus algorithm; integration tests for multi-agent flows |

**Post-Design Re-check**: All principles maintained. No new dependencies introduced outside approved stack.

## Project Structure

### Documentation (this feature)

```text
specs/001-council-ai-consensus/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 research findings
├── data-model.md        # Phase 1 entity definitions
├── quickstart.md        # Phase 1 developer guide
├── contracts/           # Phase 1 API contracts
│   └── openapi.yaml     # OpenAPI 3.1 specification
├── checklists/          # Quality checklists
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 tasks (via /speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── agents/
│   └── council/              # Council of AIs package
│       ├── coordinator.go    # Deliberation orchestration
│       ├── consensus.go      # Weighted voting + semantic clustering
│       ├── agent.go          # Agent instance wrapper
│       ├── evidence.go       # Graph-of-Thoughts retrieval
│       ├── semantic.go       # Equivalence detection (95% threshold)
│       ├── health.go         # Agent health monitoring
│       └── coordinator_test.go
├── api/
│   └── handlers/
│       └── council.go        # HTTP handlers for Council API
├── warehouse/
│   └── knowledge_graph.go    # Knowledge Graph repository
└── cache/
    └── evidence_cache.go     # 5-minute evidence caching

migrations/
└── council_ai_consensus/     # Database migrations
    ├── 001_knowledge_graph.up.sql
    ├── 002_council_tables.up.sql
    └── 003_audit_partition.up.sql

policies/
└── council.rego              # OPA policies for Council access

frontend/
└── src/
    ├── services/
    │   └── councilService.ts # API client
    ├── hooks/
    │   └── useCouncil.ts     # React hook
    └── components/
        └── council/
            ├── QueryInput.tsx
            ├── ResponseDisplay.tsx
            ├── EvidenceExplorer.tsx
            ├── ConfidenceIndicator.tsx
            └── UncertaintyDisplay.tsx
```

**Structure Decision**: Uses existing MediSync monorepo structure. Council package added to `internal/agents/` following the established module pattern (module_a through module_e). Frontend components added to dedicated `council/` directory under components.

## Complexity Tracking

> No violations requiring justification. All implementation follows established patterns.

| Concern | Approach | Rationale |
|---------|----------|-----------|
| Multi-agent orchestration | Agent ADK primitives | Aligns with existing agent architecture |
| Consensus algorithm | Custom weighted voting | Simpler than Byzantine FT, fits requirements |
| Knowledge Graph | PostgreSQL + pgvector | In approved stack, supports vector search |
| Audit retention | Partitioned tables | Native PostgreSQL, supports 7-year HIPAA |

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Research | [research.md](./research.md) | Technology decisions and rationale |
| Data Model | [data-model.md](./data-model.md) | Entity definitions and relationships |
| API Contract | [contracts/openapi.yaml](./contracts/openapi.yaml) | OpenAPI 3.1 specification |
| Quickstart | [quickstart.md](./quickstart.md) | Developer onboarding guide |

## Next Steps

Run `/speckit.tasks` to generate implementation tasks based on this plan.
