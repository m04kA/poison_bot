package core_view

import (
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	basket "poison_bot/internal/db/basket/entity"
	orders "poison_bot/internal/db/orders/entity"
)

type Sender interface {
	SendNotificationAboutNewOrder(chatId int64, orderID int) error
	SendNotificationAboutCancelOrder(chatId int64, orderID int) error
	SendOrderReport(chatId int64, order orders.Order) error
	SendStartMessage(chatId int64) error
	SendRequestUrl(chatId int64) error
	SendRequestPrice(chatId int64) error
	SendRequestQuantity(chatId int64) error
	SendUnknownMessage(chatId int64) error
}

type OrderRepository interface {
	CancelOrder(username string, orderIndex int) error
	CreateOrder(username string) (index int)
	AddItem(username string, orderIndex int, item basket.BasketItem) (err error)
	GetOrder(username string, orderIndex *int) (*orders.Order, error)
	UpdateOrder(username string, order orders.Order) error
}

type ItemProcessor interface {
	ProcessCreateItem(update tgbotapi.Update, session SessionData, isActive bool) (SessionData, error)
}

type SessionData struct {
	Url        *url.URL
	Price      *int
	Quantity   *int
	OrderIndex *int
}
