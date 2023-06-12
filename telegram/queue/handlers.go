package queue

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/storage/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create(messageId string, description string, bot *tgbotapi.BotAPI, storage storage.Storage) {
	//if err := storage.SetUserCurrentStep(message.From.ID, steps.Menu); err != nil {
	//	logger.Fatalf("Couldn't set user current step with error: %s", err.Error())
	//}

	//TODO: Also send explanation message

	if err := storage.CreateQueue(messageId, description); err != nil {
		logger.Fatalf("Couldn't create queue with error: %s", err.Error())
	}
}

func AddTo(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get description of queue with error: %s", err.Error())
	}

	if err = storage.AddUserToQueue(
		callbackQuery.InlineMessageID,
		user.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		logger.Fatalf("Couldn't add user to queue with error: %s", err.Error())
	}

	users, err := storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get users in queue with error: %s", err.Error())
	}

	updatedMessage := GetUpdatedQueueMessage(callbackQuery.InlineMessageID, description, users, true)

	_, err = bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't update message with error: %s", err.Error())
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, constants.AddedToQueueAlert)
	if _, err := bot.Request(callback); err != nil {
		logger.Panicf("Couldn't process callback with error: %s", err.Error())
	}
}
