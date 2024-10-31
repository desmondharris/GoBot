package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var appKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("portal", "http://127.0.0.1:9090/portal"),
	),
)

func main() {
	// i already removed it, don't even try :)
	bot, err := tgbotapi.NewBotAPI("6446614126:AAHqSlZTpPNiTP3ZNjgigfcKvsjZkhNjiWA")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			// confighelper, will echo message
			var msg tgbotapi.MessageConfig

			switch update.Message.Text {
			case "open":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Opening bot portal...")
				msg.ReplyMarkup = appKeyboard
				if _, err = bot.Send(msg); err != nil {
					log.Println(err)
				}
			}

		} else if update.CallbackQuery != nil {
			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
