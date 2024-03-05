package messages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const LogInOurOutButton = "Добавиться/выйти из очереди"

const StartQueueButton = "Старт в порядке очереди"
const StartQueueShuffleButton = "Старт в случайном порядке"

const NextButton = "Следующий"
const GoToMenuButton = "Перейти в меню"
const FinishQueueButton = "Закончить"

const LogInOurOutData = "log_in_our_out"
const StartQueueData = "start_queue"
const StartQueueShuffleData = "start_queue_shuffle"
const NextData = "next_user"
const GoToMenuData = "go_to_menu"
const FinishQueueData = "finish_queue"

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
