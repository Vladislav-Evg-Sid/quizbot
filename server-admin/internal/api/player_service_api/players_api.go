package playerserviceapi

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
)

type playerService interface {
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
	GetTenQuestionsByTopic(ctx context.Context, topic_name string) ([]*models.Question, int, error)
}

type PlayerServiceAPI struct {
	players_api.UnimplementedPlayersServiceServer
	playerService playerService
}

func NewPlayerServiceAPI(playerService playerService) *PlayerServiceAPI {
	return &PlayerServiceAPI{
		playerService: playerService,
	}
}
