package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"quiz-bot-client/models"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем токен бота
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in .env file")
	}

	adminAPIURL := os.Getenv("ADMIN_API_URL")
	if adminAPIURL == "" {
		panic("Требуется ввессти URL для админского микросервиса")
	}
	playerAPIURL := os.Getenv("PLAYER_API_URL")
	if playerAPIURL == "" {
		panic("Требуется ввессти URL для игрового микросервиса")
	}

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Конфигурируем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем сообщения
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch {
		case update.Message.Text == "/start":
			handleStartCommand(bot, update.Message, adminAPIURL)
		case strings.HasSuffix(strings.ToLower(update.Message.Text), "выбрать тему викторины"):
			handleChooseThemeCommand(bot, update.Message, playerAPIURL)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понял. Выберите с клавиатуры или пропишите /start"))
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
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки ответа"))
			return
		}

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🎯 Выбрать тему викторины"),
				// tgbotapi.NewKeyboardButton("📊 Рейтинги"),
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
	resp, err := http.Get(playerAPIURL + "/api/users/topics")

	if err != nil {
		log.Printf("Error calling player API: %v", err)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка соединения с сервером"+err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result models.AllTopicsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Ошибка обработки ответа"))
			return
		}

		// keyboard := tgbotapi.NewReplyKeyboard(
		// 	tgbotapi.NewKeyboardButtonRow(
		// 		tgbotapi.NewKeyboardButton("🎯 Выбрать тему викторины"),
		// 		tgbotapi.NewKeyboardButton("📊 Рейтинги"),
		// 	),
		// )
		msg_text := "Список имеющихся тем викторины:"
		for _, topic := range result.Topics {
			msg_text = msg_text + "\n" + topic.Title
		}

		msg := tgbotapi.NewMessage(msg.Chat.ID, msg_text)
		// msg.ReplyMarkup = keyboard

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
