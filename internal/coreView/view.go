package coreView

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type View struct {
	log     *log.Logger
	sender  Sender
	updates tgbotapi.UpdatesChannel
	wg      *sync.WaitGroup
}

func New(l *log.Logger, updates tgbotapi.UpdatesChannel, sender Sender, wg *sync.WaitGroup) *View {
	return &View{
		log:     l,
		sender:  sender,
		updates: updates,
		wg:      wg,
	}
}

func (v *View) Process() error {
	defer v.wg.Done()

	for update := range v.updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && (update.Message.Command() == "start" || update.Message.Command() == "help") {
			err := v.sender.SendStartMessage(update.Message.Chat.ID)
			if err != nil {
				return err
			}
		}

	}
	return nil

}
