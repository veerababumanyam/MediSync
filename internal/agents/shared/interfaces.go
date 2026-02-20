// Package shared defines agent interfaces for the MediSync AI agent system.
//
// These interfaces define the contracts that all agents must implement to work
// within the agent supervisor orchestration framework.
package shared

import (
	"context"
)

// Agent is the base interface that all agents must implement.
type Agent interface {
	// ID returns the unique identifier for the agent.
	ID() string

	// Name returns the human-readable name of the agent.
	Name() string

	// Health checks if the agent is healthy and ready to process requests.
	Health(ctx context.Context) (*AgentHealth, error)
}

// TextToSQLAgent converts natural language queries to SQL.
type TextToSQLAgent interface {
	Agent

	// GenerateSQL generates SQL from a processed query.
	GenerateSQL(ctx context.Context, query *ProcessedQuery, schemaContext string) (*GeneratedSQL, error)

	// ValidateSQL validates that the SQL is safe and correct.
	ValidateSQL(ctx context.Context, sql string) (*SQLValidationResult, error)
}

// SQLCorrectionAgent detects and corrects SQL errors.
type SQLCorrectionAgent interface {
	Agent

	// Correct attempts to correct a failed SQL query.
	Correct(ctx context.Context, sql string, errMsg string, attempt int) (*GeneratedSQL, error)

	// CanRetry determines if the error is retryable.
	CanRetry(ctx context.Context, err error) bool
}

// VisualizationRoutingAgent recommends visualization types.
type VisualizationRoutingAgent interface {
	Agent

	// RecommendChart analyzes query results and recommends the best chart type.
	RecommendChart(ctx context.Context, result *QueryResult, query *ProcessedQuery) (*VisualizationSpec, error)

	// GetAlternatives returns alternative chart recommendations.
	GetAlternatives(ctx context.Context, result *QueryResult) ([]ChartType, error)
}

// TerminologyNormalizerAgent normalizes domain terminology.
type TerminologyNormalizerAgent interface {
	Agent

	// Normalize replaces domain terms with canonical database references.
	Normalize(ctx context.Context, query string, locale string) (*ProcessedQuery, error)

	// GetMappings returns the terminology mappings applied.
	GetMappings(ctx context.Context) ([]TermMapping, error)
}

// HallucinationGuardAgent validates AI outputs.
type HallucinationGuardAgent interface {
	Agent

	// Check validates SQL and results for hallucinations.
	Check(ctx context.Context, sql string, result *QueryResult) (*HallucinationCheckResult, error)

	// ScoreRisk calculates a risk score for the output.
	ScoreRisk(ctx context.Context, sql string, result *QueryResult) (float64, error)
}

// ConfidenceScoringAgent scores confidence of results.
type ConfidenceScoringAgent interface {
	Agent

	// Score calculates an overall confidence score.
	Score(ctx context.Context, sql *GeneratedSQL, result *QueryResult, viz *VisualizationSpec) (*ConfidenceScore, error)

	// NeedsClarification determines if user clarification is needed.
	NeedsClarification(ctx context.Context, score *ConfidenceScore) (bool, string)
}

// LanguageDetectionAgent detects query language.
type LanguageDetectionAgent interface {
	Agent

	// Detect determines the language of the query.
	Detect(ctx context.Context, query string) (*LanguageDetectionResult, error)

	// IsArabic checks if the query is primarily Arabic.
	IsArabic(ctx context.Context, query string) bool
}

// QueryTranslationAgent translates queries between languages.
type QueryTranslationAgent interface {
	Agent

	// Translate translates a query from source to target language.
	Translate(ctx context.Context, query string, sourceLang string, targetLang string) (*TranslationResult, error)

	// TranslateToEnglish translates Arabic queries to English.
	TranslateToEnglish(ctx context.Context, arabicQuery string) (*TranslationResult, error)
}

// LocalizedFormatterAgent formats values for display.
type LocalizedFormatterAgent interface {
	Agent

	// Format formats a value according to locale.
	Format(ctx context.Context, value interface{}, valueType string, locale string) (*FormattedValue, error)

	// FormatResult formats an entire query result for display.
	FormatResult(ctx context.Context, result *QueryResult, locale string) (interface{}, error)

	// GetLocaleConfig returns locale-specific configuration.
	GetLocaleConfig(locale string) *LocaleConfig
}

// QueryExecutor executes SQL queries safely.
type QueryExecutor interface {
	// ExecuteQuery executes a SQL query and returns results.
	ExecuteQuery(ctx context.Context, sql string, sessionID string) (*QueryResult, error)

	// Ping checks if the database connection is healthy.
	Ping(ctx context.Context) error
}

// EventPublisher publishes SSE events during query processing.
type EventPublisher interface {
	// Publish publishes an SSE event.
	Publish(event *SSEEvent)

	// PublishThinking publishes a thinking status event.
	PublishThinking(message string)

	// PublishSQLPreview publishes the generated SQL.
	PublishSQLPreview(sql string)

	// PublishResult publishes the final result.
	PublishResult(viz *VisualizationSpec, data interface{}, confidence float64)

	// PublishError publishes an error event.
	PublishError(message string)

	// PublishClarification requests user clarification.
	PublishClarification(message string, options []string)
}
