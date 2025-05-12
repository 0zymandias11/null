package store

import (
	"context"
	"database/sql"
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
