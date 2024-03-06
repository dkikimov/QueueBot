package usecase

import (
	"context"
	"fmt"

	"QueueBot/internal/entity"
	"QueueBot/internal/usecase/storage"
)

type Bot interface {
	CreateQueue(ctx context.Context, messageID string, description string) error
	LogInOutToQueue(ctx context.Context, messageID string, user entity.User) error
	StartQueue(ctx context.Context, messageID string, shuffle bool) error
	FinishQueue(ctx context.Context, messageID string) error
	SetNextPersonToQueue(ctx context.Context, messageID string) error
	GetQueue(ctx context.Context, messageID string) (entity.Queue, error)
}

type BotUseCase struct {
	Storage storage.Storage
}

func NewBotUseCase(storage storage.Storage) *BotUseCase {
	return &BotUseCase{Storage: storage}
}

func (b BotUseCase) CreateQueue(ctx context.Context, messageID string, description string) error {
	err := b.Storage.CreateQueue(ctx, messageID, description)
	if err != nil {
		return fmt.Errorf("couldn't create queue in storage with error: %w", err)
	}

	return nil
}

func (b BotUseCase) LogInOutToQueue(ctx context.Context, messageID string, user entity.User) error {
	err := b.Storage.LogInOutToQueue(ctx, messageID, user)
	if err != nil {
		return fmt.Errorf("couldn't add user to queue in storage with error: %w", err)
	}

	return nil
}

func (b BotUseCase) StartQueue(ctx context.Context, messageID string, shuffle bool) error {
	err := b.Storage.StartQueue(ctx, messageID, shuffle)
	if err != nil {
		return fmt.Errorf("couldn't start queue in storage with error: %w", err)
	}

	return nil
}

func (b BotUseCase) FinishQueue(ctx context.Context, messageID string) error {
	err := b.Storage.DeleteQueue(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't finish queue in storage with error: %w", err)
	}

	return nil
}

func (b BotUseCase) SetNextPersonToQueue(ctx context.Context, messageID string) (err error) {
	err = b.Storage.IncrementCurrentPerson(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't set next person to queue in storage with error: %w", err)
	}

	return nil
}

func (b BotUseCase) GetQueue(ctx context.Context, messageID string) (entity.Queue, error) {
	queue, err := b.Storage.GetQueue(ctx, messageID)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't get queue from storage with error: %w", err)
	}

	return queue, nil
}
