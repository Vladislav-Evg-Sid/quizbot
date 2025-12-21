package bootstrap

import (
	server "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/api/admin_service_api"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/services/adminService"
)

func InitAdminServiceAPI(adminService *adminService.AdminService) *server.AdminServiceAPI {
	return server.NewAdminServiceAPI(adminService)
}
