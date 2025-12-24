package admininfoupsertconsumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *AdminInfoUpsertConsumer) Consume(ctx context.Context) { // TODO: Добавить кафку
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "QuizResult_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("StudentInfoUpsertConsumer.consume error", "error", err.Error())
		}
		var quizResult *models.QuizRequest
		err = json.Unmarshal(msg.Value, &quizResult)
		if err != nil {
			slog.Error("parce", "error", err)
			continue
		}
		err = c.adminInfoProcessor.Handle(ctx, quizResult)
		if err != nil {
			slog.Error("Handle", "error", err)
		}
	}
}
