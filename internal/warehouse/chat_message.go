// Package warehouse provides database repository implementations.
package warehouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/warehouse/models"
)

// ChatMessageRepository handles database operations for chat messages.
type ChatMessageRepository struct {
	db *Repo
}

// NewChatMessageRepository creates a new chat message repository.
func NewChatMessageRepository(db *Repo) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

// Create inserts a new chat message.
func (r *ChatMessageRepository) Create(ctx context.Context, msg *models.ChatMessage) error {
	query := `
		INSERT INTO app.chat_messages (session_id, user_id, role, content, chart_spec, table_data, drilldown_query, confidence_score, locale)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	err := r.db.pool.QueryRow(ctx, query,
		msg.SessionID, msg.UserID, msg.Role, msg.Content,
		msg.ChartSpec, msg.TableData, msg.DrilldownQuery,
		msg.ConfidenceScore, msg.Locale,
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("warehouse: failed to create chat message: %w", err)
	}

	return nil
}

// GetBySessionID retrieves messages for a session, ordered by creation time.
func (r *ChatMessageRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, session_id, user_id, role, content, chart_spec, table_data, drilldown_query, confidence_score, locale, created_at
		FROM app.chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.pool.Query(ctx, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		msg := &models.ChatMessage{}
		err := rows.Scan(
			&msg.ID, &msg.SessionID, &msg.UserID, &msg.Role, &msg.Content,
			&msg.ChartSpec, &msg.TableData, &msg.DrilldownQuery,
			&msg.ConfidenceScore, &msg.Locale, &msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetRecentByUserID retrieves recent messages for a user across all sessions.
func (r *ChatMessageRepository) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, session_id, user_id, role, content, chart_spec, table_data, drilldown_query, confidence_score, locale, created_at
		FROM app.chat_messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("warehouse: failed to get recent chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		msg := &models.ChatMessage{}
		err := rows.Scan(
			&msg.ID, &msg.SessionID, &msg.UserID, &msg.Role, &msg.Content,
			&msg.ChartSpec, &msg.TableData, &msg.DrilldownQuery,
			&msg.ConfidenceScore, &msg.Locale, &msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("warehouse: failed to scan chat message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
