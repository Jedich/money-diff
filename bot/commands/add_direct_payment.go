package commands

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/model"
	r "money-diff/repository"
	"strconv"
	"strings"
	"unicode/utf8"
)

func AddDirectPayment(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	if bot.Update.Message.ReplyToMessage == nil {
		return bot.SendMessage("Please reply to a target who you paid for.")
	}
	target := bot.Update.Message.ReplyToMessage.From
	if target.UserName == bot.SenderName || target.IsBot {
		return bot.SendMessage("Please choose a valid target.")
	}

	args := strings.Split(arguments, " ")
	x, args := args[0], args[1:]

	comment := strings.Join(args, " ")
	n := utf8.RuneCountInString(comment)
	if n > 50 {
		return bot.SendMessage("Please provide a shorter description. (%s > 50)", n)
	}

	value, err := strconv.ParseFloat(x, 32)
	if err != nil {
		return bot.SendMessage("Please input a correct float value.")
	}

	payment := &model.DirectPayment{
		ID:           primitive.NewObjectID(),
		ChatID:       bot.ChatID,
		FromUsername: bot.SenderName,
		ToUsername:   target.UserName,
		Value:        float32(value),
		Comment:      strings.Join(args, " "),
	}

	paymentRepo := r.NewDirectPaymentRepo(client)
	err = paymentRepo.Create(payment)
	if err != nil {
		return err
	}
	return bot.SendMessage("Payment added to the vault!")
}
