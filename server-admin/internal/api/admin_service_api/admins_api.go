package adminserviceapi

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/pb/admins_api"
)

type adminService interface {
	SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error)
	GetUserPermissions(ctx context.Context, tg_id int64) (*models.Permissions, error)
}

type AdminServiceAPI struct {
	admins_api.UnimplementedAdminsServiceServer
	adminService adminService
}

func NewAdminServiceAPI(adminService adminService) *AdminServiceAPI {
	return &AdminServiceAPI{
		adminService: adminService,
	}
}
