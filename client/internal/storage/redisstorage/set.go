package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
)

func (redisClient *RediStorage) CreateGameSession(ctx context.Context, userID int64, topicID int, questions []models.Question) (*models.GameSession, error) {
	sessionID := fmt.Sprintf("%d_%d_%d", userID, topicID, time.Now().Unix())

	session := &models.GameSession{
		SessionID:            sessionID,
		UserID:               userID,
		TopicID:              topicID,
		CurrentQuestionIndex: 0,
		Score:                0,
		TotalTime:            0,
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
	userSessionKey := fmt.Sprintf("us:%d", userID)
	err = redisClient.redis.Set(ctx, userSessionKey, sessionID, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (redisClient *RediStorage) UpdateGameSession(ctx context.Context, session *models.GameSession) error {
	sessionKey := fmt.Sprintf("s:%s", session.SessionID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return redisClient.redis.Set(ctx, sessionKey, sessionData, time.Hour).Err()
}

func (redisClient *RediStorage) DeleteGameSession(ctx context.Context, userID int64, sessionID string) error {
	// Удаляем основную сессию
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	err := redisClient.redis.Del(ctx, sessionKey).Err()
	if err != nil {
		return err
	}

	// Удаляем связь user_id -> session_id
	userSessionKey := fmt.Sprintf("us:%d", userID)
	return redisClient.redis.Del(ctx, userSessionKey).Err()
}
