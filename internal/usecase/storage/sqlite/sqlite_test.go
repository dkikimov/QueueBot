package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"QueueBot/internal/entity"
)

var errAlreadyExists = errors.New("already exists")
var errReference = errors.New("foreign key constraint failed")

func TestDatabase_CreateQueue(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID   string
		description string
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				messageID:   "123",
				description: "Test",
			},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("INSERT INTO queues").WillBeClosed()

				mock.
					ExpectExec("INSERT INTO queues").
					WithArgs("123", "Test").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Adding existing message ID",
			args: args{
				messageID:   "1234",
				description: "Test",
			},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("INSERT INTO queues").WillBeClosed()

				mock.
					ExpectExec("INSERT INTO queues").
					WithArgs("1234", "Test").
					WillReturnError(errAlreadyExists)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			err := db.CreateQueue(context.Background(), tt.args.messageID, tt.args.description)
			if tt.wantErr && err == nil {
				t.Errorf("Expected CreateQueue() to return error = %v, returned %v", tt.wantErr, err)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("CreateQueue() returned error = %v, expected %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_DeleteQueue(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("DELETE FROM queues WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("DELETE FROM queues WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectPrepare("DELETE FROM participants WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("DELETE FROM participants WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			if err := db.DeleteQueue(context.Background(), tt.args.messageID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteQueue() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_GetQueue(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		want          entity.Queue
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
			},
			want: entity.Queue{
				MessageID:        "123",
				Description:      "Test",
				CurrentPersonIdx: 0,
				Users: []entity.User{
					{
						ID:   1,
						Name: "Test",
					},
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("SELECT description, current_user_index FROM queues WHERE message_id = ?").WillBeClosed()
				mock.ExpectPrepare("SELECT user_id, user_name FROM participants WHERE message_id = ? and isDeleted = 0 ORDER BY order_number").WillBeClosed()

				rows := sqlmock.NewRows([]string{"description", "current_user_index"}).
					AddRow("Test", 0)

				mock.ExpectQuery("SELECT description, current_user_index FROM queues WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"user_id", "user_name"}).
					AddRow(1, "Test")

				mock.ExpectQuery("SELECT user_id, user_name FROM participants WHERE message_id = ? and isDeleted = 0 ORDER BY order_number").
					WithArgs(args.messageID).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Queue not found",
			args: args{
				messageID: "123",
			},
			want: entity.Queue{},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("SELECT description, current_user_index FROM queues WHERE message_id = ?").WillBeClosed()
				mock.ExpectPrepare("SELECT user_id, user_name FROM participants WHERE message_id = ? and isDeleted = 0 ORDER BY order_number").WillBeClosed()

				mock.ExpectQuery("SELECT description, current_user_index FROM queues WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			got, err := db.GetQueue(context.Background(), tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueue() got = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_IncrementCurrentPerson(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		wantErr       bool
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
			},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Queue not found",
			args: args{
				messageID: "123",
			},
			mockBehaviour: func(args args) {
				mock.ExpectPrepare("UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = current_user_index + 1 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			if err := db.IncrementCurrentPerson(context.Background(), tt.args.messageID); (err != nil) != tt.wantErr {
				t.Errorf("IncrementCurrentPerson() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_LogInOutToQueue(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
		user      entity.User
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		wantErr       bool
		mockBehaviour mockBehaviour
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
				user: entity.User{
					ID:   1,
					Name: "Test",
				},
			},
			wantErr: false,
			mockBehaviour: func(args args) {
				mock.ExpectPrepare(`INSERT INTO participants(message_id, user_id, user_name)
												VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP`).
					WillBeClosed()

				mock.ExpectExec(`INSERT INTO participants(message_id, user_id, user_name)
												VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP`).
					WithArgs(args.messageID, args.user.ID, args.user.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Unknown message ID",
			args: args{
				messageID: "1234",
				user: entity.User{
					ID:   1,
					Name: "Test",
				},
			},
			wantErr: true,
			mockBehaviour: func(args args) {
				mock.ExpectPrepare(`INSERT INTO participants(message_id, user_id, user_name)
												VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP`).
					WillBeClosed()

				mock.ExpectExec(`INSERT INTO participants(message_id, user_id, user_name)
												VALUES (?, ?, ?) on conflict do update set isDeleted=not isDeleted, joined_at=CURRENT_TIMESTAMP`).
					WithArgs(args.messageID, args.user.ID, args.user.Name).
					WillReturnError(errReference)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			if err := db.LogInOutToQueue(context.Background(), tt.args.messageID, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("LogInOutToQueue() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_StartQueue(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
		isShuffle bool
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK straight order",
			args: args{
				messageID: "123",
				isShuffle: false,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectPrepare(`UPDATE participants SET order_number = dense_rank FROM 
                                                      (SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
                                                       FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
                                                                                                        AND message_id = $1`).WillBeClosed()
				mock.ExpectExec(`UPDATE participants SET order_number = dense_rank FROM 
                                                      (SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
                                                       FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
                                                                                                        AND message_id = $1`).
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "OK shuffle order",
			args: args{
				messageID: "123",
				isShuffle: true,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectPrepare(`UPDATE participants
														SET order_number = random()
														WHERE message_id = ?;`).WillBeClosed()
				mock.ExpectExec(`UPDATE participants
														SET order_number = random()
														WHERE message_id = ?;`).
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Unknown message ID",
			args: args{
				messageID: "1234",
				isShuffle: false,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnError(errReference)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Shuffle error",
			args: args{
				messageID: "123",
				isShuffle: true,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectPrepare("UPDATE participants SET order_number = random() WHERE message_id = ?;").WillBeClosed()
				mock.ExpectExec("UPDATE participants SET order_number = random() WHERE message_id = ?;").
					WithArgs(args.messageID).
					WillReturnError(errReference)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Straight order error",
			args: args{
				messageID: "123",
				isShuffle: false,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").WillBeClosed()
				mock.ExpectExec("UPDATE queues SET current_user_index = 0 WHERE message_id = ?").
					WithArgs(args.messageID).
					WillReturnError(errReference)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			if err := db.StartQueue(context.Background(), tt.args.messageID, tt.args.isShuffle); (err != nil) != tt.wantErr {
				t.Errorf("StartQueue() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_updateParticipantsInorder(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)

	type args struct {
		messageID string
		isShuffle bool
	}

	type mockBehaviour func(args args)

	tests := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK direct order",
			args: args{
				messageID: "123",
				isShuffle: false,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare(`UPDATE participants
														SET order_number = dense_rank FROM 
														(SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
														FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
														AND message_id = $1`).WillBeClosed()
				mock.ExpectExec(`UPDATE participants
														SET order_number = dense_rank FROM 
														(SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
														FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
														AND message_id = $1`).
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Unknown message ID direct order",
			args: args{
				messageID: "1234",
				isShuffle: false,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare(`UPDATE participants
														SET order_number = dense_rank FROM 
														(SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
														FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
														AND message_id = $1`).WillBeClosed()
				mock.ExpectExec(`UPDATE participants
														SET order_number = dense_rank FROM 
														(SELECT dense_rank() OVER (ORDER BY joined_at) AS dense_rank, user_id 
														FROM participants WHERE message_id = $1) AS sub WHERE participants.user_id = sub.user_id 
														AND message_id = $1`).
					WithArgs(args.messageID).
					WillReturnError(errReference)
			},
			wantErr: true,
		},
		{
			name: "OK Shuffle",
			args: args{
				messageID: "123",
				isShuffle: true,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE participants SET order_number = random() WHERE message_id = ?;").WillBeClosed()
				mock.ExpectExec("UPDATE participants SET order_number = random() WHERE message_id = ?;").
					WithArgs(args.messageID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Unknown message ID shuffle",
			args: args{
				messageID: "1234",
				isShuffle: true,
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectPrepare("UPDATE participants SET order_number = random() WHERE message_id = ?;").WillBeClosed()
				mock.ExpectExec("UPDATE participants SET order_number = random() WHERE message_id = ?;").
					WithArgs(args.messageID).
					WillReturnError(errReference)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{})
			assert.NoError(t, err)

			if err := setParticipantsOrder(context.Background(), tx, tt.args.messageID, tt.args.isShuffle); (err != nil) != tt.wantErr {
				t.Errorf("setParticipantsOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	db := NewDatabaseFromDB(mockDB)
	type mockBehaviour func()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name: "OK",
			mockBehaviour: func() {
				mock.ExpectClose()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			if err = db.Close(); (err != nil) != tt.wantErr {
				t.Errorf("db.Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
