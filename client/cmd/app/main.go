package main

import (
	"fmt"
	"os"

	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига: %v", err))
	}

	clientRedis := bootstrap.InitRedis(cfg)
	telegramBot := bootstrap.InitTgBot(cfg)

	bootstrap.AppRun(cfg, clientRedis, telegramBot)
}
