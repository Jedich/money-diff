package commands

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"money-diff/bot/helpers"
	"money-diff/db"
	"money-diff/db/models"
	"strconv"
	"strings"
	"time"
)

func AddPayment(conn *db.Connection, user *helpers.User, bot *tgbotapi.BotAPI, arguments string) error {
	fmt.Println("args:" + arguments)
	value, err := strconv.ParseFloat(strings.Split(arguments, " ")[0], 32)
	if err != nil {
		return err
	}

	collection := conn.Client.Database("money").Collection("payments")
	ctx, cancel := context.WithTimeout(conn.Ctx, 2*time.Second)
	defer cancel()

	req, err := collection.InsertOne(ctx, &models.Payment{
		ID:     primitive.NewObjectID(),
		ChatID: user.ChatID,
		UserID: user.UserID,
		Value:  float32(value),
	})
	fmt.Println(req.InsertedID)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	msg := tgbotapi.NewMessage(user.ChatID, "")
	msg.Text = "Payment added to the vault!"

	_, err = bot.Send(msg)

	return fmt.Errorf("error sending: %s", err)
}
