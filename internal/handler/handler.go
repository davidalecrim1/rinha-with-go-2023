package handler

import (
	"encoding/json"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"net/http"
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
	var request CreatePersonRequest
	var err error

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Debug("error decoding request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	err = h.svc.CreatePerson(person)

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

	h.logger.Debug("successful request")
	w.Header().Set("Location", "/pessoas/"+person.ID)
	w.WriteHeader(http.StatusCreated)
}

// GET /pessoas/[:id]
func (h *PersonHandler) GetPersonById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	if idString == "" {
		h.logger.Debug("id is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := h.svc.GetPersonById(idString)

	if errors.Is(err, domain.ErrPersonNotFound) {
		h.logger.Debug("person not found", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.Debug("error getting person", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(person); err != nil {
		h.logger.Debug("error encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PersonHandler) SearchPersons(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	if term == "" {
		h.logger.Debug("term is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debug("searching persons", "term", term)

	persons, err := h.svc.SearchPersons(term)

	if err != nil {
		h.logger.Debug("error searching persons", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(persons) == 0 {
		h.logger.Debug("no persons found")
	}

	if err = json.NewEncoder(w).Encode(persons); err != nil {
		h.logger.Debug("error encoding response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
