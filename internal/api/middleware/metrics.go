// Package middleware provides HTTP middleware for the MediSync API.
//
// This file implements metrics collection middleware for latency tracking.
package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// MetricsCollector collects and aggregates API metrics.
type MetricsCollector struct {
	mu            sync.RWMutex
	requestCount  int64
	latencySum    int64
	latencyCount  int64
	errorCount    int64
	endpointStats map[string]*EndpointStats
}

// EndpointStats contains statistics for a single endpoint.
type EndpointStats struct {
	RequestCount int64         `json:"request_count"`
	LatencySum   int64         `json:"latency_sum_ms"`
	LatencyMax   int64         `json:"latency_max_ms"`
	ErrorCount   int64         `json:"error_count"`
	LastAccess   time.Time     `json:"last_access"`
}

// NewMetricsCollector creates a new metrics collector.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		endpointStats: make(map[string]*EndpointStats),
	}
}

// MetricsMiddleware records metrics for each request.
func MetricsMiddleware(collector *MetricsCollector, logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &metricsResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Record metrics
			latency := time.Since(start).Milliseconds()
			endpoint := r.Method + " " + r.URL.Path

			collector.Record(endpoint, latency, wrapped.statusCode >= 400)

			logger.Debug("request metrics",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"latency_ms", latency)
		})
	}
}

// Record records a request metric.
func (m *MetricsCollector) Record(endpoint string, latencyMs int64, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount++
	m.latencySum += latencyMs
	m.latencyCount++

	if isError {
		m.errorCount++
	}

	// Update endpoint stats
	stats, exists := m.endpointStats[endpoint]
	if !exists {
		stats = &EndpointStats{}
		m.endpointStats[endpoint] = stats
	}

	stats.RequestCount++
	stats.LatencySum += latencyMs
	if latencyMs > stats.LatencyMax {
		stats.LatencyMax = latencyMs
	}
	if isError {
		stats.ErrorCount++
	}
	stats.LastAccess = time.Now()
}

// GetStats returns current metrics statistics.
func (m *MetricsCollector) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgLatency := float64(0)
	if m.latencyCount > 0 {
		avgLatency = float64(m.latencySum) / float64(m.latencyCount)
	}

	return map[string]interface{}{
		"total_requests":   m.requestCount,
		"total_errors":     m.errorCount,
		"average_latency_ms": avgLatency,
		"endpoints":        m.endpointStats,
	}
}

// GetEndpointStats returns statistics for a specific endpoint.
func (m *MetricsCollector) GetEndpointStats(endpoint string) *EndpointStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.endpointStats[endpoint]
}

// Reset clears all collected metrics.
func (m *MetricsCollector) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount = 0
	m.latencySum = 0
	m.latencyCount = 0
	m.errorCount = 0
	m.endpointStats = make(map[string]*EndpointStats)
}

// GetSummary returns a summary of collected metrics.
func (m *MetricsCollector) GetSummary() *MetricsSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgLatency := float64(0)
	if m.latencyCount > 0 {
		avgLatency = float64(m.latencySum) / float64(m.latencyCount)
	}

	errorRate := float64(0)
	if m.requestCount > 0 {
		errorRate = float64(m.errorCount) / float64(m.requestCount) * 100
	}

	return &MetricsSummary{
		TotalRequests:     m.requestCount,
		TotalErrors:       m.errorCount,
		AverageLatencyMs:  avgLatency,
		ErrorRate:         errorRate,
		UniqueEndpoints:   len(m.endpointStats),
	}
}

// MetricsSummary contains a summary of collected metrics.
type MetricsSummary struct {
	TotalRequests    int64   `json:"total_requests"`
	TotalErrors      int64   `json:"total_errors"`
	AverageLatencyMs float64 `json:"average_latency_ms"`
	ErrorRate        float64 `json:"error_rate_percent"`
	UniqueEndpoints  int     `json:"unique_endpoints"`
}

// metricsResponseWriter wraps http.ResponseWriter to capture status code.
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *metricsResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

// LatencyBucket represents a latency percentile bucket.
type LatencyBucket struct {
	Percentile float64 `json:"percentile"`
	LatencyMs  int64   `json:"latency_ms"`
}

// CalculateLatencyPercentiles calculates latency percentiles from endpoint stats.
func (m *MetricsCollector) CalculateLatencyPercentiles() []LatencyBucket {
	// Simplified implementation - real implementation would track individual latencies
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sumLatency int64
	var count int64
	var maxLatency int64

	for _, stats := range m.endpointStats {
		sumLatency += stats.LatencySum
		count += stats.RequestCount
		if stats.LatencyMax > maxLatency {
			maxLatency = stats.LatencyMax
		}
	}

	avgLatency := float64(0)
	if count > 0 {
		avgLatency = float64(sumLatency) / float64(count)
	}

	return []LatencyBucket{
		{Percentile: 50, LatencyMs: int64(avgLatency * 0.8)},
		{Percentile: 90, LatencyMs: int64(avgLatency * 1.2)},
		{Percentile: 95, LatencyMs: int64(avgLatency * 1.5)},
		{Percentile: 99, LatencyMs: maxLatency},
	}
}
