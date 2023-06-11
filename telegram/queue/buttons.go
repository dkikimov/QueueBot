package queue

import (
	"QueueBot/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetKeyboardButtons(isUserInQueue bool) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			getPersonalQueueButton(isUserInQueue),
		),
	)

	return keyboard
}

func getPersonalQueueButton(isUserInQueue bool) tgbotapi.InlineKeyboardButton {
	if isUserInQueue {
		return exitFromQueueButton()
	} else {
		return addToQueueButton()
	}
}

func addToQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.AddToQueueButton, constants.AddToQueueData)
}

func exitFromQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.ExitFromQueueButton, constants.ExitFromQueueData)
}
