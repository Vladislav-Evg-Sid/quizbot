package adminService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

func (s *AdminService) SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error) {
	return s.adminStorage.SetUserData(ctx, tg_id, name, username)
}

func (s *AdminService) GetUserPermissions(ctx context.Context, tg_id int64) (*models.Permissions, error) {
	return s.adminStorage.GetUserPermissions(ctx, tg_id)
}
