package entity

type Queue struct {
	MessageID        string `bson:"messageId"`
	Description      string `bson:"description"`
	Users            []User `bson:"users"`
	CurrentPersonIdx int    `bson:"currentUserIndex"`
}
