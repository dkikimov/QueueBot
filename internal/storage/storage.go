package storage

import (
	"QueueBot/internal/models"
	"QueueBot/internal/steps"
)

type Storage interface {
	GetUsersInQueue(messageId string) ([]models.User, error)
	GetUsersInQueueCheckShuffle(messageId string) ([]models.User, error)

	CreateQueue(messageId string, description string) error
	SetUserCurrentStep(userId int64, currentStep steps.ChatStep) error
	CreateUser(userId int64) error
	GetUserCurrentStep(userId int64) (steps.ChatStep, error)
	GetDescriptionOfQueue(messageId string) (string, error)
	LogInOurOutQueue(messageId string, user models.User) error

	StartQueue(messageId string, isShuffle bool) (err error, wasUpdated bool)
	IncrementCurrentPerson(messageId string) (err error, currentPerson int)
	GoToMenu(messageId string) error
	ShuffleUsers(messageId string) error
	FinishQueueDeleteParticipants(messageId string) error

	Close() error
}
