package playerinfoprocessor

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

type playerService interface {
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
}

type PlayerInfoProcessor struct {
	playerService playerService
}

func NewPlayersInfoProcessor(playerService playerService) *PlayerInfoProcessor {
	return &PlayerInfoProcessor{
		playerService: playerService,
	}
}
