package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	r "money-diff/repository"
)

func GetHistory(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	paymentRepo := r.NewPaymentRepo(client)
	payments, err := paymentRepo.GetByChatID(bot.ChatID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(bot.ChatID, "")
	for _, payment := range payments {
		msg.Text += fmt.Sprintf("%s: %.2f <i>%s</i>\n", payment.Username, payment.Value, payment.Comment)
	}
	directPaymentRepo := r.NewDirectPaymentRepo(client)
	dPayments, err := directPaymentRepo.GetByChatID(bot.ChatID)
	if err != nil {
		return err
	}
	for _, payment := range dPayments {
		msg.Text += fmt.Sprintf("%v\n", payment)
	}
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
