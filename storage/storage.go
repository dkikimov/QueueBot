package storage

import (
	"QueueBot/telegram/steps"
	"QueueBot/user"
)

type Storage interface {
	GetUsersInQueue(messageId string) ([]user.User, error)
	CreateQueue(messageId string, description string) error
	SetUserCurrentStep(userId int64, currentStep steps.Step) error
	CreateUser(userId int64) error
	GetUserCurrentStep(userId int64) (steps.Step, error)
	GetDescriptionOfQueue(messageId string) (string, error)
	LogInOurOutQueue(messageId string, user user.User) error

	StartQueue(messageId string, isShuffle bool) (err error, wasUpdated bool)
	IncrementCurrentPerson(messageId string) (err error, currentPerson int)
	GoToMenu(messageId string) error
}
