package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/controller/telegram/messages"
	"QueueBot/internal/entity"
	"QueueBot/internal/usecase"
)

type Bot struct {
	TgBot *tgbotapi.BotAPI
	u     usecase.Bot
}

func NewAppBot(tgBot *tgbotapi.BotAPI, u usecase.Bot) *Bot {
	return &Bot{TgBot: tgBot, u: u}
}

func (b Bot) SendHelloMessage(ctx context.Context, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, HelloMessage)

	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send hello message in telegram with error: %s", err)
	}

	return nil
}

func (b Bot) SendForwardMessageButton(ctx context.Context, message *tgbotapi.Message) error {
	msg := messages.GetForwardMessage(message.Chat.ID, message.Text)
	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send forward to message in telegram with error: %s", err)
	}
	return nil
}

func (b Bot) CreateQueue(ctx context.Context, messageId string, description string) error {
	if err := b.u.CreateQueue(ctx, messageId, description); err != nil {
		return fmt.Errorf("couldn't create queue with error: %s", err)
	}

	slog.Info("Queue created successfully", "messageId", messageId, "description", description)

	return nil
}

func (b Bot) LogInOurOut(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	startTime := time.Now()

	if err := b.u.LogInOutToQueue(
		ctx,
		callbackQuery.InlineMessageID,
		entity.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		return fmt.Errorf("couldn't add user to queue with error: %s", err)
	}

	slog.Debug("Logged in/out locally", "elapsed", time.Now().Sub(startTime).String())

	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %s", err)
	}

	slog.Debug("Got queue", "elapsed", time.Now().Sub(startTime).String())

	updatedMessage := messages.GetUpdatedQueueMessage(callbackQuery.InlineMessageID, queue.Description, queue.Users)

	slog.Debug("Got updated queue message", "elapsed", time.Now().Sub(startTime).String())

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message with error: %s", err)
	}

	slog.Debug(
		"Logged in/out and sent updated message",
		"messageId", callbackQuery.InlineMessageID,
		"userId", callbackQuery.From.ID,
		"elapsed", time.Now().Sub(startTime).String(),
	)

	return nil
}

func (b Bot) Start(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery, isShuffled bool) error {
	err := b.u.StartQueue(ctx, callbackQuery.InlineMessageID, isShuffled)
	if err != nil {
		return fmt.Errorf("couldn't start queue with error: %s", err)
	}

	slog.Info("Started queue", "messageId", callbackQuery.InlineMessageID, "isShuffled", isShuffled)

	return b.sendQueueStatusMessage(ctx, callbackQuery)
}

func (b Bot) Next(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	err := b.u.SetNextPersonToQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't increment current person in queue %s with error: %s", callbackQuery.InlineMessageID, err)
	}

	slog.Info("Set next person", "messageId", callbackQuery.InlineMessageID)

	return b.sendQueueStatusMessage(ctx, callbackQuery)
}

func (b Bot) GoToMenu(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %s", err)
	}

	updatedMessage := messages.GetQueueMessage(callbackQuery.InlineMessageID, queue.Users, queue.Description)
	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't go to menu with error: %s", err)
	}

	slog.Info("Went to menu", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b Bot) FinishQueue(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	updatedMessage := messages.GetFinishedMessage(callbackQuery.InlineMessageID)
	_, err := b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't send finish queue with error: %s", err)
	}

	if err = b.u.FinishQueue(ctx, callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't finish queue with error: %s", err)
	}

	slog.Info("Finished queue", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b Bot) sendQueueStatusMessage(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {
	queue, err := b.u.GetQueue(ctx, callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue with error: %s", err)
	}

	var updatedMessage tgbotapi.EditMessageTextConfig
	if queue.CurrentPersonIdx == len(queue.Users) {
		updatedMessage = messages.GetEndQueueMessage(callbackQuery.InlineMessageID)
	} else {
		updatedMessage = messages.GetQueueAfterStartMessage(callbackQuery.InlineMessageID, queue.Description, queue.Users, queue.CurrentPersonIdx)
	}

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message after starting queue with error: %s", err.Error())
	}

	return nil
}
