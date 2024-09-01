package repository

import (
	"database/sql"
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
