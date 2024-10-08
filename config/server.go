package config

import (
	"go-rinha-de-backend-2023/config/env"
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/internal/repository"
)

func InitializeServer() {
	db := InitializeDatabase()
	defer db.Close()

	logger := NewLogger()
	repo := repository.NewPersonPostgreSqlRepository(db) // TODO: add logger
	service := domain.NewPersonService(repo)             // TODO: add logger
	handler := handler.NewPersonHandler(logger, service)

	port := env.GetEnvOrSetDefault("PORT", "8080")
	err := InitializeRouter(handler, port, logger)

	if err != nil {
		logger.Error("error initializing server", "error:", err)
	}
}
