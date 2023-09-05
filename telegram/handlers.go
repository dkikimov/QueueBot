package telegram

import (
	"QueueBot/constants"
	"QueueBot/storage"
	"QueueBot/telegram/queue"
	"QueueBot/telegram/steps"
	"QueueBot/ui"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI, storage storage.Storage, errChan chan<- error) {
	// Проверяем, если сообщение - команда.
	// Если да, отправляем соотвутствующее сообщение
	switch message.Command() {
	case constants.StartCommand:
		if err := SendHelloMessage(message, bot, storage); err != nil {
			errChan <- fmt.Errorf("sendHelloMessage error occured: %s", err)
			return
		}
		return
	case constants.CreateQueueCommand:
		if err := SendMessageToCreateQueue(message, bot, storage); err != nil {
			errChan <- fmt.Errorf("sendMessageToCreateMessage error occured: %s", err)
			return
		}
		return
	}

	// В случае, если сообщение не команда
	// Получаем текущее состояние пользователя для понимания какое действие ожидается быть следующим
	currentStep, err := storage.GetUserCurrentStep(message.From.ID)
	if err != nil {
		errChan <- fmt.Errorf("couldn't get current user step with error: %s", err)
	}

	switch currentStep {
	case steps.Menu:
		if err := SendHelloMessage(message, bot, storage); err != nil {
			errChan <- fmt.Errorf("sendHelloMessage error occured: %s", err)
			return
		}
	case steps.EnteringDescription:
		if err := SendForwardToMessage(message, bot); err != nil {
			errChan <- fmt.Errorf("sendForwardMessage error occured: %s", err)
			return
		}
	default:
		// TODO: panic?
		errChan <- fmt.Errorf("got current step (%v) that is not implemented", currentStep)
	}
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage, errChan chan<- error) {
	// Сверяемся со скрытыми данными, заложенными в сообщении для определения команды
	wasError := false
	switch callbackQuery.Data {
	case constants.LogInOurOutData:
		if err := queue.LogInOurOut(callbackQuery, bot, storage); err != nil {
			errChan <- fmt.Errorf("couldn't login or logout with error: %s", err)
			wasError = true
		}
	case constants.StartQueueData:
		if err := queue.Start(callbackQuery, bot, storage, false); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with error: %s", err)
			wasError = true
		}
	case constants.StartQueueShuffleData:
		if err := queue.Start(callbackQuery, bot, storage, true); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with shuffle with error: %s", err)
			wasError = true
		}
	case constants.NextData:
		if err := queue.Next(callbackQuery, bot, storage); err != nil {
			errChan <- fmt.Errorf("couldn't go to next person with error: %s", err)
			wasError = true
		}
	case constants.GoToMenuData:
		if err := queue.GoToMenu(callbackQuery, bot, storage); err != nil {
			errChan <- fmt.Errorf("couldn't go to menu with error: %s", err)
			wasError = true
		}
	case constants.FinishQueueData:
		if err := queue.FinishQueue(callbackQuery, bot, storage); err != nil {
			errChan <- fmt.Errorf("couldn't finish queue with error: %s", err)
			wasError = true
		}
	}

	var callback tgbotapi.CallbackConfig
	if wasError {
		callback = tgbotapi.NewCallback(callbackQuery.ID, constants.ActionError)
	} else {
		callback = tgbotapi.NewCallback(callbackQuery.ID, constants.ActionCompleted)
	}
	if _, err := bot.Request(callback); err != nil {
		errChan <- fmt.Errorf("couldn't process next_data callback with error: %s", err)
		return
	}
}

func HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult, storage storage.Storage, errChan chan<- error) {
	// Обрубаем слишком длинные описания
	if len(chosenInlineResult.Query) > 100 {
		chosenInlineResult.Query = chosenInlineResult.Query[:100]
	}

	if err := queue.Create(chosenInlineResult.InlineMessageID, chosenInlineResult.Query, storage); err != nil {
		errChan <- err
	}
}

func HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI, errChan chan<- error) {
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
		errChan <- fmt.Errorf("couldn't handle inline query with error: %s", err)
	}
}
