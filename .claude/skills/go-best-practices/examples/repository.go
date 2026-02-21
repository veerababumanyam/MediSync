// Package repository demonstrates the repository pattern in Go.
// This example uses PostgreSQL with pgx for database operations.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents the database model for a user.
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostgresRepository implements the repository interface using PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create inserts a new user into the database.
func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if isUniqueViolation(err) {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID.
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user, err := r.scanUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user, err := r.scanUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil without error for not found
		}
		return nil, fmt.Errorf("query user by email: %w", err)
	}

	return user, nil
}

// Update modifies an existing user.
func (r *PostgresRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %s", user.ID)
	}

	return nil
}

// Delete removes a user from the database.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}

// List retrieves users with pagination.
func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := r.scanUserFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

// scanUser scans a single user from a query row.
func (r *PostgresRepository) scanUser(row pgx.Row) (*User, error) {
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// scanUserFromRows scans a user from rows iterator.
func (r *PostgresRepository) scanUserFromRows(rows pgx.Rows) (*User, error) {
	user := &User{}
	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// isUniqueViolation checks if the error is a unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgx.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

// =============================================================================
// TRANSACTION SUPPORT
// =============================================================================

// TxRepository provides transaction support.
type TxRepository struct {
	tx pgx.Tx
}

// BeginTx starts a new transaction.
func (r *PostgresRepository) BeginTx(ctx context.Context) (*TxRepository, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return &TxRepository{tx: tx}, nil
}

// Commit commits the transaction.
func (r *TxRepository) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

// Rollback rolls back the transaction.
// Always safe to call - will be no-op if already committed.
func (r *TxRepository) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}

// CreateWithinTx creates a user within an existing transaction.
func (r *TxRepository) CreateWithinTx(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.tx.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert user in tx: %w", err)
	}

	return nil
}

// Example of transaction usage:
//
// func (s *Service) CreateWithAudit(ctx context.Context, user *User) error {
//     txRepo, err := s.repo.BeginTx(ctx)
//     if err != nil {
//         return err
//     }
//     defer txRepo.Rollback(ctx) // Safe rollback if not committed
//
//     if err := txRepo.CreateWithinTx(ctx, user); err != nil {
//         return err
//     }
//
//     if err := txRepo.CreateAuditLog(ctx, "user_created", user.ID); err != nil {
//         return err
//     }
//
//     return txRepo.Commit(ctx)
// }
