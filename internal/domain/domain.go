package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidNickname     = errors.New("invalid nickname")
	ErrInvalidDate         = errors.New("invalid date")
	ErrInvalidStack        = errors.New("invalid stack")
	ErrNicknameNotFound    = errors.New("nickname not found")
	ErrPersonAlreadyExists = errors.New("person already exists")
	ErrPersonNotFound      = errors.New("person not found")
)

type Person struct {
	ID       string   `json:"id"`
	Nickname string   `json:"apelido"`
	Name     string   `json:"nome"`
	Dob      string   `json:"nascimento"`
	Stack    []string `json:"stack"`
}

func NewPerson(nickname string, name string, dob string, stack []string) (*Person, error) {
	person := &Person{
		ID:       uuid.New().String(),
		Nickname: nickname,
		Name:     name,
		Dob:      dob,
		Stack:    stack,
	}

	err := person.Validate()
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *Person) Validate() error {
	var err error

	if err = p.validateName(); err != nil {
		return err
	}

	if err = p.validateNickname(); err != nil {
		return err
	}

	if err = p.validateDate(); err != nil {
		return err
	}

	if err = p.validateStack(); err != nil {
		return err
	}

	return nil
}

func (p *Person) validateName() error {
	if p.Name == "" {
		return ErrInvalidName
	}

	if len(p.Name) > 100 {
		return ErrInvalidName
	}

	return nil
}

func (p *Person) validateNickname() error {
	if p.Nickname == "" {
		return ErrInvalidNickname
	}

	if len(p.Nickname) > 32 {
		return ErrInvalidNickname
	}

	return nil
}

func (p *Person) validateDate() error {
	if len(p.Dob) != 10 {
		return ErrInvalidDate
	}

	return nil
}

func (p *Person) validateStack() error {
	for _, stack := range p.Stack {
		if stack == "" {
			return ErrInvalidStack
		}

		if len(stack) > 32 {
			return ErrInvalidStack
		}
	}

	return nil
}

type PersonRepository interface {
	CreatePerson(person *Person) error
	GetPersonByNickname(nickname string) (*Person, error)
	GetPersonById(id string) (*Person, error)
	SearchPersons(term string) ([]Person, error)
	GetPersonsCount() (int, error)
}

type PersonService struct {
	repo PersonRepository
}

func NewPersonService(repo PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (svc *PersonService) CreatePerson(p *Person) error {
	_, err := svc.repo.GetPersonByNickname(p.Nickname)

	if errors.Is(err, ErrNicknameNotFound) {
		return svc.repo.CreatePerson(p)
	}

	return ErrPersonAlreadyExists
}

func (svc *PersonService) GetPersonById(id string) (*Person, error) {
	return svc.repo.GetPersonById(id)
}

func (svc *PersonService) SearchPersons(term string) ([]Person, error) {
	return svc.repo.SearchPersons(term)
}

func (svc *PersonService) GetPersonsCount() (int, error) {
	return svc.repo.GetPersonsCount()
}
