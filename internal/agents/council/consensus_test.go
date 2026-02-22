// Package council_test provides unit tests for the Council of AIs consensus system.
//
// These tests verify the consensus algorithm behavior including:
//   - Agreement calculation with 3 agents
//   - Semantic equivalence grouping (95% threshold)
//   - Weighted voting based on confidence scores
//   - Consensus threshold verification (default 80%)
package council_test

import (
	"testing"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConsensusAlgorithm_ThreeAgentsAgreement tests that 3 agents reaching agreement
// results in a consensus with high confidence.
func TestConsensusAlgorithm_ThreeAgentsAgreement(t *testing.T) {
	tests := []struct {
		name              string
		threshold         float64
		responses         []*council.AgentResponse
		expectConsensus   bool
		expectedScore     float64
		expectedGroups    int // number of equivalence groups
	}{
		{
			name:      "all_agents_agree_full_consensus",
			threshold: 0.80,
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Aspirin is effective for reducing fever and mild pain.",
					Confidence:   95.0,
					Embedding:    createSimilarEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Aspirin works well for fever reduction and pain relief.",
					Confidence:   92.0,
					Embedding:    createSimilarEmbedding(0.97),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "For fever and mild pain, aspirin is an effective treatment.",
					Confidence:   93.0,
					Embedding:    createSimilarEmbedding(0.96),
				},
			},
			expectConsensus: true,
			expectedScore:   93.33, // Average confidence
			expectedGroups:  1,     // All in one group
		},
		{
			name:      "two_agents_agree_one_dissents",
			threshold: 0.80,
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Aspirin is effective for reducing fever.",
					Confidence:   95.0,
					Embedding:    createSimilarEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Aspirin works well for fever reduction.",
					Confidence:   92.0,
					Embedding:    createSimilarEmbedding(0.97),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Ibuprofen is preferred over aspirin for fever.",
					Confidence:   88.0,
					Embedding:    createDissimilarEmbedding(),
				},
			},
			expectConsensus: true, // 66.7% agreement but weighted by confidence
			expectedScore:   80.0, // Adjusted for dissent
			expectedGroups:  2,    // Two equivalence groups
		},
		{
			name:      "no_consensus_below_threshold",
			threshold: 0.80,
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Aspirin is the best choice.",
					Confidence:   85.0,
					Embedding:    createDissimilarEmbedding(),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Ibuprofen is more effective.",
					Confidence:   82.0,
					Embedding:    createDissimilarEmbedding(),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Paracetamol has fewer side effects.",
					Confidence:   80.0,
					Embedding:    createDissimilarEmbedding(),
				},
			},
			expectConsensus: false,
			expectedScore:   0.0, // No majority
			expectedGroups:  3,   // Three different positions
		},
		{
			name:      "weighted_voting_high_confidence_wins",
			threshold: 0.75,
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Standard dose is 500mg.",
					Confidence:   98.0,
					Embedding:    createSimilarEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "The standard dose is 500mg.",
					Confidence:   97.0,
					Embedding:    createSimilarEmbedding(0.99),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "I think it might be 250mg.",
					Confidence:   60.0,
					Embedding:    createDissimilarEmbedding(),
				},
			},
			expectConsensus: true, // High confidence agents agree
			expectedScore:   90.0, // Weighted average favors high confidence
			expectedGroups:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create consensus calculator
			calculator := council.NewConsensusCalculator(tt.threshold)

			// Calculate consensus
			result, err := calculator.Calculate(tt.responses)

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectConsensus, result.ThresholdMet,
				"Consensus threshold met should match expected")
			assert.InDelta(t, tt.expectedScore, result.AgreementScore, 5.0,
				"Agreement score should be within delta")
			assert.Len(t, result.EquivalenceGroups, tt.expectedGroups,
				"Number of equivalence groups should match")

			if tt.expectConsensus {
				assert.Equal(t, council.StatusConsensus, result.Status,
					"Status should be consensus")
				assert.NotEmpty(t, result.FinalResponse,
					"Final response should not be empty when consensus reached")
			}
		})
	}
}

// TestConsensusAlgorithm_MinimumAgents tests that at least 3 agents are required.
func TestConsensusAlgorithm_MinimumAgents(t *testing.T) {
	calculator := council.NewConsensusCalculator(0.80)

	t.Run("insufficient_agents_returns_error", func(t *testing.T) {
		responses := []*council.AgentResponse{
			{ID: "resp-1", AgentID: "agent-1", ResponseText: "Test", Confidence: 90.0},
		}

		_, err := calculator.Calculate(responses)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum 3 agents required")
	})

	t.Run("exactly_three_agents_allowed", func(t *testing.T) {
		responses := []*council.AgentResponse{
			{ID: "resp-1", AgentID: "agent-1", ResponseText: "Test 1", Confidence: 90.0, Embedding: createSimilarEmbedding(0.98)},
			{ID: "resp-2", AgentID: "agent-2", ResponseText: "Test 2", Confidence: 90.0, Embedding: createSimilarEmbedding(0.98)},
			{ID: "resp-3", AgentID: "agent-3", ResponseText: "Test 3", Confidence: 90.0, Embedding: createSimilarEmbedding(0.98)},
		}

		_, err := calculator.Calculate(responses)
		assert.NoError(t, err)
	})
}

// TestConsensusAlgorithm_ConfidenceScoring tests weighted confidence calculation.
func TestConsensusAlgorithm_ConfidenceScoring(t *testing.T) {
	calculator := council.NewConsensusCalculator(0.80)

	responses := []*council.AgentResponse{
		{ID: "resp-1", AgentID: "agent-1", ResponseText: "Same", Confidence: 100.0, Embedding: createSimilarEmbedding(0.98)},
		{ID: "resp-2", AgentID: "agent-2", ResponseText: "Same", Confidence: 100.0, Embedding: createSimilarEmbedding(0.98)},
		{ID: "resp-3", AgentID: "agent-3", ResponseText: "Same", Confidence: 100.0, Embedding: createSimilarEmbedding(0.98)},
	}

	result, err := calculator.Calculate(responses)
	require.NoError(t, err)

	// With 100% confidence on all agents, overall confidence should be 100
	assert.Equal(t, 100.0, result.ConfidenceScore)
}

// TestConsensusAlgorithm_DissentTracking tests tracking of dissenting agents.
func TestConsensusAlgorithm_DissentTracking(t *testing.T) {
	calculator := council.NewConsensusCalculator(0.80)

	responses := []*council.AgentResponse{
		{ID: "resp-1", AgentID: "agent-1", ResponseText: "A", Confidence: 90.0, Embedding: createSimilarEmbedding(0.98)},
		{ID: "resp-2", AgentID: "agent-2", ResponseText: "A", Confidence: 90.0, Embedding: createSimilarEmbedding(0.98)},
		{ID: "resp-3", AgentID: "agent-dissent", ResponseText: "B", Confidence: 85.0, Embedding: createDissimilarEmbedding()},
	}

	result, err := calculator.Calculate(responses)
	require.NoError(t, err)

	assert.Len(t, result.DissentingAgents, 1)
	assert.Contains(t, result.DissentingAgents, "agent-dissent")
}

// Helper functions

func createSimilarEmbedding(similarity float64) pgvector.Vector {
	// Create a base embedding vector (simplified for testing)
	// In production, these would be actual embeddings from the LLM
	base := make([]float32, 1536)
	for i := range base {
		base[i] = float32(similarity)
	}
	return pgvector.NewVector(base)
}

func createDissimilarEmbedding() pgvector.Vector {
	// Create an embedding that would be semantically different
	vec := make([]float32, 1536)
	for i := range vec {
		vec[i] = float32(0.1) // Low similarity
	}
	return pgvector.NewVector(vec)
}
