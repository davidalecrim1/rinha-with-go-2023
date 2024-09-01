package mock

import "go-rinha-de-backend-2023/internal/domain"

type MockRepository struct {
	persons []domain.Person
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		persons: make([]domain.Person, 0),
	}
}

func (m *MockRepository) CreatePerson(person *domain.Person) error {
	m.persons = append(m.persons, *person)
	return nil
}

func (m *MockRepository) GetPersonByNickname(nickname string) (*domain.Person, error) {
	for _, person := range m.persons {
		if person.Nickname == nickname {
			return &person, nil
		}
	}
	return nil, domain.ErrNicknameNotFound
}
