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
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		Delete(ctx context.Context, postID int64) error
		Put(ctx context.Context, postID int64, post *Post) (*Post, error)
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostsWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		Put(ctx context.Context, user *User) (*User, error)
		Get(ctx context.Context, handle string) (*User, error)
		Follow(ctx context.Context, userID int64, followerUsername string) error
		// Unfollow(ctx context.Context, userID, followerUsername string) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
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
