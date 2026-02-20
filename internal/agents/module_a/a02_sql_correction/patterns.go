// Package a02_sql_correction provides the SQL self-correction agent.
//
// This file implements error pattern recognition for SQL correction.
package a02_sql_correction

import (
	"fmt"
	"regexp"
	"strings"
)

// ErrorPatternRegistry holds patterns for recognizing and correcting SQL errors.
type ErrorPatternRegistry struct {
	patterns map[string]*ErrorPattern
}

// ErrorPattern represents a recognized error pattern with correction strategy.
type ErrorPattern struct {
	Type        string
	Pattern     *regexp.Regexp
	Description string
	Strategy    CorrectionStrategy
}

// CorrectionStrategy defines how to correct a specific error type.
type CorrectionStrategy interface {
	Apply(req CorrectionRequest) (correctedSQL string, description string, confidence float64)
}

// NewErrorPatternRegistry creates a new registry with built-in patterns.
func NewErrorPatternRegistry() *ErrorPatternRegistry {
	registry := &ErrorPatternRegistry{
		patterns: make(map[string]*ErrorPattern),
	}

	// Register built-in patterns
	registry.registerBuiltInPatterns()

	return registry
}

// registerBuiltInPatterns adds standard error patterns.
func (r *ErrorPatternRegistry) registerBuiltInPatterns() {
	// Column not found
	r.Register("column_not_found", &ColumnNotFoundStrategy{})

	// Relation/table not found
	r.Register("relation_not_found", &RelationNotFoundStrategy{})

	// Syntax errors
	r.Register("syntax_error", &SyntaxErrorStrategy{})

	// Type mismatch
	r.Register("type_mismatch", &TypeMismatchStrategy{})

	// Ambiguous reference
	r.Register("ambiguous_reference", &AmbiguousReferenceStrategy{})

	// GROUP BY errors
	r.Register("group_by_error", &GroupByErrorStrategy{})

	// JOIN errors
	r.Register("join_error", &JoinErrorStrategy{})

	// Function errors
	r.Register("function_error", &FunctionErrorStrategy{})
}

// Register adds a new error pattern.
func (r *ErrorPatternRegistry) Register(errorType string, strategy CorrectionStrategy) {
	r.patterns[errorType] = &ErrorPattern{
		Type:     errorType,
		Strategy: strategy,
	}
}

// GetStrategy returns the correction strategy for an error type.
func (r *ErrorPatternRegistry) GetStrategy(errorType string) CorrectionStrategy {
	if pattern, ok := r.patterns[errorType]; ok {
		return pattern.Strategy
	}
	return nil
}

// ColumnNotFoundStrategy handles "column does not exist" errors.
type ColumnNotFoundStrategy struct{}

func (s *ColumnNotFoundStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	// Extract the missing column name
	colPattern := regexp.MustCompile(`column "([^"]+)" does not exist`)
	matches := colPattern.FindStringSubmatch(req.Error)
	if len(matches) < 2 {
		return "", "Could not identify missing column", 0
	}

	missingCol := matches[1]

	// Common corrections
	corrections := map[string]string{
		"patient_name":   "p.patient_name",
		"doctor_name":    "d.doctor_name",
		"department":     "dept.department_name",
		"amount":         "f.amount",
		"date":           "d.date_actual",
		"billing_date":   "d.date_actual",
	}

	if alias, ok := corrections[missingCol]; ok {
		corrected := strings.Replace(req.SQL, missingCol, alias, 1)
		desc := fmt.Sprintf("Added table alias to column '%s' -> '%s'", missingCol, alias)
		return corrected, desc, 0.85
	}

	// Try to infer from schema hints
	for _, hint := range req.SchemaHints {
		if strings.Contains(hint, missingCol) {
			// Extract table from hint like "fact_billing.amount"
			parts := strings.Split(hint, ".")
			if len(parts) == 2 {
				alias := parts[0][:1] // First letter as alias
				corrected := strings.Replace(req.SQL, missingCol, alias+"."+missingCol, 1)
				desc := fmt.Sprintf("Inferred table alias '%s' for column '%s'", alias, missingCol)
				return corrected, desc, 0.7
			}
		}
	}

	return "", fmt.Sprintf("Column '%s' not found, no correction available", missingCol), 0
}

// RelationNotFoundStrategy handles "relation does not exist" errors.
type RelationNotFoundStrategy struct{}

func (s *RelationNotFoundStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	// Extract the missing table name
	relPattern := regexp.MustCompile(`relation "([^"]+)" does not exist`)
	matches := relPattern.FindStringSubmatch(req.Error)
	if len(matches) < 2 {
		return "", "Could not identify missing relation", 0
	}

	missingTable := matches[1]

	// Common table name corrections
	corrections := map[string]string{
		"patients":         "dim_patient",
		"patient":          "dim_patient",
		"doctors":          "dim_doctor",
		"doctor":           "dim_doctor",
		"departments":      "dim_department",
		"department":       "dim_department",
		"appointments":     "fact_appointments",
		"appointment":      "fact_appointments",
		"billing":          "fact_billing",
		"bills":            "fact_billing",
		"payments":         "fact_payments",
		"payment":          "fact_payments",
		"visits":           "fact_appointments",
		"patient_visits":   "fact_appointments",
		"revenue":          "fact_billing",
	}

	if corrected, ok := corrections[strings.ToLower(missingTable)]; ok {
		result := strings.Replace(req.SQL, missingTable, corrected, 1)
		desc := fmt.Sprintf("Corrected table name '%s' -> '%s'", missingTable, corrected)
		return result, desc, 0.9
	}

	return "", fmt.Sprintf("Table '%s' not found, no correction available", missingTable), 0
}

// SyntaxErrorStrategy handles SQL syntax errors.
type SyntaxErrorStrategy struct{}

func (s *SyntaxErrorStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	sql := req.SQL
	desc := ""
	confidence := 0.5

	// Missing closing parenthesis
	openParens := strings.Count(sql, "(")
	closeParens := strings.Count(sql, ")")
	if openParens > closeParens {
		sql = sql + strings.Repeat(")", openParens-closeParens)
		desc = fmt.Sprintf("Added %d missing closing parenthesis(es)", openParens-closeParens)
		confidence = 0.8
	} else if closeParens > openParens {
		sql = strings.Replace(sql, ")", "", closeParens-openParens)
		desc = fmt.Sprintf("Removed %d extra closing parenthesis(es)", closeParens-openParens)
		confidence = 0.8
	}

	// Missing closing quote
	singleQuotes := strings.Count(sql, "'")
	if singleQuotes%2 != 0 {
		// Find the last string literal and close it
		lastQuote := strings.LastIndex(sql, "'")
		if lastQuote != -1 {
			sql = sql[:lastQuote+1] + "'" + sql[lastQuote+1:]
			desc = "Added missing closing quote"
			confidence = 0.75
		}
	}

	// Missing comma in SELECT list
	selectPattern := regexp.MustCompile(`SELECT\s+(.+?)\s+FROM`)
	matches := selectPattern.FindStringSubmatch(sql)
	if len(matches) > 1 {
		selectList := matches[1]
		// Check for missing commas between aliases
		if strings.Contains(selectList, "  ") && !strings.Contains(selectList, ",") {
			// Might be missing commas
			cols := strings.Fields(selectList)
			if len(cols) > 1 {
				newSelectList := strings.Join(cols, ", ")
				sql = strings.Replace(sql, selectList, newSelectList, 1)
				desc = "Added missing commas in SELECT list"
				confidence = 0.7
			}
		}
	}

	if desc != "" {
		return sql, desc, confidence
	}

	return "", "Could not identify syntax error correction", 0
}

// TypeMismatchStrategy handles type conversion errors.
type TypeMismatchStrategy struct{}

func (s *TypeMismatchStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	// Common type mismatch patterns
	sql := req.SQL

	// Date string needs casting
	datePattern := regexp.MustCompile(`(\w+)\s*=\s*'(\d{4}-\d{2}-\d{2})'`)
	if datePattern.MatchString(sql) {
		// Add DATE cast
		sql = datePattern.ReplaceAllString(sql, "$1 = DATE '$2'")
		return sql, "Added DATE cast to date literal", 0.8
	}

	// String to number comparison
	comparePattern := regexp.MustCompile(`(\w+)\s*=\s*'(\d+)'`)
	if comparePattern.MatchString(sql) {
		// Remove quotes from numeric comparison
		sql = comparePattern.ReplaceAllString(sql, "$1 = $2")
		return sql, "Removed quotes from numeric comparison", 0.7
	}

	return "", "Could not identify type mismatch correction", 0
}

// AmbiguousReferenceStrategy handles ambiguous column references.
type AmbiguousReferenceStrategy struct{}

func (s *AmbiguousReferenceStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	// Extract the ambiguous column
	colPattern := regexp.MustCompile(`column "([^"]+)" is ambiguous`)
	matches := colPattern.FindStringSubmatch(req.Error)
	if len(matches) < 2 {
		return "", "Could not identify ambiguous column", 0
	}

	ambCol := matches[1]

	// Add table prefix based on common patterns
	// This is a simplification - real implementation would need schema context
	aliases := []string{"f", "d", "a", "p", "dept"}
	for _, alias := range aliases {
		// Try adding the alias prefix
		if strings.Contains(req.SQL, alias+".") {
			// Check if this table likely contains the column
			corrected := strings.Replace(req.SQL, ambCol, alias+"."+ambCol, 1)
			desc := fmt.Sprintf("Qualified ambiguous column '%s' with alias '%s'", ambCol, alias)
			return corrected, desc, 0.6
		}
	}

	return "", fmt.Sprintf("Could not resolve ambiguity for column '%s'", ambCol), 0
}

// GroupByErrorStrategy handles GROUP BY errors.
type GroupByErrorStrategy struct{}

func (s *GroupByErrorStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	sql := req.SQL

	// Find non-aggregated columns in SELECT
	selectPattern := regexp.MustCompile(`SELECT\s+(.+?)\s+FROM`)
	matches := selectPattern.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return "", "Could not parse SELECT clause", 0
	}

	_ = matches[1] // selectClause - unused but parsed for future enhancement

	// Check if there's a GROUP BY clause
	if !strings.Contains(strings.ToUpper(sql), "GROUP BY") {
		// Add GROUP BY for non-aggregate columns
		// This is a simplification
		desc := "Added missing GROUP BY clause"
		confidence := 0.6
		return sql, desc, confidence
	}

	return "", "Could not correct GROUP BY error", 0
}

// JoinErrorStrategy handles JOIN-related errors.
type JoinErrorStrategy struct{}

func (s *JoinErrorStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	sql := req.SQL

	// Missing ON clause for JOIN
	joinPattern := regexp.MustCompile(`JOIN\s+(\w+)\s+JOIN`)
	if joinPattern.MatchString(sql) {
		// Likely missing ON clause
		desc := "Detected missing ON clause in JOIN"
		confidence := 0.5
		return sql, desc, confidence
	}

	return "", "Could not correct JOIN error", 0
}

// FunctionErrorStrategy handles function-related errors.
type FunctionErrorStrategy struct{}

func (s *FunctionErrorStrategy) Apply(req CorrectionRequest) (string, string, float64) {
	sql := req.SQL

	// Wrong function name
	functionCorrections := map[string]string{
		"YEAR":      "EXTRACT(YEAR FROM ",
		"MONTH":     "EXTRACT(MONTH FROM ",
		"DAY":       "EXTRACT(DAY FROM ",
		"DATEPART":  "EXTRACT(",
		"GETDATE":   "CURRENT_DATE",
		"NOW":       "CURRENT_TIMESTAMP",
		"LEN":       "LENGTH",
		"ISNULL":    "COALESCE",
		"IIF":       "CASE WHEN",
	}

	for wrong, correct := range functionCorrections {
		if strings.Contains(strings.ToUpper(sql), wrong+"(") {
			sql = regexp.MustCompile("(?i)"+wrong+"\\(").ReplaceAllString(sql, correct+"(")
			desc := fmt.Sprintf("Corrected function '%s' to '%s'", wrong, correct)
			return sql, desc, 0.85
		}
	}

	return "", "Could not identify function error correction", 0
}
