package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Intent constants define the types of query intents that can be detected
// from natural language queries.
const (
	// IntentTrend indicates a query about trends over time
	IntentTrend = "trend"
	// IntentComparison indicates a query comparing two or more entities
	IntentComparison = "comparison"
	// IntentBreakdown indicates a query asking for a breakdown or distribution
	IntentBreakdown = "breakdown"
	// IntentKPI indicates a query asking for key performance indicators
	IntentKPI = "kpi"
	// IntentTable indicates a query requesting tabular data output
	IntentTable = "table"
)

// ValidIntents contains all supported query intent types.
var ValidIntents = map[string]bool{
	IntentTrend:      true,
	IntentComparison: true,
	IntentBreakdown:  true,
	IntentKPI:        true,
	IntentTable:      true,
}

// Query constraints
const (
	// MaxRawTextLength is the maximum allowed length for raw query text
	MaxRawTextLength = 2000
)

// NaturalLanguageQuery represents a user's natural language query within a session.
// It captures the raw input text along with detected locale, intent, and metadata.
type NaturalLanguageQuery struct {
	// ID is the unique identifier for the query
	ID uuid.UUID `json:"id" db:"id"`
	// SessionID references the parent query session
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	// RawText contains the original natural language query from the user
	RawText string `json:"raw_text" db:"raw_text"`
	// DetectedLocale is the language code detected from the query ("en" or "ar")
	DetectedLocale string `json:"detected_locale" db:"detected_locale"`
	// DetectedIntent is the classified intent of the query
	DetectedIntent string `json:"detected_intent" db:"detected_intent"`
	// CreatedAt is the timestamp when the query was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate checks if the NaturalLanguageQuery has valid field values.
// It ensures the raw text is present, within length limits, and locale/intent are valid.
func (q *NaturalLanguageQuery) Validate() error {
	if q.RawText == "" {
		return errors.New("raw_text is required and cannot be empty")
	}

	if len(q.RawText) > MaxRawTextLength {
		return fmt.Errorf("raw_text exceeds maximum length of %d characters", MaxRawTextLength)
	}

	if !ValidLocales[q.DetectedLocale] {
		return fmt.Errorf("invalid detected_locale '%s': must be one of 'en' or 'ar'", q.DetectedLocale)
	}

	if q.DetectedIntent != "" && !ValidIntents[q.DetectedIntent] {
		return fmt.Errorf("invalid detected_intent '%s': must be one of 'trend', 'comparison', 'breakdown', 'kpi', or 'table'", q.DetectedIntent)
	}

	return nil
}

// NewNaturalLanguageQuery creates a new NaturalLanguageQuery with the provided parameters.
// It generates a new UUID and sets the creation timestamp.
func NewNaturalLanguageQuery(sessionID uuid.UUID, rawText, locale, intent string) *NaturalLanguageQuery {
	return &NaturalLanguageQuery{
		ID:             uuid.New(),
		SessionID:      sessionID,
		RawText:        rawText,
		DetectedLocale: locale,
		DetectedIntent: intent,
		CreatedAt:      time.Now(),
	}
}
