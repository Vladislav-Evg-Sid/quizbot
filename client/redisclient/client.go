package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"quiz-bot-client/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // нет пароля
		DB:       0,  // база по умолчанию
	})

	// Проверяем подключение
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v; %v", err, redisURL)
	}

	log.Println("✅ Connected to Redis")
	return nil
}

func CreateGameSession(userID int64, topicID int, questions []models.Question) (*models.GameSession, error) {
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

	err = redisClient.Set(ctx, sessionKey, sessionData, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	// Также сохраняем связь user_id -> session_id для быстрого поиска
	userSessionKey := fmt.Sprintf("us:%d", userID)
	err = redisClient.Set(ctx, userSessionKey, sessionID, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetUserSession(userID int64) (*models.GameSession, error) {
	// Получаем session_id по user_id
	userSessionKey := fmt.Sprintf("us:%d", userID)
	sessionID, err := redisClient.Get(ctx, userSessionKey).Result()
	if err == redis.Nil {
		return nil, nil // Сессия не найдена
	} else if err != nil {
		return nil, err
	}

	// Получаем данные сессии
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	sessionData, err := redisClient.Get(ctx, sessionKey).Result()
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

func UpdateGameSession(session *models.GameSession) error {
	sessionKey := fmt.Sprintf("s:%s", session.SessionID)
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, sessionKey, sessionData, time.Hour).Err()
}

func DeleteGameSession(userID int64, sessionID string) error {
	// Удаляем основную сессию
	sessionKey := fmt.Sprintf("s:%s", sessionID)
	err := redisClient.Del(ctx, sessionKey).Err()
	if err != nil {
		return err
	}

	// Удаляем связь user_id -> session_id
	userSessionKey := fmt.Sprintf("us:%d", userID)
	return redisClient.Del(ctx, userSessionKey).Err()
}
