// Package events provides NATS messaging infrastructure for MediSync
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// =============================================================================
// TYPES
// =============================================================================

// Event represents a domain event
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Payload   map[string]interface{} `json:"payload"`
}

// Config holds NATS connection configuration
type Config struct {
	URL          string
	ClientName   string
	MaxReconnect int
	ReconnectWait time.Duration
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		URL:          "nats://localhost:4222",
		ClientName:   "medisync",
		MaxReconnect: 10,
		ReconnectWait: 2 * time.Second,
	}
}

// Publisher publishes events to NATS
type Publisher struct {
	nc  *nats.Conn
	js  nats.JetStreamContext
	cfg Config
	mu  sync.RWMutex
}

// =============================================================================
// CONSTRUCTOR
// =============================================================================

// NewPublisher creates a new NATS publisher
func NewPublisher(cfg Config) (*Publisher, error) {
	nc, err := nats.Connect(cfg.URL,
		nats.Name(cfg.ClientName+"-publisher"),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnect),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Warn("NATS disconnected", "error", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("NATS reconnected", "url", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			slog.Error("NATS connection closed", "error", nc.LastError())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("get jetstream context: %w", err)
	}

	return &Publisher{
		nc:  nc,
		js:  js,
		cfg: cfg,
	}, nil
}

// =============================================================================
// PUBLISHING METHODS
// =============================================================================

// Publish publishes an event to a subject
func (p *Publisher) Publish(ctx context.Context, subject string, event Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := p.nc.Publish(subject, data); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	slog.Debug("published event",
		"subject", subject,
		"event_id", event.ID,
		"event_type", event.Type,
	)

	return nil
}

// PublishWithAck publishes an event and waits for acknowledgement
func (p *Publisher) PublishWithAck(ctx context.Context, subject string, event Event) (*nats.PubAck, error) {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshal event: %w", err)
	}

	ack, err := p.js.Publish(subject, data,
		nats.Context(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("publish with ack: %w", err)
	}

	slog.Debug("published event with ack",
		"subject", subject,
		"event_id", event.ID,
		"stream", ack.Stream),
		"sequence", ack.Sequence,
	)

	return ack, nil
}

// Close closes the NATS connection
func (p *Publisher) Close() {
	p.nc.Close()
}

// =============================================================================
// DOMAIN EVENT PUBLISHERS
// =============================================================================

// PublishETLCompleted publishes an ETL completion event
func (p *Publisher) PublishETLCompleted(ctx context.Context, source string, stats ETLStats) error {
	subject := fmt.Sprintf("medisync.etl.%s.completed", source)
	return p.Publish(ctx, subject, Event{
		ID:        uuid.New().String(),
		Type:      "etl.completed",
		Timestamp: time.Now(),
		Source:    source,
		Payload: map[string]interface{}{
			"records_read":    stats.RecordsRead,
			"records_written": stats.RecordsWritten,
			"records_skipped": stats.RecordsSkipped,
			"duration_ms":     stats.Duration.Milliseconds(),
			"errors":          stats.Errors,
		},
	})
}

// ETLStats holds ETL execution statistics
type ETLStats struct {
	RecordsRead    int
	RecordsWritten int
	RecordsSkipped int
	Duration       time.Duration
	Errors         []string
}

// PublishDocumentProcessed publishes a document processing event
func (p *Publisher) PublishDocumentProcessed(ctx context.Context, doc DocumentEvent) error {
	subject := fmt.Sprintf("medisync.document.%s.%s", doc.DocumentType, doc.Status)
	return p.Publish(ctx, subject, Event{
		ID:        uuid.New().String(),
		Type:      "document.processed",
		Timestamp: time.Now(),
		Source:    "ocr-service",
		Payload: map[string]interface{}{
			"document_id":   doc.DocumentID,
			"document_type": doc.DocumentType,
			"status":        doc.Status,
			"confidence":    doc.Confidence,
			"extracted":     doc.ExtractedData,
			"error":         doc.Error,
		},
	})
}

// DocumentEvent represents a document processing event
type DocumentEvent struct {
	DocumentID   string
	DocumentType string
	Status       string
	Confidence   float64
	ExtractedData map[string]interface{}
	Error        string
}

// PublishTallySync publishes a Tally sync event
func (p *Publisher) PublishTallySync(ctx context.Context, sync TallySyncEvent) error {
	subject := fmt.Sprintf("medisync.tally.sync.%s", sync.Status)
	return p.PublishWithAck(ctx, subject, Event{
		ID:        uuid.New().String(),
		Type:      "tally.sync",
		Timestamp: time.Now(),
		Source:    "tally-service",
		Payload: map[string]interface{}{
			"sync_id":     sync.SyncID,
			"entry_ids":   sync.EntryIDs,
			"company":     sync.Company,
			"status":      sync.Status,
			"voucher_nos": sync.VoucherNos,
			"error":       sync.Error,
		},
	})
}

// TallySyncEvent represents a Tally sync event
type TallySyncEvent struct {
	SyncID     string
	EntryIDs   []string
	Company    string
	Status     string
	VoucherNos []string
	Error      string
}

// PublishAlert publishes an alert event
func (p *Publisher) PublishAlert(ctx context.Context, alert AlertEvent) error {
	subject := fmt.Sprintf("medisync.alert.%s.%s", alert.Severity, alert.Type)
	return p.Publish(ctx, subject, Event{
		ID:        uuid.New().String(),
		Type:      "alert.created",
		Timestamp: time.Now(),
		Source:    "alert-service",
		Payload: map[string]interface{}{
			"alert_id":  alert.ID,
			"type":      alert.Type,
			"severity":  alert.Severity,
			"title":     alert.Title,
			"message":   alert.Message,
			"company_id": alert.CompanyID,
			"metadata":  alert.Metadata,
		},
	})
}

// AlertEvent represents an alert event
type AlertEvent struct {
	ID        string
	Type      string
	Severity  string // "info", "warning", "error", "critical"
	Title     string
	Message   string
	CompanyID string
	Metadata  map[string]interface{}
}

// PublishAudit publishes an audit log event
func (p *Publisher) PublishAudit(ctx context.Context, audit AuditEvent) error {
	subject := fmt.Sprintf("medisync.audit.%s.%s", audit.EntityType, audit.Action)
	_, err := p.PublishWithAck(ctx, subject, Event{
		ID:        uuid.New().String(),
		Type:      "audit.entry",
		Timestamp: time.Now(),
		Source:    "audit-service",
		Payload: map[string]interface{}{
			"user_id":      audit.UserID,
			"company_id":   audit.CompanyID,
			"action":       audit.Action,
			"entity_type":  audit.EntityType,
			"entity_id":    audit.EntityID,
			"old_values":   audit.OldValues,
			"new_values":   audit.NewValues,
			"ip_address":   audit.IPAddress,
		},
	})
	return err
}

// AuditEvent represents an audit log event
type AuditEvent struct {
	UserID     string
	CompanyID  string
	Action     string
	EntityType string
	EntityID   string
	OldValues  map[string]interface{}
	NewValues  map[string]interface{}
	IPAddress  string
}
