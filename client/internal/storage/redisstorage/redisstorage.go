package redisstorage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	redis *redis.Client
}

func NewRedisStore(ctx context.Context, redis_URL, redis_password string, redis_db int) (*RedisStorage, error) {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     redis_URL,
		Password: redis_password,
		DB:       redis_db,
	})

	_, err := redis_client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v; %v", err, redis_URL)
	}

	return &RedisStorage{redis: redis_client}, nil
}
