// Package council implements the Council of AIs consensus system.
package council

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/medisync/medisync/internal/warehouse"
	"github.com/pgvector/pgvector-go"
)

// Repository defines the interface for Council deliberation data access.
type Repository interface {
	// Deliberation operations
	CreateDeliberation(ctx context.Context, query string, userID string, threshold float64) (*CouncilDeliberation, error)
	GetDeliberation(ctx context.Context, id string) (*CouncilDeliberation, error)
	GetDeliberationWithResponses(ctx context.Context, id string) (*DeliberationResult, error)
	ListDeliberations(ctx context.Context, userID string, isAdmin bool, opts ListOptions) ([]*CouncilDeliberation, int, error)
	UpdateDeliberationStatus(ctx context.Context, id string, status DeliberationStatus, response string, confidence float64) error

	// Agent response operations
	CreateAgentResponse(ctx context.Context, resp *AgentResponse) error
	GetAgentResponses(ctx context.Context, deliberationID string) ([]*AgentResponse, error)

	// Consensus record operations
	CreateConsensusRecord(ctx context.Context, record *ConsensusRecord) error
	GetConsensusRecord(ctx context.Context, deliberationID string) (*ConsensusRecord, error)

	// Evidence trail operations
	CreateEvidenceTrail(ctx context.Context, trail *EvidenceTrail) error
	GetEvidenceTrail(ctx context.Context, deliberationID string) (*EvidenceTrail, error)

	// Agent instance operations
	ListHealthyAgents(ctx context.Context) ([]*AgentInstance, error)
	UpdateAgentHeartbeat(ctx context.Context, agentID string) error
	UpdateAgentHealth(ctx context.Context, agentID string, status AgentHealthStatus) error

	// Audit operations
	CreateAuditEntry(ctx context.Context, entry *AuditEntry) error
	FlagDeliberation(ctx context.Context, deliberationID string, userID string, req FlagDeliberationRequest) error
}

// ListOptions represents options for listing deliberations.
type ListOptions struct {
	Status   DeliberationStatus
	FromDate time.Time
	ToDate   time.Time
	Flagged  *bool
	Limit    int
	Offset   int
}

// CouncilRepo implements Repository using PostgreSQL.
type CouncilRepo struct {
	db   warehouse.DBTX
	kg   warehouse.KnowledgeGraphRepository
}

// NewCouncilRepo creates a new Council repository.
func NewCouncilRepo(db warehouse.DBTX, kgRepo warehouse.KnowledgeGraphRepository) *CouncilRepo {
	return &CouncilRepo{db: db, kg: kgRepo}
}

// hashQuery generates a SHA-256 hash of the query for deduplication.
func hashQuery(query string) string {
	h := sha256.Sum256([]byte(query))
	return hex.EncodeToString(h[:])
}

// CreateDeliberation creates a new deliberation record.
func (r *CouncilRepo) CreateDeliberation(ctx context.Context, query string, userID string, threshold float64) (*CouncilDeliberation, error) {
	d := &CouncilDeliberation{
		QueryText:          query,
		QueryHash:          hashQuery(query),
		UserID:             userID,
		Status:             StatusPending,
		ConsensusThreshold: threshold,
		CreatedAt:          time.Now(),
	}

	querySQL := `
		INSERT INTO council_deliberations (query_text, query_hash, user_id, status, consensus_threshold, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, querySQL,
		d.QueryText, d.QueryHash, d.UserID, d.Status, d.ConsensusThreshold, d.CreatedAt,
	).Scan(&d.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create deliberation: %w", err)
	}

	return d, nil
}

// GetDeliberation retrieves a deliberation by ID.
func (r *CouncilRepo) GetDeliberation(ctx context.Context, id string) (*CouncilDeliberation, error) {
	query := `
		SELECT id, query_text, query_hash, user_id, status, consensus_threshold,
		       final_response, confidence_score, error_message, created_at, completed_at
		FROM council_deliberations
		WHERE id = $1
	`

	d := &CouncilDeliberation{}
	var completedAt *time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.QueryText, &d.QueryHash, &d.UserID, &d.Status, &d.ConsensusThreshold,
		&d.FinalResponse, &d.ConfidenceScore, &d.ErrorMessage, &d.CreatedAt, &completedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliberation: %w", err)
	}
	d.CompletedAt = completedAt

	return d, nil
}

// GetDeliberationWithResponses retrieves a deliberation with all related data.
func (r *CouncilRepo) GetDeliberationWithResponses(ctx context.Context, id string) (*DeliberationResult, error) {
	d, err := r.GetDeliberation(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &DeliberationResult{Deliberation: d}

	// Get agent responses
	responses, err := r.GetAgentResponses(ctx, id)
	if err == nil {
		result.AgentResponses = responses
	}

	// Get consensus record
	record, err := r.GetConsensusRecord(ctx, id)
	if err == nil {
		result.ConsensusRecord = record
	}

	// Get evidence trail
	trail, err := r.GetEvidenceTrail(ctx, id)
	if err == nil {
		result.EvidenceTrail = trail
	}

	return result, nil
}

// ListDeliberations lists deliberations with filtering and pagination.
func (r *CouncilRepo) ListDeliberations(ctx context.Context, userID string, isAdmin bool, opts ListOptions) ([]*CouncilDeliberation, int, error) {
	// Build query based on user role
	whereClause := "WHERE 1=1"
	args := []any{opts.Limit, opts.Offset}
	argIdx := 3

	if !isAdmin {
		whereClause += " AND user_id = $2"
		args = append([]any{userID, opts.Limit, opts.Offset}, args[3:]...)
		argIdx = 4
	}

	if opts.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, opts.Status)
		argIdx++
	}

	if !opts.FromDate.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, opts.FromDate)
		argIdx++
	}

	if !opts.ToDate.IsZero() {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, opts.ToDate)
		argIdx++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM council_deliberations %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args[2:]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deliberations: %w", err)
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, query_text, query_hash, user_id, status, consensus_threshold,
		       final_response, confidence_score, error_message, created_at, completed_at
		FROM council_deliberations
		%s
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deliberations: %w", err)
	}
	defer rows.Close()

	var deliberations []*CouncilDeliberation
	for rows.Next() {
		d := &CouncilDeliberation{}
		var completedAt *time.Time
		err := rows.Scan(
			&d.ID, &d.QueryText, &d.QueryHash, &d.UserID, &d.Status, &d.ConsensusThreshold,
			&d.FinalResponse, &d.ConfidenceScore, &d.ErrorMessage, &d.CreatedAt, &completedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deliberation: %w", err)
		}
		d.CompletedAt = completedAt
		deliberations = append(deliberations, d)
	}

	return deliberations, total, rows.Err()
}

// UpdateDeliberationStatus updates the status and result of a deliberation.
func (r *CouncilRepo) UpdateDeliberationStatus(ctx context.Context, id string, status DeliberationStatus, response string, confidence float64) error {
	query := `
		UPDATE council_deliberations
		SET status = $2, final_response = $3, confidence_score = $4, completed_at = $5
		WHERE id = $1
	`

	var completedAt *time.Time
	if status == StatusConsensus || status == StatusUncertain || status == StatusFailed {
		now := time.Now()
		completedAt = &now
	}

	_, err := r.db.ExecContext(ctx, query, id, status, response, confidence, completedAt)
	if err != nil {
		return fmt.Errorf("failed to update deliberation status: %w", err)
	}

	return nil
}

// CreateAgentResponse creates a new agent response record.
func (r *CouncilRepo) CreateAgentResponse(ctx context.Context, resp *AgentResponse) error {
	query := `
		INSERT INTO agent_responses (id, deliberation_id, agent_id, response_text, evidence_ids, confidence, embedding, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	resp.CreatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		resp.ID, resp.DeliberationID, resp.AgentID, resp.ResponseText,
		resp.EvidenceIDs, resp.Confidence, resp.Embedding, resp.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create agent response: %w", err)
	}

	return nil
}

// GetAgentResponses retrieves all agent responses for a deliberation.
func (r *CouncilRepo) GetAgentResponses(ctx context.Context, deliberationID string) ([]*AgentResponse, error) {
	query := `
		SELECT id, deliberation_id, agent_id, response_text, evidence_ids, confidence, embedding, created_at
		FROM agent_responses
		WHERE deliberation_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, deliberationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent responses: %w", err)
	}
	defer rows.Close()

	var responses []*AgentResponse
	for rows.Next() {
		r := &AgentResponse{}
		err := rows.Scan(
			&r.ID, &r.DeliberationID, &r.AgentID, &r.ResponseText,
			&r.EvidenceIDs, &r.Confidence, &r.Embedding, &r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent response: %w", err)
		}
		responses = append(responses, r)
	}

	return responses, rows.Err()
}

// CreateConsensusRecord creates a new consensus record.
func (r *CouncilRepo) CreateConsensusRecord(ctx context.Context, record *ConsensusRecord) error {
	query := `
		INSERT INTO consensus_records (id, deliberation_id, agreement_score, equivalence_groups, threshold_met, dissenting_agents, consensus_method, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	record.CreatedAt = time.Now()
	groupsJSON, _ := json.Marshal(record.EquivalenceGroups)
	_, err := r.db.ExecContext(ctx, query,
		record.ID, record.DeliberationID, record.AgreementScore,
		groupsJSON, record.ThresholdMet, record.DissentingAgents,
		record.ConsensusMethod, record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create consensus record: %w", err)
	}

	return nil
}

// GetConsensusRecord retrieves the consensus record for a deliberation.
func (r *CouncilRepo) GetConsensusRecord(ctx context.Context, deliberationID string) (*ConsensusRecord, error) {
	query := `
		SELECT id, deliberation_id, agreement_score, equivalence_groups, threshold_met, dissenting_agents, consensus_method, created_at
		FROM consensus_records
		WHERE deliberation_id = $1
	`

	record := &ConsensusRecord{}
	var groupsJSON []byte
	err := r.db.QueryRowContext(ctx, query, deliberationID).Scan(
		&record.ID, &record.DeliberationID, &record.AgreementScore,
		&groupsJSON, &record.ThresholdMet, &record.DissentingAgents,
		&record.ConsensusMethod, &record.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus record: %w", err)
	}

	json.Unmarshal(groupsJSON, &record.EquivalenceGroups)
	return record, nil
}

// CreateEvidenceTrail creates a new evidence trail record.
func (r *CouncilRepo) CreateEvidenceTrail(ctx context.Context, trail *EvidenceTrail) error {
	query := `
		INSERT INTO evidence_trails (id, deliberation_id, node_ids, traversal_path, relevance_scores, hop_count, cached_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	trail.CachedAt = time.Now()
	pathJSON, _ := json.Marshal(trail.TraversalPath)
	scoresJSON, _ := json.Marshal(trail.RelevanceScore)
	_, err := r.db.ExecContext(ctx, query,
		trail.ID, trail.DeliberationID, trail.NodeIDs,
		pathJSON, scoresJSON, trail.HopCount,
		trail.CachedAt, trail.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create evidence trail: %w", err)
	}

	return nil
}

// GetEvidenceTrail retrieves the evidence trail for a deliberation.
func (r *CouncilRepo) GetEvidenceTrail(ctx context.Context, deliberationID string) (*EvidenceTrail, error) {
	query := `
		SELECT id, deliberation_id, node_ids, traversal_path, relevance_scores, hop_count, cached_at, expires_at
		FROM evidence_trails
		WHERE deliberation_id = $1
	`

	trail := &EvidenceTrail{}
	var pathJSON, scoresJSON []byte
	err := r.db.QueryRowContext(ctx, query, deliberationID).Scan(
		&trail.ID, &trail.DeliberationID, &trail.NodeIDs,
		&pathJSON, &scoresJSON, &trail.HopCount,
		&trail.CachedAt, &trail.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get evidence trail: %w", err)
	}

	json.Unmarshal(pathJSON, &trail.TraversalPath)
	json.Unmarshal(scoresJSON, &trail.RelevanceScore)
	return trail, nil
}

// ListHealthyAgents retrieves all healthy agent instances.
func (r *CouncilRepo) ListHealthyAgents(ctx context.Context) ([]*AgentInstance, error) {
	query := `
		SELECT id, name, health_status, last_heartbeat, config, timeout_seconds, created_at
		FROM agent_instances
		WHERE health_status IN ('healthy', 'degraded')
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list healthy agents: %w", err)
	}
	defer rows.Close()

	var agents []*AgentInstance
	for rows.Next() {
		a := &AgentInstance{}
		var configJSON []byte
		err := rows.Scan(
			&a.ID, &a.Name, &a.HealthStatus, &a.LastHeartbeat,
			&configJSON, &a.TimeoutSecs, &a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		json.Unmarshal(configJSON, &a.Config)
		agents = append(agents, a)
	}

	return agents, rows.Err()
}

// UpdateAgentHeartbeat updates the last heartbeat timestamp for an agent.
func (r *CouncilRepo) UpdateAgentHeartbeat(ctx context.Context, agentID string) error {
	query := `UPDATE agent_instances SET last_heartbeat = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, agentID)
	return err
}

// UpdateAgentHealth updates the health status of an agent.
func (r *CouncilRepo) UpdateAgentHealth(ctx context.Context, agentID string, status AgentHealthStatus) error {
	query := `UPDATE agent_instances SET health_status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, agentID, status)
	return err
}

// CreateAuditEntry creates a new audit entry.
func (r *CouncilRepo) CreateAuditEntry(ctx context.Context, entry *AuditEntry) error {
	query := `
		INSERT INTO audit_entries (id, deliberation_id, user_id, action, details, ip_address, created_at, partition_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	entry.CreatedAt = time.Now()
	entry.PartitionDate = entry.CreatedAt
	detailsJSON, _ := json.Marshal(entry.Details)
	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.DeliberationID, entry.UserID, entry.Action,
		detailsJSON, entry.IPAddress, entry.CreatedAt, entry.PartitionDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit entry: %w", err)
	}

	return nil
}

// FlagDeliberation flags a deliberation for review.
func (r *CouncilRepo) FlagDeliberation(ctx context.Context, deliberationID string, userID string, req FlagDeliberationRequest) error {
	// Create audit entry for the flag action
	entry := &AuditEntry{
		DeliberationID: deliberationID,
		UserID:         userID,
		Action:         AuditActionFlag,
		Details: map[string]any{
			"reason":            req.Reason,
			"hallucination_type": req.HallucinationType,
			"severity":          req.Severity,
		},
	}

	return r.CreateAuditEntry(ctx, entry)
}

// Compile-time interface compliance check
var _ Repository = (*CouncilRepo)(nil)

// Helper function to create pgvector from float slice
func NewVector(vals []float32) pgvector.Vector {
	return pgvector.NewVector(vals)
}
