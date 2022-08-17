package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	callback "money-diff/bot/callbacks"
	"money-diff/bot/helpers"
	"money-diff/repository"
	"time"
)

func Finish(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	if _, ok := callback.GetCallback(bot.ChatID); !ok {
		msg := tgbotapi.NewMessage(bot.ChatID, "All participants must confirm!")

		keyboard := tgbotapi.InlineKeyboardMarkup{}
		for _, btnText := range []string{"Confirm", "Cancel"} {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(btnText, btnText)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			return err
		}

		cbSuccess := func() error {
			return bot.SendMessage("Request accepted! Calculating...")
		}

		cbFailure := func() error {
			return bot.SendMessage("Request was cancelled.")
		}

		repo := repository.NewConfirmRepo(client)
		err = repo.Add(bot.ChatID)
		if err != nil {
			return err
		}

		cb := callback.NewCallback(bot.ChatID, cbSuccess, cbFailure)

		err = cb.Start(5 * time.Second)
		repo.Finish(bot.ChatID)
		if err != nil {
			return err
		}

		_, err = bot.DeleteMessage(tgbotapi.DeleteMessageConfig{MessageID: sent.MessageID, ChatID: bot.ChatID})
		if err != nil {
			return err
		}
	}
	return nil
}
