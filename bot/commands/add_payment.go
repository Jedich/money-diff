package commands

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/dao/impl"
	"money-diff/dao/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

func AddPayment(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	args := strings.Split(arguments, " ")
	x, args := args[0], args[1:]

	value, err := strconv.ParseFloat(x, 32)
	if err != nil {
		return bot.SendMessage("Please input a correct float value.")
	}

	comment := strings.Join(args, " ")
	n := utf8.RuneCountInString(comment)
	if n > 50 {
		return bot.SendMessage("Please provide a shorter description. (%s > 50)", n)
	}

	payment := &models.Payment{
		ID:       primitive.NewObjectID(),
		ChatID:   bot.ChatID,
		Username: bot.SenderName,
		Value:    float32(value),
		Comment:  comment,
	}

	paymentDao := impl.NewPaymentDao(client)
	err = paymentDao.Create(payment)
	if err != nil {
		return err
	}

	return bot.SendMessage("Payment added to the vault!")
}
