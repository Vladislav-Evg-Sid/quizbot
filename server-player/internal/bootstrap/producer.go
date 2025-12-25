package bootstrap

import (
	"fmt"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/config"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/kafka"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/playerService"
)

func InitNewKafkaQuizResultProducer(cfg *config.Config, producer kafka.Producer) *playerService.KafkaQuizResultProducer {
	return playerService.NewKafkaQuizResultProducer(producer, cfg.Kafka.TopicNameAddQuizResult)
}

func InitKafkaProducer(cfg *config.Config) kafka.Producer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return kafka.NewProducer(brokers)
}
