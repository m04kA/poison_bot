package sender

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	basket "poison_bot/internal/db/basket/entity"
	orders "poison_bot/internal/domain"
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

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/help"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/create_order"),
	),
)

var editKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/help"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/send_order_to_manage"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/remove_item_data"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/cancel_order"),
	),
)

var addNewItemKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/help"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/send_order_to_manage"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/add_new_item_to_order"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/cancel_order"),
	),
)

func (s *Sender) SendStartMessage(chatId int64) error {
	text := `
	Привет! Это бот по заказу одежды и обуви с пойзона.
	Используйте команды, которые вам предлагает бот.
	
	Чтоб сделать новый заказ, выполни команду - /create_order.

	После создания заказа, вы сможете заполнить данные по товарам, которые хотите заказать.
	Если заполнили данные по товару, но понимаете что нужно что-то поменять, используйте команду - /remove_item_data
	Эта команда позволит заполнить все данные заново.

	Если заполнили все данные по товару и хотите добавить в товар ещё один, используйте команду - /add_new_item_to_order

	При завершении оформления заказа используйте команду - /send_order_to_manage
	После этого с вами свяжется наш менеджер.

	Если хотите отменить свой заказ, нужно выполнить команду /cancel_order.

	
	`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = startKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendNotificationAboutNewOrder(chatId int64, orderID int) error {
	text := `Ваш заказ с номером - ` + strconv.Itoa(orderID)
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = editKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendNotificationAboutCancelOrder(chatId int64, orderID int) error {
	text := `Ваш заказ с номером - ` + strconv.Itoa(orderID) + ` отменён.`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = startKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendRequestUrl(chatId int64) error {
	text := `
	Отправь ссылку на товар. 
	Формат ввода: https://dw4.co/t/A/1qtUwh81O
`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = editKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendRequestPrice(chatId int64) error {
	text := `
	Сколько стоит товар в CNY (¥)? 
	Формат ввода: 186
`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = editKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendRequestQuantity(chatId int64) error {
	text := `
	Какое количество нужно заказать? 
	Формат ввода: 2
`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = editKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendChoiceToAddItem(chatId int64) error {
	text := `
	Хотите добавить ещё один товар, введите команду /add_new_item_to_order
`
	msg := tgbotapi.NewMessage(chatId, text)

	msg.ReplyMarkup = addNewItemKeyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendUnknownMessage(chatId int64) error {
	text := `
	Непонятный формат ввода
	Если хотите всё сбросить и начать сначала, введите команду /cancel_order
`
	msg := tgbotapi.NewMessage(chatId, text)

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendUserOrderReport(chatId int64, order orders.Order, totalPrice float64) error {
	itemsText := ""
	for _, item := range order.Items {
		itemsText += getItemText(item)
	}

	orderText := fmt.Sprintf(
		`
User: @%s
Total price:	%.2f RUB (₽)
Items:
%s
`,
		order.UserName,
		totalPrice,
		itemsText,
	)

	msg := tgbotapi.NewMessage(chatId, orderText)
	msg.DisableWebPagePreview = true

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func (s *Sender) SendAdminOrderReport(chatId int64, order orders.Order, exchangeRate float64, totalPrice float64) error {
	itemsText := ""
	for _, item := range order.Items {
		itemsText += getItemText(item)
	}

	orderText := fmt.Sprintf(
		`
User: @%s
Exchange rate:	%.2f RUB (₽) = 1 CNY (¥)
Total price:	%.2f RUB (₽)
Items:
%s
`,
		order.UserName,
		exchangeRate,
		totalPrice,
		itemsText,
	)

	msg := tgbotapi.NewMessage(chatId, orderText)
	msg.DisableWebPagePreview = true

	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Printf("Error sending start message: %v", err)
		return err
	}
	return nil
}

func getItemText(item basket.BasketItem) string {
	return fmt.Sprintf(
		`

-----------------------------
Item № %d
Link: %s
Price: %d CNY (¥)
Quantity: %d
-----------------------------

				`,
		item.ID,
		item.Url,
		item.Price,
		item.Quantity,
	)
}
