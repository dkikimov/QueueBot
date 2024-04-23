package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"QueueBot/internal/entity"
	"QueueBot/internal/usecase/storage/mongodb"
)

var (
	testDBInstance *mongodb.Database
)

func TestMain(m *testing.M) {
	log.Println("setup is running")
	testDB, err := SetupTestDatabase()
	if err != nil {
		panic(fmt.Sprintf("couldn't setup test db: %v", err))
	}

	testDBInstance = testDB.Instance

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err = populateDB(ctx, testDB.Client); err != nil {
		panic(fmt.Sprintf("couldn't populate db: %v", err))
	}
	cancel()

	exitVal := m.Run()
	defer os.Exit(exitVal)

	log.Println("teardown is running")
	_ = testDB.container.Terminate(context.Background())
}

func TestGetQueue(t *testing.T) {
	type args struct {
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Queue
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
			},
			want: entity.Queue{
				MessageID:        "123",
				Description:      "123",
				CurrentPersonIdx: 0,
				Users: []entity.User{
					{
						ID:   1,
						Name: "Username",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not found",
			args: args{
				messageID: "1234",
			},
			want:    entity.Queue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			got, err := testDBInstance.GetQueue(ctx, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueue() got = %v, want %v", got, tt.want)
			}
		})
	}

	queue, err := testDBInstance.GetQueue(context.Background(), "123")
	assert.NoError(t, err)

	assert.Equal(t, "123", queue.MessageID)
	assert.Equal(t, "123", queue.Description)
	assert.Equal(t, 0, queue.CurrentPersonIdx)
	assert.Equal(t, len(queue.Users), 1)

	assert.EqualValues(t, 1, queue.Users[0].ID)
	assert.Equal(t, "Username", queue.Users[0].Name)
}

func TestIncrementCurrentPerson(t *testing.T) {
	type args struct {
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				messageID: "123",
			},
			wantErr: false,
		},
		{
			name: "Not found",
			args: args{
				messageID: "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			queue, err := testDBInstance.GetQueue(ctx, tt.args.messageID)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			oldIdx := queue.CurrentPersonIdx

			err = testDBInstance.IncrementCurrentPerson(ctx, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrementCurrentPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			queue, err = testDBInstance.GetQueue(ctx, tt.args.messageID)
			assert.NoError(t, err)

			assert.Equal(t, oldIdx+1, queue.CurrentPersonIdx)
		})
	}

	queue, err := testDBInstance.GetQueue(context.Background(), "123")
	assert.NoError(t, err)

	assert.Equal(t, 1, queue.CurrentPersonIdx)
}

func TestCreateQueue(t *testing.T) {
	type args struct {
		messageID   string
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				messageID:   "create",
				description: "1234",
			},
			wantErr: false,
		},
		{
			name: "Duplicate",
			args: args{
				messageID:   "create",
				description: "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := testDBInstance.CreateQueue(ctx, tt.args.messageID, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			queue, err := testDBInstance.GetQueue(ctx, tt.args.messageID)
			assert.NoError(t, err)

			assert.Equal(t, tt.args.messageID, queue.MessageID)
			assert.Equal(t, tt.args.description, queue.Description)
			assert.Equal(t, 0, queue.CurrentPersonIdx)
			assert.Equal(t, 0, len(queue.Users))
		})
	}
}

func TestLogInOutToQueue(t *testing.T) {
	type args struct {
		messageID string
		user      entity.User
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantToExist bool
	}{
		{
			name: "Add user",
			args: args{
				messageID: "456",
				user: entity.User{
					ID:   1,
					Name: "Username",
				},
			},
			wantErr:     false,
			wantToExist: true,
		},
		{
			name: "Remove user",
			args: args{
				messageID: "123",
				user: entity.User{
					ID:   1,
					Name: "Username",
				},
			},
			wantErr:     false,
			wantToExist: false,
		},
		{
			name: "Not found",
			args: args{
				messageID: "1234",
				user: entity.User{
					ID:   1,
					Name: "Username",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := testDBInstance.LogInOutToQueue(ctx, tt.args.messageID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogInOutToQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			queue, err := testDBInstance.GetQueue(ctx, tt.args.messageID)
			assert.NoError(t, err)

			exists := slices.Contains(queue.Users, tt.args.user)
			assert.Equal(t, tt.wantToExist, exists)
		})
	}
}

func TestStartQueue(t *testing.T) {
	type args struct {
		messageID string
		shuffle   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK Shuffle",
			args: args{
				messageID: "789",
				shuffle:   true,
			},
			wantErr: false,
		},
		{
			name: "OK Not Shuffle",
			args: args{
				messageID: "789",
				shuffle:   false,
			},
			wantErr: false,
		},
		{
			name: "Not found",
			args: args{
				messageID: "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			queue, err := testDBInstance.GetQueue(ctx, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var users []entity.User
			if !tt.wantErr {
				users = queue.Users
			} else {
				users = []entity.User{}
			}

			wasShuffled := false

			for i := 0; i < 50; i++ {
				err := testDBInstance.StartQueue(ctx, tt.args.messageID, tt.args.shuffle)
				if (err != nil) != tt.wantErr {
					t.Errorf("StartQueue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err != nil {
					return
				}

				queue, err = testDBInstance.GetQueue(ctx, tt.args.messageID)
				assert.NoError(t, err)
				assert.Equal(t, 0, queue.CurrentPersonIdx)

				wasShuffled = !reflect.DeepEqual(users, queue.Users)
				if wasShuffled {
					break
				}
			}

			if wasShuffled != tt.args.shuffle {
				t.Errorf("StartQueue(), shuffle = %v, got = %v, want %v", tt.args.shuffle, queue.Users, users)
			}
		})
	}
}
