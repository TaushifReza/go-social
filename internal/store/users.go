package store

import (
	"context"
	"database/sql"

	"github.com/TaushifReza/go-social/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	query := `
    INSERT INTO users (username, password, email)
    VALUES ($1, $2, $3) RETURNING id, created_at
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.UserName,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
