package client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const LogInOurOutButton = "Добавиться/выйти из очереди"

const (
	StartQueueButton        = "Старт в порядке очереди"
	StartQueueShuffleButton = "Старт в случайном порядке"
)

const (
	NextButton        = "Следующий"
	GoToMenuButton    = "Перейти в меню"
	FinishQueueButton = "Закончить"
)

const (
	LogInOurOutData       = "log_in_our_out"
	StartQueueData        = "start_queue"
	StartQueueShuffleData = "start_queue_shuffle"
	NextData              = "next_user"
	GoToMenuData          = "go_to_menu"
	FinishQueueData       = "finish_queue"
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
