package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	msg := tgbotapi.NewMessage(bot.ChatID, "Total:\n")
	resString := ""
	for _, payment := range payments {
		resString += fmt.Sprintf("%v: %.2f\n", payment.Username, payment.Total)
	}
	for _, payment := range directPayments {
		resString += fmt.Sprintf("%v\n", payment)
	}
	if resString == "" {
		resString = "It's empty.."
	}
	msg.Text += resString
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
