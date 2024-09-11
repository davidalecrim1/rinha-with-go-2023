package config

import (
	"context"
	"fmt"
	"go-rinha-de-backend-2023/config/env"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeDatabase() *pgxpool.Pool {
	db_host := env.GetEnvOrSetDefault("DB_HOST", "localhost")
	db_user := env.GetEnvOrSetDefault("DB_USER", "admin")
	db_password := env.GetEnvOrSetDefault("DB_PASSWORD", "password")
	db_port := env.GetEnvOrSetDefault("DB_PORT", "5432")
	db_schema := env.GetEnvOrSetDefault("DB_SCHEMA", "people")
	db_conn_string := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_schema)

	max_conn, err := strconv.Atoi(env.GetEnvOrSetDefault("DB_MAX_CONN", "50"))

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	config, err := pgxpool.ParseConfig(db_conn_string)

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	config.MaxConns = int32(max_conn)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("error loading database configuration: %v", err)
	}

	return pool
}
