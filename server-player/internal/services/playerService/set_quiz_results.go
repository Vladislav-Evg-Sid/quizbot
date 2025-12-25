package playerService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

func (s *PlayerService) SetResultsByQuiz(ctx context.Context, req *models.QuizRequest) error {
	return s.kafkaProducer.SendQuizResult(ctx, req)
}
