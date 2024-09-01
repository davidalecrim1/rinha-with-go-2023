package config

import (
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/internal/repository"
	"net/http"
	"os"
)

// TODO: This has multiple reasons to change, fix it
func InitializeRouter() error {
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

	router := http.NewServeMux()
	router.HandleFunc("POST /pessoas", handler.CreatePerson)

	return http.ListenAndServe(":"+port, router)
}
