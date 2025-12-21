package bootstrap

import (
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/services/adminService"
	admininfoprocessor "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/services/processors/admin_info_processor"
)

func InitAdminInfoProcessor(adminService *adminService.AdminService) *admininfoprocessor.AdminInfoProcessor {
	return admininfoprocessor.NewAdminsInfoProcessor(adminService)
}
