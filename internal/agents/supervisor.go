// Package agents provides AI agent implementations for MediSync.
//
// This file implements the agent supervisor that orchestrates all agents.
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/medisync/medisync/internal/agents/shared"
)

// Supervisor orchestrates all AI agents in the MediSync platform.
type Supervisor struct {
	agents     map[string]Agent
	logger     *slog.Logger
	mu         sync.RWMutex
	health     *HealthStatus
}

// Agent defines the interface for all agents.
type Agent interface {
	// AgentCard returns the agent's discovery metadata
	AgentCard() shared.AgentCard
}

// AgentWithHealth extends Agent with health check capability.
type AgentWithHealth interface {
	Agent
	// HealthCheck performs a health check
	HealthCheck(ctx context.Context) error
}

// SupervisorConfig holds configuration for the supervisor.
type SupervisorConfig struct {
	Logger *slog.Logger
}

// NewSupervisor creates a new agent supervisor.
func NewSupervisor(cfg SupervisorConfig) *Supervisor {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return &Supervisor{
		agents: make(map[string]Agent),
		logger: cfg.Logger.With("component", "supervisor"),
		health: &HealthStatus{
			Agents:    []AgentHealth{},
			UpdatedAt: time.Now(),
		},
	}
}

// Register adds an agent to the supervisor.
func (s *Supervisor) Register(agent Agent) error {
	card := agent.AgentCard()
	id := card.ID

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.agents[id]; exists {
		return fmt.Errorf("agent %s already registered", id)
	}

	s.agents[id] = agent
	s.logger.Info("agent registered", "agent_id", id, "name", card.Name)

	return nil
}

// Unregister removes an agent from the supervisor.
func (s *Supervisor) Unregister(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.agents[agentID]; !exists {
		return fmt.Errorf("agent %s not found", agentID)
	}

	delete(s.agents, agentID)
	s.logger.Info("agent unregistered", "agent_id", agentID)

	return nil
}

// GetAgent retrieves an agent by ID.
func (s *Supervisor) GetAgent(agentID string) (Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agent, exists := s.agents[agentID]
	return agent, exists
}

// GetAllAgents returns all registered agents.
func (s *Supervisor) GetAllAgents() []Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agents := make([]Agent, 0, len(s.agents))
	for _, agent := range s.agents {
		agents = append(agents, agent)
	}
	return agents
}

// GetAllAgentCards returns all agent cards.
func (s *Supervisor) GetAllAgentCards() []shared.AgentCard {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cards := make([]shared.AgentCard, 0, len(s.agents))
	for _, agent := range s.agents {
		cards = append(cards, agent.AgentCard())
	}
	return cards
}

// HealthStatus represents the health status of all agents.
type HealthStatus struct {
	Status      string        `json:"status"`
	Agents      []AgentHealth `json:"agents"`
	LLMProvider *LLMHealth    `json:"llm_provider,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// AgentHealth represents the health of a single agent.
type AgentHealth struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	LastCheck   time.Time `json:"last_check"`
	Error       string    `json:"error,omitempty"`
	Latency     int64     `json:"latency_ms"`
}

// LLMHealth represents the health of the LLM provider.
type LLMHealth struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Status  string `json:"status"`
	Latency int64  `json:"latency_ms,omitempty"`
}

// Health check constants
const (
	StatusHealthy   = "healthy"
	StatusDegraded  = "degraded"
	StatusUnhealthy = "unhealthy"
)

// CheckHealth performs health checks on all agents.
func (s *Supervisor) CheckHealth(ctx context.Context) *HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	health := &HealthStatus{
		Agents:    make([]AgentHealth, 0, len(s.agents)),
		UpdatedAt: time.Now(),
	}

	healthyCount := 0
	degradedCount := 0
	unhealthyCount := 0

	for id, agent := range s.agents {
		card := agent.AgentCard()
		agentHealth := AgentHealth{
			ID:        id,
			Name:      card.Name,
			LastCheck: time.Now(),
		}

		startTime := time.Now()

		// Check if agent supports health checks
		if ha, ok := agent.(AgentWithHealth); ok {
			if err := ha.HealthCheck(ctx); err != nil {
				agentHealth.Status = StatusUnhealthy
				agentHealth.Error = err.Error()
				unhealthyCount++
			} else {
				agentHealth.Status = StatusHealthy
				healthyCount++
			}
		} else {
			// Agent doesn't support health check, assume healthy
			agentHealth.Status = StatusHealthy
			healthyCount++
		}

		agentHealth.Latency = time.Since(startTime).Milliseconds()
		health.Agents = append(health.Agents, agentHealth)
	}

	// Determine overall status
	if unhealthyCount > 0 {
		health.Status = StatusUnhealthy
	} else if degradedCount > 0 {
		health.Status = StatusDegraded
	} else {
		health.Status = StatusHealthy
	}

	s.health = health
	return health
}

// GetCachedHealth returns the last cached health status.
func (s *Supervisor) GetCachedHealth() *HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.health
}

// ExecuteRequest represents a request to execute an agent.
type ExecuteRequest struct {
	AgentID   string          `json:"agent_id"`
	Operation string          `json:"operation"`
	Input     json.RawMessage `json:"input"`
}

// ExecuteResponse represents the response from agent execution.
type ExecuteResponse struct {
	AgentID    string          `json:"agent_id"`
	Operation  string          `json:"operation"`
	Output     json.RawMessage `json:"output"`
	Error      string          `json:"error,omitempty"`
	Duration   int64           `json:"duration_ms"`
	Confidence float64         `json:"confidence,omitempty"`
}

// Execute orchestrates agent execution.
func (s *Supervisor) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	s.mu.RLock()
	_, exists := s.agents[req.AgentID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent %s not found", req.AgentID)
	}

	startTime := time.Now()
	response := &ExecuteResponse{
		AgentID:   req.AgentID,
		Operation: req.Operation,
	}

	// The actual execution depends on the agent type
	// This is a simplified version - real implementation would
	// use reflection or type assertions to call the appropriate method

	response.Duration = time.Since(startTime).Milliseconds()

	s.logger.Info("agent executed",
		"agent_id", req.AgentID,
		"operation", req.Operation,
		"duration_ms", response.Duration)

	return response, nil
}

// GetAgentCount returns the number of registered agents.
func (s *Supervisor) GetAgentCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.agents)
}

// GetAgentsByCapability returns agents that have a specific capability.
func (s *Supervisor) GetAgentsByCapability(capability string) []Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matching []Agent
	for _, agent := range s.agents {
		card := agent.AgentCard()
		for _, cap := range card.Capabilities {
			if cap == capability {
				matching = append(matching, agent)
				break
			}
		}
	}
	return matching
}

// ToJSON serializes the health status.
func (h *HealthStatus) ToJSON() string {
	data, _ := json.Marshal(h)
	return string(data)
}

// SupervisorMetrics contains metrics about the supervisor.
type SupervisorMetrics struct {
	TotalAgents    int            `json:"total_agents"`
	HealthyAgents  int            `json:"healthy_agents"`
	UnhealthyAgents int           `json:"unhealthy_agents"`
	AgentsByModule map[string]int `json:"agents_by_module"`
	LastHealthCheck time.Time     `json:"last_health_check"`
}

// GetMetrics returns supervisor metrics.
func (s *Supervisor) GetMetrics() *SupervisorMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := &SupervisorMetrics{
		TotalAgents:     len(s.agents),
		AgentsByModule:  make(map[string]int),
		LastHealthCheck: s.health.UpdatedAt,
	}

	for _, agent := range s.agents {
		card := agent.AgentCard()
		// Extract module from ID (e.g., "a-01-text-to-sql" -> "module_a")
		module := extractModule(card.ID)
		metrics.AgentsByModule[module]++
	}

	for _, agentHealth := range s.health.Agents {
		if agentHealth.Status == StatusHealthy {
			metrics.HealthyAgents++
		} else {
			metrics.UnhealthyAgents++
		}
	}

	return metrics
}

// extractModule extracts the module from an agent ID.
func extractModule(agentID string) string {
	if len(agentID) > 1 && agentID[1] == '-' {
		switch agentID[0] {
		case 'a':
			return "module_a"
		case 'b':
			return "module_b"
		case 'c':
			return "module_c"
		case 'd':
			return "module_d"
		case 'e':
			return "module_e"
		}
	}
	return "unknown"
}
