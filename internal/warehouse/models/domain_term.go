package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Domain term category constants.
const (
	// CategoryHealthcare indicates a healthcare-related domain term
	CategoryHealthcare = "healthcare"
	// CategoryAccounting indicates an accounting-related domain term
	CategoryAccounting = "accounting"
	// CategoryGeneral indicates a general domain term
	CategoryGeneral = "general"
)

// ValidCategories contains all supported domain term categories.
var ValidCategories = map[string]bool{
	CategoryHealthcare: true,
	CategoryAccounting: true,
	CategoryGeneral:    true,
}

// DomainTerm represents a mapping between a synonym (user-facing term) and its
// canonical database representation. This enables the AI agent to understand
// domain-specific terminology used in healthcare and accounting contexts.
type DomainTerm struct {
	// ID is the unique identifier for the domain term
	ID int `json:"id" db:"id"`
	// Synonym is the user-facing term that might appear in natural language queries
	Synonym string `json:"synonym" db:"synonym"`
	// CanonicalTerm is the standardized term used in the database schema
	CanonicalTerm string `json:"canonical_term" db:"canonical_term"`
	// Category indicates which domain this term belongs to
	Category string `json:"category" db:"category"`
	// SQLFragment is the SQL snippet that represents this term in queries
	SQLFragment string `json:"sql_fragment" db:"sql_fragment"`
	// LocaleVariants contains locale-specific synonyms mapped by language code
	LocaleVariants map[string][]string `json:"locale_variants" db:"locale_variants"`
	// CreatedAt is the timestamp when the term was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate checks if the DomainTerm has valid field values.
// It ensures required fields are populated and category is valid.
func (t *DomainTerm) Validate() error {
	if t.Synonym == "" {
		return errors.New("synonym is required and cannot be empty")
	}

	if t.CanonicalTerm == "" {
		return errors.New("canonical_term is required and cannot be empty")
	}

	if !ValidCategories[t.Category] {
		return fmt.Errorf("invalid category '%s': must be one of 'healthcare', 'accounting', or 'general'", t.Category)
	}

	// Validate locale variants if present
	for locale := range t.LocaleVariants {
		if !ValidLocales[locale] {
			return fmt.Errorf("invalid locale '%s' in locale_variants: must be one of 'en' or 'ar'", locale)
		}
	}

	return nil
}

// GetSynonyms returns the list of synonyms for a specific locale.
// If the locale has no variants, it returns the main synonym in a slice.
func (t *DomainTerm) GetSynonyms(locale string) []string {
	if variants, ok := t.LocaleVariants[locale]; ok && len(variants) > 0 {
		return variants
	}
	// Return the main synonym as fallback
	if t.Synonym != "" {
		return []string{t.Synonym}
	}
	return []string{}
}

// Matches performs a case-insensitive check if the given text matches
// the domain term's synonym or any of its locale-specific variants.
func (t *DomainTerm) Matches(text string) bool {
	lowerText := strings.ToLower(text)
	lowerSynonym := strings.ToLower(t.Synonym)

	// Check main synonym
	if lowerText == lowerSynonym {
		return true
	}

	// Check if text contains the synonym
	if strings.Contains(lowerText, lowerSynonym) {
		return true
	}

	// Check locale variants
	for _, variants := range t.LocaleVariants {
		for _, variant := range variants {
			lowerVariant := strings.ToLower(variant)
			if lowerText == lowerVariant || strings.Contains(lowerText, lowerVariant) {
				return true
			}
		}
	}

	return false
}

// NewDomainTerm creates a new DomainTerm with the provided parameters.
// It initializes empty locale variants and sets the creation timestamp.
func NewDomainTerm(synonym, canonicalTerm, category, sqlFragment string) *DomainTerm {
	return &DomainTerm{
		ID:             0, // ID will be set by database
		Synonym:        synonym,
		CanonicalTerm:  canonicalTerm,
		Category:       category,
		SQLFragment:    sqlFragment,
		LocaleVariants: make(map[string][]string),
		CreatedAt:      time.Now(),
	}
}

// AddLocaleVariant adds a locale-specific synonym to the term.
func (t *DomainTerm) AddLocaleVariant(locale, variant string) {
	if t.LocaleVariants == nil {
		t.LocaleVariants = make(map[string][]string)
	}
	t.LocaleVariants[locale] = append(t.LocaleVariants[locale], variant)
}
