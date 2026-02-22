# Research: Council of AIs Consensus System

**Feature**: 001-council-ai-consensus
**Date**: 2026-02-22

## Research Questions

### RQ-001: Multi-Agent Consensus Algorithm

**Decision**: Weighted voting with semantic similarity clustering

**Rationale**:
- Simple majority voting fails when agents produce semantically equivalent but syntactically different responses
- Weighted voting allows incorporating agent confidence scores into consensus calculation
- Semantic clustering groups equivalent responses before counting votes
- Aligns with MediSync's existing Agent ADK framework for multi-agent orchestration

**Alternatives Considered**:
- **Unanimous consensus**: Too restrictive, would block valid responses too frequently
- **Byzantine fault tolerance**: Overkill for this use case; agents are trusted but may hallucinate
- **Ranking-based aggregation**: More complex, harder to explain to users

**Implementation Notes**:
- Use Agent ADK's built-in agent coordination primitives
- Semantic similarity via embedding comparison (pgvector cosine similarity)
- 95% similarity threshold for equivalence (per FR-015)

---

### RQ-002: Graph-of-Thoughts Retrieval Pattern

**Decision**: Multi-hop graph traversal with relevance scoring

**Rationale**:
- Medical knowledge is highly interconnected (diseases → symptoms → treatments → medications)
- Multi-hop traversal captures reasoning chains, not just direct matches
- Relevance scoring filters noise from large knowledge graphs
- pgvector supports efficient similarity search on graph embeddings

**Alternatives Considered**:
- **Single-hop retrieval**: Misses indirect relationships critical for medical reasoning
- **Full graph traversal**: Too expensive; latency exceeds 10s target
- **Subgraph extraction + LLM reasoning**: Adds unnecessary complexity

**Implementation Notes**:
- Store Knowledge Graph nodes with embeddings in PostgreSQL
- Use recursive CTEs for multi-hop traversal (max 3 hops)
- Cache traversal results for 5 minutes (per FR-017)
- Integrate with MediSync's existing medisync_readonly role

---

### RQ-003: Semantic Equivalence Detection

**Decision**: Embedding-based similarity with 95% threshold

**Rationale**:
- Embedding models capture semantic meaning regardless of phrasing
- 95% threshold balances false positives (grouping non-equivalent responses) vs false negatives (splitting equivalent responses)
- Can leverage same embeddings used for Knowledge Graph retrieval
- Aligns with FR-015 requirements

**Alternatives Considered**:
- **Exact string matching**: Fails on paraphrasing
- **LLM-based comparison**: Adds latency and cost; less deterministic
- **Lexical overlap (TF-IDF)**: Fails on synonyms and rephrasing

**Implementation Notes**:
- Pre-compute embeddings for agent responses
- Use pgvector's cosine similarity operator (`<=>`)
- Batch comparison for efficiency with 3+ agents

---

### RQ-004: Healthcare Knowledge Graph Schema

**Decision**: Entity-Relationship model with typed edges and provenance

**Rationale**:
- Healthcare domain requires tracking source provenance for trust
- Typed edges enable domain-specific reasoning (treats, causes, contraindicates)
- Entity types align with MediSync's existing data model (patients, medications, procedures)
- Supports both clinical and financial subdomains per edge case requirements

**Alternatives Considered**:
- **RDF/OWL ontology**: Standard but adds semantic web complexity
- **Property graph (Neo4j)**: Not in approved technology stack
- **Flat document store**: Loses relationship semantics

**Implementation Notes**:
- Store in PostgreSQL with adjacency list pattern
- Node types: Concept, Medication, Procedure, Condition, Organization
- Edge types: TREATS, CAUSES, CONTRAINDICATES, RELATED_TO, SUBSUMES
- Include `source`, `confidence`, `last_verified` metadata per node

---

### RQ-005: Audit Trail Storage Pattern

**Decision**: Append-only ledger with partitioning by date

**Rationale**:
- Append-only ensures immutability for compliance (HIPAA)
- Date partitioning supports 7-year retention with efficient archival
- PostgreSQL native partitioning (pg_partman) handles lifecycle
- Queryable for evidence exploration and admin review

**Alternatives Considered**:
- **Event sourcing with snapshots**: More complex; overkill for read-heavy audit access
- **Dedicated time-series DB (TimescaleDB)**: Not in approved stack
- **Object storage (S3)**: Less queryable; harder to enforce RBAC

**Implementation Notes**:
- Partition by month for efficient querying and archival
- Include deliberation_id as partition key for co-location
- Soft-delete for user-facing deletion, hard delete only after retention period
- Integrate with MediSync's existing audit logging infrastructure

---

### RQ-006: Agent Health Monitoring

**Decision**: Heartbeat-based health checks with circuit breaker pattern

**Rationale**:
- Heartbeats detect failed agents quickly (30s interval per NFR-004)
- Circuit breaker prevents cascading failures when agents are unhealthy
- Aligns with FR-009 graceful degradation requirements
- NATS can distribute health status across the cluster

**Alternatives Considered**:
- **Passive failure detection**: Slower; delays graceful degradation
- **Redundant agent instances**: Increases cost; complexity
- **External orchestration (K8s)**: Doesn't capture agent-level health

**Implementation Notes**:
- Each agent publishes heartbeat to NATS subject `agent.health.<agent_id>`
- Council coordinator tracks quorum (minimum 2 healthy agents)
- Circuit breaker opens after 3 consecutive failures
- Alert via NATS when quorum at risk

---

## Technology Decisions Summary

| Concern | Decision | Rationale |
|---------|----------|-----------|
| Consensus Algorithm | Weighted voting + semantic clustering | Handles equivalence, integrates confidence |
| Graph Retrieval | Multi-hop traversal with pgvector | Captures reasoning chains, efficient |
| Equivalence Detection | Embedding similarity (95%) | Deterministic, reuses infrastructure |
| Knowledge Graph | PostgreSQL adjacency list | In approved stack, queryable |
| Audit Storage | Partitioned append-only | Compliance, retention support |
| Health Monitoring | Heartbeat + circuit breaker | Fast detection, prevents cascade |

## Open Questions Resolved

All NEEDS CLARIFICATION items from Technical Context have been resolved through this research phase.
