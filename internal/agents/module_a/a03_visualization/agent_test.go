// Package a03_visualization provides the A-03 Visualization Routing Agent.
//
// This file contains tests for the visualization routing agent to validate
// the 98% accuracy target for chart type assignment.
package a03_visualization

import (
	"context"
	"testing"

	"github.com/medisync/medisync/internal/warehouse/models"
)

// TestClassifierQueryIntent tests query intent classification.
func TestClassifierQueryIntent(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		name           string
		query          string
		expectedType   models.ChartType
		minConfidence  float64
	}{
		// Trend patterns
		{
			name:          "trend_over_time",
			query:         "Show me revenue trends over time",
			expectedType:  models.ChartTypeLine,
			minConfidence: 80,
		},
		{
			name:          "trend_monthly",
			query:         "What is the monthly growth rate?",
			expectedType:  models.ChartTypeLine,
			minConfidence: 80,
		},
		{
			name:          "trend_last_months",
			query:         "Show sales for the last 6 months",
			expectedType:  models.ChartTypeLine,
			minConfidence: 80,
		},

		// Comparison patterns
		{
			name:          "comparison_departments",
			query:         "Compare revenue by department",
			expectedType:  models.ChartTypeBar,
			minConfidence: 80,
		},
		{
			name:          "comparison_top",
			query:         "Show top 10 doctors by revenue",
			expectedType:  models.ChartTypeBar,
			minConfidence: 80,
		},
		{
			name:          "comparison_versus",
			query:         "Show clinic versus pharmacy revenue",
			expectedType:  models.ChartTypeBar,
			minConfidence: 80,
		},

		// Breakdown patterns
		{
			name:          "breakdown_distribution",
			query:         "Show the distribution of patients by age group",
			expectedType:  models.ChartTypePie,
			minConfidence: 80,
		},
		{
			name:          "breakdown_percentage",
			query:         "What percentage of revenue comes from each department?",
			expectedType:  models.ChartTypePie,
			minConfidence: 80,
		},
		{
			name:          "breakdown_segment",
			query:         "Segment sales by product category",
			expectedType:  models.ChartTypePie,
			minConfidence: 80,
		},

		// KPI patterns
		{
			name:          "kpi_total",
			query:         "What is the total revenue?",
			expectedType:  models.ChartTypeKPICard,
			minConfidence: 80,
		},
		{
			name:          "kpi_average",
			query:         "Show average daily sales",
			expectedType:  models.ChartTypeKPICard,
			minConfidence: 80,
		},
		{
			name:          "kpi_count",
			query:         "How many patients visited today?",
			expectedType:  models.ChartTypeKPICard,
			minConfidence: 80,
		},

		// Correlation patterns
		{
			name:          "correlation",
			query:         "Show correlation between patient age and spending",
			expectedType:  models.ChartTypeScatter,
			minConfidence: 80,
		},
		{
			name:          "correlation_relationship",
			query:         "Is there a relationship between price and sales?",
			expectedType:  models.ChartTypeScatter,
			minConfidence: 80,
		},

		// Table patterns
		{
			name:          "table_list",
			query:         "List all patients",
			expectedType:  models.ChartTypeTable,
			minConfidence: 60,
		},
		{
			name:          "table_details",
			query:         "Show detailed billing records",
			expectedType:  models.ChartTypeTable,
			minConfidence: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.ClassifyQuery(tt.query)

			if result.ChartType != tt.expectedType {
				t.Errorf("expected chart type %s, got %s", tt.expectedType, result.ChartType)
			}

			if result.Confidence < tt.minConfidence {
				t.Errorf("expected confidence >= %.1f, got %.1f", tt.minConfidence, result.Confidence)
			}

			if result.Reasoning == "" {
				t.Error("expected non-empty reasoning")
			}
		})
	}
}

// TestClassifierDataStructure tests data structure classification.
func TestClassifierDataStructure(t *testing.T) {
	classifier := NewClassifier()

	tests := []struct {
		name          string
		columns       []models.ColumnMeta
		rowCount      int
		expectedType  models.ChartType
		minConfidence float64
	}{
		{
			name: "single_row_aggregation",
			columns: []models.ColumnMeta{
				{Name: "total_revenue", Type: "float"},
			},
			rowCount:      1,
			expectedType:  models.ChartTypeKPICard,
			minConfidence: 90,
		},
		{
			name: "time_series_data",
			columns: []models.ColumnMeta{
				{Name: "month", Type: "timestamp"},
				{Name: "revenue", Type: "float"},
			},
			rowCount:      20,
			expectedType:  models.ChartTypeLine,
			minConfidence: 85,
		},
		{
			name: "categorical_single_metric",
			columns: []models.ColumnMeta{
				{Name: "department", Type: "string"},
				{Name: "revenue", Type: "float"},
			},
			rowCount:      5,
			expectedType:  models.ChartTypePie,
			minConfidence: 80,
		},
		{
			name: "categorical_many_categories",
			columns: []models.ColumnMeta{
				{Name: "department", Type: "string"},
				{Name: "revenue", Type: "float"},
			},
			rowCount:      15,
			expectedType:  models.ChartTypeBar,
			minConfidence: 85,
		},
		{
			name: "multiple_numeric_correlation",
			columns: []models.ColumnMeta{
				{Name: "age", Type: "integer"},
				{Name: "spending", Type: "float"},
			},
			rowCount:      50,
			expectedType:  models.ChartTypeScatter,
			minConfidence: 80,
		},
		{
			name: "complex_data",
			columns: []models.ColumnMeta{
				{Name: "id", Type: "integer"},
				{Name: "name", Type: "string"},
				{Name: "department", Type: "string"},
				{Name: "revenue", Type: "float"},
				{Name: "cost", Type: "float"},
				{Name: "date", Type: "timestamp"},
			},
			rowCount:      100,
			expectedType:  models.ChartTypeTable,
			minConfidence: 85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.ClassifyDataStructure(tt.columns, tt.rowCount)

			if result.ChartType != tt.expectedType {
				t.Errorf("expected chart type %s, got %s", tt.expectedType, result.ChartType)
			}

			if result.Confidence < tt.minConfidence {
				t.Errorf("expected confidence >= %.1f, got %.1f", tt.minConfidence, result.Confidence)
			}
		})
	}
}

// TestAgentRouting tests the full routing pipeline.
func TestAgentRouting(t *testing.T) {
	agent := NewVisualizationRoutingAgent(&AgentConfig{})
	ctx := context.Background()

	tests := []struct {
		name          string
		query         string
		columns       []models.ColumnMeta
		rowCount      int
		expectedType  models.ChartType
		minConfidence float64
	}{
		{
			name:  "kpi_total_revenue",
			query: "What is the total revenue for this month?",
			columns: []models.ColumnMeta{
				{Name: "total_revenue", Type: "float"},
			},
			rowCount:      1,
			expectedType:  models.ChartTypeKPICard,
			minConfidence: 90,
		},
		{
			name:  "trend_monthly_revenue",
			query: "Show monthly revenue trends for the last year",
			columns: []models.ColumnMeta{
				{Name: "month", Type: "timestamp"},
				{Name: "revenue", Type: "float"},
			},
			rowCount:      12,
			expectedType:  models.ChartTypeLine,
			minConfidence: 85,
		},
		{
			name:  "comparison_department_revenue",
			query: "Compare revenue across departments",
			columns: []models.ColumnMeta{
				{Name: "department", Type: "string"},
				{Name: "revenue", Type: "float"},
			},
			rowCount:      5,
			expectedType:  models.ChartTypeBar,
			minConfidence: 85,
		},
		{
			name:  "breakdown_patient_distribution",
			query: "Show distribution of patients by department",
			columns: []models.ColumnMeta{
				{Name: "department", Type: "string"},
				{Name: "patient_count", Type: "integer"},
			},
			rowCount:      4,
			expectedType:  models.ChartTypePie,
			minConfidence: 85,
		},
		{
			name:  "table_patient_list",
			query: "List all patients with their details",
			columns: []models.ColumnMeta{
				{Name: "patient_id", Type: "integer"},
				{Name: "name", Type: "string"},
				{Name: "email", Type: "string"},
				{Name: "phone", Type: "string"},
				{Name: "department", Type: "string"},
			},
			rowCount:      50,
			expectedType:  models.ChartTypeTable,
			minConfidence: 70,
		},
		{
			name:  "correlation_analysis",
			query: "Show correlation between patient visits and revenue",
			columns: []models.ColumnMeta{
				{Name: "visit_count", Type: "integer"},
				{Name: "total_revenue", Type: "float"},
			},
			rowCount:      30,
			expectedType:  models.ChartTypeScatter,
			minConfidence: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.QueryResult{
				RowCount: tt.rowCount,
				Columns:  tt.columns,
			}

			routing, err := agent.Route(ctx, tt.query, result)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if routing.ChartType != tt.expectedType {
				t.Errorf("expected chart type %s, got %s (reasoning: %s)",
					tt.expectedType, routing.ChartType, routing.Reasoning)
			}

			if routing.Confidence < tt.minConfidence {
				t.Errorf("expected confidence >= %.1f, got %.1f", tt.minConfidence, routing.Confidence)
			}

			if routing.Reasoning == "" {
				t.Error("expected non-empty reasoning")
			}
		})
	}
}

// TestAccuracySuite tests the 98% accuracy target with comprehensive test cases.
func TestAccuracySuite(t *testing.T) {
	agent := NewVisualizationRoutingAgent(&AgentConfig{})
	ctx := context.Background()

	// Comprehensive test cases representing real-world queries
	testCases := []struct {
		query         string
		columns       []models.ColumnMeta
		rowCount      int
		expectedType  models.ChartType
	}{
		// KPI queries (15 cases)
		{"What is the total revenue?", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Show total sales for today", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"How many patients do we have?", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"What is the average bill amount?", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Show current inventory value", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"What is the profit margin?", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Total appointments this week", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Sum of all invoices", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Count of active patients", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Average daily revenue", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"What is the revenue", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Show me total collections", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Current month profit", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Total outstanding bills", singleNumericColumn(), 1, models.ChartTypeKPICard},
		{"Number of doctors", singleNumericColumn(), 1, models.ChartTypeKPICard},

		// Trend queries (15 cases)
		{"Show revenue trend over time", timeSeriesColumns(), 24, models.ChartTypeLine},
		{"Monthly sales growth", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Weekly patient visits history", timeSeriesColumns(), 52, models.ChartTypeLine},
		{"Daily revenue for last 30 days", timeSeriesColumns(), 30, models.ChartTypeLine},
		{"Yearly growth progression", timeSeriesColumns(), 5, models.ChartTypeBar},
		{"Revenue over the past year", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Trend of appointments", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Sales evolution by month", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Show me trends in patient count", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Revenue history for last 6 months", timeSeriesColumns(), 6, models.ChartTypeLine},
		{"Track monthly expenses", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Progression of profit", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"Last 3 months daily trend", timeSeriesColumns(), 90, models.ChartTypeLine},
		{"Weekly trend of new patients", timeSeriesColumns(), 12, models.ChartTypeLine},
		{"How has revenue changed over time", timeSeriesColumns(), 12, models.ChartTypeLine},

		// Comparison queries (15 cases)
		{"Compare revenue by department", categoricalNumericColumns(), 5, models.ChartTypeBar},
		{"Top 10 doctors by patient count", categoricalNumericColumns(), 10, models.ChartTypeBar},
		{"Revenue comparison by clinic", categoricalNumericColumns(), 3, models.ChartTypeBar},
		{"Sales by region", categoricalNumericColumns(), 8, models.ChartTypeBar},
		{"Compare doctor versus nurse salaries", categoricalNumericColumns(), 2, models.ChartTypeBar},
		{"Department wise performance", categoricalNumericColumns(), 6, models.ChartTypeBar},
		{"Highest revenue departments", categoricalNumericColumns(), 5, models.ChartTypeBar},
		{"Bottom 5 products by sales", categoricalNumericColumns(), 5, models.ChartTypeBar},
		{"Revenue by doctor", categoricalNumericColumns(), 15, models.ChartTypeBar},
		{"Compare this month vs last month", categoricalNumericColumns(), 2, models.ChartTypeBar},
		{"Sales difference between branches", categoricalNumericColumns(), 4, models.ChartTypeBar},
		{"Most profitable categories", categoricalNumericColumns(), 8, models.ChartTypeBar},
		{"Department revenue ranking", categoricalNumericColumns(), 6, models.ChartTypeBar},
		{"Revenue by service type", categoricalNumericColumns(), 7, models.ChartTypeBar},
		{"Top performing doctors", categoricalNumericColumns(), 10, models.ChartTypeBar},

		// Breakdown queries (15 cases)
		{"Breakdown of expenses by category", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"Distribution of patients by age", categoricalNumericColumns(), 6, models.ChartTypePie},
		{"Revenue share by department", categoricalNumericColumns(), 4, models.ChartTypePie},
		{"What percentage of sales is each product", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"Split revenue by region", categoricalNumericColumns(), 4, models.ChartTypePie},
		{"Segment patients by gender", categoricalNumericColumns(), 3, models.ChartTypePie},
		{"Proportion of services used", categoricalNumericColumns(), 6, models.ChartTypePie},
		{"Parts of total expenses", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"Composition of revenue streams", categoricalNumericColumns(), 4, models.ChartTypePie},
		{"How is budget distributed", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"Patient breakdown by insurance", categoricalNumericColumns(), 4, models.ChartTypePie},
		{"Revenue distribution by payment method", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"Share of each product category", categoricalNumericColumns(), 6, models.ChartTypePie},
		{"Expense distribution", categoricalNumericColumns(), 5, models.ChartTypePie},
		{"What percent of total is each branch", categoricalNumericColumns(), 4, models.ChartTypePie},

		// Correlation queries (10 cases)
		{"Correlation between price and sales", numericColumns(), 50, models.ChartTypeScatter},
		{"Relationship between age and spending", numericColumns(), 100, models.ChartTypeScatter},
		{"Scatter plot of visits vs revenue", numericColumns(), 30, models.ChartTypeScatter},
		{"Plot hours worked against productivity", numericColumns(), 40, models.ChartTypeScatter},
		{"Compare price against demand trends", numericColumns(), 50, models.ChartTypeScatter},
		{"Is there a correlation between marketing and sales", numericColumns(), 25, models.ChartTypeScatter},
		{"Show relationship of experience to salary", numericColumns(), 60, models.ChartTypeScatter},
		{"Analysis of quantity vs discount", numericColumns(), 40, models.ChartTypeScatter},
		{"Patient age versus visit frequency correlation", numericColumns(), 80, models.ChartTypeScatter},
		{"Revenue versus expense relationship", numericColumns(), 30, models.ChartTypeScatter},

		// Table queries (10 cases)
		{"List all patients", manyColumns(), 100, models.ChartTypeTable},
		{"Show all billing records", manyColumns(), 200, models.ChartTypeTable},
		{"Detailed appointment information", manyColumns(), 50, models.ChartTypeTable},
		{"Show rows from invoice table", manyColumns(), 75, models.ChartTypeTable},
		{"All doctors with their details", manyColumns(), 30, models.ChartTypeTable},
		{"Complete patient list with contact info", manyColumns(), 150, models.ChartTypeTable},
		{"Show me the full inventory list", manyColumns(), 80, models.ChartTypeTable},
		{"Display all transactions", manyColumns(), 300, models.ChartTypeTable},
		{"Get all employee records", manyColumns(), 40, models.ChartTypeTable},
		{"Every bill in the system", manyColumns(), 500, models.ChartTypeTable},
	}

	correct := 0
	total := len(testCases)

	for _, tc := range testCases {
		result := &models.QueryResult{
			RowCount: tc.rowCount,
			Columns:  tc.columns,
		}

		routing, err := agent.Route(ctx, tc.query, result)
		if err != nil {
			t.Errorf("error routing query '%s': %v", tc.query, err)
			continue
		}

		if routing.ChartType == tc.expectedType {
			correct++
		} else {
			t.Logf("MISMATCH: query='%s' expected=%s got=%s confidence=%.1f reasoning=%s",
				tc.query, tc.expectedType, routing.ChartType, routing.Confidence, routing.Reasoning)
		}
	}

	accuracy := float64(correct) / float64(total) * 100
	t.Logf("Accuracy: %.2f%% (%d/%d correct)", accuracy, correct, total)

	if accuracy < 98.0 {
		t.Errorf("Accuracy %.2f%% is below the 98%% target", accuracy)
	}
}

// Helper functions for test data

func singleNumericColumn() []models.ColumnMeta {
	return []models.ColumnMeta{
		{Name: "value", Type: "float"},
	}
}

func timeSeriesColumns() []models.ColumnMeta {
	return []models.ColumnMeta{
		{Name: "date", Type: "timestamp"},
		{Name: "value", Type: "float"},
	}
}

func categoricalNumericColumns() []models.ColumnMeta {
	return []models.ColumnMeta{
		{Name: "category", Type: "string"},
		{Name: "value", Type: "float"},
	}
}

func numericColumns() []models.ColumnMeta {
	return []models.ColumnMeta{
		{Name: "x", Type: "float"},
		{Name: "y", Type: "float"},
	}
}

func manyColumns() []models.ColumnMeta {
	return []models.ColumnMeta{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "string"},
		{Name: "email", Type: "string"},
		{Name: "phone", Type: "string"},
		{Name: "address", Type: "string"},
		{Name: "status", Type: "string"},
		{Name: "created_at", Type: "timestamp"},
	}
}

// TestBuildVisualizationSpec tests the visualization spec builder.
func TestBuildVisualizationSpec(t *testing.T) {
	agent := NewVisualizationRoutingAgent(&AgentConfig{})

	tests := []struct {
		name       string
		result     *RoutingResult
		data       []map[string]any
		columns    []models.ColumnMeta
		expectType models.ChartType
	}{
		{
			name: "kpi_card_spec",
			result: &RoutingResult{
				ChartType:  models.ChartTypeKPICard,
				Confidence: 95,
				Intent:     "kpi",
			},
			data: []map[string]any{
				{"total_revenue": 1250000.00},
			},
			columns:    singleNumericColumn(),
			expectType: models.ChartTypeKPICard,
		},
		{
			name: "line_chart_spec",
			result: &RoutingResult{
				ChartType:  models.ChartTypeLine,
				Confidence: 90,
				Intent:     "trend",
				XAxis:      "month",
				YAxis:      "revenue",
			},
			data: []map[string]any{
				{"month": "Jan", "revenue": 100000.0},
				{"month": "Feb", "revenue": 120000.0},
			},
			columns:    timeSeriesColumns(),
			expectType: models.ChartTypeLine,
		},
		{
			name: "table_spec",
			result: &RoutingResult{
				ChartType:  models.ChartTypeTable,
				Confidence: 85,
				Intent:     "table",
			},
			data: []map[string]any{
				{"id": 1, "name": "John"},
				{"id": 2, "name": "Jane"},
			},
			columns:    manyColumns(),
			expectType: models.ChartTypeTable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := agent.BuildVisualizationSpec(tt.result, tt.data, tt.columns)

			if spec.Type != tt.expectType {
				t.Errorf("expected type %s, got %s", tt.expectType, spec.Type)
			}

			if spec.Data == nil {
				t.Error("expected non-nil data")
			}
		})
	}
}
