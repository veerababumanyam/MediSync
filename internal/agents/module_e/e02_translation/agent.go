// Package e02_translation provides the E-02 Query Translation Agent.
//
// This agent translates Arabic queries to English intent for SQL generation,
// while preserving domain terminology and ensuring SQL-friendly output.
//
// Translation Approach:
// 1. Domain term extraction and preservation
// 2. Intent translation to English
// 3. SQL-friendly output normalization
//
// Usage:
//
//	agent := e02_translation.New(llmClient, logger)
//	result, err := agent.Translate(ctx, "أظهر إيرادات العيادة", "ar")
//	// result.TranslatedText == "Show clinic revenue"
package e02_translation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

const (
	// AgentID is the identifier for this agent.
	AgentID = "E-02"

	// AgentName is the human-readable name.
	AgentName = "Query Translation Agent"

	// DefaultTargetLocale is the default target locale for translation.
	DefaultTargetLocale = "en"
)

// LLMClient is the interface for LLM operations.
type LLMClient interface {
	// Generate generates text from a prompt.
	Generate(ctx context.Context, prompt string) (string, error)

	// GenerateWithSystem generates text with a system prompt.
	GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// TranslationResult represents the result of query translation.
type TranslationResult struct {
	// OriginalText is the original input text.
	OriginalText string `json:"original_text"`

	// TranslatedText is the translated English intent.
	TranslatedText string `json:"translated_text"`

	// SourceLocale is the source language locale.
	SourceLocale string `json:"source_locale"`

	// TargetLocale is the target language locale.
	TargetLocale string `json:"target_locale"`

	// PreservedTerms are domain terms that were preserved unchanged.
	PreservedTerms []string `json:"preserved_terms,omitempty"`

	// Confidence is the translation confidence (0.0-1.0).
	Confidence float64 `json:"confidence"`

	// RequiresHumanReview indicates if translation quality is uncertain.
	RequiresHumanReview bool `json:"requires_human_review,omitempty"`
}

// QueryTranslationAgent translates queries to English intent.
type QueryTranslationAgent struct {
	llm    LLMClient
	logger *slog.Logger

	// Arabic to English domain term mappings
	arabicToEnglish map[string]string

	// English to canonical term mappings (normalization)
	englishNormalization map[string]string

	// Terms that should never be translated
	preservedTerms map[string]bool
}

// Config holds configuration for the agent.
type Config struct {
	// LLM is the LLM client for translation.
	LLM LLMClient

	// Logger is the structured logger.
	Logger *slog.Logger
}

// New creates a new E-02 Query Translation Agent.
func New(cfg *Config) *QueryTranslationAgent {
	logger := slog.Default()
	if cfg != nil && cfg.Logger != nil {
		logger = cfg.Logger
	}

	// Initialize Arabic to English domain mappings
	arabicToEnglish := map[string]string{
		// Healthcare terms
		"العيادة":     "clinic",
		"المرضى":      "patients",
		"المريض":      "patient",
		"الأطباء":     "doctors",
		"الطبيب":      "doctor",
		"المواعيد":    "appointments",
		"الموعد":      "appointment",
		"الزيارات":    "visits",
		"الزيارة":     "visit",
		"الفواتير":    "invoices",
		"الفاتورة":    "invoice",
		"الدفع":       "payment",
		"المدفوع":     "paid",
		"غير المدفوع": "unpaid",
		"المستحق":     "due",

		// Pharmacy terms
		"الصيدلية":  "pharmacy",
		"الأدوية":   "medicines",
		"الدواء":    "medicine",
		"الوصفات":   "prescriptions",
		"الوصفة":    "prescription",
		"المخزون":   "inventory",
		"المخزونات": "inventory",

		// Financial terms
		"الإيرادات": "revenue",
		"إيرادات":   "revenue",
		"المبيعات":  "sales",
		"الأرباح":   "profits",
		"الربح":     "profit",
		"التكاليف":  "costs",
		"التكلفة":   "cost",
		"المصاريف":  "expenses",
		"المصروفات": "expenses",
		"الهامش":    "margin",
		"الميزانية": "budget",

		// Time periods
		"اليوم":         "today",
		"أمس":           "yesterday",
		"هذا الأسبوع":   "this week",
		"الشهر الحالي":  "current month",
		"هذا الشهر":     "this month",
		"الشهر الماضي":  "last month",
		"الربع الأول":   "first quarter",
		"الربع الثاني":  "second quarter",
		"الربع الثالث":  "third quarter",
		"الربع الأخير":  "last quarter",
		"السنة الحالية": "current year",
		"هذه السنة":     "this year",
		"السنة الماضية": "last year",
		"يناير":         "January",
		"فبراير":        "February",
		"مارس":          "March",
		"أبريل":         "April",
		"مايو":          "May",
		"يونيو":         "June",
		"يوليو":         "July",
		"أغسطس":         "August",
		"سبتمبر":        "September",
		"أكتوبر":        "October",
		"نوفمبر":        "November",
		"ديسمبر":        "December",

		// Query verbs
		"أظهر": "show",
		"اعرض": "show",
		"ما":   "what",
		"كم":   "how much",
		"كيف":  "how",
		"قارن": "compare",
		"احسب": "calculate",
		"أوجد": "find",

		// Aggregation terms
		"إجمالي":  "total",
		"المجموع": "total",
		"معدل":    "average",
		"متوسط":   "average",
		"عدد":     "count",
		"أعلى":    "highest",
		"أقل":     "lowest",
		"أكبر":    "largest",
		"أصغر":    "smallest",

		// Departments/Categories
		"القسم":   "department",
		"الأقسام": "departments",
		"الفئة":   "category",
		"الفئات":  "categories",
		"نوع":     "type",
		"الأنواع": "types",
	}

	// English synonym normalization (non-canonical -> canonical)
	englishNormalization := map[string]string{
		"footfall":    "patient_visits",
		"footfalls":   "patient_visits",
		"visits":      "patient_visits",
		"walkins":     "patient_visits",
		"walk-ins":    "patient_visits",
		"income":      "revenue",
		"earnings":    "revenue",
		"turnover":    "revenue",
		"billing":     "invoices",
		"bills":       "invoices",
		"charges":     "fees",
		"doctors":     "physicians",
		"meds":        "medicines",
		"drugs":       "medicines",
		"medications": "medicines",
		"scripts":     "prescriptions",
		"stock":       "inventory",
		"supplies":    "inventory",
		"spending":    "expenses",
		"expenditure": "expenses",
		"quarter":     "Q",
		"Q1":          "first quarter",
		"Q2":          "second quarter",
		"Q3":          "third quarter",
		"Q4":          "fourth quarter",
	}

	// Terms to preserve exactly (SQL keywords, table names, etc.)
	preservedTerms := map[string]bool{
		"SELECT": true, "FROM": true, "WHERE": true, "JOIN": true,
		"GROUP BY": true, "ORDER BY": true, "HAVING": true,
		"AND": true, "OR": true, "NOT": true, "IN": true,
		"ASC": true, "DESC": true, "LIMIT": true, "OFFSET": true,
		"SUM": true, "COUNT": true, "AVG": true, "MAX": true, "MIN": true,
		"true": true, "false": true, "null": true,
	}

	return &QueryTranslationAgent{
		llm:                  cfg.LLM,
		logger:               logger.With(slog.String("agent", AgentID)),
		arabicToEnglish:      arabicToEnglish,
		englishNormalization: englishNormalization,
		preservedTerms:       preservedTerms,
	}
}

// Translate translates a query from the source locale to English intent.
func (a *QueryTranslationAgent) Translate(ctx context.Context, text string, fromLocale string) (*TranslationResult, error) {
	if text == "" {
		return &TranslationResult{
			OriginalText:   text,
			TranslatedText: text,
			SourceLocale:   fromLocale,
			TargetLocale:   DefaultTargetLocale,
			Confidence:     1.0,
		}, nil
	}

	// If already English, normalize and return
	if fromLocale == "en" {
		return a.normalizeEnglish(text)
	}

	// Translate Arabic to English
	result := &TranslationResult{
		OriginalText: text,
		SourceLocale: fromLocale,
		TargetLocale: DefaultTargetLocale,
	}

	// Step 1: Extract and preserve domain terms
	preservedTerms, preprocessed := a.extractDomainTerms(text)
	result.PreservedTerms = preservedTerms

	// Step 2: Translate using dictionary first
	dictTranslation := a.dictionaryTranslate(preprocessed)

	// Step 3: Use LLM for remaining translation if available
	if a.llm != nil {
		llmTranslation, err := a.llmTranslate(ctx, dictTranslation)
		if err != nil {
			a.logger.Warn("LLM translation failed, using dictionary result",
				slog.String("error", err.Error()),
			)
			result.TranslatedText = dictTranslation
			result.Confidence = 0.7
		} else {
			result.TranslatedText = llmTranslation
			result.Confidence = 0.95
		}
	} else {
		result.TranslatedText = dictTranslation
		result.Confidence = 0.8
	}

	// Step 4: Normalize for SQL generation
	result.TranslatedText = a.normalizeForSQL(result.TranslatedText)

	a.logger.Debug("translation completed",
		slog.String("original", text),
		slog.String("translated", result.TranslatedText),
		slog.Float64("confidence", result.Confidence),
		slog.Any("preserved_terms", result.PreservedTerms),
	)

	return result, nil
}

// normalizeEnglish normalizes English text and returns translation result.
func (a *QueryTranslationAgent) normalizeEnglish(text string) (*TranslationResult, error) {
	normalized := a.applyEnglishNormalization(text)
	return &TranslationResult{
		OriginalText:   text,
		TranslatedText: normalized,
		SourceLocale:   "en",
		TargetLocale:   DefaultTargetLocale,
		Confidence:     1.0,
	}, nil
}

// extractDomainTerms identifies and extracts domain terms to preserve.
func (a *QueryTranslationAgent) extractDomainTerms(text string) ([]string, string) {
	preserved := []string{}
	result := text

	// Check for domain terms that should be preserved
	for arTerm, enTerm := range a.arabicToEnglish {
		if strings.Contains(text, arTerm) {
			preserved = append(preserved, enTerm)
		}
	}

	return preserved, result
}

// dictionaryTranslate performs translation using the built-in dictionary.
func (a *QueryTranslationAgent) dictionaryTranslate(text string) string {
	result := text

	// Replace Arabic terms with English equivalents
	for arTerm, enTerm := range a.arabicToEnglish {
		if strings.Contains(result, arTerm) {
			result = strings.ReplaceAll(result, arTerm, enTerm)
		}
	}

	return result
}

// llmTranslate uses the LLM for high-quality translation.
func (a *QueryTranslationAgent) llmTranslate(ctx context.Context, text string) (string, error) {
	systemPrompt := `You are a professional translator specializing in Arabic to English translation for healthcare and finance business intelligence queries.

Your task is to translate Arabic queries to English intent that will be used for SQL generation.

Rules:
1. Preserve all domain terminology exactly as provided
2. Output only the translated English text, no explanations
3. Use simple, SQL-friendly language
4. Maintain the original query intent
5. Use standard business terms (revenue, patients, appointments, etc.)
6. Preserve any numbers, dates, or specific values`

	userPrompt := fmt.Sprintf("Translate this Arabic query to English intent:\n\n%s", text)

	response, err := a.llm.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("LLM translation failed: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// applyEnglishNormalization normalizes English synonyms to canonical terms.
func (a *QueryTranslationAgent) applyEnglishNormalization(text string) string {
	result := text
	lowerText := strings.ToLower(text)

	for synonym, canonical := range a.englishNormalization {
		if strings.Contains(lowerText, synonym) {
			// Case-insensitive replacement
			result = regexpReplace(result, synonym, canonical)
		}
	}

	return result
}

// normalizeForSQL prepares translated text for SQL generation.
func (a *QueryTranslationAgent) normalizeForSQL(text string) string {
	result := text

	// Normalize whitespace
	result = strings.Join(strings.Fields(result), " ")

	// Apply English normalization
	result = a.applyEnglishNormalization(result)

	// Ensure common SQL-friendly patterns
	// (Additional normalization rules can be added here)

	return result
}

// TranslateBatch translates multiple queries in a batch.
func (a *QueryTranslationAgent) TranslateBatch(ctx context.Context, queries []string, fromLocale string) ([]*TranslationResult, error) {
	results := make([]*TranslationResult, len(queries))

	for i, query := range queries {
		result, err := a.Translate(ctx, query, fromLocale)
		if err != nil {
			return nil, fmt.Errorf("failed to translate query %d: %w", i, err)
		}
		results[i] = result
	}

	return results, nil
}

// GetCanonicalTerm returns the canonical form of a term.
func (a *QueryTranslationAgent) GetCanonicalTerm(term string) string {
	lower := strings.ToLower(term)
	if canonical, ok := a.englishNormalization[lower]; ok {
		return canonical
	}
	return term
}

// regexpReplace performs case-insensitive replacement.
func regexpReplace(text, pattern, replacement string) string {
	// Simple case-insensitive replacement
	lowerText := strings.ToLower(text)
	lowerPattern := strings.ToLower(pattern)

	result := text
	if strings.Contains(lowerText, lowerPattern) {
		// Find the actual occurrence and replace
		idx := strings.Index(lowerText, lowerPattern)
		if idx >= 0 {
			actualPattern := text[idx : idx+len(pattern)]
			result = strings.ReplaceAll(text, actualPattern, replacement)
		}
	}

	return result
}

// AddCustomTerm adds a custom term mapping.
func (a *QueryTranslationAgent) AddCustomTerm(arabicTerm, englishTerm string) {
	a.arabicToEnglish[arabicTerm] = englishTerm
}

// AddCustomNormalization adds a custom English normalization rule.
func (a *QueryTranslationAgent) AddCustomNormalization(synonym, canonical string) {
	a.englishNormalization[strings.ToLower(synonym)] = canonical
}
