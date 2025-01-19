package coreView

type Sender interface {
	SendStartMessage(chatId int64) error
}
