package entity

import (
	"time"

	basket "poison_bot/internal/db/basket/entity"
)

type Username string // format: username

type Order struct { // TODO: Вынести в domain
	ID        int                 `json:"id"`
	UserName  Username            `json:"userName"` // format: username
	Items     []basket.BasketItem `json:"items"`
	Status    OrderStatus         `json:"status"` // TODO: сделать enum new / inProgress / done
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "New"
	OrderStatusInProcess OrderStatus = "InProcess"
	OrderStatusComplete  OrderStatus = "Complete"
	OrderStatusCancelled OrderStatus = "Cancelled"
)