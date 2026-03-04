package store

import (
	"context"
	"database/sql"

	"github.com/TaushifReza/go-social/internal/model"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetCommentByPostID(ctx context.Context, postID int64) ([]*model.Comment, error) {
	query := `
	SELECT *
	FROM comments
	WHERE post_id = $1
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return comments, nil
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
