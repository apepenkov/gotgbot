package actions

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteMessageAction struct {
	chatID    int64
	messageID int

	executed   bool
	rolledBack bool
}

func (a *DeleteMessageAction) WithChatAndMessageId(chatID int64, messageID int) *DeleteMessageAction {
	a.chatID = chatID
	a.messageID = messageID
	return a
}

func (a *DeleteMessageAction) Execute(api *tgbotapi.BotAPI) (interface{}, error) {
	if a.chatID == 0 {
		return nil, errors.New("no chat id")
	}
	if a.messageID == 0 {
		return nil, errors.New("no message id")
	}

	_, err := api.Send(tgbotapi.NewDeleteMessage(a.chatID, a.messageID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete message: %w", err)
	}

	a.executed = true
	return nil, nil
}

func (a *DeleteMessageAction) RollBack(api *tgbotapi.BotAPI) error {
	return nil
}

func (a *DeleteMessageAction) IsExecuted() bool {
	return a.executed
}

func NewDeleteMessageAction() *DeleteMessageAction {
	return &DeleteMessageAction{}
}
