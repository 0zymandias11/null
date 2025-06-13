package main

import (
	"log"
	"net/http"
	"os"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	data := map[string]string{
		"status":  "ok",
		"env":     os.Getenv("ENV"),
		"version": app.config.version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		log.Fatal("Failed to write JSON response:", err)
		// Handle the error as needed, e.g., log it or return an error response
	}
}

func (app *application) writeJSONError(w http.ResponseWriter, status int, err error) {
	data := map[string]string{
		"error": err.Error(),
	}
	if err := writeJSON(w, status, data); err != nil {
		log.Println("Failed to write JSON error response:", err)
	}
}
