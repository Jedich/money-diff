package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DirectPayment struct {
	ID           primitive.ObjectID `bson:"_id"`
	ChatID       int64              `bson:"chat_id,omitempty"`
	FromUsername string             `bson:"from,omitempty"`
	ToUsername   string             `bson:"to,omitempty"`
	Value        float32            `bson:"value"`
	Comment      string             `bson:"comment"`
}
