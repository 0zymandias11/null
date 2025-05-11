package main

import (
	"log"
	"os"
	"path/filepath"

	"example.com/Go_Land/internal/env"
	"example.com/Go_Land/internal/env/db"
	"example.com/Go_Land/internal/env/store"
	"github.com/joho/godotenv"
)

func main() {
	// Get the executable's directory
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	exeDir := filepath.Dir(exe)

	// Look for .env in the project root (two levels up from executable)
	envPath := filepath.Join(exeDir, "..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	dsn := env.GetString("DB_DSN", "")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}
	log.Printf("Using database connection string: %s", dsn)

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn:          dsn,
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
	}

	db, err := db.New(cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Printf("Starting server on %s", cfg.addr)
	if err := app.run(mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
