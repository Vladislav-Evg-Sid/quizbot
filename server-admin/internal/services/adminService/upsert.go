package adminService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

func (s *AdminService) UpsertQuizResult(ctx context.Context, quizResult *models.QuizRequest) error {
	user_id, err := s.adminStorage.GetUserIDByTelegramID(ctx, quizResult.TgID)
	if err != nil {
		return err
	}
	return s.adminStorage.UpsertQuizResult(ctx, quizResult, user_id)
}
