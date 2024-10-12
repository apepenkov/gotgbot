package events

import (
	"github.com/apepenkov/gotgbot/actions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewMessageEvent struct {
	CommonEvent
	senderId  int64
	senderOk  bool
	Text      string
	Message   *tgbotapi.Message
	IsPrivate bool
}

func (e *NewMessageEvent) SetSenderId(senderId int64) {
	e.senderId = senderId
	e.senderOk = true
}

func (e *NewMessageEvent) SenderId() (int64, bool) {
	return e.senderId, e.senderOk
}

func (e *NewMessageEvent) ImplementsEvent() {}

func (e *NewMessageEvent) ReplyAction() *actions.SendMessageAction {
	return actions.NewSendMessageAction().WithChatID(e.ChatId).WithReply(e.Message.MessageID)
}
