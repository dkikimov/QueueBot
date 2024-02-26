package main

import (
	"os"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/logger"
	"QueueBot/internal/storage/sqlite"
	"QueueBot/internal/telegram"
)

func main() {
	botToken, exists := os.LookupEnv("BOT_TOKEN")
	if exists == false {
		logger.Fatalf("Bot token is not provided")
	}

	tgBot, err := tgBotApi.NewBotAPI(botToken)
	if err != nil {
		logger.Fatalf("Couldn't initialize bot with error: %s", err.Error())
	}

	if os.Getenv("DEBUG") == "true" {
		tgBot.Debug = true
	}

	storage, err := sqlite.NewDatabase()
	if err != nil {
		logger.Fatalf("Couldn't initialize storage: %s", err)
	}

	defer func(storage *sqlite.SQLite) {
		err := storage.Close()
		if err != nil {
			logger.Fatalf("couldn't close storage")
		}
	}(storage)

	bot := telegram.NewAppBot(tgBot, storage)

	server := telegram.NewBotServer(bot)

	updateConfig := tgBotApi.NewUpdate(0)
	updateConfig.Timeout = 30

	errChan := make(chan error)

	go server.Listen(updateConfig, errChan)

	for err := range errChan {
		if err != nil {
			logger.Println(err.Error())
		}
	}
}
