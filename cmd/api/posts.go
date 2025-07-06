package main

import (
	"net/http"
	"strconv"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
}

// CreatePostHandler godoc
// @Summary      Create a new post
// @Description  Create a new post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post  body      store.Post  true  "Post data"
// @Success      201   {object}  store.Post
// @Router       /api/v1/posts [post]

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  payload.UserID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetPostHandler godoc
// @Summary      Get a post
// @Description  Get a post by ID
// @Tags         posts
// @Produce      json
// @Param        postID  path      int  true  "Post ID"
// @Success      200     {object}  store.Post
// @Router       /api/v1/{postID}/ [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Assuming postID is passed as a query parameter
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		app.notFound(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetPostById(ctx, postID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePostHandler godoc
// @Summary      Delete a post
// @Description  Delete a post by ID
// @Tags         posts
// @Param        postID  path  int  true  "Post ID"
// @Success      204
// @Router       /api/v1/{postID}/ [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.store.Posts.Delete(ctx, postID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdatePostHandler godoc
// @Summary      Update a post
// @Description  Update a post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID  path      int   true  "Post ID"
// @Param        post    body      store.Post  true  "Post data"
// @Success      200     {object}  store.Post
// @Router       /api/v1/{postID}/ [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post := &store.Post{}
	post, err = app.store.Posts.Put(ctx, postID, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
