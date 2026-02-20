// Package models provides data models for the warehouse layer.
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ValidLocales contains the supported locale codes.
var ValidLocales = map[string]bool{
	"en": true,
	"ar": true,
}

// QuerySession represents a user's query session within the MediSync platform.
// A session groups multiple natural language queries and their results together,
// maintaining context and user preferences throughout the conversation.
type QuerySession struct {
	// ID is the unique identifier for the session
	ID uuid.UUID `json:"id" db:"id"`
	// UserID is the identifier of the user who owns this session
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	// TenantID is the identifier of the tenant (organization) this session belongs to
	TenantID uuid.UUID `json:"tenant_id" db:"tenant_id"`
	// Locale is the user's preferred language code ("en" for English, "ar" for Arabic)
	Locale string `json:"locale" db:"locale"`
	// CreatedAt is the timestamp when the session was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is the timestamp when the session was last modified
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// Metadata contains additional session-specific data as key-value pairs
	Metadata map[string]any `json:"metadata" db:"metadata"`
}

// Validate checks if the QuerySession has valid field values.
// It ensures the locale is one of the supported values and required fields are populated.
func (s *QuerySession) Validate() error {
	if s.UserID == uuid.Nil {
		return errors.New("user_id is required and cannot be empty")
	}

	if s.TenantID == uuid.Nil {
		return errors.New("tenant_id is required and cannot be empty")
	}

	if !ValidLocales[s.Locale] {
		return fmt.Errorf("invalid locale '%s': must be one of 'en' or 'ar'", s.Locale)
	}

	return nil
}

// NewSession creates a new QuerySession with the provided parameters.
// It generates a new UUID, sets the timestamps, and initializes empty metadata.
func NewSession(userID, tenantID uuid.UUID, locale string) *QuerySession {
	now := time.Now()
	return &QuerySession{
		ID:        uuid.New(),
		UserID:    userID,
		TenantID:  tenantID,
		Locale:    locale,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]any),
	}
}

// ToJSON serializes the QuerySession to JSON format.
// Returns an error if serialization fails.
func (s *QuerySession) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON deserializes JSON data into the QuerySession.
// Returns an error if deserialization fails or the JSON is malformed.
func (s *QuerySession) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}
