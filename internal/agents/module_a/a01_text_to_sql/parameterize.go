// Package a01_text_to_sql provides the text-to-SQL agent subcomponents.
//
// This file implements SQL parameterization for injection prevention.
package a01_text_to_sql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parameterizer converts raw SQL values to parameterized queries.
type Parameterizer struct {
	paramPrefix string
}

// ParameterizerConfig holds configuration for the parameterizer.
type ParameterizerConfig struct {
	ParamPrefix string // e.g., "$" for PostgreSQL or "?" for MySQL
}

// NewParameterizer creates a new SQL parameterizer.
func NewParameterizer(cfg ParameterizerConfig) *Parameterizer {
	if cfg.ParamPrefix == "" {
		cfg.ParamPrefix = "$"
	}
	return &Parameterizer{
		paramPrefix: cfg.ParamPrefix,
	}
}

// ParameterizedQuery contains a parameterized SQL statement with its values.
type ParameterizedQuery struct {
	SQL        string        `json:"sql"`
	Parameters []interface{} `json:"parameters"`
	IsSafe     bool          `json:"is_safe"`
	Warnings   []string      `json:"warnings,omitempty"`
}

// Parameterize converts literal values in SQL to parameters.
func (p *Parameterizer) Parameterize(sql string) *ParameterizedQuery {
	result := &ParameterizedQuery{
		SQL:        sql,
		Parameters: []interface{}{},
		IsSafe:     true,
		Warnings:   []string{},
	}

	// Extract and replace string literals
	result.SQL = p.replaceStringLiterals(result)

	// Extract and replace numeric literals
	result.SQL = p.replaceNumericLiterals(result)

	// Extract and replace date/timestamp literals
	result.SQL = p.replaceDateLiterals(result)

	// Check for dangerous patterns
	result.Warnings = p.checkDangerousPatterns(result.SQL)
	if len(result.Warnings) > 0 {
		result.IsSafe = false
	}

	return result
}

// replaceStringLiterals finds string literals and replaces them with parameters.
func (p *Parameterizer) replaceStringLiterals(result *ParameterizedQuery) string {
	sql := result.SQL

	// Match string literals in single quotes, handling escaped quotes
	stringPattern := regexp.MustCompile(`'(?:[^'\\]|\\.)*'`)

	paramIndex := len(result.Parameters) + 1

	matches := stringPattern.FindAllString(sql, -1)
	for _, match := range matches {
		// Remove surrounding quotes and unescape
		value := match[1 : len(match)-1]
		value = strings.ReplaceAll(value, "''", "'")
		value = strings.ReplaceAll(value, "\\'", "'")
		value = strings.ReplaceAll(value, "\\\\", "\\")

		result.Parameters = append(result.Parameters, value)
		placeholder := fmt.Sprintf("%s%d", p.paramPrefix, paramIndex)
		sql = strings.Replace(sql, match, placeholder, 1)
		paramIndex++
	}

	return sql
}

// replaceNumericLiterals finds numeric literals and replaces them with parameters.
func (p *Parameterizer) replaceNumericLiterals(result *ParameterizedQuery) string {
	sql := result.SQL

	// Match numeric literals (integers and decimals)
	// Exclude numbers that are part of identifiers (e.g., table_2026)
	numericPattern := regexp.MustCompile(`\b(\d+(?:\.\d+)?)\b`)

	paramIndex := len(result.Parameters) + 1

	// Find all matches first to preserve order
	matches := numericPattern.FindAllStringSubmatchIndex(sql, -1)

	// Process in reverse order to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		fullMatch := sql[match[0]:match[1]]

		// Skip if this looks like part of an identifier
		if match[0] > 0 {
			prevChar := sql[match[0]-1]
			if prevChar == '_' || (prevChar >= 'a' && prevChar <= 'z') || (prevChar >= 'A' && prevChar <= 'Z') {
				continue
			}
		}

		// Parse the number
		var value interface{}
		if strings.Contains(fullMatch, ".") {
			f, err := strconv.ParseFloat(fullMatch, 64)
			if err != nil {
				continue
			}
			value = f
		} else {
			n, err := strconv.ParseInt(fullMatch, 10, 64)
			if err != nil {
				continue
			}
			value = n
		}

		result.Parameters = append(result.Parameters, value)
		placeholder := fmt.Sprintf("%s%d", p.paramPrefix, paramIndex)
		sql = sql[:match[0]] + placeholder + sql[match[1]:]
		paramIndex++
	}

	return sql
}

// replaceDateLiterals finds date/timestamp literals and replaces them with parameters.
func (p *Parameterizer) replaceDateLiterals(result *ParameterizedQuery) string {
	sql := result.SQL

	// Match DATE 'YYYY-MM-DD' or TIMESTAMP 'YYYY-MM-DD HH:MM:SS' patterns
	datePattern := regexp.MustCompile(`(?i)(DATE|TIMESTAMP)\s+'([^']+)'`)

	paramIndex := len(result.Parameters) + 1

	matches := datePattern.FindAllStringSubmatchIndex(sql, -1)

	// Process in reverse order
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		dateType := sql[match[2]:match[3]]
		dateValue := sql[match[4]:match[5]]

		var value interface{}
		var err error

		if strings.ToUpper(dateType) == "TIMESTAMP" {
			value, err = time.Parse("2006-01-02 15:04:05", dateValue)
			if err != nil {
				value, err = time.Parse(time.RFC3339, dateValue)
			}
		} else {
			value, err = time.Parse("2006-01-02", dateValue)
		}

		if err != nil {
			// If we can't parse it, keep it as a string parameter
			value = dateValue
		}

		result.Parameters = append(result.Parameters, value)
		placeholder := fmt.Sprintf("%s%d", p.paramPrefix, paramIndex)
		fullMatch := sql[match[0]:match[1]]
		sql = strings.Replace(sql, fullMatch, placeholder, 1)
		paramIndex++
	}

	return sql
}

// checkDangerousPatterns checks for SQL injection patterns.
func (p *Parameterizer) checkDangerousPatterns(sql string) []string {
	warnings := []string{}

	// Convert to lowercase for checking
	sqlLower := strings.ToLower(sql)

	// Check for common injection patterns
	dangerousPatterns := []struct {
		pattern  string
		message  string
	}{
		{`;\\s*drop\\s+`, "Possible DROP statement detected"},
		{`;\\s*delete\\s+`, "Possible DELETE statement detected"},
		{`;\\s*truncate\\s+`, "Possible TRUNCATE statement detected"},
		{`;\\s*insert\\s+`, "Possible INSERT statement detected"},
		{`;\\s*update\\s+`, "Possible UPDATE statement detected"},
		{`;\\s*alter\\s+`, "Possible ALTER statement detected"},
		{`;\\s*create\\s+`, "Possible CREATE statement detected"},
		{`--\\s*$`, "Possible SQL comment at end of statement"},
		{`/\\*.*\\*/`, "Possible SQL comment block"},
		{`union\\s+select`, "Possible UNION-based injection"},
		{`or\\s+1\\s*=\\s*1`, "Possible boolean injection"},
		{`or\\s+''\\s*=\\s*'`, "Possible boolean injection"},
		{`exec\\s*\\(`, "Possible EXEC call"},
		{`xp_cmdshell`, "Possible xp_cmdshell call"},
	}

	for _, dp := range dangerousPatterns {
		matched, _ := regexp.MatchString(dp.pattern, sqlLower)
		if matched {
			warnings = append(warnings, dp.message)
		}
	}

	return warnings
}

// ValidateForReadOnly checks if a query is safe for read-only execution.
func (p *Parameterizer) ValidateForReadOnly(sql string) error {
	sqlUpper := strings.ToUpper(sql)

	// List of forbidden keywords for read-only queries
	forbiddenKeywords := []string{
		"INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"TRUNCATE", "GRANT", "REVOKE", "EXEC", "EXECUTE",
	}

	for _, keyword := range forbiddenKeywords {
		// Use word boundary matching
		pattern := regexp.MustCompile(`\b` + keyword + `\b`)
		if pattern.MatchString(sqlUpper) {
			return fmt.Errorf("query contains forbidden keyword: %s (read-only queries only)", keyword)
		}
	}

	// Check for multiple statements (potential injection)
	if strings.Contains(sql, ";") {
		// Allow semicolons only at the end
		trimmed := strings.TrimSpace(sql)
		if !strings.HasSuffix(trimmed, ";") || strings.Count(trimmed, ";") > 1 {
			return fmt.Errorf("multiple SQL statements not allowed")
		}
	}

	return nil
}

// SanitizeIdentifier sanitizes a table or column name.
func (p *Parameterizer) SanitizeIdentifier(identifier string) (string, error) {
	// Only allow alphanumeric and underscore
	validPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !validPattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid identifier: %s", identifier)
	}

	// Check against reserved words
	reservedWords := map[string]bool{
		"SELECT": true, "FROM": true, "WHERE": true, "JOIN": true,
		"INNER": true, "OUTER": true, "LEFT": true, "RIGHT": true,
		"ON": true, "AND": true, "OR": true, "NOT": true,
		"IN": true, "LIKE": true, "BETWEEN": true, "IS": true,
		"NULL": true, "TRUE": true, "FALSE": true, "AS": true,
		"ORDER": true, "BY": true, "GROUP": true, "HAVING": true,
		"LIMIT": true, "OFFSET": true, "UNION": true, "ALL": true,
		"DISTINCT": true, "COUNT": true, "SUM": true, "AVG": true,
		"MIN": true, "MAX": true, "CASE": true, "WHEN": true,
		"THEN": true, "ELSE": true, "END": true,
	}

	if reservedWords[strings.ToUpper(identifier)] {
		return "", fmt.Errorf("identifier is a reserved word: %s", identifier)
	}

	return identifier, nil
}

// ExtractParameters extracts parameters from a pre-parameterized query.
func (p *Parameterizer) ExtractParameters(sql string) []string {
	paramPattern := regexp.MustCompile(`\$\d+`)
	return paramPattern.FindAllString(sql, -1)
}

// RebuildQuery reconstructs a query with inline parameters (for logging/debugging).
func (p *Parameterizer) RebuildQuery(sql string, params []interface{}) string {
	result := sql
	for i, param := range params {
		placeholder := fmt.Sprintf("%s%d", p.paramPrefix, i+1)
		var replacement string
		switch v := param.(type) {
		case string:
			replacement = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
		case time.Time:
			replacement = fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
		default:
			replacement = fmt.Sprintf("%v", v)
		}
		result = strings.Replace(result, placeholder, replacement, 1)
	}
	return result
}
