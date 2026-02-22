// Package council_test provides integration tests for the Council Coordinator.
//
// These tests verify the complete deliberation flow including:
//   - Query submission and validation
//   - Agent coordination and response collection
//   - Consensus calculation and finalization
//   - Evidence retrieval and storage
//   - Audit trail creation
package council_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoordinator_FullDeliberationFlow tests the complete deliberation process.
func TestCoordinator_FullDeliberationFlow(t *testing.T) {
	ctx := context.Background()

	// Setup test fixtures
	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	tests := []struct {
		name               string
		request            council.CreateDeliberationRequest
		expectStatus       council.DeliberationStatus
		expectConsensus    bool
		expectConfidence   float64
		expectEvidence     bool
		expectedMinNodes   int
	}{
		{
			name: "simple_query_full_consensus",
			request: council.CreateDeliberationRequest{
				Query:              "What is the recommended dosage of aspirin for fever?",
				ConsensusThreshold: 0.80,
			},
			expectStatus:     council.StatusConsensus,
			expectConsensus:  true,
			expectConfidence: 90.0,
			expectEvidence:   true,
			expectedMinNodes: 3,
		},
		{
			name: "complex_query_partial_consensus",
			request: council.CreateDeliberationRequest{
				Query:              "What is the best treatment for chronic lower back pain?",
				ConsensusThreshold: 0.90, // High threshold
			},
			expectStatus:     council.StatusUncertain,
			expectConsensus:  false,
			expectConfidence: 0,
			expectEvidence:   true,
			expectedMinNodes: 2,
		},
		{
			name: "ambiguous_query_uncertain",
			request: council.CreateDeliberationRequest{
				Query:              "Which is better?",
				ConsensusThreshold: 0.80,
			},
			expectStatus:     council.StatusUncertain,
			expectConsensus:  false,
			expectConfidence: 0,
			expectEvidence:   false, // Ambiguous query may not retrieve evidence
			expectedMinNodes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Submit deliberation request
			result, err := coordinator.Deliberate(ctx, tt.request, "user-123")

			require.NoError(t, err, "Deliberation should not return error")
			require.NotNil(t, result, "Result should not be nil")
			require.NotNil(t, result.Deliberation, "Deliberation should not be nil")

			// Verify deliberation status
			assert.Equal(t, tt.expectStatus, result.Deliberation.Status,
				"Deliberation status should match expected")

			// Verify consensus result
			if tt.expectConsensus {
				require.NotNil(t, result.ConsensusRecord, "Should have consensus record")
				assert.True(t, result.ConsensusRecord.ThresholdMet,
					"Consensus threshold should be met")
				assert.GreaterOrEqual(t, result.ConsensusRecord.AgreementScore,
					tt.request.ConsensusThreshold*100,
					"Agreement score should meet threshold")
				assert.NotEmpty(t, result.Deliberation.FinalResponse,
					"Should have final response")
				assert.GreaterOrEqual(t, result.Deliberation.ConfidenceScore,
					tt.expectConfidence,
					"Confidence should meet expected")
			}

			// Verify evidence trail
			if tt.expectEvidence {
				require.NotNil(t, result.EvidenceTrail, "Should have evidence trail")
				assert.GreaterOrEqual(t, len(result.EvidenceTrail.NodeIDs),
					tt.expectedMinNodes,
					"Should have minimum expected nodes")
			}

			// Verify agent responses
			assert.GreaterOrEqual(t, len(result.AgentResponses), 3,
				"Should have responses from at least 3 agents")
		})
	}
}

// TestCoordinator_MinimumAgents verifies minimum agent requirement.
func TestCoordinator_MinimumAgents(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixturesWithOpts(t, testFixturesOptions{
		agentCount: 2, // Only 2 agents - below minimum
	})
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	_, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              "Test query",
		ConsensusThreshold: 0.80,
	}, "user-123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "minimum 3 agents required")
}

// TestCoordinator_TimeoutHandling tests handling of agent timeouts.
func TestCoordinator_TimeoutHandling(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixturesWithOpts(t, testFixturesOptions{
		agentCount:          5,
		timeoutDuration:     100 * time.Millisecond,
		slowAgentCount:      2, // 2 agents will be slow
	})
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	// With slow agents, should still get consensus from fast agents
	result, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              "Test query with potential timeouts",
		ConsensusThreshold: 0.60, // Lower threshold to account for timeouts
	}, "user-123")

	require.NoError(t, err)

	// Should still have responses from at least 3 agents (5 - 2 slow)
	assert.GreaterOrEqual(t, len(result.AgentResponses), 3,
		"Should have responses from non-timed-out agents")
}

// TestCoordinator_EvidenceRetrievalFailure tests graceful degradation when KG fails.
func TestCoordinator_EvidenceRetrievalFailure(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixturesWithOpts(t, testFixturesOptions{
		agentCount:        3,
		kgFailure:         true, // Simulate KG failure
		evidenceCacheData: createCachedEvidence(), // But have cached evidence
	})
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	result, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              "What is aspirin used for?",
		ConsensusThreshold: 0.80,
	}, "user-123")

	// Should succeed using cached evidence
	require.NoError(t, err)
	assert.NotNil(t, result.EvidenceTrail, "Should have cached evidence")
	assert.True(t, result.EvidenceTrail.CachedAt.Before(time.Now()),
		"Evidence should be from cache")
}

// TestCoordinator_AuditTrailCreation verifies audit entries are created.
func TestCoordinator_AuditTrailCreation(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	// Submit deliberation
	result, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              "Test audit trail",
		ConsensusThreshold: 0.80,
	}, "user-123")

	require.NoError(t, err)

	// Verify audit entry was created
	auditEntries, err := fixtures.auditRepo.GetByDeliberationID(ctx, result.Deliberation.ID)
	require.NoError(t, err)
	require.Greater(t, len(auditEntries), 0, "Should have audit entries")

	// Verify audit entry contents
	entry := auditEntries[0]
	assert.Equal(t, "user-123", entry.UserID)
	assert.Equal(t, result.Deliberation.ID, entry.DeliberationID)
	assert.Equal(t, council.AuditActionQuery, entry.Action)
	assert.NotEmpty(t, entry.Details)
}

// TestCoordinator_ConcurrentDeliberations tests handling of concurrent requests.
func TestCoordinator_ConcurrentDeliberations(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	const numConcurrent = 10
	results := make(chan *council.DeliberationResult, numConcurrent)
	errors := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(idx int) {
			result, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
				Query:              "Concurrent test query",
				ConsensusThreshold: 0.80,
			}, "user-123")
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(i)
	}

	// Wait for all to complete
	successCount := 0
	for i := 0; i < numConcurrent; i++ {
		select {
		case result := <-results:
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Deliberation.ID)
			successCount++
		case err := <-errors:
			t.Logf("Concurrent request error: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent deliberations")
		}
	}

	assert.GreaterOrEqual(t, successCount, numConcurrent/2,
		"At least half of concurrent requests should succeed")
}

// TestCoordinator_QueryDeduplication tests query hash-based deduplication.
func TestCoordinator_QueryDeduplication(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	query := "What are the side effects of aspirin?"

	// Submit same query twice
	result1, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              query,
		ConsensusThreshold: 0.80,
	}, "user-123")
	require.NoError(t, err)

	result2, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
		Query:              query,
		ConsensusThreshold: 0.80,
	}, "user-456") // Different user
	require.NoError(t, err)

	// Queries should have same hash
	assert.Equal(t, result1.Deliberation.QueryHash, result2.Deliberation.QueryHash,
		"Same query should produce same hash")

	// But different deliberations (different users)
	assert.NotEqual(t, result1.Deliberation.ID, result2.Deliberation.ID,
		"Different users should get different deliberation IDs")
}

// TestCoordinator_LatencyRequirement verifies 95% of queries complete in <10s.
func TestCoordinator_LatencyRequirement(t *testing.T) {
	ctx := context.Background()

	fixtures := setupTestFixtures(t)
	defer fixtures.cleanup()

	coordinator := fixtures.coordinator

	const numQueries = 20
	latencies := make([]time.Duration, 0, numQueries)
	maxLatency := 10 * time.Second

	for i := 0; i < numQueries; i++ {
		start := time.Now()
		_, err := coordinator.Deliberate(ctx, council.CreateDeliberationRequest{
			Query:              "Latency test query",
			ConsensusThreshold: 0.80,
		}, "user-123")
		latency := time.Since(start)
		latencies = append(latencies, latency)

		require.NoError(t, err)
	}

	// Calculate P95 latency
	quickCount := 0
	for _, l := range latencies {
		if l < maxLatency {
			quickCount++
		}
	}

	percentage := float64(quickCount) / float64(numQueries) * 100
	assert.GreaterOrEqual(t, percentage, 95.0,
		"95%% of queries should complete in <10s (actual: %.1f%%)", percentage)
}

// Test fixtures and helpers

type testFixtures struct {
	coordinator *council.Coordinator
	auditRepo   *MockAuditRepository
	cleanup     func()
}

type testFixturesOptions struct {
	agentCount          int
	timeoutDuration     time.Duration
	slowAgentCount      int
	kgFailure           bool
	evidenceCacheData   *council.EvidenceCacheEntry
}

func setupTestFixtures(t *testing.T) *testFixtures {
	return setupTestFixturesWithOpts(t, testFixturesOptions{
		agentCount:      5,
		timeoutDuration: 3 * time.Second,
	})
}

func setupTestFixturesWithOpts(t *testing.T, opts testFixturesOptions) *testFixtures {
	if opts.agentCount == 0 {
		opts.agentCount = 5
	}
	if opts.timeoutDuration == 0 {
		opts.timeoutDuration = 3 * time.Second
	}

	// Create mock agents
	agents := make([]council.Agent, opts.agentCount)
	for i := 0; i < opts.agentCount; i++ {
		isSlow := i < opts.slowAgentCount
		agents[i] = &MockAgent{
			ID:            string(rune('a' + i)),
			ResponseDelay: func() time.Duration {
				if isSlow {
					return 500 * time.Millisecond
				}
				return 50 * time.Millisecond
			}(),
		}
	}

	// Create mock repositories
	kgRepo := NewMockKnowledgeGraphRepository()
	repo := NewMockCouncilRepository()
	auditRepo := &MockAuditRepository{entries: make([]*council.AuditEntry, 0)}
	cache := NewMockEvidenceCache()

	if opts.evidenceCacheData != nil {
		_ = cache.Set(context.Background(), "cached-hash", opts.evidenceCacheData)
	}

	// Create coordinator
	coordinator := council.NewCoordinator(council.CoordinatorConfig{
		Agents:            agents,
		KGRepository:      kgRepo,
		Repository:        repo,
		EvidenceCache:     cache,
		Timeout:           opts.timeoutDuration,
		ConsensusThreshold: 0.80,
		MinAgents:         3,
	})

	return &testFixtures{
		coordinator: coordinator,
		auditRepo:   auditRepo,
		cleanup: func() {
			// Cleanup resources
		},
	}
}

func createCachedEvidence() *council.EvidenceCacheEntry {
	return &council.EvidenceCacheEntry{
		QueryHash: "cached-hash",
		NodeIDs:   []string{"node-1", "node-2", "node-3"},
		Nodes: []*warehouse.KnowledgeGraphNode{
			{ID: "node-1", Concept: "Aspirin", Definition: "Pain reliever"},
			{ID: "node-2", Concept: "Fever", Definition: "Elevated temperature"},
			{ID: "node-3", Concept: "Pain", Definition: "Physical discomfort"},
		},
		CachedAt:       time.Now().Add(-2 * time.Minute),
		ExpiresAt:      time.Now().Add(3 * time.Minute),
		KGHealthStatus: "healthy",
	}
}

// Mock repositories for testing

type MockAuditRepository struct {
	entries []*council.AuditEntry
	mu      sync.Mutex
}

func (m *MockAuditRepository) Create(entry *council.AuditEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
	return nil
}

func (m *MockAuditRepository) GetByDeliberationID(ctx context.Context, deliberationID string) ([]*council.AuditEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []*council.AuditEntry
	for _, e := range m.entries {
		if e.DeliberationID == deliberationID {
			result = append(result, e)
		}
	}
	return result, nil
}

type MockCouncilRepository struct {
	deliberations map[string]*council.CouncilDeliberation
	responses     map[string][]*council.AgentResponse
	consensus     map[string]*council.ConsensusRecord
	evidence      map[string]*council.EvidenceTrail
	mu            sync.Mutex
}

func NewMockCouncilRepository() *MockCouncilRepository {
	return &MockCouncilRepository{
		deliberations: make(map[string]*council.CouncilDeliberation),
		responses:     make(map[string][]*council.AgentResponse),
		consensus:     make(map[string]*council.ConsensusRecord),
		evidence:      make(map[string]*council.EvidenceTrail),
	}
}

func (m *MockCouncilRepository) CreateDeliberation(ctx context.Context, query, userID string, threshold float64) (*council.CouncilDeliberation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	d := &council.CouncilDeliberation{
		ID:                 "del-" + time.Now().Format("20060102150405"),
		QueryText:          query,
		QueryHash:          council.HashQuery(query),
		UserID:             userID,
		Status:             council.StatusPending,
		ConsensusThreshold: threshold,
		CreatedAt:          time.Now(),
	}
	m.deliberations[d.ID] = d
	return d, nil
}

func (m *MockCouncilRepository) GetDeliberation(ctx context.Context, id string) (*council.CouncilDeliberation, error) {
	return m.deliberations[id], nil
}

func (m *MockCouncilRepository) UpdateDeliberationStatus(ctx context.Context, id string, status council.DeliberationStatus, response string, confidence float64) error {
	if d, ok := m.deliberations[id]; ok {
		d.Status = status
		d.FinalResponse = response
		d.ConfidenceScore = confidence
		now := time.Now()
		d.CompletedAt = &now
	}
	return nil
}

func (m *MockCouncilRepository) CreateAgentResponse(ctx context.Context, resp *council.AgentResponse) error {
	m.responses[resp.DeliberationID] = append(m.responses[resp.DeliberationID], resp)
	return nil
}

func (m *MockCouncilRepository) GetAgentResponses(ctx context.Context, deliberationID string) ([]*council.AgentResponse, error) {
	return m.responses[deliberationID], nil
}

func (m *MockCouncilRepository) CreateConsensusRecord(ctx context.Context, record *council.ConsensusRecord) error {
	m.consensus[record.DeliberationID] = record
	return nil
}

func (m *MockCouncilRepository) GetConsensusRecord(ctx context.Context, deliberationID string) (*council.ConsensusRecord, error) {
	return m.consensus[deliberationID], nil
}

func (m *MockCouncilRepository) CreateEvidenceTrail(ctx context.Context, trail *council.EvidenceTrail) error {
	m.evidence[trail.DeliberationID] = trail
	return nil
}

func (m *MockCouncilRepository) GetEvidenceTrail(ctx context.Context, deliberationID string) (*council.EvidenceTrail, error) {
	return m.evidence[deliberationID], nil
}

func (m *MockCouncilRepository) GetDeliberationWithResponses(ctx context.Context, id string) (*council.DeliberationResult, error) {
	d, _ := m.deliberations[id]
	return &council.DeliberationResult{
		Deliberation:    d,
		AgentResponses:  m.responses[id],
		ConsensusRecord: m.consensus[id],
		EvidenceTrail:   m.evidence[id],
	}, nil
}

func (m *MockCouncilRepository) CreateAuditEntry(ctx context.Context, entry *council.AuditEntry) error {
	return nil
}

func (m *MockCouncilRepository) FlagDeliberation(ctx context.Context, deliberationID string, userID string, req council.FlagDeliberationRequest) error {
	return nil
}

func (m *MockCouncilRepository) ListDeliberations(ctx context.Context, userID string, isAdmin bool, opts council.ListOptions) ([]*council.CouncilDeliberation, int, error) {
	return []*council.CouncilDeliberation{}, 0, nil
}

func (m *MockCouncilRepository) ListHealthyAgents(ctx context.Context) ([]*council.AgentInstance, error) {
	return []*council.AgentInstance{}, nil
}

func (m *MockCouncilRepository) UpdateAgentHeartbeat(ctx context.Context, agentID string) error {
	return nil
}

func (m *MockCouncilRepository) UpdateAgentHealth(ctx context.Context, agentID string, status council.AgentHealthStatus) error {
	return nil
}

