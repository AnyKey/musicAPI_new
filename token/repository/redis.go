package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"musicAPI/model"
	"time"
)

type Repository struct {
	Redis *redis.Client
}

func NewTokenRepository(redis *redis.Client) *Repository {
	return &Repository{
		Redis: redis,
	}
}
func (repo *Repository) GetToken(ctx context.Context, user string) *model.Tokens {
	var tokens model.Tokens
	res := repo.Redis.Get(ctx, "JWT:"+user)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		return nil
	}
	return &tokens

}
func (repo *Repository) SetToken(ctx context.Context, user string, tokens model.Tokens) error {
	bytes, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	res := repo.Redis.Set(ctx, "JWT:"+user, bytes, 6*time.Hour)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
