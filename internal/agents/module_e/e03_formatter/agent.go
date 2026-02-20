// Package e03_formatter provides the E-03 Localized Formatter Agent.
//
// This agent formats numbers, currency, dates, and responses according to the
// user's locale, including Eastern Arabic numeral conversion for Arabic locale.
//
// Formatting Rules:
// - English: 1,234.56 (Western Arabic numerals, comma thousands, period decimal)
// - Arabic: ١٬٢٣٤٫٥٦ (Eastern Arabic numerals, ٬ thousands, ٫ decimal)
//
// Usage:
//
//	agent := e03_formatter.New(logger)
//	formatted := agent.FormatNumber(1234.56, "ar") // "١٬٢٣٤٫٥٦"
//	formatted := agent.FormatCurrency(125000.00, "INR", "ar") // "₹١٬٢٥٬٠٠٠٫٠٠"
package e03_formatter

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// AgentID is the identifier for this agent.
	AgentID = "E-03"

	// AgentName is the human-readable name.
	AgentName = "Localized Formatter Agent"

	// DefaultLocale is the default locale for formatting.
	DefaultLocale = "en"
)

// Eastern Arabic numerals mapping (Western -> Eastern)
var easternArabicNumerals = map[rune]rune{
	'0': '٠',
	'1': '١',
	'2': '٢',
	'3': '٣',
	'4': '٤',
	'5': '٥',
	'6': '٦',
	'7': '٧',
	'8': '٨',
	'9': '٩',
}

// Arabic punctuation
const (
	// ArabicDecimalSeparator is the Arabic decimal separator (U+066B)
	ArabicDecimalSeparator = '٫'
	// ArabicThousandsSeparator is the Arabic thousands separator (U+066C)
	ArabicThousandsSeparator = '٬'
	// WesternDecimalSeparator is the Western decimal separator
	WesternDecimalSeparator = '.'
	// WesternThousandsSeparator is the Western thousands separator
	WesternThousandsSeparator = ','
)

// FormattedResult represents the result of formatting.
type FormattedResult struct {
	// Value is the original value.
	Value any `json:"value"`

	// Formatted is the formatted string.
	Formatted string `json:"formatted"`

	// Locale is the locale used for formatting.
	Locale string `json:"locale"`

	// FormatType is the type of formatting applied.
	FormatType string `json:"format_type,omitempty"`
}

// Date formats for different locales
var dateFormats = map[string]string{
	"en": "January 2, 2006",
	"ar": "٢ يناير ٢٠٠٦", // Will be post-processed for Arabic numerals
}

// Short date formats
var shortDateFormats = map[string]string{
	"en": "2006-01-02",
	"ar": "٢٠٠٦-٠١-٠٢", // Will be post-processed
}

// Month names in Arabic
var arabicMonths = map[string]string{
	"January":   "يناير",
	"February":  "فبراير",
	"March":     "مارس",
	"April":     "أبريل",
	"May":       "مايو",
	"June":      "يونيو",
	"July":      "يوليو",
	"August":    "أغسطس",
	"September": "سبتمبر",
	"October":   "أكتوبر",
	"November":  "نوفمبر",
	"December":  "ديسمبر",
}

// LocalizedFormatterAgent formats values according to locale.
type LocalizedFormatterAgent struct {
	logger *slog.Logger

	// Currency symbols
	currencySymbols map[string]string

	// Currency decimal places
	currencyDecimals map[string]int
}

// Config holds configuration for the agent.
type Config struct {
	// Logger is the structured logger.
	Logger *slog.Logger
}

// New creates a new E-03 Localized Formatter Agent.
func New(cfg *Config) *LocalizedFormatterAgent {
	logger := slog.Default()
	if cfg != nil && cfg.Logger != nil {
		logger = cfg.Logger
	}

	// Currency symbols (can be extended)
	currencySymbols := map[string]string{
		"INR": "₹",
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"AED": "د.إ", // UAE Dirham
		"SAR": "﷼",   // Saudi Riyal
	}

	// Default decimal places for currencies
	currencyDecimals := map[string]int{
		"INR": 2,
		"USD": 2,
		"EUR": 2,
		"GBP": 2,
		"AED": 2,
		"SAR": 2,
		"JPY": 0,
		"KWD": 3,
	}

	return &LocalizedFormatterAgent{
		logger:           logger.With(slog.String("agent", AgentID)),
		currencySymbols:  currencySymbols,
		currencyDecimals: currencyDecimals,
	}
}

// FormatNumber formats a number according to locale.
func (a *LocalizedFormatterAgent) FormatNumber(n float64, locale string) string {
	return a.FormatNumberWithPrecision(n, locale, 2)
}

// FormatNumberWithPrecision formats a number with specified decimal places.
func (a *LocalizedFormatterAgent) FormatNumberWithPrecision(n float64, locale string, decimals int) string {
	// Format with grouping (thousands separator)
	format := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(format, n)

	// Add thousands separators
	formatted = a.addThousandsSeparator(formatted, locale)

	// Convert to Eastern Arabic numerals if Arabic locale
	if locale == "ar" {
		formatted = a.toEasternArabic(formatted)
	}

	return formatted
}

// FormatCurrency formats a currency amount according to locale.
func (a *LocalizedFormatterAgent) FormatCurrency(amount float64, currency, locale string) string {
	// Get decimal places for currency
	decimals := 2
	if d, ok := a.currencyDecimals[currency]; ok {
		decimals = d
	}

	// Format the number
	format := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(format, math.Abs(amount))

	// Add thousands separators
	formatted = a.addThousandsSeparator(formatted, locale)

	// Get currency symbol
	symbol := currency
	if s, ok := a.currencySymbols[currency]; ok {
		symbol = s
	}

	// Add negative sign if needed
	negativePrefix := ""
	if amount < 0 {
		negativePrefix = "-"
	}

	// Format based on locale
	var result string
	if locale == "ar" {
		// Arabic: symbol on right, convert numerals
		formatted = a.toEasternArabic(formatted)
		result = fmt.Sprintf("%s%s %s", negativePrefix, formatted, symbol)
	} else {
		// English: symbol on left
		result = fmt.Sprintf("%s%s%s", negativePrefix, symbol, formatted)
	}

	return result
}

// FormatCurrencyWithCode formats currency with the currency code.
func (a *LocalizedFormatterAgent) FormatCurrencyWithCode(amount float64, currency, locale string) string {
	formatted := a.FormatCurrency(amount, currency, locale)

	// Append currency code
	if locale == "ar" {
		return fmt.Sprintf("%s (%s)", formatted, currency)
	}
	return fmt.Sprintf("%s (%s)", formatted, currency)
}

// FormatDate formats a date according to locale.
func (a *LocalizedFormatterAgent) FormatDate(t time.Time, locale string) string {
	if locale == "ar" {
		return a.formatDateArabic(t)
	}
	return t.Format("January 2, 2006")
}

// FormatDateShort formats a date in short format according to locale.
func (a *LocalizedFormatterAgent) FormatDateShort(t time.Time, locale string) string {
	if locale == "ar" {
		// Format as YYYY-MM-DD with Eastern Arabic numerals
		formatted := t.Format("2006-01-02")
		return a.toEasternArabic(formatted)
	}
	return t.Format("2006-01-02")
}

// FormatTime formats a time according to locale.
func (a *LocalizedFormatterAgent) FormatTime(t time.Time, locale string) string {
	if locale == "ar" {
		// 24-hour format with Eastern Arabic numerals
		formatted := t.Format("15:04")
		return a.toEasternArabic(formatted)
	}
	return t.Format("3:04 PM")
}

// FormatDateTime formats date and time according to locale.
func (a *LocalizedFormatterAgent) FormatDateTime(t time.Time, locale string) string {
	if locale == "ar" {
		date := a.formatDateArabic(t)
		time := a.FormatTime(t, locale)
		return fmt.Sprintf("%s، %s", date, time) // Arabic comma
	}
	return t.Format("January 2, 2006 at 3:04 PM")
}

// FormatPercentage formats a percentage according to locale.
func (a *LocalizedFormatterAgent) FormatPercentage(value float64, locale string) string {
	return a.FormatPercentageWithPrecision(value, locale, 1)
}

// FormatPercentageWithPrecision formats a percentage with specified decimal places.
func (a *LocalizedFormatterAgent) FormatPercentageWithPrecision(value float64, locale string, decimals int) string {
	format := fmt.Sprintf("%%.%df%%%%", decimals)
	formatted := fmt.Sprintf(format, value)

	if locale == "ar" {
		// Convert numerals and reorder percentage sign
		formatted = a.toEasternArabic(strings.TrimSuffix(formatted, "%"))
		formatted = "٪" + formatted // Arabic percent sign on the left
	}

	return formatted
}

// FormatResponse formats an entire response object according to locale.
func (a *LocalizedFormatterAgent) FormatResponse(data any, locale string) (string, error) {
	switch v := data.(type) {
	case float64:
		return a.FormatNumber(v, locale), nil
	case int:
		return a.FormatNumber(float64(v), locale), nil
	case int64:
		return a.FormatNumber(float64(v), locale), nil
	case string:
		if locale == "ar" {
			return a.toEasternArabic(v), nil
		}
		return v, nil
	case time.Time:
		return a.FormatDate(v, locale), nil
	case map[string]any:
		return a.formatMapResponse(v, locale)
	case []any:
		return a.formatSliceResponse(v, locale)
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// formatMapResponse formats a map response.
func (a *LocalizedFormatterAgent) formatMapResponse(data map[string]any, locale string) (string, error) {
	var parts []string

	for key, value := range data {
		formatted, err := a.FormatResponse(value, locale)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%s: %s", key, formatted))
	}

	return strings.Join(parts, ", "), nil
}

// formatSliceResponse formats a slice response.
func (a *LocalizedFormatterAgent) formatSliceResponse(data []any, locale string) (string, error) {
	var parts []string

	for _, value := range data {
		formatted, err := a.FormatResponse(value, locale)
		if err != nil {
			return "", err
		}
		parts = append(parts, formatted)
	}

	return strings.Join(parts, ", "), nil
}

// formatDateArabic formats a date in Arabic.
func (a *LocalizedFormatterAgent) formatDateArabic(t time.Time) string {
	day := strconv.Itoa(t.Day())
	month := arabicMonths[t.Format("January")]
	year := strconv.Itoa(t.Year())

	// Convert numerals to Eastern Arabic
	day = a.toEasternArabic(day)
	year = a.toEasternArabic(year)

	return fmt.Sprintf("%s %s %s", day, month, year)
}

// addThousandsSeparator adds thousands separators to a number string.
func (a *LocalizedFormatterAgent) addThousandsSeparator(numStr string, locale string) string {
	// Split into integer and decimal parts
	parts := strings.Split(numStr, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Add thousands separators to integer part
	var result []rune
	for i, r := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 && r != '-' {
			if locale == "ar" {
				result = append(result, ArabicThousandsSeparator)
			} else {
				result = append(result, WesternThousandsSeparator)
			}
		}
		result = append(result, r)
	}

	// Reconstruct with decimal part
	if decimalPart != "" {
		if locale == "ar" {
			return string(result) + string(ArabicDecimalSeparator) + decimalPart
		}
		return string(result) + string(WesternDecimalSeparator) + decimalPart
	}

	return string(result)
}

// toEasternArabic converts Western Arabic numerals to Eastern Arabic numerals.
func (a *LocalizedFormatterAgent) toEasternArabic(s string) string {
	result := make([]rune, 0, len(s))

	for _, r := range s {
		if eastern, ok := easternArabicNumerals[r]; ok {
			result = append(result, eastern)
		} else if r == WesternDecimalSeparator {
			result = append(result, ArabicDecimalSeparator)
		} else if r == WesternThousandsSeparator {
			result = append(result, ArabicThousandsSeparator)
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

// ToWesternArabic converts Eastern Arabic numerals back to Western Arabic numerals.
func (a *LocalizedFormatterAgent) ToWesternArabic(s string) string {
	// Create reverse mapping
	reverseMapping := make(map[rune]rune)
	for western, eastern := range easternArabicNumerals {
		reverseMapping[eastern] = western
	}

	result := make([]rune, 0, len(s))

	for _, r := range s {
		if western, ok := reverseMapping[r]; ok {
			result = append(result, western)
		} else if r == ArabicDecimalSeparator {
			result = append(result, WesternDecimalSeparator)
		} else if r == ArabicThousandsSeparator {
			result = append(result, WesternThousandsSeparator)
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

// FormatWithContext formats a value using context for cancellation.
func (a *LocalizedFormatterAgent) FormatWithContext(ctx context.Context, value any, locale string, formatType string) (*FormattedResult, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var formatted string
	var err error

	switch formatType {
	case "number":
		formatted = a.FormatNumber(value.(float64), locale)
	case "currency":
		// Expecting map with "amount" and "currency" keys
		if m, ok := value.(map[string]any); ok {
			amount := m["amount"].(float64)
			currency := m["currency"].(string)
			formatted = a.FormatCurrency(amount, currency, locale)
		} else {
			err = fmt.Errorf("invalid currency format input")
		}
	case "date":
		formatted = a.FormatDate(value.(time.Time), locale)
	case "percentage":
		formatted = a.FormatPercentage(value.(float64), locale)
	default:
		formatted, err = a.FormatResponse(value, locale)
	}

	if err != nil {
		return nil, err
	}

	return &FormattedResult{
		Value:      value,
		Formatted:  formatted,
		Locale:     locale,
		FormatType: formatType,
	}, nil
}

// AddCurrencySymbol adds a custom currency symbol.
func (a *LocalizedFormatterAgent) AddCurrencySymbol(code, symbol string) {
	a.currencySymbols[code] = symbol
}

// SetCurrencyDecimals sets the decimal places for a currency.
func (a *LocalizedFormatterAgent) SetCurrencyDecimals(code string, decimals int) {
	a.currencyDecimals[code] = decimals
}

// GetSupportedLocales returns the list of supported locales.
func (a *LocalizedFormatterAgent) GetSupportedLocales() []string {
	return []string{"en", "ar"}
}

// FormatInteger formats an integer with locale-aware thousands separator.
func (a *LocalizedFormatterAgent) FormatInteger(n int64, locale string) string {
	// Format without decimals
	formatted := fmt.Sprintf("%d", n)
	formatted = a.addThousandsSeparator(formatted, locale)

	if locale == "ar" {
		formatted = a.toEasternArabic(formatted)
	}

	return formatted
}

// FormatCompactNumber formats large numbers in compact form (e.g., 1.2K, 1.5M).
func (a *LocalizedFormatterAgent) FormatCompactNumber(n float64, locale string) string {
	absN := math.Abs(n)
	var value float64
	var suffix string

	switch {
	case absN >= 1e9:
		value = n / 1e9
		suffix = "B"
		if locale == "ar" {
			suffix = "مليار"
		}
	case absN >= 1e6:
		value = n / 1e6
		suffix = "M"
		if locale == "ar" {
			suffix = "مليون"
		}
	case absN >= 1e3:
		value = n / 1e3
		suffix = "K"
		if locale == "ar" {
			suffix = "ألف"
		}
	default:
		return a.FormatNumber(n, locale)
	}

	formatted := a.FormatNumberWithPrecision(value, locale, 1)
	return fmt.Sprintf("%s %s", formatted, suffix)
}
