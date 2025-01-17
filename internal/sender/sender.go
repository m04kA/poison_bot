package sender

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	log *log.Logger
	bot *tgbotapi.BotAPI
}

func NewSender(l *log.Logger, bot *tgbotapi.BotAPI) *Sender {
	return &Sender{
		log: l,
		bot: bot,
	}
}

func (s *Sender) SendStartMessage(chatId int64) error {
	text := `
	Привет! Это бот по заказу одежды и обуви с пойзона.
	
	Чтоб сделать новый заказ, выполни команду /createOrder.
	После создания заказа, с вами свяжется наш менеджер.
	`
	msg := tgbotapi.NewMessage(chatId, text)

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}
