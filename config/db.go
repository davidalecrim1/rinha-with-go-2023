package config

import (
	"database/sql"
	"fmt"
	"go-rinha-de-backend-2023/config/env"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

func InitializeDatabase() *sql.DB {
	db_host := env.GetEnvOrSetDefault("DB_HOST", "localhost")
	db_user := env.GetEnvOrSetDefault("DB_USER", "admin")
	db_password := env.GetEnvOrSetDefault("DB_PASSWORD", "password")
	db_port := env.GetEnvOrSetDefault("DB_PORT", "5432")
	db_schema := env.GetEnvOrSetDefault("DB_SCHEMA", "people")
	db_conn_string := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_schema)
	db, err := sql.Open("postgres", db_conn_string)

	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	max_conn, err := strconv.Atoi(env.GetEnvOrSetDefault("DB_MAX_CONN", "50"))

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	max_idle_conn, err := strconv.Atoi(env.GetEnvOrSetDefault("DB_MAX_IDLE_CONN", "25"))

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	db.SetMaxOpenConns(max_conn)
	db.SetMaxIdleConns(max_idle_conn)

	return db
}
