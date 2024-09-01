package integration_test

import (
	"bytes"
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

func TestPersonHandler_CreatePerson_201(t *testing.T) {
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

func InitializeHandler() *handler.PersonHandler {
	logger := mock.NewLogger()

	repo := mock.NewMockRepository()
	svc := domain.NewPersonService(repo)
	h := handler.NewPersonHandler(logger, svc)

	return h
}
