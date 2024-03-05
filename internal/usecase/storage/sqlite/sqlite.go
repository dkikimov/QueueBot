package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"QueueBot/internal/entity"
)

type Database struct {
	db *sql.DB
}

func (s Database) Close() error {
	return s.db.Close()
}

func (s Database) CreateQueue(ctx context.Context, messageId string, description string) error {
	createQueueStmt, err := s.db.PrepareContext(ctx, "INSERT INTO queues (message_id, description) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("couldn't prepare create queue statement: %w", err)
	}
	_, err = createQueueStmt.ExecContext(ctx, messageId, description)

	return err
}

func (s Database) LogInOutToQueue(ctx context.Context, messageId string, user entity.User) error {
	logInOutStmt, err := s.db.PrepareContext(ctx, `INSERT INTO participants(message_id, user_id, user_name)
	VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP;`)

	if err != nil {
		return fmt.Errorf("couldn't prepare log in/out to queue statement: %w", err)
	}

	_, err = logInOutStmt.ExecContext(ctx, messageId, user.Id, user.Name)

	return err
}

func (s Database) GetQueue(ctx context.Context, messageId string) (entity.Queue, error) {
	descriptionStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT description, current_user_index FROM queues WHERE message_id = ?",
	)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't prepare get queue description statement: %w", err)
	}

	getUsersStmt, err := s.db.PrepareContext(
		ctx,
		`SELECT user_id, user_name FROM participants WHERE message_id = ? and isDeleted = 0 
                                            			 ORDER BY order_number`,
	)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't prepare get queue users statement: %w", err)
	}

	var description string
	var currentUserIndex int
	queryResult := descriptionStmt.QueryRowContext(ctx, messageId)
	if err = queryResult.Scan(&description, &currentUserIndex); err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't scan description row in queue %s: %s", messageId, err)
	}

	rows, err := getUsersStmt.QueryContext(ctx, messageId)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't get users from queue %s: %s", messageId, err)
	}

	defer rows.Close()
	var users []entity.User

	for rows.Next() {
		var user entity.User
		if err = rows.Scan(&user.Id, &user.Name); err != nil {
			return entity.Queue{}, fmt.Errorf("couldn't scan user row in queue %s: %s", messageId, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return entity.Queue{}, fmt.Errorf("error during iterating rows from queue %s: %s", messageId, err)
	}

	return entity.Queue{
		MessageID:        messageId,
		Description:      description,
		Users:            users,
		CurrentPersonIdx: currentUserIndex,
	}, nil
}

func (s Database) StartQueue(ctx context.Context, messageId string, isShuffle bool) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	setCurrentUserIndexStmt, err := tx.PrepareContext(ctx, "UPDATE queues SET current_user_index = 0 WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare set current user index statement: %w", err)
	}

	_, err = setCurrentUserIndexStmt.ExecContext(ctx, messageId)
	if err != nil {
		return fmt.Errorf("couldn't set current user index: %w", err)
	}

	if isShuffle == false {
		startStmt, err := tx.PrepareContext(ctx, `UPDATE participants SET order_number = dense_rank FROM 
                                                      (SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id FROM participants WHERE message_id = ?)
                                                          AS sub WHERE participants.user_id = sub.user_id AND message_id = ?`)
		if err != nil {
			return fmt.Errorf("couldn't prepare shuffle statement: %w", err)
		}

		_, err = startStmt.ExecContext(ctx, messageId, messageId)
		if err != nil {
			return fmt.Errorf("couldn't start queue: %w", err)
		}
	} else {
		startStmt, err := tx.PrepareContext(ctx, `UPDATE participants
														SET order_number = random()
														WHERE message_id = ?;`)
		if err != nil {
			return fmt.Errorf("couldn't prepare inorder queue start statement: %w", err)
		}

		_, err = startStmt.ExecContext(ctx, messageId)
		if err != nil {
			return fmt.Errorf("couldn't start queue: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}

func (s Database) IncrementCurrentPerson(ctx context.Context, messageId string) error {
	incrementStmt, err := s.db.PrepareContext(ctx, "UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare increment current person statement: %w", err)
	}

	_, err = incrementStmt.ExecContext(ctx, messageId)

	return err
}

func (s Database) DeleteQueue(ctx context.Context, messageId string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}

	deleteQueueStmt, err := tx.PrepareContext(ctx, "DELETE FROM queues WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare finish queue statement: %w", err)
	}

	_, err = deleteQueueStmt.ExecContext(ctx, messageId)
	if err != nil {
		return fmt.Errorf("couldn't finish queue: %w", err)
	}

	deleteParticipantsStmt, err := tx.PrepareContext(ctx, "DELETE FROM participants WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare delete participants statement: %w", err)
	}

	_, err = deleteParticipantsStmt.ExecContext(ctx, messageId)
	if err != nil {
		return fmt.Errorf("couldn't delete participants: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}

func NewDatabase(databasePath string) (*Database, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", databasePath))

	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %s", err)
	}

	if _, err := db.Exec(CreateTables); err != nil {
		return nil, fmt.Errorf("couldn't create default sqlite tables: %s", err)
	}

	return &Database{
		db: db,
	}, nil
}
