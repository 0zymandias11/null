// Package main My Awesome API.
//
//     Schemes: http, https
//     Host:     localhost:8080
//     BasePath: /api/v1
//     Version:  1.0
//     Title:    Go Land API
//     Description: This is a sample server for Go Land.
//
// swagger:meta

package main

import (
	"log"
	"os"
	"path/filepath"

	_ "example.com/Go_Land/docs"
	"example.com/Go_Land/internal/env"
	"example.com/Go_Land/internal/env/db"
	"example.com/Go_Land/internal/env/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	log.Printf("Working directory: %s", workingDir)
	//Logger
	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("Failed to initialize zap-logger %v", err)
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	// Navigate to the .env file based on the working directory
	envPath := filepath.Join(workingDir, "..", "..", ".env") // Adjust based on your directory structure
	log.Printf("Looking for .env at: %s", envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	dsn := env.GetString("DB_DSN", "")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}
	log.Printf("Using database connection string: %s", dsn)

	secret := env.GetString("JWT_SECRET", "GET RAILED DICK_BRAIN")
	cfg := config{
		addr:    env.GetString("ADDR", ":8080"),
		version: env.GetString("VERSION", "1.0.0"),
		db: dbConfig{
			dsn:          dsn,
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
		jwtSecret: secret,
	}

	defer logger.Sync()

	db, err := db.New(cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger.Info("Sugar Logs: database has connected")

	log.Default().Printf("Connected to database: %s", cfg.db.dsn)
	store := store.NewPostgresStorage(db, sugar)

	app := &application{
		config: cfg,
		store:  store,
		logger: sugar,
		dbConnector:     db,
	}

	mux := app.mount()

	log.Default().Printf("Using JWT SECRET: %s", app.config.jwtSecret)
	// http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Printf("Starting server on %s", cfg.addr)
	if err := app.run(mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
