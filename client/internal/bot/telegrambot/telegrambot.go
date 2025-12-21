package telegrambot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

func NewTgBot(token string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("ошибка при создании бота: %v", err)
		panic(err)
	}
	return &TelegramBot{Bot: bot}
}
