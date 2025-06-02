package main

import (
	"net/http"
	"strconv"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
)

type CreateCommentPayload struct {
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (app *application) createCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  payload.PostID,
		Content: payload.Content,
		UserID:  payload.UserID,
		// Assuming CreatedAt and UpdatedAt are handled by the store
	}
	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	comments, err := app.store.Comments.GetPostById(r.Context(), postID)
	if err != nil {
		app.notFound(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, comments); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
