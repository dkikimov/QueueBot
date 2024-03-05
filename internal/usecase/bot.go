package usecase

import (
	"context"
	"fmt"

	"QueueBot/internal/entity"
	"QueueBot/internal/usecase/storage"
)

type Bot interface {
	CreateQueue(ctx context.Context, messageId string, description string) error
	LogInOutToQueue(ctx context.Context, messageId string, user entity.User) error
	StartQueue(ctx context.Context, messageId string, shuffle bool) error
	FinishQueue(ctx context.Context, messageId string) error
	SetNextPersonToQueue(ctx context.Context, messageId string) error
	GetQueue(ctx context.Context, messageId string) (entity.Queue, error)
}

type BotUseCase struct {
	Storage storage.Storage
}

func NewBotUseCase(storage storage.Storage) *BotUseCase {
	return &BotUseCase{Storage: storage}
}

func (b BotUseCase) CreateQueue(ctx context.Context, messageId string, description string) error {
	err := b.Storage.CreateQueue(ctx, messageId, description)
	if err != nil {
		return fmt.Errorf("couldn't create queue in storage with error: %s", err)
	}

	return nil
}

func (b BotUseCase) LogInOutToQueue(ctx context.Context, messageId string, user entity.User) error {
	err := b.Storage.LogInOutToQueue(ctx, messageId, user)
	if err != nil {
		return fmt.Errorf("couldn't add user to queue in storage with error: %s", err)
	}

	return nil
}

func (b BotUseCase) StartQueue(ctx context.Context, messageId string, shuffle bool) error {
	err := b.Storage.StartQueue(ctx, messageId, shuffle)
	if err != nil {
		return fmt.Errorf("couldn't start queue in storage with error: %s", err)
	}

	return nil
}

func (b BotUseCase) FinishQueue(ctx context.Context, messageId string) error {
	err := b.Storage.DeleteQueue(ctx, messageId)
	if err != nil {
		return fmt.Errorf("couldn't finish queue in storage with error: %s", err)
	}

	return nil
}

func (b BotUseCase) SetNextPersonToQueue(ctx context.Context, messageId string) (err error) {
	err = b.Storage.IncrementCurrentPerson(ctx, messageId)
	if err != nil {
		return fmt.Errorf("couldn't set next person to queue in storage with error: %s", err)
	}

	return nil
}

func (b BotUseCase) GetQueue(ctx context.Context, messageId string) (entity.Queue, error) {
	queue, err := b.Storage.GetQueue(ctx, messageId)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't get queue from storage with error: %s", err)
	}

	return queue, nil
}
