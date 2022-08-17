package bot

import (
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
		errs := make(chan error, 1)
		go func(u tgbotapi.Update) {
			defer close(errs)
			if u.CallbackQuery != nil {
				if err = callbackHandler(client, u, bot); err != nil { // Handle the command message
					errs <- err
				}
				return
			}

			if u.Message == nil {
				return
			}

			log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)

			if u.Message.IsCommand() {
				if err = commandHandler(client, u, bot); err != nil { // Handle the command message
					errs <- err
				}
				return
			}

			// Do something else with the message received
			if err = messageHandler(client, u, bot); err != nil { // Handle the single message
				errs <- err
			}
		}(update)
	}

	return nil
}
