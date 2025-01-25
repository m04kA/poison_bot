package main

import (
	"log"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	coreview "poison_bot/internal/core_view"
	orderrepo "poison_bot/internal/db/orders/repository"
	pricecalculator "poison_bot/internal/price_calculator"
	"poison_bot/internal/sender"
	createitem "poison_bot/internal/usecase/create_item"
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

const (
	success = 0
	fail    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	var err error

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Printf("panic error: %v", panicErr)
			exitCode = fail
		}
	}()

	err = tgbotapi.SetLogger(logger)
	if err != nil {
		logger.Fatal("Failed to set Telegram bot logger; error: ", err)
		return fail
	}

	bot, err := tgbotapi.NewBotAPI(TgToken)
	if err != nil {
		logger.Fatal("Error creating bot", err)
		return fail
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = TimeoutForUpdate

	updates := bot.GetUpdatesChan(updateConfig)

	s := sender.NewSender(logger, bot)
	or := orderrepo.NewOrderRepository()
	ip := createitem.NewProcessor(or, s)
	pc := pricecalculator.New(CNYToRUB)
	waitGroup := sync.WaitGroup{}
	worker := coreview.New(logger, s, or, pc, updates, &waitGroup, ip, ChannelForOrdersReports)

	waitGroup.Add(Workers) // TODO: Сделать нормальные воркеры, чтоб это работало по назначению
	err = worker.Process()
	waitGroup.Wait()

	if err != nil {
		logger.Fatal("Error processing updates", err)
		return fail
	}
	return success
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
