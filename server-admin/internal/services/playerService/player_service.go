package playerService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

type PlayerStorage interface {
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
	GetTopicIdByName(ctx context.Context, topic_name string) (int, error)
	GetTenQuestionsByTopicID(ctx context.Context, topic_id int) ([]*models.Question, error)
}

type PlayerService struct {
	playerStorage PlayerStorage
}

func NewPlayerService(ctx context.Context, playerStorage PlayerStorage) *PlayerService {
	return &PlayerService{
		playerStorage: playerStorage,
	}
}
