package config

import (
	"context"
	"go-rinha-de-backend-2023/config/env"
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/internal/repository"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitializeServer() {
	db := InitializeDatabase()
	defer db.Close()

	cache := InitializeCache()
	defer cache.Close()

	logger := NewLogger()

	ctx, cancel := context.WithCancel(context.Background())
	worker := repository.NewPersonAsyncRepository(logger, ctx, db)
	go worker.Start()

	cacheRepo := repository.NewPersonCacheRepository(logger, cache)
	repo := repository.NewPersonRepository(logger, db, cacheRepo, worker.NewCreatePersonChannel())
	service := domain.NewPersonService(logger, repo)
	handler := handler.NewPersonHandler(logger, service)

	app := fiber.New()

	InitializeRouter(app, handler, logger)

	go func() {
		err := app.Listen(":" + env.GetEnvOrSetDefault("PORT", "8080"))

		if err != nil {
			log.Fatalf("error server configuration: %v", err)
		}
	}()

	GracefulShutdown(logger, cancel, app)
}

func GracefulShutdown(logger *slog.Logger, cancel context.CancelFunc, app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server...")
	cancel()

	time.Sleep(10 * time.Second)

	if err := app.Shutdown(); err != nil {
		logger.Error("error shutting down server properly", "error", err)
	}

	logger.Info("server gracefully stopped!")
}
