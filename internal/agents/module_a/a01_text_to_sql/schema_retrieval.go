// Package a01_text_to_sql provides the text-to-SQL agent subcomponents.
//
// This file implements schema context retrieval from pgvector embeddings.
package a01_text_to_sql

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pgvector/pgvector-go"
)

// SchemaRetriever retrieves relevant schema context from pgvector.
type SchemaRetriever struct {
	vectorStore  VectorStore
	embedder     Embedder
	logger       *slog.Logger
	MaxTables    int
	MaxPatterns  int
	MinSimilarity float64
}

// VectorStore defines the interface for vector storage operations.
type VectorStore interface {
	// SearchSimilar finds embeddings similar to the query vector.
	SearchSimilar(ctx context.Context, vector pgvector.Vector, limit int, filters map[string]string) ([]SchemaEmbedding, error)
	// GetByType retrieves embeddings by type.
	GetByType(ctx context.Context, embeddingType string, limit int) ([]SchemaEmbedding, error)
}

// Embedder defines the interface for creating embeddings.
type Embedder interface {
	// Embed creates an embedding vector from text.
	Embed(ctx context.Context, text string) (pgvector.Vector, error)
}

// SchemaEmbedding represents a schema element with its vector embedding.
type SchemaEmbedding struct {
	ID             int             `json:"id"`
	EmbeddingType  string          `json:"embedding_type"`
	EntityName     string          `json:"entity_name"`
	Description    string          `json:"description"`
	Metadata       json.RawMessage `json:"metadata"`
	Embedding      pgvector.Vector `json:"embedding"`
	Similarity     float64         `json:"similarity,omitempty"`
}

// TableMetadata contains metadata for table embeddings.
type TableMetadata struct {
	Columns           []string `json:"columns"`
	RowCountEstimate  int      `json:"row_count_estimate"`
	Grain             string   `json:"grain"`
	PrimaryKeys       []string `json:"primary_keys"`
	ForeignKeys       []FKRef  `json:"foreign_keys"`
}

// FKRef represents a foreign key reference.
type FKRef struct {
	Column      string `json:"column"`
	RefTable    string `json:"ref_table"`
	RefColumn   string `json:"ref_column"`
}

// QueryPatternMetadata contains metadata for query pattern embeddings.
type QueryPatternMetadata struct {
	TemplateSQL string   `json:"template_sql"`
	Metrics     []string `json:"metrics"`
	Dimensions  []string `json:"dimensions"`
}

// SchemaContext contains the retrieved schema information for SQL generation.
type SchemaContext struct {
	Tables        []TableContext    `json:"tables"`
	QueryPatterns []PatternContext  `json:"query_patterns"`
	RelevanceScore float64          `json:"relevance_score"`
}

// TableContext contains context for a single table.
type TableContext struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Columns     []string `json:"columns"`
	PrimaryKeys []string `json:"primary_keys"`
	ForeignKeys []FKRef  `json:"foreign_keys"`
	Relevance   float64  `json:"relevance"`
}

// PatternContext contains context for a query pattern.
type PatternContext struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TemplateSQL string   `json:"template_sql"`
	Metrics     []string `json:"metrics"`
	Dimensions  []string `json:"dimensions"`
	Relevance   float64  `json:"relevance"`
}

// SchemaRetrieverConfig holds configuration for the retriever.
type SchemaRetrieverConfig struct {
	VectorStore       VectorStore
	Embedder          Embedder
	Logger            *slog.Logger
	MaxTables         int
	MaxPatterns       int
	MinSimilarity     float64
}

// NewSchemaRetriever creates a new schema retriever.
func NewSchemaRetriever(cfg SchemaRetrieverConfig) *SchemaRetriever {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.MaxTables == 0 {
		cfg.MaxTables = 5
	}
	if cfg.MaxPatterns == 0 {
		cfg.MaxPatterns = 3
	}
	if cfg.MinSimilarity == 0 {
		cfg.MinSimilarity = 0.7
	}

	return &SchemaRetriever{
		vectorStore:   cfg.VectorStore,
		embedder:      cfg.Embedder,
		logger:        cfg.Logger.With("component", "schema_retriever"),
		MaxTables:     cfg.MaxTables,
		MaxPatterns:   cfg.MaxPatterns,
		MinSimilarity: cfg.MinSimilarity,
	}
}

// Retrieve retrieves relevant schema context for a natural language query.
func (r *SchemaRetriever) Retrieve(ctx context.Context, query string) (*SchemaContext, error) {
	r.logger.Debug("retrieving schema context", "query", query)

	// Create embedding for the query
	queryVector, err := r.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Search for relevant tables
	tables, err := r.vectorStore.SearchSimilar(ctx, queryVector, r.MaxTables*2, map[string]string{
		"embedding_type": "table",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search tables: %w", err)
	}

	// Search for relevant query patterns
	patterns, err := r.vectorStore.SearchSimilar(ctx, queryVector, r.MaxPatterns*2, map[string]string{
		"embedding_type": "query_pattern",
	})
	if err != nil {
		r.logger.Warn("failed to search patterns", "error", err)
		patterns = []SchemaEmbedding{} // Continue without patterns
	}

	// Build the context
	context := &SchemaContext{
		Tables:        []TableContext{},
		QueryPatterns: []PatternContext{},
	}

	// Process tables
	for _, emb := range tables {
		if emb.Similarity < r.MinSimilarity {
			continue
		}
		if len(context.Tables) >= r.MaxTables {
			break
		}

		tableCtx, err := r.parseTableEmbedding(emb)
		if err != nil {
			r.logger.Warn("failed to parse table embedding", "entity", emb.EntityName, "error", err)
			continue
		}
		context.Tables = append(context.Tables, *tableCtx)
	}

	// Process query patterns
	for _, emb := range patterns {
		if emb.Similarity < r.MinSimilarity {
			continue
		}
		if len(context.QueryPatterns) >= r.MaxPatterns {
			break
		}

		patternCtx, err := r.parsePatternEmbedding(emb)
		if err != nil {
			r.logger.Warn("failed to parse pattern embedding", "entity", emb.EntityName, "error", err)
			continue
		}
		context.QueryPatterns = append(context.QueryPatterns, *patternCtx)
	}

	// Calculate overall relevance score
	context.RelevanceScore = r.calculateRelevanceScore(context)

	r.logger.Info("schema context retrieved",
		"tables", len(context.Tables),
		"patterns", len(context.QueryPatterns),
		"relevance", context.RelevanceScore)

	return context, nil
}

// parseTableEmbedding converts a schema embedding to table context.
func (r *SchemaRetriever) parseTableEmbedding(emb SchemaEmbedding) (*TableContext, error) {
	ctx := &TableContext{
		Name:        emb.EntityName,
		Description: emb.Description,
		Relevance:   emb.Similarity,
	}

	var metadata TableMetadata
	if err := json.Unmarshal(emb.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	ctx.Columns = metadata.Columns
	ctx.PrimaryKeys = metadata.PrimaryKeys
	ctx.ForeignKeys = metadata.ForeignKeys

	return ctx, nil
}

// parsePatternEmbedding converts a schema embedding to pattern context.
func (r *SchemaRetriever) parsePatternEmbedding(emb SchemaEmbedding) (*PatternContext, error) {
	ctx := &PatternContext{
		Name:        emb.EntityName,
		Description: emb.Description,
		Relevance:   emb.Similarity,
	}

	var metadata QueryPatternMetadata
	if err := json.Unmarshal(emb.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	ctx.TemplateSQL = metadata.TemplateSQL
	ctx.Metrics = metadata.Metrics
	ctx.Dimensions = metadata.Dimensions

	return ctx, nil
}

// calculateRelevanceScore calculates an overall relevance score for the context.
func (r *SchemaRetriever) calculateRelevanceScore(context *SchemaContext) float64 {
	if len(context.Tables) == 0 {
		return 0.0
	}

	totalRelevance := 0.0
	for _, t := range context.Tables {
		totalRelevance += t.Relevance
	}

	avgRelevance := totalRelevance / float64(len(context.Tables))

	// Boost if we have patterns that match
	if len(context.QueryPatterns) > 0 {
		avgRelevance = avgRelevance * 1.1
		if avgRelevance > 1.0 {
			avgRelevance = 1.0
		}
	}

	return avgRelevance
}

// ToPrompt formats the schema context for LLM prompting.
func (c *SchemaContext) ToPrompt() string {
	var sb strings.Builder

	sb.WriteString("## Available Tables\n\n")
	for _, t := range c.Tables {
		sb.WriteString(fmt.Sprintf("### %s\n", t.Name))
		sb.WriteString(fmt.Sprintf("%s\n\n", t.Description))
		if len(t.Columns) > 0 {
			sb.WriteString("Columns: ")
			sb.WriteString(strings.Join(t.Columns, ", "))
			sb.WriteString("\n")
		}
		if len(t.PrimaryKeys) > 0 {
			sb.WriteString("Primary Keys: ")
			sb.WriteString(strings.Join(t.PrimaryKeys, ", "))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if len(c.QueryPatterns) > 0 {
		sb.WriteString("## Similar Query Patterns\n\n")
		for _, p := range c.QueryPatterns {
			sb.WriteString(fmt.Sprintf("### %s\n", p.Name))
			sb.WriteString(fmt.Sprintf("%s\n\n", p.Description))
			if p.TemplateSQL != "" {
				sb.WriteString(fmt.Sprintf("Template: %s\n\n", p.TemplateSQL))
			}
		}
	}

	return sb.String()
}

// GetTableNames returns a list of table names in the context.
func (c *SchemaContext) GetTableNames() []string {
	names := make([]string, len(c.Tables))
	for i, t := range c.Tables {
		names[i] = t.Name
	}
	return names
}

// GetColumnNames returns all unique column names across tables.
func (c *SchemaContext) GetColumnNames() []string {
	seen := make(map[string]bool)
	var columns []string
	for _, t := range c.Tables {
		for _, col := range t.Columns {
			if !seen[col] {
				seen[col] = true
				columns = append(columns, col)
			}
		}
	}
	return columns
}
