package playerserviceapi

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
)

type playerService interface { // TODO: Прописать интерфейс для взаимодействия с БД
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
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
