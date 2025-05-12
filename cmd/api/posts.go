package main

import (
	"net/http"

	"example.com/Go_Land/internal/env/store"
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
