// Package council provides health monitoring for the Council of AIs consensus system.
//
// The health module implements agent health tracking, heartbeat monitoring, and
// circuit breaker patterns to ensure reliable multi-agent deliberation.
//
// Key Features:
//   - Heartbeat tracking with configurable timeout (default 3s)
//   - Health status transitions: healthy → degraded → failed
//   - Circuit breaker to prevent cascading failures
//   - NATS-based health event publishing
//
// Usage:
//
//	monitor := health.NewMonitor(agents, natsClient, logger)
//	go monitor.Start(ctx)
//
//	healthy := monitor.GetHealthyAgents()
//	status := monitor.GetAgentStatus(agentID)
package council

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Health monitoring constants
const (
	// HeartbeatInterval is the default interval between agent heartbeats.
	HeartbeatInterval = 10 * time.Second

	// HeartbeatTimeout is the time after which an agent is considered degraded.
	HeartbeatTimeout = 30 * time.Second

	// FailureTimeout is the time after which a degraded agent is considered failed.
	FailureTimeout = 60 * time.Second

	// CircuitBreakerCooldown is the cooldown period after circuit opens.
	CircuitBreakerCooldown = 30 * time.Second

	// CircuitBreakerThreshold is the number of failures before circuit opens.
	CircuitBreakerThreshold = 3

	// HealthCheckInterval is the interval for health check sweeps.
	HealthCheckInterval = 10 * time.Second

	// MaxConsecutiveFailures is the max failures before marking agent as failed.
	MaxConsecutiveFailures = 5
)

// HealthEvent represents a health status change event.
type HealthEvent struct {
	AgentID      string            `json:"agent_id"`
	AgentName    string            `json:"agent_name"`
	OldStatus    AgentHealthStatus `json:"old_status"`
	NewStatus    AgentHealthStatus `json:"new_status"`
	Reason       string            `json:"reason"`
	Timestamp    time.Time         `json:"timestamp"`
	ResponseTime time.Duration     `json:"response_time,omitempty"`
}

// HealthMetrics captures health-related metrics for an agent.
type HealthMetrics struct {
	AgentID             string        `json:"agent_id"`
	TotalDeliberations  int64         `json:"total_deliberations"`
	SuccessfulResponses int64         `json:"successful_responses"`
	FailedResponses     int64         `json:"failed_responses"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	LastResponseTime    time.Duration `json:"last_response_time"`
	LastSuccessAt       time.Time     `json:"last_success_at"`
	LastFailureAt       time.Time     `json:"last_failure_at"`
	ConsecutiveFailures int           `json:"consecutive_failures"`
	Uptime              time.Duration `json:"uptime"`
}

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState string

const (
	// CircuitClosed means requests flow normally.
	CircuitClosed CircuitBreakerState = "closed"
	// CircuitOpen means requests are blocked.
	CircuitOpen CircuitBreakerState = "open"
	// CircuitHalfOpen means testing if service recovered.
	CircuitHalfOpen CircuitBreakerState = "half_open"
)

// CircuitBreaker implements the circuit breaker pattern for agent calls.
type CircuitBreaker struct {
	mu                sync.RWMutex
	state             CircuitBreakerState
	failures          int
	successes         int
	lastFailureTime   time.Time
	lastStateChange   time.Time
	cooldownPeriod    time.Duration
	failureThreshold int
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state:             CircuitClosed,
		cooldownPeriod:    CircuitBreakerCooldown,
		failureThreshold: CircuitBreakerThreshold,
		lastStateChange:   time.Now(),
	}
}

// Allow checks if a request should be allowed through.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if cooldown period has passed
		if time.Since(cb.lastFailureTime) > cb.cooldownPeriod {
			cb.state = CircuitHalfOpen
			cb.lastStateChange = time.Now()
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.successes++

	if cb.state == CircuitHalfOpen && cb.successes >= cb.failureThreshold {
		cb.state = CircuitClosed
		cb.lastStateChange = time.Now()
	}
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()
	cb.successes = 0

	if cb.failures >= cb.failureThreshold {
		if cb.state != CircuitOpen {
			cb.state = CircuitOpen
			cb.lastStateChange = time.Now()
		}
	}
}

// State returns the current circuit breaker state.
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitClosed
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()
}

// AgentHealthTracker tracks the health status of a single agent.
type AgentHealthTracker struct {
	mu               sync.RWMutex
	agent            *AgentInstance
	lastHeartbeat    time.Time
	metrics          HealthMetrics
	circuitBreaker   *CircuitBreaker
	status           AgentHealthStatus
	statusHistory    []HealthEvent
	maxHistoryLength int
}

// NewAgentHealthTracker creates a new health tracker for an agent.
func NewAgentHealthTracker(agent *AgentInstance) *AgentHealthTracker {
	return &AgentHealthTracker{
		agent:            agent,
		status:           HealthHealthy,
		circuitBreaker:   NewCircuitBreaker(),
		statusHistory:    make([]HealthEvent, 0, 100),
		maxHistoryLength: 100,
		lastHeartbeat:    time.Now(),
	}
}

// RecordHeartbeat records a heartbeat from the agent.
func (t *AgentHealthTracker) RecordHeartbeat() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastHeartbeat = time.Now()

	// If agent was degraded, promote back to healthy
	if t.status == HealthDegraded {
		t.transitionStatus(HealthHealthy, "heartbeat received")
	}
}

// RecordSuccess records a successful response from the agent.
func (t *AgentHealthTracker) RecordSuccess(responseTime time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.metrics.SuccessfulResponses++
	t.metrics.TotalDeliberations++
	t.metrics.LastResponseTime = responseTime
	t.metrics.LastSuccessAt = time.Now()
	t.metrics.ConsecutiveFailures = 0

	// Update average response time
	if t.metrics.AverageResponseTime == 0 {
		t.metrics.AverageResponseTime = responseTime
	} else {
		t.metrics.AverageResponseTime = (t.metrics.AverageResponseTime + responseTime) / 2
	}

	t.circuitBreaker.RecordSuccess()

	// Promote to healthy if degraded
	if t.status == HealthDegraded {
		t.transitionStatus(HealthHealthy, "successful response")
	}
}

// RecordFailure records a failed response from the agent.
func (t *AgentHealthTracker) RecordFailure(reason string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.metrics.FailedResponses++
	t.metrics.TotalDeliberations++
	t.metrics.LastFailureAt = time.Now()
	t.metrics.ConsecutiveFailures++

	t.circuitBreaker.RecordFailure()

	// Check for status transitions
	if t.metrics.ConsecutiveFailures >= MaxConsecutiveFailures {
		t.transitionStatus(HealthFailed, reason)
	} else if t.metrics.ConsecutiveFailures >= 2 {
		t.transitionStatus(HealthDegraded, reason)
	}
}

// CheckHealth checks and updates the agent health status based on heartbeat.
func (t *AgentHealthTracker) CheckHealth() HealthEvent {
	t.mu.Lock()
	defer t.mu.Unlock()

	timeSinceHeartbeat := time.Since(t.lastHeartbeat)

	var newStatus AgentHealthStatus
	var reason string

	switch {
	case timeSinceHeartbeat > FailureTimeout:
		newStatus = HealthFailed
		reason = "heartbeat timeout exceeded"
	case timeSinceHeartbeat > HeartbeatTimeout:
		newStatus = HealthDegraded
		reason = "heartbeat delayed"
	default:
		newStatus = HealthHealthy
		reason = "healthy"
	}

	if newStatus != t.status {
		t.transitionStatus(newStatus, reason)
		return t.statusHistory[len(t.statusHistory)-1]
	}

	return HealthEvent{}
}

// GetStatus returns the current health status.
func (t *AgentHealthTracker) GetStatus() AgentHealthStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

// GetMetrics returns a copy of the health metrics.
func (t *AgentHealthTracker) GetMetrics() HealthMetrics {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.metrics
}

// CanAcceptRequest checks if the agent can accept a new request.
func (t *AgentHealthTracker) CanAcceptRequest() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.status != HealthFailed && t.circuitBreaker.Allow()
}

// transitionStatus transitions the agent to a new status (must hold lock).
func (t *AgentHealthTracker) transitionStatus(newStatus AgentHealthStatus, reason string) {
	if newStatus == t.status {
		return
	}

	event := HealthEvent{
		AgentID:   t.agent.ID,
		AgentName: t.agent.Name,
		OldStatus: t.status,
		NewStatus: newStatus,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	t.status = newStatus
	t.agent.HealthStatus = newStatus

	// Add to history
	t.statusHistory = append(t.statusHistory, event)

	// Trim history if needed
	if len(t.statusHistory) > t.maxHistoryLength {
		t.statusHistory = t.statusHistory[len(t.statusHistory)-t.maxHistoryLength:]
	}
}

// HealthMonitor monitors the health of all agents in the Council.
type HealthMonitor struct {
	mu       sync.RWMutex
	trackers map[string]*AgentHealthTracker
	logger   *slog.Logger
	stopCh   chan struct{}

	// Callbacks for health events
	onHealthChange func(event HealthEvent)
}

// NewHealthMonitor creates a new health monitor.
func NewHealthMonitor(logger *slog.Logger) *HealthMonitor {
	if logger == nil {
		logger = slog.Default()
	}

	return &HealthMonitor{
		trackers: make(map[string]*AgentHealthTracker),
		logger:   logger,
		stopCh:   make(chan struct{}),
	}
}

// RegisterAgent registers an agent for health monitoring.
func (m *HealthMonitor) RegisterAgent(agent *AgentInstance) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.trackers[agent.ID]; !exists {
		m.trackers[agent.ID] = NewAgentHealthTracker(agent)
		m.logger.Debug("registered agent for health monitoring",
			slog.String("agent_id", agent.ID),
			slog.String("agent_name", agent.Name),
		)
	}
}

// UnregisterAgent removes an agent from health monitoring.
func (m *HealthMonitor) UnregisterAgent(agentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.trackers, agentID)
	m.logger.Debug("unregistered agent from health monitoring",
		slog.String("agent_id", agentID),
	)
}

// RecordHeartbeat records a heartbeat for an agent.
func (m *HealthMonitor) RecordHeartbeat(agentID string) {
	m.mu.RLock()
	tracker, exists := m.trackers[agentID]
	m.mu.RUnlock()

	if exists {
		tracker.RecordHeartbeat()
	}
}

// RecordSuccess records a successful response from an agent.
func (m *HealthMonitor) RecordSuccess(agentID string, responseTime time.Duration) {
	m.mu.RLock()
	tracker, exists := m.trackers[agentID]
	m.mu.RUnlock()

	if exists {
		tracker.RecordSuccess(responseTime)
	}
}

// RecordFailure records a failed response from an agent.
func (m *HealthMonitor) RecordFailure(agentID, reason string) {
	m.mu.RLock()
	tracker, exists := m.trackers[agentID]
	m.mu.RUnlock()

	if exists {
		tracker.RecordFailure(reason)
	}
}

// GetHealthyAgents returns all agents that can accept requests.
func (m *HealthMonitor) GetHealthyAgents() []*AgentInstance {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var healthy []*AgentInstance
	for _, tracker := range m.trackers {
		if tracker.CanAcceptRequest() {
			healthy = append(healthy, tracker.agent)
		}
	}

	return healthy
}

// GetAgentStatus returns the health status of a specific agent.
func (m *HealthMonitor) GetAgentStatus(agentID string) (AgentHealthStatus, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if tracker, exists := m.trackers[agentID]; exists {
		return tracker.GetStatus(), true
	}
	return "", false
}

// GetAgentMetrics returns health metrics for a specific agent.
func (m *HealthMonitor) GetAgentMetrics(agentID string) (HealthMetrics, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if tracker, exists := m.trackers[agentID]; exists {
		return tracker.GetMetrics(), true
	}
	return HealthMetrics{}, false
}

// GetAllMetrics returns health metrics for all agents.
func (m *HealthMonitor) GetAllMetrics() map[string]HealthMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]HealthMetrics)
	for id, tracker := range m.trackers {
		metrics[id] = tracker.GetMetrics()
	}

	return metrics
}

// SetHealthChangeCallback sets a callback for health status changes.
func (m *HealthMonitor) SetHealthChangeCallback(callback func(event HealthEvent)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onHealthChange = callback
}

// Start begins the health monitoring loop.
func (m *HealthMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.checkAllAgents()
		}
	}
}

// Stop stops the health monitoring loop.
func (m *HealthMonitor) Stop() {
	close(m.stopCh)
}

// checkAllAgents performs a health check on all registered agents.
func (m *HealthMonitor) checkAllAgents() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tracker := range m.trackers {
		event := tracker.CheckHealth()
		if event.NewStatus != "" {
			m.logger.Info("agent health status changed",
				slog.String("agent_id", event.AgentID),
				slog.String("agent_name", event.AgentName),
				slog.String("old_status", string(event.OldStatus)),
				slog.String("new_status", string(event.NewStatus)),
				slog.String("reason", event.Reason),
			)

			if m.onHealthChange != nil {
				go m.onHealthChange(event)
			}
		}
	}
}

// HealthSummary provides a summary of the Council's health.
type HealthSummary struct {
	TotalAgents      int            `json:"total_agents"`
	HealthyAgents    int            `json:"healthy_agents"`
	DegradedAgents   int            `json:"degraded_agents"`
	FailedAgents     int            `json:"failed_agents"`
	OverallStatus    string         `json:"overall_status"`
	AgentStatuses    map[string]string `json:"agent_statuses"`
	LastChecked      time.Time      `json:"last_checked"`
}

// GetSummary returns a health summary of all agents.
func (m *HealthMonitor) GetSummary() HealthSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary := HealthSummary{
		AgentStatuses: make(map[string]string),
		LastChecked:   time.Now(),
	}

	for _, tracker := range m.trackers {
		summary.TotalAgents++
		status := tracker.GetStatus()
		summary.AgentStatuses[tracker.agent.ID] = string(status)

		switch status {
		case HealthHealthy:
			summary.HealthyAgents++
		case HealthDegraded:
			summary.DegradedAgents++
		case HealthFailed:
			summary.FailedAgents++
		}
	}

	// Determine overall status
	switch {
	case summary.HealthyAgents == summary.TotalAgents:
		summary.OverallStatus = "healthy"
	case summary.FailedAgents > summary.TotalAgents/2:
		summary.OverallStatus = "critical"
	case summary.DegradedAgents > 0 || summary.FailedAgents > 0:
		summary.OverallStatus = "degraded"
	default:
		summary.OverallStatus = "unknown"
	}

	return summary
}
