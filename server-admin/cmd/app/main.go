package main

import (
	"fmt"
	"os"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/config"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига: %v", err))
	}

	adminStorage := bootstrap.InitPGStorage(cfg)
	adminService := bootstrap.InitAdminService(adminStorage, cfg)
	adminInfoProcessor := bootstrap.InitAdminInfoProcessor(adminService)
	adminInfoUpsertConsumer := bootstrap.InitAdminInfoUpsertConsumer(cfg, adminInfoProcessor)
	adminAPI := bootstrap.InitAdminServiceAPI(adminService)

	bootstrap.AppRun(*adminAPI, adminInfoUpsertConsumer, cfg)
}
