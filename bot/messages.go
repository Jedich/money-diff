package bot

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	command "money-diff/bot/commands"
	"money-diff/bot/helpers"
	"regexp"
	"strings"

	msg "money-diff/bot/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var regularMsgs = map[string]interface{}{
	"hi":        msg.Greet,
	"бот додай": command.AddPayment,
	"bot add":   command.AddPayment,
}

func messageHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	if update.Message.IsCommand() {
		log.Println("no")
	}
	message := update.Message.Text
	botData := &helpers.BotUpdateData{
		Instance: bot,
		Update:   update,
		ChatID:   update.Message.Chat.ID,
		Username: update.Message.From.UserName,
	}

	var responseQ []interface{} // contain the message received

	// Store all the possible response that can be performed based in the regex
	for keyRegex, valueRegex := range regularMsgs {
		if ok, _ := regexp.MatchString(keyRegex, strings.ToLower(message)); ok {
			responseQ = append(responseQ, valueRegex)
		}
	}
	args := strings.Split(message, " ")

	// last action found will be performed.
	if resQLen := len(responseQ); resQLen > 0 {
		// Set last word as argument
		err := responseQ[resQLen-1].(func(*mongo.Client, *helpers.BotUpdateData, string) error)(client, botData, args[len(args)-1])

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
