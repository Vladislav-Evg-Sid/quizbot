package main

import (
	"fmt"
	"os"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/config"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига: %v", err))
	}

	playerStorage := bootstrap.InitPGStorage(cfg)
	quizResultProducer := bootstrap.InitKafkaProducer(cfg)
	quizResultKafkaProducer := bootstrap.InitNewKafkaQuizResultProducer(cfg, quizResultProducer)
	playerService := bootstrap.InitPlayerService(playerStorage, quizResultKafkaProducer, cfg)
	playerInfoProcessor := bootstrap.InitPlayerInfoProcessor(playerService)
	playerInfoUpsertConsumer := bootstrap.InitPlayerInfoUpsertConsumer(cfg, playerInfoProcessor)
	playerAPI := bootstrap.InitPlayerServiceAPI(playerService, quizResultKafkaProducer)

	bootstrap.AppRun(*playerAPI, playerInfoUpsertConsumer, cfg)
}
