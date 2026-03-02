package store

import (
	"context"
	"database/sql"

	"github.com/TaushifReza/go-social/internal/model"
	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, posts *model.Posts) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		posts.Content,
		posts.Title,
		posts.UserID,
		pq.Array(posts.Tags),
	).Scan(
		&posts.ID,
		&posts.CreatedAt,
		&posts.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
