package price_calculator

import (
	domain "poison_bot/internal/domain"
)

type Calculator struct {
	exchangeRate float64
}

func New(er float64) *Calculator {
	return &Calculator{
		exchangeRate: er,
	}
}

func (c *Calculator) Calculate(order domain.Order) float64 {
	totalPrice := 0.0
	for _, item := range order.Items {
		totalPrice += float64(item.Price) * float64(item.Quantity)
	}
	return totalPrice * c.exchangeRate
}

func (c *Calculator) GetExchangeRate() float64 {
	return c.exchangeRate
}
