package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6446614126:AAHqSlZTpPNiTP3ZNjgigfcKvsjZkhNjiWA")
	if err != nil {
		log.Panic(err)
	}
	var updateConfig tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		bot.Send(msg)
	}

}
