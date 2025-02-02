package domain

import "net/url"

type SessionData struct {
	Url        *url.URL
	Price      *int
	Quantity   *int
	OrderIndex *int
	Type       *ItemType
	Size       *string
}
