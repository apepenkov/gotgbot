package handlers

import (
	"context"
	"fmt"
	"github.com/apepenkov/gotgbot/events"
	"regexp"
)

type NewMessageFunc func(e *events.NewMessageEvent, ctx context.Context) error

type HandlerNewMessage struct {
	Pattern   regexp.Regexp
	ByPattern bool

	State   string
	ByState bool

	StringPrefix   string
	ByStringPrefix bool

	FullStringMatch   string
	ByFullStringMatch bool

	All bool

	Func NewMessageFunc
}

func (h *HandlerNewMessage) Matches(e events.Event) bool {
	event, ok := e.(*events.NewMessageEvent)
	if !ok {
		return false
	}

	if h.All {
		return true
	}

	if h.ByPattern && h.Pattern.MatchString(event.Text) {
		return true
	}
	if h.ByState && event.InitialState == h.State {
		return true
	}
	if h.ByStringPrefix && len(event.Text) >= len(h.StringPrefix) && event.Text[:len(h.StringPrefix)] == h.StringPrefix {
		return true
	}
	if h.ByFullStringMatch && h.FullStringMatch == event.Text {
		return true
	}
	return false
}

func (h *HandlerNewMessage) Call(e events.Event) error {
	event, ok := e.(*events.NewMessageEvent)
	if !ok {
		return fmt.Errorf("expected *events.NewMessageEvent, got %T", e)
	}

	return h.Func(event, event.Context)
}

func NewMessageByPattern(pattern regexp.Regexp, fn NewMessageFunc) *HandlerNewMessage {
	return &HandlerNewMessage{
		Pattern:   pattern,
		ByPattern: true,
		Func:      fn,
	}
}

func NewMessageByState(state string, fn NewMessageFunc) *HandlerNewMessage {
	return &HandlerNewMessage{
		State:   state,
		ByState: true,
		Func:    fn,
	}
}

func NewMessageByStringPrefix(prefix string, fn NewMessageFunc) *HandlerNewMessage {
	return &HandlerNewMessage{
		StringPrefix:   prefix,
		ByStringPrefix: true,
		Func:           fn,
	}
}

func NewMessageByFullStringMatch(match string, fn NewMessageFunc) *HandlerNewMessage {
	return &HandlerNewMessage{
		FullStringMatch:   match,
		ByFullStringMatch: true,
		Func:              fn,
	}
}

func NewMessageAll(fn NewMessageFunc) *HandlerNewMessage {
	return &HandlerNewMessage{
		All:  true,
		Func: fn,
	}
}
