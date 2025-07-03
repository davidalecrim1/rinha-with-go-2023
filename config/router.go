package config

import (
	"go-rinha-de-backend-2023/internal/handler"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func InitializeRouter(app *fiber.App, h *handler.PersonHandler, logger *slog.Logger) {
	app.Post("/pessoas", h.CreatePerson)
	app.Get("/pessoas/:id", h.GetPersonById)
	app.Get("/pessoas", h.SearchPeople)
	app.Get("/contagem-pessoas", h.GetPeopleCount)
}
