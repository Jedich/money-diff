package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	callback "money-diff/bot/callbacks"
	"money-diff/bot/helpers"
)

var callbackList = map[string]interface{}{
	"Confirm": callback.ProcessAcceptCallback,
	"Cancel":  callback.ProcessAcceptCallback,
}

func callbackHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	query := update.CallbackQuery.Data
	botData := &helpers.BotUpdateData{
		BotAPI:     bot,
		Update:     update,
		SenderName: update.CallbackQuery.From.UserName,
		ChatID:     update.CallbackQuery.Message.Chat.ID,
	}
	cb, _ := callback.GetCallback(botData.ChatID)

	if _, ok := callbackList[query]; !ok {
		// Make a default action ?
		return nil
	}
	err := callbackList[query].(func(*mongo.Client, *helpers.BotUpdateData, callback.Callback, string) error)(client, botData, cb, query)

	if err != nil {
		log.Println(err)
		err := botData.SendMessage("Toaster is broken")
		if err != nil {
			return err
		}
	}
	return nil
}
