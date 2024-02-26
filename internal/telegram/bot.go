package telegram

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/models"
	"QueueBot/internal/steps"
	"QueueBot/internal/storage"
	"QueueBot/internal/telegram/messages"
)

type Bot struct {
	TgBot   *tgbotapi.BotAPI
	Storage storage.Storage
}

func NewAppBot(tgBot *tgbotapi.BotAPI, storage storage.Storage) *Bot {
	return &Bot{TgBot: tgBot, Storage: storage}
}

func (b Bot) SendHelloMessage(message *tgbotapi.Message) error {
	if err := b.Storage.CreateUser(message.From.ID); err != nil {
		return fmt.Errorf("couldn't create user in db with error: %s", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, HelloMessage)
	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send hello message in telegram with error: %s", err)
	}

	return nil
}

func (b Bot) SendMessageToCreateQueue(message *tgbotapi.Message) error {
	if err := b.Storage.SetUserCurrentStep(message.From.ID, steps.EnteringDescription); err != nil {
		return fmt.Errorf("couldn't set user current step in db with error: %s", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, CreateQueueMessage)
	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send create queue message with error: %s", err)
	}

	return nil
}

func (b Bot) SendForwardToMessage(message *tgbotapi.Message) error {
	msg := messages.GetForwardMessage(message.Chat.ID, message.Text)
	if _, err := b.TgBot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send forward to message in telegram with error: %s", err)
	}
	return nil
}

func (b Bot) CreateQueue(messageId string, description string) error {
	if err := b.Storage.CreateQueue(messageId, description); err != nil {
		return fmt.Errorf("couldn't create queue with error: %s", err)
	}

	slog.Info("Queue created successfully", "messageId", messageId, "description", description)

	return nil
}

func (b Bot) LogInOurOut(callbackQuery *tgbotapi.CallbackQuery) error {
	description, err := b.Storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}

	if err = b.Storage.LogInOurOutQueue(
		callbackQuery.InlineMessageID,
		models.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		return fmt.Errorf("couldn't add user to queue with error: %s", err)
	}

	users, err := b.Storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	updatedMessage := messages.GetUpdatedQueueMessage(callbackQuery.InlineMessageID, description, users)

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message with error: %s", err)
	}

	slog.Info("Logged in/out", "messageId", callbackQuery.InlineMessageID, "userId", callbackQuery.From.ID)

	return nil
}

func (b Bot) Start(callbackQuery *tgbotapi.CallbackQuery, isShuffled bool) error {
	err, wasUpdated := b.Storage.StartQueue(callbackQuery.InlineMessageID, isShuffled)
	if err != nil {
		return fmt.Errorf("couldn't start queue with error: %s", err)
	}

	if !wasUpdated {
		return nil
	}

	slog.Info("Started queue", "messageId", callbackQuery.InlineMessageID, "isShuffled", isShuffled)

	return b.sendQueueAfterStartMessage(callbackQuery, 0)
}

func (b Bot) Next(callbackQuery *tgbotapi.CallbackQuery) error {
	err, currentPersonIndex := b.Storage.IncrementCurrentPerson(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't increment current person in queue %s with error: %s", callbackQuery.InlineMessageID, err)
	}

	slog.Info("Set next person", "messageId", callbackQuery.InlineMessageID, "currentPersonIndex", currentPersonIndex)

	return b.sendQueueAfterStartMessage(callbackQuery, currentPersonIndex)
}

func (b Bot) GoToMenu(callbackQuery *tgbotapi.CallbackQuery) error {
	description, err := b.Storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}
	if err = b.Storage.GoToMenu(callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't reset current person. Queue id: %s", callbackQuery.InlineMessageID)
	}

	users, err := b.Storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	updatedMessage := messages.GetQueueMessage(callbackQuery.InlineMessageID, users, description)
	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't go to menu with error: %s", err)
	}

	slog.Info("Went to menu", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b Bot) FinishQueue(callbackQuery *tgbotapi.CallbackQuery) error {
	updatedMessage := messages.GetFinishedMessage(callbackQuery.InlineMessageID)
	_, err := b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't send finish queue with error: %s", err)
	}

	if err = b.Storage.FinishQueueDeleteParticipants(callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't finish queue with error: %s", err)
	}

	slog.Info("Finished queue", "messageId", callbackQuery.InlineMessageID)

	return nil
}

func (b Bot) sendQueueAfterStartMessage(callbackQuery *tgbotapi.CallbackQuery, currentPersonIndex int) error {
	description, err := b.Storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}

	users, err := b.Storage.GetUsersInQueueCheckShuffle(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	var updatedMessage tgbotapi.EditMessageTextConfig
	if currentPersonIndex == len(users) {
		updatedMessage = messages.GetEndQueueMessage(callbackQuery.InlineMessageID)
	} else {
		updatedMessage = messages.GetQueueAfterStartMessage(callbackQuery.InlineMessageID, description, users, currentPersonIndex)
	}

	_, err = b.TgBot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message after starting queue with error: %s", err.Error())
	}

	return nil
}
