package main

import (
	"log"
	"money-diff/bot"
	"money-diff/bot/helpers"
	"money-diff/db"
)

// main function start the application.
func main() {
	if err := helpers.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	token := helpers.GetBotToken()
	conn := db.OpenConnection()

	log.Fatal(bot.StartBot(token, conn))
}
