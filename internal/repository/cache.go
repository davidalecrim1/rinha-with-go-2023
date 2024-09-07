package repository

import (
	"context"
	"encoding/json"
	"go-rinha-de-backend-2023/internal/domain"
	"log/slog"
	"time"

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
	jsonData, _ := json.Marshal(person)
	command := c.client.B().Set().Key("person:" + person.ID).Value(string(jsonData)).ExSeconds(StandardExpiration).Build()
	err := c.client.Do(ctx, command).Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *PersonCacheRepository) CreateNickname(ctx context.Context, nickname string) error {
	command := c.client.B().Set().Key("person:nickname:" + nickname).Value("true").ExSeconds(StandardExpiration).Build()
	err := c.client.Do(ctx, command).Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *PersonCacheRepository) GetPersonById(ctx context.Context, id string) (*domain.Person, error) {
	cmd := c.client.B().Get().Key("person:" + id).Build()
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
	err = json.Unmarshal(result, &person)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (c *PersonCacheRepository) CheckNicknameExists(ctx context.Context, nickname string) (bool, error) {
	cmd := c.client.B().Exists().Key("person:nickname:" + nickname).Build()
	result, err := c.client.Do(ctx, cmd).AsBool()

	if !result {
		c.logger.Debug("nickname not found in cache", "nickname", nickname)
		return false, nil
	}

	if err != nil {
		c.logger.Error("error getting nickname from cache", "nickname", nickname, "error", err)
		return false, err
	}

	return result, nil
}
