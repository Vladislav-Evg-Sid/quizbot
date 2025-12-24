package processors

import "github.com/Vladislav-Evg-Sid/quizbot/client/internal/storage/redisstorage"

type ClientBotHandler struct {
	redis *redisstorage.RedisStorage
}

func NewClientBotHandler(redis *redisstorage.RedisStorage) *ClientBotHandler {
	return &ClientBotHandler{
		redis: redis,
	}
}
