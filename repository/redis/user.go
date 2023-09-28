package redis

import (
	"context"
	"encoding/json"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
	"strconv"
	"time"
)

type userRepositoryRDS struct {
	*Redis
}

func NewUserRepositoryRedis(redis *Redis) repository.User {
	return &userRepositoryRDS{
		redis,
	}
}

func (c *userRepositoryRDS) Create(ctx context.Context, user *model.User) (*model.User, error) {
	bytesUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	key := strconv.Itoa(int(user.ID))
	err = c.Rds.Set(ctx, key, bytesUser, 360*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *userRepositoryRDS) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := new(model.User)

	key := strconv.Itoa(int(id))
	userBytes, err := c.Rds.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(userBytes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
