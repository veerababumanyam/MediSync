// Package warehouse provides database utilities for the MediSync data warehouse.
//
// This file provides functionality to seed the domain glossary with
// healthcare and accounting terminology.
package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GlossarySeeder handles seeding domain terms into the database.
type GlossarySeeder struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewGlossarySeeder creates a new glossary seeder.
func NewGlossarySeeder(pool *pgxpool.Pool, logger *slog.Logger) *GlossarySeeder {
	if logger == nil {
		logger = slog.Default()
	}
	return &GlossarySeeder{
		pool:   pool,
		logger: logger.With("component", "glossary_seeder"),
	}
}

// GlossaryEntry represents a single term mapping.
type GlossaryEntry struct {
	ID             int                  `json:"id"`
	Synonym        string               `json:"synonym"`
	CanonicalTerm  string               `json:"canonical_term"`
	Category       string               `json:"category"`
	SQLFragment    string               `json:"sql_fragment"`
	LocaleVariants map[string][]string  `json:"locale_variants"`
	Description    string               `json:"description"`
}

// GlossaryFile represents the structure of the glossary JSON file.
type GlossaryFile struct {
	Metadata struct {
		Version     string `json:"version"`
		Created     string `json:"created"`
		Description string `json:"description"`
		TotalTerms  int    `json:"total_terms"`
		Categories  []string `json:"categories"`
	} `json:"metadata"`
	Terms []GlossaryEntry `json:"terms"`
}

// SeedFromFile loads glossary terms from a JSON file and inserts them into the database.
func (s *GlossarySeeder) SeedFromFile(ctx context.Context, filePath string) (int, error) {
	s.logger.Info("loading glossary from file", "path", filePath)

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read glossary file: %w", err)
	}

	// Parse JSON
	var glossary GlossaryFile
	if err := json.Unmarshal(data, &glossary); err != nil {
		return 0, fmt.Errorf("failed to parse glossary JSON: %w", err)
	}

	s.logger.Info("parsed glossary file", "terms", len(glossary.Terms))

	// Seed the terms
	return s.SeedTerms(ctx, glossary.Terms)
}

// SeedTerms inserts glossary terms into the database.
func (s *GlossarySeeder) SeedTerms(ctx context.Context, terms []GlossaryEntry) (int, error) {
	inserted := 0

	for _, term := range terms {
		// Convert locale variants to JSONB
		variantsJSON, err := json.Marshal(term.LocaleVariants)
		if err != nil {
			s.logger.Warn("failed to marshal locale variants", "term", term.Synonym, "error", err)
			continue
		}

		// Insert or update the term
		query := `
			INSERT INTO app.domain_terms (synonym, canonical_term, category, sql_fragment, locale_variants)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (synonym) DO UPDATE SET
				canonical_term = EXCLUDED.canonical_term,
				category = EXCLUDED.category,
				sql_fragment = EXCLUDED.sql_fragment,
				locale_variants = EXCLUDED.locale_variants
		`

		_, err = s.pool.Exec(ctx, query,
			term.Synonym,
			term.CanonicalTerm,
			term.Category,
			term.SQLFragment,
			variantsJSON,
		)

		if err != nil {
			s.logger.Warn("failed to insert term", "term", term.Synonym, "error", err)
			continue
		}

		inserted++
	}

	s.logger.Info("glossary seeding complete", "inserted", inserted, "total", len(terms))
	return inserted, nil
}

// Clear removes all glossary terms from the database.
func (s *GlossarySeeder) Clear(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM app.domain_terms")
	if err != nil {
		return fmt.Errorf("failed to clear glossary: %w", err)
	}
	s.logger.Info("glossary cleared")
	return nil
}

// GetTerm retrieves a single term by synonym.
func (s *GlossarySeeder) GetTerm(ctx context.Context, synonym string) (*GlossaryEntry, error) {
	query := `
		SELECT synonym, canonical_term, category, sql_fragment, locale_variants
		FROM app.domain_terms
		WHERE synonym = $1
	`

	var term GlossaryEntry
	var variantsJSON []byte

	err := s.pool.QueryRow(ctx, query, synonym).Scan(
		&term.Synonym,
		&term.CanonicalTerm,
		&term.Category,
		&term.SQLFragment,
		&variantsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("term not found: %w", err)
	}

	if len(variantsJSON) > 0 {
		if err := json.Unmarshal(variantsJSON, &term.LocaleVariants); err != nil {
			s.logger.Warn("failed to unmarshal locale variants", "error", err)
		}
	}

	return &term, nil
}

// GetAllTerms retrieves all glossary terms.
func (s *GlossarySeeder) GetAllTerms(ctx context.Context) ([]GlossaryEntry, error) {
	query := `
		SELECT synonym, canonical_term, category, sql_fragment, locale_variants
		FROM app.domain_terms
		ORDER BY category, synonym
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query terms: %w", err)
	}
	defer rows.Close()

	var terms []GlossaryEntry
	for rows.Next() {
		var term GlossaryEntry
		var variantsJSON []byte

		if err := rows.Scan(
			&term.Synonym,
			&term.CanonicalTerm,
			&term.Category,
			&term.SQLFragment,
			&variantsJSON,
		); err != nil {
			s.logger.Warn("failed to scan term", "error", err)
			continue
		}

		if len(variantsJSON) > 0 {
			if err := json.Unmarshal(variantsJSON, &term.LocaleVariants); err != nil {
				s.logger.Warn("failed to unmarshal locale variants", "error", err)
			}
		}

		terms = append(terms, term)
	}

	return terms, nil
}

// GetTermsByCategory retrieves terms filtered by category.
func (s *GlossarySeeder) GetTermsByCategory(ctx context.Context, category string) ([]GlossaryEntry, error) {
	query := `
		SELECT synonym, canonical_term, category, sql_fragment, locale_variants
		FROM app.domain_terms
		WHERE category = $1
		ORDER BY synonym
	`

	rows, err := s.pool.Query(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query terms by category: %w", err)
	}
	defer rows.Close()

	var terms []GlossaryEntry
	for rows.Next() {
		var term GlossaryEntry
		var variantsJSON []byte

		if err := rows.Scan(
			&term.Synonym,
			&term.CanonicalTerm,
			&term.Category,
			&term.SQLFragment,
			&variantsJSON,
		); err != nil {
			s.logger.Warn("failed to scan term", "error", err)
			continue
		}

		if len(variantsJSON) > 0 {
			if err := json.Unmarshal(variantsJSON, &term.LocaleVariants); err != nil {
				s.logger.Warn("failed to unmarshal locale variants", "error", err)
			}
		}

		terms = append(terms, term)
	}

	return terms, nil
}

// Stats returns statistics about the glossary.
func (s *GlossarySeeder) Stats(ctx context.Context) (*GlossaryStats, error) {
	stats := &GlossaryStats{}

	// Total count
	err := s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM app.domain_terms").Scan(&stats.TotalTerms)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Count by category
	rows, err := s.pool.Query(ctx, `
		SELECT category, COUNT(*)
		FROM app.domain_terms
		GROUP BY category
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get category counts: %w", err)
	}
	defer rows.Close()

	stats.ByCategory = make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			continue
		}
		stats.ByCategory[category] = count
	}

	return stats, nil
}

// GlossaryStats contains statistics about the glossary.
type GlossaryStats struct {
	TotalTerms int            `json:"total_terms"`
	ByCategory map[string]int `json:"by_category"`
}
