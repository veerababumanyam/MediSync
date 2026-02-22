// Package council_test provides unit tests for Graph-of-Thoughts evidence retrieval.
//
// These tests verify the evidence retrieval system including:
//   - Multi-hop traversal through the Knowledge Graph
//   - Relevance scoring for retrieved nodes
//   - Cycle detection and prevention
//   - Evidence caching with TTL
package council_test

import (
	"context"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGraphOfThoughts_MultiHopTraversal tests multi-hop KG traversal.
func TestGraphOfThoughts_MultiHopTraversal(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name              string
		initialNodes      []string
		maxHops           int
		expectedNodes     int // Minimum expected nodes
		expectedMaxHops   int // Expected max hop count
	}{
		{
			name: "single_hop_expansion",
			initialNodes: []string{"node-aspirin"},
			maxHops:       1,
			expectedNodes: 3, // Initial + related nodes
			expectedMaxHops: 1,
		},
		{
			name: "two_hop_expansion",
			initialNodes: []string{"node-fever"},
			maxHops:       2,
			expectedNodes: 5,
			expectedMaxHops: 2,
		},
		{
			name: "three_hop_expansion",
			initialNodes: []string{"node-headache"},
			maxHops:       3,
			expectedNodes: 7,
			expectedMaxHops: 3,
		},
		{
			name: "multiple_initial_nodes",
			initialNodes: []string{"node-aspirin", "node-ibuprofen"},
			maxHops:       2,
			expectedNodes: 8,
			expectedMaxHops: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock knowledge graph
			kgRepo := NewMockKnowledgeGraphRepository()

			// Create evidence retriever
			retriever := council.NewEvidenceRetriever(kgRepo, 3) // max 3 hops

			// Retrieve evidence
			result, err := retriever.Retrieve(ctx, tt.initialNodes, tt.maxHops)

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.GreaterOrEqual(t, len(result.NodeIDs), tt.expectedNodes,
				"Should retrieve at least expected number of nodes")
			assert.LessOrEqual(t, result.HopCount, tt.expectedMaxHops,
				"Hop count should not exceed max hops")

			// Verify no duplicate nodes
			uniqueNodes := make(map[string]bool)
			for _, id := range result.NodeIDs {
				assert.False(t, uniqueNodes[id], "Node should not be duplicated")
				uniqueNodes[id] = true
			}
		})
	}
}

// TestGraphOfThoughts_CyclePrevention tests that traversal doesn't loop infinitely.
func TestGraphOfThoughts_CyclePrevention(t *testing.T) {
	ctx := context.Background()

	// Create a mock KG with cycles
	kgRepo := NewMockCyclicKnowledgeGraph()

	retriever := council.NewEvidenceRetriever(kgRepo, 10) // High hop limit

	// This should complete without timing out
	result, err := retriever.Retrieve(ctx, []string{"node-a"}, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.LessOrEqual(t, result.HopCount, 10, "Should not exceed max hops")
	assert.LessOrEqual(t, len(result.NodeIDs), 100, "Should have reasonable node count")
}

// TestGraphOfThoughts_RelevanceScoring tests relevance score calculation.
func TestGraphOfThoughts_RelevanceScoring(t *testing.T) {
	ctx := context.Background()

	kgRepo := NewMockKnowledgeGraphRepository()
	retriever := council.NewEvidenceRetriever(kgRepo, 3)

	result, err := retriever.Retrieve(ctx, []string{"node-aspirin"}, 2)
	require.NoError(t, err)

	// Initial nodes should have highest relevance
	for _, initialID := range []string{"node-aspirin"} {
		score, exists := result.RelevanceScore[initialID]
		assert.True(t, exists, "Initial node should have relevance score")
		assert.GreaterOrEqual(t, score, 0.8, "Initial node should have high relevance")
	}

	// All nodes should have relevance scores between 0 and 1
	for nodeID, score := range result.RelevanceScore {
		assert.GreaterOrEqual(t, score, 0.0, "Score for %s should be >= 0", nodeID)
		assert.LessOrEqual(t, score, 1.0, "Score for %s should be <= 1", nodeID)
	}
}

// TestGraphOfThoughts_EdgeTypeFiltering tests filtering by edge types.
func TestGraphOfThoughts_EdgeTypeFiltering(t *testing.T) {
	ctx := context.Background()

	kgRepo := NewMockKnowledgeGraphRepository()
	retriever := council.NewEvidenceRetriever(kgRepo, 3)

	// Retrieve only TREATS relationships
	result, err := retriever.RetrieveWithFilters(ctx,
		[]string{"node-fever"},
		2,
		[]warehouse.KGEdgeType{warehouse.EdgeTreats},
	)

	require.NoError(t, err)

	// Verify all traversal steps use TREATS edge
	for _, step := range result.TraversalPath {
		assert.Equal(t, warehouse.EdgeTreats, step.EdgeType,
			"Should only traverse TREATS edges")
	}
}

// TestGraphOfThoughts_EmptyResult tests handling of empty results.
func TestGraphOfThoughts_EmptyResult(t *testing.T) {
	ctx := context.Background()

	kgRepo := NewMockEmptyKnowledgeGraph()
	retriever := council.NewEvidenceRetriever(kgRepo, 3)

	result, err := retriever.Retrieve(ctx, []string{"nonexistent-node"}, 2)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.NodeIDs)
	assert.Empty(t, result.TraversalPath)
}

// TestGraphOfThoughts_InsufficientKnowledge tests detection of insufficient evidence.
func TestGraphOfThoughts_InsufficientKnowledge(t *testing.T) {
	_ = context.Background() // ctx not needed for this test

	tests := []struct {
		name            string
		result          *council.EvidenceTrail
		expectedStatus  council.EvidenceStatus
	}{
		{
			name: "no_evidence_found",
			result: &council.EvidenceTrail{
				NodeIDs:        []string{},
				RelevanceScore: map[string]float64{},
				HopCount:       0,
			},
			expectedStatus: council.EvidenceStatusInsufficient,
		},
		{
			name: "single_low_relevance_node",
			result: &council.EvidenceTrail{
				NodeIDs: []string{"node-weak"},
				RelevanceScore: map[string]float64{
					"node-weak": 0.2,
				},
				HopCount: 1,
			},
			expectedStatus: council.EvidenceStatusInsufficient,
		},
		{
			name: "sufficient_evidence",
			result: &council.EvidenceTrail{
				NodeIDs: []string{"node-1", "node-2", "node-3"},
				RelevanceScore: map[string]float64{
					"node-1": 0.9,
					"node-2": 0.85,
					"node-3": 0.8,
				},
				HopCount: 2,
			},
			expectedStatus: council.EvidenceStatusSufficient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := council.AssessEvidenceSufficiency(tt.result, 3, 0.7)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

// TestGraphOfThoughts_TraversalPath tests that traversal path is recorded correctly.
func TestGraphOfThoughts_TraversalPath(t *testing.T) {
	ctx := context.Background()

	kgRepo := NewMockKnowledgeGraphRepository()
	retriever := council.NewEvidenceRetriever(kgRepo, 2)

	result, err := retriever.Retrieve(ctx, []string{"node-aspirin"}, 2)
	require.NoError(t, err)

	// Verify traversal path entries are valid
	for _, step := range result.TraversalPath {
		assert.NotEmpty(t, step.FromNodeID, "FromNodeID should not be empty")
		assert.NotEmpty(t, step.ToNodeID, "ToNodeID should not be empty")
		assert.NotEmpty(t, step.EdgeType, "EdgeType should not be empty")
		assert.GreaterOrEqual(t, step.Weight, 0.0, "Weight should be >= 0")
		assert.LessOrEqual(t, step.Weight, 1.0, "Weight should be <= 1")
	}
}

// TestEvidenceCache_TTL tests evidence caching with TTL.
func TestEvidenceCache_TTL(t *testing.T) {
	ctx := context.Background()

	mockCache := NewMockEvidenceCache()

	entry := &council.EvidenceCacheEntry{
		QueryHash: "hash123",
		NodeIDs:   []string{"node-1", "node-2"},
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Set cache entry
	err := mockCache.Set(ctx, "hash123", entry)
	require.NoError(t, err)

	// Retrieve immediately - should succeed
	cached, err := mockCache.Get(ctx, "hash123")
	require.NoError(t, err)
	assert.NotNil(t, cached)

	// Simulate TTL expiry
	mockCache.ExpireAll()

	// Retrieve after expiry - should return nil
	cached, err = mockCache.Get(ctx, "hash123")
	require.NoError(t, err)
	assert.Nil(t, cached)
}

// TestEvidenceCache_Deduplication tests query hash-based deduplication.
func TestEvidenceCache_Deduplication(t *testing.T) {
	ctx := context.Background()

	mockCache := NewMockEvidenceCache()

	query1 := "What are the side effects of aspirin?"
	query2 := "What are the side effects of aspirin?" // Same query
	query3 := "What are side effects of ibuprofen?"   // Different query

	hash1 := council.HashQuery(query1)
	hash2 := council.HashQuery(query2)
	hash3 := council.HashQuery(query3)

	// Same query should produce same hash
	assert.Equal(t, hash1, hash2, "Same queries should have same hash")

	// Different query should produce different hash
	assert.NotEqual(t, hash1, hash3, "Different queries should have different hashes")

	// Cache entry for first query
	entry := &council.EvidenceCacheEntry{
		QueryHash: hash1,
		NodeIDs:   []string{"node-1"},
		CachedAt:  time.Now(),
	}

	err := mockCache.Set(ctx, hash1, entry)
	require.NoError(t, err)

	// Should retrieve same entry for both identical queries
	cached1, _ := mockCache.Get(ctx, hash1)
	cached2, _ := mockCache.Get(ctx, hash2)

	assert.Equal(t, cached1.QueryHash, cached2.QueryHash)
	assert.Equal(t, cached1.NodeIDs, cached2.NodeIDs)
}

// Mock types for testing

type MockKnowledgeGraphRepository struct {
	nodes map[string]*warehouse.KnowledgeGraphNode
}

func NewMockKnowledgeGraphRepository() *MockKnowledgeGraphRepository {
	return &MockKnowledgeGraphRepository{
		nodes: createMockNodes(),
	}
}

func createMockNodes() map[string]*warehouse.KnowledgeGraphNode {
	return map[string]*warehouse.KnowledgeGraphNode{
		"node-aspirin": {
			ID:         "node-aspirin",
			NodeType:   warehouse.NodeTypeMedication,
			Concept:    "Aspirin",
			Definition: "A salicylate drug used to treat pain and fever",
			Edges:      []string{"node-fever", "node-pain", "node-bleeding"},
			EdgeTypes:  []warehouse.KGEdgeType{warehouse.EdgeTreats, warehouse.EdgeTreats, warehouse.EdgeCauses},
			Embedding:  createMockEmbedding(0.9),
		},
		"node-fever": {
			ID:         "node-fever",
			NodeType:   warehouse.NodeTypeCondition,
			Concept:    "Fever",
			Definition: "Elevated body temperature",
			Edges:      []string{"node-aspirin", "node-ibuprofen"},
			EdgeTypes:  []warehouse.KGEdgeType{warehouse.EdgeTreats, warehouse.EdgeTreats},
			Embedding:  createMockEmbedding(0.85),
		},
		"node-ibuprofen": {
			ID:         "node-ibuprofen",
			NodeType:   warehouse.NodeTypeMedication,
			Concept:    "Ibuprofen",
			Definition: "NSAID used for pain and inflammation",
			Edges:      []string{"node-fever", "node-pain"},
			EdgeTypes:  []warehouse.KGEdgeType{warehouse.EdgeTreats, warehouse.EdgeTreats},
			Embedding:  createMockEmbedding(0.88),
		},
	}
}

func (m *MockKnowledgeGraphRepository) GetNode(ctx context.Context, id string) (*warehouse.KnowledgeGraphNode, error) {
	return m.nodes[id], nil
}

func (m *MockKnowledgeGraphRepository) GetNodes(ctx context.Context, ids []string) ([]*warehouse.KnowledgeGraphNode, error) {
	var nodes []*warehouse.KnowledgeGraphNode
	for _, id := range ids {
		if node, ok := m.nodes[id]; ok {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func (m *MockKnowledgeGraphRepository) GetRelatedNodes(ctx context.Context, nodeID string, edgeTypes []warehouse.KGEdgeType, limit int) ([]*warehouse.KnowledgeGraphNode, error) {
	node, ok := m.nodes[nodeID]
	if !ok {
		return nil, nil
	}

	var related []*warehouse.KnowledgeGraphNode
	for i, edgeID := range node.Edges {
		if len(related) >= limit {
			break
		}
		if relatedNode, ok := m.nodes[edgeID]; ok {
			// Filter by edge types if specified
			if len(edgeTypes) == 0 {
				related = append(related, relatedNode)
			} else {
				for _, et := range edgeTypes {
					if i < len(node.EdgeTypes) && node.EdgeTypes[i] == et {
						related = append(related, relatedNode)
						break
					}
				}
			}
		}
	}
	return related, nil
}

func (m *MockKnowledgeGraphRepository) TraverseMultiHop(ctx context.Context, initialNodeIDs []string, maxHops int) (*warehouse.TraversalResult, error) {
	result := &warehouse.TraversalResult{
		Nodes:          []*warehouse.KnowledgeGraphNode{},
		TraversalPath:  []*warehouse.TraversalStep{},
		RelevanceScore: make(map[string]float64),
	}

	visited := make(map[string]bool)
	queue := make([]struct {
		id       string
		hop      int
		fromID   string
		edgeType warehouse.KGEdgeType
		weight   float64
	}, 0)

	// Initialize queue with starting nodes
	for _, id := range initialNodeIDs {
		queue = append(queue, struct {
			id       string
			hop      int
			fromID   string
			edgeType warehouse.KGEdgeType
			weight   float64
		}{id: id, hop: 0, weight: 1.0})
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.id] || current.hop > maxHops {
			continue
		}
		visited[current.id] = true

		node, ok := m.nodes[current.id]
		if !ok {
			continue
		}

		result.Nodes = append(result.Nodes, node)
		result.RelevanceScore[current.id] = 1.0 - (float64(current.hop) * 0.2)

		if current.hop > 0 {
			result.TraversalPath = append(result.TraversalPath, &warehouse.TraversalStep{
				FromNodeID: current.fromID,
				ToNodeID:   current.id,
				EdgeType:   current.edgeType,
				Weight:     current.weight,
			})
		}

		// Add related nodes to queue
		for i, edgeID := range node.Edges {
			if !visited[edgeID] && current.hop < maxHops {
				edgeType := warehouse.EdgeRelatedTo
				if i < len(node.EdgeTypes) {
					edgeType = node.EdgeTypes[i]
				}
				queue = append(queue, struct {
					id       string
					hop      int
					fromID   string
					edgeType warehouse.KGEdgeType
					weight   float64
				}{
					id:       edgeID,
					hop:      current.hop + 1,
					fromID:   current.id,
					edgeType: edgeType,
					weight:   0.8,
				})
			}
		}
	}

	return result, nil
}

func (m *MockKnowledgeGraphRepository) FindSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]*warehouse.KnowledgeGraphNode, error) {
	return []*warehouse.KnowledgeGraphNode{}, nil
}

func (m *MockKnowledgeGraphRepository) HealthCheck(ctx context.Context) error {
	return nil
}

// Mock cyclic knowledge graph
type MockCyclicKnowledgeGraph struct {
	*MockKnowledgeGraphRepository
}

func NewMockCyclicKnowledgeGraph() *MockCyclicKnowledgeGraph {
	return &MockCyclicKnowledgeGraph{
		MockKnowledgeGraphRepository: NewMockKnowledgeGraphRepository(),
	}
}

// Mock empty knowledge graph
type MockEmptyKnowledgeGraph struct {
	*MockKnowledgeGraphRepository
}

func NewMockEmptyKnowledgeGraph() *MockEmptyKnowledgeGraph {
	return &MockEmptyKnowledgeGraph{
		MockKnowledgeGraphRepository: &MockKnowledgeGraphRepository{nodes: make(map[string]*warehouse.KnowledgeGraphNode)},
	}
}

// Mock evidence cache
type MockEvidenceCache struct {
	entries map[string]*council.EvidenceCacheEntry
	expired map[string]bool
}

func NewMockEvidenceCache() *MockEvidenceCache {
	return &MockEvidenceCache{
		entries: make(map[string]*council.EvidenceCacheEntry),
		expired: make(map[string]bool),
	}
}

func (m *MockEvidenceCache) Set(ctx context.Context, key string, entry *council.EvidenceCacheEntry) error {
	m.entries[key] = entry
	return nil
}

func (m *MockEvidenceCache) Get(ctx context.Context, key string) (*council.EvidenceCacheEntry, error) {
	if m.expired[key] {
		return nil, nil
	}
	return m.entries[key], nil
}

func (m *MockEvidenceCache) ExpireAll() {
	for key := range m.entries {
		m.expired[key] = true
	}
}

func createMockEmbedding(value float32) pgvector.Vector {
	vec := make([]float32, 1536)
	for i := range vec {
		vec[i] = value
	}
	return pgvector.NewVector(vec)
}
