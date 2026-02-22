// Package council provides the Council Coordinator for orchestrating deliberations.
//
// The coordinator manages the complete deliberation flow:
// 1. Receive query and create deliberation record
// 2. Dispatch query to all healthy agents
// 3. Collect responses with timeout handling
// 4. Calculate consensus using semantic equivalence
// 5. Retrieve evidence from Knowledge Graph
// 6. Store results and create audit trail
//
// Key Features:
//   - Concurrent agent querying
//   - Graceful degradation on failures
//   - Evidence caching
//   - Structured logging
package council

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Coordinator orchestrates Council deliberations.
type Coordinator struct {
	agents             []Agent
	agentWrappers      []*AgentWrapper
	kgRepo             KnowledgeGraphRepository
	repo               Repository
	evidenceCache      EvidenceCache
	healthMonitor      *HealthMonitor
	consensusCalculator *ConsensusCalculator
	evidenceRetriever  *EvidenceRetriever
	config             CoordinatorConfig
	logger             *slog.Logger
}

// CoordinatorConfig holds configuration for the coordinator.
type CoordinatorConfig struct {
	Agents             []Agent
	KGRepository       KnowledgeGraphRepository
	Repository         Repository
	EvidenceCache      EvidenceCache
	Timeout            time.Duration
	ConsensusThreshold float64
	MinAgents          int
	MaxHops            int
	Logger             *slog.Logger
}

// NewCoordinator creates a new Council coordinator.
func NewCoordinator(config CoordinatorConfig) *Coordinator {
	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultAgentTimeoutSecs * time.Second
	}
	if config.ConsensusThreshold == 0 {
		config.ConsensusThreshold = DefaultConsensusThreshold
	}
	if config.MinAgents == 0 {
		config.MinAgents = DefaultMinAgents
	}
	if config.MaxHops == 0 {
		config.MaxHops = DefaultMaxHops
	}

	// Create agent wrappers
	wrappers := make([]*AgentWrapper, len(config.Agents))
	for i, agent := range config.Agents {
		wrappers[i] = NewAgentWrapper(agent, config.Timeout)
	}

	// Create health monitor
	monitor := NewHealthMonitor(config.Logger)
	for i, agent := range config.Agents {
		instance := &AgentInstance{
			ID:       agent.GetID(),
			Name:     agent.GetName(),
			HealthStatus: HealthHealthy,
		}
		monitor.RegisterAgent(instance)
		_ = wrappers[i] // Use wrapper for health tracking
	}

	return &Coordinator{
		agents:             config.Agents,
		agentWrappers:      wrappers,
		kgRepo:             config.KGRepository,
		repo:               config.Repository,
		evidenceCache:      config.EvidenceCache,
		healthMonitor:      monitor,
		consensusCalculator: NewConsensusCalculator(config.ConsensusThreshold),
		evidenceRetriever:  NewEvidenceRetriever(config.KGRepository, config.MaxHops),
		config:             config,
		logger:             config.Logger,
	}
}

// Deliberate processes a query and returns the deliberation result.
func (c *Coordinator) Deliberate(ctx context.Context, req CreateDeliberationRequest, userID string) (*DeliberationResult, error) {
	startTime := time.Now()

	// 1. Create deliberation record
	threshold := req.ConsensusThreshold
	if threshold == 0 {
		threshold = c.config.ConsensusThreshold
	}

	deliberation, err := c.repo.CreateDeliberation(ctx, req.Query, userID, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to create deliberation: %w", err)
	}

	c.logger.Info("deliberation started",
		slog.String("deliberation_id", deliberation.ID),
		slog.String("user_id", userID),
		slog.String("query_hash", deliberation.QueryHash),
	)

	// Update status to deliberating
	if err := c.repo.UpdateDeliberationStatus(ctx, deliberation.ID, StatusDeliberating, "", 0); err != nil {
		c.logger.Error("failed to update status", slog.Any("error", err))
	}

	// 2. Get healthy agents
	healthyAgents := c.getHealthyAgents()
	if len(healthyAgents) < c.config.MinAgents {
		return nil, fmt.Errorf("minimum %d agents required, only %d healthy", c.config.MinAgents, len(healthyAgents))
	}

	// 3. Query agents concurrently
	responses := c.queryAgents(ctx, healthyAgents, req.Query)

	if len(responses) < c.config.MinAgents {
		return nil, fmt.Errorf("insufficient agent responses: %d < %d", len(responses), c.config.MinAgents)
	}

	// 4. Store agent responses
	for _, resp := range responses {
		resp.DeliberationID = deliberation.ID
		if err := c.repo.CreateAgentResponse(ctx, resp); err != nil {
			c.logger.Error("failed to store agent response",
				slog.String("agent_id", resp.AgentID),
				slog.Any("error", err),
			)
		}
	}

	// 5. Calculate consensus
	consensusResult, err := c.consensusCalculator.Calculate(responses)
	if err != nil {
		return nil, fmt.Errorf("consensus calculation failed: %w", err)
	}

	// 6. Retrieve evidence
	evidenceTrail, err := c.retrieveEvidence(ctx, deliberation.ID, req.Query)
	if err != nil {
		c.logger.Warn("evidence retrieval failed, continuing without evidence",
			slog.Any("error", err),
		)
	}

	// 7. Store consensus record
	consensusRecord := &ConsensusRecord{
		ID:                generateID(),
		DeliberationID:    deliberation.ID,
		AgreementScore:    consensusResult.AgreementScore,
		EquivalenceGroups: consensusResult.EquivalenceGroups,
		ThresholdMet:      consensusResult.ThresholdMet,
		DissentingAgents:  consensusResult.DissentingAgents,
		ConsensusMethod:   consensusResult.Method,
	}
	if err := c.repo.CreateConsensusRecord(ctx, consensusRecord); err != nil {
		c.logger.Error("failed to store consensus record", slog.Any("error", err))
	}

	// 8. Store evidence trail
	if evidenceTrail != nil {
		evidenceTrail.DeliberationID = deliberation.ID
		if err := c.repo.CreateEvidenceTrail(ctx, evidenceTrail); err != nil {
			c.logger.Error("failed to store evidence trail", slog.Any("error", err))
		}
	}

	// 9. Update deliberation status
	finalStatus := consensusResult.Status
	if err := c.repo.UpdateDeliberationStatus(ctx, deliberation.ID, finalStatus,
		consensusResult.FinalResponse, consensusResult.ConfidenceScore); err != nil {
		c.logger.Error("failed to update deliberation status", slog.Any("error", err))
	}

	// 10. Create audit entry
	auditEntry := &AuditEntry{
		ID:             generateID(),
		DeliberationID: deliberation.ID,
		UserID:         userID,
		Action:         AuditActionQuery,
		Details: map[string]any{
			"query_length":       len(req.Query),
			"agent_count":        len(responses),
			"consensus_score":    consensusResult.AgreementScore,
			"evidence_node_count": len(evidenceTrail.NodeIDs),
			"duration_ms":        time.Since(startTime).Milliseconds(),
		},
	}
	if err := c.repo.CreateAuditEntry(ctx, auditEntry); err != nil {
		c.logger.Error("failed to create audit entry", slog.Any("error", err))
	}

	c.logger.Info("deliberation completed",
		slog.String("deliberation_id", deliberation.ID),
		slog.String("status", string(finalStatus)),
		slog.Float64("consensus_score", consensusResult.AgreementScore),
		slog.Duration("duration", time.Since(startTime)),
	)

	return &DeliberationResult{
		Deliberation:    deliberation,
		ConsensusRecord: consensusRecord,
		EvidenceTrail:   evidenceTrail,
		AgentResponses:  responses,
	}, nil
}

// getHealthyAgents returns agents that can accept requests.
func (c *Coordinator) getHealthyAgents() []*AgentWrapper {
	var healthy []*AgentWrapper
	for _, wrapper := range c.agentWrappers {
		if wrapper.CanAcceptRequest() {
			healthy = append(healthy, wrapper)
		}
	}
	return healthy
}

// queryAgents queries all agents concurrently and collects responses.
func (c *Coordinator) queryAgents(ctx context.Context, wrappers []*AgentWrapper, query string) []*AgentResponse {
	var wg sync.WaitGroup
	var mu sync.Mutex
	responses := make([]*AgentResponse, 0)

	for _, wrapper := range wrappers {
		wg.Add(1)
		go func(w *AgentWrapper) {
			defer wg.Done()

			startTime := time.Now()
			resp, err := w.Query(ctx, query)
			duration := time.Since(startTime)

			if err != nil {
				c.logger.Warn("agent query failed",
					slog.String("agent_id", w.agent.GetID()),
					slog.Any("error", err),
					slog.Duration("duration", duration),
				)
				c.healthMonitor.RecordFailure(w.agent.GetID(), err.Error())
				return
			}

			c.healthMonitor.RecordSuccess(w.agent.GetID(), duration)

			mu.Lock()
			responses = append(responses, resp)
			mu.Unlock()
		}(wrapper)
	}

	wg.Wait()
	return responses
}

// retrieveEvidence retrieves evidence from the Knowledge Graph.
func (c *Coordinator) retrieveEvidence(ctx context.Context, deliberationID, query string) (*EvidenceTrail, error) {
	queryHash := HashQuery(query)

	// Try cache first
	if c.evidenceCache != nil {
		cached, err := c.evidenceCache.Get(ctx, queryHash)
		if err == nil && cached != nil {
			c.logger.Debug("using cached evidence",
				slog.String("query_hash", queryHash),
				slog.Int("node_count", len(cached.NodeIDs)),
			)
			return cached.ToEvidenceTrail(deliberationID), nil
		}
	}

	// Check KG health
	if err := c.kgRepo.HealthCheck(ctx); err != nil {
		c.logger.Warn("knowledge graph health check failed",
			slog.Any("error", err),
		)
		// Try to use stale cache if available
		if c.evidenceCache != nil {
			cached, _ := c.evidenceCache.Get(ctx, queryHash)
			if cached != nil && cached.IsTrustworthyDuringOutage() {
				return cached.ToEvidenceTrail(deliberationID), nil
			}
		}
		return nil, fmt.Errorf("knowledge graph unavailable: %w", err)
	}

	// Retrieve fresh evidence
	// For now, we use empty initial nodes (would normally be derived from query embedding)
	trail, err := c.evidenceRetriever.Retrieve(ctx, []string{}, c.config.MaxHops)
	if err != nil {
		return nil, fmt.Errorf("evidence retrieval failed: %w", err)
	}

	// Cache the evidence
	if c.evidenceCache != nil {
		entry := &EvidenceCacheEntry{
			QueryHash:       queryHash,
			QueryText:       query,
			NodeIDs:         trail.NodeIDs,
			TraversalPath:   trail.TraversalPath,
			RelevanceScores: trail.RelevanceScore,
			HopCount:        trail.HopCount,
			KGHealthStatus:  "healthy",
		}
		if err := c.evidenceCache.Set(ctx, queryHash, entry); err != nil {
			c.logger.Warn("failed to cache evidence", slog.Any("error", err))
		}
	}

	return trail, nil
}

// GetHealthMonitor returns the health monitor for external health checks.
func (c *Coordinator) GetHealthMonitor() *HealthMonitor {
	return c.healthMonitor
}

// GetHealthSummary returns a summary of agent health.
func (c *Coordinator) GetHealthSummary() HealthSummary {
	return c.healthMonitor.GetSummary()
}
