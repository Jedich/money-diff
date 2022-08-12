package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID     primitive.ObjectID `bson:"_id"`
	ChatID int64              `json:"chat_id,omitempty"`
	UserID int                `json:"user_id,omitempty"`
	Value  float32            `json:"value"`
}
