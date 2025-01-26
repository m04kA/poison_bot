package core_view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"poison_bot/internal/domain"
)

type Sender interface {
	SendNotificationAboutNewOrder(chatId int64, orderID int) error
	SendNotificationAboutCancelOrder(chatId int64, orderID int) error
	SendUserOrderReport(chatId int64, order domain.Order, totalPrice float64) error
	SendAdminOrderReport(chatId int64, order domain.Order, exchangeRate float64, totalPrice float64) error
	SendStartMessage(chatId int64) error
	SendRequestUrl(chatId int64) error
	SendRequestPrice(chatId int64) error
	SendRequestQuantity(chatId int64) error
	SendUnknownMessage(chatId int64) error
	SendCallback(callbackID, data string) error
}

type OrderRepository interface {
	CancelOrder(username string, orderIndex int) error
	CreateOrder(username string) (index int)
	AddItem(username string, orderIndex int, item domain.BasketItem) (err error)
	GetOrder(username string, orderIndex *int) (*domain.Order, error)
	UpdateOrder(username string, order domain.Order) error
}

type ItemProcessor interface {
	ProcessCreateItem(update tgbotapi.Update, session domain.SessionData, isActive bool, userName string, chatID int64) (domain.SessionData, error)
}

type PriceCalculator interface {
	Calculate(order domain.Order) float64
	GetExchangeRate() float64
}
