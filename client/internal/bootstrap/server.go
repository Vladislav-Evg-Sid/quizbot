package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vladislav-Evg-Sid/quizbot/client/config"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/bot/processors"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/bot/telegrambot"
	"github.com/Vladislav-Evg-Sid/quizbot/client/internal/storage/redisstorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AppRun(cfg *config.Config, clientRedis *redisstorage.RedisStorage, tgBot *telegrambot.TelegramBot) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := processors.NewClientBotHandler(clientRedis)

	botHandler := telegrambot.NewBotHandler(clientRedis, handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() {
		log.Println("Запускам Telegram бота...")
		done <- runBot(ctx, cfg, tgBot, botHandler)
	}()

	log.Println("🚀 Бот запущен. Нажмите Ctrl+C для остановки.")

	select {
	case <-sigChan:
		log.Println("Получен сигнал остановки")
		cancel()
	case err := <-done:
		if err != nil {
			log.Printf("Бот завершился с ошибкой: %v", err)
		}
	}

	log.Println("⏳ Ожидаем завершения работы...")
	// Даём время на завершение операций
	time.Sleep(2 * time.Second)
	log.Println("✅ Бот остановлен")
}

func runBot(ctx context.Context, cfg *config.Config, telegrambot *telegrambot.TelegramBot, handler *telegrambot.BotHandler) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := telegrambot.Bot.GetUpdatesChan(u)

	log.Printf("Бот @%s запущен", telegrambot.Bot.Self.UserName)

	for {
		select {
		case <-ctx.Done():
			log.Println("Получен сигнал остановки бота")
			return nil
		case update := <-updates:
			go handler.HandleUpdate(telegrambot, update, cfg)
		}
	}
}
