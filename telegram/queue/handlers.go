package queue

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/storage/user"
	"QueueBot/telegram/steps"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) {
	if err := storage.SetUserCurrentStep(message.From.ID, steps.Menu); err != nil {
		logger.Fatalf("Couldn't set user current step with error: %s", err.Error())
	}

	//TODO: Also send explanation message
	answer := GetQueueMessage(message.Chat.ID, message.Text, nil)
	sentMessage, err := bot.Send(answer)
	if err != nil {
		logger.Fatalf("Couldn't send queue message with error: %s", err.Error())
	}

	if err = storage.CreateQueue(sentMessage.MessageID, message.Text); err != nil {
		logger.Fatalf("Couldn't create queue with error: %s", err.Error())
	}
}

func AddTo(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.Message.MessageID)
	if err != nil {
		logger.Fatalf("Couldn't get description of queue with error: %s", err.Error())
	}

	if err := storage.AddUserToQueue(
		callbackQuery.Message.MessageID,
		user.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		logger.Fatalf("Couldn't add user to queue with error: %s", err.Error())
	}

	users, err := storage.GetUsersInQueue(callbackQuery.Message.MessageID)
	if err != nil {
		logger.Fatalf("Couldn't get users in queue with error: %s", err.Error())
	}

	updatedMessage := GetUpdatedQueueMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, description, users, true)

	_, err = bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't update message with error: %s", err.Error())
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, constants.AddedToQueueAlert)
	if _, err := bot.Request(callback); err != nil {
		logger.Panicf("Couldn't process callback with error: %s", err.Error())
	}
}
