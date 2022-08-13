package helpers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotUpdateData struct {
	Instance *tgbotapi.BotAPI
	Update   tgbotapi.Update
	Username string
	ChatID   int64
}

func (bot *BotUpdateData) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(bot.ChatID, text)
	_, err := bot.Instance.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending: %s", err)
	}
	return nil
}
