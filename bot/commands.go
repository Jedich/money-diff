package bot

import (
	"fmt"
	"money-diff/bot/helpers"
	"money-diff/dao/db"

	command "money-diff/bot/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var commandList = map[string]interface{}{
	"help":  command.Help,
	"ap":    command.AddPayment,
	"total": command.GetTotal,
}

func cmdHandler(conn *db.Connection, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {

	if update.Message.IsCommand() {
		commandReq := update.Message.Command()
		commandArgs := update.Message.CommandArguments()
		user := &helpers.User{
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
		}
		if _, ok := commandList[commandReq]; !ok {
			// Make a default action ?
			return nil
		}
		err := commandList[commandReq].(func(*db.Connection, *helpers.User, *tgbotapi.BotAPI, string) error)(conn, user, bot, commandArgs)

		if err != nil {
			fmt.Println(err)
		}

	}
	return nil
}

func commandHandler(conn *db.Connection, update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	return cmdHandler(conn, update, bot)
}
