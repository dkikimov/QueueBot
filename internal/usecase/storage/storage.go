package storage

import (
	"context"

	"QueueBot/internal/entity"
)

type Storage interface {
	CreateQueue(ctx context.Context, messageId string, description string) error
	LogInOutToQueue(ctx context.Context, messageId string, user entity.User) error
	GetQueue(ctx context.Context, messageId string) (entity.Queue, error)

	StartQueue(ctx context.Context, messageId string, isShuffle bool) error
	IncrementCurrentPerson(ctx context.Context, messageId string) error
	DeleteQueue(ctx context.Context, messageId string) error

	Close() error
}
