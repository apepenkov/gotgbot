package actions

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendMessageAction struct {
	chatID           int64
	replyToMessageID int
	text             string
	asMarkdown       bool
	asHTML           bool
	markup           interface{}
	file             tgbotapi.RequestFileData
	fileType         FileType

	sentMessage *tgbotapi.Message
	rolledBack  bool
}

func (a *SendMessageAction) WithChatID(id int64) *SendMessageAction {
	a.chatID = id
	return a
}

func (a *SendMessageAction) WithReply(id int) *SendMessageAction {
	a.replyToMessageID = id
	return a
}

func (a *SendMessageAction) WithText(text string) *SendMessageAction {
	a.text = text
	return a
}

func (a *SendMessageAction) WithMarkdown() *SendMessageAction {
	a.asMarkdown = true
	return a
}

func (a *SendMessageAction) WithHTML() *SendMessageAction {
	a.asHTML = true
	return a
}

func (a *SendMessageAction) WithFile(file tgbotapi.RequestFileData) *SendMessageAction {
	a.file = file
	a.fileType = FileTypeDocument
	return a
}

func (a *SendMessageAction) WithPhoto(file tgbotapi.RequestFileData) *SendMessageAction {
	a.file = file
	a.fileType = FileTypePhoto
	return a
}

func (a *SendMessageAction) WithVideo(file tgbotapi.RequestFileData) *SendMessageAction {
	a.file = file
	a.fileType = FileTypeVideo
	return a
}

func (a *SendMessageAction) WithAudio(file tgbotapi.RequestFileData) *SendMessageAction {
	a.file = file
	a.fileType = FileTypeAudio
	return a
}

func (a *SendMessageAction) WithMarkup(markup interface{}) *SendMessageAction {
	a.markup = markup
	return a
}

func (a *SendMessageAction) parseMode() string {
	if a.asMarkdown {
		return "Markdown"
	}
	if a.asHTML {
		return "HTML"
	}
	return ""
}

func (a *SendMessageAction) Execute(api *tgbotapi.BotAPI) (interface{}, error) {
	if a.text == "" && a.file == nil {
		return nil, errors.New("no text or files to send")
	}
	if a.chatID == 0 {
		return nil, errors.New("no chat id")
	}
	if a.markup != nil && !IsMarkup(a.markup) {
		return nil, errors.New("invalid markup")
	}

	var chattable tgbotapi.Chattable

	if a.file != nil {
		switch a.fileType {
		case FileTypeDocument:
			doc := tgbotapi.NewDocument(a.chatID, a.file)
			doc.Caption = a.text
			doc.ParseMode = a.parseMode()
			doc.ReplyToMessageID = a.replyToMessageID
			doc.ReplyMarkup = a.markup
			chattable = doc

		case FileTypePhoto:
			photo := tgbotapi.NewPhoto(a.chatID, a.file)
			photo.Caption = a.text
			photo.ParseMode = a.parseMode()
			photo.ReplyToMessageID = a.replyToMessageID
			photo.ReplyMarkup = a.markup
			chattable = photo

		case FileTypeVideo:
			video := tgbotapi.NewVideo(a.chatID, a.file)
			video.Caption = a.text
			video.ParseMode = a.parseMode()
			video.ReplyToMessageID = a.replyToMessageID
			video.ReplyMarkup = a.markup
			chattable = video

		case FileTypeAudio:
			audio := tgbotapi.NewAudio(a.chatID, a.file)
			audio.Caption = a.text
			audio.ParseMode = a.parseMode()
			audio.ReplyToMessageID = a.replyToMessageID
			audio.ReplyMarkup = a.markup
			chattable = audio

		default:
			return nil, errors.New("unknown file type")
		}
	} else {
		msg := tgbotapi.NewMessage(a.chatID, a.text)
		msg.ParseMode = a.parseMode()
		msg.ReplyToMessageID = a.replyToMessageID
		msg.ReplyMarkup = a.markup
		chattable = msg
	}

	sentMessage, err := api.Send(chattable)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	a.sentMessage = &sentMessage
	return &sentMessage, nil
}

// ExecuteR is a helper function for Execute that returns *tgbotapi.Message instead of interface{}
func (a *SendMessageAction) ExecuteR(api *tgbotapi.BotAPI) (*tgbotapi.Message, error) {
	res, err := a.Execute(api)
	if err != nil {
		return nil, err
	}
	return res.(*tgbotapi.Message), nil
}

func (a *SendMessageAction) RollBack(api *tgbotapi.BotAPI) error {
	if a.rolledBack {
		return AlreadyRollBackedError{}
	}
	if a.sentMessage == nil {
		return nil
	}
	_, err := api.Send(tgbotapi.NewDeleteMessage(a.chatID, a.sentMessage.MessageID))
	a.rolledBack = true
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (a *SendMessageAction) IsExecuted() bool {
	return a.sentMessage != nil
}

func NewSendMessageAction() *SendMessageAction {
	return &SendMessageAction{}
}
