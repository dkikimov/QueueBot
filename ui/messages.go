package ui

import (
	"QueueBot/constants"
	"QueueBot/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getMessageContentBeforeStart(title string, users []user.User) string {
	return fmt.Sprintf("*%s*\n%s\n%s", title, constants.QueueDescription, user.ListToString(users))
}

func getMessageContentAfterStart(title string, users []user.User, currentPersonIndex int) string {
	return fmt.Sprintf("*%s*\n%s\n%s", title, constants.QueueDescription, user.ListToStringWithCurrent(users, currentPersonIndex))
}

func GetQueueMessageContent(description string) tgbotapi.InputTextMessageContent {
	answer := tgbotapi.InputTextMessageContent{
		Text:      getMessageContentBeforeStart(description, nil),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetQueueMessage(messageID string, users []user.User, description string) tgbotapi.EditMessageTextConfig {
	keyboard := GetBeforeStartKeyboard()
	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
			ReplyMarkup:     &keyboard,
		},
		Text:      getMessageContentBeforeStart(description, users),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetUpdatedQueueMessage(messageID string, description string, users []user.User) tgbotapi.EditMessageTextConfig {
	keyboard := GetBeforeStartKeyboard()
	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
			ReplyMarkup:     &keyboard,
		},
		Text:      getMessageContentBeforeStart(description, users),
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

func GetQueueAfterStartMessage(messageID string, description string, users []user.User, currentPersonIndex int) tgbotapi.EditMessageTextConfig {
	keyboard := GetAfterStartKeyboard()

	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
			ReplyMarkup:     &keyboard,
		},
		Text:      getMessageContentAfterStart(description, users, currentPersonIndex),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetEndQueueMessage(messageID string) tgbotapi.EditMessageTextConfig {
	keyboard := GetEndedQueueKeyboard()

	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
			ReplyMarkup:     &keyboard,
		},
		Text: constants.EndedQueue,
	}
	return answer
}
