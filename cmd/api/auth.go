package main

import (
	"log"
	"net/http"

	"example.com/Go_Land/internal/env/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type loginUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// registerUserHandler godoc
// @Summary      Register a new user
// @Description  Register a new user
// @Tags         auth
// @Accept      json
// @Produce     json
// @Param      payload body RegisterUserPayload true "User credential"
// @Success     201   {object}  User
// @Failure     400   {object}  ErrorResponse
// @Failure     500   {object}  ErrorResponse
// @Router      /auth/user [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload
	if err := readJSON(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := app.store.Users.GetUserCred(r.Context(), payload.Username)
	if err != nil {
		app.notFound(w, r, err)
		return
	}

	log.Printf("User %s is trying to login with hashed password %+v with input payload of: %s", user.Username, user.Password, payload.Password)

	if err := user.Password.CompareHash(payload.Password); err != nil {
		app.unauthorized(w, r, err)
		return
	}

	// Generate JWT token for the user
	token, err := app.JwtHandler(user.Username)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send the token back to the user
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}
