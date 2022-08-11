package main

import (
	"log"

	bot "money-diff/bot"
	helper "money-diff/bot/helpers"
)

// main function start the application.
func main() {

	if err := helper.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	token := helper.GetBotToken()

	log.Fatal(bot.StartBot(token))
}
