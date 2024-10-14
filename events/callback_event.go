package events

import (
	"github.com/apepenkov/gotgbot/actions"
	"github.com/apepenkov/gotgbot/cb_data"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackEvent struct {
	CommonEvent
	CbData    cb_data.CallbackData
	Params    []string
	QueryId   string
	MessageId int
	Callback  *tgbotapi.CallbackQuery
	Answered  bool
}

func (c *CallbackEvent) answer(cb tgbotapi.CallbackConfig) error {
	c.Answered = true
	_, err := c.Bot.Request(cb)
	return err
}

func (c *CallbackEvent) AnswerEmpty() error {
	return c.answer(tgbotapi.NewCallback(c.QueryId, ""))
}

func (c *CallbackEvent) AnswerText(text string) error {
	return c.answer(tgbotapi.NewCallback(c.QueryId, text))
}

func (c *CallbackEvent) AnswerAlert(text string) error {
	return c.answer(tgbotapi.NewCallbackWithAlert(c.QueryId, text))
}

func (c *CallbackEvent) EditAction() *actions.EditMessageAction {
	if c.Update.Message != nil {
		return actions.NewEditMessageAction().WithMessage(c.Update.Message)
	} else if c.Update.CallbackQuery != nil && c.Update.CallbackQuery.Message != nil {
		return actions.NewEditMessageAction().WithMessage(c.Update.CallbackQuery.Message)
	} else if c.Update.CallbackQuery != nil && c.Update.CallbackQuery.InlineMessageID != "" {
		return actions.NewEditMessageAction().WithInlineMessageID(c.Update.CallbackQuery.InlineMessageID)
	}
	return actions.NewEditMessageAction()
}

func (c *CallbackEvent) ImplementsEvent() {}

func NewCallbackEvent(event CommonEvent) *CallbackEvent {
	cbd, args := cb_data.GetCallbackData(event.Update.CallbackQuery.Data)

	e := &CallbackEvent{
		CommonEvent: event,
		CbData:      cbd,
		Params:      args,
		QueryId:     event.Update.CallbackQuery.ID,
		MessageId:   event.Update.CallbackQuery.Message.MessageID,
		Callback:    event.Update.CallbackQuery,
	}

	return e
}
