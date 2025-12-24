package adminService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

type AdminStorage interface {
	SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error)
	GetUserPermissions(ctx context.Context, tg_id int64) (*models.Permissions, error)
	UpsertQuizResult(ctx context.Context, quizResults *models.QuizRequest, user_id int) error
	GetUserIDByTelegramID(ctx context.Context, tg_id int64) (int, error)
}

type AdminService struct {
	adminStorage AdminStorage
}

func NewAdminService(ctx context.Context, adminStorage AdminStorage) *AdminService {
	return &AdminService{
		adminStorage: adminStorage,
	}
}
