package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"QueueBot/internal/usecase/storage/mongodb"
)

type TestDatabase struct {
	Instance  *mongodb.Database
	Client    *mongo.Client
	Address   string
	container testcontainers.Container
}

func SetupTestDatabase() (*TestDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, dbInstance, connection, dbAddr, err := createMongoContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup test: %w", err)
	}

	return &TestDatabase{
		container: container,
		Client:    connection,
		Instance:  dbInstance,
		Address:   dbAddr,
	}, nil
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func createMongoContainer(ctx context.Context) (testcontainers.Container, *mongodb.Database, *mongo.Client, string, error) {
	var env = map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": "root",
		"MONGO_INITDB_ROOT_PASSWORD": "pass",
		"MONGO_INITDB_DATABASE":      "testdb",
	}
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:7.0.2",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, nil, "", fmt.Errorf("failed to start container: %w", err)
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, nil, nil, "", fmt.Errorf("failed to get container external port: %w", err)
	}

	log.Println("mongo container ready and running at port: ", p.Port())

	uri := fmt.Sprintf("mongodb://root:pass@localhost:%s", p.Port())
	connection, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return container, nil, nil, "", fmt.Errorf("failed to connect to mongo: %w", err)
	}

	db, err := mongodb.NewDatabaseFromClient(ctx, connection)
	if err != nil {
		return container, db, nil, uri, fmt.Errorf("failed to establish database connection: %w", err)
	}

	return container, db, connection, uri, nil
}
