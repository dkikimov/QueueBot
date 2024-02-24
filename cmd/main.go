package main

import (
	"log"
	"os"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"QueueBot/internal/logger"
	"QueueBot/internal/storage/sqlite"
	"QueueBot/internal/telegram"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Panicf("Couldn't open .env file %s", err)
	}

	tgBot, err := tgBotApi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logger.Panicf("Couldn't initialize bot with error: %s", err.Error())
	}

	if os.Getenv("DEBUG") == "true" {
		tgBot.Debug = true
	}

	storage := sqlite.NewDatabase()
	defer func(storage *sqlite.SQLite) {
		err := storage.Close()
		if err != nil {
			log.Fatalf("couldn't close storage")
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
