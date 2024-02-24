package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/constants"
	"QueueBot/internal/storage"
	"QueueBot/internal/telegram/steps"
	"QueueBot/internal/telegram/ui"
)

func SendHelloMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	if err := storage.CreateUser(message.From.ID); err != nil {
		return fmt.Errorf("couldn't create user in db with error: %s", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, constants.HelloMessage)
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send hello message in telegram with error: %s", err)
	}

	return nil
}

func SendMessageToCreateQueue(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) error {
	if err := storage.SetUserCurrentStep(message.From.ID, steps.EnteringDescription); err != nil {
		return fmt.Errorf("couldn't set user current step in db with error: %s", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, constants.CreateQueueMessage)
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send create queue message with error: %s", err)
	}

	return nil
}

func SendForwardToMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) error {
	msg := ui.GetForwardMessage(message.Chat.ID, message.Text)
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("couldn't send forward to message in telegram with error: %s", err)
	}
	return nil
}
