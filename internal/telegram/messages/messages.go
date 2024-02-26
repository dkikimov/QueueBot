package messages

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/models"
)

const ForwardQueueToMessage = "Отлично! Теперь с помощью кнопки ниже вы можете переслать свою 'очередь'"
const QueueDescription = "В очереди состоят:"
const EndedQueue = "Участники закончились, значит и очередь тоже. Что делаем дальше?"
const FinishedQueue = "'Очередь' окончена 🎉"
const ForwardQueueButton = "Переслать 'очередь'"

func getMessageContentBeforeStart(title string, users []models.User) string {
	return fmt.Sprintf("*%s*\n%s\n%s", title, QueueDescription, models.ListToString(users))
}

func getMessageContentAfterStart(title string, users []models.User, currentPersonIndex int) string {
	return fmt.Sprintf("*%s*\n%s\n%s", title, QueueDescription, models.ListToStringWithCurrent(users, currentPersonIndex))
}

func GetQueueMessageContent(description string) tgbotapi.InputTextMessageContent {
	answer := tgbotapi.InputTextMessageContent{
		Text:      getMessageContentBeforeStart(description, nil),
		ParseMode: tgbotapi.ModeMarkdown,
	}
	return answer
}

func GetQueueMessage(messageID string, users []models.User, description string) tgbotapi.EditMessageTextConfig {
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

func GetUpdatedQueueMessage(messageID string, description string, users []models.User) tgbotapi.EditMessageTextConfig {
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
	answer := tgbotapi.NewMessage(chatId, ForwardQueueToMessage)
	answer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonSwitch(ForwardQueueButton, description),
	))
	answer.ParseMode = tgbotapi.ModeMarkdown
	return answer
}

func GetQueueAfterStartMessage(messageID string, description string, users []models.User, currentPersonIndex int) tgbotapi.EditMessageTextConfig {
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
		Text: EndedQueue,
	}
	return answer
}

func GetFinishedMessage(messageID string) tgbotapi.EditMessageTextConfig {
	answer := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			InlineMessageID: messageID,
		},
		Text: FinishedQueue,
	}
	return answer
}
