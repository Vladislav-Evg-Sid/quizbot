package bootstrap

import (
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/config"
	playerinfoupsertconsumer "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/consumer/player_Info_upsert_consumer"
	playerinfoprocessor "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/processors/player_info_processor"
)

func InitPlayerInfoUpsertConsumer(cfg *config.Config, playerInfoProcessor *playerinfoprocessor.PlayerInfoProcessor) *playerinfoupsertconsumer.PlayerInfoUpsertConsumer {
	return playerinfoupsertconsumer.NewPlayerInfoUpsertConsumer(playerInfoProcessor)
}
