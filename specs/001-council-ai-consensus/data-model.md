# Data Model: Council of AIs Consensus System

**Feature**: 001-council-ai-consensus
**Date**: 2026-02-22

## Entity Relationship Diagram

```
┌─────────────────────┐       ┌─────────────────────┐
│ CouncilDeliberation │       │    AgentInstance    │
├─────────────────────┤       ├─────────────────────┤
│ id (PK)             │       │ id (PK)             │
│ query_text          │       │ name                │
│ query_hash          │       │ health_status       │
│ user_id (FK)        │       │ last_heartbeat      │
│ status              │───────│ config              │
│ consensus_threshold │       │ created_at          │
│ final_response      │       └─────────────────────┘
│ confidence_score    │                  │
│ created_at          │                  │ participates in
│ completed_at        │                  ▼
└─────────────────────┘       ┌─────────────────────┐
          │                   │  AgentResponse      │
          │                   ├─────────────────────┤
          │                   │ id (PK)             │
          │                   │ deliberation_id(FK) │
          │                   │ agent_id (FK)       │
          │                   │ response_text       │
          │                   │ evidence_ids[]      │
          │                   │ confidence          │
          │                   │ created_at          │
          │                   └─────────────────────┘
          │
          │ produces
          ▼
┌─────────────────────┐       ┌─────────────────────┐
│  ConsensusRecord    │       │   EvidenceTrail     │
├─────────────────────┤       ├─────────────────────┤
│ id (PK)             │       │ id (PK)             │
│ deliberation_id(FK) │       │ deliberation_id(FK) │
│ agreement_score     │       │ node_ids[] (FK)     │
│ equivalence_groups  │       │ traversal_path      │
│ threshold_met       │       │ relevance_scores{}  │
│ dissenting_agents   │       │ cached_at           │
│ created_at          │       │ expires_at          │
└─────────────────────┘       └─────────────────────┘
                                        │
                                        │ references
                                        ▼
                              ┌─────────────────────┐
                              │ KnowledgeGraphNode  │
                              ├─────────────────────┤
                              │ id (PK)             │
                              │ node_type           │
                              │ concept             │
                              │ definition          │
                              │ embedding vector    │
                              │ source              │
                              │ confidence          │
                              │ last_verified       │
                              │ edges[] (FK self)   │
                              │ edge_types[]        │
                              │ created_at          │
                              └─────────────────────┘

┌─────────────────────┐
│    AuditEntry       │
├─────────────────────┤
│ id (PK)             │
│ deliberation_id(FK) │
│ user_id (FK)        │
│ action              │
│ details             │
│ created_at          │
│ partition_date      │
└─────────────────────┘
```

## Entity Definitions

### CouncilDeliberation

Primary entity representing a single multi-agent deliberation session.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique deliberation identifier |
| query_text | TEXT | NOT NULL | Original user query |
| query_hash | VARCHAR(64) | NOT NULL, INDEX | SHA-256 hash for deduplication |
| user_id | UUID | FK → users.id, NOT NULL | User who submitted query |
| status | ENUM | NOT NULL, DEFAULT 'pending' | pending, deliberating, consensus, uncertain, failed |
| consensus_threshold | DECIMAL(3,2) | NOT NULL, DEFAULT 0.80 | Required agreement level |
| final_response | TEXT | NULLABLE | Consensus response text |
| confidence_score | DECIMAL(5,2) | NULLABLE, CHECK (0-100) | Overall confidence percentage |
| error_message | TEXT | NULLABLE | Error details if failed |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Deliberation start time |
| completed_at | TIMESTAMPTZ | NULLABLE | Deliberation end time |

**State Transitions**:
```
pending → deliberating → consensus → (complete)
                      ↘ uncertain → (complete)
                      ↘ failed → (complete)
```

**Indexes**:
- `idx_deliberation_user_id` on (user_id)
- `idx_deliberation_status` on (status)
- `idx_deliberation_created_at` on (created_at)
- `idx_deliberation_query_hash` on (query_hash)

---

### AgentInstance

Represents an independent AI reasoning unit in the Council.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique agent identifier |
| name | VARCHAR(100) | NOT NULL, UNIQUE | Human-readable agent name |
| health_status | ENUM | NOT NULL, DEFAULT 'healthy' | healthy, degraded, failed |
| last_heartbeat | TIMESTAMPTZ | NOT NULL | Last health check timestamp |
| config | JSONB | NOT NULL, DEFAULT '{}' | Agent configuration (model, temperature, etc.) |
| timeout_seconds | INTEGER | NOT NULL, DEFAULT 3 | Response timeout |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Registration timestamp |

**Health Status Rules**:
- `healthy`: Last heartbeat within 60 seconds
- `degraded`: Last heartbeat within 120 seconds
- `failed`: Last heartbeat > 120 seconds or explicit failure

---

### AgentResponse

Individual response from a single agent in a deliberation.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique response identifier |
| deliberation_id | UUID | FK → council_deliberations.id, NOT NULL | Parent deliberation |
| agent_id | UUID | FK → agent_instances.id, NOT NULL | Responding agent |
| response_text | TEXT | NOT NULL | Agent's proposed response |
| evidence_ids | UUID[] | NOT NULL, DEFAULT '{}' | Knowledge nodes referenced |
| confidence | DECIMAL(5,2) | NOT NULL, CHECK (0-100) | Agent's confidence in response |
| embedding | vector(1536) | NULLABLE | Response embedding for similarity |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Response timestamp |

**Indexes**:
- `idx_agent_response_deliberation` on (deliberation_id)
- `idx_agent_response_agent` on (agent_id)

---

### ConsensusRecord

Captures the consensus calculation results for a deliberation.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique record identifier |
| deliberation_id | UUID | FK → council_deliberations.id, NOT NULL, UNIQUE | Parent deliberation |
| agreement_score | DECIMAL(5,2) | NOT NULL, CHECK (0-100) | Calculated agreement percentage |
| equivalence_groups | JSONB | NOT NULL, DEFAULT '[]' | Groups of semantically equivalent responses |
| threshold_met | BOOLEAN | NOT NULL | Whether consensus threshold was met |
| dissenting_agents | UUID[] | NOT NULL, DEFAULT '{}' | Agents not in consensus |
| consensus_method | VARCHAR(50) | NOT NULL, DEFAULT 'weighted_vote' | Algorithm used |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Calculation timestamp |

---

### EvidenceTrail

Records the Knowledge Graph traversal path for a deliberation.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique trail identifier |
| deliberation_id | UUID | FK → council_deliberations.id, NOT NULL, UNIQUE | Parent deliberation |
| node_ids | UUID[] | NOT NULL | Ordered list of traversed nodes |
| traversal_path | JSONB | NOT NULL | Full path with edges and scores |
| relevance_scores | JSONB | NOT NULL, DEFAULT '{}' | Node ID → relevance score mapping |
| hop_count | INTEGER | NOT NULL, DEFAULT 0 | Number of hops in traversal |
| cached_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | When evidence was cached |
| expires_at | TIMESTAMPTZ | NOT NULL | Cache expiration (5 minutes) |

**Cache Rules**:
- Evidence valid for 5 minutes (FR-017)
- Expired trails must be re-fetched from Knowledge Graph

---

### KnowledgeGraphNode

A unit of verified medical/healthcare knowledge.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique node identifier |
| node_type | ENUM | NOT NULL | concept, medication, procedure, condition, organization |
| concept | VARCHAR(255) | NOT NULL | Primary concept name |
| definition | TEXT | NOT NULL | Full definition/description |
| embedding | vector(1536) | NOT NULL | Semantic embedding for similarity |
| source | VARCHAR(255) | NOT NULL | Provenance (e.g., "ICD-10", "RxNorm") |
| source_id | VARCHAR(255) | NULLABLE | External system reference |
| confidence | DECIMAL(5,2) | NOT NULL, DEFAULT 100.00 | Knowledge reliability score |
| last_verified | TIMESTAMPTZ | NOT NULL | Last verification date |
| edges | UUID[] | NOT NULL, DEFAULT '{}' | Connected node IDs |
| edge_types | VARCHAR(50)[] | NOT NULL, DEFAULT '{}' | Edge types parallel to edges |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Node creation time |
| updated_at | TIMESTAMPTZ | NOT NULL | Last update time |

**Edge Types**:
- `TREATS`: Treatment relationship
- `CAUSES`: Causal relationship
- `CONTRAINDICATES`: Safety warning
- `RELATED_TO`: General association
- `SUBSUMES`: Hierarchical relationship
- `PART_OF`: Composition relationship

**Indexes**:
- `idx_kg_node_type` on (node_type)
- `idx_kg_node_concept` on (concept)
- `idx_kg_node_embedding` on (embedding) USING ivfflat (vector_cosine_ops)

---

### AuditEntry

Immutable audit log for compliance (HIPAA).

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK, NOT NULL | Unique entry identifier |
| deliberation_id | UUID | FK → council_deliberations.id, NOT NULL | Related deliberation |
| user_id | UUID | FK → users.id, NOT NULL | Acting user |
| action | VARCHAR(100) | NOT NULL | Action type (query, review, flag, export) |
| details | JSONB | NOT NULL, DEFAULT '{}' | Full action context |
| ip_address | INET | NULLABLE | Client IP for access logs |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Action timestamp |
| partition_date | DATE | NOT NULL, DEFAULT CURRENT_DATE | Partition key |

**Partitioning**:
- Partitioned by month using `partition_date`
- Retain for 7 years (HIPAA compliance)
- Archive to cold storage after 1 year

**Indexes**:
- `idx_audit_deliberation` on (deliberation_id)
- `idx_audit_user` on (user_id)
- `idx_audit_created_at` on (created_at)

---

## Validation Rules

### Deliberation Status Transitions

```
VALID TRANSITIONS:
  pending → deliberating
  deliberating → consensus
  deliberating → uncertain
  deliberating → failed
  consensus → (terminal)
  uncertain → (terminal)
  failed → (terminal)

INVALID TRANSITIONS:
  Any → pending (no backwards)
  consensus → uncertain (no reconsideration)
  failed → deliberating (no retry)
```

### Consensus Score Calculation

```
consensus_score = (
  Σ (agent_confidence × equivalence_weight)
  / Σ equivalence_weight
) × 100

WHERE:
  - equivalence_weight = count of semantically equivalent responses
  - 95% similarity threshold for equivalence
  - minimum 2 agents required for valid consensus
```

### Cache Expiration

```
evidence_valid = (
  expires_at > NOW()
  AND cached_at >= (NOW() - INTERVAL '5 minutes')
)
```

## Data Volume Estimates

| Entity | Growth Rate | 7-Year Volume | Storage Estimate |
|--------|-------------|---------------|------------------|
| CouncilDeliberation | 1,000/day | 2.5M records | ~5 GB |
| AgentResponse | 3,000/day | 7.7M records | ~15 GB |
| ConsensusRecord | 1,000/day | 2.5M records | ~2 GB |
| EvidenceTrail | 1,000/day | 2.5M records | ~10 GB |
| AuditEntry | 2,000/day | 5.1M records | ~8 GB |
| KnowledgeGraphNode | 100/month | 8,400 records | ~500 MB |

**Total 7-Year Estimate**: ~40 GB (before compression)
