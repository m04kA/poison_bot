package domain

import "time"

type Item struct {
	ID        int64
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
