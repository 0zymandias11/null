package main

import (
	"net/http"

	"example.com/Go_Land/internal/env/store"
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
