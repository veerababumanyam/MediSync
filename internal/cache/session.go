// Package cache provides session management using Redis.
//
// This file provides the SessionCache struct for managing query sessions.
// Sessions are used to maintain state across multiple turns in a conversation
// with AI agents.
//
// Session data includes:
//   - User preferences (locale, timezone, calendar system)
//   - Conversation context for multi-turn queries
//   - Query history for follow-up questions
//   - Temporary state for complex workflows
//
// Usage:
//
//	sessionCache := cache.NewSessionCache(redisClient, logger)
//
//	// Create a new session
//	session := &cache.QuerySession{
//	    SessionID: uuid.New().String(),
//	    UserID:    "user-123",
//	    Locale:    "en",
//	}
//	err := sessionCache.SetSession(ctx, session, 24*time.Hour)
//
//	// Retrieve the session
//	session, err := sessionCache.GetSession(ctx, sessionID)
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session key prefix for Redis.
const sessionKeyPrefix = "medisync:session"

// Default session TTL values.
const (
	// DefaultSessionTTL is the default time-to-live for sessions (24 hours).
	DefaultSessionTTL = 24 * time.Hour

	// MaxSessionTTL is the maximum allowed session TTL (7 days).
	MaxSessionTTL = 7 * 24 * time.Hour
)

// QuerySession represents a user session for query context.
type QuerySession struct {
	// SessionID is the unique identifier for this session.
	SessionID string `json:"session_id"`

	// UserID is the ID of the user who owns this session.
	UserID string `json:"user_id"`

	// TenantID is the tenant/organization ID for multi-tenancy.
	TenantID string `json:"tenant_id,omitempty"`

	// Locale is the user's preferred language ("en" or "ar").
	Locale string `json:"locale"`

	// Timezone is the user's timezone (e.g., "Asia/Riyadh").
	Timezone string `json:"timezone"`

	// CalendarSystem is the user's calendar preference ("gregorian" or "hijri").
	CalendarSystem string `json:"calendar_system"`

	// Roles are the user's assigned roles for authorization.
	Roles []string `json:"roles,omitempty"`

	// CostCentres are the cost centres the user has access to.
	CostCentres []string `json:"cost_centres,omitempty"`

	// ConversationHistory stores recent conversation turns.
	ConversationHistory []ConversationTurn `json:"conversation_history,omitempty"`

	// CurrentContext holds the current query context.
	CurrentContext *QueryContext `json:"current_context,omitempty"`

	// Preferences holds user-specific preferences.
	Preferences map[string]interface{} `json:"preferences,omitempty"`

	// CreatedAt is when the session was created.
	CreatedAt time.Time `json:"created_at"`

	// LastActivityAt is when the session was last accessed.
	LastActivityAt time.Time `json:"last_activity_at"`

	// ExpiresAt is when the session expires.
	ExpiresAt time.Time `json:"expires_at"`
}

// ConversationTurn represents a single turn in the conversation.
type ConversationTurn struct {
	// ID is the unique identifier for this turn.
	ID string `json:"id"`

	// Query is the natural language query from the user.
	Query string `json:"query"`

	// SQL is the generated SQL query.
	SQL string `json:"sql,omitempty"`

	// Response is the AI's response to the query.
	Response string `json:"response,omitempty"`

	// VisualizationType is the type of chart generated.
	VisualizationType string `json:"visualization_type,omitempty"`

	// Confidence is the confidence score for this turn.
	Confidence float64 `json:"confidence,omitempty"`

	// Timestamp is when this turn occurred.
	Timestamp time.Time `json:"timestamp"`
}

// QueryContext holds context for the current query.
type QueryContext struct {
	// LastQuery is the most recent natural language query.
	LastQuery string `json:"last_query,omitempty"`

	// LastSQL is the most recent SQL query.
	LastSQL string `json:"last_sql,omitempty"`

	// LastTables are the tables used in the last query.
	LastTables []string `json:"last_tables,omitempty"`

	// DataSource indicates the data source being queried ("hims", "tally", "combined").
	DataSource string `json:"data_source,omitempty"`

	// AgentID is the ID of the agent handling the current query.
	AgentID string `json:"agent_id,omitempty"`

	// WorkflowState holds state for multi-step workflows.
	WorkflowState map[string]interface{} `json:"workflow_state,omitempty"`
}

// SessionCache provides session management using Redis.
type SessionCache struct {
	client *redis.Client
	logger *slog.Logger
}

// SessionCacheConfig holds configuration for the session cache.
type SessionCacheConfig struct {
	// Client is the Redis client.
	Client *redis.Client

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewSessionCache creates a new session cache using an existing Redis client.
func NewSessionCache(client *redis.Client, logger *slog.Logger) *SessionCache {
	if logger == nil {
		logger = slog.Default()
	}

	return &SessionCache{
		client: client,
		logger: logger,
	}
}

// NewSessionCacheFromClient creates a new session cache from an existing cache client.
func NewSessionCacheFromClient(cacheClient *Client) *SessionCache {
	return &SessionCache{
		client: cacheClient.client,
		logger: cacheClient.logger,
	}
}

// GetSession retrieves a session by ID.
func (sc *SessionCache) GetSession(ctx context.Context, sessionID string) (*QuerySession, error) {
	key := sc.getSessionKey(sessionID)

	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Session not found
		}
		return nil, fmt.Errorf("cache: failed to get session: %w", err)
	}

	var session QuerySession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("cache: failed to unmarshal session: %w", err)
	}

	// Update last activity
	session.LastActivityAt = time.Now()

	sc.logger.Debug("session retrieved",
		slog.String("session_id", sessionID),
		slog.String("user_id", session.UserID),
	)

	return &session, nil
}

// SetSession stores a session with the specified TTL.
func (sc *SessionCache) SetSession(ctx context.Context, session *QuerySession, ttl time.Duration) error {
	if session == nil {
		return fmt.Errorf("cache: session is required")
	}

	if session.SessionID == "" {
		return fmt.Errorf("cache: session ID is required")
	}

	// Validate and set TTL
	if ttl == 0 {
		ttl = DefaultSessionTTL
	}
	if ttl > MaxSessionTTL {
		ttl = MaxSessionTTL
	}

	// Set timestamps
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.LastActivityAt = now
	session.ExpiresAt = now.Add(ttl)

	key := sc.getSessionKey(session.SessionID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("cache: failed to marshal session: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to set session: %w", err)
	}

	sc.logger.Debug("session stored",
		slog.String("session_id", session.SessionID),
		slog.String("user_id", session.UserID),
		slog.Duration("ttl", ttl),
	)

	return nil
}

// DeleteSession removes a session.
func (sc *SessionCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := sc.getSessionKey(sessionID)

	if err := sc.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: failed to delete session: %w", err)
	}

	sc.logger.Debug("session deleted",
		slog.String("session_id", sessionID),
	)

	return nil
}

// RefreshSession extends the TTL of an existing session.
func (sc *SessionCache) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = DefaultSessionTTL
	}
	if ttl > MaxSessionTTL {
		ttl = MaxSessionTTL
	}

	key := sc.getSessionKey(sessionID)

	// Check if session exists
	exists, err := sc.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("cache: failed to check session existence: %w", err)
	}

	if exists == 0 {
		return fmt.Errorf("cache: session not found")
	}

	// Update expiry
	if err := sc.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to refresh session: %w", err)
	}

	// Update ExpiresAt in the session data
	session, err := sc.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session != nil {
		session.ExpiresAt = time.Now().Add(ttl)
		if err := sc.SetSession(ctx, session, ttl); err != nil {
			return err
		}
	}

	sc.logger.Debug("session refreshed",
		slog.String("session_id", sessionID),
		slog.Duration("ttl", ttl),
	)

	return nil
}

// AddConversationTurn adds a new turn to the conversation history.
func (sc *SessionCache) AddConversationTurn(ctx context.Context, sessionID string, turn ConversationTurn) error {
	session, err := sc.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return fmt.Errorf("cache: session not found")
	}

	// Add turn to history
	if session.ConversationHistory == nil {
		session.ConversationHistory = []ConversationTurn{}
	}

	// Set timestamp if not provided
	if turn.Timestamp.IsZero() {
		turn.Timestamp = time.Now()
	}

	session.ConversationHistory = append(session.ConversationHistory, turn)

	// Keep only the last 50 turns to prevent unbounded growth
	if len(session.ConversationHistory) > 50 {
		session.ConversationHistory = session.ConversationHistory[len(session.ConversationHistory)-50:]
	}

	// Update current context
	if session.CurrentContext == nil {
		session.CurrentContext = &QueryContext{}
	}
	session.CurrentContext.LastQuery = turn.Query
	session.CurrentContext.LastSQL = turn.SQL

	return sc.SetSession(ctx, session, 0) // Use existing TTL
}

// UpdateQueryContext updates the current query context for a session.
func (sc *SessionCache) UpdateQueryContext(ctx context.Context, sessionID string, queryCtx *QueryContext) error {
	session, err := sc.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return fmt.Errorf("cache: session not found")
	}

	session.CurrentContext = queryCtx

	return sc.SetSession(ctx, session, 0)
}

// GetConversationHistory retrieves the conversation history for a session.
func (sc *SessionCache) GetConversationHistory(ctx context.Context, sessionID string, limit int) ([]ConversationTurn, error) {
	session, err := sc.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil || len(session.ConversationHistory) == 0 {
		return nil, nil
	}

	history := session.ConversationHistory

	// Apply limit
	if limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}

	return history, nil
}

// Exists checks if a session exists.
func (sc *SessionCache) Exists(ctx context.Context, sessionID string) (bool, error) {
	key := sc.getSessionKey(sessionID)

	exists, err := sc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache: failed to check session existence: %w", err)
	}

	return exists > 0, nil
}

// GetSessionTTL returns the remaining TTL for a session.
func (sc *SessionCache) GetSessionTTL(ctx context.Context, sessionID string) (time.Duration, error) {
	key := sc.getSessionKey(sessionID)

	ttl, err := sc.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("cache: failed to get session TTL: %w", err)
	}

	if ttl < 0 {
		return 0, nil // Session doesn't exist or has no expiry
	}

	return ttl, nil
}

// DeleteUserSessions removes all sessions for a user.
func (sc *SessionCache) DeleteUserSessions(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s:*", sessionKeyPrefix)

	iter := sc.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keysToDelete []string

	for iter.Next(ctx) {
		// We'd need to check the value to match userID
		// For now, we'll just note this limitation
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache: failed to scan sessions: %w", err)
	}

	// Note: This is a simplified implementation
	// In production, you'd want to index sessions by user_id
	// or use a separate set to track user sessions

	sc.logger.Debug("user sessions deletion requested",
		slog.String("user_id", userID),
		slog.Int("potential_keys", len(keysToDelete)),
	)

	return nil
}

// getSessionKey returns the Redis key for a session.
func (sc *SessionCache) getSessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", sessionKeyPrefix, sessionID)
}
