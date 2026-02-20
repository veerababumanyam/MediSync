// Package middleware provides HTTP middleware for the MediSync API.
//
// This file implements the LocaleMiddleware that extracts user locale from
// the request with the priority: URL param > JWT claim > Accept-Language header > default "en".
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// LocaleKey is the context key for user locale.
const LocaleKey contextKey = "locale"

// Supported locales
var supportedLocales = map[string]bool{
	"en": true, // English (default)
	"ar": true, // Arabic (RTL)
}

// DefaultLocale is the default locale when none is specified.
const DefaultLocale = "en"

// LocaleMiddleware extracts user locale from the request.
// Priority: URL param > JWT claim > Accept-Language header > default "en"
// The locale is validated to be either "en" or "ar".
func LocaleMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var locale string

			// 1. Check URL parameter (highest priority)
			if l := r.URL.Query().Get("lang"); l != "" {
				locale = normalizeLocale(l)
				logger.Debug("locale from URL parameter",
					slog.String("locale", locale),
				)
			}

			// 2. Check JWT claim
			if locale == "" {
				if claims := GetClaims(r.Context()); claims != nil && claims.Locale != "" {
					locale = normalizeLocale(claims.Locale)
					logger.Debug("locale from JWT claim",
						slog.String("locale", locale),
					)
				}
			}

			// 3. Check Accept-Language header
			if locale == "" {
				acceptLanguage := r.Header.Get("Accept-Language")
				if acceptLanguage != "" {
					locale = parseAcceptLanguage(acceptLanguage)
					logger.Debug("locale from Accept-Language header",
						slog.String("locale", locale),
						slog.String("accept_language", acceptLanguage),
					)
				}
			}

			// 4. Fall back to default
			if locale == "" {
				locale = DefaultLocale
				logger.Debug("using default locale",
					slog.String("locale", locale),
				)
			}

			// Validate locale is supported
			if !supportedLocales[locale] {
				logger.Warn("unsupported locale, falling back to default",
					slog.String("requested_locale", locale),
					slog.String("default_locale", DefaultLocale),
				)
				locale = DefaultLocale
			}

			// Add locale to context
			ctx := context.WithValue(r.Context(), LocaleKey, locale)

			// Set Content-Language header for response
			w.Header().Set("Content-Language", locale)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetLocale retrieves the locale from the request context.
func GetLocale(ctx context.Context) string {
	if locale, ok := ctx.Value(LocaleKey).(string); ok {
		return locale
	}
	return DefaultLocale
}

// IsRTL returns true if the locale is a right-to-left language.
func IsRTL(ctx context.Context) bool {
	locale := GetLocale(ctx)
	return locale == "ar"
}

// GetTextDirection returns "rtl" or "ltr" based on the locale.
func GetTextDirection(ctx context.Context) string {
	if IsRTL(ctx) {
		return "rtl"
	}
	return "ltr"
}

// normalizeLocale normalizes a locale string to the standard format.
// Examples: "en-US" -> "en", "ar-SA" -> "ar", "EN" -> "en"
func normalizeLocale(locale string) string {
	// Convert to lowercase
	locale = strings.ToLower(locale)

	// Extract language part from locale (e.g., "en-us" -> "en")
	if idx := strings.Index(locale, "-"); idx > 0 {
		locale = locale[:idx]
	}

	// Extract language part from underscore format (e.g., "en_US" -> "en")
	if idx := strings.Index(locale, "_"); idx > 0 {
		locale = locale[:idx]
	}

	return locale
}

// parseAcceptLanguage parses the Accept-Language header and returns the
// first supported locale found.
func parseAcceptLanguage(acceptLanguage string) string {
	// Accept-Language format: "en-US,en;q=0.9,ar;q=0.8"
	languages := strings.Split(acceptLanguage, ",")

	for _, lang := range languages {
		// Remove quality value if present
		if idx := strings.Index(lang, ";"); idx > 0 {
			lang = lang[:idx]
		}

		lang = strings.TrimSpace(lang)
		normalized := normalizeLocale(lang)

		if supportedLocales[normalized] {
			return normalized
		}
	}

	return ""
}

// GetLocaleName returns the human-readable name for a locale.
func GetLocaleName(locale string) string {
	switch locale {
	case "en":
		return "English"
	case "ar":
		return "العربية" // Arabic
	default:
		return locale
	}
}

// GetAllSupportedLocales returns all supported locales.
func GetAllSupportedLocales() []string {
	locales := make([]string, 0, len(supportedLocales))
	for locale := range supportedLocales {
		locales = append(locales, locale)
	}
	return locales
}
