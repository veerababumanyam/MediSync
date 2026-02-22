# Go Error Handling Patterns

This reference provides comprehensive patterns for handling errors in Go applications.

## Core Principles

1. **Errors are values** - Treat them like any other value
2. **Return early** - Check errors immediately and return
3. **Add context** - Wrap errors with meaningful information
4. **Don't panic** - Use error returns, not panics, for expected failures

## Basic Patterns

### Immediate Check Pattern

```go
func ProcessFile(path string) ([]byte, error) {
    // Check error immediately after the operation
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    result, err := transform(data)
    if err != nil {
        return nil, err
    }

    return result, nil
}
```

### Error Wrapping with Context

```go
import "fmt"

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        // Wrap with context about what operation failed
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    return user, nil
}

// Error chain example:
// "get user 123: find by id: query: connection refused"
```

### %w vs %v in fmt.Errorf

```go
// Use %w when you want callers to be able to unwrap the error
if err != nil {
    return fmt.Errorf("operation failed: %w", err)  // Allows errors.Is/As
}

// Use %v when you just want to add context without unwrapping
if err != nil {
    return fmt.Errorf("operation failed: %v", err)  // No unwrapping
}

// Use %w only once per error chain to avoid ambiguity
```

## Custom Error Types

### Sentinel Errors

```go
package errors

// Define at package level for comparison
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrAlreadyExists = errors.New("already exists")
)

// Usage
func (r *Repository) Find(id string) (*Entity, error) {
    entity, ok := r.data[id]
    if !ok {
        return nil, ErrNotFound
    }
    return entity, nil
}
```

### Structured Error Types

```go
package errors

type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: field %s has invalid value %v: %s",
        e.Field, e.Value, e.Message)
}

func (e *ValidationError) Unwrap() error {
    return nil // No underlying error
}

// Constructor
func NewValidationError(field string, value interface{}, message string) error {
    return &ValidationError{
        Field:   field,
        Value:   value,
        Message: message,
    }
}
```

### Error with Code

```go
package errors

type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Predefined error codes
const (
    CodeInvalidInput = "INVALID_INPUT"
    CodeNotFound     = "NOT_FOUND"
    CodeInternal     = "INTERNAL_ERROR"
)

// Constructor
func NewAppError(code, message string, err error) error {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}
```

## Error Inspection

### errors.Is - Check Exact Error

```go
import "errors"

func HandleRequest(id string) error {
    user, err := svc.GetUser(id)
    if err != nil {
        // Check if it's a specific error (works through wrapped errors)
        if errors.Is(err, ErrNotFound) {
            return fmt.Errorf("user not found: %s", id)
        }
        return fmt.Errorf("failed to get user: %w", err)
    }
    return nil
}
```

### errors.As - Extract Error Type

```go
func HandleRequest(data string) error {
    err := validate(data)
    if err != nil {
        // Try to extract ValidationError
        var valErr *ValidationError
        if errors.As(err, &valErr) {
            // Now have access to valErr.Field, valErr.Message
            log.Printf("Invalid field: %s - %s", valErr.Field, valErr.Message)
            return fmt.Errorf("invalid input on field %s", valErr.Field)
        }
        return err
    }
    return nil
}
```

### Custom Unwrap for Complex Chains

```go
type MultiError struct {
    Errors []error
}

func (e *MultiError) Error() string {
    messages := make([]string, len(e.Errors))
    for i, err := range e.Errors {
        messages[i] = err.Error()
    }
    return strings.Join(messages, "; ")
}

func (e *MultiError) Unwrap() error {
    if len(e.Errors) == 0 {
        return nil
    }
    return e.Errors[0] // Return first error for errors.Is/As
}

// Check if contains specific error
func (e *MultiError) Is(target error) bool {
    for _, err := range e.Errors {
        if errors.Is(err, target) {
            return true
        }
    }
    return false
}
```

## Logging Errors

### Structured Logging Pattern

```go
import "log/slog"

func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    order, err := s.repo.GetOrder(ctx, orderID)
    if err != nil {
        // Log with structured attributes
        slog.Error("failed to get order",
            "order_id", orderID,
            "error", err,
        )
        return fmt.Errorf("get order %s: %w", orderID, err)
    }

    if err := s.validate(order); err != nil {
        // Log at appropriate level
        slog.Warn("order validation failed",
            "order_id", orderID,
            "error", err,
        )
        return fmt.Errorf("validate order: %w", err)
    }

    return nil
}
```

### Log Levels Guide

```go
// Debug: Detailed info for debugging (usually disabled in prod)
slog.Debug("processing item", "item_id", id, "status", status)

// Info: General operational events
slog.Info("order processed", "order_id", orderID, "duration", dur)

// Warn: Unexpected but handled situations
slog.Warn("retry attempt", "attempt", attempt, "error", err)

// Error: Failures that affect the operation
slog.Error("failed to connect", "endpoint", endpoint, "error", err)
```

## HTTP Error Responses

### API Error Response Pattern

```go
package api

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error: message,
        Code:  code,
    })
}

// Handler with error handling
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    user, err := h.service.GetUser(r.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            WriteError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
        case errors.Is(err, ErrUnauthorized):
            WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
        default:
            slog.Error("get user failed", "error", err, "user_id", id)
            WriteError(w, http.StatusInternalServerError, "INTERNAL", "Internal server error")
        }
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

## Retry with Error Handling

```go
func WithRetry(ctx context.Context, maxAttempts int, fn func() error) error {
    var lastErr error
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }

        lastErr = err

        // Check if error is retryable
        if !isRetryable(err) {
            return fmt.Errorf("non-retryable error: %w", err)
        }

        // Check context
        if ctx.Err() != nil {
            return ctx.Err()
        }

        // Log retry
        slog.Warn("retry attempt",
            "attempt", attempt,
            "max_attempts", maxAttempts,
            "error", err,
        )

        // Exponential backoff
        time.Sleep(time.Duration(attempt) * time.Second)
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryable(err error) bool {
    // Network errors are usually retryable
    var netErr net.Error
    if errors.As(err, &netErr) {
        return netErr.Timeout() || netErr.Temporary()
    }

    // Check for specific error types
    if errors.Is(err, ErrConnectionRefused) {
        return true
    }

    return false
}
```

## Deferred Error Handling

### Close with Error Check

```go
func ProcessFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open file: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil {
            // If no other error, return close error
            if err == nil {
                err = fmt.Errorf("close file: %w", cerr)
            } else {
                // Log the close error but don't overwrite original
                slog.Warn("file close failed", "error", cerr)
            }
        }
    }()

    // Process file...
    return err
}
```

## Best Practices Summary

1. **Always check errors** - Never ignore returned errors
2. **Return early** - Check and return immediately
3. **Add context** - Use `fmt.Errorf` with `%w` to wrap
4. **Use typed errors** - For domain-specific error handling
5. **Log at appropriate levels** - Don't log and return the same error
6. **Don't panic** - Use error returns for expected failures
7. **Make errors actionable** - Include enough context to debug
8. **Keep error messages useful** - No "an error occurred"
