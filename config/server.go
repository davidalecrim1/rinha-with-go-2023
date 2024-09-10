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
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	app := fiber.New(fiber.Config{
		AppName:     "rinha-go-app by @davidalecrim1",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	InitializeRouter(app, handler, logger)

	go func() {
		err := app.Listen(":" + env.GetEnvOrSetDefault("PORT", "8080"))

		if err != nil {
			log.Fatalf("error server configuration: %v", err)
		}
	}()

	WarmUpUuid()
	PerformProfiling(logger)
	GracefulShutdown(logger, cancel, app)
}

func WarmUpUuid() {
	uuid.EnableRandPool()
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

func PerformProfiling(logger *slog.Logger) {
	if os.Getenv("ENABLE_PROFILING") != "1" {
		logger.Info("profiling is disabled")
		return
	}

	err := os.MkdirAll("prof", 0755)
	if err != nil {
		logger.Error("failed to create profiling directory", "error", err)
	}

	logger.Info("profiling is enabled")

	cf, err := os.Create(os.Getenv("CPU_PROF_OUT"))
	if err != nil {
		logger.Error("failed to start CPU profiling", "error", err)
	}
	pprof.StartCPUProfile(cf)

	mf, err := os.Create(os.Getenv("MEMORY_PROF_OUT"))
	if err != nil {
		logger.Error("failed to start memory profiling", "error", err)
	}
	pprof.WriteHeapProfile(mf)

	tc, err := os.Create(os.Getenv("TRACE_PROF_OUT"))
	if err != nil {
		logger.Error("failed to start memory profiling", "error", err)
	}
	trace.Start(tc)

	stop := time.After(time.Minute * 4)

	go func() {
		<-stop
		pprof.StopCPUProfile()
		trace.Stop()
		cf.Close()
		mf.Close()
		tc.Close()
		logger.Info("CPU profiling stopped")
	}()
}
