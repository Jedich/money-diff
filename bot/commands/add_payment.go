package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"money-diff/bot/helpers"
	"money-diff/dao/db"
	"money-diff/dao/impl"
	"money-diff/dao/models"
	"strconv"
	"strings"
)

func AddPayment(conn *db.Connection, user *helpers.User, bot *tgbotapi.BotAPI, arguments string) error {
	fmt.Println("args:" + arguments)
	value, err := strconv.ParseFloat(strings.Split(arguments, " ")[0], 32)
	if err != nil {
		return err
	}
	payment := &models.Payment{
		ID:       primitive.NewObjectID(),
		ChatID:   user.ChatID,
		Username: user.Username,
		Value:    float32(value),
	}

	paymentDao := impl.NewPaymentDao(conn)
	err = paymentDao.Create(payment)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(user.ChatID, "")
	msg.Text = "Payment added to the vault!"
	_, err = bot.Send(msg)

	return fmt.Errorf("error sending: %s", err)
}
