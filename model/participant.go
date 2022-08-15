package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Participant struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID int                `bson:"user_id"`
	ChatID int64              `bson:"chat_id"`
}
