package repository

import (
	"errors"
	"strconv"
	"time"

	basket "poison_bot/internal/db/basket/entity"
	"poison_bot/internal/db/orders/entity"
)

type OrderRepository struct {
	orders map[string][]entity.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string][]entity.Order),
	}
}

func (r *OrderRepository) CreateOrder(username string) (index int) {
	user, ok := r.orders[username]
	if !ok {
		r.orders[username] = make([]entity.Order, 0, 1)
		user = r.orders[username]
	}

	order := entity.Order{}
	if len(user) != 0 {
		order = user[len(user)-1]
		if order.Status != entity.OrderStatusNew {
			order = entity.Order{
				ID:        len(r.orders[username]), // TODO: сделать нормальную логику
				UserName:  entity.Username(username),
				Items:     make([]basket.BasketItem, 0),
				Status:    entity.OrderStatusNew,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			user = append(user, order)
		}
	}

	return len(user) - 1
}

func (r *OrderRepository) GetOrder(username string, orderIndex *int) (*entity.Order, error) {
	user, ok := r.orders[username]
	if !ok {
		return nil, errors.New("user with username: " + username + " not found")
	}

	if orderIndex != nil && (*orderIndex >= len(user) || *orderIndex < 0) {
		return nil, errors.New("for user: " + username + "; invalid index order: " + strconv.Itoa(*orderIndex))
	}

	if len(user) == 0 {
		return nil, nil
	}

	if orderIndex == nil {
		lastIndex := len(user) - 1
		orderIndex = &lastIndex
	}
	return &user[*orderIndex], nil
}

func (r *OrderRepository) CancelOrder(username string, orderIndex int) error {
	user, ok := r.orders[username]
	if !ok {
		return errors.New("user with username: " + username + " not found")
	}

	if len(user) == 0 || orderIndex >= len(user) || orderIndex < 0 {
		return errors.New("for user: " + username + "; not found order with index: " + strconv.Itoa(orderIndex))
	}

	user[orderIndex].Status = entity.OrderStatusCancelled
	return nil
}

func (r *OrderRepository) UpdateOrder(username string, order entity.Order) error {
	user, ok := r.orders[username]
	if !ok {
		return errors.New("user with username: " + username + " not found")
	}

	if len(user) == 0 {
		return errors.New("for user: " + username + "; not found orders")
	}

	user[order.ID] = order
	return nil
}

func (r *OrderRepository) AddItem(username string, orderIndex int, item basket.BasketItem) (err error) {
	user, ok := r.orders[username]
	if !ok {
		return errors.New("user with username: " + username + " not found")
	}

	if orderIndex >= len(user) || orderIndex < 0 {
		return errors.New("for user: " + username + "; invalid indx order: " + strconv.Itoa(orderIndex))
	}

	order := user[orderIndex]
	order.Items = append(order.Items, item)
	order.UpdatedAt = time.Now()
	user[orderIndex] = order

	return nil
}