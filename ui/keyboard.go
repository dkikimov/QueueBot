package ui

import (
	"QueueBot/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetBeforeStartKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			logInOurOutQueueButton(),
		),
		tgbotapi.NewInlineKeyboardRow(
			startQueueButton(),
		),
	)

	return keyboard
}

func GetAfterStartKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			nextButton(),
		),
	)

	return keyboard
}

func GetEndedQueueKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			goToMenuButton(),
		),
		tgbotapi.NewInlineKeyboardRow(
			endQueueButton(),
		),
	)

	return keyboard
}

func logInOurOutQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.LogInOurOutButton, constants.LogInOurOutData)
}

func startQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.StartQueueButton, constants.StartQueueData)
}

func nextButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.NextButton, constants.NextData)
}

func goToMenuButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.GoToMenuButton, constants.GoToMenuData)
}

func endQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(constants.FinishQueueButton, constants.FinishQueueData)
}
