package main

import (
	"log"
	"net/http"

	"example.com/Go_Land/internal/env/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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
	addr string
	db   dbConfig
	//env     string
	version string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve Swagger UI at /swagger/*
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		w.Write([]byte(reqId))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Post("/posts", app.createPostHandler)
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Get("/{userID}", app.getUserHandler)
			r.Put("/{userID}", app.updateUserHandler)
			r.Post("/{userID}/follow", app.followUserHandler)
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		r.Route("/{postID}", func(r chi.Router) {
			r.Get("/", app.getPostHandler)
			r.Put("/", app.updatePostHandler)
			r.Delete("/", app.deletePostHandler)
			r.Get("/comments", app.getCommentsHandler)
			r.Post("/comments", app.createCommentsHandler)
		})
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
