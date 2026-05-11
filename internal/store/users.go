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
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. Find the invitation
		invite, err := s.getUserInvitations(ctx, tx, token) // Note: pass tx here
		if err != nil {
			return err
		}

		// 2. Validate expiry
		if time.Now().UTC().After(invite.Expire.UTC()) {
			return errors.New("invitation has expired")
		}

		// 3. Update user status
		if err := s.updateIsActive(ctx, tx, invite.UserID); err != nil {
			return err
		}

		// 4. Delete the token
		if err := s.deleteInvitation(ctx, tx, token); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) getUserInvitations(ctx context.Context, tx *sql.Tx, token string) (*model.UserInvitations, error) {
	// 1. Explicitly name columns to match your struct scan order
	query := `
        SELECT token, user_id, expiry
        FROM user_invitations
        WHERE token = $1
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var invitation model.UserInvitations

	// 2. Use 'tx' instead of 's.db' so it participates in the transaction
	err := tx.QueryRowContext(ctx, query, token).Scan(
		&invitation.Token,
		&invitation.UserID,
		&invitation.Expire,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invitation not found or invalid")
		}
		return nil, err
	}

	return &invitation, nil
}

func (s *UserStore) updateIsActive(ctx context.Context, db DBQueryer, userID int64) error {
	// Now this works whether you pass a transaction or the standard DB pool
	_, err := db.ExecContext(ctx, "UPDATE users SET is_active = true WHERE id = $1", userID)
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
