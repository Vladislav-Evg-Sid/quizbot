package telegrambot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/storage/redisstorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	redis *redisstorage.RediStorage
}

func NewBotHandler(redis *redisstorage.RediStorage) *BotHandler {
	return &BotHandler{
		redis: redis,
	}
}

func (h *BotHandler) HandleUpdate(tgBot *TelegramBot, update tgbotapi.Update, cfg *config.Config) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	switch {
	case update.Message.Text == "/start":
		handleStartCommand(tgBot.Bot, update.Message, cfg.Network.AdminREST)
	case strings.HasSuffix(strings.ToLower(update.Message.Text), "выбрать тему викторины"):
		handleChooseThemeCommand(tgBot.Bot, update.Message, cfg.Network.PlayerREST)
	default:
		session, err := h.redis.GetUserSession(context.Background(), update.Message.From.ID)
		if err != nil {
			log.Printf("Redis error: %v", err)
			tgBot.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Ошибка сервера"))
			return
		}

		if session == nil {
			// Нет игровой сессии: создаём новую
			h.handleGetQuestionsForQuiz(tgBot.Bot, update.Message, cfg.Network.PlayerREST)
		} else {
			// Есть игровая сессия: продолжаем играть
			h.handleProcessAnswer(tgBot.Bot, update.Message, session)
		}
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, adminAPIURL string) {
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

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🎯 Выбрать тему викторины"),
			),
		)

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

func handleChooseThemeCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, playerAPIURL string) {
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

		msg := tgbotapi.NewMessage(msg.Chat.ID, "Выбирите тему из списка")
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

func (h *BotHandler) handleGetQuestionsForQuiz(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, playerAPIURL string) {
	topicName := msg.Text
	parts := strings.SplitN(topicName, ". ", 2)

	if len(parts) < 2 {
		log.Print("Topic processing error: bad topic's name")
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки темы викторины"))
		return
	}

	topicName = parts[1]
	resp, err := http.Get(fmt.Sprintf("%s/api/player/tenquestions/%s", playerAPIURL, url.PathEscape(topicName)))
	if err != nil {
		log.Printf("Error calling player API: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка соединения с сервером"))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var result models.TenQuestionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки ответа"))
			return
		}

		_, err := h.redis.CreateGameSession(context.Background(), msg.From.ID, result.TopicId, result.Questions)
		if err != nil {
			log.Printf("Error create play session: %v", err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка создания игровой сессии"))
			return
		}

		var keyboardButtons [][]tgbotapi.KeyboardButton
		for _, answer := range result.Questions[0].Options {
			keyboardButtons = append(keyboardButtons,
				tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer)),
			)
		}
		keyboard := tgbotapi.NewReplyKeyboard(keyboardButtons...)

		msg := tgbotapi.NewMessage(msg.Chat.ID, result.Questions[0].Level+" вопрос:\n"+result.Questions[0].Text)
		msg.ReplyMarkup = keyboard

		bot.Send(msg)
	} else {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		errorMsg := "❌ Ошибка получения вопросов"
		if errMsg, ok := errorResp["error"].(string); ok {
			errorMsg += ": " + errMsg
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, errorMsg))
	}
}

func (h *BotHandler) handleProcessAnswer(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, session *models.GameSession) {
	hardLevel2Score := map[string]int{
		"простой": 1,
		"средний": 2,
		"сложный": 4,
	}

	userAnswer := msg.Text

	correctAnswerIndex := session.Questions[session.CurrentQuestionIndex].CorrectIndex
	correctAnswer := session.Questions[session.CurrentQuestionIndex].Options[correctAnswerIndex]

	if userAnswer == correctAnswer {
		session.Score += hardLevel2Score[session.Questions[session.CurrentQuestionIndex].Level]
	}
	session.CurrentQuestionIndex++

	if err := h.redis.UpdateGameSession(context.Background(), session); err != nil {
		log.Printf("Error update play session: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обновления игровой сессии"))
		return
	}

	if session.CurrentQuestionIndex == len(session.Questions) {
		h.redis.DeleteGameSession(context.Background(), msg.From.ID, session.SessionID)
		// Добавить запись в БД

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🎯 Выбрать тему викторины"),
			),
		)

		new_msg := tgbotapi.NewMessage(msg.Chat.ID, "Спасибо за игру!\nВаш результат: "+strconv.Itoa(session.Score)) // TODO: Добавить обработку времени
		new_msg.ReplyMarkup = keyboard

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

	new_msg := tgbotapi.NewMessage(msg.Chat.ID, session.Questions[session.CurrentQuestionIndex].Level+"вопрос:\n"+session.Questions[session.CurrentQuestionIndex].Text)
	new_msg.ReplyMarkup = keyboard

	bot.Send(new_msg)
}
