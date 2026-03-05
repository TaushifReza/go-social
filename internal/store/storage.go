package store

import (
	"context"
	"database/sql"

	"github.com/TaushifReza/go-social/internal/model"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *model.Posts) error
		GetByID(ctx context.Context, id int64) (*model.Posts, error)
		DeletePostByID(ctx context.Context, id int64) error
	}
	Users interface {
		Create(context.Context, *model.User) error
	}
	Comments interface {
		GetCommentByPostID(ctx context.Context, postID int64) ([]*model.Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
