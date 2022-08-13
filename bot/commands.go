package bot

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"

	command "money-diff/bot/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var commandList = map[string]interface{}{
	"help":  command.Help,
	"ap":    command.AddPayment,
	"total": command.GetTotal,
}

func cmdHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {

	if update.Message.IsCommand() {
		commandReq := update.Message.Command()
		commandArgs := update.Message.CommandArguments()
		botData := &helpers.BotUpdateData{
			Instance: bot,
			Update:   update,
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.UserName,
		}
		if _, ok := commandList[commandReq]; !ok {
			// Make a default action ?
			return nil
		}
		err := commandList[commandReq].(func(*mongo.Client, *helpers.BotUpdateData, string) error)(client, botData, commandArgs)

		if err != nil {
			fmt.Println(err)
		}

	}
	return nil
}

func commandHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	return cmdHandler(client, update, bot)
}
