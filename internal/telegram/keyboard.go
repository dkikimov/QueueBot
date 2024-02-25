package telegram

import (
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
		tgbotapi.NewInlineKeyboardRow(
			startQueueShuffleButton(),
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
	return tgbotapi.NewInlineKeyboardButtonData(LogInOurOutButton, LogInOurOutData)
}

func startQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(StartQueueButton, StartQueueData)
}

func startQueueShuffleButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(StartQueueShuffleButton, StartQueueShuffleData)
}

func nextButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(NextButton, NextData)
}

func goToMenuButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(GoToMenuButton, GoToMenuData)
}

func endQueueButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(FinishQueueButton, FinishQueueData)
}
