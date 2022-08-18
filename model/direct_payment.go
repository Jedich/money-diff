package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirectPayment struct {
	ID      primitive.ObjectID `bson:"_id"`
	ChatID  int64              `bson:"chat_id,omitempty"`
	WhoOwes string             `bson:"from,omitempty"`
	Whom    string             `bson:"to,omitempty"`
	Value   float64            `bson:"value"`
	Comment string             `bson:"comment"`
}

type GroupedDirectPayment struct {
	User  DirectUserDTO `bson:"_id,omitempty"`
	Total float64       `bson:"value"`
}

type DirectUserDTO struct {
	WhoOwes string `bson:"from,omitempty"`
	Whom    string `bson:"to,omitempty"`
}

func (d DirectPayment) String() string {
	return fmt.Sprintf("%s винен %s: %.2f <i>%s</i>", d.WhoOwes, d.Whom, d.Value, d.Comment)
}

func (d GroupedDirectPayment) String() string {
	return fmt.Sprintf("%s винен %s: %.2f", d.User.WhoOwes, d.User.Whom, d.Total)
}
