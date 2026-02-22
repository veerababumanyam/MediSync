// Package council provides consensus calculation for multi-agent deliberations.
//
// The consensus module implements weighted voting and agreement calculation to
// determine if agents have reached sufficient consensus on a response.
//
// Key Features:
//   - Weighted voting based on agent confidence
//   - 80% default threshold for consensus
//   - Semantic equivalence grouping integration
//   - Dissenting agent tracking
package council

import (
	"fmt"
)

// ConsensusCalculator calculates consensus from agent responses.
type ConsensusCalculator struct {
	threshold float64 // Agreement threshold (default 0.80)
	detector  *SemanticDetector
}

// NewConsensusCalculator creates a new consensus calculator.
func NewConsensusCalculator(threshold float64) *ConsensusCalculator {
	if threshold <= 0 || threshold > 1 {
		threshold = DefaultConsensusThreshold
	}
	return &ConsensusCalculator{
		threshold: threshold,
		detector:  NewSemanticDetector(DefaultSemanticThreshold),
	}
}

// ConsensusResult represents the outcome of consensus calculation.
type ConsensusResult struct {
	AgreementScore    float64            `json:"agreement_score"`
	ThresholdMet      bool               `json:"threshold_met"`
	Status            DeliberationStatus `json:"status"`
	FinalResponse     string             `json:"final_response"`
	ConfidenceScore   float64            `json:"confidence_score"`
	EquivalenceGroups []Group            `json:"equivalence_groups"`
	DissentingAgents  []string           `json:"dissenting_agents"`
	Method            string             `json:"method"`
}

// Calculate calculates consensus from agent responses.
func (c *ConsensusCalculator) Calculate(responses []*AgentResponse) (*ConsensusResult, error) {
	if len(responses) < DefaultMinAgents {
		return nil, fmt.Errorf("minimum %d agents required, got %d", DefaultMinAgents, len(responses))
	}

	// Group equivalent responses
	groups := c.detector.GroupEquivalentResponses(responses)

	// Calculate weighted agreement
	result := &ConsensusResult{
		EquivalenceGroups: groups,
		Method:            "weighted_vote",
	}

	// Find the largest group (main consensus)
	if len(groups) == 0 {
		result.Status = StatusUncertain
		result.AgreementScore = 0
		return result, nil
	}

	largestGroup := groups[0]

	// Calculate agreement score (weighted by confidence)
	totalWeight := 0.0
	consensusWeight := 0.0
	var dissentingAgents []string

	for _, resp := range responses {
		weight := resp.Confidence / 100.0
		totalWeight += weight

		// Check if this agent is in the largest group
		isInLargestGroup := false
		for _, agentID := range largestGroup.AgentIDs {
			if agentID == resp.AgentID {
				isInLargestGroup = true
				break
			}
		}

		if isInLargestGroup {
			consensusWeight += weight
		} else {
			dissentingAgents = append(dissentingAgents, resp.AgentID)
		}
	}

	// Calculate agreement percentage
	if totalWeight > 0 {
		result.AgreementScore = (consensusWeight / totalWeight) * 100
	}

	result.DissentingAgents = dissentingAgents

	// Determine if threshold is met
	result.ThresholdMet = result.AgreementScore >= c.threshold*100

	if result.ThresholdMet {
		result.Status = StatusConsensus
		result.FinalResponse = largestGroup.Canonical
		result.ConfidenceScore = calculateGroupConfidence(largestGroup, responses)
	} else {
		result.Status = StatusUncertain
		result.ConfidenceScore = 0
	}

	return result, nil
}

// calculateGroupConfidence calculates the average confidence for agents in a group.
func calculateGroupConfidence(group Group, responses []*AgentResponse) float64 {
	var total float64
	var count int

	for _, resp := range responses {
		for _, agentID := range group.AgentIDs {
			if resp.AgentID == agentID {
				total += resp.Confidence
				count++
				break
			}
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// HashQuery generates a SHA-256 hash of a query for deduplication.
func HashQuery(query string) string {
	return hashQuery(query)
}
