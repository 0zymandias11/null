package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	text *string
	hash []byte
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

func (p *password) CompareHash(text string) error {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(text)); err != nil {
		return err
	}
	return nil
}

type UserStore struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func (s *UserStore) Get(ctx context.Context, handle string) (*User, error) {
	query := `Select users.Email, users.ID, users.Username from users where username = $1 or email = $1`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, handle).Scan(
		&user.Email,
		&user.ID,
		&user.Username,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) GetUserCred(ctx context.Context, handle string) (*User, error) {
	s.logger.Logln(zap.DebugLevel, "Getting User Creds", handle)
	query := `SELECT users.Username, users.Password FROM users WHERE username = $1`
	user := &User{}
	var hash []byte
	err := s.db.QueryRowContext(ctx, query, handle).Scan(
		&user.Username,
		&hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	user.Password.hash = hash
	return user, nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	s.logger.Logln(zap.DebugLevel, "Creating User", "user:", user.Username, "email", user.Email, "Password: ", user.Password.text)
	query := `INSERT INTO users (email, password, username)
              VALUES ($1, $2, $3)
              ON CONFLICT (email) DO NOTHING
              RETURNING id, created_at, updated_at`

	// Try to insert, but if the user already exists, just return nil (no error)
	err := tx.QueryRowContext(ctx,
		query,
		user.Email,
		user.Password.hash,
		user.Username,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		// User already exists, treat as success for seeding
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Put(ctx context.Context, tx *sql.Tx, user *User) (*User, error) {
	query := `UPDATE users 
			  SET email = $1, password = $2, username = $3, updated_at = NOW() 
			  WHERE id = $4`

	_, err := tx.ExecContext(ctx,
		query,
		user.Email,
		user.Password,
		user.Username,
		user.ID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) Follow(ctx context.Context, tx *sql.Tx, userID int64, followerUsername string) error {
	query := `INSERT INTO followers (user_id, follower_username) 
			  VALUES ($1, $2) 
			  ON CONFLICT (user_id, follower_username) DO NOTHING`

	_, err := tx.ExecContext(ctx, query, userID, followerUsername)
	if err != nil {
		return err
	}
	return nil
}
