package main

import (
	"database/sql"
	"log"
	"time"

	"example.com/Go_Land/internal/env"
	"example.com/Go_Land/internal/env/store"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
	}

	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxIdleTime(duration)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	store := store.NewPostgresStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
