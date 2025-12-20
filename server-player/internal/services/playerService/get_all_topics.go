package playerService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

func (s *PlayerService) GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error) {
	return s.playerStorage.GetAllTopics(ctx)
}
