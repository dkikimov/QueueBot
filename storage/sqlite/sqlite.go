package sqlite

import (
	"QueueBot/logger"
	"QueueBot/telegram/steps"
	"QueueBot/user"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type Commands struct {
	createQueueStmt,
	createUserStmt,
	setUserCurrentStepStmt,
	getUserCurrentStepStmt,
	getUsersInQueueStmt,
	getDescriptionOfQueueStmt,
	addUserToQueueStmt,
	removeUserFromQueueStmt,
	countMatchesInParticipantsStmt,
	startQueueStmt,
	incCurrentPersonStmt *sql.Stmt
}

type SQLite struct {
	db       *sql.DB
	mu       sync.Mutex
	commands Commands
}

func (sqlite *SQLite) IncrementCurrentPerson(messageId string) (err error, currentPerson int) {
	row := sqlite.commands.incCurrentPersonStmt.QueryRow(messageId)
	if err = row.Scan(&currentPerson); err != nil {
		return err, 0
	}

	return err, currentPerson
}

func (sqlite *SQLite) StartQueue(messageId string) (error, bool) {
	row := sqlite.commands.startQueueStmt.QueryRow(messageId)
	var wasUpdated int
	if err := row.Scan(&wasUpdated); err != nil {
		return err, false
	}

	return nil, wasUpdated == 1
}

func (sqlite *SQLite) LogInOurOutQueue(messageId string, user user.User) (err error) {
	row := sqlite.commands.countMatchesInParticipantsStmt.QueryRow(user.Id, messageId)
	var count int
	if err = row.Scan(&count); err != nil {
		return err
	}

	if count == 1 {
		_, err = sqlite.commands.removeUserFromQueueStmt.Exec(user.Id, messageId)
		logger.Printf("Removed user with id %s", user.Id)
	} else {
		_, err = sqlite.commands.addUserToQueueStmt.Exec(messageId, user.Id, user.Name)
		logger.Printf("Added user user with id %s", user.Id)

	}
	return err
}

func (sqlite *SQLite) GetDescriptionOfQueue(messageId string) (description string, err error) {
	result := sqlite.commands.getDescriptionOfQueueStmt.QueryRow(messageId)
	if err = result.Scan(&description); err != nil {
		return "", err
	}
	return description, err
}

func (sqlite *SQLite) GetUserCurrentStep(userId int64) (currentStep steps.Step, err error) {
	result := sqlite.commands.getUserCurrentStepStmt.QueryRow(userId)
	if err = result.Scan(&currentStep); err != nil {
		return 0, err
	}
	return currentStep, err
}

func (sqlite *SQLite) CreateUser(userId int64) error {
	_, err := sqlite.commands.createUserStmt.Exec(userId)
	return err
}

func (sqlite *SQLite) SetUserCurrentStep(userId int64, currentStep steps.Step) error {
	_, err := sqlite.commands.setUserCurrentStepStmt.Exec(int(currentStep), userId)
	return err
}

func (sqlite *SQLite) GetUsersInQueue(messageId string) ([]user.User, error) {
	rows, err := sqlite.commands.getUsersInQueueStmt.Query(messageId)
	if err != nil {
		return nil, err
	}
	//TODO: handle error
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var user user.User
		if err = rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, err
}

func (sqlite *SQLite) CreateQueue(messageId string, description string) error {
	_, err := sqlite.commands.createQueueStmt.Exec(messageId, description)
	return err
}

func NewDatabase() *SQLite {
	//db, err := sql.Open("sqlite3", "file:./database.sqlite3:memory:?cache=shared")
	db, err := sql.Open("sqlite3", "./database.sqlite3?cache=shared")

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
		logger.Panicf("Couldn't prepare create queue command with error: %s", err.Error())
	}

	createUserStmt, err := db.Prepare(CreateUser)
	if err != nil {
		logger.Panicf("Couldn't prepare create user command with error: %s", err.Error())
	}

	setUserCurrentStepStmt, err := db.Prepare(SetUserCurrentStep)
	if err != nil {
		logger.Panicf("Couldn't prepare set user next step command with error: %s", err.Error())
	}

	getUserCurrentStepStmt, err := db.Prepare(GetUserCurrentStep)
	if err != nil {
		logger.Panicf("Couldn't prepare get user next step command with error: %s", err.Error())
	}

	getUsersInQueueStmt, err := db.Prepare(GetUsersInQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare get users in queue command with error: %s", err.Error())
	}

	getDescriptionOfQueueStmt, err := db.Prepare(GetDescriptionOfQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare get description of queue command with error: %s", err.Error())
	}

	addUserToQueueStmt, err := db.Prepare(AddUserToQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare add user to queue command with error: %s", err.Error())
	}

	countMatchesInParticipantsStmt, err := db.Prepare(CountMatchesInParticipants)
	if err != nil {
		logger.Panicf("Couldn't prepare count matches in participants command with error: %s", err.Error())
	}

	removeUserFromQueueStmt, err := db.Prepare(RemoveUserFromQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare remove user from queue command with error: %s", err.Error())
	}

	startQueueStmt, err := db.Prepare(StartQueue)
	if err != nil {
		logger.Panicf("Couldn't prepare start queue command with error: %s", err.Error())
	}

	incCurrentPersonStmt, err := db.Prepare(IncrementCurrentPerson)
	if err != nil {
		logger.Panicf("Couldn't prepare increment current person command with error: %s", err.Error())
	}

	return Commands{
		createQueueStmt:                createQueueStmt,
		createUserStmt:                 createUserStmt,
		setUserCurrentStepStmt:         setUserCurrentStepStmt,
		getUserCurrentStepStmt:         getUserCurrentStepStmt,
		getUsersInQueueStmt:            getUsersInQueueStmt,
		getDescriptionOfQueueStmt:      getDescriptionOfQueueStmt,
		addUserToQueueStmt:             addUserToQueueStmt,
		countMatchesInParticipantsStmt: countMatchesInParticipantsStmt,
		removeUserFromQueueStmt:        removeUserFromQueueStmt,
		startQueueStmt:                 startQueueStmt,
		incCurrentPersonStmt:           incCurrentPersonStmt,
	}
}
