package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db}
}

func (s *UserStore) Get(ctx context.Context, handle string) (*User, error) {
	query := `Select * from users where username = $1 or email = $1`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, handle).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password, username) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx,
		query,
		user.Email,
		user.Password,
		user.Username,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Put(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE users 
			  SET email = $1, password = $2, username = $3, updated_at = NOW() 
			  WHERE id = $4`

	_, err := s.db.ExecContext(ctx,
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
