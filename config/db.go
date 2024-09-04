package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitializeDatabase() *sql.DB {
	var err error
	db, err := sql.Open("postgres", "postgres://admin:password@postgres-db:5432/persons?sslmode=disable")

	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	return db
}
