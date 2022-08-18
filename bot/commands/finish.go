package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	callback "money-diff/bot/callbacks"
	"money-diff/bot/helpers"
	"money-diff/db"
	"money-diff/model"
	"money-diff/repository"
	"time"
)

func Finish(client *mongo.Client, bot *helpers.BotUpdateData, arguments string) error {
	if _, ok := callback.GetCallback(bot.ChatID); !ok {
		res, err := repository.NewParticipantRepo(client).GetByChatID(bot.ChatID)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			err := bot.SendMessage("Nothing to confirm yet.")
			if err != nil {
				return err
			}
			return nil
		}

		msg := tgbotapi.NewMessage(bot.ChatID, "")
		msg.Text = fmt.Sprintf("All participants must confirm! You have %d participant(s) and *5 minutes* to do this.", len(res))
		msg.ParseMode = tgbotapi.ModeMarkdown
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
			paymentRepo := repository.NewPaymentRepo(client)
			directPaymentRepo := repository.NewDirectPaymentRepo(client)
			payments, err := paymentRepo.GetGroupByChatID(bot.ChatID)
			if err != nil {
				return nil
			}
			directPayments, err := directPaymentRepo.GetGroupByChatID(bot.ChatID)
			if err != nil {
				return nil
			}
			res, err := Calculate(payments, directPayments)
			resString := ""
			for _, r := range res {
				resString += fmt.Sprintf("%v\n", r)
			}
			if resString == "" {
				resString = "All clear!"
			}
			err = bot.SendMessage("Result:\n" + resString)
			if err != nil {
				return err
			}
			err = db.WithTransaction(client, func(ctx mongo.SessionContext, client *mongo.Client) error {
				err := repository.NewPaymentRepo(client).DeleteByChatID(ctx, bot.ChatID)
				if err != nil {
					return err
				}
				err = repository.NewDirectPaymentRepo(client).DeleteByChatID(ctx, bot.ChatID)
				if err != nil {
					return err
				}
				err = repository.NewParticipantRepo(client).DeleteByChatID(ctx, bot.ChatID)
				if err != nil {
					return err
				}
				return nil
			})
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

		err = cb.Start(5 * time.Minute)
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

type debt struct {
	who   string
	whom  string
	value float64
}

func (d debt) String() string {
	return fmt.Sprintf("%s -> %s: %.2f", d.who, d.whom, d.value)
}

func Calculate(payments []model.GroupedPayment, directPayments []model.GroupedDirectPayment) (map[string]debt, error) {
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
	results := make(map[string]debt, len(payments))
	for key1, val1 := range positiveValues {
		v1 := val1
		for key2, val2 := range negativeValues {
			if v1 > 0 {
				if v1+val2 < 0 {
					results[key2+key1] = debt{key2, key1, v1}
					negativeValues[key2] = val2 + v1
				} else {
					results[key2+key1] = debt{key2, key1, -val2}
					negativeValues[key2] = 0
				}
			} else {
				break
			}
		}
	}

	for _, currentDebt := range directPayments {
		name := currentDebt.User.WhoOwes + currentDebt.User.Whom
		invName := currentDebt.User.Whom + currentDebt.User.WhoOwes
		if existingDebtInv, ok := results[invName]; ok {
			if currentDebt.Total == existingDebtInv.value {
				delete(results, invName)
			} else if currentDebt.Total > existingDebtInv.value {
				currentDebt.Total -= existingDebtInv.value
				delete(results, invName)
			} else {
				existingDebtInv.value -= currentDebt.Total
				results[invName] = existingDebtInv
				continue
			}
		}
		if t, ok := results[name]; ok {
			t.value += currentDebt.Total
			results[name] = t
		} else {
			results[name] = debt{currentDebt.User.WhoOwes, currentDebt.User.Whom, currentDebt.Total}
		}
	}
	return results, nil
}
