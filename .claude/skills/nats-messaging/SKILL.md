---
name: nats-messaging
description: This skill should be used when the user asks to "implement NATS messaging", "publish NATS events", "NATS JetStream", "event-driven architecture", "message queues", "NATS subscriptions", "event publishing", or mentions NATS-specific concepts like subjects, consumers, streams, or ack.
---

# NATS Messaging Patterns for MediSync

NATS JetStream provides event-driven messaging for MediSync's ETL pipeline, notifications, and inter-service communication. This skill covers publishing, subscribing, and stream management.

★ Insight ─────────────────────────────────────
MediSync uses NATS for:
1. **ETL Events** - Data sync completion notifications
2. **Agent Coordination** - Inter-agent communication
3. **Alerts** - Real-time notification delivery
4. **Audit Logs** - Event sourcing for audit trail

JetStream provides durability and replay capabilities.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Convention |
|--------|------------|
| **Subjects** | Dot-notation: `domain.entity.action` |
| **Encoding** | JSON for structured data |
| **Ack** | Always acknowledge processed messages |
| **Durable** | Use durable consumers for reliability |
| **Queues** | Use queue groups for load balancing |

## Subject Naming Convention

```
medisync.{domain}.{entity}.{action}

Examples:
- medisync.etl.tally.completed
- medisync.etl.hims.synced
- medisync.document.ocr.extracted
- medisync.tally.sync.approved
- medisync.alert.notification.created
```

## Publishing Events

### Basic Publishing (Go)

```go
package events

import (
    "encoding/json"
    "fmt"

    "github.com/nats-io/nats.go"
)

type Publisher struct {
    nc *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
    nc, err := nats.Connect(url,
        nats.Name("medisync-publisher"),
        nats.ReconnectWait(2*time.Second),
        nats.MaxReconnects(10),
    )
    if err != nil {
        return nil, fmt.Errorf("connect to NATS: %w", err)
    }
    return &Publisher{nc: nc}, nil
}

// Event represents a domain event
type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Source    string                 `json:"source"`
    Payload   map[string]interface{} `json:"payload"`
}

// Publish publishes an event to a subject
func (p *Publisher) Publish(ctx context.Context, subject string, event Event) error {
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }

    if err := p.nc.Publish(subject, data); err != nil {
        return fmt.Errorf("publish event: %w", err)
    }

    return nil
}

// Common event types
func (p *Publisher) PublishETLCompleted(ctx context.Context, source string, recordsProcessed int) error {
    return p.Publish(ctx, "medisync.etl."+source+".completed", Event{
        ID:        uuid.New().String(),
        Type:      "etl.completed",
        Timestamp: time.Now(),
        Source:    source,
        Payload: map[string]interface{}{
            "records_processed": recordsProcessed,
        },
    })
}
```

### JetStream Publishing

```go
func (p *Publisher) PublishToStream(ctx context.Context, stream, subject string, event Event) (*nats.PubAck, error) {
    js, err := p.nc.JetStream()
    if err != nil {
        return nil, fmt.Errorf("get jetstream: %w", err)
    }

    data, err := json.Marshal(event)
    if err != nil {
        return nil, fmt.Errorf("marshal event: %w", err)
    }

    ack, err := js.Publish(subject, data)
    if err != nil {
        return nil, fmt.Errorf("publish to stream: %w", err)
    }

    return ack, nil
}
```

## Subscribing to Events

### Queue Subscription (Load Balancing)

```go
type Subscriber struct {
    nc *nats.Conn
    js nats.JetStreamContext
}

func NewSubscriber(url string) (*Subscriber, error) {
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }

    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    return &Subscriber{nc: nc, js: js}, nil
}

// SubscribeQueue joins a queue group for load balancing
func (s *Subscriber) SubscribeQueue(subject, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
    sub, err := s.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
        // Handle message
        handler(msg)
        // Always acknowledge
        msg.Ack()
    })
    if err != nil {
        return nil, fmt.Errorf("subscribe: %w", err)
    }
    return sub, nil
}
```

### Durable Consumer (JetStream)

```go
func (s *Subscriber) CreateDurableConsumer(stream, subject, durable string, handler func(Event) error) (*nats.Subscription, error) {
    sub, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
        var event Event
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            // Invalid message, nack and move on
            msg.Nak()
            return
        }

        if err := handler(event); err != nil {
            // Processing failed, requeue after delay
            msg.NakWithDelay(30 * time.Second)
            return
        }

        // Success
        msg.Ack()
    },
        nats.Durable(durable),
        nats.ManualAck(),
        nats.DeliverAll(),
        nats.MaxDeliver(5),
    )
    if err != nil {
        return nil, fmt.Errorf("create consumer: %w", err)
    }
    return sub, nil
}
```

## Stream Management

### Creating Streams

```go
func (s *Subscriber) EnsureStream(streamName string, subjects []string) error {
    // Check if stream exists
    info, err := s.js.StreamInfo(streamName)
    if err == nil {
        return nil // Stream exists
    }

    // Create stream
    _, err = s.js.AddStream(&nats.StreamConfig{
        Name:     streamName,
        Subjects: subjects,
        Retention: nats.LimitsPolicy,
        MaxAge:   7 * 24 * time.Hour, // 7 days retention
        Replicas: 1,
        Storage:  nats.FileStorage,
    })
    if err != nil {
        return fmt.Errorf("create stream: %w", err)
    }

    return nil
}

// Common streams for MediSync
func (s *Subscriber) SetupMediSyncStreams() error {
    streams := []struct {
        name     string
        subjects []string
    }{
        {
            name:     "EVENTS",
            subjects: []string{"medisync.>"},
        },
        {
            name:     "ALERTS",
            subjects: []string{"medisync.alert.>"},
        },
        {
            name:     "AUDIT",
            subjects: []string{"medisync.audit.>"},
        },
    }

    for _, s := range streams {
        if err := s.EnsureStream(s.name, s.subjects); err != nil {
            return err
        }
    }
    return nil
}
```

## Event Patterns

### ETL Completion Event

```go
type ETLEvent struct {
    Source          string    `json:"source"`
    RecordsRead     int       `json:"records_read"`
    RecordsWritten  int       `json:"records_written"`
    RecordsSkipped  int       `json:"records_skipped"`
    Duration        int64     `json:"duration_ms"`
    CompletedAt     time.Time `json:"completed_at"`
    Error           string    `json:"error,omitempty"`
}

func (p *Publisher) PublishETLEvent(ctx context.Context, event ETLEvent) error {
    subject := fmt.Sprintf("medisync.etl.%s.completed", event.Source)
    return p.Publish(ctx, subject, Event{
        ID:        uuid.New().String(),
        Type:      "etl.completed",
        Timestamp: time.Now(),
        Source:    "etl-service",
        Payload:   event,
    })
}
```

### Document Processing Event

```go
type DocumentEvent struct {
    DocumentID   string  `json:"document_id"`
    DocumentType string  `json:"document_type"`
    Status       string  `json:"status"`
    Confidence   float64 `json:"confidence"`
    ExtractedData any     `json:"extracted_data,omitempty"`
    Error        string  `json:"error,omitempty"`
}

func (p *Publisher) PublishDocumentProcessed(ctx context.Context, event DocumentEvent) error {
    subject := fmt.Sprintf("medisync.document.%s.%s", event.DocumentType, event.Status)
    return p.Publish(ctx, subject, Event{
        ID:        uuid.New().String(),
        Type:      "document.processed",
        Timestamp: time.Now(),
        Source:    "ocr-service",
        Payload:   event,
    })
}
```

### Tally Sync Event

```go
type TallySyncEvent struct {
    SyncID       string   `json:"sync_id"`
    EntryIDs     []string `json:"entry_ids"`
    Company      string   `json:"company"`
    Status       string   `json:"status"`
    SyncedAt     time.Time `json:"synced_at"`
    VoucherNos   []string `json:"voucher_nos,omitempty"`
    Error        string   `json:"error,omitempty"`
}

func (p *Publisher) PublishTallySync(ctx context.Context, event TallySyncEvent) error {
    subject := fmt.Sprintf("medisync.tally.sync.%s", event.Status)
    return p.Publish(ctx, subject, Event{
        ID:        uuid.New().String(),
        Type:      "tally.sync",
        Timestamp: time.Now(),
        Source:    "tally-service",
        Payload:   event,
    })
}
```

## Error Handling

```go
func (s *Subscriber) SubscribeWithRetry(subject string, maxRetries int, handler func(Event) error) error {
    _, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
        var event Event
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            log.Error("unmarshal event", "error", err)
            msg.Term() // Don't redeliver invalid messages
            return
        }

        // Get metadata for retry count
        meta, _ := msg.Metadata()

        if err := handler(event); err != nil {
            if meta.NumDelivered >= int64(maxRetries) {
                log.Error("max retries exceeded", "event_id", event.ID)
                msg.Term() // Give up
                return
            }

            // Retry with exponential backoff
            delay := time.Duration(math.Pow(2, float64(meta.NumDelivered))) * time.Second
            msg.NakWithDelay(delay)
            return
        }

        msg.Ack()
    },
        nats.ManualAck(),
        nats.MaxDeliver(maxRetries),
        nats.AckWait(30*time.Second),
    )

    return err
}
```

## Configuration

### NATS Config (nats.conf)

```
# Server configuration
server_name: medisync-nats
listen: 0.0.0.0:4222

# JetStream
jetstream {
    store_dir: /data/jetstream
    max_mem: 1G
    max_file: 10G
}

# Authentication
authorization {
    users: [
        {
            user: medisync,
            password: $2a$11$...  # bcrypt hash
        }
    ]
}

# Clustering (for production)
cluster {
    name: medisync-cluster
    listen: 0.0.0.0:6222
    routes: [
        nats://nats-1:6222
        nats://nats-2:6222
    ]
}
```

## Additional Resources

### Reference Files
- **`references/jetstream.md`** - Advanced JetStream patterns
- **`references/patterns.md`** - Event-driven architecture patterns

### Example Files
- **`examples/publisher.go`** - Complete publisher implementation
- **`examples/subscriber.go`** - Complete subscriber implementation
