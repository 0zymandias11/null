package main

import (
	"net/http"
	"strconv"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(r, &payload); err != nil {
		app.writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  payload.UserID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Assuming postID is passed as a query parameter
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}
