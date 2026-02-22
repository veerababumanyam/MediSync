// Package cache provides evidence caching for the Council of AIs consensus system.
//
// The EvidenceCache stores Knowledge Graph traversal results to support brief KG
// outages without compromising accuracy. Per FR-017, evidence is cached for up to
// 5 minutes with automatic expiration.
//
// Cache Structure:
//   - Key: "evidence:{query_hash}" - SHA-256 hash of the query for deduplication
//   - Value: EvidenceCacheEntry containing nodes, traversal path, and relevance scores
//   - TTL: 5 minutes (configurable via COUNCIL_CACHE_TTL)
//
// Usage:
//
//	cache := cache.NewEvidenceCache(redisClient, logger)
//	defer cache.Close()
//
//	ctx := context.Background()
//
//	// Cache evidence for a query
//	entry := &cache.EvidenceCacheEntry{...}
//	err := cache.Set(ctx, queryHash, entry)
//
//	// Retrieve cached evidence
//	cached, err := cache.Get(ctx, queryHash)
//	if err == nil && cached != nil {
//	    // Use cached evidence
//	}
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/medisync/medisync/internal/agents/council"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/redis/go-redis/v9"
)

// Evidence cache key prefix and defaults
const (
	// KeyEvidence is the prefix for evidence cache keys.
	KeyEvidence = "evidence"

	// TTLEvidence is the default TTL for evidence cache (5 minutes per FR-017).
	TTLEvidence = 5 * time.Minute

	// TTLEvidenceExtended is the extended TTL during KG outages (15 minutes).
	TTLEvidenceExtended = 15 * time.Minute
)

// EvidenceCacheEntry represents cached Knowledge Graph traversal results.
type EvidenceCacheEntry struct {
	// QueryHash is the SHA-256 hash of the original query.
	QueryHash string `json:"query_hash"`

	// QueryText is the original query (optional, for debugging).
	QueryText string `json:"query_text,omitempty"`

	// NodeIDs are the unique IDs of all KG nodes in the evidence.
	NodeIDs []string `json:"node_ids"`

	// Nodes contains the actual Knowledge Graph node data.
	Nodes []*warehouse.KnowledgeGraphNode `json:"nodes"`

	// TraversalPath records the multi-hop path through the graph.
	TraversalPath []*warehouse.TraversalStep `json:"traversal_path"`

	// RelevanceScores maps node IDs to their relevance scores (0-1).
	RelevanceScores map[string]float64 `json:"relevance_scores"`

	// HopCount is the number of hops in the traversal.
	HopCount int `json:"hop_count"`

	// CachedAt is when this evidence was cached.
	CachedAt time.Time `json:"cached_at"`

	// ExpiresAt is when this cache entry expires.
	ExpiresAt time.Time `json:"expires_at"`

	// Source identifies the source of the evidence (e.g., "knowledge_graph").
	Source string `json:"source"`

	// KGHealthStatus captures the health of the KG when evidence was cached.
	// Used to determine if cached evidence is trustworthy during outages.
	KGHealthStatus string `json:"kg_health_status"`
}

// EvidenceCache provides Redis-based caching for Council evidence.
type EvidenceCache struct {
	client *redis.Client
	logger *slog.Logger
	ttl    time.Duration
}

// NewEvidenceCache creates a new evidence cache.
func NewEvidenceCache(client *redis.Client, logger *slog.Logger, ttl time.Duration) *EvidenceCache {
	if ttl == 0 {
		ttl = TTLEvidence
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &EvidenceCache{
		client: client,
		logger: logger,
		ttl:    ttl,
	}
}

// HashQuery generates a SHA-256 hash of the query for cache key.
func HashQuery(query string) string {
	h := sha256.Sum256([]byte(query))
	return hex.EncodeToString(h[:])
}

// buildKey constructs the full cache key for a query hash.
func (c *EvidenceCache) buildKey(queryHash string) string {
	return fmt.Sprintf("%s:%s", KeyEvidence, queryHash)
}

// Set stores evidence in the cache.
func (c *EvidenceCache) Set(ctx context.Context, queryHash string, entry *EvidenceCacheEntry) error {
	key := c.buildKey(queryHash)

	// Set timestamps
	now := time.Now()
	entry.CachedAt = now
	entry.ExpiresAt = now.Add(c.ttl)
	entry.QueryHash = queryHash

	// Marshal entry
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("evidence_cache: failed to marshal entry: %w", err)
	}

	// Store with TTL
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("evidence_cache: failed to set entry: %w", err)
	}

	c.logger.Debug("cached evidence",
		slog.String("query_hash", queryHash),
		slog.Int("node_count", len(entry.NodeIDs)),
		slog.Int("hop_count", entry.HopCount),
		slog.Duration("ttl", c.ttl),
	)

	return nil
}

// Get retrieves evidence from the cache.
func (c *EvidenceCache) Get(ctx context.Context, queryHash string) (*EvidenceCacheEntry, error) {
	key := c.buildKey(queryHash)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("evidence_cache: failed to get entry: %w", err)
	}

	var entry EvidenceCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("evidence_cache: failed to unmarshal entry: %w", err)
	}

	// Check if expired (belt and suspenders - Redis TTL should handle this)
	if time.Now().After(entry.ExpiresAt) {
		// Clean up expired entry
		_ = c.client.Del(ctx, key)
		return nil, nil
	}

	c.logger.Debug("retrieved cached evidence",
		slog.String("query_hash", queryHash),
		slog.Int("node_count", len(entry.NodeIDs)),
		slog.String("age", time.Since(entry.CachedAt).String()),
	)

	return &entry, nil
}

// GetOrCompute retrieves cached evidence or computes and caches new evidence.
// This is useful for graceful degradation during KG outages.
func (c *EvidenceCache) GetOrCompute(
	ctx context.Context,
	query string,
	computeFn func() (*EvidenceCacheEntry, error),
) (*EvidenceCacheEntry, error) {
	queryHash := HashQuery(query)

	// Try to get from cache first
	cached, err := c.Get(ctx, queryHash)
	if err != nil {
		c.logger.Warn("cache read error, computing fresh",
			slog.String("query_hash", queryHash),
			slog.Any("error", err),
		)
	} else if cached != nil {
		return cached, nil
	}

	// Compute fresh evidence
	entry, err := computeFn()
	if err != nil {
		// If compute fails but we have stale cache, return it with a warning
		if cached != nil {
			c.logger.Warn("using stale cached evidence due to compute failure",
				slog.String("query_hash", queryHash),
				slog.Any("compute_error", err),
				slog.Duration("age", time.Since(cached.CachedAt)),
			)
			return cached, nil
		}
		return nil, fmt.Errorf("evidence_cache: compute failed: %w", err)
	}

	// Cache the new entry
	if err := c.Set(ctx, queryHash, entry); err != nil {
		c.logger.Warn("failed to cache computed evidence",
			slog.String("query_hash", queryHash),
			slog.Any("error", err),
		)
		// Still return the computed entry even if caching failed
	}

	return entry, nil
}

// Delete removes evidence from the cache.
func (c *EvidenceCache) Delete(ctx context.Context, queryHash string) error {
	key := c.buildKey(queryHash)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("evidence_cache: failed to delete entry: %w", err)
	}

	c.logger.Debug("deleted cached evidence",
		slog.String("query_hash", queryHash),
	)

	return nil
}

// Exists checks if evidence exists in the cache.
func (c *EvidenceCache) Exists(ctx context.Context, queryHash string) (bool, error) {
	key := c.buildKey(queryHash)

	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("evidence_cache: failed to check existence: %w", err)
	}

	return count > 0, nil
}

// ExtendTTL extends the TTL for an entry (useful during KG outages).
func (c *EvidenceCache) ExtendTTL(ctx context.Context, queryHash string, extension time.Duration) error {
	key := c.buildKey(queryHash)

	// Check if key exists first
	exists, err := c.Exists(ctx, queryHash)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("evidence_cache: entry not found")
	}

	// Extend TTL
	newTTL := c.ttl + extension
	if newTTL > TTLEvidenceExtended {
		newTTL = TTLEvidenceExtended
	}

	if err := c.client.Expire(ctx, key, newTTL).Err(); err != nil {
		return fmt.Errorf("evidence_cache: failed to extend TTL: %w", err)
	}

	c.logger.Debug("extended evidence TTL",
		slog.String("query_hash", queryHash),
		slog.Duration("new_ttl", newTTL),
	)

	return nil
}

// InvalidatePattern removes all evidence entries matching a pattern.
func (c *EvidenceCache) InvalidatePattern(ctx context.Context, pattern string) error {
	key := fmt.Sprintf("%s:%s", KeyEvidence, pattern)

	iter := c.client.Scan(ctx, 0, key, 100).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("evidence_cache: failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("evidence_cache: failed to delete keys: %w", err)
		}

		c.logger.Debug("invalidated evidence pattern",
			slog.String("pattern", pattern),
			slog.Int("count", len(keys)),
		)
	}

	return nil
}

// InvalidateAll removes all evidence entries.
func (c *EvidenceCache) InvalidateAll(ctx context.Context) error {
	return c.InvalidatePattern(ctx, "*")
}

// GetStats returns evidence cache statistics.
type EvidenceCacheStats struct {
	TotalEntries int64 `json:"total_entries"`
	HitCount     int64 `json:"hit_count"`
	MissCount    int64 `json:"miss_count"`
	HitRate      float64 `json:"hit_rate"`
}

// GetStats retrieves evidence cache statistics.
func (c *EvidenceCache) GetStats(ctx context.Context) (*EvidenceCacheStats, error) {
	// Count evidence keys
	var cursor uint64
	var count int64

	for {
		var keys []string
		var err error
		keys, cursor, err = c.client.Scan(ctx, cursor, KeyEvidence+":*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("evidence_cache: failed to scan keys: %w", err)
		}
		count += int64(len(keys))

		if cursor == 0 {
			break
		}
	}

	return &EvidenceCacheStats{
		TotalEntries: count,
	}, nil
}

// IsExpired checks if a cached entry is expired.
func (e *EvidenceCacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Age returns how long ago the entry was cached.
func (e *EvidenceCacheEntry) Age() time.Duration {
	return time.Since(e.CachedAt)
}

// RemainingTTL returns the time until the entry expires.
func (e *EvidenceCacheEntry) RemainingTTL() time.Duration {
	return time.Until(e.ExpiresAt)
}

// IsFresh checks if the entry is fresh (cached within the last minute).
func (e *EvidenceCacheEntry) IsFresh() bool {
	return e.Age() < time.Minute
}

// IsTrustworthyDuringOutage checks if cached evidence is trustworthy during a KG outage.
// Evidence is trustworthy if:
// - It was cached when KG was healthy
// - It hasn't expired
// - It's not too old (within extended TTL)
func (e *EvidenceCacheEntry) IsTrustworthyDuringOutage() bool {
	return e.KGHealthStatus == "healthy" &&
		!e.IsExpired() &&
		e.Age() < TTLEvidenceExtended
}

// ToEvidenceTrail converts the cache entry to an EvidenceTrail for the Council.
func (e *EvidenceCacheEntry) ToEvidenceTrail(deliberationID string) *council.EvidenceTrail {
	return &council.EvidenceTrail{
		ID:             "", // Will be set by repository
		DeliberationID: deliberationID,
		NodeIDs:        e.NodeIDs,
		TraversalPath:  e.TraversalPath,
		RelevanceScore: e.RelevanceScores,
		HopCount:       e.HopCount,
		CachedAt:       e.CachedAt,
		ExpiresAt:      e.ExpiresAt,
	}
}
