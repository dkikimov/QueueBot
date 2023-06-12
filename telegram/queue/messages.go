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

func GetQueueMessage(description string) tgbotapi.InputTextMessageContent {
	answer := tgbotapi.InputTextMessageContent{
		Text:      getMessageContent(description, nil),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetUpdatedQueueMessage(messageID string, description string, users []user.User) tgbotapi.EditMessageTextConfig {
	keyboard := GetKeyboard()
	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
			ReplyMarkup:     &keyboard,
		},
		Text:      getMessageContent(description, users),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetForwardMessage(chatId int64, description string) tgbotapi.MessageConfig {
	answer := tgbotapi.NewMessage(chatId, constants.ForwardQueueToMessage)
	answer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonSwitch(constants.ForwardQueueButton, description),
	))
	answer.ParseMode = tgbotapi.ModeMarkdown
	return answer
}
