package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidNickname = errors.New("invalid nickname")
	ErrInvalidDate     = errors.New("invalid date")
	ErrInvalidStack    = errors.New("invalid stack")
)

type Person struct {
	ID       string
	Nickname string
	Name     string
	Dob      string
	Stack    []string
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
}

type PersonService struct {
	repo PersonRepository
}

func NewPersonService(repo PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (svc *PersonService) CreatePerson(p *Person) error {
	return svc.repo.CreatePerson(p)
}
