// Example: Testing a MediSync Agent
// Demonstrates unit testing with mocks and integration testing patterns

package module_a

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockModel is a mock LLM for testing
type MockModel struct {
	mock.Mock
}

func (m *MockModel) Generate(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

// MockWarehouse is a mock database for testing
type MockWarehouse struct {
	mock.Mock
}

func (m *MockWarehouse) Query(ctx context.Context, sql, role string) (any, error) {
	args := m.Called(ctx, sql, role)
	return args.Get(0), args.Error(1)
}

// MockOPA is a mock OPA client for testing
type MockOPA struct {
	mock.Mock
}

func (m *MockOPA) Allow(ctx context.Context, action string, input map[string]interface{}) (bool, error) {
	args := m.Called(ctx, action, input)
	return args.Bool(0), args.Error(1)
}

// TestSimpleQueryAgent_Success tests a successful query flow
func TestSimpleQueryAgent_Success(t *testing.T) {
	// Setup mocks
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	// Create agent
	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	// Setup expectations
	mockOPA.On("Allow", mock.Anything, "warehouse_query", mock.Anything).Return(true, nil)
	mockModel.On("Generate", mock.Anything, mock.Anything).Return("SELECT * FROM patients", nil)
	mockDB.On("Query", mock.Anything, "SELECT * FROM patients", "medisync_readonly").Return([]map[string]any{"id": 1}, nil)

	// Execute
	req := SimpleQueryRequest{
		Query:     "Show all patients",
		Locale:    "en",
		UserID:    "user-123",
		CompanyID: "company-456",
	}

	resp, err := agent.ProcessFlow(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM patients", resp.SQL)
	assert.Greater(t, resp.Confidence, 0.0)
	assert.NotEmpty(t, resp.Explanation)

	// Verify mocks were called
	mockOPA.AssertExpectations(t)
	mockModel.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// TestSimpleQueryAgent_EmptyQuery tests validation failure
func TestSimpleQueryAgent_EmptyQuery(t *testing.T) {
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	req := SimpleQueryRequest{
		Query:     "",
		Locale:    "en",
		UserID:    "user-123",
		CompanyID: "company-456",
	}

	resp, err := agent.ProcessFlow(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "query cannot be empty")
}

// TestSimpleQueryAgent_NonSelectSQL_Rejected tests security validation
func TestSimpleQueryAgent_NonSelectSQL_Rejected(t *testing.T) {
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	mockOPA.On("Allow", mock.Anything, "warehouse_query", mock.Anything).Return(true, nil)
	mockModel.On("Generate", mock.Anything, mock.Anything).Return("DELETE FROM patients", nil)

	req := SimpleQueryRequest{
		Query:     "Delete all patients",
		Locale:    "en",
		UserID:    "user-123",
		CompanyID: "company-456",
	}

	resp, err := agent.ProcessFlow(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not SELECT-only")
}

// TestSimpleQueryAgent_UnauthorizedUser tests OPA enforcement
func TestSimpleQueryAgent_UnauthorizedUser(t *testing.T) {
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	// OPA denies access
	mockOPA.On("Allow", mock.Anything, "warehouse_query", mock.Anything).Return(false, nil)

	req := SimpleQueryRequest{
		Query:     "Show patients",
		Locale:    "en",
		UserID:    "unauthorized-user",
		CompanyID: "company-456",
	}

	resp, err := agent.ProcessFlow(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not authorized")
}

// TestSimpleQueryAgent_ArabicLocale tests i18n support
func TestSimpleQueryAgent_ArabicLocale(t *testing.T) {
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	mockOPA.On("Allow", mock.Anything, "warehouse_query", mock.Anything).Return(true, nil)
	mockModel.On("Generate", mock.Anything, mock.MatchedBy(func(p string) bool {
		return strings.Contains(p, "Respond in Arabic")
	})).Return("SELECT * FROM patients", nil)
	mockDB.On("Query", mock.Anything, "SELECT * FROM patients", "medisync_readonly").Return([]map[string]any{}, nil)

	req := SimpleQueryRequest{
		Query:     "كم عدد المرضى؟",
		Locale:    "ar",
		UserID:    "user-123",
		CompanyID: "company-456",
	}

	resp, err := agent.ProcessFlow(context.Background(), req)

	require.NoError(t, err)
	assert.Contains(t, resp.Explanation, "تم تنفيذ")
}

// TestIsSelectOnly tests the SQL validation function
func TestIsSelectOnly(t *testing.T) {
	agent := &SimpleQueryAgent{}

	tests := []struct {
		name     string
		sql      string
		expected bool
	}{
		{"valid select", "SELECT * FROM patients", true},
		{"select with join", "SELECT p.* FROM patients p JOIN doctors d ON p.doctor_id = d.id", true},
		{"lowercase select", "select from patients", false},
		{"insert", "INSERT INTO patients VALUES (...)", false},
		{"update", "UPDATE patients SET name = 'John'", false},
		{"delete", "DELETE FROM patients", false},
		{"drop", "DROP TABLE patients", false},
		{"create", "CREATE TABLE new_table", false},
		{"select with comment", "-- COMMENT\nSELECT * FROM patients", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.isSelectOnly(tt.sql)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// BenchmarkProcessFlow benchmarks the agent flow
func BenchmarkProcessFlow(b *testing.B) {
	mockModel := new(MockModel)
	mockDB := new(MockWarehouse)
	mockOPA := new(MockOPA)

	agent := NewSimpleQueryAgent(mockModel, mockDB, mockOPA)

	mockOPA.On("Allow", mock.Anything, "warehouse_query", mock.Anything).Return(true, nil)
	mockModel.On("Generate", mock.Anything, mock.Anything).Return("SELECT * FROM patients", nil)
	mockDB.On("Query", mock.Anything, "SELECT * FROM patients", "medisync_readonly").Return([]map[string]any{}, nil)

	req := SimpleQueryRequest{
		Query:     "Show patients",
		Locale:    "en",
		UserID:    "user-123",
		CompanyID: "company-456",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = agent.ProcessFlow(context.Background(), req)
	}
}
