package repository

import (
	"context"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/rueidis"
)

var StandardExpiration = int64(time.Second * 60 * 60) // 1 hour

type PersonCacheRepository struct {
	logger *slog.Logger
	client rueidis.Client
}

func NewPersonCacheRepository(logger *slog.Logger, client rueidis.Client) *PersonCacheRepository {
	return &PersonCacheRepository{
		logger: logger,
		client: client,
	}
}

func (c *PersonCacheRepository) CreatePerson(ctx context.Context, person *domain.Person) error {
	p, _ := sonic.MarshalString(person)
	command := c.client.
		B().
		Set().
		Key("person:" + person.ID).
		Value(p).
		ExSeconds(StandardExpiration).
		Build()
	err := c.client.Do(ctx, command).Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *PersonCacheRepository) CreateNickname(ctx context.Context, nickname string) error {
	command := c.client.
		B().
		Set().
		Key("person:nickname:" + nickname).
		Value("true").
		ExSeconds(StandardExpiration).
		Build()
	err := c.client.Do(ctx, command).Error()

	if err != nil {
		c.logger.Error("error setting nickname in cache", "nickname", nickname, "error", err)
		return err
	}

	return nil
}

func (c *PersonCacheRepository) GetPersonById(ctx context.Context, id string) (*domain.Person, error) {
	cmd := c.client.
		B().
		Get().
		Key("person:" + id).
		Build()
	result, err := c.client.Do(ctx, cmd).AsBytes()

	if rueidis.IsRedisNil(err) {
		c.logger.Debug("person not found in cache", "id", id)
		return nil, nil
	}

	if err != nil {
		c.logger.Error("error getting person from cache", "id", id, "error", err)
		return nil, err
	}

	var person domain.Person
	err = sonic.Unmarshal(result, &person)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (c *PersonCacheRepository) CheckNicknameExists(ctx context.Context, nickname string) (bool, error) {
	cmd := c.client.
		B().
		Exists().
		Key("person:nickname:" + nickname).
		Build()
	result, err := c.client.Do(ctx, cmd).AsBool()

	if !result {
		c.logger.Debug("nickname not found in cache", "nickname", nickname)
		return false, nil
	}

	if err != nil {
		c.logger.Error("error getting nickname from cache", "nickname", nickname, "error", err)
		return false, err
	}

	c.logger.Debug("nickname found in cache", "nickname", nickname)
	return result, nil
}

func (c *PersonCacheRepository) SetSearchPeople(ctx context.Context, term string, people *[]domain.Person) error {
	peopleStr, _ := sonic.MarshalString(people)
	cmd := c.client.
		B().
		Set().
		Key("search:" + term).
		Value(peopleStr).
		ExSeconds(20).
		Build()
	err := c.client.Do(ctx, cmd).Error()

	if err != nil {
		c.logger.Debug("error setting search people in cache", "term", term, "error", err)
		return err
	}

	return nil
}

func (c *PersonCacheRepository) GetSearchPeople(ctx context.Context, term string) (*[]domain.Person, error) {
	cmd := c.client.
		B().
		Get().
		Key("search:" + term).
		Build()
	results, err := c.client.Do(ctx, cmd).AsBytes()

	if rueidis.IsRedisNil(err) {
		c.logger.Debug("search people not found in cache", "term", term, "error", err)
		return nil, nil
	}

	if err != nil {
		c.logger.Error("error getting search people from cache", "term", term, "error", err)
		return nil, err
	}

	var people []domain.Person
	sonic.Unmarshal(results, &people)

	return &people, nil
}
