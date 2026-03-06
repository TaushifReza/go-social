package store

import (
	"context"
	"database/sql"
	"log"

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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
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

func (s *PostStore) GetByID(ctx context.Context, id int64) (*model.Posts, error) {
	query := `
	SELECT id, content, title, user_id, tags, created_at, updated_at
	FROM posts
	WHERE ID = ($1)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var post model.Posts
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostStore) DeletePostByID(ctx context.Context, id int64) error {
	query := `
	DELETE FROM posts
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// execute the DELETE query with context
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// optionally, check number of row affected
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Deleted %v rows", rowAffected)
	if rowAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *model.Posts) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		return err
	}

	return nil
}
