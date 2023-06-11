package queue

import (
	"QueueBot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
)

func Create(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	if err := storage.CreateQueue(inlineQuery.ID, inlineQuery.Query); err != nil {
		// TODO:
		return
	}

	message := tgbotapi.InlineQueryResultArticle{
		Type:                "article",
		ID:                  strconv.FormatInt(rand.Int63(), 10),
		Title:               "Создать новую очередь",
		InputMessageContent: GetMessageContent(inlineQuery.Query),
		ReplyMarkup:         GetKeyboardButtons(),
	}

	answer := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       []interface{}{message},
		CacheTime:     1000000000000,
	}

	_, err := bot.Send(answer)
	if err != nil {
		return
	}
}

func AddTo(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, storage storage.Storage) {
	//if err := storage.CreateQueue(inlineQuery.ID, inlineQuery.Query); err != nil {
	//	// TODO:
	//	return
	//}
	//
	//message := tgbotapi.InlineQueryResultArticle{
	//	Type:                "article",
	//	ID:                  strconv.FormatInt(rand.Int63(), 10),
	//	Title:               "Создать новую очередь",
	//	InputMessageContent: GetMessageContent(inlineQuery.Query),
	//	ReplyMarkup:         GetKeyboardButtons(),
	//}
	//
	//answer := tgbotapi.InlineConfig{
	//	InlineQueryID: inlineQuery.ID,
	//	Results:       []interface{}{message},
	//	CacheTime:     1000000000000,
	//}
	//
	//_, err := bot.Send(answer)
	//if err != nil {
	//	return
	//}
}
