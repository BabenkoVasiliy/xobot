package main

import (
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var gameState = make(map[int64]bool)

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
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			if text == "Играть с ботом" {
				gameState[chatID] = true

				btn := tgbotapi.NewInlineKeyboardButtonWebApp("Играть", "https://BabenkoVasiliy.github.io/")
				kbInline := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

				msg := tgbotapi.NewMessage(chatID, "Нажми кнопку ниже чтобы начать игру:")
				msg.ReplyMarkup = kbInline
				bot.Send(msg)
			} else if text == "Найти соперника" {
				msg := tgbotapi.NewMessage(chatID, "Поиск соперника временно недоступен")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Выбери режим игры:")
				msg.ReplyMarkup = kb
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			bot.Request(callback)

			log.Printf("Callback from %d: %s", update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		}
	}

	_ = strconv.Mkinter
}