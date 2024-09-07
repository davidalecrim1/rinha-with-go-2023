package handler

import (
	"context"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PersonHandler struct {
	logger *slog.Logger
	svc    *domain.PersonService
}

func NewPersonHandler(logger *slog.Logger, svc *domain.PersonService) *PersonHandler {
	return &PersonHandler{
		logger: logger,
		svc:    svc,
	}
}

type CreatePersonRequest struct {
	Nickname string   `json:"apelido" validate:"required, max=32"`
	Name     string   `json:"nome" validate:"required, max=100"`
	Dob      string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack    []string `json:"stack" validate:"dive,nonempty, max=32"`
}

// POST /pessoas
func (h *PersonHandler) CreatePerson(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	var request CreatePersonRequest
	var err error

	if err = c.BodyParser(&request); err != nil {
		h.logger.Debug("error decoding request body", "error", err)
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	person, err := domain.NewPerson(
		request.Nickname,
		request.Name,
		request.Dob,
		request.Stack,
	)

	if err != nil {
		h.logger.Debug("error creating person", "error", err)
		c.Status(fiber.StatusUnprocessableEntity)
		return nil
	}

	err = h.svc.CreatePerson(ctx, person)

	if errors.Is(err, domain.ErrPersonAlreadyExists) {
		h.logger.Debug("this person already exists", "error", err)
		c.Status(fiber.StatusUnprocessableEntity)
		return nil
	}

	if err != nil {
		h.logger.Debug("error creating person", "error", err)
		c.Status(fiber.StatusUnprocessableEntity)
		return nil
	}

	c.Set("Location", "/pessoas/"+person.ID)
	c.Status(fiber.StatusCreated)
	return nil
}

// GET /pessoas/[:id]
func (h *PersonHandler) GetPersonById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	idString := c.Params("id")

	if idString == "" {
		h.logger.Debug("id is required")
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	person, err := h.svc.GetPersonById(ctx, idString)

	if errors.Is(err, domain.ErrPersonNotFound) {
		h.logger.Debug("person not found", "error", err)
		c.Status(fiber.StatusNotFound)
		return nil
	}

	if err != nil {
		h.logger.Info("error getting person", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	if err = c.JSON(person); err != nil {
		h.logger.Info("error encoding response", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	return nil
}

func (h *PersonHandler) SearchPersons(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	term := c.Query("t")

	if term == "" {
		h.logger.Debug("term is required")
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	h.logger.Debug("searching people", "term", term)

	people, err := h.svc.SearchPersons(ctx, term)

	if err != nil {
		h.logger.Info("error searching people", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	if len(people) == 0 {
		h.logger.Debug("no people found")
	}

	if err = c.JSON(people); err != nil {
		h.logger.Info("error encoding response", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}
	return nil
}

// GET /contagem-pessoas
func (h *PersonHandler) GetPersonsCount(c *fiber.Ctx) error {
	count, err := h.svc.GetPersonsCount()

	if err != nil {
		h.logger.Info("error getting people count", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Write([]byte(strconv.Itoa(count)))
	return nil
}
