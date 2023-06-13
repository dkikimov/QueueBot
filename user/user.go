package user

import (
	"fmt"
	"strings"
)

type User struct {
	Id   int64
	Name string
}

func New(id int64, lastName string, firstName string) User {
	return User{Id: id, Name: strings.TrimSpace(strings.Join([]string{lastName, firstName}, " "))}
}

// TODO: Make more effective

func ListToString(users []User) (result string) {
	for _, user := range users {
		result += user.Name + "\n"
	}
	return result
}

func ListToStringWithCurrent(users []User, currentUser int) (result string) {
	for idx, user := range users {
		if currentUser == idx {
			result += fmt.Sprintf("-> %s <-\n", user.Name)
		} else {
			result += user.Name + "\n"
		}

	}
	return result
}
