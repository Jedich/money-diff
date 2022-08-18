package commands

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/db"
	"money-diff/model"
	r "money-diff/repository"
	"strings"
)

func Involve(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	var toInvolve string
	participantRepo := r.NewParticipantRepo(client)
	args := strings.Split(arguments, " ")
	if bot.ChatID == int64(bot.Update.Message.From.ID) {
		if args[0] == "" {
			return bot.SendMessageMarkdown("Please input a correct command.\n" +
				"`/involve[who]`")
		}
		toInvolve, args = args[0], args[1:]
	} else {
		if bot.Update.Message.ReplyToMessage == nil {
			return bot.SendMessage("Please reply to a target who you want to involve.")
		}
		target := bot.Update.Message.ReplyToMessage.From
		if target.UserName == bot.SenderName || target.IsBot {
			return bot.SendMessage("Please choose a valid target.")
		}
		toInvolve = target.UserName

		participant := &model.Participant{
			UserID: target.ID,
			ChatID: bot.ChatID,
		}
		err := participantRepo.Create(context.Background(), participant)
		if err != nil {
			return err
		}

	}
	err := db.WithTransaction(client, func(ctx mongo.SessionContext, client *mongo.Client) error {
		payment := &model.Payment{
			ID:       primitive.NewObjectID(),
			ChatID:   bot.ChatID,
			Username: toInvolve,
			Value:    0,
			Comment:  "",
		}

		paymentRepo := r.NewPaymentRepo(client)
		err := paymentRepo.Create(ctx, payment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	err = bot.SendMessage("Participant involved!")
	if err != nil {
		return err
	}
	return nil
}
