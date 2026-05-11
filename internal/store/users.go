package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
	"github.com/lib/pq"
)

var (
	ErrDuplicateEmail    = errors.New("a user with this email already exists")
	ErrDuplicateUsername = errors.New("a user with this username already exists")
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
    INSERT INTO users (username, password, email)
    VALUES ($1, $2, $3) RETURNING id, created_at
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := tx.QueryRowContext(
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
		if pqErr, ok := err.(*pq.Error); ok {
			// 23505 is the code for unique_violation
			if pqErr.Code == "23505" {
				switch pqErr.Constraint {
				case "users_email_key":
					return ErrDuplicateEmail
				case "users_username_key":
					return ErrDuplicateUsername
				}
			}
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("User not found")
		}
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *model.User, token string, invitationExp time.Duration) error {
	// transaction wrapper
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		// create the user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().UTC().Add(exp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// 1. Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case something fails before we commit
	defer tx.Rollback()

	// 2. Get the invitation (pass 'tx' if your query supports it, or use db)
	invite, err := s.getUserInvitations(ctx, token)
	if err != nil {
		return err
	}

	// 3. Validate expiry
	if time.Now().UTC().After(invite.Expire.UTC()) {
		return errors.New("invitation has expired")
	}

	// 4. Update user to active
	if err := s.updateIsActive(ctx, tx, invite.UserID); err != nil {
		return err
	}

	// 5. Delete the token
	if err := s.deleteInvitation(ctx, tx, token); err != nil {
		return err
	}

	// 6. Commit the transaction
	return tx.Commit()
}

func (s *UserStore) getUserInvitations(ctx context.Context, token string) (*model.UserInvitations, error) {
	query := `
	SELECT *
	FROM user_invitations
	WHERE token = ($1)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user_invitations model.UserInvitations
	if err := s.db.QueryRowContext(
		ctx,
		query,
		token,
	).Scan(
		&user_invitations.Token,
		&user_invitations.UserID,
		&user_invitations.Expire,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, errors.New("invitation not found or invalid")
		}
		return nil, err
	}

	return &user_invitations, nil
}

// Change s.db.ExecContext to tx.ExecContext
func (s *UserStore) updateIsActive(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `UPDATE users SET is_active = true WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, token string) error {
	query := `DELETE FROM user_invitations WHERE token = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	return nil
}
