package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Open(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	createSchema(db)
	return db
}

func createSchema(db *sql.DB) {
	const schema = `
CREATE TABLE IF NOT EXISTS videos (
	id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	filename TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);
`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
}
