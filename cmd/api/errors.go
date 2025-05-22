package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error %s path: %s  error %s", r.Method, r.URL.Path, err)

	app.writeJSONError(w, http.StatusInternalServerError, errors.New("the server encountered a problem"))

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request %s path: %s  error %s", r.Method, r.URL.Path, err)

	app.writeJSONError(w, http.StatusBadRequest, errors.New(err.Error()))

}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found %s path: %s  error %s", r.Method, r.URL.Path, err)

	app.writeJSONError(w, http.StatusNotFound, errors.New("the requested resource was not found"))

}
