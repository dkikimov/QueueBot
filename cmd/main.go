package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/config"
	"QueueBot/internal/controller/telegram"
	"QueueBot/internal/controller/telegram/client"
	"QueueBot/internal/usecase"
	"QueueBot/internal/usecase/storage"
	"QueueBot/internal/usecase/storage/mongodb"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Couldn't create config: %s", err)
	}

	programLevel := new(slog.LevelVar)
	if cfg.IsAppDebug {
		programLevel.Set(slog.LevelDebug)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Couldn't initialize bot with error: %s", err.Error())
	}

	botAPI.Debug = cfg.IsTelegramDebug

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	repository, err := mongodb.NewDatabase(timeoutCtx, cfg.DatabasePath)
	if err != nil {
		log.Panicf("Couldn't initialize repository: %s", err)
	}

	defer func(storage storage.Storage) {
		err := storage.Close()
		if err != nil {
			log.Fatalf("couldn't close repository")
		}
	}(repository)

	botUseCase := usecase.NewBotUseCase(repository)
	bot := client.NewTelegramBot(botAPI, botUseCase)
	server := telegram.NewBotServer(bot)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	errChan := make(chan error)

	go server.Listen(updateConfig, errChan)

	for err := range errChan {
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
