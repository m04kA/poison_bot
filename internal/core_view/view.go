package core_view

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	orders "poison_bot/internal/db/orders/entity"
)

type View struct {
	log                       *log.Logger
	sender                    Sender
	orderRepo                 OrderRepository
	updates                   tgbotapi.UpdatesChannel
	wg                        *sync.WaitGroup
	itemProcessor             ItemProcessor
	idChannelForOrdersReports int64
}

func New(l *log.Logger, sender Sender, or OrderRepository, updates tgbotapi.UpdatesChannel, wg *sync.WaitGroup, ip ItemProcessor, idReportChanel int64) *View {
	return &View{
		log:                       l,
		sender:                    sender,
		orderRepo:                 or,
		updates:                   updates,
		wg:                        wg,
		itemProcessor:             ip,
		idChannelForOrdersReports: idReportChanel,
	}
}

var userSession = make(map[string]SessionData) // TODO: сделать независисое хранилище сессий

func (v *View) Process() (err error) {
	defer v.wg.Done()

	for update := range v.updates {
		if update.Message == nil {
			continue
		}

		userName := update.Message.From.UserName
		ChatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start", "help":
				err = v.sender.SendStartMessage(ChatID)
				if err != nil {
					return err
				}
			case "create_order":
				orderIndex := v.orderRepo.CreateOrder(userName)

				data, ok := userSession[userName]
				if !ok {
					userSession[userName] = SessionData{
						OrderIndex: &orderIndex,
					}
				} else {
					data.OrderIndex = &orderIndex
					userSession[userName] = data
				}

				err = v.sender.SendNotificationAboutNewOrder(ChatID, orderIndex)
				if err != nil {
					return err
				}

				err = v.sender.SendRequestUrl(ChatID)
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
					err = v.sender.SendNotificationAboutCancelOrder(ChatID, *orderIndex)
					if err != nil {
						return err
					}
				}
				delete(userSession, userName)
				err = v.sender.SendStartMessage(ChatID)
				if err != nil {
					return err
				}
			case "remove_item_data", "add_new_item_to_order":
				data, ok := userSession[userName]
				if !ok {
					userSession[userName] = SessionData{}
					err = v.sender.SendUnknownMessage(ChatID)
					if err != nil {
						return err
					}
				} else {
					data.Url = nil
					data.Price = nil
					data.Quantity = nil
					userSession[userName] = data

					err = v.sender.SendRequestUrl(ChatID)
					if err != nil {
						return err
					}
				}
			case "send_order_to_manage":
				order := &orders.Order{}
				data, ok := userSession[userName]
				if !ok {
					order, err = v.orderRepo.GetOrder(userName, nil)
					if err != nil {
						return err
					}
					userSession[userName] = SessionData{}
					if order == nil || order.Status != orders.OrderStatusNew {
						err = v.sender.SendUnknownMessage(ChatID)
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
					if order == nil || order.Status != orders.OrderStatusNew {
						err = v.sender.SendUnknownMessage(ChatID)
						if err != nil {
							return err
						}
						continue
					}
				}

				err = v.sender.SendOrderReport(ChatID, *order)
				if err != nil {
					return err
				}
				err = v.sender.SendOrderReport(v.idChannelForOrdersReports, *order)
				if err != nil {
					return err
				}
				order.Status = orders.OrderStatusInProcess
				err = v.orderRepo.UpdateOrder(userName, *order)
				if err != nil {
					return err
				}

			default:
				err = v.sender.SendUnknownMessage(ChatID)
				if err != nil {
					return err
				}
			}

		} else {
			sessionData, ok := userSession[userName]
			if !ok {
				sessionData = SessionData{}
				userSession[userName] = sessionData
			}
			sessionData, err = v.itemProcessor.ProcessCreateItem(update, sessionData, ok)
			if err != nil {
				return err
			}
			userSession[userName] = sessionData
		}
		//if update.Message.IsCommand() && update.Message.Command() == "create_item" {

		//orderEntity := orderRepo.CreateOrder(update.userName)
		//sender.SendMessege("Скинь ссылку на товар Брат")
		//url := await recipient.GetMessage()
		//}

	}
	return nil

}
