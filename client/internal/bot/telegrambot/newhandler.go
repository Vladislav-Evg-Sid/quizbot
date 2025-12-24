package telegrambot

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/storage/redisstorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClientBotHandler interface {
	HandlerStartCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, adminAPIURL string)
	HandleSelectThemsProcessing(ctx context.Context, thema string, tg_id int64, session *models.UserSession, bot *tgbotapi.BotAPI)
	HandleChooseThemeCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, playerAPIURL string, session *models.UserSession)
	HandleGetQuestionsForQuiz(ctx context.Context, bot *tgbotapi.BotAPI, topicName string, tg_id int64, session *models.UserSession, playerAPIURL string)
	HandleProcessAnswer(ctx context.Context, bot *tgbotapi.BotAPI, userAnswer string, tg_id int64, playerAPIURL string)
}

type BotHandler struct {
	redis   *redisstorage.RedisStorage
	handler ClientBotHandler
}

func NewBotHandler(redis *redisstorage.RedisStorage, handler ClientBotHandler) *BotHandler {
	return &BotHandler{
		redis:   redis,
		handler: handler,
	}
}
