package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	r "money-diff/repository"
)

func GetTotal(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	paymentRepo := r.NewPaymentRepo(client)
	directPaymentRepo := r.NewDirectPaymentRepo(client)
	payments, err := paymentRepo.GetGroupByChatID(bot.ChatID)
	if err != nil {
		return err
	}
	directPayments, err := directPaymentRepo.GetGroupByChatID(bot.ChatID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(bot.ChatID, "")
	for _, payment := range payments {
		msg.Text += fmt.Sprintf("%v: %.2f\n", payment["_id"], payment["value"])
	}
	for _, payment := range directPayments {
		direct := payment["_id"].(bson.M)
		msg.Text += fmt.Sprintf("%v -> %v: %.2f\n", direct["from"], direct["to"], payment["value"])
	}
	_, err = bot.Instance.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
