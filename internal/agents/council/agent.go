// Package council provides agent instance management for the Council of AIs.
//
// The agent module implements the agent wrapper pattern with timeout handling,
// circuit breaker, and Genkit integration for LLM-based reasoning.
//
// Key Features:
//   - 3-second timeout for agent responses
//   - Circuit breaker pattern for failure handling
//   - Health status tracking
//   - Genkit flow integration
package council

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Error definitions for agent operations
var (
	ErrAgentTimeout     = errors.New("agent response timeout")
	ErrCircuitOpen      = errors.New("circuit breaker is open")
	ErrAgentUnavailable = errors.New("agent unavailable")
)

// Agent defines the interface for an AI agent instance.
type Agent interface {
	Query(ctx context.Context, query string) (*AgentResponse, error)
	GetID() string
	GetName() string
}

// AgentWrapper wraps an Agent with timeout and circuit breaker functionality.
type AgentWrapper struct {
	agent              Agent
	timeout            time.Duration
	circuitBreaker     *CircuitBreaker
	healthTracker      *AgentHealthTracker
	mu                 sync.RWMutex
}

// NewAgentWrapper creates a new agent wrapper.
func NewAgentWrapper(agent Agent, timeout time.Duration) *AgentWrapper {
	if timeout <= 0 {
		timeout = DefaultAgentTimeoutSecs * time.Second
	}

	instance := &AgentInstance{
		ID:       agent.GetID(),
		Name:     agent.GetName(),
		HealthStatus: HealthHealthy,
	}

	return &AgentWrapper{
		agent:         agent,
		timeout:       timeout,
		circuitBreaker: NewCircuitBreaker(),
		healthTracker: NewAgentHealthTracker(instance),
	}
}

// SetCircuitBreakerThreshold sets the failure threshold for the circuit breaker.
func (w *AgentWrapper) SetCircuitBreakerThreshold(threshold int) {
	w.circuitBreaker.failureThreshold = threshold
}

// Query sends a query to the agent with timeout and circuit breaker.
func (w *AgentWrapper) Query(ctx context.Context, query string) (*AgentResponse, error) {
	// Check circuit breaker
	if !w.circuitBreaker.Allow() {
		return nil, ErrCircuitOpen
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	// Channel for response
	type result struct {
		resp *AgentResponse
		err  error
	}
	resultCh := make(chan result, 1)

	// Execute query in goroutine
	go func() {
		resp, err := w.agent.Query(ctx, query)
		select {
		case resultCh <- result{resp: resp, err: err}:
		case <-ctx.Done():
		}
	}()

	// Wait for response or timeout
	select {
	case r := <-resultCh:
		if r.err != nil {
			w.circuitBreaker.RecordFailure()
			w.healthTracker.RecordFailure(r.err.Error())
			return nil, r.err
		}
		w.circuitBreaker.RecordSuccess()
		w.healthTracker.RecordSuccess(time.Since(time.Now())) // Approximate
		return r.resp, nil
	case <-ctx.Done():
		w.circuitBreaker.RecordFailure()
		w.healthTracker.RecordFailure("timeout")
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrAgentTimeout
		}
		return nil, ctx.Err()
	}
}

// RecordFailure records a failure manually.
func (w *AgentWrapper) RecordFailure(err error) {
	w.circuitBreaker.RecordFailure()
	w.healthTracker.RecordFailure(err.Error())
}

// IsCircuitOpen checks if the circuit breaker is open.
func (w *AgentWrapper) IsCircuitOpen() bool {
	return w.circuitBreaker.State() == CircuitOpen
}

// GetHealthStatus returns the current health status.
func (w *AgentWrapper) GetHealthStatus() AgentHealthStatus {
	return w.healthTracker.GetStatus()
}

// CanAcceptRequest checks if the agent can accept a new request.
func (w *AgentWrapper) CanAcceptRequest() bool {
	return w.healthTracker.CanAcceptRequest()
}

// GetMetrics returns the agent health metrics.
func (w *AgentWrapper) GetMetrics() HealthMetrics {
	return w.healthTracker.GetMetrics()
}

// GenkitAgent implements the Agent interface using Genkit.
type GenkitAgent struct {
	id       string
	name     string
	flowName string
	client   GenkitClient
	config   map[string]any
}

// GenkitClient defines the interface for Genkit operations.
type GenkitClient interface {
	RunFlow(ctx context.Context, flowName string, input map[string]any) (map[string]any, error)
}

// NewGenkitAgent creates a new Genkit-based agent.
func NewGenkitAgent(id, name, flowName string, client GenkitClient, config map[string]any) *GenkitAgent {
	return &GenkitAgent{
		id:       id,
		name:     name,
		flowName: flowName,
		client:   client,
		config:   config,
	}
}

// Query implements the Agent interface.
func (a *GenkitAgent) Query(ctx context.Context, query string) (*AgentResponse, error) {
	input := map[string]any{
		"query": query,
	}
	for k, v := range a.config {
		input[k] = v
	}

	output, err := a.client.RunFlow(ctx, a.flowName, input)
	if err != nil {
		return nil, fmt.Errorf("genkit flow error: %w", err)
	}

	response := &AgentResponse{
		ID:           generateID(),
		AgentID:      a.id,
		ResponseText: getString(output, "response"),
		Confidence:   getFloat64(output, "confidence"),
	}

	if embedding, ok := output["embedding"].([]float32); ok {
		response.Embedding = NewVector(embedding)
	}

	return response, nil
}

// GetID returns the agent ID.
func (a *GenkitAgent) GetID() string {
	return a.id
}

// GetName returns the agent name.
func (a *GenkitAgent) GetName() string {
	return a.name
}

// MockAgent implements Agent for testing purposes.
type MockAgent struct {
	ID            string
	Name          string
	ResponseDelay time.Duration
	ShouldError   bool
	Response      string
	Confidence    float64
}

// Query implements the Agent interface.
func (a *MockAgent) Query(ctx context.Context, query string) (*AgentResponse, error) {
	if a.ResponseDelay > 0 {
		select {
		case <-time.After(a.ResponseDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if a.ShouldError {
		return nil, errors.New("mock agent error")
	}

	response := a.Response
	if response == "" {
		response = "Mock response for: " + query
	}

	confidence := a.Confidence
	if confidence == 0 {
		confidence = 95.0
	}

	return &AgentResponse{
		ID:           generateID(),
		AgentID:      a.ID,
		ResponseText: response,
		Confidence:   confidence,
	}, nil
}

// GetID returns the agent ID.
func (a *MockAgent) GetID() string {
	return a.ID
}

// GetName returns the agent name.
func (a *MockAgent) GetName() string {
	return a.Name
}

// Helper functions

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat64(m map[string]any, key string) float64 {
	switch v := m[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

// Embedding helper for pgvector
type Embedding []float32

// Slice returns the embedding as a float32 slice.
func (e Embedding) Slice() []float32 {
	return []float32(e)
}
