package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// StartBot this function to start the bot
func StartBot(token string, client *mongo.Client) error {

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			//_ := update.CallbackQuery.Data
			//classesMap[update.CallbackQuery.From.ID] = class
			_, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
				"Ok, I remember"))
			fmt.Println(update.CallbackQuery.Data)
			if err != nil {
				return err
			}
		}

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			if err = commandHandler(client, update, bot); err != nil { // Handle the command message
				return err
			}
			continue
		}

		// Do something else with the message received
		if err = messageHandler(client, update, bot); err != nil { // Handle the single message
			return err
		}

	}

	return nil
}
