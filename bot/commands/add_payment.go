package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/dao/impl"
	"money-diff/dao/models"
	"strconv"
	"strings"
)

func AddPayment(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	fmt.Println("args:" + arguments)
	value, err := strconv.ParseFloat(strings.Split(arguments, " ")[0], 32)
	if err != nil {
		return err
	}
	payment := &models.Payment{
		ID:       primitive.NewObjectID(),
		ChatID:   bot.ChatID,
		Username: bot.Username,
		Value:    float32(value),
	}

	paymentDao := impl.NewPaymentDao(client)
	err = paymentDao.Create(payment)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(bot.Update.Message.Chat.ID, "")
	msg.Text = "Payment added to the vault!"
	_, err = bot.Instance.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
