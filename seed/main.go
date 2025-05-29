package main

import (
	"log"

	"example.com/Go_Land/internal/env"
	"example.com/Go_Land/internal/env/db"
	"example.com/Go_Land/internal/env/store"
)

func main() {
	addr := env.GetString("DB_DSN", "")
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
