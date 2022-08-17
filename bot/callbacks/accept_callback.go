package callbacks

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"money-diff/bot/helpers"
	"money-diff/repository"
)

func ProcessAcceptCallback(client *mongo.Client, bot *helpers.BotUpdateData, callback Callback, query string) error {
	repo := repository.NewConfirmRepo(client)
	if _, ok := repo.GetUsers(bot.ChatID)[bot.Update.CallbackQuery.From.ID]; !ok {
		return nil
	}
	issuerID := bot.Update.CallbackQuery.From.ID
	switch query {
	case "Confirm":
		err := repo.DeleteUser(bot.ChatID, issuerID)
		log.Println(repo.GetUsers(bot.ChatID))
		if err != nil {
			return err
		}
		if len(repo.GetUsers(bot.ChatID)) == 0 {
			callback.Finish()
		}
		return nil
	case "Cancel":
		callback.Cancel()
	default:
		return errors.New("incorrect callback processor")
	}
	return nil
}
