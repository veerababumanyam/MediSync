---
name: go-best-practices
description: This skill should be used when the user asks to "write Go code", "create a Go package", "implement Go patterns", "Go idiomatic code", "Golang best practices", "Go error handling", "Go concurrency patterns", "Go project structure", or mentions Go-specific concepts like goroutines, channels, interfaces, or struct embedding.
---

# Go (Golang) Best Practices

Go is a statically typed, compiled programming language designed for simplicity, reliability, and efficiency. This skill provides idiomatic patterns, conventions, and best practices for writing clean, maintainable Go code.

★ Insight ─────────────────────────────────────
Go's philosophy emphasizes:
1. **Simplicity** - Few keywords, clear semantics, no hidden magic
2. **Readability** - Code is read more than written
3. **Pragmatism** - Practical solutions over theoretical purity

When in doubt, choose the simpler approach.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Convention |
|--------|------------|
| **Formatting** | `gofmt` / `go fmt` - always run before commit |
| **Imports** | Standard library, then external, then local (grouped) |
| **Naming** | MixedCaps/mixedCaps (not snake_case), exported = Capitalized |
| **Error Handling** | Return errors as last return value, never panic |
| **Comments** | Start with the name of the thing being described |

## Project Structure

```
my-project/
├── cmd/                    # Main applications
│   ├── api/
│   │   └── main.go
│   └── etl/
│       └── main.go
├── internal/               # Private application code
│   ├── service/
│   └── repository/
├── pkg/                    # Public library code (if any)
│   └── mylib/
├── api/                    # API definitions (OpenAPI, proto, etc.)
├── configs/                # Configuration files
├── scripts/                # Build/deploy scripts
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
└── Makefile                # Build automation
```

**Key directories:**
- `cmd/` - Entry points, one subdirectory per executable
- `internal/` - Compiler-enforced private code
- `pkg/` - Public packages (importable by external projects)

## Naming Conventions

### Packages
```go
// Package names: lowercase, single word, no underscores
package repository   // Good
package data_access  // Bad
```

### Variables and Functions
```go
// Local variables: mixedCaps
var userProfile User    // Good
var user_profile User   // Bad

// Exported: Capitalized
func GetUser(id int) (*User, error) {}  // Exported
func getUser(id int) (*User, error) {}  // Unexported

// Acronyms: consistent case
var urlURL string    // Good
var UrlURL string    // Bad
var HTTPClient       // Good
var HttpClient       // Bad
```

### Interfaces
```go
// Single-method interfaces: -er suffix
type Reader interface { Read([]byte) (int, error) }
type Writer interface { Write([]byte) (int, error) }
type Repository interface { Get(id string) (*Entity, error) }
```

## Error Handling

### Basic Pattern
```go
// Return error as last value, check immediately
func GetUser(id int) (*User, error) {
    user, err := db.QueryUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    return user, nil
}
```

### Wrapping Errors
```go
import "errors"

// Use %w for wrapping (allows errors.Is/As)
if err != nil {
    return fmt.Errorf("failed to process order %d: %w", orderID, err)
}

// Use %v for just adding context (no unwrapping)
if err != nil {
    return fmt.Errorf("failed to process order: %v", err)
}
```

### Custom Errors
```go
// Sentinel errors for specific conditions
var ErrNotFound = errors.New("not found")

// Custom error types for structured errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}
```

### Error Checking
```go
// Check specific errors
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

// Check error type
var valErr *ValidationError
if errors.As(err, &valErr) {
    // Handle validation error
}
```

## Interface Design

### Accept Interfaces, Return Structs
```go
// Good: Accept interface for flexibility
func ProcessData(r io.Reader) error {
    // Works with any io.Reader
}

// Good: Return concrete type for clarity
func NewUser(name string) *User {
    return &User{Name: name}
}
```

### Small Interfaces
```go
// Good: Focused, single responsibility
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Bad: Too many methods
type UserManager interface {
    Create() error
    Update() error
    Delete() error
    SendEmail() error
    GenerateReport() error
}
```

### Interface Satisfaction
```go
// Implicit implementation - no "implements" keyword
type Logger interface {
    Log(message string)
}

// Any type with Log(string) satisfies Logger
type ConsoleLogger struct{}
func (l ConsoleLogger) Log(message string) {
    fmt.Println(message)
}
```

## Concurrency Patterns

### Goroutines and Channels
```go
// Worker pool pattern
func processJobs(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        results <- process(job)
    }
}

// Launch workers
jobs := make(chan Job, 100)
results := make(chan Result, 100)
var wg sync.WaitGroup

for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go processJobs(jobs, results, &wg)
}

// Close jobs channel when done sending
go func() {
    wg.Wait()
    close(results)
}()
```

### Context for Cancellation
```go
func (s *Service) FetchData(ctx context.Context, id string) (*Data, error) {
    // Check context before starting
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Pass context to downstream operations
    resp, err := s.client.Get(ctx, "/data/"+id)
    if err != nil {
        return nil, fmt.Errorf("fetch data: %w", err)
    }
    return resp, nil
}
```

## Testing

### Table-Driven Tests
```go
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name     string
        items    []Item
        expected float64
        wantErr  bool
    }{
        {
            name:     "empty items",
            items:    []Item{},
            expected: 0,
            wantErr:  false,
        },
        {
            name: "multiple items",
            items: []Item{
                {Price: 10.0, Qty: 2},
                {Price: 5.0, Qty: 3},
            },
            expected: 35.0,
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CalculateTotal(tt.items)
            if (err != nil) != tt.wantErr {
                t.Errorf("CalculateTotal() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.expected {
                t.Errorf("CalculateTotal() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Test Helpers
```go
// Use t.Helper() for better error reporting
func testUser(t *testing.T, name string) *User {
    t.Helper() // Marks this as a helper
    return &User{Name: name, ID: uuid.New()}
}
```

## Common Patterns

### Functional Options
```go
type Server struct {
    port    int
    timeout time.Duration
    logger  Logger
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) { s.timeout = timeout }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        port:    8080,           // default
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(WithPort(9000), WithTimeout(60*time.Second))
```

### Zero Values
```go
// Go types have useful zero values
var s string      // "" (empty string)
var i int         // 0
var f float64     // 0.0
var b bool        // false
var ptr *int      // nil
var slice []int   // nil (len=0, cap=0)
var m map[K]V     // nil (len=0)
var ch chan int   // nil

// Design structs to work with zero values
type Config struct {
    Port    int    // 0 means use default
    Timeout int    // 0 means no timeout
    Debug   bool   // false is the safe default
}
```

## Code Organization

### Package Organization
```go
// Group related types and functions
// user/user.go
package user

type User struct { /* ... */ }
type Service struct { /* ... */ }

func NewService(repo Repository) *Service { /* ... */ }
func (s *Service) Get(id string) (*User, error) { /* ... */ }
```

### Import Organization
```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // External packages
    "github.com/google/uuid"
    "go.uber.org/zap"

    // Internal packages
    "myapp/internal/repository"
    "myapp/internal/service"
)
```

## Common Gotchas

| Gotcha | Issue | Solution |
|--------|-------|----------|
| Loop variable capture | All goroutines share same variable | Use local copy or pass as argument |
| Nil map write | Panic on write to nil map | Initialize with `make(map[K]V)` |
| Unused imports | Compile error | Remove or use blank import `_` |
| Shadowing | Inner variable shadows outer | Use `go vet` to detect |

## Additional Resources

### Reference Files
- **`references/error-handling.md`** - Comprehensive error patterns
- **`references/concurrency.md`** - Advanced concurrency patterns
- **`references/testing.md`** - Testing strategies and patterns
- **`references/project-layout.md`** - Detailed project structure

### Example Files
- **`examples/service.go`** - Complete service implementation
- **`examples/repository.go`** - Repository pattern example
- **`examples/handler.go`** - HTTP handler patterns

### Utility Scripts
- **`scripts/lint.sh`** - Run golangci-lint with recommended config

## Key Commands

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter (install first)
golangci-lint run

# Update dependencies
go mod tidy

# Vendor dependencies
go mod vendor
```
