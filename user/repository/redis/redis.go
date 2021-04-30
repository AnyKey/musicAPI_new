package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"musicAPI/user"
	"time"
)

type Repository struct {
	Redis *redis.Client
}

func New(redis *redis.Client) *Repository {
	return &Repository{
		Redis: redis,
	}
}
func (repo *Repository) GetToken(ctx context.Context, username string) *user.Tokens {
	var tokens user.Tokens
	res := repo.Redis.Get(ctx, "JWT:"+username)
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
func (repo *Repository) SetToken(ctx context.Context, username string, tokens user.Tokens) error {
	bytes, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	res := repo.Redis.Set(ctx, "JWT:"+username, bytes, 6*time.Hour)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
