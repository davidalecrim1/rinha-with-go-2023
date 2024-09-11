package mock_repository

import "go-rinha-de-backend-2023/internal/domain"

type MockPersonRepository struct {
	people []*domain.Person
}

func NewMockPersonRepository() *MockPersonRepository {
	return &MockPersonRepository{
		people: []*domain.Person{},
	}
}

func (m *MockPersonRepository) CreatePerson(person *domain.Person) error {
	m.people = append(m.people, person)
	return nil
}
