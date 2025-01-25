package create_item

import (
	"net/url"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	coreView "poison_bot/internal/core_view"
	basket "poison_bot/internal/db/basket/entity"
	orders "poison_bot/internal/domain"
)

type Processor struct {
	orderRepo OrderRepository
	sender    Sender
}

func NewProcessor(or OrderRepository, s Sender) *Processor {
	return &Processor{
		orderRepo: or,
		sender:    s,
	}
}

func (p *Processor) ProcessCreateItem(update tgbotapi.Update, session coreView.SessionData, isActive bool) (coreView.SessionData, error) {
	// TODO: Добавить нормальных логов
	// TODO: Убрать связность между пакетами в передаваемых и отдаваемых структурах

	chatID := update.Message.Chat.ID
	userName := update.Message.From.UserName

	if !isActive {
		err := p.sender.SendUnknownMessage(chatID)
		if err != nil {
			return session, err
		}
		return session, nil
	}

	order, err := p.orderRepo.GetOrder(userName, session.OrderIndex)
	if err != nil {
		return session, err
	}

	if order == nil || order.Status != orders.OrderStatusNew {
		newOrderInd := p.orderRepo.CreateOrder(userName)
		session.OrderIndex = &newOrderInd
		order, err = p.orderRepo.GetOrder(userName, session.OrderIndex)
		if err != nil {
			return session, err
		}
	}

	if session.OrderIndex == nil {
		id := order.ID
		session.OrderIndex = &id
	}

	switch {
	case session.Url == nil:
		urlText := update.Message.Text
		urlData, err := url.ParseRequestURI(urlText)
		if err != nil {
			err = p.sender.SendUnknownMessage(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		}

		session.Url = urlData

		err = p.sender.SendRequestPrice(chatID)
		if err != nil {
			return session, err
		}
	case session.Price == nil:
		price, errParce := strconv.Atoi(update.Message.Text)
		if errParce != nil {
			err = p.sender.SendUnknownMessage(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		}

		session.Price = &price

		err = p.sender.SendRequestQuantity(chatID)
		if err != nil {
			return session, err
		}

	case session.Quantity == nil:
		quantity, errParce := strconv.Atoi(update.Message.Text)
		if errParce != nil {
			err = p.sender.SendUnknownMessage(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}
			return session, nil
		}
		session.Quantity = &quantity
		item := p.createItemBasket(*order, session.Url.String(), *session.Price, *session.Quantity)
		err = p.orderRepo.AddItem(userName, order.ID, *item)
		if err != nil {
			return session, err
		}

		err = p.sender.SendChoiceToAddItem(chatID)
		if err != nil {
			return session, err
		}
	}

	return session, nil
}

func (p *Processor) createItemBasket(order orders.Order, url string, price int, quantity int) *basket.BasketItem {
	return &basket.BasketItem{
		ID:        len(order.Items),
		OrderID:   order.ID,
		Url:       url,
		Price:     price,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
