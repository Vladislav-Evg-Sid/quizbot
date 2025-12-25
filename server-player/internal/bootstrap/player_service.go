package bootstrap

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/config"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/playerService"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/storage/pgstorage"
)

func InitPlayerService(storage *pgstorage.PGstorage, kafakaProducer *playerService.KafkaQuizResultProducer, cfg *config.Config) *playerService.PlayerService {
	return playerService.NewPlayerService(context.Background(), storage, kafakaProducer)
}
