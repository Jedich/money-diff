package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BotUpdateData struct {
	Instance *tgbotapi.BotAPI
	Update   tgbotapi.Update
	Username string
	ChatID   int64
}
