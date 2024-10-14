package actions

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EditMessageAction struct {
	chatID          int64
	messageID       int
	inlineMessageID string
	text            *string
	markup          *tgbotapi.InlineKeyboardMarkup
	asMarkdown      bool
	asHTML          bool
	wasFile         bool
	entities        []tgbotapi.MessageEntity
	rolledBack      bool

	// prevMessage, when set, would let us rollback to the previous message
	prevMessage *tgbotapi.Message

	editedMessage *tgbotapi.Message
}

func (a *EditMessageAction) WithChatAndMessageId(chatID int64, messageID int, isFile bool) *EditMessageAction {
	a.chatID = chatID
	a.messageID = messageID
	a.wasFile = isFile
	return a
}

func (a *EditMessageAction) WithMessage(message *tgbotapi.Message) *EditMessageAction {
	a.chatID = message.Chat.ID
	a.messageID = message.MessageID
	a.prevMessage = message
	if message.Document != nil || message.Photo != nil || message.Video != nil || message.Audio != nil {
		a.wasFile = true
	}
	return a
}

func (a *EditMessageAction) WithInlineMessageID(inlineMessageID string) *EditMessageAction {
	a.inlineMessageID = inlineMessageID
	return a
}

func (a *EditMessageAction) WithText(text string) *EditMessageAction {
	a.text = &text
	return a
}

func (a *EditMessageAction) WithMarkdown() *EditMessageAction {
	a.asMarkdown = true
	return a
}

func (a *EditMessageAction) WithHTML() *EditMessageAction {
	a.asHTML = true
	return a
}

func (a *EditMessageAction) WithMarkup(markup *tgbotapi.InlineKeyboardMarkup) *EditMessageAction {
	a.markup = markup
	return a
}

func (a *EditMessageAction) WithEntities(entities []tgbotapi.MessageEntity) *EditMessageAction {
	a.entities = entities
	return a
}

func (a *EditMessageAction) parseMode() string {
	if a.asMarkdown {
		return "Markdown"
	}
	if a.asHTML {
		return "HTML"
	}
	return ""
}

func (a *EditMessageAction) Execute(api *tgbotapi.BotAPI) (interface{}, error) {
	if a.text == nil && a.markup == nil {
		return nil, errors.New("no text or markup to edit")
	}
	if a.chatID == 0 && a.inlineMessageID == "" {
		return nil, errors.New("no chat id")
	}
	if a.messageID == 0 {
		return nil, errors.New("no message id")
	}
	if a.markup != nil && !IsMarkup(a.markup) {
		return nil, errors.New("invalid markup")
	}
	var chattable tgbotapi.Chattable

	baseEdit := tgbotapi.BaseEdit{
		ChatID:          a.chatID,
		MessageID:       a.messageID,
		InlineMessageID: a.inlineMessageID,
		ReplyMarkup:     a.markup,
	}

	if a.text == nil && a.markup != nil {
		chattable = &tgbotapi.EditMessageReplyMarkupConfig{
			BaseEdit: baseEdit,
		}
	} else {
		if a.wasFile {
			chattable = &tgbotapi.EditMessageCaptionConfig{
				BaseEdit:  baseEdit,
				ParseMode: a.parseMode(),
			}
			if a.text != nil {
				chattable.(*tgbotapi.EditMessageCaptionConfig).Caption = *a.text
			}
			if len(a.entities) > 0 {
				chattable.(*tgbotapi.EditMessageCaptionConfig).CaptionEntities = a.entities
			}
		} else {
			chattable = &tgbotapi.EditMessageTextConfig{
				BaseEdit:  baseEdit,
				ParseMode: a.parseMode(),
			}
			if a.text != nil {
				chattable.(*tgbotapi.EditMessageTextConfig).Text = *a.text
			}
			if len(a.entities) > 0 {
				chattable.(*tgbotapi.EditMessageTextConfig).Entities = a.entities
			}
		}
	}

	editedMessage, err := api.Send(chattable)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	a.editedMessage = &editedMessage
	return &editedMessage, nil
}

// ExecuteR is a helper function for Execute that returns *tgbotapi.Message instead of interface{}
func (a *EditMessageAction) ExecuteR(api *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	res, err := a.Execute(api)
	if err != nil {
		return nil, err
	}
	return res.(*tgbotapi.Message), nil
}

func (a *EditMessageAction) RollBack(api *tgbotapi.BotAPI) error {
	if a.rolledBack {
		return AlreadyRollBackedError{}
	}
	if a.editedMessage == nil {
		return nil
	}
	if a.prevMessage == nil {
		return CantRollBackError{}
	}
	action := EditMessageAction{}
	_, err := action.WithMessage(a.prevMessage).WithMarkup(a.prevMessage.ReplyMarkup).WithEntities(a.prevMessage.Entities).Execute(api)
	a.rolledBack = true
	if err != nil {
		return fmt.Errorf("failed to rollback edit message: %w", err)
	}
	return nil
}

func (a *EditMessageAction) IsExecuted() bool {
	return a.editedMessage != nil
}

func NewEditMessageAction() *EditMessageAction {
	return &EditMessageAction{}
}
