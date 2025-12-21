package bootstrap

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/config"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/services/adminService"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/storage/pgstorage"
)

func InitAdminService(storage *pgstorage.PGstorage, cfg *config.Config) *adminService.AdminService {
	return adminService.NewAdminService(context.Background(), storage)
}
