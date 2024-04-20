package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"QueueBot/internal/apperrors"
	"QueueBot/internal/entity"
)

type Database struct {
	mongoClient *mongo.Client
	queueMutex  sync.Map
}

func (db *Database) IncrementCurrentPerson(ctx context.Context, messageID string) error {
	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	filter := bson.D{{Key: "messageId", Value: messageID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "currentUserIndex", Value: 1}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("couldn't increment current person: %w", err)
	}

	return nil
}

func (db *Database) DeleteQueue(ctx context.Context, messageID string) error {
	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	filter := bson.D{{Key: "messageId", Value: messageID}}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("couldn't delete queue: %w", err)
	}

	return nil
}

func (db *Database) CreateQueue(ctx context.Context, messageID string, description string) error {
	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	queue := bson.D{
		{Key: "messageId", Value: messageID},
		{Key: "description", Value: description},
		{Key: "currentUserIndex", Value: 0},
		{Key: "users", Value: []entity.User{}},
	}

	_, err := collection.InsertOne(ctx, queue)
	if err != nil {
		return fmt.Errorf("couldn't create queue: %w", err)
	}

	return nil
}

func (db *Database) LogInOutToQueue(ctx context.Context, messageID string, user entity.User) error {
	if _, ok := db.queueMutex.Load(messageID); ok {
		return apperrors.NewCallbackError(nil, "queue is locked")
	}

	session, err := db.mongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("couldn't start session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) error {
		collection := db.mongoClient.Database("queue-bot").Collection("queues")
		filter := bson.D{{Key: "messageId", Value: messageID}}

		add := bson.D{
			{Key: "$addToSet", Value: bson.D{{Key: "users", Value: user}}},
		}

		addResult, err := collection.UpdateOne(sessionContext, filter, add)
		if err != nil {
			return fmt.Errorf("couldn't log add to queue: %w", err)
		}

		if addResult.ModifiedCount == 1 {
			return nil
		}

		remove := bson.D{
			{Key: "$pull", Value: bson.D{{Key: "users", Value: bson.D{{Key: "id", Value: user.ID}}}}},
		}

		_, err = collection.UpdateOne(sessionContext, filter, remove)
		if err != nil {
			return fmt.Errorf("couldn't log in/out to queue: %w", err)
		}

		return nil
	}

	err = mongo.WithSession(ctx, session, callback)
	if err != nil {
		return fmt.Errorf("error in transaction: %w", err)
	}

	return nil
}

func (db *Database) GetQueue(ctx context.Context, messageID string) (entity.Queue, error) {
	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	filter := bson.D{{Key: "messageId", Value: messageID}}

	var queue entity.Queue
	err := collection.FindOne(ctx, filter).Decode(&queue)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't get queue: %w", err)
	}

	return queue, nil
}

func (db *Database) StartQueue(ctx context.Context, messageID string, isShuffle bool) error {
	db.queueMutex.Store(messageID, struct{}{})
	defer db.queueMutex.Delete(messageID)

	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	filter := bson.D{{Key: "messageId", Value: messageID}}

	if !isShuffle {
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "currentUserIndex", Value: 0}}},
		}

		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("couldn't start queue without shuffle: %w", err)
		}

		return nil
	}

	queue, err := db.GetQueue(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't get queue: %w", err)
	}

	rand.New(rand.NewSource(time.Now().Unix()))

	rand.Shuffle(len(queue.Users), func(i, j int) {
		queue.Users[i], queue.Users[j] = queue.Users[j], queue.Users[i]
	})

	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "currentUserIndex", Value: 0},
				{Key: "users", Value: queue.Users},
			},
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("couldn't start queue with shuffle: %w", err)
	}

	return nil
}

func (db *Database) Close() error {
	db.mongoClient.Database("queue-bot").Collection("queues").Indexes()
	return db.mongoClient.Disconnect(context.Background())
}

func (db *Database) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "messageId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	collection := db.mongoClient.Database("queue-bot").Collection("queues")
	_, err := collection.Indexes().CreateOne(
		ctx,
		indexModel,
	)

	if err != nil {
		return fmt.Errorf("couldn't create index: %w", err)
	}

	return nil
}

func NewDatabase(ctx context.Context, connectString string) (*Database, error) {
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %w", err)
	}

	if err := connect.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("couldn't ping to database: %w", err)
	}
	slog.Info("Connected to MongoDB")

	mongoDB := &Database{mongoClient: connect}
	err = mongoDB.EnsureIndexes(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't ensure indexes: %w", err)
	}

	return mongoDB, nil
}
