package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID       primitive.ObjectID `bson:"_id"`
	ChatID   int64              `bson:"chat_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Value    float64            `bson:"value"`
	Comment  string             `bson:"comment"`
}

type GroupedPayment struct {
	Username string  `bson:"_id,omitempty"`
	Total    float64 `bson:"value"`
}
