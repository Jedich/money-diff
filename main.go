package main

import (
	"context"
	"fmt"
	"log"
	"money-diff/bot"
	"money-diff/bot/helpers"
	"money-diff/dao/db"
)

// main function start the application.
func main() {
	if err := helpers.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	token := helpers.GetBotToken()
	client := db.OpenConnection(ctx)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("disconnected from dao")
	}()

	log.Fatal(bot.StartBot(token, &db.Connection{
		Client: client,
		Ctx:    ctx,
	}))

}
