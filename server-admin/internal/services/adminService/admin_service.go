package adminService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

type AdminStorage interface {
	SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error)
}

type AdminService struct {
	adminStorage AdminStorage
}

func NewAdminService(ctx context.Context, adminStorage AdminStorage) *AdminService {
	return &AdminService{
		adminStorage: adminStorage,
	}
}
