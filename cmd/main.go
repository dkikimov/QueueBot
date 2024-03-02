package main

import (
	"log"
	"log/slog"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/config"
	"QueueBot/internal/storage/sqlite"
	"QueueBot/internal/telegram"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Couldn't create config: %s", err)
	}

	tgBot, err := tgBotApi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Couldn't initialize bot with error: %s", err.Error())
	}

	tgBot.Debug = cfg.IsDebug

	storage, err := sqlite.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Couldn't initialize storage: %s", err)
	}

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
			slog.Error(err.Error())
		}
	}
}
