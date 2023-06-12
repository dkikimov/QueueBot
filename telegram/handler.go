package telegram

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/telegram/queue"
	"QueueBot/telegram/steps"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, storage storage.Storage) {
	switch update.Message.Command() {
	case constants.StartCommand:
		SendHelloMessage(update.Message, bot, storage)
		return
	case constants.CreateQueueCommand:
		SendMessageToCreateQueue(update.Message, bot, storage)
		return
	}

	currentStep, err := storage.GetUserCurrentStep(update.Message.From.ID)
	if err != nil {
		logger.Fatalf("Couldn't get current user step with error: %s", err.Error())
	}

	switch currentStep {
	case steps.Menu:
		SendHelloMessage(update.Message, bot, storage)
		break
	case steps.EnteringDescription:
		SendForwardToMessage(update.Message, bot)
		break
	}
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	switch callbackQuery.Data {
	case constants.AddToQueueData:
		queue.AddTo(callbackQuery, bot, storage)
	}
}

func HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult, bot *tgbotapi.BotAPI, storage storage.Storage) {
	queue.Create(chosenInlineResult.InlineMessageID, chosenInlineResult.Query, bot, storage)
}

func HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	article := tgbotapi.NewInlineQueryResultArticle(inlineQuery.ID, constants.CreateQueue, fmt.Sprintf("С описанием: %s", inlineQuery.Query))
	article.InputMessageContent = queue.GetQueueMessage(inlineQuery.Query)

	keyboard := queue.GetKeyboardButtons(false)
	article.ReplyMarkup = &keyboard

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		CacheTime:     9999999,
		Results:       []interface{}{article},
	}

	_, err := bot.Request(inlineConf)
	if err != nil {
		logger.Fatalf("Couldn't handle inline query with error: %s", err.Error())
	}

}
