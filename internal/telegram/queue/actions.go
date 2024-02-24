package queue

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/models"
	"QueueBot/internal/storage"
	"QueueBot/internal/telegram/ui"
)

func Create(messageId string, description string, storage storage.Storage) error {
	if err := storage.CreateQueue(messageId, description); err != nil {
		return fmt.Errorf("couldn't create queue with error: %s", err)
	}
	return nil
}

func LogInOurOut(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}

	if err = storage.LogInOurOutQueue(
		callbackQuery.InlineMessageID,
		models.New(callbackQuery.From.ID, callbackQuery.From.LastName, callbackQuery.From.FirstName),
	); err != nil {
		return fmt.Errorf("couldn't add user to queue with error: %s", err)
	}

	users, err := storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	updatedMessage := ui.GetUpdatedQueueMessage(callbackQuery.InlineMessageID, description, users)

	_, err = bot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message with error: %s", err)
	}

	return nil
}

func Start(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage, isShuffled bool) error {
	err, wasUpdated := storage.StartQueue(callbackQuery.InlineMessageID, isShuffled)
	if err != nil {
		return fmt.Errorf("couldn't start queue with error: %s", err)
	}

	if !wasUpdated {
		return nil
	}

	return sendQueueAfterStartMessage(callbackQuery, bot, storage, 0)
}

func Next(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	err, currentPersonIndex := storage.IncrementCurrentPerson(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't increment current person in queue %s with error: %s", callbackQuery.InlineMessageID, err)

	}

	return sendQueueAfterStartMessage(callbackQuery, bot, storage, currentPersonIndex)
}

func GoToMenu(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}
	if err = storage.GoToMenu(callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't reset current person. Queue id: %s", callbackQuery.InlineMessageID)
	}

	users, err := storage.GetUsersInQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	updatedMessage := ui.GetQueueMessage(callbackQuery.InlineMessageID, users, description)
	_, err = bot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't go to menu with error: %s", err)
	}
	return nil
}

func FinishQueue(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	updatedMessage := ui.GetFinishedMessage(callbackQuery.InlineMessageID)
	_, err := bot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't send finish queue with error: %s", err)
	}

	if err = storage.FinishQueueDeleteParticipants(callbackQuery.InlineMessageID); err != nil {
		return fmt.Errorf("couldn't finish queue with error: %s", err)
	}

	return nil
}

func sendQueueAfterStartMessage(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage, currentPersonIndex int) error {
	description, err := storage.GetDescriptionOfQueue(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get description of queue with error: %s", err)
	}

	users, err := storage.GetUsersInQueueCheckShuffle(callbackQuery.InlineMessageID)
	if err != nil {
		return fmt.Errorf("couldn't get users in queue with error: %s", err)
	}

	var updatedMessage tgbotapi.EditMessageTextConfig
	if currentPersonIndex == len(users) {
		updatedMessage = ui.GetEndQueueMessage(callbackQuery.InlineMessageID)
	} else {
		updatedMessage = ui.GetQueueAfterStartMessage(callbackQuery.InlineMessageID, description, users, currentPersonIndex)
	}

	_, err = bot.Request(updatedMessage)
	if err != nil {
		return fmt.Errorf("couldn't update message after starting queue with error: %s", err.Error())
	}

	return nil
}
