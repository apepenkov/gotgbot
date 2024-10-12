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
	if err == nil {
		return nil
	}
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
	errorsOccurred := make([]error, 0)

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
			errorsOccurred = append(errorsOccurred, err)
		}
	}
	if len(errorsOccurred) > 0 {
		return errors.Join(errorsOccurred...)
	}
	return nil
}

func (e *CommonEvent) RespondAction() *actions.SendMessageAction {
	return actions.NewSendMessageAction().WithChatID(e.ChatId)
}

func NewCommonEvent(bot *tgbotapi.BotAPI, update *tgbotapi.Update, ctx context.Context) CommonEvent {
	chatId := int64(0)
	if update.Message != nil && update.Message.Chat != nil {
		chatId = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	}
	return CommonEvent{
		Bot:     bot,
		Update:  update,
		ChatId:  chatId,
		actions: make([]actions.Action, 0),
		Context: ctx,
	}
}
