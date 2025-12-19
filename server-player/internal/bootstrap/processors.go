package bootstrap

import (
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/playerService"
	playerinfoprocessor "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/processors/player_info_processor"
)

func InitPlayerInfoProcessor(playerService *playerService.PlayerService) *playerinfoprocessor.PlayerInfoProcessor {
	return playerinfoprocessor.NewStudentsInfoProcessor(playerService)
}
