// Package models provides data models for the MediSync warehouse.
//
// This file defines the VisualizationSpec model and related types for
// chart type routing and visualization generation in the A-03 agent.
package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ChartType defines the supported visualization types for query results.
type ChartType string

const (
	// ChartTypeLine is for time series and trend visualization.
	ChartTypeLine ChartType = "lineChart"
	// ChartTypeBar is for comparison and categorical data visualization.
	ChartTypeBar ChartType = "barChart"
	// ChartTypePie is for breakdown and distribution visualization.
	ChartTypePie ChartType = "pieChart"
	// ChartTypeKPICard is for single value/KPI display.
	ChartTypeKPICard ChartType = "kpiCard"
	// ChartTypeTable is for detailed tabular data display.
	ChartTypeTable ChartType = "dataTable"
	// ChartTypeScatter is for correlation and relationship visualization.
	ChartTypeScatter ChartType = "scatterChart"
)

// ValidChartTypes contains all supported chart types.
var ValidChartTypes = map[ChartType]bool{
	ChartTypeLine:    true,
	ChartTypeBar:     true,
	ChartTypePie:     true,
	ChartTypeKPICard: true,
	ChartTypeTable:   true,
	ChartTypeScatter: true,
}

// VisualizationSpec represents a complete visualization specification
// for rendering charts in the frontend.
type VisualizationSpec struct {
	// Type is the chart type to render.
	Type ChartType `json:"type" db:"type"`
	// Data contains the chart data (structure varies by chart type).
	Data any `json:"data" db:"data"`
	// Options contains chart-specific configuration options.
	Options map[string]any `json:"options,omitempty" db:"options"`
	// XAxis is the label for the x-axis (for line, bar, scatter charts).
	XAxis string `json:"x_axis,omitempty" db:"x_axis"`
	// YAxis is the label for the y-axis (for line, bar, scatter charts).
	YAxis string `json:"y_axis,omitempty" db:"y_axis"`
	// Title is the chart title.
	Title string `json:"title,omitempty" db:"title"`
	// Confidence is the routing confidence score (0-100).
	Confidence float64 `json:"confidence,omitempty" db:"confidence"`
	// Reasoning explains why this chart type was selected.
	Reasoning string `json:"reasoning,omitempty" db:"reasoning"`
}

// Validate checks if the VisualizationSpec has valid field values.
func (v *VisualizationSpec) Validate() error {
	if !ValidChartTypes[v.Type] {
		return fmt.Errorf("invalid chart type '%s': must be one of lineChart, barChart, pieChart, kpiCard, dataTable, scatterChart", v.Type)
	}

	if v.Confidence < 0 || v.Confidence > 100 {
		return errors.New("confidence must be between 0 and 100")
	}

	return nil
}

// ToJSON serializes the VisualizationSpec to JSON format.
func (v *VisualizationSpec) ToJSON() ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON deserializes JSON data into the VisualizationSpec.
func (v *VisualizationSpec) FromJSON(data []byte) error {
	return json.Unmarshal(data, v)
}

// NewVisualizationSpec creates a new VisualizationSpec with default values.
func NewVisualizationSpec(chartType ChartType) *VisualizationSpec {
	return &VisualizationSpec{
		Type:    chartType,
		Options: make(map[string]any),
	}
}

// SetData sets the data for the visualization.
func (v *VisualizationSpec) SetData(data any) *VisualizationSpec {
	v.Data = data
	return v
}

// SetAxes sets the x and y axis labels.
func (v *VisualizationSpec) SetAxes(xAxis, yAxis string) *VisualizationSpec {
	v.XAxis = xAxis
	v.YAxis = yAxis
	return v
}

// SetTitle sets the chart title.
func (v *VisualizationSpec) SetTitle(title string) *VisualizationSpec {
	v.Title = title
	return v
}

// SetConfidence sets the routing confidence score.
func (v *VisualizationSpec) SetConfidence(confidence float64) *VisualizationSpec {
	v.Confidence = confidence
	return v
}

// SetReasoning sets the reasoning for the chart selection.
func (v *VisualizationSpec) SetReasoning(reasoning string) *VisualizationSpec {
	v.Reasoning = reasoning
	return v
}

// SetOption sets a chart-specific option.
func (v *VisualizationSpec) SetOption(key string, value any) *VisualizationSpec {
	if v.Options == nil {
		v.Options = make(map[string]any)
	}
	v.Options[key] = value
	return v
}

// IsTimeSeries returns true if the visualization is suitable for time series data.
func (v *VisualizationSpec) IsTimeSeries() bool {
	return v.Type == ChartTypeLine || v.Type == ChartTypeBar
}

// IsComparison returns true if the visualization is suitable for comparison data.
func (v *VisualizationSpec) IsComparison() bool {
	return v.Type == ChartTypeBar || v.Type == ChartTypeTable
}

// IsBreakdown returns true if the visualization is suitable for breakdown data.
func (v *VisualizationSpec) IsBreakdown() bool {
	return v.Type == ChartTypePie || v.Type == ChartTypeBar
}

// IsSingleValue returns true if the visualization displays a single value.
func (v *VisualizationSpec) IsSingleValue() bool {
	return v.Type == ChartTypeKPICard
}

// ChartCapabilities describes what a chart type can visualize.
type ChartCapabilities struct {
	// SupportsTimeSeries indicates if the chart can show trends over time.
	SupportsTimeSeries bool
	// SupportsComparison indicates if the chart can compare categories.
	SupportsComparison bool
	// SupportsBreakdown indicates if the chart can show distribution.
	SupportsBreakdown bool
	// SupportsSingleValue indicates if the chart can display KPIs.
	SupportsSingleValue bool
	// SupportsCorrelation indicates if the chart can show relationships.
	SupportsCorrelation bool
	// MaxDataPoints is the recommended maximum data points for clarity.
	MaxDataPoints int
	// MinDataPoints is the minimum data points needed.
	MinDataPoints int
}

// ChartCapabilitiesMap maps chart types to their capabilities.
var ChartCapabilitiesMap = map[ChartType]ChartCapabilities{
	ChartTypeLine: {
		SupportsTimeSeries: true,
		SupportsComparison: true,
		MaxDataPoints:      100,
		MinDataPoints:      2,
	},
	ChartTypeBar: {
		SupportsTimeSeries: true,
		SupportsComparison: true,
		SupportsBreakdown:  true,
		MaxDataPoints:      50,
		MinDataPoints:      2,
	},
	ChartTypePie: {
		SupportsBreakdown: true,
		MaxDataPoints:     8,
		MinDataPoints:     2,
	},
	ChartTypeKPICard: {
		SupportsSingleValue: true,
		MaxDataPoints:       1,
		MinDataPoints:       1,
	},
	ChartTypeTable: {
		SupportsTimeSeries: true,
		SupportsComparison: true,
		SupportsBreakdown:  true,
		MaxDataPoints:      1000,
		MinDataPoints:      1,
	},
	ChartTypeScatter: {
		SupportsCorrelation: true,
		MaxDataPoints:       500,
		MinDataPoints:       5,
	},
}

// GetCapabilities returns the capabilities for a chart type.
func (v *VisualizationSpec) GetCapabilities() ChartCapabilities {
	if caps, ok := ChartCapabilitiesMap[v.Type]; ok {
		return caps
	}
	return ChartCapabilities{}
}
