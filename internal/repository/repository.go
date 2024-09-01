package repository

import (
	"database/sql"
	"errors"
	"go-rinha-de-backend-2023/internal/domain"

	"github.com/lib/pq"
)

type PersonPostgreSqlRepository struct {
	db *sql.DB
}

func NewPersonPostgreSqlRepository(db *sql.DB) *PersonPostgreSqlRepository {
	return &PersonPostgreSqlRepository{db: db}
}

func (repo *PersonPostgreSqlRepository) CreatePerson(person *domain.Person) error {
	query := "INSERT INTO persons (id, nickname, name, dob, stack) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.db.Exec(query, person.ID, person.Nickname, person.Name, person.Dob, pq.Array(person.Stack))
	return err
}

func (repo *PersonPostgreSqlRepository) GetPersonByNickname(nickname string) (*domain.Person, error) {
	query := "SELECT id, nickname, name, dob, stack FROM persons WHERE nickname = $1"
	row := repo.db.QueryRow(query, nickname)

	var person domain.Person
	err := row.Scan(&person.ID, &person.Nickname, &person.Name, &person.Dob, pq.Array(&person.Stack))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNicknameNotFound
	}

	if err != nil {
		return nil, err
	}

	return &person, nil
}
