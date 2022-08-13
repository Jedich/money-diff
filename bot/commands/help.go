package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
)

// Help display the help for the bot
func Help(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {

	msg := tgbotapi.NewMessage(bot.ChatID, "")
	msg.Text = "this is the help"

	_, err := bot.Instance.Send(msg)

	return err
}
