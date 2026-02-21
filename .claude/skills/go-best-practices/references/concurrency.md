# Go Concurrency Patterns

This reference provides comprehensive patterns for concurrent programming in Go.

## Core Concepts

### Goroutines

Lightweight threads managed by the Go runtime.

```go
// Start a goroutine
go func() {
    fmt.Println("Running in background")
}()

// With parameters (capture by value)
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)  // Pass i as argument
}
```

### Channels

Typed conduits for communication between goroutines.

```go
// Unbuffered channel - blocks until both sides ready
ch := make(chan int)

// Buffered channel - blocks when buffer full
ch := make(chan int, 10)

// Sending (blocks if channel full/unbuffered and no receiver)
ch <- value

// Receiving (blocks if channel empty)
value := <-ch

// Non-blocking receive with ok check
value, ok := <-ch
if !ok {
    // Channel closed
}
```

## Basic Patterns

### Worker Pool

```go
func WorkerPool[T any, R any](ctx context.Context, items []T, workers int, fn func(T) R) []R {
    jobs := make(chan T, len(items))
    results := make(chan R, len(items))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                select {
                case <-ctx.Done():
                    return
                default:
                    results <- fn(item)
                }
            }
        }()
    }

    // Send jobs
    go func() {
        for _, item := range items {
            jobs <- item
        }
        close(jobs)
    }()

    // Wait and close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    var output []R
    for result := range results {
        output = append(output, result)
    }
    return output
}
```

### Fan-Out, Fan-In

```go
// Fan-out: Multiple goroutines read from same channel
func fanOut(source <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = processWorker(source)
    }
    return channels
}

func processWorker(source <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range source {
            out <- v * 2 // Process
        }
    }()
    return out
}

// Fan-in: Merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Pipeline

```go
// Stage 1: Generate
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Stage 2: Transform
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Stage 3: Filter
func filter(in <-chan int, predicate func(int) bool) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if predicate(n) {
                out <- n
            }
        }
    }()
    return out
}

// Usage
nums := generate(1, 2, 3, 4, 5)
squares := square(nums)
evens := filter(squares, func(n int) bool { return n%2 == 0 })
for n := range evens {
    fmt.Println(n)
}
```

## Context Patterns

### Timeout Pattern

```go
func FetchWithTimeout(ctx context.Context, url string) (*Response, error) {
    // Create child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel() // Always call cancel

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("request timed out: %w", err)
        }
        return nil, err
    }

    return resp, nil
}
```

### Cancellation Pattern

```go
func ProcessStream(ctx context.Context, stream <-chan Event) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case event, ok := <-stream:
            if !ok {
                return nil // Stream closed
            }
            if err := process(event); err != nil {
                return fmt.Errorf("process event: %w", err)
            }
        }
    }
}
```

### Propagating Context

```go
type Service struct {
    db *sql.DB
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // Context propagates cancellation and timeouts
    row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

    var user User
    if err := row.Scan(&user.ID, &user.Name); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("query user: %w", err)
    }
    return &user, nil
}
```

## Synchronization Patterns

### Mutex for Shared State

```go
type SafeCounter struct {
    mu    sync.RWMutex
    count map[string]int
}

func (c *SafeCounter) Increment(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count[key]++
}

func (c *SafeCounter) Get(key string) int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count[key]
}
```

### WaitGroup for Coordination

```go
func ProcessAll(items []string) []Result {
    var wg sync.WaitGroup
    results := make([]Result, len(items))

    for i, item := range items {
        wg.Add(1)
        go func(idx int, s string) {
            defer wg.Done()
            results[idx] = process(s)
        }(i, item)
    }

    wg.Wait()
    return results
}
```

### Once for Initialization

```go
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "initialized"}
    })
    return instance
}
```

### Pool for Resource Reuse

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    buf.Write(data)
    return buf.String()
}
```

## Advanced Patterns

### Select with Default (Non-blocking)

```go
func TrySend(ch chan<- int, value int) bool {
    select {
    case ch <- value:
        return true
    default:
        return false // Channel full
    }
}

func TryReceive(ch <-chan int) (int, bool) {
    select {
    case v := <-ch:
        return v, true
    default:
        return 0, false // Channel empty
    }
}
```

### Timeout in Select

```go
func WaitForResult(ch <-chan Result) (Result, error) {
    select {
    case result := <-ch:
        return result, nil
    case <-time.After(5 * time.Second):
        return Result{}, errors.New("timeout waiting for result")
    }
}
```

### Rate Limiting

```go
func RateLimited(throttle <-chan time.Time, items []string, fn func(string)) {
    for _, item := range items {
        <-throttle // Wait for throttle tick
        go fn(item)
    }
}

// Usage: 10 requests per second
throttle := time.Tick(100 * time.Millisecond)
RateLimited(throttle, urls, fetchURL)
```

### Graceful Shutdown

```go
type Server struct {
    shutdown chan struct{}
    wg       sync.WaitGroup
}

func (s *Server) Start() {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        for {
            select {
            case <-s.shutdown:
                return
            case work := <-s.workQueue:
                s.process(work)
            }
        }
    }()
}

func (s *Server) Stop() {
    close(s.shutdown)
    s.wg.Wait() // Wait for all goroutines to finish
}
```

### errgroup for Error Handling

```go
import "golang.org/x/sync/errgroup"

func ProcessFiles(ctx context.Context, files []string) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, file := range files {
        file := file // Capture loop variable
        g.Go(func() error {
            return processFile(ctx, file)
        })
    }

    return g.Wait() // Returns first error
}
```

## Common Pitfalls

### Loop Variable Capture

```go
// WRONG: All goroutines see the same i
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // Will print 5,5,5,5,5 or similar
    }()
}

// CORRECT: Pass as argument
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}

// CORRECT: Create new variable
for i := 0; i < 5; i++ {
    i := i // New variable for each iteration
    go func() {
        fmt.Println(i)
    }()
}
```

### Nil Channel in Select

```go
// Nil channels block forever, useful for disabling cases
func Process(ctx context.Context, input <-chan int, enabled bool) {
    var ch <-chan int
    if enabled {
        ch = input
    } // else ch is nil, that case is disabled

    select {
    case <-ctx.Done():
        return
    case v := <-ch: // Only active if enabled
        process(v)
    }
}
```

### Goroutine Leak Prevention

```go
// Always ensure goroutines can terminate
func Process(ctx context.Context, ch <-chan int) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return // Exit when context cancelled
            case v := <-ch:
                handle(v)
            }
        }
    }()
}
```

## Best Practices Summary

1. **Don't communicate by sharing memory; share memory by communicating**
2. **Always handle context cancellation**
3. **Close channels from the sender side**
4. **Use defer to ensure cleanup**
5. **Prefer channels over mutexes for coordination**
6. **Avoid goroutine leaks - ensure termination**
7. **Use buffered channels when you know the capacity**
8. **Start goroutines only when you have a clear need**
