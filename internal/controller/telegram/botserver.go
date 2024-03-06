package telegram

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/controller/telegram/messages"
)

const StartCommand = "start"

const HelloMessage = `Привет! Я бот, предназначенный для создания очередей. 
Введи описание своей очереди, а я тебе ее создам`

const (
	ActionCompleted = "Действие выполнено!"
	ActionError     = "Произошла ошибка"
)

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
		go func(update tgbotapi.Update) {
			switch {
			case update.Message != nil:
				if err := s.HandleMessage(update.Message); err != nil {
					errChan <- fmt.Errorf("couldn't handle message: %w", err)
				}
			case update.CallbackQuery != nil:

				if err := s.HandleCallbackQuery(update.CallbackQuery); err != nil {
					errChan <- fmt.Errorf("couldn't handle callback query: %w", err)
				}
			case update.InlineQuery != nil:
				if err := s.HandleInlineQuery(update.InlineQuery); err != nil {
					errChan <- fmt.Errorf("couldn't handle inline query: %w", err)
				}
			case update.ChosenInlineResult != nil:
				if err := s.HandleChosenInlineResult(update.ChosenInlineResult); err != nil {
					errChan <- fmt.Errorf("couldn't handle chosen inline result: %w", err)
				}
			}
		}(update)
	}
}

func (s BotServer) HandleMessage(message *tgbotapi.Message) error {
	// Проверяем, если сообщение - команда.
	// Если да, отправляем соотвутствующее сообщение
	if message.Command() == StartCommand {
		if err := s.bot.SendHelloMessage(message); err != nil {
			return fmt.Errorf("sendHelloMessage error occurred: %w", err)
		}
	}

	if err := s.bot.SendForwardMessageButton(message); err != nil {
		return fmt.Errorf("sendMessageToCreateMessage error occurred: %w", err)
	}

	return nil
}

func (s BotServer) HandleChosenInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	// Обрубаем слишком длинные описания
	if len(chosenInlineResult.Query) > 100 {
		chosenInlineResult.Query = chosenInlineResult.Query[:100]
	}

	if err := s.bot.CreateQueue(context.Background(), chosenInlineResult.InlineMessageID, chosenInlineResult.Query); err != nil {
		return fmt.Errorf("couldn't create queue: %w", err)
	}

	return nil
}

func (s BotServer) HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery) error {
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
		return fmt.Errorf("couldn't handle inline query with error: %w", err)
	}

	return nil
}
