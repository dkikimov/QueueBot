package entity

type Queue struct {
	MessageID        string `bson:"message_id"`
	Description      string `bson:"description"`
	Users            []User `bson:"users"`
	CurrentPersonIdx int    `bson:"current_user_index"`
}
