package storage

import (
	"QueueBot/storage/user"
	"QueueBot/telegram/steps"
)

type Storage interface {
	GetUsersInQueue(messageId string) ([]user.User, error)
	AddUserToQueue(messageId string, user user.User) error
	DeleteUserFromQueueById(messageId string, userId int64) error
	CreateQueue(messageId string, description string) error
	SetUserCurrentStep(userId int64, currentStep steps.Step) error
	CreateUser(userId int64) error
	GetUserCurrentStep(userId int64) (steps.Step, error)
	GetDescriptionOfQueue(messageId string) (string, error)
}
