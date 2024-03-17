package entity

type Queue struct {
	MessageID        string
	Description      string
	Users            []User
	CurrentPersonIdx int
}
