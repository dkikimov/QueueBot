package entity

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

func ListToString(users []User) (result string) {
	sb := strings.Builder{}
	for _, user := range users {
		sb.WriteString(user.Name + "\n")
	}
	return sb.String()
}

func ListToStringWithCurrent(users []User, currentUser int) string {
	sb := strings.Builder{}
	for idx, user := range users {
		if currentUser == idx {
			sb.WriteString(fmt.Sprintf("-> %s <-\n", user.Name))
		} else {
			sb.WriteString(user.Name + "\n")
		}
	}

	return sb.String()
}
