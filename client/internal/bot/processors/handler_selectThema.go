package processors

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *ClientBotHandler) HandleSelectThemsProcessing(ctx context.Context, thema string, tg_id int64, session *models.UserSession, bot *tgbotapi.BotAPI) {
	var keyboardButtons [][]tgbotapi.KeyboardButton
	if slices.Contains(session.Permissions, "canReadQuestions") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Посмотреть список вопросов темы"),
		))
	}
	if slices.Contains(session.Permissions, "canPlayInQuiz") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Начать игру"),
		))
	}
	if slices.Contains(session.Permissions, "canReadRatings") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Посмотреть рейтинг"),
		))
	}
	if slices.Contains(session.Permissions, "canUpdateTopics") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Обновить название темы"),
		))
	}
	if slices.Contains(session.Permissions, "canDeleteTopics") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить тему"),
		))
	}
	if slices.Contains(session.Permissions, "canCreateQuestions") {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить новые вопросы"),
		))
	}

	keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

	msg := tgbotapi.NewMessage(tg_id, fmt.Sprintf("Выбирите действие для темы \"%s\"", thema))
	msg.ReplyMarkup = keyboard

	session.CurrentStep = fmt.Sprintf("thema:%s", thema)
	err := h.redis.UpdateUserSession(ctx, session)
	if err != nil {
		log.Printf("Error update redis: %v", err)
		bot.Send(tgbotapi.NewMessage(tg_id, "Ошибка сервера"))
		return
	}

	bot.Send(msg)
}
