package main

import (
	"log"
	"log/slog"
	"os"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/config"
	"QueueBot/internal/controller/telegram"
	"QueueBot/internal/usecase"
	"QueueBot/internal/usecase/storage/sqlite"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Couldn't create config: %s", err)
	}

	var programLevel = new(slog.LevelVar)
	if cfg.IsAppDebug {
		programLevel.Set(slog.LevelDebug)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	botApi, err := tgBotApi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Couldn't initialize bot with error: %s", err.Error())
	}

	botApi.Debug = cfg.IsTelegramDebug

	storage, err := sqlite.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Couldn't initialize storage: %s", err)
	}

	defer func(storage *sqlite.Database) {
		err := storage.Close()
		if err != nil {
			log.Fatalf("couldn't close storage")
		}
	}(storage)

	botUseCase := usecase.NewBotUseCase(storage)
	bot := telegram.NewAppBot(botApi, botUseCase)
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
