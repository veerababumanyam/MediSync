// Package middleware provides HTTP middleware for the MediSync API.
//
// This file implements the TimeoutMiddleware for request timeout handling.
// It ensures that long-running requests are properly cancelled and return
// appropriate 504 Gateway Timeout responses.
//
// The middleware:
//   - Sets a deadline for all requests
//   - Cancels the context on timeout for proper cleanup
//   - Returns a localized 504 response on timeout
//   - Works with query operations that may take longer
//
// Usage:
//
//	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
package middleware

import (
	"context"
	"net/http"
	"time"
)

// DefaultTimeout is the default request timeout duration.
const DefaultTimeout = 30 * time.Second

// TimeoutContextKey is the context key for the timeout deadline.
const TimeoutContextKey contextKey = "timeout_deadline"

// TimeoutMiddleware returns a middleware that enforces a request timeout.
// If the request exceeds the timeout, it returns 504 Gateway Timeout.
// The context is cancelled to allow handlers to clean up resources.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Store the deadline in context for handlers to check
			deadline, ok := ctx.Deadline()
			if ok {
				ctx = context.WithValue(ctx, TimeoutContextKey, deadline)
			}

			// Create a response wrapper to detect if response was written
			tw := &timeoutWriter{
				ResponseWriter: w,
				h:              make(http.Header),
			}

			// Channel to signal handler completion
			done := make(chan struct{})

			// Run handler in goroutine
			go func() {
				defer close(done)
				next.ServeHTTP(tw, r.WithContext(ctx))
			}()

			select {
			case <-done:
				// Handler completed normally
				tw.copyHeaders(w)
				if tw.code > 0 {
					w.WriteHeader(tw.code)
				}
				w.Write(tw.buf)
				return

			case <-ctx.Done():
				// Context was cancelled (timeout or client disconnect)
				if ctx.Err() == context.DeadlineExceeded {
					// Request timed out
					writeTimeoutResponse(w, r)
					return
				}
				// Client disconnected - nothing to do
				return
			}
		})
	}
}

// TimeoutMiddlewareWithConfig returns a timeout middleware with custom configuration.
func TimeoutMiddlewareWithConfig(config TimeoutConfig) func(http.Handler) http.Handler {
	if config.Timeout <= 0 {
		config.Timeout = DefaultTimeout
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determine timeout based on request path
			timeout := config.Timeout
			if pathTimeout, ok := config.PathTimeouts[r.URL.Path]; ok {
				timeout = pathTimeout
			}

			// Apply timeout middleware with determined duration
			timeoutMW := TimeoutMiddleware(timeout)
			timeoutMW(next).ServeHTTP(w, r)
		})
	}
}

// TimeoutConfig holds configuration for the timeout middleware.
type TimeoutConfig struct {
	// Timeout is the default request timeout.
	Timeout time.Duration

	// PathTimeouts allows different timeouts for specific paths.
	PathTimeouts map[string]time.Duration

	// SkippedPaths are paths that should not have timeout enforcement.
	SkippedPaths map[string]bool
}

// timeoutWriter wraps http.ResponseWriter to capture response state.
type timeoutWriter struct {
	http.ResponseWriter
	h       http.Header
	buf     []byte
	code    int
	written bool
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	if !tw.written {
		tw.written = true
	}
	tw.buf = append(tw.buf, p...)
	return len(p), nil
}

func (tw *timeoutWriter) WriteHeader(code int) {
	if !tw.written {
		tw.code = code
		tw.written = true
	}
}

func (tw *timeoutWriter) copyHeaders(dst http.ResponseWriter) {
	for k, vv := range tw.h {
		dst.Header()[k] = vv
	}
}

// writeTimeoutResponse writes a 504 Gateway Timeout response.
func writeTimeoutResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusGatewayTimeout)

	// Localized error message
	message := "The request took too long to process. Please try again later."
	if locale == "ar" {
		message = "استغرق الطلب وقتًا طويلاً للمعالجة. يرجى المحاولة مرة أخرى لاحقًا."
	}

	// Simple JSON encoding
	w.Write([]byte(`{"error":{"code":"GATEWAY_TIMEOUT","message":"` + message + `"}}`))
}

// getLocaleFromContext extracts locale from the request context.
func getLocaleFromContext(ctx context.Context) string {
	if locale, ok := ctx.Value(LocaleKey).(string); ok {
		return locale
	}
	return "en"
}

// QueryTimeoutMiddleware returns a middleware specifically for query endpoints.
// Query operations may need longer timeouts due to database operations.
func QueryTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	if timeout <= 0 {
		timeout = 45 * time.Second // Longer default for queries
	}

	return TimeoutMiddleware(timeout)
}

// GetTimeoutDeadline returns the timeout deadline from the context.
func GetTimeoutDeadline(ctx context.Context) (time.Time, bool) {
	if deadline, ok := ctx.Deadline(); ok {
		return deadline, true
	}
	if deadline, ok := ctx.Value(TimeoutContextKey).(time.Time); ok {
		return deadline, true
	}
	return time.Time{}, false
}

// TimeRemaining returns the remaining time before the request times out.
func TimeRemaining(ctx context.Context) time.Duration {
	deadline, ok := GetTimeoutDeadline(ctx)
	if !ok {
		return 0 // No deadline set
	}
	remaining := time.Until(deadline)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsTimedOut checks if the context has timed out.
func IsTimedOut(ctx context.Context) bool {
	return ctx.Err() == context.DeadlineExceeded
}

// IsCancelled checks if the context was cancelled (client disconnect or timeout).
func IsCancelled(ctx context.Context) bool {
	return ctx.Err() != nil
}
