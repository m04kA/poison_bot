package create_item

import (
	"net/url"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"poison_bot/internal/domain"
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

func (p *Processor) ProcessCreateItem(
	update tgbotapi.Update,
	session domain.SessionData,
	isActive bool,
	userName string,
	chatID int64,
) (domain.SessionData, error) {
	// TODO: Добавить нормальных логов
	// TODO: Убрать связность между пакетами в передаваемых и отдаваемых структурах

	//chatID := update.Message.Chat.ID
	//userName := update.Message.From.UserName

	if !isActive {
		err := p.sender.SendUnknownMessage(chatID)
		if err != nil {
			return session, err
		}
		return session, nil
	}

	order, err := p.orderRepo.GetOrder(userName, session.OrderIndex)
	if err != nil {
		err = p.sender.SendUnknownMessage(chatID)
		if err != nil {
			return session, err // TODO: Обернуть ошибку нормально
		}
	}

	if order == nil || order.Status != domain.OrderStatusNew {
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

	if update.CallbackQuery != nil {

		if session.Type != nil {
			err = p.sender.SendUnknownMessage(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		}

		data := domain.ItemType(update.CallbackQuery.Data)
		session.Type = &data

		switch data {
		case domain.ItemTypeShoes:
			err = p.sender.SendRequestShoesSize(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		case domain.ItemTypeOuterwear, domain.ItemTypeCloth:
			err = p.sender.SendRequestClosesSize(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		}
	}

	switch {
	case session.Url == nil:
		urlText := update.Message.Text
		urlData, errUrl := url.ParseRequestURI(urlText)
		if errUrl != nil {
			err = p.sender.SendUnknownMessage(chatID)
			if err != nil {
				return session, err // TODO: Обернуть ошибку нормально
			}

			return session, nil
		}

		session.Url = urlData

		err = p.sender.SendRequestThinkType(chatID)
		if err != nil {
			return session, err
		}
	case session.Size == nil:
		size := update.Message.Text
		session.Size = &size

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
		item := p.createItemBasket(*order, session)
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

func (p *Processor) createItemBasket(order domain.Order, session domain.SessionData) *domain.BasketItem {
	return &domain.BasketItem{
		ID:        len(order.Items),
		OrderID:   order.ID,
		Url:       session.Url.String(),
		Price:     *session.Price,
		Quantity:  *session.Quantity,
		Type:      *session.Type,
		Size:      *session.Size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
