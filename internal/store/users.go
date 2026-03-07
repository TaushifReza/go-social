package store

import (
	"context"
	"database/sql"

	"github.com/TaushifReza/go-social/internal/dto"
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

func (s *UserStore) GetUserbyID(ctx context.Context, id int64) (*dto.UserResponseDto, error) {

	query := `
    	SELECT id, username, email, created_at
    	FROM users
    	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var user dto.UserResponseDto
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Follow(ctx context.Context, userID int64, followUserID int64) error {
	query := `
    INSERT INTO followers (user_id, follower_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, follower_id) DO NOTHING
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followUserID, userID)
	return err
}

func (s *UserStore) UnFollow(ctx context.Context, userID int64, followUserID int64) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followUserID, userID)
	return err
}
