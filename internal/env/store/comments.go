package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type CommentsStore struct {
	db *sql.DB
}

type Comment struct{
	ID int64 `json:"id"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *CommentsStore) GetPostById(ctx context.Context, postID int64) ([]*Comment, error){
	query := "Select * from comments Join users on users.id = comments.user_id where comments.post_id = $1 order by comments.created_at desc"
	comments := []*Comment{}
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next(){
		comment := &Comment{}
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
