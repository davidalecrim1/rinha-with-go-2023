package mock

import (
	"context"
	"go-rinha-de-backend-2023/internal/domain"
	"strings"
)

type MockRepository struct {
	persons []domain.Person
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		persons: make([]domain.Person, 0),
	}
}

func (m *MockRepository) CreatePerson(_ context.Context, person *domain.Person) error {
	m.persons = append(m.persons, *person)
	return nil
}

func (m *MockRepository) GetPersonByNickname(_ context.Context, nickname string) (*domain.Person, error) {
	for _, person := range m.persons {
		if person.Nickname == nickname {
			return &person, nil
		}
	}
	return nil, domain.ErrNicknameNotFound
}

func (m *MockRepository) GetPersonById(_ context.Context, id string) (*domain.Person, error) {
	for _, person := range m.persons {
		if person.ID == id {
			return &person, nil
		}
	}
	return nil, domain.ErrPersonNotFound
}

func (m *MockRepository) SearchPersons(_ context.Context, term string) ([]domain.Person, error) {
	var persons []domain.Person = make([]domain.Person, 0)
	for _, person := range m.persons {
		if strings.Contains(person.Nickname, term) || strings.Contains(person.Name, term) {
			persons = append(persons, person)
		} else {
			for _, stack := range person.Stack {
				if strings.Contains(stack, term) {
					persons = append(persons, person)
					break
				}
			}
		}
	}
	return persons, nil
}

func (m *MockRepository) GetPersonsCount() (int, error) {
	return len(m.persons), nil
}
