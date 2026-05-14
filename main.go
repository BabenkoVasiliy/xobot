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

	log.Printf("Authorized on account %s", bot.Self.UserName)

	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Найти соперника"),
			tgbotapi.NewKeyboardButton("Играть с ботом"),
		),
	)

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери режим игры:")
			msg.ReplyMarkup = kb

			bot.Send(msg)
			log.Printf("Sent keyboard to chat %d", update.Message.Chat.ID)
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			bot.Request(callback)

			log.Printf("Callback: %s", update.CallbackQuery.Data)
		}
	}
}