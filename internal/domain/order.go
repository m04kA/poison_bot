package domain

import "time"

type Order struct {
	ID        int
	UserName  string // format: @username
	Items     []Item
	CreatedAt time.Time
	UpdatedAt time.Time
}
