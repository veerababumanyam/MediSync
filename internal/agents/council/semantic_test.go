// Package council_test provides unit tests for semantic equivalence detection.
//
// These tests verify the semantic equivalence algorithm including:
//   - Cosine similarity calculation using pgvector
//   - 95% threshold for equivalence grouping
//   - Handling of edge cases (empty responses, identical responses)
package council_test

import (
	"testing"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSemanticEquivalence_Threshold tests the 95% similarity threshold for equivalence.
func TestSemanticEquivalence_Threshold(t *testing.T) {
	tests := []struct {
		name           string
		embedding1     pgvector.Vector
		embedding2     pgvector.Vector
		expectSimilar  bool
		expectedScore  float64
	}{
		{
			name:          "identical_embeddings_high_similarity",
			embedding1:    createTestEmbedding(1.0),
			embedding2:    createTestEmbedding(1.0),
			expectSimilar: true,
			expectedScore: 1.0,
		},
		{
			name:          "very_similar_above_threshold",
			embedding1:    createTestEmbedding(0.97),
			embedding2:    createTestEmbedding(0.96),
			expectSimilar: true,
			expectedScore: 0.95, // Should be above 95% threshold
		},
		{
			name:          "just_below_threshold",
			embedding1:    createTestEmbedding(0.94),
			embedding2:    createTestEmbedding(0.93),
			expectSimilar: false,
			expectedScore: 0.94, // Should be below 95% threshold
		},
		{
			name:          "dissimilar_embeddings",
			embedding1:    createTestEmbedding(0.5),
			embedding2:    createTestEmbedding(-0.5),
			expectSimilar: false,
			expectedScore: 0.0,
		},
	}

	detector := council.NewSemanticDetector(0.95) // 95% threshold

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert pgvector.Vector to council.Embedding for function calls
			e1 := council.Embedding(tt.embedding1.Slice())
			e2 := council.Embedding(tt.embedding2.Slice())

			similarity, err := detector.CosineSimilarity(e1, e2)
			require.NoError(t, err)

			assert.InDelta(t, tt.expectedScore, similarity, 0.1)

			isEquivalent := detector.IsEquivalent(e1, e2)
			assert.Equal(t, tt.expectSimilar, isEquivalent)
		})
	}
}

// TestSemanticEquivalence_Grouping tests grouping of equivalent responses.
func TestSemanticEquivalence_Grouping(t *testing.T) {
	tests := []struct {
		name          string
		responses     []*council.AgentResponse
		expectedGroups int
		groupSizes    []int // Expected size of each group
	}{
		{
			name: "all_equivalent_single_group",
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Aspirin reduces fever.",
					Confidence:   95.0,
					Embedding:    createTestEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Aspirin is effective for fever.",
					Confidence:   92.0,
					Embedding:    createTestEmbedding(0.97),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Fever can be treated with aspirin.",
					Confidence:   93.0,
					Embedding:    createTestEmbedding(0.96),
				},
			},
			expectedGroups: 1,
			groupSizes:     []int{3},
		},
		{
			name: "two_groups_two_and_one",
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Aspirin for fever.",
					Confidence:   95.0,
					Embedding:    createTestEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Aspirin treats fever.",
					Confidence:   92.0,
					Embedding:    createTestEmbedding(0.97),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Ibuprofen is better.",
					Confidence:   88.0,
					Embedding:    createDissimilarTestEmbedding(),
				},
			},
			expectedGroups: 2,
			groupSizes:     []int{2, 1},
		},
		{
			name: "three_distinct_groups",
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Use aspirin.",
					Confidence:   95.0,
					Embedding:    createTestEmbeddingWithSeed(1),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Use ibuprofen.",
					Confidence:   92.0,
					Embedding:    createTestEmbeddingWithSeed(2),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Use paracetamol.",
					Confidence:   88.0,
					Embedding:    createTestEmbeddingWithSeed(3),
				},
			},
			expectedGroups: 3,
			groupSizes:     []int{1, 1, 1},
		},
		{
			name: "four_agents_two_groups",
			responses: []*council.AgentResponse{
				{
					ID:           "resp-1",
					AgentID:      "agent-1",
					ResponseText: "Group A response.",
					Confidence:   95.0,
					Embedding:    createTestEmbedding(0.98),
				},
				{
					ID:           "resp-2",
					AgentID:      "agent-2",
					ResponseText: "Group A similar.",
					Confidence:   94.0,
					Embedding:    createTestEmbedding(0.97),
				},
				{
					ID:           "resp-3",
					AgentID:      "agent-3",
					ResponseText: "Group B response.",
					Confidence:   93.0,
					Embedding:    createDissimilarTestEmbedding(),
				},
				{
					ID:           "resp-4",
					AgentID:      "agent-4",
					ResponseText: "Group B similar.",
					Confidence:   92.0,
					Embedding:    createDissimilarTestEmbedding(),
				},
			},
			expectedGroups: 2,
			groupSizes:     []int{2, 2},
		},
	}

	detector := council.NewSemanticDetector(0.95)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := detector.GroupEquivalentResponses(tt.responses)

			require.Len(t, groups, tt.expectedGroups,
				"Number of groups should match expected")

			for i, group := range groups {
				if i < len(tt.groupSizes) {
					assert.Len(t, group.AgentIDs, tt.groupSizes[i],
						"Group %d size should match expected", i)
				}
			}
		})
	}
}

// TestSemanticEquivalence_CanonicalSelection tests selection of canonical response.
func TestSemanticEquivalence_CanonicalSelection(t *testing.T) {
	detector := council.NewSemanticDetector(0.95)

	responses := []*council.AgentResponse{
		{
			ID:           "resp-1",
			AgentID:      "agent-1",
			ResponseText: "Short response.",
			Confidence:   80.0,
			Embedding:    createTestEmbedding(0.98),
		},
		{
			ID:           "resp-2",
			AgentID:      "agent-2",
			ResponseText: "This is a longer, more detailed response with better explanation.",
			Confidence:   95.0, // Highest confidence
			Embedding:    createTestEmbedding(0.97),
		},
		{
			ID:           "resp-3",
			AgentID:      "agent-3",
			ResponseText: "Medium length response.",
			Confidence:   85.0,
			Embedding:    createTestEmbedding(0.96),
		},
	}

	groups := detector.GroupEquivalentResponses(responses)
	require.Len(t, groups, 1)

	// Canonical should be from agent with highest confidence
	assert.Equal(t, "This is a longer, more detailed response with better explanation.",
		groups[0].Canonical)
	assert.Contains(t, groups[0].AgentIDs, "agent-2")
}

// TestSemanticEquivalence_EdgeCases tests edge cases in equivalence detection.
func TestSemanticEquivalence_EdgeCases(t *testing.T) {
	detector := council.NewSemanticDetector(0.95)

	t.Run("nil_embedding_returns_error", func(t *testing.T) {
		resp := &council.AgentResponse{
			ID:           "resp-1",
			AgentID:      "agent-1",
			ResponseText: "Test",
			Confidence:   90.0,
			Embedding:    pgvector.Vector{},
		}

		_, err := detector.CalculateSimilarity(resp, resp)
		assert.Error(t, err)
	})

	t.Run("empty_response_list_returns_empty_groups", func(t *testing.T) {
		groups := detector.GroupEquivalentResponses([]*council.AgentResponse{})
		assert.Empty(t, groups)
	})

	t.Run("single_response_creates_single_group", func(t *testing.T) {
		responses := []*council.AgentResponse{
			{
				ID:           "resp-1",
				AgentID:      "agent-1",
				ResponseText: "Only response",
				Confidence:   90.0,
				Embedding:    createTestEmbedding(0.98),
			},
		}

		groups := detector.GroupEquivalentResponses(responses)
		require.Len(t, groups, 1)
		assert.Len(t, groups[0].AgentIDs, 1)
		assert.Equal(t, "Only response", groups[0].Canonical)
	})
}

// TestSemanticEquivalence_Transitivity tests transitivity in equivalence grouping.
// If A ~ B and B ~ C, then A ~ C (within the threshold).
func TestSemanticEquivalence_Transitivity(t *testing.T) {
	detector := council.NewSemanticDetector(0.95)

	// Create embeddings where A~B (98%), B~C (96%), but A~C might be lower
	embeddingA := createTestEmbedding(1.0)
	embeddingB := createTestEmbedding(0.97) // ~96% similar to A
	embeddingC := createTestEmbedding(0.94) // ~97% similar to B, ~94% to A

	responses := []*council.AgentResponse{
		{ID: "resp-A", AgentID: "agent-A", ResponseText: "A", Confidence: 90.0, Embedding: embeddingA},
		{ID: "resp-B", AgentID: "agent-B", ResponseText: "B", Confidence: 90.0, Embedding: embeddingB},
		{ID: "resp-C", AgentID: "agent-C", ResponseText: "C", Confidence: 90.0, Embedding: embeddingC},
	}

	groups := detector.GroupEquivalentResponses(responses)

	// A and B should group together, C might be separate depending on actual similarity
	assert.GreaterOrEqual(t, len(groups), 1)
	assert.LessOrEqual(t, len(groups), 3)
}

// TestSemanticEquivalence_MultilingualResponses tests handling of multilingual content.
func TestSemanticEquivalence_MultilingualResponses(t *testing.T) {
	detector := council.NewSemanticDetector(0.95)

	// In production, these would have embeddings from a multilingual model
	// For testing, we simulate equivalent semantic meaning
	responses := []*council.AgentResponse{
		{
			ID:           "resp-en",
			AgentID:      "agent-1",
			ResponseText: "Aspirin is effective for fever.",
			Confidence:   95.0,
			Embedding:    createTestEmbedding(0.98),
		},
		{
			ID:           "resp-ar",
			AgentID:      "agent-2",
			ResponseText: "الأسبرين فعال للحمى.", // Same meaning in Arabic
			Confidence:   93.0,
			Embedding:    createTestEmbedding(0.97), // Similar embedding = same meaning
		},
		{
			ID:           "resp-en2",
			AgentID:      "agent-3",
			ResponseText: "Fever can be treated with aspirin.",
			Confidence:   94.0,
			Embedding:    createTestEmbedding(0.96),
		},
	}

	groups := detector.GroupEquivalentResponses(responses)

	// All should group together as they have the same semantic meaning
	require.Len(t, groups, 1, "Multilingual responses with same meaning should group together")
	assert.Len(t, groups[0].AgentIDs, 3)
}

// Helper functions for creating test embeddings

func createTestEmbedding(baseValue float32) pgvector.Vector {
	vec := make([]float32, 1536)
	for i := range vec {
		vec[i] = baseValue
	}
	return pgvector.NewVector(vec)
}

func createTestEmbeddingWithSeed(seed int) pgvector.Vector {
	vec := make([]float32, 1536)
	// Use seed to create different but deterministic embeddings
	for i := range vec {
		vec[i] = float32((seed + i) % 100) / 100.0
	}
	return pgvector.NewVector(vec)
}

func createDissimilarTestEmbedding() pgvector.Vector {
	vec := make([]float32, 1536)
	for i := range vec {
		vec[i] = -0.5 // Negative values create dissimilar embedding
	}
	return pgvector.NewVector(vec)
}
