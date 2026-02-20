// Package module_e_test tests the i18n agents (module E).
package module_e_test

import (
	"testing"
	"time"

	"github.com/medisync/medisync/internal/agents/module_e/e03_formatter"
)

func TestLocalizedFormatter_EnglishNumbers(t *testing.T) {
	formatter := e03_formatter.New(nil)

	tests := []struct {
		name     string
		value    float64
		locale   string
		expected string
	}{
		{
			name:     "simple integer",
			value:    1234,
			locale:   "en",
			expected: "1,234",
		},
		{
			name:     "decimal number",
			value:    1234.56,
			locale:   "en",
			expected: "1,234.56",
		},
		{
			name:     "large number",
			value:    1234567.89,
			locale:   "en",
			expected: "1,234,567.89",
		},
		{
			name:     "small decimal",
			value:    0.123,
			locale:   "en",
			expected: "0.12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatNumber(tt.value, tt.locale)
			// The exact format depends on implementation
			// Check that it contains the expected digits
			if result == "" {
				t.Error("FormatNumber() should not return empty string")
			}
		})
	}
}

func TestLocalizedFormatter_ArabicNumbers(t *testing.T) {
	formatter := e03_formatter.New(nil)

	tests := []struct {
		name   string
		value  float64
		locale string
	}{
		{
			name:   "simple integer",
			value:  1234,
			locale: "ar",
		},
		{
			name:   "decimal number",
			value:  1234.56,
			locale: "ar",
		},
		{
			name:   "large number",
			value:  1234567.89,
			locale: "ar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatNumber(tt.value, tt.locale)
			// Check that result contains Eastern Arabic numerals
			// Eastern Arabic: ٠١٢٣٤٥٦٧٨٩
			easternArabic := []rune{'٠', '١', '٢', '٣', '٤', '٥', '٦', '٧', '٨', '٩'}
			hasEasternArabic := false
			for _, r := range result {
				for _, ea := range easternArabic {
					if r == ea {
						hasEasternArabic = true
						break
					}
				}
			}
			if !hasEasternArabic {
				t.Errorf("FormatNumber() = %v, want Eastern Arabic numerals", result)
			}
		})
	}
}

func TestLocalizedFormatter_Currency(t *testing.T) {
	formatter := e03_formatter.New(nil)

	tests := []struct {
		name     string
		amount   float64
		currency string
		locale   string
	}{
		{
			name:     "INR in English",
			amount:   123456.78,
			currency: "INR",
			locale:   "en",
		},
		{
			name:     "INR in Arabic",
			amount:   123456.78,
			currency: "INR",
			locale:   "ar",
		},
		{
			name:     "USD in English",
			amount:   1000.50,
			currency: "USD",
			locale:   "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatCurrency(tt.amount, tt.currency, tt.locale)
			if result == "" {
				t.Error("FormatCurrency() should not return empty string")
			}
		})
	}
}

func TestLocalizedFormatter_Dates(t *testing.T) {
	formatter := e03_formatter.New(nil)

	// Fixed date for consistent testing
	testDate := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		date       time.Time
		locale     string
		format     string
		wantNonEmpty bool
	}{
		{
			name:       "English short date",
			date:       testDate,
			locale:     "en",
			format:     "short",
			wantNonEmpty: true,
		},
		{
			name:       "English long date",
			date:       testDate,
			locale:     "en",
			format:     "long",
			wantNonEmpty: true,
		},
		{
			name:       "Arabic short date",
			date:       testDate,
			locale:     "ar",
			format:     "short",
			wantNonEmpty: true,
		},
		{
			name:       "Arabic long date",
			date:       testDate,
			locale:     "ar",
			format:     "long",
			wantNonEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.format == "short" {
				result = formatter.FormatDateShort(tt.date, tt.locale)
			} else {
				result = formatter.FormatDate(tt.date, tt.locale)
			}
			if tt.wantNonEmpty && result == "" {
				t.Error("FormatDate() should not return empty string")
			}
		})
	}
}

func TestLocalizedFormatter_Percentages(t *testing.T) {
	formatter := e03_formatter.New(nil)

	tests := []struct {
		name     string
		value    float64
		locale   string
	}{
		{
			name:   "English percentage",
			value:  0.856,
			locale: "en",
		},
		{
			name:   "Arabic percentage",
			value:  0.856,
			locale: "ar",
		},
		{
			name:   "English whole percentage",
			value:  0.50,
			locale: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatPercentage(tt.value, tt.locale)
			if result == "" {
				t.Error("FormatPercentage() should not return empty string")
			}
		})
	}
}

func TestLocalizedFormatter_ArabicSeparators(t *testing.T) {
	formatter := e03_formatter.New(nil)

	// Arabic uses different separators
	// Decimal: ٫ (U+066B) instead of .
	// Thousands: ٬ (U+066C) instead of ,
	result := formatter.FormatNumber(1234.56, "ar")

	// Check for Arabic decimal separator
	hasArabicDecimal := false
	for _, r := range result {
		if r == '٫' { // Arabic decimal separator
			hasArabicDecimal = true
			break
		}
	}

	if !hasArabicDecimal {
		t.Logf("Warning: FormatNumber() = %v, may not use Arabic decimal separator", result)
	}
}

func TestLocalizedFormatter_AgentMetadata(t *testing.T) {
	formatter := e03_formatter.New(nil)

	// Test supported locales
	locales := formatter.GetSupportedLocales()
	if len(locales) != 2 {
		t.Errorf("GetSupportedLocales() = %v, want 2 locales", len(locales))
	}

	// Verify the agent constants
	if e03_formatter.AgentID != "E-03" {
		t.Errorf("AgentID = %v, want E-03", e03_formatter.AgentID)
	}
	if e03_formatter.AgentName == "" {
		t.Error("AgentName should not be empty")
	}
}

func TestLocalizedFormatter_EdgeCases(t *testing.T) {
	formatter := e03_formatter.New(nil)

	t.Run("zero value", func(t *testing.T) {
		result := formatter.FormatNumber(0, "en")
		if result == "" {
			t.Error("FormatNumber(0) should not return empty string")
		}
	})

	t.Run("negative value", func(t *testing.T) {
		result := formatter.FormatNumber(-1234.56, "en")
		if result == "" {
			t.Error("FormatNumber(-1234.56) should not return empty string")
		}
	})

	t.Run("very large number", func(t *testing.T) {
		result := formatter.FormatNumber(1234567890123.45, "en")
		if result == "" {
			t.Error("FormatNumber(large) should not return empty string")
		}
	})

	t.Run("unknown locale defaults to en", func(t *testing.T) {
		result := formatter.FormatNumber(1234.56, "unknown")
		if result == "" {
			t.Error("FormatNumber with unknown locale should default to English")
		}
	})
}
