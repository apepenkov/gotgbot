package handlers

import (
	"github.com/apepenkov/gotgbot/events"
)

type Handler interface {
	Matches(e events.Event) bool
	Call(e events.Event) error
}
