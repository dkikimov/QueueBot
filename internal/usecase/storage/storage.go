package storage

import (
	"context"

	"QueueBot/internal/entity"
)

type Storage interface {
	CreateQueue(ctx context.Context, messageID string, description string) error
	LogInOutToQueue(ctx context.Context, messageID string, user entity.User) error
	GetQueue(ctx context.Context, messageID string) (entity.Queue, error)

	StartQueue(ctx context.Context, messageID string, isShuffle bool) error
	IncrementCurrentPerson(ctx context.Context, messageID string) error
	DeleteQueue(ctx context.Context, messageID string) error

	Close() error
}
