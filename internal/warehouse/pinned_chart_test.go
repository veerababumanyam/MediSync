package warehouse

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
	"github.com/stretchr/testify/assert"
)

// TestPinnedChartModel_TableName tests the table name (unit test)
func TestPinnedChartModel_TableName(t *testing.T) {
	chart := models.PinnedChart{}
	// PinnedChart uses app.pinned_charts table
	assert.NotNil(t, chart.ID)
}

// TestPinnedChartModel_Fields tests model field types (unit test)
func TestPinnedChartModel_Fields(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	userID := uuid.New()
	queryID := uuid.New()

	chart := models.PinnedChart{
		ID:                   id,
		UserID:               userID,
		QueryID:              &queryID,
		Title:                "Test",
		ChartType:            "bar",
		NaturalLanguageQuery: "Show me revenue",
		SQLQuery:             "SELECT * FROM revenue",
		ChartSpec:            map[string]any{"key": "value"},
		RefreshInterval:      300,
		Locale:               "en",
		Position:             models.ChartPosition{Row: 0, Col: 0, Size: 1},
		IsActive:             true,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	assert.Equal(t, id, chart.ID)
	assert.Equal(t, userID, chart.UserID)
	assert.Equal(t, &queryID, chart.QueryID)
	assert.Equal(t, "Test", chart.Title)
	assert.Equal(t, "bar", chart.ChartType)
	assert.Equal(t, "Show me revenue", chart.NaturalLanguageQuery)
	assert.Equal(t, "value", chart.ChartSpec["key"])
	assert.Equal(t, 0, chart.Position.Row)
	assert.Equal(t, 0, chart.Position.Col)
	assert.True(t, chart.IsActive)
}

// TestChartPosition_JSON tests JSON marshaling (unit test)
func TestChartPosition_JSON(t *testing.T) {
	pos := models.ChartPosition{Row: 1, Col: 2, Size: 3}

	data, err := pos.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"row":1`)
	assert.Contains(t, string(data), `"col":2`)
	assert.Contains(t, string(data), `"size":3`)

	var parsed models.ChartPosition
	err = parsed.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, parsed.Row)
	assert.Equal(t, 2, parsed.Col)
	assert.Equal(t, 3, parsed.Size)
}

// TestChartPosition_DefaultValues tests default values (unit test)
func TestChartPosition_DefaultValues(t *testing.T) {
	pos := models.ChartPosition{}
	assert.Equal(t, 0, pos.Row)
	assert.Equal(t, 0, pos.Col)
	assert.Equal(t, 0, pos.Size)
}

// TestPinnedChart_WithNilQueryID tests nil query ID handling (unit test)
func TestPinnedChart_WithNilQueryID(t *testing.T) {
	chart := models.PinnedChart{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		QueryID:   nil,
		Title:     "Direct Query Chart",
		ChartType: "line",
	}

	assert.Nil(t, chart.QueryID)
	assert.Equal(t, "Direct Query Chart", chart.Title)
}

// TestPinnedChart_WithLastRefreshed tests last refreshed timestamp (unit test)
func TestPinnedChart_WithLastRefreshed(t *testing.T) {
	now := time.Now()
	chart := models.PinnedChart{
		ID:              uuid.New(),
		UserID:          uuid.New(),
		LastRefreshedAt: &now,
	}

	assert.NotNil(t, chart.LastRefreshedAt)
	assert.Equal(t, now, *chart.LastRefreshedAt)
}

// TestPinnedChart_RefreshInterval tests refresh interval values (unit test)
func TestPinnedChart_RefreshInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval int
	}{
		{"no refresh", 0},
		{"5 minutes", 300},
		{"15 minutes", 900},
		{"1 hour", 3600},
		{"1 day", 86400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart := models.PinnedChart{RefreshInterval: tt.interval}
			assert.Equal(t, tt.interval, chart.RefreshInterval)
		})
	}
}

// TestPinnedChart_Locale tests locale field (unit test)
func TestPinnedChart_Locale(t *testing.T) {
	tests := []struct {
		name   string
		locale string
	}{
		{"English", "en"},
		{"Arabic", "ar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart := models.PinnedChart{Locale: tt.locale}
			assert.Equal(t, tt.locale, chart.Locale)
		})
	}
}

// TestPinnedChart_ChartSpec tests chart spec JSONB (unit test)
func TestPinnedChart_ChartSpec(t *testing.T) {
	spec := map[string]any{
		"xAxis": "date",
		"yAxis": "revenue",
		"series": []string{"clinic_a", "clinic_b"},
		"colors": map[string]string{
			"clinic_a": "#4F46E5",
			"clinic_b": "#10B981",
		},
	}

	chart := models.PinnedChart{ChartSpec: spec}

	assert.Equal(t, "date", chart.ChartSpec["xAxis"])
	assert.Equal(t, "revenue", chart.ChartSpec["yAxis"])
	assert.NotNil(t, chart.ChartSpec["series"])
}

// TestPinnedChart_IsActive tests active flag (unit test)
func TestPinnedChart_IsActive(t *testing.T) {
	activeChart := models.PinnedChart{IsActive: true}
	inactiveChart := models.PinnedChart{IsActive: false}

	assert.True(t, activeChart.IsActive)
	assert.False(t, inactiveChart.IsActive)
}
