// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the WebSocket handler for real-time bidirectional chat.
// It handles GET /v1/ws/chat requests for WebSocket connections.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/medisync/medisync/internal/warehouse"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader configuration.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

// WSHandler handles WebSocket connections for chat.
type WSHandler struct {
	logger       *slog.Logger
	db           *warehouse.ReadOnlyClient
	chatHandler  *ChatHandler
	pingInterval time.Duration
	pongWait     time.Duration
	writeWait    time.Duration
}

// WSHandlerConfig holds configuration for the WSHandler.
type WSHandlerConfig struct {
	Logger       *slog.Logger
	DB           *warehouse.ReadOnlyClient
	ChatHandler  *ChatHandler
	PingInterval time.Duration // Ping interval (default: 30s)
	PongWait     time.Duration // Max wait for pong (default: 60s)
	WriteWait    time.Duration // Write timeout (default: 10s)
}

// NewWSHandler creates a new WebSocket handler.
func NewWSHandler(cfg WSHandlerConfig) *WSHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	if cfg.PingInterval == 0 {
		cfg.PingInterval = 30 * time.Second
	}

	if cfg.PongWait == 0 {
		cfg.PongWait = 60 * time.Second
	}

	if cfg.WriteWait == 0 {
		cfg.WriteWait = 10 * time.Second
	}

	return &WSHandler{
		logger:       cfg.Logger,
		db:           cfg.DB,
		chatHandler:  cfg.ChatHandler,
		pingInterval: cfg.PingInterval,
		pongWait:     cfg.PongWait,
		writeWait:    cfg.WriteWait,
	}
}

// HandleWebSocket handles WebSocket upgrade and message processing.
// It expects a JWT token in the query parameter.
func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract and validate JWT token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		h.logger.Warn("websocket connection attempted without token")
		http.Error(w, "token query parameter required", http.StatusUnauthorized)
		return
	}

	// In production, validate the JWT token here
	// For now, we'll just log it
	h.logger.Debug("websocket connection with token", slog.String("token_len", fmt.Sprintf("%d", len(token))))

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", slog.Any("error", err))
		return
	}
	defer conn.Close()

	h.logger.Info("websocket connection established",
		slog.String("remote_addr", r.RemoteAddr),
	)

	// Create client session
	client := &WSClient{
		conn:      conn,
		handler:   h,
		sessionID: uuid.New().String(),
		sendCh:    make(chan *WSMessage, 256),
		mu:        sync.Mutex{},
	}

	// Start read and write pumps
	go client.writePump()
	client.readPump()
}

// WSClient represents a connected WebSocket client.
type WSClient struct {
	conn      *websocket.Conn
	handler   *WSHandler
	sessionID string
	locale    string
	sendCh    chan *WSMessage
	mu        sync.Mutex
}

// readPump handles incoming messages from the WebSocket connection.
func (c *WSClient) readPump() {
	defer func() {
		c.handler.logger.Info("websocket client disconnected",
			slog.String("session_id", c.sessionID),
		)
		c.conn.Close()
	}()

	// Set read limits and deadlines
	c.conn.SetReadLimit(8192) // Max message size
	c.conn.SetReadDeadline(time.Now().Add(c.handler.pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.handler.pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.handler.logger.Error("websocket read error",
					slog.Any("error", err),
					slog.String("session_id", c.sessionID),
				)
			}
			break
		}

		c.handleMessage(message)
	}
}

// writePump handles outgoing messages to the WebSocket connection.
func (c *WSClient) writePump() {
	ticker := time.NewTicker(c.handler.pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCh:
			c.conn.SetWriteDeadline(time.Now().Add(c.handler.writeWait))
			if !ok {
				// Channel closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.writeJSON(message); err != nil {
				c.handler.logger.Error("websocket write error",
					slog.Any("error", err),
					slog.String("session_id", c.sessionID),
				)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.handler.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.handler.logger.Error("websocket ping error",
					slog.Any("error", err),
					slog.String("session_id", c.sessionID),
				)
				return
			}
		}
	}
}

// writeJSON safely writes a JSON message to the connection.
func (c *WSClient) writeJSON(message *WSMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(message)
}

// handleMessage processes an incoming WebSocket message.
func (c *WSClient) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("Invalid message format")
		return
	}

	switch msg.Type {
	case WSMessageTypeQuery:
		c.handleQuery(msg.Payload)
	case WSMessageTypePong:
		// Pong handled by SetPongHandler
	default:
		c.sendError(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleQuery processes a query message.
func (c *WSClient) handleQuery(payload *WSMessagePayload) {
	if payload == nil || payload.Query == "" {
		c.sendError("Query is required")
		return
	}

	// Set locale from payload or default
	if payload.Locale != "" {
		c.locale = payload.Locale
	} else if c.locale == "" {
		c.locale = "en"
	}

	c.handler.logger.Info("processing websocket query",
		slog.String("session_id", c.sessionID),
		slog.String("query", payload.Query),
		slog.String("locale", c.locale),
	)

	// Send thinking message
	c.send(&WSMessage{
		Type: WSMessageTypeThinking,
		Payload: &WSMessagePayload{
			Message: "Processing your query...",
		},
	})

	// Process the query (reuse ChatHandler logic)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Generate SQL (placeholder - uses same logic as ChatHandler)
	sql, chartType, confidence, err := c.handler.chatHandler.generateSQL(ctx, payload.Query, c.locale)
	if err != nil {
		c.sendError(fmt.Sprintf("Failed to generate query: %s", err.Error()))
		return
	}

	// Check confidence for clarification
	if confidence < 70 {
		c.send(&WSMessage{
			Type: WSMessageTypeClarification,
			Payload: &WSMessagePayload{
				Message: "I'm not sure I understood your query correctly. Could you clarify?",
				Options: []string{
					"Show clinic revenue",
					"Show pharmacy revenue",
					"Show total revenue",
				},
			},
		})
		return
	}

	// Send SQL preview
	c.send(&WSMessage{
		Type: WSMessageTypeSQLPreview,
		Payload: &WSMessagePayload{
			SQL: sql,
		},
	})

	// Get data (mock for now)
	data := c.handler.chatHandler.getMockData(chartType)

	// Send result
	c.send(&WSMessage{
		Type: WSMessageTypeResult,
		Payload: &WSMessagePayload{
			ChartType:  string(chartType),
			Data:       data.(map[string]interface{}),
			Confidence: confidence,
		},
	})
}

// send sends a message to the client.
func (c *WSClient) send(msg *WSMessage) {
	select {
	case c.sendCh <- msg:
	default:
		c.handler.logger.Warn("websocket send buffer full",
			slog.String("session_id", c.sessionID),
		)
	}
}

// sendError sends an error message to the client.
func (c *WSClient) sendError(message string) {
	c.send(&WSMessage{
		Type: WSMessageTypeError,
		Payload: &WSMessagePayload{
			Message: message,
		},
	})
}
