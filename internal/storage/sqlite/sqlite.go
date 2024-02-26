package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"QueueBot/internal/models"
	"QueueBot/internal/steps"
)

type Commands struct {
	createQueueStmt,
	createUserStmt,
	setUserCurrentStepStmt,
	getUserCurrentStepStmt,
	getUsersInQueueStmt,
	getUsersInQueueShuffledStmt,
	getDescriptionOfQueueStmt,
	addUserToQueueStmt,
	removeUserFromQueueStmt,
	countMatchesInParticipantsStmt,
	startQueueStmt,
	incCurrentPersonStmt,
	isQueueShuffledStmt,
	shuffleUsersStmt,
	endQueueStmt,
	deleteRelativeParticipants,
	goToMenuStmt *sql.Stmt
}

type SQLite struct {
	db       *sql.DB
	commands *Commands
}

func (sqlite *SQLite) Close() error {
	if err := sqlite.commands.createQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.createUserStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.setUserCurrentStepStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.getUsersInQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.getUsersInQueueShuffledStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.getDescriptionOfQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.addUserToQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.countMatchesInParticipantsStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.incCurrentPersonStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.isQueueShuffledStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.shuffleUsersStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.removeUserFromQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.endQueueStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.deleteRelativeParticipants.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.goToMenuStmt.Close(); err != nil {
		return err
	}
	if err := sqlite.commands.startQueueStmt.Close(); err != nil {
		return err
	}
	return sqlite.db.Close()
}

func (sqlite *SQLite) FinishQueueDeleteParticipants(messageId string) error {
	_, err := sqlite.commands.endQueueStmt.Exec(messageId)
	if err != nil {
		return fmt.Errorf("couldn't end queue id %s: %s", messageId, err)
	}

	_, err = sqlite.commands.deleteRelativeParticipants.Exec(messageId)
	if err != nil {
		return fmt.Errorf("couldn't delete participants from ended queue id %s: %s", messageId, err)
	}

	return err
}

func (sqlite *SQLite) ShuffleUsers(messageId string) error {
	_, err := sqlite.commands.shuffleUsersStmt.Exec(messageId)
	if err != nil {
		return fmt.Errorf("couldn't shuffle users in queue with id %s: %s", messageId, err)
	}
	return nil
}

func (sqlite *SQLite) GetUsersInQueueCheckShuffle(messageId string) (users []models.User, err error) {
	row := sqlite.commands.isQueueShuffledStmt.QueryRow(messageId)

	var isShuffled int
	if err := row.Scan(&isShuffled); err != nil {
		return nil, fmt.Errorf("couldn't scan isShuffled row in queue with id %s: %s", messageId, err)
	}

	var rows *sql.Rows
	if isShuffled == 1 {
		rows, err = sqlite.commands.getUsersInQueueShuffledStmt.Query(messageId)
		if err != nil {
			return nil, fmt.Errorf("couldn't get users in shuffled queue with id %s: %s", messageId, err)
		}
	} else {
		rows, err = sqlite.commands.getUsersInQueueStmt.Query(messageId)
		if err != nil {
			return nil, fmt.Errorf("couldn't get users in queue with id %s: %s", messageId, err)
		}
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("couldn't close rows in GetUsersInQueueCheckShuffle, messageId %s, isShuffled %t", messageId, isShuffled == 1)
		}
	}(rows)

	for rows.Next() {
		var currentUser models.User
		if err = rows.Scan(&currentUser.Id, &currentUser.Name); err != nil {
			return nil, fmt.Errorf("couldn't scan users in queue with id %s: %s", messageId, err)
		}
		users = append(users, currentUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iterating over users in queue with id %s: %s", messageId, err)
	}

	return users, err
}

func (sqlite *SQLite) StartQueue(messageId string, isShuffle bool) (error, bool) {
	row := sqlite.commands.startQueueStmt.QueryRow(isShuffle, messageId)
	var wasUpdated int
	if err := row.Scan(&wasUpdated); err != nil {
		return err, false
	}

	if isShuffle && wasUpdated == 1 {
		if err := sqlite.ShuffleUsers(messageId); err != nil {
			return fmt.Errorf("couldn't shuffle users in queue with id %s: %s", messageId, err), false
		}
	}
	return nil, wasUpdated == 1
}

func (sqlite *SQLite) GoToMenu(messageId string) error {
	_, err := sqlite.commands.goToMenuStmt.Exec(messageId)
	if err != nil {
		return fmt.Errorf("couldn't go to menu in queue with id %s: %s", messageId, err)
	}
	return nil
}

func (sqlite *SQLite) IncrementCurrentPerson(messageId string) (err error, currentPerson int) {
	row := sqlite.commands.incCurrentPersonStmt.QueryRow(messageId)
	if err = row.Scan(&currentPerson); err != nil {
		return fmt.Errorf("couldn't increment current person in queue %s: %s", messageId, err), 0
	}

	return err, currentPerson
}

func (sqlite *SQLite) LogInOurOutQueue(messageId string, user models.User) error {
	row := sqlite.commands.countMatchesInParticipantsStmt.QueryRow(user.Id, messageId)
	var count int
	if err := row.Scan(&count); err != nil {
		return fmt.Errorf("couldn't scan row: %s", err)
	}

	if count == 1 {
		_, err := sqlite.commands.removeUserFromQueueStmt.Exec(user.Id, messageId)
		if err != nil {
			return fmt.Errorf("couldn't remove user from queue: %s", err)
		}
	} else {
		_, err := sqlite.commands.addUserToQueueStmt.Exec(messageId, user.Id, user.Name)
		if err != nil {
			return fmt.Errorf("couldn't add user from queue: %s", err)
		}
	}
	return nil
}

func (sqlite *SQLite) GetDescriptionOfQueue(messageId string) (description string, err error) {
	result := sqlite.commands.getDescriptionOfQueueStmt.QueryRow(messageId)
	if err = result.Scan(&description); err != nil {
		return "", fmt.Errorf("couldn't scan description row in queue %s: %s", messageId, err)
	}
	return description, nil
}

func (sqlite *SQLite) GetUserCurrentStep(userId int64) (currentStep steps.ChatStep, err error) {
	result := sqlite.commands.getUserCurrentStepStmt.QueryRow(userId)
	if err = result.Scan(&currentStep); err != nil {
		return 0, fmt.Errorf("couldn't get current user %d step queue: %s", currentStep, err)
	}

	return currentStep, err
}

func (sqlite *SQLite) CreateUser(userId int64) error {
	_, err := sqlite.commands.createUserStmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("couldn't create user %d: %s", userId, err)
	}

	return nil
}

func (sqlite *SQLite) SetUserCurrentStep(userId int64, currentStep steps.ChatStep) error {
	_, err := sqlite.commands.setUserCurrentStepStmt.Exec(int(currentStep), userId)
	if err != nil {
		return fmt.Errorf("couldn't set person %d current step to %d: %s", userId, currentStep, err)
	}
	return err
}

func (sqlite *SQLite) GetUsersInQueue(messageId string) (users []models.User, err error) {
	rows, err := sqlite.commands.getUsersInQueueStmt.Query(messageId)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("couldn't close rows in GetUsersInQueue, messageId %s", messageId)
		}
	}(rows)

	for rows.Next() {
		var currentUser models.User
		if err = rows.Scan(&currentUser.Id, &currentUser.Name); err != nil {
			return nil, fmt.Errorf("couldn't get users in queue %s: %s", messageId, err)
		}
		users = append(users, currentUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iterating over users in queue with id %s: %s", messageId, err)
	}

	return
}

func (sqlite *SQLite) CreateQueue(messageId string, description string) error {
	_, err := sqlite.commands.createQueueStmt.Exec(messageId, description)
	if err != nil {
		return fmt.Errorf("couldn't create queue for message %s: %s", messageId, err)
	}
	return nil
}

func NewDatabase() (*SQLite, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite3?cache=shared")

	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %s", err)
	}

	if _, err := db.Exec(CreateTables); err != nil {
		return nil, fmt.Errorf("couldn't create default sqlite tables: %s", err)
	}

	commands, err := getPreparedCommands(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't get prepared commands")
	}

	return &SQLite{
		db:       db,
		commands: commands,
	}, nil
}

func getPreparedCommands(db *sql.DB) (*Commands, error) {
	createQueueStmt, err := db.Prepare(CreateQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare create queue command with error: %s", err)
	}

	createUserStmt, err := db.Prepare(CreateUser)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare create user command with error: %s", err)
	}

	setUserCurrentStepStmt, err := db.Prepare(SetUserCurrentStep)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare set user next step command with error: %s", err)
	}

	getUserCurrentStepStmt, err := db.Prepare(GetUserCurrentStep)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare get user next step command with error: %s", err)
	}

	getUsersInQueueStmt, err := db.Prepare(GetUsersInQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare get users in queue command with error: %s", err)
	}

	getUsersInQueueShuffledStmt, err := db.Prepare(GetUsersInQueueShuffled)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare get users in queue shuffled command with error: %s", err)
	}

	getDescriptionOfQueueStmt, err := db.Prepare(GetDescriptionOfQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare get description of queue command with error: %s", err)
	}

	addUserToQueueStmt, err := db.Prepare(AddUserToQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare add user to queue command with error: %s", err)
	}

	countMatchesInParticipantsStmt, err := db.Prepare(CountMatchesInParticipants)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare count matches in participants command with error: %s", err)
	}

	removeUserFromQueueStmt, err := db.Prepare(RemoveUserFromQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare remove user from queue command with error: %s", err)
	}

	startQueueStmt, err := db.Prepare(StartQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare start queue command with error: %s", err)
	}

	incCurrentPersonStmt, err := db.Prepare(IncrementCurrentPerson)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare increment current person command with error: %s", err)
	}

	resetCurrentPersonStmt, err := db.Prepare(GoToMenu)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare reset current person command with error: %s", err)
	}

	isQueueShuffledStmt, err := db.Prepare(IsQueueShuffled)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare is queue shuffled command with error: %s", err)
	}

	shuffleUsersStmt, err := db.Prepare(ShuffleUsers)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare shuffle users command with error: %s", err)
	}

	endQueueStmt, err := db.Prepare(EndQueue)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare end queue command with error: %s", err)
	}

	deleteRelativeParticipants, err := db.Prepare(DeleteRelativeParticipants)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare delete relative participants command with error: %s", err)
	}

	return &Commands{
		createQueueStmt:                createQueueStmt,
		createUserStmt:                 createUserStmt,
		setUserCurrentStepStmt:         setUserCurrentStepStmt,
		getUserCurrentStepStmt:         getUserCurrentStepStmt,
		getUsersInQueueStmt:            getUsersInQueueStmt,
		getUsersInQueueShuffledStmt:    getUsersInQueueShuffledStmt,
		getDescriptionOfQueueStmt:      getDescriptionOfQueueStmt,
		addUserToQueueStmt:             addUserToQueueStmt,
		countMatchesInParticipantsStmt: countMatchesInParticipantsStmt,
		removeUserFromQueueStmt:        removeUserFromQueueStmt,
		startQueueStmt:                 startQueueStmt,
		incCurrentPersonStmt:           incCurrentPersonStmt,
		goToMenuStmt:                   resetCurrentPersonStmt,
		isQueueShuffledStmt:            isQueueShuffledStmt,
		shuffleUsersStmt:               shuffleUsersStmt,
		endQueueStmt:                   endQueueStmt,
		deleteRelativeParticipants:     deleteRelativeParticipants,
	}, nil
}
