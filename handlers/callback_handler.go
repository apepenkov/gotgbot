package handlers

import (
	"context"
	"fmt"
	"github.com/apepenkov/gotgbot/cb_data"
	"github.com/apepenkov/gotgbot/events"
)

type CallbackFunc func(e *events.CallbackEvent, ctx context.Context) error

type CallbackHandler struct {
	CallbackData   cb_data.CallbackData
	ArgumentsCheck func([]string) bool

	Func CallbackFunc
}

func (h *CallbackHandler) Matches(e events.Event) bool {
	event, ok := e.(*events.CallbackEvent)
	if !ok {
		return false
	}

	if event.CbData != h.CallbackData {
		return false
	}

	if h.ArgumentsCheck != nil && !h.ArgumentsCheck(event.Params) {
		return false
	}

	return true
}

func (h *CallbackHandler) Call(e events.Event) error {
	event, ok := e.(*events.CallbackEvent)
	if !ok {
		return fmt.Errorf("expected *events.CallbackEvent, got %T", e)
	}

	return h.Func(event, event.Context)
}
