package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	Likes     int64      `json:"likes"`
	Dislikes  int64      `json:"dislikes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Comments  []*Comment `json:"comments"`
	Version   int        `json:"version"`
	User      User       `json:"user"`
}

type PostsWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}


func (s *PostStore) Create(ctx context.Context, tx *sql.Tx, post *Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags, likes, dislikes)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
		post.Likes,
		post.Dislikes,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err == sql.ErrNoRows {
		s.logger.Infow("Post already exists", "title", post.Title, "user_id", post.UserID)
		return nil
	}
	if err != nil {
		s.logger.Errorw("Failed to create post", "error", err, "title", post.Title, "user_id", post.UserID)
		return err
	}
	s.logger.Infow("Post created", "id", post.ID, "title", post.Title, "user_id", post.UserID)
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := "SELECT * from posts where id =$1"
	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Dislikes, &post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) Delete(ctx context.Context, tx *sql.Tx, postID int64) error {
	query := "Delete from posts where id = $1"
	res, err := tx.ExecContext(ctx, query, postID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Put(ctx context.Context, tx *sql.Tx, postID int64, post *Post) (*Post, error) {
	query := `UPDATE posts SET title = $1, content = $2, user_id = $3, tags = $4, updated_at = NOW() WHERE id = $5 AND version = $6 RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
		postID,
		post.Version,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostsWithMetadata, error) {
	query := `p.id, p.user_id, p.title, p.content, p.tags, p.likes, p.dislikes, p.created_at, p.updated_at, p.version, count(c.id) 
				as comment_count from posts p left join comments c on c.post_id = p.id left join 
				users u on u.id = p.user_id 
				join
				followers f 
				ON f.follower_id = p.user_id OR p.user_id = $1
				where f.user_id = $1 or p.user_id = $1
				and p.deleted_at is null and p.version > 0 
				group by p.id, u.username 
				order by p.created_at` + fq.Sort +
		`LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostsWithMetadata
	for rows.Next() {
		var post PostsWithMetadata
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags), &post.Likes, &post.Dislikes, &post.CreatedAt, &post.UpdatedAt, &post.Version, &post.CommentCount); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
