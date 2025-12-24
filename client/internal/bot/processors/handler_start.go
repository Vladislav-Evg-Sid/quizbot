package processors

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *ClientBotHandler) HandlerStartCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, adminAPIURL string) {
	userData := models.StartRequest{
		TelegramID: msg.From.ID,
		Name:       msg.From.FirstName + " " + msg.From.LastName,
		Username:   msg.From.UserName,
	}

	jsonData, _ := json.Marshal(userData)
	resp, err := http.Post(adminAPIURL+"/api/users/start", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("Error calling admin API: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка соединения с сервером"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result models.StartResponse
		err := json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки ответа"))
			return
		}

		_, err = h.redis.CreateUserSession(context.Background(), msg.From.ID, result.User.Permissions, "root")
		if err != nil {
			log.Printf("Error crating user session in redis: %v", err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка создания сессии"))
		}

		var keyboardButtons [][]tgbotapi.KeyboardButton
		if slices.Contains(result.User.Permissions, "canReadTopics") {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Посмотреть список тем викторины"),
			))
		}

		keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

		msg := tgbotapi.NewMessage(msg.Chat.ID, "✅ Добро пожаловать в викторину! Вы успешно зарегистрированы. Выберите действие:")
		msg.ReplyMarkup = keyboard

		bot.Send(msg)
	} else {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		errorMsg := "❌ Ошибка регистрации"
		if errMsg, ok := errorResp["error"].(string); ok {
			errorMsg += ": " + errMsg
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, errorMsg))
	}
}
