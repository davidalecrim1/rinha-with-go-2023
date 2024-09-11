package repository

import (
	"context"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonPostgreSqlRepository struct {
	db *pgxpool.Pool
}

func NewPersonPostgreSqlRepository(db *pgxpool.Pool) *PersonPostgreSqlRepository {
	return &PersonPostgreSqlRepository{db: db}
}

func (repo *PersonPostgreSqlRepository) CreatePerson(ctx context.Context, person *domain.Person) error {
	query := "INSERT INTO people (id, nickname, name, dob, stack) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.db.Exec(ctx, query, person.ID, person.Nickname, person.Name, person.Dob, strings.Join(person.Stack, " | "))

	if err != nil {
		if pgxErr, ok := err.(*pgconn.PgError); ok {
			if pgxErr.Code == "23505" {
				return domain.ErrPersonAlreadyExists
			}
		}
	}

	return err
}

func (repo *PersonPostgreSqlRepository) GetPersonById(ctx context.Context, id string) (*domain.Person, error) {
	query := "SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack FROM people WHERE id = $1"
	row := repo.db.QueryRow(ctx, query, id)

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

func (repo *PersonPostgreSqlRepository) SearchPersons(ctx context.Context, term string) ([]domain.Person, error) {
	query := `
	SELECT id, nickname, name, dob, string_to_array(stack, ' | ') as stack 
	FROM people
	WHERE searchable LIKE $1`

	rows, err := repo.db.Query(ctx, query, "%"+term+"%")

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

func (repo *PersonPostgreSqlRepository) GetPersonsCount() (int, error) {
	query := "SELECT COUNT(id) FROM people"
	var count int
	err := repo.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}
