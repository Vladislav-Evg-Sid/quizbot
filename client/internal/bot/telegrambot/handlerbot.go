package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) HandleUpdate(tgBot *TelegramBot, update tgbotapi.Update, cfg *config.Config) {
	if update.Message == nil {
		return
	}
	ctx := context.Background()

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.Text == "/start" {
		h.handler.HandlerStartCommand(tgBot.Bot, update.Message, cfg.Network.AdminREST)
		return
	}

	session, err := h.getCurrentUserSession(ctx, update.Message.From.ID, tgBot, cfg.Network.AdminREST)
	if err != nil {
		log.Printf("Redis error: %v", err)
		tgBot.Bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "❌ Ошибка сервера"))
	}

	switch {
	case session.CurrentStep == "root":
		switch {
		case strings.HasSuffix(strings.ToLower(update.Message.Text), "посмотреть список тем викторины"):
			if slices.Contains(session.Permissions, "canReadTopics") {
				h.handler.HandleChooseThemeCommand(ctx, tgBot.Bot, update.Message, cfg.Network.PlayerREST, session)
			}
		default:
			tgBot.Bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Выбирите действие из кнопок"))
		}
	case session.CurrentStep == "selectThema":
		h.handler.HandleSelectThemsProcessing(ctx, update.Message.Text, update.Message.From.ID, session, tgBot.Bot)
	case strings.HasPrefix(session.CurrentStep, "thema:"):
		switch {
		case update.Message.Text == "Посмотреть список вопросов темы":
			if slices.Contains(session.Permissions, "canReadQuestions") {
				// TODO: Добавить получение всех вопросов темы
			}
		case update.Message.Text == "Начать игру":
			if slices.Contains(session.Permissions, "canPlayInQuiz") {
				h.handler.HandleGetQuestionsForQuiz(ctx, tgBot.Bot, session.CurrentStep[6:], update.Message.From.ID, session, cfg.Network.PlayerREST)
			}
		case update.Message.Text == "Посмотреть рейтинг":
			if slices.Contains(session.Permissions, "canReadRatings") {
				// TODO: Добавить получение рейтингов
			}
		case update.Message.Text == "Обновить название темы":
			if slices.Contains(session.Permissions, "canUpdateTopics") {
				// TODO: Добавить обновление названия темы
			}
		case update.Message.Text == "Удалить тему":
			if slices.Contains(session.Permissions, "canDeleteTopics") {
				// TODO: Добавить возможность удаления темы
			}
		case update.Message.Text == "Добавить новые вопросы":
			if slices.Contains(session.Permissions, "canCreateQuestions") {
				// TODO: Добавить добавление вопросов
			}
		default:
			tgBot.Bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Выбирите действие из кнопок"))
		}
	case session.CurrentStep == "gameProcess":
		h.handler.HandleProcessAnswer(ctx, tgBot.Bot, update.Message.Text, update.Message.From.ID, cfg.Network.PlayerREST)
	}
}

func (h *BotHandler) getCurrentUserSession(ctx context.Context, tg_id int64, tgBot *TelegramBot, adminAPIURL string) (*models.UserSession, error) {
	session, err := h.redis.GetUserSession(ctx, tg_id)
	if err != nil {
		log.Printf("Redis error: %v", err)
		tgBot.Bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка сервера"))
		return nil, err
	}
	if session == nil {
		resp, err := http.Get(fmt.Sprintf("%s/api/users/%d/permissions", adminAPIURL, tg_id))
		if err != nil {
			log.Printf("Error calling admin API: %v", err)
			tgBot.Bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка соединения с сервером"))
			return nil, err
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var result models.Permissions
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Printf("Error processing servers`s answers: %v", err)
				tgBot.Bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка обработки ответа"))
				return nil, err
			}
			session, err := h.redis.CreateUserSession(ctx, tg_id, result.Permissions, "root")
			if err != nil {
				log.Printf("Redis error: %v", err)
				tgBot.Bot.Send(tgbotapi.NewMessage(tg_id, "❌ Ошибка сервера"))
				return nil, err
			}
			return session, nil
		} else {
			var errorResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errorResp)
			errorMsg := "❌ Ошибка регистрации"
			if errMsg, ok := errorResp["error"].(string); ok {
				errorMsg += ": " + errMsg
			}
			tgBot.Bot.Send(tgbotapi.NewMessage(tg_id, errorMsg))
		}
	}
	return session, nil
}
