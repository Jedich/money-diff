package helpers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotUpdateData struct {
	*tgbotapi.BotAPI
	Update     tgbotapi.Update
	SenderName string
	ChatID     int64
}

func (bot *BotUpdateData) SendMessage(text string, a ...any) error {
	msg := tgbotapi.NewMessage(bot.ChatID, "")
	var err error
	if len(a) != 0 {
		err = bot.sendMessage(msg, text, a)
	} else {
		err = bot.sendMessage(msg, text)
	}
	return err
}

func (bot *BotUpdateData) SendMessageMarkdown(text string, a ...any) error {
	msg := tgbotapi.NewMessage(bot.ChatID, "")
	msg.ParseMode = tgbotapi.ModeMarkdown
	var err error
	if len(a) != 0 {
		err = bot.sendMessage(msg, text, a)
	} else {
		err = bot.sendMessage(msg, text)
	}
	return err
}

func (bot *BotUpdateData) sendMessage(msg tgbotapi.MessageConfig, text string, a ...any) error {
	msg.Text = ""
	if len(a) != 0 {
		msg.Text = fmt.Sprintf(text, a)
	} else {
		msg.Text = fmt.Sprintf(text)
	}
	if msg.Text != "" {
		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("error sending: %s", err)
		}
	}
	return nil
}
