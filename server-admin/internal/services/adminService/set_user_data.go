package adminService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

func (s *AdminService) SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error) {
	// TODO: Сделать обработку входных параметров
	return s.adminStorage.SetUserData(ctx, tg_id, name, username)
}
