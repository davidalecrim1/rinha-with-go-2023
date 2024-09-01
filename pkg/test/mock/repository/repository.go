package mock_repository

import "go-rinha-de-backend-2023/internal/domain"

type MockPersonRepository struct {
	persons []*domain.Person
}

func NewMockPersonRepository() *MockPersonRepository {
	return &MockPersonRepository{
		persons: []*domain.Person{},
	}
}

func (m *MockPersonRepository) CreatePerson(person *domain.Person) error {
	m.persons = append(m.persons, person)
	return nil
}
