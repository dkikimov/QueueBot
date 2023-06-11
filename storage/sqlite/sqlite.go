package sqlite

import (
	"QueueBot/logger"
	"QueueBot/storage"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type Commands struct {
	createQueueStmt *sql.Stmt
}

type SQLite struct {
	db       *sql.DB
	mu       sync.Mutex
	commands Commands
}

func (sqlite *SQLite) GetUsersInQueue(queueId string) ([]storage.User, error) {
	//TODO implement me
	panic("implement me")
}

func (sqlite *SQLite) AddUserToQueue(queueId string, user storage.User) error {
	//TODO implement me
	panic("implement me")
}

func (sqlite *SQLite) DeleteUserFromQueueById(queueId string, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (sqlite *SQLite) CreateQueue(queryID string, description string) error {
	_, err := sqlite.commands.createQueueStmt.Exec(queryID, description)
	return err
}
func NewDatabase() *SQLite {
	db, err := sql.Open("sqlite3", "./database.sqlite3")
	if err != nil {
		logger.Panicf("Couldn't open database with error %s", err.Error())
	}

	if _, err := db.Exec(CreateTables); err != nil {
		logger.Panicf("Couldn't create default sqlite tables with error %s", err.Error())
	}

	return &SQLite{
		db:       db,
		commands: getPreparedCommands(db),
	}
}

func getPreparedCommands(db *sql.DB) Commands {
	createQueueStmt, err := db.Prepare(CreateQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare create queue command with error %s", err.Error())
	}

	return Commands{
		createQueueStmt: createQueueStmt,
	}
}
