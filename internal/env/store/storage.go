package store

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, tx *sql.Tx, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		Delete(ctx context.Context, tx *sql.Tx, postID int64) error
		Put(ctx context.Context, tx *sql.Tx, postID int64, post *Post) (*Post, error)
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostsWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		Put(ctx context.Context, tx *sql.Tx, user *User) (*User, error)
		Get(ctx context.Context, handle string) (*User, error)
		Follow(ctx context.Context, tx *sql.Tx, userID int64, followerUsername string) error
		GetUserCred(ctx context.Context, handle string) (*User, error)
		// Unfollow(ctx context.Context, userID, followerUsername string) error
	}
	Comments interface {
		Create(ctx context.Context, tx *sql.Tx, comment *Comment) error
		GetPostById(ctx context.Context, postID int64) ([]*Comment, error)
	}
}

func NewPostgresStorage(db *sql.DB, logger *zap.SugaredLogger) Storage {
	return Storage{
		Posts:    &PostStore{db, logger},
		Users:    &UserStore{db, logger},
		Comments: &CommentsStore{db, logger},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	return fn(tx)
}
