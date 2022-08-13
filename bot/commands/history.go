package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
	"money-diff/dao/impl"
)

func GetHistory(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	paymentDao := impl.NewPaymentDao(client)
	payments, err := paymentDao.GetByChatID(bot.ChatID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(bot.ChatID, "")
	for _, payment := range payments {
		msg.Text += fmt.Sprintf("%s: %.2f <i>%s</i>\n", payment.Username, payment.Value, payment.Comment)
		if payment.Comment != "" {

		}
	}
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = bot.Instance.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
