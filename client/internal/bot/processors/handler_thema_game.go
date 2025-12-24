package processors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *ClientBotHandler) HandleGetQuestionsForQuiz(ctx context.Context, bot *tgbotapi.BotAPI, topicName string, tg_id int64, session *models.UserSession, playerAPIURL string) {
	parts := strings.SplitN(topicName, ". ", 2)

	if len(parts) < 2 {
		log.Print("Topic processing error: bad topic's name")
		bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка обработки темы викторины"))
		return
	}

	topicName = parts[1]
	resp, err := http.Get(fmt.Sprintf("%s/api/player/tenquestions/%s", playerAPIURL, url.PathEscape(topicName)))
	if err != nil {
		log.Printf("Error calling player API: %v", err)
		bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка соединения с сервером"))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var result models.TenQuestionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка обработки ответа"))
			return
		}

		_, err := h.redis.CreateGameSession(ctx, tg_id, result.TopicId, result.Questions)
		if err != nil {
			log.Printf("Error create play session: %v", err)
			bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка создания игровой сессии"))
			return
		}

		var keyboardButtons [][]tgbotapi.KeyboardButton
		for _, answer := range result.Questions[0].Options {
			keyboardButtons = append(keyboardButtons,
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer)),
			)
		}
		keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

		msg := tgbotapi.NewMessage(tg_id, result.Questions[0].Level+" вопрос:\n"+result.Questions[0].Text)
		msg.ReplyMarkup = keyboard

		session.CurrentStep = "gameProcess"
		err = h.redis.UpdateUserSession(ctx, session)

		bot.Send(msg)
	} else {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		errorMsg := "❌ Ошибка получения вопросов"
		if errMsg, ok := errorResp["error"].(string); ok {
			errorMsg += ": " + errMsg
		}
		bot.Send(tgbotapi.NewMessage(tg_id, errorMsg))
	}
}

func (h *ClientBotHandler) HandleProcessAnswer(ctx context.Context, bot *tgbotapi.BotAPI, userAnswer string, tg_id int64, playerAPI string) {
	hardLevel2Score := map[string]int{
		"простой": 1,
		"средний": 2,
		"сложный": 4,
	}

	session, err := h.redis.GetPlayerSession(ctx, tg_id)
	if err != nil {
		log.Printf("Redis error: %v", err)
		bot.Send(tgbotapi.NewMessage(tg_id, "Ошибка сервера"))
		return
	}

	correctAnswerIndex := session.Questions[session.CurrentQuestionIndex].CorrectIndex
	correctAnswer := session.Questions[session.CurrentQuestionIndex].Options[correctAnswerIndex]

	if userAnswer == correctAnswer {
		session.Score += hardLevel2Score[session.Questions[session.CurrentQuestionIndex].Level]
	}
	session.CurrentQuestionIndex++

	if err := h.redis.UpdateGameSession(ctx, session); err != nil {
		log.Printf("Error update play session: %v", err)
		bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка обновления игровой сессии"))
		return
	}

	if session.CurrentQuestionIndex == len(session.Questions) {
		h.redis.DeleteGameSession(ctx, tg_id, session.SessionID)
		userSession, err := h.redis.GetUserSession(ctx, tg_id)
		if err != nil {
			log.Printf("Error read redis: %v", err)
			bot.Send(tgbotapi.NewMessage(tg_id, "Ошибка сервера"))
			return
		}

		var keyboardButtons [][]tgbotapi.KeyboardButton
		if slices.Contains(userSession.Permissions, "canReadTopics") {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Посмотреть список тем викторины"),
			))
		}

		keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

		new_msg := tgbotapi.NewMessage(tg_id, "Спасибо за игру!\nВаш результат: "+strconv.Itoa(session.Score)+"/22")
		new_msg.ReplyMarkup = keyboard

		quizResults := models.FinishQuizRequest{
			TgID:    tg_id,
			ThemaID: int32(session.TopicID),
			Score:   int32(session.Score),
			Time:    int32(time.Since(session.StartedAt).Seconds()),
		}
		jsonData, _ := json.Marshal(quizResults)
		req, _ := http.NewRequest("PUT", playerAPI+"/api/player/quiz/finish", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error calling admin API: %v", err)
			bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка соединения с сервером"))
			return
		}
		defer resp.Body.Close()

		userSession.CurrentStep = "root"
		err = h.redis.UpdateUserSession(ctx, userSession)
		if err != nil {
			log.Printf("Error update redis: %v", err)
			bot.Send(tgbotapi.NewMessage(tg_id, "Ошибка сервера"))
			return
		}

		bot.Send(new_msg)
		return
	}

	var keyboardButtons [][]tgbotapi.KeyboardButton
	for _, answer := range session.Questions[session.CurrentQuestionIndex].Options {
		keyboardButtons = append(keyboardButtons,
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer)),
		)
	}
	keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

	new_msg := tgbotapi.NewMessage(tg_id, session.Questions[session.CurrentQuestionIndex].Level+"вопрос:\n"+session.Questions[session.CurrentQuestionIndex].Text)
	new_msg.ReplyMarkup = keyboard

	bot.Send(new_msg)
}
