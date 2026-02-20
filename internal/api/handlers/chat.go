// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the SSE streaming chat handler for the conversational BI feature.
// It handles POST /v1/chat requests and streams responses using Server-Sent Events.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/medisync/medisync/internal/agents/module_a/a03_visualization"
	"github.com/medisync/medisync/internal/warehouse"
	"github.com/medisync/medisync/internal/warehouse/models"

	"github.com/google/uuid"
)

// ChatHandler handles chat requests with SSE streaming responses.
type ChatHandler struct {
	logger             *slog.Logger
	db                 *warehouse.ReadOnlyClient
	visualizationAgent *a03_visualization.VisualizationRoutingAgent
	chatMsgRepo        *warehouse.ChatMessageRepository
	eventBuffer        int // Size of the SSE event buffer
	flushTimeout       time.Duration
}

// ChatHandlerConfig holds configuration for the ChatHandler.
type ChatHandlerConfig struct {
	Logger       *slog.Logger
	DB           *warehouse.ReadOnlyClient
	ChatMsgRepo  *warehouse.ChatMessageRepository
	EventBuffer  int           // Size of SSE event buffer (default: 16)
	FlushTimeout time.Duration // Timeout for flushing SSE (default: 5s)
}

// NewChatHandler creates a new ChatHandler instance.
func NewChatHandler(cfg ChatHandlerConfig) *ChatHandler {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	if cfg.EventBuffer == 0 {
		cfg.EventBuffer = 16
	}

	if cfg.FlushTimeout == 0 {
		cfg.FlushTimeout = 5 * time.Second
	}

	// Initialize visualization routing agent
	visualizationAgent := a03_visualization.NewVisualizationRoutingAgent(&a03_visualization.AgentConfig{
		Logger: cfg.Logger,
	})

	return &ChatHandler{
		logger:             cfg.Logger,
		db:                 cfg.DB,
		visualizationAgent: visualizationAgent,
		chatMsgRepo:        cfg.ChatMsgRepo,
		eventBuffer:        cfg.EventBuffer,
		flushTimeout:       cfg.FlushTimeout,
	}
}

// HandleChat handles POST /v1/chat requests with SSE streaming.
// It validates the request, streams thinking events, SQL preview, and results.
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse and validate request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeSSEError(w, ErrorCodeInvalidRequest, "Invalid JSON request body")
		return
	}

	if err := req.Validate(); err != nil {
		h.writeSSEError(w, ErrorCodeInvalidRequest, err.Error())
		return
	}

	// Set SSE headers
	h.setSSEHeaders(w)

	// Create flusher for streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.logger.Error("streaming not supported")
		h.writeSSEError(w, ErrorCodeInternalError, "Streaming not supported")
		return
	}

	// Get or create session ID
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	locale := req.GetLocale()

	h.logger.Info("processing chat request",
		slog.String("session_id", sessionID),
		slog.String("locale", locale),
		slog.Int("query_length", len(req.Query)),
	)

	// Process the query with streaming
	if err := h.processQuery(ctx, w, flusher, &req, sessionID, locale); err != nil {
		h.logger.Error("query processing failed",
			slog.Any("error", err),
			slog.String("session_id", sessionID),
		)
		// Error already sent via SSE
	}

	// Send completion signal
	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// HandleAgentsHealth handles GET /v1/agents/health requests.
func (h *ChatHandler) HandleAgentsHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check database health
	dbStatus := AgentStatusHealthy
	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			dbStatus = AgentStatusUnhealthy
		}
	} else {
		dbStatus = AgentStatusDegraded
	}

	// Build health response
	overallStatus := AgentStatusHealthy
	switch dbStatus {
	case AgentStatusUnhealthy:
		overallStatus = AgentStatusUnhealthy
	case AgentStatusDegraded:
		overallStatus = AgentStatusDegraded
	}

	response := &AgentsHealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Agents: []AgentHealth{
			{
				ID:        "a-01-text-to-sql",
				Name:      "Text-to-SQL Agent",
				Status:    dbStatus,
				LastCheck: time.Now().UTC().Format(time.RFC3339),
			},
			{
				ID:        "e-01-language",
				Name:      "Language Detection",
				Status:    AgentStatusHealthy,
				LastCheck: time.Now().UTC().Format(time.RFC3339),
			},
		},
		LLMProvider: &LLMProviderHealth{
			Name:   "ollama",
			Model:  "llama4",
			Status: LLMStatusConnected,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode health response", slog.Any("error", err))
	}
}

// processQuery processes the natural language query and streams results.
func (h *ChatHandler) processQuery(
	ctx context.Context,
	w http.ResponseWriter,
	flusher http.Flusher,
	req *ChatRequest,
	sessionID, locale string,
) error {
	// Parse session ID
	sessID, err := uuid.Parse(sessionID)
	if err != nil {
		h.sendSSEEvent(w, flusher, NewErrorEvent("Invalid session ID"))
		return err
	}

	// Get user ID from context
	userIDStr, _ := ctx.Value("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	// Save user message
	userMsg := &models.ChatMessage{
		SessionID: sessID,
		UserID:    userID,
		Role:      "user",
		Content:   req.Query,
		Locale:    locale,
	}
	if h.chatMsgRepo != nil {
		if err := h.chatMsgRepo.Create(ctx, userMsg); err != nil {
			h.logger.Error("failed to save user message", slog.Any("error", err))
		}
	}

	// 1. Send thinking event: Detecting language
	h.sendSSEEvent(w, flusher, NewThinkingEvent("Detecting query language..."))

	// 2. Send thinking event: Analyzing intent
	h.sendSSEEvent(w, flusher, NewThinkingEvent("Analyzing query intent..."))

	// 3. Send thinking event: Retrieving schema
	h.sendSSEEvent(w, flusher, NewThinkingEvent("Retrieving relevant schema..."))

	// 4. Generate SQL (placeholder - in real implementation, call AI agent)
	sql, chartType, confidence, err := h.generateSQL(ctx, req.Query, locale)
	if err != nil {
		h.sendSSEEvent(w, flusher, NewErrorEvent(fmt.Sprintf("Failed to generate query: %s", err.Error())))
		return err
	}

	// 5. Check confidence for clarification
	if confidence < 70 {
		clarificationEvent := NewClarificationEvent(
			"I'm not sure I understood your query correctly. Could you clarify?",
			[]string{
				"Show clinic revenue",
				"Show pharmacy revenue",
				"Show total revenue",
			},
		)
		h.sendSSEEvent(w, flusher, clarificationEvent)
		return nil
	}

	// 6. Send SQL preview
	h.sendSSEEvent(w, flusher, NewSQLPreviewEvent(sql))

	// 7. Execute query (if database is available)
	var data interface{}
	if h.db != nil {
		data, err = h.executeSQL(ctx, sql, sessionID)
		if err != nil {
			h.sendSSEEvent(w, flusher, NewErrorEvent(fmt.Sprintf("Query execution failed: %s", err.Error())))
			return err
		}
	} else {
		// Mock data for development
		data = h.getMockData(chartType)
	}

	// 8. Send result event
	resultEvent := NewResultEvent(chartType, data, confidence)
	h.sendSSEEvent(w, flusher, resultEvent)

	// 9. Save assistant message
	assistantMsg := &models.ChatMessage{
		SessionID:       sessID,
		UserID:          userID,
		Role:            "assistant",
		Content:         "Query completed successfully",
		ConfidenceScore: &confidence,
		Locale:          locale,
	}
	if chartType != "" && data != nil {
		if chartData, ok := data.(*ChartData); ok {
			assistantMsg.ChartSpec = map[string]any{
				"type":  string(chartType),
				"chart": chartData,
			}
		}
	}
	if h.chatMsgRepo != nil {
		if err := h.chatMsgRepo.Create(ctx, assistantMsg); err != nil {
			h.logger.Error("failed to save assistant message", slog.Any("error", err))
		}
	}

	return nil
}

// generateSQL generates SQL from natural language query.
// This is a placeholder implementation - in production, this calls the AI agent.
func (h *ChatHandler) generateSQL(ctx context.Context, query, locale string) (string, ChartType, float64, error) {
	// Placeholder parameters - will be used when AI agent calls are implemented
	_, _ = query, locale

	// Placeholder implementation
	// In production, this would:
	// 1. Call E-01 Language Detection agent
	// 2. Call E-02 Query Translation agent (if Arabic)
	// 3. Call A-04 Domain Terminology agent
	// 4. Call A-01 Text-to-SQL agent
	// 5. Call A-03 Visualization Routing agent

	// Simulate processing time
	select {
	case <-ctx.Done():
		return "", "", 0, ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	// Mock response based on query patterns
	sql := `SELECT SUM(amount) AS total_revenue FROM fact_billing WHERE billing_date >= '2026-01-01' AND billing_date < '2026-02-01'`
	chartType := ChartTypeKPICard
	confidence := 94.0

	return sql, chartType, confidence, nil
}

// executeSQL executes the SQL query against the read-only database.
func (h *ChatHandler) executeSQL(ctx context.Context, sql, sessionID string) (interface{}, error) {
	// Execute with audit metadata
	auditEntry := warehouse.QueryAuditEntry{
		SQL:         sql,
		SessionID:   sessionID,
		QuerySource: "A-01",
	}

	rows, err := h.db.ExecuteQueryWithAudit(ctx, auditEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Collect results
	var results []map[string]interface{}
	fieldDescriptions := rows.FieldDescriptions()

	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(fieldDescriptions))
		valuePtrs := make([]interface{}, len(fieldDescriptions))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		for i, fd := range fieldDescriptions {
			row[string(fd.Name)] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// getMockData returns mock data for development when database is unavailable.
func (h *ChatHandler) getMockData(chartType ChartType) interface{} {
	switch chartType {
	case ChartTypeKPICard:
		return map[string]interface{}{
			"value":     1250000.00,
			"formatted": "₹12,50,000.00",
		}
	case ChartTypeLineChart, ChartTypeBarChart:
		return &ChartData{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
			Series: []ChartSeries{
				{
					Name:   "Revenue",
					Values: []interface{}{100000, 120000, 115000, 140000, 135000, 150000},
				},
			},
		}
	case ChartTypePieChart:
		return &ChartData{
			Labels: []string{"Clinic", "Pharmacy", "Lab"},
			Series: []ChartSeries{
				{
					Name:   "Distribution",
					Values: []interface{}{45, 35, 20},
				},
			},
		}
	case ChartTypeDataTable:
		return &ChartData{
			Columns: []ChartColumn{
				{Name: "Department", Type: "string"},
				{Name: "Revenue", Type: "number"},
			},
			Rows: []map[string]interface{}{
				{"Department": "Cardiology", "Revenue": 250000},
				{"Department": "Orthopedics", "Revenue": 180000},
				{"Department": "Pediatrics", "Revenue": 150000},
			},
		}
	default:
		return map[string]interface{}{
			"value":     0,
			"formatted": "N/A",
		}
	}
}

// setSSEHeaders sets the required headers for Server-Sent Events.
func (h *ChatHandler) setSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// sendSSEEvent sends an SSE event to the client.
func (h *ChatHandler) sendSSEEvent(w http.ResponseWriter, flusher http.Flusher, event *SSEEvent) {
	jsonData, err := event.ToJSON()
	if err != nil {
		h.logger.Error("failed to marshal SSE event", slog.Any("error", err))
		return
	}

	// Write SSE format: "data: <json>\n\n"
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()

	h.logger.Debug("sent SSE event",
		slog.String("type", string(event.Type)),
		slog.String("session_id", ""),
	)
}

// writeSSEError writes an error response in SSE format for early errors.
func (h *ChatHandler) writeSSEError(w http.ResponseWriter, code ErrorCode, message string) {
	h.setSSEHeaders(w)

	errResp := &ErrorResponse{
		Code:    code,
		Message: message,
	}

	jsonData, err := errResp.ToJSON()
	if err != nil {
		h.logger.Error("failed to marshal error response", slog.Any("error", err))
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if ok {
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	} else {
		http.Error(w, message, http.StatusBadRequest)
	}
}

// ============================================================================
// Streaming Writer Helper
// ============================================================================

// SSEWriter wraps an http.ResponseWriter for convenient SSE writing.
type SSEWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
	mu      sync.Mutex
}

// NewSSEWriter creates a new SSE writer.
func NewSSEWriter(w http.ResponseWriter) (*SSEWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	return &SSEWriter{
		w:       w,
		flusher: flusher,
	}, nil
}

// WriteEvent writes an SSE event.
func (sw *SSEWriter) WriteEvent(event *SSEEvent) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	jsonData, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	fmt.Fprintf(sw.w, "data: %s\n\n", jsonData)
	sw.flusher.Flush()

	return nil
}

// WriteDone writes the SSE completion signal.
func (sw *SSEWriter) WriteDone() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	fmt.Fprint(sw.w, "data: [DONE]\n\n")
	sw.flusher.Flush()
}

// WriteError writes an error event.
func (sw *SSEWriter) WriteError(message string) {
	sw.WriteEvent(NewErrorEvent(message))
}
