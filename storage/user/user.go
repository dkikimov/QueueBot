package user

import "strings"

type User struct {
	Id   int64
	Name string
}

func New(id int64, lastName string, firstName string) User {
	return User{Id: id, Name: strings.TrimSpace(strings.Join([]string{lastName, firstName}, " "))}
}

// TODO: Make more effective

func UsersToString(users []User) (result string) {
	for _, user := range users {
		result += user.Name + "\n"
	}
	return result
}
