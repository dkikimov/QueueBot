package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"QueueBot/internal/controller/telegram/client"
)

const CreateQueue = "Создать очередь"

func (s BotServer) HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery) error {
	article := tgbotapi.NewInlineQueryResultArticle(inlineQuery.ID, CreateQueue, fmt.Sprintf("С описанием: %s", inlineQuery.Query))
	article.InputMessageContent = client.GetQueueMessageContent(inlineQuery.Query)

	keyboard := client.GetBeforeStartKeyboard()
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
