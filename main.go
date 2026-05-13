package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := "8669066608:AAHQeaPBVCT_khKKTWSEZsDXWe7pAWjuoMo"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
			log.Printf("Echo: %s", update.Message.Text)
		}
	}
}