package admininfoprocessor

type adminService interface {
}

type AdminInfoProcessor struct {
	adminService adminService
}

func NewAdminsInfoProcessor(adminService adminService) *AdminInfoProcessor {
	return &AdminInfoProcessor{
		adminService: adminService,
	}
}
