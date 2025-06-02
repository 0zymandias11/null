package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		Delete(ctx context.Context, postID int64) error
		Put(ctx context.Context, postID int64, post *Post) (*Post, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		Put(ctx context.Context, user *User) (*User, error)
		Get(ctx context.Context, handle string) (*User, error)
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetPostById(ctx context.Context, postID int64) ([]*Comment, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentsStore{db},
	}
}
