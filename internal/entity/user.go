package entity

import (
	"fmt"
	"strings"
)

type User struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
}

func New(id int64, lastName string, firstName string) User {
	return User{ID: id, Name: strings.TrimSpace(strings.Join([]string{lastName, firstName}, " "))}
}

func ListToString(users []User) (result string) {
	sb := strings.Builder{}
	for i, user := range users {
		sb.WriteString(user.Name)
		if i < len(users)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func ListToStringWithCurrent(users []User, currentUser int) string {
	sb := strings.Builder{}
	for idx, user := range users {
		if currentUser == idx {
			sb.WriteString(fmt.Sprintf("-> %s <-", user.Name))
		} else {
			sb.WriteString(user.Name)
		}

		if idx < len(users)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
