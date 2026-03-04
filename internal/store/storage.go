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
	}
	User interface {
		Create(context.Context, *model.User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		User:  &UserStore{db: db},
	}
}
