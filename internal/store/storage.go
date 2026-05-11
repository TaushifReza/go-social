package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
)

var QueryTimeoutDuration = time.Second * 5

type Storage struct {
	Posts interface {
		Create(context.Context, *model.Posts) error
		GetByID(ctx context.Context, id int64) (*model.Posts, error)
		DeletePostByID(ctx context.Context, id int64) error
		Update(context.Context, *model.Posts) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]dto.PostWithMetaData, error)
	}
	Users interface {
		CreateAndInvite(ctx context.Context, user *model.User, token string, invitationExp time.Duration) error
		Create(ctx context.Context, tx *sql.Tx, user *model.User) error
		GetUserbyID(ctx context.Context, id int64) (*dto.UserResponseDto, error)
		Follow(ctx context.Context, userID int64, followUserID int64) error
		UnFollow(ctx context.Context, userID int64, followUserID int64) error
		Activate(context.Context, string) error
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

// This interface works for both *sql.DB and *sql.Tx
type DBQueryer interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// This ensures we rollback if the function panics
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rolling back
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback() // Rollback on explicit error
		return err
	}

	return tx.Commit()
}
