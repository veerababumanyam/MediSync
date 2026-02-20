// Package cache provides Redis-based caching for MediSync.
//
// This package handles caching of schema context, session data, and query results
// to improve performance and reduce database load. It uses go-redis/v9 for Redis operations.
//
// Cache keys follow a naming convention: `namespace:id` for easy identification.
// All cached values have TTL (time-to-live) to prevent stale data.
//
// Usage:
//
//	cfg := config.MustLoad()
//	cache := cache.New(cfg.Redis, logger)
//	defer cache.Close()
//
//	ctx := context.Background()
//	err := cache.SetSchemaContext(ctx, "hims_analytics", "dim_patients", contextData)
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache key prefixes
const (
	// KeySchema is the prefix for schema context cache keys.
	KeySchema = "schema"
	// KeySession is the prefix for session cache keys.
	KeySession = "session"
	// KeyQuery is the prefix for query result cache keys.
	KeyQuery = "query"
	// KeyAgent is the prefix for agent state cache keys.
	KeyAgent = "agent"
)

// Default TTL values for different cache types.
const (
	// TTLSchema is the default TTL for schema context (1 hour).
	TTLSchema = 1 * time.Hour
	// TTLSession is the default TTL for session data (24 hours).
	TTLSession = 24 * time.Hour
	// TTLQuery is the default TTL for query results (5 minutes).
	TTLQuery = 5 * time.Minute
	// TTLAgent is the default TTL for agent state (30 minutes).
	TTLAgent = 30 * time.Minute
)

// Client provides Redis caching operations.
type Client struct {
	client *redis.Client
	logger *slog.Logger
}

// ClientConfig holds configuration for creating a new Redis client.
type ClientConfig struct {
	// Addr is the Redis server address (host:port).
	Addr string

	// Password is the Redis password (optional).
	Password string

	// DB is the Redis database number.
	DB int

	// PoolSize is the connection pool size.
	PoolSize int

	// MinIdleConns is the minimum number of idle connections.
	MinIdleConns int

	// DialTimeout is the connection timeout.
	DialTimeout time.Duration

	// ReadTimeout is the read operation timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the write operation timeout.
	WriteTimeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewClient creates a new Redis cache client.
func NewClient(cfg interface{}, logger *slog.Logger) (*Client, error) {
	// Parse config
	var addr, password string
	var db int

	switch c := cfg.(type) {
	case map[string]interface{}:
		if url, ok := c["url"].(string); ok && url != "" {
			// Parse from URL
			addr, password, db = parseRedisURL(url)
		} else {
			if host, ok := c["host"].(string); ok {
				addr = fmt.Sprintf("%s:%d", host, int(c["port"].(float64)))
			}
			if pass, ok := c["password"].(string); ok {
				password = pass
			}
			if dbNum, ok := c["database"].(float64); ok {
				db = int(dbNum)
			}
		}
	case string:
		addr, password, db = parseRedisURL(c)
	}

	if addr == "" {
		addr = "localhost:6379"
	}

	if logger == nil {
		logger = slog.Default()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("cache: failed to connect to Redis: %w", err)
	}

	logger.Info("connected to Redis",
		slog.String("addr", addr),
		slog.Int("db", db),
	)

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

// parseRedisURL parses a Redis URL and returns addr, password, and db.
func parseRedisURL(url string) (addr, password string, db int) {
	// Simple URL parsing for redis://[:password@]host:port[/db]
	// In production, use a proper URL parser
	if len(url) > 9 && url[:9] == "redis://" {
		url = url[9:]
	}

	// Extract password if present
	passwordEnd := -1
	for idx := 0; idx < len(url) && url[idx] != '@'; idx++ {
		if url[idx] == ':' && passwordEnd == -1 {
			passwordEnd = idx
		}
	}
	if passwordEnd >= 0 && passwordEnd < len(url) && url[passwordEnd+1] == '@' {
		password = url[:passwordEnd]
		url = url[passwordEnd+2:]
	}

	// Extract db if present
	dbStart := -1
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			dbStr := url[i+1:]
			if dbStr != "" {
				fmt.Sscanf(dbStr, "%d", &db)
			}
			dbStart = i
			break
		}
	}
	if dbStart >= 0 {
		url = url[:dbStart]
	}

	addr = url
	return
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Ping checks if the Redis connection is alive.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// FlushAll clears all cached data (use with caution).
func (c *Client) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

// ============================================================================
// Schema Context Caching (for A-01 Text-to-SQL)
// ============================================================================

// SchemaContext represents cached schema metadata for a table.
type SchemaContext struct {
	SchemaName      string                   `json:"schema_name"`
	TableName       string                   `json:"table_name"`
	Columns         []ColumnContext          `json:"columns"`
	Relationships   []RelationshipContext    `json:"relationships"`
	DescriptionEN   string                   `json:"description_en"`
	DescriptionAR   string                   `json:"description_ar"`
	BusinessContext string                   `json:"business_context"`
	SampleRows      []map[string]interface{} `json:"sample_rows,omitempty"`
	LastModified    time.Time                `json:"last_modified"`
}

// ColumnContext represents a column in a schema.
type ColumnContext struct {
	Name          string   `json:"name"`
	DataType      string   `json:"data_type"`
	IsNullable    bool     `json:"is_nullable"`
	DescriptionEN string   `json:"description_en"`
	DescriptionAR string   `json:"description_ar"`
	SampleValues  []string `json:"sample_values,omitempty"`
}

// RelationshipContext represents a relationship between tables.
type RelationshipContext struct {
	FromTable  string `json:"from_table"`
	FromColumn string `json:"from_column"`
	ToTable    string `json:"to_table"`
	ToColumn   string `json:"to_column"`
	Type       string `json:"type"` // one_to_one, one_to_many, many_to_many
}

// SetSchemaContext stores schema context for a table.
func (c *Client) SetSchemaContext(ctx context.Context, schema, table string, context *SchemaContext) error {
	key := fmt.Sprintf("%s:%s:%s", KeySchema, schema, table)

	if err := c.SetStruct(ctx, key, context, TTLSchema); err != nil {
		return err
	}

	c.logger.Debug("cached schema context",
		slog.String("schema", schema),
		slog.String("table", table),
	)

	return nil
}

// GetSchemaContext retrieves schema context for a table.
func (c *Client) GetSchemaContext(ctx context.Context, schema, table string) (*SchemaContext, error) {
	key := fmt.Sprintf("%s:%s:%s", KeySchema, schema, table)

	var context SchemaContext
	if err := c.GetStruct(ctx, key, &context); err != nil {
		return nil, err
	}

	if context.TableName == "" {
		return nil, nil // Not found
	}

	return &context, nil
}

// InvalidateSchemaContext removes cached schema context for a table.
func (c *Client) InvalidateSchemaContext(ctx context.Context, schema, table string) error {
	key := fmt.Sprintf("%s:%s:%s", KeySchema, schema, table)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: failed to invalidate schema context: %w", err)
	}

	c.logger.Debug("invalidated schema context",
		slog.String("schema", schema),
		slog.String("table", table),
	)

	return nil
}

// InvalidateSchemaPattern removes all cached schema context matching a pattern.
func (c *Client) InvalidateSchemaPattern(ctx context.Context, pattern string) error {
	key := fmt.Sprintf("%s:%s", KeySchema, pattern)

	iter := c.client.Scan(ctx, 0, key, 100).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache: failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache: failed to delete keys: %w", err)
		}

		c.logger.Debug("invalidated schema pattern",
			slog.String("pattern", pattern),
			slog.Int("count", len(keys)),
		)
	}

	return nil
}

// ============================================================================
// Session Caching
// ============================================================================

// SessionData represents cached user session data.
type SessionData struct {
	SessionID      string                 `json:"session_id"`
	UserID         string                 `json:"user_id"`
	Locale         string                 `json:"locale"`
	Timezone       string                 `json:"timezone"`
	CalendarSystem string                 `json:"calendar_system"`
	Role           string                 `json:"role"`
	CostCentres    []string               `json:"cost_centres"`
	Preferences    map[string]interface{} `json:"preferences"`
	CreatedAt      time.Time              `json:"created_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
	LastActivity   time.Time              `json:"last_activity"`
}

// SetSession stores session data.
func (c *Client) SetSession(ctx context.Context, sessionID string, data *SessionData) error {
	key := fmt.Sprintf("%s:%s", KeySession, sessionID)

	data.LastActivity = time.Now()

	// Calculate TTL based on expires_at
	ttl := time.Until(data.ExpiresAt)
	if ttl <= 0 {
		ttl = TTLSession
	}

	if err := c.SetStruct(ctx, key, data, ttl); err != nil {
		return err
	}

	c.logger.Debug("cached session",
		slog.String("session_id", sessionID),
		slog.Duration("ttl", ttl),
	)

	return nil
}

// GetSession retrieves session data.
func (c *Client) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("%s:%s", KeySession, sessionID)

	var session SessionData
	if err := c.GetStruct(ctx, key, &session); err != nil {
		return nil, err
	}

	if session.SessionID == "" {
		return nil, nil // Not found
	}

	// Update last activity
	session.LastActivity = time.Now()

	return &session, nil
}

// DeleteSession removes session data.
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("%s:%s", KeySession, sessionID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: failed to delete session: %w", err)
	}

	c.logger.Debug("deleted session",
		slog.String("session_id", sessionID),
	)

	return nil
}

// ============================================================================
// Query Result Caching
// ============================================================================

// QueryResult represents a cached query result.
type QueryResult struct {
	Query      string                   `json:"query"`
	Params     map[string]interface{}   `json:"params"`
	Columns    []string                 `json:"columns"`
	Rows       []map[string]interface{} `json:"rows"`
	RowCount   int                      `json:"row_count"`
	ExecutedAt time.Time                `json:"executed_at"`
	DurationMs int64                    `json:"duration_ms"`
}

// SetQueryResult stores a query result.
func (c *Client) SetQueryResult(ctx context.Context, queryID string, result *QueryResult, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", KeyQuery, queryID)

	if ttl == 0 {
		ttl = TTLQuery
	}

	if err := c.SetStruct(ctx, key, result, ttl); err != nil {
		return err
	}

	c.logger.Debug("cached query result",
		slog.String("query_id", queryID),
		slog.Duration("ttl", ttl),
	)

	return nil
}

// GetQueryResult retrieves a cached query result.
func (c *Client) GetQueryResult(ctx context.Context, queryID string) (*QueryResult, error) {
	key := fmt.Sprintf("%s:%s", KeyQuery, queryID)

	var result QueryResult
	if err := c.GetStruct(ctx, key, &result); err != nil {
		return nil, err
	}

	if result.Query == "" {
		return nil, nil // Not found
	}

	return &result, nil
}

// InvalidateQueryPattern removes cached query results matching a pattern.
func (c *Client) InvalidateQueryPattern(ctx context.Context, pattern string) error {
	key := fmt.Sprintf("%s:%s", KeyQuery, pattern)

	iter := c.client.Scan(ctx, 0, key, 100).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache: failed to scan query keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache: failed to delete query keys: %w", err)
		}

		c.logger.Debug("invalidated query pattern",
			slog.String("pattern", pattern),
			slog.Int("count", len(keys)),
		)
	}

	return nil
}

// ============================================================================
// Agent State Caching
// ============================================================================

// AgentState represents cached agent state for multi-turn conversations.
type AgentState struct {
	AgentID   string                 `json:"agent_id"`
	SessionID string                 `json:"session_id"`
	State     map[string]interface{} `json:"state"`
	Context   []string               `json:"context"` // Conversation history
	Metadata  map[string]interface{} `json:"metadata"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// SetAgentState stores agent state.
func (c *Client) SetAgentState(ctx context.Context, agentID, sessionID string, state *AgentState) error {
	key := fmt.Sprintf("%s:%s:%s", KeyAgent, agentID, sessionID)

	state.UpdatedAt = time.Now()

	if err := c.SetStruct(ctx, key, state, TTLAgent); err != nil {
		return err
	}

	c.logger.Debug("cached agent state",
		slog.String("agent_id", agentID),
		slog.String("session_id", sessionID),
	)

	return nil
}

// GetAgentState retrieves agent state.
func (c *Client) GetAgentState(ctx context.Context, agentID, sessionID string) (*AgentState, error) {
	key := fmt.Sprintf("%s:%s:%s", KeyAgent, agentID, sessionID)

	var state AgentState
	if err := c.GetStruct(ctx, key, &state); err != nil {
		return nil, err
	}

	if state.AgentID == "" {
		return nil, nil // Not found
	}

	return &state, nil
}

// DeleteAgentState removes agent state.
func (c *Client) DeleteAgentState(ctx context.Context, agentID, sessionID string) error {
	key := fmt.Sprintf("%s:%s:%s", KeyAgent, agentID, sessionID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: failed to delete agent state: %w", err)
	}

	c.logger.Debug("deleted agent state",
		slog.String("agent_id", agentID),
		slog.String("session_id", sessionID),
	)

	return nil
}

// ============================================================================
// Generic Operations
// ============================================================================

// Set stores a string value with a key and TTL.
func (c *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to set value: %w", err)
	}
	return nil
}

// SetStruct stores an object as JSON with a key and TTL.
func (c *Client) SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to set value: %w", err)
	}

	return nil
}

// Get retrieves a string value by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Not found
		}
		return "", fmt.Errorf("cache: failed to get value: %w", err)
	}
	return val, nil
}

// GetStruct retrieves a value by key and unmarshals it into dest.
func (c *Client) GetStruct(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Not found
		}
		return fmt.Errorf("cache: failed to get value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("cache: failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key.
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: failed to delete key: %w", err)
	}
	return nil
}

// Increment increments a key and returns the new value.
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("cache: failed to increment key: %w", err)
	}
	return val, nil
}

// Expire sets a TTL on a key.
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := c.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to set expiry: %w", err)
	}
	return nil
}

// Exists checks if a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache: failed to check key existence: %w", err)
	}
	return count > 0, nil
}

// GetStats returns Redis cache statistics.
type CacheStats struct {
	TotalKeys      int64   `json:"total_keys"`
	MemoryUsed     string  `json:"memory_used"`
	MemoryPeak     string  `json:"memory_peak"`
	HitRate        float64 `json:"hit_rate"`
	KeyspaceHits   int64   `json:"keyspace_hits"`
	KeyspaceMisses int64   `json:"keyspace_misses"`
}

// GetStats retrieves cache statistics.
func (c *Client) GetStats(ctx context.Context) (*CacheStats, error) {
	info, err := c.client.Info(ctx, "memory", "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("cache: failed to get info: %w", err)
	}

	// Parse info response (simplified)
	stats := &CacheStats{}

	// Get total keys for current database
	dbSize := c.client.DBSize(ctx)
	stats.TotalKeys, _ = dbSize.Result()

	// Parse memory info (simplified - in production use proper parsing)
	lines := parseInfo(info)
	if mem, ok := lines["used_memory_human"]; ok {
		stats.MemoryUsed = mem
	}
	if mem, ok := lines["used_memory_peak_human"]; ok {
		stats.MemoryPeak = mem
	}

	// Get hit/miss stats
	if hits, ok := lines["keyspace_hits"]; ok {
		fmt.Sscanf(hits, "%d", &stats.KeyspaceHits)
	}
	if misses, ok := lines["keyspace_misses"]; ok {
		fmt.Sscanf(misses, "%d", &stats.KeyspaceMisses)
	}

	total := stats.KeyspaceHits + stats.KeyspaceMisses
	if total > 0 {
		stats.HitRate = float64(stats.KeyspaceHits) / float64(total) * 100
	}

	return stats, nil
}

// parseInfo parses Redis INFO response into a map.
func parseInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := splitLines(info)

	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := splitColon(line, 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			if i > start {
				lines = append(lines, s[start:i])
			}
			start = i + 1
		}
	}

	if start < len(s) {
		lines = append(lines, s[start:])
	}

	return lines
}

func splitColon(s string, n int) []string {
	var parts []string
	start := 0
	count := 0

	for i := 0; i < len(s) && count < n-1; i++ {
		if s[i] == ':' {
			parts = append(parts, s[start:i])
			start = i + 1
			count++
		}
	}

	parts = append(parts, s[start:])
	return parts
}
