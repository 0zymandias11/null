package main

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
)

type CreateUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(r, &payload); err != nil {
		app.writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	user := &store.User{
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "userID")
	if handle == "" {
		app.writeJSONError(w, http.StatusBadRequest, errors.New("userID is required"))
		return
	}
	var payload CreateUserPayload
	if err := readJSON(r, &payload); err != nil {
		app.writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	user := &store.User{
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
	}

	if _, err := app.store.Users.Put(r.Context(), user); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "userID")
	if handle == "" {
		app.writeJSONError(w, http.StatusBadRequest, errors.New("userID is required"))
		return
	}

	user, err := app.store.Users.Get(r.Context(), handle)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.writeJSONError(w, http.StatusNotFound, err)
			return
		}
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	handle, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.writeJSONError(w, http.StatusBadRequest, errors.New("userID is required"))
		return
	}

	var payload CreateUserPayload
	if err := readJSON(r, &payload); err != nil {
		app.writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	if payload.Username != "" {
		if err := app.store.Users.Follow(r.Context(), handle, payload.Username); err != nil {
			app.writeJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if err := writeJSON(w, http.StatusOK, map[string]string{"status": "followed"}); err != nil {
			app.writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}
}
