package adminserviceapi

import (
	"context"
	"log"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/pb/admins_api"
)

func (s *AdminServiceAPI) SetUserData(ctx context.Context, req *admins_api.SetUserDataRequest) (*admins_api.SetUserDataResponse, error) {
	log.Print("Received request")

	responce, err := s.adminService.SetUserData(ctx, req.TelegramId, req.Name, req.Username)
	if err != nil {
		return &admins_api.SetUserDataResponse{}, err
	}
	return &admins_api.SetUserDataResponse{
		User: mapUserByResponce(responce),
	}, err
}

func mapUserByResponce(userInfo *models.User) *admins_api.User {
	return &admins_api.User{
		Id:         userInfo.ID,
		TelegramId: userInfo.TelegramID,
		Name:       userInfo.Name,
		Username:   userInfo.Username,
		IsAdmin:    userInfo.IsAdmin,
	}
}
