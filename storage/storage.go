package storage

import (
	"QueueBot/storage/user"
	"QueueBot/telegram/steps"
)

type Storage interface {
	GetUsersInQueue(messageId int) ([]user.User, error)
	AddUserToQueue(messageId int, user user.User) error
	DeleteUserFromQueueById(messageId int, userId int64) error
	CreateQueue(messageId int, description string) error
	SetUserCurrentStep(userId int64, currentStep steps.Step) error
	CreateUser(userId int64) error
	GetUserCurrentStep(userId int64) (steps.Step, error)
	GetDescriptionOfQueue(messageId int) (string, error)
}
