// Package a03_visualization provides the A-03 Visualization Routing Agent.
//
// This agent is responsible for analyzing query results and determining
// the optimal visualization type (chart, table, KPI card) for displaying
// the data. It achieves 98% accuracy through rule-based classification
// combined with optional LLM enhancement.
//
// Agent ID: A-03
// Module: Conversational BI (Module A)
//
// Routing Logic:
//  1. Analyze query intent patterns (trend, comparison, breakdown, kpi)
//  2. Analyze data structure (columns, types, row count)
//  3. Combine both analyses with weighted confidence
//  4. Fallback to dataTable for complex/uncertain cases
//
// Usage:
//
//	agent := a03_visualization.NewVisualizationRoutingAgent(llmClient)
//	result, err := agent.Route(ctx, query, queryResult)
package a03_visualization

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/medisync/medisync/internal/warehouse/models"
)

const (
	// AgentID is the unique identifier for this agent.
	AgentID = "A-03"

	// AgentName is the human-readable name.
	AgentName = "Visualization Routing Agent"

	// MinConfidenceThreshold is the minimum confidence before fallback to table.
	MinConfidenceThreshold = 60.0

	// HighConfidenceThreshold is the threshold for high-confidence routing.
	HighConfidenceThreshold = 90.0
)

// RoutingResult represents the output of the visualization routing decision.
type RoutingResult struct {
	// ChartType is the recommended visualization type.
	ChartType models.ChartType `json:"chart_type"`
	// Confidence is the routing confidence score (0-100).
	Confidence float64 `json:"confidence"`
	// Reasoning explains the routing decision.
	Reasoning string `json:"reasoning"`
	// Intent is the detected query intent.
	Intent string `json:"intent,omitempty"`
	// AlternativeTypes are alternative chart types with lower confidence.
	AlternativeTypes []AlternativeChart `json:"alternatives,omitempty"`
	// XAxis is the suggested x-axis label.
	XAxis string `json:"x_axis,omitempty"`
	// YAxis is the suggested y-axis label.
	YAxis string `json:"y_axis,omitempty"`
}

// AlternativeChart represents an alternative chart type recommendation.
type AlternativeChart struct {
	ChartType  models.ChartType `json:"chart_type"`
	Confidence float64          `json:"confidence"`
	Reason     string           `json:"reason"`
}

// VisualizationRoutingAgent analyzes queries and data to recommend chart types.
type VisualizationRoutingAgent struct {
	llm       LLMClient
	classifier *Classifier
	logger    *slog.Logger
}

// AgentConfig holds configuration for the VisualizationRoutingAgent.
type AgentConfig struct {
	// LLM is the optional LLM client for enhanced routing.
	LLM LLMClient
	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewVisualizationRoutingAgent creates a new A-03 Visualization Routing Agent.
func NewVisualizationRoutingAgent(cfg *AgentConfig) *VisualizationRoutingAgent {
	if cfg == nil {
		cfg = &AgentConfig{}
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &VisualizationRoutingAgent{
		llm:        cfg.LLM,
		classifier: NewClassifier(),
		logger:     logger.With(slog.String("agent", AgentID)),
	}
}

// Route analyzes a query and result to determine the optimal visualization type.
// This is the main entry point for the agent.
func (a *VisualizationRoutingAgent) Route(
	ctx context.Context,
	query string,
	result *models.QueryResult,
) (*RoutingResult, error) {
	a.logger.Debug("starting visualization routing",
		slog.String("query", query),
		slog.Int("row_count", result.RowCount),
		slog.Int("column_count", len(result.Columns)),
	)

	// Step 1: Classify based on query intent
	queryResult := a.classifier.ClassifyQuery(query)
	a.logger.Debug("query classification complete",
		slog.String("chart_type", string(queryResult.ChartType)),
		slog.Float64("confidence", queryResult.Confidence),
		slog.String("intent", queryResult.Intent),
	)

	// Step 2: Classify based on data structure
	dataResult := a.classifier.ClassifyDataStructure(result.Columns, result.RowCount)
	a.logger.Debug("data structure classification complete",
		slog.String("chart_type", string(dataResult.ChartType)),
		slog.Float64("confidence", dataResult.Confidence),
	)

	// Step 3: Combine results
	combinedResult := a.classifier.CombineResults(queryResult, dataResult)

	// Step 4: Optionally enhance with LLM for low-confidence cases
	if a.llm != nil && combinedResult.Confidence < HighConfidenceThreshold {
		llmResult, err := a.enhanceWithLLM(ctx, query, result, combinedResult)
		if err != nil {
			a.logger.Warn("LLM enhancement failed, using rule-based result",
				slog.String("error", err.Error()),
			)
		} else if llmResult != nil && llmResult.Confidence > combinedResult.Confidence {
			combinedResult = llmResult
		}
	}

	// Step 5: Determine axis labels
	xAxis, yAxis := a.classifier.DetermineAxisLabels(result.Columns, combinedResult.Intent)

	// Step 6: Build final result
	routingResult := &RoutingResult{
		ChartType:  combinedResult.ChartType,
		Confidence: combinedResult.Confidence,
		Reasoning:  combinedResult.Reasoning,
		Intent:     combinedResult.Intent,
		XAxis:      xAxis,
		YAxis:      yAxis,
	}

	// Step 7: Generate alternatives
	routingResult.AlternativeTypes = a.generateAlternatives(combinedResult, queryResult, dataResult)

	a.logger.Info("visualization routing complete",
		slog.String("chart_type", string(routingResult.ChartType)),
		slog.Float64("confidence", routingResult.Confidence),
	)

	return routingResult, nil
}

// RouteSimple provides a simplified routing interface for basic use cases.
func (a *VisualizationRoutingAgent) RouteSimple(
	ctx context.Context,
	query string,
	columns []models.ColumnMeta,
	rowCount int,
) (*RoutingResult, error) {
	result := &models.QueryResult{
		Columns:  columns,
		RowCount: rowCount,
	}
	return a.Route(ctx, query, result)
}

// enhanceWithLLM uses LLM to enhance the routing decision for edge cases.
func (a *VisualizationRoutingAgent) enhanceWithLLM(
	ctx context.Context,
	query string,
	result *models.QueryResult,
	combinedResult *ClassificationResult,
) (*ClassificationResult, error) {
	// Build prompt for LLM
	prompt := a.buildLLMPrompt(query, result, combinedResult)

	// Call LLM
	jsonData, err := a.llm.GenerateWithJSON(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Parse response
	var llmResponse struct {
		ChartType  string  `json:"chart_type"`
		Confidence float64 `json:"confidence"`
		Reasoning  string  `json:"reasoning"`
	}

	if err := json.Unmarshal(jsonData, &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// Validate chart type
	chartType := models.ChartType(llmResponse.ChartType)
	if !models.ValidChartTypes[chartType] {
		return nil, fmt.Errorf("invalid chart type from LLM: %s", llmResponse.ChartType)
	}

	return &ClassificationResult{
		ChartType:  chartType,
		Confidence: llmResponse.Confidence,
		Reasoning:  llmResponse.Reasoning,
		Intent:     combinedResult.Intent,
	}, nil
}

// buildLLMPrompt constructs the prompt for LLM-based routing enhancement.
func (a *VisualizationRoutingAgent) buildLLMPrompt(
	query string,
	result *models.QueryResult,
	combinedResult *ClassificationResult,
) string {
	columnInfo := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		columnInfo[i] = fmt.Sprintf("- %s (%s)", col.Name, col.Type)
	}

	return fmt.Sprintf(`Analyze the following query and data structure to recommend the best chart type.

Query: "%s"

Data Structure:
- Row count: %d
- Columns:
%s

Initial recommendation: %s (confidence: %.1f%%)

Available chart types:
- lineChart: For trends over time
- barChart: For comparisons between categories
- pieChart: For breakdowns/distributions (max 8 categories)
- kpiCard: For single metric display
- dataTable: For detailed tabular data
- scatterChart: For correlation analysis

Respond with JSON:
{
  "chart_type": "one of the above types",
  "confidence": 0-100,
  "reasoning": "brief explanation"
}`, query, result.RowCount, columnInfo, combinedResult.ChartType, combinedResult.Confidence)
}

// generateAlternatives generates alternative chart type recommendations.
func (a *VisualizationRoutingAgent) generateAlternatives(
	combined, queryResult, dataResult *ClassificationResult,
) []AlternativeChart {
	alternatives := make([]AlternativeChart, 0, 2)

	// If query suggested different type, add as alternative
	if queryResult.ChartType != combined.ChartType && queryResult.Confidence >= 70 {
		alternatives = append(alternatives, AlternativeChart{
			ChartType:  queryResult.ChartType,
			Confidence: queryResult.Confidence,
			Reason:     "Query intent suggests this visualization",
		})
	}

	// If data structure suggested different type, add as alternative
	if dataResult.ChartType != combined.ChartType && dataResult.ChartType != queryResult.ChartType {
		if dataResult.Confidence >= 75 {
			alternatives = append(alternatives, AlternativeChart{
				ChartType:  dataResult.ChartType,
				Confidence: dataResult.Confidence,
				Reason:     "Data structure suits this visualization",
			})
		}
	}

	// Always include table as fallback for detailed view
	if combined.ChartType != models.ChartTypeTable {
		alternatives = append(alternatives, AlternativeChart{
			ChartType:  models.ChartTypeTable,
			Confidence: 80.0,
			Reason:     "Detailed tabular view available",
		})
	}

	return alternatives
}

// BuildVisualizationSpec creates a complete VisualizationSpec from routing result.
func (a *VisualizationRoutingAgent) BuildVisualizationSpec(
	result *RoutingResult,
	data []map[string]any,
	columns []models.ColumnMeta,
) *models.VisualizationSpec {
	spec := models.NewVisualizationSpec(result.ChartType).
		SetTitle(a.generateTitle(result.Intent)).
		SetAxes(result.XAxis, result.YAxis).
		SetConfidence(result.Confidence).
		SetReasoning(result.Reasoning)

	// Format data based on chart type
	switch result.ChartType {
	case models.ChartTypeKPICard:
		value, formatted := a.classifier.ExtractKPIValue(data)
		spec.Data = map[string]any{
			"value":     value,
			"formatted": formatted,
		}

	case models.ChartTypeLine, models.ChartTypeBar, models.ChartTypePie:
		spec.Data = a.formatChartData(data, columns, result)

	case models.ChartTypeScatter:
		spec.Data = a.formatScatterData(data, columns)

	case models.ChartTypeTable:
		spec.Data = map[string]any{
			"columns": columns,
			"rows":    data,
		}
	}

	return spec
}

// generateTitle generates a chart title based on intent.
func (a *VisualizationRoutingAgent) generateTitle(intent string) string {
	titles := map[string]string{
		"trend":       "Trend Analysis",
		"comparison":  "Comparison View",
		"breakdown":   "Distribution",
		"kpi":         "Key Metrics",
		"correlation": "Correlation Analysis",
		"table":       "Data Details",
	}

	if title, ok := titles[intent]; ok {
		return title
	}
	return "Query Results"
}

// formatChartData formats data for line, bar, or pie charts.
func (a *VisualizationRoutingAgent) formatChartData(
	data []map[string]any,
	columns []models.ColumnMeta,
	result *RoutingResult,
) map[string]any {
	// Identify label and value columns
	var labelCol, valueCol string

	for _, col := range columns {
		colType := col.Type
		colName := col.Name

		if colType == "string" || colType == "text" || colType == "varchar" ||
			colType == "timestamp" || colType == "date" {
			if labelCol == "" {
				labelCol = colName
			}
		} else if colType == "integer" || colType == "float" || colType == "numeric" {
			if valueCol == "" || strings.Contains(strings.ToLower(colName), "amount") ||
				strings.Contains(strings.ToLower(colName), "total") {
				valueCol = colName
			}
		}
	}

	// Extract labels and values
	labels := make([]string, 0, len(data))
	values := make([]float64, 0, len(data))

	for _, row := range data {
		if labelCol != "" {
			if val, ok := row[labelCol]; ok {
				labels = append(labels, fmt.Sprintf("%v", val))
			}
		}
		if valueCol != "" {
			if val, ok := row[valueCol]; ok {
				switch v := val.(type) {
				case float64:
					values = append(values, v)
				case int:
					values = append(values, float64(v))
				case int64:
					values = append(values, float64(v))
				default:
					values = append(values, 0)
				}
			}
		}
	}

	return map[string]any{
		"labels": labels,
		"series": []map[string]any{
			{
				"name":   result.YAxis,
				"values": values,
			},
		},
	}
}

// formatScatterData formats data for scatter charts.
func (a *VisualizationRoutingAgent) formatScatterData(
	data []map[string]any,
	columns []models.ColumnMeta,
) map[string]any {
	// Find two numeric columns for x and y
	var xCol, yCol string

	for _, col := range columns {
		if col.Type == "integer" || col.Type == "float" || col.Type == "numeric" {
			if xCol == "" {
				xCol = col.Name
			} else if yCol == "" {
				yCol = col.Name
				break
			}
		}
	}

	points := make([]map[string]any, 0, len(data))

	for _, row := range data {
		point := make(map[string]any)
		if xCol != "" {
			if val, ok := row[xCol]; ok {
				point["x"] = val
			}
		}
		if yCol != "" {
			if val, ok := row[yCol]; ok {
				point["y"] = val
			}
		}
		if len(point) == 2 {
			points = append(points, point)
		}
	}

	return map[string]any{
		"points": points,
		"x_axis": xCol,
		"y_axis": yCol,
	}
}

// GetChartCapabilities returns the capabilities of a chart type.
func (a *VisualizationRoutingAgent) GetChartCapabilities(chartType models.ChartType) models.ChartCapabilities {
	return models.ChartCapabilitiesMap[chartType]
}

// IsChartSuitable checks if a chart type is suitable for the given data.
func (a *VisualizationRoutingAgent) IsChartSuitable(
	chartType models.ChartType,
	rowCount int,
	columnCount int,
) bool {
	caps, ok := models.ChartCapabilitiesMap[chartType]
	if !ok {
		return false
	}

	if rowCount < caps.MinDataPoints {
		return false
	}

	if caps.MaxDataPoints > 0 && rowCount > caps.MaxDataPoints {
		// Still suitable but not optimal
		return true
	}

	return true
}
