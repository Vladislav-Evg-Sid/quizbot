package bootstrap

import (
	"context"
	"log"

	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/storage/redisstorage"
)

func InitRedis(cfg *config.Config) *redisstorage.RediStorage {
	redis_URL := cfg.Redis.Url
	redis_password := cfg.Redis.Password
	redis_db := cfg.Redis.Database
	redis_client, err := redisstorage.NewRedisStore(context.Background(), redis_URL, redis_password, redis_db)
	if err != nil {
		log.Panicf("ошибка подключения к redis: %v", err)
		panic(err)
	}
	return redis_client
}
