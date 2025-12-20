package playerserviceapi

import (
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
)

type playerService interface { // ПРОписать интерфейс для взаимодействия с БД
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
