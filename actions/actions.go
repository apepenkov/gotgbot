package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// all handlers will queue up actions for themselves. If handler is executed without errors, actions will be executed
// actions are sending messages, editing messages, deleting messages, etc.
// some actions can be rolled back, some can't. If action can't be rolled back, it should return CantRollBackError.
// Some actions will be executed in-flight (during handling, for example, sending message like "please wait").

type CantRollBackError struct{}

func (e CantRollBackError) Error() string {
	return "Can't roll back this action"
}

type AlreadyRollBackedError struct{}

func (e AlreadyRollBackedError) Error() string {
	return "Action already rolled back"
}

func IsMarkup(markup interface{}) bool {
	switch markup.(type) {
	case tgbotapi.InlineKeyboardMarkup, tgbotapi.ReplyKeyboardMarkup, tgbotapi.ReplyKeyboardRemove:
		return true
	default:
		return false
	}
}

type Action interface {
	Execute(api *tgbotapi.BotAPI) (interface{}, error)
	RollBack(api *tgbotapi.BotAPI) error
	IsExecuted() bool
}

type FileType int

const (
	FileTypeDocument FileType = iota
	FileTypePhoto
	FileTypeVideo
	FileTypeAudio
)
