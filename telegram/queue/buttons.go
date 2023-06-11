package queue

import (
	"QueueBot/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKeyboardButtons() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			addToQueueButton(),
		),
	)

	return &keyboard
}

func addToQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.AddToQueueButton, constants.AddToQueueData)
}
