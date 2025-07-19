package main

import (
	"errors"
	"time"

	"example.com/Go_Land/internal/env/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) tokenSpawn(handle string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": handle,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (app *application) JwtHandler(handle string) (string, error) {
	if handle == "" {
		return "", errors.New("userID is required")
	}

	tokenString, err := app.tokenSpawn(handle)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", err
		}
		app.logger.Errorw("failed to generate token",
			"handle", handle,
		)
		return "", err
	}
	return tokenString, nil
}
