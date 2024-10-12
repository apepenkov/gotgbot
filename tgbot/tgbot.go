package tgbot

import (
	"context"
	"errors"
	"fmt"
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

	DefAnErrorOccurredFunc func(event events.Event, err error) error

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
			event, err := b.innerHandleUpdate(&update)
			if err != nil {
				if b.DefAnErrorOccurredFunc != nil {
					_ = b.DefAnErrorOccurredFunc(event, err)
				}
				log.Printf("Error handling update: %v", err)
			}
		}
	}
}

func (b *TgBot) innerHandleUpdate(update *tgbotapi.Update) (event events.Event, err error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "bot", b.Bot)

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("%v", x)
			}
		}
	}()

	cmn := events.NewCommonEvent(b.Bot, update, ctx)

	if update.Message != nil {
		state := ""
		if b.StateGetter != nil {
			state = b.StateGetter.GetState(cmn.ChatId)
		}
		newMsgEvent := events.NewNewMessageEvent(cmn, state)
		event = newMsgEvent

		if len(newMsgEvent.Text) > 0 && newMsgEvent.Text[0] == '/' {
			commandEvent := events.NewCommandEvent(*newMsgEvent)
			event = commandEvent

			for _, handler := range b.CommandHandlers {
				if handler.Matches(commandEvent) {
					return event, handler.Call(commandEvent)
				}
			}
			if b.UnknownCommandFunc != nil {
				return event, b.UnknownCommandFunc(commandEvent, commandEvent.Context)
			}
			return event, nil
		}

		for _, handler := range b.MessageHandlers {
			if handler.Matches(newMsgEvent) {
				return event, handler.Call(newMsgEvent)
			}
		}
	} else if update.CallbackQuery != nil {
		callbackEvent := events.NewCallbackEvent(cmn)
		event = callbackEvent
		for _, handler := range b.CallbackHandlers {
			if handler.Matches(callbackEvent) {
				return event, handler.Call(callbackEvent)
			}
		}
		if b.UnknownCallbackFunc != nil {
			return event, b.UnknownCallbackFunc(callbackEvent, callbackEvent.Context)
		}
		return
	}

	return
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

func (b *TgBot) RegisterCommandHandler(handler handlers.CommandHandler) {
	b.CommandHandlers = append(b.CommandHandlers, handler)
}

func (b *TgBot) RegisterMessageHandler(handler handlers.NewMessageHandler) {
	b.MessageHandlers = append(b.MessageHandlers, handler)
}

func (b *TgBot) RegisterCallbackHandler(handler handlers.CallbackHandler) {
	b.CallbackHandlers = append(b.CallbackHandlers, handler)
}

func (b *TgBot) SetUnknownCommandFunc(f handlers.CommandFunc) {
	b.UnknownCommandFunc = f
}

func (b *TgBot) SetUnknownMessageFunc(f handlers.NewMessageFunc) {
	b.UnknownMessageFunc = f
}

func (b *TgBot) SetUnknownCallbackFunc(f handlers.CallbackFunc) {
	b.UnknownCallbackFunc = f
}
