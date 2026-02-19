// Package events provides NATS messaging for MediSync event-driven architecture.
//
// This package handles publishing and subscribing to events for the ETL pipeline,
// including sync completion, data quality alerts, and general system notifications.
//
// Usage:
//
//	cfg := config.MustLoad()
//	publisher, err := events.NewPublisher(cfg.NATS, logger)
//	if err != nil {
//	    log.Fatal("Failed to create NATS publisher:", err)
//	}
//	defer publisher.Close()
//
//	err = publisher.PublishSyncCompleted(ctx, &events.SyncCompletedEvent{...})
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Event subjects
const (
	// SubjectSyncCompleted is published after a successful sync.
	SubjectSyncCompleted = "etl.sync.completed"
	// SubjectSyncFailed is published on sync failure.
	SubjectSyncFailed = "etl.sync.failed"
	// SubjectDataQualityAlert is published by C-06 agent for quality issues.
	SubjectDataQualityAlert = "etl.data.quality.alert"
	// SubjectAlert is for general system alerts.
	SubjectAlert = "etl.alert"
)

// Publisher provides NATS publishing functionality.
type Publisher struct {
	conn   *nats.Conn
	logger *slog.Logger
	mu     sync.Mutex
}

// PublisherConfig holds configuration for creating a Publisher.
type PublisherConfig struct {
	// URL is the NATS server URL.
	URL string

	// Name is the client connection name.
	Name string

	// MaxReconnects is the maximum reconnection attempts.
	MaxReconnects int

	// ReconnectWait is the wait duration between reconnection attempts.
	ReconnectWait time.Duration

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewPublisher creates a new NATS event publisher.
func NewPublisher(cfg interface{}, logger *slog.Logger) (*Publisher, error) {
	// Parse config
	var url, name string
	var maxReconnects int
	var reconnectWait time.Duration

	switch c := cfg.(type) {
	case map[string]interface{}:
		if u, ok := c["url"].(string); ok {
			url = u
		} else if host, ok := c["host"].(string); ok {
			port := int(c["port"].(float64))
			url = fmt.Sprintf("nats://%s:%d", host, port)
		}
		if n, ok := c["name"].(string); ok {
			name = n
		}
		if mr, ok := c["max_reconnects"].(float64); ok {
			maxReconnects = int(mr)
		}
		if rw, ok := c["reconnect_wait"].(float64); ok {
			reconnectWait = time.Duration(rw) * time.Second
		}
	}

	if url == "" {
		url = nats.DefaultURL
	}

	if name == "" {
		name = "medisync-publisher"
	}

	if maxReconnects == 0 {
		maxReconnects = 10
	}

	if reconnectWait == 0 {
		reconnectWait = 2 * time.Second
	}

	if logger == nil {
		logger = slog.Default()
	}

	// Connect to NATS
	nc, err := nats.Connect(url,
		nats.Name(name),
		nats.MaxReconnects(maxReconnects),
		nats.ReconnectWait(reconnectWait),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				logger.Warn("NATS disconnected",
					slog.String("error", err.Error()),
				)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected",
				slog.String("url", nc.ConnectedUrl()),
			)
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Info("NATS connection closed")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("events: failed to connect to NATS: %w", err)
	}

	publisher := &Publisher{
		conn:   nc,
		logger: logger,
	}

	logger.Info("connected to NATS",
		slog.String("url", url),
	)

	return publisher, nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
	return nil
}

// Publish publishes a message to a NATS subject.
func (p *Publisher) Publish(ctx context.Context, subject string, data interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil {
		return fmt.Errorf("events: publisher is closed")
	}

	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("events: failed to marshal event data: %w", err)
	}

	// Publish message
	err = p.conn.Publish(subject, jsonData)
	if err != nil {
		return fmt.Errorf("events: failed to publish to %s: %w", subject, err)
	}

	p.logger.Debug("published event",
		slog.String("subject", subject),
		slog.Int("size", len(jsonData)),
	)

	return nil
}

// PublishAsync publishes a message asynchronously (simplified - uses sync publish).
func (p *Publisher) PublishAsync(ctx context.Context, subject string, data interface{}) error {
	return p.Publish(ctx, subject, data)
}

// ============================================================================
// Event Types
// ============================================================================

// SyncCompletedEvent is published after a successful sync.
type SyncCompletedEvent struct {
	EventID             string                 `json:"event_id"`
	Source              string                 `json:"source"`
	Entity              string                 `json:"entity"`
	StartedAt           time.Time              `json:"started_at"`
	CompletedAt         time.Time              `json:"completed_at"`
	RecordsProcessed    int                    `json:"records_processed"`
	RecordsInserted     int                    `json:"records_inserted"`
	RecordsUpdated      int                    `json:"records_updated"`
	RecordsQuarantined  int                    `json:"records_quarantined"`
	DurationSec         float64                `json:"duration_sec"`
	Status              string                 `json:"status"`
	BatchID             *string                `json:"batch_id,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// SyncFailedEvent is published on sync failure.
type SyncFailedEvent struct {
	EventID           string                 `json:"event_id"`
	Source            string                 `json:"source"`
	Entity            string                 `json:"entity"`
	StartedAt         time.Time              `json:"started_at"`
	FailedAt          time.Time              `json:"failed_at"`
	Error             string                 `json:"error"`
	ErrorCode         string                 `json:"error_code"`
	RecordsProcessed  int                    `json:"records_processed"`
	DurationSec       float64                `json:"duration_sec"`
	Retryable         bool                   `json:"retryable"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// DataQualityAlertEvent is published for data quality issues.
type DataQualityAlertEvent struct {
	EventID      string                 `json:"event_id"`
	BatchID      string                 `json:"batch_id"`
	Source       string                 `json:"source"`
	Entity       *string                `json:"entity,omitempty"`
	AlertLevel   string                 `json:"alert_level"` // info, warning, error, critical
	AlertType    string                 `json:"alert_type"`
	Message      string                 `json:"message"`
	Details      map[string]interface{} `json:"details,omitempty"`
	QualityScore *float64              `json:"quality_score,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// AlertEvent is for general system alerts.
type AlertEvent struct {
	EventID    string                 `json:"event_id"`
	Level      string                 `json:"level"` // info, warning, error, critical
	Type       string                 `json:"type"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ============================================================================
// Convenience Methods
// ============================================================================

// PublishSyncCompleted publishes a sync completion event.
func (p *Publisher) PublishSyncCompleted(ctx context.Context, event *SyncCompletedEvent) error {
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	return p.Publish(ctx, SubjectSyncCompleted, event)
}

// PublishSyncFailed publishes a sync failure event.
func (p *Publisher) PublishSyncFailed(ctx context.Context, event *SyncFailedEvent) error {
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	return p.Publish(ctx, SubjectSyncFailed, event)
}

// PublishDataQualityAlert publishes a data quality alert.
func (p *Publisher) PublishDataQualityAlert(ctx context.Context, event *DataQualityAlertEvent) error {
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	event.CreatedAt = time.Now()
	return p.Publish(ctx, SubjectDataQualityAlert, event)
}

// PublishAlert publishes a general alert.
func (p *Publisher) PublishAlert(ctx context.Context, event *AlertEvent) error {
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	event.CreatedAt = time.Now()
	return p.Publish(ctx, SubjectAlert, event)
}

// ============================================================================
// Subscriber
// ============================================================================

// SubscriptionOptions configures a subscription.
type SubscriptionOptions struct {
	// Queue is the queue group name.
	Queue string

	// Durable is the durable subscription name.
	Durable string

	// AutoAck automatically acknowledges messages.
	AutoAck bool
}

// MessageHandler is a function that handles incoming messages.
type MessageHandler func(msg *Message) error

// Message represents a received NATS message.
type Message struct {
	Subject string
	Data    []byte
	Reply   string
}

// Subscriber provides NATS subscription functionality.
type Subscriber struct {
	pub     *Publisher
	sub     *nats.Subscription
	handler MessageHandler
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *slog.Logger
}

// NewSubscriber creates a new NATS subscriber.
func NewSubscriber(publisher *Publisher, subject string, handler MessageHandler, opts *SubscriptionOptions) (*Subscriber, error) {
	if publisher == nil {
		return nil, fmt.Errorf("events: publisher is required")
	}

	if publisher.conn == nil {
		return nil, fmt.Errorf("events: publisher connection is closed")
	}

	if opts == nil {
		opts = &SubscriptionOptions{}
	}

	sub := &Subscriber{
		pub:     publisher,
		handler: handler,
		logger:  publisher.logger,
	}

	sub.ctx, sub.cancel = context.WithCancel(context.Background())

	// Create subscription
	natsSub, err := publisher.conn.Subscribe(subject, func(msg *nats.Msg) {
		sub.handleMessage(msg)
	})
	if err != nil {
		return nil, fmt.Errorf("events: failed to create subscription: %w", err)
	}
	sub.sub = natsSub

	publisher.logger.Info("created NATS subscription",
		slog.String("subject", subject),
		slog.String("queue", opts.Queue),
	)

	return sub, nil
}

// handleMessage processes a standard NATS message.
func (s *Subscriber) handleMessage(msg *nats.Msg) {
	wrapped := &Message{
		Subject: msg.Subject,
		Data:    msg.Data,
		Reply:   msg.Reply,
	}

	if err := s.handler(wrapped); err != nil {
		s.logger.Error("message handler error",
			slog.String("subject", msg.Subject),
			slog.String("error", err.Error()),
		)
		// Negative acknowledge for retry
		msg.Nak()
	} else {
		msg.Ack()
	}
}

// Close closes the subscriber.
func (s *Subscriber) Close() error {
	s.cancel()

	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// Utility Functions
// ============================================================================

// generateEventID generates a unique event ID.
func generateEventID() string {
	return fmt.Sprintf("evt-%d", time.Now().UnixNano())
}

// PublishSyncCompletedBatch publishes multiple sync completion events efficiently.
func (p *Publisher) PublishSyncCompletedBatch(ctx context.Context, events []*SyncCompletedEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if event.EventID == "" {
			event.EventID = generateEventID()
		}
		if err := p.PublishAsync(ctx, SubjectSyncCompleted, event); err != nil {
			return err
		}
	}

	// Flush to ensure all messages are sent
	if p.conn != nil {
		return p.conn.Flush()
	}

	return nil
}

// IsValidSubject checks if a subject is valid for publishing.
func IsValidSubject(subject string) bool {
	// Basic NATS subject validation
	if subject == "" {
		return false
	}

	// Check for invalid characters
	for _, c := range subject {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			return false
		}
	}

	return true
}
