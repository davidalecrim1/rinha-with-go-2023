package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-rinha-de-backend-2023/internal/domain"
	"go-rinha-de-backend-2023/internal/handler"
	"go-rinha-de-backend-2023/test/integration/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonHandler_CreatePerson(t *testing.T) {
	t.Run("valid person with stack", func(t *testing.T) {
		h := InitializeHandler()

		personRequest := handler.CreatePersonRequest{
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    []string{"Go", "Docker"},
		}

		body, err := json.Marshal(personRequest)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreatePerson(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		location := rr.Header().Get("Location")
		assert.Contains(t, location, "/pessoas/")
	})

	t.Run("valid person without stack", func(t *testing.T) {
		h := InitializeHandler()

		personRequest := handler.CreatePersonRequest{
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    nil,
		}

		body, err := json.Marshal(personRequest)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		h.CreatePerson(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		location := rr.Header().Get("Location")
		assert.Contains(t, location, "/pessoas/")
	})

	t.Run("valid person that already exists", func(t *testing.T) {
		h := InitializeHandler()

		personRequest := handler.CreatePersonRequest{
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    nil,
		}

		sameBody, err := json.Marshal(personRequest)
		require.NoError(t, err)

		req1 := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(sameBody))
		rr1 := httptest.NewRecorder()
		h.CreatePerson(rr1, req1)

		req2 := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(sameBody))
		rr2 := httptest.NewRecorder()
		h.CreatePerson(rr2, req2)

		assert.Equal(t, http.StatusUnprocessableEntity, rr2.Code)
	})
}

func TestPersonHandler_GetPersonById(t *testing.T) {

	// Not working currently because the r.pathValue(id)
	// cannot match given is not using servermux. Fix it later
	t.Run("get someone by uuid", func(t *testing.T) {
		logger := mock.NewLogger()

		repo := mock.NewMockRepository()
		svc := domain.NewPersonService(repo)
		h := handler.NewPersonHandler(logger, svc)

		id := "5ce4668c-4710-4cfb-ae5f-38988d6d49cb"

		err := repo.CreatePerson(context.Background(), &domain.Person{
			ID:       id,
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    []string{"Go", "Docker"},
		})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/pessoas/"+id, nil)
		rr := httptest.NewRecorder()

		h.GetPersonById(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "johndoe")
	})
}

func TestPersonHandler_SeachPersons(t *testing.T) {
	t.Run("search persons by term with results", func(t *testing.T) {
		logger := mock.NewLogger()

		repo := mock.NewMockRepository()
		svc := domain.NewPersonService(repo)
		h := handler.NewPersonHandler(logger, svc)

		var err error
		err = repo.CreatePerson(context.Background(), &domain.Person{
			ID:       "5ce4668c-4710-4cfb-ae5f-38988d6d49cb",
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    []string{"Go", "Docker"},
		})
		require.NoError(t, err)

		err = repo.CreatePerson(context.Background(), &domain.Person{
			ID:       "f7379ae8-8f9b-4cd5-8221-51efe19e721b",
			Nickname: "janedoe",
			Name:     "Jane Doe",
			Dob:      "1990-01-01",
			Stack:    []string{"Python", "Ruby", "Angular"},
		})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/pessoas?t=Doe", nil)
		rr := httptest.NewRecorder()

		h.SearchPersons(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "johndoe")
	})

	t.Run("search persons with empty results", func(t *testing.T) {
		logger := mock.NewLogger()

		repo := mock.NewMockRepository()
		svc := domain.NewPersonService(repo)
		h := handler.NewPersonHandler(logger, svc)

		req := httptest.NewRequest("GET", "/pessoas?t=Doe", nil)
		rr := httptest.NewRecorder()

		h.SearchPersons(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "[]")
	})
}

func TestPersonHandler_GetPersonsCount(t *testing.T) {
	t.Run("get persons count", func(t *testing.T) {
		logger := mock.NewLogger()

		repo := mock.NewMockRepository()
		svc := domain.NewPersonService(repo)
		h := handler.NewPersonHandler(logger, svc)

		err := repo.CreatePerson(context.Background(), &domain.Person{
			ID:       "5ce4668c-4710-4cfb-ae5f-38988d6d49cb",
			Nickname: "johndoe",
			Name:     "John Doe",
			Dob:      "1990-01-01",
			Stack:    []string{"Go", "Docker"},
		})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/contagem-pessoas", nil)
		rr := httptest.NewRecorder()

		h.GetPersonsCount(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "1")
	})
}

func InitializeHandler() *handler.PersonHandler {
	logger := mock.NewLogger()

	repo := mock.NewMockRepository()
	svc := domain.NewPersonService(repo)
	h := handler.NewPersonHandler(logger, svc)

	return h
}
