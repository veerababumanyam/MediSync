// Package council implements the Council of AIs consensus system for hallucination eradication.
//
// The Council coordinates multiple independent AI agent instances to reach consensus
// on responses, grounded by Graph-of-Thoughts retrieval from a Medical Knowledge Graph.
//
// Key Concepts:
//   - Deliberation: A single query processing session with multiple agents
//   - Consensus: Agreement threshold (default 80%) required to release responses
//   - Evidence Trail: Knowledge Graph traversal path supporting a response
//   - Semantic Equivalence: 95% similarity threshold for grouping equivalent responses
//
// Usage:
//
//	coordinator := council.NewCoordinator(agents, kgRepo, cache)
//	result, err := coordinator.Deliberate(ctx, query, userID)
package council

import (
	"time"

	"github.com/medisync/medisync/internal/warehouse"
	"github.com/pgvector/pgvector-go"
)

// DeliberationStatus represents the current state of a deliberation.
type DeliberationStatus string

const (
	StatusPending     DeliberationStatus = "pending"     // Query received, not yet processing
	StatusDeliberating DeliberationStatus = "deliberating" // Agents are analyzing
	StatusConsensus   DeliberationStatus = "consensus"   // Consensus reached
	StatusUncertain   DeliberationStatus = "uncertain"   // No consensus, uncertainty signaled
	StatusFailed      DeliberationStatus = "failed"      // Processing failed
)

// AgentHealthStatus represents the health state of an agent instance.
type AgentHealthStatus string

const (
	HealthHealthy  AgentHealthStatus = "healthy"  // Agent responding normally
	HealthDegraded AgentHealthStatus = "degraded" // Agent slow or intermittent
	HealthFailed   AgentHealthStatus = "failed"   // Agent not responding
)

// AuditActionType represents the type of action being audited.
type AuditActionType string

const (
	AuditActionQuery  AuditActionType = "query"
	AuditActionReview AuditActionType = "review"
	AuditActionFlag   AuditActionType = "flag"
	AuditActionExport AuditActionType = "export"
	AuditActionAccess AuditActionType = "access"
)

// CouncilDeliberation represents a single query processing session.
type CouncilDeliberation struct {
	ID                 string             `json:"id"`
	QueryText          string             `json:"query_text"`
	QueryHash          string             `json:"query_hash"`
	UserID             string             `json:"user_id"`
	Status             DeliberationStatus `json:"status"`
	ConsensusThreshold float64            `json:"consensus_threshold"`
	FinalResponse      string             `json:"final_response,omitempty"`
	ConfidenceScore    float64            `json:"confidence_score,omitempty"`
	ErrorMessage       string             `json:"error_message,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	CompletedAt        *time.Time         `json:"completed_at,omitempty"`
}

// AgentInstance represents an independent AI reasoning unit in the Council.
type AgentInstance struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	HealthStatus  AgentHealthStatus `json:"health_status"`
	LastHeartbeat time.Time         `json:"last_heartbeat"`
	Config        map[string]any    `json:"config"`
	TimeoutSecs   int               `json:"timeout_seconds"`
	CreatedAt     time.Time         `json:"created_at"`
}

// AgentResponse represents an individual response from a single agent.
type AgentResponse struct {
	ID             string         `json:"id"`
	DeliberationID string         `json:"deliberation_id"`
	AgentID        string         `json:"agent_id"`
	ResponseText   string         `json:"response_text"`
	EvidenceIDs    []string       `json:"evidence_ids"`
	Confidence     float64        `json:"confidence"`
	Embedding      pgvector.Vector `json:"embedding,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

// ConsensusRecord captures the consensus calculation results for a deliberation.
type ConsensusRecord struct {
	ID                string    `json:"id"`
	DeliberationID    string    `json:"deliberation_id"`
	AgreementScore    float64   `json:"agreement_score"`
	EquivalenceGroups []Group   `json:"equivalence_groups"`
	ThresholdMet      bool      `json:"threshold_met"`
	DissentingAgents  []string  `json:"dissenting_agents"`
	ConsensusMethod   string    `json:"consensus_method"`
	CreatedAt         time.Time `json:"created_at"`
}

// Group represents a cluster of semantically equivalent agent responses.
type Group struct {
	GroupID     int      `json:"group_id"`
	AgentIDs    []string `json:"agent_ids"`
	Canonical   string   `json:"canonical"` // Representative response text
	Similarity  float64  `json:"similarity"`
}

// EvidenceTrail records the Knowledge Graph traversal path for a deliberation.
type EvidenceTrail struct {
	ID             string                     `json:"id"`
	DeliberationID string                     `json:"deliberation_id"`
	NodeIDs        []string                   `json:"node_ids"`
	TraversalPath  []*warehouse.TraversalStep `json:"traversal_path"`
	RelevanceScore map[string]float64         `json:"relevance_scores"` // nodeID → score
	HopCount       int                        `json:"hop_count"`
	CachedAt       time.Time                  `json:"cached_at"`
	ExpiresAt      time.Time                  `json:"expires_at"`
}

// AuditEntry represents an immutable audit log entry for compliance.
type AuditEntry struct {
	ID             string           `json:"id"`
	DeliberationID string           `json:"deliberation_id"`
	UserID         string           `json:"user_id"`
	Action         AuditActionType  `json:"action"`
	Details        map[string]any   `json:"details"`
	IPAddress      string           `json:"ip_address,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	PartitionDate  time.Time        `json:"partition_date"`
}

// CreateDeliberationRequest represents a request to create a new deliberation.
type CreateDeliberationRequest struct {
	Query              string  `json:"query"`
	ConsensusThreshold float64 `json:"consensus_threshold,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

// DeliberationResult represents the outcome of a deliberation.
type DeliberationResult struct {
	Deliberation    *CouncilDeliberation `json:"deliberation"`
	ConsensusRecord *ConsensusRecord     `json:"consensus_record,omitempty"`
	EvidenceTrail   *EvidenceTrail       `json:"evidence_trail,omitempty"`
	AgentResponses  []*AgentResponse     `json:"agent_responses,omitempty"`
}

// FlagDeliberationRequest represents a request to flag a deliberation for review.
type FlagDeliberationRequest struct {
	Reason           string `json:"reason"`
	HallucinationType string `json:"hallucination_type,omitempty"` // factual_error, fabricated_data, misleading_context, other
	Severity         string `json:"severity"` // low, medium, high, critical
}

// Constants for default configuration values
const (
	DefaultMinAgents           = 3
	DefaultConsensusThreshold  = 0.80
	DefaultAgentTimeoutSecs    = 3
	DefaultCacheTTLMinutes     = 5
	DefaultHealthCheckInterval = 30 * time.Second
	DefaultSemanticThreshold   = 0.95
	DefaultMaxHops             = 3
	DefaultAuditRetentionYears = 7
)
