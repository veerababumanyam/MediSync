# Agentic RAG Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement an Agentic RAG system with query decomposition, parallel retrieval, and multi-hop capabilities to solve semantic collapse at 50K+ document scale.

**Architecture:** Four specialized indices (clauses, parties, financial, summaries) with a Query Planner Agent that decomposes complex queries into sub-queries, executes parallel retrieval, aggregates evidence, and synthesizes responses with mandatory citations.

**Tech Stack:** Go 1.26, PostgreSQL 18.2 + pgvector, Ollama (Llama 3.1 + BGE-large), NATS JetStream, Genkit

---

## Prerequisites

Before starting, ensure:
- [ ] PostgreSQL 18.2 with pgvector extension installed
- [ ] Ollama running locally with `llama3.1:8b` and `bge-large-en-v1.5` models
- [ ] NATS JetStream running
- [ ] Go 1.26 installed

---

## Phase 1: Database Schema & Specialized Indices

### Task 1.1: Create Contract Clauses Index

**Files:**
- Create: `migrations/060_create_contract_indices.up.sql`
- Create: `migrations/060_create_contract_indices.down.sql`

**Step 1: Write the migration (up)**

```sql
-- migrations/060_create_contract_indices.up.sql

-- Enable pgvector if not already enabled
CREATE EXTENSION IF NOT EXISTS vector;

-- Index 1: CLAUSES (for compliance/audit queries)
CREATE TABLE contract_clauses (
    clause_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    clause_type VARCHAR(100) NOT NULL,
    clause_title TEXT,
    clause_text TEXT NOT NULL,
    embedding vector(1024),
    page_number INT,
    position_in_doc INT,
    parent_section TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Hybrid search support
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', clause_text)) STORED
);

-- Indexes for clauses
CREATE INDEX idx_clauses_embedding ON contract_clauses
    USING hnsw (embedding vector_cosine_ops) WITH (m=16, ef_construction=64);
CREATE INDEX idx_clauses_fts ON contract_clauses USING gin(search_vector);
CREATE INDEX idx_clauses_type ON contract_clauses(clause_type);
CREATE INDEX idx_clauses_document ON contract_clauses(document_id);

-- Index 2: PARTIES (for relationship queries)
CREATE TABLE contract_parties (
    party_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    party_name TEXT NOT NULL,
    party_name_normalized TEXT,
    party_role VARCHAR(50),
    jurisdiction VARCHAR(100),
    embedding vector(1024),
    entity_type VARCHAR(50),
    identifiers JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- Auto-normalize party name for fuzzy matching
    CONSTRAINT chk_party_name CHECK (length(party_name) > 0)
);

CREATE INDEX idx_parties_embedding ON contract_parties
    USING hnsw (embedding vector_cosine_ops) WITH (m=16, ef_construction=64);
CREATE INDEX idx_parties_name ON contract_parties(party_name_normalized text_pattern_ops);
CREATE INDEX idx_parties_document ON contract_parties(document_id);

-- Index 3: FINANCIAL (for aggregation/analysis queries)
CREATE TABLE contract_financials (
    financial_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    metric_type VARCHAR(100) NOT NULL,
    amount NUMERIC(18,2),
    currency VARCHAR(3) DEFAULT 'USD',
    amount_normalized NUMERIC(18,2),
    effective_date DATE,
    expiry_date DATE,
    payment_terms TEXT,
    embedding vector(1024),
    fiscal_year INT,
    fiscal_quarter INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_financial_embedding ON contract_financials
    USING hnsw (embedding vector_cosine_ops) WITH (m=16, ef_construction=64);
CREATE INDEX idx_financial_amount ON contract_financials(amount_normalized DESC NULLS LAST);
CREATE INDEX idx_financial_dates ON contract_financials(effective_date, expiry_date);
CREATE INDEX idx_financial_document ON contract_financials(document_id);
CREATE INDEX idx_financial_type ON contract_financials(metric_type);

-- Index 4: DOCUMENT_SUMMARIES (for high-level queries)
CREATE TABLE contract_summaries (
    summary_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    summary_type VARCHAR(50) NOT NULL DEFAULT 'executive',
    summary_text TEXT NOT NULL,
    key_entities JSONB DEFAULT '{}',
    key_risks JSONB DEFAULT '[]',
    embedding vector(1024),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(document_id, summary_type)
);

CREATE INDEX idx_summaries_embedding ON contract_summaries
    USING hnsw (embedding vector_cosine_ops) WITH (m=16, ef_construction=64);
CREATE INDEX idx_summaries_document ON contract_summaries(document_id);
CREATE INDEX idx_summaries_type ON contract_summaries(summary_type);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_clauses_updated_at
    BEFORE UPDATE ON contract_clauses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_summaries_updated_at
    BEFORE UPDATE ON contract_summaries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant readonly access to medisync_readonly role
GRANT SELECT ON contract_clauses TO medisync_readonly;
GRANT SELECT ON contract_parties TO medisync_readonly;
GRANT SELECT ON contract_financials TO medisync_readonly;
GRANT SELECT ON contract_summaries TO medisync_readonly;

-- Grant full access to app role
GRANT ALL ON contract_clauses TO medisync_app;
GRANT ALL ON contract_parties TO medisync_app;
GRANT ALL ON contract_financials TO medisync_app;
GRANT ALL ON contract_summaries TO medisync_app;
```

**Step 2: Write the migration (down)**

```sql
-- migrations/060_create_contract_indices.down.sql

DROP TRIGGER IF EXISTS update_summaries_updated_at ON contract_summaries;
DROP TRIGGER IF EXISTS update_clauses_updated_at ON contract_clauses;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS contract_summaries CASCADE;
DROP TABLE IF EXISTS contract_financials CASCADE;
DROP TABLE IF EXISTS contract_parties CASCADE;
DROP TABLE IF EXISTS contract_clauses CASCADE;
```

**Step 3: Run the migration**

```bash
go run ./cmd/migrate up
```

**Expected output:**
```
Applying migration 060_create_contract_indices.up.sql...
Migration applied successfully.
```

**Step 4: Verify tables exist**

```bash
psql -d medisync -c "\dt contract_*"
```

**Expected output:**
```
                List of relations
 Schema |         Name          | Type  |   Owner
--------+-----------------------+-------+------------
 public | contract_clauses      | table | medisync
 public | contract_parties      | table | medisync
 public | contract_financials   | table | medisync
 public | contract_summaries    | table | medisync
```

**Step 5: Commit**

```bash
git add migrations/060_create_contract_indices.up.sql migrations/060_create_contract_indices.down.sql
git commit -m "feat(rag): create specialized contract indices for agentic retrieval"
```

---

### Task 1.2: Create RAG Query Audit Log

**Files:**
- Create: `migrations/061_create_rag_audit_log.up.sql`
- Create: `migrations/061_create_rag_audit_log.down.sql`

**Step 1: Write the migration (up)**

```sql
-- migrations/061_create_rag_audit_log.up.sql

-- Audit log for RAG queries
CREATE TABLE rag_query_log (
    query_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    tenant_id UUID REFERENCES tenants(id),
    original_query TEXT NOT NULL,
    query_type VARCHAR(50),
    sub_queries JSONB DEFAULT '[]',
    execution_plan JSONB,
    total_hops INT DEFAULT 1,
    contracts_matched INT DEFAULT 0,
    confidence_score NUMERIC(5,4),
    latency_ms BIGINT,
    status VARCHAR(30) DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Audit log for retrieved evidence
CREATE TABLE rag_evidence_log (
    evidence_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    query_id UUID NOT NULL REFERENCES rag_query_log(query_id) ON DELETE CASCADE,
    document_id UUID REFERENCES documents(id),
    chunk_ids JSONB DEFAULT '[]',
    relevance_score NUMERIC(5,4),
    was_cited BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for audit queries
CREATE INDEX idx_rag_query_user ON rag_query_log(user_id);
CREATE INDEX idx_rag_query_tenant ON rag_query_log(tenant_id);
CREATE INDEX idx_rag_query_created ON rag_query_log(created_at DESC);
CREATE INDEX idx_rag_query_status ON rag_query_log(status);
CREATE INDEX idx_rag_evidence_query ON rag_evidence_log(query_id);

-- Grant access
GRANT SELECT ON rag_query_log TO medisync_readonly;
GRANT SELECT ON rag_evidence_log TO medisync_readonly;
GRANT ALL ON rag_query_log TO medisync_app;
GRANT ALL ON rag_evidence_log TO medisync_app;
```

**Step 2: Write the migration (down)**

```sql
-- migrations/061_create_rag_audit_log.down.sql

DROP TABLE IF EXISTS rag_evidence_log CASCADE;
DROP TABLE IF EXISTS rag_query_log CASCADE;
```

**Step 3: Run the migration**

```bash
go run ./cmd/migrate up
```

**Step 4: Commit**

```bash
git add migrations/061_create_rag_audit_log.up.sql migrations/061_create_rag_audit_log.down.sql
git commit -m "feat(rag): add audit logging for RAG queries and evidence"
```

---

## Phase 2: Core Types and Interfaces

### Task 2.1: Define RAG Domain Types

**Files:**
- Create: `internal/agents/module_rag/types.go`

**Step 1: Write the types**

```go
// internal/agents/module_rag/types.go

package module_rag

import (
	"time"

	"github.com/google/uuid"
)

// QueryType represents the classification of a user query
type QueryType int

const (
	QueryTypeLookup     QueryType = iota // "What's the termination clause in contract #123?"
	QueryTypeAggregate                   // "Total value of all ACME contracts"
	QueryTypeCompliance                  // "Which contracts lack force majeure?"
	QueryTypeComparison                  // "Compare liability caps across vendors"
)

func (qt QueryType) String() string {
	switch qt {
	case QueryTypeLookup:
		return "lookup"
	case QueryTypeAggregate:
		return "aggregate"
	case QueryTypeCompliance:
		return "compliance"
	case QueryTypeComparison:
		return "comparison"
	default:
		return "unknown"
	}
}

// IndexType represents which specialized index to query
type IndexType int

const (
	IndexTypeClauses IndexType = iota
	IndexTypeParties
	IndexTypeFinancial
	IndexTypeSummaries
)

func (it IndexType) String() string {
	switch it {
	case IndexTypeClauses:
		return "clauses"
	case IndexTypeParties:
		return "parties"
	case IndexTypeFinancial:
		return "financial"
	case IndexTypeSummaries:
		return "summaries"
	default:
		return "unknown"
	}
}

// RetrievalType specifies how to retrieve from the index
type RetrievalType int

const (
	RetrievalTypeExact    RetrievalType = iota // Exact filter matching
	RetrievalTypeSemantic                      // Vector similarity search
	RetrievalTypeHybrid                        // BM25 + Vector combined
)

// Filter represents a filter condition for retrieval
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, neq, gt, gte, lt, lte, in, contains
	Value    interface{} `json:"value"`
}

// SubQuery represents a decomposed sub-query
type SubQuery struct {
	ID             string         `json:"id"`
	QueryText      string         `json:"query_text"`
	TargetIndex    IndexType      `json:"target_index"`
	Filters        []Filter       `json:"filters"`
	RetrievalType  RetrievalType  `json:"retrieval_type"`
	Dependencies   []string       `json:"dependencies"` // IDs of sub-queries that must complete first
	Aggregation    *Aggregation   `json:"aggregation,omitempty"`
	Limit          int            `json:"limit"`
}

// Aggregation represents an aggregation operation
type Aggregation struct {
	Operation string `json:"operation"` // SUM, COUNT, AVG, MIN, MAX
	Field     string `json:"field"`
	GroupBy   string `json:"group_by,omitempty"`
}

// QueryPlan represents the complete execution plan
type QueryPlan struct {
	QueryID        string       `json:"query_id"`
	OriginalQuery  string       `json:"original_query"`
	QueryType      QueryType    `json:"query_type"`
	SubQueries     []SubQuery   `json:"sub_queries"`
	ExecutionPlan  string       `json:"execution_plan"` // parallel, sequential, parallel_then_merge
	Reasoning      string       `json:"reasoning"`
}

// ChunkMatch represents a retrieved text chunk
type ChunkMatch struct {
	ChunkID      uuid.UUID `json:"chunk_id"`
	DocumentID   uuid.UUID `json:"document_id"`
	IndexSource  string    `json:"index_source"`
	Text         string    `json:"text"`
	Score        float64   `json:"score"`
	PageNumber   int       `json:"page_number,omitempty"`
	Highlights   []Range   `json:"highlights,omitempty"`
}

// Range represents a text range for highlighting
type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// ContractMatch represents aggregated matches for a contract
type ContractMatch struct {
	DocumentID     uuid.UUID     `json:"document_id"`
	DocumentTitle  string        `json:"document_title"`
	MatchScore     float64       `json:"match_score"`
	MatchReason    string        `json:"match_reason"`
	RelevantChunks []ChunkMatch  `json:"relevant_chunks"`
	Metadata       ContractMeta  `json:"metadata"`
}

// ContractMeta contains metadata about a contract
type ContractMeta struct {
	Parties       []string  `json:"parties"`
	EffectiveDate *time.Time `json:"effective_date,omitempty"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	TotalValue    float64   `json:"total_value,omitempty"`
}

// IndexResult represents results from a single index
type IndexResult struct {
	SubQueryID string                   `json:"sub_query_id"`
	Chunks     []ChunkMatch             `json:"chunks"`
	ByContract map[uuid.UUID][]ChunkMatch `json:"by_contract"`
	Error      string                   `json:"error,omitempty"`
	LatencyMs  int64                    `json:"latency_ms"`
}

// AggregatedEvidence represents merged evidence from all sub-queries
type AggregatedEvidence struct {
	QueryID          string           `json:"query_id"`
	ContractMatches  []ContractMatch  `json:"contract_matches"`
	TotalContracts   int              `json:"total_contracts"`
	Confidence       float64          `json:"confidence"`
	EvidenceGaps     []string         `json:"evidence_gaps"`
	SourceCitations  []Citation       `json:"source_citations"`
}

// Citation represents a source citation
type Citation struct {
	CitationID  string    `json:"citation_id"`
	DocumentID  uuid.UUID `json:"document_id"`
	ChunkID     uuid.UUID `json:"chunk_id"`
	DocumentTitle string  `json:"document_title"`
	PageNumber  int       `json:"page_number,omitempty"`
	Text        string    `json:"text"` // Snippet for display
}

// HopAnalysis represents analysis of a single retrieval hop
type HopAnalysis struct {
	Confidence      float64  `json:"confidence"`
	InformationGaps []string `json:"information_gaps"`
	MissingContext  []string `json:"missing_context"`
	RecommendedHops []string `json:"recommended_hops"`
}

// HopResult represents results from a single hop
type HopResult struct {
	HopNumber      int                    `json:"hop_number"`
	SubQueries     []SubQuery             `json:"sub_queries"`
	Results        map[string]*IndexResult `json:"results"`
	Analysis       HopAnalysis            `json:"analysis"`
	ShouldContinue bool                   `json:"should_continue"`
	NextQueries    []SubQuery             `json:"next_queries,omitempty"`
}

// FinalResult represents the complete RAG response
type FinalResult struct {
	QueryID        string       `json:"query_id"`
	Query          string       `json:"query"`
	Synthesis      *Synthesis   `json:"synthesis"`
	HopResults     []*HopResult `json:"hop_results"`
	TotalHops      int          `json:"total_hops"`
	Confidence     float64      `json:"confidence"`
	TotalLatencyMs int64        `json:"total_latency_ms"`
}

// Synthesis represents the synthesized response
type Synthesis struct {
	Summary            string    `json:"summary"`
	Details            string    `json:"details"`
	EvidenceQuality    string    `json:"evidence_quality"`
	ContractsAnalyzed  int       `json:"contracts_analyzed"`
	Confidence         string    `json:"confidence"` // HIGH, MEDIUM, LOW
	Gaps               []string  `json:"gaps"`
	Sources            []Citation `json:"sources"`
}

// RAGError represents an error in the RAG pipeline
type RAGError struct {
	Code        string `json:"code"`
	Stage       string `json:"stage"`
	Message     string `json:"message"`
	Recoverable bool   `json:"recoverable"`
	Fallback    string `json:"fallback,omitempty"`
}

func (e *RAGError) Error() string {
	return e.Message
}

// Error codes
const (
	ErrCodeQueryTooVague    = "E001"
	ErrCodeNoResults        = "E002"
	ErrCodeIndexUnavailable = "E003"
	ErrCodeSynthesisFailed  = "E004"
	ErrCodeTimeoutExceeded  = "E005"
	ErrCodeConfidenceTooLow = "E006"
)
```

**Step 2: Commit**

```bash
git add internal/agents/module_rag/types.go
git commit -m "feat(rag): define core types and interfaces for agentic RAG"
```

---

### Task 2.2: Create RAG Repository Interface

**Files:**
- Create: `internal/agents/module_rag/repository.go`

**Step 1: Write the repository interface**

```go
// internal/agents/module_rag/repository.go

package module_rag

import (
	"context"

	"github.com/google/uuid"
)

// IndexRepository defines the interface for specialized index access
type IndexRepository interface {
	// SearchClauses searches the clauses index
	SearchClauses(ctx context.Context, req ClauseSearchRequest) (*IndexResult, error)

	// SearchParties searches the parties index
	SearchParties(ctx context.Context, req PartySearchRequest) (*IndexResult, error)

	// SearchFinancials searches the financial index
	SearchFinancials(ctx context.Context, req FinancialSearchRequest) (*IndexResult, error)

	// SearchSummaries searches the summaries index
	SearchSummaries(ctx context.Context, req SummarySearchRequest) (*IndexResult, error)

	// GetContractMeta retrieves metadata for a contract
	GetContractMeta(ctx context.Context, documentID uuid.UUID) (*ContractMeta, error)

	// InsertClause inserts a clause into the index
	InsertClause(ctx context.Context, clause *ClauseRecord) error

	// InsertParty inserts a party into the index
	InsertParty(ctx context.Context, party *PartyRecord) error

	// InsertFinancial inserts a financial record into the index
	InsertFinancial(ctx context.Context, financial *FinancialRecord) error

	// InsertSummary inserts a summary into the index
	InsertSummary(ctx context.Context, summary *SummaryRecord) error
}

// ClauseSearchRequest defines search parameters for clauses
type ClauseSearchRequest struct {
	QueryText   string   `json:"query_text"`
	QueryVector []float32 `json:"query_vector"`
	Filters     []Filter `json:"filters"`
	Limit       int      `json:"limit"`
	UseHybrid   bool     `json:"use_hybrid"`
}

// PartySearchRequest defines search parameters for parties
type PartySearchRequest struct {
	QueryText       string   `json:"query_text"`
	QueryVector     []float32 `json:"query_vector"`
	PartyName       string   `json:"party_name,omitempty"`
	PartyRole       string   `json:"party_role,omitempty"`
	DocumentIDs     []uuid.UUID `json:"document_ids,omitempty"`
	Limit           int      `json:"limit"`
}

// FinancialSearchRequest defines search parameters for financials
type FinancialSearchRequest struct {
	MetricType      string     `json:"metric_type,omitempty"`
	MinAmount       float64    `json:"min_amount,omitempty"`
	MaxAmount       float64    `json:"max_amount,omitempty"`
	EffectiveAfter  *time.Time `json:"effective_after,omitempty"`
	ExpiresBefore   *time.Time `json:"expires_before,omitempty"`
	DocumentIDs     []uuid.UUID `json:"document_ids,omitempty"`
	Aggregation     *Aggregation `json:"aggregation,omitempty"`
	Limit           int        `json:"limit"`
}

// SummarySearchRequest defines search parameters for summaries
type SummarySearchRequest struct {
	QueryText   string   `json:"query_text"`
	QueryVector []float32 `json:"query_vector"`
	SummaryType string   `json:"summary_type,omitempty"`
	DocumentIDs []uuid.UUID `json:"document_ids,omitempty"`
	Limit       int      `json:"limit"`
}

// Record types for insertion

type ClauseRecord struct {
	DocumentID    uuid.UUID
	ClauseType    string
	ClauseTitle   string
	ClauseText    string
	Embedding     []float32
	PageNumber    int
	PositionInDoc int
	ParentSection string
	Metadata      map[string]interface{}
}

type PartyRecord struct {
	DocumentID          uuid.UUID
	PartyName           string
	PartyNameNormalized string
	PartyRole           string
	Jurisdiction        string
	Embedding           []float32
	EntityType          string
	Identifiers         map[string]string
}

type FinancialRecord struct {
	DocumentID       uuid.UUID
	MetricType       string
	Amount           float64
	Currency         string
	AmountNormalized float64
	EffectiveDate    *time.Time
	ExpiryDate       *time.Time
	PaymentTerms     string
	Embedding        []float32
}

type SummaryRecord struct {
	DocumentID  uuid.UUID
	SummaryType string
	SummaryText string
	KeyEntities map[string]interface{}
	KeyRisks    []interface{}
	Embedding   []float32
}

// AuditRepository defines the interface for audit logging
type AuditRepository interface {
	// CreateQueryLog creates a new query log entry
	CreateQueryLog(ctx context.Context, log *QueryLogRecord) error

	// UpdateQueryLog updates a query log entry
	UpdateQueryLog(ctx context.Context, queryID uuid.UUID, updates map[string]interface{}) error

	// CreateEvidenceLog creates an evidence log entry
	CreateEvidenceLog(ctx context.Context, log *EvidenceLogRecord) error

	// GetQueryLog retrieves a query log by ID
	GetQueryLog(ctx context.Context, queryID uuid.UUID) (*QueryLogRecord, error)
}

type QueryLogRecord struct {
	QueryID         uuid.UUID
	UserID          uuid.UUID
	TenantID        uuid.UUID
	OriginalQuery   string
	QueryType       string
	SubQueries      []SubQuery
	ExecutionPlan   *QueryPlan
	TotalHops       int
	ContractsMatched int
	ConfidenceScore float64
	LatencyMs       int64
	Status          string
	ErrorMessage    string
	CreatedAt       time.Time
	CompletedAt     *time.Time
}

type EvidenceLogRecord struct {
	EvidenceID     uuid.UUID
	QueryID        uuid.UUID
	DocumentID     uuid.UUID
	ChunkIDs       []uuid.UUID
	RelevanceScore float64
	WasCited       bool
}
```

**Step 2: Commit**

```bash
git add internal/agents/module_rag/repository.go
git commit -m "feat(rag): add repository interfaces for specialized indices"
```

---

## Phase 3: PostgreSQL Repository Implementation

### Task 3.1: Implement Index Repository

**Files:**
- Create: `internal/agents/module_rag/pg_repository.go`
- Create: `internal/agents/module_rag/pg_repository_test.go`

**Step 1: Write the failing test**

```go
// internal/agents/module_rag/pg_repository_test.go

package module_rag

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClauseSearch(t *testing.T) {
	// Skip if no database URL
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://medisync:medisync@localhost:5432/medisync_test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	repo := NewPGIndexRepository(pool)

	t.Run("search by clause type filter", func(t *testing.T) {
		req := ClauseSearchRequest{
			Filters: []Filter{
				{Field: "clause_type", Operator: "eq", Value: "termination"},
			},
			Limit: 10,
		}

		result, err := repo.SearchClauses(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "sq_test", result.SubQueryID)
	})

	t.Run("search with semantic similarity", func(t *testing.T) {
		// Create a dummy embedding (1024 dimensions for BGE-large)
		queryVector := make([]float32, 1024)
		for i := range queryVector {
			queryVector[i] = 0.01
		}

		req := ClauseSearchRequest{
			QueryText:   "termination notice period",
			QueryVector: queryVector,
			Limit:       5,
			UseHybrid:   true,
		}

		result, err := repo.SearchClauses(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestPartySearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://medisync:medisync@localhost:5432/medisync_test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	repo := NewPGIndexRepository(pool)

	t.Run("search by party name", func(t *testing.T) {
		req := PartySearchRequest{
			PartyName: "ACME Corp",
			Limit:     10,
		}

		result, err := repo.SearchParties(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestFinancialAggregation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://medisync:medisync@localhost:5432/medisync_test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	repo := NewPGIndexRepository(pool)

	t.Run("aggregate contract values", func(t *testing.T) {
		req := FinancialSearchRequest{
			MetricType: "contract_value",
			Aggregation: &Aggregation{
				Operation: "SUM",
				Field:     "amount_normalized",
			},
		}

		result, err := repo.SearchFinancials(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/agents/module_rag/... -v -run TestClauseSearch
```

**Expected:** `undefined: NewPGIndexRepository`

**Step 3: Implement the repository**

```go
// internal/agents/module_rag/pg_repository.go

package module_rag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGIndexRepository implements IndexRepository using PostgreSQL + pgvector
type PGIndexRepository struct {
	pool *pgxpool.Pool
}

// NewPGIndexRepository creates a new PostgreSQL index repository
func NewPGIndexRepository(pool *pgxpool.Pool) *PGIndexRepository {
	return &PGIndexRepository{pool: pool}
}

// SearchClauses searches the clauses index
func (r *PGIndexRepository) SearchClauses(ctx context.Context, req ClauseSearchRequest) (*IndexResult, error) {
	startTime := time.Now()
	result := &IndexResult{
		SubQueryID: "sq_result",
		Chunks:     []ChunkMatch{},
		ByContract: make(map[uuid.UUID][]ChunkMatch),
	}

	// Build query based on search type
	var query string
	var args []interface{}
	argPos := 1

	if len(req.QueryVector) > 0 && req.UseHybrid {
		// Hybrid search: combine vector similarity with full-text search
		query = fmt.Sprintf(`
			SELECT
				clause_id, document_id, clause_type, clause_title, clause_text,
				page_number, position_in_doc, parent_section, metadata,
				-- Combine vector similarity (70%) and FTS rank (30%)
				(0.7 * (1 - (embedding <=> $%d)) + 0.3 * ts_rank(search_vector, websearch_to_tsquery($%d))) as score
			FROM contract_clauses
			WHERE 1=1
		`, argPos, argPos+1)
		args = append(args, pgvector(req.QueryVector), req.QueryText)
		argPos += 2
	} else if len(req.QueryVector) > 0 {
		// Pure semantic search
		query = fmt.Sprintf(`
			SELECT
				clause_id, document_id, clause_type, clause_title, clause_text,
				page_number, position_in_doc, parent_section, metadata,
				(1 - (embedding <=> $%d)) as score
			FROM contract_clauses
			WHERE 1=1
		`, argPos)
		args = append(args, pgvector(req.QueryVector))
		argPos++
	} else {
		// Filter-only search
		query = `
			SELECT
				clause_id, document_id, clause_type, clause_title, clause_text,
				page_number, position_in_doc, parent_section, metadata,
				1.0 as score
			FROM contract_clauses
			WHERE 1=1
		`
	}

	// Add filters
	query, args = r.addFilters(query, args, req.Filters, &argPos)

	// Add ordering and limit
	if len(req.QueryVector) > 0 {
		query += " ORDER BY score DESC"
	}
	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	// Execute query
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("clause search failed: %w", err)
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var chunk ChunkMatch
		var metadataBytes []byte

		err := rows.Scan(
			&chunk.ChunkID,
			&chunk.DocumentID,
			&chunk.IndexSource, // clause_type stored here temporarily
			&chunk.MatchReason, // clause_title stored here temporarily
			&chunk.Text,
			&chunk.PageNumber,
			&chunk.Score,
			&chunk.PageNumber,
			&metadataBytes,
			&chunk.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning clause row: %w", err)
		}

		chunk.IndexSource = "clauses"
		result.Chunks = append(result.Chunks, chunk)
		result.ByContract[chunk.DocumentID] = append(result.ByContract[chunk.DocumentID], chunk)
	}

	result.LatencyMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// SearchParties searches the parties index
func (r *PGIndexRepository) SearchParties(ctx context.Context, req PartySearchRequest) (*IndexResult, error) {
	startTime := time.Now()
	result := &IndexResult{
		SubQueryID: "sq_result",
		Chunks:     []ChunkMatch{},
		ByContract: make(map[uuid.UUID][]ChunkMatch),
	}

	query := `
		SELECT party_id, document_id, party_name, party_role, jurisdiction, identifiers
		FROM contract_parties
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	// Add party name filter (case-insensitive, partial match)
	if req.PartyName != "" {
		query += fmt.Sprintf(" AND party_name_normalized ILIKE $%d", argPos)
		args = append(args, "%"+strings.ToLower(req.PartyName)+"%")
		argPos++
	}

	// Add party role filter
	if req.PartyRole != "" {
		query += fmt.Sprintf(" AND party_role = $%d", argPos)
		args = append(args, req.PartyRole)
		argPos++
	}

	// Add document ID filter
	if len(req.DocumentIDs) > 0 {
		query += fmt.Sprintf(" AND document_id = ANY($%d)", argPos)
		args = append(args, req.DocumentIDs)
		argPos++
	}

	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("party search failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var partyID, documentID uuid.UUID
		var partyName, partyRole, jurisdiction string
		var identifiers map[string]string

		err := rows.Scan(&partyID, &documentID, &partyName, &partyRole, &jurisdiction, &identifiers)
		if err != nil {
			return nil, fmt.Errorf("scanning party row: %w", err)
		}

		chunk := ChunkMatch{
			ChunkID:     partyID,
			DocumentID:  documentID,
			IndexSource: "parties",
			Text:        fmt.Sprintf("%s (%s)", partyName, partyRole),
			Score:       1.0,
		}
		result.Chunks = append(result.Chunks, chunk)
		result.ByContract[documentID] = append(result.ByContract[documentID], chunk)
	}

	result.LatencyMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// SearchFinancials searches the financial index
func (r *PGIndexRepository) SearchFinancials(ctx context.Context, req FinancialSearchRequest) (*IndexResult, error) {
	startTime := time.Now()
	result := &IndexResult{
		SubQueryID: "sq_result",
		Chunks:     []ChunkMatch{},
		ByContract: make(map[uuid.UUID][]ChunkMatch),
	}

	// Handle aggregation queries differently
	if req.Aggregation != nil {
		return r.executeFinancialAggregation(ctx, req, startTime)
	}

	query := `
		SELECT financial_id, document_id, metric_type, amount, currency,
			   amount_normalized, effective_date, expiry_date, payment_terms
		FROM contract_financials
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if req.MetricType != "" {
		query += fmt.Sprintf(" AND metric_type = $%d", argPos)
		args = append(args, req.MetricType)
		argPos++
	}

	if req.MinAmount > 0 {
		query += fmt.Sprintf(" AND amount_normalized >= $%d", argPos)
		args = append(args, req.MinAmount)
		argPos++
	}

	if req.MaxAmount > 0 {
		query += fmt.Sprintf(" AND amount_normalized <= $%d", argPos)
		args = append(args, req.MaxAmount)
		argPos++
	}

	if req.EffectiveAfter != nil {
		query += fmt.Sprintf(" AND effective_date >= $%d", argPos)
		args = append(args, req.EffectiveAfter)
		argPos++
	}

	if req.ExpiresBefore != nil {
		query += fmt.Sprintf(" AND expiry_date <= $%d", argPos)
		args = append(args, req.ExpiresBefore)
		argPos++
	}

	if len(req.DocumentIDs) > 0 {
		query += fmt.Sprintf(" AND document_id = ANY($%d)", argPos)
		args = append(args, req.DocumentIDs)
		argPos++
	}

	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("financial search failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var financialID, documentID uuid.UUID
		var metricType, currency, paymentTerms string
		var amount, amountNormalized float64
		var effectiveDate, expiryDate *time.Time

		err := rows.Scan(&financialID, &documentID, &metricType, &amount, &currency,
			&amountNormalized, &effectiveDate, &expiryDate, &paymentTerms)
		if err != nil {
			return nil, fmt.Errorf("scanning financial row: %w", err)
		}

		chunk := ChunkMatch{
			ChunkID:     financialID,
			DocumentID:  documentID,
			IndexSource: "financial",
			Text:        fmt.Sprintf("%s: %.2f %s", metricType, amount, currency),
			Score:       1.0,
		}
		result.Chunks = append(result.Chunks, chunk)
		result.ByContract[documentID] = append(result.ByContract[documentID], chunk)
	}

	result.LatencyMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// executeFinancialAggregation handles aggregation queries
func (r *PGIndexRepository) executeFinancialAggregation(ctx context.Context, req FinancialSearchRequest, startTime time.Time) (*IndexResult, error) {
	result := &IndexResult{
		SubQueryID: "sq_result",
		Chunks:     []ChunkMatch{},
		ByContract: make(map[uuid.UUID][]ChunkMatch),
	}

	// Build aggregation query
	var selectClause string
	switch strings.ToUpper(req.Aggregation.Operation) {
	case "SUM":
		selectClause = fmt.Sprintf("SUM(%s) as aggregate_value", req.Aggregation.Field)
	case "AVG":
		selectClause = fmt.Sprintf("AVG(%s) as aggregate_value", req.Aggregation.Field)
	case "COUNT":
		selectClause = "COUNT(*) as aggregate_value"
	case "MIN":
		selectClause = fmt.Sprintf("MIN(%s) as aggregate_value", req.Aggregation.Field)
	case "MAX":
		selectClause = fmt.Sprintf("MAX(%s) as aggregate_value", req.Aggregation.Field)
	default:
		return nil, fmt.Errorf("unsupported aggregation operation: %s", req.Aggregation.Operation)
	}

	query := fmt.Sprintf("SELECT %s FROM contract_financials WHERE 1=1", selectClause)
	args := []interface{}{}
	argPos := 1

	// Add filters (same as non-aggregation query)
	if req.MetricType != "" {
		query += fmt.Sprintf(" AND metric_type = $%d", argPos)
		args = append(args, req.MetricType)
		argPos++
	}

	if len(req.DocumentIDs) > 0 {
		query += fmt.Sprintf(" AND document_id = ANY($%d)", argPos)
		args = append(args, req.DocumentIDs)
		argPos++
	}

	var aggregateValue float64
	err := r.pool.QueryRow(ctx, query, args...).Scan(&aggregateValue)
	if err != nil {
		return nil, fmt.Errorf("financial aggregation failed: %w", err)
	}

	// Return aggregate as a single "chunk"
	chunk := ChunkMatch{
		ChunkID:     uuid.New(),
		DocumentID:  uuid.Nil, // No specific document for aggregates
		IndexSource: "financial_aggregate",
		Text:        fmt.Sprintf("%s: %.2f", req.Aggregation.Operation, aggregateValue),
		Score:       1.0,
	}
	result.Chunks = append(result.Chunks, chunk)

	result.LatencyMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// SearchSummaries searches the summaries index
func (r *PGIndexRepository) SearchSummaries(ctx context.Context, req SummarySearchRequest) (*IndexResult, error) {
	startTime := time.Now()
	result := &IndexResult{
		SubQueryID: "sq_result",
		Chunks:     []ChunkMatch{},
		ByContract: make(map[uuid.UUID][]ChunkMatch),
	}

	var query string
	var args []interface{}
	argPos := 1

	if len(req.QueryVector) > 0 {
		query = fmt.Sprintf(`
			SELECT summary_id, document_id, summary_type, summary_text, key_entities, key_risks,
				   (1 - (embedding <=> $%d)) as score
			FROM contract_summaries
			WHERE 1=1
		`, argPos)
		args = append(args, pgvector(req.QueryVector))
		argPos++
	} else {
		query = `
			SELECT summary_id, document_id, summary_type, summary_text, key_entities, key_risks, 1.0 as score
			FROM contract_summaries
			WHERE 1=1
		`
	}

	if req.SummaryType != "" {
		query += fmt.Sprintf(" AND summary_type = $%d", argPos)
		args = append(args, req.SummaryType)
		argPos++
	}

	if len(req.DocumentIDs) > 0 {
		query += fmt.Sprintf(" AND document_id = ANY($%d)", argPos)
		args = append(args, req.DocumentIDs)
		argPos++
	}

	if len(req.QueryVector) > 0 {
		query += " ORDER BY score DESC"
	}

	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("summary search failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk ChunkMatch
		var keyEntities, keyRisks []byte

		err := rows.Scan(
			&chunk.ChunkID,
			&chunk.DocumentID,
			&chunk.MatchReason, // summary_type temporarily
			&chunk.Text,
			&keyEntities,
			&keyRisks,
			&chunk.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning summary row: %w", err)
		}

		chunk.IndexSource = "summaries"
		result.Chunks = append(result.Chunks, chunk)
		result.ByContract[chunk.DocumentID] = append(result.ByContract[chunk.DocumentID], chunk)
	}

	result.LatencyMs = time.Since(startTime).Milliseconds()
	return result, nil
}

// GetContractMeta retrieves metadata for a contract
func (r *PGIndexRepository) GetContractMeta(ctx context.Context, documentID uuid.UUID) (*ContractMeta, error) {
	// Get parties
	parties, err := r.getContractParties(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Get financial summary
	financial, err := r.getContractFinancialSummary(ctx, documentID)
	if err != nil {
		return nil, err
	}

	return &ContractMeta{
		Parties:    parties,
		TotalValue: financial,
	}, nil
}

func (r *PGIndexRepository) getContractParties(ctx context.Context, documentID uuid.UUID) ([]string, error) {
	query := `SELECT party_name FROM contract_parties WHERE document_id = $1`
	rows, err := r.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parties []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		parties = append(parties, name)
	}
	return parties, nil
}

func (r *PGIndexRepository) getContractFinancialSummary(ctx context.Context, documentID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount_normalized), 0)
		FROM contract_financials
		WHERE document_id = $1 AND metric_type = 'contract_value'
	`
	var total float64
	err := r.pool.QueryRow(ctx, query, documentID).Scan(&total)
	return total, err
}

// InsertClause inserts a clause into the index
func (r *PGIndexRepository) InsertClause(ctx context.Context, clause *ClauseRecord) error {
	metadataJSON, _ := json.Marshal(clause.Metadata)

	query := `
		INSERT INTO contract_clauses (document_id, clause_type, clause_title, clause_text,
			embedding, page_number, position_in_doc, parent_section, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		clause.DocumentID, clause.ClauseType, clause.ClauseTitle, clause.ClauseText,
		pgvector(clause.Embedding), clause.PageNumber, clause.PositionInDoc,
		clause.ParentSection, metadataJSON)
	return err
}

// InsertParty inserts a party into the index
func (r *PGIndexRepository) InsertParty(ctx context.Context, party *PartyRecord) error {
	identifiersJSON, _ := json.Marshal(party.Identifiers)

	query := `
		INSERT INTO contract_parties (document_id, party_name, party_name_normalized,
			party_role, jurisdiction, embedding, entity_type, identifiers)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		party.DocumentID, party.PartyName, party.PartyNameNormalized,
		party.PartyRole, party.Jurisdiction, pgvector(party.Embedding),
		party.EntityType, identifiersJSON)
	return err
}

// InsertFinancial inserts a financial record into the index
func (r *PGIndexRepository) InsertFinancial(ctx context.Context, financial *FinancialRecord) error {
	query := `
		INSERT INTO contract_financials (document_id, metric_type, amount, currency,
			amount_normalized, effective_date, expiry_date, payment_terms, embedding)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		financial.DocumentID, financial.MetricType, financial.Amount, financial.Currency,
		financial.AmountNormalized, financial.EffectiveDate, financial.ExpiryDate,
		financial.PaymentTerms, pgvector(financial.Embedding))
	return err
}

// InsertSummary inserts a summary into the index
func (r *PGIndexRepository) InsertSummary(ctx context.Context, summary *SummaryRecord) error {
	entitiesJSON, _ := json.Marshal(summary.KeyEntities)
	risksJSON, _ := json.Marshal(summary.KeyRisks)

	query := `
		INSERT INTO contract_summaries (document_id, summary_type, summary_text,
			key_entities, key_risks, embedding)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (document_id, summary_type) DO UPDATE SET
			summary_text = EXCLUDED.summary_text,
			key_entities = EXCLUDED.key_entities,
			key_risks = EXCLUDED.key_risks,
			embedding = EXCLUDED.embedding,
			updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, query,
		summary.DocumentID, summary.SummaryType, summary.SummaryText,
		entitiesJSON, risksJSON, pgvector(summary.Embedding))
	return err
}

// addFilters adds filter conditions to a query
func (r *PGIndexRepository) addFilters(query string, args []interface{}, filters []Filter, argPos *int) (string, []interface{}) {
	for _, f := range filters {
		switch f.Operator {
		case "eq":
			query += fmt.Sprintf(" AND %s = $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "neq":
			query += fmt.Sprintf(" AND %s != $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "gt":
			query += fmt.Sprintf(" AND %s > $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "gte":
			query += fmt.Sprintf(" AND %s >= $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "lt":
			query += fmt.Sprintf(" AND %s < $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "lte":
			query += fmt.Sprintf(" AND %s <= $%d", f.Field, *argPos)
			args = append(args, f.Value)
		case "in":
			query += fmt.Sprintf(" AND %s = ANY($%d)", f.Field, *argPos)
			args = append(args, f.Value)
		case "contains":
			query += fmt.Sprintf(" AND %s ILIKE $%d", f.Field, *argPos)
			args = append(args, "%"+fmt.Sprintf("%v", f.Value)+"%")
		}
		*argPos++
	}
	return query, args
}

// pgvector is a helper type for pgvector compatibility
type pgvector []float32

func (v pgvector) Scan(src interface{}) error {
	// Implementation for scanning from database
	return nil
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/agents/module_rag/... -v -short
```

**Step 5: Commit**

```bash
git add internal/agents/module_rag/pg_repository.go internal/agents/module_rag/pg_repository_test.go
git commit -m "feat(rag): implement PostgreSQL index repository with hybrid search"
```

---

## Phase 4: Query Planner Agent

### Task 4.1: Implement Query Planner

**Files:**
- Create: `internal/agents/module_rag/query_planner.go`
- Create: `internal/agents/module_rag/query_planner_test.go`

**Step 1: Write the failing test**

```go
// internal/agents/module_rag/query_planner_test.go

package module_rag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryPlanner_DecomposeComplianceQuery(t *testing.T) {
	planner := NewQueryPlanner(nil) // nil for mock LLM

	ctx := context.Background()

	t.Run("decomposes multi-condition compliance query", func(t *testing.T) {
		query := "Find all ACME contracts with termination clauses and value > $100K"

		plan, err := planner.Plan(ctx, query)
		require.NoError(t, err)

		assert.Equal(t, QueryTypeCompliance, plan.QueryType)
		assert.Len(t, plan.SubQueries, 4) // parties, clauses, financial, merge

		// Verify first sub-query targets parties index
		assert.Equal(t, IndexTypeParties, plan.SubQueries[0].TargetIndex)
		assert.Contains(t, plan.SubQueries[0].QueryText, "ACME")

		// Verify second sub-query targets clauses index
		assert.Equal(t, IndexTypeClauses, plan.SubQueries[1].TargetIndex)

		// Verify third sub-query targets financial index
		assert.Equal(t, IndexTypeFinancial, plan.SubQueries[2].TargetIndex)
	})

	t.Run("decomposes aggregate query", func(t *testing.T) {
		query := "Total value of contracts expiring in Q4 2024"

		plan, err := planner.Plan(ctx, query)
		require.NoError(t, err)

		assert.Equal(t, QueryTypeAggregate, plan.QueryType)
		assert.Len(t, plan.SubQueries, 1)
		assert.NotNil(t, plan.SubQueries[0].Aggregation)
		assert.Equal(t, "SUM", plan.SubQueries[0].Aggregation.Operation)
	})

	t.Run("decomposes lookup query", func(t *testing.T) {
		query := "What is the termination clause in contract #123?"

		plan, err := planner.Plan(ctx, query)
		require.NoError(t, err)

		assert.Equal(t, QueryTypeLookup, plan.QueryType)
	})
}

func TestQueryPlanner_BuildDependencyGraph(t *testing.T) {
	planner := NewQueryPlanner(nil)

	subQueries := []SubQuery{
		{ID: "sq_1", Dependencies: []string{}},
		{ID: "sq_2", Dependencies: []string{}},
		{ID: "sq_3", Dependencies: []string{"sq_1", "sq_2"}},
		{ID: "sq_4", Dependencies: []string{"sq_3"}},
	}

	levels := planner.BuildDependencyLevels(subQueries)

	require.Len(t, levels, 3)
	assert.Len(t, levels[0], 2) // sq_1, sq_2 (parallel)
	assert.Len(t, levels[1], 1) // sq_3
	assert.Len(t, levels[2], 1) // sq_4
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/agents/module_rag/... -v -run TestQueryPlanner
```

**Step 3: Implement the query planner**

```go
// internal/agents/module_rag/query_planner.go

package module_rag

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// QueryPlannerAgent decomposes complex queries into sub-queries
type QueryPlannerAgent struct {
	llmClient LLMClient
}

// LLMClient interface for LLM operations
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// NewQueryPlanner creates a new query planner
func NewQueryPlanner(llmClient LLMClient) *QueryPlannerAgent {
	return &QueryPlannerAgent{
		llmClient: llmClient,
	}
}

// Plan decomposes a query into a structured execution plan
func (p *QueryPlannerAgent) Plan(ctx context.Context, query string) (*QueryPlan, error) {
	queryID := uuid.New().String()

	// Classify query type
	queryType := p.classifyQuery(query)

	// Decompose based on query type
	var subQueries []SubQuery
	var err error

	switch queryType {
	case QueryTypeCompliance:
		subQueries, err = p.decomposeCompliance(query)
	case QueryTypeAggregate:
		subQueries, err = p.decomposeAggregate(query)
	case QueryTypeComparison:
		subQueries, err = p.decomposeComparison(query)
	case QueryTypeLookup:
		subQueries, err = p.decomposeLookup(query)
	}

	if err != nil {
		return nil, fmt.Errorf("query decomposition failed: %w", err)
	}

	// Determine execution plan
	executionPlan := p.determineExecutionPlan(subQueries)

	return &QueryPlan{
		QueryID:       queryID,
		OriginalQuery: query,
		QueryType:     queryType,
		SubQueries:    subQueries,
		ExecutionPlan: executionPlan,
		Reasoning:     p.generateReasoning(queryType, subQueries),
	}, nil
}

// classifyQuery determines the query type
func (p *QueryPlannerAgent) classifyQuery(query string) QueryType {
	queryLower := strings.ToLower(query)

	// Check for aggregate patterns
	aggregatePatterns := []string{"total", "sum", "average", "avg", "count", "maximum", "minimum", "how many"}
	for _, pattern := range aggregatePatterns {
		if strings.Contains(queryLower, pattern) {
			return QueryTypeAggregate
		}
	}

	// Check for comparison patterns
	comparisonPatterns := []string{"compare", "difference", "versus", "vs", "between"}
	for _, pattern := range comparisonPatterns {
		if strings.Contains(queryLower, pattern) {
			return QueryTypeComparison
		}
	}

	// Check for compliance patterns
	compliancePatterns := []string{"which contracts", "which documents", "missing", "lack", "without", "find all"}
	for _, pattern := range compliancePatterns {
		if strings.Contains(queryLower, pattern) {
			return QueryTypeCompliance
		}
	}

	// Default to lookup
	return QueryTypeLookup
}

// decomposeCompliance handles compliance queries
func (p *QueryPlannerAgent) decomposeCompliance(query string) ([]SubQuery, error) {
	var subQueries []SubQuery

	// Extract party names
	parties := p.extractEntities(query, "party")
	if len(parties) > 0 {
		subQueries = append(subQueries, SubQuery{
			ID:            "sq_parties",
			QueryText:     fmt.Sprintf("contracts with %s", strings.Join(parties, ", ")),
			TargetIndex:   IndexTypeParties,
			Filters:       []Filter{{Field: "party_name", Operator: "contains", Value: parties[0]}},
			RetrievalType: RetrievalTypeExact,
			Limit:         100,
		})
	}

	// Extract clause types
	clauseTypes := p.extractClauseTypes(query)
	if len(clauseTypes) > 0 {
		subQueries = append(subQueries, SubQuery{
			ID:            "sq_clauses",
			QueryText:     fmt.Sprintf("%s clauses", strings.Join(clauseTypes, ", ")),
			TargetIndex:   IndexTypeClauses,
			Filters:       []Filter{{Field: "clause_type", Operator: "eq", Value: clauseTypes[0]}},
			RetrievalType: RetrievalTypeSemantic,
			Limit:         100,
		})
	}

	// Extract financial conditions
	amountFilter := p.extractAmountFilter(query)
	if amountFilter != nil {
		subQueries = append(subQueries, SubQuery{
			ID:            "sq_financial",
			QueryText:     "contract values",
			TargetIndex:   IndexTypeFinancial,
			Filters:       []Filter{*amountFilter},
			RetrievalType: RetrievalTypeExact,
			Limit:         100,
		})
	}

	// Add final merge query
	dependencies := make([]string, len(subQueries))
	for i, sq := range subQueries {
		dependencies[i] = sq.ID
	}

	subQueries = append(subQueries, SubQuery{
		ID:           "sq_merge",
		QueryText:    "merge and retrieve full context",
		TargetIndex:  IndexTypeSummaries,
		Dependencies: dependencies,
		Limit:        50,
	})

	return subQueries, nil
}

// decomposeAggregate handles aggregate queries
func (p *QueryPlannerAgent) decomposeAggregate(query string) ([]SubQuery, error) {
	// Extract date range
	dateRange := p.extractDateRange(query)

	// Extract metric type
	metricType := "contract_value"
	if strings.Contains(strings.ToLower(query), "value") || strings.Contains(strings.ToLower(query), "amount") {
		metricType = "contract_value"
	}

	// Determine aggregation operation
	operation := "SUM"
	queryLower := strings.ToLower(query)
	if strings.Contains(queryLower, "average") || strings.Contains(queryLower, "avg") {
		operation = "AVG"
	} else if strings.Contains(queryLower, "count") || strings.Contains(queryLower, "how many") {
		operation = "COUNT"
	} else if strings.Contains(queryLower, "maximum") || strings.Contains(queryLower, "max") {
		operation = "MAX"
	} else if strings.Contains(queryLower, "minimum") || strings.Contains(queryLower, "min") {
		operation = "MIN"
	}

	filters := []Filter{}
	if dateRange != nil {
		filters = append(filters, *dateRange)
	}

	return []SubQuery{
		{
			ID:      "sq_aggregate",
			QueryText: query,
			TargetIndex: IndexTypeFinancial,
			Filters: filters,
			Aggregation: &Aggregation{
				Operation: operation,
				Field:     "amount_normalized",
			},
			Limit: 1,
		},
	}, nil
}

// decomposeComparison handles comparison queries
func (p *QueryPlannerAgent) decomposeComparison(query string) ([]SubQuery, error) {
	// Extract entities to compare
	parties := p.extractEntities(query, "party")

	var subQueries []SubQuery

	// Query for all parties mentioned
	if len(parties) > 0 {
		subQueries = append(subQueries, SubQuery{
			ID:            "sq_parties",
			QueryText:     fmt.Sprintf("contracts for %s", strings.Join(parties, ", ")),
			TargetIndex:   IndexTypeParties,
			Filters:       []Filter{{Field: "party_name", Operator: "in", Value: parties}},
			RetrievalType: RetrievalTypeExact,
			Limit:         100,
		})
	}

	// Query for the comparison dimension
	clauseTypes := p.extractClauseTypes(query)
	if len(clauseTypes) > 0 {
		subQueries = append(subQueries, SubQuery{
			ID:            "sq_clauses",
			QueryText:     fmt.Sprintf("%s clauses", clauseTypes[0]),
			TargetIndex:   IndexTypeClauses,
			Filters:       []Filter{{Field: "clause_type", Operator: "eq", Value: clauseTypes[0]}},
			RetrievalType: RetrievalTypeSemantic,
			Dependencies:  []string{"sq_parties"},
			Limit:         100,
		})
	}

	return subQueries, nil
}

// decomposeLookup handles simple lookup queries
func (p *QueryPlannerAgent) decomposeLookup(query string) ([]SubQuery, error) {
	return []SubQuery{
		{
			ID:            "sq_lookup",
			QueryText:     query,
			TargetIndex:   IndexTypeClauses,
			RetrievalType: RetrievalTypeSemantic,
			Limit:         10,
		},
	}, nil
}

// extractEntities extracts named entities from query
func (p *QueryPlannerAgent) extractEntities(query, entityType string) []string {
	// Simple pattern matching for demo
	// In production, use NER model
	patterns := map[string][]string{
		"party": {"ACME", "Corp", "Inc", "LLC", "Ltd", "Company"},
	}

	var entities []string
	words := strings.Fields(query)
	for _, word := range words {
		for _, suffix := range patterns[entityType] {
			if strings.Contains(word, suffix) {
				entities = append(entities, word)
			}
		}
	}
	return entities
}

// extractClauseTypes extracts clause types from query
func (p *QueryPlannerAgent) extractClauseTypes(query string) []string {
	clauseTypes := []string{
		"termination", "liability", "payment", "confidentiality",
		"intellectual property", "ip", "force majeure", "indemnification",
	}

	queryLower := strings.ToLower(query)
	var found []string
	for _, ct := range clauseTypes {
		if strings.Contains(queryLower, ct) {
			found = append(found, ct)
		}
	}
	return found
}

// extractAmountFilter extracts amount filters from query
func (p *QueryPlannerAgent) extractAmountFilter(query string) *Filter {
	// Pattern: "> $100K", ">100000", "exceeds 100k"
	re := regexp.MustCompile(`(?:>|greater than|exceeds?|above)\s*\$?(\d+(?:,\d+)*(?:\.\d+)?)[KkMmBb]?`)
	matches := re.FindStringSubmatch(query)

	if len(matches) > 1 {
		amountStr := strings.ReplaceAll(matches[1], ",", "")
		var amount float64
		fmt.Sscanf(amountStr, "%f", &amount)

		// Handle K/M/B suffixes
		queryLower := strings.ToLower(query)
		if idx := strings.Index(queryLower, matches[1]); idx > 0 {
			suffixIdx := idx + len(matches[1])
			if suffixIdx < len(queryLower) {
				suffix := queryLower[suffixIdx]
				if suffix == 'k' {
					amount *= 1000
				} else if suffix == 'm' {
					amount *= 1000000
				} else if suffix == 'b' {
					amount *= 1000000000
				}
			}
		}

		return &Filter{
			Field:    "amount_normalized",
			Operator: "gt",
			Value:    amount,
		}
	}
	return nil
}

// extractDateRange extracts date range from query
func (p *QueryPlannerAgent) extractDateRange(query string) *Filter {
	// Pattern: "Q4 2024", "2024", "January 2024", etc.
	// Simplified for demo
	if strings.Contains(query, "Q4 2024") || strings.Contains(query, "q4 2024") {
		return &Filter{
			Field:    "expiry_date",
			Operator: "gte",
			Value:    "2024-10-01",
		}
	}
	return nil
}

// determineExecutionPlan determines how to execute the sub-queries
func (p *QueryPlannerAgent) determineExecutionPlan(subQueries []SubQuery) string {
	hasDependencies := false
	for _, sq := range subQueries {
		if len(sq.Dependencies) > 0 {
			hasDependencies = true
			break
		}
	}

	if hasDependencies {
		return "parallel_then_merge"
	}
	return "parallel"
}

// generateReasoning explains the decomposition
func (p *QueryPlannerAgent) generateReasoning(queryType QueryType, subQueries []SubQuery) string {
	return fmt.Sprintf("Query classified as %s. Decomposed into %d sub-queries across %d indices.",
		queryType.String(), len(subQueries), p.countUniqueIndices(subQueries))
}

func (p *QueryPlannerAgent) countUniqueIndices(subQueries []SubQuery) int {
	indices := make(map[IndexType]bool)
	for _, sq := range subQueries {
		indices[sq.TargetIndex] = true
	}
	return len(indices)
}

// BuildDependencyLevels organizes sub-queries by dependency level for parallel execution
func (p *QueryPlannerAgent) BuildDependencyLevels(subQueries []SubQuery) [][]SubQuery {
	// Build dependency graph
	dependencyMap := make(map[string][]string)
	queryMap := make(map[string]SubQuery)

	for _, sq := range subQueries {
		dependencyMap[sq.ID] = sq.Dependencies
		queryMap[sq.ID] = sq
	}

	// Topological sort with level tracking
	levels := [][]SubQuery{}
	completed := make(map[string]bool)

	for len(completed) < len(subQueries) {
		level := []SubQuery{}

		for _, sq := range subQueries {
			if completed[sq.ID] {
				continue
			}

			// Check if all dependencies are completed
			allDepsCompleted := true
			for _, dep := range dependencyMap[sq.ID] {
				if !completed[dep] {
					allDepsCompleted = false
					break
				}
			}

			if allDepsCompleted {
				level = append(level, sq)
			}
		}

		if len(level) == 0 {
			break // Cycle detected or error
		}

		levels = append(levels, level)
		for _, sq := range level {
			completed[sq.ID] = true
		}
	}

	return levels
}

// PlanWithContext creates a plan considering previous hop results
func (p *QueryPlannerAgent) PlanWithContext(ctx context.Context, query string, previousHops []*HopResult) (*QueryPlan, error) {
	// For now, delegate to basic Plan
	// In production, consider previous results for refinement
	return p.Plan(ctx, query)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/agents/module_rag/... -v -run TestQueryPlanner
```

**Step 5: Commit**

```bash
git add internal/agents/module_rag/query_planner.go internal/agents/module_rag/query_planner_test.go
git commit -m "feat(rag): implement query planner agent with decomposition logic"
```

---

## Phase 5: Parallel Retrieval Executor

### Task 5.1: Implement Retrieval Executor

**Files:**
- Create: `internal/agents/module_rag/retrieval_executor.go`
- Create: `internal/agents/module_rag/retrieval_executor_test.go`

**Step 1: Write the failing test**

```go
// internal/agents/module_rag/retrieval_executor_test.go

package module_rag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrievalExecutor_ExecutePlan(t *testing.T) {
	// Mock repository for testing
	executor := NewRetrievalExecutor(&mockIndexRepository{}, nil)

	ctx := context.Background()

	t.Run("executes parallel sub-queries", func(t *testing.T) {
		plan := &QueryPlan{
			QueryID: "test-123",
			SubQueries: []SubQuery{
				{ID: "sq_1", TargetIndex: IndexTypeParties, Limit: 10},
				{ID: "sq_2", TargetIndex: IndexTypeClauses, Limit: 10},
			},
		}

		result, err := executor.ExecutePlan(ctx, plan)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("respects dependencies", func(t *testing.T) {
		plan := &QueryPlan{
			QueryID: "test-456",
			SubQueries: []SubQuery{
				{ID: "sq_1", TargetIndex: IndexTypeParties, Limit: 10},
				{ID: "sq_2", TargetIndex: IndexTypeClauses, Dependencies: []string{"sq_1"}, Limit: 10},
			},
		}

		result, err := executor.ExecutePlan(ctx, plan)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

// mockIndexRepository implements IndexRepository for testing
type mockIndexRepository struct{}

func (m *mockIndexRepository) SearchClauses(ctx context.Context, req ClauseSearchRequest) (*IndexResult, error) {
	return &IndexResult{SubQueryID: "test", Chunks: []ChunkMatch{}, ByContract: make(map[uuid.UUID][]ChunkMatch)}, nil
}

func (m *mockIndexRepository) SearchParties(ctx context.Context, req PartySearchRequest) (*IndexResult, error) {
	return &IndexResult{SubQueryID: "test", Chunks: []ChunkMatch{}, ByContract: make(map[uuid.UUID][]ChunkMatch)}, nil
}

func (m *mockIndexRepository) SearchFinancials(ctx context.Context, req FinancialSearchRequest) (*IndexResult, error) {
	return &IndexResult{SubQueryID: "test", Chunks: []ChunkMatch{}, ByContract: make(map[uuid.UUID][]ChunkMatch)}, nil
}

func (m *mockIndexRepository) SearchSummaries(ctx context.Context, req SummarySearchRequest) (*IndexResult, error) {
	return &IndexResult{SubQueryID: "test", Chunks: []ChunkMatch{}, ByContract: make(map[uuid.UUID][]ChunkMatch)}, nil
}

func (m *mockIndexRepository) GetContractMeta(ctx context.Context, documentID uuid.UUID) (*ContractMeta, error) {
	return &ContractMeta{}, nil
}

func (m *mockIndexRepository) InsertClause(ctx context.Context, clause *ClauseRecord) error { return nil }
func (m *mockIndexRepository) InsertParty(ctx context.Context, party *PartyRecord) error     { return nil }
func (m *mockIndexRepository) InsertFinancial(ctx context.Context, financial *FinancialRecord) error {
	return nil
}
func (m *mockIndexRepository) InsertSummary(ctx context.Context, summary *SummaryRecord) error { return nil }
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/agents/module_rag/... -v -run TestRetrievalExecutor
```

**Step 3: Implement the executor**

```go
// internal/agents/module_rag/retrieval_executor.go

package module_rag

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RetrievalExecutor executes query plans with parallel retrieval
type RetrievalExecutor struct {
	repo       IndexRepository
	embedder   Embedder
	maxWorkers int
	timeout    time.Duration
}

// Embedder interface for generating embeddings
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// NewRetrievalExecutor creates a new retrieval executor
func NewRetrievalExecutor(repo IndexRepository, embedder Embedder) *RetrievalExecutor {
	return &RetrievalExecutor{
		repo:       repo,
		embedder:   embedder,
		maxWorkers: 4,
		timeout:    5 * time.Second,
	}
}

// ExecutePlan executes all sub-queries in a plan
func (e *RetrievalExecutor) ExecutePlan(ctx context.Context, plan *QueryPlan) (map[string]*IndexResult, error) {
	planner := NewQueryPlanner(nil)
	levels := planner.BuildDependencyLevels(plan.SubQueries)

	results := make(map[string]*IndexResult)
	mu := sync.Mutex{}

	for levelNum, level := range levels {
		// Execute all queries at this level in parallel
		var wg sync.WaitGroup
		errChan := make(chan error, len(level))

		for _, sq := range level {
			wg.Add(1)
			go func(subQuery SubQuery) {
				defer wg.Done()

				// Apply timeout
				queryCtx, cancel := context.WithTimeout(ctx, e.timeout)
				defer cancel()

				// Get document IDs from dependencies
				var docIDs []uuid.UUID
				for _, depID := range subQuery.Dependencies {
					mu.Lock()
					if depResult, ok := results[depID]; ok {
						for docID := range depResult.ByContract {
							docIDs = append(docIDs, docID)
						}
					}
					mu.Unlock()
				}

				// Execute the sub-query
				result, err := e.executeSubQuery(queryCtx, subQuery, docIDs)
				if err != nil {
					errChan <- fmt.Errorf("sub-query %s failed: %w", subQuery.ID, err)
					return
				}

				mu.Lock()
				results[subQuery.ID] = result
				mu.Unlock()
			}(sq)
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		for err := range errChan {
			if err != nil {
				return nil, err
			}
		}

		// Log level completion
		fmt.Printf("Completed level %d with %d queries\n", levelNum, len(level))
	}

	return results, nil
}

// executeSubQuery executes a single sub-query against the appropriate index
func (e *RetrievalExecutor) executeSubQuery(ctx context.Context, sq SubQuery, docIDs []uuid.UUID) (*IndexResult, error) {
	switch sq.TargetIndex {
	case IndexTypeClauses:
		return e.executeClauseQuery(ctx, sq, docIDs)
	case IndexTypeParties:
		return e.executePartyQuery(ctx, sq, docIDs)
	case IndexTypeFinancial:
		return e.executeFinancialQuery(ctx, sq, docIDs)
	case IndexTypeSummaries:
		return e.executeSummaryQuery(ctx, sq, docIDs)
	default:
		return nil, fmt.Errorf("unknown index type: %d", sq.TargetIndex)
	}
}

func (e *RetrievalExecutor) executeClauseQuery(ctx context.Context, sq SubQuery, docIDs []uuid.UUID) (*IndexResult, error) {
	req := ClauseSearchRequest{
		Filters:   sq.Filters,
		Limit:     sq.Limit,
		UseHybrid: sq.RetrievalType == RetrievalTypeHybrid,
	}

	// Generate embedding for semantic search
	if sq.RetrievalType == RetrievalTypeSemantic || sq.RetrievalType == RetrievalTypeHybrid {
		if e.embedder != nil {
			embedding, err := e.embedder.Embed(ctx, sq.QueryText)
			if err != nil {
				return nil, fmt.Errorf("embedding failed: %w", err)
			}
			req.QueryText = sq.QueryText
			req.QueryVector = embedding
		}
	}

	return e.repo.SearchClauses(ctx, req)
}

func (e *RetrievalExecutor) executePartyQuery(ctx context.Context, sq SubQuery, docIDs []uuid.UUID) (*IndexResult, error) {
	req := PartySearchRequest{
		DocumentIDs: docIDs,
		Limit:       sq.Limit,
	}

	// Extract party name from filters
	for _, f := range sq.Filters {
		if f.Field == "party_name" {
			if str, ok := f.Value.(string); ok {
				req.PartyName = str
			}
		}
	}

	return e.repo.SearchParties(ctx, req)
}

func (e *RetrievalExecutor) executeFinancialQuery(ctx context.Context, sq SubQuery, docIDs []uuid.UUID) (*IndexResult, error) {
	req := FinancialSearchRequest{
		DocumentIDs: docIDs,
		Limit:       sq.Limit,
		Aggregation: sq.Aggregation,
	}

	// Extract filters
	for _, f := range sq.Filters {
		switch f.Field {
		case "metric_type":
			if str, ok := f.Value.(string); ok {
				req.MetricType = str
			}
		case "amount_normalized":
			switch f.Operator {
			case "gt":
				if num, ok := f.Value.(float64); ok {
					req.MinAmount = num
				}
			case "lt":
				if num, ok := f.Value.(float64); ok {
					req.MaxAmount = num
				}
			}
		case "expiry_date":
			if str, ok := f.Value.(string); ok {
				t, err := time.Parse("2006-01-02", str)
				if err == nil {
					if f.Operator == "gte" {
						req.EffectiveAfter = &t
					} else if f.Operator == "lte" {
						req.ExpiresBefore = &t
					}
				}
			}
		}
	}

	return e.repo.SearchFinancials(ctx, req)
}

func (e *RetrievalExecutor) executeSummaryQuery(ctx context.Context, sq SubQuery, docIDs []uuid.UUID) (*IndexResult, error) {
	req := SummarySearchRequest{
		DocumentIDs: docIDs,
		Limit:       sq.Limit,
	}

	return e.repo.SearchSummaries(ctx, req)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/agents/module_rag/... -v -run TestRetrievalExecutor
```

**Step 5: Commit**

```bash
git add internal/agents/module_rag/retrieval_executor.go internal/agents/module_rag/retrieval_executor_test.go
git commit -m "feat(rag): implement parallel retrieval executor with dependency handling"
```

---

## Phase 6: Evidence Aggregator

### Task 6.1: Implement Evidence Aggregator

**Files:**
- Create: `internal/agents/module_rag/evidence_aggregator.go`
- Create: `internal/agents/module_rag/evidence_aggregator_test.go`

**Step 1: Write the failing test**

```go
// internal/agents/module_rag/evidence_aggregator_test.go

package module_rag

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvidenceAggregator_Aggregate(t *testing.T) {
	aggregator := NewEvidenceAggregator(0.95)

	t.Run("aggregates results from multiple indices", func(t *testing.T) {
		docID := uuid.New()

		results := map[string]*IndexResult{
			"sq_1": {
				SubQueryID: "sq_1",
				Chunks: []ChunkMatch{
					{ChunkID: uuid.New(), DocumentID: docID, Text: "Party: ACME Corp", Score: 0.95},
				},
				ByContract: map[uuid.UUID][]ChunkMatch{
					docID: {{ChunkID: uuid.New(), DocumentID: docID, Text: "Party: ACME Corp", Score: 0.95}},
				},
			},
			"sq_2": {
				SubQueryID: "sq_2",
				Chunks: []ChunkMatch{
					{ChunkID: uuid.New(), DocumentID: docID, Text: "Termination clause...", Score: 0.88},
				},
				ByContract: map[uuid.UUID][]ChunkMatch{
					docID: {{ChunkID: uuid.New(), DocumentID: docID, Text: "Termination clause...", Score: 0.88}},
				},
			},
		}

		plan := &QueryPlan{
			QueryID:       "test-123",
			OriginalQuery: "Find ACME contracts with termination clauses",
			QueryType:     QueryTypeCompliance,
		}

		evidence, err := aggregator.Aggregate(results, plan)
		require.NoError(t, err)

		assert.Equal(t, "test-123", evidence.QueryID)
		assert.Len(t, evidence.ContractMatches, 1)
		assert.Equal(t, docID, evidence.ContractMatches[0].DocumentID)
		assert.GreaterOrEqual(t, evidence.Confidence, 0.5)
	})

	t.Run("deduplicates similar chunks", func(t *testing.T) {
		docID := uuid.New()
		chunkID := uuid.New()

		results := map[string]*IndexResult{
			"sq_1": {
				ByContract: map[uuid.UUID][]ChunkMatch{
					docID: {
						{ChunkID: chunkID, DocumentID: docID, Text: "Same text", Score: 0.9},
					},
				},
			},
			"sq_2": {
				ByContract: map[uuid.UUID][]ChunkMatch{
					docID: {
						{ChunkID: chunkID, DocumentID: docID, Text: "Same text", Score: 0.9},
					},
				},
			},
		}

		plan := &QueryPlan{QueryID: "test-dedup"}

		evidence, err := aggregator.Aggregate(results, plan)
		require.NoError(t, err)

		// Should deduplicate identical chunks
		assert.LessOrEqual(t, len(evidence.ContractMatches[0].RelevantChunks), 1)
	})
}
```

**Step 2: Implement the aggregator**

```go
// internal/agents/module_rag/evidence_aggregator.go

package module_rag

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

// EvidenceAggregator merges and scores evidence from multiple retrieval results
type EvidenceAggregator struct {
	deduplicationThreshold float64
}

// NewEvidenceAggregator creates a new evidence aggregator
func NewEvidenceAggregator(deduplicationThreshold float64) *EvidenceAggregator {
	return &EvidenceAggregator{
		deduplicationThreshold: deduplicationThreshold,
	}
}

// Aggregate merges results from multiple sub-queries into unified evidence
func (a *EvidenceAggregator) Aggregate(
	results map[string]*IndexResult,
	plan *QueryPlan,
) (*AggregatedEvidence, error) {
	evidence := &AggregatedEvidence{
		QueryID:         plan.QueryID,
		ContractMatches: []ContractMatch{},
		EvidenceGaps:    []string{},
		SourceCitations: []Citation{},
	}

	// 1. Collect all contract IDs
	contractIDs := a.collectContractIDs(results)

	// 2. Merge evidence for each contract
	for contractID := range contractIDs {
		match := ContractMatch{
			DocumentID:     contractID,
			RelevantChunks: []ChunkMatch{},
		}

		// Collect chunks from all sub-query results
		for sqID, result := range results {
			if chunks, ok := result.ByContract[contractID]; ok {
				match.RelevantChunks = append(match.RelevantChunks, chunks...)
				_ = sqID // Used for tracking source
			}
		}

		// 3. Deduplicate chunks
		match.RelevantChunks = a.deduplicateChunks(match.RelevantChunks)

		// 4. Compute aggregate match score
		match.MatchScore = a.computeMatchScore(match.RelevantChunks, plan)

		evidence.ContractMatches = append(evidence.ContractMatches, match)
	}

	// 5. Sort by relevance score
	sort.Slice(evidence.ContractMatches, func(i, j int) bool {
		return evidence.ContractMatches[i].MatchScore > evidence.ContractMatches[j].MatchScore
	})

	// 6. Identify evidence gaps
	evidence.EvidenceGaps = a.identifyGaps(results, plan)

	// 7. Compute overall confidence
	evidence.Confidence = a.computeOverallConfidence(evidence, plan)

	// 8. Build citations
	evidence.SourceCitations = a.buildCitations(evidence.ContractMatches)

	evidence.TotalContracts = len(evidence.ContractMatches)

	return evidence, nil
}

// collectContractIDs gathers all unique contract IDs from results
func (a *EvidenceAggregator) collectContractIDs(results map[string]*IndexResult) map[uuid.UUID]bool {
	contractIDs := make(map[uuid.UUID]bool)

	for _, result := range results {
		for docID := range result.ByContract {
			contractIDs[docID] = true
		}
	}

	return contractIDs
}

// deduplicateChunks removes duplicate or near-duplicate chunks
func (a *EvidenceAggregator) deduplicateChunks(chunks []ChunkMatch) []ChunkMatch {
	if len(chunks) <= 1 {
		return chunks
	}

	seen := make(map[uuid.UUID]bool)
	deduped := []ChunkMatch{}

	for _, chunk := range chunks {
		if !seen[chunk.ChunkID] {
			seen[chunk.ChunkID] = true
			deduped = append(deduped, chunk)
		}
	}

	return deduped
}

// computeMatchScore calculates relevance score for a contract match
func (a *EvidenceAggregator) computeMatchScore(chunks []ChunkMatch, plan *QueryPlan) float64 {
	if len(chunks) == 0 {
		return 0
	}

	// Weighted average of chunk scores
	totalScore := 0.0
	for _, chunk := range chunks {
		totalScore += chunk.Score
	}

	return totalScore / float64(len(chunks))
}

// identifyGaps finds missing information
func (a *EvidenceAggregator) identifyGaps(results map[string]*IndexResult, plan *QueryPlan) []string {
	gaps := []string{}

	// Check if any sub-queries returned no results
	for _, sq := range plan.SubQueries {
		if result, ok := results[sq.ID]; ok {
			if len(result.Chunks) == 0 {
				gaps = append(gaps, fmt.Sprintf("No results from %s index", sq.TargetIndex))
			}
		}
	}

	return gaps
}

// computeOverallConfidence calculates confidence score for the entire result
func (a *EvidenceAggregator) computeOverallConfidence(evidence *AggregatedEvidence, plan *QueryPlan) float64 {
	if len(evidence.ContractMatches) == 0 {
		return 0
	}

	// Base confidence on number of contracts found and their scores
	avgScore := 0.0
	for _, match := range evidence.ContractMatches {
		avgScore += match.MatchScore
	}
	avgScore /= float64(len(evidence.ContractMatches))

	// Penalize for gaps
	gapPenalty := float64(len(evidence.EvidenceGaps)) * 0.1

	confidence := avgScore - gapPenalty
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

// buildCitations creates citation objects from contract matches
func (a *EvidenceAggregator) buildCitations(matches []ContractMatch) []Citation {
	citations := []Citation{}

	for _, match := range matches {
		for i, chunk := range match.RelevantChunks {
			if i >= 10 { // Limit citations
				break
			}

			citation := Citation{
				CitationID:    fmt.Sprintf("cite_%s_%d", match.DocumentID.String()[:8], i),
				DocumentID:    match.DocumentID,
				ChunkID:       chunk.ChunkID,
				PageNumber:    chunk.PageNumber,
				Text:          truncateText(chunk.Text, 200),
			}
			citations = append(citations, citation)
		}
	}

	return citations
}

// truncateText truncates text to maxLen characters
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// AggregateMultiHop merges evidence from multiple hops
func (a *EvidenceAggregator) AggregateMultiHop(hopResults []*HopResult) *AggregatedEvidence {
	// Merge all results from all hops
	allResults := make(map[string]*IndexResult)

	for _, hop := range hopResults {
		for sqID, result := range hop.Results {
			// Prefix with hop number to avoid collisions
			key := fmt.Sprintf("hop%d_%s", hop.HopNumber, sqID)
			allResults[key] = result
		}
	}

	// Create a combined plan for context
	combinedPlan := &QueryPlan{
		QueryID:       "multi_hop",
		OriginalQuery: "Multi-hop query",
	}

	evidence, _ := a.Aggregate(allResults, combinedPlan)
	return evidence
}
```

**Step 3: Run tests**

```bash
go test ./internal/agents/module_rag/... -v -run TestEvidenceAggregator
```

**Step 4: Commit**

```bash
git add internal/agents/module_rag/evidence_aggregator.go internal/agents/module_rag/evidence_aggregator_test.go
git commit -m "feat(rag): implement evidence aggregator with deduplication and scoring"
```

---

## Phase 7-10: Remaining Components (Summary)

The remaining phases follow the same TDD pattern:

### Phase 7: Synthesis Agent
- `internal/agents/module_rag/synthesis_agent.go` - LLM-based synthesis with citations
- Response formatting with mandatory source references
- Citation validation to prevent hallucination

### Phase 8: Multi-Hop Controller
- `internal/agents/module_rag/multihop_controller.go` - Iterative retrieval orchestration
- Confidence-based early termination
- Follow-up query generation

### Phase 9: API Layer
- `internal/api/handlers/rag_handler.go` - HTTP endpoints
- `policies/rag.rego` - OPA authorization policies

### Phase 10: Testing & Evaluation
- `tests/rag/test_cases/` - Test dataset
- `internal/agents/module_rag/evaluation/evaluator.go` - Quality metrics

---

## Verification

After completing all phases, verify the implementation:

```bash
# Run all tests
go test ./internal/agents/module_rag/... -v

# Run migrations
go run ./cmd/migrate up

# Start the API
go run ./cmd/api

# Test the RAG endpoint
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Find all ACME contracts with termination clauses"}'
```

---

## Success Metrics

| Metric | Target | How to Verify |
|--------|--------|---------------|
| Retrieval accuracy | ≥ 85% | Run evaluation suite |
| Latency P95 | ≤ 5s | Load testing |
| Citation accuracy | ≥ 95% | Manual spot check |
| Test coverage | ≥ 80% | `go test -cover` |

---

*End of Implementation Plan*
