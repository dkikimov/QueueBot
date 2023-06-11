package telegram_errors

import "QueueBot/logger"

func HandleSendMessage(err error) {
	logger.PrintfError("Couldn't send message with error %s", err.Error())
}
