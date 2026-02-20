// Package warehouse provides database utilities for the MediSync data warehouse.
//
// This file provides functionality to seed schema embeddings for AI context.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

// EmbeddingSeeder handles seeding schema embeddings into the database.
type EmbeddingSeeder struct {
	pool     *pgxpool.Pool
	embedder Embedder
	logger   *slog.Logger
}

// Embedder creates embeddings from text.
type Embedder interface {
	Embed(ctx context.Context, text string) (pgvector.Vector, error)
}

// NewEmbeddingSeeder creates a new embedding seeder.
func NewEmbeddingSeeder(pool *pgxpool.Pool, embedder Embedder, logger *slog.Logger) *EmbeddingSeeder {
	if logger == nil {
		logger = slog.Default()
	}
	return &EmbeddingSeeder{
		pool:     pool,
		embedder: embedder,
		logger:   logger.With("component", "embedding_seeder"),
	}
}

// SchemaDefinition represents a table or column definition for embedding.
type SchemaDefinition struct {
	Type        string                 `json:"type"`        // "table", "column", "query_pattern"
	Name        string                 `json:"name"`        // Table/column name
	Description string                 `json:"description"` // Natural language description
	Metadata    map[string]interface{} `json:"metadata"`    // Additional context
}

// SeedTables seeds embeddings for warehouse tables.
func (s *EmbeddingSeeder) SeedTables(ctx context.Context, tables []SchemaDefinition) (int, error) {
	s.logger.Info("seeding table embeddings", "count", len(tables))

	inserted := 0
	for _, table := range tables {
		if table.Type != "table" {
			continue
		}

		// Create embedding from description
		embedding, err := s.embedder.Embed(ctx, table.Description)
		if err != nil {
			s.logger.Warn("failed to embed table", "table", table.Name, "error", err)
			continue
		}

		// Serialize metadata
		metadataJSON, err := json.Marshal(table.Metadata)
		if err != nil {
			s.logger.Warn("failed to marshal metadata", "table", table.Name, "error", err)
			continue
		}

		// Insert into database
		query := `
			INSERT INTO vectors.schema_embeddings (embedding_type, entity_name, description, metadata, embedding)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (entity_name) WHERE embedding_type = 'table' DO UPDATE SET
				description = EXCLUDED.description,
				metadata = EXCLUDED.metadata,
				embedding = EXCLUDED.embedding,
				updated_at = NOW()
		`

		_, err = s.pool.Exec(ctx, query,
			"table",
			table.Name,
			table.Description,
			metadataJSON,
			embedding,
		)

		if err != nil {
			s.logger.Warn("failed to insert table embedding", "table", table.Name, "error", err)
			continue
		}

		inserted++
	}

	s.logger.Info("table embeddings seeded", "inserted", inserted, "total", len(tables))
	return inserted, nil
}

// SeedQueryPatterns seeds embeddings for common query patterns.
func (s *EmbeddingSeeder) SeedQueryPatterns(ctx context.Context, patterns []SchemaDefinition) (int, error) {
	s.logger.Info("seeding query pattern embeddings", "count", len(patterns))

	inserted := 0
	for _, pattern := range patterns {
		if pattern.Type != "query_pattern" {
			continue
		}

		// Create embedding from description
		embedding, err := s.embedder.Embed(ctx, pattern.Description)
		if err != nil {
			s.logger.Warn("failed to embed pattern", "pattern", pattern.Name, "error", err)
			continue
		}

		// Serialize metadata
		metadataJSON, err := json.Marshal(pattern.Metadata)
		if err != nil {
			s.logger.Warn("failed to marshal metadata", "pattern", pattern.Name, "error", err)
			continue
		}

		// Insert into database
		query := `
			INSERT INTO vectors.schema_embeddings (embedding_type, entity_name, description, metadata, embedding)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (entity_name) WHERE embedding_type = 'query_pattern' DO UPDATE SET
				description = EXCLUDED.description,
				metadata = EXCLUDED.metadata,
				embedding = EXCLUDED.embedding,
				updated_at = NOW()
		`

		_, err = s.pool.Exec(ctx, query,
			"query_pattern",
			pattern.Name,
			pattern.Description,
			metadataJSON,
			embedding,
		)

		if err != nil {
			s.logger.Warn("failed to insert pattern embedding", "pattern", pattern.Name, "error", err)
			continue
		}

		inserted++
	}

	s.logger.Info("query pattern embeddings seeded", "inserted", inserted, "total", len(patterns))
	return inserted, nil
}

// SeedFromFile loads schema definitions from a JSON file and seeds them.
func (s *EmbeddingSeeder) SeedFromFile(ctx context.Context, filePath string) (int, error) {
	s.logger.Info("loading schema definitions from file", "path", filePath)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Parse JSON
	var definitions struct {
		Tables        []SchemaDefinition `json:"tables"`
		QueryPatterns []SchemaDefinition `json:"query_patterns"`
	}
	if err := json.Unmarshal(data, &definitions); err != nil {
		return 0, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	// Seed tables
	tableCount, err := s.SeedTables(ctx, definitions.Tables)
	if err != nil {
		return tableCount, err
	}

	// Seed query patterns
	patternCount, err := s.SeedQueryPatterns(ctx, definitions.QueryPatterns)
	if err != nil {
		return tableCount + patternCount, err
	}

	return tableCount + patternCount, nil
}

// GetEmbedding retrieves an embedding by entity name.
func (s *EmbeddingSeeder) GetEmbedding(ctx context.Context, embeddingType, entityName string) (pgvector.Vector, error) {
	query := `
		SELECT embedding FROM vectors.schema_embeddings
		WHERE embedding_type = $1 AND entity_name = $2
	`

	var embedding pgvector.Vector
	err := s.pool.QueryRow(ctx, query, embeddingType, entityName).Scan(&embedding)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("embedding not found: %w", err)
	}

	return embedding, nil
}

// SearchSimilar finds similar embeddings using cosine similarity.
func (s *EmbeddingSeeder) SearchSimilar(ctx context.Context, embedding pgvector.Vector, limit int, embeddingType string) ([]SearchResult, error) {
	query := `
		SELECT entity_name, description, metadata, 1 - (embedding <=> $1) as similarity
		FROM vectors.schema_embeddings
		WHERE ($2 = '' OR embedding_type = $2)
		ORDER BY embedding <=> $1
		LIMIT $3
	`

	rows, err := s.pool.Query(ctx, query, embedding, embeddingType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search embeddings: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var result SearchResult
		var metadataJSON []byte

		if err := rows.Scan(&result.EntityName, &result.Description, &metadataJSON, &result.Similarity); err != nil {
			s.logger.Warn("failed to scan search result", "error", err)
			continue
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &result.Metadata)
		}

		results = append(results, result)
	}

	return results, nil
}

// SearchResult represents a similarity search result.
type SearchResult struct {
	EntityName  string                 `json:"entity_name"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	Similarity  float64                `json:"similarity"`
}

// ClearAll removes all embeddings from the database.
func (s *EmbeddingSeeder) ClearAll(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM vectors.schema_embeddings")
	if err != nil {
		return fmt.Errorf("failed to clear embeddings: %w", err)
	}
	s.logger.Info("all embeddings cleared")
	return nil
}

// GetStats returns statistics about stored embeddings.
func (s *EmbeddingSeeder) GetStats(ctx context.Context) (*EmbeddingStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE embedding_type = 'table') as tables,
			COUNT(*) FILTER (WHERE embedding_type = 'column') as columns,
			COUNT(*) FILTER (WHERE embedding_type = 'query_pattern') as patterns
		FROM vectors.schema_embeddings
	`

	var stats EmbeddingStats
	err := s.pool.QueryRow(ctx, query).Scan(
		&stats.Total,
		&stats.Tables,
		&stats.Columns,
		&stats.Patterns,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get embedding stats: %w", err)
	}

	return &stats, nil
}

// EmbeddingStats contains statistics about stored embeddings.
type EmbeddingStats struct {
	Total    int `json:"total"`
	Tables   int `json:"tables"`
	Columns  int `json:"columns"`
	Patterns int `json:"patterns"`
}
