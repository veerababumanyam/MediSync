// Package a04_terminology provides the domain terminology normalizer agent.
//
// This agent maps user vocabulary (synonyms, colloquialisms) to canonical
// database terminology, supporting both English and Arabic locales.
//
// Usage:
//
//	agent := a04_terminology.New(glossaryRepo)
//	result, err := agent.Normalize(ctx, query, locale)
package a04_terminology

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"github.com/medisync/medisync/internal/agents/shared"
)

// AgentID is the unique identifier for this agent.
const AgentID = "a-04-terminology"

// GlossaryEntry represents a single term mapping in the glossary.
type GlossaryEntry struct {
	ID             int               `json:"id"`
	Synonym        string            `json:"synonym"`
	CanonicalTerm  string            `json:"canonical_term"`
	Category       string            `json:"category"`
	SQLFragment    string            `json:"sql_fragment"`
	LocaleVariants map[string][]string `json:"locale_variants"`
	Description    string            `json:"description"`
}

// GlossaryRepository provides access to the domain glossary.
type GlossaryRepository interface {
	// GetAll returns all glossary entries.
	GetAll(ctx context.Context) ([]GlossaryEntry, error)
	// GetBySynonym returns entries matching a synonym.
	GetBySynonym(ctx context.Context, synonym string) (*GlossaryEntry, error)
	// GetByCategory returns entries in a category.
	GetByCategory(ctx context.Context, category string) ([]GlossaryEntry, error)
}

// NormalizationResult contains the result of terminology normalization.
type NormalizationResult struct {
	// OriginalQuery is the user's original query text.
	OriginalQuery string `json:"original_query"`
	// NormalizedQuery is the query with canonical terms.
	NormalizedQuery string `json:"normalized_query"`
	// AppliedMappings lists the term mappings that were applied.
	AppliedMappings []TermMapping `json:"applied_mappings"`
	// Confidence indicates how confident the normalization is.
	Confidence float64 `json:"confidence"`
	// Locale is the detected/specified locale.
	Locale string `json:"locale"`
}

// TermMapping represents a single term substitution.
type TermMapping struct {
	// Original is the user's term.
	Original string `json:"original"`
	// Canonical is the canonical database term.
	Canonical string `json:"canonical"`
	// SQLFragment is the SQL mapping hint.
	SQLFragment string `json:"sql_fragment,omitempty"`
	// Category is the term category.
	Category string `json:"category"`
}

// Agent implements the domain terminology normalizer.
type Agent struct {
	id       string
	glossary GlossaryRepository
	cache    map[string]*GlossaryEntry
	cacheMu  sync.RWMutex
	logger   *slog.Logger
}

// AgentConfig holds configuration for the agent.
type AgentConfig struct {
	Glossary GlossaryRepository
	Logger   *slog.Logger
}

// New creates a new terminology normalizer agent.
func New(cfg AgentConfig) *Agent {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &Agent{
		id:       AgentID,
		glossary: cfg.Glossary,
		cache:    make(map[string]*GlossaryEntry),
		logger:   cfg.Logger.With("agent", AgentID),
	}
}

// AgentCard returns the ADK agent card for discovery.
func (a *Agent) AgentCard() shared.AgentCard {
	return shared.AgentCard{
		ID:          AgentID,
		Name:        "Domain Terminology Normalizer",
		Description: "Maps user vocabulary to canonical database terminology for healthcare and accounting domains",
		Capabilities: []string{
			"terminology-normalization",
			"synonym-resolution",
			"i18n-terminology",
		},
		Version: "1.0.0",
	}
}

// NormalizeRequest is the input for normalization.
type NormalizeRequest struct {
	Query string `json:"query"`
	Locale string `json:"locale,omitempty"`
}

// NormalizeResponse is the output from normalization.
type NormalizeResponse struct {
	Result NormalizationResult `json:"result"`
	Error  string              `json:"error,omitempty"`
}

// Normalize applies terminology normalization to a query.
func (a *Agent) Normalize(ctx context.Context, req NormalizeRequest) (*NormalizationResult, error) {
	a.logger.Debug("normalizing query", "query", req.Query, "locale", req.Locale)

	// Load glossary if not cached
	if err := a.ensureGlossaryLoaded(ctx); err != nil {
		return nil, fmt.Errorf("failed to load glossary: %w", err)
	}

	result := &NormalizationResult{
		OriginalQuery:   req.Query,
		NormalizedQuery: req.Query,
		AppliedMappings: []TermMapping{},
		Confidence:      1.0,
		Locale:          req.Locale,
	}

	if req.Locale == "" {
		result.Locale = "en"
	}

	// Apply normalization for all locale variants
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	for _, entry := range a.cache {
		// Check all variants for this term
		variants := a.getVariantsForLocale(entry, result.Locale)
		for _, variant := range variants {
			if strings.Contains(strings.ToLower(result.NormalizedQuery), strings.ToLower(variant)) {
				// Apply the mapping
				mapping := TermMapping{
					Original:    variant,
					Canonical:   entry.CanonicalTerm,
					SQLFragment: entry.SQLFragment,
					Category:    entry.Category,
				}
				result.AppliedMappings = append(result.AppliedMappings, mapping)

				// Log the mapping
				a.logger.Debug("applied term mapping",
					"original", variant,
					"canonical", entry.CanonicalTerm,
					"category", entry.Category)
			}
		}
	}

	// Adjust confidence based on number of mappings
	if len(result.AppliedMappings) > 0 {
		// Each mapping slightly reduces confidence due to potential ambiguity
		result.Confidence = 1.0 - (float64(len(result.AppliedMappings)) * 0.02)
		if result.Confidence < 0.7 {
			result.Confidence = 0.7
		}
	}

	a.logger.Info("normalization complete",
		"mappings_count", len(result.AppliedMappings),
		"confidence", result.Confidence)

	return result, nil
}

// getVariantsForLocale returns the term variants for a specific locale.
func (a *Agent) getVariantsForLocale(entry *GlossaryEntry, locale string) []string {
	variants := []string{entry.Synonym}
	if localeVariants, ok := entry.LocaleVariants[locale]; ok {
		variants = append(variants, localeVariants...)
	}
	return variants
}

// ensureGlossaryLoaded loads the glossary into cache if not already loaded.
func (a *Agent) ensureGlossaryLoaded(ctx context.Context) error {
	a.cacheMu.RLock()
	if len(a.cache) > 0 {
		a.cacheMu.RUnlock()
		return nil
	}
	a.cacheMu.RUnlock()

	// Load from repository
	entries, err := a.glossary.GetAll(ctx)
	if err != nil {
		return err
	}

	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()

	// Build cache with all variants as keys
	for i := range entries {
		entry := &entries[i]
		// Cache by primary synonym
		a.cache[strings.ToLower(entry.Synonym)] = entry

		// Cache by all locale variants
		for locale, variants := range entry.LocaleVariants {
			for _, variant := range variants {
				key := fmt.Sprintf("%s:%s", locale, strings.ToLower(variant))
				a.cache[key] = entry
			}
		}
	}

	a.logger.Debug("glossary loaded", "entries", len(entries))
	return nil
}

// GetSQLHints returns SQL fragment hints for terms found in the query.
func (a *Agent) GetSQLHints(ctx context.Context, query string) ([]string, error) {
	if err := a.ensureGlossaryLoaded(ctx); err != nil {
		return nil, err
	}

	hints := []string{}
	queryLower := strings.ToLower(query)

	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	for key, entry := range a.cache {
		if !strings.Contains(key, ":") && strings.Contains(queryLower, key) {
			if entry.SQLFragment != "" {
				hints = append(hints, entry.SQLFragment)
			}
		}
	}

	return hints, nil
}

// ExtractDomainContext extracts domain context from a query for SQL generation.
func (a *Agent) ExtractDomainContext(ctx context.Context, query string) (*DomainContext, error) {
	if err := a.ensureGlossaryLoaded(ctx); err != nil {
		return nil, err
	}

	context := &DomainContext{
		HealthcareTerms: []string{},
		AccountingTerms: []string{},
		GeneralTerms:    []string{},
		Tables:          []string{},
		Columns:         []string{},
	}

	queryLower := strings.ToLower(query)
	seenTables := make(map[string]bool)
	seenColumns := make(map[string]bool)

	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()

	for key, entry := range a.cache {
		if strings.Contains(key, ":") {
			continue // Skip locale-prefixed keys
		}
		if strings.Contains(queryLower, key) {
			switch entry.Category {
			case "healthcare":
				context.HealthcareTerms = append(context.HealthcareTerms, entry.CanonicalTerm)
			case "accounting":
				context.AccountingTerms = append(context.AccountingTerms, entry.CanonicalTerm)
			case "general":
				context.GeneralTerms = append(context.GeneralTerms, entry.CanonicalTerm)
			}

			// Extract table/column hints from SQL fragment
			tables, columns := parseSQLFragment(entry.SQLFragment)
			for _, t := range tables {
				if !seenTables[t] {
					seenTables[t] = true
					context.Tables = append(context.Tables, t)
				}
			}
			for _, c := range columns {
				if !seenColumns[c] {
					seenColumns[c] = true
					context.Columns = append(context.Columns, c)
				}
			}
		}
	}

	return context, nil
}

// DomainContext contains extracted domain information for query context.
type DomainContext struct {
	HealthcareTerms []string `json:"healthcare_terms"`
	AccountingTerms []string `json:"accounting_terms"`
	GeneralTerms    []string `json:"general_terms"`
	Tables          []string `json:"tables"`
	Columns         []string `json:"columns"`
}

// parseSQLFragment extracts table and column names from an SQL fragment.
func parseSQLFragment(fragment string) (tables []string, columns []string) {
	// Simple regex patterns for table.column references
	tablePattern := regexp.MustCompile(`(\w+)\.(\w+)`)
	matches := tablePattern.FindAllStringSubmatch(fragment, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			tables = append(tables, match[1])
		}
		if len(match) >= 3 {
			columns = append(columns, match[2])
		}
	}
	return tables, columns
}

// ToJSON serializes the normalization result.
func (r *NormalizationResult) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// FromJSON deserializes a normalization result.
func FromJSON(data string) (*NormalizationResult, error) {
	var result NormalizationResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
