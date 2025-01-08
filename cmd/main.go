package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	err := tgbotapi.SetLogger(logger)
	if err != nil {
		logger.Fatal("Failed to set Telegram bot logger; error: ", err)
	}

	bot, err := tgbotapi.NewBotAPI(TgToken)
	if err != nil {
		logger.Fatal("Error creating bot", err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = TimeoutForUpdate

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err = bot.Send(msg); err != nil {
			logger.Fatal("Error sending message", err)
		}
	}
}
