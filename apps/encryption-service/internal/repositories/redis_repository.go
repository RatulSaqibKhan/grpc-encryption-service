package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisRepository) Set(ctx context.Context, key, value string) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}
