package telegram

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/telegram/queue"
	"QueueBot/telegram/steps"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, storage storage.Storage) {
	switch update.Message.Command() {
	case constants.StartCommand:
		SendHelloMessage(update.Message, bot, storage)
		return
	case constants.CreateQueueCommand:
		SendMessageToCreateQueue(update.Message, bot, storage)
		return
	}

	currentStep, err := storage.GetUserCurrentStep(update.Message.From.ID)
	if err != nil {
		logger.Fatalf("Couldn't get current user step with error: %s", err.Error())
	}

	switch currentStep {
	case steps.Menu:
		SendHelloMessage(update.Message, bot, storage)
		break
	case steps.EnteringDescription:
		queue.Create(update.Message, bot, storage)
	}
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	switch callbackQuery.Data {
	case constants.AddToQueueData:
		queue.AddTo(callbackQuery, bot, storage)
	}
}
