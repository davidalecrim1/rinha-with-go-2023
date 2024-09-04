package config

import (
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/internal/repository"
	"os"
)

func InitializeServer() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	db := InitializeDatabase()
	defer db.Close()

	logger := NewLogger()
	repo := repository.NewPersonPostgreSqlRepository(db) // TODO: add logger
	service := domain.NewPersonService(repo)             // TODO: add logger
	handler := handler.NewPersonHandler(logger, service)

	err := InitializeRouter(handler, port, logger)

	if err != nil {
		logger.Error("error initializing server", "error:", err)
	}
}
