// Package council_test provides unit tests for agent instance management.
//
// These tests verify agent instance behavior including:
//   - Timeout handling (3s default)
//   - Circuit breaker pattern
//   - Health status transitions
//   - Graceful degradation on failures
package council_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAgentInstance_TimeoutHandling tests the 3-second timeout for agent responses.
func TestAgentInstance_TimeoutHandling(t *testing.T) {
	tests := []struct {
		name           string
		timeout        time.Duration
		responseDelay  time.Duration
		expectTimeout  bool
	}{
		{
			name:          "response_within_timeout",
			timeout:       3 * time.Second,
			responseDelay: 100 * time.Millisecond,
			expectTimeout: false,
		},
		{
			name:          "response_exactly_at_timeout",
			timeout:       100 * time.Millisecond,
			responseDelay: 100 * time.Millisecond,
			expectTimeout: true, // Boundary case - timeout wins
		},
		{
			name:          "response_exceeds_timeout",
			timeout:       50 * time.Millisecond,
			responseDelay: 200 * time.Millisecond,
			expectTimeout: true,
		},
		{
			name:          "very_fast_response",
			timeout:       3 * time.Second,
			responseDelay: 1 * time.Millisecond,
			expectTimeout: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			agent := &MockAgent{
				ID:           "agent-1",
				ResponseDelay: tt.responseDelay,
				ShouldError:   false,
			}

			wrapper := council.NewAgentWrapper(agent, tt.timeout)

			start := time.Now()
			resp, err := wrapper.Query(ctx, "test query")
			elapsed := time.Since(start)

			if tt.expectTimeout {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, context.DeadlineExceeded) || errors.Is(err, council.ErrAgentTimeout),
					"Error should be timeout")
				assert.Nil(t, resp)
				assert.Less(t, elapsed, tt.timeout+100*time.Millisecond,
					"Should timeout within reasonable time")
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.GreaterOrEqual(t, elapsed, tt.responseDelay,
					"Should wait for response")
			}
		})
	}
}

// TestAgentInstance_ConcurrentQueries tests handling of concurrent requests.
func TestAgentInstance_ConcurrentQueries(t *testing.T) {
	ctx := context.Background()

	agent := &MockAgent{
		ID:            "agent-1",
		ResponseDelay: 50 * time.Millisecond,
	}

	wrapper := council.NewAgentWrapper(agent, 1*time.Second)

	var wg sync.WaitGroup
	results := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := wrapper.Query(ctx, "concurrent query")
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}

	assert.Equal(t, 10, successCount, "All concurrent queries should succeed")
}

// TestAgentInstance_CircuitBreaker tests circuit breaker behavior.
func TestAgentInstance_CircuitBreaker(t *testing.T) {
	ctx := context.Background()

	agent := &MockAgent{
		ID:           "agent-1",
		ShouldError:  true,
		ResponseDelay: 0,
	}

	wrapper := council.NewAgentWrapper(agent, 1*time.Second)
	wrapper.SetCircuitBreakerThreshold(3)

	// First few failures should be allowed through
	for i := 0; i < 3; i++ {
		_, err := wrapper.Query(ctx, "test query")
		assert.Error(t, err, "Query %d should fail", i+1)
		assert.False(t, wrapper.IsCircuitOpen(),
			"Circuit should still be closed after %d failures", i+1)
	}

	// After threshold, circuit should open
	wrapper.RecordFailure(errors.New("failure"))
	assert.True(t, wrapper.IsCircuitOpen(), "Circuit should be open after threshold failures")

	// Requests should be blocked immediately
	_, err := wrapper.Query(ctx, "test query")
	assert.True(t, errors.Is(err, council.ErrCircuitOpen),
		"Should return circuit open error")

	// Wait for cooldown
	time.Sleep(council.CircuitBreakerCooldown + 100*time.Millisecond)

	// Circuit should allow test request
	agent.ShouldError = false
	_, err = wrapper.Query(ctx, "recovery test")
	// After cooldown, circuit enters half-open state
	assert.False(t, wrapper.IsCircuitOpen() || err == nil,
		"Circuit should be in half-open state or recovered")
}

// TestAgentInstance_HealthTransitions tests health status changes.
func TestAgentInstance_HealthTransitions(t *testing.T) {
	monitor := council.NewHealthMonitor(nil)

	agent := &council.AgentInstance{
		ID:           "agent-1",
		Name:         "Test Agent",
		HealthStatus: council.HealthHealthy,
	}

	monitor.RegisterAgent(agent)

	// Initial status should be healthy
	status, exists := monitor.GetAgentStatus("agent-1")
	assert.True(t, exists)
	assert.Equal(t, council.HealthHealthy, status)

	// Record a failure - should transition to degraded
	monitor.RecordFailure("agent-1", "timeout")
	status, _ = monitor.GetAgentStatus("agent-1")
	assert.Equal(t, council.HealthDegraded, status)

	// Record more failures - should transition to failed
	monitor.RecordFailure("agent-1", "timeout")
	monitor.RecordFailure("agent-1", "timeout")
	monitor.RecordFailure("agent-1", "timeout")
	monitor.RecordFailure("agent-1", "timeout")

	status, _ = monitor.GetAgentStatus("agent-1")
	assert.Equal(t, council.HealthFailed, status)

	// Record heartbeat - should not immediately recover failed agent
	monitor.RecordHeartbeat("agent-1")
	status, _ = monitor.GetAgentStatus("agent-1")
	// Failed status persists until explicit recovery
	assert.Equal(t, council.HealthFailed, status)
}

// TestAgentInstance_GracefulDegradation tests behavior with failing agents.
func TestAgentInstance_GracefulDegradation(t *testing.T) {
	monitor := council.NewHealthMonitor(nil)

	// Register 5 agents
	for i := 1; i <= 5; i++ {
		agent := &council.AgentInstance{
			ID:           string(rune('a' + i - 1)),
			Name:         "Agent " + string(rune('A'+i-1)),
			HealthStatus: council.HealthHealthy,
		}
		monitor.RegisterAgent(agent)
	}

	// All agents should be available
	healthy := monitor.GetHealthyAgents()
	assert.Len(t, healthy, 5)

	// Fail 2 agents
	monitor.RecordFailure("a", "error")
	monitor.RecordFailure("b", "error")
	monitor.RecordFailure("b", "error")
	monitor.RecordFailure("b", "error")
	monitor.RecordFailure("b", "error")

	healthy = monitor.GetHealthyAgents()
	assert.Len(t, healthy, 4, "4 agents should still be healthy")

	// Fail more agents - only 2 should remain
	monitor.RecordFailure("c", "error")
	monitor.RecordFailure("c", "error")
	monitor.RecordFailure("c", "error")
	monitor.RecordFailure("c", "error")
	monitor.RecordFailure("c", "error")

	healthy = monitor.GetHealthyAgents()
	assert.Len(t, healthy, 3, "3 agents should remain healthy")

	// Verify minimum agents requirement (3)
	assert.GreaterOrEqual(t, len(healthy), 3,
		"Should have at least 3 healthy agents for consensus")
}

// TestAgentInstance_TimeoutContext tests context cancellation propagation.
func TestAgentInstance_TimeoutContext(t *testing.T) {
	agent := &MockAgent{
		ID:            "agent-1",
		ResponseDelay: 10 * time.Second, // Long delay
	}

	wrapper := council.NewAgentWrapper(agent, 5*time.Second)

	// Create context with shorter timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := wrapper.Query(ctx, "test query")
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded),
		"Should return context deadline exceeded")
	assert.Less(t, elapsed, 200*time.Millisecond,
		"Should respect context timeout")
}

// TestAgentInstance_MetricsCollection tests health metrics tracking.
func TestAgentInstance_MetricsCollection(t *testing.T) {
	monitor := council.NewHealthMonitor(nil)

	agent := &council.AgentInstance{
		ID:           "agent-1",
		Name:         "Test Agent",
		HealthStatus: council.HealthHealthy,
	}
	monitor.RegisterAgent(agent)

	// Record some activity
	monitor.RecordSuccess("agent-1", 100*time.Millisecond)
	monitor.RecordSuccess("agent-1", 150*time.Millisecond)
	monitor.RecordSuccess("agent-1", 120*time.Millisecond)
	monitor.RecordFailure("agent-1", "timeout")

	metrics, exists := monitor.GetAgentMetrics("agent-1")
	require.True(t, exists)

	assert.Equal(t, int64(3), metrics.SuccessfulResponses)
	assert.Equal(t, int64(1), metrics.FailedResponses)
	assert.Equal(t, int64(4), metrics.TotalDeliberations)
	assert.Greater(t, metrics.AverageResponseTime, time.Duration(0))
	assert.Equal(t, 1, metrics.ConsecutiveFailures)
}

// TestAgentInstance_HealthSummary tests overall health summary.
func TestAgentInstance_HealthSummary(t *testing.T) {
	monitor := council.NewHealthMonitor(nil)

	// Register agents with different health states
	agents := []*council.AgentInstance{
		{ID: "a", Name: "A", HealthStatus: council.HealthHealthy},
		{ID: "b", Name: "B", HealthStatus: council.HealthHealthy},
		{ID: "c", Name: "C", HealthStatus: council.HealthDegraded},
		{ID: "d", Name: "D", HealthStatus: council.HealthFailed},
		{ID: "e", Name: "E", HealthStatus: council.HealthHealthy},
	}

	for _, agent := range agents {
		monitor.RegisterAgent(agent)
	}

	// Manually trigger health check
	summary := monitor.GetSummary()

	assert.Equal(t, 5, summary.TotalAgents)
	assert.Equal(t, 3, summary.HealthyAgents)
	assert.Equal(t, 1, summary.DegradedAgents)
	assert.Equal(t, 1, summary.FailedAgents)
	assert.Equal(t, "degraded", summary.OverallStatus) // Has degraded and failed agents
}

// MockAgent implements council.Agent interface for testing
type MockAgent struct {
	ID            string
	ResponseDelay time.Duration
	ShouldError   bool
	mu            sync.Mutex
	callCount     int
}

func (m *MockAgent) Query(ctx context.Context, query string) (*council.AgentResponse, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	if m.ResponseDelay > 0 {
		select {
		case <-time.After(m.ResponseDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if m.ShouldError {
		return nil, errors.New("agent error")
	}

	return &council.AgentResponse{
		ID:           "resp-" + m.ID,
		AgentID:      m.ID,
		ResponseText: "Test response for: " + query,
		Confidence:   95.0,
	}, nil
}

func (m *MockAgent) GetID() string {
	return m.ID
}

func (m *MockAgent) GetName() string {
	return "Mock Agent " + m.ID
}

func (m *MockAgent) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}
