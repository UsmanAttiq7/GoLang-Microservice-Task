package store

import (
	"context"
	"fmt"

	"github.com/golang_falcon_task/user-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGUserStore struct {
	db *pgxpool.Pool
}

// NewPGUserStore creates a new PGUserStore instance.
func NewPGUserStore(db *pgxpool.Pool) *PGUserStore {
	return &PGUserStore{db: db}
}

// GetUser retrieves a user by ID.
func (s *PGUserStore) GetUser(ctx context.Context, id int32) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow(ctx, `SELECT user_id, name FROM users WHERE user_id = $1`, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}
	return &user, nil
}

// CreateUser creates a new user and returns their ID.
func (s *PGUserStore) CreateUser(ctx context.Context, name string) (int32, error) {
	var userID int32
	err := s.db.QueryRow(ctx, `INSERT INTO users (name) VALUES ($1) RETURNING user_id`, name).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}
	return userID, nil
}

// DeleteUser deletes a user by ID.
func (s *PGUserStore) DeleteUser(ctx context.Context, id int32) error {
	result, err := s.db.Exec(ctx, `DELETE FROM users WHERE user_id = $1`, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	// Check the number of rows affected
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
