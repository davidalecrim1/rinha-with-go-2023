package mock

import (
	"context"
	"go-rinha-de-backend-2023/internal/domain"
	"strings"
)

type MockRepository struct {
	people []domain.Person
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		people: make([]domain.Person, 0),
	}
}

func (m *MockRepository) CreatePerson(_ context.Context, person *domain.Person) error {
	for _, p := range m.people {
		if p.Nickname == person.Nickname {
			return domain.ErrPersonAlreadyExists
		}
	}

	m.people = append(m.people, *person)
	return nil
}

func (m *MockRepository) GetPersonById(_ context.Context, id string) (*domain.Person, error) {
	for _, person := range m.people {
		if person.ID == id {
			return &person, nil
		}
	}
	return nil, domain.ErrPersonNotFound
}

func (m *MockRepository) SearchPersons(_ context.Context, term string) ([]domain.Person, error) {
	var people []domain.Person = make([]domain.Person, 0)
	for _, person := range m.people {
		if strings.Contains(person.Nickname, term) || strings.Contains(person.Name, term) {
			people = append(people, person)
		} else {
			for _, stack := range person.Stack {
				if strings.Contains(stack, term) {
					people = append(people, person)
					break
				}
			}
		}
	}
	return people, nil
}

func (m *MockRepository) GetPersonsCount() (int, error) {
	return len(m.people), nil
}
