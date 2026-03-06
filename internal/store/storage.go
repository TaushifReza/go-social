package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/TaushifReza/go-social/internal/model"
)

var QueryTimeoutDuration = time.Second * 5

type Storage struct {
	Posts interface {
		Create(context.Context, *model.Posts) error
		GetByID(ctx context.Context, id int64) (*model.Posts, error)
		DeletePostByID(ctx context.Context, id int64) error
		Update(context.Context, *model.Posts) error
	}
	Users interface {
		Create(context.Context, *model.User) error
	}
	Comments interface {
		GetCommentByPostID(ctx context.Context, postID int64) ([]*model.Comment, error)
		Create(context.Context, *model.Comment) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
