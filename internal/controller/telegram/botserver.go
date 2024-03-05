package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/controller/telegram/messages"
)

const StartCommand = "start"

const CreateQueueMessage = "Окей. Теперь введи для чего предназначена эта очередь"
const HelloMessage = `Привет! Я бот, предназначенный для создания очередей. 
Введи описание своей очереди, а я тебе ее создам`

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
		if err := s.bot.SendHelloMessage(context.Background(), message); err != nil {
			errChan <- fmt.Errorf("sendHelloMessage error occured: %s", err)
			return
		}
		return
	}

	if err := s.bot.SendMessageToCreateQueue(context.Background(), message); err != nil {
		errChan <- fmt.Errorf("sendMessageToCreateMessage error occured: %s", err)
	}
}

func (s BotServer) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, errChan chan<- error) {
	// Сверяемся со скрытыми данными, заложенными в сообщении для определения команды
	startTime := time.Now()
	slog.Debug("Got callback query with data: ", "data", callbackQuery.Data)

	wasError := false
	switch callbackQuery.Data {
	case messages.LogInOurOutData:
		if err := s.bot.LogInOurOut(context.Background(), callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't login or logout with error: %s", err)
			wasError = true
		}
	case messages.StartQueueData:
		if err := s.bot.Start(context.Background(), callbackQuery, false); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with error: %s", err)
			wasError = true
		}
	case messages.StartQueueShuffleData:
		if err := s.bot.Start(context.Background(), callbackQuery, true); err != nil {
			errChan <- fmt.Errorf("couldn't start queue with shuffle with error: %s", err)
			wasError = true
		}
	case messages.NextData:
		if err := s.bot.Next(context.Background(), callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't go to next person with error: %s", err)
			wasError = true
		}
	case messages.GoToMenuData:
		if err := s.bot.GoToMenu(context.Background(), callbackQuery); err != nil {
			errChan <- fmt.Errorf("couldn't go to menu with error: %s", err)
			wasError = true
		}
	case messages.FinishQueueData:
		if err := s.bot.FinishQueue(context.Background(), callbackQuery); err != nil {
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

	slog.Debug("Processed callback query with data: ", "data", callbackQuery.Data, "elapsed", time.Now().Sub(startTime).String())
}

func (s BotServer) HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult, errChan chan<- error) {
	// Обрубаем слишком длинные описания
	if len(chosenInlineResult.Query) > 100 {
		chosenInlineResult.Query = chosenInlineResult.Query[:100]
	}

	if err := s.bot.CreateQueue(context.Background(), chosenInlineResult.InlineMessageID, chosenInlineResult.Query); err != nil {
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
