package messages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
)

// Greet this function send a random
// greet to the bot
func Greet(client *mongo.Client, bot *helpers.BotUpdateData, message string) error {

	msg := tgbotapi.NewMessage(bot.ChatID, "")
	msg.Text = "Hiiiii!"

	_, err := bot.Send(msg)

	return err
}
