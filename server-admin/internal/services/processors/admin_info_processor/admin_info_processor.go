package admininfoprocessor

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

type adminService interface {
	UpsertQuizResult(ctx context.Context, quizResult *models.QuizRequest) error
}

type AdminInfoProcessor struct {
	adminService adminService
}

func NewAdminsInfoProcessor(adminService adminService) *AdminInfoProcessor {
	return &AdminInfoProcessor{
		adminService: adminService,
	}
}
