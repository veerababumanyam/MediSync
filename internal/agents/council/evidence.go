// Package council provides evidence retrieval from the Knowledge Graph.
//
// The evidence module implements Graph-of-Thoughts retrieval with multi-hop
// traversal through the medical knowledge graph.
//
// Key Features:
//   - Multi-hop graph traversal (default 3 hops)
//   - Relevance score calculation
//   - Cycle detection and prevention
//   - Evidence sufficiency assessment
package council

import (
	"context"
	"fmt"
	"time"

	"github.com/medisync/medisync/internal/warehouse"
	"github.com/pgvector/pgvector-go"
)

// EvidenceStatus represents the sufficiency of retrieved evidence.
type EvidenceStatus string

const (
	EvidenceStatusSufficient   EvidenceStatus = "sufficient"
	EvidenceStatusInsufficient EvidenceStatus = "insufficient"
	EvidenceStatusNone         EvidenceStatus = "none"
)

// EvidenceRetriever retrieves evidence from the Knowledge Graph.
type EvidenceRetriever struct {
	kgRepo  KnowledgeGraphRepository
	maxHops int
}

// KnowledgeGraphRepository defines the interface for KG operations.
type KnowledgeGraphRepository interface {
	GetNode(ctx context.Context, id string) (*warehouse.KnowledgeGraphNode, error)
	GetNodes(ctx context.Context, ids []string) ([]*warehouse.KnowledgeGraphNode, error)
	GetRelatedNodes(ctx context.Context, nodeID string, edgeTypes []warehouse.KGEdgeType, limit int) ([]*warehouse.KnowledgeGraphNode, error)
	TraverseMultiHop(ctx context.Context, initialNodeIDs []string, maxHops int) (*warehouse.TraversalResult, error)
	FindSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]*warehouse.KnowledgeGraphNode, error)
	HealthCheck(ctx context.Context) error
}

// TraversalResult represents the result of a multi-hop traversal.
type TraversalResult struct {
	Nodes           []*warehouse.KnowledgeGraphNode `json:"nodes"`
	TraversalPath   []warehouse.TraversalStep       `json:"traversal_path"`
	RelevanceScores map[string]float64              `json:"relevance_scores"`
}

// EvidenceCacheEntry represents cached evidence data.
type EvidenceCacheEntry struct {
	QueryHash       string                          `json:"query_hash"`
	QueryText       string                          `json:"query_text,omitempty"`
	NodeIDs         []string                        `json:"node_ids"`
	Nodes           []*warehouse.KnowledgeGraphNode `json:"nodes"`
	TraversalPath   []*warehouse.TraversalStep      `json:"traversal_path"`
	RelevanceScores map[string]float64              `json:"relevance_scores"`
	HopCount        int                             `json:"hop_count"`
	CachedAt        time.Time                       `json:"cached_at"`
	ExpiresAt       time.Time                       `json:"expires_at"`
	Source          string                          `json:"source"`
	KGHealthStatus  string                          `json:"kg_health_status"`
}

// IsTrustworthyDuringOutage checks if cached evidence is trustworthy during a KG outage.
func (e *EvidenceCacheEntry) IsTrustworthyDuringOutage() bool {
	return e.KGHealthStatus == "healthy" && !e.IsExpired()
}

// IsExpired checks if the entry has expired.
func (e *EvidenceCacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// ToEvidenceTrail converts the cache entry to an EvidenceTrail.
func (e *EvidenceCacheEntry) ToEvidenceTrail(deliberationID string) *EvidenceTrail {
	return &EvidenceTrail{
		ID:             "",
		DeliberationID: deliberationID,
		NodeIDs:        e.NodeIDs,
		TraversalPath:  e.TraversalPath,
		RelevanceScore: e.RelevanceScores,
		HopCount:       e.HopCount,
		CachedAt:       e.CachedAt,
		ExpiresAt:      e.ExpiresAt,
	}
}

// EvidenceCache defines the interface for caching evidence.
type EvidenceCache interface {
	Get(ctx context.Context, queryHash string) (*EvidenceCacheEntry, error)
	Set(ctx context.Context, queryHash string, entry *EvidenceCacheEntry) error
}

// NewEvidenceRetriever creates a new evidence retriever.
func NewEvidenceRetriever(kgRepo KnowledgeGraphRepository, maxHops int) *EvidenceRetriever {
	if maxHops <= 0 {
		maxHops = DefaultMaxHops
	}
	return &EvidenceRetriever{
		kgRepo:  kgRepo,
		maxHops: maxHops,
	}
}

// Retrieve retrieves evidence starting from initial node IDs.
func (r *EvidenceRetriever) Retrieve(ctx context.Context, initialNodeIDs []string, maxHops int) (*EvidenceTrail, error) {
	if maxHops <= 0 {
		maxHops = r.maxHops
	}

	result, err := r.kgRepo.TraverseMultiHop(ctx, initialNodeIDs, maxHops)
	if err != nil {
		return nil, fmt.Errorf("traversal failed: %w", err)
	}

	nodeIDs := make([]string, 0, len(result.Nodes))
	for _, node := range result.Nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}

	return &EvidenceTrail{
		ID:             "",
		DeliberationID: "",
		NodeIDs:        nodeIDs,
		TraversalPath:  result.TraversalPath,
		RelevanceScore: result.RelevanceScore,
		HopCount:       maxHops,
		CachedAt:       time.Now(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}, nil
}

// RetrieveWithFilters retrieves evidence with edge type filtering.
func (r *EvidenceRetriever) RetrieveWithFilters(ctx context.Context, initialNodeIDs []string, maxHops int, edgeTypes []warehouse.KGEdgeType) (*EvidenceTrail, error) {
	// For filtered retrieval, we do a custom traversal
	if maxHops <= 0 {
		maxHops = r.maxHops
	}

	visited := make(map[string]bool)
	path := make([]*warehouse.TraversalStep, 0)
	relevance := make(map[string]float64)
	var allNodes []string

	queue := make([]struct {
		id     string
		hop    int
		from   string
		edge   warehouse.KGEdgeType
		weight float64
	}, 0)

	// Initialize with starting nodes
	for _, id := range initialNodeIDs {
		queue = append(queue, struct {
			id     string
			hop    int
			from   string
			edge   warehouse.KGEdgeType
			weight float64
		}{id: id, hop: 0, weight: 1.0})
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.id] || current.hop > maxHops {
			continue
		}
		visited[current.id] = true
		allNodes = append(allNodes, current.id)

		// Calculate relevance (decreases with hop count)
		relevance[current.id] = 1.0 - (float64(current.hop) * 0.2)

		if current.hop > 0 {
			path = append(path, &warehouse.TraversalStep{
				FromNodeID: current.from,
				ToNodeID:   current.id,
				EdgeType:   current.edge,
				Weight:     current.weight,
			})
		}

		// Get related nodes with filtering
		related, err := r.kgRepo.GetRelatedNodes(ctx, current.id, edgeTypes, 10)
		if err != nil {
			continue
		}

		for _, node := range related {
			if !visited[node.ID] && current.hop < maxHops {
				// Determine edge type (simplified)
				edgeType := warehouse.EdgeRelatedTo
				if len(edgeTypes) > 0 {
					edgeType = edgeTypes[0]
				}
				queue = append(queue, struct {
					id     string
					hop    int
					from   string
					edge   warehouse.KGEdgeType
					weight float64
				}{
					id:     node.ID,
					hop:    current.hop + 1,
					from:   current.id,
					edge:   edgeType,
					weight: 0.8,
				})
			}
		}
	}

	return &EvidenceTrail{
		NodeIDs:        allNodes,
		TraversalPath:  path,
		RelevanceScore: relevance,
		HopCount:       maxHops,
		CachedAt:       time.Now(),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
	}, nil
}

// AssessEvidenceSufficiency determines if evidence is sufficient for consensus.
func AssessEvidenceSufficiency(trail *EvidenceTrail, minNodes int, minRelevance float64) EvidenceStatus {
	if trail == nil || len(trail.NodeIDs) == 0 {
		return EvidenceStatusNone
	}

	if len(trail.NodeIDs) < minNodes {
		return EvidenceStatusInsufficient
	}

	// Check average relevance
	var totalRelevance float64
	for _, score := range trail.RelevanceScore {
		totalRelevance += score
	}

	avgRelevance := totalRelevance / float64(len(trail.RelevanceScore))
	if avgRelevance < minRelevance {
		return EvidenceStatusInsufficient
	}

	return EvidenceStatusSufficient
}
