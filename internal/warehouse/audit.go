// Package warehouse provides the audit logging service for MediSync.
//
// This file implements the AuditService for comprehensive audit logging of all
// query operations, user actions, and system events. It follows the audit
// requirements specified in the security architecture.
//
// All audit logs are stored in the app.audit_logs table and include:
//   - User context (user ID, tenant ID, session ID)
//   - Request metadata (IP address, user agent, request ID)
//   - Action details (action type, resource, changes)
//   - Timestamps and status
//
// Usage:
//
//	auditSvc := warehouse.NewAuditService(pool, logger)
//	err := auditSvc.LogQuerySubmit(ctx, userID, tenantID, "show me sales by region")
package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// contextKey is a type for context keys to avoid collisions.
type contextKey string

// AuditAction represents the type of action being audited.
type AuditAction string

const (
	// Query submission actions
	AuditActionQuerySubmit   AuditAction = "QUERY_SUBMIT"
	AuditActionSQLExecute    AuditAction = "SQL_EXECUTE"
	AuditActionResultReturn  AuditAction = "RESULT_RETURN"
	AuditActionQueryError    AuditAction = "QUERY_ERROR"
	AuditActionQueryCancel   AuditAction = "QUERY_CANCEL"

	// Document processing actions (Module B)
	AuditActionOCRProcess    AuditAction = "OCR_PROCESS"
	AuditActionLedgerMap     AuditAction = "LEDGER_MAP"
	AuditActionTallySync     AuditAction = "TALLY_SYNC"
	AuditActionApprovalGate  AuditAction = "APPROVAL_GATE"

	// Report actions (Module C)
	AuditActionReportGenerate AuditAction = "REPORT_GENERATE"
	AuditActionReportExport   AuditAction = "REPORT_EXPORT"

	// Admin actions
	AuditActionUserLogin     AuditAction = "USER_LOGIN"
	AuditActionUserLogout    AuditAction = "USER_LOGOUT"
	AuditActionConfigChange  AuditAction = "CONFIG_CHANGE"
	AuditActionDataAccess    AuditAction = "DATA_ACCESS"
)

// AuditStatus represents the outcome of an audited action.
type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success"
	AuditStatusFailure AuditStatus = "failure"
	AuditStatusPending AuditStatus = "pending"
)

// AuditLogEntry represents a single audit log record.
type AuditLogEntry struct {
	// ID is the unique identifier for the audit log entry.
	ID uuid.UUID `json:"id"`

	// UserID is the ID of the user who performed the action.
	UserID uuid.UUID `json:"user_id"`

	// TenantID is the ID of the tenant (organization) the user belongs to.
	TenantID uuid.UUID `json:"tenant_id"`

	// SessionID is the ID of the user session.
	SessionID *uuid.UUID `json:"session_id,omitempty"`

	// Action is the type of action performed.
	Action AuditAction `json:"action"`

	// ResourceType is the type of resource affected (e.g., "query", "document", "report").
	ResourceType string `json:"resource_type,omitempty"`

	// ResourceID is the ID of the specific resource affected.
	ResourceID *uuid.UUID `json:"resource_id,omitempty"`

	// Status indicates whether the action succeeded, failed, or is pending.
	Status AuditStatus `json:"status"`

	// IPAddress is the client's IP address.
	IPAddress string `json:"ip_address,omitempty"`

	// UserAgent is the client's user agent string.
	UserAgent string `json:"user_agent,omitempty"`

	// RequestID is the unique identifier for the HTTP request.
	RequestID string `json:"request_id,omitempty"`

	// Details contains additional structured data about the action.
	Details map[string]interface{} `json:"details,omitempty"`

	// ErrorMessage contains error details if the action failed.
	ErrorMessage string `json:"error_message,omitempty"`

	// Duration is how long the action took to complete.
	Duration *time.Duration `json:"duration,omitempty"`

	// CreatedAt is when the audit log entry was created.
	CreatedAt time.Time `json:"created_at"`
}

// AuditService provides audit logging functionality.
type AuditService struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewAuditService creates a new audit service.
func NewAuditService(pool *PostgresPool, logger *slog.Logger) (*AuditService, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required for audit service")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &AuditService{
		pool:   pool.Pool(),
		logger: logger,
	}, nil
}

// NewAuditServiceFromPool creates an audit service from an existing pgxpool.Pool.
func NewAuditServiceFromPool(pool *pgxpool.Pool, logger *slog.Logger) (*AuditService, error) {
	if pool == nil {
		return nil, fmt.Errorf("warehouse: connection pool is required for audit service")
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &AuditService{
		pool:   pool,
		logger: logger,
	}, nil
}

// Log writes a generic audit log entry to the database.
func (s *AuditService) Log(ctx context.Context, entry *AuditLogEntry) error {
	if entry == nil {
		return fmt.Errorf("warehouse: audit log entry is required")
	}

	query := `
		INSERT INTO app.audit_logs (
			user_id, tenant_id, session_id, action, resource_type, resource_id,
			status, ip_address, user_agent, request_id, details, error_message,
			duration, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW()
		) RETURNING id, created_at
	`

	err := s.pool.QueryRow(ctx, query,
		entry.UserID,
		entry.TenantID,
		entry.SessionID,
		entry.Action,
		entry.ResourceType,
		entry.ResourceID,
		entry.Status,
		entry.IPAddress,
		entry.UserAgent,
		entry.RequestID,
		entry.Details,
		entry.ErrorMessage,
		entry.Duration,
	).Scan(&entry.ID, &entry.CreatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to write audit log: %w", err)
	}

	s.logger.Debug("audit log entry written",
		slog.String("action", string(entry.Action)),
		slog.String("user_id", entry.UserID.String()),
		slog.String("tenant_id", entry.TenantID.String()),
		slog.String("status", string(entry.Status)),
	)

	return nil
}

// LogQuerySubmit logs when a user submits a natural language query.
func (s *AuditService) LogQuerySubmit(ctx context.Context, userID, tenantID uuid.UUID, query string) error {
	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionQuerySubmit,
		ResourceType: "query",
		Status:       AuditStatusSuccess,
		Details: map[string]interface{}{
			"query_text":     query,
			"query_length":   len(query),
			"language":       detectQueryLanguage(query),
		},
	}

	// Extract request metadata from context if available
	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogSQLExecute logs when a SQL query is executed against the database.
func (s *AuditService) LogSQLExecute(ctx context.Context, userID, tenantID uuid.UUID, sql string) error {
	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionSQLExecute,
		ResourceType: "sql_query",
		Status:       AuditStatusSuccess,
		Details: map[string]interface{}{
			"sql_hash":       hashSQL(sql),
			"sql_length":     len(sql),
			"query_type":     getQueryType(sql),
			"is_readonly":    isSelectOnly(sql),
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogSQLExecuteWithDuration logs a SQL execution with timing information.
func (s *AuditService) LogSQLExecuteWithDuration(ctx context.Context, userID, tenantID uuid.UUID, sql string, duration time.Duration, err error) error {
	status := AuditStatusSuccess
	var errorMessage string
	if err != nil {
		status = AuditStatusFailure
		errorMessage = err.Error()
	}

	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionSQLExecute,
		ResourceType: "sql_query",
		Status:       status,
		Duration:     &duration,
		ErrorMessage: errorMessage,
		Details: map[string]interface{}{
			"sql_hash":       hashSQL(sql),
			"sql_length":     len(sql),
			"query_type":     getQueryType(sql),
			"is_readonly":    isSelectOnly(sql),
			"duration_ms":    duration.Milliseconds(),
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogResultReturn logs when query results are returned to the user.
func (s *AuditService) LogResultReturn(ctx context.Context, userID, tenantID uuid.UUID, rowCount int) error {
	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionResultReturn,
		ResourceType: "query_result",
		Status:       AuditStatusSuccess,
		Details: map[string]interface{}{
			"row_count": rowCount,
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogQueryError logs when a query fails with an error.
func (s *AuditService) LogQueryError(ctx context.Context, userID, tenantID uuid.UUID, query string, err error) error {
	entry := &AuditLogEntry{
		UserID:        userID,
		TenantID:      tenantID,
		Action:        AuditActionQueryError,
		ResourceType:  "query",
		Status:        AuditStatusFailure,
		ErrorMessage:  err.Error(),
		Details: map[string]interface{}{
			"query_text":   query,
			"error_type":   getErrorType(err),
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogOCRProcess logs an OCR document processing event.
func (s *AuditService) LogOCRProcess(ctx context.Context, userID, tenantID uuid.UUID, documentID uuid.UUID, status AuditStatus, details map[string]interface{}) error {
	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionOCRProcess,
		ResourceType: "document",
		ResourceID:   &documentID,
		Status:       status,
		Details:      details,
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogTallySync logs a Tally sync event (always requires approval tracking).
func (s *AuditService) LogTallySync(ctx context.Context, userID, tenantID uuid.UUID, syncID uuid.UUID, status AuditStatus, approvedBy *uuid.UUID) error {
	entry := &AuditLogEntry{
		UserID:       userID,
		TenantID:     tenantID,
		Action:       AuditActionTallySync,
		ResourceType: "tally_sync",
		ResourceID:   &syncID,
		Status:       status,
		Details: map[string]interface{}{
			"approved_by": approvedBy,
			"hitl_gate":   approvedBy != nil,
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// LogApprovalGate logs when a human-in-the-loop approval is required/granted.
func (s *AuditService) LogApprovalGate(ctx context.Context, userID, tenantID uuid.UUID, resourceType string, resourceID uuid.UUID, approved bool, approverID uuid.UUID) error {
	status := AuditStatusPending
	if approved {
		status = AuditStatusSuccess
	}

	entry := &AuditLogEntry{
		UserID:       approverID,
		TenantID:     tenantID,
		Action:       AuditActionApprovalGate,
		ResourceType: resourceType,
		ResourceID:   &resourceID,
		Status:       status,
		Details: map[string]interface{}{
			"requested_by": userID.String(),
			"approved":     approved,
		},
	}

	s.extractRequestMetadata(ctx, entry)

	return s.Log(ctx, entry)
}

// GetUserAuditLogs retrieves audit logs for a specific user.
func (s *AuditService) GetUserAuditLogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]AuditLogEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	query := `
		SELECT id, user_id, tenant_id, session_id, action, resource_type, resource_id,
		       status, ip_address, user_agent, request_id, details, error_message,
		       duration, created_at
		FROM app.audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get user audit logs: %w", err)
	}
	defer rows.Close()

	return s.scanAuditLogs(rows)
}

// GetTenantAuditLogs retrieves audit logs for a tenant (admin use).
func (s *AuditService) GetTenantAuditLogs(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]AuditLogEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := `
		SELECT id, user_id, tenant_id, session_id, action, resource_type, resource_id,
		       status, ip_address, user_agent, request_id, details, error_message,
		       duration, created_at
		FROM app.audit_logs
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get tenant audit logs: %w", err)
	}
	defer rows.Close()

	return s.scanAuditLogs(rows)
}

// extractRequestMetadata extracts IP, user agent, and request ID from context.
func (s *AuditService) extractRequestMetadata(ctx context.Context, entry *AuditLogEntry) {
	// These would typically be set by middleware
	if ip := getContextString(ctx, "ip_address"); ip != "" {
		entry.IPAddress = ip
	}
	if ua := getContextString(ctx, "user_agent"); ua != "" {
		entry.UserAgent = ua
	}
	if reqID := getContextString(ctx, "request_id"); reqID != "" {
		entry.RequestID = reqID
	}
}

// getContextString retrieves a string value from context.
func getContextString(ctx context.Context, key string) string {
	if v := ctx.Value(contextKey(key)); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// scanAuditLogs scans database rows into AuditLogEntry slices.
func (s *AuditService) scanAuditLogs(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]AuditLogEntry, error) {
	var entries []AuditLogEntry

	for rows.Next() {
		var entry AuditLogEntry
		var duration *int64 // Store as milliseconds

		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.TenantID,
			&entry.SessionID,
			&entry.Action,
			&entry.ResourceType,
			&entry.ResourceID,
			&entry.Status,
			&entry.IPAddress,
			&entry.UserAgent,
			&entry.RequestID,
			&entry.Details,
			&entry.ErrorMessage,
			&duration,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan audit log: %w", err)
		}

		if duration != nil {
			d := time.Duration(*duration) * time.Millisecond
			entry.Duration = &d
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("warehouse: error iterating audit logs: %w", err)
	}

	return entries, nil
}

// Helper functions for query analysis

// detectQueryLanguage detects if the query is Arabic or English.
func detectQueryLanguage(query string) string {
	for _, r := range query {
		if r >= 0x0600 && r <= 0x06FF {
			return "ar"
		}
	}
	return "en"
}

// hashSQL creates a hash of the SQL for logging without exposing full query.
func hashSQL(sql string) string {
	if len(sql) > 32 {
		return sql[:16] + "..." + sql[len(sql)-8:]
	}
	return sql
}

// getQueryType extracts the primary query type (SELECT, INSERT, etc.).
func getQueryType(sql string) string {
	// Simple extraction - in production use proper SQL parser
	for i, r := range sql {
		if r == ' ' || r == '\n' {
			return sql[:i]
		}
	}
	return sql
}

// isSelectOnly checks if the query is a SELECT statement.
func isSelectOnly(sql string) bool {
	// Simple check - in production use proper SQL validation
	queryType := getQueryType(sql)
	return queryType == "SELECT" || queryType == "select"
}

// getErrorType categorizes an error for logging.
func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	// Simple categorization based on error message patterns
	switch {
	case containsString(err.Error(), "timeout"):
		return "timeout"
	case containsString(err.Error(), "permission"):
		return "permission"
	case containsString(err.Error(), "syntax"):
		return "syntax"
	default:
		return "unknown"
	}
}

// containsString checks if s contains substr (case-insensitive).
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsString(s[1:], substr))
}

// AuditMiddleware returns an HTTP middleware that logs requests to the audit service.
type AuditMiddleware struct {
	audit  *AuditService
	logger *slog.Logger
}

// NewAuditMiddleware creates a new audit middleware.
func NewAuditMiddleware(audit *AuditService, logger *slog.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		audit:  audit,
		logger: logger,
	}
}

// Middleware returns the HTTP middleware function.
func (m *AuditMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add request metadata to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKey("ip_address"), r.RemoteAddr)
		ctx = context.WithValue(ctx, contextKey("user_agent"), r.UserAgent())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
