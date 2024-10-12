package events

import (
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

func (c *CallbackEvent) ImplementsEvent() {}
