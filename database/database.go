package database

import (
	"database/sql"
	"log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	return db
}