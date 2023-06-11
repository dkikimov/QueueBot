package storage

type Queue struct {
	description string
}

func NewQueue(description string) *Queue {
	return &Queue{description: description}
}

type User struct {
	id   int64
	name string
}

type Storage interface {
	GetUsersInQueue(queueId string) ([]User, error)
	AddUserToQueue(queueId string, user User) error
	DeleteUserFromQueueById(queueId string, userId int64) error
	CreateQueue(queryID string, description string) error
}
