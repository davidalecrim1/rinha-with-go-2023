package repository

import (
	"context"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var initialBufferSize = 1000
var tickerDelay = 3 * time.Second

type PersonAsyncRepository struct {
	ctx                 context.Context
	db                  *pgxpool.Pool
	buffer              []*domain.Person
	logger              *slog.Logger
	createPersonChannel <-chan *domain.Person
}

func NewPersonAsyncRepository(
	logger *slog.Logger,
	ctx context.Context,
	db *pgxpool.Pool,
) *PersonAsyncRepository {

	return &PersonAsyncRepository{
		logger: logger,
		ctx:    ctx,
		db:     db,
		buffer: make([]*domain.Person, 0, initialBufferSize),
	}
}

func (w *PersonAsyncRepository) Start() {
	w.logger.Debug("worker started")
	ticker := time.NewTicker(tickerDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(w.buffer) > 0 {
				w.bulkCreatePeople()
			}

		case person := <-w.createPersonChannel:
			w.buffer = append(w.buffer, person)

		case <-w.ctx.Done():
			// ensure we don't miss inserting people if the app is stopped
			if len(w.buffer) > 0 {
				w.bulkCreatePeople()
			}
			return
		}
	}
}

func (w *PersonAsyncRepository) bulkCreatePeople() {
	batch := &pgx.Batch{}
	query := "INSERT INTO people (id, nickname, name, dob, stack) VALUES ($1, $2, $3, $4, $5)"

	for _, person := range w.buffer {
		batch.Queue(query, person.ID, person.Nickname, person.Name, person.Dob, strings.Join(person.Stack, " | "))
	}

	br := w.db.SendBatch(context.Background(), batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		w.logger.Error("error while bulk inserting people", "error", err)
		return
	}

	w.cleanBuffer()
}

func (w *PersonAsyncRepository) cleanBuffer() {
	w.buffer = make([]*domain.Person, 0, initialBufferSize)
}

func (w *PersonAsyncRepository) NewCreatePersonChannel() chan<- *domain.Person {
	ch := make(chan *domain.Person)
	w.createPersonChannel = ch
	return ch
}
