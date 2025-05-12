package main

import (
	"log"
	"net/http"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		w.Write([]byte(reqId))
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Post("/posts", app.createPostHandler)
		r.Post("/users", app.createUserHandler)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	log.Printf("Server has started running at Port: %s", app.config.addr)

	return srv.ListenAndServe()
}
