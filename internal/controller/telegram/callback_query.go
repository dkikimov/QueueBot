package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/controller/telegram/client"
)

const (
	ActionCompleted = "Действие выполнено!"
	ActionError     = "Произошла ошибка"
)

func (s BotServer) handleCallbackData(callbackQuery *tgbotapi.CallbackQuery) error {
	switch callbackQuery.Data {
	case client.LogInOurOutData:
		if err := s.bot.LogInOurOut(context.Background(), callbackQuery); err != nil {
			return fmt.Errorf("couldn't login or logout with error: %w", err)
		}
	case client.StartQueueData:
		if err := s.bot.Start(context.Background(), callbackQuery, false); err != nil {
			return fmt.Errorf("couldn't start queue with error: %w", err)
		}
	case client.StartQueueShuffleData:
		if err := s.bot.Start(context.Background(), callbackQuery, true); err != nil {
			return fmt.Errorf("couldn't start queue with shuffle with error: %w", err)
		}
	case client.NextData:
		if err := s.bot.Next(context.Background(), callbackQuery); err != nil {
			return fmt.Errorf("couldn't go to next person with error: %w", err)
		}
	case client.GoToMenuData:
		if err := s.bot.GoToMenu(context.Background(), callbackQuery); err != nil {
			return fmt.Errorf("couldn't go to menu with error: %w", err)
		}
	case client.FinishQueueData:
		if err := s.bot.FinishQueue(context.Background(), callbackQuery); err != nil {
			return fmt.Errorf("couldn't finish queue with error: %w", err)
		}
	}

	return nil
}

func (s BotServer) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	// Сверяемся со скрытыми данными, заложенными в сообщении для определения команды
	startTime := time.Now()
	slog.Debug("Got callback query with data: ", "data", callbackQuery.Data)

	err := s.handleCallbackData(callbackQuery)
	if err != nil {
		slog.Error(
			"Couldn't handle callback query",
			"reason",
			err, "data",
			callbackQuery.Data, "user_id",
			callbackQuery.From.ID,
		)
	}

	var callback tgbotapi.CallbackConfig
	if err != nil {
		callback = tgbotapi.NewCallback(callbackQuery.ID, ActionError)
	} else {
		callback = tgbotapi.NewCallback(callbackQuery.ID, ActionCompleted)
	}

	if _, err = s.bot.TgBot.Request(callback); err != nil {
		return fmt.Errorf("couldn't process next_data callback with error: %w", err)
	}

	slog.Debug("Processed callback query with data: ", "data", callbackQuery.Data, "elapsed", time.Since(startTime).String())

	return nil
}
