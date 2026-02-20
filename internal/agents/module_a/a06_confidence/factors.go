// Package a06_confidence provides the confidence scoring agent.
//
// This file implements factor calculation for confidence scoring.
package a06_confidence

import (
	"regexp"
	"strings"
	"time"

	"github.com/medisync/medisync/internal/warehouse/models"
)

// FactorCalculator calculates individual confidence factors.
type FactorCalculator struct {
	complexityWeights ComplexityWeights
}

// ComplexityWeights holds weights for SQL complexity calculation.
type ComplexityWeights struct {
	JoinWeight      float64
	SubqueryWeight  float64
	AggregateWeight float64
	WindowWeight    float64
	CTEWeight       float64
}

// NewFactorCalculator creates a new factor calculator.
func NewFactorCalculator() *FactorCalculator {
	return &FactorCalculator{
		complexityWeights: ComplexityWeights{
			JoinWeight:      0.05,
			SubqueryWeight:  0.08,
			AggregateWeight: 0.02,
			WindowWeight:    0.10,
			CTEWeight:       0.07,
		},
	}
}

// Calculate calculates all confidence factors for a request.
func (c *FactorCalculator) Calculate(req ScoreRequest) models.ConfidenceFactors {
	return models.ConfidenceFactors{
		IntentClarity:         c.CalculateIntentClarity(req),
		SchemaMatchQuality:    c.CalculateSchemaMatchQuality(req),
		SQLComplexityPenalty:  c.CalculateSQLComplexityPenalty(req.GeneratedSQL),
		RetryPenalty:          c.CalculateRetryPenalty(req.RetryCount),
		HallucinationRisk:     c.CalculateHallucinationRisk(req),
	}
}

// CalculateIntentClarity measures how clear the user's intent was.
func (c *FactorCalculator) CalculateIntentClarity(req ScoreRequest) float64 {
	// Base clarity from intent detection confidence
	clarity := req.IntentConfidence

	// Adjust based on query length (very short queries are less clear)
	queryLen := len(strings.TrimSpace(req.UserQuery))
	if queryLen < 10 {
		clarity *= 0.7
	} else if queryLen < 20 {
		clarity *= 0.9
	}

	// Adjust based on detected intent
	intentBonus := map[string]float64{
		"kpi":        0.1,
		"trend":      0.05,
		"comparison": 0.05,
		"breakdown":  0.05,
		"table":      0.0,
	}
	if bonus, ok := intentBonus[req.DetectedIntent]; ok {
		clarity = minFloat(1.0, clarity+bonus)
	}

	// Penalize if validation failed
	if !req.ValidationPassed {
		clarity *= 0.6
	}

	return minFloat(1.0, maxFloat(0.0, clarity))
}

// CalculateSchemaMatchQuality measures how well the schema matched the query.
func (c *FactorCalculator) CalculateSchemaMatchQuality(req ScoreRequest) float64 {
	if len(req.SchemaMatches) == 0 {
		return 0.3 // Low confidence with no schema context
	}

	// Quality based on number of relevant tables found
	matchCount := len(req.SchemaMatches)

	quality := 0.5 // Base quality
	if matchCount >= 1 {
		quality += 0.2
	}
	if matchCount >= 2 {
		quality += 0.1
	}
	if matchCount >= 3 {
		quality += 0.1
	}
	if matchCount >= 5 {
		quality -= 0.1 // Too many matches might indicate ambiguity
	}
	if matchCount >= 8 {
		quality -= 0.2 // Very ambiguous
	}

	// Bonus for fast execution (indicates good query plan)
	if req.ExecutionTime > 0 && req.ExecutionTime < 1*time.Second {
		quality += 0.05
	}

	// Penalty for very slow execution
	if req.ExecutionTime > 10*time.Second {
		quality -= 0.1
	}

	return minFloat(1.0, maxFloat(0.0, quality))
}

// CalculateSQLComplexityPenalty penalizes complex SQL queries.
func (c *FactorCalculator) CalculateSQLComplexityPenalty(sql string) (penalty float64) {
	sqlUpper := strings.ToUpper(sql)

	// Count JOINs
	joinCount := strings.Count(sqlUpper, " JOIN ")
	penalty += float64(joinCount) * c.complexityWeights.JoinWeight

	// Count subqueries
	subqueryPattern := regexp.MustCompile(`\(\s*SELECT`)
	subqueryCount := len(subqueryPattern.FindAllString(sqlUpper, -1))
	penalty += float64(subqueryCount) * c.complexityWeights.SubqueryWeight

	// Count aggregations
	aggPattern := regexp.MustCompile(`\b(SUM|AVG|COUNT|MIN|MAX|STDDEV|VARIANCE)\s*\(`)
	aggCount := len(aggPattern.FindAllString(sqlUpper, -1))
	penalty += float64(aggCount) * c.complexityWeights.AggregateWeight

	// Check for window functions
	if strings.Contains(sqlUpper, " OVER ") {
		penalty += c.complexityWeights.WindowWeight
	}

	// Check for CTEs
	if strings.Contains(sqlUpper, " WITH ") {
		penalty += c.complexityWeights.CTEWeight
	}

	// Cap the penalty
	return minFloat(0.3, penalty)
}

// CalculateRetryPenalty penalizes queries that needed self-correction.
func (c *FactorCalculator) CalculateRetryPenalty(retryCount int) float64 {
	// 10% penalty per retry, max 30%
	return minFloat(0.3, float64(retryCount)*0.1)
}

// CalculateHallucinationRisk measures the risk of hallucinated data.
func (c *FactorCalculator) CalculateHallucinationRisk(req ScoreRequest) float64 {
	risk := 0.0

	// Higher risk for queries with no results
	if req.RowCount == 0 {
		risk += 0.2
	}

	// Higher risk for very high result counts (might be wrong join)
	if req.RowCount > 100000 {
		risk += 0.15
	}

	// Higher risk for complex queries that had retries
	if req.RetryCount > 0 && c.CalculateSQLComplexityPenalty(req.GeneratedSQL) > 0.15 {
		risk += 0.1
	}

	// Higher risk if intent detection was uncertain
	if req.IntentConfidence < 0.7 {
		risk += 0.15
	}

	// Higher risk for queries with many schema matches (ambiguity)
	if len(req.SchemaMatches) > 6 {
		risk += 0.1
	}

	// Lower risk if validation passed
	if req.ValidationPassed {
		risk -= 0.1
	}

	return minFloat(1.0, maxFloat(0.0, risk))
}

// minFloat returns the smaller of two floats.
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// maxFloat returns the larger of two floats.
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
