package queue

import (
	"QueueBot/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			logInOurOutQueueButton(),
		),
	)

	return keyboard
}

func logInOurOutQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.LogInOurOutButton, constants.LogInOurOutData)
}
