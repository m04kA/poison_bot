package core_view

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"poison_bot/internal/domain"
)

type View struct {
	log                       *log.Logger
	sender                    Sender
	orderRepo                 OrderRepository
	updates                   tgbotapi.UpdatesChannel
	wg                        *sync.WaitGroup
	itemProcessor             ItemProcessor
	priceCalculator           PriceCalculator
	idChannelForOrdersReports int64
}

func New(l *log.Logger, sender Sender, or OrderRepository, pc PriceCalculator, updates tgbotapi.UpdatesChannel, wg *sync.WaitGroup, ip ItemProcessor, idReportChanel int64) *View {
	return &View{
		log:                       l,
		sender:                    sender,
		orderRepo:                 or,
		updates:                   updates,
		wg:                        wg,
		itemProcessor:             ip,
		priceCalculator:           pc,
		idChannelForOrdersReports: idReportChanel,
	}
}

var userSession = make(map[string]domain.SessionData) // TODO: сделать независисое хранилище сессий

func (v *View) Process() (err error) {
	defer v.wg.Done()

	for update := range v.updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var (
			userName string
			chatID   int64
		)

		switch {
		case update.Message != nil:
			userName = update.Message.From.UserName
			chatID = update.Message.Chat.ID
		case update.CallbackQuery != nil:
			userName = update.CallbackQuery.From.UserName
			chatID = update.CallbackQuery.Message.Chat.ID

			err = v.sender.SendCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if err != nil {
				return err
			}
		}

		if update.Message == nil || !update.Message.IsCommand() {
			sessionData, ok := userSession[userName]
			if !ok {
				sessionData = domain.SessionData{}
				userSession[userName] = sessionData
			}
			sessionData, err = v.itemProcessor.ProcessCreateItem(update, sessionData, ok, userName, chatID)
			if err != nil {
				return err
			}
			userSession[userName] = sessionData
			continue
		}

		switch update.Message.Command() {
		case "start", "help":
			err = v.sender.SendStartMessage(chatID)
			if err != nil {
				return err
			}
		case "create_order":
			orderIndex := v.orderRepo.CreateOrder(userName)

			data, ok := userSession[userName]
			if !ok {
				userSession[userName] = domain.SessionData{
					OrderIndex: &orderIndex,
				}
			} else {
				data.OrderIndex = &orderIndex
				userSession[userName] = data
			}

			err = v.sender.SendNotificationAboutNewOrder(chatID, orderIndex)
			if err != nil {
				return err
			}

			err = v.sender.SendRequestUrl(chatID)
			if err != nil {
				return err
			}
		case "cancel_order":
			orderIndex := userSession[userName].OrderIndex
			if orderIndex != nil {
				err = v.orderRepo.CancelOrder(userName, *orderIndex)
				if err != nil {
					return err
				}
				err = v.sender.SendNotificationAboutCancelOrder(chatID, *orderIndex)
				if err != nil {
					return err
				}
			}
			delete(userSession, userName)
			err = v.sender.SendStartMessage(chatID)
			if err != nil {
				return err
			}
		case "remove_item_data", "add_new_item_to_order":
			data, ok := userSession[userName]
			if !ok {
				userSession[userName] = domain.SessionData{}
				err = v.sender.SendUnknownMessage(chatID)
				if err != nil {
					return err
				}
			} else {
				data = domain.SessionData{
					OrderIndex: data.OrderIndex,
				}
				userSession[userName] = data

				err = v.sender.SendRequestUrl(chatID)
				if err != nil {
					return err
				}
			}
		case "send_order_to_manage":
			order := &domain.Order{}
			data, ok := userSession[userName]
			if !ok {
				order, err = v.orderRepo.GetOrder(userName, nil)
				if err != nil {
					return err
				}
				userSession[userName] = domain.SessionData{}
				if order == nil || order.Status != domain.OrderStatusNew {
					err = v.sender.SendUnknownMessage(chatID)
					if err != nil {
						return err
					}
					continue
				}
			} else {
				order, err = v.orderRepo.GetOrder(userName, data.OrderIndex)
				if err != nil {
					return err
				}
				if order == nil || order.Status != domain.OrderStatusNew {
					err = v.sender.SendUnknownMessage(chatID)
					if err != nil {
						return err
					}
					continue
				}
			}

			exchangeRate := v.priceCalculator.GetExchangeRate()
			totalPRice := v.priceCalculator.Calculate(*order)
			err = v.sender.SendUserOrderReport(chatID, *order, totalPRice)
			if err != nil {
				return err
			}
			err = v.sender.SendAdminOrderReport(v.idChannelForOrdersReports, *order, exchangeRate, totalPRice)
			if err != nil {
				return err
			}
			order.Status = domain.OrderStatusInProcess
			err = v.orderRepo.UpdateOrder(userName, *order)
			if err != nil {
				return err
			}

		case "view_order":
			order := &domain.Order{}
			data, ok := userSession[userName]
			if !ok {
				order, err = v.orderRepo.GetOrder(userName, nil)
				if err != nil {
					return err
				}
				userSession[userName] = domain.SessionData{}
				if order == nil || order.Status != domain.OrderStatusNew {
					err = v.sender.SendUnknownMessage(chatID)
					if err != nil {
						return err
					}
					continue
				}
			} else {
				order, err = v.orderRepo.GetOrder(userName, data.OrderIndex)
				if err != nil {
					return err
				}
				if order == nil || order.Status != domain.OrderStatusNew {
					err = v.sender.SendUnknownMessage(chatID)
					if err != nil {
						return err
					}
					continue
				}
			}
			totalPRice := v.priceCalculator.Calculate(*order)
			err = v.sender.SendUserOrderReport(chatID, *order, totalPRice)
			if err != nil {
				return err
			}

		default:
			err = v.sender.SendUnknownMessage(chatID)
			if err != nil {
				return err
			}
		}

	}

	return nil

}
