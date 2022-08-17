package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DirectPayment struct {
	ID           primitive.ObjectID `bson:"_id"`
	ChatID       int64              `bson:"chat_id,omitempty"`
	FromUsername string             `bson:"from,omitempty"`
	ToUsername   string             `bson:"to,omitempty"`
	Value        float64            `bson:"value"`
	Comment      string             `bson:"comment"`
}

type GroupedDirectPayment struct {
	User  DirectUserDTO `bson:"_id,omitempty"`
	Total float64       `bson:"value"`
}

type DirectUserDTO struct {
	From string `bson:"from,omitempty"`
	To   string `bson:"to,omitempty"`
}
