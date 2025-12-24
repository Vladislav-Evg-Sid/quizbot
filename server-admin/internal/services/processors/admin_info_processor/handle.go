package admininfoprocessor

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

func (p *AdminInfoProcessor) Handle(ctx context.Context, quizResult *models.QuizRequest) error {
	return p.adminService.UpsertQuizResult(ctx, quizResult)
}
