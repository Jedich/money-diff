package commands

import (
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/bot/helpers"
)

// Help display the help for the bot
func Help(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {

	err := bot.SendMessage("Hello! My name is Gustavo, but you can call me Gus. \n" +
		"I can split your shared spendings! List of commands:\n" +
		"/help - Show help\n" +
		"/ap - Add payment for yourself. In private: `/ap[name][value]{desc}`\n" +
		"/adp - Add who you paid for. Reply to target. In private: `/adp[who][to][value]{desc}`\n" +
		"/involve - Add participant to split share. In private `/involve[who]`\n" +
		"/history - show all additions to payments in chat\n" +
		"/total - show summarized total of each member\n" +
		"/finish - Finish and get results. All participants must confirm!")

	return err
}
