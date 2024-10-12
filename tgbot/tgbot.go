package tgbot

import (
	"context"
	"github.com/apepenkov/gotgbot/events"
	"github.com/apepenkov/gotgbot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TgBot struct {
	Bot         *tgbotapi.BotAPI
	StateGetter handlers.StateGettable

	CommandHandlers    []handlers.CommandHandler
	UnknownCommandFunc handlers.CommandFunc

	MessageHandlers    []handlers.NewMessageHandler
	UnknownMessageFunc handlers.NewMessageFunc

	CallbackHandlers    []handlers.CallbackHandler
	UnknownCallbackFunc handlers.CallbackFunc

	updateGoroutines int
	cancelUpdateChan chan struct{}
}

func NewTgBot(token string, stateGetter handlers.StateGettable) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TgBot{
		Bot:         bot,
		StateGetter: stateGetter,

		CommandHandlers:  make([]handlers.CommandHandler, 0),
		MessageHandlers:  make([]handlers.NewMessageHandler, 0),
		CallbackHandlers: make([]handlers.CallbackHandler, 0),

		updateGoroutines: 5,
	}, nil
}

func (b *TgBot) updateGoroutine(updatesChan tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-b.cancelUpdateChan:
			return
		case update := <-updatesChan:
			err := b.innerHandleUpdate(&update)
			if err != nil {
				log.Printf("Error handling update: %v", err)
			}
		}
	}
}

func (b *TgBot) innerHandleUpdate(update *tgbotapi.Update) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "bot", b.Bot)

	cmn := events.NewCommonEvent(b.Bot, update, ctx)

	if update.Message != nil {
		state := ""
		if b.StateGetter != nil {
			state = b.StateGetter.GetState(cmn.ChatId)
		}
		newMsgEvent := events.NewNewMessageEvent(cmn, state)

		if len(newMsgEvent.Text) > 0 && newMsgEvent.Text[0] == '/' {
			event := events.NewCommandEvent(*newMsgEvent)

			for _, handler := range b.CommandHandlers {
				if handler.Matches(event) {
					return handler.Call(event)
				}
			}
			if b.UnknownCommandFunc != nil {
				return b.UnknownCommandFunc(event, event.Context)
			}
			return nil
		}

		for _, handler := range b.MessageHandlers {
			if handler.Matches(newMsgEvent) {
				return handler.Call(newMsgEvent)
			}
		}
	} else if update.CallbackQuery != nil {
		event := events.NewCallbackEvent(cmn)
		for _, handler := range b.CallbackHandlers {
			if handler.Matches(event) {
				return handler.Call(event)
			}
		}
		if b.UnknownCallbackFunc != nil {
			return b.UnknownCallbackFunc(event, event.Context)
		}
		return nil
	}

	return nil
}

func (b *TgBot) RunLongPooling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Bot.GetUpdatesChan(u)
	b.cancelUpdateChan = make(chan struct{})

	for i := 0; i < b.updateGoroutines; i++ {
		go b.updateGoroutine(updates)
	}
}

func (b *TgBot) Stop() {
	close(b.cancelUpdateChan)
}
