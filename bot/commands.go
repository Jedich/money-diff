package bot

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"money-diff/bot/helpers"

	command "money-diff/bot/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var commandList = map[string]interface{}{
	"help":    command.Help,
	"ap":      command.AddPayment,
	"adp":     command.AddDirectPayment,
	"total":   command.GetTotal,
	"history": command.GetHistory,
	"finish":  command.Finish,
}

func cmdHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	commandReq := update.Message.Command()
	commandArgs := update.Message.CommandArguments()
	botData := &helpers.BotUpdateData{
		BotAPI:     bot,
		Update:     update,
		SenderName: update.Message.From.UserName,
		ChatID:     update.Message.Chat.ID,
	}

	if _, ok := commandList[commandReq]; !ok {
		// Make a default action ?
		return nil
	}
	err := commandList[commandReq].(func(*mongo.Client, *helpers.BotUpdateData, string) error)(client, botData, commandArgs)

	if err != nil {
		log.Println(err)
		err := botData.SendMessage("Toaster is broken")
		if err != nil {
			return err
		}
	}
	return nil
}

func commandHandler(client *mongo.Client, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	return cmdHandler(client, update, bot)
}
