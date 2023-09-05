package telegram

import (
	"QueueBot/constants"
	"QueueBot/logger"
	"QueueBot/storage"
	"QueueBot/telegram/queue"
	"QueueBot/telegram/steps"
	"QueueBot/ui"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage) {
	// Проверяем, если сообщение - команда.
	// Если да, отправляем соотвутствующее сообщение
	switch message.Command() {
	case constants.StartCommand:
		SendHelloMessage(message, bot, storage)
		return
	case constants.CreateQueueCommand:
		SendMessageToCreateQueue(message, bot, storage)
		return
	}

	// Получаем текущее состояние пользователя для понимания какое действие ожидается быть следующим
	currentStep, err := storage.GetUserCurrentStep(message.From.ID)
	if err != nil {
		logger.Fatalf("Couldn't get current user step with error: %s", err.Error())
	}

	switch currentStep {
	case steps.Menu:
		SendHelloMessage(message, bot, storage)
		break
	case steps.EnteringDescription:
		SendForwardToMessage(message, bot)
		break
	}
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	// Сверяемся со скрытыми данными, заложенными в сообщении для определения команды
	switch callbackQuery.Data {
	case constants.LogInOurOutData:
		queue.LogInOurOut(callbackQuery, bot, storage)
	case constants.StartQueueData:
		queue.Start(callbackQuery, bot, storage, false)
	case constants.StartQueueShuffleData:
		queue.Start(callbackQuery, bot, storage, true)
	case constants.NextData:
		queue.Next(callbackQuery, bot, storage)
	case constants.GoToMenuData:
		queue.GoToMenu(callbackQuery, bot, storage)
	case constants.FinishQueueData:
		queue.FinishQueue(callbackQuery, bot, storage)
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, constants.ActionCompleted)
	if _, err := bot.Request(callback); err != nil {
		logger.Panicf("Couldn't process next_data callback with error: %s", err.Error())
	}
}

func HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult, storage storage.Storage) {
	//TODO: Middleware for query
	queue.Create(chosenInlineResult.InlineMessageID, chosenInlineResult.Query, storage)
}

func HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI) {
	article := tgbotapi.NewInlineQueryResultArticle(inlineQuery.ID, constants.CreateQueue, fmt.Sprintf("С описанием: %s", inlineQuery.Query))
	article.InputMessageContent = ui.GetQueueMessageContent(inlineQuery.Query)

	keyboard := ui.GetBeforeStartKeyboard()
	article.ReplyMarkup = &keyboard

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	_, err := bot.Request(inlineConf)
	if err != nil {
		logger.Fatalf("Couldn't handle inline query with error: %s", err.Error())
	}

}
