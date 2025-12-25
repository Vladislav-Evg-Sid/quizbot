package playerService

import (
	"context"
	"encoding/json"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/kafka"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

type KafkaQuizResultProducer struct {
	Producer kafka.Producer
	Topic    string
}

func NewKafkaQuizResultProducer(producer kafka.Producer, topic string) *KafkaQuizResultProducer {
	return &KafkaQuizResultProducer{
		Producer: producer,
		Topic:    topic,
	}
}

func (k *KafkaQuizResultProducer) SendQuizResult(ctx context.Context, result *models.QuizRequest) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return k.Producer.SendMessage(ctx, k.Topic, nil, data)
}
