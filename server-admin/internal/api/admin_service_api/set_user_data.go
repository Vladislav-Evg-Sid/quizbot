package adminserviceapi

import (
	"context"
	"log"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/pb/admins_api"
)

func (s *AdminServiceAPI) SetUserData(ctx context.Context, req *admins_api.SetUserDataRequest) (*admins_api.SetUserDataResponce, error) {
	log.Print("Received request")

	responce, err := s.adminService.SetUserData(ctx, req.TelegramId, req.Name, req.Username)
	if err != nil {
		return &admins_api.SetUserDataResponce{}, err
	}

	responcePermissions, err := s.adminService.GetUserPermissions(ctx, req.TelegramId)
	if err != nil {
		return &admins_api.SetUserDataResponce{}, err
	}

	return &admins_api.SetUserDataResponce{
		User: mapUserByResponce(responce, &responcePermissions.Permissions),
	}, err
}

func mapUserByResponce(userInfo *models.User, permissionsInfo *[]string) *admins_api.User {
	return &admins_api.User{
		Id:          userInfo.ID,
		TelegramId:  userInfo.TelegramID,
		Name:        userInfo.Name,
		Username:    userInfo.Username,
		Permissions: *permissionsInfo,
	}
}
