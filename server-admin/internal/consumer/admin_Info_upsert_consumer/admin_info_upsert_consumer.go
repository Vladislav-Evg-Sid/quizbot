package admininfoupsertconsumer

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
)

type adminInfoProcessor interface {
	Handle(ctx context.Context, quizResult *models.QuizRequest) error
}

type AdminInfoUpsertConsumer struct {
	adminInfoProcessor adminInfoProcessor
	kafkaBroker        []string
	topicName          string
}

func NewAdminInfoUpsertConsumer(adminInfoProcessor adminInfoProcessor, kafkaBroker []string, topicName string) *AdminInfoUpsertConsumer {
	return &AdminInfoUpsertConsumer{
		adminInfoProcessor: adminInfoProcessor,
		kafkaBroker:        kafkaBroker,
		topicName:          topicName,
	}
}
