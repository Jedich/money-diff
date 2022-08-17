package commands

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/db"
	"money-diff/model"
	r "money-diff/repository"
	"strconv"
	"strings"
	"unicode/utf8"
)

func AddPayment(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	var valueString string
	var username string
	args := strings.Split(arguments, " ")
	if bot.ChatID == int64(bot.Update.Message.From.ID) {
		if len(args) < 2 {
			return bot.SendMessage("Please input a correct command.")
		}
		username, valueString, args = args[0], args[1], args[2:]
	} else {
		if len(args) < 1 {
			return bot.SendMessage("Please input a correct command.")
		}
		valueString, args = args[0], args[1:]
		username = bot.SenderName
	}

	value, err := strconv.ParseFloat(valueString, 32)
	if err != nil {
		return bot.SendMessage("Please input a correct float value.")
	}

	comment := strings.Join(args, " ")
	n := utf8.RuneCountInString(comment)
	if n > 50 {
		return bot.SendMessage("Please provide a shorter description. (%s > 50)", n)
	}
	err = db.WithTransaction(client, func(ctx mongo.SessionContext, client *mongo.Client) error {
		payment := &model.Payment{
			ID:       primitive.NewObjectID(),
			ChatID:   bot.ChatID,
			Username: username,
			Value:    float32(value),
			Comment:  comment,
		}

		paymentRepo := r.NewPaymentRepo(client)
		err = paymentRepo.Create(ctx, payment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	participant := &model.Participant{
		UserID: bot.Update.Message.From.ID,
		ChatID: bot.ChatID,
	}

	participantRepo := r.NewParticipantRepo(client)
	err = participantRepo.Create(context.Background(), participant)
	if err != nil {
		return err
	}

	return bot.SendMessage("Payment added to the vault!")
}
