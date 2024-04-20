package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/apperrors"
	"QueueBot/internal/entity"
	"QueueBot/internal/usecase"
)

const HelloMessage = `Привет! Я бот, предназначенный для создания очередей. 
Введи описание своей очереди, а я тебе ее создам`

type TelegramBot struct {
	TgBot *tgbotapi.BotAPI
	u     usecase.Bot
}

func NewTelegramBot(tgBot *tgbotapi.BotAPI, u usecase.Bot) *TelegramBot {
	return &TelegramBot{TgBot: tgBot, u: u}
}

func (b TelegramBot) SendHelloMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, HelloMessage)

	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send hello message in telegram with error: %w", err)
	}

	return nil
}

func (b TelegramBot) SendForwardMessageButton(message *tgbotapi.Message) error {
	msg := GetForwardMessage(message.Chat.ID, message.Text)
	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send forward to message in telegram with error: %w", err)
	}

	return nil
}

func (b TelegramBot) CreateQueue(ctx context.Context, messageID string, description string) error {
	if err := b.u.CreateQueue(ctx, messageID, description); err != nil {
		return fmt.Errorf("couldn't create queue with error: %w", err)
	}

	slog.Info("Queue created successfully", "messageID", messageID, "description", description)

	return nil
}

func (b TelegramBot) LogInOurOut(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	startTime := time.Now()

	err := b.u.LogInOutToQueue(
		ctx,
		callbackQuery.InlineMessageID,
		entity.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	)

	if err != nil {
		var callbackError *apperrors.CallbackError
		if errors.As(err, &callbackError) {
			callback := tgbotapi.NewCallback(callbackQuery.ID, callbackError.Message)
			_, err := b.TgBot.Request(callback)
			if err != nil {
				return fmt.Errorf("couldn't send callback with message %s error: %w", callbackError.Message, err)
			}

		} else {
			return fmt.Errorf("couldn't add user to queue with error: %w", err)
		}
	}

	slog.Debug("Logged in/out locally", "elapsed", time.Since(startTime).String())

	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %w", err)
	}

	slog.Debug("Got queue", "elapsed", time.Since(startTime).String())

	updatedMessage := GetUpdatedQueueMessage(callbackQuery.InlineMessageID, queue.Description, queue.Users)

	slog.Debug("Got updated queue message", "elapsed", time.Since(startTime).String())

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message with error: %w", err)
	}

	slog.Debug(
		"Logged in/out and sent updated message",
		"messageId", callbackQuery.InlineMessageID,
		"userId", callbackQuery.From.ID,
		"elapsed", time.Since(startTime).String(),
	)

	return nil
}

func (b TelegramBot) Start(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery, isShuffled bool) error {
	err := b.u.StartQueue(ctx, callbackQuery.InlineMessageID, isShuffled)
	if err != nil {
		return fmt.Errorf("couldn't start queue with error: %w", err)
	}

	slog.Info("Started queue", "messageId", callbackQuery.InlineMessageID, "isShuffled", isShuffled)

	return b.sendQueueStatusMessage(ctx, callbackQuery)
}

func (b TelegramBot) Next(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	err := b.u.SetNextPersonToQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't increment current person in queue %s with error: %w", callbackQuery.InlineMessageID, err)
	}

	slog.Info("Set next person", "messageId", callbackQuery.InlineMessageID)

	return b.sendQueueStatusMessage(ctx, callbackQuery)
}

func (b TelegramBot) GoToMenu(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %w", err)
	}

	updatedMessage := GetQueueMessage(callbackQuery.InlineMessageID, queue.Users, queue.Description)
	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't go to menu with error: %w", err)
	}

	slog.Info("Went to menu", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b TelegramBot) FinishQueue(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	updatedMessage := GetFinishedMessage(callbackQuery.InlineMessageID)
	_, err := b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't send finish queue with error: %w", err)
	}

	if err = b.u.FinishQueue(ctx, callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't finish queue with error: %w", err)
	}

	slog.Info("Finished queue", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b TelegramBot) sendQueueStatusMessage(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %w", err)
	}

	var updatedMessage tgbotapi.EditMessageTextConfig
	if queue.CurrentPersonIdx == len(queue.Users) {
		updatedMessage = GetEndQueueMessage(callbackQuery.InlineMessageID)
	} else {
		updatedMessage = GetQueueAfterStartMessage(callbackQuery.InlineMessageID, queue.Description, queue.Users, queue.CurrentPersonIdx)
	}

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message after starting queue with error: %w", err)
	}

	return nil
}
