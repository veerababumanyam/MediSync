// Package council provides semantic equivalence detection for agent responses.
//
// The semantic module implements equivalence detection using pgvector cosine
// similarity to group responses that express the same information differently.
//
// Key Features:
//   - Cosine similarity calculation using pgvector
//   - 95% threshold for equivalence grouping
//   - Canonical response selection based on confidence
//   - Support for multilingual equivalence detection
package council

import (
	"fmt"
	"sort"
)

// SemanticDetector provides semantic equivalence detection for agent responses.
type SemanticDetector struct {
	threshold float64 // Similarity threshold (default 0.95)
}

// NewSemanticDetector creates a new semantic detector with the given threshold.
func NewSemanticDetector(threshold float64) *SemanticDetector {
	if threshold <= 0 || threshold > 1 {
		threshold = DefaultSemanticThreshold
	}
	return &SemanticDetector{threshold: threshold}
}

// CosineSimilarity calculates the cosine similarity between two embedding vectors.
func (d *SemanticDetector) CosineSimilarity(e1, e2 Embedding) (float64, error) {
	if len(e1) == 0 || len(e2) == 0 {
		return 0, fmt.Errorf("empty embedding")
	}
	if len(e1) != len(e2) {
		return 0, fmt.Errorf("embedding dimensions don't match: %d vs %d", len(e1), len(e2))
	}

	var dotProduct, norm1, norm2 float64
	for i := range e1 {
		dotProduct += float64(e1[i]) * float64(e2[i])
		norm1 += float64(e1[i]) * float64(e1[i])
		norm2 += float64(e2[i]) * float64(e2[i])
	}

	if norm1 == 0 || norm2 == 0 {
		return 0, nil
	}

	similarity := dotProduct / (sqrt(norm1) * sqrt(norm2))
	return similarity, nil
}

// IsEquivalent checks if two responses are semantically equivalent.
func (d *SemanticDetector) IsEquivalent(e1, e2 Embedding) bool {
	similarity, err := d.CosineSimilarity(e1, e2)
	if err != nil {
		return false
	}
	return similarity >= d.threshold
}

// CalculateSimilarity calculates similarity between two responses.
func (d *SemanticDetector) CalculateSimilarity(r1, r2 *AgentResponse) (float64, error) {
	if len(r1.Embedding.Slice()) == 0 || len(r2.Embedding.Slice()) == 0 {
		return 0, fmt.Errorf("response missing embedding")
	}
	return d.CosineSimilarity(r1.Embedding.Slice(), r2.Embedding.Slice())
}

// GroupEquivalentResponses groups responses that are semantically equivalent.
func (d *SemanticDetector) GroupEquivalentResponses(responses []*AgentResponse) []Group {
	if len(responses) == 0 {
		return nil
	}

	// Build similarity matrix
	n := len(responses)
	similarity := make([][]float64, n)
	for i := range similarity {
		similarity[i] = make([]float64, n)
		for j := range similarity[i] {
			if i == j {
				similarity[i][j] = 1.0
			} else {
				sim, _ := d.CalculateSimilarity(responses[i], responses[j])
				similarity[i][j] = sim
			}
		}
	}

	// Union-Find for grouping
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) {
		px, py := find(x), find(y)
		if px != py {
			parent[px] = py
		}
	}

	// Group equivalent responses
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if similarity[i][j] >= d.threshold {
				union(i, j)
			}
		}
	}

	// Collect groups
	groupMap := make(map[int][]int)
	for i := 0; i < n; i++ {
		root := find(i)
		groupMap[root] = append(groupMap[root], i)
	}

	// Convert to Group structs
	groups := make([]Group, 0, len(groupMap))
	for _, indices := range groupMap {
		group := Group{
			GroupID:  len(groups) + 1,
			AgentIDs: make([]string, 0, len(indices)),
		}

		// Find canonical (highest confidence) response
		var maxConfidence float64
		var canonicalIdx int
		for _, idx := range indices {
			resp := responses[idx]
			group.AgentIDs = append(group.AgentIDs, resp.AgentID)
			if resp.Confidence > maxConfidence {
				maxConfidence = resp.Confidence
				canonicalIdx = idx
			}
		}

		group.Canonical = responses[canonicalIdx].ResponseText
		group.Similarity = calculateAvgSimilarity(indices, similarity)

		groups = append(groups, group)
	}

	// Sort groups by size (largest first)
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i].AgentIDs) > len(groups[j].AgentIDs)
	})

	return groups
}

// calculateAvgSimilarity calculates the average similarity within a group.
func calculateAvgSimilarity(indices []int, similarity [][]float64) float64 {
	if len(indices) < 2 {
		return 1.0
	}

	var sum float64
	var count int
	for i := 0; i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			sum += similarity[indices[i]][indices[j]]
			count++
		}
	}

	if count == 0 {
		return 1.0
	}
	return sum / float64(count)
}

// sqrt computes the square root of a float64.
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}

	// Newton's method for square root
	z := x
	for i := 0; i < 100; i++ {
		z = z - (z*z-x)/(2*z)
		if z*z-x < 1e-10 && z*z-x > -1e-10 {
			break
		}
	}
	return z
}
