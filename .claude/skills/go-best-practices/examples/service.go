// Package service demonstrates idiomatic Go patterns for a service layer.
// This is a complete, runnable example of best practices.
package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// DOMAIN TYPES
// =============================================================================

// User represents a user in the system.
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserRequest contains data needed to create a user.
type CreateUserRequest struct {
	Email string
	Name  string
}

// UpdateUserRequest contains data needed to update a user.
type UpdateUserRequest struct {
	Email *string
	Name  *string
}

// =============================================================================
// ERRORS
// =============================================================================

// Sentinel errors for specific conditions.
var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrInvalidEmail      = fmt.Errorf("invalid email")
)

// ValidationError represents a field-level validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// =============================================================================
// REPOSITORY INTERFACE
// =============================================================================

// Repository defines the data access interface.
// Keep interfaces small and focused.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// =============================================================================
// SERVICE IMPLEMENTATION
// =============================================================================

// Config holds service configuration.
type Config struct {
	MaxListLimit int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxListLimit: 100,
	}
}

// Option is a functional option for configuring the service.
type Option func(*UserService)

// WithConfig sets a custom configuration.
func WithConfig(cfg Config) Option {
	return func(s *UserService) {
		s.config = cfg
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(s *UserService) {
		s.logger = logger
	}
}

// UserService provides business logic for user operations.
type UserService struct {
	repo   Repository
	config Config
	logger *slog.Logger
}

// NewUserService creates a new UserService with the given repository and options.
func NewUserService(repo Repository, opts ...Option) *UserService {
	s := &UserService{
		repo:   repo,
		config: DefaultConfig(),
		logger: slog.Default(),
	}

	// Apply functional options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// Validate input
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !IsNotFoundError(err) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// Create user
	now := time.Now()
	user := &User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	s.logger.Info("user created", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// Get retrieves a user by ID.
func (s *UserService) Get(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, &ValidationError{Field: "id", Message: "cannot be empty"}
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}

	return user, nil
}

// Update modifies an existing user.
func (s *UserService) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	// Get existing user
	user, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates (only update fields that are provided)
	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return nil, err
		}
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user %s: %w", id, err)
	}

	s.logger.Info("user updated", "user_id", user.ID)
	return user, nil
}

// Delete removes a user.
func (s *UserService) Delete(ctx context.Context, id string) error {
	// Verify user exists first
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user %s: %w", id, err)
	}

	s.logger.Info("user deleted", "user_id", id)
	return nil
}

// List retrieves users with pagination.
func (s *UserService) List(ctx context.Context, limit, offset int) ([]*User, error) {
	// Apply constraints
	if limit <= 0 {
		limit = s.config.MaxListLimit
	}
	if limit > s.config.MaxListLimit {
		limit = s.config.MaxListLimit
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

// =============================================================================
// HELPERS
// =============================================================================

// validateEmail checks if an email is valid.
func validateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "cannot be empty"}
	}
	// Add more validation as needed
	if len(email) < 5 || len(email) > 255 {
		return &ValidationError{Field: "email", Message: "must be between 5 and 255 characters"}
	}
	return nil
}

// IsNotFoundError checks if an error indicates not found.
func IsNotFoundError(err error) bool {
	return err == ErrUserNotFound ||
		err.Error() == "not found" ||
		err.Error() == "sql: no rows in result set"
}
