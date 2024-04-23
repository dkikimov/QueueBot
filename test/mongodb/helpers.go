package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"QueueBot/internal/entity"
	"QueueBot/internal/usecase/storage/mongodb"
)

func populateDB(ctx context.Context, client *mongo.Client) error {
	collection := client.Database(mongodb.DatabaseName).Collection(mongodb.CollectionName)
	if _, err := collection.InsertOne(ctx, bson.M{
		"messageId":        "123",
		"description":      "123",
		"currentPersonIdx": 0,
		"users": []entity.User{
			{
				ID:   1,
				Name: "Username",
			},
		},
	}, nil); err != nil {
		return fmt.Errorf("couldn't insert document: %w", err)
	}

	if _, err := collection.InsertOne(ctx, bson.M{
		"messageId":        "456",
		"description":      "456",
		"currentPersonIdx": 0,
		"users":            []entity.User{},
	}, nil); err != nil {
		return fmt.Errorf("couldn't insert document: %w", err)
	}

	if _, err := collection.InsertOne(ctx, bson.M{
		"messageId":        "789",
		"description":      "789",
		"currentPersonIdx": 0,
		"users": []entity.User{
			{
				ID:   1,
				Name: "Username",
			},
			{
				ID:   2,
				Name: "Username2",
			},
			{
				ID:   3,
				Name: "Username3",
			},
		},
	}, nil); err != nil {
		return fmt.Errorf("couldn't insert document: %w", err)
	}

	return nil
}
