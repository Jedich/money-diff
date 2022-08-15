package model

type Participant struct {
	UserID int   `bson:"user_id"`
	ChatID int64 `bson:"chat_id"`
}
