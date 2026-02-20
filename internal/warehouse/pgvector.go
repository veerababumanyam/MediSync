// Package warehouse provides pgvector operations for schema embeddings.
//
// This file provides the PgVectorClient struct for semantic search over database
// schema embeddings. It uses cosine similarity with the HNSW index for fast
// similarity searches.
//
// The schema_embeddings table stores vector representations of:
//   - Tables (business context, column summaries)
//   - Columns (data types, sample values, descriptions)
//   - Relationships (foreign keys, joins)
//
// Usage:
//
//	client, err := warehouse.NewPgVectorClient(pool, logger)
//	if err != nil {
//	    log.Fatal("Failed to create pgvector client:", err)
//	}
//
//	// Search for relevant schema elements
//	results, err := client.SearchSchema(ctx, queryEmbedding, 10)
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SchemaEmbedding represents a row in the vectors.schema_embeddings table.
type SchemaEmbedding struct {
	// ID is the unique identifier for the embedding record.
	ID uuid.UUID `json:"id"`

	// SchemaName is the database schema name (e.g., "hims_analytics", "tally_analytics").
	SchemaName string `json:"schema_name"`

	// TableName is the table name.
	TableName string `json:"table_name"`

	// ColumnName is the column name (null for table-level embeddings).
	ColumnName *string `json:"column_name,omitempty"`

	// ObjectType indicates the type of schema object: "table", "column", "relationship".
	ObjectType string `json:"object_type"`

	// DescriptionEN is the English description of the schema element.
	DescriptionEN string `json:"description_en"`

	// DescriptionAR is the Arabic description of the schema element.
	DescriptionAR *string `json:"description_ar,omitempty"`

	// DataType is the PostgreSQL data type (for columns).
	DataType *string `json:"data_type,omitempty"`

	// SampleValues contains example values for the column.
	SampleValues []string `json:"sample_values,omitempty"`

	// BusinessContext explains the business meaning of this schema element.
	BusinessContext string `json:"business_context"`

	// Embedding is the vector representation (1536 dimensions for OpenAI ada-002).
	Embedding []float32 `json:"embedding,omitempty"`

	// ModelName is the name of the embedding model used.
	ModelName string `json:"model_name"`

	// Similarity is the cosine similarity score (populated during search).
	Similarity float64 `json:"similarity,omitempty"`

	// CreatedAt is when the embedding was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the embedding was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// PgVectorClient provides operations for schema embedding search.
type PgVectorClient struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// PgVectorClientConfig holds configuration for the pgvector client.
type PgVectorClientConfig struct {
	// Pool is the PostgreSQL connection pool.
	Pool *pgxpool.Pool

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewPgVectorClient creates a new pgvector client for schema embedding operations.
func NewPgVectorClient(pool *PostgresPool, logger *slog.Logger) (*PgVectorClient, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &PgVectorClient{
		pool:   pool.Pool(),
		logger: logger,
	}, nil
}

// NewPgVectorClientFromPool creates a pgvector client from an existing pool.
func NewPgVectorClientFromPool(pool *pgxpool.Pool, logger *slog.Logger) (*PgVectorClient, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &PgVectorClient{
		pool:   pool,
		logger: logger,
	}, nil
}

// SearchSchema searches for schema embeddings using cosine similarity.
// It uses the HNSW index for fast approximate nearest neighbor search.
func (c *PgVectorClient) SearchSchema(ctx context.Context, queryEmbedding []float32, limit int) ([]SchemaEmbedding, error) {
	if len(queryEmbedding) == 0 {
		return nil, fmt.Errorf("warehouse: query embedding is required")
	}

	if limit <= 0 {
		limit = 10
	}

	// Build the vector literal for PostgreSQL
	vectorLiteral := c.buildVectorLiteral(queryEmbedding)

	query := fmt.Sprintf(`
		SELECT
			id,
			schema_name,
			table_name,
			column_name,
			object_type,
			description_en,
			description_ar,
			data_type,
			sample_values,
			business_context,
			model_name,
			1 - (embedding <=> %s) AS similarity,
			created_at,
			updated_at
		FROM vectors.schema_embeddings
		ORDER BY embedding <=> %s
		LIMIT %d
	`, vectorLiteral, vectorLiteral, limit)

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to search schema embeddings: %w", err)
	}
	defer rows.Close()

	var results []SchemaEmbedding
	for rows.Next() {
		var emb SchemaEmbedding
		var sampleValues []byte // pgx reads text[] as []byte

		err := rows.Scan(
			&emb.ID,
			&emb.SchemaName,
			&emb.TableName,
			&emb.ColumnName,
			&emb.ObjectType,
			&emb.DescriptionEN,
			&emb.DescriptionAR,
			&emb.DataType,
			&sampleValues,
			&emb.BusinessContext,
			&emb.ModelName,
			&emb.Similarity,
			&emb.CreatedAt,
			&emb.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan schema embedding: %w", err)
		}

		// Parse sample values from PostgreSQL array format
		if len(sampleValues) > 0 {
			emb.SampleValues = c.parsePostgresArray(string(sampleValues))
		}

		results = append(results, emb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating schema embeddings: %w", err)
	}

	c.logger.Debug("schema embedding search completed",
		slog.Int("results", len(results)),
		slog.Int("limit", limit),
	)

	return results, nil
}

// SearchSchemaByType searches for schema embeddings filtered by object type.
func (c *PgVectorClient) SearchSchemaByType(ctx context.Context, queryEmbedding []float32, objectType string, limit int) ([]SchemaEmbedding, error) {
	if len(queryEmbedding) == 0 {
		return nil, fmt.Errorf("warehouse: query embedding is required")
	}

	if limit <= 0 {
		limit = 10
	}

	vectorLiteral := c.buildVectorLiteral(queryEmbedding)

	query := fmt.Sprintf(`
		SELECT
			id,
			schema_name,
			table_name,
			column_name,
			object_type,
			description_en,
			description_ar,
			data_type,
			sample_values,
			business_context,
			model_name,
			1 - (embedding <=> %s) AS similarity,
			created_at,
			updated_at
		FROM vectors.schema_embeddings
		WHERE object_type = $1
		ORDER BY embedding <=> %s
		LIMIT %d
	`, vectorLiteral, vectorLiteral, limit)

	rows, err := c.pool.Query(ctx, query, objectType)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to search schema embeddings by type: %w", err)
	}
	defer rows.Close()

	return c.scanEmbeddings(rows)
}

// SearchSchemaByTable searches for schema embeddings within a specific table.
func (c *PgVectorClient) SearchSchemaByTable(ctx context.Context, queryEmbedding []float32, schemaName, tableName string, limit int) ([]SchemaEmbedding, error) {
	if len(queryEmbedding) == 0 {
		return nil, fmt.Errorf("warehouse: query embedding is required")
	}

	if limit <= 0 {
		limit = 10
	}

	vectorLiteral := c.buildVectorLiteral(queryEmbedding)

	query := fmt.Sprintf(`
		SELECT
			id,
			schema_name,
			table_name,
			column_name,
			object_type,
			description_en,
			description_ar,
			data_type,
			sample_values,
			business_context,
			model_name,
			1 - (embedding <=> %s) AS similarity,
			created_at,
			updated_at
		FROM vectors.schema_embeddings
		WHERE schema_name = $1 AND table_name = $2
		ORDER BY embedding <=> %s
		LIMIT %d
	`, vectorLiteral, vectorLiteral, limit)

	rows, err := c.pool.Query(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to search schema embeddings by table: %w", err)
	}
	defer rows.Close()

	return c.scanEmbeddings(rows)
}

// GetTableSchema retrieves all embeddings for a specific table.
func (c *PgVectorClient) GetTableSchema(ctx context.Context, schemaName, tableName string) ([]SchemaEmbedding, error) {
	query := `
		SELECT
			id,
			schema_name,
			table_name,
			column_name,
			object_type,
			description_en,
			description_ar,
			data_type,
			sample_values,
			business_context,
			model_name,
			created_at,
			updated_at
		FROM vectors.schema_embeddings
		WHERE schema_name = $1 AND table_name = $2
		ORDER BY
			CASE object_type
				WHEN 'table' THEN 1
				WHEN 'column' THEN 2
				WHEN 'relationship' THEN 3
			END,
			column_name NULLS FIRST
	`

	rows, err := c.pool.Query(ctx, query, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get table schema: %w", err)
	}
	defer rows.Close()

	return c.scanEmbeddingsWithoutSimilarity(rows)
}

// UpsertSchemaEmbedding inserts or updates a schema embedding.
func (c *PgVectorClient) UpsertSchemaEmbedding(ctx context.Context, emb *SchemaEmbedding) error {
	if emb == nil {
		return fmt.Errorf("warehouse: embedding is required")
	}

	query := `
		INSERT INTO vectors.schema_embeddings (
			schema_name, table_name, column_name, object_type,
			description_en, description_ar, data_type, sample_values,
			business_context, embedding, model_name
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (schema_name, table_name, column_name)
		DO UPDATE SET
			description_en = EXCLUDED.description_en,
			description_ar = EXCLUDED.description_ar,
			data_type = EXCLUDED.data_type,
			sample_values = EXCLUDED.sample_values,
			business_context = EXCLUDED.business_context,
			embedding = EXCLUDED.embedding,
			model_name = EXCLUDED.model_name,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	// Build vector literal for embedding
	vectorLiteral := c.buildVectorLiteral(emb.Embedding)

	err := c.pool.QueryRow(ctx, query,
		emb.SchemaName,
		emb.TableName,
		emb.ColumnName,
		emb.ObjectType,
		emb.DescriptionEN,
		emb.DescriptionAR,
		emb.DataType,
		emb.SampleValues,
		emb.BusinessContext,
		vectorLiteral,
		emb.ModelName,
	).Scan(&emb.ID, &emb.CreatedAt, &emb.UpdatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to upsert schema embedding: %w", err)
	}

	c.logger.Debug("schema embedding upserted",
		slog.String("schema", emb.SchemaName),
		slog.String("table", emb.TableName),
		slog.String("type", emb.ObjectType),
	)

	return nil
}

// DeleteTableEmbeddings removes all embeddings for a specific table.
func (c *PgVectorClient) DeleteTableEmbeddings(ctx context.Context, schemaName, tableName string) error {
	query := `
		DELETE FROM vectors.schema_embeddings
		WHERE schema_name = $1 AND table_name = $2
	`

	result, err := c.pool.Exec(ctx, query, schemaName, tableName)
	if err != nil {
		return fmt.Errorf("warehouse: failed to delete table embeddings: %w", err)
	}

	c.logger.Debug("table embeddings deleted",
		slog.String("schema", schemaName),
		slog.String("table", tableName),
		slog.Int64("count", result.RowsAffected()),
	)

	return nil
}

// GetSchemaStats returns statistics about the schema embeddings.
func (c *PgVectorClient) GetSchemaStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total_embeddings,
			COUNT(DISTINCT schema_name) as schema_count,
			COUNT(DISTINCT table_name) as table_count,
			COUNT(*) FILTER (WHERE object_type = 'table') as table_embeddings,
			COUNT(*) FILTER (WHERE object_type = 'column') as column_embeddings,
			COUNT(*) FILTER (WHERE object_type = 'relationship') as relationship_embeddings,
			MAX(updated_at) as last_updated
		FROM vectors.schema_embeddings
	`

	var stats struct {
		TotalEmbeddings      int64
		SchemaCount          int64
		TableCount           int64
		TableEmbeddings      int64
		ColumnEmbeddings     int64
		RelationshipEmbeddings int64
		LastUpdated          time.Time
	}

	err := c.pool.QueryRow(ctx, query).Scan(
		&stats.TotalEmbeddings,
		&stats.SchemaCount,
		&stats.TableCount,
		&stats.TableEmbeddings,
		&stats.ColumnEmbeddings,
		&stats.RelationshipEmbeddings,
		&stats.LastUpdated,
	)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get schema stats: %w", err)
	}

	return map[string]interface{}{
		"total_embeddings":       stats.TotalEmbeddings,
		"schema_count":           stats.SchemaCount,
		"table_count":            stats.TableCount,
		"table_embeddings":       stats.TableEmbeddings,
		"column_embeddings":      stats.ColumnEmbeddings,
		"relationship_embeddings": stats.RelationshipEmbeddings,
		"last_updated":           stats.LastUpdated,
	}, nil
}

// buildVectorLiteral converts a float32 slice to PostgreSQL vector literal format.
func (c *PgVectorClient) buildVectorLiteral(embedding []float32) string {
	if len(embedding) == 0 {
		return "'[]'::vector"
	}

	// Build array string representation
	result := "["
	for i, v := range embedding {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", v)
	}
	result += "]"

	return fmt.Sprintf("'%s'::vector", result)
}

// parsePostgresArray parses a PostgreSQL array string into a Go slice.
func (c *PgVectorClient) parsePostgresArray(arrayStr string) []string {
	if arrayStr == "" || arrayStr == "{}" || arrayStr == "NULL" {
		return nil
	}

	// Remove surrounding braces
	if len(arrayStr) >= 2 && arrayStr[0] == '{' && arrayStr[len(arrayStr)-1] == '}' {
		arrayStr = arrayStr[1 : len(arrayStr)-1]
	}

	if arrayStr == "" {
		return nil
	}

	// Simple split by comma (doesn't handle quoted values with commas)
	// For production, consider using pgtype.Array
	var results []string
	var current string
	inQuotes := false

	for _, r := range arrayStr {
		switch r {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				if current != "" {
					results = append(results, current)
				}
				current = ""
				continue
			}
			fallthrough
		default:
			current += string(r)
		}
	}

	if current != "" {
		results = append(results, current)
	}

	return results
}

// scanEmbeddings scans rows into SchemaEmbedding slices with similarity scores.
func (c *PgVectorClient) scanEmbeddings(rows pgx.Rows) ([]SchemaEmbedding, error) {
	var results []SchemaEmbedding

	for rows.Next() {
		var emb SchemaEmbedding
		var sampleValues []byte

		err := rows.Scan(
			&emb.ID,
			&emb.SchemaName,
			&emb.TableName,
			&emb.ColumnName,
			&emb.ObjectType,
			&emb.DescriptionEN,
			&emb.DescriptionAR,
			&emb.DataType,
			&sampleValues,
			&emb.BusinessContext,
			&emb.ModelName,
			&emb.Similarity,
			&emb.CreatedAt,
			&emb.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan schema embedding: %w", err)
		}

		if len(sampleValues) > 0 {
			emb.SampleValues = c.parsePostgresArray(string(sampleValues))
		}

		results = append(results, emb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating schema embeddings: %w", err)
	}

	return results, nil
}

// scanEmbeddingsWithoutSimilarity scans rows without similarity scores.
func (c *PgVectorClient) scanEmbeddingsWithoutSimilarity(rows pgx.Rows) ([]SchemaEmbedding, error) {
	var results []SchemaEmbedding

	for rows.Next() {
		var emb SchemaEmbedding
		var sampleValues []byte

		err := rows.Scan(
			&emb.ID,
			&emb.SchemaName,
			&emb.TableName,
			&emb.ColumnName,
			&emb.ObjectType,
			&emb.DescriptionEN,
			&emb.DescriptionAR,
			&emb.DataType,
			&sampleValues,
			&emb.BusinessContext,
			&emb.ModelName,
			&emb.CreatedAt,
			&emb.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan schema embedding: %w", err)
		}

		if len(sampleValues) > 0 {
			emb.SampleValues = c.parsePostgresArray(string(sampleValues))
		}

		results = append(results, emb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating schema embeddings: %w", err)
	}

	return results, nil
}
