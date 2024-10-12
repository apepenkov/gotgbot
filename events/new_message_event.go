package events

import (
	"github.com/apepenkov/gotgbot/actions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewMessageEvent struct {
	CommonEvent
	InitialState string
	senderId     int64
	senderOk     bool
	Text         string
	Message      *tgbotapi.Message
	IsPrivate    bool
}

func (e *NewMessageEvent) SetSenderId(senderId int64) {
	e.senderId = senderId
	e.senderOk = true
}

func (e *NewMessageEvent) SenderId() (int64, bool) {
	return e.senderId, e.senderOk
}

func (e *NewMessageEvent) MustSenderId() int64 {
	if e.senderOk {
		return e.senderId
	}
	panic("senderId is not set")
}

func (e *NewMessageEvent) ImplementsEvent() {}

func (e *NewMessageEvent) ReplyAction() *actions.SendMessageAction {
	return actions.NewSendMessageAction().WithChatID(e.ChatId).WithReply(e.Message.MessageID)
}

func NewNewMessageEvent(cmn CommonEvent, initialState string) *NewMessageEvent {
	e := &NewMessageEvent{
		CommonEvent:  cmn,
		InitialState: initialState,
	}
	if cmn.Update.Message != nil {
		if cmn.Update.Message.Chat != nil {
			e.ChatId = cmn.Update.Message.Chat.ID
			e.IsPrivate = cmn.Update.Message.Chat.IsPrivate()
		}
		e.Text = cmn.Update.Message.Text
		e.Message = cmn.Update.Message
	}

	return e
}
