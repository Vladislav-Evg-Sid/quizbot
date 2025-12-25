package adminserviceapi

import (
	"context"
	"log"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/pb/admins_api"
)

func (s *AdminServiceAPI) GetUserPermissions(ctx context.Context, req *admins_api.GetUserPermissionsRequest) (*admins_api.GetUserPermissionsResponce, error) {
	log.Printf("Received request, %d", req.TelegramId)

	responce, err := s.adminService.GetUserPermissions(ctx, req.TelegramId)
	if err != nil {
		return &admins_api.GetUserPermissionsResponce{}, err
	}

	return &admins_api.GetUserPermissionsResponce{
		Permissions: responce.Permissions,
	}, nil
}
