package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
)

func Finish(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	msg := tgbotapi.NewMessage(bot.ChatID, "All participants must confirm!")

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, btnText := range []string{"Confirm", "Cancel"} {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, btnText)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	_, err := bot.Instance.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
