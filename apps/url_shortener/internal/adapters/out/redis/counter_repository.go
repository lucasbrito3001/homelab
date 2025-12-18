package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCounterRepository struct {
	client *redis.Client
}

func NewRedisCounterRepository(client *redis.Client) *RedisCounterRepository {
	return &RedisCounterRepository{
		client: client,
	}
}

func (r *RedisCounterRepository) Increment(ctx context.Context) (int64, error) {
	return r.client.Incr(ctx, "counter").Result()
}
