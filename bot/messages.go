package bot

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"money-diff/bot/helpers"
	"regexp"

	msg "money-diff/bot/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var regularMsgs = map[string]interface{}{
	"hi":   msg.Greet,
	"ALLO": msg.Greet,
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
		if ok, _ := regexp.MatchString(keyRegex, message); ok {
			responseQ = append(responseQ, valueRegex)
		}
	}

	// Take a decision about which actions will be performed.
	// For this example I will only executed the last action found.
	if resQLen := len(responseQ); resQLen > 0 {
		err := responseQ[resQLen-1].(func(*mongo.Client, *helpers.BotUpdateData, string) error)(client, botData, message)

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
