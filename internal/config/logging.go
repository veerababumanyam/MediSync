// Package config provides configuration management for MediSync.
// This file handles structured logging with slog.
package config

import (
	"context"
	"log/slog"
	"os"
)

// contextKey is a type for context keys in this package.
type contextKey string

const (
	// RequestIDKey is the context key for request ID.
	RequestIDKey contextKey = "request_id"
)

// Logger wraps slog.Logger with additional functionality.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger based on the environment.
// In production, it outputs JSON format. In development, it outputs text format.
func NewLogger(env, level string) *Logger {
	var handler slog.Handler

	// Parse log level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Use JSON format for production, text for development
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithRequestID adds a request ID to the logger context.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("request_id", requestID),
	}
}

// WithContext creates a new logger with context values extracted.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return l.WithRequestID(requestID)
	}
	return l
}

// WithModule adds a module name to the logger.
func (l *Logger) WithModule(module string) *Logger {
	return &Logger{
		Logger: l.Logger.With("module", module),
	}
}

// WithAgent adds an agent ID to the logger.
func (l *Logger) WithAgent(agentID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("agent_id", agentID),
	}
}

// WithError adds an error to the logger.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With("error", err.Error()),
	}
}

// WithField adds a single field to the logger.
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		Logger: l.Logger.With(key, value),
	}
}

// WithFields adds multiple fields to the logger.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// LogAgentStart logs the start of an agent operation.
func (l *Logger) LogAgentStart(ctx context.Context, agentID, operation string, params map[string]any) {
	logger := l.WithContext(ctx).WithAgent(agentID)
	logger.Info("agent operation started",
		"operation", operation,
		"params", params,
	)
}

// LogAgentComplete logs the completion of an agent operation.
func (l *Logger) LogAgentComplete(ctx context.Context, agentID, operation string, confidence float64) {
	logger := l.WithContext(ctx).WithAgent(agentID)
	logger.Info("agent operation completed",
		"operation", operation,
		"confidence", confidence,
	)
}

// LogAgentError logs an error during an agent operation.
func (l *Logger) LogAgentError(ctx context.Context, agentID, operation string, err error) {
	logger := l.WithContext(ctx).WithAgent(agentID).WithError(err)
	logger.Error("agent operation failed",
		"operation", operation,
	)
}

// LogQuery logs a database query with context.
func (l *Logger) LogQuery(ctx context.Context, query string, args []any, durationMs int64) {
	logger := l.WithContext(ctx)
	logger.Debug("database query executed",
		"query", query,
		"args_count", len(args),
		"duration_ms", durationMs,
	)
}

// LogHTTPRequest logs an HTTP request.
func (l *Logger) LogHTTPRequest(ctx context.Context, method, path, statusCode string, durationMs int64) {
	logger := l.WithContext(ctx)
	logger.Info("http request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", durationMs,
	)
}

// Global logger instance
var globalLogger *Logger

// InitLogger initializes the global logger.
func InitLogger(env, level string) {
	globalLogger = NewLogger(env, level)
}

// L returns the global logger.
func L() *Logger {
	if globalLogger == nil {
		InitLogger("development", "info")
	}
	return globalLogger
}
