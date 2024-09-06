package config

import (
	"go-rinha-de-backend-2023/config/env"
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func InitializeServer() {
	db := InitializeDatabase()
	defer db.Close()

	logger := NewLogger()
	repo := repository.NewPersonPostgreSqlRepository(db) // TODO: add logger
	service := domain.NewPersonService(repo)             // TODO: add logger
	handler := handler.NewPersonHandler(logger, service)

	app := fiber.New()
	InitializeRouter(app, handler, logger)

	err := app.Listen(":" + env.GetEnvOrSetDefault("PORT", "8080"))

	if err != nil {
		logger.Error("error initializing server", "error:", err)
		return
	}
}
