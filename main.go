package main

import (
	"QueueBot/logger"
	"QueueBot/storage/sqlite"
	"QueueBot/telegram"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Couldn't open .env file %s", err)
	}

	bot, err := tgBotApi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	storage := sqlite.NewDatabase()

	if err != nil {
		logger.Fatalf("Couldn't initialize bot with error: %s", err.Error())
	}
	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
	}

	updateConfig := tgBotApi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	logger.Printf("Bot started")

	for update := range updates {
		switch {
		case update.Message != nil:
			telegram.HandleMessage(&update, bot)
		case update.InlineQuery != nil:
			telegram.HandleInlineQuery(update.InlineQuery, bot, storage)
		case update.CallbackQuery != nil:
			telegram.HandleCallbackQuery(update.CallbackQuery, bot, storage)
		}

	}
}
