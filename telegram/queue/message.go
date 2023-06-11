package queue

import (
	"QueueBot/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMessageContent(title string) tgbotapi.InputTextMessageContent {
	return tgbotapi.InputTextMessageContent{
		Text:      fmt.Sprintf("*%s*\n%s", title, constants.QueueDescription),
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
