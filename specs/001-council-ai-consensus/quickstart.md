# Quickstart: Council of AIs Consensus System

This guide provides a rapid onboarding for developers implementing the Council of AIs consensus system.

## Prerequisites

- Go 1.26+ installed
- PostgreSQL 18.2+ with pgvector extension
- Redis 8+ for caching
- NATS JetStream running
- Keycloak for authentication
- Access to MediSync monorepo

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Layer                                │
│                    (go-chi HTTP handlers)                        │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Council Coordinator                           │
│         (orchestrates agent instances, manages consensus)        │
└───────────────────────────┬─────────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            ▼               ▼               ▼
     ┌──────────┐    ┌──────────┐    ┌──────────┐
     │ Agent 1  │    │ Agent 2  │    │ Agent 3  │
     │(Genkit)  │    │(Genkit)  │    │(Genkit)  │
     └────┬─────┘    └────┬─────┘    └────┬─────┘
          │               │               │
          └───────────────┼───────────────┘
                          ▼
              ┌───────────────────────┐
              │   Knowledge Graph     │
              │   (PostgreSQL +       │
              │    pgvector)          │
              └───────────────────────┘
```

## Quick Setup

### 1. Database Migration

```bash
# Run migrations for Council tables
cd /path/to/medisync
go run ./cmd/migrate --target council_ai_consensus
```

### 2. Environment Variables

```bash
# Add to .env or environment
COUNCIL_MIN_AGENTS=3
COUNCIL_DEFAULT_THRESHOLD=0.80
COUNCIL_AGENT_TIMEOUT=3s
COUNCIL_CACHE_TTL=5m
COUNCIL_HEALTH_CHECK_INTERVAL=30s
```

### 3. Start Services

```bash
# Start infrastructure
docker-compose up -d nats postgres redis keycloak

# Start API server with Council routes
go run ./cmd/api
```

## Key Implementation Files

### Backend (Go)

| File | Purpose |
|------|---------|
| `internal/agents/council/coordinator.go` | Council deliberation orchestration |
| `internal/agents/council/consensus.go` | Consensus algorithm implementation |
| `internal/agents/council/agent.go` | Individual agent instance wrapper |
| `internal/agents/council/evidence.go` | Graph-of-Thoughts retrieval |
| `internal/agents/council/semantic.go` | Semantic equivalence detection |
| `internal/api/handlers/council.go` | HTTP handlers for API |
| `internal/warehouse/knowledge_graph.go` | Knowledge Graph repository |

### Frontend (React)

| File | Purpose |
|------|---------|
| `frontend/src/services/councilService.ts` | API client |
| `frontend/src/hooks/useCouncil.ts` | React hook for deliberations |
| `frontend/src/components/council/QueryInput.tsx` | Query submission UI |
| `frontend/src/components/council/ResponseDisplay.tsx` | Consensus response display |
| `frontend/src/components/council/EvidenceExplorer.tsx` | Evidence trail visualization |
| `frontend/src/components/council/ConfidenceIndicator.tsx` | Confidence score display |

## Core Flows

### Submit Query Flow

```go
// 1. Handler receives query
func (h *Handler) CreateDeliberation(w http.ResponseWriter, r *http.Request) {
    var req CreateDeliberationRequest
    json.NewDecoder(r.Body).Decode(&req)

    // 2. Create deliberation record
    deliberation := coordinator.NewDeliberation(req.Query, req.Threshold)

    // 3. Start async deliberation
    go coordinator.Deliberate(ctx, deliberation)

    // 4. Return immediately with ID
    json.NewEncoder(w).Encode(DeliberationResponse{ID: deliberation.ID})
}
```

### Consensus Algorithm

```go
func (c *Coordinator) CalculateConsensus(responses []AgentResponse) ConsensusResult {
    // 1. Group semantically equivalent responses
    groups := c.semanticClustering(responses)

    // 2. Calculate weighted agreement
    agreement := c.weightedVote(groups)

    // 3. Check threshold
    if agreement >= c.threshold {
        return ConsensusResult{Met: true, Score: agreement}
    }

    return ConsensusResult{Met: false, Score: agreement}
}
```

### Graph-of-Thoughts Retrieval

```go
func (r *Retriever) Traverse(ctx context.Context, query string) (*EvidenceTrail, error) {
    // 1. Embed query
    queryEmbedding := r.embedder.Embed(query)

    // 2. Find initial nodes via similarity
    initialNodes := r.repo.FindSimilar(ctx, queryEmbedding, limit=10)

    // 3. Multi-hop traversal (max 3 hops)
    trail := r.multiHopTraverse(ctx, initialNodes, maxHops=3)

    // 4. Cache with 5-minute TTL
    r.cache.Set(trail.ID, trail, 5*time.Minute)

    return trail, nil
}
```

## Testing

### Unit Tests

```bash
# Run Council package tests
go test ./internal/agents/council/... -v

# Run with coverage
go test ./internal/agents/council/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run against local services
go test ./internal/agents/council/... -tags=integration -v
```

### API Tests

```bash
# Using the OpenAPI contract
curl -X POST http://localhost:8080/api/v1/council/deliberations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "What were patient visits last month?"}'
```

## Common Tasks

### Add New Agent Type

1. Implement `Agent` interface in `internal/agents/council/agents/`
2. Register in `coordinator.go` agent pool
3. Add configuration in `config/agents.yaml`
4. Add unit tests for agent behavior

### Modify Consensus Threshold

```go
// Per-request threshold in API call
{
  "query": "...",
  "consensusThreshold": 0.90
}

// Or update default in config
COUNCIL_DEFAULT_THRESHOLD=0.85
```

### Query Evidence Trail

```bash
# Get evidence for a deliberation
curl http://localhost:8080/api/v1/council/deliberations/{id}/evidence \
  -H "Authorization: Bearer $TOKEN"
```

## Monitoring

### Key Metrics

- `council_deliberations_total` - Total deliberations started
- `council_consensus_achieved` - Deliberations reaching consensus
- `council_uncertainty_rate` - Deliberations with no consensus
- `council_latency_seconds` - Deliberation duration histogram
- `council_agent_health` - Agent health gauge by ID

### Health Check

```bash
# System health endpoint
curl http://localhost:8080/api/v1/council/health
```

## Troubleshooting

### Knowledge Graph Unavailable

```
Error: Knowledge Graph unavailable
Solution: Check PostgreSQL connection, verify pgvector extension loaded
```

### Quorum Not Met

```
Error: Insufficient healthy agents (1/3)
Solution: Check agent health endpoint, restart failed agents
```

### High Latency

```
Issue: Deliberations exceeding 10s target
Checks:
  1. Knowledge Graph query performance (check slow query log)
  2. Agent response times (check agent health metrics)
  3. Semantic clustering efficiency (check embedding batch size)
```

## Next Steps

1. Review [data-model.md](./data-model.md) for entity details
2. Review [contracts/openapi.yaml](./contracts/openapi.yaml) for API specs
3. Review [research.md](./research.md) for design decisions
4. Run `/speckit.tasks` to generate implementation tasks
