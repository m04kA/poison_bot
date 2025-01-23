package entity

import (
	"time"
)

type BasketItem struct { // TODO: Вынести в domain
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Url       string    `json:"url"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
