package telegram

import (
	"QueueBot/constants"
	"QueueBot/storage"
	"QueueBot/telegram/queue"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.Message.Command() {
	//case CreateQueue:
	//	queue.Create(update.InlineQuery, bot)
	}
}

func HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	queue.Create(inlineQuery, bot, storage)
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	switch callbackQuery.Data {
	case constants.AddToQueueData:

	}
}
