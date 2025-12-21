package bootstrap

import (
	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/bot/telegrambot"
)

func InitTgBot(cfg *config.Config) *telegrambot.TelegramBot {
	return telegrambot.NewTgBot(cfg.TgBot.Token)
}
