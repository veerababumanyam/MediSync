# Tasks: Council of AIs Consensus System

**Input**: Design documents from `/specs/001-council-ai-consensus/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml

**Tests**: Included per MediSync Constitution Principle V (TDD requirement)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `internal/` at repository root
- **Frontend**: `frontend/src/` at repository root
- **Migrations**: `migrations/council_ai_consensus/`
- **Policies**: `policies/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory structure

- [x] T001 Create directory structure for Council package in internal/agents/council/
- [x] T002 [P] Create migrations directory at migrations/council_ai_consensus/
- [x] T003 [P] Create frontend component directory at frontend/src/components/council/
- [x] T004 [P] Add Council environment variables to config (COUNCIL_MIN_AGENTS, COUNCIL_DEFAULT_THRESHOLD, COUNCIL_AGENT_TIMEOUT, COUNCIL_CACHE_TTL)
- [x] T005 [P] Add i18n keys for Council UI in frontend/src/i18n/locales/en/council.json and frontend/src/i18n/locales/ar/council.json

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Database Migrations

- [ ] T006 Create Knowledge Graph nodes table migration in migrations/council_ai_consensus/001_knowledge_graph.up.sql
- [ ] T007 [P] Create Council core tables migration (council_deliberations, agent_instances, agent_responses) in migrations/council_ai_consensus/002_council_tables.up.sql
- [ ] T008 [P] Create Consensus and Evidence tables migration (consensus_records, evidence_trails) in migrations/council_ai_consensus/003_consensus_evidence.up.sql
- [ ] T009 [P] Create Audit partition table migration in migrations/council_ai_consensus/004_audit_partition.up.sql
- [ ] T010 Run migrations and verify schema in database

### Core Types and Interfaces

- [ ] T011 [P] Define CouncilDeliberation struct and status enum in internal/agents/council/types.go
- [ ] T012 [P] Define AgentInstance struct and health status enum in internal/agents/council/types.go
- [ ] T013 [P] Define AgentResponse struct with embedding field in internal/agents/council/types.go
- [ ] T014 [P] Define ConsensusRecord struct in internal/agents/council/types.go
- [ ] T015 [P] Define EvidenceTrail struct in internal/agents/council/types.go
- [ ] T016 [P] Define KnowledgeGraphNode struct with edge types in internal/agents/council/types.go
- [ ] T017 [P] Define AuditEntry struct in internal/agents/council/types.go

### Repository Layer

- [ ] T018 [P] Implement Knowledge Graph repository interface in internal/warehouse/knowledge_graph.go
- [ ] T019 [P] Implement Council Deliberation repository in internal/agents/council/repository.go
- [ ] T020 [P] Implement Evidence cache with Redis in internal/cache/evidence_cache.go

### OPA Policy

- [ ] T021 Create Council access policy (role-based: admin=all, user=own) in policies/council.rego

### Health Infrastructure

- [ ] T022 Implement agent health monitoring types and constants in internal/agents/council/health.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Verified Response to Medical Query (Priority: P1) 🎯 MVP

**Goal**: Users can submit queries and receive hallucination-free responses backed by multi-agent consensus

**Independent Test**: Submit a healthcare query via API, verify response includes confidence score ≥80% and evidence attribution

### Tests for User Story 1 (TDD)

- [ ] T023 [P] [US1] Write unit test for consensus algorithm with 3 agents reaching agreement in internal/agents/council/consensus_test.go
- [ ] T024 [P] [US1] Write unit test for semantic equivalence detection (95% threshold) in internal/agents/council/semantic_test.go
- [ ] T025 [P] [US1] Write unit test for Graph-of-Thoughts retrieval in internal/agents/council/evidence_test.go
- [ ] T026 [P] [US1] Write unit test for agent instance timeout handling (3s) in internal/agents/council/agent_test.go
- [ ] T027 [US1] Write integration test for full deliberation flow in internal/agents/council/coordinator_test.go

### Backend Implementation for User Story 1

- [ ] T028 [P] [US1] Implement AgentInstance wrapper with Genkit integration in internal/agents/council/agent.go
- [ ] T029 [P] [US1] Implement semantic equivalence detection using pgvector in internal/agents/council/semantic.go
- [ ] T030 [US1] Implement weighted voting consensus algorithm in internal/agents/council/consensus.go (depends on T029)
- [ ] T031 [US1] Implement Graph-of-Thoughts retrieval with multi-hop traversal in internal/agents/council/evidence.go
- [ ] T032 [US1] Implement Council Coordinator orchestration in internal/agents/council/coordinator.go (depends on T028, T030, T031)
- [ ] T033 [US1] Implement POST /deliberations handler in internal/api/handlers/council.go
- [ ] T034 [US1] Implement GET /deliberations handler (list) with RBAC filtering in internal/api/handlers/council.go
- [ ] T035 [US1] Implement GET /deliberations/{id} handler with ownership check in internal/api/handlers/council.go
- [ ] T036 [US1] Register Council routes in internal/api/routes.go
- [ ] T037 [US1] Add structured logging for deliberation events in internal/agents/council/coordinator.go

### Frontend Implementation for User Story 1

- [ ] T038 [P] [US1] Create Council API client in frontend/src/services/councilService.ts
- [ ] T039 [P] [US1] Create useCouncil React hook in frontend/src/hooks/useCouncil.ts
- [ ] T040 [US1] Create QueryInput component with threshold configuration in frontend/src/components/council/QueryInput.tsx
- [ ] T041 [US1] Create ResponseDisplay component showing consensus response in frontend/src/components/council/ResponseDisplay.tsx
- [ ] T042 [US1] Create ConfidenceIndicator component (0-100% visual) in frontend/src/components/council/ConfidenceIndicator.tsx
- [ ] T043 [US1] Add Council components to i18n translation files

**Checkpoint**: At this point, User Story 1 should be fully functional - users can submit queries and receive consensus responses

---

## Phase 4: User Story 2 - Disagreement Detection and Uncertainty Signaling (Priority: P2)

**Goal**: When consensus cannot be reached, system transparently communicates uncertainty with range of positions

**Independent Test**: Submit ambiguous query, verify system displays uncertainty indicator and shows range of agent positions

### Tests for User Story 2 (TDD)

- [ ] T044 [P] [US2] Write unit test for consensus threshold not met scenario in internal/agents/council/consensus_test.go
- [ ] T045 [P] [US2] Write unit test for insufficient knowledge detection in internal/agents/council/evidence_test.go
- [ ] T046 [US2] Write integration test for partial consensus display in internal/agents/council/coordinator_test.go

### Backend Implementation for User Story 2

- [ ] T047 [US2] Implement uncertainty calculation in consensus algorithm in internal/agents/council/consensus.go
- [ ] T048 [US2] Implement evidence relevance scoring for "insufficient knowledge" detection in internal/agents/council/evidence.go
- [ ] T049 [US2] Add uncertain status handling to Coordinator in internal/agents/council/coordinator.go
- [ ] T050 [US2] Update GET /deliberations/{id} to include dissenting agent info in internal/api/handlers/council.go

### Frontend Implementation for User Story 2

- [ ] T051 [US2] Create UncertaintyDisplay component showing agent position range in frontend/src/components/council/UncertaintyDisplay.tsx
- [ ] T052 [US2] Update ResponseDisplay to show partial consensus indicators in frontend/src/components/council/ResponseDisplay.tsx
- [ ] T053 [US2] Add uncertainty i18n keys in frontend/public/locales/{en,ar}/council.json

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Knowledge Graph Evidence Exploration (Priority: P2)

**Goal**: Users can explore the evidence trail and Knowledge Graph nodes that supported each response

**Independent Test**: Submit query, retrieve evidence trail via API, verify nodes and traversal paths are displayed

### Tests for User Story 3 (TDD)

- [ ] T054 [P] [US3] Write unit test for evidence trail retrieval in internal/agents/council/evidence_test.go
- [ ] T055 [P] [US3] Write unit test for Knowledge Graph node expansion in internal/warehouse/knowledge_graph_test.go
- [ ] T056 [US3] Write integration test for evidence API endpoint in internal/api/handlers/council_test.go

### Backend Implementation for User Story 3

- [ ] T057 [US3] Implement GET /deliberations/{id}/evidence handler in internal/api/handlers/council.go
- [ ] T058 [US3] Implement GET /deliberations/{id}/evidence/nodes/{nodeId} handler in internal/api/handlers/council.go
- [ ] T059 [US3] Add node relationship expansion to Knowledge Graph repository in internal/warehouse/knowledge_graph.go

### Frontend Implementation for User Story 3

- [ ] T060 [US3] Create EvidenceExplorer component with graph visualization in frontend/src/components/council/EvidenceExplorer.tsx
- [ ] T061 [US3] Create EvidenceNode component for expandable node display in frontend/src/components/council/EvidenceNode.tsx
- [ ] T062 [US3] Create TraversalPath component showing reasoning chain in frontend/src/components/council/TraversalPath.tsx
- [ ] T063 [US3] Add evidence exploration i18n keys in frontend/public/locales/{en,ar}/council.json

**Checkpoint**: At this point, Users can explore evidence supporting their responses

---

## Phase 6: User Story 4 - Response Accuracy Audit Trail (Priority: P3)

**Goal**: Administrators can review all deliberations, flag potential hallucinations, and maintain compliance records

**Independent Test**: Access audit endpoint as admin, verify complete deliberation records; verify non-admin can only see own records

### Tests for User Story 4 (TDD)

- [ ] T064 [P] [US4] Write unit test for audit trail creation in internal/agents/council/repository_test.go
- [ ] T065 [P] [US4] Write unit test for RBAC audit access (admin vs user) in internal/api/handlers/council_test.go
- [ ] T066 [US4] Write integration test for flag deliberation flow in internal/api/handlers/council_test.go

### Backend Implementation for User Story 4

- [ ] T067 [US4] Implement audit entry creation on deliberation completion in internal/agents/council/coordinator.go
- [ ] T068 [US4] Implement GET /audit/deliberations handler (admin only) in internal/api/handlers/council.go
- [ ] T069 [US4] Implement POST /audit/deliberations/{id}/flag handler in internal/api/handlers/council.go
- [ ] T070 [US4] Add OPA policy for audit access (admin role required) in policies/council.rego

### Frontend Implementation for User Story 4

- [ ] T071 [US4] Create AuditTrailList component for admin view in frontend/src/components/council/AuditTrailList.tsx
- [ ] T072 [US4] Create AuditDetail component with full deliberation view in frontend/src/components/council/AuditDetail.tsx
- [ ] T073 [US4] Create FlagDialog component for hallucination reporting in frontend/src/components/council/FlagDialog.tsx
- [ ] T074 [US4] Add audit i18n keys in frontend/public/locales/{en,ar}/council.json

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: Health Monitoring & Observability

**Purpose**: System health endpoints and metrics for production readiness

- [ ] T075 [P] Implement agent heartbeat publishing to NATS in internal/agents/council/health.go
- [ ] T076 [P] Implement circuit breaker for agent failures in internal/agents/council/agent.go
- [ ] T077 Implement GET /health system health endpoint in internal/api/handlers/council.go
- [ ] T078 Implement GET /health/agents endpoint in internal/api/handlers/council.go
- [ ] T079 Add Prometheus metrics for consensus (agreement rates, latency, uncertainty frequency) in internal/agents/council/metrics.go
- [ ] T080 Implement Knowledge Graph health check (30s interval) in internal/warehouse/knowledge_graph.go

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T081 [P] Add Knowledge Graph graceful degradation (503 on unavailability) in internal/agents/council/coordinator.go
- [ ] T082 [P] Implement evidence cache expiration (5-minute TTL) in internal/cache/evidence_cache.go
- [ ] T083 Verify all i18n translations complete (EN/AR) in frontend/public/locales/
- [ ] T084 Run quickstart.md validation - verify all code examples work
- [ ] T085 Add API documentation comments (godoc) to all exported functions
- [ ] T086 Performance test deliberation latency (<10s for 95% of queries)
- [ ] T087 Security review: verify medisync_readonly role usage, no SQL injection vectors

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1 (P1): Core functionality - implement first
  - US2 (P2): Enhances US1 but independently testable
  - US3 (P2): Independent feature, can parallel with US2
  - US4 (P3): Admin features, depends on US1 deliberation data existing
- **Health & Observability (Phase 7)**: Can start after Phase 2, but should complete before production
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Extends US1 types but independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Uses US1 evidence data but independently testable
- **User Story 4 (P3)**: Depends on US1 deliberation data existing - best after US1 complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD per Constitution)
- Models before services
- Services before handlers
- Backend before frontend (API must exist for frontend to consume)
- Core implementation before integration

### Parallel Opportunities

- All Setup tasks (T001-T005) can run in parallel
- All migration files (T006-T009) can run in parallel
- All type definitions (T011-T017) can run in parallel
- All repository implementations (T018-T020) can run in parallel
- All US1 tests (T023-T026) can run in parallel
- All US2 tests (T044-T045) can run in parallel
- US2 and US3 can be developed in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task T023: "Write unit test for consensus algorithm"
Task T024: "Write unit test for semantic equivalence detection"
Task T025: "Write unit test for Graph-of-Thoughts retrieval"
Task T026: "Write unit test for agent instance timeout handling"

# Launch backend implementations that don't depend on each other:
Task T028: "Implement AgentInstance wrapper"
Task T029: "Implement semantic equivalence detection"

# Launch frontend tasks that don't depend on each other:
Task T038: "Create Council API client"
Task T039: "Create useCouncil React hook"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready - users can already submit queries and get consensus responses

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo (Uncertainty signaling)
4. Add User Story 3 → Test independently → Deploy/Demo (Evidence exploration)
5. Add User Story 4 → Test independently → Deploy/Demo (Audit compliance)
6. Add Health + Polish → Production ready

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (MVP - highest priority)
   - Developer B: User Story 3 (Evidence - can parallel)
3. After US1 complete:
   - Developer A: User Story 2 (Uncertainty)
   - Developer B: User Story 4 (Audit)
4. Stories complete and integrate independently

---

## Task Summary

| Phase | Task Count | Parallel Tasks | Description |
|-------|------------|----------------|-------------|
| 1. Setup | 5 | 4 | Project structure and configuration |
| 2. Foundational | 17 | 14 | Database, types, repositories, policies |
| 3. US1 (P1) 🎯 | 21 | 8 | Core consensus - MVP |
| 4. US2 (P2) | 10 | 3 | Uncertainty signaling |
| 5. US3 (P2) | 10 | 3 | Evidence exploration |
| 6. US4 (P3) | 11 | 3 | Audit trail |
| 7. Health | 6 | 2 | Observability |
| 8. Polish | 7 | 2 | Cross-cutting concerns |
| **Total** | **87** | **39** | |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Tests written FIRST per TDD (Constitution Principle V)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
