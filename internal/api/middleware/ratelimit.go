// Package middleware provides HTTP middleware for the MediSync API.
//
// This file implements the RateLimitMiddleware that rate limits requests per user
// using Redis or in-memory storage as a fallback.
package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// RateLimiter defines the interface for rate limiting.
type RateLimiter interface {
	// Allow checks if a request is allowed for the given key.
	// Returns (true, 0) if allowed, (false, retryAfterSeconds) if not.
	Allow(ctx context.Context, key string) (bool, int, error)

	// Increment increments the counter for the given key.
	Increment(ctx context.Context, key string, window time.Duration) error
}

// RateLimitMiddleware rate limits requests per user based on their user ID.
// It uses Redis for distributed rate limiting or falls back to in-memory.
// Default: 60 requests per minute per user.
func RateLimitMiddleware(cache CacheClient, logger *slog.Logger, requestsPerMinute int) func(http.Handler) http.Handler {
	// Create rate limiter
	var limiter RateLimiter
	if cache != nil {
		limiter = NewRedisRateLimiter(cache, requestsPerMinute, time.Minute)
		logger.Info("using Redis rate limiter",
			slog.Int("requests_per_minute", requestsPerMinute),
		)
	} else {
		limiter = NewMemoryRateLimiter(requestsPerMinute, time.Minute)
		logger.Warn("Redis not available, using in-memory rate limiter",
			slog.Int("requests_per_minute", requestsPerMinute),
		)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip rate limiting for health and ready endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}

			// Get user ID from context (set by AuthMiddleware)
			userID := GetUserID(r.Context())
			if userID == "" {
				// If no user ID, use IP address as fallback
				userID = r.RemoteAddr
			}

			// Create rate limit key
			key := fmt.Sprintf("ratelimit:%s", userID)

			// Check rate limit
			allowed, retryAfter, err := limiter.Allow(r.Context(), key)
			if err != nil {
				logger.Error("rate limit check failed",
					slog.Any("error", err),
					slog.String("user_id", userID),
				)
				// On error, allow the request (fail open)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				logger.Warn("rate limit exceeded",
					slog.String("user_id", userID),
					slog.Int("retry_after", retryAfter),
				)
				writeRateLimited(w, retryAfter)
				return
			}

			// Increment counter
			if err := limiter.Increment(r.Context(), key, time.Minute); err != nil {
				logger.Error("rate limit increment failed",
					slog.Any("error", err),
					slog.String("user_id", userID),
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// writeRateLimited writes a 429 Too Many Requests response.
func writeRateLimited(w http.ResponseWriter, retryAfter int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(fmt.Sprintf(`{"error":{"code":"rate_limited","message":"rate limit exceeded, retry after %d seconds"}}`, retryAfter)))
}

// CacheClient defines the interface for cache operations.
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

// RedisRateLimiter implements rate limiting using Redis.
type RedisRateLimiter struct {
	cache       CacheClient
	maxRequests int
	window      time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter.
func NewRedisRateLimiter(cache CacheClient, maxRequests int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		cache:       cache,
		maxRequests: maxRequests,
		window:      window,
	}
}

// Allow checks if a request is allowed.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
	count, err := r.cache.Increment(ctx, key)
	if err != nil {
		return false, 0, err
	}

	// Set expiry on first increment
	if count == 1 {
		if err := r.cache.Expire(ctx, key, r.window); err != nil {
			return false, 0, err
		}
	}

	if count > int64(r.maxRequests) {
		// Calculate retry after
		retryAfter := int(r.window.Seconds())
		return false, retryAfter, nil
	}

	return true, 0, nil
}

// Increment increments the counter.
func (r *RedisRateLimiter) Increment(ctx context.Context, key string, window time.Duration) error {
	// Already incremented in Allow
	return nil
}

// MemoryRateLimiter implements rate limiting using in-memory storage.
type MemoryRateLimiter struct {
	counters    map[string]*counter
	maxRequests int
	window      time.Duration
}

type counter struct {
	count     int
	expiresAt time.Time
}

// NewMemoryRateLimiter creates a new in-memory rate limiter.
func NewMemoryRateLimiter(maxRequests int, window time.Duration) *MemoryRateLimiter {
	limiter := &MemoryRateLimiter{
		counters:    make(map[string]*counter),
		maxRequests: maxRequests,
		window:      window,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Allow checks if a request is allowed.
func (m *MemoryRateLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
	c, exists := m.counters[key]

	now := time.Now()

	if !exists || now.After(c.expiresAt) {
		m.counters[key] = &counter{
			count:     0,
			expiresAt: now.Add(m.window),
		}
		c = m.counters[key]
	}

	if c.count >= m.maxRequests {
		retryAfter := int(time.Until(c.expiresAt).Seconds())
		if retryAfter < 0 {
			retryAfter = 1
		}
		return false, retryAfter, nil
	}

	return true, 0, nil
}

// Increment increments the counter.
func (m *MemoryRateLimiter) Increment(ctx context.Context, key string, window time.Duration) error {
	c, exists := m.counters[key]
	if !exists {
		return nil
	}
	c.count++
	return nil
}

// cleanup periodically removes expired counters.
func (m *MemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, c := range m.counters {
			if now.After(c.expiresAt) {
				delete(m.counters, key)
			}
		}
	}
}
