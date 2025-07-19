package main

import (
	"errors"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("internal error %s path: %s  error %s", r.Method, r.URL.Path, err)
	app.logger.Errorw("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusInternalServerError, errors.New("the server encountered a problem"))

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("bad request %s path: %s  error %s", r.Method, r.URL.Path, err)
	app.logger.Errorw("bad request",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusBadRequest, errors.New(err.Error()))

}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("not found %s path: %s  error %s", r.Method, r.URL.Path, err)
	app.logger.Errorw("not found",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusNotFound, err)

}

func (app *application) unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized",
		"method", r.Method,
		"path", r.URL.Path,
	)
	app.writeJSONError(w, http.StatusUnauthorized, errors.New("unauthorized access"))
}
