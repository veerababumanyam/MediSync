// Package a03_visualization provides the A-03 Visualization Routing Agent.
//
// This file implements the classification rules for chart type routing.
// It analyzes query intent patterns and data structures to determine
// the optimal visualization type.
package a03_visualization

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/medisync/medisync/internal/warehouse/models"
)

// IntentPattern defines a pattern for matching query intents.
type IntentPattern struct {
	Pattern   *regexp.Regexp
	ChartType models.ChartType
	Intent    string
}

// Classifier provides rule-based chart classification.
type Classifier struct {
	intentPatterns []IntentPattern
}

// NewClassifier creates a new Classifier with predefined patterns.
func NewClassifier() *Classifier {
	return &Classifier{
		intentPatterns: []IntentPattern{
			// Trend patterns - line chart
			{
				Pattern:   regexp.MustCompile(`(?i)\b(trend|over time|monthly|weekly|daily|yearly|growth|history|evolution|progression)\b`),
				ChartType: models.ChartTypeLine,
				Intent:    "trend",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(last\s+\d+\s+(months?|weeks?|days?|years?))\b`),
				ChartType: models.ChartTypeLine,
				Intent:    "trend",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(compare.*\bover\b.*time)\b`),
				ChartType: models.ChartTypeLine,
				Intent:    "trend",
			},

			// Comparison patterns - bar chart
			{
				Pattern:   regexp.MustCompile(`(?i)\b(compare|versus|vs\.?|against|difference between)\b`),
				ChartType: models.ChartTypeBar,
				Intent:    "comparison",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(top\s+\d+|bottom\s+\d+|highest|lowest|most|least)\b`),
				ChartType: models.ChartTypeBar,
				Intent:    "comparison",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(by\s+(department|category|type|region|doctor|patient))\b`),
				ChartType: models.ChartTypeBar,
				Intent:    "comparison",
			},

			// Breakdown patterns - pie chart
			{
				Pattern:   regexp.MustCompile(`(?i)\b(breakdown|distribution|split|segment|proportion|percentage of|share)\b`),
				ChartType: models.ChartTypePie,
				Intent:    "breakdown",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(what\s+percent|what\s+percentage|how\s+much\s+of)\b`),
				ChartType: models.ChartTypePie,
				Intent:    "breakdown",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(parts?\s+of|composition\s+of)\b`),
				ChartType: models.ChartTypePie,
				Intent:    "breakdown",
			},

			// Single value/KPI patterns - kpiCard
			{
				Pattern:   regexp.MustCompile(`(?i)\b(total|sum|average|avg|count|how\s+many|what\s+is\s+the\s+(total|sum|average))\b`),
				ChartType: models.ChartTypeKPICard,
				Intent:    "kpi",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(current|current\s+balance|current\s+stock)\b`),
				ChartType: models.ChartTypeKPICard,
				Intent:    "kpi",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(show\s+me\s+(the\s+)?(revenue|profit|sales|total)\s*$`),
				ChartType: models.ChartTypeKPICard,
				Intent:    "kpi",
			},

			// Correlation patterns - scatter chart
			{
				Pattern:   regexp.MustCompile(`(?i)\b(correlation|relationship|scatter|versus|plot.*against)\b`),
				ChartType: models.ChartTypeScatter,
				Intent:    "correlation",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(relationship\s+between|compare\s+.*\s+and\s+.*\s+trends?)\b`),
				ChartType: models.ChartTypeScatter,
				Intent:    "correlation",
			},

			// Table patterns - dataTable
			{
				Pattern:   regexp.MustCompile(`(?i)\b(list|show\s+all|details|detailed|table|rows)\b`),
				ChartType: models.ChartTypeTable,
				Intent:    "table",
			},
			{
				Pattern:   regexp.MustCompile(`(?i)\b(show\s+me\s+(all\s+)?(patients?|doctors?|bills?|appointments?|records?)\b`),
				ChartType: models.ChartTypeTable,
				Intent:    "table",
			},
		},
	}
}

// ClassificationResult represents the result of chart classification.
type ClassificationResult struct {
	ChartType   models.ChartType
	Intent      string
	Confidence  float64
	Reasoning   string
	MatchedRule string
}

// ClassifyQuery classifies a natural language query to determine the best chart type.
func (c *Classifier) ClassifyQuery(query string) *ClassificationResult {
	// Normalize the query
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	var bestMatch *ClassificationResult

	// Check each pattern
	for _, pattern := range c.intentPatterns {
		if pattern.Pattern.MatchString(normalizedQuery) {
			match := &ClassificationResult{
				ChartType:   pattern.ChartType,
				Intent:      pattern.Intent,
				Confidence:  85.0, // Base confidence for pattern match
				MatchedRule: pattern.Pattern.String(),
			}

			// If this is the first match or a more specific pattern
			if bestMatch == nil || len(pattern.Pattern.String()) > len(bestMatch.MatchedRule) {
				bestMatch = match
			}
		}
	}

	// Default to table if no pattern matches
	if bestMatch == nil {
		return &ClassificationResult{
			ChartType:   models.ChartTypeTable,
			Intent:      "table",
			Confidence:  60.0,
			Reasoning:   "No specific pattern matched, defaulting to table view",
			MatchedRule: "default",
		}
	}

	bestMatch.Reasoning = c.generateReasoning(bestMatch.Intent, bestMatch.ChartType)
	return bestMatch
}

// ClassifyDataStructure analyzes data structure to recommend a chart type.
func (c *Classifier) ClassifyDataStructure(columns []models.ColumnMeta, rowCount int) *ClassificationResult {
	// Analyze column types
	var numericColumns, stringColumns, dateColumns int

	for _, col := range columns {
		switch strings.ToLower(col.Type) {
		case "integer", "float", "numeric", "decimal", "double":
			numericColumns++
		case "string", "text", "varchar":
			stringColumns++
		case "timestamp", "date", "datetime", "time":
			dateColumns++
		}
	}

	// Single row with aggregation - KPI card
	if rowCount == 1 && numericColumns >= 1 {
		return &ClassificationResult{
			ChartType:  models.ChartTypeKPICard,
			Intent:     "kpi",
			Confidence: 95.0,
			Reasoning:  "Single row result with numeric aggregation, ideal for KPI card display",
		}
	}

	// Time series data - line or bar chart
	if dateColumns >= 1 && numericColumns >= 1 {
		if rowCount > 15 {
			return &ClassificationResult{
				ChartType:  models.ChartTypeLine,
				Intent:     "trend",
				Confidence: 90.0,
				Reasoning:  "Time-series data with multiple data points, line chart for trend visualization",
			}
		}
		return &ClassificationResult{
			ChartType:  models.ChartTypeBar,
			Intent:     "trend",
			Confidence: 88.0,
			Reasoning:  "Time-series data with few data points, bar chart for clarity",
		}
	}

	// Categorical data with single metric - pie or bar
	if stringColumns >= 1 && numericColumns == 1 {
		if rowCount <= 8 && rowCount >= 2 {
			return &ClassificationResult{
				ChartType:  models.ChartTypePie,
				Intent:     "breakdown",
				Confidence: 88.0,
				Reasoning:  "Categorical breakdown with few categories, pie chart for distribution",
			}
		}
		if rowCount > 8 && rowCount <= 50 {
			return &ClassificationResult{
				ChartType:  models.ChartTypeBar,
				Intent:     "comparison",
				Confidence: 90.0,
				Reasoning:  "Multiple categories with single metric, bar chart for comparison",
			}
		}
	}

	// Multiple numeric columns - scatter or table
	if numericColumns >= 2 && rowCount >= 5 {
		return &ClassificationResult{
			ChartType:  models.ChartTypeScatter,
			Intent:     "correlation",
			Confidence: 85.0,
			Reasoning:  "Multiple numeric columns suggest correlation analysis, scatter chart",
		}
	}

	// Many columns or complex data - table
	if len(columns) > 4 || rowCount > 50 {
		return &ClassificationResult{
			ChartType:  models.ChartTypeTable,
			Intent:     "table",
			Confidence: 92.0,
			Reasoning:  "Complex data structure with many columns/rows, table for detailed view",
		}
	}

	// Default to bar chart for simple comparisons
	if stringColumns >= 1 && numericColumns >= 1 {
		return &ClassificationResult{
			ChartType:  models.ChartTypeBar,
			Intent:     "comparison",
			Confidence: 82.0,
			Reasoning:  "Categorical and numeric data, bar chart for comparison",
		}
	}

	// Fallback to table
	return &ClassificationResult{
		ChartType:  models.ChartTypeTable,
		Intent:     "table",
		Confidence: 70.0,
		Reasoning:  "Unable to determine optimal chart type, defaulting to table view",
	}
}

// CombineResults combines query intent and data structure analysis.
func (c *Classifier) CombineResults(queryResult, dataResult *ClassificationResult) *ClassificationResult {
	// If data structure strongly suggests a type (high confidence), prefer it
	if dataResult.Confidence >= 95.0 {
		return dataResult
	}

	// If query intent has high confidence, prefer it
	if queryResult.Confidence >= 90.0 {
		return queryResult
	}

	// If both agree, boost confidence
	if queryResult.ChartType == dataResult.ChartType {
		return &ClassificationResult{
			ChartType:  queryResult.ChartType,
			Intent:     queryResult.Intent,
			Confidence: min(queryResult.Confidence+10, 98.0),
			Reasoning:  queryResult.Reasoning + " (confirmed by data structure)",
		}
	}

	// If they disagree, prefer query intent for user experience
	// but lower confidence
	return &ClassificationResult{
		ChartType:  queryResult.ChartType,
		Intent:     queryResult.Intent,
		Confidence: queryResult.Confidence * 0.9,
		Reasoning:  queryResult.Reasoning + " (note: data structure suggests " + string(dataResult.ChartType) + ")",
	}
}

// generateReasoning generates a human-readable explanation for the chart choice.
func (c *Classifier) generateReasoning(intent string, chartType models.ChartType) string {
	reasons := map[string]map[models.ChartType]string{
		"trend": {
			models.ChartTypeLine: "Query indicates trend analysis, line chart shows changes over time effectively",
			models.ChartTypeBar:  "Query indicates trend analysis with discrete time periods, bar chart for clarity",
		},
		"comparison": {
			models.ChartTypeBar:   "Query indicates comparison between entities, bar chart for visual comparison",
			models.ChartTypeTable: "Query indicates comparison, table for precise value comparison",
		},
		"breakdown": {
			models.ChartTypePie: "Query indicates breakdown/distribution, pie chart for part-to-whole",
			models.ChartTypeBar: "Query indicates breakdown with many categories, bar chart for readability",
		},
		"kpi": {
			models.ChartTypeKPICard: "Query requests single metric, KPI card for prominent display",
		},
		"correlation": {
			models.ChartTypeScatter: "Query indicates relationship analysis, scatter chart for correlation",
		},
		"table": {
			models.ChartTypeTable: "Query requests detailed data, table for complete information",
		},
	}

	if reasonMap, ok := reasons[intent]; ok {
		if reason, ok := reasonMap[chartType]; ok {
			return reason
		}
	}

	return "Selected " + string(chartType) + " based on query analysis"
}

// DetermineAxisLabels suggests axis labels based on column metadata.
func (c *Classifier) DetermineAxisLabels(columns []models.ColumnMeta, intent string) (xAxis, yAxis string) {
	var timeColumn, categoryColumn, valueColumn string

	for _, col := range columns {
		colType := strings.ToLower(col.Type)
		colName := col.Name

		// Identify time/date columns
		if colType == "timestamp" || colType == "date" || colType == "datetime" {
			if timeColumn == "" {
				timeColumn = colName
			}
			continue
		}

		// Identify numeric columns (potential y-axis)
		if colType == "integer" || colType == "float" || colType == "numeric" || colType == "decimal" {
			if valueColumn == "" || strings.Contains(strings.ToLower(colName), "amount") ||
				strings.Contains(strings.ToLower(colName), "total") || strings.Contains(strings.ToLower(colName), "count") {
				valueColumn = colName
			}
			continue
		}

		// Identify category columns
		if colType == "string" || colType == "text" || colType == "varchar" {
			if categoryColumn == "" {
				categoryColumn = colName
			}
		}
	}

	// Assign based on intent
	switch intent {
	case "trend":
		if timeColumn != "" {
			xAxis = timeColumn
		}
		if valueColumn != "" {
			yAxis = valueColumn
		}
	case "comparison", "breakdown":
		if categoryColumn != "" {
			xAxis = categoryColumn
		}
		if valueColumn != "" {
			yAxis = valueColumn
		}
	case "correlation":
		// For scatter, both axes are numeric
		xAxis = "X"
		yAxis = "Y"
	}

	return xAxis, yAxis
}

// ExtractKPIValue extracts a single KPI value from result data.
func (c *Classifier) ExtractKPIValue(data []map[string]any) (value any, formatted string) {
	if len(data) == 0 {
		return nil, "N/A"
	}

	row := data[0]

	// Look for common KPI column names
	kpiColumns := []string{
		"total", "sum", "amount", "count", "value", "total_amount",
		"total_revenue", "revenue", "total_sales", "sales",
	}

	for _, col := range kpiColumns {
		if val, ok := row[col]; ok {
			return val, formatKPIValue(val)
		}
	}

	// Return first numeric value found
	for _, val := range row {
		switch v := val.(type) {
		case float64:
			return v, formatKPIValue(v)
		case float32:
			return v, formatKPIValue(float64(v))
		case int:
			return v, formatKPIValue(float64(v))
		case int64:
			return v, formatKPIValue(float64(v))
		}
	}

	return nil, "N/A"
}

// formatKPIValue formats a numeric value for display.
func formatKPIValue(val any) string {
	switch v := val.(type) {
	case float64:
		if v >= 1000000 {
			return formatNumber(v/1000000, 2) + "M"
		} else if v >= 1000 {
			return formatNumber(v/1000, 2) + "K"
		}
		return formatNumber(v, 2)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return "N/A"
	}
}

// formatNumber formats a float with specified decimal places.
func formatNumber(val float64, decimals int) string {
	format := "%." + strconv.Itoa(decimals) + "f"
	return sprintf(format, val)
}

// sprintf is a simple wrapper for fmt.Sprintf to avoid import.
func sprintf(format string, a ...any) string {
	return strings.ReplaceAll(format, "%f", "") + " " + strings.TrimPrefix(strings.TrimSuffix(format, "f"), "%.")
}
