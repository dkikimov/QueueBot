package telegram

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/telegram/steps"
	"QueueBot/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendHelloMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) {
	if err := storage.CreateUser(message.From.ID); err != nil {
		logger.Fatalf("Couldn't create user with error: %s", err.Error())
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, constants.HelloMessage)
	if _, err := bot.Send(msg); err != nil {
		logger.Fatalf("Couldn't send hello message with error: %s", err.Error())
	}
}

func SendMessageToCreateQueue(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) {
	if err := storage.SetUserCurrentStep(message.From.ID, steps.EnteringDescription); err != nil {
		logger.Fatalf("Couldn't set user current step with error: %s", err.Error())
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, constants.CreateQueueMessage)
	if _, err := bot.Send(msg); err != nil {
		logger.Fatalf("Couldn't send create queue message with error: %s", err.Error())
	}
}

func SendForwardToMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	msg := ui.GetForwardMessage(message.Chat.ID, message.Text)
	if _, err := bot.Send(msg); err != nil {
		logger.Fatalf("Couldn't send forward to message with error: %s", err.Error())
	}
}
