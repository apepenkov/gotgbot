package handlers

import (
	"context"
	"fmt"
	"github.com/apepenkov/gotgbot/events"
)

type CommandFunc func(e *events.CommandEvent, ctx context.Context) error

type CommandHandler struct {
	Command                string
	NeedsBotMentionInGroup bool
	ArgumentsCheck         func([]string) bool

	Func CommandFunc
}

func (h *CommandHandler) Matches(e events.Event) bool {
	event, ok := e.(*events.CommandEvent)
	if !ok {
		return false
	}

	if h.Command != event.Command {
		return false
	}

	if h.NeedsBotMentionInGroup && !event.NewMessageEvent.IsPrivate && event.MentionedBot != event.Bot.Self.UserName {
		return false
	}

	if h.ArgumentsCheck != nil && !h.ArgumentsCheck(event.Arguments) {
		return false
	}

	return true
}

func (h *CommandHandler) Call(e events.Event) error {
	event, ok := e.(*events.CommandEvent)
	if !ok {
		return fmt.Errorf("expected *events.CommandEvent, got %T", e)
	}

	return h.Func(event, event.Context)
}
