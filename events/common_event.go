package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/apepenkov/gotgbot/actions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommonEvent struct {
	Bot    *tgbotapi.BotAPI
	Update *tgbotapi.Update

	ChatId  int64
	actions []actions.Action
	Context context.Context
}

func (e *CommonEvent) AddAction(action actions.Action) {
	e.actions = append(e.actions, action)
}

func (e *CommonEvent) ExecuteAction(action actions.Action) error {
	_, err := action.Execute(e.Bot)
	e.AddAction(action)
	return fmt.Errorf("failed to execute action: %w", err)
}

func (e *CommonEvent) ExecuteActions() error {
	for _, action := range e.actions {
		if action.IsExecuted() {
			continue
		}
		_, err := action.Execute(e.Bot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *CommonEvent) RollbackActions() error {
	for i := len(e.actions) - 1; i >= 0; i-- {
		action := e.actions[i]
		if !action.IsExecuted() {
			continue
		}
		err := action.RollBack(e.Bot)
		if err != nil {
			if errors.Is(err, actions.AlreadyRollBackedError{}) {
				continue
			}
			if errors.Is(err, actions.CantRollBackError{}) {
				continue
			}
			return fmt.Errorf("failed to rollback action: %w", err)
		}
	}
	return nil
}

func (e *CommonEvent) RespondAction() *actions.SendMessageAction {
	return actions.NewSendMessageAction().WithChatID(e.ChatId)
}

func NewCommonEvent(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctx context.Context) CommonEvent {
	return CommonEvent{
		Bot:     bot,
		Update:  update,
		Context: ctx,
	}
}
