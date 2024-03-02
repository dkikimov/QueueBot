package telegram

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/steps"
	"QueueBot/internal/telegram/messages"
)

const CreateQueueCommand = "create"
const StartCommand = "start"

const CreateQueueMessage = "Окей. Теперь введи для чего предназначена эта очередь"
const HelloMessage = "Привет! Я бот, предназначенный для создания очередей. \nДля этого введи команду /create"

const ActionCompleted = "Действие выполнено!"
const ActionError = "Произошла ошибка"

const CreateQueue = "Создать очередь"

type BotServer struct {
	bot *Bot
}

func NewBotServer(bot *Bot) *BotServer {
	return &BotServer{bot: bot}
}

func (s BotServer) Listen(config tgbotapi.UpdateConfig, errChan chan<- error) {
	updates := s.bot.TgBot.GetUpdatesChan(config)
	slog.Info("Started listening update channel")

	for update := range updates {
		switch {
		case update.Message != nil:
			go s.HandleMessage(update.Message, errChan)
		case update.CallbackQuery != nil:
			go s.HandleCallbackQuery(update.CallbackQuery, errChan)
		case update.InlineQuery != nil:
			go s.HandleInlineQuery(update.InlineQuery, errChan)
		case update.ChosenInlineResult != nil:
			go s.HandleChosenInlineResult(update.ChosenInlineResult, errChan)
		}
	}
}

func (s BotServer) HandleMessage(message *tgbotapi.Message, errChan chan<- error) {
	// Проверяем, если сообщение - команда.
	// Если да, отправляем соотвутствующее сообщение
	switch message.Command() {
	case StartCommand:
		if err := s.bot.SendHelloMessage(message); err != nil {
			errChan <- fmt.Errorf("sendHelloMessage error occured: %s", err)
			return
		}
		return
	case CreateQueueCommand:
		if err := s.bot.SendMessageToCreateQueue(message); err != nil {
			errChan <- fmt.Errorf("sendMessageToCreateMessage error occured: %s", err)
			return
		}
		return
	}

	// В случае, если сообщение не команда
	// Получаем текущее состояние пользователя для понимания какое действие ожидается быть следующим
	currentStep, err := s.bot.Storage.GetUserCurrentStep(message.From.ID)
	if err != nil {
		errChan <- fmt.Errorf("couldn't get current user step with error: %s", err)
	}

	switch currentStep {
	case steps.Menu:
		if err := s.bot.SendHelloMessage(message); err != nil {
			errChan <- fmt.Errorf("sendHelloMessage error occured: %s", err)
			return
		}
	case steps.EnteringDescription:
		if err := s.bot.SendForwardToMessage(message); err != nil {
			errChan <- fmt.Errorf("sendForwardMessage error occured: %s", err)
			return
		}
	default:
		errChan <- fmt.Errorf("got current step (%v) that is not implemented", currentStep)
	}
}

func (s BotServer) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, errChan chan<- error) {
	// Сверяемся со скрытыми данными, заложенными в сообщении для определения команды
	wasError := false
	switch callbackQuery.Data {
	case messages.LogInOurOutData:
		if err := s.bot.LogInOurOut(callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't login or logout with error: %s", err)
			wasError = true
		}
	case messages.StartQueueData:
		if err := s.bot.Start(callbackQuery, false); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with error: %s", err)
			wasError = true
		}
	case messages.StartQueueShuffleData:
		if err := s.bot.Start(callbackQuery, true); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with shuffle with error: %s", err)
			wasError = true
		}
	case messages.NextData:
		if err := s.bot.Next(callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't go to next person with error: %s", err)
			wasError = true
		}
	case messages.GoToMenuData:
		if err := s.bot.GoToMenu(callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't go to menu with error: %s", err)
			wasError = true
		}
	case messages.FinishQueueData:
		if err := s.bot.FinishQueue(callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't finish queue with error: %s", err)
			wasError = true
		}
	}

	var callback tgbotapi.CallbackConfig
	if wasError {
		callback = tgbotapi.NewCallback(callbackQuery.ID, ActionError)
	} else {
		callback = tgbotapi.NewCallback(callbackQuery.ID, ActionCompleted)
	}
	if _, err := s.bot.TgBot.Request(callback); err != nil {
		errChan <- fmt.Errorf("couldn't process next_data callback with error: %s", err)
		return
	}
}

func (s BotServer) HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult, errChan chan<- error) {
	// Обрубаем слишком длинные описания
	if len(chosenInlineResult.Query) > 100 {
		chosenInlineResult.Query = chosenInlineResult.Query[:100]
	}

	if err := s.bot.CreateQueue(chosenInlineResult.InlineMessageID, chosenInlineResult.Query); err != nil {
		errChan <- err
	}
}

func (s BotServer) HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, errChan chan<- error) {
	article := tgbotapi.NewInlineQueryResultArticle(inlineQuery.ID, CreateQueue, fmt.Sprintf("С описанием: %s", inlineQuery.Query))
	article.InputMessageContent = messages.GetQueueMessageContent(inlineQuery.Query)

	keyboard := messages.GetBeforeStartKeyboard()
	article.ReplyMarkup = &keyboard

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		CacheTime:     9999,
		Results:       []interface{}{article},
	}

	_, err := s.bot.TgBot.Request(inlineConf)
	if err != nil {
		errChan <- fmt.Errorf("couldn't handle inline query with error: %s", err)
	}
}
