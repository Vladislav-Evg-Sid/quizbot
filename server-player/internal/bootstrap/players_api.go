package bootstrap

import (
	server "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/api/player_service_api"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/playerService"
)

func InitPlayerServiceAPI(playerService *playerService.PlayerService, quizResultProducer *playerService.KafkaQuizResultProducer) *server.PlayerServiceAPI {
	return server.NewPlayerServiceAPI(playerService, quizResultProducer)
}
