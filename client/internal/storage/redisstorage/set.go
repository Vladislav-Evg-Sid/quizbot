package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
)

func (redisClient *RedisStorage) CreateGameSession(ctx context.Context, userID int64, topicID int, questions []models.Question) (*models.GameSession, error) {
	sessionID := fmt.Sprintf("%d_%d_%d", userID, topicID, time.Now().Unix())

	session := &models.GameSession{
		SessionID:            sessionID,
		UserID:               userID,
		TopicID:              topicID,
		CurrentQuestionIndex: 0,
		Score:                0,
		StartedAt:            time.Now(),
		Questions:            questions,
	}

	// Сохраняем в Redis
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	err = redisClient.redis.Set(ctx, sessionKey, sessionData, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	// Также сохраняем связь user_id -> session_id для быстрого поиска
	userSessionKey := fmt.Sprintf("ps:%d", userID)
	err = redisClient.redis.Set(ctx, userSessionKey, sessionID, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (redisClient *RedisStorage) UpdateGameSession(ctx context.Context, session *models.GameSession) error {
	sessionKey := fmt.Sprintf("s:%s", session.SessionID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return redisClient.redis.Set(ctx, sessionKey, sessionData, time.Hour).Err()
}

func (redisClient *RedisStorage) DeleteGameSession(ctx context.Context, userID int64, sessionID string) error {
	// Удаляем основную сессию
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	err := redisClient.redis.Del(ctx, sessionKey).Err()
	if err != nil {
		return err
	}

	// Удаляем связь user_id -> session_id
	userSessionKey := fmt.Sprintf("ps:%d", userID)
	return redisClient.redis.Del(ctx, userSessionKey).Err()
}

func (redisClient *RedisStorage) CreateUserSession(ctx context.Context, userID int64, permissions []string, currentStep string) (*models.UserSession, error) {
	session := &models.UserSession{
		UserID:      userID,
		Permissions: permissions,
		CurrentStep: currentStep,
	}

	sessionKey := fmt.Sprintf("us:%d", userID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	err = redisClient.redis.Set(ctx, sessionKey, sessionData, time.Hour*24).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (redisClient *RedisStorage) UpdateUserSession(ctx context.Context, session *models.UserSession) error {
	sessionKey := fmt.Sprintf("us:%d", session.UserID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return redisClient.redis.Set(ctx, sessionKey, sessionData, time.Hour*24).Err()
}

func (redisClient *RedisStorage) DeleteUserSession(ctx context.Context, userID int64) error {
	sessionKey := fmt.Sprintf("us:%d", userID)
	return redisClient.redis.Del(ctx, sessionKey).Err()
}
