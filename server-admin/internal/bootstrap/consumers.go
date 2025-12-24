package bootstrap

import (
	"fmt"

	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/config"
	admininfoupsertconsumer "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/consumer/admin_Info_upsert_consumer"
	admininfoprocessor "github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/services/processors/admin_info_processor"
)

func InitAdminInfoUpsertConsumer(cfg *config.Config, adminInfoProcessor *admininfoprocessor.AdminInfoProcessor) *admininfoupsertconsumer.AdminInfoUpsertConsumer {
	kafkaBrockers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return admininfoupsertconsumer.NewAdminInfoUpsertConsumer(adminInfoProcessor, kafkaBrockers, cfg.Kafka.TopicNameAddQuizResult)
}
