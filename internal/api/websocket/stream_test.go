// Package websocket provides WebSocket utilities for the MediSync API.
//
// This file contains tests for the stream handler.
package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testWSServer creates a test WebSocket server that echoes messages.
func testWSServer(t *testing.T, handler func(*websocket.Conn)) (*httptest.Server, string) {
	t.Helper()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		defer conn.Close()

		handler(conn)
	}))

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	return srv, wsURL
}

// TestStreamWriter_WriteChunk tests writing chunks to a WebSocket.
func TestStreamWriter_WriteChunk(t *testing.T) {
	var receivedMessages []*StreamMessage
	var mu sync.Mutex

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			mu.Lock()
			receivedMessages = append(receivedMessages, &msg)
			mu.Unlock()
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test-stream-1",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	ctx := context.Background()

	// Write chunks
	chunk1 := &StreamChunk{Index: 0, Data: []string{"item1", "item2"}, IsLast: false}
	err = writer.WriteChunk(ctx, chunk1)
	require.NoError(t, err)

	chunk2 := &StreamChunk{Index: 1, Data: []string{"item3", "item4"}, IsLast: true}
	err = writer.WriteChunk(ctx, chunk2)
	require.NoError(t, err)

	// Complete stream
	err = writer.Complete(ctx)
	require.NoError(t, err)

	// Wait for messages to be received
	time.Sleep(100 * time.Millisecond)

	// Verify messages
	mu.Lock()
	defer mu.Unlock()

	require.Len(t, receivedMessages, 3) // 2 chunks + 1 complete

	// Check first chunk
	assert.Equal(t, StreamMessageTypeChunk, receivedMessages[0].Type)
	assert.Equal(t, "test-stream-1", receivedMessages[0].StreamID)
	assert.Equal(t, 0, receivedMessages[0].ChunkIndex)

	// Check second chunk
	assert.Equal(t, StreamMessageTypeChunk, receivedMessages[1].Type)
	assert.Equal(t, 1, receivedMessages[1].ChunkIndex)

	// Check completion
	assert.Equal(t, StreamMessageTypeComplete, receivedMessages[2].Type)
	assert.Equal(t, float64(100), receivedMessages[2].Progress)
}

// TestStreamWriter_Error tests sending error messages.
func TestStreamWriter_Error(t *testing.T) {
	var receivedMessage *StreamMessage

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		var msg StreamMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		receivedMessage = &msg
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test-stream-error",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	ctx := context.Background()

	// Send error
	err = writer.Error(ctx, "TEST_ERROR", "Something went wrong", true)
	require.NoError(t, err)

	// Wait for message
	time.Sleep(100 * time.Millisecond)

	// Verify error message
	require.NotNil(t, receivedMessage)
	assert.Equal(t, StreamMessageTypeError, receivedMessage.Type)
	assert.Equal(t, "test-stream-error", receivedMessage.StreamID)
	require.NotNil(t, receivedMessage.Error)
	assert.Equal(t, "TEST_ERROR", receivedMessage.Error.Code)
	assert.Equal(t, "Something went wrong", receivedMessage.Error.Message)
	assert.True(t, receivedMessage.Error.Recoverable)
}

// TestStreamWriter_ThreadSafety tests concurrent writes.
func TestStreamWriter_ThreadSafety(t *testing.T) {
	var receivedCount int
	var mu sync.Mutex

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			mu.Lock()
			receivedCount++
			mu.Unlock()
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:        conn,
		StreamID:    "test-stream-concurrent",
		Logger:      slog.Default(),
		WriteWait:   5 * time.Second,
		TotalChunks: 10,
	})

	ctx := context.Background()

	// Write chunks concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			chunk := &StreamChunk{
				Index: index,
				Data:  map[string]int{"index": index},
				IsLast: index == 9,
			}
			_ = writer.WriteChunk(ctx, chunk)
		}(i)
	}

	wg.Wait()

	// Complete stream
	err = writer.Complete(ctx)
	require.NoError(t, err)

	// Wait for messages
	time.Sleep(200 * time.Millisecond)

	// Verify all chunks were written
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 11, receivedCount) // 10 chunks + 1 complete
}

// TestStreamWriter_ClosedWriter tests that a closed writer rejects new writes.
func TestStreamWriter_ClosedWriter(t *testing.T) {
	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test-stream-closed",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	ctx := context.Background()

	// Complete stream
	err = writer.Complete(ctx)
	require.NoError(t, err)

	// Try to write after close
	chunk := &StreamChunk{Index: 0, Data: "test", IsLast: true}
	err = writer.WriteChunk(ctx, chunk)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "closed")

	// Try to complete again
	err = writer.Complete(ctx)
	assert.Error(t, err)
}

// TestStreamWriter_ContextCancellation tests context cancellation handling.
func TestStreamWriter_ContextCancellation(t *testing.T) {
	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		// Keep connection open but don't read
		time.Sleep(2 * time.Second)
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream writer
	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test-stream-context",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Try to write with cancelled context
	chunk := &StreamChunk{Index: 0, Data: "test", IsLast: true}
	err = writer.WriteChunk(ctx, chunk)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cancelled")
}

// TestStreamHandler_StreamData tests the StreamHandler with data streaming.
func TestStreamHandler_StreamData(t *testing.T) {
	var receivedMessages []*StreamMessage
	var mu sync.Mutex

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			mu.Lock()
			receivedMessages = append(receivedMessages, &msg)
			mu.Unlock()
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream handler
	handler := NewStreamHandler(StreamHandlerConfig{
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
		ChunkSize: 3, // Small chunk size for testing
	})

	ctx := context.Background()

	// Create test data (10 items)
	data := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		data[i] = map[string]int{"id": i, "value": i * 10}
	}

	// Stream data
	err = handler.StreamData(ctx, conn, "test-handler-stream", data)
	require.NoError(t, err)

	// Wait for messages
	time.Sleep(200 * time.Millisecond)

	// Verify messages
	mu.Lock()
	defer mu.Unlock()

	// 4 chunks (3, 3, 3, 1) + 1 complete = 5 messages
	require.Len(t, receivedMessages, 5)

	// Check chunk messages
	for i := 0; i < 4; i++ {
		assert.Equal(t, StreamMessageTypeChunk, receivedMessages[i].Type)
		assert.Equal(t, "test-handler-stream", receivedMessages[i].StreamID)
		assert.Equal(t, i, receivedMessages[i].ChunkIndex)
		assert.Equal(t, 4, receivedMessages[i].TotalChunks)
	}

	// Check completion
	assert.Equal(t, StreamMessageTypeComplete, receivedMessages[4].Type)
	assert.Equal(t, float64(100), receivedMessages[4].Progress)
}

// TestStreamHandler_StreamData_Empty tests streaming empty data.
func TestStreamHandler_StreamData_Empty(t *testing.T) {
	var receivedMessage *StreamMessage

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		var msg StreamMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		receivedMessage = &msg
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream handler
	handler := NewStreamHandler(StreamHandlerConfig{
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
		ChunkSize: 10,
	})

	ctx := context.Background()

	// Stream empty data
	err = handler.StreamData(ctx, conn, "test-empty-stream", []interface{}{})
	require.NoError(t, err)

	// Wait for message
	time.Sleep(100 * time.Millisecond)

	// Verify completion message was sent
	require.NotNil(t, receivedMessage)
	assert.Equal(t, StreamMessageTypeComplete, receivedMessage.Type)
	assert.Equal(t, "test-empty-stream", receivedMessage.StreamID)
}

// TestStreamHandler_StreamJSONArray tests streaming JSON arrays.
func TestStreamHandler_StreamJSONArray(t *testing.T) {
	var receivedMessages []*StreamMessage
	var mu sync.Mutex

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			mu.Lock()
			receivedMessages = append(receivedMessages, &msg)
			mu.Unlock()
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream handler
	handler := NewStreamHandler(StreamHandlerConfig{
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
		ChunkSize: 2,
	})

	ctx := context.Background()

	// Create JSON array
	jsonData := []byte(`[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5}]`)

	// Stream JSON
	err = handler.StreamJSONArray(ctx, conn, "test-json-stream", jsonData)
	require.NoError(t, err)

	// Wait for messages
	time.Sleep(100 * time.Millisecond)

	// Verify messages
	mu.Lock()
	defer mu.Unlock()

	// 3 chunks (2, 2, 1) + 1 complete = 4 messages
	require.Len(t, receivedMessages, 4)
}

// TestStreamHandler_StreamRows tests streaming rows with a callback.
func TestStreamHandler_StreamRows(t *testing.T) {
	var receivedMessages []*StreamMessage
	var mu sync.Mutex

	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			mu.Lock()
			receivedMessages = append(receivedMessages, &msg)
			mu.Unlock()
		}
	})
	defer srv.Close()

	// Connect to server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Create stream handler
	handler := NewStreamHandler(StreamHandlerConfig{
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
		ChunkSize: 5,
	})

	ctx := context.Background()

	// Mock row fetcher
	totalRows := 12
	rowsFunc := func(start, end int) ([]interface{}, error) {
		rows := make([]interface{}, end-start)
		for i := start; i < end; i++ {
			rows[i-start] = map[string]int{"row": i}
		}
		return rows, nil
	}

	// Stream rows
	err = handler.StreamRows(ctx, conn, "test-rows-stream", totalRows, rowsFunc)
	require.NoError(t, err)

	// Wait for messages
	time.Sleep(100 * time.Millisecond)

	// Verify messages
	mu.Lock()
	defer mu.Unlock()

	// 3 chunks (5, 5, 2) + 1 complete = 4 messages
	require.Len(t, receivedMessages, 4)

	// Verify progress increases
	for i := 0; i < 3; i++ {
		assert.Equal(t, StreamMessageTypeChunk, receivedMessages[i].Type)
		if i > 0 {
			assert.Greater(t, receivedMessages[i].Progress, receivedMessages[i-1].Progress)
		}
	}
}

// TestGenerateStreamID tests stream ID generation.
func TestGenerateStreamID(t *testing.T) {
	id1 := GenerateStreamID()
	id2 := GenerateStreamID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.True(t, strings.HasPrefix(id1, "stream_"))
	assert.True(t, strings.HasPrefix(id2, "stream_"))
}

// TestParseStreamMessage tests parsing stream messages.
func TestParseStreamMessage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *StreamMessage)
	}{
		{
			name: "chunk message",
			input: `{"type":"stream_chunk","stream_id":"test","chunk_index":0,"total_chunks":5,"data":[1,2,3],"progress":20}`,
			wantErr: false,
			check: func(t *testing.T, msg *StreamMessage) {
				assert.Equal(t, StreamMessageTypeChunk, msg.Type)
				assert.Equal(t, "test", msg.StreamID)
				assert.Equal(t, 0, msg.ChunkIndex)
				assert.Equal(t, 5, msg.TotalChunks)
				assert.Equal(t, float64(20), msg.Progress)
			},
		},
		{
			name: "complete message",
			input: `{"type":"stream_complete","stream_id":"test","progress":100}`,
			wantErr: false,
			check: func(t *testing.T, msg *StreamMessage) {
				assert.Equal(t, StreamMessageTypeComplete, msg.Type)
				assert.Equal(t, float64(100), msg.Progress)
			},
		},
		{
			name: "error message",
			input: `{"type":"stream_error","stream_id":"test","error":{"code":"ERR001","message":"test error","recoverable":true}}`,
			wantErr: false,
			check: func(t *testing.T, msg *StreamMessage) {
				assert.Equal(t, StreamMessageTypeError, msg.Type)
				require.NotNil(t, msg.Error)
				assert.Equal(t, "ERR001", msg.Error.Code)
				assert.Equal(t, "test error", msg.Error.Message)
				assert.True(t, msg.Error.Recoverable)
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseStreamMessage([]byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, msg)
			if tt.check != nil {
				tt.check(t, msg)
			}
		})
	}
}

// TestStreamMessage_JSON tests JSON marshaling/unmarshaling of StreamMessage.
func TestStreamMessage_JSON(t *testing.T) {
	original := &StreamMessage{
		Type:       StreamMessageTypeChunk,
		StreamID:   "test-stream",
		ChunkIndex: 2,
		TotalChunks: 5,
		Data:       map[string]interface{}{"key": "value"},
		Progress:   40.5,
		Timestamp:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	// Marshal
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var parsed StreamMessage
	err = json.Unmarshal(jsonData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.StreamID, parsed.StreamID)
	assert.Equal(t, original.ChunkIndex, parsed.ChunkIndex)
	assert.Equal(t, original.TotalChunks, parsed.TotalChunks)
	assert.Equal(t, original.Progress, parsed.Progress)
}

// TestStreamWriter_ChunkCount tests the ChunkCount method.
func TestStreamWriter_ChunkCount(t *testing.T) {
	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
		}
	})
	defer srv.Close()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	ctx := context.Background()

	// Initially zero
	assert.Equal(t, 0, writer.ChunkCount())

	// Write 3 chunks
	for i := 0; i < 3; i++ {
		chunk := &StreamChunk{Index: i, Data: i, IsLast: i == 2}
		err := writer.WriteChunk(ctx, chunk)
		require.NoError(t, err)
		assert.Equal(t, i+1, writer.ChunkCount())
	}
}

// TestStreamHandler_ChunkSize tests the ChunkSize method.
func TestStreamHandler_ChunkSize(t *testing.T) {
	handler := NewStreamHandler(StreamHandlerConfig{
		ChunkSize: 50,
	})
	assert.Equal(t, 50, handler.ChunkSize())

	// Default chunk size
	handler2 := NewStreamHandler(StreamHandlerConfig{})
	assert.Equal(t, 100, handler2.ChunkSize())
}

// TestStreamWriter_IsClosed tests the IsClosed method.
func TestStreamWriter_IsClosed(t *testing.T) {
	srv, wsURL := testWSServer(t, func(conn *websocket.Conn) {
		for {
			var msg StreamMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
		}
	})
	defer srv.Close()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	writer := NewStreamWriter(StreamWriterConfig{
		Conn:      conn,
		StreamID:  "test",
		Logger:    slog.Default(),
		WriteWait: 5 * time.Second,
	})

	assert.False(t, writer.IsClosed())

	err = writer.Complete(context.Background())
	require.NoError(t, err)

	assert.True(t, writer.IsClosed())
}
