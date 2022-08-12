package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"money-diff/bot/helpers"
	"money-diff/dao/db"
	"money-diff/dao/impl"
)

func GetTotal(conn *db.Connection, user *helpers.User, bot *tgbotapi.BotAPI, arguments string) error {
	paymentDao := impl.NewPaymentDao(conn)
	payments, err := paymentDao.GetByChatID(user.ChatID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(user.ChatID, "")
	for _, payment := range payments {
		msg.Text += fmt.Sprintf("%v: %.2f\n", payment["_id"], payment["value"])
	}
	_, err = bot.Send(msg)
	return nil
}
