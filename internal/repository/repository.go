package repository

import (
	"context"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CacheRepository interface {
	CreatePerson(ctx context.Context, person *domain.Person) error
	GetPersonById(ctx context.Context, id string) (*domain.Person, error)
	CreateNickname(ctx context.Context, nickname string) error
	CheckNicknameExists(ctx context.Context, nickname string) (bool, error)
}

type PersonRepository struct {
	logger              *slog.Logger
	db                  *pgxpool.Pool
	cache               CacheRepository
	createPersonChannel chan<- *domain.Person
}

func NewPersonRepository(
	logger *slog.Logger,
	db *pgxpool.Pool,
	cache CacheRepository,
	createPersonChannel chan<- *domain.Person) *PersonRepository {

	return &PersonRepository{logger: logger, db: db, cache: cache, createPersonChannel: createPersonChannel}
}

func (r *PersonRepository) CreatePerson(ctx context.Context, person *domain.Person) error {
	err := r.cache.CreatePerson(ctx, person)

	if err != nil {
		return err
	}

	err = r.cache.CreateNickname(ctx, person.Nickname)

	if err != nil {
		return err
	}

	r.logger.Debug("person sent for creation", "person", person)
	r.createPersonAsync(person)
	return nil
}

func (r *PersonRepository) createPersonAsync(person *domain.Person) {
	r.createPersonChannel <- person
}

func (r *PersonRepository) CheckNicknameExists(ctx context.Context, nickname string) (bool, error) {
	nicknameCached, err := r.cache.CheckNicknameExists(ctx, nickname)

	if err != nil {
		return false, err
	}

	if nicknameCached {
		return nicknameCached, nil
	}

	return r.checkNicknameExistsInDatabase(ctx, nickname)
}

func (r *PersonRepository) checkNicknameExistsInDatabase(ctx context.Context, nickname string) (bool, error) {
	query := "SELECT COUNT(id) FROM people WHERE nickname = $1"
	var id int
	err := r.db.QueryRow(ctx, query, nickname).Scan(&id)

	if err != nil {
		return false, err
	}

	if id > 0 {
		return true, nil
	}

	return false, nil
}

func (r *PersonRepository) GetPersonById(ctx context.Context, id string) (*domain.Person, error) {
	cachedPerson, err := r.cache.GetPersonById(ctx, id)

	if err != nil {
		return nil, err
	}

	if cachedPerson != nil {
		return cachedPerson, nil
	}

	return r.getPersonByIdFromDatabase(ctx, id)
}

func (r *PersonRepository) getPersonByIdFromDatabase(ctx context.Context, id string) (*domain.Person, error) {
	query := "SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack FROM people WHERE id = $1"
	row := r.db.QueryRow(ctx, query, id)

	var person domain.Person
	err := row.Scan(&person.ID, &person.Nickname, &person.Name, &person.Dob, &person.Stack)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPersonNotFound
	}

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *PersonRepository) SearchPersons(ctx context.Context, term string) ([]domain.Person, error) {
	query := `
	SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack 
	FROM people
	WHERE searchable LIKE $1`

	rows, err := r.db.Query(ctx, query, "%"+term+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var people []domain.Person = make([]domain.Person, 0)
	for rows.Next() {
		var person domain.Person
		err := rows.Scan(&person.ID, &person.Nickname, &person.Name, &person.Dob, &person.Stack)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

func (r *PersonRepository) GetPersonsCount() (int, error) {
	query := "SELECT COUNT(id) FROM people"
	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}
