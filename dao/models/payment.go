package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID       primitive.ObjectID `bson:"_id"`
	ChatID   int64              `bson:"chat_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Value    float32            `bson:"value"`
	Comment  string             `bson:"comment"`
}

type GroupedPayment struct {
	Username string  `bson:"username,omitempty"`
	Total    float32 `bson:"value"`
}
