package domain

import (
	"time"
)

type BasketItem struct { // TODO: Вынести в domain
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Url       string    `json:"url"`
	Price     int       `json:"price"`
	Type      ItemType  `json:"type"`
	Size      string    `json:"size"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemType string

const (
	ItemTypeShoes     ItemType = "Обувь"
	ItemTypeOuterwear ItemType = "Верхняя одежда"
	ItemTypeCloth     ItemType = "Футболки или кофты"
)
