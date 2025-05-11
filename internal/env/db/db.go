package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// Ensure SSL is disabled in the connection string
	if addr == "" {
		return nil, sql.ErrConnDone
	}

	// Only append sslmode if it's not already in the connection string
	if !strings.Contains(addr, "sslmode=") {
		if addr[len(addr)-1] != '?' {
			addr += "?"
		}
		addr += "sslmode=disable"
	}

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
