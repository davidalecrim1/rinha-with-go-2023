package domain

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName           = errors.New("invalid name")
	ErrInvalidNickname       = errors.New("invalid nickname")
	ErrInvalidDate           = errors.New("invalid date")
	ErrInvalidStack          = errors.New("invalid stack")
	ErrPersonAlreadyExists   = errors.New("person already exists")
	ErrNicknameAlreadyExists = errors.New("nickname already exists")
	ErrPersonNotFound        = errors.New("person not found")
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
	if p.Name == "" || len(p.Name) > 100 {
		return ErrInvalidName
	}
	return nil
}

func (p *Person) validateNickname() error {
	if p.Nickname == "" || len(p.Nickname) > 32 {
		return ErrInvalidNickname
	}
	return nil
}

func (p *Person) validateDate() error {
	_, err := time.Parse("2006-01-02", p.Dob)
	return err
}

func (p *Person) validateStack() error {
	if p.Stack != nil || len(p.Stack) > 0 {
		for _, stack := range p.Stack {
			if stack == "" {
				return ErrInvalidStack
			}

			if len(stack) > 32 {
				return ErrInvalidStack
			}
		}
	}

	return nil
}

type Repository interface {
	CreatePerson(ctx context.Context, person *Person) error
	GetPersonById(ctx context.Context, id string) (*Person, error)
	SearchPeople(ctx context.Context, term string) (*[]Person, error)
	GetPeopleCount() (int, error)
	CheckNicknameExists(ctx context.Context, nickname string) (bool, error)
}

type PersonService struct {
	logger *slog.Logger
	repo   Repository
}

func NewPersonService(logger *slog.Logger, repo Repository) *PersonService {
	return &PersonService{logger: logger, repo: repo}
}

func (svc *PersonService) CreatePerson(ctx context.Context, p *Person) error {
	exists, err := svc.repo.CheckNicknameExists(ctx, p.Nickname)

	if err != nil {
		return err
	}

	if exists {
		return ErrPersonAlreadyExists
	}

	return svc.repo.CreatePerson(ctx, p)
}

func (svc *PersonService) GetPersonById(ctx context.Context, id string) (*Person, error) {
	return svc.repo.GetPersonById(ctx, id)
}

func (svc *PersonService) SearchPeople(ctx context.Context, term string) (*[]Person, error) {
	return svc.repo.SearchPeople(ctx, term)
}

func (svc *PersonService) GetPeopleCount() (int, error) {
	return svc.repo.GetPeopleCount()
}
