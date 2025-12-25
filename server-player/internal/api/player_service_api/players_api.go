package playerserviceapi

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
)

type QuizResultProducer interface {
	SendQuizResult(ctx context.Context, result *models.QuizRequest) error
}

type playerService interface {
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
	GetTenQuestionsByTopic(ctx context.Context, topic_name string) ([]*models.Question, int, error)
	SetResultsByQuiz(ctx context.Context, req *models.QuizRequest) error
}

type PlayerServiceAPI struct {
	players_api.UnimplementedPlayersServiceServer
	playerService      playerService
	QuizResultProducer QuizResultProducer
}

func NewPlayerServiceAPI(playerService playerService, producer QuizResultProducer) *PlayerServiceAPI {
	return &PlayerServiceAPI{
		playerService:      playerService,
		QuizResultProducer: producer,
	}
}
