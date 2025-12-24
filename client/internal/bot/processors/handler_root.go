package processors

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *ClientBotHandler) HandleChooseThemeCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, playerAPIURL string, session *models.UserSession) {
	resp, err := http.Get(playerAPIURL + "/api/players/topics")
	if err != nil {
		log.Printf("Error calling player API: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка соединения с сервером"))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var result models.AllTopicsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки ответа"))
			return
		}

		var keyboardButtons [][]tgbotapi.KeyboardButton
		var keyboardButtonsRow []tgbotapi.KeyboardButton
		colCountInKeyboard := 2
		for topic_number, topic := range result.Topics {
			keyboardButtonsRow = append(keyboardButtonsRow, tgbotapi.NewKeyboardButton(strconv.Itoa(topic_number+1)+". "+topic.Title))
			if len(keyboardButtonsRow) == colCountInKeyboard {
				keyboardButtons = append(keyboardButtons,
					tgbotapi.NewKeyboardButtonRow(keyboardButtonsRow...),
				)
				keyboardButtonsRow = nil
			}
		}
		if len(keyboardButtonsRow) > 0 {
			keyboardButtons = append(keyboardButtons,
				tgbotapi.NewKeyboardButtonRow(keyboardButtonsRow...),
			)
		}
		keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

		new_msg := tgbotapi.NewMessage(msg.Chat.ID, "Выбирите тему из списка")
		new_msg.ReplyMarkup = keyboard

		session.CurrentStep = "selectThema"
		err = h.redis.UpdateUserSession(ctx, session)
		if err != nil {
			log.Printf("Error update redis: %v", err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ошибка сервера"))
			return
		}

		bot.Send(new_msg)
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
