package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"money-diff/dao/models"
)

type PaymentDao interface {
	Create(p *models.Payment) error
	/*Delete(id primitive.ObjectID) error
	GetAll() ([]models.Payment, error)*/
	GetByChatID(chatID int64) ([]bson.M, error)
}
