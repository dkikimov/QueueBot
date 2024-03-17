package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	// Sqlite driver...
	_ "github.com/mattn/go-sqlite3"

	"QueueBot/internal/entity"
)

type Database struct {
	db *sql.DB
}

func (s Database) Close() error {
	return s.db.Close()
}

func (s Database) CreateQueue(ctx context.Context, messageID string, description string) error {
	createQueueStmt, err := s.db.PrepareContext(ctx, "INSERT INTO queues (message_id, description) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("couldn't prepare create queue statement: %w", err)
	}
	defer createQueueStmt.Close()

	_, err = createQueueStmt.ExecContext(ctx, messageID, description)

	return err
}

func (s Database) LogInOutToQueue(ctx context.Context, messageID string, user entity.User) error {
	logInOutStmt, err := s.db.PrepareContext(ctx, `INSERT INTO participants(message_id, user_id, user_name)
	VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP`)
	if err != nil {
		return fmt.Errorf("couldn't prepare log in/out to queue statement: %w", err)
	}
	defer logInOutStmt.Close()

	_, err = logInOutStmt.ExecContext(ctx, messageID, user.ID, user.Name)

	return err
}

func (s Database) GetQueue(ctx context.Context, messageID string) (entity.Queue, error) {
	descriptionStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT description, current_user_index FROM queues WHERE message_id = ?",
	)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't prepare get queue description statement: %w", err)
	}
	defer descriptionStmt.Close()

	getUsersStmt, err := s.db.PrepareContext(
		ctx,
		"SELECT user_id, user_name FROM participants WHERE message_id = ? and isDeleted = 0 ORDER BY order_number",
	)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't prepare get queue users statement: %w", err)
	}
	defer getUsersStmt.Close()

	var description string
	var currentUserIndex int
	queryResult := descriptionStmt.QueryRowContext(ctx, messageID)
	if err = queryResult.Scan(&description, &currentUserIndex); err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't scan description row in queue %s: %w", messageID, err)
	}

	rows, err := getUsersStmt.QueryContext(ctx, messageID)
	if err != nil {
		return entity.Queue{}, fmt.Errorf("couldn't get users from queue %s: %w", messageID, err)
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User
		if err = rows.Scan(&user.ID, &user.Name); err != nil {
			return entity.Queue{}, fmt.Errorf("couldn't scan user row in queue %s: %w", messageID, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return entity.Queue{}, fmt.Errorf("error during iterating rows from queue %s: %w", messageID, err)
	}

	return entity.Queue{
		MessageID:        messageID,
		Description:      description,
		Users:            users,
		CurrentPersonIdx: currentUserIndex,
	}, nil
}

func updateParticipantsInorder(ctx context.Context, tx *sql.Tx, messageID string) error {
	startStmt, err := tx.PrepareContext(ctx, `UPDATE participants SET order_number = dense_rank FROM 
                                                      (SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id FROM participants WHERE message_id = ?)
                                                          AS sub WHERE participants.user_id = sub.user_id AND message_id = ?`)
	if err != nil {
		return fmt.Errorf("couldn't prepare shuffle statement: %w", err)
	}
	defer startStmt.Close()

	_, err = startStmt.ExecContext(ctx, messageID, messageID)
	if err != nil {
		return fmt.Errorf("couldn't start queue: %w", err)
	}

	return nil
}

func updateParticipantsShuffle(ctx context.Context, tx *sql.Tx, messageID string) error {
	startStmt, err := tx.PrepareContext(ctx, `UPDATE participants
														SET order_number = random()
														WHERE message_id = ?;`)
	if err != nil {
		return fmt.Errorf("couldn't prepare inorder queue start statement: %w", err)
	}
	defer startStmt.Close()

	_, err = startStmt.ExecContext(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't start queue: %w", err)
	}

	return nil
}

func (s Database) StartQueue(ctx context.Context, messageID string, isShuffle bool) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}
	defer tx.Rollback()

	setCurrentUserIndexStmt, err := tx.PrepareContext(ctx, "UPDATE queues SET current_user_index = 0 WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare set current user index statement: %w", err)
	}
	defer setCurrentUserIndexStmt.Close()

	_, err = setCurrentUserIndexStmt.ExecContext(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't set current user index: %w", err)
	}

	if isShuffle {
		err := updateParticipantsShuffle(ctx, tx, messageID)
		if err != nil {
			return fmt.Errorf("couldn't update participant in shuffle order: %w", err)
		}
	} else {
		err := updateParticipantsInorder(ctx, tx, messageID)
		if err != nil {
			return fmt.Errorf("couldn't update participant inorder: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}

func (s Database) IncrementCurrentPerson(ctx context.Context, messageID string) error {
	incrementStmt, err := s.db.PrepareContext(ctx, "UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare increment current person statement: %w", err)
	}
	defer incrementStmt.Close()

	_, err = incrementStmt.ExecContext(ctx, messageID)

	return err
}

func (s Database) DeleteQueue(ctx context.Context, messageID string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteQueueStmt, err := tx.PrepareContext(ctx, "DELETE FROM queues WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare finish queue statement: %w", err)
	}
	defer deleteQueueStmt.Close()

	_, err = deleteQueueStmt.ExecContext(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't finish queue: %w", err)
	}

	deleteParticipantsStmt, err := tx.PrepareContext(ctx, "DELETE FROM participants WHERE message_id = ?")
	if err != nil {
		return fmt.Errorf("couldn't prepare delete participants statement: %w", err)
	}
	defer deleteParticipantsStmt.Close()

	_, err = deleteParticipantsStmt.ExecContext(ctx, messageID)
	if err != nil {
		return fmt.Errorf("couldn't delete participants: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return nil
}

func NewDatabase(databasePath string) (*Database, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared", databasePath))
	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %w", err)
	}

	if _, err := db.Exec(CreateTables); err != nil {
		return nil, fmt.Errorf("couldn't create default sqlite tables: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func NewDatabaseFromDB(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}
