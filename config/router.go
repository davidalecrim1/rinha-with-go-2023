package config

import (
	"go-rinha-de-backend-2023/internal/handler"
	"log/slog"
	"net/http"
)

func InitializeRouter(h *handler.PersonHandler, port string, logger *slog.Logger) error {
	router := http.NewServeMux()
	router.HandleFunc("POST /pessoas", h.CreatePerson)
	router.HandleFunc("GET /pessoas/{id}", h.GetPersonById)
	router.HandleFunc("GET /pessoas", h.SearchPersons)
	router.HandleFunc("GET /contagem-pessoas", h.GetPersonsCount)

	logger.Info("server is running on port " + port)
	return http.ListenAndServe(":"+port, router)
}
