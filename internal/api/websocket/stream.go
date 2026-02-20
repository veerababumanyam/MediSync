// Package websocket provides WebSocket utilities for the MediSync API.
//
// This file implements chunked streaming support for WebSocket connections,
// enabling efficient transfer of large datasets with progress tracking.
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================================
// Stream Message Types
// ============================================================================

// StreamMessageType defines the types of streaming messages.
type StreamMessageType string

const (
	// StreamMessageTypeChunk indicates a data chunk in the stream.
	StreamMessageTypeChunk StreamMessageType = "stream_chunk"

	// StreamMessageTypeComplete indicates the stream has completed successfully.
	StreamMessageTypeComplete StreamMessageType = "stream_complete"

	// StreamMessageTypeError indicates an error occurred during streaming.
	StreamMessageTypeError StreamMessageType = "stream_error"
)

// StreamMessage represents a message in the streaming protocol.
type StreamMessage struct {
	// Type is the message type (stream_chunk, stream_complete, stream_error).
	Type StreamMessageType `json:"type"`

	// StreamID is a unique identifier for the stream session.
	StreamID string `json:"stream_id"`

	// ChunkIndex is the 0-based index of this chunk in the stream.
	ChunkIndex int `json:"chunk_index,omitempty"`

	// TotalChunks is the total number of chunks (if known).
	TotalChunks int `json:"total_chunks,omitempty"`

	// Data contains the chunk payload.
	Data interface{} `json:"data,omitempty"`

	// Progress is the percentage complete (0-100).
	Progress float64 `json:"progress,omitempty"`

	// Error contains error details for stream_error messages.
	Error *StreamError `json:"error,omitempty"`

	// Timestamp is when the message was created.
	Timestamp time.Time `json:"timestamp"`
}

// StreamError represents an error in the streaming protocol.
type StreamError struct {
	// Code is the error code.
	Code string `json:"code"`

	// Message is the human-readable error message.
	Message string `json:"message"`

	// Recoverable indicates if the stream can be resumed.
	Recoverable bool `json:"recoverable"`
}

// StreamChunk represents a single chunk of streaming data.
type StreamChunk struct {
	// Index is the 0-based chunk index.
	Index int `json:"index"`

	// Data is the chunk payload.
	Data interface{} `json:"data"`

	// IsLast indicates if this is the final chunk.
	IsLast bool `json:"is_last"`
}

// ============================================================================
// StreamWriter
// ============================================================================

// StreamWriter provides thread-safe writing of chunked data to a WebSocket connection.
type StreamWriter struct {
	conn        *websocket.Conn
	streamID    string
	mu          sync.Mutex
	logger      *slog.Logger
	writeWait   time.Duration
	chunkCount  int
	totalChunks int
	closed      bool
}

// StreamWriterConfig holds configuration for the StreamWriter.
type StreamWriterConfig struct {
	// Conn is the WebSocket connection to write to.
	Conn *websocket.Conn

	// StreamID is a unique identifier for this stream.
	StreamID string

	// Logger is the structured logger (defaults to slog.Default()).
	Logger *slog.Logger

	// WriteWait is the write timeout (defaults to 10s).
	WriteWait time.Duration

	// TotalChunks is the expected total number of chunks (0 if unknown).
	TotalChunks int
}

// NewStreamWriter creates a new thread-safe StreamWriter.
func NewStreamWriter(cfg StreamWriterConfig) *StreamWriter {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	if cfg.WriteWait == 0 {
		cfg.WriteWait = 10 * time.Second
	}

	return &StreamWriter{
		conn:        cfg.Conn,
		streamID:    cfg.StreamID,
		logger:      cfg.Logger,
		writeWait:   cfg.WriteWait,
		totalChunks: cfg.TotalChunks,
	}
}

// WriteChunk writes a data chunk to the WebSocket connection.
// It is safe to call from multiple goroutines.
func (sw *StreamWriter) WriteChunk(ctx context.Context, chunk *StreamChunk) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.closed {
		return fmt.Errorf("stream writer is closed")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
	}

	// Calculate progress
	var progress float64
	if sw.totalChunks > 0 {
		progress = float64(chunk.Index+1) / float64(sw.totalChunks) * 100
	} else if chunk.IsLast {
		progress = 100
	}

	// Create stream message
	msg := &StreamMessage{
		Type:       StreamMessageTypeChunk,
		StreamID:   sw.streamID,
		ChunkIndex: chunk.Index,
		TotalChunks: sw.totalChunks,
		Data:       chunk.Data,
		Progress:   progress,
		Timestamp:  time.Now(),
	}

	// Set write deadline
	if err := sw.conn.SetWriteDeadline(time.Now().Add(sw.writeWait)); err != nil {
		sw.logger.Error("failed to set write deadline",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
		)
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write message
	if err := sw.conn.WriteJSON(msg); err != nil {
		sw.logger.Error("failed to write chunk",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
			slog.Int("chunk_index", chunk.Index),
		)
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	sw.chunkCount++

	sw.logger.Debug("chunk written",
		slog.String("stream_id", sw.streamID),
		slog.Int("chunk_index", chunk.Index),
		slog.Float64("progress", progress),
	)

	return nil
}

// Complete marks the stream as complete and sends a completion message.
func (sw *StreamWriter) Complete(ctx context.Context) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.closed {
		return fmt.Errorf("stream writer is closed")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
	}

	// Create completion message
	msg := &StreamMessage{
		Type:       StreamMessageTypeComplete,
		StreamID:   sw.streamID,
		ChunkIndex: sw.chunkCount - 1,
		TotalChunks: sw.chunkCount,
		Progress:   100,
		Timestamp:  time.Now(),
	}

	// Set write deadline
	if err := sw.conn.SetWriteDeadline(time.Now().Add(sw.writeWait)); err != nil {
		sw.logger.Error("failed to set write deadline for completion",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
		)
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write completion message
	if err := sw.conn.WriteJSON(msg); err != nil {
		sw.logger.Error("failed to write completion message",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
		)
		return fmt.Errorf("failed to write completion message: %w", err)
	}

	sw.closed = true

	sw.logger.Info("stream completed",
		slog.String("stream_id", sw.streamID),
		slog.Int("total_chunks", sw.chunkCount),
	)

	return nil
}

// Error sends an error message and marks the stream as closed.
func (sw *StreamWriter) Error(ctx context.Context, code, message string, recoverable bool) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.closed {
		return nil // Already closed, don't send duplicate error
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
	}

	// Create error message
	msg := &StreamMessage{
		Type:     StreamMessageTypeError,
		StreamID: sw.streamID,
		Error: &StreamError{
			Code:        code,
			Message:     message,
			Recoverable: recoverable,
		},
		Timestamp: time.Now(),
	}

	// Set write deadline
	if err := sw.conn.SetWriteDeadline(time.Now().Add(sw.writeWait)); err != nil {
		sw.logger.Error("failed to set write deadline for error",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
		)
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write error message
	if err := sw.conn.WriteJSON(msg); err != nil {
		sw.logger.Error("failed to write error message",
			slog.Any("error", err),
			slog.String("stream_id", sw.streamID),
		)
		return fmt.Errorf("failed to write error message: %w", err)
	}

	sw.closed = true

	sw.logger.Error("stream error sent",
		slog.String("stream_id", sw.streamID),
		slog.String("error_code", code),
		slog.String("error_message", message),
		slog.Bool("recoverable", recoverable),
	)

	return nil
}

// IsClosed returns whether the stream writer is closed.
func (sw *StreamWriter) IsClosed() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.closed
}

// ChunkCount returns the number of chunks written.
func (sw *StreamWriter) ChunkCount() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.chunkCount
}

// ============================================================================
// StreamHandler
// ============================================================================

// StreamHandler handles chunked streaming for large datasets over WebSocket.
type StreamHandler struct {
	logger      *slog.Logger
	writeWait   time.Duration
	chunkSize   int
}

// StreamHandlerConfig holds configuration for the StreamHandler.
type StreamHandlerConfig struct {
	// Logger is the structured logger (defaults to slog.Default()).
	Logger *slog.Logger

	// WriteWait is the write timeout (defaults to 10s).
	WriteWait time.Duration

	// ChunkSize is the maximum number of items per chunk (defaults to 100).
	ChunkSize int
}

// NewStreamHandler creates a new StreamHandler.
func NewStreamHandler(cfg StreamHandlerConfig) *StreamHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	if cfg.WriteWait == 0 {
		cfg.WriteWait = 10 * time.Second
	}

	if cfg.ChunkSize == 0 {
		cfg.ChunkSize = 100
	}

	return &StreamHandler{
		logger:    cfg.Logger,
		writeWait: cfg.WriteWait,
		chunkSize: cfg.ChunkSize,
	}
}

// StreamData streams a slice of data as chunks over the WebSocket connection.
// It creates a new StreamWriter for the stream and handles chunking automatically.
func (h *StreamHandler) StreamData(ctx context.Context, conn *websocket.Conn, streamID string, data []interface{}) error {
	if len(data) == 0 {
		// Send completion for empty data
		writer := NewStreamWriter(StreamWriterConfig{
			Conn:      conn,
			StreamID:  streamID,
			Logger:    h.logger,
			WriteWait: h.writeWait,
		})
		return writer.Complete(ctx)
	}

	// Calculate total chunks
	totalChunks := (len(data) + h.chunkSize - 1) / h.chunkSize

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:        conn,
		StreamID:    streamID,
		Logger:      h.logger,
		WriteWait:   h.writeWait,
		TotalChunks: totalChunks,
	})

	// Stream data in chunks
	for i := 0; i < len(data); i += h.chunkSize {
		end := i + h.chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := &StreamChunk{
			Index: i / h.chunkSize,
			Data:  data[i:end],
			IsLast: end >= len(data),
		}

		if err := writer.WriteChunk(ctx, chunk); err != nil {
			// Try to send error message
			_ = writer.Error(ctx, "STREAM_ERROR", err.Error(), false)
			return fmt.Errorf("failed to stream data: %w", err)
		}
	}

	// Mark stream as complete
	if err := writer.Complete(ctx); err != nil {
		return fmt.Errorf("failed to complete stream: %w", err)
	}

	h.logger.Info("data streamed successfully",
		slog.String("stream_id", streamID),
		slog.Int("total_items", len(data)),
		slog.Int("total_chunks", totalChunks),
	)

	return nil
}

// StreamJSONArray streams a JSON array as chunks over the WebSocket connection.
// This is useful for streaming large JSON responses from database queries.
func (h *StreamHandler) StreamJSONArray(ctx context.Context, conn *websocket.Conn, streamID string, jsonData []byte) error {
	// Parse JSON array
	var data []interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON array: %w", err)
	}

	return h.StreamData(ctx, conn, streamID, data)
}

// StreamRows streams database rows as chunks over the WebSocket connection.
// The rowsFunc callback is called for each row to convert to interface{}.
func (h *StreamHandler) StreamRows(ctx context.Context, conn *websocket.Conn, streamID string, totalRows int, rowsFunc func(start, end int) ([]interface{}, error)) error {
	if totalRows == 0 {
		// Send completion for empty data
		writer := NewStreamWriter(StreamWriterConfig{
			Conn:      conn,
			StreamID:  streamID,
			Logger:    h.logger,
			WriteWait: h.writeWait,
		})
		return writer.Complete(ctx)
	}

	// Calculate total chunks
	totalChunks := (totalRows + h.chunkSize - 1) / h.chunkSize

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:        conn,
		StreamID:    streamID,
		Logger:      h.logger,
		WriteWait:   h.writeWait,
		TotalChunks: totalChunks,
	})

	// Stream rows in chunks
	for offset := 0; offset < totalRows; offset += h.chunkSize {
		end := offset + h.chunkSize
		if end > totalRows {
			end = totalRows
		}

		// Fetch rows for this chunk
		rows, err := rowsFunc(offset, end)
		if err != nil {
			_ = writer.Error(ctx, "ROW_FETCH_ERROR", err.Error(), false)
			return fmt.Errorf("failed to fetch rows: %w", err)
		}

		chunk := &StreamChunk{
			Index:  offset / h.chunkSize,
			Data:   rows,
			IsLast: end >= totalRows,
		}

		if err := writer.WriteChunk(ctx, chunk); err != nil {
			_ = writer.Error(ctx, "STREAM_ERROR", err.Error(), false)
			return fmt.Errorf("failed to stream rows: %w", err)
		}
	}

	// Mark stream as complete
	if err := writer.Complete(ctx); err != nil {
		return fmt.Errorf("failed to complete stream: %w", err)
	}

	h.logger.Info("rows streamed successfully",
		slog.String("stream_id", streamID),
		slog.Int("total_rows", totalRows),
		slog.Int("total_chunks", totalChunks),
	)

	return nil
}

// CreateStreamWriter creates a new StreamWriter for custom streaming scenarios.
func (h *StreamHandler) CreateStreamWriter(conn *websocket.Conn, streamID string, totalChunks int) *StreamWriter {
	return NewStreamWriter(StreamWriterConfig{
		Conn:        conn,
		StreamID:    streamID,
		Logger:      h.logger,
		WriteWait:   h.writeWait,
		TotalChunks: totalChunks,
	})
}

// ChunkSize returns the configured chunk size.
func (h *StreamHandler) ChunkSize() int {
	return h.chunkSize
}

// ============================================================================
// Utility Functions
// ============================================================================

// GenerateStreamID generates a unique stream ID.
func GenerateStreamID() string {
	return fmt.Sprintf("stream_%d", time.Now().UnixNano())
}

// ParseStreamMessage parses a JSON byte array into a StreamMessage.
func ParseStreamMessage(data []byte) (*StreamMessage, error) {
	var msg StreamMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse stream message: %w", err)
	}
	return &msg, nil
}
