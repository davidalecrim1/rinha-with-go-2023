package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"strings"

	"github.com/lib/pq"
)

type PersonPostgreSqlRepository struct {
	db *sql.DB
}

func NewPersonPostgreSqlRepository(db *sql.DB) *PersonPostgreSqlRepository {
	return &PersonPostgreSqlRepository{db: db}
}

func (repo *PersonPostgreSqlRepository) CreatePerson(ctx context.Context, person *domain.Person) error {
	query := "INSERT INTO persons (id, nickname, name, dob, stack) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.db.ExecContext(ctx, query, person.ID, person.Nickname, person.Name, person.Dob, strings.Join(person.Stack, " | "))

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return domain.ErrPersonAlreadyExists
			}
		}
	}

	return err
}

func (repo *PersonPostgreSqlRepository) GetPersonById(ctx context.Context, id string) (*domain.Person, error) {
	query := "SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack FROM persons WHERE id = $1"
	row := repo.db.QueryRowContext(ctx, query, id)

	var person domain.Person
	err := row.Scan(&person.ID, &person.Nickname, &person.Name, &person.Dob, pq.Array(&person.Stack))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPersonNotFound
	}

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (repo *PersonPostgreSqlRepository) SearchPersons(ctx context.Context, term string) ([]domain.Person, error) {
	query := `
	SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack 
	FROM persons
	WHERE searchable LIKE $1`

	rows, err := repo.db.QueryContext(ctx, query, "%"+term+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var persons []domain.Person = make([]domain.Person, 0)
	for rows.Next() {
		var person domain.Person
		err := rows.Scan(&person.ID, &person.Nickname, &person.Name, &person.Dob, pq.Array(&person.Stack))
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (repo *PersonPostgreSqlRepository) GetPersonsCount() (int, error) {
	query := "SELECT COUNT(id) FROM persons"
	var count int
	err := repo.db.QueryRow(query).Scan(&count)
	return count, err
}
