package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	"github.com/redis/go-redis/v9"
)

func (redisClient *RedisStorage) GetPlayerSession(ctx context.Context, userID int64) (*models.GameSession, error) {
	// Получаем session_id по user_id
	GameSessionKey := fmt.Sprintf("ps:%d", userID)
	sessionID, err := redisClient.redis.Get(ctx, GameSessionKey).Result()
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

func (redisClient *RedisStorage) GetUserSession(ctx context.Context, userID int64) (*models.UserSession, error) {
	sessionID := fmt.Sprintf("us:%d", userID)
	sessionData, err := redisClient.redis.Get(ctx, sessionID).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var session models.UserSession
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
