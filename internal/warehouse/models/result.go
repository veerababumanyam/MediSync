package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ColumnMeta represents metadata about a column in a query result set.
type ColumnMeta struct {
	// Name is the name of the column
	Name string `json:"name" db:"name"`
	// Type is the data type of the column (e.g., "string", "integer", "float", "timestamp")
	Type string `json:"type" db:"type"`
}

// QueryResult represents the result of executing a SQL statement.
// It contains the row data, column metadata, execution metrics, and any error information.
type QueryResult struct {
	// ID is the unique identifier for the query result
	ID uuid.UUID `json:"id" db:"id"`
	// StatementID references the SQL statement that produced this result
	StatementID uuid.UUID `json:"statement_id" db:"statement_id"`
	// RowCount is the number of rows returned by the query
	RowCount int `json:"row_count" db:"row_count"`
	// Columns contains metadata for each column in the result set
	Columns []ColumnMeta `json:"columns" db:"columns"`
	// Data contains the actual row data as a slice of maps
	Data []map[string]any `json:"data" db:"data"`
	// ExecutionTimeMs is the query execution time in milliseconds
	ExecutionTimeMs int `json:"execution_time_ms" db:"execution_time_ms"`
	// ErrorMessage contains any error that occurred during execution
	ErrorMessage string `json:"error_message" db:"error_message"`
	// CreatedAt is the timestamp when the result was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate checks if the QueryResult has valid field values.
// It ensures the statement ID is present and row count is non-negative.
func (r *QueryResult) Validate() error {
	if r.StatementID == uuid.Nil {
		return errors.New("statement_id is required and cannot be empty")
	}

	if r.RowCount < 0 {
		return errors.New("row_count cannot be negative")
	}

	if r.ExecutionTimeMs < 0 {
		return errors.New("execution_time_ms cannot be negative")
	}

	// Verify column count matches data structure if data exists
	if len(r.Data) > 0 && len(r.Columns) > 0 {
		for i, row := range r.Data {
			if len(row) != len(r.Columns) {
				return fmt.Errorf("row %d has %d columns but expected %d", i, len(row), len(r.Columns))
			}
		}
	}

	return nil
}

// HasError returns true if the query resulted in an error.
func (r *QueryResult) HasError() bool {
	return r.ErrorMessage != ""
}

// ToJSON serializes the QueryResult to JSON format.
// Returns an error if serialization fails.
func (r *QueryResult) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the QueryResult.
// Returns an error if deserialization fails or the JSON is malformed.
func (r *QueryResult) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// NewQueryResult creates a new QueryResult with the provided parameters.
// It generates a new UUID, initializes empty slices, and sets the creation timestamp.
func NewQueryResult(statementID uuid.UUID) *QueryResult {
	return &QueryResult{
		ID:              uuid.New(),
		StatementID:     statementID,
		RowCount:        0,
		Columns:         make([]ColumnMeta, 0),
		Data:            make([]map[string]any, 0),
		ExecutionTimeMs: 0,
		ErrorMessage:    "",
		CreatedAt:       time.Now(),
	}
}

// SetError sets the error message for the query result.
func (r *QueryResult) SetError(err error) {
	if err != nil {
		r.ErrorMessage = err.Error()
	} else {
		r.ErrorMessage = ""
	}
}

// SetData sets the result data and updates the row count accordingly.
func (r *QueryResult) SetData(columns []ColumnMeta, data []map[string]any) {
	r.Columns = columns
	r.Data = data
	r.RowCount = len(data)
}
