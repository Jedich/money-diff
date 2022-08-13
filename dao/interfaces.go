package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"money-diff/dao/models"
)

type PaymentDao interface {
	Create(p *models.Payment) error
	/*Delete(id primitive.ObjectID) error
	GetAll() ([]models.Payment, error)*/
	GetGroupByChatID(chatID int64) ([]bson.M, error)
	GetByChatID(chatID int64) ([]models.Payment, error)
}

type DirectPaymentDao interface {
	Create(p *models.DirectPayment) error
	/*Delete(id primitive.ObjectID) error
	GetAll() ([]models.Payment, error)*/
	GetGroupByChatID(chatID int64) ([]bson.M, error)
	GetByChatID(chatID int64) ([]models.DirectPayment, error)
}
