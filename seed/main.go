package main

import (
	"log"
	"path/filepath"

	"example.com/Go_Land/internal/env"
	"example.com/Go_Land/internal/env/db"
	"example.com/Go_Land/internal/env/store"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	addr := env.GetString("DB_DSN", "")
	log.Println("Using database connection string:", addr)
	// Check if the DB_DSN environment variable is set
	if addr == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	store := store.NewPostgresStorage(conn)
	db.Seed(store)
	log.Println("Database seeded successfully")
}
