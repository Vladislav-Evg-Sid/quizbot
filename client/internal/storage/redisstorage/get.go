package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	"github.com/redis/go-redis/v9"
)

func (redisClient *RediStorage) GetUserSession(ctx context.Context, userID int64) (*models.GameSession, error) {
	// Получаем session_id по user_id
	userSessionKey := fmt.Sprintf("us:%d", userID)
	sessionID, err := redisClient.redis.Get(ctx, userSessionKey).Result()
	if err == redis.Nil {
		return nil, nil // Сессия не найдена
	} else if err != nil {
		return nil, err
	}

	// Получаем данные сессии
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	sessionData, err := redisClient.redis.Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	var session models.GameSession
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
