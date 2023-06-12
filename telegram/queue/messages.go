package queue

import (
	"QueueBot/constants"
	"QueueBot/storage/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getMessageContent(title string, users []user.User) string {
	return fmt.Sprintf("*%s*\n%s\n%s", title, constants.QueueDescription, user.UsersToString(users))
}

func GetQueueMessage(chatId int64, description string, users []user.User) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(chatId, getMessageContent(description, users))
	answer.ParseMode = tgbotapi.ModeMarkdown
	answer.ReplyMarkup = GetKeyboardButtons(false)

	return answer
}

func GetUpdatedQueueMessage(chatID int64, messageID int, description string, users []user.User, isUserInQueue bool) tgbotapi.EditMessageTextConfig {
	keyboard := GetKeyboardButtons(isUserInQueue)
	answer := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, getMessageContent(description, users), keyboard)
	answer.ParseMode = tgbotapi.ModeMarkdown
	return answer
}

func GetForwardMessage(chatId int64, description string) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(chatId, constants.ForwardQueueToMessage)
	answer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonSwitch(constants.ForwardQueueTo, description),
	))
	answer.ParseMode = tgbotapi.ModeMarkdown
	return answer
}
