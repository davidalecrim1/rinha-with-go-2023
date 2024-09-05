package config

import (
	"database/sql"
	"fmt"
	"go-rinha-de-backend-2023/config/env"
	"log"

	_ "github.com/lib/pq"
)

func InitializeDatabase() *sql.DB {
	db_host := env.GetEnvOrSetDefault("DB_HOST", "localhost")
	db_conn_string := fmt.Sprintf("postgres://admin:password@%s:5432/persons?sslmode=disable", db_host)
	db, err := sql.Open("postgres", db_conn_string)

	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)

	return db
}
