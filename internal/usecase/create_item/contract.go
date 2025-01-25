package create_item

import (
	basket "poison_bot/internal/db/basket/entity"
	orders "poison_bot/internal/domain"
)

type OrderRepository interface {
	CreateOrder(username string) (index int)
	AddItem(username string, orderIndex int, item basket.BasketItem) (err error)
	GetOrder(username string, orderIndex *int) (*orders.Order, error)
}

type Sender interface {
	SendStartMessage(chatId int64) error
	SendRequestUrl(chatId int64) error
	SendRequestPrice(chatId int64) error
	SendRequestQuantity(chatId int64) error
	SendUnknownMessage(chatId int64) error
	SendChoiceToAddItem(chatId int64) error
}
