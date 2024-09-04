package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
	"time"
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
	Nickname string   `json:"apelido"`
	Name     string   `json:"nome"`
	Dob      string   `json:"nascimento"`
	Stack    []string `json:"stack"`
}

// POST /pessoas
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request CreatePersonRequest
	var err error

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Debug("error decoding request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	person, err := domain.NewPerson(
		request.Nickname,
		request.Name,
		request.Dob,
		request.Stack,
	)

	if err != nil {
		h.logger.Debug("error creating person", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.svc.CreatePerson(ctx, person)

	if errors.Is(err, domain.ErrPersonAlreadyExists) {
		h.logger.Debug("this person already exists", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		h.logger.Debug("error creating person", "error", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Location", "/pessoas/"+person.ID)
	w.WriteHeader(http.StatusCreated)
}

// GET /pessoas/[:id]
func (h *PersonHandler) GetPersonById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")

	if idString == "" {
		h.logger.Debug("id is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := h.svc.GetPersonById(ctx, idString)

	if errors.Is(err, domain.ErrPersonNotFound) {
		h.logger.Debug("person not found", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.Info("error getting person", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(person); err != nil {
		h.logger.Info("error encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *PersonHandler) SearchPersons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	term := r.URL.Query().Get("t")

	if term == "" {
		h.logger.Debug("term is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug("searching persons", "term", term)

	persons, err := h.svc.SearchPersons(ctx, term)

	if err != nil {
		h.logger.Info("error searching persons", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(persons) == 0 {
		h.logger.Debug("no persons found")
	}

	if err = json.NewEncoder(w).Encode(persons); err != nil {
		h.logger.Info("error encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contagem-pessoas
func (h *PersonHandler) GetPersonsCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.GetPersonsCount()

	if err != nil {
		h.logger.Info("error getting persons count", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.Itoa(count)))
}
