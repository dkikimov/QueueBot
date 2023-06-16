package queue

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/ui"
	"QueueBot/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create(messageId string, description string, storage storage.Storage) {
	if err := storage.CreateQueue(messageId, description); err != nil {
		logger.Fatalf("Couldn't create queue with error: %s", err.Error())
	}
}

func LogInOurOut(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get description of queue with error: %s", err.Error())
	}

	if err = storage.LogInOurOutQueue(
		callbackQuery.InlineMessageID,
		user.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		logger.Fatalf("Couldn't add user to queue with error: %s", err.Error())
	}

	users, err := storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get users in queue with error: %s", err.Error())
	}

	updatedMessage := ui.GetUpdatedQueueMessage(callbackQuery.InlineMessageID, description, users)

	_, err = bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't update message with error: %s", err.Error())
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, constants.LogInOurOutAlert)
	if _, err = bot.Request(callback); err != nil {
		logger.Panicf("Couldn't process callback with error: %s", err.Error())
	}
}

func Start(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage, isShuffled bool) {
	err, wasUpdated := storage.StartQueue(callbackQuery.InlineMessageID, isShuffled)
	if err != nil {
		logger.Fatalf("Couldn't start queue with error: %s", err.Error())
	}

	if !wasUpdated {
		return
	}

	sendQueueAfterStartMessage(callbackQuery, bot, storage, 0)
}

func Next(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	err, currentPersonIndex := storage.IncrementCurrentPerson(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't increment current person in queue %s with error: %s", callbackQuery.InlineMessageID, err.Error())
	}
	sendQueueAfterStartMessage(callbackQuery, bot, storage, currentPersonIndex)
}

func GoToMenu(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err = storage.GoToMenu(callbackQuery.InlineMessageID); err != nil {
		logger.Fatalf("Couldn't reset current person. Queue id: %s", callbackQuery.InlineMessageID)
	}

	users, err := storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get users in queue with error: %s", err.Error())
	}

	updatedMessage := ui.GetQueueMessage(callbackQuery.InlineMessageID, users, description)
	_, err = bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't go to menu with error: %s", err.Error())
	}
}

func FinishQueue(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	updatedMessage := ui.GetFinishedMessage(callbackQuery.InlineMessageID)
	_, err := bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't finish queue with error: %s", err.Error())
	}
}

func sendQueueAfterStartMessage(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage, currentPersonIndex int) {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get description of queue with error: %s", err.Error())
	}

	users, err := storage.GetUsersInQueueCheckShuffle(callbackQuery.InlineMessageID)
	if err != nil {
		logger.Fatalf("Couldn't get users in queue with error: %s", err.Error())
	}

	var updatedMessage tgbotapi.EditMessageTextConfig
	if currentPersonIndex == len(users) {
		updatedMessage = ui.GetEndQueueMessage(callbackQuery.InlineMessageID)
	} else {
		updatedMessage = ui.GetQueueAfterStartMessage(callbackQuery.InlineMessageID, description, users, currentPersonIndex)
	}

	_, err = bot.Request(updatedMessage)
	if err != nil {
		logger.Fatalf("Couldn't update message after starting queue with error: %s", err.Error())
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, constants.NextData)
	if _, err = bot.Request(callback); err != nil {
		logger.Panicf("Couldn't process next_data callback with error: %s", err.Error())
	}
}
