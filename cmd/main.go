package main

import (
	"log"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"poison_bot/internal/coreView"
	"poison_bot/internal/sender"
)

//var numericKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("1"),
//		tgbotapi.NewKeyboardButton("2"),
//		tgbotapi.NewKeyboardButton("3"),
//	),
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("4"),
//		tgbotapi.NewKeyboardButton("5"),
//		tgbotapi.NewKeyboardButton("6"),
//	),
//)

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

	s := sender.NewSender(logger, bot)
	waitGroup := sync.WaitGroup{}
	worker := coreView.New(logger, updates, s, &waitGroup)

	waitGroup.Add(Workers)
	err = worker.Process()
	waitGroup.Wait()

	//for update := range updates {
	//	if update.Message == nil { // ignore any non-Message updates
	//		continue
	//	}
	//
	//	//if !update.Message.IsCommand() { // ignore any non-command Messages
	//	//	continue
	//	//}
	//
	//	// Create a new MessageConfig. We don't have text yet,
	//	// so we leave it empty.
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	//chanID := int64(-1002430802369)
	//	msgForChannel := tgbotapi.NewMessage(ChannelForOrdersID, msg.Text)
	//
	//	switch update.Message.Text {
	//	case "open":
	//		msg.ReplyMarkup = numericKeyboard
	//	case "close":
	//		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	//	}
	//
	//	data, err := bot.Send(msgForChannel)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	logger.Println(data)
	//}
}
