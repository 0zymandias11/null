package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockPostStore struct {
	CreateFn  func(ctx context.Context, post *store.Post) error
	GetByIDFn func(ctx context.Context, postID int64) (*store.Post, error)
	DeleteFn  func(ctx context.Context, postID int64) error
	PutFn     func(ctx context.Context, postID int64, post *store.Post) (*store.Post, error)
}

// Use the existing CreatePostPayload from posts.go

type UpdatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (m *mockPostStore) Create(ctx context.Context, post *store.Post) error {
	return m.CreateFn(ctx, post)
}
func (m *mockPostStore) GetByID(ctx context.Context, postID int64) (*store.Post, error) {
	return m.GetByIDFn(ctx, postID)
}
func (m *mockPostStore) Delete(ctx context.Context, postID int64) error {
	return m.DeleteFn(ctx, postID)
}
func (m *mockPostStore) Put(ctx context.Context, postID int64, post *store.Post) (*store.Post, error) {
	return m.PutFn(ctx, postID, post)
}

type mockCommentsStore struct{}

func (m *mockCommentsStore) GetPostById(ctx context.Context, postID int64) ([]*store.Comment, error) {
	return []*store.Comment{}, nil
}

func (m *mockCommentsStore) Create(ctx context.Context, comment *store.Comment) error {
	return nil
}

func TestCreatePostHandler(t *testing.T) {
	app := &application{
		store: store.Storage{
			Posts: &mockPostStore{
				CreateFn: func(ctx context.Context, post *store.Post) error {
					post.ID = 1
					return nil
				},
			},
		},
	}
	payload := CreatePostPayload{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  1,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	app.createPostHandler(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}
func TestGetPostHandler(t *testing.T) {
	app := &application{
		store: store.Storage{
			Posts: &mockPostStore{
				GetByIDFn: func(ctx context.Context, postID int64) (*store.Post, error) {
					return &store.Post{ID: postID, Title: "Test", Content: "Test Content", UserID: 1}, nil
				},
			},
			Comments: &mockCommentsStore{},
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/1", nil)
	recorder := httptest.NewRecorder()
	// Set chi URLParam for postID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	app.getPostHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
func TestDeletePostHandler(t *testing.T) {
	app := &application{
		store: store.Storage{
			Posts: &mockPostStore{
				DeleteFn: func(ctx context.Context, postID int64) error {
					return nil
				},
			},
		},
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/1", nil)
	recorder := httptest.NewRecorder()

	// Set chi URLParam for postID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	app.deletePostHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdatePostHandler(t *testing.T) {
	app := &application{
		store: store.Storage{
			Posts: &mockPostStore{
				PutFn: func(ctx context.Context, postID int64, post *store.Post) (*store.Post, error) {
					return &store.Post{ID: postID, Title: "Updated", Content: "Updated Content", UserID: 1}, nil
				},
			},
		},
	}

	payload := UpdatePostPayload{
		Title:   "Updated",
		Content: "Updated Content",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/posts/1", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	// Set chi URLParam for postID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("postID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	app.updatePostHandler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
