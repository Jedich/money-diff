package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	callback "money-diff/bot/callbacks"
	"money-diff/bot/helpers"
	"money-diff/repository"
	"time"
)

func Finish(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	if _, ok := callback.GetCallback(bot.ChatID); !ok {
		msg := tgbotapi.NewMessage(bot.ChatID, "All participants must confirm!")

		keyboard := tgbotapi.InlineKeyboardMarkup{}
		for _, btnText := range []string{"Confirm", "Cancel"} {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(btnText, btnText)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			return err
		}

		cbSuccess := func() error {
			err := bot.SendMessage("Request accepted! Calculating...")
			if err != nil {
				return err
			}
			err = GetTotal(client, bot, "")
			if err != nil {
				return err
			}
			res, err := Calculate(client, bot)
			resString := "Result:\n"
			for _, r := range res {
				resString += fmt.Sprintf("%v\n", r)
			}
			err = bot.SendMessage(resString)
			if err != nil {
				return err
			}
			return nil
		}

		cbFailure := func() error {
			return bot.SendMessage("Request was cancelled.")
		}

		repo := repository.NewConfirmRepo(client)
		err = repo.Add(bot.ChatID)
		if err != nil {
			return err
		}

		cb := callback.NewCallback(bot.ChatID, cbSuccess, cbFailure)

		err = cb.Start(5 * time.Second)
		repo.Finish(bot.ChatID)
		if err != nil {
			return err
		}

		_, err = bot.DeleteMessage(tgbotapi.DeleteMessageConfig{MessageID: sent.MessageID, ChatID: bot.ChatID})
		if err != nil {
			return err
		}
	}
	return nil
}

type Debt struct {
	from  string
	to    string
	Value float64
}

func (d Debt) String() string {
	return fmt.Sprintf("%s -> %s: %.2f", d.from, d.to, d.Value)
}

func Calculate(client *mongo.Client, bot *helpers.BotUpdateData) (map[string]Debt, error) {
	paymentRepo := repository.NewPaymentRepo(client)
	directPaymentRepo := repository.NewDirectPaymentRepo(client)
	payments, err := paymentRepo.GetGroupByChatID(bot.ChatID)
	if err != nil {
		return nil, err
	}

	var totalSum float64
	for _, payment := range payments {
		totalSum += payment.Total
	}
	part := totalSum / float64(len(payments))
	negativeValues := make(map[string]float64)
	positiveValues := make(map[string]float64)

	for _, payment := range payments {
		res := payment.Total - part
		if res >= 0 {
			positiveValues[payment.Username] = res
			continue
		}
		negativeValues[payment.Username] = res
	}
	results := make(map[string]Debt, len(payments))
	for key1, val1 := range positiveValues {
		v1 := val1
		for key2, val2 := range negativeValues {
			if v1 > 0 {
				if v1+val2 < 0 {
					results[key1+key2] = Debt{key1, key2, -v1}
					negativeValues[key2] = val2 + v1
				} else {
					results[key1+key2] = Debt{key1, key2, v1}
					negativeValues[key2] = 0
				}
			} else {
				break
			}
		}
	}

	directPayments, err := directPaymentRepo.GetGroupByChatID(bot.ChatID)
	if err != nil {
		return nil, err
	}
	for _, payment := range directPayments {
		name := payment.User.From + payment.User.To
		invName := payment.User.To + payment.User.From
		if t, ok := results[invName]; ok {
			if payment.Total == t.Value {
				delete(results, invName)
			} else if payment.Total > t.Value {
				payment.Total -= t.Value
				delete(results, invName)
			} else {
				t.Value -= payment.Total
				results[invName] = t
			}
		}
		if t, ok := results[name]; ok {
			t.Value += payment.Total
			results[name] = t
		} else {
			results[name] = Debt{payment.User.From, payment.User.To, payment.Total}
		}
	}
	return results, nil
}
